package main

import (
	"context"
	"gophermart/internal/adapters/logging"
	"gophermart/internal/application"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	logging.InitLogger()
	go application.Start(ctx)
	<-ctx.Done()
	application.Stop()
}
