# Шаг 3: Создание .env.example

## Цель
Создать файл с примерами переменных окружения для документирования конфигурации и упрощения развертывания.

## Задачи

### 3.1 Создание основного .env.example

Файл: `/frontend/svetu/.env.example`

```bash
# ==========================================
# Frontend Environment Variables
# ==========================================
# Copy this file to .env.local for local development
# All NEXT_PUBLIC_* variables are exposed to the browser

# API Configuration
# ==========================================
# Public API URL used by the browser
NEXT_PUBLIC_API_URL=http://localhost:3000

# Internal API URL for SSR (Docker/Kubernetes only)
# Used for server-side requests to avoid going through public network
INTERNAL_API_URL=http://backend:3000

# Storage Configuration (MinIO/S3)
# ==========================================
# Public MinIO URL for accessing images
NEXT_PUBLIC_MINIO_URL=http://localhost:9000

# Comma-separated list of allowed image hosts
# Format: protocol:hostname:port (port is optional for 80/443)
NEXT_PUBLIC_IMAGE_HOSTS=http:localhost:9000,https:svetu.rs:443,http:localhost:3000

# Pattern for image paths (glob pattern)
NEXT_PUBLIC_IMAGE_PATH_PATTERN=/listings/**

# WebSocket Configuration
# ==========================================
# WebSocket URL for real-time features (chat, notifications)
# Leave empty to disable real-time features
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:3000

# Feature Flags
# ==========================================
# Enable/disable specific features
NEXT_PUBLIC_ENABLE_PAYMENTS=false
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_DEBUG=false

# Third-party Services
# ==========================================
# Google OAuth (configured in Google Cloud Console)
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com

# Sentry error tracking (optional)
NEXT_PUBLIC_SENTRY_DSN=

# Environment
# ==========================================
# DO NOT CHANGE IN PRODUCTION
NODE_ENV=development

# Additional Settings
# ==========================================
# Disable Next.js telemetry
NEXT_TELEMETRY_DISABLED=1

# Port for development server
PORT=3001
```

### 3.2 Создание .env.production.example

Файл: `/frontend/svetu/.env.production.example`

```bash
# ==========================================
# Production Environment Variables
# ==========================================
# These values should be set in your deployment environment
# DO NOT commit actual production values to git

# API Configuration
# ==========================================
NEXT_PUBLIC_API_URL=https://api.svetu.rs
INTERNAL_API_URL=http://api-service:3000

# Storage Configuration
# ==========================================
NEXT_PUBLIC_MINIO_URL=https://svetu.rs
NEXT_PUBLIC_IMAGE_HOSTS=https:svetu.rs:443
NEXT_PUBLIC_IMAGE_PATH_PATTERN=/listings/**

# WebSocket Configuration
# ==========================================
NEXT_PUBLIC_WEBSOCKET_URL=wss://api.svetu.rs

# Feature Flags
# ==========================================
NEXT_PUBLIC_ENABLE_PAYMENTS=true
NEXT_PUBLIC_ENABLE_ANALYTICS=true
NEXT_PUBLIC_ENABLE_DEBUG=false

# Third-party Services
# ==========================================
NEXT_PUBLIC_GOOGLE_CLIENT_ID=production-client-id.apps.googleusercontent.com
NEXT_PUBLIC_SENTRY_DSN=https://your-sentry-dsn@sentry.io/project-id

# Environment
# ==========================================
NODE_ENV=production
NEXT_TELEMETRY_DISABLED=1
```

### 3.3 Создание документации по переменным

Файл: `/frontend/svetu/docs/ENVIRONMENT.md`

```markdown
# Environment Variables Documentation

## Overview
This document describes all environment variables used by the frontend application.

## Variable Categories

### 1. Public Variables (NEXT_PUBLIC_*)
These variables are exposed to the browser and can be accessed in client-side code.

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | Yes | `http://localhost:3000` | Public API endpoint URL |
| `NEXT_PUBLIC_MINIO_URL` | Yes | `http://localhost:9000` | MinIO/S3 storage URL |
| `NEXT_PUBLIC_IMAGE_HOSTS` | No | See .env.example | Allowed image host domains |
| `NEXT_PUBLIC_IMAGE_PATH_PATTERN` | No | `/listings/**` | Valid image path patterns |
| `NEXT_PUBLIC_WEBSOCKET_URL` | No | - | WebSocket endpoint for real-time features |
| `NEXT_PUBLIC_ENABLE_PAYMENTS` | No | `false` | Enable payment features |
| `NEXT_PUBLIC_GOOGLE_CLIENT_ID` | No | - | Google OAuth client ID |

