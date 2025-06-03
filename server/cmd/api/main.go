package main

import (
	"beerbux/internal/api"
	"beerbux/internal/api/config"
	"log/slog"
	"net/http"
	"os"

	_ "beerbux/cmd/api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Beerbux API
// @version 1.0
// @description API for the Beerbux application.
// @host localhost:42069
// @BasePath /api
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.LogLevel,
	}))

	app, err := api.NewApp(cfg, logger)
	if err != nil {
		logger.Error("Failed to create API application", "error", err)
		os.Exit(1)
	}

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	if err := app.Start(); err != nil {
		logger.Error("Failed to start API", "error", err)
		os.Exit(1)
	}
}
