package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Server   ServerConfig
	DB       DBConfig
	Redis    RedisConfig
	Search   SearchConfig
	Storage  StorageConfig
	Auth     AuthConfig
	Delivery DeliveryConfig
	Worker   WorkerConfig
	Features FeatureFlags
	Tracing  TracingConfig
	CORS     CORSConfig
	Health   HealthConfig
	WMS      WMSConfig
}

// AppConfig contains general application settings
type AppConfig struct {
	Env       string `envconfig:"VONDILISTINGS_ENV" default:"development"`
	LogLevel  string `envconfig:"VONDILISTINGS_LOG_LEVEL" default:"info"`
	LogFormat string `envconfig:"VONDILISTINGS_LOG_FORMAT" default:"json"`
}

// ServerConfig contains server ports and settings
type ServerConfig struct {
	GRPCHost    string `envconfig:"VONDILISTINGS_GRPC_HOST" default:"0.0.0.0"`
	GRPCPort    int    `envconfig:"VONDILISTINGS_GRPC_PORT" default:"50053"`
	HTTPHost    string `envconfig:"VONDILISTINGS_HTTP_HOST" default:"0.0.0.0"`
	HTTPPort    int    `envconfig:"VONDILISTINGS_HTTP_PORT" default:"8086"`
	MetricsHost string `envconfig:"VONDILISTINGS_METRICS_HOST" default:"0.0.0.0"`
	MetricsPort int    `envconfig:"VONDILISTINGS_METRICS_PORT" default:"9093"`
}

// DBConfig contains PostgreSQL database configuration
type DBConfig struct {
	Host     string `envconfig:"VONDILISTINGS_DB_HOST" default:"localhost"`
	Port     int    `envconfig:"VONDILISTINGS_DB_PORT" default:"35434"`
	User     string `envconfig:"VONDILISTINGS_DB_USER" default:"listings_user"`
	Password string `envconfig:"VONDILISTINGS_DB_PASSWORD" default:"listings_password"`
	Name     string `envconfig:"VONDILISTINGS_DB_NAME" default:"listings_db"`
	SSLMode  string `envconfig:"VONDILISTINGS_DB_SSLMODE" default:"disable"`

	// Connection pool settings
	MaxOpenConns    int           `envconfig:"VONDILISTINGS_DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"VONDILISTINGS_DB_MAX_IDLE_CONNS" default:"10"`
	ConnMaxLifetime time.Duration `envconfig:"VONDILISTINGS_DB_CONN_MAX_LIFETIME" default:"5m"`
	ConnMaxIdleTime time.Duration `envconfig:"VONDILISTINGS_DB_CONN_MAX_IDLE_TIME" default:"10m"`
}

// DSN returns PostgreSQL connection string
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `envconfig:"VONDILISTINGS_REDIS_HOST" default:"localhost"`
	Port     int    `envconfig:"VONDILISTINGS_REDIS_PORT" default:"36380"`
	Password string `envconfig:"VONDILISTINGS_REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"VONDILISTINGS_REDIS_DB" default:"0"`

	// Connection pool
	PoolSize     int `envconfig:"VONDILISTINGS_REDIS_POOL_SIZE" default:"10"`
	MinIdleConns int `envconfig:"VONDILISTINGS_REDIS_MIN_IDLE_CONNS" default:"5"`

	// Cache TTL
	ListingTTL time.Duration `envconfig:"VONDILISTINGS_CACHE_LISTING_TTL" default:"5m"`
	SearchTTL  time.Duration `envconfig:"VONDILISTINGS_CACHE_SEARCH_TTL" default:"2m"`
}

// Addr returns Redis address
func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// SearchConfig contains OpenSearch configuration
type SearchConfig struct {
	Addresses []string `envconfig:"VONDILISTINGS_OPENSEARCH_ADDRESSES" default:"http://localhost:9200"`
	Username  string   `envconfig:"VONDILISTINGS_OPENSEARCH_USERNAME" default:"admin"`
	Password  string   `envconfig:"VONDILISTINGS_OPENSEARCH_PASSWORD" default:"admin"`
	Index     string   `envconfig:"VONDILISTINGS_OPENSEARCH_INDEX" default:"marketplace_listings"`
}

// StorageConfig contains MinIO (S3-compatible) configuration
type StorageConfig struct {
	Endpoint      string `envconfig:"VONDILISTINGS_MINIO_ENDPOINT" default:"localhost:9000"`
	AccessKey     string `envconfig:"VONDILISTINGS_MINIO_ACCESS_KEY" default:"minioadmin"`
	SecretKey     string `envconfig:"VONDILISTINGS_MINIO_SECRET_KEY" default:"minioadmin"`
	UseSSL        bool   `envconfig:"VONDILISTINGS_MINIO_USE_SSL" default:"false"`
	Bucket        string `envconfig:"VONDILISTINGS_MINIO_BUCKET" default:"listings-images"`
	PublicBaseURL string `envconfig:"VONDILISTINGS_MINIO_PUBLIC_BASE_URL" default:""`
}

// GetPublicBaseURL returns the public base URL for objects.
// If PublicBaseURL is configured, uses it. Otherwise, constructs from Endpoint.
func (c StorageConfig) GetPublicBaseURL() string {
	if c.PublicBaseURL != "" {
		return c.PublicBaseURL
	}
	// Construct from endpoint
	protocol := "http"
	if c.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s", protocol, c.Endpoint, c.Bucket)
}

