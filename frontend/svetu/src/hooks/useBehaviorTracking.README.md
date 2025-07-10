# useBehaviorTracking Hook

React hook для поведенческого трекинга пользователей в проекте hostel-booking-system.

## Описание

Hook `useBehaviorTracking` предоставляет комплексное решение для отслеживания поведения пользователей с автоматической батчевой отправкой событий, retry логикой и поддержкой анонимных пользователей.

## Основные возможности

- ✅ **Автоматическая генерация session_id** - уникальный идентификатор сессии с сохранением в localStorage
- ✅ **Батчевая отправка событий** - события отправляются пакетами каждые 5 секунд или при накоплении 10 событий
- ✅ **Retry логика** - автоматические повторные попытки при сетевых ошибках
- ✅ **Поддержка анонимных пользователей** - работает без авторизации
- ✅ **Дебаунс событий** - предотвращение спама одинаковых событий
- ✅ **Автоматическое определение устройства** - device_type, browser, user_agent
- ✅ **Graceful shutdown** - отправка событий при закрытии страницы через navigator.sendBeacon
- ✅ **TypeScript поддержка** - полная типизация всех событий и параметров

## Типы событий

### 1. search_performed

Выполнен поиск пользователем

```typescript
await trackSearchPerformed({
  search_query: 'iPhone 15',
  search_filters: { category: 'electronics', price_min: 1000 },
  search_sort: 'price-asc',
  results_count: 25,
  search_duration_ms: 120,
});
```

### 2. result_clicked

Клик по результату поиска

```typescript
await trackResultClicked({
  search_query: 'iPhone 15',
  clicked_item_id: 'item-123',
  click_position: 3,
  total_results: 25,
  click_time_from_search_ms: 5000,
});
```

### 3. item_viewed

Просмотр конкретного товара

```typescript
await trackItemViewed({
  item_id: 'item-123',
  item_title: 'iPhone 15 Pro Max',
  item_category: 'electronics',
  item_price: 1199.99,
  previous_page: 'search-results',
});
```

### 4. item_purchased

Покупка товара

```typescript
await trackItemPurchased({
  item_id: 'item-123',
  item_title: 'iPhone 15 Pro Max',
  item_category: 'electronics',
  purchase_amount: 1199.99,
  currency: 'USD',
  payment_method: 'credit_card',
  time_from_first_view_ms: 300000, // 5 минут
});
```

### 5. search_filter_applied

Применение фильтра в поиске

```typescript
await trackSearchFilterApplied({
  search_query: 'iPhone 15',
  filter_type: 'price_range',
  filter_value: '1000-1500',
  results_count_before: 25,
  results_count_after: 8,
});
```

### 6. search_sort_changed

Изменение сортировки результатов

```typescript
await trackSearchSortChanged({
  search_query: 'iPhone 15',
  sort_type: 'price-desc',
  previous_sort: 'price-asc',
  results_count: 8,
});
```

## Использование

### Базовое использование

```typescript
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';

function MyComponent() {
  const {
    state,
    trackSearchPerformed,
    trackResultClicked,
    trackItemViewed,
    // ... другие методы
  } = useBehaviorTracking({
    enabled: true,
    userId: 'user-123', // опционально
    batchSize: 10,
    batchTimeout: 5000
  });

  const handleSearch = async (query: string) => {
    await trackSearchPerformed({
      search_query: query,
      results_count: searchResults.length,
      search_duration_ms: searchTime
    });
  };

  return (
    <div>
      <p>Tracking enabled: {state.isEnabled}</p>
      <p>Pending events: {state.pendingEvents.length}</p>
      {/* UI компоненты */}
    </div>
  );
}
```

### Расширенное использование с контекстом

```typescript
function SearchComponent() {
  const tracking = useBehaviorTracking({ userId: user?.id });

  const handleSearch = async (query: string, filters: any, sort: string) => {
    // Начинаем отслеживание поиска
    tracking.startSearch(query, filters, sort);

    const startTime = Date.now();
    const results = await performSearch(query, filters, sort);
    const searchDuration = Date.now() - startTime;

    // Отправляем событие поиска
    await tracking.trackSearchPerformed({
      search_query: query,
      search_filters: filters,
      search_sort: sort,
      results_count: results.length,
      search_duration_ms: searchDuration
    });
  };

  const handleResultClick = async (item: any, position: number) => {
    await tracking.trackResultClicked({
      search_query: tracking.state.context.currentSearchQuery || '',
      clicked_item_id: item.id,
      click_position: position,
      total_results: tracking.state.context.lastSearchResults || 0,
      click_time_from_search_ms: Date.now() - (tracking.state.context.currentSearchStartTime || Date.now())
    });
  };

  return <SearchInterface onSearch={handleSearch} onResultClick={handleResultClick} />;
}
```

