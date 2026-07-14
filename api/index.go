package handler

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

	db := openDatabase(cfg)
	router = sharedRouter.New(db, cfg)
}

func openDatabase(cfg config.Config) *gorm.DB {
	dsn := databaseDSN(cfg)
	if dsn == "" {
		log.Println("WARNING: konfigurasi database belum disetting di Vercel Environment Variables")
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

func databaseDSN(cfg config.Config) string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return normalizeMySQLDSN(dsn)
	}

	if os.Getenv("DB_HOST") == "" || os.Getenv("DB_NAME") == "" || os.Getenv("DB_USER") == "" {
		return ""
	}

	return cfg.DSN()
}

func normalizeMySQLDSN(dsn string) string {
	if !strings.Contains(dsn, "://") {
		return dsn
	}

	parsed, err := url.Parse(dsn)
	if err != nil || parsed.Scheme != "mysql" {
		return dsn
	}

	password, _ := parsed.User.Password()
	query := parsed.Query()
	if query.Get("parseTime") == "" {
		query.Set("parseTime", "true")
	}

	return parsed.User.Username() + ":" + password + "@tcp(" + parsed.Host + ")/" + strings.TrimPrefix(parsed.Path, "/") + "?" + query.Encode()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if router == nil {
		http.Error(w, `{"error":"Router tidak terinisialisasi."}`, http.StatusInternalServerError)
		return
	}

	router.ServeHTTP(w, r)
}
