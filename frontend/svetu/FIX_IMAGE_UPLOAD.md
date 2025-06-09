# Исправление загрузки изображений

## Проблема

При создании объявления фотографии не загружались с ошибкой:

```
Ошибка получения MultipartForm: request Content-Type has bad boundary or is not multipart/form-data
```

## Причина

1. В `listings.ts` использовался метод `apiClient.post()`, который вызывал `JSON.stringify()` на FormData
2. В `api-client.ts` в методе `request()` всегда устанавливался `Content-Type: application/json`, даже для FormData

## Решение

### 1. Изменения в `src/services/listings.ts`

```typescript
// Было:
const response = await apiClient.post<UploadImagesResponse>(
  `/api/v1/marketplace/listings/${listingId}/images`,
  formData
);

// Стало:
const response = await apiClient.upload<UploadImagesResponse>(
  `/api/v1/marketplace/listings/${listingId}/images`,
  formData
);
```

### 2. Изменения в `src/services/api-client.ts`

```typescript
// Было:
const headers: Record<string, string> = {
  'Content-Type': 'application/json',
  ...(options?.headers as Record<string, string>),
};

// Стало:
const headers: Record<string, string> = {
  ...(options?.headers as Record<string, string>),
};

// Устанавливаем Content-Type только если это не FormData
if (!(options?.body instanceof FormData)) {
  headers['Content-Type'] = 'application/json';
}
```

## Результат

Теперь при загрузке FormData:

- Не вызывается JSON.stringify()
- Не устанавливается Content-Type header
- Браузер автоматически устанавливает правильный Content-Type с boundary для multipart/form-data

Загрузка изображений работает корректно!
