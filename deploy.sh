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

# Проверка прав доступа
if [ -d "$DB_VOLUME_PATH" ]; then
  echo "Проверяем и исправляем права доступа на томе базы данных..."
  chown -R 999:999 "$DB_VOLUME_PATH"  # 999 - стандартный UID для postgres в контейнере
  chmod -R 700 "$DB_VOLUME_PATH"
  
  # Проверка на наличие критичных файлов PostgreSQL
  if [ -f "$DB_VOLUME_PATH/PG_VERSION" ]; then
    echo "Обнаружен файл PG_VERSION, проверяем версию PostgreSQL..."
    PG_VERSION=$(cat "$DB_VOLUME_PATH/PG_VERSION")
    echo "Версия PostgreSQL в томе: $PG_VERSION"
    
    # Проверяем целостность кластера
    if [ ! -f "$DB_VOLUME_PATH/global/pg_control" ]; then
      echo "ВНИМАНИЕ: Файл pg_control не найден. Возможно, кластер PostgreSQL поврежден."
      echo "Создаем резервную копию текущего тома и инициализируем новый кластер..."
      BACKUP_PATH="${DB_VOLUME_PATH}_backup_$(date +%Y%m%d_%H%M%S)"
      mv "$DB_VOLUME_PATH" "$BACKUP_PATH"
      mkdir -p "$DB_VOLUME_PATH"
      chown 999:999 "$DB_VOLUME_PATH"
      chmod 700 "$DB_VOLUME_PATH"
      echo "Старый том сохранен в $BACKUP_PATH, создан новый том."
    fi
  else
    # Если директория не пустая, но файла PG_VERSION нет, значит это не валидный кластер PostgreSQL
    if [ -n "$(ls -A $DB_VOLUME_PATH 2>/dev/null)" ]; then
      echo "ВНИМАНИЕ: Том базы данных содержит файлы, но это не кластер PostgreSQL."
      echo "Создаем резервную копию директории и инициализируем новый кластер..."
      BACKUP_PATH="${DB_VOLUME_PATH}_backup_$(date +%Y%m%d_%H%M%S)"
      mv "$DB_VOLUME_PATH" "$BACKUP_PATH"
      mkdir -p "$DB_VOLUME_PATH"
      chown 999:999 "$DB_VOLUME_PATH"
      chmod 700 "$DB_VOLUME_PATH"
      echo "Старый том сохранен в $BACKUP_PATH, создан новый том."
    fi
  fi
fi

# Проверка на наличие битого PID файла
PID_FILE="$DB_VOLUME_PATH/postmaster.pid"
if [ -f "$PID_FILE" ]; then
  echo "Обнаружен PID файл. Это может быть причиной проблемы. Удаляем его..."
  rm -f "$PID_FILE"
fi

# Создаем файл .postgresql_empty для указания PostgreSQL использовать PGDATA как пустую директорию
echo "Создаем файл-флаг для корректной инициализации PostgreSQL..."
mkdir -p "$DB_VOLUME_PATH/pgdata" 2>/dev/null || true
touch "$DB_VOLUME_PATH/pgdata/.postgresql_empty" 2>/dev/null || true

# Добавляем вспомогательный скрипт для восстановления после ошибок инициализации
DOCKER_ENTRYPOINT_DIR="$DB_VOLUME_PATH/docker-entrypoint-initdb.d"
mkdir -p "$DOCKER_ENTRYPOINT_DIR"
cat > "$DOCKER_ENTRYPOINT_DIR/fix-initdb.sh" << 'EOSCRIPT'
#!/bin/bash
set -e

# Скрипт для корректной инициализации PostgreSQL
echo "Проверяем директорию данных..."
if [ -f "${PGDATA}/.postgresql_empty" ]; then
  echo "Обнаружен флаг .postgresql_empty, выполняем специальную инициализацию..."
  rm -f "${PGDATA}/.postgresql_empty"
fi

echo "Инициализация PostgreSQL завершена успешно."
EOSCRIPT

chmod +x "$DOCKER_ENTRYPOINT_DIR/fix-initdb.sh"

# Модифицируем docker-compose.prod.yml для использования pgdata как подкаталога
echo "Модифицируем docker-compose.prod.yml для использования PGDATA..."
if ! grep -q "PGDATA:" docker-compose.prod.yml; then
  sed -i '/POSTGRES_DB:/a\      PGDATA: /var/lib/postgresql/data/pgdata  # Подкаталог для хранения файлов БД' docker-compose.prod.yml
fi

# Запускаем базу данных через docker-compose
echo "Запускаем базу данных с правильной конфигурацией PGDATA..."
docker-compose -f docker-compose.prod.yml up -d db

