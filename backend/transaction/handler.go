package transaction

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"budgeting-app/golang/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) Index(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	transactions, err := h.service.List(c.GetUint("userID"), c.Query("type"), uint(categoryID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch transactions", gin.H{"request": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Transactions fetched", transactions)
}

func (h Handler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", gin.H{"id": []string{"invalid id"}})
		return
	}

	transactionRecord, err := h.service.Get(c.GetUint("userID"), uint(id))
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Error(c, http.StatusForbidden, "Forbidden", gin.H{"transaction": []string{err.Error()}})
			return
		}
		response.Error(c, http.StatusNotFound, "Transaction not found", gin.H{"transaction": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Transaction fetched", transactionRecord)
}

func (h Handler) Store(c *gin.Context) {
	h.writeTransaction(c, 0)
}

func (h Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", gin.H{"id": []string{"invalid id"}})
		return
	}

	h.writeTransaction(c, uint(id))
}

func (h Handler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", gin.H{"id": []string{"invalid id"}})
		return
	}

	if err := h.service.Delete(c.GetUint("userID"), uint(id)); err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Error(c, http.StatusForbidden, "Forbidden", gin.H{"transaction": []string{err.Error()}})
			return
		}
		response.Error(c, http.StatusNotFound, "Transaction not found", gin.H{"transaction": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Transaction deleted successfully", gin.H{"id": uint(id)})
}

func (h Handler) writeTransaction(c *gin.Context, id uint) {
	var input struct {
		CategoryID  uint    `json:"category_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Type        string  `json:"type" binding:"required"`
		Date        string  `json:"date" binding:"required"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"date": []string{"date must use YYYY-MM-DD"}})
		return
	}

	var transaction *Transaction
	if id == 0 {
		transaction, err = h.service.Create(c.GetUint("userID"), input.CategoryID, input.Amount, input.Type, date, input.Description)
	} else {
		transaction, err = h.service.Update(c.GetUint("userID"), id, input.CategoryID, input.Amount, input.Type, date, input.Description)
	}
	if err != nil {
		status := http.StatusUnprocessableEntity
		if errors.Is(err, ErrForbidden) {
			status = http.StatusForbidden
		}
		response.Error(c, status, "Validation failed", gin.H{"transaction": []string{err.Error()}})
		return
	}

	message := "Transaction created successfully"
	code := http.StatusCreated
	if id > 0 {
		message = "Transaction updated successfully"
		code = http.StatusOK
	}
	response.Success(c, code, message, transaction)
}
