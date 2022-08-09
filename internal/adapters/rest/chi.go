package rest

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gophermart/internal/domain/repo"
	"gophermart/internal/ports"
	"net"
	"net/http"
	"time"
)

type ChiServer struct {
	store    ports.Store
	http     *http.Server
	listener net.Listener
	logger   ports.Logger
}

func NewChiServer(address string, store ports.Store, logger ports.Logger) (*ChiServer, error) {
	var (
		server ChiServer
		err    error
	)

	server.logger = logger
	server.store = store

	server.listener, err = net.Listen("tcp", address)
	if err != nil {
		server.logger.Error("failed listen port", err)
	}

	server.http = &http.Server{
		Handler: server.routes(),
	}

	return &server, nil
}

func (s *ChiServer) Start() error {
	if err := s.http.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *ChiServer) Stop(ctx context.Context) error {
	err := s.http.Shutdown(ctx)
	s.logger.Info("http server stopped")
	return err
}

func (s *ChiServer) routes() http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Compress(5))
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", s.register)
		r.Post("/login", s.login)
		r.With(s.ValidateSession).Post("/orders", s.uploadOrder)
	})
	return r
}

func (s *ChiServer) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := Validate(req)
		if err != nil {
			s.logger.Error("jwt validation failed:", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func (s *ChiServer) register(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("register http request")
	b, err := ReadBody(req)
	if err != nil {
		s.logger.Error("no request body for user registration:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := repo.UnmarshalUser(b)
	if err != nil {
		s.logger.Error("bad request for user registration:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.store.RegisterUser(req.Context(), u)
	if err == repo.ErrUserExists {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		s.logger.Error("store error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Login(w, u.Login, 24*time.Hour)
	w.WriteHeader(http.StatusOK)
}

func (s *ChiServer) login(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("login http request")
	b, err := ReadBody(req)
	if err != nil {
		s.logger.Error("no request body for login:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := repo.UnmarshalUser(b)
	if err != nil {
		s.logger.Error("bad request for login:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.store.Validation(req.Context(), u)
	if err == repo.ErrUserValidation {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		s.logger.Error("store error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Login(w, u.Login, 24*time.Hour)
	w.WriteHeader(http.StatusOK)
}

func (s *ChiServer) uploadOrder(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