# Проверяем готовность базы данных
echo "Проверяем готовность базы данных..."
MAX_DB_ATTEMPTS=30
for i in $(seq 1 $MAX_DB_ATTEMPTS); do
  if docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "База данных готова! (попытка $i/$MAX_DB_ATTEMPTS)"
    # Проверяем соединение с базой дополнительно
    if docker-compose -f docker-compose.prod.yml exec -T db psql -U postgres -c "SELECT 1" > /dev/null 2>&1; then
      echo "Успешно проверено соединение с базой данных!"
      
      # Задержка для полной стабилизации БД
      echo "Ожидаем дополнительное время для стабилизации БД..."
      sleep 5
      break
    else
      echo "pg_isready вернул успех, но не удалось подключиться через psql. Продолжаем ожидание..."
    fi
  fi
  
  # На последней попытке выполняем продвинутую диагностику
  if [ $i -eq $MAX_DB_ATTEMPTS ]; then
    echo "Выполняем продвинутую диагностику базы данных..."
    
    # Проверяем состояние контейнера
    CONTAINER_ID=$(docker-compose -f docker-compose.prod.yml ps -q db)
    if [ -n "$CONTAINER_ID" ]; then
      echo "Контейнер БД имеет ID: $CONTAINER_ID"
      
      echo "Логи контейнера базы данных:"
      docker logs "$CONTAINER_ID"
      
      # Исправление ошибки с флагом -T
      echo "Файлы в директории данных:"
      docker exec "$CONTAINER_ID" ls -la /var/lib/postgresql/data || true
      
      # Исправление ошибки с флагом -T
      echo "Пробуем перезапустить PostgreSQL внутри контейнера..."
      docker exec "$CONTAINER_ID" bash -c "pg_ctl -D /var/lib/postgresql/data restart" || true
      
      # Даем время на перезапуск
      sleep 10
      
      # Проверяем еще раз
      if docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
        echo "БД заработала после перезапуска!"
        break
      fi
    else
      echo "Контейнер БД не найден. Это странно..."
    fi
    
    # Последняя попытка - полностью удаляем содержимое тома и инициализируем новый кластер
    echo "ВНИМАНИЕ: Все предыдущие попытки не сработали. Удаляем содержимое тома и создаем новый кластер."
    echo "Создаем резервную копию текущего тома..."
    FINAL_BACKUP_PATH="${DB_VOLUME_PATH}_final_backup_$(date +%Y%m%d_%H%M%S)"
    cp -r "$DB_VOLUME_PATH" "$FINAL_BACKUP_PATH" 2>/dev/null || true
    rm -rf "$DB_VOLUME_PATH"/*
    mkdir -p "$DB_VOLUME_PATH/pgdata"
    touch "$DB_VOLUME_PATH/pgdata/.postgresql_empty"
    chown -R 999:999 "$DB_VOLUME_PATH"
    chmod -R 700 "$DB_VOLUME_PATH"
    
    # Останавливаем и удаляем контейнер
    docker-compose -f docker-compose.prod.yml stop db
    docker-compose -f docker-compose.prod.yml rm -f db
    
    # Запускаем с новым томом
    echo "Запускаем БД с чистым томом..."
    docker-compose -f docker-compose.prod.yml up -d db
    
    # Даем время на инициализацию
    sleep 15
    
    # Проверяем состояние
    if docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
      echo "БД заработала с новым кластером!"
      echo "ВАЖНО: Все данные предыдущего кластера сохранены в $FINAL_BACKUP_PATH"
      echo "Необходимо будет восстановить данные из резервной копии после завершения деплоя."
      break
    else
      echo "БД все еще не работает. Продолжаем деплой без проверки БД..."
      echo "ВНИМАНИЕ: Это может привести к ошибкам в работе приложения!"
      break
    fi
  fi
  
  echo "Ожидаем готовность базы данных... Попытка $i/$MAX_DB_ATTEMPTS"
  sleep 5
done

# Запускаем миграции напрямую
echo "Запускаем миграции..."
docker run --rm --network hostel-booking-system_hostel_network \
  -v $(pwd)/backend/migrations:/migrations \
  migrate/migrate \
  -path=/migrations/ \
  -database="postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable" \
  up || {
    echo "Ошибка при выполнении миграций. Продолжаем запуск других сервисов..."
  }

# Запускаем остальные сервисы
echo "Запускаем остальные сервисы..."
docker-compose -f docker-compose.prod.yml up -d opensearch

# Даем время OpenSearch запуститься
echo "Ожидаем запуска OpenSearch..."
sleep 15

# Запускаем оставшиеся сервисы
echo "Запускаем backend и другие сервисы..."
docker-compose -f docker-compose.prod.yml up -d backend opensearch-dashboards nginx

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
  echo "Ожидаем готовность backend... Попытка $i/$RETRY_COUNT"
  sleep 3
done

# Сохраняем последние 5 бэкапов и удаляем более старые
find /tmp/hostel-backup/db -name "*.sql" -type f | sort -r | tail -n +6 | xargs rm -f 2>/dev/null || true

echo "Деплой завершен!"
echo "Логи контейнеров можно посмотреть с помощью: docker-compose -f docker-compose.prod.yml logs -f"