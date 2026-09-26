package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"car-api/internal/models"
)

const toyotaCatalogURL = "https://www.toyota.com/all-vehicles/"

type ToyotaCar struct {
	Model      string `json:"model"`
	ModelYear  int    `json:"model_year"`
	Trim       string `json:"trim"`
	BodyStyle  string `json:"body_style"`
	Powertrain string `json:"powertrain"`
}

type FirecrawlClient struct {
	APIKey     string
	HTTPClient *http.Client
	Endpoint   string
}

var allowedBodyStyles = map[string]bool{
	"sedan": true, "hatchback": true, "suv": true, "truck": true,
	"coupe": true, "convertible": true, "wagon": true, "van": true,
	"other": true,
}

var allowedPowertrains = map[string]bool{
	"gas": true, "diesel": true, "hybrid": true, "electric": true,
	"other": true,
}

func NewFirecrawlClient(apiKey string) *FirecrawlClient {
	return &FirecrawlClient{
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 3 * time.Minute},
		Endpoint:   "https://api.firecrawl.dev/v2/scrape",
	}
}

func (client *FirecrawlClient) ScrapeToyotaCatalog(ctx context.Context) ([]ToyotaCar, error) {
	if strings.TrimSpace(client.APIKey) == "" {
		return nil, fmt.Errorf("FIRECRAWL_API_KEY is required")
	}

	if client.HTTPClient == nil {
		return nil, fmt.Errorf("Firecrawl HTTP client is not configured")
	}

	currentYear := time.Now().Year()
	requestBody := map[string]any{
		"url": toyotaCatalogURL,
		"formats": []map[string]any{
			{
				"type":   "json",
				"prompt": fmt.Sprintf("Extract every vehicle model and every distinct trim in Toyota's catalog for model years %d and %d shown on this page. Use only facts explicitly present on the page; do not invent trims or infer a model year that is not shown. Include the exact model year, model name, trim, body style, and powertrain. Use one record per trim. Map body_style to one of: sedan, hatchback, suv, truck, coupe, convertible, wagon, van, other. Map powertrain to one of: gas, diesel, hybrid, electric, other. Do not include accessories or non-vehicle products.", currentYear, currentYear+1),
				"schema": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"cars": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"model":      map[string]any{"type": "string"},
									"model_year": map[string]any{"type": "integer"},
									"trim":       map[string]any{"type": "string"},
									"body_style": map[string]any{"type": "string", "enum": []string{"sedan", "hatchback", "suv", "truck", "coupe", "convertible", "wagon", "van", "other"}},
									"powertrain": map[string]any{"type": "string", "enum": []string{"gas", "diesel", "hybrid", "electric", "other"}},
								},
								"required": []string{"model", "model_year", "trim", "body_style", "powertrain"},
							},
						},
					},
					"required": []string{"cars"},
				},
			},
		},
		"onlyMainContent": false,
		"timeout":         120000,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal Firecrawl request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create Firecrawl request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+client.APIKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("call Firecrawl scrape API: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("read Firecrawl response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Firecrawl returned HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Data    struct {
			JSON     json.RawMessage `json:"json"`
			Metadata struct {
				StatusCode int    `json:"statusCode"`
				Error      string `json:"error"`
			} `json:"metadata"`
		} `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("decode Firecrawl response: %w", err)
	}
	if !result.Success {
		return nil, fmt.Errorf("Firecrawl scrape failed: %s", result.Error)
	}
	if result.Data.Metadata.StatusCode < http.StatusOK || result.Data.Metadata.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Toyota catalog returned HTTP %d: %s", result.Data.Metadata.StatusCode, result.Data.Metadata.Error)
	}

	var extraction struct {
		Cars []ToyotaCar `json:"cars"`
	}
	if len(result.Data.JSON) == 0 || string(result.Data.JSON) == "null" {
		return nil, fmt.Errorf("Firecrawl returned no structured Toyota catalog data")
	}
	if err := json.Unmarshal(result.Data.JSON, &extraction); err != nil {
		return nil, fmt.Errorf("decode extracted Toyota catalog: %w", err)
	}
	if len(extraction.Cars) == 0 {
		return nil, fmt.Errorf("Firecrawl extracted no Toyota cars from %s", toyotaCatalogURL)
	}

	return extraction.Cars, nil
}

func ValidateToyotaCatalog(catalog []ToyotaCar) ([]models.Car, error) {
	if len(catalog) == 0 {
		return nil, fmt.Errorf("Toyota catalog is empty")
	}

	cars := make([]models.Car, 0, len(catalog))
	identities := make(map[string]struct{}, len(catalog))
	currentYear := time.Now().Year()

	for i, entry := range catalog {
		entry.Model = strings.TrimSpace(entry.Model)
		entry.Trim = strings.TrimSpace(entry.Trim)
		if entry.Model == "" {
			return nil, fmt.Errorf("Toyota catalog record %d has no model", i+1)
		}
		if entry.ModelYear < currentYear || entry.ModelYear > currentYear+1 {
			return nil, fmt.Errorf("Toyota catalog record %q has model year %d; only %d and %d are allowed", entry.Model, entry.ModelYear, currentYear, currentYear+1)
		}
		if !allowedBodyStyles[entry.BodyStyle] {
			return nil, fmt.Errorf("Toyota catalog record %q has unsupported body style %q", entry.Model, entry.BodyStyle)
		}
		if !allowedPowertrains[entry.Powertrain] {
			return nil, fmt.Errorf("Toyota catalog record %q has unsupported powertrain %q", entry.Model, entry.Powertrain)
		}

		identity := strings.ToLower(entry.Model) + "\x00" +
			strings.ToLower(entry.Trim) + "\x00" +
			fmt.Sprint(entry.ModelYear)
		if _, exists := identities[identity]; exists {
			return nil, fmt.Errorf("Toyota catalog has duplicate model/trim record %q %q (%d)", entry.Model, entry.Trim, entry.ModelYear)
		}
		identities[identity] = struct{}{}

		cars = append(cars, models.Car{
			Make:       "Toyota",
			Model:      entry.Model,
			Year:       entry.ModelYear,
			Trim:       entry.Trim,
			BodyStyle:  entry.BodyStyle,
			Powertrain: entry.Powertrain,
		})
	}

	return cars, nil
}
