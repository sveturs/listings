# Лучшие практики создания универсального Docker-образа Next.js с runtime конфигурацией

Next.js традиционно фиксирует переменные окружения на этапе сборки, что противоречит принципу "build once, deploy many". Для создания универсальных Docker-образов требуется использование специальных техник runtime конфигурации - от инъекции переменных через API routes до использования библиотек вроде `next-runtime-env`. Правильная архитектура позволяет одному образу работать в разных средах, настраивая API endpoints динамически: серверная часть обращается к внутренним Docker-сетям, а клиентская получает публичные URL через браузер.

## Современные подходы к runtime configuration

### Рекомендуемое решение: next-runtime-env

Наиболее элегантный подход для Next.js 13+ - использование библиотеки `next-runtime-env`, которая решает проблему runtime конфигурации без дополнительных HTTP запросов:

```typescript
// app/layout.tsx
import { PublicEnvScript } from 'next-runtime-env';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html>
      <head>
        <PublicEnvScript />
      </head>
      <body>{children}</body>
    </html>
  );
}

// app/client-component.tsx
'use client';
import { env } from 'next-runtime-env';

export default function ClientComponent() {
  const apiUrl = env('NEXT_PUBLIC_API_URL');
  return <div>API: {apiUrl}</div>;
}
```

### Server-side injection для App Router

Альтернативный подход через серверные компоненты Next.js 13+:

```typescript
// app/config-provider.tsx
import { unstable_noStore as noStore } from 'next/cache';

export default function ConfigProvider({ children }: { children: React.ReactNode }) {
  noStore(); // Обеспечивает динамический рендеринг
  
  const config = {
    apiUrl: process.env.NEXT_PUBLIC_API_URL,
    appEnv: process.env.NODE_ENV,
  };
  
  return (
    <>
      <script
        dangerouslySetInnerHTML={{
          __html: `window.__CONFIG__ = ${JSON.stringify(config)};`,
        }}
      />
      {children}
    </>
  );
}
```

### API Route для sensitive конфигурации

Для безопасной передачи конфигурации без экспонирования в HTML:

```typescript
// app/api/config/route.ts
import { NextResponse } from 'next/server';

export const dynamic = 'force-dynamic';

export async function GET() {
  const config = {
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'https://api.example.com',
    appEnv: process.env.NODE_ENV || 'production',
    // НЕ включаем секретные переменные
  };
  
  return NextResponse.json(config);
}
```

## Паттерны разделения серверных и клиентских переменных

### Архитектурное разделение API endpoints

```javascript
// lib/api-config.ts
export function getApiConfig(isServer: boolean = typeof window === 'undefined') {
  if (isServer) {
    // Серверная конфигурация - internal Docker URLs
    return {
      baseURL: process.env.INTERNAL_API_URL || 'http://api:3001',
      timeout: 30000,
    };
  } else {
    // Клиентская конфигурация - публичные URLs
    return {
      baseURL: process.env.NEXT_PUBLIC_API_URL || 'https://api.example.com',
      timeout: 5000,
    };
  }
}
```

### Docker Compose с разными URL для сервера и клиента

```yaml
version: '3.8'
services:
  nextjs:
    environment:
      - INTERNAL_API_URL=http://api:3001          # Внутренний URL для SSR
      - NEXT_PUBLIC_API_URL=https://api.example.com  # Публичный URL для браузера
  api:
    hostname: api
    ports:
      - "3001:3001"
```

## Production-ready Dockerfile с оптимизациями

### Multi-stage build с security best practices

```dockerfile
# syntax=docker.io/docker/dockerfile:1

ARG NODE_VERSION=22
ARG ALPINE_VERSION=3.19

# Base image
FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk add --no-cache libc6-compat dumb-init

# Dependencies stage
FROM base AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev && npm cache clean --force

# Builder stage
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# Placeholder для runtime замены
ENV NEXT_PUBLIC_API_URL="__NEXT_PUBLIC_API_URL__"

RUN npm run build

# Production stage
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Security: non-root user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Копирование с правильными permissions
COPY --from=builder --chown=nextjs:nodejs /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Entrypoint script для runtime config
COPY --chown=nextjs:nodejs entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

USER nextjs
EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

ENTRYPOINT ["dumb-init", "--", "entrypoint.sh"]
CMD ["node", "server.js"]
```

### Entrypoint script для инъекции переменных

```bash
#!/bin/sh
# entrypoint.sh

set -e

# Runtime замена переменных в статических файлах
replace_env_vars() {
  echo "🔄 Injecting environment variables..."
  
  find /app/.next/static -name "*.js" -type f -exec \
    sed -i "s|__NEXT_PUBLIC_API_URL__|${NEXT_PUBLIC_API_URL:-https://api.example.com}|g" {} \;
  
  echo "✅ Environment variables injected"
}

# Graceful shutdown
graceful_shutdown() {
  echo "🛑 Starting graceful shutdown..."
  kill -TERM "$MAIN_PID" 2>/dev/null || true
  wait "$MAIN_PID" 2>/dev/null || true
  exit 0
}

trap graceful_shutdown SIGTERM SIGINT

# Инъекция переменных
replace_env_vars

# Запуск приложения
echo "🚀 Starting Next.js application..."
node server.js &
MAIN_PID=$!

wait $MAIN_PID
```

