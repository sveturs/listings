# Listings Microservice

A high-performance, scalable microservice for managing marketplace listings with dual protocol support (gRPC + HTTP REST).

## Features

- **Dual Protocol Support**: gRPC for internal service communication + HTTP REST for external API
- **Distributed Architecture**: Isolated database, independent deployment
- **Async Processing**: Background worker for OpenSearch indexing
- **Caching Layer**: Redis for performance optimization
- **S3 Storage**: MinIO integration for image management
- **Full-Text Search**: OpenSearch integration
- **Observability**: Prometheus metrics, structured logging
- **Production Ready**: Docker support, CI/CD pipeline, comprehensive testing

## Architecture

```
┌─────────────────┐
│  Other Services │
└────────┬────────┘
         │ gRPC (50053)
         │
┌────────▼─────────┐      ┌──────────────┐
│ Listings Service │◄────►│  PostgreSQL  │
│  (Port 8086 HTTP)│      │  (Port 35433)│
└────────┬─────────┘      └──────────────┘
         │
         │                ┌──────────────┐
         ├───────────────►│    Redis     │
         │                │  (Port 36380)│
         │                └──────────────┘
         │
         │                ┌──────────────┐
         ├───────────────►│  OpenSearch  │
         │                │  (Port 9200) │
         │                └──────────────┘
         │
         │                ┌──────────────┐
         └───────────────►│    MinIO     │
                          │  (Port 9000) │
                          └──────────────┘
```

## Tech Stack

- **Language**: Go 1.23+
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Search**: OpenSearch
- **Storage**: MinIO (S3-compatible)
- **RPC**: gRPC
- **HTTP**: Fiber
- **Migrations**: golang-migrate
- **Configuration**: envconfig
- **Logging**: zerolog
- **Metrics**: Prometheus

## Project Structure

```
listings/
├── cmd/server/                 # Application entry point
│   └── main.go
├── internal/                   # Private application code
│   ├── config/                 # Configuration management
│   ├── service/                # Business logic layer
│   │   └── listings/
│   ├── repository/             # Data access layer
│   │   ├── postgres/
│   │   ├── opensearch/
│   │   └── minio/
│   ├── transport/              # HTTP + gRPC handlers
│   │   ├── http/
│   │   └── grpc/
│   └── worker/                 # Async indexing worker
├── pkg/                        # PUBLIC LIBRARY (importable by other services)
│   ├── service/                # Go client
│   ├── http/fiber/middleware/  # Fiber middleware
│   └── grpc/                   # gRPC client
├── api/proto/listings/v1/      # Protobuf definitions
├── migrations/                 # Database migrations
├── docker-compose.yml          # Local development setup
├── Dockerfile                  # Production image
├── Makefile                    # Build automation
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.23 or higher
- Docker & Docker Compose
- Make (optional, but recommended)
- golang-migrate (for migrations)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/sveturs/listings.git
   cd listings
   ```

2. **Copy environment configuration**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start infrastructure** (PostgreSQL + Redis):
   ```bash
   make docker-up
   ```

4. **Run database migrations**:
   ```bash
   make migrate-up
   ```

5. **Install dependencies**:
   ```bash
   make deps
   ```

6. **Build the application**:
   ```bash
   make build
   ```

7. **Run the service**:
   ```bash
   make run
   ```

The service will start on:
- **gRPC**: `localhost:50053`
- **HTTP REST**: `localhost:8086`
- **Metrics**: `localhost:9093/metrics`

### One-Command Setup

For development environment setup:

```bash
make dev
```

This will:
1. Start Docker Compose services
2. Run database migrations
3. Download dependencies

## Development

### Available Commands

```bash
make help                # Show all available commands

# Building
make build               # Build application binary
make build-all           # Build for all platforms

# Testing
make test                # Run unit tests
make test-coverage       # Run tests with coverage report
make test-integration    # Run integration tests
make bench               # Run benchmarks

# Code Quality
make lint                # Run linter
make format              # Format code
make tidy                # Tidy Go modules

# Docker
make docker-build        # Build Docker image
make docker-up           # Start services
make docker-down         # Stop services
make docker-logs         # View logs

# Database
make migrate-up          # Apply migrations
make migrate-down        # Rollback migration
make migrate-reset       # Reset database
make migrate-create      # Create new migration

# Development
make dev                 # Setup dev environment
make dev-reset           # Reset dev environment
make ci                  # Run CI pipeline locally
```

### Creating a New Migration

```bash
make migrate-create NAME=add_listing_tags
```

This creates two files:
- `migrations/000002_add_listing_tags.up.sql`
- `migrations/000002_add_listing_tags.down.sql`

### Running Tests

```bash
# Unit tests
make test

# With coverage
make test-coverage
open coverage.html

# Integration tests (requires running services)
make test-integration
```

### Code Quality Checks

Before committing:

```bash
make pre-commit
```

This runs:
1. Code formatting
2. Linter
3. Tests

## Configuration

All configuration is done via environment variables with the prefix `SVETULISTINGS_`.

Key configuration sections:

