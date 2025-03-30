package main

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog" // logrus - Slower due to reflection & interfaces slog - Faster due to structured logging with zero allocations
	"net/http"
	"os"
	"sync"
	"time"
	"web-analyzer/internal/monitoring"

	"github.com/go-chi/chi/v5"
	_ "net/http/pprof" // Enables pprof handlers

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
		server.RegisterRoutes(logger)
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

	// Start a separate server for pprof using goroutine
	go func() {
		logger.Info("Starting pprof server on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil &&
			!errors.Is(err, http.ErrServerClosed) { // Default pprof routes
			logger.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	monitoring.InitMetrics()
	// Start separate Prometheus metrics endpoint using goroutine
	go func() {
		logger.Info("Starting Prometheus metrics server on :9090")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil &&
			!errors.Is(err, http.ErrServerClosed) { // Default monitoring route
			logger.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	srv := initServer(logger)

	quit := make(chan os.Signal, 1)
	shutdownComplete := make(chan struct{})
	go func() {
		<-quit
		logger.Info("Shutting down server gracefully...")
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
