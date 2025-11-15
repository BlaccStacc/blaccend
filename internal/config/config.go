package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "admin")
	pass := getEnv("DB_PASS", "admin")
	name := getEnv("DB_NAME", "firstdb")

	dbURL := getEnv("DB_URL",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, pass, host, port, name,
		),
	)

	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBURL: dbURL,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
