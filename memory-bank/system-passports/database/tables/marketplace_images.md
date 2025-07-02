# Паспорт таблицы: marketplace_images

## Назначение
Хранение информации об изображениях объявлений. Поддерживает множественные изображения для каждого объявления с интеграцией MinIO для физического хранения файлов.

## Структура таблицы

```sql
CREATE TABLE marketplace_images (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    storage_type VARCHAR(20) DEFAULT 'local',
    storage_bucket VARCHAR(100),
    public_url TEXT
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор изображения (SERIAL)
- `listing_id` - привязка к объявлению (FK с CASCADE DELETE)

### Информация о файле
- `file_path` - путь к файлу в хранилище (до 255 символов)
- `file_name` - оригинальное имя файла (до 255 символов)
- `file_size` - размер файла в байтах
- `content_type` - MIME тип файла (например, 'image/jpeg')

### Параметры отображения
- `is_main` - является ли главным изображением (по умолчанию false)

### Хранилище
- `storage_type` - тип хранилища: 'local', 'minio' (по умолчанию 'local')
- `storage_bucket` - название bucket в MinIO
- `public_url` - публичный URL для доступа к изображению

### Системные поля
- `created_at` - дата загрузки изображения

## Индексы

1. **idx_marketplace_images_listing_main** - составной индекс по listing_id и is_main для быстрого поиска главного изображения

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- `listing_id` → `marketplace_listings.id` - объявление (CASCADE DELETE)

### Обратные связи
- Нет прямых ссылок из других таблиц

## Бизнес-правила

### Ограничения
1. **Каскадное удаление** - при удалении объявления автоматически удаляются все изображения
2. **Одно главное изображение** - только одно изображение может быть is_main = true для объявления
3. **Обязательные поля** - file_path, file_name, file_size, content_type обязательны

### Поддерживаемые форматы
- `image/jpeg` - JPEG изображения
- `image/png` - PNG изображения
- `image/webp` - WebP изображения
- `image/gif` - GIF изображения

### Ограничения размера
- Максимальный размер файла: 10 MB (проверяется на уровне приложения)
- Рекомендуемые размеры: 1200x900 px

## Примеры использования

### Добавление изображения к объявлению
```sql
INSERT INTO marketplace_images (
    listing_id, file_path, file_name, file_size, 
    content_type, is_main, storage_type, storage_bucket, public_url
) VALUES (
    123, 
    '/listings/2024/01/uuid-filename.jpg',
    'product-photo.jpg',
    524288,
    'image/jpeg',
    true,
    'minio',
    'listings',
    'https://cdn.svetu.rs/listings/2024/01/uuid-filename.jpg'
);
```

### Получение всех изображений объявления
```sql
SELECT * FROM marketplace_images 
WHERE listing_id = 123 
ORDER BY is_main DESC, created_at ASC;
```

### Установка главного изображения
```sql
-- Сначала сбрасываем текущее главное
UPDATE marketplace_images 
SET is_main = false 
WHERE listing_id = 123;

-- Затем устанавливаем новое
UPDATE marketplace_images 
SET is_main = true 
WHERE id = 456 AND listing_id = 123;
```

### Получение главных изображений для списка объявлений
```sql
SELECT l.id, l.title, i.public_url
FROM marketplace_listings l
LEFT JOIN marketplace_images i ON l.id = i.listing_id AND i.is_main = true
WHERE l.status = 'active';
```

## Структура хранения в MinIO

### Путь к файлам
```
/listings/
  /2024/
    /01/
      /15/
        /{uuid}-{timestamp}.{extension}
```

### Naming convention
- UUID для уникальности
- Timestamp для сортировки
- Оригинальное расширение файла

## Известные особенности

1. **CASCADE DELETE** - автоматическое удаление при удалении объявления
2. **MinIO интеграция** - физические файлы хранятся в MinIO
3. **Public URL** - прямые ссылки для CDN
4. **Множественные изображения** - нет ограничения на количество
5. **Storage migration** - поддержка миграции с local на minio

## Миграции

- **000001** - создание таблицы
- **000025** - добавление storage_type, storage_bucket, public_url
- **000039** - добавление индекса для главного изображения

## Интеграция с другими компонентами

### Backend
1. **Upload handler** - загрузка и валидация изображений
2. **MinIO service** - сохранение файлов
3. **Image processor** - создание миниатюр и оптимизация

### Frontend
1. **Image gallery** - компонент просмотра изображений
2. **Upload widget** - drag&drop загрузка
3. **Image preview** - предпросмотр в списках

### Инфраструктура
1. **MinIO** - объектное хранилище
2. **CDN** - кеширование и раздача
3. **Nginx** - проксирование запросов

## Процесс работы с изображениями

1. **Загрузка**
   - Frontend отправляет файл на `/api/v1/marketplace/listings/{id}/images`
   - Backend валидирует формат и размер
   - Файл сохраняется в MinIO
   - Запись создается в БД с public_url

2. **Отображение**
   - Frontend запрашивает public_url
   - CDN кеширует изображение
   - Пользователь получает оптимизированную версию

3. **Удаление**
   - При удалении записи из БД
   - Асинхронно удаляется файл из MinIO
   - Очищается кеш CDN