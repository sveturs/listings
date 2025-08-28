# Environment Variables Documentation for SveTu Dev

## Backend Environment Variables

### Core Settings
- **APP_MODE**: production
- **ENV_FILE**: .env (not used in container, variables passed directly)
- **WS_ENABLED**: "true" - WebSocket support
- **MIGRATIONS_ON_API**: full - Run migrations and fixtures on startup

### Authentication
- **JWT_SECRET**: Secret key for JWT tokens
- **GOOGLE_CLIENT_ID**: Google OAuth client ID
- **GOOGLE_CLIENT_SECRET**: Google OAuth client secret  
- **GOOGLE_OAUTH_REDIRECT_URL**: OAuth callback URL

### URLs
- **FRONTEND_URL**: https://dev.svetu.rs
- **BACKEND_URL**: https://devapi.svetu.rs
- **FILE_STORAGE_PUBLIC_URL**: Public URL for file access

### Database
- **DATABASE_URL**: PostgreSQL connection string
  - Default: postgres://postgres:password@db:5432/svetu_db?sslmode=disable

### Redis
- **REDIS_URL**: redis://redis:6379
- **REDIS_HOST**: redis
- **REDIS_PORT**: 6379

### OpenSearch
- **OPENSEARCH_URL**: http://opensearch:9200
- **OPENSEARCH_MARKETPLACE_INDEX**: marketplace

### MinIO (S3 Storage)
- **FILE_STORAGE_PROVIDER**: minio
- **MINIO_ENDPOINT**: minio:9000
- **MINIO_ACCESS_KEY**: minioadmin
- **MINIO_SECRET_KEY**: 1321321321321
- **MINIO_USE_SSL**: "false"
- **MINIO_BUCKET_NAME**: listings
- **MINIO_LOCATION**: eu-central-1

### AI Integration
- **OPENAI_API_KEY**: OpenAI API key (use real key for AI features)

## Frontend Environment Variables

### Core
- **NODE_ENV**: production

### API Configuration
- **NEXT_PUBLIC_API_URL**: https://devapi.svetu.rs
- **NEXT_PUBLIC_WEBSOCKET_URL**: wss://devapi.svetu.rs

### Storage
- **NEXT_PUBLIC_MINIO_URL**: http://localhost:9002
- **NEXT_PUBLIC_IMAGE_HOSTS**: Allowed image hosts
- **NEXT_PUBLIC_IMAGE_PATH_PATTERN**: /listings/**

### Features
- **NEXT_PUBLIC_ENABLE_PAYMENTS**: false

## How to Manage Variables

### 1. Using .env file
Create or edit :
```bash
cd /opt/svetu-dev
nano .env
```

Variables in .env are used as defaults if not set in docker-compose.dev.yml

### 2. Updating docker-compose.dev.yml
Edit the environment section for each service:
```bash
cd /opt/svetu-dev
nano docker-compose.dev.yml
```

### 3. Apply changes
After changing variables:
```bash
# Restart specific service
docker-compose -f docker-compose.dev.yml restart backend
docker-compose -f docker-compose.dev.yml restart frontend

# Or recreate with new variables
docker-compose -f docker-compose.dev.yml up -d --force-recreate backend frontend
```

### 4. Check current variables
```bash
# Check backend variables
docker exec svetu-dev_backend_1 env | sort

# Check frontend variables  
docker exec svetu-dev_frontend_1 env | grep NEXT_PUBLIC
```

## Important Notes

1. **Migrations**: Backend automatically runs migrations and fixtures on startup when MIGRATIONS_ON_API=full

2. **SSL/HTTPS**: Dev environment uses HTTPS via nginx proxy:
   - Frontend: https://dev.svetu.rs
   - Backend: https://devapi.svetu.rs

3. **Dummy values**: Many API keys use dummy values for dev. Replace with real keys for full functionality:
   - OPENAI_API_KEY
   - GOOGLE_CLIENT_ID/SECRET

4. **Ports**:
   - Backend: 3002 (host) -> 3000 (container)
   - Frontend: 3003 (host) -> 3000 (container)
   - PostgreSQL: 5433
   - MinIO: 9002 (API), 9003 (Console)
   - OpenSearch: 9201
   - Redis: 6380

## Reset Everything

To completely reset the dev environment with fresh database:
```bash
cd /opt/svetu-dev
./reset-dev.sh
```

This will:
- Stop all containers
- Delete all volumes and data
- Recreate everything from scratch
- Apply migrations and fixtures
- Start all services
