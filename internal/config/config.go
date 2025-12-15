package config

import (
	"os"
)

type Config struct {
	HTTPPort string
	DBPath   string
}

func Load() *Config {
	return &Config{
		HTTPPort: getEnv("HTTP_PORT", ":8098"),
		DBPath:   getEnv("DB_PATH", "data/app.db"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
