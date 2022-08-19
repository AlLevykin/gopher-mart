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
	DatabaseURL          = "DATABASE_URI"
	RunAddress           = "RUN_ADDRESS"
	AccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
)

var (
	server ports.RESTServer
)

func Start(ctx context.Context) {

	logger := logging.InitLogger()

	readEnv(logger)

	pgconn := viper.GetString(DatabaseURL)
	store, err := db.NewPostgresStore(ctx, pgconn, logger)
	if err != nil {
		logger.Fatal("store creation filed", err)
	}

	accrualconn := viper.GetString(AccrualSystemAddress)
	accrual := accrualdispatcher.NewGopherAccrualDispatcher(accrualconn, store, logger)

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
	err := server.Stop(context.Background())
	if err != nil {
		panic("app stop failed: " + err.Error())
	}
}

func readEnv(logger ports.Logger) {
	var (
		err        error
		addr       string
		dbURL      string
		accrualSys string
	)

	flag.StringVar(&addr, "a", "127.0.0.1:8080", "")
	flag.StringVar(&dbURL, "d", "postgres://postgres:qwerty@localhost:5432/gophermart?sslmode=disable", "")
	flag.StringVar(&accrualSys, "r", "127.0.0.1:8080", "")

	err = viper.BindEnv(DatabaseURL)
	if err != nil {
		logger.Fatal("database url env: %v", err)
	}
	viper.SetDefault(DatabaseURL, dbURL)

	err = viper.BindEnv(RunAddress)
	if err != nil {
		logger.Fatal("run address env: %v", err)
	}
	viper.SetDefault(RunAddress, addr)

	err = viper.BindEnv(AccrualSystemAddress)
	if err != nil {
		logger.Fatal("accrual system address env: %v", err)
	}
	viper.SetDefault(AccrualSystemAddress, accrualSys)
}
