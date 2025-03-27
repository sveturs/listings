#!/bin/bash
set -e # Останавливаем выполнение при ошибках

# Функция логирования с отметками времени
log() {
    echo "[$(date +"%Y-%m-%d %H:%M:%S")] $1"
}

# Настройка логирования
LOG_FILE="/var/log/deploy.log"
# Проверка наличия и доступности директории для логов
if [ ! -d "/var/log" ] || [ ! -w "/var/log" ]; then
    LOG_FILE="/tmp/deploy.log"
    log "ВНИМАНИЕ: Невозможно записать в /var/log, используем $LOG_FILE"
fi

# Проверка наличия лог-файла
if [ ! -f "$LOG_FILE" ]; then
    touch "$LOG_FILE" 2>/dev/null || true
fi

# Проверка прав на запись
if [ ! -w "$LOG_FILE" ] && [ -f "$LOG_FILE" ]; then
    chmod 666 "$LOG_FILE" 2>/dev/null || true
    if [ ! -w "$LOG_FILE" ]; then
        LOG_FILE="/tmp/deploy.log"
        touch "$LOG_FILE" 2>/dev/null || true
        log "ВНИМАНИЕ: Проблемы с правами доступа, используем $LOG_FILE"
    fi
fi

# Ротация логов деплоя
if [ -f "$LOG_FILE" ] && [ -w "$LOG_FILE" ]; then
    log "Ротация логов деплоя..."
    if [ -w "$(dirname "$LOG_FILE")" ]; then
        cp -f "$LOG_FILE" "${LOG_FILE}.$(date +%Y%m%d)" 2>/dev/null || true
        find "$(dirname "$LOG_FILE")" -name "$(basename "$LOG_FILE").*" -type f -mtime +7 -delete 2>/dev/null || true
        : > "$LOG_FILE" # Очищаем текущий лог-файл
    else
        log "ВНИМАНИЕ: Нет прав на ротацию логов в директории $(dirname "$LOG_FILE")"
    fi
fi

log "Начинаем деплой..."

# Функция для принудительной очистки Docker
docker_force_cleanup() {
    log "=== Выполняем принудительную очистку Docker ==="
    
    # Останавливаем все запущенные контейнеры, связанные с проектом
    log "Останавливаем контейнеры, связанные с проектом..."
    docker ps -a | grep -E 'hostel_|hostel-booking|opensearch' | awk '{print $1}' | xargs -r docker stop
    
    # Удаляем все остановленные контейнеры, связанные с проектом
    log "Удаляем контейнеры, связанные с проектом..."
    docker ps -a | grep -E 'hostel_|hostel-booking|opensearch' | awk '{print $1}' | xargs -r docker rm -f
    
    # Обработка проблемной сети
    local network_id=$(docker network ls | grep "hostel-booking-system_hostel_network" | awk '{print $1}')
    if [ -n "$network_id" ]; then
        log "Найдена сеть с ID: $network_id, пробуем удалить..."
        
        # Отключение всех контейнеров от сети
        local endpoints=$(docker network inspect $network_id -f '{{range $k, $v := .Containers}}{{$k}} {{end}}' 2>/dev/null || echo "")
        for ep in $endpoints; do
            log "Отключение контейнера $ep от сети..."
            docker network disconnect -f $network_id $ep 2>/dev/null || true
        done
        
        # Удаление сети
        docker network rm $network_id 2>/dev/null || true
        
        # Если сеть не удалось удалить, использовать новую сеть с временной меткой
        if docker network ls | grep -q $network_id; then
            log "ВНИМАНИЕ: Не удалось удалить сеть. Будет создана новая сеть с другим именем."
            NETWORK_SUFFIX=$(date +%s)
            NETWORK_NAME="hostel-booking-system_hostel_network_${NETWORK_SUFFIX}"
            log "Создаем новую сеть с именем $NETWORK_NAME"
            docker network create "$NETWORK_NAME"
            
            # Обновляем docker-compose.prod.yml для использования новой сети
            sed -i "s/hostel_network:$/hostel_network_${NETWORK_SUFFIX}:/g" docker-compose.prod.yml
            sed -i "s/hostel_network$/hostel_network_${NETWORK_SUFFIX}/g" docker-compose.prod.yml
        fi
    fi
    
    # Очистка неиспользуемых ресурсов Docker
    log "Очистка неиспользуемых ресурсов Docker..."
    docker system prune -f
    
    log "=== Очистка Docker завершена ==="
}

