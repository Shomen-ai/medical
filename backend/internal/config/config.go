// Файл: internal/config/config.go
// Назначение: чтение конфигурации приложения из переменных окружения (DATABASE_URL, REDIS_URL, JWT_SECRET, SMSC-учётные данные).
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит все runtime-настройки сервиса BeautyMed.
type Config struct {
	DatabaseURL  string
	RedisURL     string
	JWTSecret    string
	SMSCLogin    string
	SMSCPassword string
	Port         string
	Env          string
}

// Load читает .env (если есть) и собирает структуру Config из переменных окружения.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}
	return &Config{
		DatabaseURL:  mustEnv("DATABASE_URL"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:    mustEnv("JWT_SECRET"),
		SMSCLogin:    getEnv("SMSC_LOGIN", ""),
		SMSCPassword: getEnv("SMSC_PASSWORD", ""),
		Port:         getEnv("PORT", "8080"),
		Env:          getEnv("ENV", "development"),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
