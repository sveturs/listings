# Руководство по тестированию загрузки изображений

Это руководство описывает как самостоятельно протестировать систему загрузки изображений для всех типов сущностей: товаров витрин, объявлений маркетплейса, файлов чатов и фотографий отзывов.

## Содержание
1. [Получение JWT токена](#получение-jwt-токена)
2. [Тестирование товаров витрин](#тестирование-товаров-витрин)
3. [Тестирование объявлений маркетплейса](#тестирование-объявлений-маркетплейса)
4. [Тестирование файлов чатов](#тестирование-файлов-чатов)
5. [Тестирование фотографий отзывов](#тестирование-фотографий-отзывов)
6. [Проверка через браузер](#проверка-через-браузер)
7. [Проверка доступности изображений](#проверка-доступности-изображений)
8. [Автоматизированное тестирование](#автоматизированное-тестирование)
9. [Проверка через базу данных](#проверка-через-базу-данных)
10. [Типичные проблемы и решения](#типичные-проблемы-и-решения)

## Получение JWT токена

**ВАЖНО**: Используйте bash для корректной работы с токеном в ZSH!

### Шаг 1: Получение токена с auth сервера

```bash
# Получить токен и сохранить в файл
ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && sed 's|/data/auth_svetu/keys/private.pem|./keys/private.pem|g' scripts/create_admin_jwt.go > /tmp/create_jwt_fixed.go && go run /tmp/create_jwt_fixed.go" > /tmp/jwt_token.txt

# Проверить что токен получен (должен начинаться с eyJ)
head -c 50 /tmp/jwt_token.txt
```

### Шаг 2: Проверка работы токена

```bash
# Проверить что backend запущен
netstat -tlnp 2>/dev/null | grep :3000

# Если не запущен, запустить:
/home/dim/.local/bin/kill-port-3000.sh
screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'

# Подождать 5 секунд для запуска
sleep 5

# Проверить токен
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me | jq ".data.email"'
```

Должен вернуться email пользователя (например: "voroshilovdo@gmail.com")

## Тестирование товаров витрин

### Подготовка тестового изображения

```bash
# Создать красное тестовое изображение 200x200
convert -size 200x200 xc:red -pointsize 60 -fill white -gravity center -annotate +0+0 "TEST" /tmp/test_red.jpg

# Проверить что файл создан
ls -lh /tmp/test_red.jpg
```

### Создание товара и загрузка изображения

```bash
# 1. Получить slug своей витрины
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/storefronts/my | jq ".data[0].slug"'

# Заменить YOUR_SLUG на полученное значение (например: "shop")
SLUG="YOUR_SLUG"

# 2. Создать тестовый товар
bash -c "TOKEN=\$(cat /tmp/jwt_token.txt); curl -s -X POST \"http://localhost:3000/api/v1/storefronts/slug/\$SLUG/products\" \
  -H \"Authorization: Bearer \$TOKEN\" \
  -H \"Content-Type: application/json\" \
  -d '{
    \"name\": \"Тестовый товар с изображением\",
    \"description\": \"Проверка загрузки изображений в правильный bucket\",
    \"price\": 1000,
    \"currency\": \"RSD\",
    \"category_id\": 10207,
    \"stock_quantity\": 10,
    \"is_active\": true
  }' | jq '.'"

# Сохранить PRODUCT_ID из ответа (поле data.id)

# 3. Загрузить изображение к товару
PRODUCT_ID="YOUR_PRODUCT_ID"  # Заменить на реальный ID

bash -c "TOKEN=\$(cat /tmp/jwt_token.txt); curl -s -X POST \"http://localhost:3000/api/v1/storefronts/slug/\$SLUG/products/\$PRODUCT_ID/images\" \
  -H \"Authorization: Bearer \$TOKEN\" \
  -F \"images=@/tmp/test_red.jpg\" | jq '.'"
```

### Проверка результата

**Что проверить в ответе:**

1. **Поле `data[0].image_url` должно начинаться с** `https://s3.svetu.rs/dimalocal-storefront-products/`
2. **Формат URL:** `https://s3.svetu.rs/dimalocal-storefront-products/{product_id}/{timestamp}_{filename}.jpg`
3. **Поле `data[0].thumbnail_url`** должно иметь такой же формат с суффиксом `_thumb`

**Пример правильного ответа:**

```json
{
  "success": true,
  "data": [
    {
      "id": 123,
      "image_url": "https://s3.svetu.rs/dimalocal-storefront-products/340/1738264850_test_red.jpg",
      "thumbnail_url": "https://s3.svetu.rs/dimalocal-storefront-products/340/1738264850_test_red_thumb.jpg",
      "display_order": 0
    }
  ]
}
```

## Тестирование объявлений маркетплейса

### Создание объявления C2C с изображением

```bash
# 1. Создать объявление
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -X POST "http://localhost:3000/api/v1/marketplace" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Тестовое объявление с фото\",
    \"description\": \"Проверка загрузки изображений\",
    \"price\": 5000,
    \"currency\": \"RSD\",
    \"category_id\": 1301,
    \"location\": {
      \"latitude\": 44.7866,
      \"longitude\": 20.4489,
      \"address\": \"Belgrade, Serbia\"
    }
  }" | jq "."'

# Сохранить LISTING_ID из ответа

# 2. Загрузить изображения
LISTING_ID="YOUR_LISTING_ID"

bash -c "TOKEN=\$(cat /tmp/jwt_token.txt); curl -s -X POST \"http://localhost:3000/api/v1/marketplace/\$LISTING_ID/images\" \
  -H \"Authorization: Bearer \$TOKEN\" \
  -F \"images=@/tmp/test_red.jpg\" | jq '.'"
```

**Ожидаемый формат URL:** `https://s3.svetu.rs/dimalocal-listings/{listing_id}/{timestamp}_{filename}.jpg`

## Тестирование файлов чатов

```bash
# 1. Получить список своих чатов
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/marketplace/chats | jq "."'

# Сохранить CHAT_ID

# 2. Отправить файл в чат
CHAT_ID="YOUR_CHAT_ID"

bash -c "TOKEN=\$(cat /tmp/jwt_token.txt); curl -s -X POST \"http://localhost:3000/api/v1/marketplace/chats/\$CHAT_ID/messages\" \
  -H \"Authorization: Bearer \$TOKEN\" \
  -F \"attachment=@/tmp/test_red.jpg\" \
  -F 'content={\"text\":\"Тестовый файл\"}' | jq '.'"
```

**Ожидаемый формат URL:** `https://s3.svetu.rs/dimalocal-chat-files/{chat_id}/{timestamp}_{filename}.jpg`

## Тестирование фотографий отзывов

```bash
# 1. Найти заказ для создания отзыва
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/orders | jq ".data[] | select(.status == \"delivered\") | {id, storefront_id}"'

# Сохранить ORDER_ID и STOREFRONT_ID

# 2. Создать отзыв с фотографией
ORDER_ID="YOUR_ORDER_ID"
STOREFRONT_ID="YOUR_STOREFRONT_ID"

bash -c "TOKEN=\$(cat /tmp/jwt_token.txt); curl -s -X POST \"http://localhost:3000/api/v1/storefronts/\$STOREFRONT_ID/reviews\" \
  -H \"Authorization: Bearer \$TOKEN\" \
  -F \"order_id=\$ORDER_ID\" \
  -F \"rating=5\" \
  -F \"comment=Отличный товар!\" \
  -F \"images=@/tmp/test_red.jpg\" | jq '.'"
```

**Ожидаемый формат URL:** `https://s3.svetu.rs/dimalocal-review-photos/{review_id}/{timestamp}_{filename}.jpg`

## Проверка через браузер

### 1. Запустить frontend (если не запущен)

```bash
/home/dim/.local/bin/start-frontend-screen.sh
```

### 2. Открыть в браузере

- **Витрины:** http://localhost:3001/ru/storefronts/YOUR_SLUG
- **Маркетплейс:** http://localhost:3001/ru/marketplace
- **Объявление:** http://localhost:3001/ru/marketplace/YOUR_LISTING_ID

### 3. Проверить DevTools

1. Открыть DevTools (F12)
2. Перейти на вкладку Network
3. Обновить страницу
4. Найти запросы к изображениям
5. Проверить что URLs начинаются с `https://s3.svetu.rs/dimalocal-*`
6. Убедиться что статус ответа **200 OK**, а не 400 или 404

## Проверка доступности изображений

### Прямая проверка URL

```bash
# Взять image_url из предыдущих тестов и проверить доступность
IMAGE_URL="https://s3.svetu.rs/dimalocal-storefront-products/340/1738264850_test_red.jpg"

# Проверить доступность
curl -I "$IMAGE_URL"

# Должен вернуться статус 200 OK
# Загрузить изображение
curl -o /tmp/downloaded_image.jpg "$IMAGE_URL"

# Проверить что файл загружен
file /tmp/downloaded_image.jpg
# Должно показать: JPEG image data
```

### Проверка через MinIO Client

```bash
# Установить mc (MinIO Client) если не установлен
# brew install minio/stable/mc  # macOS
# или скачать с https://min.io/docs/minio/linux/reference/minio-mc.html

# Настроить alias для S3
mc alias set svetu https://s3.svetu.rs YOUR_ACCESS_KEY YOUR_SECRET_KEY

# Проверить список файлов в bucket
mc ls svetu/dimalocal-storefront-products/

# Проверить конкретный файл
mc stat svetu/dimalocal-storefront-products/340/1738264850_test_red.jpg
```

## Проверка через базу данных

### Проверка URLs в БД

```bash
# Подключиться к БД
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
```

```sql
-- Проверить последние загруженные изображения товаров витрин
SELECT
    id,
    entity_type,
    entity_id,
    LEFT(image_url, 80) as image_url,
    LEFT(thumbnail_url, 80) as thumbnail_url,
    created_at
FROM images
WHERE entity_type = 'storefront_product'
ORDER BY created_at DESC
LIMIT 10;

-- Проверить что все URLs содержат правильный bucket
SELECT
    entity_type,
    COUNT(*) as total,
    COUNT(CASE WHEN image_url LIKE 'https://s3.svetu.rs/dimalocal-%' THEN 1 END) as correct_format,
    COUNT(CASE WHEN image_url NOT LIKE 'https://s3.svetu.rs/dimalocal-%' THEN 1 END) as wrong_format
FROM images
GROUP BY entity_type;

-- Найти изображения с неправильным форматом URL
SELECT
    id,
    entity_type,
    entity_id,
    image_url
FROM images
WHERE image_url NOT LIKE 'https://s3.svetu.rs/dimalocal-%'
ORDER BY created_at DESC
LIMIT 20;
```

### Ожидаемые паттерны URLs по типам

| Тип сущности | Паттерн URL |
|--------------|-------------|
| storefront_product | `https://s3.svetu.rs/dimalocal-storefront-products/{product_id}/{timestamp}_{filename}.jpg` |
| marketplace_listing | `https://s3.svetu.rs/dimalocal-listings/{listing_id}/{timestamp}_{filename}.jpg` |
| chat_file | `https://s3.svetu.rs/dimalocal-chat-files/{chat_id}/{timestamp}_{filename}.jpg` |
| review_photo | `https://s3.svetu.rs/dimalocal-review-photos/{review_id}/{timestamp}_{filename}.jpg` |

## Типичные проблемы и решения

### Проблема 1: 401 Unauthorized при загрузке

**Причина:** Невалидный или истекший JWT токен

**Решение:**
```bash
# Получить новый токен
ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && sed 's|/data/auth_svetu/keys/private.pem|./keys/private.pem|g' scripts/create_admin_jwt.go > /tmp/create_jwt_fixed.go && go run /tmp/create_jwt_fixed.go" > /tmp/jwt_token.txt
```

### Проблема 2: 403 Access Denied

**Причина:** Пытаетесь загрузить изображение к чужому ресурсу (например, товару другого пользователя)

**Решение:** Используйте свои собственные ресурсы или токен владельца ресурса

### Проблема 3: 400 Bad Request при отображении изображения

**Причина:** URL в базе данных содержит относительный путь или неправильный bucket

**Решение:**
```bash
# Перезагрузить изображение с исправленным кодом
# Или обновить URL в базе данных вручную (НЕ рекомендуется)
```

### Проблема 4: Backend не запущен

**Симптомы:**
```bash
curl: (7) Failed to connect to localhost port 3000: Connection refused
```

**Решение:**
```bash
# Запустить backend
/home/dim/.local/bin/kill-port-3000.sh
screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'

# Подождать запуска и проверить
sleep 5
curl http://localhost:3000/health
```

### Проблема 5: Изображение загрузилось, но не отображается

**Проверить:**
1. URL в базе данных правильный (начинается с https://s3.svetu.rs/dimalocal-)
2. Файл доступен на S3 (curl -I $IMAGE_URL возвращает 200)
3. Frontend использует правильный URL (проверить в Network DevTools)

**Решение:** Проверить настройки CORS в MinIO/S3

## Резюме проверок

После внесения изменений в систему загрузки изображений **обязательно проверьте:**

✅ Новые загрузки товаров витрин используют `dimalocal-storefront-products`
✅ Новые загрузки объявлений используют `dimalocal-listings`
✅ Новые файлы чатов используют `dimalocal-chat-files`
✅ Новые фотографии отзывов используют `dimalocal-review-photos`
✅ Все URLs имеют формат `https://s3.svetu.rs/{bucket}/{path}`
✅ Изображения доступны по прямой ссылке (curl -I возвращает 200)
✅ Frontend корректно отображает изображения без ошибок 400/404
✅ В базе данных нет записей с относительными путями `/storefront-products/...`

## Конфигурация buckets

Все bucket names настраиваются через `.env` файл:

```bash
MINIO_BUCKET_NAME=dimalocal-listings
MINIO_CHAT_BUCKET=dimalocal-chat-files
MINIO_STOREFRONT_BUCKET=dimalocal-storefront-products
MINIO_REVIEW_PHOTOS_BUCKET=dimalocal-review-photos
```

После изменения `.env` необходимо перезапустить backend.