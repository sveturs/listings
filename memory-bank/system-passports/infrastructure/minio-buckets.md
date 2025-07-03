# 📋 Паспорт MinIO Buckets

## 🏷️ Метаданные
- **Назначение:** S3-совместимое хранилище файлов и изображений
- **Тип компонента:** Инфраструктура / Object Storage
- **Статус:** Активный, используется в production
- **Версия MinIO:** Latest
- **Порт:** 9000 (API), 9001 (Console)

## 🎯 Назначение
MinIO buckets обеспечивают надежное хранение файлов различных типов для платформы Sve Tu: изображения объявлений, файлы чатов, фотографии отзывов с автоматическим управлением доступом и оптимизированной структурой организации.

## 🔧 Конфигурация подключения

### Переменные окружения
```bash
MINIO_ENDPOINT=localhost:9000         # Адрес MinIO сервера
MINIO_ACCESS_KEY=minioadmin           # Ключ доступа
MINIO_SECRET_KEY=1321321321321       # Секретный ключ
MINIO_USE_SSL=false                   # Использование SSL
MINIO_BUCKET_NAME=listings            # Основной bucket
MINIO_LOCATION=eu-central-1           # Регион
```

### Структура Config
```go
type MinIOConfig struct {
    Endpoint   string // Адрес MinIO сервера
    AccessKey  string // Ключ доступа
    SecretKey  string // Секретный ключ
    UseSSL     bool   // Использование SSL
    BucketName string // Основной bucket
    Location   string // Регион размещения
}
```

## 📦 Структура buckets

### 1. Bucket: `listings` (основной)
**Назначение:** Хранение изображений объявлений маркетплейса

#### Структура папок
```
listings/
├── listing_<id>_<timestamp>_<original_name>    # Новый формат
└── <filename>                                  # Legacy файлы
```

#### Примеры файлов
```
listing_123_1640995200_photo.jpg
listing_456_1640995201_car_image.png
listing_789_1640995202_apartment_view.webp
```

#### Настройки доступа
- **Политика:** Публичное чтение (`s3:GetObject`)
- **Загрузка:** Только через аутентифицированный API
- **URL доступа:** `http://localhost:9000/listings/<filename>`

### 2. Bucket: `chat-files`
**Назначение:** Файлы и вложения в чатах между пользователями

#### Структура папок
```
chat-files/
├── images/
│   └── YYYY/MM/DD/
│       └── <messageID>_<timestamp>_<filename>
├── videos/
│   └── YYYY/MM/DD/
│       └── <messageID>_<timestamp>_<filename>
├── documents/
│   └── YYYY/MM/DD/
│       └── <messageID>_<timestamp>_<filename>
└── temp/
    └── temp_<userID>_<timestamp>_<filename>
```

#### Примеры файлов
```
images/2024/12/29/456_1640995200_screenshot.jpg
videos/2024/12/29/789_1640995201_product_demo.mp4
documents/2024/12/29/321_1640995202_contract.pdf
temp/temp_123_1640995203_uploading.jpg
```

#### Особенности
- **Организация по датам** - автоматическое создание папок по дням
- **Временные файлы** - в папке `temp/` до привязки к сообщению
- **Тайм-код в имени** - для уникальности и сортировки

### 3. Bucket: `review-photos`
**Назначение:** Фотографии в отзывах пользователей

#### Структура папок
```
review-photos/
├── reviews/
│   └── review_<reviewID>_<timestamp>_<filename>
└── temp/
    └── temp_<userID>_<timestamp>_<filename>
```

#### Примеры файлов
```
reviews/review_789_1640995200_product_photo.jpg
reviews/review_456_1640995201_quality_image.png
temp/temp_123_1640995202_upload_pending.webp
```

#### Особенности
- **Связь с отзывами** - ID отзыва в имени файла
- **Временное хранение** - до подтверждения отзыва
- **Высокое качество** - для демонстрации продуктов

## 🔧 API методы работы с файлами

### MinioClient (основной интерфейс)
```go
// Работа с основным bucket
UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
DeleteFile(ctx context.Context, objectName string) error
GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
GetObject(ctx context.Context, objectName string) (io.ReadCloser, error)

// Работа с кастомными buckets
UploadToCustomBucket(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (string, error)
DeleteFileFromCustomBucket(ctx context.Context, bucketName, objectName string) error
GetPresignedURLFromCustomBucket(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
GetObjectFromCustomBucket(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
```

### FileStorageInterface (универсальный)
```go
UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
DeleteFile(ctx context.Context, objectName string) error
GetURL(ctx context.Context, objectName string) (string, error)
GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
GetFile(ctx context.Context, objectName string) (io.ReadCloser, error)
```

### Специализированные обертки
```go
// ChatFilesWrapper - для файлов чата
type ChatFilesWrapper struct {
    client MinioClient
}

// chatFileStorageWrapper - внутренняя обертка
type chatFileStorageWrapper struct {
    storage storage.FileStorage
}
```

## 📄 Ограничения и валидация файлов

### Размеры файлов
```go
MaxImageSize:    10 MB     // Максимальный размер изображения
MaxVideoSize:    100 MB    // Максимальный размер видео
MaxDocumentSize: 20 MB     // Максимальный размер документа
```

