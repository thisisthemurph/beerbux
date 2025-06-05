package api

import (
	"beerbux/internal/api/middleware"
	"beerbux/internal/api/webapp"
	"beerbux/internal/auth/command"
	authQueries "beerbux/internal/auth/db"
	authHandler "beerbux/internal/auth/handler"
	friendsHandler "beerbux/internal/friends/handler"
	sessionHandler "beerbux/internal/session/handler"
	"beerbux/internal/sse"
	streamHandler "beerbux/internal/streamer/handler"
	userHandler "beerbux/internal/user/handler"
	"beerbux/pkg/email"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Server struct {
	*http.Server
}

// NewServer creates a server for API endpoints.
// When running in development environment, this server also serves the frontend.
func (app *App) NewServer(streamServer *sse.Server) (*Server, error) {
	rootMux := http.NewServeMux()
	apiMux := http.NewServeMux()
	webMux := http.NewServeMux()

	emailSender := email.New(app.Config.Resend)
	if app.Config.Environment.IsDevelopment() && app.Config.Resend.DevelopmentSendToEmail == "" {
		emailSender = email.NewTerminalEmailLogger(app.Logger)
	}

	// Build and handle API routes
	apiMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	authHandler.BuildRoutes(app.Config, app.Logger, app.DB, emailSender, apiMux)
	sessionHandler.BuildRoutes(app.Logger, app.DB, apiMux, app.MessageReceiver())
	userHandler.BuildRoutes(app.Logger, app.DB, apiMux)
	friendsHandler.BuildRoutes(app.Logger, app.DB, apiMux)
	apiMux.Handle("/events/session", streamHandler.NewSessionTransactionCreatedHandler(app.Logger, streamServer))

	// Construct middleware for API routes
	authenticationQueries := authQueries.New(app.DB)
	refreshTokenCommand := command.NewRefreshTokenCommand(authenticationQueries, app.Config.GetAuthOptions())
	authMiddleware := middleware.NewAuthMiddleware(refreshTokenCommand, app.Config.Secrets.JWTSecret)
	recoverMiddleware := middleware.NewRecoverMiddleware(app.Logger)

	var apiHandler http.Handler
	if app.Config.Environment.IsDevelopment() {
		// We only want to run the CORS middleware when we run React separate
		// from the API in development mode.
		app.Logger.Info("Setting up API with CORS middleware")
		apiHandler = recoverMiddleware.Recover(
			authMiddleware.WithJWT(
				middleware.CORS(apiMux, app.Config.CORSClientBaseURL),
			),
		)
	} else {
		app.Logger.Info("Setting up API without CORS middleware")
		apiHandler = recoverMiddleware.Recover(
			authMiddleware.WithJWT(apiMux),
		)
	}

	// In development, we would like the swagger documentation to be available.
	if app.Config.Environment.IsDevelopment() {
		app.Logger.Info("Hosting swagger documentation")
		rootMux.Handle("/swagger/", httpSwagger.WrapHandler)
	}

	// In development, we want to run React webapp separate from the API.
	// In production, the React webapp will be served from the API.
	if !app.Config.Environment.IsDevelopment() {
		if err := webapp.New(webMux); err != nil {
			return nil, err
		}
		rootMux.Handle("/", webMux)
	}

	rootMux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	return &Server{
		Server: &http.Server{
			Addr:    app.Config.Address,
			Handler: rootMux,
		},
	}, nil
}
