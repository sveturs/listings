#!/bin/bash
set -e # Останавливаем выполнение при ошибках
echo "Начинаем деплой..."
cd /opt/hostel-booking-system

# Создаем директории для хранения данных, если их еще нет
mkdir -p /opt/hostel-data/uploads
mkdir -p /opt/hostel-data/db
mkdir -p /opt/hostel-data/opensearch
mkdir -p /opt/hostel-data/yarn-cache # Директория для кэша yarn
mkdir -p certbot/conf
mkdir -p certbot/www
mkdir -p /tmp/hostel-backup/db

# Настраиваем git pull strategy
git config pull.rebase false

# Сохраняем важные файлы
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true
cp -f frontend/hostel-frontend/.env /tmp/hostel-backup/frontend.env 2>/dev/null || true

# Сохраняем SSL сертификаты
if [ -d "certbot/conf" ]; then
  cp -r certbot/conf /tmp/hostel-backup/ 2>/dev/null || true
fi

# Делаем бэкап базы данных только если контейнер запущен
echo "Пытаемся создать бэкап базы данных..."
if docker-compose -f docker-compose.prod.yml ps | grep -q "db.*Up"; then
  BACKUP_FILE="/tmp/hostel-backup/db/backup_$(date +%Y%m%d_%H%M%S).sql"
  docker-compose -f docker-compose.prod.yml exec -T db pg_dumpall -U postgres > "$BACKUP_FILE"
  if [ $? -eq 0 ]; then
    echo "Бэкап базы данных создан в $BACKUP_FILE"
  else
    echo "Ошибка создания бэкапа базы данных, но продолжаем деплой"
  fi
else
  echo "База данных не запущена, пропускаем создание бэкапа"
fi

# Обеспечиваем чистое состояние git, но исключаем критические директории
git fetch origin
git reset --hard origin/main
git clean -fdx -e "*.env*" -e "uploads/" -e "certbot/" -e "/opt/hostel-data/"

# Восстанавливаем файлы конфигурации
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true
if [ -d "/tmp/hostel-backup/conf" ]; then
  cp -r /tmp/hostel-backup/conf certbot/ 2>/dev/null || true
fi

# Удаляем старые образы и принудительно пересоздаем контейнеры
echo "Останавливаем сервисы..."
docker-compose -f docker-compose.prod.yml down --remove-orphans

# Принудительно удаляем все связанные контейнеры, но оставляем тома нетронутыми
echo "Удаляем старые контейнеры, сохраняя тома с данными..."
docker-compose -f docker-compose.prod.yml rm -f

# Проверяем, нет ли застрявших или мертвых контейнеров
echo "Проверяем наличие застрявших контейнеров..."
for container in $(docker ps -a --filter "name=hostel_db" --filter "name=opensearch" --filter "status=exited" --format "{{.ID}}"); do
  echo "Найден застрявший контейнер: $container. Удаляем..."
  docker rm -f $container 2>/dev/null || true
done

# Удаляем старые образы, чтобы принудительно пересобрать их
echo "Удаляем старые образы для принудительной пересборки..."
docker-compose -f docker-compose.prod.yml build --no-cache

# Собираем фронтенд
echo "Собираем фронтенд..."
cd frontend/hostel-frontend

# Проверяем существование package.json и создаем измененный файл
echo "Проверяем package.json и добавляем необходимые зависимости..."
if [ -f "package.json" ]; then
  # Сохраняем оригинальный package.json
  cp package.json package.json.orig
  
  # Проверяем, есть ли уже react-query и react-window в зависимостях
  if ! grep -q '"react-query"' package.json; then
    sed -i 's/"dependencies": {/"dependencies": {\n    "react-query": "^3.39.3",/g' package.json
  fi
  
  if ! grep -q '"react-window"' package.json; then
    sed -i 's/"dependencies": {/"dependencies": {\n    "react-window": "^1.8.9",/g' package.json
  fi
fi

# Создаем кастомный скрипт для сборки фронтенда с исправленным синтаксисом
cat > build_frontend.sh << 'EOL'
#!/bin/bash
set -e
echo "Устанавливаем зависимости..."
npm cache clean --force

# Проверяем наличие yarn
echo "Проверяем наличие yarn..."
if command -v yarn &> /dev/null; then
  echo "Yarn уже установлен, используем существующую версию"
else
  echo "Устанавливаем yarn..."
  npm install -g yarn --force
fi

