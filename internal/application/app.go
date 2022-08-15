package application

import (
	"context"
	"flag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"gophermart/internal/adapters/accrualdispatcher"
	"gophermart/internal/adapters/db"
	"gophermart/internal/adapters/logging"
	"gophermart/internal/adapters/rest"
	"gophermart/internal/ports"
)

const (
	DatabaseURL          = "DATABASE_URL"
	RunAddress           = "RUN_ADDRESS"
	AccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
)

var (
	server  ports.RESTServer
	accrual *accrualdispatcher.GopherAccrualDispatcher
)

func Start(ctx context.Context) {

	readEnv()

	logger := logging.GetLogger()

	pgconn := viper.GetString(DatabaseURL)
	store, err := db.NewPostgresStore(ctx, pgconn, logger)
	if err != nil {
		logger.Fatal("store creation filed", err)
	}

	accrualconn := viper.GetString(AccrualSystemAddress)
	accrual = accrualdispatcher.NewGopherAccrualDispatcher(accrualconn, store, logger)
	accrual.Start()

	address := viper.GetString(RunAddress)
	server, err = rest.NewChiServer(address, store, accrual, logger)
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
	accrual.Stop()
	logger.Info("app stopped")
}

func readEnv() {
	var (
		err        error
		addr       string
		dbURL      string
		accrualSys string
	)
	logger := logging.GetLogger()

	flag.StringVar(&addr, "a", "127.0.0.1:8080", "")
	flag.StringVar(&dbURL, "d", "postgres://postgres:qwerty@localhost:5432/gophermart?sslmode=disable", "")
	flag.StringVar(&accrualSys, "r", "127.0.0.1:8080", "")

	err = viper.BindEnv(DatabaseURL)
	if err != nil {
		logger.Fatalf("database url env: %v", err)
	}
	viper.SetDefault(DatabaseURL, dbURL)

	err = viper.BindEnv(RunAddress)
	if err != nil {
		logger.Fatalf("run address env: %v", err)
	}
	viper.SetDefault(RunAddress, addr)

	err = viper.BindEnv(AccrualSystemAddress)
	if err != nil {
		logger.Fatalf("accrual system address env: %v", err)
	}
	viper.SetDefault(AccrualSystemAddress, accrualSys)
}
