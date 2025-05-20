package main

import (
	"beerbux/internal/api"
	"beerbux/internal/api/config"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	app, err := api.NewApp(cfg, logger)
	if err != nil {
		logger.Error("Failed to create API application", "error", err)
		os.Exit(1)
	}

	if err := app.Start(); err != nil {
		logger.Error("Failed to start API", "error", err)
		os.Exit(1)
	}
}
