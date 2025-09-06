# Проблема с загрузкой изображений витрины

## Контекст
При создании витрины реализована возможность загрузки логотипа и баннера, но изображения не сохраняются корректно.

## Детали реализации

### Frontend

1. **BasicInfoStep.tsx** - добавлены поля для загрузки:
   - Логотип: input file с превью, макс 5MB
   - Баннер: input file с превью, макс 10MB
   - Файлы сохраняются в formData.logoFile и formData.bannerFile

2. **CreateStorefrontContext.tsx** - обновлен submitStorefront:
   ```typescript
   // После создания витрины
   if (formData.logoFile) {
     storefrontApi.uploadLogo(response.id, formData.logoFile)
   }
   if (formData.bannerFile) {
     storefrontApi.uploadBanner(response.id, formData.bannerFile)
   }
   ```

3. **storefrontApi.ts** - методы загрузки:
   ```typescript
   async uploadLogo(storefrontId: number, file: File) {
     const formData = new FormData();
     formData.append('logo', file);
     // POST /api/v1/storefronts/{id}/logo
   }
   ```

### Backend

1. **Endpoints** (module.go):
   - POST /api/v1/storefronts/:id/logo
   - POST /api/v1/storefronts/:id/banner

2. **Handler** ожидает файл в поле 'logo' или 'banner' соответственно

## Проблема

В логах видно:
```
10:31PM INF REQUEST method=POST path=/api/v1/storefronts/19/logo
10:31PM INF RESPONSE duration=0.949309 method=POST path=/api/v1/storefronts/19/logo status=200
10:31PM ERR Error in handler error="Cannot POST /api/v1/storefronts/19/logo"
```

Статус 200, но потом ошибка "Cannot POST".

## Что проверить в следующей сессии

1. **Console.log в браузере** - проверить, что файлы действительно выбраны
2. **Network tab** - посмотреть заголовки и тело запроса
3. **Backend логи** - добавить больше логирования в UploadLogo/UploadBanner
4. **Проверить MinIO** - работает ли загрузка файлов вообще
5. **CORS/Middleware** - возможно, проблема с multipart/form-data

## Временное решение

Можно реализовать загрузку изображений через отдельную страницу редактирования витрины после создания.