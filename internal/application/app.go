package application

import (
	"context"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"gophermart/internal/adapters/db"
	"gophermart/internal/adapters/logging"
	"gophermart/internal/adapters/rest"
	"gophermart/internal/ports"
)

var (
	server ports.RESTServer
)

func Start(ctx context.Context) {
	logger := logging.GetLogger()

	pgconn := viper.GetString("DATABASE_URL")
	store, err := db.NewPostgresStore(ctx, pgconn, logger)
	if err != nil {
		logger.Fatal("store creation filed", err)
	}

	port := viper.GetString("APP_PORT")
	server, err := rest.NewChiServer(port, store, logger)
	if err != nil {
		logger.Fatal("rest server creation filed", err)
	}

	var g errgroup.Group
	g.Go(func() error {
		return server.Start()
	})

	logger.Info("app started")
	err = g.Wait()
	if err != nil {
		logger.Fatalf("app start failed: %v", err)
	}
}

func Stop() {
	logger := logging.GetLogger()
	err := server.Stop(context.Background())
	if err != nil {
		logger.Fatalf("app stop failed: %v", err)
	}
	logger.Info("app stopped")
}
