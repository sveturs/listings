# Инструкции по деплою с Harbor

## Локальная разработка

Для локальной разработки используйте стандартный docker-compose.yml, который сохранен в корне проекта. Этот файл настроен на локальную сборку образов с использованием директив `build`:

```bash
# Запуск локальной среды разработки
docker-compose up -d
```

Это запустит все службы с локальной сборкой образов для бэкенда и фронтенда.

## Продакшн деплой с использованием Harbor

Для развертывания в продакшн с использованием Harbor-реестра следуйте этим инструкциям:

### 1. Подготовка образов

Перед деплоем необходимо собрать и отправить образы в Harbor:

```bash
# Сборка бэкенд образа
cd /data/hostel-booking-system
docker build -t 207.180.197.172/svetu/backend/api:latest ./backend

# Сборка фронтенд образа
docker build -t 207.180.197.172/svetu/frontend/app:latest ./frontend/hostel-frontend

# Аутентификация в Harbor
docker login -u admin -p SveTu2025 207.180.197.172

# Отправка образов в Harbor
docker push 207.180.197.172/svetu/backend/api:latest
docker push 207.180.197.172/svetu/frontend/app:latest
```

### 2. Настройка деплоя

На продакшн-сервере используйте специальный файл для деплоя с Harbor:

```bash
# Развертывание с использованием Harbor
docker-compose -f docker-compose.prod.yml.harbor up -d
```

Файл `docker-compose.prod.yml.harbor` уже настроен для использования образов из Harbor вместо локальной сборки.

### 3. Обновление deploy.sh для использования Harbor

Для автоматизации процесса деплоя с использованием Harbor, обновите скрипт deploy.sh:

```bash
#!/bin/bash
set -e

echo "==== Деплой Sve Tu Platform с использованием Harbor ===="

# Резервное копирование данных
echo "Создание резервной копии базы данных..."
# ... (код резервного копирования)

# Аутентификация в Harbor (если нужно)
echo "Аутентификация в Harbor..."
docker login -u admin -p SveTu2025 207.180.197.172

# Загрузка последних образов из Harbor
echo "Получение последних образов из Harbor..."
docker pull 207.180.197.172/svetu/backend/api:latest
docker pull 207.180.197.172/svetu/frontend/app:latest
docker pull 207.180.197.172/svetu/db/postgres:15
docker pull 207.180.197.172/svetu/opensearch/opensearch:2.11.0
docker pull 207.180.197.172/svetu/opensearch/dashboards:2.11.0
docker pull 207.180.197.172/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
docker pull 207.180.197.172/svetu/minio/mc:latest
docker pull 207.180.197.172/svetu/tools/migrate:latest

# Запуск с использованием Harbor образов
echo "Запуск сервисов с использованием образов из Harbor..."
docker-compose -f docker-compose.prod.yml.harbor up -d

echo "Деплой успешно завершен!"
```

### 4. Обновление сертификатов на продакшн-сервере

Для работы с Harbor, на продакшн-сервере необходимо настроить доверие к сертификатам Harbor:

```bash
# На продакшн-сервере
sudo mkdir -p /etc/docker/certs.d/207.180.197.172
sudo scp dima@207.180.197.172:/opt/harbor/certs/server.crt /etc/docker/certs.d/207.180.197.172/ca.crt
sudo systemctl restart docker
```

## Выбор стратегии для разных окружений

- **Разработка**: Используйте стандартный docker-compose.yml с локальной сборкой образов.
- **Тестирование**: Можно использовать Harbor-образы для проверки стабильности перед продакшн.
- **Продакшн**: Используйте docker-compose.prod.yml.harbor с образами из Harbor.

## Преимущества этого подхода

1. **Разделение окружений**: Локальная разработка остается простой и быстрой.
2. **Централизованное хранение образов**: Продакшн-образы хранятся и версионируются в Harbor.
3. **Согласованность продакшн-деплоев**: Все продакшн-серверы используют одинаковые образы из Harbor.
4. **Безопасность**: Образы в Harbor можно сканировать на уязвимости перед деплоем.

## Рекомендуемый рабочий процесс

1. Разработка с использованием локального docker-compose.yml
2. Тестирование изменений локально
3. Сборка и отправка образов в Harbor 
4. Деплой в продакшн с использованием docker-compose.prod.yml.harbor