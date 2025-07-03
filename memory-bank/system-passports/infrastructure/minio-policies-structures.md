# 📋 Паспорт MinIO Policies и структур доступа

## 🏷️ Метаданные
- **Назначение:** Политики безопасности и структуры доступа для MinIO хранилища
- **Тип компонента:** Инфраструктура / Security & Access Control
- **Статус:** Активный, используется в production
- **Версия:** IAM Policy v2012-10-17
- **Файлы:** `backend/internal/storage/minio/client.go`

## 🎯 Назначение
MinIO policies и структуры доступа определяют правила безопасности, авторизации и управления доступом к файловому хранилищу, обеспечивая контролируемый публичный доступ для чтения и защищенную загрузку через аутентифицированный API.

## 🔒 Структуры политик доступа

### 1. Политика публичного чтения (Public Read)
**Применяется ко всем buckets**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": ["*"]
      },
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::bucket-name/*"]
    }
  ]
}
```

#### Характеристики
- **Effect:** Allow - разрешающая политика
- **Principal:** `*` - доступ для всех пользователей (анонимный)
- **Action:** `s3:GetObject` - только операции чтения
- **Resource:** `bucket-name/*` - все объекты в bucket

### 2. Приватная политика (для конфиденциальных данных)
**Потенциальная конфигурация для private buckets**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": ["arn:aws:s3:::private-bucket/*"]
    }
  ]
}
```

### 3. Политика для аутентифицированных пользователей
**Расширенные права для авторизованных операций**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": ["arn:aws:iam::account:user/api-service"]},
      "Action": [
        "s3:GetObject",
        "s3:PutObject", 
        "s3:DeleteObject"
      ],
      "Resource": ["arn:aws:s3:::user-content/*"]
    }
  ]
}
```

## 🛠️ Автоматическая установка политик

### Процедура создания и настройки bucket
```go
func (m *MinioClient) ensureBucketExists(bucketName string) error {
    ctx := context.Background()
    
    // Проверяем существование bucket
    exists, err := m.client.BucketExists(ctx, bucketName)
    if err != nil {
        return fmt.Errorf("ошибка проверки bucket: %w", err)
    }
    
    // Создаем bucket если не существует
    if !exists {
        err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
            Region: m.location,
        })
        if err != nil {
            return fmt.Errorf("ошибка создания bucket: %w", err)
        }
        
        // Устанавливаем политику публичного чтения
        policy := fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": ["*"]},
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::%s/*"]
                }
            ]
        }`, bucketName)
        
        err = m.client.SetBucketPolicy(ctx, bucketName, policy)
        if err != nil {
            return fmt.Errorf("ошибка установки политики: %w", err)
        }
    }
    
    return nil
}
```

### Инициализация всех buckets при запуске
```go
func (m *MinioClient) InitializeBuckets() error {
    buckets := []string{
        "listings",      // Изображения объявлений
        "chat-files",    // Файлы чатов
        "review-photos", // Фотографии отзывов
    }
    
    for _, bucket := range buckets {
        if err := m.ensureBucketExists(bucket); err != nil {
            return fmt.Errorf("ошибка инициализации bucket %s: %w", bucket, err)
        }
    }
    
    return nil
}
```

## 🌐 CORS конфигурация для веб-интеграции

### Backend CORS настройки
```go
type CORSConfig struct {
    AllowOrigins: []string{
        "https://svetu.rs",
        "https://www.svetu.rs",
        "http://localhost:3000",  // Frontend dev
        "http://localhost:3001",  // Frontend Turbopack
    },
    AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
    AllowHeaders:     []string{
        "Origin", 
        "Content-Type", 
        "Accept", 
        "Authorization", 
        "X-Requested-With", 
        "X-CSRF-Token",
    },
    AllowCredentials: true,
    MaxAge:          86400, // 24 hours
}
```

### Nginx проксирование файлов
```nginx
# Прямой доступ к изображениям объявлений
location ~ ^/listings/(.+)$ {
    proxy_pass http://minio:9000/listings/$1;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_buffering off;
    expires 7d;
    add_header Cache-Control "public, immutable";
}

# Прокси для файлов чата
location ~ ^/chat-files/(.+)$ {
    proxy_pass http://minio:9000/chat-files/$1;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_buffering off;
    expires 1d;
}

# Прокси для фото отзывов
location ~ ^/review-photos/(.+)$ {
    proxy_pass http://minio:9000/review-photos/$1;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_buffering off;
    expires 7d;
}
```

## 🔐 Временный доступ через Presigned URLs

### Интерфейс для создания подписанных URL
```go
type FileStorageInterface interface {
    // Базовый метод для основного bucket
    GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
    
    // Метод для кастомных buckets
    GetPresignedURLFromCustomBucket(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
}
```

### Реализация для основного bucket
```go
func (m *MinioClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
    // Нормализация имени объекта
    if strings.HasPrefix(objectName, "/") {
        objectName = objectName[1:]
    }
    
    // Создание подписанного URL
    presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, nil)
    if err != nil {
        return "", fmt.Errorf("ошибка создания presigned URL: %w", err)
    }
    
    return presignedURL.String(), nil
}
```

### Реализация для кастомных buckets
```go
func (m *MinioClient) GetPresignedURLFromCustomBucket(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
    presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
    if err != nil {
        return "", fmt.Errorf("ошибка создания presigned URL для bucket %s: %w", bucketName, err)
    }
    
    return presignedURL.String(), nil
}
```

### Примеры использования Presigned URLs
```go
// Временный доступ к приватному файлу на 1 час
privateURL, err := minioClient.GetPresignedURL(ctx, "private/document.pdf", 1*time.Hour)
// Результат: http://localhost:9000/bucket/private/document.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...

// Временный доступ для загрузки в чат на 24 часа
chatUploadURL, err := minioClient.GetPresignedURLFromCustomBucket(ctx, "chat-files", "temp/upload.jpg", 24*time.Hour)
```

## 🛡️ Безопасность и контроль доступа

### Конфигурация аутентификации
```go
type MinioConfig struct {
    Endpoint        string // localhost:9000 или production URL
    AccessKeyID     string // minioadmin или production key
    SecretAccessKey string // Секретный ключ
    UseSSL          bool   // false для dev, true для production
    BucketName      string // Основной bucket (listings)
    Location        string // Регион (eu-central-1)
}
```

### Валидация прав доступа
```go
// Проверка прав при удалении файла
func (s *ChatService) DeleteAttachment(ctx context.Context, attachmentID int, userID int) error {
    // Получаем вложение
    attachment, err := s.storage.GetChatAttachmentByID(ctx, attachmentID)
    if err != nil {
        return err
    }
    
    // Получаем сообщение для проверки прав
    message, err := s.storage.GetMessageByID(ctx, attachment.MessageID)
    if err != nil {
        return err
    }
    
    // Проверяем, что пользователь - автор сообщения
    if message.SenderID != userID {
        return fmt.Errorf("permission denied: user %d cannot delete attachment from message by user %d", 
            userID, message.SenderID)
    }
    
    // Удаляем файл из MinIO
    return s.fileStorage.DeleteFile(ctx, attachment.FilePath)
}
```

### Ограничения файлов и безопасность
```go
type FileUploadLimits struct {
    MaxImageSize         int64    // 10 MB
    MaxVideoSize         int64    // 100 MB
    MaxDocumentSize      int64    // 20 MB
    AllowedImageTypes    []string // ["image/jpeg", "image/png", "image/gif", "image/webp"]
    AllowedVideoTypes    []string // ["video/mp4", "video/webm", "video/quicktime"]
    AllowedDocumentTypes []string // ["application/pdf", "text/plain", "application/msword"]
}

// Проверка безопасности файла
func ValidateFileUpload(file multipart.File, header *multipart.FileHeader, limits FileUploadLimits) error {
    // Проверка размера
    if header.Size > getMaxSizeForType(header.Header.Get("Content-Type"), limits) {
        return fmt.Errorf("file size exceeds limit")
    }
    
    // Проверка MIME типа
    contentType := header.Header.Get("Content-Type")
    if !isAllowedContentType(contentType, limits) {
        return fmt.Errorf("file type not allowed: %s", contentType)
    }
    
    // Санитизация имени файла
    sanitizedName := sanitizeFileName(header.Filename)
    if sanitizedName != header.Filename {
        return fmt.Errorf("filename contains invalid characters")
    }
    
    return nil
}
```

## 🔄 Структуры доступа и архитектура

### Иерархия доступа
```
1. Public Read Access (все пользователи)
   ├── GET запросы к файлам
   ├── Прямые URL через Nginx
   └── CDN кеширование

2. Authenticated Write Access (авторизованные пользователи)
   ├── POST /upload/* endpoints
   ├── DELETE /files/* endpoints
   └── PUT /files/* endpoints (обновление)

3. Administrative Access (администраторы)
   ├── Управление bucket policies
   ├── Мониторинг использования
   └── Массовые операции
```

### Абстракция доступа через интерфейсы
```go
// Универсальный интерфейс файлового хранилища
type FileStorageInterface interface {
    UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
    DeleteFile(ctx context.Context, objectName string) error
    GetURL(ctx context.Context, objectName string) (string, error)
    GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
    GetFile(ctx context.Context, objectName string) (io.ReadCloser, error)
}

// Расширенный интерфейс для множественных buckets
type MultiBucketStorageInterface interface {
    FileStorageInterface
    UploadToCustomBucket(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (string, error)
    DeleteFileFromCustomBucket(ctx context.Context, bucketName, objectName string) error
    GetPresignedURLFromCustomBucket(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
}
```

## 🚀 Оптимизации производительности

### Прямой доступ через Nginx
```nginx
# Минимизация нагрузки на backend
location ~ ^/(?:listings|chat-files|review-photos)/(.+)$ {
    # Прямое проксирование к MinIO
    proxy_pass http://minio:9000$request_uri;
    
    # Оптимизация
    proxy_buffering off;           # Потоковая передача
    proxy_request_buffering off;   # Для больших файлов
    
    # Кеширование
    expires 7d;
    add_header Cache-Control "public, immutable";
    
    # Безопасность
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
}
```

### Публичные URL для максимальной производительности
```go
// Формирование прямых URL без проксирования через backend
func (m *MinioClient) GetPublicURL(objectName string) string {
    protocol := "http"
    if m.useSSL {
        protocol = "https"
    }
    
    // Прямой URL к MinIO
    return fmt.Sprintf("%s://%s/%s/%s", protocol, m.endpoint, m.bucketName, objectName)
}

// Для production с Nginx проксированием
func (m *MinioClient) GetProxiedURL(objectName string) string {
    return fmt.Sprintf("https://svetu.rs/%s/%s", m.bucketName, objectName)
}
```

## 📊 Мониторинг и аудит

### Логирование операций доступа
```go
func (m *MinioClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
    logger := log.Ctx(ctx)
    
    // Логируем начало операции
    logger.Info().
        Str("operation", "upload").
        Str("bucket", m.bucketName).
        Str("object", objectName).
        Int64("size", size).
        Str("content_type", contentType).
        Msg("Starting file upload")
    
    // Выполняем загрузку
    _, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, minio.PutObjectOptions{
        ContentType: contentType,
    })
    
    if err != nil {
        logger.Error().Err(err).Msg("File upload failed")
        return "", err
    }
    
    publicURL := m.GetPublicURL(objectName)
    logger.Info().
        Str("public_url", publicURL).
        Msg("File upload completed successfully")
    
    return publicURL, nil
}
```

### Метрики и алерты
```go
type StorageMetrics struct {
    TotalBuckets       int64
    TotalObjects       int64
    TotalSize          int64
    UploadRate         float64  // uploads per minute
    DownloadRate       float64  // downloads per minute
    ErrorRate          float64  // errors per minute
    AverageFileSize    int64
}

// Мониторинг состояния buckets
func (m *MinioClient) GetBucketMetrics(ctx context.Context) (*StorageMetrics, error) {
    // Сбор метрик по всем buckets
    // Подсчет объектов, размеров, статистики использования
}
```

---
**Паспорт создан:** 2025-06-29  
**Компонент:** MinIO Policies и структуры доступа  
**Статус:** Активный в production