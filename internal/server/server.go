package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"web-analyzer/internal/analyzer"
	"web-analyzer/internal/utils"
)

// AnalyzeHandler processes the URL and returns analysis results
func AnalyzeHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, parsedURL := utils.ValidateRequest(w, r, logger)
		if resp == nil || parsedURL == nil {
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
