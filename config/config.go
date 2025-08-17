package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	SimServiceAPIKey     string
	SimUrlDefault        string
	AppPort              string
	JWTSecretKey         string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	DatabaseURL          string
	RedisURL             string
	RedisPassword        string
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Printf("no .env file found, assuming environment variables are set: %v", err)
	}

	cfg := &AppConfig{
		SimServiceAPIKey: os.Getenv("SERVICE_API_KEY"),
		SimUrlDefault:    os.Getenv("SERVICE_API_URL"),
		AppPort:          os.Getenv("APP_PORT"),
		JWTSecretKey:     os.Getenv("JWT_SECRET_KEY"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		RedisURL:         os.Getenv("REDIS_URL"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("XENDIT_API_KEY environment variable not set")
	}
	if cfg.SimServiceAPIKey == "" {
		log.Fatal("SIM_SERVICE_API_KEY environment variable not set")
	}
	if cfg.SimUrlDefault == "" {
		log.Fatal("SIM_API_URL_SERVICE environment variable not set")
	}
	if cfg.RedisURL == "" {
		log.Fatal("REDIS_URL environment variable not set")
	}

	return cfg
}
