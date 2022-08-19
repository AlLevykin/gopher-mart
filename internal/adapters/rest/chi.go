package rest

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gophermart/internal/domain/luhn"
	"gophermart/internal/domain/repo"
	"gophermart/internal/ports"
	"net"
	"net/http"
	"time"
)

type ContextKey string

type ChiServer struct {
	store    ports.Store
	accrual  ports.AccrualDispatcher
	http     *http.Server
	listener net.Listener
	logger   ports.Logger
}

func NewChiServer(address string, store ports.Store, accrual ports.AccrualDispatcher, logger ports.Logger) (*ChiServer, error) {
	var (
		server ChiServer
		err    error
	)

	server.logger = logger
	server.store = store
	server.accrual = accrual

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
	s.accrual.Start()
	if err := s.http.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *ChiServer) Stop(ctx context.Context) error {
	err := s.http.Shutdown(ctx)
	s.logger.Info("http server stopped")
	s.accrual.Stop()
	return err
}

func (s *ChiServer) startOrderProcessing(order string) {
	s.accrual.Dispatch(order)
}

func (s *ChiServer) routes() http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Compress(5))
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", s.register)
		r.Post("/login", s.login)
		r.With(s.ValidateSession).Group(func(r chi.Router) {
			r.Post("/orders", s.uploadOrder)
			r.Get("/orders", s.getOrders)
			r.Get("/balance", s.getBalance)
			r.Post("/balance/withdraw", s.sendWithdraw)
			r.Get("/balance/withdrawals", s.getWithdrawals)
			r.Get("/withdrawals", s.getWithdrawals)
		})
	})
	return r
}

func (s *ChiServer) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l, err := Validate(req)
		if err != nil {
			s.logger.Error("jwt validation failed:", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(req.Context(), ContextKey("LOGIN"), l)
		next.ServeHTTP(w, req.WithContext(ctx))
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
	v := req.Context().Value(ContextKey("LOGIN"))
	if v == nil {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	l, ok := v.(string)
	if !ok {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	num, err := ReadBody(req)
	if err != nil {
		s.logger.Error("can't get order number:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !luhn.Valid(num) {
		s.logger.Error("order number not valid")
		http.Error(w, "order number not valid", http.StatusUnprocessableEntity)
		return
	}
	accepted, err := s.store.IsOrderAcceptedByUser(req.Context(), num, l)
	if err != nil {
		s.logger.Error("can't check order acceptation by user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accepted {
		s.logger.Info("order already accepted by user:", num, ", ", l)
		s.startOrderProcessing(num)
		w.WriteHeader(http.StatusOK)
		return
	}
	accepted, err = s.store.IsOrderAccepted(req.Context(), num)
	if err != nil {
		s.logger.Error("can't check order acceptation:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accepted {
		s.logger.Info("order already accepted:", num)
		w.WriteHeader(http.StatusConflict)
		return
	}
	err = s.store.SaveOrder(req.Context(), num, l)
	if err != nil {
		s.logger.Error("order uploading failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.startOrderProcessing(num)
	w.WriteHeader(http.StatusAccepted)
}

func (s *ChiServer) getOrders(w http.ResponseWriter, req *http.Request) {
	v := req.Context().Value(ContextKey("LOGIN"))
	if v == nil {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	l, ok := v.(string)
	if !ok {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	orders, err := s.store.GetOrders(req.Context(), l)
	if err == sql.ErrNoRows {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		s.logger.Error("orders selecting failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write([]byte(orders))
	if err != nil {
		s.logger.Error("data sending failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *ChiServer) getBalance(w http.ResponseWriter, req *http.Request) {
	v := req.Context().Value(ContextKey("LOGIN"))
	if v == nil {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	l, ok := v.(string)
	if !ok {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	balance, err := s.store.GetBalance(req.Context(), l)
	if err == sql.ErrNoRows {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		s.logger.Error("orders selecting failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write([]byte(balance))
	if err != nil {
		s.logger.Error("data sending failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *ChiServer) sendWithdraw(w http.ResponseWriter, req *http.Request) {
	v := req.Context().Value(ContextKey("LOGIN"))
	if v == nil {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	l, ok := v.(string)
	if !ok {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	b, err := ReadBody(req)
	if err != nil {
		s.logger.Error("no request body for withdrawing:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wd, err := repo.UnmarshalWithdraw(b)
	if err != nil {
		s.logger.Error("bad request for withdrawing:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !luhn.Valid(wd.Order) {
		s.logger.Error("order number not valid")
		http.Error(w, "order number not valid", http.StatusUnprocessableEntity)
		return
	}
	accepted, err := s.store.IsOrderAccepted(req.Context(), wd.Order)
	if err != nil {
		s.logger.Error("can't check order acceptation:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !accepted {
		err = s.store.SaveOrder(req.Context(), wd.Order, l)
		if err != nil {
			s.logger.Error("order uploading failed:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json, err := s.store.GetBalance(req.Context(), l)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	if err != nil {
		s.logger.Error("balance calculation failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	balance, err := repo.UnmarshalBalance(json)
	if err != nil {
		s.logger.Error("balance calculation failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if (balance.Current - wd.Sum) < 0 {
		s.logger.Info("insufficient funds:", l)
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	err = s.store.SaveWithdraw(req.Context(), wd)
	if err != nil {
		s.logger.Error("orders selecting failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *ChiServer) getWithdrawals(w http.ResponseWriter, req *http.Request) {
	v := req.Context().Value(ContextKey("LOGIN"))
	if v == nil {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	l, ok := v.(string)
	if !ok {
		s.logger.Error("can't get context data")
		http.Error(w, "can't get context data", http.StatusInternalServerError)
		return
	}
	withdrawals, err := s.store.GetWithdrawals(req.Context(), l)
	if err == sql.ErrNoRows {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		s.logger.Error("withdrawals selecting failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write([]byte(withdrawals))
	if err != nil {
		s.logger.Error("data sending failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
