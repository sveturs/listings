# План улучшения системы конфигурации Frontend

## 1. Внедрение Runtime конфигурации с next-runtime-env

### Установка и настройка:
```bash
yarn add next-runtime-env
```

### Изменения в layout.tsx:
```typescript
// app/[locale]/layout.tsx
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
```

## 2. Улучшенная структура конфигурации

### src/config/types.ts:
```typescript
import { z } from 'zod';

// Схема валидации для публичных переменных
export const publicEnvSchema = z.object({
  NEXT_PUBLIC_API_URL: z.string().url().min(1),
  NEXT_PUBLIC_MINIO_URL: z.string().url().min(1),
  NEXT_PUBLIC_IMAGE_HOSTS: z.string().optional(),
  NEXT_PUBLIC_IMAGE_PATH_PATTERN: z.string().optional(),
  NEXT_PUBLIC_WEBSOCKET_URL: z.string().url().optional(),
});

// Схема для серверных переменных
export const serverEnvSchema = z.object({
  INTERNAL_API_URL: z.string().url().optional(),
  NODE_ENV: z.enum(['development', 'production', 'test']),
});

// Типы из схем
export type PublicEnvVariables = z.infer<typeof publicEnvSchema>;
export type ServerEnvVariables = z.infer<typeof serverEnvSchema>;

// Расширенный интерфейс конфигурации
export interface Config {
  api: {
    url: string;
    internalUrl?: string; // Для SSR запросов
    websocketUrl?: string;
  };
  storage: {
    minioUrl: string;
    imageHosts: ImageHost[];
    imagePathPattern: string;
  };
  env: {
    isProduction: boolean;
    isDevelopment: boolean;
    isServer: boolean;
  };
  features?: {
    enableChat: boolean;
    enablePayments: boolean;
  };
}
```

### src/config/index.ts с runtime поддержкой:
```typescript
import { env } from 'next-runtime-env';
import { publicEnvSchema, serverEnvSchema, Config } from './types';

class ConfigManager {
  private config: Config | null = null;
  
  private loadConfig(): Config {
    const isServer = typeof window === 'undefined';
    
    // Для клиента используем runtime env
    const getEnvValue = (key: string, defaultValue: string = ''): string => {
      if (isServer) {
        return process.env[key] || defaultValue;
      }
      return env(key) || defaultValue;
    };
    
    // Валидация публичных переменных
    const publicEnv = {
      NEXT_PUBLIC_API_URL: getEnvValue('NEXT_PUBLIC_API_URL', 'http://localhost:3000'),
      NEXT_PUBLIC_MINIO_URL: getEnvValue('NEXT_PUBLIC_MINIO_URL', 'http://localhost:9000'),
      NEXT_PUBLIC_IMAGE_HOSTS: getEnvValue('NEXT_PUBLIC_IMAGE_HOSTS'),
      NEXT_PUBLIC_IMAGE_PATH_PATTERN: getEnvValue('NEXT_PUBLIC_IMAGE_PATH_PATTERN'),
      NEXT_PUBLIC_WEBSOCKET_URL: getEnvValue('NEXT_PUBLIC_WEBSOCKET_URL'),
    };
    
    // Валидация только в production
    if (process.env.NODE_ENV === 'production') {
      const result = publicEnvSchema.safeParse(publicEnv);
      if (!result.success) {
        console.error('Config validation error:', result.error.format());
        throw new Error('Invalid configuration');
      }
    }
    
    return {
      api: {
        url: publicEnv.NEXT_PUBLIC_API_URL,
        internalUrl: isServer ? process.env.INTERNAL_API_URL : undefined,
        websocketUrl: publicEnv.NEXT_PUBLIC_WEBSOCKET_URL,
      },
      storage: {
        minioUrl: publicEnv.NEXT_PUBLIC_MINIO_URL,
        imageHosts: this.parseImageHosts(publicEnv.NEXT_PUBLIC_IMAGE_HOSTS),
        imagePathPattern: publicEnv.NEXT_PUBLIC_IMAGE_PATH_PATTERN || '/listings/**',
      },
      env: {
        isProduction: process.env.NODE_ENV === 'production',
        isDevelopment: process.env.NODE_ENV === 'development',
        isServer,
      },
      features: {
        enableChat: !!publicEnv.NEXT_PUBLIC_WEBSOCKET_URL,
        enablePayments: getEnvValue('NEXT_PUBLIC_ENABLE_PAYMENTS') === 'true',
      },
    };
  }
  
  public getConfig(): Config {
    if (!this.config) {
      this.config = this.loadConfig();
    }
    return this.config;
  }
  
  // Метод для получения правильного API URL в зависимости от контекста
  public getApiUrl(options?: { internal?: boolean }): string {
    const config = this.getConfig();
    
    // Для серверных запросов используем внутренний URL если доступен
    if (config.env.isServer && options?.internal && config.api.internalUrl) {
      return config.api.internalUrl;
    }
    
    // В разработке на клиенте используем пустую строку для proxy
    if (config.env.isDevelopment && !config.env.isServer) {
      return '';
    }
    
    return config.api.url;
  }
}

export default new ConfigManager();
```

