package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	// Database
	DBURL  string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	// JWT / security
	JWTSecret string
	AppURL    string // e.g. https://yourapp.com (used for email verification links)

	// SMTP email
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string // FROM: noreply@yourapp.com
}

func Load() *Config {
	cfg := &Config{}

	// App settings
	cfg.Port = getEnv("PORT", "8080")
	cfg.AppURL = getEnv("APP_URL", "http://localhost:8080")

	// DB settings
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("DB_PORT", "5432")
	cfg.DBUser = getEnv("DB_USER", "admin")
	cfg.DBPass = getEnv("DB_PASS", "admin")
	cfg.DBName = getEnv("DB_NAME", "firstdb")

	// DB URL (can override the above entirely)
	cfg.DBURL = getEnv("DB_URL",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName),
	)

	// JWT
	cfg.JWTSecret = getEnv("JWT_SECRET", "dev-secret-change-me")

	// SMTP
	cfg.SMTPHost = getEnv("SMTP_HOST", "localhost")
	cfg.SMTPUser = getEnv("SMTP_USER", "")
	cfg.SMTPPass = getEnv("SMTP_PASS", "")
	cfg.SMTPFrom = getEnv("SMTP_FROM", "noreply@example.com")

	smtpPortStr := getEnv("SMTP_PORT", "1025")
	cfg.SMTPPort, _ = strconv.Atoi(smtpPortStr)

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
