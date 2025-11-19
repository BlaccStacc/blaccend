package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	DBURL  string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	StorageEndpoint  string
	StorageRegion    string
	StorageAccessKey string
	StorageSecretKey string
	StorageBucket    string

	JWTSecret string
	AppURL    string

	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string
}

var cfg *Config

func Init() {
	c := &Config{}

	c.Port = getEnv("PORT", "8080")
	c.AppURL = getEnv("APP_URL", "http://localhost:8080")

	c.DBHost = getEnv("DB_HOST", "localhost")
	c.DBPort = getEnv("DB_PORT", "5432")
	c.DBUser = getEnv("DB_USER", "admin")
	c.DBPass = getEnv("DB_PASS", "admin")
	c.DBName = getEnv("DB_NAME", "firstdb")

	c.StorageEndpoint = getEnv("S3_ENDPOINT", "http://garage:3900")
	c.StorageRegion = getEnv("S3_REGION", "garage")
	c.StorageAccessKey = getEnv("S3_ACCESS_KEY", "admin-token")
	c.StorageSecretKey = getEnv("S3_SECRET_KEY", "admin-token")
	c.StorageBucket = getEnv("S3_BUCKET", "app")

	c.DBURL = getEnv("DB_URL",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName),
	)

	c.JWTSecret = getEnv("JWT_SECRET", "dev-secret-change-me")

	c.SMTPHost = getEnv("SMTP_HOST", "localhost")
	c.SMTPUser = getEnv("SMTP_USER", "")
	c.SMTPPass = getEnv("SMTP_PASS", "")
	c.SMTPFrom = getEnv("SMTP_FROM", "noreply@example.com")

	smtpPortStr := getEnv("SMTP_PORT", "1025")
	c.SMTPPort, _ = strconv.Atoi(smtpPortStr)

	cfg = c
}

func Load() *Config {
	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
