package rest

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"gophermart/internal/ports"
	"net"
	"net/http"
)

type ChiServer struct {
	store    ports.Store
	http     *http.Server
	listener net.Listener
	logger   ports.Logger
}

func NewChiServer(address string, s ports.Store, logger ports.Logger) (*ChiServer, error) {
	var (
		server ChiServer
		err    error
	)

	server.logger = logger

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
	return r
}
