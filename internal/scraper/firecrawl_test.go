package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestValidateToyotaCatalog(t *testing.T) {
	cars, err := ValidateToyotaCatalog([]ToyotaCar{
		{
			Model:      " RAV4 ",
			ModelYear:  2026,
			Trim:       " XLE ",
			BodyStyle:  "suv",
			Powertrain: "hybrid",
		},
	})
	if err != nil {
		t.Fatalf("ValidateToyotaCatalog returned error: %v", err)
	}
	if len(cars) != 1 {
		t.Fatalf("got %d cars, want 1", len(cars))
	}
	if cars[0].Make != "Toyota" || cars[0].Model != "RAV4" || cars[0].Trim != "XLE" {
		t.Fatalf("unexpected mapped car: %#v", cars[0])
	}
}

func TestValidateToyotaCatalogAllowsCurrentAndNextModelYears(t *testing.T) {
	currentYear := time.Now().Year()
	cars, err := ValidateToyotaCatalog([]ToyotaCar{
		{Model: "Camry", ModelYear: currentYear, Trim: "LE", BodyStyle: "sedan", Powertrain: "gas"},
		{Model: "RAV4", ModelYear: currentYear + 1, Trim: "LE", BodyStyle: "suv", Powertrain: "gas"},
	})
	if err != nil {
		t.Fatalf("ValidateToyotaCatalog returned error: %v", err)
	}
	if len(cars) != 2 {
		t.Fatalf("got %d cars, want 2", len(cars))
	}
}

func TestValidateToyotaCatalogRejectsInvalidRecords(t *testing.T) {
	currentYear := time.Now().Year()
	tests := []struct {
		name    string
		catalog []ToyotaCar
	}{
		{
			name: "unsupported body style",
			catalog: []ToyotaCar{
				{Model: "Camry", ModelYear: currentYear, Trim: "LE", BodyStyle: "car", Powertrain: "gas"},
			},
		},
		{
			name: "duplicate model trim",
			catalog: []ToyotaCar{
				{Model: "Camry", ModelYear: currentYear, Trim: "LE", BodyStyle: "sedan", Powertrain: "gas"},
				{Model: " camry ", ModelYear: currentYear, Trim: " le ", BodyStyle: "sedan", Powertrain: "gas"},
			},
		},
		{
			name: "unsupported model year",
			catalog: []ToyotaCar{
				{Model: "Camry", ModelYear: currentYear + 2, Trim: "LE", BodyStyle: "sedan", Powertrain: "gas"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ValidateToyotaCatalog(test.catalog); err == nil {
				t.Fatal("ValidateToyotaCatalog succeeded, want an error")
			}
		})
	}
}

func TestScrapeToyotaCatalog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("request method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want bearer test-key", got)
		}
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"json": {
					"cars": [{
						"model": "Camry",
						"model_year": 2026,
						"trim": "LE",
						"body_style": "sedan",
						"powertrain": "hybrid"
					}]
				},
				"metadata": {"statusCode": 200}
			}
		}`))
	}))
	defer server.Close()

	client := NewFirecrawlClient("test-key")
	client.HTTPClient = server.Client()
	client.Endpoint = server.URL

	cars, err := client.ScrapeToyotaCatalog(context.Background())
	if err != nil {
		t.Fatalf("ScrapeToyotaCatalog returned error: %v", err)
	}
	if len(cars) != 1 || cars[0].Model != "Camry" {
		t.Fatalf("unexpected scraped cars: %#v", cars)
	}
}

func TestScrapeToyotaCatalogRequiresAPIKey(t *testing.T) {
	client := NewFirecrawlClient(" ")
	_, err := client.ScrapeToyotaCatalog(context.Background())
	if err == nil || !strings.Contains(err.Error(), "FIRECRAWL_API_KEY") {
		t.Fatalf("error = %v, want missing API key error", err)
	}
}