### Разрешенные типы
```go
// Изображения
AllowedImageTypes: [
    "image/jpeg", "image/jpg", "image/png", 
    "image/gif", "image/webp", "image/svg+xml"
]

// Видео
AllowedVideoTypes: [
    "video/mp4", "video/mpeg", "video/quicktime", 
    "video/webm", "video/x-msvideo"
]

// Документы
AllowedDocumentTypes: [
    "application/pdf", "application/msword",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/plain", "application/rtf"
]
```

### Процесс валидации
1. **Проверка размера** - до и после загрузки
2. **Валидация MIME типа** - по заголовкам HTTP
3. **Проверка расширения** - сопоставление с типом
4. **Санитизация имени** - очистка специальных символов

## 🔗 Формирование URL и доступ

### Локальная разработка
```bash
# Прямой доступ к MinIO
http://localhost:9000/listings/listing_123_photo.jpg

# Через backend API (проксирование)
http://localhost:3000/listings/listing_123_photo.jpg
```

### Production (через Nginx)
```bash
# Проксирование MinIO через Nginx
https://svetu.rs/listings/listing_123_photo.jpg
https://svetu.rs/chat-files/images/2024/12/29/456_photo.jpg
https://svetu.rs/review-photos/reviews/review_789_photo.jpg
```

### Подписанные URL (временный доступ)
```go
// Создание подписанного URL на 24 часа
signedURL, err := minioClient.GetPresignedURL(ctx, "private_file.jpg", 24*time.Hour)
// Результат: http://localhost:9000/bucket/file.jpg?X-Amz-Algorithm=...
```

## 🚀 Автоматическая инициализация

### Создание buckets при запуске
```go
func (m *MinioClient) EnsureBucketsExist(ctx context.Context) error {
    buckets := []string{"listings", "chat-files", "review-photos"}
    
    for _, bucket := range buckets {
        exists, err := m.client.BucketExists(ctx, bucket)
        if !exists {
            err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
                Region: m.location,
            })
            // Установка публичной политики для чтения
            err = m.setBucketPolicy(ctx, bucket)
        }
    }
    return nil
}
```

### Политики доступа
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": ["*"]},
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::bucket-name/*"]
    }
  ]
}
```

## 🔍 Примеры использования

### Загрузка изображения объявления
```go
// Генерация имени файла
objectName := fmt.Sprintf("listing_%d_%d_%s", listingID, time.Now().Unix(), originalName)

// Загрузка в bucket listings
publicURL, err := minioClient.UploadFile(ctx, objectName, file, fileSize, "image/jpeg")
// Результат: http://localhost:9000/listings/listing_123_1640995200_photo.jpg
```

### Загрузка файла чата
```go
// Путь с организацией по датам
now := time.Now()
objectName := fmt.Sprintf("images/%d/%02d/%02d/%d_%d_%s", 
    now.Year(), now.Month(), now.Day(), messageID, now.Unix(), fileName)

// Загрузка в bucket chat-files
publicURL, err := minioClient.UploadToCustomBucket(ctx, "chat-files", objectName, file, fileSize, contentType)
// Результат: http://localhost:3000/chat-files/images/2024/12/29/456_1640995200_photo.jpg
```

### Загрузка фото отзыва
```go
// Временная загрузка
tempName := fmt.Sprintf("temp/temp_%d_%d_%s", userID, time.Now().Unix(), fileName)
tempURL, err := minioClient.UploadToCustomBucket(ctx, "review-photos", tempName, file, fileSize, "image/jpeg")

// После подтверждения отзыва - перемещение
finalName := fmt.Sprintf("reviews/review_%d_%d_%s", reviewID, time.Now().Unix(), fileName)
// Копирование из temp в reviews и удаление temp файла
```

## 🔧 Интеграция с другими компонентами

### Frontend интеграция
```typescript
// Получение URL изображений через API
const listingImages = listing.images.map(img => 
  `${API_BASE_URL}/listings/${img.file_path}`
);

// Отображение в компонентах
<OptimizedImage 
  src={imageUrl} 
  alt={listing.title}
  sizes="(max-width: 768px) 100vw, 50vw"
/>
```

### Backend API endpoints
```go
// GET /listings/{filename} - проксирование файлов listings
// GET /chat-files/{path} - проксирование файлов чата  
// GET /review-photos/{path} - проксирование фото отзывов
// POST /upload/listing-image - загрузка изображения объявления
// POST /upload/chat-file - загрузка файла чата
// POST /upload/review-photo - загрузка фото отзыва
```

### Очистка временных файлов
```go
// Cron задача для удаления старых temp файлов
func CleanupTempFiles(ctx context.Context, olderThan time.Duration) error {
    // Удаление файлов в temp/ папках старше указанного времени
}
```

## 🚨 Особенности и ограничения

### Безопасность
1. **Публичное чтение** - все файлы доступны по прямым URL
2. **Аутентифицированная запись** - загрузка только через API
3. **Валидация на backend** - проверка прав и типов файлов
4. **Санитизация имен** - предотвращение path traversal

### Производительность
1. **CDN интеграция** - возможность подключения CloudFront/CloudFlare
2. **Nginx кеширование** - статическое кеширование файлов
3. **Оптимизация изображений** - автоматическое сжатие WebP
4. **Lazy loading** - отложенная загрузка на frontend

### Мониторинг
1. **Метрики хранилища** - размер buckets, количество файлов
2. **Логирование операций** - все upload/delete операции
3. **Контроль доступа** - аудит обращений к файлам
4. **Бэкапы** - регулярное резервное копирование

---
**Паспорт создан:** 2025-06-29  
**Компонент:** MinIO Buckets  
**Статус:** Активный в production