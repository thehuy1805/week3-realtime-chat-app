package config

import (
	"os"
)

type Config struct {
	PostgresURL string
	RedisAddr   string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		PostgresURL: getEnv("POSTGRES_URL", "postgres://postgres:0937491454az@localhost:5432/chat_app?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
