package health

import (
	"net/http"

	"budgeting-app/golang/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Ping(c *gin.Context) {
	response.Success(c, http.StatusOK, "API is running", gin.H{
		"status": "ok",
	})
}
