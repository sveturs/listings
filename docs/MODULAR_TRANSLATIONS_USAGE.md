# 🚀 Руководство по использованию модульной системы переводов

## 📋 Быстрый старт

### 1. Активация системы

```bash
# В файле .env.local
USE_MODULAR_I18N=true
```

### 2. Запуск приложения

```bash
yarn dev
```

## 🎯 Использование в компонентах

### Базовый пример

```typescript
// Компонент использует переводы из модуля marketplace
'use client';

import { useTranslations } from 'next-intl';

export function MarketplaceComponent() {
  const t = useTranslations('marketplace');
  
  return (
    <div>
      <h1>{t('title')}</h1>
      <button>{t('createListing')}</button>
    </div>
  );
}
```

### Использование вложенных ключей

```typescript
// ВАЖНО: В модульной системе НЕ используются вложенные пути в useTranslations!
// Правильно:
const t = useTranslations('marketplace');
t('listing.title') // для доступа к вложенному ключу
t('listing.price')

// Неправильно:
// const t = useTranslations('marketplace.listing'); // Это вызовет ошибку!
```

## 📄 Обновление страниц

### Server Component с модулями

```typescript
// app/[locale]/marketplace/page.tsx
import { NextIntlClientProvider } from 'next-intl';
import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function MarketplacePage({ 
  params: { locale } 
}) {
  // Загружаем только нужные модули
  const messages = await loadMessages(locale, ['marketplace']);
  
  return (
    <NextIntlClientProvider messages={messages}>
      <MarketplaceContent />
    </NextIntlClientProvider>
  );
}
```

### Client Component с lazy loading

```typescript
'use client';

import { useState, useEffect } from 'react';
import { loadMessages } from '@/lib/i18n/loadMessages';
import { useLocale } from 'next-intl';

export function DynamicFeature() {
  const locale = useLocale();
  const [carsModule, setCarsModule] = useState(null);
  
  useEffect(() => {
    // Загружаем модуль когда он нужен
    loadMessages(locale, ['cars']).then(messages => {
      setCarsModule(messages.cars);
    });
  }, [locale]);
  
  if (!carsModule) return <div>Loading...</div>;
  
  return <div>{carsModule.title}</div>;
}
```

## 🗂️ Структура модулей

```
src/messages/
├── ru/
│   ├── common.json      # Базовые UI элементы
│   ├── auth.json        # Авторизация
│   ├── marketplace.json # Маркетплейс
│   ├── admin.json       # Админ панель
│   ├── storefront.json  # Витрины
│   ├── cars.json        # Автомобили
│   ├── cart.json        # Корзина
│   └── chat.json        # Чат
├── en/
│   └── ... (та же структура)
└── sr/
    └── ... (та же структура)
```

## 🔍 Определение нужных модулей

### Автоматическое определение по URL

```typescript
import { getRequiredModules } from '@/lib/i18n/loadMessages';

// Для /ru/marketplace/listings
const modules = getRequiredModules(pathname);
// Вернет: ['common', 'marketplace']

// Для /ru/admin/users  
const modules = getRequiredModules(pathname);
// Вернет: ['common', 'admin']
```

### Ручное указание модулей

```typescript
// Загружаем конкретные модули
const messages = await loadMessages(locale, [
  'common',      // Всегда нужен
  'marketplace', // Основной функционал
  'cart'         // Дополнительный
]);
```

## ⚡ Оптимизация производительности

### Предзагрузка модулей

```typescript
// В компоненте навигации
import { preloadModules } from '@/lib/i18n/loadMessages';

function Navigation() {
  const locale = useLocale();
  
  const handleHover = (modules) => {
    preloadModules(locale, modules);
  };
  
  return (
    <nav>
      <Link 
        href="/marketplace"
        onMouseEnter={() => handleHover(['marketplace'])}
      >
        Маркетплейс
      </Link>
    </nav>
  );
}
```

### Service Worker (offline поддержка)

Service Worker автоматически кэширует загруженные модули для offline доступа.

```typescript
// Предзагрузка критичных модулей
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.ready.then(registration => {
    registration.active.postMessage({
      type: 'PRELOAD_MODULES',
      locale: 'ru',
      modules: ['common', 'marketplace']
    });
  });
}
```

## 📊 Мониторинг

### Отслеживание загрузки модулей

```typescript
// В консоли браузера
window.__TRANSLATION_MODULES_LOADED__
// ['common', 'marketplace', 'auth']

// Размер загруженных модулей
window.__TRANSLATION_MODULES_SIZE__
// { common: 12800, marketplace: 38420, auth: 17430 }
```

### Performance метрики

```typescript
// Измерение времени загрузки
performance.mark('translation-module-start');
await loadMessages(locale, ['admin']);
performance.mark('translation-module-end');
performance.measure(
  'translation-module-load',
  'translation-module-start',
  'translation-module-end'
);
```

## 🐛 Отладка

### Проверка загруженных модулей

```bash
# В DevTools Console
localStorage.getItem('translation-modules-loaded')
```

### Очистка кэша

```typescript
import { clearModuleCache } from '@/lib/i18n/loadMessages';

// При смене языка
clearModuleCache();
```

### Логирование

```bash
# Включить debug логи
DEBUG=translations:* yarn dev
```

## ✅ Чеклист интеграции

- [ ] Установить `USE_MODULAR_I18N=true` в .env.local
- [ ] Обновить layout.tsx для использования ModularIntlProvider
- [ ] Обновить страницы для загрузки нужных модулей
- [ ] Добавить предзагрузку в навигацию
- [ ] Протестировать offline режим
- [ ] Проверить метрики производительности

## 📈 Ожидаемые результаты

- **Initial JS**: -70-85% (с 52KB до 10-15KB)
- **FCP**: -20-30%
- **TTI**: -15-25%
- **Lighthouse Performance**: +5-10 баллов