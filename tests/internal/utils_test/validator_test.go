package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"web-analyzer/internal/utils"
)

func TestValidateRequest(t *testing.T) {
	logger := utils.InitLogger()

	// invalid http method to validate first logic of validator
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/validate", nil)
		rec := httptest.NewRecorder()

		resp, parsedURL := utils.ValidateRequest(rec, req, logger)
		assert.Nil(t, resp)
		assert.Nil(t, parsedURL)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	// request without url query parameter to check query param missing check
	t.Run("Missing URL Parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/validate", nil)
		rec := httptest.NewRecorder()

		resp, parsedURL := utils.ValidateRequest(rec, req, logger)
		assert.Nil(t, resp)
		assert.Nil(t, parsedURL)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// request with invalid url to check url parser
	t.Run("Invalid URL Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/validate?url=invalid_url", nil)
		rec := httptest.NewRecorder()

		resp, parsedURL := utils.ValidateRequest(rec, req, logger)
		assert.Nil(t, resp)
		assert.Nil(t, parsedURL)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// request with invalid url domain to check regex
	t.Run("Invalid URL Format 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/validate?url=ftp://example.com", nil)
		rec := httptest.NewRecorder()

		resp, parsedURL := utils.ValidateRequest(rec, req, logger)
		assert.Nil(t, resp)
		assert.Nil(t, parsedURL)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// request with invalid url domain to check regex
	t.Run("Invalid URL Format 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/validate?url=https://-example.com", nil)
		rec := httptest.NewRecorder()

		resp, parsedURL := utils.ValidateRequest(rec, req, logger)
		assert.Nil(t, resp)
		assert.Nil(t, parsedURL)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Failed to fetch URL
	t.Run("Failed to fetch URL", func(t *testing.T) {
		// Mock server to return a 200 response
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close() // defer schedules this mock server close after function exited to save memory

		req := httptest.NewRequest(http.MethodGet, "/validate?url=http://localhost.com", nil)
		rec := httptest.NewRecorder()

		resp, parsedURL := utils.ValidateRequest(rec, req, logger)
		assert.Nil(t, resp)
		assert.Nil(t, parsedURL)
	})

	// successful request to check happy path
	t.Run("Successful Request", func(t *testing.T) {
		// Mock server to return a 200 response
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close() // defer schedules this mock server close after function exited to save memory

		req := httptest.NewRequest(http.MethodGet, "/validate?url="+mockServer.URL, nil)
		rec := httptest.NewRecorder()

		resp, parsedURL := utils.ValidateRequest(rec, req, logger)
		assert.NotNil(t, resp)
		assert.NotNil(t, parsedURL)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
