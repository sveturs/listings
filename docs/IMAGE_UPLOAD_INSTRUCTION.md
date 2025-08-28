# 📸 Инструкция по добавлению изображений в систему SveTu

## 🎯 Цель
Добавить реальные изображения из интернета для всех объявлений и товаров в системе маркетплейса.

## 🛠️ Необходимые инструменты
- Python 3 с библиотеками: PIL (Pillow), requests
- MinIO Client (mc) - установлен в Docker контейнере
- Docker для работы с MinIO
- Доступ к интернету для загрузки изображений

## 📋 Пошаговая инструкция

### 1. Проверка текущего состояния

```bash
# Проверить количество объявлений без изображений
docker exec hostel_db psql -U postgres -d svetubd -c "
    SELECT COUNT(*) as total_listings,
           COUNT(DISTINCT mi.listing_id) as listings_with_images,
           COUNT(*) - COUNT(DISTINCT mi.listing_id) as listings_without_images
    FROM marketplace_listings ml
    LEFT JOIN marketplace_images mi ON ml.id = mi.listing_id
    WHERE ml.status = 'active';"

# Проверить какие объявления без изображений
docker exec hostel_db psql -U postgres -d svetubd -c "
    SELECT l.id, l.category_id, l.title
    FROM marketplace_listings l
    LEFT JOIN marketplace_images i ON l.id = i.listing_id
    WHERE i.id IS NULL
    ORDER BY l.id;"
```

### 2. Источники изображений

#### Бесплатные стоковые фото:
- **Unsplash**: `https://source.unsplash.com/random/800x600/?{keyword}`
- **Pexels API**: Требует API ключ, но дает более точные результаты
- **Pixabay**: Бесплатные изображения с API
- **Lorem Picsum**: `https://picsum.photos/800/600` - случайные изображения

#### Прямые ссылки на изображения по категориям:
```python
IMAGE_SOURCES = {
    'apartment': [
        'https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?w=800&h=600',
        'https://images.unsplash.com/photo-1560448204-e02f11c3d0e2?w=800&h=600',
        # ... больше ссылок
    ],
    'car': [
        'https://images.unsplash.com/photo-1555215858-9dc80e68c8e8?w=800&h=600',
        # ... больше ссылок
    ]
    # ... другие категории
}
```

### 3. Скрипт для загрузки изображений

Создать файл `/tmp/add_more_images.py`:

```python
#!/usr/bin/env python3
import subprocess
import os
import requests
import time
import random

def get_listings_without_images():
    """Получить список объявлений без изображений"""
    cmd = """docker exec hostel_db psql -U postgres -d svetubd -t -c "
        SELECT l.id, l.category_id, l.title
        FROM marketplace_listings l
        LEFT JOIN marketplace_images i ON l.id = i.listing_id
        WHERE i.id IS NULL
        ORDER BY l.id;"
    """
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    return result.stdout.strip().split('\n')

def download_image_from_unsplash(keyword, output_path):
    """Загрузить изображение с Unsplash по ключевому слову"""
    url = f"https://source.unsplash.com/random/800x600/?{keyword}"
    try:
        response = requests.get(url, timeout=10)
        if response.status_code == 200:
            with open(output_path, 'wb') as f:
                f.write(response.content)
            return True
    except:
        pass
    return False

def upload_to_minio(local_path, minio_path):
    """Загрузить в MinIO"""
    # Копировать в контейнер
    subprocess.run(f"docker cp {local_path} minio:/tmp/", shell=True)
    filename = os.path.basename(local_path)
    
    # Загрузить через mc
    cmd = f"docker exec minio mc cp /tmp/{filename} myminio/listings/{minio_path}"
    result = subprocess.run(cmd, shell=True, capture_output=True)
    
    # Очистить временный файл
    subprocess.run(f"docker exec minio rm /tmp/{filename}", shell=True)
    return result.returncode == 0

# Основная логика
listings = get_listings_without_images()
for listing in listings:
    # Обработать каждое объявление
    # ... загрузить и добавить изображения
```

### 4. Добавление записей в базу данных

После загрузки изображений в MinIO, нужно добавить записи в таблицу `marketplace_images`:

```sql
-- Добавить записи об изображениях
INSERT INTO marketplace_images (listing_id, file_name, file_path, is_main, storage_type, created_at)
VALUES 
    (184, 'listing_184_main.jpg', 'listings/184/main.jpg', true, 'minio', NOW()),
    (184, 'listing_184_2.jpg', 'listings/184/image2.jpg', false, 'minio', NOW());
```

### 5. Массовое добавление изображений

Для добавления изображений ко всем объявлениям:

```bash
# Запустить скрипт
python3 /tmp/add_more_images.py

# Проверить результат
docker exec minio mc ls myminio/listings/ --recursive | wc -l
```

### 6. Проверка результата

```bash
# Проверить доступность изображений
curl -I http://localhost:9000/listings/184/main.jpg

# Проверить в базе данных
docker exec hostel_db psql -U postgres -d svetubd -c "
    SELECT listing_id, COUNT(*) as image_count
    FROM marketplace_images
    GROUP BY listing_id
    ORDER BY listing_id;"
```

## 🔧 Решение проблем

### Если изображение не загружается:
1. Проверить доступность источника
2. Использовать альтернативный источник
3. Увеличить timeout в requests

### Если MinIO не принимает файлы:
1. Проверить права доступа к bucket
2. Проверить размер файла (не более 10MB)
3. Проверить формат изображения (JPEG/PNG)

## 📝 Примечания

- Всегда добавляйте минимум 1 главное изображение (is_main=true)
- Рекомендуется 3-5 изображений на объявление
- Изображения должны соответствовать категории товара
- Используйте задержку между запросами (time.sleep(0.5)) чтобы не перегружать источники

## 🚀 Быстрый старт для следующей сессии

```bash
# 1. Проверить объявления без фото
docker exec hostel_db psql -U postgres -d svetubd -c "SELECT id, title FROM marketplace_listings WHERE id NOT IN (SELECT DISTINCT listing_id FROM marketplace_images);"

# 2. Запустить скрипт добавления изображений
python3 /tmp/add_more_images.py

# 3. Проверить результат
curl -s http://localhost:3001/ru | grep -c "img"
```