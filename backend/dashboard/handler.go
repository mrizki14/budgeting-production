package dashboard

import (
	"net/http"

	"budgeting-app/golang/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) Show(c *gin.Context) {
	summary, err := h.service.Summary(c.GetUint("userID"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to build dashboard", gin.H{"dashboard": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Dashboard fetched", summary)
}