# Функция для обработки прерывания выполнения скрипта
cleanup_on_interrupt() {
    log "Получен сигнал прерывания. Выполняем безопасную остановку..."
    # Остановка всех контейнеров из docker-compose
    docker-compose -f docker-compose.prod.yml down || true
    log "Деплой прерван пользователем."
    exit 1
}

# Устанавливаем обработчик сигналов прерывания
trap cleanup_on_interrupt SIGINT SIGTERM


log "Начинаем деплой..."
cd /opt/hostel-booking-system

# Создаем директории для хранения данных, если их еще нет
mkdir -p /opt/hostel-data/uploads
mkdir -p /opt/hostel-data/db
mkdir -p /opt/hostel-data/opensearch
mkdir -p /opt/hostel-data/yarn-cache # Директория для кэша yarn
mkdir -p certbot/conf
mkdir -p certbot/www
mkdir -p /tmp/hostel-backup/db

# Ротация логов деплоя
if [ -f "/var/log/deploy.log" ]; then
    log "Ротация логов деплоя..."
    mv /var/log/deploy.log /var/log/deploy.log.$(date +%Y%m%d)
    find /var/log -name "deploy.log.*" -type f -mtime +7 -delete
fi

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
log "Пытаемся создать бэкап базы данных..."
if docker-compose -f docker-compose.prod.yml ps | grep -q "db.*Up"; then
    BACKUP_FILE="/tmp/hostel-backup/db/backup_$(date +%Y%m%d_%H%M%S).sql"
    docker-compose -f docker-compose.prod.yml exec -T db pg_dumpall -U postgres > "$BACKUP_FILE"
    if [ $? -eq 0 ]; then
        log "Бэкап базы данных создан в $BACKUP_FILE"
    else
        log "Ошибка создания бэкапа базы данных, но продолжаем деплой"
    fi
else
    log "База данных не запущена, пропускаем создание бэкапа"
fi

# Перед сбросом изменений выполняем первичную очистку Docker
log "Выполняем первичную очистку Docker перед обновлением кода..."
docker_force_cleanup

# Обеспечиваем чистое состояние git, но исключаем критические директории
log "Получаем последние изменения из репозитория..."
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
log "Останавливаем сервисы..."
docker-compose -f docker-compose.prod.yml down --remove-orphans

# Принудительно удаляем все связанные контейнеры, но оставляем тома нетронутыми
log "Удаляем старые контейнеры, сохраняя тома с данными..."
docker-compose -f docker-compose.prod.yml rm -f

# Вторая очистка Docker после остановки контейнеров
log "Выполняем вторую очистку Docker после остановки контейнеров..."
docker_force_cleanup

# Удаляем старые образы, чтобы принудительно пересобрать их
log "Пересобираем образы..."
docker-compose -f docker-compose.prod.yml build --no-cache

# Собираем фронтенд
log "Собираем фронтенд..."
cd frontend/hostel-frontend

# Проверяем существование package.json и создаем измененный файл
log "Проверяем package.json и добавляем необходимые зависимости..."
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
log "Запускаем сборку фронтенда в контейнере с увеличенными ресурсами..."
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
    log "Сборка фронтенда не удалась. Проверяем логи для анализа ошибок..."
    
    # Анализируем логи сборки
    if [ -f "build_error.log" ]; then
        ERROR_LOG=$(cat build_error.log)
        if echo "$ERROR_LOG" | grep -q "out of memory"; then
            log "Ошибка сборки: Недостаточно памяти. Попробуйте увеличить значение NODE_OPTIONS."
        elif echo "$ERROR_LOG" | grep -q "Cannot find module"; then
            log "Ошибка сборки: Не найден модуль. Проверьте зависимости."
        else
            log "Ошибка сборки: Неизвестная проблема. Проверьте полные логи."
        fi
    fi
    
    # Использовать предыдущую сборку, если она есть
    if [ -d "build" ] && [ -n "$(ls -A build 2>/dev/null)" ]; then
        log "Используем существующую сборку фронтенда..."
    else
        log "Деплой прерван из-за проблем сборки фронтенда и отсутствия существующей сборки."
        exit 1
    fi
fi

cd ../..

# Третья очистка Docker перед запуском сервисов
log "Выполняем очистку Docker перед запуском сервисов..."
docker_force_cleanup

# Подготовка директорий для базы данных
DB_VOLUME_PATH="/opt/hostel-data/db"

