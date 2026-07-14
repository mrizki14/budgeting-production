package category

import (
	"errors"
	"net/http"
	"strconv"

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
	categories, err := h.service.List(c.GetUint("userID"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch categories", gin.H{"request": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Categories fetched", categories)
}

func (h Handler) Show(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	categoryRecord, err := h.service.Get(c.GetUint("userID"), id)
	if err != nil {
		if errors.Is(err, ErrCategoryForbidden) {
			response.Error(c, http.StatusForbidden, "Forbidden", gin.H{"category": []string{err.Error()}})
			return
		}
		response.Error(c, http.StatusNotFound, "Category not found", gin.H{"category": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Category fetched", categoryRecord)
}

func (h Handler) Store(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	category, err := h.service.Create(c.GetUint("userID"), input.Name, input.Type)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"type": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusCreated, "Category created successfully", category)
}

func (h Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	category, err := h.service.Update(c.GetUint("userID"), id, input.Name, input.Type)
	if err != nil {
		if errors.Is(err, ErrCategoryForbidden) {
			response.Error(c, http.StatusForbidden, "Forbidden", gin.H{"category": []string{err.Error()}})
			return
		}
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"category": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Category updated successfully", category)
}

func (h Handler) Destroy(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.GetUint("userID"), id); err != nil {
		if errors.Is(err, ErrCategoryForbidden) {
			response.Error(c, http.StatusForbidden, "Forbidden", gin.H{"category": []string{err.Error()}})
			return
		}
		response.Error(c, http.StatusNotFound, "Category not found", gin.H{"category": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Category deleted successfully", gin.H{"id": id})
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", gin.H{"id": []string{"invalid id"}})
		return 0, false
	}

	return uint(id), true
}
