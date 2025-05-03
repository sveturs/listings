#!/bin/bash
set -e # Останавливаем выполнение при ошибках
echo "Начинаем полный деплой с использованием Harbor..."
cd /opt/hostel-booking-system

# Останавливаем текущие службы
echo "Останавливаем текущие сервисы..."
docker-compose -f docker-compose.prod.yml down

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u admin -p SveTu2025 harbor.svetu.rs

# Загрузка всех образов из Harbor
echo "Загрузка образов из Harbor..."
docker pull harbor.svetu.rs/svetu/db/postgres:15
docker pull harbor.svetu.rs/svetu/opensearch/opensearch:2.11.0
docker pull harbor.svetu.rs/svetu/tools/migrate:latest
docker pull harbor.svetu.rs/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
docker pull harbor.svetu.rs/svetu/minio/mc:latest
docker pull harbor.svetu.rs/svetu/backend/api:latest
docker pull harbor.svetu.rs/svetu/tools/certbot:latest
docker pull harbor.svetu.rs/svetu/mail/server:latest
docker pull harbor.svetu.rs/svetu/mail/webui:latest
docker pull harbor.svetu.rs/svetu/frontend/app:latest
docker pull harbor.svetu.rs/svetu/nginx/nginx:latest

# Создаем необходимые директории
mkdir -p /tmp/hostel-backup/db
mkdir -p backend/uploads
mkdir -p /opt/hostel-data/{opensearch,uploads,minio}

# Делаем бэкап базы данных
echo "Создание бэкапа базы данных..."
if docker ps -a | grep -q hostel_db; then
  BACKUP_FILE="/tmp/hostel-backup/db/backup_$(date +%Y%m%d_%H%M%S).sql"
  docker exec -t hostel_db pg_dumpall -U postgres > "$BACKUP_FILE"
  if [ $? -eq 0 ]; then
    echo "Бэкап базы данных создан в $BACKUP_FILE"
  else
    echo "Ошибка создания бэкапа базы данных, но продолжаем деплой"
  fi
else
  echo "База данных не найдена, пропускаем создание бэкапа"
fi