### 2. Server-only Variables
These variables are only available in server-side code (API routes, SSR).

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `INTERNAL_API_URL` | No | - | Internal API URL for Docker/K8s |
| `NODE_ENV` | Yes | `development` | Node environment |
| `PORT` | No | `3000` | Server port |

### 3. Build-time Variables
These affect the build process and optimization.

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEXT_TELEMETRY_DISABLED` | No | `0` | Disable Next.js telemetry |

## Environment-specific Configurations

### Development
```bash
NEXT_PUBLIC_API_URL=http://localhost:3000
NEXT_PUBLIC_ENABLE_DEBUG=true
NODE_ENV=development
```

### Staging
```bash
NEXT_PUBLIC_API_URL=https://staging-api.svetu.rs
NEXT_PUBLIC_ENABLE_DEBUG=true
NODE_ENV=production
```

### Production
```bash
NEXT_PUBLIC_API_URL=https://api.svetu.rs
NEXT_PUBLIC_ENABLE_DEBUG=false
NODE_ENV=production
```

## Security Considerations

1. **Never commit real values** for production environment variables
2. **Use secrets management** for sensitive values (API keys, tokens)
3. **Validate all variables** at runtime to catch configuration errors early
4. **Minimize public variables** to reduce attack surface

## Docker/Kubernetes Setup

### Docker Compose
```yaml
services:
  frontend:
    environment:
      - NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
      - INTERNAL_API_URL=http://backend:3000
```

### Kubernetes ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: frontend-config
data:
  NEXT_PUBLIC_API_URL: "https://api.svetu.rs"
  NEXT_PUBLIC_ENABLE_PAYMENTS: "true"
```

## Troubleshooting

### Variable not updating
1. Restart the development server
2. Clear Next.js cache: `rm -rf .next`
3. Check variable name starts with `NEXT_PUBLIC_`

### Variable undefined in production
1. Ensure variable is set in deployment environment
2. Check Docker/K8s configuration
3. Verify build logs for warnings

### Type errors with env variables
1. Update type definitions in `config/types.ts`
2. Run validation to catch errors early
3. Use default values for optional variables
```

### 3.4 Добавление скрипта для проверки переменных

Файл: `/frontend/svetu/scripts/check-env.js`

```javascript
#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Required environment variables
const requiredVars = [
  'NEXT_PUBLIC_API_URL',
  'NEXT_PUBLIC_MINIO_URL',
];

// Optional but recommended
const recommendedVars = [
  'NEXT_PUBLIC_WEBSOCKET_URL',
  'NEXT_PUBLIC_GOOGLE_CLIENT_ID',
];

console.log('🔍 Checking environment variables...\n');

let hasErrors = false;
let hasWarnings = false;

// Check required variables
console.log('Required variables:');
requiredVars.forEach(varName => {
  if (process.env[varName]) {
    console.log(`✅ ${varName}: ${process.env[varName]}`);
  } else {
    console.log(`❌ ${varName}: NOT SET`);
    hasErrors = true;
  }
});

console.log('\nRecommended variables:');
recommendedVars.forEach(varName => {
  if (process.env[varName]) {
    console.log(`✅ ${varName}: ${process.env[varName]}`);
  } else {
    console.log(`⚠️  ${varName}: not set (optional)`);
    hasWarnings = true;
  }
});

// Check for .env.local file
const envLocalPath = path.join(process.cwd(), '.env.local');
if (!fs.existsSync(envLocalPath)) {
  console.log('\n⚠️  No .env.local file found. Using defaults or system environment.');
  hasWarnings = true;
}

