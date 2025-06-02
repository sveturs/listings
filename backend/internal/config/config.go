// backend/internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleRedirectURL     string
	FrontendURL           string
	Environment           string
	OpenAIAPIKey          string
	GoogleTranslateAPIKey string
	StripeAPIKey          string
	StripeWebhookSecret   string
	JWTSecret             string
	JWTExpirationHours    int               `yaml:"jwt_expiration_hours"`
	OpenSearch            OpenSearchConfig  `yaml:"opensearch"`
	FileStorage           FileStorageConfig `yaml:"file_storage"`
}

type FileStorageConfig struct {
	Provider        string `yaml:"provider"` // "local" или "minio"
	LocalBasePath   string `yaml:"local_base_path"`
	PublicBaseURL   string `yaml:"public_base_url"`
	MinioEndpoint   string `yaml:"minio_endpoint"`
	MinioAccessKey  string `yaml:"minio_access_key"`
	MinioSecretKey  string `yaml:"minio_secret_key"`
	MinioUseSSL     bool   `yaml:"minio_use_ssl"`
	MinioBucketName string `yaml:"minio_bucket_name"`
	MinioLocation   string `yaml:"minio_location"`
}

type OpenSearchConfig struct {
	URL              string `yaml:"url"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
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

	// Получаем JWT секретный ключ
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-jwt-secret-change-in-production" // Дефолтное значение для разработки
	}
	config.JWTSecret = jwtSecret

	// Получаем время жизни JWT токена (в часах)
	jwtExpirationStr := os.Getenv("JWT_EXPIRATION_HOURS")
	jwtExpirationHours := 1 // По умолчанию 1 час
	if jwtExpirationStr != "" {
		if parsed, err := strconv.Atoi(jwtExpirationStr); err == nil && parsed > 0 {
			jwtExpirationHours = parsed
		}
	}
	config.JWTExpirationHours = jwtExpirationHours

	// Получаем ключ Google Translate API (необязательный)
	config.GoogleTranslateAPIKey = os.Getenv("GOOGLE_TRANSLATE_API_KEY")

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

	// Настройки хранилища файлов
	provider := os.Getenv("FILE_STORAGE_PROVIDER")
	if provider == "" {
		provider = "minio" // По умолчанию используем MinIO
	}
	
	config.FileStorage = FileStorageConfig{
		Provider:        provider,
		LocalBasePath:   os.Getenv("FILE_STORAGE_LOCAL_PATH"),
		PublicBaseURL:   os.Getenv("FILE_STORAGE_PUBLIC_URL"),
		MinioEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		MinioUseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
		MinioBucketName: os.Getenv("MINIO_BUCKET_NAME"),
		MinioLocation:   os.Getenv("MINIO_LOCATION"),
	}
	
	// Валидация на основе выбранного провайдера
	if config.FileStorage.Provider == "minio" {
		if config.FileStorage.MinioEndpoint == "" {
			return nil, fmt.Errorf("MINIO_ENDPOINT is not set")
		}
		if config.FileStorage.MinioAccessKey == "" {
			return nil, fmt.Errorf("MINIO_ACCESS_KEY is not set")
		}
		if config.FileStorage.MinioSecretKey == "" {
			return nil, fmt.Errorf("MINIO_SECRET_KEY is not set")
		}
		if config.FileStorage.MinioBucketName == "" {
			config.FileStorage.MinioBucketName = "listings" // По умолчанию
		}
	}

	return &Config{
		Port:                  port,
		DatabaseURL:           dbURL,
		GoogleClientID:        googleClientID,
		GoogleClientSecret:    googleClientSecret,
		GoogleRedirectURL:     googleRedirectURL,
		FrontendURL:           frontendURL,
		Environment:           environment,
		OpenAIAPIKey:          openAIAPIKey,
		GoogleTranslateAPIKey: config.GoogleTranslateAPIKey,
		StripeAPIKey:          config.StripeAPIKey,
		StripeWebhookSecret:   config.StripeWebhookSecret,
		JWTSecret:             config.JWTSecret,
		JWTExpirationHours:    config.JWTExpirationHours,
		OpenSearch:            config.OpenSearch,
		FileStorage:           config.FileStorage,
	}, nil
}

// GetJWTDuration возвращает время жизни JWT токена как time.Duration
func (c *Config) GetJWTDuration() time.Duration {
	return time.Duration(c.JWTExpirationHours) * time.Hour
}