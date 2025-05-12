#!/bin/bash
# Улучшенный скрипт для blue-green deployment с поддержкой всех сервисов
# Использование: ./blue_green_deploy_on_svetu.rs.sh [backend|frontend|all] [-m]

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Функция для логирования с разным уровнем важности
log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${RED}[WARNING]${NC} $1"
}

log_debug() {
    if [ "${DEBUG:-false}" = "true" ]; then
        echo -e "${CYAN}[DEBUG]${NC} $1"
    fi
}

# Функция фильтрации нежелательных логов
filter_logs() {
    grep -v -E "upstream server temporarily disabled|No route to host|connect\(\) failed|redirection cycle|docker-entrypoint" || true
}

# Проверяем аргументы
if [ -z "$1" ]; then
  log_error "Укажите название сервиса (backend, frontend, all)"
  exit 1
fi

SERVICE="$1"
HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROD_SERVER="161.97.89.28"
PROD_SERVER_USER="root"
PROD_SERVER_PATH="/opt/hostel-booking-system"
LOCAL_BACKEND_DIR="./backend"
LOCAL_FRONTEND_DIR="/data/hostel-booking-system/frontend/hostel-frontend"

# Создаем директорию для логов
LOG_DIR="./logs"
mkdir -p $LOG_DIR
LOG_FILE="${LOG_DIR}/blue_green_deploy_$(date +%Y%m%d%H%M%S).log"

# Проверяем флаг миграций
RUN_MIGRATIONS=false
if [[ "$2" == "-m" ]]; then
  RUN_MIGRATIONS=true
  log_info "Включен режим миграций. После деплоя будут применены миграции базы данных."
fi

# Текущая дата и время для тегов
TIMESTAMP=$(date +%Y%m%d%H%M%S)

# Функция для сборки и загрузки backend
build_backend() {
  log_info "Сборка и загрузка backend..."
  cd $LOCAL_BACKEND_DIR

  # Сборка образа с тегом релиза
  docker build -t backend:$TIMESTAMP .

  # Тегирование для Harbor
  docker tag backend:$TIMESTAMP $HARBOR_URL/svetu/backend/api:$TIMESTAMP
  docker tag backend:$TIMESTAMP $HARBOR_URL/svetu/backend/api:latest

  # Загрузка в Harbor
  docker push $HARBOR_URL/svetu/backend/api:$TIMESTAMP
  docker push $HARBOR_URL/svetu/backend/api:latest

  log_success "Backend успешно загружен в Harbor с тегами: latest и $TIMESTAMP"
}

# Функция для сборки и загрузки frontend
build_frontend() {
  log_info "Сборка и загрузка frontend..."
  cd $LOCAL_FRONTEND_DIR

  # Создаем env.js с динамической конфигурацией для production
  cat > public/env.js << EOT
window.ENV = {
  REACT_APP_BACKEND_URL: 'https://svetu.rs',
  REACT_APP_API_PREFIX: '/api',
  REACT_APP_AUTH_PREFIX: '/auth',
  REACT_APP_WEBSOCKET_URL: 'wss://svetu.rs/ws',
  REACT_APP_WS_URL: 'wss://svetu.rs/ws/chat',
  REACT_APP_MINIO_URL: 'https://svetu.rs',
  REACT_APP_HOST: 'https://svetu.rs'
};
EOT

  # Сборка проекта
  npm run build

  # Проверяем наличие env.js в сборке
  if [ ! -f "build/env.js" ]; then
    cp public/env.js build/env.js
  fi

  # Создаем копию как process-env.js для обратной совместимости
  cp build/env.js build/process-env.js

  # Сборка образа с тегом релиза
  docker build -t frontend:$TIMESTAMP -f Dockerfile .

  # Тегирование для Harbor
  docker tag frontend:$TIMESTAMP $HARBOR_URL/svetu/frontend/app:$TIMESTAMP
  docker tag frontend:$TIMESTAMP $HARBOR_URL/svetu/frontend/app:latest

  # Загрузка в Harbor
  docker push $HARBOR_URL/svetu/frontend/app:$TIMESTAMP
  docker push $HARBOR_URL/svetu/frontend/app:latest

  log_success "Frontend успешно загружен в Harbor с тегами: latest и $TIMESTAMP"
}

# Функция для бесшовного обновления backend
blue_green_backend() {
  log_info "Бесшовное обновление backend..."

  # Создаем единый скрипт для выполнения на сервере с улучшенной поддержкой всех сервисов
  cat > /tmp/server_deploy.sh << 'EOT'
#!/bin/bash
set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Функция для логирования с разным уровнем важности
log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${RED}[WARNING]${NC} $1"
}

log_debug() {
    if [ "${DEBUG:-false}" = "true" ]; then
        echo -e "${CYAN}[DEBUG]${NC} $1"
    fi
}

# Полностью отключаем вывод логов nginx кроме критических
filter_nginx_logs() {
    # Выводим только критические ошибки уровня [emerg] и [crit]
    grep -E "\[emerg\]|\[crit\]" 2>/dev/null || true
}

# Функция для запуска команд с подавлением лишнего вывода
run_quietly() {
    "$@" > /dev/null 2>&1 || true
}

# Проверяем параметры
RUN_MIGRATIONS=false
if [[ "$1" == "run-migrations" ]]; then
  RUN_MIGRATIONS=true
  log_info "Включен режим миграций. После деплоя будут применены миграции базы данных."
fi

# Создаем рабочую директорию
WORK_DIR="/tmp/bluegreen"
rm -rf $WORK_DIR
mkdir -p $WORK_DIR

# Определение портов для blue/green
BLUE_PORT=8081
GREEN_PORT=8082

