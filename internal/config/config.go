package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode        string
	DatabaseURL string
	Port        string
}

// Load loads the configuration from environment variables or other sources. caps makes it
// easy to know that is meant to be accessible from outside
func Load() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		Mode:        os.Getenv("MODE"),
	}

	return config, nil
}
