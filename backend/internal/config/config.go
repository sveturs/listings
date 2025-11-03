// Package config
// backend/internal/config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	envValueTrue = "true"
)

type Config struct {
	Port                     string
	DatabaseURL              string
	GoogleClientID           string
	GoogleClientSecret       string
	GoogleRedirectURL        string
	FrontendURL              string
	AuthServiceURL           string
	AuthServicePublicKeyPath string // Path to auth service public key for JWT validation
	BackendURL               string
	Environment              string
	LogLevel                 string `yaml:"log_level"`
	OpenAIAPIKey             string
	GoogleTranslateAPIKey    string
	ClaudeAPIKey             string
	DeepLAPIKey              string
	DeepLUseFreeAPI          bool `yaml:"deepl_use_free_api"`
	StripeAPIKey             string
	StripeWebhookSecret      string
	OpenSearch               OpenSearchConfig  `yaml:"opensearch"`
	FileStorage              FileStorageConfig `yaml:"file_storage"`
	FileUpload               FileUploadConfig  `yaml:"file_upload"`
	MinIOPublicURL           string
	Docs                     DocsConfig        `yaml:"docs"`
	AllSecure                AllSecureConfig   `yaml:"allsecure"`
	PostExpress              PostExpressConfig `yaml:"postexpress"`
	BEXAuthToken             string            `yaml:"bex_auth_token"`
	BEXClientID              string            `yaml:"bex_client_id"`
	BEXAPIURL                string            `yaml:"bex_api_url"`
	SearchWeights            *SearchWeights    `yaml:"search_weights"`
	Redis                    RedisConfig       `yaml:"redis"`
	MigrationsOnAPI          string            `yaml:"migrations_on_api"` // off, schema, full
	ReindexOnAPI             string            `yaml:"reindex_on_api"`    // on
	FeatureFlags             *FeatureFlags     `yaml:"feature_flags"`
	Currency                 CurrencyConfig    `yaml:"currency"`
	DeliveryGRPCURL          string            `yaml:"delivery_grpc_url"`         // URL для delivery микросервиса через gRPC
	ListingsGRPCURL          string            `yaml:"listings_grpc_url"`         // URL для listings микросервиса через gRPC
	UseListingsMicroservice  bool              `yaml:"use_listings_microservice"` // Feature flag для включения listings микросервиса
	Marketplace              MarketplaceConfig `yaml:"marketplace"`               // Marketplace microservice migration config
}

