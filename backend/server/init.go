package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/auth"
	"github.com/meh-hackathon/meh/config"
	"github.com/meh-hackathon/meh/httpx"
	"github.com/meh-hackathon/meh/logger"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config ../codegen-conf.yml ../openapi.yml

var (
	ErrNotFound            = apperror.Define("api:not_found", "API endpoint not found").WithStatus(http.StatusNotFound)
	ErrInternalServerError = apperror.Define("api:internal_server_error", "Internal Server Error").WithStatus(http.StatusInternalServerError)
	ErrForbidden           = apperror.Define("api:forbidden", "Access denied").WithStatus(http.StatusForbidden)
)

type Server struct {
	oauthHandler *auth.OAuthHandler
	HttpServer   *http.Server
	ApiMux       *http.ServeMux
}

func New() (*Server, error) {
	mux := http.NewServeMux()
	srv := &Server{}

	// UI
	uiHandler, err := GetUiHandler()
	if err != nil {
		return nil, fmt.Errorf("Failed to get UI handler: %w", err)
	}
	mux.Handle("/", uiHandler)

	// health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// api
	apiMux := http.NewServeMux()
	apiMux.Handle("/", Handler(srv))
	mux.Handle("/api/", http.StripPrefix("/api", auth.Middleware(config.Secret)(apiMux)))

	// auth
	srv.oauthHandler, err = auth.NewOAuthHandler(auth.WithLocalAuthenticator)
	if err != nil {
		return nil, err
	}

	addr := "0.0.0.0:"
	if runtime.GOOS == "windows" {
		addr = "127.0.0.1:"
	}

	srv.HttpServer = &http.Server{
		Addr:    addr + fmt.Sprintf("%d", config.Port),
		Handler: CorsMiddleware(mux),
	}
	srv.ApiMux = apiMux
	return srv, nil
}

func (srv *Server) Start() {
	go func() {
		logger.Info("Starting HTTP server", "address", srv.HttpServer.Addr)
		err := srv.HttpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.HttpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP Server shutdown failed: %w", err)
	}

	logger.Debug("HTTP Server gracefully stopped")
	return nil
}
