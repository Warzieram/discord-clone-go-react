package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRetrieveMessagesParameterValidation(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Invalid limit parameter",
			queryParams: map[string]string{
				"limit":  "invalid",
				"offset": "0",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Invalid offset parameter",
			queryParams: map[string]string{
				"limit":  "10",
				"offset": "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing limit parameter",
			queryParams: map[string]string{
				"offset": "0",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing offset parameter",
			queryParams: map[string]string{
				"limit": "10",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
			q := url.Values{}
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			// Create response recorder
			rr := httptest.NewRecorder()

			// Test parameter parsing - should fail before database call
			RetrieveMessages(rr, req)

			// Check status code for parameter validation errors
			if tt.expectError && rr.Code != http.StatusBadRequest {
				t.Errorf("Expected bad request status for invalid parameters, got %d", rr.Code)
			}
		})
	}
}

func TestRetrieveMessagesParameterParsing(t *testing.T) {
	t.Run("Parameter extraction", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/messages?limit=5&offset=10", nil)

		limitParam := req.URL.Query().Get("limit")
		offsetParam := req.URL.Query().Get("offset")

		if limitParam != "5" {
			t.Errorf("Expected limit '5', got '%s'", limitParam)
		}

		if offsetParam != "10" {
			t.Errorf("Expected offset '10', got '%s'", offsetParam)
		}
	})

	t.Run("Missing parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/messages", nil)

		limitParam := req.URL.Query().Get("limit")
		offsetParam := req.URL.Query().Get("offset")

		if limitParam != "" {
			t.Errorf("Expected empty limit, got '%s'", limitParam)
		}

		if offsetParam != "" {
			t.Errorf("Expected empty offset, got '%s'", offsetParam)
		}
	})
}

// Note: Full integration tests would require database setup using testcontainers
// These tests focus on parameter validation and basic handler behavior
