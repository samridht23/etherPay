package portal

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func InitializeServer(router *chi.Mux) {

	server := &http.Server{
		Addr:         ":" + os.Getenv("SERVER_PORT"),
		Handler:      router,
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	go func() {
		zap.S().Info("Server listening on port: ", os.Getenv("SERVER_PORT"))
		err := server.ListenAndServe()
		if err != nil {
			zap.S().Error(err)
		}
	}()
	<-sigCh
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	zap.S().Info("Gracefully shutting down server")
	if err := server.Shutdown(ctx); err != nil {
		zap.S().Errorf("Failed to gracefully shutdown server: %v", err)
		zap.S().Info("Attempting to close server forcefully")
		// Attempt to close the server forcefully if shutdown fails
		if err := server.Close(); err != nil {
			zap.S().Errorf("Failed to close server's underlying listener: %v", err)
		}
		// Exit with an error status
		os.Exit(1)
	}
	// Log successful graceful shutdown
	zap.S().Info("Server shutdown gracefully")
	// Exit with success status
	os.Exit(0)
}
