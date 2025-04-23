package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	RedisAddr  string
	JWTSecret  string
	RateLimit  int // Maximum rate limit (requests per window)
	WindowSec  int // Window duration in seconds
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8000"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		RateLimit:  getIntEnv("RATE_LIMIT", 5),
		WindowSec:  getIntEnv("WINDOW_SEC", 60),
	}
}

// getEnv retrieves a string environment variable value or returns a fallback default if it's not set
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// getIntEnv retrieves an integer environment variable value or returns a fallback default if it's not set
// It logs an error if the conversion fails and returns the fallback value
func getIntEnv(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	// Convert string to int, log error if conversion fails
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Invalid value for %s, using fallback: %d", key, fallback)
		return fallback
	}
	return intVal
}