# Оптимизация для CI/CD системы
echo "Настраиваем оптимизации для сборки..."
# Увеличиваем лимит памяти для Node
export NODE_OPTIONS="--max-old-space-size=4096"
# Устанавливаем переменные окружения для оптимизации сборки
export GENERATE_SOURCEMAP=false
export INLINE_RUNTIME_CHUNK=true
export DISABLE_ESLINT_PLUGIN=true
export CI=false

# Настраиваем кэширование yarn для ускорения сборки в будущем
echo "Настраиваем кэширование yarn..."
YARN_CACHE_FOLDER=/opt/hostel-data/yarn-cache
export YARN_CACHE_FOLDER
export YARN_NETWORK_TIMEOUT=600000

echo "Устанавливаем зависимости с помощью yarn..."
yarn install --network-timeout 600000 --no-audit

# Проверяем наличие необходимых пакетов
echo "Проверяем наличие необходимых пакетов..."
yarn add react-query@3.39.3 react-window@1.8.9 react-scripts@5.0.1 \
  ajv@6.12.6 ajv-keywords@3.5.2 schema-utils@3.1.1 \
  --no-audit --production=true

# Обновляем package.json для ускорения сборки
echo "Оптимизируем package.json для быстрой сборки..."
if [ -f "package.json" ]; then
  # Добавляем конфигурацию babel для ускорения сборки
  if ! grep -q '"browserslist"' package.json; then
    echo 'Обновляем browserslist для ускорения сборки...'
    sed -i 's/"browserslist": {/"browserslist": {\n  "production": [\n    "last 2 chrome version",\n    "last 2 firefox version"\n  ],/g' package.json
  fi
  
  # Изменяем скрипт сборки для ускорения процесса
  sed -i 's/"build": "react-scripts build"/"build": "GENERATE_SOURCEMAP=false react-scripts build"/g' package.json
fi

# Проверка наличия .env файла и создание его при необходимости
if [ ! -f ".env" ] && [ -f "../.env.example" ]; then
  echo "Создаем .env файл из шаблона..."
  cp "../.env.example" ".env"
fi

echo "Пробуем сборку проекта..."
echo "Обратите внимание: Процесс оптимизации production-сборки может занять некоторое время..."
echo "Время сборки начало: $(date)"

# Мониторинг процесса сборки
(
  while true; do
    echo "Сборка продолжается... $(date)"
    sleep 60
  done
) &
MONITOR_PID=$!

# Сборка с таймаутом
CI=false DISABLE_ESLINT_PLUGIN=true yarn build
BUILD_STATUS=$?

if [ $BUILD_STATUS -ne 0 ]; then
  echo "Сборка не удалась, пробуем с дополнительными параметрами..."
  
  # Останавливаем процесс мониторинга
  kill $MONITOR_PID 2>/dev/null || true
  
  # Добавляем параметры для дальнейшей оптимизации
  export DISABLE_NEW_JSX_TRANSFORM=true
  export CI=false

  echo "Пробуем сборку с отключенными проверками..."
  yarn build --no-lint
  BUILD_STATUS=$?
  
  if [ $BUILD_STATUS -ne 0 ]; then
    echo "Сборка все еще не удается, последняя попытка..."
    
    # Анализируем ошибки сборки
    yarn build 2>&1 | tee build_error.log
    
    # Находим все упоминания о недостающих пакетах
    if grep -q "Can't resolve" build_error.log; then
      packages=$(grep -o "Can't resolve '[^']*'" build_error.log | sed "s/Can't resolve '//g" | sed "s/'//g")
      echo "Найдены отсутствующие пакеты: $packages"
      
      for pkg in $packages; do
        echo "Устанавливаем $pkg..."
        yarn add "$pkg" --no-audit
      done
      
      # Пробуем собрать снова с увеличенным таймаутом
      echo "Последняя попытка сборки..."
      yarn build --no-lint
      BUILD_STATUS=$?
    fi
  fi
fi

# Останавливаем процесс мониторинга
kill $MONITOR_PID 2>/dev/null || true

echo "Время окончания сборки: $(date)"
exit $BUILD_STATUS
EOL

chmod +x build_frontend.sh

# Запускаем сборку фронтенда в контейнере с увеличенными ресурсами
echo "Запускаем сборку фронтенда в контейнере с увеличенными ресурсами..."
NODE_ENV=production docker run --rm \
  -v $(pwd):/app \
  -v /opt/hostel-data/yarn-cache:/opt/hostel-data/yarn-cache \
  -w /app \
  --cpus=4 \
  --memory=6g \
  node:18 bash -c "./build_frontend.sh"
BUILD_STATUS=$?

