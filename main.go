package main

import (
	"context"
	"log/slog" // logrus - Slower due to reflection & interfaces slog - Faster due to structured logging with zero allocations
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	_ "web-analyzer/internal/analyzer"
	"web-analyzer/internal/server"
	"web-analyzer/internal/utils"
)

var (
	serverInstance *http.Server
	serverOnce     sync.Once
)

func initServer(logger *slog.Logger) *http.Server {
	serverOnce.Do(func() {
		r := chi.NewRouter()
		server.RegisterRoutes(r, logger)
		serverInstance = &http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 25 * time.Second,
		}
	})
	return serverInstance
}

func main() {
	logger := utils.InitLogger()
	logger.Info("Starting Web Analyzer API...")

	srv := initServer(logger)
	go func() {
		logger.Info("Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	shutdownComplete := make(chan struct{})
	go func() {
		<-quit
		logger.Info("Shutting down server_test...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown failed", slog.String("error", err.Error()))
		}
		close(shutdownComplete)
	}()

	<-shutdownComplete
	logger.Info("Server exited properly")
}
