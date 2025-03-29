package server_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyzer/internal/server"
)

func TestAnalyzeHandler(t *testing.T) {
	logger := slog.Default()
	analyzeHandler := server.AnalyzeHandler(logger)

	t.Run("OPTIONS Request for CORS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/analyze", nil)
		w := httptest.NewRecorder()

		analyzeHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Direct Call to AnalyzeHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/analyze?url=https://example.com", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		analyzeHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response)
	})

	t.Run("Invalid Query param - Missing URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/analyze", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		analyzeHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
