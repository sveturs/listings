# 🚀 РУКОВОДСТВО РАЗРАБОТЧИКА
## Интеграция с S3 хранилищем SveTu.rs

---

## 📋 О СИСТЕМЕ

### Что это?
Централизованное S3-совместимое хранилище для всех окружений проекта SveTu.rs - маркетплейса объявлений для Балканского региона (типа OLX/Avito).

### Для чего используется:
- 📸 **Фотографии объявлений** - от частных объявлений о продаже варенья до недвижимости
- 🏪 **Витрины магазинов** - изображения и медиа для B2C сегмента
- 💬 **Файлы из чатов** - документы и изображения, которыми обмениваются пользователи
- 📄 **Документы пользователей** - верификационные документы, договоры

### Характеристики системы:
- **Емкость:** 1.5TB (1TB SSD + 500GB Object Storage)
- **Производительность:** до 500 запросов/сек
- **Максимальный размер файла:** 5TB (рекомендуется до 100MB)
- **Доступность:** 99.9% SLA
- **Пропускная способность:** 1 Gbps
- **Латентность:** <50ms для региона Балкан
- **API:** 100% совместимость с AWS S3
- **CDN Ready:** подготовлено для интеграции с CloudFlare

### Возможности:
✅ Multipart upload для больших файлов  
✅ Presigned URLs для безопасной загрузки  
✅ Публичный доступ для изображений  
✅ Версионирование объектов (по запросу)  
✅ Метаданные объектов  
✅ Поиск по префиксам  
✅ Streaming загрузка/выгрузка  
✅ Batch операции  
✅ События (webhooks) - планируется  

---

## 🔑 CREDENTIALS ДЛЯ ПОДКЛЮЧЕНИЯ

### Production окружение:
```yaml
endpoint: https://s3.svetu.rs
access_key: production_backend
secret_key: xK9mNjR3tP5wQ2aLbV7cH8dS
region: us-east-1  # используйте этот регион для совместимости
use_ssl: true
```

### Development окружение:
```yaml
endpoint: https://s3.svetu.rs
access_key: dev_backend
secret_key: pL4kJ8nM2qR6tY9wX5zC3vB7
region: us-east-1
use_ssl: true
```

### Публичные URL для изображений:
```
Production: https://s3.svetu.rs/production-listings/
Development: https://s3.svetu.rs/development-listings/
```

---

## 📦 СТРУКТУРА БАКЕТОВ

### Production бакеты:
| Бакет | Назначение | Доступ | Lifecycle |
|-------|------------|--------|-----------|
| `production-listings` | Фото объявлений | Публичный | 30 дней → архив |
| `production-storefronts` | Витрины магазинов | Публичный | Постоянно |
| `production-chat-files` | Файлы из чатов | Приватный | 90 дней → удаление |
| `production-user-documents` | Документы | Приватный | Постоянно |

### Development бакеты:
| Бакет | Назначение | Доступ | Lifecycle |
|-------|------------|--------|-----------|
| `development-listings` | Тестовые фото | Публичный | 30 дней → удаление |
| `development-storefronts` | Тест витрины | Публичный | 7 дней → удаление |
| `development-chat-files` | Тест файлы | Приватный | 7 дней → удаление |

---

## 💻 ПРИМЕРЫ КОДА

### Node.js / JavaScript

#### Установка:
```bash
npm install @aws-sdk/client-s3
# или
npm install minio
```

#### AWS SDK v3:
```javascript
import { S3Client, PutObjectCommand, GetObjectCommand } from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";

// Конфигурация клиента
const s3Client = new S3Client({
  endpoint: "http://194.163.132.116:9000",
  region: "us-east-1",
  credentials: {
    accessKeyId: "production_backend",
    secretAccessKey: "xK9mNjR3tP5wQ2aLbV7cH8dS"
  },
  forcePathStyle: true // Важно для MinIO!
});

// Загрузка файла
async function uploadFile(file, listingId) {
  const key = `listings/${listingId}/${Date.now()}-${file.name}`;
  
  const command = new PutObjectCommand({
    Bucket: "production-listings",
    Key: key,
    Body: file.buffer,
    ContentType: file.mimetype,
    Metadata: {
      'listing-id': listingId,
      'uploaded-by': userId,
      'original-name': file.name
    }
  });
  
  await s3Client.send(command);
  
  // Публичный URL
  return `http://194.163.132.116:9000/production-listings/${key}`;
}

