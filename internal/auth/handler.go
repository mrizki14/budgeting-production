package auth

import (
	"errors"
	"net/http"

	"budgeting-app/golang/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	user, token, err := h.service.Register(input.Name, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"email": []string{err.Error()}})
			return
		}

		response.Error(c, http.StatusInternalServerError, "Register failed", gin.H{"request": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusCreated, "Register success", gin.H{"user": user, "token": token})
}

func (h Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	user, token, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Invalid credentials", gin.H{"email": []string{"invalid credentials"}})
			return
		}

		response.Error(c, http.StatusInternalServerError, "Login failed", gin.H{"request": []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Login success", gin.H{"user": user, "token": token})
}

func (h Handler) Me(c *gin.Context) {
	user, err := h.service.Me(c.GetUint("userID"))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", gin.H{"user": []string{"user not found"}})
		return
	}

	response.Success(c, http.StatusOK, "User fetched", user)
}

func (h Handler) UpdateProfile(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required,max=255"`
		Email string `json:"email" binding:"required,email,max=255"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	user, err := h.service.UpdateProfile(c.GetUint("userID"), input.Name, input.Email)
	if err != nil {
		field := "request"
		if errors.Is(err, ErrEmailExists) {
			field = "email"
		} else if errors.Is(err, ErrNameRequired) {
			field = "name"
		}
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{field: []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Profil berhasil diperbarui.", user)
}

func (h Handler) UpdatePassword(c *gin.Context) {
	var input struct {
		CurrentPassword      string `json:"current_password" binding:"required"`
		Password             string `json:"password" binding:"required,min=8"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{"request": []string{err.Error()}})
		return
	}

	err := h.service.UpdatePassword(c.GetUint("userID"), input.CurrentPassword, input.Password, input.PasswordConfirmation)
	if err != nil {
		field := "request"
		if errors.Is(err, ErrCurrentPassword) {
			field = "current_password"
		} else if errors.Is(err, ErrPasswordTooShort) || errors.Is(err, ErrPasswordConfirmation) {
			field = "password"
		}
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{field: []string{err.Error()}})
		return
	}

	response.Success(c, http.StatusOK, "Password berhasil diperbarui.", gin.H{})
}
