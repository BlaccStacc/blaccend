package config

import (
	"os"
)

type Config struct {
	Port string
	DBURL string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		DBURL: getEnv("DB_URL", "postgres://admin:admin@localhost:5432/firstdb"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

//ff nice ca aparent se pot apela functii post definite