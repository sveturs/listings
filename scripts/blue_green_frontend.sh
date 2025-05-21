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