# Ищем правильную сеть - сначала пробуем найти сеть hostel_network
HOSTEL_NETWORK="hostel-booking-system_hostel_network"
if ! docker network inspect $HOSTEL_NETWORK >/dev/null 2>&1; then
  log_info "Сеть $HOSTEL_NETWORK не найдена, ищем альтернативы..."
  # Ищем любую сеть с hostel в имени
  HOSTEL_NETWORK=$(docker network ls --format "{{.Name}}" | grep -i hostel | head -n 1)
fi

# Если сеть не найдена, используем default
if [ -z "$HOSTEL_NETWORK" ]; then
  log_info "Не найдена сеть с hostel в имени, используем default..."
  HOSTEL_NETWORK="hostel-booking-system_default"
  if ! docker network inspect $HOSTEL_NETWORK >/dev/null 2>&1; then
    # Ищем любую default сеть
    HOSTEL_NETWORK=$(docker network ls --format "{{.Name}}" | grep -i default | head -n 1)
  fi
fi

log_info "Будем использовать сеть: $HOSTEL_NETWORK"

# Автоматическое определение контейнера с базой данных
DB_CONTAINER=$(docker ps --format "{{.Names}}" | grep -i "db\|postgres" | head -n 1)
if [ -z "$DB_CONTAINER" ]; then
  log_info "Контейнер базы данных не найден, используем 'hostel_db'"
  DB_CONTAINER="hostel_db"
else
  log_info "Найден контейнер базы данных: $DB_CONTAINER"
fi

# Определение минио контейнера
MINIO_CONTAINER=$(docker ps --format "{{.Names}}" | grep -i "minio" | head -n 1)
if [ -z "$MINIO_CONTAINER" ]; then
  log_info "Контейнер MinIO не найден, используем 'minio'"
  MINIO_CONTAINER="minio"
else
  log_info "Найден контейнер MinIO: $MINIO_CONTAINER"
fi

# Определение opensearch контейнера
OPENSEARCH_CONTAINER=$(docker ps --format "{{.Names}}" | grep -i "opensearch" | head -n 1)
if [ -z "$OPENSEARCH_CONTAINER" ]; then
  log_info "Контейнер OpenSearch не найден, используем 'opensearch'"
  OPENSEARCH_CONTAINER="opensearch"
else
  log_info "Найден контейнер OpenSearch: $OPENSEARCH_CONTAINER"
fi

# Получаем IP-адреса сервисов (для резервного использования)
MINIO_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $MINIO_CONTAINER 2>/dev/null || echo "")
DB_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $DB_CONTAINER 2>/dev/null || echo "")
OPENSEARCH_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $OPENSEARCH_CONTAINER 2>/dev/null || echo "")

log_info "IP-адрес MinIO: $MINIO_IP"
log_info "IP-адрес DB: $DB_IP"
log_info "IP-адрес OpenSearch: $OPENSEARCH_IP"

# Получаем IP-подсеть этой сети
NETWORK_SUBNET=$(docker network inspect $HOSTEL_NETWORK 2>/dev/null | grep -oP '"Subnet": "\K[^"]+' || echo "")
log_info "Сеть '$HOSTEL_NETWORK' использует подсеть: $NETWORK_SUBNET"

log_info "Определяем текущую конфигурацию системы..."

# Проверяем текущие контейнеры
BLUE_RUNNING=$(docker ps --filter name=backend-blue --filter status=running -q | wc -l)
GREEN_RUNNING=$(docker ps --filter name=backend-green --filter status=running -q | wc -l)
ORIGINAL_RUNNING=$(docker ps --filter name=backend --filter status=running -q | wc -l)

# Вывод информации о существующих контейнерах
log_info "Текущие контейнеры:"
echo -e "BLUE_RUNNING: $BLUE_RUNNING"
echo -e "GREEN_RUNNING: $GREEN_RUNNING"
echo -e "ORIGINAL_RUNNING: $ORIGINAL_RUNNING"

# Определяем текущий и новый цвет
if [ "$BLUE_RUNNING" -gt "0" ]; then
  CURRENT_COLOR="blue"
  NEW_COLOR="green"
  CURRENT_CONTAINER="backend-blue"
  NEW_CONTAINER="backend-green"
  NEW_PORT=$GREEN_PORT
elif [ "$GREEN_RUNNING" -gt "0" ]; then
  CURRENT_COLOR="green"
  NEW_COLOR="blue"
  CURRENT_CONTAINER="backend-green"
  NEW_CONTAINER="backend-blue"
  NEW_PORT=$BLUE_PORT
elif [ "$ORIGINAL_RUNNING" -gt "0" ]; then
  CURRENT_COLOR="original"
  NEW_COLOR="blue"
  CURRENT_CONTAINER="backend"
  NEW_CONTAINER="backend-blue"
  NEW_PORT=$BLUE_PORT
else
  CURRENT_COLOR="none"
  NEW_COLOR="blue"
  CURRENT_CONTAINER="none"
  NEW_CONTAINER="backend-blue"
  NEW_PORT=$BLUE_PORT
fi

log_info "Текущий активный контейнер: $CURRENT_COLOR ($CURRENT_CONTAINER)"
log_info "Будет создан новый контейнер: $NEW_COLOR ($NEW_CONTAINER) на порту $NEW_PORT"

# Останавливаем существующий контейнер для нового цвета, если он есть
if docker ps -a --filter name=$NEW_CONTAINER -q | grep -q .; then
  log_info "Останавливаем существующий контейнер $NEW_CONTAINER..."
  docker stop $NEW_CONTAINER 2>/dev/null || true
  docker rm $NEW_CONTAINER 2>/dev/null || true
fi

# Авторизация в Harbor
log_info "Авторизация в Harbor..."
run_quietly docker login -u admin -p SveTu2025 harbor.svetu.rs

