package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv            string
	AppPort           string
	JWTSecret         string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPass            string
	DBParams          string
	CORSAllowedOrigin string
}

func Load() Config {
	return Config{
		AppEnv:            envOr("APP_ENV", "development"),
		AppPort:           envOr("APP_PORT", "8080"),
		JWTSecret:         envOr("APP_JWT_SECRET", "change-me-in-local"),
		DBHost:            envOr("DB_HOST", "127.0.0.1"),
		DBPort:            envOr("DB_PORT", "3306"),
		DBName:            envOr("DB_NAME", "budgeting_app"),
		DBUser:            envOr("DB_USER", "root"),
		DBPass:            os.Getenv("DB_PASSWORD"),
		DBParams:          envOr("DB_PARAMS", "parseTime=true"),
		CORSAllowedOrigin: envOr("CORS_ALLOWED_ORIGIN", "http://127.0.0.1:5173"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName, c.DBParams)
}

func envOr(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
