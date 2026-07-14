package handler

import (
	"log"
	"net/http"
	"os"

	// Import module auth Anda (sesuaikan path dengan module di go.mod)
	"budgeting-app/golang/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Instance Gin global agar tidak dibuat ulang pada setiap request (optimasi serverless)
var router *gin.Engine

// Fungsi init() akan dijalankan otomatis saat serverless function Vercel di-start
func init() {
	// 1. Mode rilis agar log tidak terlalu berantakan di console Vercel
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(gin.Recovery())

	// 2. Setup Koneksi Database (menggunakan environment variable dari Vercel)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("WARNING: DATABASE_URL belum disetting di Vercel Environment Variables")
		// Jika DSN kosong, kita hentikan inisialisasi route database untuk mencegah panic,
		// tapi router tetap jalan agar API bisa merespon error dengan anggun.
		return
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Optimasi koneksi GORM khusus untuk lingkungan Serverless
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(2) // Jangan terlalu besar agar batas TiDB tidak penuh
		sqlDB.SetMaxIdleConns(1)
	}

	// 3. Setup Dependency Injection (Asumsi nama fungsi konstruktor standar)
	// Jika nama konstruktor Repository/Service Anda berbeda, sesuaikan di bawah ini:
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// 4. Setup Routing
	apiGroup := router.Group("/api")
	{
		// Cek status API
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "API is running on Vercel"})
		})

		// Group Auth
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)

			// Catatan: Jika Anda menggunakan middleware JWT untuk autentikasi,
			// Anda bisa menyisipkannya di sini. Contoh:
			// authGroup.Use(auth.JWTMiddleware()) 
			
			authGroup.GET("/me", authHandler.Me)
			authGroup.PUT("/profile", authHandler.UpdateProfile)
			authGroup.PUT("/password", authHandler.UpdatePassword)
		}
	}
}

// Handler ini adalah jembatan utama yang wajib ada. 
// Vercel mengharapkan signature standar Go (http.ResponseWriter, *http.Request)
func Handler(w http.ResponseWriter, r *http.Request) {
	if router == nil {
		http.Error(w, `{"error": "Router tidak terinisialisasi. Periksa DATABASE_URL."}`, http.StatusInternalServerError)
		return
	}

	// Gin mengambil alih HTTP request mentah dari Vercel
	router.ServeHTTP(w, r)
}