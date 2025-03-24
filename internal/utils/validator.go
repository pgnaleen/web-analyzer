package utils

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func ValidateRequest(w http.ResponseWriter, r *http.Request, logger *slog.Logger) (*http.Response, *url.URL) {
	urlQuery := r.URL.Query().Get("url")

	if urlQuery == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return nil, nil
	}

	parsedURL, err := url.Parse(urlQuery)
	if err != nil || !parsedURL.IsAbs() {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, nil
	}

	// Fetch the web page content
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(urlQuery)
	if err != nil {
		logger.Error("Failed to fetch URL", slog.String("error", err.Error()))
		http.Error(w, "Unable to fetch the URL", http.StatusInternalServerError)
		return nil, nil
	}

	return resp, parsedURL
}