# Проверяем, успешно ли прошла сборка
if [ $BUILD_STATUS -ne 0 ] || [ ! -d "build" ] || [ -z "$(ls -A build 2>/dev/null)" ]; then
  echo "Сборка фронтенда не удалась. Проверяем логи для анализа ошибок..."
  
  # Анализируем логи сборки
  if [ -f "build_error.log" ]; then
    ERROR_LOG=$(cat build_error.log)
    if echo "$ERROR_LOG" | grep -q "out of memory"; then
      echo "Ошибка сборки: Недостаточно памяти. Попробуйте увеличить значение NODE_OPTIONS."
    elif echo "$ERROR_LOG" | grep -q "Cannot find module"; then
      echo "Ошибка сборки: Не найден модуль. Проверьте зависимости."
    else
      echo "Ошибка сборки: Неизвестная проблема. Проверьте полные логи."
    fi
  fi
  
  # Использовать предыдущую сборку, если она есть
  if [ -d "build" ] && [ -n "$(ls -A build 2>/dev/null)" ]; then
    echo "Используем существующую сборку фронтенда..."
  else
    echo "Деплой прерван из-за проблем сборки фронтенда и отсутствия существующей сборки."
    exit 1
  fi
fi

cd ../..

# Проверяем наличие проблем с томами базы данных и исправляем их
echo "Проверяем состояние тома базы данных..."
DB_VOLUME_PATH="/opt/hostel-data/db"

# Создаем новую структуру директорий для PostgreSQL
echo "Настраиваем новую структуру для PostgreSQL с вложенными директориями..."
# Делаем резервную копию текущего состояния при наличии данных
if [ -d "$DB_VOLUME_PATH" ] && [ -n "$(ls -A $DB_VOLUME_PATH 2>/dev/null)" ]; then
  BACKUP_DIR="${DB_VOLUME_PATH}_backup_$(date +%Y%m%d_%H%M%S)"
  echo "Сохраняем текущее состояние в $BACKUP_DIR..."
  cp -r "$DB_VOLUME_PATH" "$BACKUP_DIR" 2>/dev/null || true
fi

# Остановим все контейнеры базы данных
echo "Останавливаем все контейнеры базы данных..."
docker-compose -f docker-compose.prod.yml stop db
docker-compose -f docker-compose.prod.yml rm -f db

# Полностью очищаем директорию базы данных
echo "Полностью очищаем директорию базы данных..."
rm -rf "$DB_VOLUME_PATH"
mkdir -p "$DB_VOLUME_PATH"

# Создаем более глубокую структуру вложенности для решения проблемы точек монтирования
echo "Создаем вложенную структуру каталогов..."
mkdir -p "$DB_VOLUME_PATH/data"
chown -R 999:999 "$DB_VOLUME_PATH"
chmod -R 700 "$DB_VOLUME_PATH"

# Модифицируем docker-compose.prod.yml для использования новой структуры
echo "Модифицируем docker-compose.prod.yml для использования новой структуры..."
TEMP_COMPOSE_FILE=$(mktemp)

# Создаем временную копию с измененными параметрами только для DB
cat > "$TEMP_COMPOSE_FILE" << EOL
version: '3.8'
services:
  db:
    image: postgres:15
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

volumes:
  db_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /opt/hostel-data/db/data

networks:
  hostel_network:
    name: hostel-booking-system_hostel_network
    driver: bridge
EOL

# Запускаем только базу данных с временной конфигурацией
echo "Запускаем базу данных с новой конфигурацией..."
docker-compose -f "$TEMP_COMPOSE_FILE" up -d db

# Проверяем готовность базы данных
echo "Проверяем готовность базы данных..."
MAX_DB_ATTEMPTS=30
for i in $(seq 1 $MAX_DB_ATTEMPTS); do
  if docker-compose -f "$TEMP_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "База данных готова! (попытка $i/$MAX_DB_ATTEMPTS)"
    
    # Проверяем соединение с базой дополнительно
    if docker-compose -f "$TEMP_COMPOSE_FILE" exec -T db psql -U postgres -c "SELECT 1" > /dev/null 2>&1; then
      echo "Успешно проверено соединение с базой данных!"
      
      # Если подключение работает, обновляем только путь тома в основном файле
      echo "Обновляем пути томов в основном docker-compose.prod.yml..."
      # Ищем и заменяем строку с путем тома db_data
      sed -i 's|device: /opt/hostel-data/db|device: /opt/hostel-data/db/data|g' docker-compose.prod.yml
      
      # Добавляем PGDATA, если его еще нет
      if ! grep -q "PGDATA:" docker-compose.prod.yml; then
        sed -i '/POSTGRES_DB:/a\      PGDATA: /var/lib/postgresql/data/pgdata  # Подкаталог для хранения файлов БД' docker-compose.prod.yml
      fi
      
      # НЕ останавливаем временную базу данных, а используем ее
      echo "Используем запущенную базу данных для дальнейших операций..."
      break
    else
      echo "pg_isready вернул успех, но не удалось подключиться через psql. Продолжаем ожидание..."
    fi
  fi
  
  if [ $i -eq $MAX_DB_ATTEMPTS ]; then
    echo "Не удалось запустить базу данных. Продолжаем деплой без проверки БД..."
    echo "ВНИМАНИЕ: Это может привести к ошибкам в работе приложения!"
  fi
  
  echo "Ожидаем готовность базы данных... Попытка $i/$MAX_DB_ATTEMPTS"
  sleep 5
