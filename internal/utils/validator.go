package utils

import (
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func ValidateRequest(w http.ResponseWriter, r *http.Request, logger *slog.Logger) (*http.Response, *url.URL) {
	if r.Method != "GET" {
		logger.Error("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return nil, nil
	}

	urlQuery := r.URL.Query().Get("url")
	if urlQuery == "" {
		logger.Error("Missing URL parameter")
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return nil, nil
	}

	// Ensures the URL is correctly structured and Ensures the URL contains a valid scheme (http or https)
	parsedURL, err := url.Parse(urlQuery)
	if err != nil || !parsedURL.IsAbs() {
		logger.Error("Invalid URL format", slog.String("error", err.Error()))
		http.Error(w, "{\"error_code\":\"9001\",\"error_message\":\"Invalid URL format\"}", http.StatusBadRequest)
		return nil, nil
	}

	// ensure domain structure, ports, and paths
	re := regexp.MustCompile(`^(https?://)?[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+(:\d+)?(/.*)?$`)
	if !re.MatchString(urlQuery) {
		logger.Error("URL format does not match the expected pattern")
		http.Error(w, "{\"error_code\":\"9002\",\"error_message\":\"URL format does not match the expected pattern\"}", http.StatusBadRequest)
		return nil, nil
	}

	// Fetch the web page content
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(urlQuery)
	if err != nil {
		logger.Error("Failed to fetch URL", slog.String("error", err.Error()))
		http.Error(w, "{\"error_code\":\"9003\",\"error_message\":\"Unable to fetch the URL\"}", http.StatusInternalServerError)
		return nil, nil
	}

	return resp, parsedURL
}
