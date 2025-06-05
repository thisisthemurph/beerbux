package main

import (
	"beerbux/cmd/api/docs"
	"beerbux/internal/api"
	"beerbux/internal/api/config"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Beerbux API
// @version 1.0
// @description API for the Beerbux application.
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

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost%s", app.Config.Address)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	if err := app.Start(); err != nil {
		logger.Error("Failed to start API", "error", err)
		os.Exit(1)
	}
}
