package dashboard

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestShowReturnsDashboardEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(NewService(dashboardRepository{income: 500, expenses: 200}))
	router := gin.New()
	router.GET("/api/dashboard", func(c *gin.Context) {
		c.Set("userID", uint(4))
		handler.Show(c)
	})
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/dashboard", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"total_balance":300`) {
		t.Fatalf("unexpected body: %s", recorder.Body.String())
	}
}
