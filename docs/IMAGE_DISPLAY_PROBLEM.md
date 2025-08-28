# Проблемы с изображениями и AI анализом

## Дата обнаружения
2025-08-28

## Проблема 1: Claude API возвращает ошибку 401

### Симптомы
При загрузке фотографии для создания объявления с AI-анализом:
- Фото успешно загружается
- EXIF данные извлекаются корректно (GPS координаты, время съёмки)
- Геокодирование работает (определяется адрес: Васе Стајића 18, Нови Сад)
- НО: Claude API возвращает ошибку 401 (Unauthorized)

### Технические детали из консоли
```
claude.service.ts:75 Claude API error: 401 {"error":"Claude API error: 401"}
```

### Причина
API ключ Claude невалидный или истёк срок действия.

### Текущая конфигурация
1. **Backend** (.env):
   ```
   CLAUDE_API_KEY=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
   ```

2. **Frontend** (.env.local):
   ```
   NEXT_PUBLIC_CLAUDE_API_KEY=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
   ```

### Решение
1. Получить новый валидный API ключ от Anthropic: https://console.anthropic.com/
2. Обновить ключ в обоих файлах конфигурации
3. Перезапустить backend и frontend сервисы

## Проблема 2: Отображение изображений на детальных страницах

### Симптомы
1. **На детальной странице объявления** (например, http://localhost:3001/ru/marketplace/250):
   - В главной галерее изображений показывается плейсхолдер вместо реальных фото
   - В секции "Похожие объявления" внизу страницы изображения отображаются корректно

2. **На главной странице маркетплейса**:
   - После переиндексации OpenSearch изображения отображаются корректно
   - Миниатюры объявлений показывают реальные фото

### Техническая причина

Проблема связана с разными источниками данных:

1. **Главная страница и похожие объявления** получают данные из **OpenSearch**:
   - После переиндексации OpenSearch содержит URL изображений
   - Данные включают поле `imageUrl` с корректными путями к MinIO

2. **Детальная страница объявления** использует **Server Side Rendering (SSR)**:
   - Получает данные напрямую из backend API endpoint `/api/v1/marketplace/listings/{id}`
   - Backend API **НЕ возвращает** информацию об изображениях в ответе
   - Frontend получает объект без поля `images` или `imageUrl`

## Текущее состояние

### ✅ Что работает:
- Все 33 объявления имеют от 2 до 5 реальных изображений
- В MinIO загружено 166 файлов (128 для listings, 38 для products)
- В базе данных есть 127 записей в таблице `marketplace_images`
- OpenSearch проиндексирован и содержит URL изображений
- Изображения доступны по прямым ссылкам (например, http://localhost:9000/listings/250/main.jpg)
- Главная страница показывает изображения после переиндексации
- Секция "Похожие объявления" показывает изображения

### ❌ Что НЕ работает:
- Детальные страницы объявлений показывают плейсхолдеры вместо реальных изображений
- Backend API не включает данные об изображениях в ответ

## Анализ кода

### Backend endpoint (требует исправления)
Файл: `backend/internal/proj/marketplace/handler/listings.go`

Метод `GetListingByID` получает данные из базы через сервис, но:
1. Сервис `GetListingByID` не загружает связанные изображения
2. В ответе API отсутствует поле с изображениями

### Frontend компонент
Файл: `frontend/svetu/src/app/[locale]/marketplace/[id]/page.tsx`

Страница использует SSR и получает данные через:
```typescript
const response = await fetch(`${API_URL}/marketplace/listings/${params.id}`);
const listing = await response.json();
```

Полученный объект `listing` не содержит информации об изображениях.

## Решение

### Вариант 1: Исправить backend API (рекомендуется)
1. Модифицировать метод `GetListingByID` в сервисном слое для загрузки изображений
2. Добавить JOIN с таблицей `marketplace_images` в SQL запрос
3. Включить массив изображений в ответ API

### Вариант 2: Дополнительный запрос на frontend
1. После получения основных данных объявления
2. Сделать отдельный запрос для получения изображений
3. Объединить данные на стороне frontend

### Вариант 3: Использовать данные из OpenSearch
1. Изменить источник данных для детальной страницы
2. Получать данные из OpenSearch вместо прямого API
3. Это обеспечит консистентность с главной страницей

## Временное решение (workaround)

Пока backend не исправлен, можно:
1. Использовать клиентскую загрузку изображений после рендеринга страницы
2. Сделать отдельный API endpoint для получения изображений по listing_id
3. Загружать изображения через useEffect после монтирования компонента

## Файлы для проверки и исправления

### Backend:
- `/backend/internal/proj/marketplace/service/marketplace.go` - метод GetListingByID
- `/backend/internal/proj/marketplace/storage/postgres/marketplace.go` - SQL запросы
- `/backend/internal/proj/marketplace/handler/listings.go` - HTTP handler

### Frontend:
- `/frontend/svetu/src/app/[locale]/marketplace/[id]/page.tsx` - детальная страница
- `/frontend/svetu/src/components/marketplace/ListingImageGallery.tsx` - компонент галереи

## Команды для диагностики

```bash
# Проверить наличие изображений в БД для конкретного объявления
docker exec hostel_db psql -U postgres -d svetubd -c "
SELECT * FROM marketplace_images WHERE listing_id = 250;
"

# Проверить доступность изображения в MinIO
curl -I http://localhost:9000/listings/250/main.jpg

# Проверить ответ backend API
curl http://localhost:3000/api/v1/marketplace/listings/250 | jq .

# Проверить данные в OpenSearch
curl -X GET "http://localhost:9200/marketplace/_doc/250" | jq .
```

## Приоритет исправления
**ВЫСОКИЙ** - изображения критически важны для пользовательского опыта на детальных страницах объявлений.

## Что работает правильно (из логов консоли)

### ✅ EXIF обработка фотографий
- Успешное извлечение GPS координат: `{latitude: 45.25128555277778, longitude: 19.84362411472222}`
- Определение времени съёмки: `Tue Jul 29 2025 13:44:55`
- Распознавание камеры: `HUAWEI ELS-N39`
- Извлечение всех 65 EXIF параметров

### ✅ Геокодирование
- Автоматическое определение адреса по GPS: `Васе Стајића 18, Нови Сад 21101`
- Правильная работа с MapBox API
- Корректное сохранение адреса в форме

### ✅ WebSocket и чат
- WebSocket соединение установлено
- Статусы пользователей онлайн работают
- Система чатов функционирует

### ✅ Корзина и авторизация
- Сервис корзины работает
- Профиль пользователя загружается
- Система авторизации функционирует

## Связанные документы
- `/data/hostel-booking-system/docs/IMAGE_UPLOAD_INSTRUCTION.md` - инструкция по добавлению изображений
- `/tmp/images_addition_final_report.txt` - финальный отчёт о добавлении изображений
- `/data/hostel-booking-system/backend/scripts/create_test_jwt.go` - генератор JWT токенов для тестирования
- `/data/hostel-booking-system/backend/scripts/create_admin_jwt.go` - генератор JWT токенов для администратора