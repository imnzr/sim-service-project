package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	// XenditApiKey         string
	// XenditCallbackToken  string
	SimServiceAPIKey     string
	SimUrlDefault        string
	AppPort              string
	JWTSecretKey         string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	DatabaseURL          string
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Printf("no .env file found, assuming environment variables are set: %v", err)
	}

	cfg := &AppConfig{
		// XenditApiKey: os.Getenv("XENDIT_API_KEY"),
		// XenditCallbackToken: os.Getenv(),
		SimServiceAPIKey: os.Getenv("SIM_API_KEY_SERVICE"),
		SimUrlDefault:    os.Getenv("SIM_API_URL_SERVICE"),
		AppPort:          os.Getenv("APP_PORT"),
		JWTSecretKey:     os.Getenv("JWT_SECRET_KEY"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
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

	return cfg
}
