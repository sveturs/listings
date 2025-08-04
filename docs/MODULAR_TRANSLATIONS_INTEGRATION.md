# Руководство по интеграции модульной системы переводов

## 📊 Результаты анализа

### Экономия размера bundle по страницам:
- **Главная страница**: 74.8% экономии (13KB вместо 51.7KB)
- **Админ панель**: 77.7% экономии (11.5KB вместо 51.7KB)
- **Корзина**: 87.1% экономии (6.7KB вместо 51.7KB)
- **Витрина**: 85.3% экономии (7.6KB вместо 51.7KB)
- **Автомобили**: 73.3% экономии (13.8KB вместо 51.7KB)

## 🚀 Пошаговая интеграция

### 1. Обновление переменных окружения

```bash
# .env.local
USE_MODULAR_I18N=true
```

### 2. Обновление layout.tsx

```typescript
// app/[locale]/layout.tsx
import { ModularIntlProvider } from '@/providers/ModularIntlProvider';
import { loadMessages, getRequiredModules } from '@/lib/i18n/loadMessages';

export default async function LocaleLayout({
  children,
  params: { locale }
}: {
  children: React.ReactNode;
  params: { locale: string };
}) {
  // Загружаем базовые модули для layout
  const messages = await loadMessages(locale, ['common']);
  
  return (
    <html lang={locale}>
      <body>
        <ModularIntlProvider locale={locale} messages={messages}>
          {children}
        </ModularIntlProvider>
      </body>
    </html>
  );
}
```

### 3. Обновление страниц

#### Пример для страницы маркетплейса:

```typescript
// app/[locale]/marketplace/page.tsx
import { loadMessages } from '@/lib/i18n/loadMessages';
import { NextIntlClientProvider } from 'next-intl';

export default async function MarketplacePage({ 
  params: { locale } 
}: { 
  params: { locale: string } 
}) {
  // Загружаем необходимые модули
  const messages = await loadMessages(locale, ['marketplace']);
  
  return (
    <NextIntlClientProvider messages={messages}>
      {/* Компоненты страницы */}
    </NextIntlClientProvider>
  );
}
```

#### Пример для админ панели:

```typescript
// app/[locale]/admin/page.tsx
import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function AdminPage({ 
  params: { locale } 
}: { 
  params: { locale: string } 
}) {
  const messages = await loadMessages(locale, ['admin']);
  
  return (
    <NextIntlClientProvider messages={messages}>
      {/* Админ компоненты */}
    </NextIntlClientProvider>
  );
}
```

### 4. Обновление компонентов

#### Использование с namespace:

```typescript
// Старый способ
const t = useTranslations();
t('marketplace.listing.title'); // длинный ключ

// Новый способ с namespace
const t = useTranslations('marketplace.listing');
t('title'); // короткий ключ
```

#### Динамическая загрузка в клиентских компонентах:

```typescript
'use client';

import { useState, useEffect } from 'react';
import { loadMessages } from '@/lib/i18n/loadMessages';
import { useLocale } from 'next-intl';

export function DynamicComponent() {
  const locale = useLocale();
  const [messages, setMessages] = useState(null);
  
  useEffect(() => {
    // Загружаем модуль когда компонент монтируется
    loadMessages(locale, ['cars']).then(setMessages);
  }, [locale]);
  
  if (!messages) return <div>Loading...</div>;
  
  return <div>{messages.cars.title}</div>;
}
```

### 5. Middleware для автоматического определения модулей

```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { getRequiredModules } from '@/lib/i18n/loadMessages';

export function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;
  const modules = getRequiredModules(pathname);
  
  // Добавляем заголовок с необходимыми модулями
  const response = NextResponse.next();
  response.headers.set('X-Required-Modules', modules.join(','));
  
  return response;
}
```

### 6. Оптимизация с помощью предзагрузки

```typescript
// components/Navigation.tsx
'use client';

import Link from 'next/link';
import { preloadModules } from '@/lib/i18n/loadMessages';
import { useLocale } from 'next-intl';

export function Navigation() {
  const locale = useLocale();
  
  const handleMouseEnter = (modules: string[]) => {
    // Предзагружаем модули при наведении
    preloadModules(locale, modules);
  };
  
  return (
    <nav>
      <Link 
        href="/marketplace"
        onMouseEnter={() => handleMouseEnter(['marketplace'])}
      >
        Маркетплейс
      </Link>
      <Link 
        href="/admin"
        onMouseEnter={() => handleMouseEnter(['admin'])}
      >
        Админка
      </Link>
    </nav>
  );
}
```

## 📝 Чеклист миграции

- [ ] Обновить .env.local с USE_MODULAR_I18N=true
- [ ] Заменить i18n.ts на i18n-new.ts
- [ ] Обновить корневой layout с ModularIntlProvider
- [ ] Мигрировать страницы на loadMessages
- [ ] Обновить компоненты для использования namespace
- [ ] Добавить предзагрузку в навигацию
- [ ] Настроить кэширование на CDN
- [ ] Добавить мониторинг производительности

## 🔧 Настройка кэширования

### Next.js конфигурация:

```javascript
// next.config.js
module.exports = {
  // Долгосрочное кэширование для модулей переводов
  async headers() {
    return [
      {
        source: '/_next/static/chunks/messages-*.js',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
    ];
  },
};
```

### Service Worker для offline:

```javascript
// public/sw.js
self.addEventListener('fetch', (event) => {
  if (event.request.url.includes('/messages/')) {
    event.respondWith(
      caches.match(event.request).then((response) => {
        return response || fetch(event.request).then((response) => {
          return caches.open('translations-v1').then((cache) => {
            cache.put(event.request, response.clone());
            return response;
          });
        });
      })
    );
  }
});
```

## 📈 Мониторинг

### Метрики для отслеживания:

```typescript
// lib/metrics.ts
export function trackTranslationLoad(module: string, duration: number) {
  // Отправляем в аналитику
  if (typeof window !== 'undefined' && window.gtag) {
    window.gtag('event', 'translation_module_load', {
      module_name: module,
      load_duration: duration,
      locale: document.documentElement.lang,
    });
  }
}
```

## 🎯 Ожидаемые результаты

1. **Производительность**:
   - First Contentful Paint: -20-30%
   - Time to Interactive: -15-25%
   - Lighthouse Score: +5-10 баллов

2. **Размер bundle**:
   - Initial JS: -70-85%
   - Общий размер: без изменений (lazy loading)

3. **Пользовательский опыт**:
   - Быстрее начальная загрузка
   - Плавные переходы между страницами
   - Работа offline (с Service Worker)