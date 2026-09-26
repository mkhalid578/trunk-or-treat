package main

import (
	"car-api/internal/config"
	"car-api/internal/database"
	"car-api/internal/repository"
	"car-api/internal/scraper"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	if cfg.FirecrawlAPIKey == "" {
		return fmt.Errorf("FIRECRAWL_API_KEY is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client := scraper.NewFirecrawlClient(cfg.FirecrawlAPIKey)
	catalog, err := client.ScrapeToyotaCatalog(ctx)
	if err != nil {
		return fmt.Errorf("scrape Toyota catalog: %w", err)
	}
	cars, err := scraper.ValidateToyotaCatalog(catalog)
	if err != nil {
		return fmt.Errorf("validate Toyota catalog: %w", err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	if err := repository.UpsertCars(pool, cars); err != nil {
		return fmt.Errorf("import Toyota catalog: %w", err)
	}

	fmt.Printf("Imported or updated %d Toyota model/trim records.\n", len(cars))
	return nil
}
