package config

import (
	"fmt"
	"time"

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
	Worker   WorkerConfig
	Features FeatureFlags
	Tracing  TracingConfig
	CORS     CORSConfig
}

// AppConfig contains general application settings
type AppConfig struct {
	Env      string `envconfig:"SVETULISTINGS_ENV" default:"development"`
	LogLevel string `envconfig:"SVETULISTINGS_LOG_LEVEL" default:"info"`
	LogFormat string `envconfig:"SVETULISTINGS_LOG_FORMAT" default:"json"`
}

// ServerConfig contains server ports and settings
type ServerConfig struct {
	GRPCHost    string `envconfig:"SVETULISTINGS_GRPC_HOST" default:"0.0.0.0"`
	GRPCPort    int    `envconfig:"SVETULISTINGS_GRPC_PORT" default:"50053"`
	HTTPHost    string `envconfig:"SVETULISTINGS_HTTP_HOST" default:"0.0.0.0"`
	HTTPPort    int    `envconfig:"SVETULISTINGS_HTTP_PORT" default:"8086"`
	MetricsHost string `envconfig:"SVETULISTINGS_METRICS_HOST" default:"0.0.0.0"`
	MetricsPort int    `envconfig:"SVETULISTINGS_METRICS_PORT" default:"9093"`
}

// DBConfig contains PostgreSQL database configuration
type DBConfig struct {
	Host     string `envconfig:"SVETULISTINGS_DB_HOST" default:"localhost"`
	Port     int    `envconfig:"SVETULISTINGS_DB_PORT" default:"35433"`
	User     string `envconfig:"SVETULISTINGS_DB_USER" default:"listings_user"`
	Password string `envconfig:"SVETULISTINGS_DB_PASSWORD" default:"listings_password"`
	Name     string `envconfig:"SVETULISTINGS_DB_NAME" default:"listings_db"`
	SSLMode  string `envconfig:"SVETULISTINGS_DB_SSLMODE" default:"disable"`

	// Connection pool settings
	MaxOpenConns    int           `envconfig:"SVETULISTINGS_DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"SVETULISTINGS_DB_MAX_IDLE_CONNS" default:"10"`
	ConnMaxLifetime time.Duration `envconfig:"SVETULISTINGS_DB_CONN_MAX_LIFETIME" default:"5m"`
	ConnMaxIdleTime time.Duration `envconfig:"SVETULISTINGS_DB_CONN_MAX_IDLE_TIME" default:"10m"`
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
	Host     string `envconfig:"SVETULISTINGS_REDIS_HOST" default:"localhost"`
	Port     int    `envconfig:"SVETULISTINGS_REDIS_PORT" default:"36380"`
	Password string `envconfig:"SVETULISTINGS_REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"SVETULISTINGS_REDIS_DB" default:"0"`

	// Connection pool
	PoolSize     int `envconfig:"SVETULISTINGS_REDIS_POOL_SIZE" default:"10"`
	MinIdleConns int `envconfig:"SVETULISTINGS_REDIS_MIN_IDLE_CONNS" default:"5"`

	// Cache TTL
	ListingTTL time.Duration `envconfig:"SVETULISTINGS_CACHE_LISTING_TTL" default:"5m"`
	SearchTTL  time.Duration `envconfig:"SVETULISTINGS_CACHE_SEARCH_TTL" default:"2m"`
}

// Addr returns Redis address
func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// SearchConfig contains OpenSearch configuration
type SearchConfig struct {
	Addresses []string `envconfig:"SVETULISTINGS_OPENSEARCH_ADDRESSES" default:"http://localhost:9200"`
	Username  string   `envconfig:"SVETULISTINGS_OPENSEARCH_USERNAME" default:"admin"`
	Password  string   `envconfig:"SVETULISTINGS_OPENSEARCH_PASSWORD" default:"admin"`
	Index     string   `envconfig:"SVETULISTINGS_OPENSEARCH_INDEX" default:"marketplace_listings"`
}

// StorageConfig contains MinIO (S3-compatible) configuration
type StorageConfig struct {
	Endpoint  string `envconfig:"SVETULISTINGS_MINIO_ENDPOINT" default:"localhost:9000"`
	AccessKey string `envconfig:"SVETULISTINGS_MINIO_ACCESS_KEY" default:"minioadmin"`
	SecretKey string `envconfig:"SVETULISTINGS_MINIO_SECRET_KEY" default:"minioadmin"`
	UseSSL    bool   `envconfig:"SVETULISTINGS_MINIO_USE_SSL" default:"false"`
	Bucket    string `envconfig:"SVETULISTINGS_MINIO_BUCKET" default:"listings-images"`
}

// AuthConfig contains Auth Service integration settings
type AuthConfig struct {
	ServiceURL    string `envconfig:"SVETULISTINGS_AUTH_SERVICE_URL" default:"http://localhost:8081"`
	PublicKeyPath string `envconfig:"SVETULISTINGS_AUTH_PUBLIC_KEY_PATH" default:"/keys/public.pem"`
}

// WorkerConfig contains async worker settings
type WorkerConfig struct {
	Enabled     bool   `envconfig:"SVETULISTINGS_WORKER_ENABLED" default:"true"`
	Concurrency int    `envconfig:"SVETULISTINGS_WORKER_CONCURRENCY" default:"5"`
	QueueName   string `envconfig:"SVETULISTINGS_WORKER_QUEUE_NAME" default:"listings_indexing"`
}

// FeatureFlags contains feature toggle settings
type FeatureFlags struct {
	AsyncIndexing      bool `envconfig:"SVETULISTINGS_FEATURE_ASYNC_INDEXING" default:"true"`
	ImageOptimization  bool `envconfig:"SVETULISTINGS_FEATURE_IMAGE_OPTIMIZATION" default:"true"`
	CacheEnabled       bool `envconfig:"SVETULISTINGS_FEATURE_CACHE_ENABLED" default:"true"`
	RateLimitEnabled   bool `envconfig:"SVETULISTINGS_RATE_LIMIT_ENABLED" default:"true"`
	RateLimitRPS       int  `envconfig:"SVETULISTINGS_RATE_LIMIT_RPS" default:"100"`
	RateLimitBurst     int  `envconfig:"SVETULISTINGS_RATE_LIMIT_BURST" default:"200"`
}

// TracingConfig contains tracing and monitoring settings
type TracingConfig struct {
	Enabled        bool   `envconfig:"SVETULISTINGS_TRACING_ENABLED" default:"false"`
	JaegerEndpoint string `envconfig:"SVETULISTINGS_JAEGER_ENDPOINT" default:"http://localhost:14268/api/traces"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `envconfig:"SVETULISTINGS_CORS_ALLOWED_ORIGINS" default:"http://localhost:3001,http://localhost:3000"`
	AllowedMethods []string `envconfig:"SVETULISTINGS_CORS_ALLOWED_METHODS" default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `envconfig:"SVETULISTINGS_CORS_ALLOWED_HEADERS" default:"Content-Type,Authorization"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
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
