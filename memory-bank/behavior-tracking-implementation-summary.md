# Behavior Tracking Implementation Summary

## Дата выполнения
2025-07-08

## Обзор выполненной работы
Создан комплексный frontend hook `useBehaviorTracking` для поведенческого трекинга пользователей в проекте hostel-booking-system.

## Созданные файлы

### 1. Типы и интерфейсы
**`/data/hostel-booking-system/frontend/svetu/src/types/behavior.ts`**
- Определены TypeScript типы для всех типов событий
- Интерфейсы для базовых и специфичных событий
- Конфигурация hook'а и состояние
- Полная типизация всех методов и данных

### 2. Utility функции для определения устройства
**`/data/hostel-booking-system/frontend/svetu/src/utils/deviceDetection.ts`**
- Автоматическое определение типа устройства (desktop/mobile/tablet)
- Распознавание браузера и версии
- Определение операционной системы
- Дополнительные функции для получения информации об экране и возможностях браузера

### 3. Основной hook
**`/data/hostel-booking-system/frontend/svetu/src/hooks/useBehaviorTracking.ts`**
- Автоматическая генерация и сохранение session_id
- Батчевая отправка событий (каждые 5 секунд или при 10 событиях)
- Retry логика с экспоненциальной задержкой
- Дебаунс для предотвращения спама событий
- Поддержка анонимных пользователей
- Graceful shutdown с navigator.sendBeacon

### 4. Переводы
Добавлены переводы ошибок и сообщений в файлы локализации:
- **`/data/hostel-booking-system/frontend/svetu/src/messages/en.json`**
- **`/data/hostel-booking-system/frontend/svetu/src/messages/ru.json`**
- **`/data/hostel-booking-system/frontend/svetu/src/messages/sr.json`**

### 5. Примеры и тесты
**`/data/hostel-booking-system/frontend/svetu/src/hooks/__tests__/useBehaviorTracking.example.tsx`**
- Полный пример использования hook'а
- Демонстрация всех типов событий
- UI для тестирования функциональности

**`/data/hostel-booking-system/frontend/svetu/src/app/[locale]/test-behavior-tracking/page.tsx`**
- Тестовая страница для проверки работы hook'а

**`/data/hostel-booking-system/frontend/svetu/src/utils/__tests__/deviceDetection.test.ts`**
- Unit тесты для utility функций определения устройства

### 6. Документация
**`/data/hostel-booking-system/frontend/svetu/src/hooks/useBehaviorTracking.README.md`**
- Подробная документация по использованию
- Примеры кода для всех типов событий
- Описание конфигурации и API

## Реализованные типы событий

### 1. search_performed
Отслеживание выполненных поисков с фильтрами и сортировкой

### 2. result_clicked
Клики по результатам поиска с позицией и временем

### 3. item_viewed
Просмотры товаров с категорией и ценой

### 4. item_purchased
Покупки с суммой, валютой и методом оплаты

### 5. search_filter_applied
Применение фильтров в поиске

### 6. search_sort_changed
Изменение сортировки результатов

## Ключевые особенности

### Автоматизация
- Автоматическое определение device_type, browser, user_agent
- Генерация уникального session_id с сохранением в localStorage
- Автоматическая батчевая отправка событий

### Надежность
- Retry логика с настраиваемым количеством попыток
- Экспоненциальная задержка между повторными попытками
- Отправка событий при закрытии страницы через navigator.sendBeacon

### Производительность
- Дебаунс для предотвращения спама событий (100мс)
- Батчевая отправка для минимизации HTTP запросов
- Ленивая загрузка компонентов

### Гибкость
- Полная настройка через опции hook'а
- Поддержка анонимных пользователей
- Контекстное отслеживание для связанных событий

### TypeScript
- Полная типизация всех событий и параметров
- Type-safe методы для каждого типа события
- Автокомплит и проверка ошибок в IDE

## API интеграция

Hook отправляет события на endpoint `/api/v1/analytics/track` в формате:

```json
{
  "events": [
    {
      "session_id": "session_1720425600000_abc123",
      "event_type": "search_performed",
      "timestamp": "2025-07-08T08:00:00.000Z",
      "user_id": "user-123",
      "device_type": "desktop",
      "browser": "Chrome 91.0.4472.124",
      "user_agent": "Mozilla/5.0...",
      "search_query": "iPhone 15",
      "results_count": 25,
      "search_duration_ms": 120
    }
  ],
  "batch_id": "batch_1720425600000_def456",
  "created_at": "2025-07-08T08:00:00.000Z"
}
```

## Тестирование

### Unit тесты
- ✅ Тесты для utility функций определения устройства (15 тестов)
- ✅ Проверка различных браузеров и ОС
- ✅ Валидация типов устройств

### Интеграционное тестирование
- ✅ Тестовая страница с примерами всех событий
- ✅ Демонстрация состояния hook'а в реальном времени
- ✅ Проверка батчевой отправки и retry логики

### Сборка и форматирование
- ✅ Успешная сборка проекта с Next.js
- ✅ Прохождение ESLint проверок
- ✅ Автоматическое форматирование с Prettier

## Совместимость

- ✅ React 18+
- ✅ Next.js 15+ (App Router)
- ✅ TypeScript 5+
- ✅ Современные браузеры с ES2020+
- ✅ Серверный рендеринг (SSR)

## Использование

```typescript
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';

function MyComponent() {
  const {
    trackSearchPerformed,
    trackItemViewed,
    trackItemPurchased,
    state
  } = useBehaviorTracking({
    enabled: true,
    userId: user?.id,
    batchSize: 10,
    batchTimeout: 5000
  });

  // Использование методов для отслеживания событий
  await trackSearchPerformed({
    search_query: 'iPhone',
    results_count: 25,
    search_duration_ms: 120
  });
}
```

## Следующие шаги

1. **Backend интеграция**: Необходимо реализовать endpoint `/api/v1/analytics/track` на backend'е
2. **База данных**: Создать таблицы для хранения событий поведенческого трекинга
3. **Аналитика**: Добавить dashboard для просмотра собранных данных
4. **Интеграция в компоненты**: Внедрить hook в существующие компоненты поиска и просмотра товаров

## Структура файлов

```
frontend/svetu/src/
├── types/
│   └── behavior.ts
├── utils/
│   ├── deviceDetection.ts
│   └── __tests__/
│       └── deviceDetection.test.ts
├── hooks/
│   ├── useBehaviorTracking.ts
│   ├── useBehaviorTracking.README.md
│   └── __tests__/
│       └── useBehaviorTracking.example.tsx
├── app/[locale]/
│   └── test-behavior-tracking/
│       └── page.tsx
└── messages/
    ├── en.json (обновлен)
    ├── ru.json (обновлен)
    └── sr.json (обновлен)
```

## Конфигурация по умолчанию

```typescript
const DEFAULT_CONFIG = {
  enabled: true,
  batchSize: 10,
  batchTimeout: 5000, // 5 секунд
  maxRetries: 3,
  endpoint: '/api/v1/analytics/track'
};
```

Hook готов к использованию и интеграции в существующие компоненты проекта.