// Получение presigned URL для загрузки
async function getUploadUrl(listingId, fileName) {
  const key = `listings/${listingId}/${fileName}`;
  
  const command = new PutObjectCommand({
    Bucket: "production-listings",
    Key: key
  });
  
  // URL действителен 15 минут
  const url = await getSignedUrl(s3Client, command, { expiresIn: 900 });
  return url;
}

// Удаление файла
async function deleteFile(key) {
  const command = new DeleteObjectCommand({
    Bucket: "production-listings",
    Key: key
  });
  
  await s3Client.send(command);
}
```

#### MinIO SDK:
```javascript
const Minio = require('minio');

const minioClient = new Minio.Client({
  endPoint: '194.163.132.116',
  port: 9000,
  useSSL: false,
  accessKey: 'production_backend',
  secretKey: 'xK9mNjR3tP5wQ2aLbV7cH8dS'
});

// Загрузка с прогрессом
function uploadWithProgress(bucket, objectName, filePath) {
  return new Promise((resolve, reject) => {
    const file = fs.createReadStream(filePath);
    const stat = fs.statSync(filePath);
    
    minioClient.putObject(bucket, objectName, file, stat.size, {
      'Content-Type': 'image/jpeg',
      'X-Amz-Meta-Testing': 'metadata'
    }, (err, etag) => {
      if (err) return reject(err);
      resolve(etag);
    });
  });
}

// Получение списка файлов
async function listFiles(bucket, prefix) {
  const objects = [];
  const stream = minioClient.listObjectsV2(bucket, prefix, true);
  
  return new Promise((resolve, reject) => {
    stream.on('data', obj => objects.push(obj));
    stream.on('error', reject);
    stream.on('end', () => resolve(objects));
  });
}
```

### Python

#### Установка:
```bash
pip install boto3
# или
pip install minio
```

#### Boto3:
```python
import boto3
from botocore.client import Config

# Конфигурация
s3 = boto3.client(
    's3',
    endpoint_url='http://194.163.132.116:9000',
    aws_access_key_id='production_backend',
    aws_secret_access_key='xK9mNjR3tP5wQ2aLbV7cH8dS',
    config=Config(signature_version='s3v4'),
    region_name='us-east-1'
)

# Загрузка файла
def upload_listing_image(file_path, listing_id):
    key = f"listings/{listing_id}/{os.path.basename(file_path)}"
    
    s3.upload_file(
        file_path,
        'production-listings',
        key,
        ExtraArgs={
            'ContentType': 'image/jpeg',
            'Metadata': {
                'listing-id': str(listing_id)
            }
        }
    )
    
    return f"http://194.163.132.116:9000/production-listings/{key}"

# Multipart upload для больших файлов
def upload_large_file(file_path, bucket, key):
    # Автоматически использует multipart для файлов > 5MB
    config = boto3.s3.transfer.TransferConfig(
        multipart_threshold=1024 * 25,  # 25MB
        max_concurrency=10,
        multipart_chunksize=1024 * 25,
        use_threads=True
    )
    
    s3.upload_file(
        file_path, bucket, key,
        Config=config,
        Callback=ProgressPercentage(file_path)
    )

# Генерация presigned URL
def create_presigned_url(bucket, key, expiration=3600):
    return s3.generate_presigned_url(
        ClientMethod='get_object',
        Params={'Bucket': bucket, 'Key': key},
        ExpiresIn=expiration
    )
```

### Go

#### Установка:
```bash
go get github.com/minio/minio-go/v7
```

#### Пример:
```go
package main

