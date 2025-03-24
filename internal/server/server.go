package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"net/url"
	"time"
	"web-analyzer/internal/analyzer"
)

// AnalyzeHandler processes the URL and returns analysis results
func AnalyzeHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlQuery := r.URL.Query().Get("url")
		if urlQuery == "" {
			http.Error(w, "Missing URL parameter", http.StatusBadRequest)
			return
		}

		parsedURL, err := url.Parse(urlQuery)
		if err != nil || !parsedURL.IsAbs() {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		// Fetch the web page content
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(urlQuery)
		if err != nil {
			logger.Error("Failed to fetch URL", slog.String("error", err.Error()))
			http.Error(w, "Unable to fetch the URL", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Call the AnalyzeHTML function
		analysis := analyzer.AnalyzeHTML(resp, parsedURL.Scheme+"://"+parsedURL.Host, logger, w)

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analysis)
	}
}

// RegisterRoutes sets up the API endpoints
func RegisterRoutes(r *chi.Mux, logger *slog.Logger) {
	r.Get("/analyze", AnalyzeHandler(logger))
}