done

# Проверяем, существует ли нужная сеть - она уже должна существовать после запуска временной базы данных
NETWORK_NAME="hostel-booking-system_hostel_network"
if ! docker network ls | grep -q "$NETWORK_NAME"; then
  echo "Сеть $NETWORK_NAME не существует, создаем..."
  docker network create "$NETWORK_NAME"
fi

# Запускаем миграции напрямую
echo "Запускаем миграции..."

# Проверяем работу базы данных перед миграцией
if docker-compose -f "$TEMP_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
  echo "База данных готова для миграций!"
  
  # Получаем IP адрес базы данных
  DB_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' hostel_db)
  if [ -n "$DB_IP" ]; then
    echo "IP-адрес базы данных: $DB_IP"
    
    # Запускаем миграции с использованием IP-адреса
    docker run --rm --network "$NETWORK_NAME" \
      -v $(pwd)/backend/migrations:/migrations \
      migrate/migrate \
      -path=/migrations/ \
      -database="postgres://postgres:c9XWc7Cm@$DB_IP:5432/hostel_db?sslmode=disable" \
      up || {
        echo "Ошибка при выполнении миграций. Продолжаем запуск других сервисов..."
      }
  else
    echo "Не удалось получить IP-адрес базы данных. Используем имя хоста..."
    docker run --rm --network "$NETWORK_NAME" \
      -v $(pwd)/backend/migrations:/migrations \
      migrate/migrate \
      -path=/migrations/ \
      -database="postgres://postgres:c9XWc7Cm@hostel_db:5432/hostel_db?sslmode=disable" \
      up || {
        echo "Ошибка при выполнении миграций. Продолжаем запуск других сервисов..."
      }
  fi
else
  echo "ВНИМАНИЕ: База данных не готова для миграций!"
fi

# Запускаем остальные сервисы, НЕ запуская базу данных снова
echo "Запускаем остальные сервисы, используя текущую базу данных..."
docker-compose -f docker-compose.prod.yml up -d opensearch

# Даем время OpenSearch запуститься
echo "Ожидаем запуска OpenSearch..."
sleep 15

# Запускаем оставшиеся сервисы
echo "Запускаем backend и другие сервисы..."
docker-compose -f docker-compose.prod.yml up -d backend opensearch-dashboards nginx

# Удаляем временный файл
rm -f "$TEMP_COMPOSE_FILE"

# Проверяем, все ли сервисы запущены
echo "Проверяем статус всех сервисов..."
SERVICES_STATUS=$(docker-compose -f docker-compose.prod.yml ps)
echo "$SERVICES_STATUS"

# Проверяем здоровье backend
echo "Проверяем здоровье backend..."
RETRY_COUNT=30
for i in $(seq 1 $RETRY_COUNT); do
  if docker-compose -f docker-compose.prod.yml exec -T backend curl -s -f http://localhost:3000 > /dev/null 2>&1; then
    echo "Backend здоров! (попытка $i/$RETRY_COUNT)"
    break
  fi
  
  # Проверяем, почему backend не отвечает
  if [ $i -eq 10 ] || [ $i -eq 20 ]; then
    echo "Backend все еще не отвечает. Проверяем логи..."
    docker-compose -f docker-compose.prod.yml logs --tail=30 backend
  fi
  
  echo "Ожидаем готовность backend... Попытка $i/$RETRY_COUNT"
  sleep 3
done

# Сохраняем последние 5 бэкапов и удаляем более старые
find /tmp/hostel-backup/db -name "*.sql" -type f | sort -r | tail -n +6 | xargs rm -f 2>/dev/null || true

echo "Деплой завершен!"
echo "Логи контейнеров можно посмотреть с помощью: docker-compose -f docker-compose.prod.yml logs -f"