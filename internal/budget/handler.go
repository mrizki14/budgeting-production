package budget

import (
	"errors"
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

func (h Handler) Index(c *gin.Context) {
	budgets, err := h.service.List(c.GetUint("userID"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch budgets", gin.H{"request": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Budgets fetched", budgets)
}

func (h Handler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", gin.H{"id": []string{"invalid id"}})
		return
	}

	budgetRecord, err := h.service.Get(c.GetUint("userID"), uint(id))
	if err != nil {
		if errors.Is(err, ErrBudgetForbidden) {
			response.Error(c, http.StatusForbidden, "Forbidden", gin.H{"budget": []string{err.Error()}})
			return
		}
		response.Error(c, http.StatusNotFound, "Budget not found", gin.H{"budget": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Budget fetched", budgetRecord)
}

func (h Handler) Store(c *gin.Context) {
	h.writeBudget(c, 0)
}

func (h Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", gin.H{"id": []string{"invalid id"}})
		return
	}

	h.writeBudget(c, uint(id))
}

func (h Handler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", gin.H{"id": []string{"invalid id"}})
		return
	}

	if err := h.service.Delete(c.GetUint("userID"), uint(id)); err != nil {
		if errors.Is(err, ErrBudgetForbidden) {
			response.Error(c, http.StatusForbidden, "Forbidden", gin.H{"budget": []string{err.Error()}})
			return
		}
		response.Error(c, http.StatusNotFound, "Budget not found", gin.H{"budget": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Budget deleted successfully", gin.H{"id": uint(id)})
}

func (h Handler) writeBudget(c *gin.Context, id uint) {
	var input struct {
		CategoryID  uint    `json:"category_id" binding:"required"`
		Month       int     `json:"month" binding:"required"`
		Year        int     `json:"year" binding:"required"`
		LimitAmount float64 `json:"limit_amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	var (
		budget *Budget
		err    error
	)
	if id == 0 {
		budget, err = h.service.Create(c.GetUint("userID"), input.CategoryID, input.Month, input.Year, input.LimitAmount)
	} else {
		budget, err = h.service.Update(c.GetUint("userID"), id, input.CategoryID, input.Month, input.Year, input.LimitAmount)
	}
	if err != nil {
		status := http.StatusUnprocessableEntity
		if errors.Is(err, ErrBudgetForbidden) {
			status = http.StatusForbidden
		}
		response.Error(c, status, "Validation failed", gin.H{"budget": []string{err.Error()}})
		return
	}

	message := "Budget created successfully"
	code := http.StatusCreated
	if id > 0 {
		message = "Budget updated successfully"
		code = http.StatusOK
	}

	response.Success(c, code, message, budget)
}
