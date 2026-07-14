package middleware

import (
	"net/http"
	"strings"

	authshared "budgeting-app/golang/internal/shared/auth"
	"budgeting-app/golang/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "Unauthorized", gin.H{"token": []string{"missing bearer token"}})
			c.Abort()
			return
		}

		claims, err := authshared.ParseToken(secret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Unauthorized", gin.H{"token": []string{"invalid token"}})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
