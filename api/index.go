package handler

import (
	"log"
	"net/http"
	"os"

	"budgeting-app/golang/backend/shared/config"
	sharedRouter "budgeting-app/golang/backend/shared/router"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.Load()
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" && os.Getenv("APP_JWT_SECRET") == "" {
		cfg.JWTSecret = jwtSecret
	}

	db := openDatabase()
	router = sharedRouter.New(db, cfg)
}

func openDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("WARNING: DATABASE_URL belum disetting di Vercel Environment Variables")
		return nil
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("WARNING: gagal terhubung ke database: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("WARNING: gagal mengambil sql.DB dari GORM: %v", err)
		return db
	}

	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(1)

	return db
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if router == nil {
		http.Error(w, `{"error":"Router tidak terinisialisasi."}`, http.StatusInternalServerError)
		return
	}

	router.ServeHTTP(w, r)
}
