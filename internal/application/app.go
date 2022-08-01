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

const (
	DatabaseUrl          = "DATABASE_URL"
	RunAddress           = "RUN_ADDRESS"
	AccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
)

var (
	server ports.RESTServer
)

func Start(ctx context.Context) {

	readEnv()

	logger := logging.GetLogger()

	pgconn := viper.GetString(DatabaseUrl)
	store, err := db.NewPostgresStore(ctx, pgconn, logger)
	if err != nil {
		logger.Fatal("store creation filed", err)
	}

	address := viper.GetString(RunAddress)
	server, err = rest.NewChiServer(address, store, logger)
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

func readEnv() {
	var err error
	logger := logging.GetLogger()

	err = viper.BindEnv(DatabaseUrl)
	if err != nil {
		logger.Fatalf("database url env: %v", err)
	}
	viper.SetDefault(DatabaseUrl, "postgres://postgres:qwerty@localhost:5432/gophermart?sslmode=disable")

	err = viper.BindEnv(RunAddress)
	if err != nil {
		logger.Fatalf("run address env: %v", err)
	}
	viper.SetDefault(RunAddress, "127.0.0.1:8080")

	err = viper.BindEnv(AccrualSystemAddress)
	if err != nil {
		logger.Fatalf("accrual system address env: %v", err)
	}
	viper.SetDefault(AccrualSystemAddress, "127.0.0.1:8080")
}
