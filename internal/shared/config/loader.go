package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string
	GRPCPort    string
	JWTSecret   string
	APIKey      string
	Env         string
}

// Load reads .env and system variables into Config struct
func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		ServiceName: getEnv("SERVICE_NAME", "unarya-service"),
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecret"),
		APIKey:      getEnv("API_KEY", ""),
		Env:         getEnv("ENV", "development"),
	}

	log.Printf("[Config] Loaded for service: %s", cfg.ServiceName)
	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