// Summary
console.log('\n' + '='.repeat(50));
if (hasErrors) {
  console.log('❌ Environment check failed! Missing required variables.');
  process.exit(1);
} else if (hasWarnings) {
  console.log('⚠️  Environment check passed with warnings.');
} else {
  console.log('✅ Environment check passed!');
}
```

### 3.5 Обновление package.json

Добавить в `/frontend/svetu/package.json`:

```json
{
  "scripts": {
    "env:check": "node scripts/check-env.js",
    "env:create": "cp .env.example .env.local",
    "predev": "npm run env:check",
    "prebuild": "npm run env:check"
  }
}
```

### 3.6 Обновление .gitignore

Убедиться что в `/frontend/svetu/.gitignore` есть:

```gitignore
# Environment files
.env
.env.local
.env.production
.env.*.local

# Keep example files
!.env.example
!.env.*.example
```

## Использование

### Для новых разработчиков
```bash
# Клонировать репозиторий
git clone <repo>
cd frontend/svetu

# Создать локальный env файл из примера
yarn env:create

# Отредактировать .env.local с нужными значениями
nano .env.local

# Проверить конфигурацию
yarn env:check

# Запустить приложение
yarn dev
```

### Для CI/CD
```bash
# В GitHub Actions
env:
  NEXT_PUBLIC_API_URL: ${{ secrets.API_URL }}
  NEXT_PUBLIC_MINIO_URL: ${{ secrets.MINIO_URL }}

# В Docker
docker run -e NEXT_PUBLIC_API_URL=https://api.svetu.rs myapp

# В Kubernetes
kubectl create configmap frontend-env --from-env-file=.env.production
```

## Результат
После выполнения этого шага:
1. Будет создана документация по всем переменным окружения
2. Новые разработчики смогут быстро настроить окружение
3. CI/CD процессы будут иметь четкую конфигурацию
4. Появится автоматическая проверка переменных перед запуском

## Критерии приемки

### 1. Файлы созданы
- [x] `/frontend/svetu/.env.example` создан
- [x] `/frontend/svetu/.env.production.example` создан
- [x] `/frontend/svetu/docs/ENVIRONMENT.md` создан
- [x] Все файлы содержат актуальную информацию

### 2. Содержимое .env.example
- [x] Все используемые в приложении переменные документированы
- [x] Каждая переменная имеет комментарий с описанием
- [x] Примеры значений корректны и работоспособны
- [x] Переменные сгруппированы по категориям (API, Storage, Features и т.д.)

### 3. Скрипт проверки переменных
- [x] Файл `/frontend/svetu/scripts/check-env.js` создан
- [x] Скрипт исполняемый (`chmod +x scripts/check-env.js`)
- [x] Скрипт проверяет обязательные переменные (NEXT_PUBLIC_API_URL, NEXT_PUBLIC_MINIO_URL)
- [x] Скрипт возвращает exit code 1 при отсутствии обязательных переменных
- [x] Скрипт выводит понятные сообщения об ошибках и предупреждениях

### 4. Интеграция в package.json
- [x] Команда `env:check` добавлена и работает
- [x] Команда `env:create` создает .env.local из .env.example
- [x] Hook `predev` запускает проверку перед `yarn dev`
- [x] Hook `prebuild` запускает проверку перед `yarn build`
- [x] При отсутствии обязательных переменных dev/build не запускается

### 5. Git конфигурация
- [x] `.gitignore` обновлен
- [x] `.env.local` не коммитится (проверить `git status`)
- [x] `.env.example` файлы коммитятся
- [x] Нет случайно закоммиченных env файлов с реальными значениями

### 6. Использование и документация
- [x] Новый разработчик может настроить окружение следуя инструкциям:
  ```bash
  yarn env:create
  # Редактирование .env.local
  yarn env:check  # Проверка проходит
  yarn dev        # Приложение запускается
  ```
- [x] Документация ENVIRONMENT.md содержит:
  - Таблицы с описанием всех переменных
  - Примеры для разных окружений (dev, staging, prod)
  - Инструкции по настройке Docker/Kubernetes
  - Раздел по устранению неполадок

### 7. Безопасность
- [x] Production значения не указаны в примерах
- [x] Есть предупреждения о недопустимости коммита реальных значений
- [x] Sensitive переменные помечены соответствующим образом