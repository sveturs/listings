# Инструкции по тестированию Harbor в продакшн-окружении

## Подготовка к тестированию

Следуйте этим инструкциям для проверки работы Harbor в продакшн-окружении на сервере svetu.rs.

### 1. Подготовка файлов

Убедитесь, что следующие файлы скопированы на продакшн-сервер:

- `docker-compose.prod.yml.harbor` - файл настроек для деплоя с использованием Harbor
- `harbor-scripts/test_harbor_connection.sh` - скрипт для тестирования соединения с Harbor

```bash
# На локальной машине
scp /data/hostel-booking-system/docker-compose.prod.yml.harbor user@svetu.rs:/path/to/hostel-booking-system/
scp /data/hostel-booking-system/harbor-scripts/test_harbor_connection.sh user@svetu.rs:/path/to/hostel-booking-system/
```

### 2. Настройка доверия к сертификатам Harbor на продакшн-сервере

На продакшн-сервере необходимо настроить доверие к сертификатам Harbor:

```bash
# На сервере svetu.rs
sudo mkdir -p /etc/docker/certs.d/harbor.svetu.rs
```

Поскольку мы используем реальные сертификаты Let's Encrypt, дополнительной настройки может не потребоваться.

## Тестирование соединения с Harbor

### 1. Запуск скрипта тестирования

На продакшн-сервере выполните скрипт тестирования:

```bash
# На сервере svetu.rs
chmod +x test_harbor_connection.sh
./test_harbor_connection.sh
```

Этот скрипт выполнит следующие проверки:
- Доступность Harbor по HTTPS
- Авторизация в Harbor
- Загрузка тестового образа из Harbor
- Тестовый деплой с использованием образа из Harbor

### 2. Исправление возможных проблем

#### Проблемы с доступом к Harbor

Если Harbor недоступен, проверьте:
- Настройки DNS (dig harbor.svetu.rs)
- Доступность по IP (ping 207.180.197.172)
- Работу Harbor на сервере Harbor

#### Проблемы с сертификатами

Если возникают проблемы с SSL-сертификатами:

```bash
# Копирование сертификата с сервера Harbor
ssh user@harbor.svetu.rs "sudo cat /etc/letsencrypt/live/harbor.svetu.rs/fullchain.pem" | sudo tee /etc/docker/certs.d/harbor.svetu.rs/ca.crt

# Перезапуск Docker
sudo systemctl restart docker
```

#### Проблемы с авторизацией

Если не удается авторизоваться в Harbor:

```bash
# Удаление старой авторизации и повторная попытка
docker logout harbor.svetu.rs
docker login -u admin -p SveTu2025 harbor.svetu.rs
```

## Тестирование продакшн-деплоя с Harbor

### 1. Тестирование загрузки образов

Перед полным деплоем проверьте загрузку основных образов:

```bash
docker pull harbor.svetu.rs/svetu/db/postgres:15
docker pull harbor.svetu.rs/svetu/opensearch/opensearch:2.11.0
docker pull harbor.svetu.rs/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
docker pull harbor.svetu.rs/svetu/tools/migrate:latest
docker pull harbor.svetu.rs/svetu/backend/api:latest
```

### 2. Тестовый деплой компонентов

Выполните тестовый деплой отдельных компонентов, чтобы убедиться, что они работают с образами из Harbor:

```bash
# Например, тестирование компонента базы данных
docker compose -f docker-compose.prod.yml.harbor up -d db
```

### 3. Полный тестовый деплой

После успешного тестирования отдельных компонентов можно выполнить полный деплой:

```bash
docker compose -f docker-compose.prod.yml.harbor up -d
```

## Переключение на Harbor в продакшн

Когда все тесты пройдены успешно, вы можете переключить продакшн-окружение на использование Harbor:

```bash
# Создание резервной копии текущего docker-compose.prod.yml
cp docker-compose.prod.yml docker-compose.prod.yml.bak

# Копирование новой версии с использованием Harbor
cp docker-compose.prod.yml.harbor docker-compose.prod.yml

# Перезапуск сервисов
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d
```

## Инструкции по обновлению образов в Harbor

Когда вы готовы к обновлению образов в Harbor:

1. Сначала соберите образы локально:
   ```bash
   docker build -t harbor.svetu.rs/svetu/backend/api:latest ./backend
   docker build -t harbor.svetu.rs/svetu/frontend/app:latest ./frontend/hostel-frontend
   ```

2. Отправьте образы в Harbor:
   ```bash
   docker login -u admin -p SveTu2025 harbor.svetu.rs
   docker push harbor.svetu.rs/svetu/backend/api:latest
   docker push harbor.svetu.rs/svetu/frontend/app:latest
   ```

3. Обновите сервисы на продакшн-сервере:
   ```bash
   # На сервере svetu.rs
   docker compose -f docker-compose.prod.yml pull
   docker compose -f docker-compose.prod.yml up -d
   ```

## Полезные команды

### Мониторинг образов в Harbor

Просмотр списка проектов в Harbor:
```bash
curl -u admin:SveTu2025 -X GET "https://harbor.svetu.rs/api/v2.0/projects"
```

Просмотр списка репозиториев в проекте:
```bash
curl -u admin:SveTu2025 -X GET "https://harbor.svetu.rs/api/v2.0/projects/svetu/repositories"
```

### Проверка статуса контейнеров

```bash
docker compose -f docker-compose.prod.yml ps
```

### Просмотр логов контейнеров

```bash
docker compose -f docker-compose.prod.yml logs -f service_name
```