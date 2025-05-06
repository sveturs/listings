#!/bin/bash
# Улучшенный скрипт для blue-green deployment с поддержкой всех сервисов
# Использование: ./blue_green_deploy_on_svetu.rs.sh [backend|frontend|all]

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Проверяем аргументы
if [ -z "$1" ]; then
  echo -e "${RED}Ошибка: Укажите название сервиса (backend, frontend, all)${NC}"
  exit 1
fi

SERVICE="$1"
HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROD_SERVER="161.97.89.28"
PROD_SERVER_USER="root"
PROD_SERVER_PATH="/opt/hostel-booking-system"
LOCAL_BACKEND_DIR="/data/hostel-booking-system/backend"
LOCAL_FRONTEND_DIR="/data/hostel-booking-system/frontend/hostel-frontend"

# Текущая дата и время для тегов
TIMESTAMP=$(date +%Y%m%d%H%M%S)

# Функция для сборки и загрузки backend
build_backend() {
  echo -e "${YELLOW}Сборка и загрузка backend...${NC}"
  cd $LOCAL_BACKEND_DIR

  # Сборка образа с тегом релиза
  docker build -t backend:$TIMESTAMP .

  # Тегирование для Harbor
  docker tag backend:$TIMESTAMP $HARBOR_URL/svetu/backend/api:$TIMESTAMP
  docker tag backend:$TIMESTAMP $HARBOR_URL/svetu/backend/api:latest

  # Загрузка в Harbor
  docker push $HARBOR_URL/svetu/backend/api:$TIMESTAMP
  docker push $HARBOR_URL/svetu/backend/api:latest

  echo -e "${GREEN}Backend успешно загружен в Harbor с тегами: latest и $TIMESTAMP${NC}"
}

# Функция для сборки и загрузки frontend
build_frontend() {
  echo -e "${YELLOW}Сборка и загрузка frontend...${NC}"
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

  echo -e "${GREEN}Frontend успешно загружен в Harbor с тегами: latest и $TIMESTAMP${NC}"
}

