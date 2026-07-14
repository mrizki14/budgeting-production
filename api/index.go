package handler

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"budgeting-app/golang/backend/shared/config"
	sharedRouter "budgeting-app/golang/backend/shared/router"

	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	gormMySQL "gorm.io/driver/mysql"
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

	registerDatabaseTLS(dsn)

	db, err := gorm.Open(gormMySQL.Open(dsn), &gorm.Config{})
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
	if strings.Contains(parsed.Hostname(), "tidbcloud.com") && (query.Get("tls") == "" || query.Get("tls") == "true") {
		query.Set("tls", "tidb")
	}

	return parsed.User.Username() + ":" + password + "@tcp(" + parsed.Host + ")/" + strings.TrimPrefix(parsed.Path, "/") + "?" + query.Encode()
}

func registerDatabaseTLS(dsn string) {
	dbConfig, err := mysqlDriver.ParseDSN(dsn)
	if err != nil || dbConfig.TLSConfig != "tidb" {
		return
	}

	serverName := dbConfig.Addr
	if host, _, err := net.SplitHostPort(dbConfig.Addr); err == nil {
		serverName = host
	}
	if serverName == "" {
		return
	}

	if err := mysqlDriver.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: serverName,
	}); err != nil && !strings.Contains(err.Error(), "already registered") {
		log.Printf("WARNING: gagal mendaftarkan TLS config TiDB: %v", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if router == nil {
		http.Error(w, `{"error":"Router tidak terinisialisasi."}`, http.StatusInternalServerError)
		return
	}

	router.ServeHTTP(w, r)
}
