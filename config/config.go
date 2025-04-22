package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ServerPort string
	RedisAddr  string
	JWTSecret  string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8000"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