import (
    "context"
    "log"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
    // Инициализация
    minioClient, err := minio.New("194.163.132.116:9000", &minio.Options{
        Creds:  credentials.NewStaticV4("production_backend", "xK9mNjR3tP5wQ2aLbV7cH8dS", ""),
        Secure: false,
    })
    if err != nil {
        log.Fatalln(err)
    }

    // Загрузка файла
    bucketName := "production-listings"
    objectName := "listings/123/image.jpg"
    filePath := "/path/to/image.jpg"
    contentType := "image/jpeg"

    info, err := minioClient.FPutObject(
        context.Background(),
        bucketName,
        objectName,
        filePath,
        minio.PutObjectOptions{
            ContentType: contentType,
            UserMetadata: map[string]string{
                "listing-id": "123",
            },
        },
    )
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("Uploaded %s of size %d\n", objectName, info.Size)
}
```

### PHP

#### Установка:
```bash
composer require aws/aws-sdk-php
```

#### Пример:
```php
<?php
require 'vendor/autoload.php';

use Aws\S3\S3Client;
use Aws\Exception\AwsException;

// Создание клиента
$s3 = new S3Client([
    'version' => 'latest',
    'region'  => 'us-east-1',
    'endpoint' => 'http://194.163.132.116:9000',
    'use_path_style_endpoint' => true,
    'credentials' => [
        'key'    => 'production_backend',
        'secret' => 'xK9mNjR3tP5wQ2aLbV7cH8dS',
    ],
]);

// Загрузка файла
function uploadListingImage($filePath, $listingId) {
    global $s3;
    
    $key = "listings/{$listingId}/" . basename($filePath);
    
    try {
        $result = $s3->putObject([
            'Bucket' => 'production-listings',
            'Key'    => $key,
            'Body'   => fopen($filePath, 'r'),
            'ContentType' => 'image/jpeg',
            'Metadata' => [
                'listing-id' => $listingId
            ]
        ]);
        
        return "http://194.163.132.116:9000/production-listings/{$key}";
    } catch (AwsException $e) {
        echo $e->getMessage();
        return false;
    }
}

// Получение presigned URL
function getPresignedUrl($bucket, $key, $expiration = '+15 minutes') {
    global $s3;
    
    $cmd = $s3->getCommand('GetObject', [
        'Bucket' => $bucket,
        'Key'    => $key
    ]);
    
    $request = $s3->createPresignedRequest($cmd, $expiration);
    return (string) $request->getUri();
}
?>
```

---

## 🎯 BEST PRACTICES

### 1. Организация файлов
```
✅ ПРАВИЛЬНО:
production-listings/
├── 2025/09/01/listing-12345/image-1.jpg
├── 2025/09/01/listing-12345/image-2.jpg
└── 2025/09/01/listing-12346/image-1.jpg

❌ НЕПРАВИЛЬНО:
production-listings/
├── image1.jpg
├── image2.jpg
└── photo.jpg
```

### 2. Именование файлов
```javascript
// ✅ Хорошо
const key = `listings/${year}/${month}/${day}/${listingId}/${uuid}-${sanitizedName}`;

// ❌ Плохо
const key = `listings/${originalFileName}`;
```

### 3. Оптимизация изображений
```javascript
// Загружайте оптимизированные версии
const sizes = {
  thumbnail: { width: 150, height: 150 },
  medium: { width: 800, height: 600 },
  large: { width: 1920, height: 1080 }
};

// Сохраняйте разные размеры
await uploadImage(thumbnail, `${key}-thumb.jpg`);
await uploadImage(medium, `${key}-medium.jpg`);
await uploadImage(original, `${key}-original.jpg`);
```

### 4. Обработка ошибок
```javascript
async function safeUpload(file, bucket, key, retries = 3) {
  for (let i = 0; i < retries; i++) {
    try {
      return await s3.upload(file, bucket, key);
    } catch (error) {
      if (i === retries - 1) throw error;
      
      // Экспоненциальная задержка
      await new Promise(r => setTimeout(r, Math.pow(2, i) * 1000));
    }
  }
}
```

### 5. Метаданные
```javascript
// Всегда добавляйте метаданные
const metadata = {
  'content-hash': md5Hash,
  'uploaded-at': new Date().toISOString(),
  'user-id': userId,
  'listing-id': listingId,
  'environment': process.env.NODE_ENV
};
```

---

## 📊 ЛИМИТЫ И КВОТЫ

### Размеры:
- **Максимальный размер одного объекта:** 5TB
- **Рекомендуемый размер:** до 100MB
- **Multipart порог:** 5MB (автоматически)
- **Размер части multipart:** 5MB - 5GB

### Лимиты запросов:
- **Максимум запросов/сек:** 500 RPS
- **Максимум подключений:** 1000
- **Timeout загрузки:** 5 минут
- **Presigned URL срок:** максимум 7 дней

### Квоты для окружений:
| Окружение | Квота | Политика при превышении |
|-----------|-------|-------------------------|
| Production | 1TB | Уведомление админа |
| Development | 100GB | Автоочистка старых файлов |
| Staging | 50GB | Блокировка загрузки |

---

## 🔧 ОТЛАДКА И ТЕСТИРОВАНИЕ

### MinIO CLI для тестирования:
```bash
# Установка
curl https://dl.min.io/client/mc/release/linux-amd64/mc -o mc
chmod +x mc