### Отслеживание просмотра товара

```typescript
function ProductPage({ productId }: { productId: string }) {
  const tracking = useBehaviorTracking();

  useEffect(() => {
    // Начинаем отслеживание просмотра
    tracking.startItemView(productId);

    const viewStart = Date.now();

    // Отправляем событие просмотра
    tracking.trackItemViewed({
      item_id: productId,
      item_title: product.title,
      item_category: product.category,
      item_price: product.price,
      previous_page: document.referrer
    });

    // Завершаем отслеживание при unmount
    return () => {
      tracking.endItemView();
    };
  }, [productId]);

  return <ProductDetails product={product} />;
}
```

## Конфигурация

```typescript
interface UseBehaviorTrackingOptions {
  enabled?: boolean; // Включить/выключить трекинг (по умолчанию: true)
  batchSize?: number; // Размер батча для отправки (по умолчанию: 10)
  batchTimeout?: number; // Таймаут отправки в мс (по умолчанию: 5000)
  maxRetries?: number; // Максимальное количество повторов (по умолчанию: 3)
  endpoint?: string; // API endpoint (по умолчанию: '/api/v1/analytics/track')
  userId?: string; // ID пользователя (опционально)
  autoInit?: boolean; // Автоматическая инициализация (по умолчанию: true)
}
```

## API интеграция

Hook отправляет события на endpoint `/api/v1/analytics/track` в следующем формате:

```typescript
interface BehaviorEventBatch {
  events: BehaviorEvent[];
  batch_id: string;
  created_at: string;
}
```

Каждое событие содержит:

- Базовые данные: session_id, timestamp, user_id, device_type, browser, user_agent
- Специфичные данные события: search_query, item_id, purchase_amount и т.д.

## Состояние hook'а

```typescript
interface BehaviorTrackingState {
  isEnabled: boolean; // Включен ли трекинг
  isInitialized: boolean; // Инициализирован ли hook
  pendingEvents: BehaviorEvent[]; // События в очереди отправки
  context: TrackingContext; // Контекст текущей сессии
  lastError?: string; // Последняя ошибка
  retryCount: number; // Количество повторных попыток
}
```

## Автоматические возможности

### Сохранение session_id

Session ID автоматически генерируется и сохраняется в localStorage с поддержкой многопользовательского режима.

### Определение устройства и браузера

Автоматически определяется:

- `device_type`: 'desktop' | 'mobile' | 'tablet'
- `browser`: название и версия браузера
- `user_agent`: полная строка user agent

### Graceful shutdown

При закрытии страницы события отправляются через `navigator.sendBeacon` для максимальной надежности.

### Retry логика

- Автоматические повторные попытки при сетевых ошибках
- Экспоненциальная задержка между попытками
- Настраиваемое количество попыток

## Обработка ошибок

Hook автоматически обрабатывает ошибки и предоставляет информацию через `state.lastError`. Все ошибки логируются в консоль с префиксом `[BehaviorTracking]`.

## Производительность

- **Дебаунс**: Предотвращает спам событий с задержкой 100мс
- **Батчинг**: Минимизирует количество HTTP запросов
- **Lazy loading**: Компоненты загружаются только при необходимости
- **Memory efficient**: Минимальное потребление памяти

## Безопасность

- Не отправляет персональные данные автоматически
- Поддерживает анонимный режим
- Валидирует данные перед отправкой
- Использует HTTPS для всех запросов в продакшене

## Тестирование

Для тестирования hook'а создана страница `/test-behavior-tracking` с примерами всех типов событий.

## Совместимость

- React 18+
- Next.js 13+ (App Router)
- TypeScript 5+
- Современные браузеры с поддержкой ES2020+

## Лицензия

Часть проекта hostel-booking-system