## Интеграция с Docker Compose и Kubernetes

### Docker Compose для multiple environments

```yaml
# docker-compose.yml (базовая)
version: '3.8'
services:
  nextjs:
    build: 
      context: .
      dockerfile: Dockerfile
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
    networks:
      - app-network

networks:
  app-network:
    driver: bridge

# docker-compose.prod.yml (production overrides)
version: '3.8'
services:
  nextjs:
    environment:
      - NODE_ENV=production
      - INTERNAL_API_URL=http://api:3001
      - NEXT_PUBLIC_API_URL=https://api.production.com
    restart: unless-stopped
    deploy:
      replicas: 2
      resources:
        limits:
          memory: 512M
```

### Kubernetes с ConfigMaps и Secrets

```yaml
# ConfigMap для публичной конфигурации
apiVersion: v1
kind: ConfigMap
metadata:
  name: nextjs-config
data:
  NODE_ENV: "production"
  NEXT_PUBLIC_API_URL: "https://api.example.com"

---
# Secret для sensitive данных
apiVersion: v1
kind: Secret
metadata:
  name: nextjs-secrets
type: Opaque
data:
  DATABASE_URL: cG9zdGdyZXNxbDovL3VzZXI6cGFzc0Bkyi8= # base64

---
# Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nextjs-app
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: nextjs
        image: myregistry/nextjs-app:latest
        envFrom:
        - configMapRef:
            name: nextjs-config
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: nextjs-secrets
              key: DATABASE_URL
        livenessProbe:
          httpGet:
            path: /api/health
            port: 3000
```

## Настройка реверс-прокси для API Gateway паттерна

### Nginx конфигурация с path rewriting

```nginx
# nginx/nginx.conf
upstream nextjs {
    server nextjs:3000;
}

server {
    listen 80;
    server_name app.example.com;

    # Кеширование статических ассетов
    location /_next/static {
        proxy_pass http://nextjs;
        proxy_cache_valid 200 60m;
        add_header Cache-Control "public, max-age=31536000, immutable";
    }

    # API routing с разделением internal/external
    location /api/internal/ {
        internal;  # Только для внутренних запросов
        proxy_pass http://nextjs;
    }

    location /api/ {
        proxy_pass http://nextjs;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Основное приложение
    location / {
        proxy_pass http://nextjs;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Безопасность при передаче конфигурации

### Защита sensitive данных

```javascript
// lib/config-validator.js
import Joi from 'joi';

const configSchema = Joi.object({
  DATABASE_URL: Joi.string().uri().required(),
  API_KEY: Joi.string().min(32).required(),
  NODE_ENV: Joi.string().valid('development', 'production', 'test'),
});

export function validateConfig() {
  const { error, value } = configSchema.validate(process.env);
  if (error) {
    throw new Error(`Config validation error: ${error.message}`);
  }
  return value;
}
```

### Secrets management интеграция

```dockerfile
# Использование Docker BuildKit secrets
FROM node:18-alpine AS builder
RUN --mount=type=secret,id=api-key,dst=/run/secrets/api-key \
    export API_KEY=$(cat /run/secrets/api-key) && \
    npm run build
```

## Стратегии кеширования и оптимизации

### Оптимальная структура слоев

```dockerfile
# Максимизация кеширования через правильный порядок команд
FROM node:22-alpine AS deps
WORKDIR /app
# Сначала копируем только package files
COPY package*.json ./
RUN npm ci --omit=dev

FROM node:22-alpine AS builder
WORKDIR /app
# Потом копируем зависимости и код
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build
```

### Next.js standalone mode

```javascript
// next.config.js
module.exports = {
  output: 'standalone', // Минимальный production bundle
  compress: true,
  poweredByHeader: false,
  generateEtags: false,
};
```

## Заключение

Создание универсального Docker-образа Next.js требует комплексного подхода к runtime конфигурации. Ключевые рекомендации:

1. **Используйте `next-runtime-env`** для элегантного решения проблемы runtime переменных
2. **Разделяйте серверные и клиентские конфигурации** через API routes и условную логику
3. **Применяйте multi-stage builds** с non-root пользователями для безопасности
4. **Настройте реверс-прокси** для правильной маршрутизации internal/external API
5. **Внедрите proper secrets management** через Docker secrets или внешние системы

Эта архитектура обеспечивает принцип "build once, deploy many", позволяя одному Docker-образу работать во всех окружениях с динамической конфигурацией на этапе запуска контейнера.