# Обновляем файл docker-compose.prod.yml.harbor для использования всех образов из Harbor
echo "Обновление docker-compose.prod.yml.harbor..."
cat > docker-compose.prod.yml.harbor << EOL
version: '3.8'
services:
  # Авторизация в Harbor через HTTPS
  harbor-login:
    image: docker:cli
    command: sh -c "docker login -u admin -p SveTu2025 https://harbor.svetu.rs"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /etc/docker/certs.d/harbor.svetu.rs:/etc/docker/certs.d/harbor.svetu.rs
  db:
    image: harbor.svetu.rs/svetu/db/postgres:15
    container_name: hostel_db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: c9XWc7Cm
      POSTGRES_DB: hostel_db
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - db_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - hostel_network
    stop_grace_period: 10s
    stop_signal: SIGINT
    restart: unless-stopped

  opensearch:
    image: harbor.svetu.rs/svetu/opensearch/opensearch:2.11.0
    container_name: opensearch
    environment:
      - "discovery.type=single-node"
      - "bootstrap.memory_lock=true"
      - "OPENSEARCH_JAVA_OPTS=-Xms1024m -Xmx1024m"
      - "DISABLE_SECURITY_PLUGIN=true"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - opensearch_data:/usr/share/opensearch/data
    networks:
      - hostel_network
    restart: unless-stopped
  
  mailserver:
    image: harbor.svetu.rs/svetu/mail/server:latest
    container_name: mailserver
    hostname: mail.svetu.rs
    domainname: svetu.rs
    ports:
      - "25:25"      # SMTP
      - "587:587"    # SMTP с шифрованием
      - "465:465"    # SMTPS
      - "143:143"    # IMAP
      - "993:993"    # IMAPS
      - "110:110"    # POP3
      - "995:995"    # POP3S
    volumes:
      - ./mailserver/mail-data:/var/mail
      - ./mailserver/mail-state:/var/mail-state
      - ./mailserver/mail-logs:/var/log/mail
      - ./mailserver/config:/tmp/docker-mailserver
      - ./certbot/conf:/etc/letsencrypt
    environment:
      - ENABLE_SPAMASSASSIN=1
      - ENABLE_CLAMAV=0
      - ENABLE_FAIL2BAN=1
      - ENABLE_POSTGREY=1
      - SSL_TYPE=letsencrypt
      - PERMIT_DOCKER=network
      - ONE_DIR=1
      - DMS_DEBUG=0
      - ENABLE_AMAVIS=1
      - POSTMASTER_ADDRESS=postmaster@svetu.rs
      - SSL_CERT_PATH=/etc/letsencrypt/live/mail.svetu.rs/fullchain.pem
      - SSL_KEY_PATH=/etc/letsencrypt/live/mail.svetu.rs/privkey.pem
    cap_add:
      - NET_ADMIN
    networks:
      hostel_network:
        aliases:
          - mailserver
          - mail.svetu.rs
    restart: unless-stopped


  certbot:
    image: harbor.svetu.rs/svetu/tools/certbot:latest
    volumes:
      - ./certbot/conf:/etc/letsencrypt
      - ./certbot/www:/var/www/certbot
    depends_on:
      - nginx
    command: certonly --webroot --webroot-path=/var/www/certbot --email admin@svetu.rs --agree-tos --no-eff-email --force-renewal -d svetu.rs -d www.svetu.rs -d mail.svetu.rs -d autodiscover.svetu.rs -d autoconfig.svetu.rs -d klimagrad.svetu.rs


  mail-webui:
    image: harbor.svetu.rs/svetu/mail/webui:latest
    container_name: mail-webui
    volumes:
      - ./roundcube/config:/var/www/html/config
      - roundcube_data:/var/roundcube  # Добавьте этот том для постоянного хранения данных
    environment:
      - ROUNDCUBEMAIL_DEFAULT_HOST=ssl://mailserver:993
      - ROUNDCUBEMAIL_SMTP_SERVER=tls://mailserver:587
      - ROUNDCUBEMAIL_SMTP_PORT=587
      - ROUNDCUBEMAIL_PLUGINS=archive,zipdownload
      - ROUNDCUBEMAIL_DES_KEY=random-string-for-encryption
      - ROUNDCUBEMAIL_SKIN=elastic
    expose:
      - "80"
    networks:
      hostel_network:
        aliases:
          - mail-webui
    restart: unless-stopped
    depends_on:
      - mailserver


