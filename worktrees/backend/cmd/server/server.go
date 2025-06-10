package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openfoundry/oms/internal/config"
	"github.com/openfoundry/oms/internal/interfaces/rest"
	"go.uber.org/zap"
)

// createHTTPServer creates and configures HTTP server
func createHTTPServer(cfg *config.Config, services *rest.Services, logger *zap.Logger) *http.Server {
	router := rest.NewRouter(cfg, services, logger)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("HTTP server configured", zap.Int("port", cfg.Server.Port))
	return server
}

// startServer starts the HTTP server in a goroutine
func startServer(server *http.Server, logger *zap.Logger) {
	go func() {
		logger.Info("Starting HTTP server", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()
}

// waitForShutdownSignal waits for interrupt signal
func waitForShutdownSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// shutdownServer gracefully shuts down the server
func shutdownServer(server *http.Server, logger *zap.Logger) error {
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	logger.Info("Server exited gracefully")
	return nil
}