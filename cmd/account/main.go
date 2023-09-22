package main

import (
	"context"
	"errors"
	"github.com/Hvaekar/med-account/cmd/account/commands"
	"github.com/Hvaekar/med-account/pkg/server/ginmiddleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("failed load .env: %s", err.Error())
	}
}

func main() {
	ginmiddleware.MustRegisterMetrics(prometheus.DefaultRegisterer)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := commands.NewApp().RunContext(ctx, os.Args); err != nil && !errors.Is(err, context.Canceled) {
		log.Panicf("failed executing command: %s", err)
	}
}