# Получаем новый образ
log_info "Загрузка нового образа backend:latest..."
run_quietly docker pull harbor.svetu.rs/svetu/backend/api:latest

# Получаем IP-адрес сети хоста
HOST_IP=$(hostname -I | awk '{print $1}')
log_info "IP-адрес хоста: $HOST_IP"

# Запускаем новый контейнер - ВАЖНО! Сразу указываем правильную сеть!
log_info "Запуск нового контейнера $NEW_CONTAINER на порту $NEW_PORT в сети $HOSTEL_NETWORK..."

# Используем комбинацию имен и IP-адресов (IP как резервный вариант)
docker run -d --name $NEW_CONTAINER \
  --network $HOSTEL_NETWORK \
  -p $NEW_PORT:3000 \
  -v /opt/hostel-data/uploads:/app/uploads \
  -v /opt/hostel-data/minio:/data/minio \
  -v /opt/hostel-data/credentials:/app/credentials \
  -e APP_MODE=production \
  -e ENV=production \
  -e WS_ENABLED=true \
  -e FILE_STORAGE_PROVIDER=minio \
  -e MINIO_USE_SSL=false \
  -e MINIO_BUCKET_NAME=listings \
  -e MINIO_LOCATION=eu-central-1 \
  -e FILE_STORAGE_PUBLIC_URL=https://svetu.rs \
  -e PORT=3000 \
  -e SERVER_HOST=https://svetu.rs \
  -e POSTGRES_HOST=$DB_CONTAINER \
  -e POSTGRES_PASSWORD=c9XWc7Cm \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=hostel_db \
  -e DATABASE_URL="postgres://postgres:c9XWc7Cm@$DB_CONTAINER:5432/hostel_db?sslmode=disable" \
  -e DB_IP="$DB_IP" \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=5465465465465 \
  -e MINIO_ACCESS_KEY=minioadmin \
  -e MINIO_SECRET_KEY=5465465465465 \
  -e MINIO_ENDPOINT="$MINIO_CONTAINER:9000" \
  -e MINIO_IP="$MINIO_IP" \
  -e OPENSEARCH_URL="http://$OPENSEARCH_CONTAINER:9200" \
  -e OPENSEARCH_IP="$OPENSEARCH_IP" \
  -e OPENSEARCH_USERNAME=admin \
  -e OPENSEARCH_PASSWORD=admin \
  -e OPENSEARCH_MARKETPLACE_INDEX=marketplace \
  -e GOOGLE_CLIENT_ID=917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com \
  -e GOOGLE_CLIENT_SECRET=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 \
  -e GOOGLE_OAUTH_REDIRECT_URL=https://svetu.rs/auth/google/callback \
  -e FRONTEND_URL=https://svetu.rs \
  -e JWT_SECRET=yoursecretkey \
  -e GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/neat-environs-140712-40c581381093.json \
  -e GOOGLE_TRANSLATE_API_KEY=AIzaSyCBFx7eUrqoUjFzsVIDrGXgqxmtymWhlpE \
  -e OPENAI_API_KEY=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 \
  -e TELEGRAM_BOT_TOKEN=8091644299:AAHHqzkVxgwLr_E5xwaTP_UQUzKWhVPXfi8 \
  -e STRIPE_API_KEY=sk_test_51Qw4kpEFv5ruKgVJwdkgenH8UU32QLw3o5QOBiicGm40jtI33Kv6fIouhx2r6g489e9iaqwqybApFvXvA6X4BYMR00egncZOqk \
  -e STRIPE_PUBLISHABLE_KEY=pk_test_51Qw4kpEFv5ruKgVJ1cLiCo0M4ztd1PmCFPi6ZG8JROm9TncM8tVffrrkxc5mM96egSxSsErFhjIqYMNrS4iePhY7005VezvaYe \
  -e STRIPE_WEBHOOK_SECRET=whsec_25fd811f9022a61563d4acf52e4dd2ced797a4e029828502b585bcb612b9bd81 \
  -e OPENSEARCH_USERNAME=admin \
  -e OPENSEARCH_PASSWORD=admin \
  -e OPENSEARCH_MARKETPLACE_INDEX=marketplace \
  -e EMAIL_PASSWORD=Pass4ma!l \
  harbor.svetu.rs/svetu/backend/api:latest > /dev/null

# Проверяем, что контейнер подключен к нужной сети
NETWORK_CHECK=$(docker inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}}{{end}}' $NEW_CONTAINER 2>/dev/null || echo "")
log_info "Контейнер $NEW_CONTAINER подключен к сети: $NETWORK_CHECK"

# Получаем IP-адрес нового контейнера
log_info "Ожидаем запуск контейнера (30 секунд)..."
sleep 30

NEW_CONTAINER_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $NEW_CONTAINER 2>/dev/null || echo "")
log_info "IP-адрес нового контейнера $NEW_CONTAINER: $NEW_CONTAINER_IP"

# Вывод логов контейнера для диагностики - сокращаем вывод
log_info "Ключевые логи нового контейнера (только запуск):"
docker logs --tail 5 $NEW_CONTAINER 2>&1 | grep -v "warning" | grep -v "deprecated" || true

# Проверка сетевого подключения между контейнерами
log_info "Проверка связи между контейнерами..."
docker exec $NEW_CONTAINER sh -c "ping -c 1 $DB_CONTAINER || ping -c 1 $DB_IP || echo 'Ping к DB не прошел'" > /dev/null

# Проверка работоспособности напрямую
MAX_RETRIES=15
RETRY_COUNT=0
HEALTH_CHECK_OK=false