### Server Ports
```env
SVETULISTINGS_GRPC_PORT=50053          # gRPC internal API
SVETULISTINGS_HTTP_PORT=8086           # HTTP REST API
SVETULISTINGS_METRICS_PORT=9093        # Prometheus metrics
```

### Database
```env
SVETULISTINGS_DB_HOST=localhost
SVETULISTINGS_DB_PORT=35433
SVETULISTINGS_DB_USER=listings_user
SVETULISTINGS_DB_PASSWORD=listings_password
SVETULISTINGS_DB_NAME=listings_db
```

### Redis
```env
SVETULISTINGS_REDIS_HOST=localhost
SVETULISTINGS_REDIS_PORT=36380
SVETULISTINGS_REDIS_PASSWORD=redis_password
```

See `.env.example` for complete configuration.

## API Documentation

### HTTP REST Endpoints

| Method | Endpoint                      | Description              |
|--------|-------------------------------|--------------------------|
| GET    | `/health`                     | Health check             |
| GET    | `/api/v1/listings`            | List all listings        |
| GET    | `/api/v1/listings/:id`        | Get listing by ID        |
| POST   | `/api/v1/listings`            | Create new listing       |
| PUT    | `/api/v1/listings/:id`        | Update listing           |
| DELETE | `/api/v1/listings/:id`        | Delete listing           |
| GET    | `/api/v1/listings/search`     | Search listings          |
| GET    | `/metrics`                    | Prometheus metrics       |

### gRPC Service

Protobuf definitions are in `api/proto/listings/v1/`.

To regenerate code after proto changes:

```bash
make proto
```

## Deployment

### Docker

Build production image:

```bash
make docker-build
```

Run with Docker Compose:

```bash
docker-compose up -d
```

### Environment-Specific Configuration

- **Development**: `.env` file
- **Staging/Production**: Environment variables injected by orchestration platform

## Database Schema

Main tables:

- **listings**: Core listing entity
- **listing_attributes**: Flexible key-value attributes
- **listing_images**: Image metadata and URLs
- **listing_tags**: Listing tags
- **listing_locations**: Geographic data
- **listing_stats**: Cached statistics
- **indexing_queue**: Async indexing queue

See `migrations/000001_initial_schema.up.sql` for complete schema.

## OpenSearch Integration

The microservice uses OpenSearch for full-text search, faceted filtering, and geo-location queries.

### Quick Setup

1. **Create index**:
   ```bash
   python3 scripts/create_opensearch_index.py
   ```

2. **Reindex data**:
   ```bash
   python3 scripts/reindex_listings.py --target-password <password>
   ```

3. **Validate**:
   ```bash
   python3 scripts/validate_opensearch.py --target-password <password>
   ```

### Configuration

```env
SVETULISTINGS_OPENSEARCH_ADDRESSES=http://localhost:9200
SVETULISTINGS_OPENSEARCH_USERNAME=admin
SVETULISTINGS_OPENSEARCH_PASSWORD=admin
SVETULISTINGS_OPENSEARCH_INDEX=listings_microservice
```

**Note:** Change index name from `marketplace_listings` (monolith) to `listings_microservice` (microservice).

### Features

- **Full-text search** on title and description (Russian stopwords)
- **Autocomplete** support (edge n-grams)
- **Price range filtering** (scaled_float)
- **Category and status faceting**
- **Geo-location queries** (geo_point)
- **Nested objects** for images and attributes

### Documentation

See [OPENSEARCH_SETUP.md](./OPENSEARCH_SETUP.md) for comprehensive guide including:
- Index schema details
- Reindexing procedures
- Search API examples
- Troubleshooting
- Performance tuning

## Monitoring

### Prometheus Metrics

Exposed on port 9093:

```bash
curl http://localhost:9093/metrics
```

Key metrics:
- Request duration
- Request count by endpoint
- Database connection pool stats
- Cache hit/miss ratio
- Worker queue length

### Health Check

```bash
curl http://localhost:8086/health
```

## Performance

- **Connection Pooling**: Configured for optimal DB performance
- **Redis Caching**: 5-minute TTL for listings, 2-minute for search
- **Async Indexing**: Non-blocking OpenSearch updates
- **Worker Concurrency**: 5 concurrent indexing workers

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests and linting (`make pre-commit`)
4. Commit your changes
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Coding Standards

- Follow Go best practices and idioms
- Write tests for new features
- Update documentation
- Run `make pre-commit` before committing

## License

[Your License Here]

## Support

For issues and questions:
- GitHub Issues: [https://github.com/sveturs/listings/issues](https://github.com/sveturs/listings/issues)

## Roadmap

- [ ] Implement business logic (Sprint 4.2)
- [ ] gRPC service implementation
- [ ] HTTP REST API handlers
- [ ] OpenSearch integration
- [ ] MinIO image storage
- [ ] Background worker
- [ ] Comprehensive testing
- [ ] Production deployment

## Related Services

- **Auth Service**: `github.com/sveturs/auth`
- **Main Monolith**: `github.com/sveturs/svetu`
- **Delivery Service**: `github.com/sveturs/delivery`

---

**Version**: 0.1.0 (Sprint 4.1 - Initial Scaffold)
**Last Updated**: 2025-10-31
