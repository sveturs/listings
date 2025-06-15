# Шаг 6: Примеры использования

## Цель
Продемонстрировать использование новой системы конфигурации в различных сценариях и компонентах.

## Примеры

### 6.1 Настройка для различных окружений

#### Локальная разработка
```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:3000
NEXT_PUBLIC_MINIO_URL=http://localhost:9000
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:3000
NEXT_PUBLIC_ENABLE_PAYMENTS=false
NODE_ENV=development

# Запуск
yarn dev
```

#### Docker разработка
```bash
# Запуск с переменными окружения
docker run -p 3001:3000 \
  -e NEXT_PUBLIC_API_URL=http://host.docker.internal:3000 \
  -e INTERNAL_API_URL=http://backend:3000 \
  -e NEXT_PUBLIC_MINIO_URL=http://host.docker.internal:9000 \
  svetu-frontend:latest
```

#### Production deployment
```yaml
# kubernetes/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
spec:
  template:
    spec:
      containers:
      - name: frontend
        image: harbor.svetu.rs/svetu/frontend:latest
        env:
        - name: NEXT_PUBLIC_API_URL
          value: "https://api.svetu.rs"
        - name: INTERNAL_API_URL
          value: "http://api-service:3000"
        - name: NEXT_PUBLIC_MINIO_URL
          value: "https://svetu.rs"
        - name: NEXT_PUBLIC_ENABLE_PAYMENTS
          value: "true"
```

### 6.2 Использование в компонентах

#### Server Component с данными
```typescript
// app/[locale]/page.tsx
import { marketplaceApi } from '@/services/api/endpoints';
import { MarketplaceList } from '@/components/marketplace/MarketplaceList';
import configManager from '@/config';

export default async function HomePage() {
  // Использует внутренний URL автоматически
  const response = await marketplaceApi.getListings({ 
    page: 1, 
    limit: 20 
  });

  // Проверяем feature flags
  const paymentsEnabled = configManager.isFeatureEnabled('enablePayments');

  return (
    <div className="container mx-auto">
      <h1 className="text-3xl font-bold mb-6">
        Welcome to Sve Tu Marketplace
      </h1>
      
      {paymentsEnabled && (
        <div className="alert alert-info mb-4">
          🎉 Payments are now available!
        </div>
      )}
      
      <MarketplaceList 
        initialData={response.data?.items || []} 
        totalCount={response.data?.total || 0}
      />
    </div>
  );
}
```

#### Client Component с API вызовами
```typescript
// components/marketplace/MarketplaceFilters.tsx
'use client';

import { useState, useEffect } from 'react';
import { useApi } from '@/hooks/useApi';
import { marketplaceApi } from '@/services/api/endpoints';
import { useConfig } from '@/hooks/useConfig';

export function MarketplaceFilters({ onFilterChange }) {
  const [categories, setCategories] = useState([]);
  const config = useConfig();

  // Загрузка категорий
  const { data, loading } = useApi(
    () => marketplaceApi.getCategories(),
    { immediate: true }
  );

  useEffect(() => {
    if (data) {
      setCategories(data);
    }
  }, [data]);

  return (
    <div className="filters">
      <h3>Filters</h3>
      
      {/* Debug info в development */}
      {config.env.isDevelopment && (
        <div className="text-xs text-gray-500 mb-2">
          API: {config.api.url}
        </div>
      )}
      
      {loading ? (
        <div className="skeleton h-32 w-full"></div>
      ) : (
        <select 
          onChange={(e) => onFilterChange({ category: e.target.value })}
          className="select select-bordered w-full"
        >
          <option value="">All Categories</option>
          {categories.map(cat => (
            <option key={cat.id} value={cat.id}>{cat.name}</option>
          ))}
        </select>
      )}
    </div>
  );
}
```

