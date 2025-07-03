# 📋 Паспорт Docker Compose Setup

## 🏷️ Метаданные
- **Назначение:** Контейнеризация и оркестрация сервисов Sve Tu Platform
- **Тип компонента:** Инфраструктура / Container Orchestration
- **Статус:** Активный, используется в development и production
- **Версия Docker Compose:** 3.8
- **Файлы:** `docker-compose.yml`, `docker-compose.prod.yml`

## 🎯 Назначение
Docker Compose обеспечивает контейнеризацию всех компонентов платформы, управление зависимостями между сервисами, изоляцию сред разработки и production, а также упрощенное развертывание и масштабирование системы.

## 📂 Структура Compose файлов

### 1. Development Environment (`docker-compose.yml`)
**Назначение:** Локальная разработка и тестирование

### 2. Production Environment (`docker-compose.prod.yml`)
**Назначение:** Production развертывание с полной инфраструктурой

### 3. Deploy Configuration (`deploy/docker-compose.yml`)
**Назначение:** Упрощенное развертывание

### 4. Frontend Standalone (`frontend/svetu/docker-compose.yml`)
**Назначение:** Изолированная разработка фронтенда

## 🏗️ Архитектура сервисов

### Development Configuration

#### Core Services
```yaml
services:
  # База данных
  db:
    image: postgres:15
    container_name: hostel_db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: hostel_db
    ports:
      - \"5432:5432\"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: [\"CMD\", \"pg_isready\", \"-U\", \"postgres\"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - hostel_network

  # Поисковая система
  opensearch:
    image: opensearchproject/opensearch:2.11.0
    container_name: opensearch
    environment:
      - \"discovery.type=single-node\"
      - \"bootstrap.memory_lock=true\"
      - \"OPENSEARCH_JAVA_OPTS=-Xms1024m -Xmx1024m\"
      - \"DISABLE_SECURITY_PLUGIN=true\"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - opensearch-data:/usr/share/opensearch/data
    ports:
      - \"9200:9200\"
    networks:
      - hostel_network

  # OpenSearch веб-интерфейс
  opensearch-dashboards:
    image: opensearchproject/opensearch-dashboards:2.11.0
    container_name: opensearch-dashboards
    ports:
      - \"5601:5601\"
    environment:
      - \"OPENSEARCH_HOSTS=http://opensearch:9200\"
      - \"DISABLE_SECURITY_DASHBOARDS_PLUGIN=true\"
    networks:
      - hostel_network

  # Объектное хранилище
  minio:
    image: minio/minio:RELEASE.2023-09-30T07-02-29Z
    container_name: minio
    command: server /data --console-address \":9001\"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: 1321321321321
      MINIO_BROWSER_REDIRECT_URL: http://localhost:9001
      MINIO_SERVER_URL: http://localhost:9000
    ports:
      - \"9000:9000\"
      - \"9001:9001\"
    volumes:
      - ./data/minio:/data
    restart: unless-stopped
    healthcheck:
      test: [\"CMD\", \"curl\", \"-f\", \"http://localhost:9000/minio/health/live\"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - hostel_network

  # Инициализация MinIO buckets
  createbuckets:
    image: minio/mc
    container_name: minio_createbuckets
    depends_on:
      - minio
    entrypoint: >\n      /bin/sh -c \"\n      /usr/bin/mc config host add myminio http://minio:9000 minioadmin 1321321321321;\n      /usr/bin/mc mb myminio/listings --ignore-existing;\n      /usr/bin/mc policy download myminio/listings;\n      /usr/bin/mc mb myminio/chat-files --ignore-existing;\n      /usr/bin/mc policy download myminio/chat-files;\n      exit 0;\n      \"\n    networks:\n      - hostel_network\n\n  # Миграции базы данных\n  migrate:\n    image: migrate/migrate\n    container_name: hostel_migrate\n    depends_on:\n      - db\n    volumes:\n      - ./backend/migrations:/migrations\n    command: [\n      \"-path\", \"/migrations\",\n      \"-database\", \"postgres://postgres:password@db:5432/hostel_db?sslmode=disable\",\n      \"up\"\n    ]\n    networks:\n      - hostel_network\n\n  # Backend API (закомментирован в development)\n  # backend:\n  #   build: ./backend\n  #   container_name: hostel_backend\n  #   depends_on:\n  #     - db\n  #     - opensearch\n  #     - minio\n  #   ports:\n  #     - \"3000:3000\"\n  #   environment:\n  #     ENV_FILE: .env\n  #   networks:\n  #     - hostel_network\n\n  # Frontend (закомментирован в development)\n  # frontend:\n  #   build: ./frontend/svetu\n  #   container_name: hostel_frontend\n  #   ports:\n  #     - \"3001:3000\"\n  #   networks:\n  #     - hostel_network\n```\n\n### Production Configuration\n\n#### Full Production Stack\n```yaml\nservices:\n  # База данных с production credentials\n  db:\n    image: harbor.svetu.rs/svetu/db/postgres:15\n    container_name: postgres\n    environment:\n      POSTGRES_USER: postgres\n      POSTGRES_PASSWORD: c9XWc7Cm\n      POSTGRES_DB: hostel_db\n      PGDATA: /var/lib/postgresql/data/pgdata\n    volumes:\n      - db_data:/var/lib/postgresql/data\n    restart: unless-stopped\n    healthcheck:\n      test: [\"CMD\", \"pg_isready\", \"-U\", \"postgres\"]\n      interval: 10s\n      timeout: 5s\n      retries: 5\n    networks:\n      - hostel_network\n\n  # OpenSearch для production\n  opensearch:\n    image: harbor.svetu.rs/svetu/opensearch/opensearch:2.11.0\n    container_name: opensearch\n    environment:\n      - \"discovery.type=single-node\"\n      - \"bootstrap.memory_lock=true\"\n      - \"OPENSEARCH_JAVA_OPTS=-Xms1024m -Xmx1024m\"\n      - \"DISABLE_SECURITY_PLUGIN=true\"\n    ulimits:\n      memlock:\n        soft: -1\n        hard: -1\n    volumes:\n      - /opt/hostel-data/opensearch:/usr/share/opensearch/data\n    restart: unless-stopped\n    networks:\n      - hostel_network\n\n  # Почтовый сервер\n  mailserver:\n    image: harbor.svetu.rs/svetu/mail/server:latest\n    container_name: mailserver\n    hostname: mail.svetu.rs\n    ports:\n      - \"25:25\"\n      - \"587:587\"\n      - \"465:465\"\n      - \"143:143\"\n      - \"993:993\"\n      - \"110:110\"\n      - \"995:995\"\n    volumes:\n      - /opt/hostel-data/maildata:/var/mail\n      - /opt/hostel-data/mailstate:/var/mail-state\n      - /opt/hostel-data/maillogs:/var/log/mail\n      - /opt/hostel-data/certbot/conf:/etc/letsencrypt:ro\n    restart: unless-stopped\n    cap_add:\n      - NET_ADMIN\n    networks:\n      - hostel_network\n\n  # SSL сертификаты\n  certbot:\n    image: harbor.svetu.rs/svetu/tools/certbot:latest\n    container_name: certbot\n    volumes:\n      - /opt/hostel-data/certbot/conf:/etc/letsencrypt\n      - /opt/hostel-data/certbot/www:/var/www/certbot\n    command: /bin/sh -c \"trap exit TERM; while :; do sleep 12h & wait $${!}; certbot renew; done;\"\n    networks:\n      - hostel_network\n\n  # Веб-интерфейс почты (Roundcube)\n  mail-webui:\n    image: harbor.svetu.rs/svetu/mail/webui:latest\n    container_name: mail-webui\n    volumes:\n      - roundcube_data:/var/www/html/temp\n      - roundcube_data:/var/www/html/logs\n    restart: unless-stopped\n    networks:\n      - hostel_network\n\n  # MinIO для production\n  minio:\n    image: harbor.svetu.rs/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z\n    container_name: minio\n    command: server /data --console-address \":9001\"\n    environment:\n      MINIO_ROOT_USER: minioadmin\n      MINIO_ROOT_PASSWORD: 5465465465465\n    ports:\n      - \"9000:9000\"\n      - \"9001:9001\"\n    volumes:\n      - /opt/hostel-data/minio:/data\n    restart: unless-stopped\n    healthcheck:\n      test: [\"CMD\", \"curl\", \"-f\", \"http://localhost:9000/minio/health/live\"]\n      interval: 30s\n      timeout: 10s\n      retries: 3\n    networks:\n      - hostel_network\n\n  # Создание MinIO buckets для production\n  createbuckets:\n    image: harbor.svetu.rs/svetu/minio/mc:latest\n    container_name: minio_createbuckets\n    depends_on:\n      - minio\n    entrypoint: >\n      /bin/sh -c \"\n      /usr/bin/mc config host add myminio http://minio:9000 minioadmin 5465465465465;\n      /usr/bin/mc mb myminio/listings --ignore-existing;\n      /usr/bin/mc policy download myminio/listings;\n      /usr/bin/mc mb myminio/chat-files --ignore-existing;\n      /usr/bin/mc policy download myminio/chat-files;\n      exit 0;\n      \"\n    networks:\n      - hostel_network\n\n  # Backend API для production\n  backend:\n    image: harbor.svetu.rs/svetu/backend/api:latest\n    container_name: backend\n    depends_on:\n      - db\n      - opensearch\n      - minio\n    ports:\n      - \"3000:3000\"\n    environment:\n      APP_MODE: production\n      ENV_FILE: .env\n      WS_ENABLED: true\n      OPENSEARCH_URL: http://opensearch:9200\n      OPENSEARCH_MARKETPLACE_INDEX: marketplace\n      FILE_STORAGE_PROVIDER: minio\n      MINIO_ENDPOINT: minio:9000\n      MINIO_ACCESS_KEY: minioadmin\n      MINIO_SECRET_KEY: 5465465465465\n      MINIO_USE_SSL: false\n      MINIO_BUCKET_NAME: listings\n      MINIO_LOCATION: eu-central-1\n      FILE_STORAGE_PUBLIC_URL: https://svetu.rs\n      POSTGRES_USER: postgres\n      POSTGRES_PASSWORD: c9XWc7Cm\n      POSTGRES_DB: hostel_db\n      DATABASE_URL: postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable\n    restart: unless-stopped\n    healthcheck:\n      test: [\"CMD\", \"curl\", \"-f\", \"http://localhost:3000\"]\n      interval: 10s\n      timeout: 5s\n      retries: 3\n    networks:\n      - hostel_network\n\n  # Nginx reverse proxy\n  nginx:\n    image: harbor.svetu.rs/svetu/nginx/nginx:latest\n    container_name: nginx\n    depends_on:\n      - backend\n    ports:\n      - \"80:80\"\n      - \"443:443\"\n    volumes:\n      - /opt/hostel-data/certbot/conf:/etc/letsencrypt:ro\n      - /opt/hostel-data/certbot/www:/var/www/certbot:ro\n      - uploads_data:/usr/share/nginx/uploads:ro\n    restart: unless-stopped\n    healthcheck:\n      test: [\"CMD\", \"wget\", \"--spider\", \"--quiet\", \"http://localhost/\"]\n      interval: 30s\n      timeout: 10s\n      retries: 3\n      start_period: 15s\n    networks:\n      - hostel_network\n```\n\n## 🌐 Сетевая архитектура\n\n### Основная сеть\n```yaml\nnetworks:\n  hostel_network:\n    driver: bridge\n```\n\n### Межсервисное взаимодействие\n| Сервис | Внутренний адрес | Внешний порт | Назначение |\n|--------|------------------|--------------|------------|\n| PostgreSQL | `db:5432` | 5432 (dev) | База данных |\n| OpenSearch | `opensearch:9200` | 9200 (dev) | Поисковый движок |\n| OpenSearch Dashboards | `opensearch-dashboards:5601` | 5601 (dev) | Веб-интерфейс поиска |\n| MinIO API | `minio:9000` | 9000 | Объектное хранилище |\n| MinIO Console | `minio:9001` | 9001 | Веб-интерфейс MinIO |\n| Backend API | `backend:3000` | 3000 (prod) | REST API |\n| Mail Server | `mailserver` | 25,587,465,143,993,110,995 | Почтовые протоколы |\n| Roundcube | `mail-webui:80` | - | Веб-почта |\n| Nginx | `nginx` | 80,443 | Reverse proxy |\n\n## 💾 Управление данными\n\n### Development Volumes\n```yaml\nvolumes:\n  postgres_data:\n    # PostgreSQL данные в Docker volume\n  opensearch-data:\n    # OpenSearch индексы в Docker volume\n  # MinIO использует bind mount\n  ./data/minio:/data\n```\n\n### Production Volumes\n```yaml\nvolumes:\n  db_data:\n    # PostgreSQL данные в Docker volume\n  roundcube_data:\n    # Roundcube временные файлы и логи\n  uploads_data:\n    # Загруженные файлы\n\n# Bind mounts в /opt/hostel-data/\n/opt/hostel-data/opensearch:/usr/share/opensearch/data\n/opt/hostel-data/minio:/data\n/opt/hostel-data/certbot/conf:/etc/letsencrypt\n/opt/hostel-data/certbot/www:/var/www/certbot\n/opt/hostel-data/uploads:/usr/share/nginx/uploads\n/opt/hostel-data/maildata:/var/mail\n/opt/hostel-data/mailstate:/var/mail-state\n/opt/hostel-data/maillogs:/var/log/mail\n```\n\n## 🔄 Lifecycle Management\n\n### Restart Policies\n```yaml\n# Development\nrestart: unless-stopped  # Только MinIO\n# Остальные сервисы: no (default)\n\n# Production\nrestart: unless-stopped  # Все основные сервисы\nrestart: \"no\"            # migrate (разовая задача)\n```\n\n### Health Checks\n\n#### PostgreSQL\n```yaml\nhealthcheck:\n  test: [\"CMD\", \"pg_isready\", \"-U\", \"postgres\"]\n  interval: 10s\n  timeout: 5s\n  retries: 5\n```\n\n#### MinIO\n```yaml\nhealthcheck:\n  test: [\"CMD\", \"curl\", \"-f\", \"http://localhost:9000/minio/health/live\"]\n  interval: 30s\n  timeout: 10s\n  retries: 3\n```\n\n#### Backend (Production)\n```yaml\nhealthcheck:\n  test: [\"CMD\", \"curl\", \"-f\", \"http://localhost:3000\"]\n  interval: 10s\n  timeout: 5s\n  retries: 3\n```\n\n#### Nginx (Production)\n```yaml\nhealthcheck:\n  test: [\"CMD\", \"wget\", \"--spider\", \"--quiet\", \"http://localhost/\"]\n  interval: 30s\n  timeout: 10s\n  retries: 3\n  start_period: 15s\n```\n\n### Graceful Shutdown\n```yaml\nstop_grace_period: 10s\nstop_signal: SIGINT\n```\n\n## 🏷️ Container Registry\n\n### Development\n**Source:** Docker Hub (публичные образы)\n- postgres:15\n- opensearchproject/opensearch:2.11.0\n- opensearchproject/opensearch-dashboards:2.11.0\n- minio/minio:RELEASE.2023-09-30T07-02-29Z\n- minio/mc\n- migrate/migrate\n\n### Production\n**Source:** Harbor Registry (приватные образы)\n- harbor.svetu.rs/svetu/db/postgres:15\n- harbor.svetu.rs/svetu/opensearch/opensearch:2.11.0\n- harbor.svetu.rs/svetu/mail/server:latest\n- harbor.svetu.rs/svetu/tools/certbot:latest\n- harbor.svetu.rs/svetu/mail/webui:latest\n- harbor.svetu.rs/svetu/tools/migrate:latest\n- harbor.svetu.rs/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z\n- harbor.svetu.rs/svetu/minio/mc:latest\n- harbor.svetu.rs/svetu/backend/api:latest\n- harbor.svetu.rs/svetu/nginx/nginx:latest\n\n## ⚙️ Переменные окружения\n\n### Общие для всех сред\n```bash\n# OpenSearch\ndiscovery.type=single-node\nbootstrap.memory_lock=true\nOPENSEARCH_JAVA_OPTS=-Xms1024m -Xmx1024m\nDISABLE_SECURITY_PLUGIN=true\nOPENSEARCH_HOSTS=http://opensearch:9200\nDISABLE_SECURITY_DASHBOARDS_PLUGIN=true\n```\n\n### Development Specific\n```bash\n# PostgreSQL\nPOSTGRES_USER=postgres\nPOSTGRES_PASSWORD=password\nPOSTGRES_DB=hostel_db\n\n# MinIO\nMINIO_ROOT_USER=minioadmin\nMINIO_ROOT_PASSWORD=1321321321321\nMINIO_BROWSER_REDIRECT_URL=http://localhost:9001\nMINIO_SERVER_URL=http://localhost:9000\n```\n\n### Production Specific\n```bash\n# PostgreSQL\nPOSTGRES_USER=postgres\nPOSTGRES_PASSWORD=c9XWc7Cm\nPOSTGRES_DB=hostel_db\nPGDATA=/var/lib/postgresql/data/pgdata\n\n# MinIO\nMINIO_ROOT_USER=minioadmin\nMINIO_ROOT_PASSWORD=5465465465465\n\n# Backend API\nAPP_MODE=production\nENV_FILE=.env\nWS_ENABLED=true\nOPENSEARCH_URL=http://opensearch:9200\nOPENSEARCH_MARKETPLACE_INDEX=marketplace\nFILE_STORAGE_PROVIDER=minio\nMINIO_ENDPOINT=minio:9000\nMINIO_ACCESS_KEY=minioadmin\nMINIO_SECRET_KEY=5465465465465\nMINIO_USE_SSL=false\nMINIO_BUCKET_NAME=listings\nMINIO_LOCATION=eu-central-1\nFILE_STORAGE_PUBLIC_URL=https://svetu.rs\nDATABASE_URL=postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable\n```\n\n## 🚀 Deployment Strategies\n\n### Development Deployment\n```bash\n# Запуск основных сервисов\ndocker-compose up -d db opensearch opensearch-dashboards minio\n\n# Создание MinIO buckets\ndocker-compose up createbuckets\n\n# Миграции базы данных\ndocker-compose up migrate\n\n# Разработка backend/frontend локально\n# Не используются контейнеры для разработки\n```\n\n### Production Deployment\n```bash\n# Полный стек\ndocker-compose -f docker-compose.prod.yml up -d\n\n# Последовательный запуск с зависимостями\ndocker-compose -f docker-compose.prod.yml up -d db opensearch\ndocker-compose -f docker-compose.prod.yml up -d minio createbuckets\ndocker-compose -f docker-compose.prod.yml up -d migrate\ndocker-compose -f docker-compose.prod.yml up -d backend mailserver mail-webui\ndocker-compose -f docker-compose.prod.yml up -d nginx certbot\n```\n\n### Scaling Considerations\n```bash\n# Горизонтальное масштабирование backend\ndocker-compose -f docker-compose.prod.yml up -d --scale backend=3\n\n# Обновление образов без downtime\ndocker-compose -f docker-compose.prod.yml pull\ndocker-compose -f docker-compose.prod.yml up -d --no-deps backend\n```\n\n## 🔧 Service Dependencies\n\n### Порядок запуска\n```mermaid\ngraph TD\n    A[PostgreSQL] --> B[Migrations]\n    C[OpenSearch] --> D[Backend]\n    E[MinIO] --> F[CreateBuckets]\n    F --> D\n    B --> D\n    D --> G[Nginx]\n    H[MailServer] --> G\n    I[Mail-WebUI] --> G\n    J[Certbot] --> G\n```\n\n### Critical Dependencies\n1. **Database** должна быть готова перед **Migrations**\n2. **MinIO** должен быть готов перед **CreateBuckets**\n3. **Backend** зависит от **DB**, **OpenSearch**, **MinIO**\n4. **Nginx** зависит от **Backend** для проксирования\n5. **Mail services** независимы от основного приложения\n\n## 📊 Мониторинг и логирование\n\n### Health Check Endpoints\n```bash\n# PostgreSQL\ndocker-compose exec db pg_isready -U postgres\n\n# MinIO\ncurl http://localhost:9000/minio/health/live\n\n# Backend (production)\ncurl http://localhost:3000/api/health\n\n# Nginx\nwget --spider --quiet http://localhost/\n```\n\n### Container Logs\n```bash\n# Просмотр логов сервиса\ndocker-compose logs -f backend\n\n# Все логи\ndocker-compose logs\n\n# Логи с временными метками\ndocker-compose logs -t --since=\"1h\"\n```\n\n### Resource Monitoring\n```bash\n# Использование ресурсов\ndocker stats\n\n# Состояние контейнеров\ndocker-compose ps\n\n# Health status\ndocker inspect --format='{{.State.Health.Status}}' container_name\n```\n\n## 🛡️ Security Considerations\n\n### Network Isolation\n- Все сервисы изолированы в `hostel_network`\n- Внешний доступ только к необходимым портам\n- Inter-service communication через внутренние имена\n\n### Secrets Management\n- Development: простые пароли в compose файле\n- Production: сложные пароли в environment variables\n- SSL сертификаты через Let's Encrypt в volumes\n\n### User Privileges\n- Контейнеры не требуют root привилегий (кроме mailserver)\n- Mail server требует `NET_ADMIN` для сетевых операций\n\n---\n**Паспорт создан:** 2025-06-29  \n**Компонент:** Docker Compose Setup  \n**Статус:** Активный в development и production