## 3. Создание .env.example

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:3000
INTERNAL_API_URL=http://backend:3000  # Внутренний URL для SSR (только для Docker)

# Storage Configuration
NEXT_PUBLIC_MINIO_URL=http://localhost:9000
NEXT_PUBLIC_IMAGE_HOSTS=http:localhost:9000,https:svetu.rs:443

# WebSocket Configuration (optional)
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:3000

# Feature Flags
NEXT_PUBLIC_ENABLE_PAYMENTS=false

# Environment
NODE_ENV=development
```

## 4. Улучшенный Dockerfile с runtime поддержкой

```dockerfile
FROM node:22-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev

FROM node:22-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Используем плейсхолдеры для runtime замены
ENV NEXT_PUBLIC_API_URL="__NEXT_PUBLIC_API_URL__"
ENV NEXT_PUBLIC_MINIO_URL="__NEXT_PUBLIC_MINIO_URL__"

RUN npm run build

FROM node:22-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=builder --chown=nextjs:nodejs /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Entrypoint для runtime конфигурации
COPY --chown=nextjs:nodejs docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

USER nextjs
EXPOSE 3000

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["node", "server.js"]
```

## 5. API client с поддержкой разных контекстов

### src/services/api-client.ts:
```typescript
import configManager from '@/config';

class ApiClient {
  private getBaseUrl(isInternal: boolean = false): string {
    return configManager.getApiUrl({ internal: isInternal });
  }
  
  async fetch(endpoint: string, options?: RequestInit, context?: { internal?: boolean }) {
    const baseUrl = this.getBaseUrl(context?.internal);
    const url = `${baseUrl}${endpoint}`;
    
    // Добавляем заголовки для внутренних запросов
    if (context?.internal) {
      options = {
        ...options,
        headers: {
          ...options?.headers,
          'X-Internal-Request': 'true',
        },
      };
    }
    
    return fetch(url, options);
  }
}
```

## 6. Использование в компонентах

### Клиентский компонент:
```typescript
'use client';
import { env } from 'next-runtime-env';

export default function ClientComponent() {
  // Runtime переменные доступны через env()
  const apiUrl = env('NEXT_PUBLIC_API_URL');
  
  return <div>API: {apiUrl}</div>;
}
```

### Серверный компонент с SSR:
```typescript
import configManager from '@/config';

export default async function ServerComponent() {
  // Используем внутренний URL для серверных запросов
  const apiUrl = configManager.getApiUrl({ internal: true });
  const data = await fetch(`${apiUrl}/api/data`);
  
  return <div>{/* render data */}</div>;
}
```

## Преимущества предлагаемого подхода:

1. **Runtime конфигурация**: Один Docker образ для всех окружений
2. **Типобезопасность**: Zod схемы для валидации
3. **Разделение контекстов**: Разные URL для сервера и клиента
4. **Feature flags**: Динамическое включение/выключение функций
5. **Безопасность**: Серверные переменные не попадают в браузер
6. **Developer experience**: Понятные ошибки при неправильной конфигурации