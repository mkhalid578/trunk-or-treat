package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode            string
	DatabaseURL     string
	Port            string
	FirecrawlAPIKey string
}

// Load loads the configuration from environment variables or other sources. caps makes it
// easy to know that is meant to be accessible from outside
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	config := &Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		Port:            os.Getenv("PORT"),
		Mode:            os.Getenv("MODE"),
		FirecrawlAPIKey: os.Getenv("FIRECRAWL_API_KEY"),
	}

	return config, nil
}
