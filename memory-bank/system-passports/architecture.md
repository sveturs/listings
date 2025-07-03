# Архитектура системы Sve Tu Platform

## Общая схема взаимодействия слоев

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   MinIO         │    │   OpenSearch    │
│   React/Next.js │    │   Object Store  │    │   Search Engine │
│   TypeScript    │    │                 │    │                 │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ HTTP API             │ File Upload          │ Search Query
          │ snake_case           │ /listings/{file}     │ Denormalized
          │                      │                      │ Data
          ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Backend (Go)                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │   Handlers  │  │   Business  │  │   Storage   │            │
│  │   (Fiber)   │  │    Logic    │  │  Adapters   │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          │ SQL Queries
                          │ snake_case
                          ▼
                ┌─────────────────┐
                │   PostgreSQL    │
                │   Primary DB    │
                │                 │
                └─────────────────┘
```

## Слои системы

### 1. Presentation Layer (Frontend)
- **Технологии**: React 19, Next.js 15, TypeScript
- **Стили**: Tailwind CSS v4, DaisyUI
- **Состояние**: Redux Toolkit
- **Интернационализация**: next-intl (ru/en)

### 2. API Gateway (Backend)
- **Технологии**: Go, Fiber framework
- **Документация**: Swagger/OpenAPI
- **Аутентификация**: JWT + Google OAuth
- **Middleware**: CORS, Logging, Auth

### 3. Business Logic Layer
- **Модули**: Marketplace, Storefronts, Users, Payments
- **Паттерны**: Repository, Service, Handler
- **Валидация**: Go validator

### 4. Data Access Layer
- **Primary Storage**: PostgreSQL
- **Search Engine**: OpenSearch
- **File Storage**: MinIO (S3-compatible)
- **Caching**: Redis (планируется)

## Поток данных

### Создание объявления
```
Frontend → Backend API → PostgreSQL → OpenSearch Reindex → MinIO (images)
```

### Поиск объявлений
```
Frontend → Backend API → OpenSearch → Response (with PostgreSQL IDs)
```

### Загрузка файлов
```
Frontend → Backend API → MinIO → PostgreSQL (metadata) → OpenSearch Reindex
```

## Ключевые принципы архитектуры

### Consistency (Согласованность)
- **Naming Convention**: snake_case во всех слоях для API данных
- **ID Fields**: {entity}_id pattern повсеместно
- **Timestamps**: created_at, updated_at стандарт

### Scalability (Масштабируемость)
- **Separation of Concerns**: четкое разделение ответственности
- **Microservice Ready**: модульная структура backend
- **Stateless API**: RESTful дизайн

### Security (Безопасность)
- **Authentication**: JWT токены
- **Authorization**: Resource ownership validation
- **Input Validation**: На всех уровнях
- **File Upload**: Validation и sanitization

### Performance (Производительность)
- **Search**: Денормализованные данные в OpenSearch
- **Files**: CDN-ready через MinIO
- **Database**: Оптимизированные индексы
- **Frontend**: Static generation где возможно

## Интеграционные паттерны

### API Communication
- **Format**: JSON с snake_case полями
- **Versioning**: Через URL path (/api/v1/)
- **Error Handling**: Структурированные error responses
- **Documentation**: Auto-generated Swagger

### Data Synchronization
- **PostgreSQL → OpenSearch**: Manual reindex command
- **File Uploads**: Transactional (rollback on failure)
- **Cache Invalidation**: Manual после изменений

### Multi-language Support
- **Backend**: Placeholder возврат (e.g., "notifications.error")
- **Frontend**: Translation через next-intl
- **Database**: Dedicated translations table

## Технические ограничения

### File Upload Limits
- **Images**: 10 MB max
- **Videos**: 100 MB max
- **Documents**: 20 MB max

### API Rate Limits
- **Public endpoints**: 100 req/min
- **Authenticated**: 1000 req/min
- **File upload**: 10 files/min

### Database Constraints
- **Connections**: Pool size 25
- **Query timeout**: 30 seconds
- **Transaction timeout**: 60 seconds

## Monitoring & Observability

### Logging
- **Format**: Structured JSON
- **Levels**: Error, Warning, Info, Debug
- **Correlation**: Request ID трacing

### Metrics
- **Application**: Custom business metrics
- **Infrastructure**: Resource utilization
- **API**: Response times, error rates

### Health Checks
- **Database**: Connection + simple query
- **OpenSearch**: Cluster health
- **MinIO**: Bucket access
- **External Services**: Third-party dependencies

## Development & Deployment

### Local Development
- **Frontend**: `yarn dev -p 3001`
- **Backend**: `go run ./cmd/api/main.go`
- **Database**: Docker Compose
- **Services**: Local MinIO + OpenSearch

### Production Deployment
- **Infrastructure**: Docker + Kubernetes
- **Database**: Managed PostgreSQL
- **Storage**: Production MinIO cluster
- **CDN**: Static assets через nginx

### CI/CD Pipeline
- **Testing**: Unit + Integration tests
- **Linting**: Go fmt, ESLint, Prettier
- **Build**: Multi-stage Docker builds
- **Deploy**: Blue-green deployment