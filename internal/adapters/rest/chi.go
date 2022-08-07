package rest

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
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
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", s.register)
		r.Post("/login", s.login)
	})
	return r
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
	SetCookie(w, u.Login, 24*time.Hour)
	w.WriteHeader(http.StatusOK)
}

func (s *ChiServer) login(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("login http request")
	SetCookie(w, "", -1*time.Hour)
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
	if err == repo.ErrValidation {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		s.logger.Error("store error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SetCookie(w, u.Login, 24*time.Hour)
	w.WriteHeader(http.StatusOK)
}
