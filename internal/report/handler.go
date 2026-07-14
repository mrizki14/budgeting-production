package report

import (
	"net/http"
	"strconv"

	"budgeting-app/golang/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) Summary(c *gin.Context) {
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	summary, err := h.service.Summary(c.GetUint("userID"), month, year)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to build report", gin.H{"report": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Report fetched", summary)
}
