package firebase

import (
	"os"
)

type Config struct {
	ProjectID       string
	DatabaseURL     string
	CredentialsFile string
}

func LoadConfig() *Config {
	return &Config{
		ProjectID:   getEnv("FIREBASE_PROJECT_ID", "ferrous-cogency-215410"),
		DatabaseURL: getEnv("FIREBASE_DATABASE_URL", "https://s-usc1b-nss-2136.firebaseio.com"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
