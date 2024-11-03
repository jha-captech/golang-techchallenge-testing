package handlers

import (
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHealthCheck(t *testing.T) {
	tests := map[string]struct {
		wantStatus int
		wantBody   string
	}{
		"happy path": {
			wantStatus: 200,
			wantBody:   `{"status":"ok"}`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new request
			req := httptest.NewRequest("GET", "/health", nil)

			// Create a new response recorder
			rec := httptest.NewRecorder()

			// Create a new logger
			logger := slog.Default()

			// Call the handler
			HandleHealthCheck(logger)(rec, req)

			// Check the status code
			if rec.Code != tc.wantStatus {
				t.Errorf("want status %d, got %d", tc.wantStatus, rec.Code)
			}

			// Check the body
			if strings.Trim(rec.Body.String(), "\n") != tc.wantBody {
				t.Errorf("want body %q, got %q", tc.wantBody, rec.Body.String())
			}
		})
	}
}
