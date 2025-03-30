package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"web-analyzer/internal/analyzer"
	"web-analyzer/internal/monitoring"
	"web-analyzer/internal/utils"
)

// AnalyzeHandler processes the URL and returns analysis results
func AnalyzeHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request received", "request method", r.Method, "request URI", r.RequestURI, "body", r.Body)
		monitoring.RecordRequest()

		// Allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS Preflight Request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		resp, parsedURL := utils.ValidateRequest(w, r, logger)
		if resp == nil || parsedURL == nil {
			return
		}

		// closing the response body to avoid memory leaks
		// defer schedules this response close after function exited to save memory
		defer resp.Body.Close()

		// Call the AnalyzeHTML function
		analysis := analyzer.AnalyzeHTML(resp, parsedURL.Scheme+"://"+parsedURL.Host, logger, w)

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analysis)
	}
}

// RegisterRoutes sets up the API endpoints
func RegisterRoutes(logger *slog.Logger) {
	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", AnalyzeHandler(logger))
	http.ListenAndServe(":8080", mux)
}