#  opensearch-dashboards:
#    image: harbor.svetu.rs/svetu/opensearch/dashboards:2.11.0
#    container_name: opensearch-dashboards
#    ports:
#      - "5601:5601"
#    environment:
#      - "OPENSEARCH_HOSTS=http://opensearch:9200"
#      - "DISABLE_SECURITY_DASHBOARDS_PLUGIN=true"
#    networks:
#      - hostel_network
#    depends_on:
#      - opensearch
#    restart: unless-stopped

  migrate:
    image: harbor.svetu.rs/svetu/tools/migrate:latest
    volumes:
      - ./backend/migrations:/migrations
    command: [
      "-path", "/migrations",
      "-database", "postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable",
      "up"
    ]
    depends_on:
      db:
        condition: service_healthy
    networks:
      - hostel_network
    stop_grace_period: 1s
    restart: "no"

  minio:
    image: harbor.svetu.rs/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
    container_name: minio
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: 5465465465465
    volumes:
      - minio_data:/data
    ports:
      - "9000:9000"
      - "9001:9001"
    command: server /data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - hostel_network
    restart: unless-stopped

  createbuckets:
    image: harbor.svetu.rs/svetu/minio/mc:latest
    depends_on:
      - minio
    entrypoint: >
      /bin/sh -c "
      sleep 5;
      /usr/bin/mc config host add myminio http://minio:9000 minioadmin 5465465465465;
      /usr/bin/mc mb --ignore-existing myminio/listings;
      /usr/bin/mc policy set public myminio/listings;
      exit 0;
      "
    networks:
      - hostel_network


  backend:
    image: harbor.svetu.rs/svetu/backend/api:latest
    container_name: backend
    environment:
      - APP_MODE=production
      - ENV_FILE=.env
      - WS_ENABLED=true
      - OPENSEARCH_URL=http://opensearch:9200
      - OPENSEARCH_MARKETPLACE_INDEX=marketplace
      - FILE_STORAGE_PROVIDER=minio
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY=minioadmin
      - MINIO_SECRET_KEY=5465465465465
      - MINIO_USE_SSL=false
      - MINIO_BUCKET_NAME=listings
      - MINIO_LOCATION=eu-central-1
      - FILE_STORAGE_PUBLIC_URL=https://svetu.rs
    user: "101:101"
    volumes:
      - uploads_data:/app/uploads
    expose:
      - "3000"
    ports:
      - "3000:3000"
    depends_on:
      db:
        condition: service_healthy
      opensearch:
        condition: service_started
      minio:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000"]
      interval: 10s
      timeout: 5s
      retries: 3
    networks:
      - hostel_network
    stop_grace_period: 10s
    stop_signal: SIGTERM
    init: true
    restart: unless-stopped

  nginx:
    image: harbor.svetu.rs/svetu/nginx/nginx:latest
    container_name: hostel_nginx
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
      - ./frontend/hostel-frontend/build:/usr/share/nginx/html:ro
      - /opt/hostel-data/klimagrad-site:/usr/share/nginx/klimagrad-site:ro
      - uploads_data:/usr/share/nginx/uploads:ro
      - ./certbot/conf:/etc/letsencrypt
      - ./certbot/www:/var/www/certbot
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - backend
    healthcheck:
      test: ["CMD", "wget", "--spider", "--quiet", "http://localhost/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 15s
    networks:
      - hostel_network
    stop_grace_period: 5s
    stop_signal: SIGTERM
    restart: unless-stopped

volumes:
  db_data:
  opensearch_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /opt/hostel-data/opensearch
  uploads_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /opt/hostel-data/uploads
  roundcube_data:  # Том для хранения данных Roundcube
  minio_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /opt/hostel-data/minio

networks:
  hostel_network:
    driver: bridge
EOL

# Запускаем базу данных из Harbor
echo "Запускаем базу данных из Harbor..."
docker-compose -f docker-compose.prod.yml.harbor up -d db

# Проверяем базу данных
echo "Проверяем готовность базы данных..."
RETRY_COUNT=30
for i in $(seq 1 $RETRY_COUNT); do
  if docker-compose -f docker-compose.prod.yml.harbor exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "База данных готова!"
    break
  fi
  echo "Ожидаем готовность базы данных... Попытка $i/$RETRY_COUNT"
  sleep 2
done

if ! docker-compose -f docker-compose.prod.yml.harbor exec -T db pg_isready -U postgres > /dev/null 2>&1; then
  echo "Не удалось запустить базу данных"
  exit 1
fi

# Запускаем миграции
echo "Запускаем миграции..."
docker-compose -f docker-compose.prod.yml.harbor up -d migrate

# Восстанавливаем данные из бэкапа, если он есть
if [ -n "$(ls -t /tmp/hostel-backup/db/*.sql 2>/dev/null | head -1)" ]; then
  LATEST_BACKUP=$(ls -t /tmp/hostel-backup/db/*.sql | head -1)
  echo "Восстанавливаем базу данных из $LATEST_BACKUP..."
  cat "$LATEST_BACKUP" | docker-compose -f docker-compose.prod.yml.harbor exec -T db psql -U postgres
  if [ $? -eq 0 ]; then
    echo "База данных успешно восстановлена"
  else
    echo "Ошибка восстановления базы данных, но продолжаем деплой"
  fi
else
  echo "Бэкап базы данных не найден, пропускаем восстановление"
fi

# Запускаем основные сервисы из Harbor
echo "Запускаем основные сервисы из Harbor..."
docker-compose -f docker-compose.prod.yml.harbor up -d opensearch minio createbuckets

# Запускаем backend и nginx из Harbor 
echo "Запускаем backend и nginx из Harbor..."
docker-compose -f docker-compose.prod.yml.harbor up -d backend nginx

# Запускаем почтовые сервисы из Harbor
echo "Запускаем почтовые сервисы из Harbor..."
docker-compose -f docker-compose.prod.yml.harbor up -d mailserver mail-webui certbot

echo "Полный деплой с использованием Harbor завершен!"
echo "Все сервисы теперь используют образы из Harbor:"
echo "- db, backend, opensearch, minio, createbuckets, migrate, nginx, certbot, mailserver, mail-webui"