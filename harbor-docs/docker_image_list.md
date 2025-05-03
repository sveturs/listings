# Список образов Docker для Sve Tu Platform

## Основные образы

Эти образы необходимо адаптировать для хранения в Harbor.

### Бэкенд

**Текущий:** 
```
build:
  context: ./backend
  dockerfile: Dockerfile
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/backend/api:latest
```

### Фронтенд

**Текущий:**
```
build:
  context: ./frontend/hostel-frontend
  dockerfile: Dockerfile
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/frontend/app:latest
```

### База данных

**Текущий:**
```
image: postgres:15
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/db/postgres:15
```

### OpenSearch

**Текущий:**
```
image: opensearchproject/opensearch:2.11.0
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/opensearch/opensearch:2.11.0
```

### OpenSearch Dashboard

**Текущий:**
```
image: opensearchproject/opensearch-dashboards:2.11.0
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/opensearch/dashboards:2.11.0
```

### MinIO

**Текущий:**
```
image: minio/minio:RELEASE.2023-09-30T07-02-29Z
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
```

### MinIO Client

**Текущий:**
```
image: minio/mc
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/minio/mc:latest
```

### Nginx

**Текущий:**
```
# Используется в docker-compose.prod.yml
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/nginx/nginx:latest
```

### Миграция

**Текущий:**
```
image: migrate/migrate
```

**Новый в Harbor:**
```
image: 207.180.197.172/svetu/tools/migrate:latest
```

## Процесс миграции образов

Для каждого образа выполните следующие шаги:

1. **Загрузка оригинального образа:**
   ```bash
   docker pull postgres:15
   ```

2. **Переименование образа для Harbor:**
   ```bash
   docker tag postgres:15 207.180.197.172/svetu/db/postgres:15
   ```

3. **Отправка образа в Harbor:**
   ```bash
   docker push 207.180.197.172/svetu/db/postgres:15
   ```

## Создание кастомных образов

### Бэкенд

```bash
cd /path/to/backend
docker build -t 207.180.197.172/svetu/backend/api:latest .
docker push 207.180.197.172/svetu/backend/api:latest
```

### Фронтенд

```bash
cd /path/to/frontend/hostel-frontend
docker build -t 207.180.197.172/svetu/frontend/app:latest .
docker push 207.180.197.172/svetu/frontend/app:latest
```

### Nginx

```bash
cd /path/to/nginx
docker build -t 207.180.197.172/svetu/nginx/nginx:latest .
docker push 207.180.197.172/svetu/nginx/nginx:latest
```

## Автоматизация процесса с помощью скрипта

Скрипт для автоматизации процесса миграции образов:

```bash
#!/bin/bash
set -e

HARBOR_URL="207.180.197.172"
HARBOR_USER="admin"
HARBOR_PASSWORD="YourPassword"

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Список образов для миграции
declare -A images=(
  ["postgres:15"]="svetu/db/postgres:15"
  ["opensearchproject/opensearch:2.11.0"]="svetu/opensearch/opensearch:2.11.0"
  ["opensearchproject/opensearch-dashboards:2.11.0"]="svetu/opensearch/dashboards:2.11.0"
  ["minio/minio:RELEASE.2023-09-30T07-02-29Z"]="svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z"
  ["minio/mc:latest"]="svetu/minio/mc:latest"
  ["migrate/migrate:latest"]="svetu/tools/migrate:latest"
)

# Загрузка, переименование и отправка образов
for src in "${!images[@]}"; do
  dst="${HARBOR_URL}/${images[$src]}"
  echo "Миграция $src -> $dst"
  
  echo "- Загрузка $src..."
  docker pull $src
  
  echo "- Переименование в $dst..."
  docker tag $src $dst
  
  echo "- Отправка в Harbor..."
  docker push $dst
  
  echo "✓ Завершено"
  echo
done

echo "Миграция образов завершена!"
```

Сохраните этот скрипт как `migrate_images.sh` и выполните после настройки Harbor.