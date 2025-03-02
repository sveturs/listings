// backend/internal/config/config.go
package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                string
	DatabaseURL         string
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	FrontendURL         string
	Environment         string
	OpenAIAPIKey        string
	StripeAPIKey        string
	StripeWebhookSecret string
	OpenSearch          OpenSearchConfig `yaml:"opensearch"`
}
type OpenSearchConfig struct {
    URL             string `yaml:"url"`
    Username        string `yaml:"username"`
    Password        string `yaml:"password"`
    MarketplaceIndex string `yaml:"marketplace_index"`
}
func NewConfig() (*Config, error) {

	config := &Config{}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	if openAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is not set")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is not set")
	}

	googleRedirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")
	if googleRedirectURL == "" {
		return nil, fmt.Errorf("GOOGLE_OAUTH_REDIRECT_URL is not set")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		return nil, fmt.Errorf("FRONTEND_URL is not set")
	}

	environment := os.Getenv("APP_MODE")
	if environment == "" {
		environment = "development"
	}
	config.StripeAPIKey = os.Getenv("STRIPE_API_KEY")
	config.StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

    config.OpenSearch = OpenSearchConfig{
        URL:              os.Getenv("OPENSEARCH_URL"),
        Username:         os.Getenv("OPENSEARCH_USERNAME"),
        Password:         os.Getenv("OPENSEARCH_PASSWORD"),
        MarketplaceIndex: os.Getenv("OPENSEARCH_MARKETPLACE_INDEX"),
    }

   // Если индекс не указан, используем значение по умолчанию
   if config.OpenSearch.MarketplaceIndex == "" {
	config.OpenSearch.MarketplaceIndex = "marketplace"
}

return &Config{
	Port:                port,
	DatabaseURL:         dbURL,
	GoogleClientID:      googleClientID,
	GoogleClientSecret:  googleClientSecret,
	GoogleRedirectURL:   googleRedirectURL,
	FrontendURL:         frontendURL,
	Environment:         environment,
	OpenAIAPIKey:        openAIAPIKey,
	StripeAPIKey:        config.StripeAPIKey,
	StripeWebhookSecret: config.StripeWebhookSecret,
	OpenSearch:          config.OpenSearch, 
}, nil
}