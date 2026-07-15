package router

import (
	"context"
	"net/http"
	"strings"
	"time"

	"budgeting-app/golang/backend/auth"
	"budgeting-app/golang/backend/budget"
	"budgeting-app/golang/backend/category"
	"budgeting-app/golang/backend/dashboard"
	"budgeting-app/golang/backend/report"
	"budgeting-app/golang/backend/shared/config"
	"budgeting-app/golang/backend/shared/health"
	"budgeting-app/golang/backend/shared/middleware"
	"budgeting-app/golang/backend/shared/response"
	"budgeting-app/golang/backend/transaction"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(db *gorm.DB, cfg config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(cfg.CORSAllowedOrigin))

	healthHandler := health.NewHandler()
	r.GET("/api/health", healthHandler.Ping)
	r.GET("/api/ready", databaseReady(db))

	if db == nil {
		registerDatabaseUnavailableRoutes(r)
		return r
	}

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	categoryRepo := category.NewRepository(db)
	categoryHandler := category.NewHandler(category.NewService(categoryRepo))
	transactionHandler := transaction.NewHandler(transaction.NewService(transaction.NewRepository(db), categoryRepo))
	budgetHandler := budget.NewHandler(budget.NewService(budget.NewRepository(db), categoryRepo))
	reportHandler := report.NewHandler(report.NewService(report.NewRepository(db)))
	dashboardHandler := dashboard.NewHandler(dashboard.NewService(dashboard.NewRepository(db)))

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	protected := r.Group("/api")
	protected.Use(middleware.RequireAuth(cfg.JWTSecret))
	protected.GET("/auth/me", authHandler.Me)
	protected.PUT("/settings/profile", authHandler.UpdateProfile)
	protected.PUT("/settings/password", authHandler.UpdatePassword)
	protected.GET("/dashboard", dashboardHandler.Show)
	protected.GET("/categories", categoryHandler.Index)
	protected.GET("/categories/:id", categoryHandler.Show)
	protected.POST("/categories", categoryHandler.Store)
	protected.PUT("/categories/:id", categoryHandler.Update)
	protected.DELETE("/categories/:id", categoryHandler.Destroy)
	protected.GET("/transactions", transactionHandler.Index)
	protected.GET("/transactions/:id", transactionHandler.Show)
	protected.POST("/transactions", transactionHandler.Store)
	protected.PUT("/transactions/:id", transactionHandler.Update)
	protected.DELETE("/transactions/:id", transactionHandler.Destroy)
	protected.GET("/budgets", budgetHandler.Index)
	protected.GET("/budgets/:id", budgetHandler.Show)
	protected.POST("/budgets", budgetHandler.Store)
	protected.PUT("/budgets/:id", budgetHandler.Update)
	protected.DELETE("/budgets/:id", budgetHandler.Destroy)
	protected.GET("/reports/summary", reportHandler.Summary)

	return r
}

func databaseReady(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			response.Error(c, http.StatusServiceUnavailable, "Database is not configured or unavailable", gin.H{"database": []string{"database connection is not initialized"}})
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			response.Error(c, http.StatusServiceUnavailable, "Database is unavailable", gin.H{"database": []string{err.Error()}})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			response.Error(c, http.StatusServiceUnavailable, "Database is unavailable", gin.H{"database": []string{err.Error()}})
			return
		}

		response.Success(c, http.StatusOK, "Database is ready", gin.H{"database": "ok"})
	}
}

func registerDatabaseUnavailableRoutes(r *gin.Engine) {
	unavailable := func(c *gin.Context) {
		response.Error(c, http.StatusServiceUnavailable, "Database is not configured or unavailable", gin.H{"database": []string{"check DATABASE_URL and Vercel function logs"}})
	}

	r.POST("/api/auth/register", unavailable)
	r.POST("/api/auth/login", unavailable)

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			unavailable(c)
			return
		}

		c.Status(http.StatusNotFound)
	})
}