log_info "Проверка работоспособности нового контейнера по API endpoint..."
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  # Пробуем более простой эндпоинт /health сначала
  HEALTH_RESPONSE=$(curl -s --max-time 5 http://$NEW_CONTAINER_IP:3000/api/health || echo "Failed")

  if [[ "$HEALTH_RESPONSE" == *"OK"* ]]; then
    log_success "Контейнер отвечает на базовые запросы!"
    HEALTH_CHECK_OK=true
    break
  fi

  # Пробуем получить категории если /health не отвечает
  CATEGORY_RESPONSE=$(curl -s --max-time 5 http://$NEW_CONTAINER_IP:3000/api/v1/marketplace/category-tree | head -c 50 || echo "Failed")

  if [[ "$CATEGORY_RESPONSE" == *"categories"* ]] || [[ "$CATEGORY_RESPONSE" == *"id"* ]]; then
    log_success "Новый контейнер успешно запущен и отвечает на запросы API!"
    HEALTH_CHECK_OK=true
    break
  fi

  log_info "Ожидание ответа от нового контейнера... ($RETRY_COUNT/$MAX_RETRIES)"

  # Дополнительная диагностика только на 3й и 9й попытке
  if (( RETRY_COUNT == 3 )) || (( RETRY_COUNT == 9 )); then
    log_info "Диагностические данные..."
    # Проверка соединения с базой данных - без вывода
    docker exec $NEW_CONTAINER sh -c "ping -c 1 $DB_CONTAINER || echo 'Ping к DB не прошел'" > /dev/null
  fi

  sleep 7
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ "$HEALTH_CHECK_OK" != "true" ]; then
  log_error "Ошибка: Новый контейнер не отвечает на запросы API!"
  log_info "Останавливаем и удаляем новый контейнер..."
  docker stop $NEW_CONTAINER 2>/dev/null || true
  docker rm $NEW_CONTAINER 2>/dev/null || true
  exit 1
fi

# Проверка доступности по порту
log_info "Проверка доступности по порту $NEW_PORT..."
RETRY_COUNT=0
PORT_CHECK_OK=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  PORT_RESPONSE=$(curl -s --max-time 5 http://localhost:$NEW_PORT/api/health || echo "Failed")

  if [[ "$PORT_RESPONSE" == *"OK"* ]]; then
    log_success "Новый контейнер успешно доступен по порту $NEW_PORT!"
    PORT_CHECK_OK=true
    break
  fi

  # Пробуем категории
  CATEGORY_PORT_RESPONSE=$(curl -s --max-time 5 http://localhost:$NEW_PORT/api/v1/marketplace/category-tree | head -c 50 || echo "Failed")

  if [[ "$CATEGORY_PORT_RESPONSE" == *"categories"* ]] || [[ "$CATEGORY_PORT_RESPONSE" == *"id"* ]]; then
    log_success "Новый контейнер успешно доступен по порту $NEW_PORT!"
    PORT_CHECK_OK=true
    break
  fi

  log_info "Ожидание доступности по порту... ($RETRY_COUNT/$MAX_RETRIES)"

  sleep 5
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ "$PORT_CHECK_OK" != "true" ]; then
  log_error "Ошибка: Новый контейнер не доступен по порту $NEW_PORT!"
  log_info "Останавливаем и удаляем новый контейнер..."
  docker stop $NEW_CONTAINER 2>/dev/null || true
  docker rm $NEW_CONTAINER 2>/dev/null || true
  exit 1
fi

# Обновляем конфигурацию upstream в nginx
log_info "Обновляем конфигурацию upstream в nginx..."

# Сохраняем текущую конфигурацию nginx
cp /opt/hostel-booking-system/nginx.conf $WORK_DIR/nginx.conf.bak

# Проверяем текущее состояние конфигурации
CONFIG_CHECK=$(grep -A 4 "upstream api_backend\|upstream websocket_backend" /opt/hostel-booking-system/nginx.conf)
log_debug "Текущая конфигурация upstream блоков: $CONFIG_CHECK"

# Комплексное обновление upstream блоков
# Обновляем оба типа конфигурации: по имени сервера или по IP-адресу
# 1. Заменяем старое имя контейнера на новое
sed -i "s/server backend-[^:]*:3000;/server $NEW_CONTAINER:3000;/g" /opt/hostel-booking-system/nginx.conf
sed -i "s/server backend:3000;/server $NEW_CONTAINER:3000;/g" /opt/hostel-booking-system/nginx.conf

# 2. Если в конфигурации используются IP-адреса, обновляем их на новый IP
if grep -q "server 192.168." /opt/hostel-booking-system/nginx.conf; then
  log_info "Найдены IP-адреса в конфигурации upstream. Обновляем на новый IP: $NEW_CONTAINER_IP"
  # Проверка, что IP-адрес был получен
  if [ -z "$NEW_CONTAINER_IP" ]; then
    log_warn "Не удалось получить IP нового контейнера! Используем имя контейнера вместо IP."
    sed -i "s/server 192\.168\.[0-9]\+\.[0-9]\+:3000;/server $NEW_CONTAINER:3000;/g" /opt/hostel-booking-system/nginx.conf
  else
    # Заменяем любой IP из подсети 192.168.x.x на новый IP
    sed -i "s/server 192\.168\.[0-9]\+\.[0-9]\+:3000;/server $NEW_CONTAINER_IP:3000;/g" /opt/hostel-booking-system/nginx.conf
  fi
fi

# 3. Для надежности также добавляем проверку, что upstream блоки действительно обновились
NEW_CONFIG_CHECK=$(grep -A 4 "upstream api_backend\|upstream websocket_backend" /opt/hostel-booking-system/nginx.conf)
log_debug "Новая конфигурация upstream блоков: $NEW_CONFIG_CHECK"

# Проверка успешности обновления
if [[ "$NEW_CONFIG_CHECK" == *"$NEW_CONTAINER"* ]] || [[ "$NEW_CONFIG_CHECK" == *"$NEW_CONTAINER_IP"* ]]; then
  log_success "Конфигурация upstream успешно обновлена"
else
  log_warn "Возможно, конфигурация upstream не была обновлена. Проверяем вручную."

  # Определение типа upstream для ручного обновления
  if grep -q "upstream api_backend" /opt/hostel-booking-system/nginx.conf; then
    # Ручное создание upstream блоков с новым контейнером
    log_info "Ручное обновление upstream блоков..."

    # Создаем временные файлы с новыми upstream блоками
    cat > $WORK_DIR/new_api_upstream.txt << EOF
upstream api_backend {
    server $NEW_CONTAINER:3000;
    keepalive 32;
}
EOF

    cat > $WORK_DIR/new_ws_upstream.txt << EOF
upstream websocket_backend {
    server $NEW_CONTAINER:3000;
    keepalive 32;
}
EOF

    # Заменяем существующие upstream блоки
    sed -i '/upstream api_backend {/,/}/d' /opt/hostel-booking-system/nginx.conf
    sed -i '/upstream websocket_backend {/,/}/d' /opt/hostel-booking-system/nginx.conf

    # Добавляем новые upstream блоки в начало файла
    sed -i "1r $WORK_DIR/new_ws_upstream.txt" /opt/hostel-booking-system/nginx.conf
    sed -i "1r $WORK_DIR/new_api_upstream.txt" /opt/hostel-booking-system/nginx.conf

    log_info "Ручное обновление upstream блоков завершено"
  fi
fi

# Проверяем статус nginx
log_info "Проверяем статус nginx..."
NGINX_RUNNING=$(docker ps --filter name=hostel_nginx --filter status=running -q | wc -l)

if [ "$NGINX_RUNNING" -eq "0" ]; then
  log_info "Nginx не запущен, запускаем новый контейнер..."
  # Останавливаем существующий контейнер если есть
  docker stop hostel_nginx > /dev/null 2>&1 || true
  docker rm hostel_nginx > /dev/null 2>&1 || true

  # Запускаем новый контейнер Nginx
  docker run -d --name hostel_nginx \
    --network $HOSTEL_NETWORK \
    -p 80:80 -p 443:443 \
    -v /opt/hostel-booking-system/nginx.conf:/etc/nginx/conf.d/default.conf \
    -v /opt/hostel-booking-system/frontend/hostel-frontend/build:/usr/share/nginx/html \
    -v /opt/hostel-booking-system/certbot/conf:/etc/letsencrypt \
    -v /opt/hostel-booking-system/certbot/www:/var/www/certbot \
    -v /opt/hostel-data/uploads:/usr/share/nginx/uploads \
    harbor.svetu.rs/svetu/nginx/nginx:latest > /dev/null 2>&1 || log_error "Не удалось запустить контейнер Nginx"

  sleep 5
  NGINX_RUNNING=$(docker ps --filter name=hostel_nginx --filter status=running -q | wc -l)
  if [ "$NGINX_RUNNING" -gt "0" ]; then
    log_success "Nginx успешно запущен"
  else
    log_error "Не удалось запустить Nginx!"
  fi
else
  log_info "Nginx уже запущен, перезагружаем конфигурацию..."

  # Подключаем Nginx к нужной сети, если необходимо
  if ! docker network inspect $HOSTEL_NETWORK | grep -q "hostel_nginx"; then
    log_info "Подключаем Nginx к сети $HOSTEL_NETWORK..."
    docker network connect $HOSTEL_NETWORK hostel_nginx > /dev/null 2>&1 || true
  fi

  # Проверяем синтаксис конфигурации
  if ! docker exec hostel_nginx nginx -t > /dev/null 2>&1; then
    log_error "Ошибка в синтаксисе конфигурации Nginx! Восстанавливаем оригинальную конфигурацию..."
    cp $WORK_DIR/nginx.conf.bak /opt/hostel-booking-system/nginx.conf
    docker exec hostel_nginx nginx -s reload > /dev/null 2>&1 || docker restart hostel_nginx > /dev/null 2>&1
    exit 1
  fi

  # Перезагружаем конфигурацию
  log_info "Перезагружаем конфигурацию Nginx..."
  if ! docker exec hostel_nginx nginx -s reload > /dev/null 2>&1; then
    log_info "Не удалось перезагрузить конфигурацию, перезапускаем контейнер..."
    docker restart hostel_nginx > /dev/null 2>&1
    sleep 5
  fi
fi

# Выполняем тест соединения от Nginx к новому контейнеру
log_info "Проверяем связь между Nginx и новым контейнером..."
if ! docker exec hostel_nginx ping -c 1 $NEW_CONTAINER > /dev/null 2>&1; then
  log_warn "Nginx не может соединиться с контейнером по имени! Проверяем по IP..."
  if ! docker exec hostel_nginx ping -c 1 $NEW_CONTAINER_IP > /dev/null 2>&1; then
    log_error "Ошибка: Nginx не может соединиться с новым контейнером ни по имени, ни по IP!"
    log_info "Пробуем перезапустить Nginx..."
    docker restart hostel_nginx
    sleep 5
  else
    log_info "Nginx может соединиться с новым контейнером по IP, но не по имени. Это может вызвать проблемы."
  fi
else
  log_success "Nginx успешно соединяется с новым контейнером по имени!"
fi

# Проверяем доступность API
log_info "Проверка доступности API через Nginx..."
sleep 5 # Даем время на применение изменений

MAX_RETRIES=10
RETRY_COUNT=0
API_CHECK_OK=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  # Проверяем health endpoint - без полного вывода ответа
  API_HEALTH=$(curl -k -s https://svetu.rs/api/health || echo "Failed")
  if [[ "$API_HEALTH" == *"OK"* ]]; then
    log_success "API health endpoint доступен через Nginx!"
    API_CHECK_OK=true
    break
  fi

  log_info "Ожидание доступности API... ($RETRY_COUNT/$MAX_RETRIES)"
  sleep 3
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ "$API_CHECK_OK" != "true" ]; then
  log_error "API недоступен через Nginx после обновления!"
  log_info "Проверяем лог Nginx для определения проблемы..."
  NGINX_ERRORS=$(docker logs hostel_nginx 2>&1 | grep -E "error|connect\(\) failed" | tail -n 5)
  echo "$NGINX_ERRORS"

  log_info "Пробуем прямой запрос к контейнеру для диагностики..."
  curl -s http://$NEW_CONTAINER_IP:3000/api/health || echo "Прямой запрос к контейнеру также неудачен"

  log_warn "Восстанавливаем оригинальную конфигурацию Nginx..."
  cp $WORK_DIR/nginx.conf.bak /opt/hostel-booking-system/nginx.conf
  docker exec hostel_nginx nginx -s reload > /dev/null 2>&1 || docker restart hostel_nginx > /dev/null 2>&1

  log_error "Возникла проблема с доступом к API. Возможно, требуется прописать корректные upstream блоки."
  echo "Попробуйте выполнить следующие команды для исправления:"
  echo "1. Обновите upstream блоки вручную:"
  echo "   docker exec -it hostel_nginx bash -c \"sed -i 's/server [^;]*:3000;/server $NEW_CONTAINER:3000;/g' /etc/nginx/conf.d/default.conf\""
  echo "2. Перезагрузите конфигурацию:"
  echo "   docker exec hostel_nginx nginx -s reload"

  exit 1
else
  log_success "API успешно доступен через Nginx!"
fi

# Если все проверки прошли успешно, запускаем миграции если нужно
if [ "$RUN_MIGRATIONS" = "true" ]; then
  log_info "Запуск миграций..."

  # Клонируем репозиторий во временную директорию для получения актуальных миграций
  TEMP_DIR="/tmp/svetu-migrations"
  rm -rf $TEMP_DIR
  mkdir -p $TEMP_DIR

  # Проверяем наличие токена GitHub
  GITHUB_TOKEN_FILE="/root/.github_token"

  if [ -f "$GITHUB_TOKEN_FILE" ]; then
    # Используем токен из защищенного файла
    GITHUB_TOKEN=$(cat $GITHUB_TOKEN_FILE)
    log_info "Получение свежих миграций из GitHub с использованием токена..."
    git clone --depth 1 https://${GITHUB_TOKEN}@github.com/DmitruNS/hostel-booking-system.git $TEMP_DIR > /dev/null 2>&1
    CLONE_SUCCESS=$?

    if [ $CLONE_SUCCESS -ne 0 ]; then
      log_error "Ошибка при клонировании репозитория!"
      RUN_MIGRATIONS=false
    fi
  else
    log_error "Не найден токен GitHub для аутентификации"
    log_info "Для настройки доступа к GitHub, создайте файл $GITHUB_TOKEN_FILE с токеном."
    RUN_MIGRATIONS=false
  fi

  if [ "$RUN_MIGRATIONS" = "true" ]; then
    # Создаем директорию для миграций
    mkdir -p /opt/hostel-booking-system/backend/migrations

    # Проверяем наличие миграций
    if [ ! -d "$TEMP_DIR/backend/migrations" ] || [ "$(ls -A $TEMP_DIR/backend/migrations/*.up.sql 2>/dev/null | wc -l)" -eq "0" ]; then
      log_info "В репозитории нет файлов миграций. Пропускаем этап миграций."
      RUN_MIGRATIONS=false
    fi
  fi

  if [ "$RUN_MIGRATIONS" = "true" ]; then
    log_info "Копирование файлов миграций из репозитория..."
    cp -f $TEMP_DIR/backend/migrations/*.up.sql /opt/hostel-booking-system/backend/migrations/ 2>/dev/null || true

    log_info "Запуск контейнера migrate/migrate для выполнения миграций..."
    docker run --rm --network $HOSTEL_NETWORK \
      -v /opt/hostel-booking-system/backend/migrations:/migrations \
      migrate/migrate \
      -path=/migrations/ \
      -database="postgres://postgres:c9XWc7Cm@$DB_CONTAINER:5432/hostel_db?sslmode=disable" \
      up > /dev/null 2>&1

    MIGRATION_RESULT=$?
    if [ $MIGRATION_RESULT -ne 0 ]; then
      log_error "Ошибка при выполнении миграций, но продолжаем деплой."
    else
      log_success "Миграции успешно выполнены!"
    fi
  else
    log_info "Миграции пропущены."
  fi

  # Очистка временных файлов
  rm -rf $TEMP_DIR
fi

# Выполняем еще одну проверку, чтобы убедиться, что все работает
log_info "Финальная проверка API..."
FINAL_CHECK=$(curl -k -s -m 5 https://svetu.rs/api/health || echo "Failed")
if [[ "$FINAL_CHECK" == *"OK"* ]]; then
  log_success "Финальная проверка API успешна! Система готова к работе."
else
  log_warn "Финальная проверка API не прошла успешно. Возможно, требуется ручная проверка."
fi

# Если проверки прошли успешно, останавливаем старый контейнер
if [ "$CURRENT_COLOR" = "original" ]; then
  log_info "Останавливаем оригинальный контейнер backend..."
  docker stop backend > /dev/null 2>&1 || true
  docker rm backend > /dev/null 2>&1 || true
elif [ "$CURRENT_COLOR" != "none" ]; then
  log_info "Останавливаем предыдущий контейнер $CURRENT_CONTAINER..."
  docker stop $CURRENT_CONTAINER > /dev/null 2>&1 || true
  docker rm $CURRENT_CONTAINER > /dev/null 2>&1 || true
fi

log_success "Бесшовное обновление backend успешно завершено!"
log_success "Новый активный контейнер: $NEW_CONTAINER с IP: $NEW_CONTAINER_IP"

# Очистка временных файлов
rm -rf $WORK_DIR
EOT

  # Делаем скрипт исполняемым
  chmod +x /tmp/server_deploy.sh

  # Добавляем параметр миграций
  MIGRATIONS_PARAM=""
  if [ "$RUN_MIGRATIONS" = "true" ]; then
    MIGRATIONS_PARAM="run-migrations"
  fi

  # Отправляем скрипт на сервер
  scp -q /tmp/server_deploy.sh $PROD_SERVER_USER@$PROD_SERVER:$PROD_SERVER_PATH/server_deploy.sh

  # Запускаем скрипт на сервере с фильтрацией логов
  log_info "Запуск скрипта деплоя на сервере..."
  ssh $PROD_SERVER_USER@$PROD_SERVER "cd $PROD_SERVER_PATH && chmod +x server_deploy.sh && ./server_deploy.sh $MIGRATIONS_PARAM" 2>&1 |
    grep -v -E "upstream server temporarily disabled|No route to host|connect\(\) failed|redirection cycle|docker-entrypoint|Using /etc/nginx/conf.d/default.conf" ||
    true

  # Удаляем временный скрипт
  rm /tmp/server_deploy.sh
  log_success "Деплой бэкенда завершен!"
}

# Функция для бесшовного обновления frontend
blue_green_frontend() {
  log_info "Бесшовное обновление frontend..."

  # Создаем скрипт для обновления frontend с поддержкой env.js
  cat > /tmp/frontend_deploy.sh << 'EOT'
#!/bin/bash
set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Функция для логирования с разным уровнем важности
log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${RED}[WARNING]${NC} $1"
}

log_debug() {
    if [ "${DEBUG:-false}" = "true" ]; then
        echo -e "${CYAN}[DEBUG]${NC} $1"
    fi
}

# Фильтрация логов Nginx
filter_nginx_logs() {
    grep -E "\[emerg\]|\[crit\]" 2>/dev/null || true
}

# Функция для запуска команд с подавлением лишнего вывода
run_quietly() {
    "$@" > /dev/null 2>&1 || true
}

# Создаем рабочие директории
WORK_DIR="/tmp/bluegreen"
rm -rf $WORK_DIR
mkdir -p $WORK_DIR
BACKUP_DIR="$WORK_DIR/frontend_backup"
mkdir -p $BACKUP_DIR

# Сохраняем текущую версию frontend
log_info "Создаем резервную копию frontend..."
mkdir -p /opt/hostel-booking-system/frontend/hostel-frontend/build.backup/
cp -r /opt/hostel-booking-system/frontend/hostel-frontend/build/* /opt/hostel-booking-system/frontend/hostel-frontend/build.backup/ 2>/dev/null || true

# Авторизация в Harbor
log_info "Авторизация в Harbor..."
run_quietly docker login -u admin -p SveTu2025 harbor.svetu.rs

# Загружаем последнюю версию frontend
log_info "Загрузка нового образа frontend..."
run_quietly docker pull harbor.svetu.rs/svetu/frontend/app:latest

# Создаем временный контейнер для извлечения сборки
log_info "Создаем временный контейнер..."
TEMP_CONTAINER=$(docker create harbor.svetu.rs/svetu/frontend/app:latest)

# Извлекаем сборку из контейнера
log_info "Извлекаем сборку из контейнера..."
docker cp $TEMP_CONTAINER:/app/build/. $WORK_DIR/frontend/

# Проверяем наличие index.html
if [ ! -f "$WORK_DIR/frontend/index.html" ]; then
  log_error "index.html не найден! Отмена обновления..."
  docker rm $TEMP_CONTAINER > /dev/null 2>&1 || true
  exit 1
fi

# Создаем env.js с правильными настройками для динамической конфигурации
log_info "Создаем env.js с нужными настройками..."
cat > $WORK_DIR/frontend/env.js << 'EOF'
window.ENV = {
  REACT_APP_BACKEND_URL: 'https://svetu.rs',
  REACT_APP_API_PREFIX: '/api',
  REACT_APP_AUTH_PREFIX: '/auth',
  REACT_APP_WEBSOCKET_URL: 'wss://svetu.rs/ws',
  REACT_APP_WS_URL: 'wss://svetu.rs/ws/chat',
  REACT_APP_MINIO_URL: 'https://svetu.rs',
  REACT_APP_HOST: 'https://svetu.rs'
};
EOF

# Добавляем env.production.js для совместимости
cat > $WORK_DIR/frontend/env.production.js << 'EOF'
window.ENV = {
  REACT_APP_BACKEND_URL: 'https://svetu.rs',
  REACT_APP_API_PREFIX: '/api',
  REACT_APP_AUTH_PREFIX: '/auth',
  REACT_APP_WEBSOCKET_URL: 'wss://svetu.rs/ws',
  REACT_APP_WS_URL: 'wss://svetu.rs/ws/chat',
  REACT_APP_MINIO_URL: 'https://svetu.rs',
  REACT_APP_HOST: 'https://svetu.rs'
};
EOF

# Проверяем и создаем директории
log_info "Проверяем и создаем директории..."
mkdir -p /opt/hostel-booking-system/frontend/hostel-frontend/build/

# Обновляем frontend
log_info "Обновляем frontend..."
rm -rf /opt/hostel-booking-system/frontend/hostel-frontend/build/*
cp -r $WORK_DIR/frontend/* /opt/hostel-booking-system/frontend/hostel-frontend/build/

# Удостоверяемся, что env.js правильно скопирован
cp $WORK_DIR/frontend/env.js /opt/hostel-booking-system/frontend/hostel-frontend/build/env.js
cp $WORK_DIR/frontend/env.production.js /opt/hostel-booking-system/frontend/hostel-frontend/build/env.production.js

# Проверяем результат обновления
if [ ! -f "/opt/hostel-booking-system/frontend/hostel-frontend/build/index.html" ]; then
  log_error "Не удалось обновить frontend! Восстанавливаем из резервной копии..."
  rm -rf /opt/hostel-booking-system/frontend/hostel-frontend/build/*
  cp -r /opt/hostel-booking-system/frontend/hostel-frontend/build.backup/* /opt/hostel-booking-system/frontend/hostel-frontend/build/ 2>/dev/null || true
  docker rm $TEMP_CONTAINER > /dev/null 2>&1 || true
  exit 1
fi

# Создаем копию env.js как process-env.js для обратной совместимости
cp /opt/hostel-booking-system/frontend/hostel-frontend/build/env.js /opt/hostel-booking-system/frontend/hostel-frontend/build/process-env.js

# Проверяем и устанавливаем правильные права
log_info "Устанавливаем правильные права на файлы..."
chmod -R 755 /opt/hostel-booking-system/frontend/hostel-frontend/build/

# Очистка
log_info "Очистка временных файлов..."
docker rm $TEMP_CONTAINER > /dev/null 2>&1 || true
rm -rf $WORK_DIR

# Проверяем статус Nginx
log_info "Проверяем статус Nginx..."
if docker ps -q --filter "name=hostel_nginx" | grep -q .; then
  log_info "Nginx запущен, перезагружаем конфигурацию..."
  # Перезагружаем Nginx, чтобы он подхватил новые файлы
  docker exec hostel_nginx nginx -s reload > /dev/null 2>&1 || true

  sleep 3
  # Проверяем доступность frontend
  log_info "Проверяем доступность frontend..."
  FRONTEND_RESPONSE=$(curl -k -s -L https://svetu.rs/ | grep -c "id=\"root\"" || echo "0")
  if [[ "$FRONTEND_RESPONSE" -gt "0" ]]; then
    log_success "Frontend успешно доступен!"
  else
    log_warn "Предупреждение: frontend может быть недоступен, проверяем вручную..."
    # Дополнительная проверка
    NGINX_STATUS=$(docker exec hostel_nginx nginx -t 2>&1)
    if [[ "$NGINX_STATUS" == *"successful"* ]]; then
      log_info "Конфигурация Nginx корректна, проблема может быть в другом месте."
    else
      log_error "Обнаружена проблема с конфигурацией Nginx: $NGINX_STATUS"
    fi
  fi
else
  log_info "Nginx не запущен, нужно его запустить вручную."
fi

log_success "Обновление frontend успешно завершено!"
EOT

  # Делаем скрипт исполняемым
  chmod +x /tmp/frontend_deploy.sh

  # Отправляем скрипт на сервер
  scp -q /tmp/frontend_deploy.sh $PROD_SERVER_USER@$PROD_SERVER:$PROD_SERVER_PATH/frontend_deploy.sh

  # Запускаем скрипт на сервере с фильтрацией логов
  log_info "Запуск обновления фронтенда на сервере..."
  ssh $PROD_SERVER_USER@$PROD_SERVER "cd $PROD_SERVER_PATH && chmod +x frontend_deploy.sh && ./frontend_deploy.sh" 2>&1 |
    grep -v -E "upstream server temporarily disabled|No route to host|connect\(\) failed|redirection cycle|docker-entrypoint|Using /etc/nginx/conf.d/default.conf" ||
    true

  # Удаляем временный скрипт
  rm /tmp/frontend_deploy.sh
  log_success "Деплой фронтенда завершен!"
}

# Создание лог-файла
echo "Начало логирования: $(date)" | tee $LOG_FILE

# Проверка пароля Harbor
if [ -z "$HARBOR_PASSWORD" ]; then
  log_info "Используем пароль Harbor по умолчанию (рекомендуется использовать переменную окружения)" | tee -a $LOG_FILE
  HARBOR_PASSWORD="SveTu2025"
fi

# Авторизация в Harbor
log_info "Авторизация в Harbor..." | tee -a $LOG_FILE
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL | grep -v "WARNING" || true

# Проверка статуса авторизации
if [ $? -ne 0 ]; then
  log_error "Ошибка авторизации в Harbor. Проверьте учетные данные." | tee -a $LOG_FILE
  exit 1
fi

# Обработка в зависимости от сервиса
case "$SERVICE" in
  "backend")
    build_backend | tee -a $LOG_FILE | filter_logs
    blue_green_backend | tee -a $LOG_FILE | filter_logs
    ;;
  "build-backend")
    build_backend | tee -a $LOG_FILE | filter_logs
    ;;
  "deploy-backend")
    blue_green_backend | tee -a $LOG_FILE | filter_logs
    ;;
  "frontend")
    build_frontend | tee -a $LOG_FILE | filter_logs
    blue_green_frontend | tee -a $LOG_FILE | filter_logs
    ;;
  "all")
    build_backend | tee -a $LOG_FILE | filter_logs
    build_frontend | tee -a $LOG_FILE | filter_logs
    blue_green_backend | tee -a $LOG_FILE | filter_logs
    blue_green_frontend | tee -a $LOG_FILE | filter_logs
    ;;
  *)
    log_error "Неизвестный сервис '$SERVICE'. Используйте backend, frontend или all." | tee -a $LOG_FILE
    exit 1
    ;;
esac

log_success "Процесс бесшовного обновления успешно завершен!" | tee -a $LOG_FILE
echo "Завершение логирования: $(date)" | tee -a $LOG_FILE