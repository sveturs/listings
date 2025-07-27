# Listing Edit API Fixes

## Исправленные проблемы

### 1. Ошибка с полями city/country

**Проблема**: Backend ожидал поля `address_city` и `address_country`, а frontend отправлял `city` и `country`.

**Решение**: 
- Обновлен интерфейс `Listing` для поддержки обоих вариантов
- При загрузке данных проверяются оба варианта полей
- При отправке используются правильные имена `address_city` и `address_country`

### 2. Отсутствующий endpoint для изменения порядка изображений

**Проблема**: Frontend пытался вызвать несуществующий endpoint `/api/v1/marketplace/listings/{id}/images/reorder`.

**Решение**:
- Изменение главного изображения теперь происходит локально
- При сохранении объявления отправляется полный массив изображений с обновленными флагами
- Добавлено уведомление пользователю о необходимости сохранить изменения

### 3. Обновление изображений при сохранении

**Решение**: В данные для обновления добавлен массив изображений с правильным порядком (display_order).

## Структура данных

### Поля в базе данных (marketplace_listings):
- `address_city` - город
- `address_country` - страна
- `location` - адрес/местоположение
- `latitude`, `longitude` - координаты
- `show_on_map` - показывать на карте

### Обновленная структура запроса:
```typescript
{
  title: string,
  description: string,
  price: number,
  condition: string,
  address_city: string,      // Правильное имя поля
  address_country: string,   // Правильное имя поля
  location: string,
  latitude?: number,
  longitude?: number,
  show_on_map: boolean,
  category_id: number,
  attributes: Array<{
    attribute_id: number,
    value: any
  }>,
  images: Array<{           // Добавлен массив изображений
    id: number,
    display_order: number,
    is_main: boolean,
    // ... другие поля
  }>
}
```

## Рекомендации для backend

1. Добавить отдельный endpoint для изменения порядка изображений (опционально)
2. Обновить swagger документацию с точными типами для UpdateListingRequest
3. Рассмотреть возможность использования patch вместо put для частичных обновлений