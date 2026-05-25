package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Morpa/async-api/apiserver"
	"github.com/Morpa/async-api/config"
	"github.com/Morpa/async-api/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	conf, err := config.New()
	if err != nil {
		return err
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	db, err := store.NewPostgresDb(conf)
	if err != nil {
		return err
	}

	jwtManager := apiserver.NewJwtManager(conf)
	dataStore := store.New(db)
	server := apiserver.New(conf, logger, dataStore, jwtManager)
	if err := server.Start(ctx); err != nil {
		return err
	}
	return nil
}