#### Component с условным рендерингом по feature flags
```typescript
// components/PaymentButton.tsx
'use client';

import { useFeature } from '@/hooks/useConfig';

interface PaymentButtonProps {
  amount: number;
  onPayment: () => void;
}

export function PaymentButton({ amount, onPayment }: PaymentButtonProps) {
  const paymentsEnabled = useFeature('enablePayments');

  // Не рендерим если payments отключены
  if (!paymentsEnabled) {
    return null;
  }

  return (
    <button 
      onClick={onPayment}
      className="btn btn-primary"
    >
      Pay ${amount}
    </button>
  );
}
```

### 6.3 Работа с изображениями

#### Image component с динамическим URL
```typescript
// components/OptimizedImage.tsx
'use client';

import Image from 'next/image';
import configManager from '@/config';
import { useState } from 'react';

interface OptimizedImageProps {
  src: string;
  alt: string;
  width: number;
  height: number;
}

export function OptimizedImage({ src, alt, width, height }: OptimizedImageProps) {
  const [error, setError] = useState(false);
  
  // Строим полный URL для изображения
  const imageUrl = configManager.buildImageUrl(src);
  const fallbackUrl = '/placeholder-listing.jpg';

  return (
    <Image
      src={error ? fallbackUrl : imageUrl}
      alt={alt}
      width={width}
      height={height}
      onError={() => setError(true)}
      placeholder="blur"
      blurDataURL="data:image/jpeg;base64,/9j/4AAQSkZJRg..."
    />
  );
}
```

### 6.4 WebSocket с runtime конфигурацией

#### WebSocket manager
```typescript
// utils/websocket.ts
import configManager from '@/config';

class WebSocketManager {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnects = 5;

  connect() {
    const config = configManager.getConfig();
    
    // Проверяем доступность WebSocket
    if (!config.api.websocketUrl) {
      console.warn('WebSocket URL not configured');
      return;
    }

    try {
      this.ws = new WebSocket(config.api.websocketUrl);
      
      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
      };

      this.ws.onclose = () => {
        this.handleReconnect();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
    }
  }

  private handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnects) {
      this.reconnectAttempts++;
      const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
      
      setTimeout(() => {
        console.log(`Reconnecting... (attempt ${this.reconnectAttempts})`);
        this.connect();
      }, delay);
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  send(message: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected');
    }
  }
}

export const wsManager = new WebSocketManager();
```

### 6.5 Middleware с конфигурацией

#### Auth middleware
```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { env } from 'next-runtime-env';

export function middleware(request: NextRequest) {
  // Получаем конфигурацию
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  const isProduction = process.env.NODE_ENV === 'production';

  // Проверяем авторизацию для защищенных маршрутов
  if (request.nextUrl.pathname.startsWith('/profile')) {
    const token = request.cookies.get('auth-token');
    
    if (!token) {
      const loginUrl = new URL('/login', request.url);
      loginUrl.searchParams.set('from', request.nextUrl.pathname);
      return NextResponse.redirect(loginUrl);
    }
  }

  // Добавляем security headers в production
  if (isProduction) {
    const response = NextResponse.next();
    response.headers.set('X-Frame-Options', 'DENY');
    response.headers.set('X-Content-Type-Options', 'nosniff');
    response.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');
    return response;
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/profile/:path*', '/admin/:path*'],
};
```

### 6.6 Testing с различными конфигурациями

#### Component test
```typescript
// __tests__/components/PaymentButton.test.tsx
import { render, screen } from '@testing-library/react';
import { PaymentButton } from '@/components/PaymentButton';
import configManager from '@/config';

describe('PaymentButton', () => {
  beforeEach(() => {
    // Сбрасываем конфигурацию перед каждым тестом
    configManager.resetConfig();
  });

  it('should render when payments enabled', () => {
    // Mock конфигурации
    process.env.NEXT_PUBLIC_ENABLE_PAYMENTS = 'true';
    
    render(<PaymentButton amount={100} onPayment={() => {}} />);
    
    expect(screen.getByText('Pay $100')).toBeInTheDocument();
  });

  it('should not render when payments disabled', () => {
    process.env.NEXT_PUBLIC_ENABLE_PAYMENTS = 'false';
    
    render(<PaymentButton amount={100} onPayment={() => {}} />);
    
    expect(screen.queryByText('Pay $100')).not.toBeInTheDocument();
  });
});
```