# Функция для бесшовного обновления backend
blue_green_backend() {
  echo -e "${YELLOW}Бесшовное обновление backend...${NC}"

  # Создаем единый скрипт для выполнения на сервере с улучшенной поддержкой всех сервисов
  cat > /tmp/server_deploy.sh << 'EOT'
#!/bin/bash
set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Создаем рабочую директорию
WORK_DIR="/tmp/bluegreen"
rm -rf $WORK_DIR
mkdir -p $WORK_DIR

# Определение портов для blue/green
BLUE_PORT=8081
GREEN_PORT=8082

# Сохраняем текущую конфигурацию nginx для возможного восстановления
echo -e "${YELLOW}Сохраняем текущую конфигурацию nginx...${NC}"
cp /opt/hostel-booking-system/nginx.conf $WORK_DIR/nginx.conf.original

# Ищем правильную сеть - сначала пробуем найти сеть hostel_network
HOSTEL_NETWORK="hostel-booking-system_hostel_network"
if ! docker network inspect $HOSTEL_NETWORK >/dev/null 2>&1; then
  echo -e "${YELLOW}Сеть $HOSTEL_NETWORK не найдена, ищем альтернативы...${NC}"
  # Ищем любую сеть с hostel в имени
  HOSTEL_NETWORK=$(docker network ls --format "{{.Name}}" | grep -i hostel | head -n 1)
fi

# Если сеть не найдена, используем default
if [ -z "$HOSTEL_NETWORK" ]; then
  echo -e "${YELLOW}Не найдена сеть с hostel в имени, используем default...${NC}"
  HOSTEL_NETWORK="hostel-booking-system_default"
  if ! docker network inspect $HOSTEL_NETWORK >/dev/null 2>&1; then
    # Ищем любую default сеть
    HOSTEL_NETWORK=$(docker network ls --format "{{.Name}}" | grep -i default | head -n 1)
  fi
fi

echo -e "${YELLOW}Будем использовать сеть: $HOSTEL_NETWORK${NC}"

# Автоматическое определение контейнера с базой данных
DB_CONTAINER=$(docker ps --format "{{.Names}}" | grep -i "db\|postgres" | head -n 1)
if [ -z "$DB_CONTAINER" ]; then
  echo -e "${YELLOW}Контейнер базы данных не найден, используем 'hostel_db'${NC}"
  DB_CONTAINER="hostel_db"
else
  echo -e "${YELLOW}Найден контейнер базы данных: $DB_CONTAINER${NC}"
fi

# Определение минио контейнера
MINIO_CONTAINER=$(docker ps --format "{{.Names}}" | grep -i "minio" | head -n 1)
if [ -z "$MINIO_CONTAINER" ]; then
  echo -e "${YELLOW}Контейнер MinIO не найден, используем 'minio'${NC}"
  MINIO_CONTAINER="minio"
else
  echo -e "${YELLOW}Найден контейнер MinIO: $MINIO_CONTAINER${NC}"
fi

# Определение opensearch контейнера
OPENSEARCH_CONTAINER=$(docker ps --format "{{.Names}}" | grep -i "opensearch" | head -n 1)
if [ -z "$OPENSEARCH_CONTAINER" ]; then
  echo -e "${YELLOW}Контейнер OpenSearch не найден, используем 'opensearch'${NC}"
  OPENSEARCH_CONTAINER="opensearch"
else
  echo -e "${YELLOW}Найден контейнер OpenSearch: $OPENSEARCH_CONTAINER${NC}"
fi

# Получаем IP-адреса сервисов (для резервного использования)
MINIO_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $MINIO_CONTAINER)
DB_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $DB_CONTAINER)
OPENSEARCH_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $OPENSEARCH_CONTAINER)

echo -e "${YELLOW}IP-адрес MinIO: $MINIO_IP${NC}"
echo -e "${YELLOW}IP-адрес DB: $DB_IP${NC}"
echo -e "${YELLOW}IP-адрес OpenSearch: $OPENSEARCH_IP${NC}"

# Получаем IP-подсеть этой сети
NETWORK_SUBNET=$(docker network inspect $HOSTEL_NETWORK | grep -oP '"Subnet": "\K[^"]+')
echo -e "${YELLOW}Сеть '$HOSTEL_NETWORK' использует подсеть: $NETWORK_SUBNET${NC}"

echo -e "${YELLOW}Определяем текущую конфигурацию системы...${NC}"

# Проверяем текущие контейнеры
BLUE_RUNNING=$(docker ps --filter name=backend-blue --filter status=running -q | wc -l)
GREEN_RUNNING=$(docker ps --filter name=backend-green --filter status=running -q | wc -l)
ORIGINAL_RUNNING=$(docker ps --filter name=backend --filter status=running -q | wc -l)

# Вывод информации о существующих контейнерах
echo -e "${YELLOW}Текущие контейнеры:${NC}"
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

echo -e "${YELLOW}Текущий активный контейнер: $CURRENT_COLOR ($CURRENT_CONTAINER)${NC}"
echo -e "${YELLOW}Будет создан новый контейнер: $NEW_COLOR ($NEW_CONTAINER) на порту $NEW_PORT${NC}"

# Останавливаем существующий контейнер для нового цвета, если он есть
if docker ps -a --filter name=$NEW_CONTAINER -q | grep -q .; then
  echo -e "${YELLOW}Останавливаем существующий контейнер $NEW_CONTAINER...${NC}"
  docker stop $NEW_CONTAINER 2>/dev/null || true
  docker rm $NEW_CONTAINER 2>/dev/null || true
fi

# Авторизация в Harbor
echo -e "${YELLOW}Авторизация в Harbor...${NC}"
docker login -u admin -p SveTu2025 harbor.svetu.rs