# Если база была запущена ранее, сначала создаем резервную копию данных
if [ -d "$DB_VOLUME_PATH" ] && [ -n "$(ls -A $DB_VOLUME_PATH 2>/dev/null)" ]; then
    BACKUP_DIR="${DB_VOLUME_PATH}_backup_$(date +%Y%m%d_%H%M%S)"
    log "Сохраняем текущее состояние БД в $BACKUP_DIR..."
    mkdir -p "$BACKUP_DIR"
    cp -r "$DB_VOLUME_PATH/"* "$BACKUP_DIR/" 2>/dev/null || true
fi

# Создаем новую структуру директорий для PostgreSQL
log "Настраиваем структуру каталогов для PostgreSQL..."
mkdir -p "$DB_VOLUME_PATH/data"
chown -R 999:999 "$DB_VOLUME_PATH"
chmod -R 700 "$DB_VOLUME_PATH"

# Модифицируем docker-compose.prod.yml для использования новой структуры
log "Обновляем конфигурацию docker-compose.prod.yml..."
sed -i 's|device: /opt/hostel-data/db|device: /opt/hostel-data/db/data|g' docker-compose.prod.yml
if ! grep -q "PGDATA:" docker-compose.prod.yml; then
    sed -i '/POSTGRES_DB:/a\      PGDATA: /var/lib/postgresql/data/pgdata  # Подкаталог для хранения файлов БД' docker-compose.prod.yml
fi

# Запускаем сервисы с принудительным пересозданием контейнеров
log "Запускаем все сервисы..."
docker-compose -f docker-compose.prod.yml up -d --force-recreate --remove-orphans

# Выжидаем небольшую паузу для инициализации сервисов
sleep 10

# Проверяем, все ли сервисы успешно запустились
log "Проверяем статус запущенных сервисов..."
RUNNING_SERVICES=$(docker-compose -f docker-compose.prod.yml ps --services --filter "status=running")
EXPECTED_SERVICES="db opensearch opensearch-dashboards backend nginx"

# Подробная проверка статуса каждого сервиса
ALL_SERVICES_RUNNING=true
for SERVICE in $EXPECTED_SERVICES; do
    if ! echo "$RUNNING_SERVICES" | grep -q "$SERVICE"; then
        ALL_SERVICES_RUNNING=false
        log "ВНИМАНИЕ: Сервис $SERVICE не запущен!"
        # Проверяем логи неработающего сервиса
        log "Последние логи сервиса $SERVICE:"
        docker-compose -f docker-compose.prod.yml logs --tail=15 $SERVICE
    fi
done

# Проверяем здоровье backend с несколькими попытками
log "Проверяем здоровье backend..."
RETRY_COUNT=30
BACKEND_HEALTHY=false

for i in $(seq 1 $RETRY_COUNT); do
    if docker-compose -f docker-compose.prod.yml exec -T backend curl -s -f http://localhost:3000 > /dev/null 2>&1; then
        log "Backend здоров! (попытка $i/$RETRY_COUNT)"
        BACKEND_HEALTHY=true
        break
    fi
    
    # На определенных попытках проверяем логи
    if [ $i -eq 5 ] || [ $i -eq 15 ] || [ $i -eq 25 ]; then
        log "Backend все еще не отвечает (попытка $i/$RETRY_COUNT). Проверяем логи..."
        docker-compose -f docker-compose.prod.yml logs --tail=20 backend
    fi
    
    log "Ожидаем готовность backend... Попытка $i/$RETRY_COUNT"
    sleep 3
done

# Проверяем, удалось ли установить соединение с backend
if ! $BACKEND_HEALTHY; then
    log "ВНИМАНИЕ: Не удалось установить соединение с backend после $RETRY_COUNT попыток."
    log "Последние логи backend:"
    docker-compose -f docker-compose.prod.yml logs --tail=50 backend
fi

# Итоговая проверка состояния системы
if $ALL_SERVICES_RUNNING && $BACKEND_HEALTHY; then
    log "Все сервисы успешно запущены и работают!"
    log "Деплой завершен успешно!"
else
    log "ВНИМАНИЕ: Некоторые сервисы не запущены или не отвечают."
    log "Текущий статус контейнеров:"
    docker-compose -f docker-compose.prod.yml ps
    
    # Последняя попытка перезапуска проблемных сервисов
    if ! $ALL_SERVICES_RUNNING; then
        log "Пробуем перезапустить проблемные сервисы..."
        docker-compose -f docker-compose.prod.yml up -d
    fi
fi

# Сохраняем последние 5 бэкапов и удаляем более старые
find /tmp/hostel-backup/db -name "*.sql" -type f | sort -r | tail -n +6 | xargs rm -f 2>/dev/null || true

log "Деплой завершен!"
log "Логи контейнеров можно посмотреть с помощью: docker-compose -f docker-compose.prod.yml logs -f"