# Настройка
./mc alias set svetu http://194.163.132.116:9000 \
  production_backend xK9mNjR3tP5wQ2aLbV7cH8dS

# Тестовые команды
./mc ls svetu/production-listings
./mc cp test.jpg svetu/production-listings/test.jpg
./mc stat svetu/production-listings/test.jpg
```

### cURL для проверки:
```bash
# Проверка доступности
curl -I http://194.163.132.116:9000/minio/health/live

# Загрузка через presigned URL
curl -X PUT --upload-file image.jpg "PRESIGNED_URL_HERE"
```

### Логирование:
```javascript
// Включите debug логи для отладки
const s3Client = new S3Client({
  // ... config
  logger: console, // или ваш logger
});

// Логируйте все операции
s3Client.middlewareStack.add(
  (next) => async (args) => {
    console.log('S3 Request:', args);
    const result = await next(args);
    console.log('S3 Response:', result);
    return result;
  },
  { step: 'finalizeRequest' }
);
```

---

## ⚠️ ВАЖНЫЕ ЗАМЕЧАНИЯ

### Безопасность:
1. **НИКОГДА** не коммитьте credentials в репозиторий
2. Используйте переменные окружения для хранения ключей
3. Ротируйте ключи каждые 90 дней
4. Используйте presigned URLs для загрузки с клиента

### Производительность:
1. Используйте multipart upload для файлов > 5MB
2. Реализуйте retry логику с exponential backoff
3. Кэшируйте часто используемые объекты
4. Используйте CDN для раздачи статики

### Мониторинг:
1. Логируйте все операции с S3
2. Отслеживайте размер загружаемых файлов
3. Мониторьте количество запросов
4. Алертинг при ошибках загрузки

---

## 🆘 ПОДДЕРЖКА

### При проблемах проверьте:
1. Доступность сервиса: `curl http://194.163.132.116:9000/minio/health/live`
2. Правильность credentials
3. Название бакета и ключа
4. Размер файла не превышает лимиты

### Частые ошибки:

| Ошибка | Причина | Решение |
|--------|---------|---------|
| `AccessDenied` | Неверные credentials | Проверьте access/secret key |
| `NoSuchBucket` | Бакет не существует | Проверьте название бакета |
| `RequestTimeout` | Большой файл | Используйте multipart upload |
| `Connection refused` | Сервис недоступен | Обратитесь к администратору |

### Контакты:
- **Email админа:** admin@svetu.rs
- **Telegram поддержки:** @svetu_dev_support
- **Документация MinIO:** https://min.io/docs/

---

## 📈 ROADMAP

### Ближайшие улучшения:
- ✅ SSL сертификаты (в процессе)
- ⏳ CDN интеграция через CloudFlare
- ⏳ Webhooks для событий
- ⏳ Автоматическая оптимизация изображений
- ⏳ Версионирование объектов
- ⏳ Поиск по метаданным

### Планируемые фичи:
- Object tagging
- Lifecycle policies для всех бакетов
- Репликация в другой регион
- S3 Select для запросов к данным
- Encryption at rest

---

*Версия документа: 1.0*  
*Последнее обновление: Сентябрь 2025*  
*API версия: S3 v4*