type FileStorageConfig struct {
	Provider                string `yaml:"provider"` // "local" или "minio"
	LocalBasePath           string `yaml:"local_base_path"`
	PublicBaseURL           string `yaml:"public_base_url"`
	MinioEndpoint           string `yaml:"minio_endpoint"`
	MinioAccessKey          string `yaml:"minio_access_key"`
	MinioSecretKey          string `yaml:"minio_secret_key"`
	MinioUseSSL             bool   `yaml:"minio_use_ssl"`
	MinioBucketName         string `yaml:"minio_bucket_name"`          // Основной bucket для объявлений
	MinioChatBucket         string `yaml:"minio_chat_bucket"`          // Bucket для файлов чата
	MinioStorefrontBucket   string `yaml:"minio_storefront_bucket"`    // Bucket для товаров витрин
	MinioReviewPhotosBucket string `yaml:"minio_review_photos_bucket"` // Bucket для фотографий отзывов
	MinioLocation           string `yaml:"minio_location"`
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

type PostExpressConfig struct {
	BaseURL       string `yaml:"base_url"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	TestMode      bool   `yaml:"test_mode"`
	SenderName    string `yaml:"sender_name"`
	SenderAddress string `yaml:"sender_address"`
}

type RedisConfig struct {
	URL               string        `yaml:"url"`
	Password          string        `yaml:"password"`
	DB                int           `yaml:"db"`
	MaxRetries        int           `yaml:"max_retries"`
	PoolSize          int           `yaml:"pool_size"`
	CacheTTL          time.Duration `yaml:"cache_ttl"`
	CategoriesTTL     time.Duration `yaml:"categories_ttl"`
	AttributesTTL     time.Duration `yaml:"attributes_ttl"`
	AttributeGroupTTL time.Duration `yaml:"attribute_group_ttl"`
}

// MarketplaceConfig содержит настройки для миграции marketplace на microservice
type MarketplaceConfig struct {
	// Feature flag для включения microservice
	// Default: true (100% microservice for dev environment)
	UseMicroservice bool `yaml:"use_microservice" envconfig:"USE_MARKETPLACE_MICROSERVICE" default:"true"`

	// Процент трафика, который идёт на microservice (0-100)
	// Default: 100 (100% microservice for dev environment)
	RolloutPercent int `yaml:"rollout_percent" envconfig:"MARKETPLACE_ROLLOUT_PERCENT" default:"100"`

	// URL microservice через gRPC
	MicroserviceGRPCURL string `yaml:"microservice_grpc_url" envconfig:"MARKETPLACE_MICROSERVICE_GRPC_URL" default:"localhost:50053"`

	// Admin override - админы всегда используют microservice (для тестирования)
	AdminOverride bool `yaml:"admin_override" envconfig:"MARKETPLACE_ADMIN_OVERRIDE" default:"true"`

	// Список canary user IDs (через запятую) - эти пользователи всегда на microservice
	CanaryUserIDs string `yaml:"canary_user_ids" envconfig:"MARKETPLACE_CANARY_USERIDS" default:""`

	// Timeout для gRPC запросов
	GRPCTimeout time.Duration `yaml:"grpc_timeout" envconfig:"MARKETPLACE_GRPC_TIMEOUT" default:"5s"`

	// Timeout для microservice запросов (per-request context deadline)
	MicroserviceTimeout time.Duration `yaml:"microservice_timeout" envconfig:"MICROSERVICE_TIMEOUT" default:"500ms"`

	// Fallback на monolith при ошибках microservice
	FallbackToMonolith bool `yaml:"fallback_to_monolith" envconfig:"MARKETPLACE_FALLBACK_TO_MONOLITH" default:"true"`

	// Circuit breaker configuration
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker"`
}

// CircuitBreakerConfig defines circuit breaker configuration for marketplace microservice
type CircuitBreakerConfig struct {
	// Enabled включает/выключает circuit breaker
	Enabled bool `yaml:"enabled" envconfig:"MARKETPLACE_CB_ENABLED" default:"true"`

	// FailureThreshold количество consecutive failures для открытия circuit
	FailureThreshold int `yaml:"failure_threshold" envconfig:"MARKETPLACE_CB_FAILURE_THRESHOLD" default:"5"`

	// SuccessThreshold количество successful requests в HALF_OPEN для закрытия circuit
	SuccessThreshold int `yaml:"success_threshold" envconfig:"MARKETPLACE_CB_SUCCESS_THRESHOLD" default:"2"`

	// Timeout время ожидания перед переходом из OPEN в HALF_OPEN
	Timeout time.Duration `yaml:"timeout" envconfig:"MARKETPLACE_CB_TIMEOUT" default:"60s"`

	// HalfOpenMaxRequests максимальное количество requests в HALF_OPEN state
	HalfOpenMaxRequests int `yaml:"half_open_max_requests" envconfig:"MARKETPLACE_CB_HALF_OPEN_MAX" default:"3"`

	// CounterResetInterval интервал сброса счётчиков (для sliding window)
	CounterResetInterval time.Duration `yaml:"counter_reset_interval" envconfig:"MARKETPLACE_CB_COUNTER_RESET_INTERVAL" default:"60s"`
}

type OpenSearchConfig struct {
	URL            string `yaml:"url"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	UnifiedIndex   string `yaml:"unified_index"`    // Unified index (Sprint 2.2) - c2c + b2c unified
	UnifiedIndexV2 string `yaml:"unified_index_v2"` // Unified index v2 (optimized mapping, 30 fields)
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

	// Получаем уровень логирования
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // По умолчанию info
	}
	config.StripeAPIKey = os.Getenv("STRIPE_API_KEY")
	config.StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

	// Получаем URL auth-service
	config.AuthServiceURL = os.Getenv("AUTH_SERVICE_URL")
	if config.AuthServiceURL == "" {
		config.AuthServiceURL = "http://localhost:28080"
		log.Printf("Using default AUTH_SERVICE_URL: %s", config.AuthServiceURL)
	}

	// Получаем путь к публичному ключу auth service
	config.AuthServicePublicKeyPath = os.Getenv("AUTH_SERVICE_PUBLIC_KEY_PATH")
	if config.AuthServicePublicKeyPath == "" {
		config.AuthServicePublicKeyPath = "./keys/auth_service_public.pem"
		log.Printf("Using default AUTH_SERVICE_PUBLIC_KEY_PATH: %s", config.AuthServicePublicKeyPath)
	}

	// Получаем URL текущего backend для OAuth callbacks
	config.BackendURL = os.Getenv("BACKEND_URL")
	if config.BackendURL == "" {
		config.BackendURL = fmt.Sprintf("http://localhost:%s", port)
		log.Printf("Using default BACKEND_URL: %s", config.BackendURL)
	}

	// Получаем ключ Google Translate API (необязательный)
	config.GoogleTranslateAPIKey = os.Getenv("GOOGLE_TRANSLATE_API_KEY")

	// Получаем ключ Claude API (необязательный)
	config.ClaudeAPIKey = os.Getenv("CLAUDE_API_KEY")

	// Получаем ключ DeepL API (необязательный)
	config.DeepLAPIKey = os.Getenv("DEEPL_API_KEY")
	config.DeepLUseFreeAPI = os.Getenv("DEEPL_USE_FREE_API") == envValueTrue

	// BEX Express configuration
	config.BEXAuthToken = os.Getenv("BEX_AUTH_TOKEN")
	config.BEXClientID = os.Getenv("BEX_CLIENT_ID")
	bexAPIURL := os.Getenv("BEX_API_URL")
	if bexAPIURL == "" {
		bexAPIURL = "https://api.bex.rs:62502"
	}
	config.BEXAPIURL = bexAPIURL

	// Delivery gRPC service configuration (required in .env)
	config.DeliveryGRPCURL = os.Getenv("DELIVERY_GRPC_URL")

	// Listings gRPC service configuration (optional - feature flag controlled)
	config.ListingsGRPCURL = os.Getenv("LISTINGS_GRPC_URL")
	if config.ListingsGRPCURL == "" {
		config.ListingsGRPCURL = "localhost:50051" // Default gRPC port
	}

	// Feature flag для использования listings микросервиса (по умолчанию false - используем local DB)
	config.UseListingsMicroservice = os.Getenv("USE_LISTINGS_MICROSERVICE") == envValueTrue

	// Получаем публичный URL для MinIO (по умолчанию localhost)
	minioPublicURL := os.Getenv("MINIO_PUBLIC_URL")
	if minioPublicURL == "" {
		minioPublicURL = "http://localhost:9000"
	}
	config.MinIOPublicURL = minioPublicURL

	config.OpenSearch = OpenSearchConfig{
		URL:            os.Getenv("OPENSEARCH_URL"),
		Username:       os.Getenv("OPENSEARCH_USERNAME"),
		Password:       os.Getenv("OPENSEARCH_PASSWORD"),
		UnifiedIndex:   os.Getenv("OPENSEARCH_UNIFIED_INDEX"),
		UnifiedIndexV2: os.Getenv("OPENSEARCH_UNIFIED_INDEX_V2"),
	}

	// Если индекс Unified не указан, используем значение по умолчанию
	if config.OpenSearch.UnifiedIndex == "" {
		config.OpenSearch.UnifiedIndex = "marketplace_listings"
	}

	// Если индекс Unified V2 не указан, используем значение по умолчанию
	if config.OpenSearch.UnifiedIndexV2 == "" {
		config.OpenSearch.UnifiedIndexV2 = "unified_listings_v2"
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
		MinioUseSSL:     os.Getenv("MINIO_USE_SSL") == envValueTrue,
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
		// Загружаем названия дополнительных bucket'ов
		config.FileStorage.MinioChatBucket = os.Getenv("MINIO_CHAT_BUCKET")
		if config.FileStorage.MinioChatBucket == "" {
			config.FileStorage.MinioChatBucket = "chat-files" // По умолчанию
		}
		config.FileStorage.MinioStorefrontBucket = os.Getenv("MINIO_STOREFRONT_BUCKET")
		if config.FileStorage.MinioStorefrontBucket == "" {
			config.FileStorage.MinioStorefrontBucket = "storefront-products" // По умолчанию
		}
		config.FileStorage.MinioReviewPhotosBucket = os.Getenv("MINIO_REVIEW_PHOTOS_BUCKET")
		if config.FileStorage.MinioReviewPhotosBucket == "" {
			config.FileStorage.MinioReviewPhotosBucket = "review-photos" // По умолчанию
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

	// Настройки Post Express
	postExpressConfig := PostExpressConfig{
		BaseURL:       os.Getenv("POSTEXPRESS_BASE_URL"),
		Username:      os.Getenv("POSTEXPRESS_USERNAME"),
		Password:      os.Getenv("POSTEXPRESS_PASSWORD"),
		TestMode:      os.Getenv("POSTEXPRESS_TEST_MODE") == envValueTrue,
		SenderName:    os.Getenv("POSTEXPRESS_SENDER_NAME"),
		SenderAddress: os.Getenv("POSTEXPRESS_SENDER_ADDRESS"),
	}

	// Если базовый URL не указан, используем production endpoint
	if postExpressConfig.BaseURL == "" {
		postExpressConfig.BaseURL = "https://wsp.postexpress.rs/api/Transakcija"
	}

	// Значения по умолчанию для отправителя
	if postExpressConfig.SenderName == "" {
		postExpressConfig.SenderName = "Sve Tu d.o.o."
	}
	if postExpressConfig.SenderAddress == "" {
		postExpressConfig.SenderAddress = "Микија Манојловића 53, 21000 Нови Сад"
	}

	// Загружаем веса поиска (используем дефолтные значения)
	searchWeights := GetDefaultSearchWeights()

	// Настройки Redis
	redisConfig := RedisConfig{
		URL:               os.Getenv("REDIS_URL"),
		Password:          os.Getenv("REDIS_PASSWORD"),
		DB:                0,
		MaxRetries:        3,
		PoolSize:          10,
		CacheTTL:          24 * time.Hour, // По умолчанию 24 часа
		CategoriesTTL:     6 * time.Hour,  // Категории кешируем на 6 часов
		AttributesTTL:     4 * time.Hour,  // Атрибуты на 4 часа
		AttributeGroupTTL: 4 * time.Hour,  // Группы атрибутов на 4 часа
	}

	// Если Redis URL не указан, используем локальный по умолчанию
	if redisConfig.URL == "" {
		redisConfig.URL = "localhost:6379"
	}

	// Парсим настройки TTL из переменных окружения
	if ttl := os.Getenv("REDIS_CATEGORIES_TTL"); ttl != "" {
		if duration, err := time.ParseDuration(ttl); err == nil {
			redisConfig.CategoriesTTL = duration
		}
	}
	if ttl := os.Getenv("REDIS_ATTRIBUTES_TTL"); ttl != "" {
		if duration, err := time.ParseDuration(ttl); err == nil {
			redisConfig.AttributesTTL = duration
		}
	}
	if ttl := os.Getenv("REDIS_ATTRIBUTE_GROUP_TTL"); ttl != "" {
		if duration, err := time.ParseDuration(ttl); err == nil {
			redisConfig.AttributeGroupTTL = duration
		}
	}

	// Парсим настройки пула соединений
	if poolSize := os.Getenv("REDIS_POOL_SIZE"); poolSize != "" {
		if size, err := strconv.Atoi(poolSize); err == nil && size > 0 {
			redisConfig.PoolSize = size
		}
	}

	// Парсим настройки базы данных
	if db := os.Getenv("REDIS_DB"); db != "" {
		if dbNum, err := strconv.Atoi(db); err == nil && dbNum >= 0 {
			redisConfig.DB = dbNum
		}
	}

	// Получаем настройку миграций при старте API
	migrationsOnAPI := os.Getenv("MIGRATIONS_ON_API")
	if migrationsOnAPI == "" {
		migrationsOnAPI = "off" // По умолчанию выключено
	}

	// Получаем настройку reindex при старте API
	reindexOnAPI := strings.TrimSpace(os.Getenv("REINDEX_ON_API"))

	// Загружаем конфигурацию валют
	currencyConfig := GetDefaultCurrencyConfig()

	// Переопределяем из переменных окружения, если указаны
	if defaultCurrency := os.Getenv("DEFAULT_CURRENCY"); defaultCurrency != "" {
		currencyConfig.DefaultCurrency = defaultCurrency
	}
	if supportedCurrencies := os.Getenv("SUPPORTED_CURRENCIES"); supportedCurrencies != "" {
		currencyConfig.SupportedCurrencies = strings.Split(supportedCurrencies, ",")
	}
	if exchangeRatesEnabled := os.Getenv("EXCHANGE_RATES_ENABLED"); exchangeRatesEnabled == envValueTrue {
		currencyConfig.ExchangeRatesEnabled = true
	}
	if exchangeRatesProvider := os.Getenv("EXCHANGE_RATES_PROVIDER"); exchangeRatesProvider != "" {
		currencyConfig.ExchangeRatesProvider = exchangeRatesProvider
	}

	// Загружаем конфигурацию marketplace microservice migration
	marketplaceConfig := MarketplaceConfig{
		UseMicroservice:     os.Getenv("USE_MARKETPLACE_MICROSERVICE") == envValueTrue,
		RolloutPercent:      0,
		MicroserviceGRPCURL: "localhost:50053",
		AdminOverride:       true,
		CanaryUserIDs:       os.Getenv("MARKETPLACE_CANARY_USERIDS"),
		GRPCTimeout:         5 * time.Second,
		MicroserviceTimeout: 500 * time.Millisecond,
		FallbackToMonolith:  true,
	}

	// Парсим rollout percent
	if rolloutStr := os.Getenv("MARKETPLACE_ROLLOUT_PERCENT"); rolloutStr != "" {
		if rollout, err := strconv.Atoi(rolloutStr); err == nil && rollout >= 0 && rollout <= 100 {
			marketplaceConfig.RolloutPercent = rollout
		}
	}

	// Парсим microservice URL
	if grpcURL := os.Getenv("MARKETPLACE_MICROSERVICE_GRPC_URL"); grpcURL != "" {
		marketplaceConfig.MicroserviceGRPCURL = grpcURL
	}

	// Парсим admin override
	if adminOverride := os.Getenv("MARKETPLACE_ADMIN_OVERRIDE"); adminOverride != "" {
		marketplaceConfig.AdminOverride = adminOverride == envValueTrue
	}

	// Парсим gRPC timeout
	if timeoutStr := os.Getenv("MARKETPLACE_GRPC_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil && timeout > 0 {
			marketplaceConfig.GRPCTimeout = timeout
		}
	}

	// Парсим microservice timeout
	if timeoutStr := os.Getenv("MICROSERVICE_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil && timeout > 0 {
			marketplaceConfig.MicroserviceTimeout = timeout
		}
	}

	// Парсим fallback flag
	if fallbackStr := os.Getenv("MARKETPLACE_FALLBACK_TO_MONOLITH"); fallbackStr != "" {
		marketplaceConfig.FallbackToMonolith = fallbackStr == envValueTrue
	}

	return &Config{
		Port:                     port,
		DatabaseURL:              dbURL,
		GoogleClientID:           googleClientID,
		GoogleClientSecret:       googleClientSecret,
		GoogleRedirectURL:        googleRedirectURL,
		FrontendURL:              frontendURL,
		Environment:              environment,
		LogLevel:                 logLevel,
		OpenAIAPIKey:             openAIAPIKey,
		GoogleTranslateAPIKey:    config.GoogleTranslateAPIKey,
		ClaudeAPIKey:             config.ClaudeAPIKey,
		DeepLAPIKey:              config.DeepLAPIKey,
		DeepLUseFreeAPI:          config.DeepLUseFreeAPI,
		StripeAPIKey:             config.StripeAPIKey,
		StripeWebhookSecret:      config.StripeWebhookSecret,
		AuthServiceURL:           config.AuthServiceURL,
		AuthServicePublicKeyPath: config.AuthServicePublicKeyPath,
		BackendURL:               config.BackendURL,
		MinIOPublicURL:           config.MinIOPublicURL,
		DeliveryGRPCURL:          config.DeliveryGRPCURL,
		ListingsGRPCURL:          config.ListingsGRPCURL,
		UseListingsMicroservice:  config.UseListingsMicroservice,
		OpenSearch:               config.OpenSearch,
		FileStorage:              config.FileStorage,
		FileUpload:               fileUploadConfig,
		Docs:                     docsConfig,
		AllSecure:                allSecureConfig,
		PostExpress:              postExpressConfig,
		BEXAuthToken:             config.BEXAuthToken,
		BEXClientID:              config.BEXClientID,
		BEXAPIURL:                config.BEXAPIURL,
		SearchWeights:            searchWeights,
		Redis:                    redisConfig,
		MigrationsOnAPI:          migrationsOnAPI,
		ReindexOnAPI:             reindexOnAPI,
		FeatureFlags:             LoadFeatureFlags(),
		Currency:                 currencyConfig,
		Marketplace:              marketplaceConfig,
	}, nil
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
	// В development используем "None" для работы с разными доменами/IP
	// Это позволит cookies работать при доступе через Tailscale
	// В production используем "Lax" для безопасности
	if c.IsDevelopment() {
		// Если не используем HTTPS, то не можем использовать SameSite=None
		// В таком случае используем пустую строку
		if !c.GetCookieSecure() {
			return "" // Без SameSite для HTTP
		}
		return "None" // Для HTTPS разрешаем кросс-доменные запросы
	}
	return "Lax"
}

// GetCookieDomain возвращает домен для cookie в зависимости от окружения
func (c *Config) GetCookieDomain() string {
	// В development не устанавливаем домен, чтобы cookie работали
	// как для localhost, так и для IP адресов (например, через Tailscale)
	if c.IsDevelopment() {
		return "" // Пустой домен позволяет cookie работать на любом хосте
	}
	return "" // В production также используем текущий домен
}