# Получаем новый образ
echo -e "${YELLOW}Загрузка нового образа backend:latest...${NC}"
docker pull harbor.svetu.rs/svetu/backend/api:latest

# Создаем специальную конфигурацию для blue-green deployment с правильным разделением WebSocket и API
cat > $WORK_DIR/bluegreen.env << EOF
# Переменные окружения для blue-green деплоя
APP_MODE=production
WS_ENABLED=true
FILE_STORAGE_PROVIDER=minio
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=listings
MINIO_LOCATION=eu-central-1
FILE_STORAGE_PUBLIC_URL=https://svetu.rs
PORT=3000
SERVER_HOST=https://svetu.rs
EOF

# Получаем IP-адрес сети хоста
HOST_IP=$(hostname -I | awk '{print $1}')
echo -e "${YELLOW}IP-адрес хоста: $HOST_IP${NC}"

# Запускаем новый контейнер - ВАЖНО! Сразу указываем правильную сеть!
echo -e "${YELLOW}Запуск нового контейнера $NEW_CONTAINER на порту $NEW_PORT в сети $HOSTEL_NETWORK...${NC}"

# Используем комбинацию имен и IP-адресов (IP как резервный вариант)
docker run -d --name $NEW_CONTAINER \
  --network $HOSTEL_NETWORK \
  -p $NEW_PORT:3000 \
  -v /opt/hostel-data/uploads:/app/uploads \
  -v /opt/hostel-data/minio:/data/minio \
  -v /opt/hostel-data/credentials:/app/credentials \
  --env-file $WORK_DIR/bluegreen.env \
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
  harbor.svetu.rs/svetu/backend/api:latest

# Проверяем, что контейнер подключен к нужной сети
NETWORK_CHECK=$(docker inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}}{{end}}' $NEW_CONTAINER)
echo -e "${YELLOW}Контейнер $NEW_CONTAINER подключен к сети: $NETWORK_CHECK${NC}"

# Получаем IP-адрес нового контейнера
echo -e "${YELLOW}Ожидаем запуск контейнера (30 секунд)...${NC}"
sleep 30

NEW_CONTAINER_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $NEW_CONTAINER)
echo -e "${YELLOW}IP-адрес нового контейнера $NEW_CONTAINER: $NEW_CONTAINER_IP${NC}"

# Вывод логов контейнера для диагностики
echo -e "${YELLOW}Логи нового контейнера:${NC}"
docker logs --tail 30 $NEW_CONTAINER

# Проверка сетевого подключения между контейнерами
echo -e "${YELLOW}Проверка связи между backend и minio...${NC}"
docker exec $NEW_CONTAINER sh -c "ping -c 1 $MINIO_CONTAINER || ping -c 1 $MINIO_IP || echo 'Ping к Minio не прошел'"
echo -e "${YELLOW}Проверка связи между backend и db...${NC}"
docker exec $NEW_CONTAINER sh -c "ping -c 1 $DB_CONTAINER || ping -c 1 $DB_IP || echo 'Ping к DB не прошел'"

# Проверка работоспособности напрямую
MAX_RETRIES=15
RETRY_COUNT=0
HEALTH_CHECK_OK=false