// AuthConfig contains Auth Service integration settings
type AuthConfig struct {
	ServiceURL    string        `envconfig:"VONDILISTINGS_AUTH_SERVICE_URL" default:"http://localhost:8081"`
	PublicKeyPath string        `envconfig:"VONDILISTINGS_AUTH_PUBLIC_KEY_PATH" default:"/keys/public.pem"`
	Timeout       time.Duration `envconfig:"VONDILISTINGS_AUTH_TIMEOUT" default:"10s"`
	Enabled       bool          `envconfig:"VONDILISTINGS_AUTH_ENABLED" default:"false"` // Disabled for Phase 13.1.15.8 until logger adapter is fixed
}

// DeliveryConfig contains Delivery microservice integration settings
type DeliveryConfig struct {
	GRPCAddress string        `envconfig:"VONDILISTINGS_DELIVERY_GRPC_ADDRESS" default:"localhost:50052"`
	Timeout     time.Duration `envconfig:"VONDILISTINGS_DELIVERY_TIMEOUT" default:"10s"`
	MaxRetries  int           `envconfig:"VONDILISTINGS_DELIVERY_MAX_RETRIES" default:"3"`
	RetryDelay  time.Duration `envconfig:"VONDILISTINGS_DELIVERY_RETRY_DELAY" default:"100ms"`
	Enabled     bool          `envconfig:"VONDILISTINGS_DELIVERY_ENABLED" default:"false"`
}

// WorkerConfig contains async worker settings
type WorkerConfig struct {
	Enabled     bool   `envconfig:"VONDILISTINGS_WORKER_ENABLED" default:"true"`
	Concurrency int    `envconfig:"VONDILISTINGS_WORKER_CONCURRENCY" default:"5"`
	QueueName   string `envconfig:"VONDILISTINGS_WORKER_QUEUE_NAME" default:"listings_indexing"`
}

// FeatureFlags contains feature toggle settings
type FeatureFlags struct {
	AsyncIndexing     bool `envconfig:"VONDILISTINGS_FEATURE_ASYNC_INDEXING" default:"true"`
	ImageOptimization bool `envconfig:"VONDILISTINGS_FEATURE_IMAGE_OPTIMIZATION" default:"true"`
	CacheEnabled      bool `envconfig:"VONDILISTINGS_FEATURE_CACHE_ENABLED" default:"true"`
	RateLimitEnabled  bool `envconfig:"VONDILISTINGS_RATE_LIMIT_ENABLED" default:"true"`
	RateLimitRPS      int  `envconfig:"VONDILISTINGS_RATE_LIMIT_RPS" default:"100"`
	RateLimitBurst    int  `envconfig:"VONDILISTINGS_RATE_LIMIT_BURST" default:"200"`
}

// TracingConfig contains tracing and monitoring settings
type TracingConfig struct {
	Enabled        bool   `envconfig:"VONDILISTINGS_TRACING_ENABLED" default:"false"`
	JaegerEndpoint string `envconfig:"VONDILISTINGS_JAEGER_ENDPOINT" default:"http://localhost:14268/api/traces"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `envconfig:"VONDILISTINGS_CORS_ALLOWED_ORIGINS" default:"http://localhost:3001,http://localhost:3000"`
	AllowedMethods []string `envconfig:"VONDILISTINGS_CORS_ALLOWED_METHODS" default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `envconfig:"VONDILISTINGS_CORS_ALLOWED_HEADERS" default:"Content-Type,Authorization"`
}

// HealthConfig contains health check configuration
type HealthConfig struct {
	CheckTimeout     time.Duration `envconfig:"VONDILISTINGS_HEALTH_CHECK_TIMEOUT" default:"5s"`
	CheckInterval    time.Duration `envconfig:"VONDILISTINGS_HEALTH_CHECK_INTERVAL" default:"30s"`
	StartupTimeout   time.Duration `envconfig:"VONDILISTINGS_HEALTH_STARTUP_TIMEOUT" default:"60s"`
	CacheDuration    time.Duration `envconfig:"VONDILISTINGS_HEALTH_CACHE_DURATION" default:"10s"`
	EnableDeepChecks bool          `envconfig:"VONDILISTINGS_HEALTH_ENABLE_DEEP_CHECKS" default:"true"`
}

// WMSConfig contains WMS (Warehouse Management System) integration settings
type WMSConfig struct {
	Enabled            bool  `envconfig:"VONDILISTINGS_WMS_ENABLED" default:"true"`
	DefaultWarehouseID int64 `envconfig:"VONDILISTINGS_WMS_DEFAULT_WAREHOUSE_ID" default:"1"`
}

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	return godotenv.Load()
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file (ignore error if file doesn't exist - OK for production)
	_ = godotenv.Load()

	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

// Validate performs basic configuration validation
func (c *Config) Validate() error {
	if c.Server.GRPCPort < 1 || c.Server.GRPCPort > 65535 {
		return fmt.Errorf("invalid gRPC port: %d", c.Server.GRPCPort)
	}

	if c.Server.HTTPPort < 1 || c.Server.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Server.HTTPPort)
	}

	if c.DB.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.DB.Name == "" {
		return fmt.Errorf("database name is required")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development" || c.App.Env == "dev"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production" || c.App.Env == "prod"
}