#### E2E test с разными окружениями
```typescript
// e2e/config.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Configuration', () => {
  test('should use correct API URL in production', async ({ page }) => {
    // Запускаем с production конфигурацией
    await page.goto('/', {
      waitUntil: 'networkidle',
    });

    // Проверяем что запросы идут на правильный URL
    const apiRequest = await page.waitForRequest(req => 
      req.url().includes('/api/v1/marketplace')
    );
    
    expect(apiRequest.url()).toContain('https://api.svetu.rs');
  });

  test('should show debug info in development', async ({ page }) => {
    // Устанавливаем dev окружение
    process.env.NODE_ENV = 'development';
    
    await page.goto('/');
    
    // Проверяем наличие debug информации
    const debugInfo = await page.locator('.debug-info').textContent();
    expect(debugInfo).toContain('API: http://localhost:3000');
  });
});
```

### 6.7 Мониторинг конфигурации

#### Health check endpoint
```typescript
// app/api/health/route.ts
import { NextResponse } from 'next/server';
import configManager from '@/config';

export async function GET() {
  const config = configManager.getConfig();
  const errors = configManager.getValidationErrors();

  // Проверяем критичные сервисы
  const checks = {
    api: await checkApiHealth(config.api.url),
    storage: await checkStorageHealth(config.storage.minioUrl),
    config: errors.length === 0,
  };

  const isHealthy = Object.values(checks).every(v => v === true);

  return NextResponse.json({
    status: isHealthy ? 'healthy' : 'unhealthy',
    timestamp: new Date().toISOString(),
    checks,
    config: {
      environment: config.env.isProduction ? 'production' : 'development',
      features: config.features,
    },
    errors: errors.length > 0 ? errors : undefined,
  }, {
    status: isHealthy ? 200 : 503,
  });
}

async function checkApiHealth(url: string): Promise<boolean> {
  try {
    const response = await fetch(`${url}/health`, {
      method: 'GET',
      signal: AbortSignal.timeout(5000),
    });
    return response.ok;
  } catch {
    return false;
  }
}

async function checkStorageHealth(url: string): Promise<boolean> {
  try {
    const response = await fetch(`${url}/minio/health/live`, {
      method: 'GET',
      signal: AbortSignal.timeout(5000),
    });
    return response.ok;
  } catch {
    return false;
  }
}
```

## Миграция существующего кода

### Чеклист миграции
1. ✅ Установить `next-runtime-env`
2. ✅ Обновить `layout.tsx` с `PublicEnvScript`
3. ✅ Обновить `config/types.ts` с Zod схемами
4. ✅ Обновить `config/index.ts` с runtime поддержкой
5. ✅ Создать `.env.example` файлы
6. ✅ Обновить `Dockerfile` и `docker-entrypoint.sh`
7. ✅ Обновить `api-client.ts` с контекстами
8. ✅ Протестировать в разных окружениях

### Команды для тестирования
```bash
# Локальная разработка
yarn dev

# Docker с дефолтной конфигурацией
make docker-run

# Docker с кастомной конфигурацией
docker run -p 3001:3000 \
  -e NEXT_PUBLIC_API_URL=https://staging.api.svetu.rs \
  -e NEXT_PUBLIC_ENABLE_PAYMENTS=true \
  svetu-frontend:latest

# Проверка health
curl http://localhost:3001/api/health
```

## Результат
После выполнения всех шагов:
1. ✅ Runtime конфигурация работает во всех окружениях
2. ✅ Один Docker образ для dev/staging/production
3. ✅ Типобезопасная конфигурация с валидацией
4. ✅ Автоматический выбор URL для SSR/CSR
5. ✅ Feature flags для управления функциональностью
6. ✅ Мониторинг и health checks