echo -e "${YELLOW}Проверка работоспособности нового контейнера по API endpoint...${NC}"
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  # Пробуем более простой эндпоинт /health сначала
  HEALTH_RESPONSE=$(curl -s --max-time 5 http://$NEW_CONTAINER_IP:3000/api/health || echo "Failed")

  if [[ "$HEALTH_RESPONSE" == *"OK"* ]]; then
    echo -e "${GREEN}Контейнер отвечает на базовые запросы!${NC}"
    HEALTH_CHECK_OK=true
    break
  fi

  # Пробуем получить категории если /health не отвечает
  CATEGORY_RESPONSE=$(curl -s --max-time 5 http://$NEW_CONTAINER_IP:3000/api/v1/marketplace/category-tree | head -c 50 || echo "Failed")

  if [[ "$CATEGORY_RESPONSE" == *"categories"* ]] || [[ "$CATEGORY_RESPONSE" == *"id"* ]]; then
    echo -e "${GREEN}Новый контейнер успешно запущен и отвечает на запросы API!${NC}"
    HEALTH_CHECK_OK=true
    break
  fi

  echo -e "${YELLOW}Ожидание ответа от нового контейнера... ($RETRY_COUNT/$MAX_RETRIES)${NC}"
  echo -e "${YELLOW}Ответ health: $HEALTH_RESPONSE${NC}"
  echo -e "${YELLOW}Ответ категорий: $CATEGORY_RESPONSE${NC}"

  # Дополнительная диагностика на каждой итерации
  if (( RETRY_COUNT % 3 == 0 )); then
    echo -e "${YELLOW}Последние логи контейнера:${NC}"
    docker logs --tail 10 $NEW_CONTAINER

    # Проверка соединения с базой данных
    echo -e "${YELLOW}Тест соединения с базой данных:${NC}"
    docker exec $NEW_CONTAINER sh -c "echo 'Тест соединения с DB: $DB_CONTAINER'"
    docker exec $NEW_CONTAINER sh -c "ping -c 1 $DB_CONTAINER || ping -c 1 $DB_IP || echo 'Ping failed'"
  fi

  sleep 7
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ "$HEALTH_CHECK_OK" != "true" ]; then
  echo -e "${RED}Ошибка: Новый контейнер не отвечает на запросы API!${NC}"
  echo -e "${RED}Полный лог контейнера:${NC}"
  docker logs $NEW_CONTAINER
  echo -e "${YELLOW}Останавливаем и удаляем новый контейнер...${NC}"
  docker stop $NEW_CONTAINER 2>/dev/null || true
  docker rm $NEW_CONTAINER 2>/dev/null || true
  exit 1
fi

# Проверка доступности по порту
echo -e "${YELLOW}Проверка доступности по порту $NEW_PORT...${NC}"
RETRY_COUNT=0
PORT_CHECK_OK=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  PORT_RESPONSE=$(curl -s --max-time 5 http://localhost:$NEW_PORT/api/health || echo "Failed")

  if [[ "$PORT_RESPONSE" == *"OK"* ]]; then
    echo -e "${GREEN}Новый контейнер успешно доступен по порту $NEW_PORT!${NC}"
    PORT_CHECK_OK=true
    break
  fi

  # Пробуем категории
  CATEGORY_PORT_RESPONSE=$(curl -s --max-time 5 http://localhost:$NEW_PORT/api/v1/marketplace/category-tree | head -c 50 || echo "Failed")

  if [[ "$CATEGORY_PORT_RESPONSE" == *"categories"* ]] || [[ "$CATEGORY_PORT_RESPONSE" == *"id"* ]]; then
    echo -e "${GREEN}Новый контейнер успешно доступен по порту $NEW_PORT!${NC}"
    PORT_CHECK_OK=true
    break
  fi

  echo -e "${YELLOW}Ожидание доступности по порту... ($RETRY_COUNT/$MAX_RETRIES)${NC}"
  echo -e "${YELLOW}Ответ health: $PORT_RESPONSE${NC}"
  echo -e "${YELLOW}Ответ категорий: $CATEGORY_PORT_RESPONSE${NC}"

  sleep 5
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ "$PORT_CHECK_OK" != "true" ]; then
  echo -e "${RED}Ошибка: Новый контейнер не доступен по порту $NEW_PORT!${NC}"
  echo -e "${RED}Полный лог контейнера:${NC}"
  docker logs $NEW_CONTAINER
  echo -e "${YELLOW}Останавливаем и удаляем новый контейнер...${NC}"
  docker stop $NEW_CONTAINER 2>/dev/null || true
  docker rm $NEW_CONTAINER 2>/dev/null || true
  exit 1
fi

# Создаем временный файл конфигурации с заменой переменных
cp /opt/hostel-booking-system/nginx.conf $WORK_DIR/nginx.conf.template

# Заменяем переменные в конфигурации nginx
sed -i "s/\$NEW_CONTAINER_IP/$NEW_CONTAINER_IP/g" $WORK_DIR/nginx.conf.template
sed -i "s/\$MINIO_CONTAINER/$MINIO_CONTAINER/g" $WORK_DIR/nginx.conf.template
sed -i "s/\$MINIO_IP/$MINIO_IP/g" $WORK_DIR/nginx.conf.template

# Копируем обработанный файл конфигурации в целевое место
cp $WORK_DIR/nginx.conf.template $WORK_DIR/nginx.conf.new

# Проверяем статус nginx
NGINX_RUNNING=$(docker ps --filter name=hostel_nginx --filter status=running -q | wc -l)
if [ "$NGINX_RUNNING" -eq "0" ]; then
  echo -e "${YELLOW}Nginx не запущен, пытаемся запустить...${NC}"
  docker start hostel_nginx || docker run -d --name hostel_nginx \
    --network $HOSTEL_NETWORK \
    -p 80:80 -p 443:443 \
    -v $WORK_DIR/nginx.conf.new:/etc/nginx/conf.d/default.conf \
    -v /opt/hostel-booking-system/frontend/hostel-frontend/build:/usr/share/nginx/html \
    -v /opt/hostel-booking-system/certbot/conf:/etc/letsencrypt \
    -v /opt/hostel-booking-system/certbot/www:/var/www/certbot \
    -v /opt/hostel-data/uploads:/usr/share/nginx/uploads \
    harbor.svetu.rs/svetu/nginx/nginx:latest

  sleep 5
  NGINX_RUNNING=$(docker ps --filter name=hostel_nginx --filter status=running -q | wc -l)
  if [ "$NGINX_RUNNING" -eq "0" ]; then
    echo -e "${RED}Ошибка: Не удалось запустить Nginx!${NC}"
    docker logs hostel_nginx
    # Тем не менее продолжаем, т.к. мы хотим обновить конфигурацию
  fi
fi

# Обновляем конфигурацию nginx
echo -e "${YELLOW}Обновляем конфигурацию nginx...${NC}"
cp $WORK_DIR/nginx.conf.new /opt/hostel-booking-system/nginx.conf

# Если nginx запущен, проверяем конфигурацию и перезапускаем
if [ "$NGINX_RUNNING" -gt "0" ]; then
  # Проверяем синтаксис nginx
  echo -e "${YELLOW}Проверяем синтаксис nginx...${NC}"

  if ! docker exec hostel_nginx nginx -t >/dev/null 2>&1; then
    echo -e "${RED}Ошибка в конфигурации nginx! Выводим подробную информацию...${NC}"
    docker exec hostel_nginx nginx -t
    echo -e "${YELLOW}Восстанавливаем оригинальную конфигурацию...${NC}"
    cp $WORK_DIR/nginx.conf.original /opt/hostel-booking-system/nginx.conf

    echo -e "${YELLOW}Останавливаем новый контейнер...${NC}"
    docker stop $NEW_CONTAINER
    docker rm $NEW_CONTAINER
    exit 1
  fi

  # Подключаем nginx к тем же сетям, что и бэкенд для надежности
  echo -e "${YELLOW}Проверка и подключение nginx к той же сети, что и бэкенд...${NC}"
  NGINX_NETWORKS=$(docker inspect hostel_nginx -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}} {{end}}')
  echo -e "${YELLOW}Сети Nginx: $NGINX_NETWORKS${NC}"

  # Проверяем, подключен ли nginx к той же сети, что и бэкенд
  if [[ ! "$NGINX_NETWORKS" == *"$HOSTEL_NETWORK"* ]]; then
    echo -e "${YELLOW}Подключаем nginx к сети $HOSTEL_NETWORK...${NC}"
    docker network connect $HOSTEL_NETWORK hostel_nginx || true
  fi

  # Перезапускаем nginx для применения изменений
  echo -e "${YELLOW}Перезапускаем nginx...${NC}"
  docker restart hostel_nginx

  # Также перезапустим бэкенд для сброса соединений
  docker restart $NEW_CONTAINER
  echo -e "${YELLOW}Перезапустили контейнер бэкенда для сброса соединений...${NC}"
  sleep 5

  # Проверяем доступность API через nginx
  echo -e "${YELLOW}Проверяем доступность API через nginx...${NC}"
  sleep 10
  MAX_RETRIES=15
  RETRY_COUNT=0
  API_CHECK_OK=false

  while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    # Проверяем сначала /health
    API_HEALTH=$(curl -k -s https://svetu.rs/api/health || echo "Failed")

    if [[ "$API_HEALTH" == *"OK"* ]]; then
      echo -e "${GREEN}API health endpoint успешно доступен через nginx!${NC}"
      API_CHECK_OK=true
      break
    fi

    # Проверяем categories API
    API_RESPONSE=$(curl -k -s https://svetu.rs/api/v1/marketplace/category-tree | head -c 50 || echo "Failed")
    if [[ "$API_RESPONSE" == *"categories"* ]] || [[ "$API_RESPONSE" == *"id"* ]]; then
      echo -e "${GREEN}API успешно доступен через nginx!${NC}"
      API_CHECK_OK=true
      break
    fi

    echo -e "${YELLOW}Ожидание доступности API через nginx... ($RETRY_COUNT/$MAX_RETRIES)${NC}"
    echo -e "${YELLOW}Ответ health: $API_HEALTH${NC}"
    echo -e "${YELLOW}Ответ categories: $API_RESPONSE${NC}"

    # Дополнительная диагностика
    if (( RETRY_COUNT % 3 == 0 )); then
      echo -e "${YELLOW}Проверяем работающие запросы к API внутри контейнера nginx...${NC}"
      docker exec hostel_nginx curl -s $NEW_CONTAINER_IP:3000/api/v1/marketplace/category-tree | head -c 50 || echo "Failed"
    fi

    sleep 5
    RETRY_COUNT=$((RETRY_COUNT+1))
  done

  if [ "$API_CHECK_OK" != "true" ]; then
    echo -e "${RED}Ошибка: API не доступен через nginx! Восстанавливаем оригинальную конфигурацию...${NC}"
    cp $WORK_DIR/nginx.conf.original /opt/hostel-booking-system/nginx.conf
    docker restart hostel_nginx

    echo -e "${YELLOW}Останавливаем новый контейнер...${NC}"
    docker stop $NEW_CONTAINER
    docker rm $NEW_CONTAINER
    exit 1
  fi

  # Проверка доступности WebSocket (эту часть можно расширить)
  echo -e "${YELLOW}Проверка доступности WebSocket будет проведена вручную.${NC}"

else
  echo -e "${YELLOW}Nginx не запущен, сохраняем конфигурацию, но не перезапускаем.${NC}"
  echo -e "${YELLOW}Для ручного запуска выполните: docker start hostel_nginx${NC}"
fi

# Если проверки прошли успешно, останавливаем старый контейнер
if [ "$CURRENT_COLOR" = "original" ]; then
  echo -e "${YELLOW}Останавливаем оригинальный контейнер backend...${NC}"
  docker stop backend
  docker rm backend
elif [ "$CURRENT_COLOR" != "none" ]; then
  echo -e "${YELLOW}Останавливаем предыдущий контейнер $CURRENT_CONTAINER...${NC}"
  docker stop $CURRENT_CONTAINER
  docker rm $CURRENT_CONTAINER
fi

echo -e "${GREEN}Бесшовное обновление backend успешно завершено!${NC}"
echo -e "${GREEN}Новый активный контейнер: $NEW_CONTAINER с IP: $NEW_CONTAINER_IP${NC}"

# Очистка
rm -rf $WORK_DIR
EOT

  # Делаем скрипт исполняемым
  chmod +x /tmp/server_deploy.sh

  # Отправляем скрипт на сервер
  scp /tmp/server_deploy.sh $PROD_SERVER_USER@$PROD_SERVER:$PROD_SERVER_PATH/server_deploy.sh

  # Запускаем скрипт на сервере
  ssh $PROD_SERVER_USER@$PROD_SERVER "cd $PROD_SERVER_PATH && chmod +x server_deploy.sh && ./server_deploy.sh"

  # Удаляем временный скрипт
  rm /tmp/server_deploy.sh
}

# Функция для бесшовного обновления frontend
blue_green_frontend() {
  echo -e "${YELLOW}Бесшовное обновление frontend...${NC}"

  # Создаем скрипт для обновления frontend с поддержкой env.js
  cat > /tmp/frontend_deploy.sh << 'EOT'
#!/bin/bash
set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Создаем рабочие директории
WORK_DIR="/tmp/bluegreen"
rm -rf $WORK_DIR
mkdir -p $WORK_DIR
BACKUP_DIR="$WORK_DIR/frontend_backup"
mkdir -p $BACKUP_DIR

# Сохраняем текущую версию frontend
echo -e "${YELLOW}Создаем резервную копию frontend...${NC}"
mkdir -p /opt/hostel-booking-system/frontend/hostel-frontend/build.backup/
cp -r /opt/hostel-booking-system/frontend/hostel-frontend/build/* /opt/hostel-booking-system/frontend/hostel-frontend/build.backup/

# Авторизация в Harbor
echo -e "${YELLOW}Авторизация в Harbor...${NC}"
docker login -u admin -p SveTu2025 harbor.svetu.rs

# Загружаем последнюю версию frontend
echo -e "${YELLOW}Загрузка нового образа frontend...${NC}"
docker pull harbor.svetu.rs/svetu/frontend/app:latest

# Создаем временный контейнер для извлечения сборки
echo -e "${YELLOW}Создаем временный контейнер...${NC}"
TEMP_CONTAINER=$(docker create harbor.svetu.rs/svetu/frontend/app:latest)

# Извлекаем сборку из контейнера
echo -e "${YELLOW}Извлекаем сборку из контейнера...${NC}"
docker cp $TEMP_CONTAINER:/app/build/. $WORK_DIR/frontend/

# Проверяем наличие index.html
if [ ! -f "$WORK_DIR/frontend/index.html" ]; then
  echo -e "${RED}Ошибка: index.html не найден! Отмена обновления...${NC}"
  docker rm $TEMP_CONTAINER
  exit 1
fi

# Создаем env.js с правильными настройками для динамической конфигурации
echo -e "${YELLOW}Создаем env.js с нужными настройками...${NC}"
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

# Проверяем и создаем директории
echo -e "${YELLOW}Проверяем и создаем директории...${NC}"
mkdir -p /opt/hostel-booking-system/frontend/hostel-frontend/build/

# Обновляем frontend
echo -e "${YELLOW}Обновляем frontend...${NC}"
rm -rf /opt/hostel-booking-system/frontend/hostel-frontend/build/*
cp -r $WORK_DIR/frontend/* /opt/hostel-booking-system/frontend/hostel-frontend/build/

# Удостоверяемся, что env.js правильно скопирован
cp $WORK_DIR/frontend/env.js /opt/hostel-booking-system/frontend/hostel-frontend/build/env.js

# Проверяем результат обновления
if [ ! -f "/opt/hostel-booking-system/frontend/hostel-frontend/build/index.html" ]; then
  echo -e "${RED}Ошибка: Не удалось обновить frontend! Восстанавливаем из резервной копии...${NC}"
  rm -rf /opt/hostel-booking-system/frontend/hostel-frontend/build/*
  cp -r /opt/hostel-booking-system/frontend/hostel-frontend/build.backup/* /opt/hostel-booking-system/frontend/hostel-frontend/build/
  docker rm $TEMP_CONTAINER
  exit 1
fi

# Проверка правильности env.js и создание process-env.js (обратная совместимость)
if [ ! -f "/opt/hostel-booking-system/frontend/hostel-frontend/build/env.js" ]; then
  echo -e "${YELLOW}env.js не обнаружен, создаем его...${NC}"
  cat > /opt/hostel-booking-system/frontend/hostel-frontend/build/env.js << 'EOF'
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
fi

# Создаем копию env.js как process-env.js для обратной совместимости
cp /opt/hostel-booking-system/frontend/hostel-frontend/build/env.js /opt/hostel-booking-system/frontend/hostel-frontend/build/process-env.js

# Проверяем и устанавливаем правильные права
echo -e "${YELLOW}Устанавливаем правильные права на файлы...${NC}"
chmod -R 755 /opt/hostel-booking-system/frontend/hostel-frontend/build/

# Очистка
echo -e "${YELLOW}Очистка временных файлов...${NC}"
docker rm $TEMP_CONTAINER
rm -rf $WORK_DIR

# Перезапуск nginx для применения изменений
if docker ps -q --filter "name=hostel_nginx" | grep -q .; then
  echo -e "${YELLOW}Перезапуск nginx для применения изменений...${NC}"
  docker restart hostel_nginx

  # Проверка доступности frontend
  echo -e "${YELLOW}Проверка доступности frontend...${NC}"
  sleep 5
  FRONTEND_RESPONSE=$(curl -k -s https://svetu.rs/ | grep -c "id=\"root\"" || echo "0")
  if [[ "$FRONTEND_RESPONSE" -gt "0" ]]; then
    echo -e "${GREEN}Frontend успешно доступен!${NC}"
  else
    echo -e "${YELLOW}Предупреждение: возможно, frontend недоступен. Проверьте вручную.${NC}"
  fi
fi

echo -e "${GREEN}Обновление frontend успешно завершено!${NC}"
EOT

  # Делаем скрипт исполняемым
  chmod +x /tmp/frontend_deploy.sh
  
  # Отправляем скрипт на сервер
  scp /tmp/frontend_deploy.sh $PROD_SERVER_USER@$PROD_SERVER:$PROD_SERVER_PATH/frontend_deploy.sh
  
  # Запускаем скрипт на сервере
  ssh $PROD_SERVER_USER@$PROD_SERVER "cd $PROD_SERVER_PATH && chmod +x frontend_deploy.sh && ./frontend_deploy.sh"
  
  # Удаляем временный скрипт
  rm /tmp/frontend_deploy.sh
}

# Авторизация в Harbor
echo -e "${YELLOW}Авторизация в Harbor...${NC}"
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Проверка статуса авторизации
if [ $? -ne 0 ]; then
  echo -e "${RED}Ошибка авторизации в Harbor. Проверьте учетные данные.${NC}"
  exit 1
fi

# Обработка в зависимости от сервиса
case "$SERVICE" in
  "backend")
    build_backend
    blue_green_backend
    ;;
  "frontend")
    build_frontend
    blue_green_frontend
    ;;
  "all")
    build_backend
    build_frontend
    blue_green_backend
    blue_green_frontend
    ;;
  *)
    echo -e "${RED}Ошибка: Неизвестный сервис '$SERVICE'. Используйте backend, frontend или all.${NC}"
    exit 1
    ;;
esac

echo -e "${GREEN}Процесс бесшовного обновления успешно завершен!${NC}"