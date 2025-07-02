// Package config
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
	FileUpload            FileUploadConfig  `yaml:"file_upload"`
	MinIOPublicURL        string
	Docs                  DocsConfig      `yaml:"docs"`
	AllSecure             AllSecureConfig `yaml:"allsecure"`
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

type AllSecureConfig struct {
	BaseURL                   string  `yaml:"base_url"`
	Username                  string  `yaml:"username"`
	Password                  string  `yaml:"password"`
	WebhookURL                string  `yaml:"webhook_url"`
	WebhookSecret             string  `yaml:"webhook_secret"`
	Timeout                   int     `yaml:"timeout"`
	MarketplaceCommissionRate float64 `yaml:"marketplace_commission_rate"`
	EscrowReleaseDays         int     `yaml:"escrow_release_days"`
	SandboxMode               bool    `yaml:"sandbox_mode"`
}

type OpenSearchConfig struct {
	URL              string `yaml:"url"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	MarketplaceIndex string `yaml:"marketplace_index"`
}

// FileUploadConfig содержит настройки для загрузки файлов
type FileUploadConfig struct {
	MaxImageSize         int64    // Максимальный размер изображения в байтах
	MaxVideoSize         int64    // Максимальный размер видео в байтах
	MaxDocumentSize      int64    // Максимальный размер документа в байтах
	AllowedImageTypes    []string // Разрешенные MIME типы для изображений
	AllowedVideoTypes    []string // Разрешенные MIME типы для видео
	AllowedDocumentTypes []string // Разрешенные MIME типы для документов
}

// DocsConfig содержит настройки для модуля документации
type DocsConfig struct {
	RootPath string // Корневая директория с документацией
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

	// Получаем публичный URL для MinIO (по умолчанию localhost)
	minioPublicURL := os.Getenv("MINIO_PUBLIC_URL")
	if minioPublicURL == "" {
		minioPublicURL = "http://localhost:9000"
	}
	config.MinIOPublicURL = minioPublicURL

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

	// Настройки загрузки файлов
	fileUploadConfig := FileUploadConfig{
		MaxImageSize:    10 * 1024 * 1024,  // 10 MB
		MaxVideoSize:    100 * 1024 * 1024, // 100 MB
		MaxDocumentSize: 20 * 1024 * 1024,  // 20 MB
		AllowedImageTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedVideoTypes: []string{
			"video/mp4",
			"video/mpeg",
			"video/quicktime",
			"video/webm",
		},
		AllowedDocumentTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"text/plain",
			"text/markdown",
			"text/csv",
			"application/octet-stream", // Для неизвестных типов
		},
	}

	// Настройки модуля документации
	docsConfig := DocsConfig{
		RootPath: os.Getenv("DOCS_ROOT_PATH"),
	}

	// Если путь к документации не указан, используем текущую директорию
	if docsConfig.RootPath == "" {
		docsConfig.RootPath = "./docs"
	}

	// Настройки AllSecure
	allSecureConfig := AllSecureConfig{
		BaseURL:                   os.Getenv("ALLSECURE_BASE_URL"),
		Username:                  os.Getenv("ALLSECURE_USERNAME"),
		Password:                  os.Getenv("ALLSECURE_PASSWORD"),
		WebhookURL:                os.Getenv("ALLSECURE_WEBHOOK_URL"),
		WebhookSecret:             os.Getenv("ALLSECURE_WEBHOOK_SECRET"),
		Timeout:                   30,
		MarketplaceCommissionRate: 0.05, // 5% по умолчанию
		EscrowReleaseDays:         7,    // 7 дней по умолчанию
		SandboxMode:               true, // По умолчанию включен sandbox
	}

	// Если базовый URL не указан, используем production endpoint
	if allSecureConfig.BaseURL == "" {
		allSecureConfig.BaseURL = "https://asxgw.com"
	}

	// Настройка комиссии из переменной окружения
	if commissionStr := os.Getenv("ALLSECURE_COMMISSION_RATE"); commissionStr != "" {
		if commission, err := strconv.ParseFloat(commissionStr, 64); err == nil && commission > 0 && commission < 1 {
			allSecureConfig.MarketplaceCommissionRate = commission
		}
	}

	// Настройка дней удержания escrow
	if escrowDaysStr := os.Getenv("ALLSECURE_ESCROW_DAYS"); escrowDaysStr != "" {
		if escrowDays, err := strconv.Atoi(escrowDaysStr); err == nil && escrowDays > 0 {
			allSecureConfig.EscrowReleaseDays = escrowDays
		}
	}

	// Настройка sandbox режима
	if sandboxStr := os.Getenv("ALLSECURE_SANDBOX_MODE"); sandboxStr != "" {
		if sandbox, err := strconv.ParseBool(sandboxStr); err == nil {
			allSecureConfig.SandboxMode = sandbox
		}
	}

	// Настройка timeout
	if timeoutStr := os.Getenv("ALLSECURE_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
			allSecureConfig.Timeout = timeout
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
		MinIOPublicURL:        config.MinIOPublicURL,
		OpenSearch:            config.OpenSearch,
		FileStorage:           config.FileStorage,
		FileUpload:            fileUploadConfig,
		Docs:                  docsConfig,
		AllSecure:             allSecureConfig,
	}, nil
}

// GetJWTDuration возвращает время жизни JWT токена как time.Duration
func (c *Config) GetJWTDuration() time.Duration {
	return time.Duration(c.JWTExpirationHours) * time.Hour
}

// IsProduction проверяет, работает ли приложение в production режиме
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment проверяет, работает ли приложение в development режиме
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetCookieSecure возвращает значение Secure для cookie в зависимости от окружения
func (c *Config) GetCookieSecure() bool {
	// В production всегда true
	// В development false для работы с HTTP
	return c.IsProduction()
}

// GetCookieSameSite возвращает значение SameSite для cookie в зависимости от окружения
func (c *Config) GetCookieSameSite() string {
	// В development используем пустую строку для максимальной совместимости
	// API Routes в Next.js будут проксировать запросы
	// В production используем "Lax" для безопасности
	if c.IsDevelopment() {
		return "" // По умолчанию
	}
	return "Lax"
}
