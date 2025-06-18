# Исправление API endpoints системы отзывов

## Дата: 2025-06-17
## Ветка: rev_17_06

## Проблема
После интеграции системы отзывов на frontend обнаружились ошибки:
1. 404 ошибки - неправильные пути API
2. 401 Unauthorized - некоторые endpoints требуют авторизацию

## Анализ
При проверке backend обнаружено:
1. Endpoint статистики имеет путь `/api/v1/entity/{type}/{id}/stats`, а не `/api/v1/reviews/stats/{type}/{id}`
2. Endpoint `can-review` требует JWT авторизацию
3. Есть расхождение между Swagger документацией и реальными путями

## Исправления

### 1. Добавлен базовый путь API
```typescript
// src/services/reviewApi.ts
const API_BASE = config.getApiUrl() + '/api/v1';
```

### 2. Исправлен путь для статистики
```typescript
// Было:
`${API_BASE}/reviews/stats/${entityType}/${entityId}`
// Стало:
`${API_BASE}/entity/${entityType}/${entityId}/stats`
```

### 3. Добавлена проверка авторизации для can-review
```typescript
// src/hooks/useReviews.ts
export const useCanReview = (entityType: string, entityId: number, userId?: number) => {
  return useQuery({
    queryKey: QUERY_KEYS.canReview(entityType, entityId),
    queryFn: () => reviewApi.canReview(entityType, entityId),
    enabled: !!entityType && !!entityId && !!userId, // Only run if user is authenticated
    staleTime: 60 * 1000,
  });
};
```

## Результат
Теперь:
- Статистика отзывов загружается корректно (публичный endpoint)
- Проверка возможности оставить отзыв работает только для авторизованных пользователей
- Все пути API соответствуют реальным endpoints в backend