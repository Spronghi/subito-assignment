package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/handler"
)

func newTestHealthMux(t *testing.T) *http.ServeMux {
	t.Helper()

	mux := http.NewServeMux()

	handler.NewHealthHandler().RegisterRoutes(mux)

	return mux
}

func TestHealthHandler(t *testing.T) {
	mux := newTestHealthMux(t)

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Body.String() != "OK" {
		t.Errorf("Expected body '%s', got '%s'", "OK", rec.Body.String())
	}
}
