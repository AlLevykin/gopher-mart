package main

import (
	"context"
	"gophermart/internal/application"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	app := application.NewGophermartApplication()
	go app.Start(ctx)
	<-ctx.Done()
	app.Stop()
}
