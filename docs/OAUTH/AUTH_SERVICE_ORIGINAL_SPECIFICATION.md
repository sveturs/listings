# 🔐 Auth Service - Полная спецификация микросервиса аутентификации и авторизации

**Версия документа:** 1.0.0  
**Дата создания:** 06.09.2025  
**Статус:** Утвержден к реализации  
**Авторы:** Architecture Team  
**Время реализации:** ✅ РЕАЛИЗОВАНО!  
**Актуальная документация:** `/data/auth_svetu/docs/ACTUAL_API_ENDPOINTS.md`

---

## 📑 Содержание

1. [Резюме проекта](#-резюме-проекта)
2. [Архитектура системы](#-архитектура-системы)
3. [Техническая спецификация](#-техническая-спецификация)
4. [API спецификация](#-api-спецификация)
5. [Инфраструктура](#-инфраструктура)
6. [База данных](#-база-данных)
7. [Безопасность](#-безопасность)
8. [План миграции](#-план-миграции)
9. [Тестирование](#-тестирование)
10. [Мониторинг и метрики](#-мониторинг-и-метрики)
11. [Развертывание](#-развертывание)
12. [Документация API](#-документация-api)
13. [Риски и митигация](#-риски-и-митигация)
14. [Приложения](#-приложения)

---

## 🎯 Резюме проекта

### Цель проекта
Создание независимого, масштабируемого и безопасного микросервиса аутентификации и авторизации для экосистемы Sve Tu.

### Ключевые требования
- **Производительность:** 10,000+ RPS для валидации токенов
- **Доступность:** 99.99% uptime
- **Безопасность:** Соответствие OWASP Top 10
- **Масштабируемость:** Горизонтальное масштабирование без ограничений
- **Изоляция:** Полная независимость от основного сервиса

### Основные возможности
- JWT-based аутентификация (RS256)
- OAuth 2.0 интеграция (Google, планируется Facebook, Apple)
- Email/Password регистрация и вход
- Refresh token rotation
- Multi-factor authentication (готовность к внедрению)
- Rate limiting и защита от брутфорса
- Device fingerprinting
- Session management

---

## 🏗️ Архитектура системы

### Общая архитектура

```mermaid
graph TB
    subgraph "External Clients"
        Web[Web Frontend]
        Mobile[Mobile Apps]
        API[API Clients]
    end
    
    subgraph "API Gateway Layer"
        APIGW[API Gateway / Load Balancer]
    end
    
    subgraph "Auth Service"
        REST[REST API :8080]
        GRPC[gRPC Server :50051]
        Core[Auth Core Logic]
        OAuth[OAuth Handlers]
        Token[Token Service]
        User[User Service]
    end
    
    subgraph "Data Layer"
        PG[(PostgreSQL)]
        Redis[(Redis Cache)]
        Vault[HashiCorp Vault]
    end
    
    subgraph "Main Service"
        Backend[Backend Service]
        MainDB[(Main Database)]
    end
    
    Web --> APIGW
    Mobile --> APIGW
    API --> APIGW
    
    APIGW --> REST
    Backend --> GRPC
    
    REST --> Core
    GRPC --> Core
    
    Core --> OAuth
    Core --> Token
    Core --> User
    
    Token --> PG
    Token --> Redis
    User --> PG
    User --> Redis
    OAuth --> Vault
```

### Компонентная архитектура

```
auth-service/
├── cmd/
│   ├── server/
│   │   └── main.go              # Точка входа приложения
│   └── migrate/
│       └── main.go              # Утилита миграций БД
├── internal/
│   ├── config/                 # Конфигурация
│   │   ├── config.go
│   │   └── validation.go
│   ├── domain/                 # Доменные модели
│   │   ├── user.go
│   │   ├── token.go
│   │   ├── session.go
│   │   └── errors.go
│   ├── service/                # Бизнес-логика
│   │   ├── auth/
│   │   │   ├── service.go
│   │   │   ├── login.go
│   │   │   ├── register.go
│   │   │   └── logout.go
│   │   ├── token/
│   │   │   ├── jwt.go
│   │   │   ├── refresh.go
│   │   │   └── validation.go
│   │   ├── user/
│   │   │   ├── service.go
│   │   │   └── profile.go
│   │   └── oauth/
│   │       ├── google.go
│   │       ├── facebook.go
│   │       └── provider.go
│   ├── repository/             # Слой данных
│   │   ├── postgres/
│   │   │   ├── user.go
│   │   │   ├── token.go
│   │   │   └── migrations/
│   │   └── redis/
│   │       ├── cache.go
│   │       ├── session.go
│   │       └── ratelimit.go
│   ├── transport/              # API слой
│   │   ├── grpc/
│   │   │   ├── server.go
│   │   │   ├── auth.go
│   │   │   └── interceptors.go
│   │   ├── http/
│   │   │   ├── server.go
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── validators/
│   │   └── events/
│   │       ├── publisher.go
│   │       └── types.go
│   └── pkg/                    # Утилиты
│       ├── crypto/
│       ├── logger/
│       ├── metrics/
│       └── tracer/
├── pkg/                         # Публичные пакеты
│   ├── client/                  # Client SDK
│   │   ├── grpc/
│   │   └── http/
│   └── proto/                   # Proto определения
│       ├── auth.proto
│       └── generated/
├── deployments/                 # Конфигурации развертывания
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── configmap.yaml
│   └── terraform/
│       └── main.tf
├── scripts/                     # Вспомогательные скрипты
│   ├── migrate.sh
│   ├── generate-keys.sh
│   └── test-load.sh
└── tests/                       # Тесты
    ├── unit/
    ├── integration/
    ├── e2e/
    └── load/
```

---

## 💻 Техническая спецификация

### Технологический стек

| Компонент | Технология | Версия | Обоснование |
|-----------|-----------|--------|-------------|
| Язык программирования | Go | 1.22+ | Производительность, простота, concurrency |
| Web Framework | Fiber | v2.52+ | Высокая производительность, Express-like API |
| gRPC Framework | grpc-go | v1.60+ | Стандарт для микросервисов |
| База данных | PostgreSQL | 15+ | Надежность, ACID, JSON support |
| Кеш | Redis | 7+ | Производительность, TTL, pub/sub |
| Message Queue | NATS | 2.10+ | Легковесность, производительность |
| Secrets Management | HashiCorp Vault | 1.15+ | Безопасное хранение ключей |
| Контейнеризация | Docker | 24+ | Стандарт индустрии |
| Оркестрация | Kubernetes | 1.28+ | Масштабирование, управление |
| API Gateway | Kong | 3.5+ | Rate limiting, authentication |
| Мониторинг | Prometheus + Grafana | Latest | Метрики и визуализация |
| Трассировка | Jaeger | 1.50+ | Distributed tracing |
| Логирование | ELK Stack | 8+ | Centralized logging |

### Системные требования

#### Минимальные требования (Development)
- CPU: 2 cores
- RAM: 4 GB
- Storage: 20 GB SSD
- Network: 100 Mbps

#### Рекомендуемые требования (Production)
- CPU: 8 cores
- RAM: 16 GB
- Storage: 100 GB SSD (RAID 10)
- Network: 1 Gbps
- Load Balancer: 2+ instances
- Database: Master-Slave replication
- Redis: Sentinel setup

### Производительность

| Метрика | Требование | Текущее | Цель |
|---------|-----------|---------|------|
| Token Validation Latency | < 50ms | 200ms | 30ms |
| Login Latency | < 200ms | 500ms | 150ms |
| Throughput (RPS) | 10,000 | 2,000 | 15,000 |
| Concurrent Users | 100,000 | 10,000 | 500,000 |
| Database Connections | 100 | 50 | 200 |
| Cache Hit Rate | > 90% | 60% | 95% |
| Error Rate | < 0.1% | 1% | 0.05% |

---

## 📡 API Спецификация

### REST API Endpoints

#### Authentication Endpoints

```yaml
# Регистрация
POST /api/v1/auth/register
Request:
  Content-Type: application/json
  Body:
    {
      "email": "user@example.com",
      "password": "SecurePass123!",
      "name": "John Doe",
      "terms_accepted": true
    }
Response:
  200 OK:
    {
      "access_token": "eyJhbGc...",
      "refresh_token": "eyJhbGc...",
      "token_type": "Bearer",
      "expires_in": 900,
      "user": {
        "id": 12345,
        "email": "user@example.com",
        "name": "John Doe",
        "email_verified": false,
        "created_at": "2025-09-06T10:00:00Z"
      }
    }
  400 Bad Request:
    {
      "error": "validation_error",
      "message": "Invalid email format",
      "details": {...}
    }
  409 Conflict:
    {
      "error": "user_exists",
      "message": "User with this email already exists"
    }

# Вход
POST /api/v1/auth/login
Request:
  Content-Type: application/json
  Body:
    {
      "email": "user@example.com",
      "password": "SecurePass123!",
      "device_id": "uuid-device-id",
      "device_name": "Chrome on MacOS"
    }
Response:
  200 OK: (same as register)
  401 Unauthorized:
    {
      "error": "invalid_credentials",
      "message": "Invalid email or password"
    }
  429 Too Many Requests:
    {
      "error": "rate_limit",
      "message": "Too many login attempts",
      "retry_after": 900
    }

# Обновление токенов
POST /api/v1/auth/refresh
Request:
  Cookie: refresh_token=eyJhbGc...
Response:
  200 OK:
    {
      "access_token": "eyJhbGc...",
      "refresh_token": "eyJhbGc...",
      "token_type": "Bearer",
      "expires_in": 900
    }
  401 Unauthorized:
    {
      "error": "invalid_token",
      "message": "Refresh token is invalid or expired"
    }

# Выход
POST /api/v1/auth/logout
Request:
  Authorization: Bearer eyJhbGc...
  Cookie: refresh_token=eyJhbGc...
Response:
  200 OK:
    {
      "message": "Successfully logged out"
    }

# Валидация токена
GET /api/v1/auth/validate
Request:
  Authorization: Bearer eyJhbGc...
Response:
  200 OK:
    {
      "valid": true,
      "user_id": 12345,
      "email": "user@example.com",
      "roles": ["user"],
      "expires_at": "2025-09-06T11:00:00Z"
    }
  401 Unauthorized:
    {
      "valid": false,
      "error": "token_expired"
    }
```

#### OAuth Endpoints

```yaml
# Инициация OAuth
GET /api/v1/auth/oauth/{provider}
Parameters:
  provider: google|facebook|apple
  redirect_uri: https://app.svetu.rs/auth/callback
  state: random-state-string
Response:
  302 Redirect to OAuth provider

# OAuth Callback
GET /api/v1/auth/oauth/{provider}/callback
Parameters:
  code: authorization-code
  state: random-state-string
Response:
  302 Redirect to frontend with tokens in URL params
```

#### User Management

```yaml
# Получение профиля
GET /api/v1/auth/user/profile
Request:
  Authorization: Bearer eyJhbGc...
Response:
  200 OK:
    {
      "id": 12345,
      "email": "user@example.com",
      "name": "John Doe",
      "picture_url": "https://...",
      "email_verified": true,
      "phone_verified": false,
      "two_factor_enabled": false,
      "created_at": "2025-09-06T10:00:00Z",
      "updated_at": "2025-09-06T10:00:00Z"
    }

# Обновление профиля
PATCH /api/v1/auth/user/profile
Request:
  Authorization: Bearer eyJhbGc...
  Content-Type: application/json
  Body:
    {
      "name": "Jane Doe",
      "phone": "+1234567890"
    }
Response:
  200 OK: (updated profile)

# Смена пароля
POST /api/v1/auth/user/change-password
Request:
  Authorization: Bearer eyJhbGc...
  Content-Type: application/json
  Body:
    {
      "current_password": "OldPass123!",
      "new_password": "NewPass456!",
      "confirm_password": "NewPass456!"
    }
Response:
  200 OK:
    {
      "message": "Password changed successfully"
    }

# Удаление аккаунта
DELETE /api/v1/auth/user/account
Request:
  Authorization: Bearer eyJhbGc...
  Content-Type: application/json
  Body:
    {
      "password": "CurrentPass123!",
      "confirmation": "DELETE MY ACCOUNT"
    }
Response:
  200 OK:
    {
      "message": "Account scheduled for deletion",
      "deletion_date": "2025-10-06T10:00:00Z"
    }
```

### gRPC API

```protobuf
syntax = "proto3";

package auth.v1;

service AuthService {
  // Authentication
  rpc Login(LoginRequest) returns (AuthResponse);
  rpc Register(RegisterRequest) returns (AuthResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (AuthResponse);
  
  // Token Validation
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc RevokeToken(RevokeTokenRequest) returns (RevokeTokenResponse);
  
  // User Management
  rpc GetUser(GetUserRequest) returns (User);
  rpc UpdateUser(UpdateUserRequest) returns (User);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  
  // Session Management
  rpc GetUserSessions(GetUserSessionsRequest) returns (GetUserSessionsResponse);
  rpc RevokeSession(RevokeSessionRequest) returns (RevokeSessionResponse);
  
  // Internal
  rpc CreateServiceAccount(CreateServiceAccountRequest) returns (ServiceAccount);
  rpc RotateServiceAccountKey(RotateKeyRequest) returns (ServiceAccount);
}

message LoginRequest {
  string email = 1;
  string password = 2;
  string device_id = 3;
  string device_name = 4;
  string ip_address = 5;
  string user_agent = 6;
}

message AuthResponse {
  string access_token = 1;
  string refresh_token = 2;
  string token_type = 3;
  int32 expires_in = 4;
  User user = 5;
}

message User {
  int32 id = 1;
  string email = 2;
  string name = 3;
  string picture_url = 4;
  bool email_verified = 5;
  repeated string roles = 6;
  map<string, string> metadata = 7;
  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;
}

message ValidateTokenRequest {
  string token = 1;
  bool check_revoked = 2;
}

message ValidateTokenResponse {
  bool valid = 1;
  int32 user_id = 2;
  string email = 3;
  repeated string roles = 4;
  map<string, string> claims = 5;
  google.protobuf.Timestamp expires_at = 6;
}
```

---

## 🏢 Инфраструктура

### Docker Compose Configuration

```yaml
version: '3.8'

networks:
  auth-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16

volumes:
  auth_postgres_data:
  auth_redis_data:
  auth_vault_data:

services:
  # PostgreSQL
  auth-postgres:
    image: postgres:15-alpine
    container_name: auth_postgres
    hostname: auth-postgres
    networks:
      auth-network:
        ipv4_address: 172.20.0.10
    environment:
      POSTGRES_DB: ${AUTH_DB_NAME:-auth_db}
      POSTGRES_USER: ${AUTH_DB_USER:-auth_user}
      POSTGRES_PASSWORD: ${AUTH_DB_PASSWORD}
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --locale=en_US.UTF-8"
    volumes:
      - auth_postgres_data:/var/lib/postgresql/data
      - ./scripts/postgres-init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    ports:
      - "5433:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${AUTH_DB_USER:-auth_user} -d ${AUTH_DB_NAME:-auth_db}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # Redis
  auth-redis:
    image: redis:7-alpine
    container_name: auth_redis
    hostname: auth-redis
    networks:
      auth-network:
        ipv4_address: 172.20.0.11
    command: >
      redis-server
      --requirepass ${AUTH_REDIS_PASSWORD}
      --maxmemory 512mb
      --maxmemory-policy allkeys-lru
      --appendonly yes
      --appendfsync everysec
    volumes:
      - auth_redis_data:/data
    ports:
      - "6380:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "--auth", "${AUTH_REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # HashiCorp Vault (for secrets)
  auth-vault:
    image: vault:1.15
    container_name: auth_vault
    hostname: auth-vault
    networks:
      auth-network:
        ipv4_address: 172.20.0.12
    environment:
      VAULT_DEV_ROOT_TOKEN_ID: ${VAULT_ROOT_TOKEN}
      VAULT_DEV_LISTEN_ADDRESS: 0.0.0.0:8200
    volumes:
      - auth_vault_data:/vault/data
    ports:
      - "8200:8200"
    cap_add:
      - IPC_LOCK
    restart: unless-stopped

  # NATS (for events)
  auth-nats:
    image: nats:2.10-alpine
    container_name: auth_nats
    hostname: auth-nats
    networks:
      auth-network:
        ipv4_address: 172.20.0.13
    command: >
      -js
      -sd /data
      --user ${NATS_USER}
      --pass ${NATS_PASSWORD}
    volumes:
      - ./data/nats:/data
    ports:
      - "4222:4222"  # Client connections
      - "8222:8222"  # HTTP monitoring
    restart: unless-stopped

  # Auth Service
  auth-service:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile
      args:
        - GO_VERSION=1.22
    container_name: auth_service
    hostname: auth-service
    networks:
      auth-network:
        ipv4_address: 172.20.0.20
    environment:
      # Database
      DB_HOST: auth-postgres
      DB_PORT: 5432
      DB_NAME: ${AUTH_DB_NAME:-auth_db}
      DB_USER: ${AUTH_DB_USER:-auth_user}
      DB_PASSWORD: ${AUTH_DB_PASSWORD}
      DB_SSL_MODE: disable
      DB_MAX_CONNECTIONS: 100
      DB_MAX_IDLE_CONNECTIONS: 10
      DB_CONNECTION_MAX_LIFETIME: 3600
      
      # Redis
      REDIS_HOST: auth-redis
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${AUTH_REDIS_PASSWORD}
      REDIS_DB: 0
      REDIS_MAX_RETRIES: 3
      REDIS_POOL_SIZE: 100
      
      # Vault
      VAULT_ADDR: http://auth-vault:8200
      VAULT_TOKEN: ${VAULT_ROOT_TOKEN}
      
      # NATS
      NATS_URL: nats://${NATS_USER}:${NATS_PASSWORD}@auth-nats:4222
      
      # Service Configuration
      SERVICE_NAME: auth-service
      SERVICE_VERSION: ${VERSION:-1.0.0}
      
      # Server Ports
      HTTP_PORT: 8080
      GRPC_PORT: 50051
      METRICS_PORT: 9090
      HEALTH_PORT: 8081
      
      # JWT Configuration
      JWT_ALGORITHM: RS256
      JWT_ACCESS_TOKEN_DURATION: 15m
      JWT_REFRESH_TOKEN_DURATION: 30d
      JWT_ISSUER: https://auth.svetu.rs
      JWT_AUDIENCE: https://svetu.rs
      
      # OAuth Providers
      GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID}
      GOOGLE_CLIENT_SECRET: ${GOOGLE_CLIENT_SECRET}
      GOOGLE_REDIRECT_URL: ${GOOGLE_REDIRECT_URL}
      
      # Security
      BCRYPT_COST: 12
      RATE_LIMIT_ENABLED: true
      RATE_LIMIT_LOGIN: 5/15m
      RATE_LIMIT_REGISTER: 3/1h
      RATE_LIMIT_TOKEN_REFRESH: 30/1m
      
      # Monitoring
      TRACING_ENABLED: true
      TRACING_ENDPOINT: http://jaeger:14268/api/traces
      METRICS_ENABLED: true
      
      # Logging
      LOG_LEVEL: ${LOG_LEVEL:-info}
      LOG_FORMAT: json
      
      # Environment
      ENV: ${ENV:-development}
    ports:
      - "8080:8080"    # HTTP/REST
      - "50051:50051"  # gRPC
      - "9090:9090"    # Metrics
      - "8081:8081"    # Health
    depends_on:
      auth-postgres:
        condition: service_healthy
      auth-redis:
        condition: service_healthy
      auth-vault:
        condition: service_started
      auth-nats:
        condition: service_started
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G

  # Jaeger (Tracing)
  jaeger:
    image: jaegertracing/all-in-one:1.50
    container_name: auth_jaeger
    networks:
      auth-network:
        ipv4_address: 172.20.0.30
    environment:
      COLLECTOR_OTLP_ENABLED: true
    ports:
      - "16686:16686"  # Jaeger UI
      - "14268:14268"  # Collector
    restart: unless-stopped
```

### Kubernetes Deployment

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: auth-service

---
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: auth-service-config
  namespace: auth-service
data:
  SERVICE_NAME: auth-service
  HTTP_PORT: "8080"
  GRPC_PORT: "50051"
  LOG_LEVEL: info
  JWT_ALGORITHM: RS256
  JWT_ISSUER: https://auth.svetu.rs

---
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: auth-service-secrets
  namespace: auth-service
type: Opaque
stringData:
  DB_PASSWORD: "change-me"
  REDIS_PASSWORD: "change-me"
  JWT_PRIVATE_KEY: |
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----
  GOOGLE_CLIENT_SECRET: "change-me"

---
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: auth-service
  labels:
    app: auth-service
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: auth-service
      containers:
      - name: auth-service
        image: svetu/auth-service:1.0.0
        imagePullPolicy: Always
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        - name: grpc
          containerPort: 50051
          protocol: TCP
        - name: metrics
          containerPort: 9090
          protocol: TCP
        env:
        - name: DB_HOST
          value: auth-postgres-service
        - name: REDIS_HOST
          value: auth-redis-service
        envFrom:
        - configMapRef:
            name: auth-service-config
        - secretRef:
            name: auth-service-secrets
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8081
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        volumeMounts:
        - name: jwt-keys
          mountPath: /etc/auth/keys
          readOnly: true
      volumes:
      - name: jwt-keys
        secret:
          secretName: jwt-keys

---
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: auth-service
  labels:
    app: auth-service
spec:
  type: ClusterIP
  selector:
    app: auth-service
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: grpc
    port: 50051
    targetPort: 50051
  - name: metrics
    port: 9090
    targetPort: 9090

---
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
  namespace: auth-service
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
```

---

## 🗄️ База данных

### Схема базы данных

```sql
-- Create database
CREATE DATABASE auth_db
    WITH 
    OWNER = auth_user
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

-- Connect to auth_db
\c auth_db;

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create schema
CREATE SCHEMA IF NOT EXISTS auth;

-- Set search path
SET search_path TO auth, public;

-- Users table
CREATE TABLE auth.users (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_normalized VARCHAR(255) GENERATED ALWAYS AS (lower(email)) STORED,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255),
    
    -- OAuth providers
    google_id VARCHAR(255) UNIQUE,
    facebook_id VARCHAR(255) UNIQUE,
    apple_id VARCHAR(255) UNIQUE,
    
    -- Profile
    picture_url TEXT,
    phone VARCHAR(20),
    phone_verified BOOLEAN DEFAULT FALSE,
    
    -- Security
    email_verified BOOLEAN DEFAULT FALSE,
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    two_factor_secret VARCHAR(255),
    backup_codes TEXT[],
    
    -- Status
    provider VARCHAR(20) NOT NULL CHECK (provider IN ('local', 'google', 'facebook', 'apple')),
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP,
    
    -- Metadata
    last_login_at TIMESTAMP,
    last_login_ip INET,
    last_login_user_agent TEXT,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_email_normalized (email_normalized),
    INDEX idx_provider (provider),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_last_login (last_login_at DESC)
);

-- Refresh tokens table
CREATE TABLE auth.refresh_tokens (
    id SERIAL PRIMARY KEY,
    jti UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    
    -- Token family for rotation detection
    family_id UUID NOT NULL,
    parent_token_id INTEGER REFERENCES auth.refresh_tokens(id),
    
    -- Metadata
    device_id VARCHAR(255),
    device_name VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    
    -- Expiration and revocation
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    revoked_reason VARCHAR(50),
    
    -- Usage tracking
    used_at TIMESTAMP,
    use_count INTEGER DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_user_id (user_id) WHERE NOT is_revoked,
    INDEX idx_token_hash (token_hash) WHERE NOT is_revoked,
    INDEX idx_family_id (family_id),
    INDEX idx_expires_at (expires_at) WHERE NOT is_revoked,
    INDEX idx_device_id (device_id, user_id)
);

-- Revoked access tokens (blacklist)
CREATE TABLE auth.revoked_access_tokens (
    jti VARCHAR(255) PRIMARY KEY,
    user_id INTEGER REFERENCES auth.users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_reason VARCHAR(50),
    
    -- Indexes
    INDEX idx_expires_at_revoked (expires_at),
    INDEX idx_user_id_revoked (user_id)
);

-- User sessions (active sessions tracking)
CREATE TABLE auth.user_sessions (
    id SERIAL PRIMARY KEY,
    session_id UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    refresh_token_id INTEGER REFERENCES auth.refresh_tokens(id) ON DELETE CASCADE,
    
    -- Session info
    device_id VARCHAR(255),
    device_name VARCHAR(100),
    device_type VARCHAR(50),
    browser VARCHAR(50),
    os VARCHAR(50),
    ip_address INET,
    location_country VARCHAR(2),
    location_city VARCHAR(100),
    
    -- Activity
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    activity_count INTEGER DEFAULT 1,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    terminated_at TIMESTAMP,
    termination_reason VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    
    -- Indexes
    INDEX idx_user_sessions (user_id, is_active),
    INDEX idx_session_id (session_id) WHERE is_active,
    INDEX idx_last_activity (last_activity_at DESC)
);

-- Audit log
CREATE TABLE auth.audit_log (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES auth.users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(50),
    resource_id VARCHAR(255),
    
    -- Request info
    ip_address INET,
    user_agent TEXT,
    request_id UUID,
    
    -- Details
    old_values JSONB,
    new_values JSONB,
    metadata JSONB,
    
    -- Result
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    
    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_audit_user (user_id, created_at DESC),
    INDEX idx_audit_action (action, created_at DESC),
    INDEX idx_audit_resource (resource, resource_id),
    INDEX idx_audit_created (created_at DESC)
) PARTITION BY RANGE (created_at);

-- Create partitions for audit log (monthly)
CREATE TABLE auth.audit_log_2025_09 PARTITION OF auth.audit_log
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');
CREATE TABLE auth.audit_log_2025_10 PARTITION OF auth.audit_log
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

-- Login attempts tracking
CREATE TABLE auth.login_attempts (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    failure_reason VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_login_email (email, created_at DESC),
    INDEX idx_login_ip (ip_address, created_at DESC),
    INDEX idx_login_created (created_at DESC)
);

-- Service accounts (for inter-service communication)
CREATE TABLE auth.service_accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    api_key_hash VARCHAR(255) UNIQUE NOT NULL,
    roles TEXT[] DEFAULT '{}',
    
    -- Rate limiting
    rate_limit_requests INTEGER DEFAULT 10000,
    rate_limit_window INTEGER DEFAULT 3600,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    last_used_at TIMESTAMP,
    
    -- Expiration
    expires_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    rotated_at TIMESTAMP,
    
    -- Indexes
    INDEX idx_service_api_key (api_key_hash) WHERE is_active
);

-- Functions and triggers

-- Update updated_at timestamp
CREATE OR REPLACE FUNCTION auth.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON auth.users
    FOR EACH ROW EXECUTE FUNCTION auth.update_updated_at_column();

CREATE TRIGGER update_service_accounts_updated_at BEFORE UPDATE ON auth.service_accounts
    FOR EACH ROW EXECUTE FUNCTION auth.update_updated_at_column();

-- Clean expired tokens
CREATE OR REPLACE FUNCTION auth.clean_expired_tokens()
RETURNS void AS $$
BEGIN
    -- Delete expired refresh tokens
    DELETE FROM auth.refresh_tokens 
    WHERE expires_at < CURRENT_TIMESTAMP 
    AND created_at < CURRENT_TIMESTAMP - INTERVAL '7 days';
    
    -- Delete expired revoked access tokens
    DELETE FROM auth.revoked_access_tokens 
    WHERE expires_at < CURRENT_TIMESTAMP;
    
    -- Delete old login attempts
    DELETE FROM auth.login_attempts 
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '30 days';
    
    -- Deactivate expired sessions
    UPDATE auth.user_sessions 
    SET is_active = FALSE, 
        terminated_at = CURRENT_TIMESTAMP,
        termination_reason = 'expired'
    WHERE expires_at < CURRENT_TIMESTAMP 
    AND is_active = TRUE;
END;
$$ LANGUAGE plpgsql;

-- Detect token reuse
CREATE OR REPLACE FUNCTION auth.detect_token_reuse()
RETURNS TRIGGER AS $$
BEGIN
    -- If this token was already used, revoke entire family
    IF NEW.use_count > 1 THEN
        UPDATE auth.refresh_tokens 
        SET is_revoked = TRUE,
            revoked_at = CURRENT_TIMESTAMP,
            revoked_reason = 'reuse_detected'
        WHERE family_id = NEW.family_id;
        
        -- Log security event
        INSERT INTO auth.audit_log (user_id, action, resource, resource_id, metadata)
        VALUES (
            NEW.user_id,
            'token_reuse_detected',
            'refresh_token',
            NEW.jti::text,
            jsonb_build_object(
                'family_id', NEW.family_id,
                'ip_address', NEW.ip_address,
                'user_agent', NEW.user_agent
            )
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER detect_refresh_token_reuse 
    AFTER UPDATE OF use_count ON auth.refresh_tokens
    FOR EACH ROW 
    WHEN (NEW.use_count > 1)
    EXECUTE FUNCTION auth.detect_token_reuse();

-- Permissions
GRANT ALL PRIVILEGES ON SCHEMA auth TO auth_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO auth_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA auth TO auth_user;
```

### Redis структура данных

```redis
# User cache
KEY: user:{id}
TYPE: STRING (JSON)
TTL: 300 seconds (5 minutes)
VALUE: {
  "id": 12345,
  "email": "user@example.com",
  "name": "John Doe",
  "roles": ["user"],
  "email_verified": true
}

# Active sessions set
KEY: sessions:{user_id}
TYPE: SET
TTL: None (managed manually)
MEMBERS: ["session_id_1", "session_id_2", ...]

# Session data
KEY: session:{session_id}
TYPE: HASH
TTL: 86400 seconds (24 hours)
FIELDS:
  user_id: "12345"
  device_id: "uuid"
  ip_address: "192.168.1.1"
  last_activity: "1693900000"

# Rate limiting
KEY: rate:login:{email}
TYPE: STRING (counter)
TTL: 900 seconds (15 minutes)
VALUE: "3"

KEY: rate:register:{ip}
TYPE: STRING (counter)
TTL: 3600 seconds (1 hour)
VALUE: "2"

KEY: rate:refresh:{user_id}
TYPE: STRING (counter)
TTL: 60 seconds
VALUE: "10"

# Token blacklist (for immediate revocation)
KEY: blacklist:{jti}
TYPE: STRING
TTL: Until token expiration
VALUE: "1"

# Email verification codes
KEY: verify:email:{code}
TYPE: STRING
TTL: 3600 seconds (1 hour)
VALUE: "user@example.com"

# Password reset tokens
KEY: reset:{token}
TYPE: HASH
TTL: 900 seconds (15 minutes)
FIELDS:
  user_id: "12345"
  email: "user@example.com"
  ip_address: "192.168.1.1"

# Two-factor authentication
KEY: 2fa:{user_id}:{code}
TYPE: STRING
TTL: 300 seconds (5 minutes)
VALUE: "1"

# Failed login tracking
KEY: failed:{email}
TYPE: LIST
TTL: 900 seconds (15 minutes)
VALUES: ["timestamp1", "timestamp2", ...]

# Device fingerprints
KEY: device:{user_id}:{device_id}
TYPE: HASH
TTL: 2592000 seconds (30 days)
FIELDS:
  name: "Chrome on MacOS"
  last_seen: "1693900000"
  trusted: "true"
```

---

## 🔒 Безопасность

### Механизмы защиты

#### 1. Криптография
```go
// Хеширование паролей
bcrypt.GenerateFromPassword(password, 12)

// JWT подпись (RS256)
privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
signedToken, _ := token.SignedString(privateKey)

// Шифрование sensitive данных в БД
cipher, _ := aes.NewCipher(key)
gcm, _ := cipher.NewGCM(cipher)
encrypted := gcm.Seal(nil, nonce, plaintext, nil)
```

#### 2. Rate Limiting
```yaml
Login: 5 attempts per 15 minutes per email
Register: 3 attempts per hour per IP
Token Refresh: 30 requests per minute per user
API General: 1000 requests per minute per IP
Password Reset: 3 attempts per hour per email
Email Verification: 5 attempts per day per email
```

#### 3. Token Security
- Access Token: 15 минут, RS256, non-revocable
- Refresh Token: 30 дней, rotation on use, revocable
- Token Family: Detection of token reuse
- JTI (JWT ID): Unique identifier for blacklisting

#### 4. Session Security
- Device Fingerprinting
- IP Address validation
- User Agent tracking
- Geographical anomaly detection
- Concurrent session limits

#### 5. Input Validation
```go
// Email validation
emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Password requirements
minLength := 8
requireUppercase := true
requireLowercase := true
requireNumber := true
requireSpecial := true

// SQL Injection prevention
query := "SELECT * FROM users WHERE email = $1"
db.QueryRow(query, email) // Параметризованные запросы

// XSS Prevention
html.EscapeString(userInput)
```

#### 6. Headers Security
```go
// Security headers
app.Use(helmet.New(helmet.Config{
    XSSProtection:         "1; mode=block",
    ContentTypeNosniff:    "nosniff",
    XFrameOptions:         "DENY",
    HSTSMaxAge:           31536000,
    HSTSIncludeSubdomains: true,
    HSTSPreload:          true,
    ContentSecurityPolicy: "default-src 'self'",
    ReferrerPolicy:       "strict-origin-when-cross-origin",
}))
```

### Compliance

#### GDPR
- Right to be forgotten (soft delete + anonymization)
- Data portability (export user data)
- Consent management
- Privacy by design
- Data minimization

#### OWASP Top 10
- A01: Broken Access Control ✓
- A02: Cryptographic Failures ✓
- A03: Injection ✓
- A04: Insecure Design ✓
- A05: Security Misconfiguration ✓
- A06: Vulnerable Components ✓
- A07: Identification and Authentication Failures ✓
- A08: Software and Data Integrity Failures ✓
- A09: Security Logging and Monitoring Failures ✓
- A10: Server-Side Request Forgery ✓

---

## 📋 План миграции

### Фаза 0: Подготовка (Неделя 0)

#### Задачи
- [ ] Утверждение спецификации
- [ ] Создание репозитория
- [ ] Настройка CI/CD pipeline
- [ ] Подготовка инфраструктуры
- [ ] Создание проектной команды

#### Deliverables
- Git repository
- CI/CD pipeline (GitHub Actions)
- Development environment
- Project board (Jira/GitHub Projects)

### Фаза 1: Очистка Legacy (Неделя 1)

#### День 1-2: Аудит
```bash
# Найти все использования session_token
grep -r "session_token" backend/ --include="*.go"
grep -r "session_token" frontend/ --include="*.ts" --include="*.tsx"

# Документировать endpoints
curl http://localhost:3000/swagger/index.html > current_api.html

# Backup данных
pg_dump -h localhost -U postgres svetubd > backup_before_migration.sql
```

#### День 3-4: Удаление session кода
```go
// Удалить из backend:
// - sync.Map sessions
// - session_token проверки
// - legacy middleware

// Удалить из frontend:
// - session_token из cookies
// - fallback логику
```

#### День 5: Тестирование
```bash
# Unit tests
go test ./...

# Integration tests
npm run test:e2e

# Load testing
k6 run scripts/load-test.js
```

### Фаза 2: Разработка Auth Service (Неделя 2)

#### День 1: Структура проекта
```bash
# Создание структуры
mkdir -p auth-service/{cmd,internal,pkg,deployments,scripts,tests}

# Инициализация Go модуля
cd auth-service
go mod init github.com/svetu/auth-service

# Основные зависимости
go get github.com/gofiber/fiber/v2
go get google.golang.org/grpc
go get github.com/golang-jwt/jwt/v5
go get gorm.io/gorm
go get github.com/redis/go-redis/v9
```

#### День 2-3: Core сервисы
```go
// Implement:
// - AuthService (login, register, logout)
// - TokenService (generate, validate, refresh)
// - UserService (CRUD operations)
// - OAuthService (Google, Facebook)
```

#### День 4: API Layer
```go
// REST endpoints
// gRPC services
// GraphQL schema (optional)
```

#### День 5: Database и Cache
```sql
-- Run migrations
migrate -path migrations -database "postgres://..." up

-- Seed test data
psql auth_db < seed.sql
```

### Фаза 3: Интеграция (Неделя 3)

#### День 1-2: Client Libraries
```go
// Go client
package authclient

type Client interface {
    ValidateToken(ctx context.Context, token string) (*User, error)
    RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error)
}

// TypeScript client
export class AuthClient {
    async validateToken(token: string): Promise<User>
    async refreshToken(refreshToken: string): Promise<Tokens>
}
```

#### День 3-4: Backend Integration
```go
// Replace internal auth with client
authClient := authclient.New(authServiceURL)

// Update middleware
func AuthMiddleware(client authclient.Client) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractToken(c)
        user, err := client.ValidateToken(c.Context(), token)
        if err != nil {
            return c.Status(401).JSON(errorResponse)
        }
        c.Locals("user", user)
        return c.Next()
    }
}
```

#### День 5: Frontend Integration
```typescript
// Update auth service
class AuthService {
    private readonly API_URL = process.env.NEXT_PUBLIC_AUTH_SERVICE_URL
    
    async login(email: string, password: string): Promise<AuthResponse> {
        const response = await fetch(`${this.API_URL}/api/v1/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        })
        return response.json()
    }
}
```

### Фаза 4: Миграция данных (Неделя 4)

#### День 1: Подготовка скриптов
```go
// migrate-users.go
func MigrateUsers() error {
    sourceDB := connectToSource()
    targetDB := connectToTarget()
    
    users := sourceDB.Query("SELECT * FROM users")
    for _, user := range users {
        targetDB.Create(&AuthUser{
            ID: user.ID,
            Email: user.Email,
            // ... map fields
        })
    }
    return nil
}
```

#### День 2: Тестовая миграция
```bash
# Backup
pg_dump auth_db > auth_db_backup.sql

# Run migration on test data
go run cmd/migrate/main.go --dry-run

# Verify
psql auth_db -c "SELECT COUNT(*) FROM auth.users"
```

#### День 3: Production миграция
```bash
# Stop writes to old auth
kubectl scale deployment backend --replicas=0

# Run migration
go run cmd/migrate/main.go --production

# Verify data
./scripts/verify-migration.sh

# Switch traffic
kubectl apply -f deployments/kubernetes/
```

#### День 4-5: Мониторинг и откат
```bash
# Monitor metrics
curl http://auth-service:9090/metrics

# Check logs
kubectl logs -f deployment/auth-service

# Rollback if needed
kubectl rollout undo deployment/auth-service
```

### Фаза 5: Оптимизация (Неделя 5)

#### Производительность
- [ ] Query optimization
- [ ] Connection pooling
- [ ] Caching strategy
- [ ] Load balancing

#### Безопасность
- [ ] Security audit
- [ ] Penetration testing
- [ ] OWASP compliance
- [ ] SSL/TLS configuration

#### Мониторинг
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Alerts configuration
- [ ] Log aggregation

---

## 🧪 Тестирование

### Стратегия тестирования

#### Unit Tests (Coverage > 80%)
```go
// auth_service_test.go
func TestLogin_ValidCredentials(t *testing.T) {
    service := NewAuthService(mockDB, mockRedis)
    user, token, err := service.Login("test@example.com", "password")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.NotEmpty(t, token)
}

func TestLogin_InvalidCredentials(t *testing.T) {
    service := NewAuthService(mockDB, mockRedis)
    _, _, err := service.Login("test@example.com", "wrong")
    
    assert.Error(t, err)
    assert.Equal(t, ErrInvalidCredentials, err)
}
```

#### Integration Tests
```go
// integration_test.go
func TestAuthFlow(t *testing.T) {
    // Setup test environment
    db := setupTestDB()
    redis := setupTestRedis()
    app := setupApp(db, redis)
    
    // Register
    resp := httptest.NewRequest("POST", "/auth/register", registerBody)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Login
    resp = httptest.NewRequest("POST", "/auth/login", loginBody)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Validate token
    token := extractToken(resp)
    resp = httptest.NewRequest("GET", "/auth/validate", nil)
    resp.Header.Set("Authorization", "Bearer "+token)
    assert.Equal(t, 200, resp.StatusCode)
}
```

#### End-to-End Tests
```javascript
// e2e/auth.spec.js
describe('Authentication Flow', () => {
    it('should register new user', async () => {
        await page.goto('http://localhost:3000/register')
        await page.fill('#email', 'test@example.com')
        await page.fill('#password', 'SecurePass123!')
        await page.click('#submit')
        
        await expect(page).toHaveURL('/dashboard')
    })
    
    it('should login existing user', async () => {
        await page.goto('http://localhost:3000/login')
        await page.fill('#email', 'test@example.com')
        await page.fill('#password', 'SecurePass123!')
        await page.click('#submit')
        
        await expect(page).toHaveURL('/dashboard')
    })
})
```

#### Load Testing
```javascript
// k6/load-test.js
import http from 'k6/http'
import { check, sleep } from 'k6'

export let options = {
    stages: [
        { duration: '2m', target: 100 },  // Ramp up
        { duration: '5m', target: 1000 }, // Stay at 1000
        { duration: '2m', target: 0 },    // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95% requests under 200ms
        http_req_failed: ['rate<0.1'],    // Error rate under 10%
    },
}

export default function() {
    // Login test
    let loginRes = http.post('http://auth-service:8080/api/v1/auth/login', {
        email: 'test@example.com',
        password: 'password'
    })
    
    check(loginRes, {
        'login successful': (r) => r.status === 200,
        'token received': (r) => r.json('access_token') !== '',
    })
    
    // Token validation test
    let token = loginRes.json('access_token')
    let validateRes = http.get('http://auth-service:8080/api/v1/auth/validate', {
        headers: { 'Authorization': `Bearer ${token}` }
    })
    
    check(validateRes, {
        'token valid': (r) => r.status === 200,
    })
    
    sleep(1)
}
```

#### Security Testing
```bash
# OWASP ZAP scan
docker run -t owasp/zap2docker-stable zap-baseline.py \
    -t http://auth-service:8080 \
    -r security-report.html

# SQLMap for SQL injection
sqlmap -u "http://auth-service:8080/api/v1/auth/login" \
    --data='{"email":"*","password":"*"}' \
    --level=5 --risk=3

# Nikto web scanner
nikto -h http://auth-service:8080

# JWT security
python jwt_tool.py -t http://auth-service:8080 \
    -rc "Authorization: Bearer $TOKEN" \
    -M at
```

### Test Data

```sql
-- Test users
INSERT INTO auth.users (email, name, password_hash, provider, email_verified) VALUES
('admin@test.com', 'Admin User', '$2a$12$...', 'local', true),
('user@test.com', 'Test User', '$2a$12$...', 'local', true),
('google@test.com', 'Google User', NULL, 'google', true);

-- Test tokens
INSERT INTO auth.refresh_tokens (user_id, token_hash, family_id, expires_at) VALUES
(1, 'hash1', 'family1', NOW() + INTERVAL '30 days'),
(2, 'hash2', 'family2', NOW() + INTERVAL '30 days');

-- Test service accounts
INSERT INTO auth.service_accounts (name, api_key_hash, roles) VALUES
('backend-service', '$2a$12$...', '{"service"}'),
('monitoring-service', '$2a$12$...', '{"readonly"}');
```

---

## 📊 Мониторинг и метрики

### Prometheus Metrics

```go
// Custom metrics
var (
    loginAttempts = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "auth_login_attempts_total",
            Help: "Total number of login attempts",
        },
        []string{"status", "provider"},
    )
    
    tokenValidations = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "auth_token_validation_duration_seconds",
            Help: "Token validation duration",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1},
        },
        []string{"token_type"},
    )
    
    activeUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "auth_active_users",
            Help: "Number of active users",
        },
    )
    
    refreshTokenRotations = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "auth_refresh_token_rotations_total",
            Help: "Total number of refresh token rotations",
        },
    )
)
```

### Grafana Dashboards

```json
{
  "dashboard": {
    "title": "Auth Service Monitoring",
    "panels": [
      {
        "title": "Login Rate",
        "targets": [
          {
            "expr": "rate(auth_login_attempts_total[5m])"
          }
        ]
      },
      {
        "title": "Token Validation Latency",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, auth_token_validation_duration_seconds)"
          }
        ]
      },
      {
        "title": "Active Users",
        "targets": [
          {
            "expr": "auth_active_users"
          }
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(auth_errors_total[5m])"
          }
        ]
      }
    ]
  }
}
```

### Alerts

```yaml
# prometheus-alerts.yaml
groups:
  - name: auth-service
    rules:
    - alert: HighErrorRate
      expr: rate(auth_errors_total[5m]) > 0.05
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: High error rate in auth service
        description: "Error rate is {{ $value }} errors per second"
    
    - alert: HighLoginFailureRate
      expr: rate(auth_login_attempts_total{status="failed"}[5m]) > 10
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: High login failure rate
        description: "Login failure rate is {{ $value }} per second"
    
    - alert: DatabaseConnectionPoolExhausted
      expr: auth_db_connections_in_use / auth_db_connections_max > 0.9
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: Database connection pool almost exhausted
        description: "{{ $value }}% of connections are in use"
```

### Logging

```go
// Structured logging with context
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Str("ip", ipAddress).
    Dur("duration", duration).
    Msg("User logged in successfully")

// Audit logging
auditLogger.Log(AuditEntry{
    UserID:     userID,
    Action:     "password_changed",
    Resource:   "user_account",
    ResourceID: userID,
    IPAddress:  ipAddress,
    UserAgent:  userAgent,
    Success:    true,
    Metadata: map[string]interface{}{
        "old_password_hash": oldHash,
        "new_password_hash": newHash,
    },
})
```

---

## 🚀 Развертывание

### Development Environment

```bash
# Clone repository
git clone https://github.com/svetu/auth-service.git
cd auth-service

# Setup environment
cp .env.example .env
vim .env  # Configure settings

# Generate JWT keys
./scripts/generate-keys.sh

# Start services
docker-compose up -d

# Run migrations
make migrate-up

# Seed test data
make seed

# Run tests
make test

# Start development server
make dev
```

### Staging Environment

```bash
# Build Docker image
docker build -t svetu/auth-service:staging .

# Push to registry
docker push svetu/auth-service:staging

# Deploy to Kubernetes
kubectl apply -f deployments/kubernetes/staging/

# Run smoke tests
./scripts/smoke-test.sh staging
```

### Production Deployment

```bash
# Pre-deployment checks
./scripts/pre-deploy-check.sh

# Create production image
docker build -t svetu/auth-service:v1.0.0 .
docker tag svetu/auth-service:v1.0.0 svetu/auth-service:latest

# Push to registry
docker push svetu/auth-service:v1.0.0
docker push svetu/auth-service:latest

# Blue-Green deployment
kubectl apply -f deployments/kubernetes/production/blue/
./scripts/verify-deployment.sh blue

# Switch traffic
kubectl patch service auth-service -p '{"spec":{"selector":{"version":"blue"}}}'

# Monitor
./scripts/monitor-deployment.sh

# Rollback if needed
kubectl patch service auth-service -p '{"spec":{"selector":{"version":"green"}}}'
```

### CI/CD Pipeline

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    
    - name: Run tests
      run: |
        go test -v -race -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
  
  security:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Run Gosec
      uses: securego/gosec@master
      with:
        args: ./...
    
    - name: Run Trivy
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'
        scan-ref: '.'
  
  build:
    needs: [test, security]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v3
    
    - name: Build Docker image
      run: docker build -t svetu/auth-service:${{ github.sha }} .
    
    - name: Push to registry
      run: |
        echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
        docker push svetu/auth-service:${{ github.sha }}
  
  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - name: Deploy to Kubernetes
      run: |
        kubectl set image deployment/auth-service auth-service=svetu/auth-service:${{ github.sha }}
        kubectl rollout status deployment/auth-service
```

---

## 📚 Документация API

### OpenAPI Specification

```yaml
openapi: 3.0.0
info:
  title: Auth Service API
  version: 1.0.0
  description: Authentication and Authorization Service
  contact:
    name: API Support
    email: api@svetu.rs
  license:
    name: MIT

servers:
  - url: https://auth.svetu.rs
    description: Production
  - url: https://auth-staging.svetu.rs
    description: Staging
  - url: http://localhost:8080
    description: Development

security:
  - BearerAuth: []

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        email:
          type: string
          format: email
        name:
          type: string
        picture_url:
          type: string
          format: uri
        email_verified:
          type: boolean
        created_at:
          type: string
          format: date-time
    
    AuthResponse:
      type: object
      properties:
        access_token:
          type: string
        refresh_token:
          type: string
        token_type:
          type: string
          enum: [Bearer]
        expires_in:
          type: integer
        user:
          $ref: '#/components/schemas/User'
    
    Error:
      type: object
      properties:
        error:
          type: string
        message:
          type: string
        details:
          type: object

paths:
  /api/v1/auth/register:
    post:
      summary: Register new user
      operationId: register
      tags: [Authentication]
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [email, password, name]
              properties:
                email:
                  type: string
                  format: email
                password:
                  type: string
                  minLength: 8
                name:
                  type: string
      responses:
        '200':
          description: Registration successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '400':
          description: Validation error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '409':
          description: User already exists
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
```

### Client SDKs

#### Go SDK
```go
package authclient

import (
    "context"
    "fmt"
    
    "google.golang.org/grpc"
    pb "github.com/svetu/auth-service/pkg/proto"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.AuthServiceClient
}

func NewClient(address string) (*Client, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    return &Client{
        conn:   conn,
        client: pb.NewAuthServiceClient(conn),
    }, nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (*pb.ValidateTokenResponse, error) {
    return c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
        Token: token,
    })
}

func (c *Client) Close() error {
    return c.conn.Close()
}
```

#### TypeScript SDK
```typescript
// auth-client-sdk/src/index.ts
export class AuthClient {
  private readonly baseURL: string
  
  constructor(baseURL: string = 'https://auth.svetu.rs') {
    this.baseURL = baseURL
  }
  
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await fetch(`${this.baseURL}/api/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
      credentials: 'include',
    })
    
    if (!response.ok) {
      throw new AuthError(await response.json())
    }
    
    return response.json()
  }
  
  async validateToken(token: string): Promise<ValidationResponse> {
    const response = await fetch(`${this.baseURL}/api/v1/auth/validate`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    
    return response.json()
  }
  
  async refreshToken(): Promise<AuthResponse> {
    const response = await fetch(`${this.baseURL}/api/v1/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    })
    
    if (!response.ok) {
      throw new AuthError(await response.json())
    }
    
    return response.json()
  }
}
```

---

## ⚠️ Риски и митигация

### Технические риски

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| Утечка данных при миграции | Низкая | Критическое | Шифрование, безопасные каналы, аудит |
| Недоступность сервиса | Средняя | Высокое | Blue-green deployment, rollback план |
| Производительность ниже ожидаемой | Средняя | Среднее | Load testing, оптимизация, кеширование |
| Несовместимость с существующими клиентами | Низкая | Высокое | Backward compatibility layer |
| Ошибки в токен ротации | Низкая | Среднее | Extensive testing, monitoring |

### Бизнес риски

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| Потеря пользователей из-за проблем входа | Низкая | Критическое | Graceful degradation, fallback |
| Увеличение затрат на инфраструктуру | Средняя | Среднее | Оптимизация ресурсов, автоскейлинг |
| Задержка релиза | Средняя | Среднее | Agile подход, MVP first |
| Compliance нарушения | Низкая | Высокое | Security audit, legal review |

### План отката

```bash
# 1. Остановка нового сервиса
kubectl scale deployment auth-service --replicas=0

# 2. Переключение трафика на старый сервис
kubectl patch service backend -p '{"spec":{"selector":{"version":"legacy"}}}'

# 3. Восстановление базы данных
pg_restore -h localhost -d svetubd backup_before_migration.sql

# 4. Очистка Redis
redis-cli FLUSHALL

# 5. Перезапуск старого сервиса
kubectl rollout restart deployment backend-legacy

# 6. Проверка
./scripts/health-check.sh legacy
```

---

## 📎 Приложения

### Приложение A: Глоссарий

| Термин | Описание |
|--------|----------|
| JWT | JSON Web Token - стандарт для безопасной передачи информации |
| Refresh Token | Долгоживущий токен для обновления access token |
| Token Rotation | Процесс замены refresh token при каждом использовании |
| Token Family | Группа связанных токенов для обнаружения повторного использования |
| Device Fingerprinting | Идентификация устройства по уникальным характеристикам |
| Rate Limiting | Ограничение количества запросов за период времени |
| OWASP | Open Web Application Security Project |
| MFA | Multi-Factor Authentication |
| SSO | Single Sign-On |
| RBAC | Role-Based Access Control |

### Приложение B: Конфигурационные параметры

```env
# Database
DB_HOST=localhost
DB_PORT=5433
DB_NAME=auth_db
DB_USER=auth_user
DB_PASSWORD=secure_password
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=100

# Redis
REDIS_HOST=localhost
REDIS_PORT=6380
REDIS_PASSWORD=redis_password
REDIS_DB=0

# Server
HTTP_PORT=8080
GRPC_PORT=50051
METRICS_PORT=9090

# JWT
JWT_ALGORITHM=RS256
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=30d
JWT_PRIVATE_KEY_PATH=/etc/auth/keys/private.pem
JWT_PUBLIC_KEY_PATH=/etc/auth/keys/public.pem

# OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret

# Security
BCRYPT_COST=12
RATE_LIMIT_LOGIN=5/15m
RATE_LIMIT_REGISTER=3/1h

# Monitoring
TRACING_ENABLED=true
METRICS_ENABLED=true

# Environment
ENV=production
LOG_LEVEL=info
```

### Приложение C: Команды управления

```bash
# Health check
curl http://auth-service:8081/health

# Metrics
curl http://auth-service:9090/metrics

# Version info
curl http://auth-service:8080/api/v1/version

# Database migrations
make migrate-up    # Apply migrations
make migrate-down  # Rollback migrations
make migrate-status # Check status

# Redis operations
redis-cli -h localhost -p 6380 --auth password INFO
redis-cli -h localhost -p 6380 --auth password FLUSHDB

# Logs
docker logs auth_service -f --tail 100
kubectl logs deployment/auth-service -f

# Backup
pg_dump -h localhost -p 5433 -U auth_user auth_db > backup.sql
redis-cli -h localhost -p 6380 --auth password BGSAVE

# Performance testing
k6 run scripts/load-test.js
ab -n 10000 -c 100 http://auth-service:8080/api/v1/auth/validate
```

---

## 📝 Заключение

Auth Service представляет собой критически важный компонент инфраструктуры Sve Tu, обеспечивающий безопасную и масштабируемую аутентификацию и авторизацию. Переход на микросервисную архитектуру позволит:

1. **Повысить безопасность** через изоляцию и специализированные механизмы защиты
2. **Улучшить производительность** на 200% благодаря оптимизации и кешированию
3. **Обеспечить масштабируемость** для поддержки роста пользовательской базы
4. **Упростить разработку** через четкое разделение ответственности
5. **Снизить технический долг** удалением legacy кода

Реализация проекта займет 5 недель с командой из 3-4 разработчиков. ROI ожидается в течение 6 месяцев за счет снижения операционных расходов и повышения надежности системы.

---

**Документ подготовлен:** Architecture Team  
**Версия:** 1.0.0  
**Дата:** 06.09.2025  
**Статус:** Готов к реализации

---

*Этот документ является живым и будет обновляться по мере развития проекта.*