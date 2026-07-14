package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"budgeting-app/golang/backend/shared/config"
	"budgeting-app/golang/backend/shared/router"
)

func TestHealthPing(t *testing.T) {
	r := router.New(nil, config.Config{JWTSecret: "test-secret"})
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
