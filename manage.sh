#!/bin/bash
set -e # Останавливаем выполнение при ошибках

PROJECT_DIR="/opt/hostel-booking-system"
DOCKER_COMPOSE_FILE="$PROJECT_DIR/docker-compose.prod.yml"
BACKUP_DIR="/tmp/hostel-backup/db"

# Функция для создания резервной копии базы данных
backup_database() {
  echo "Создаем резервную копию базы данных..."
  
  # Проверяем, запущен ли контейнер с базой данных
  if docker-compose -f "$DOCKER_COMPOSE_FILE" ps | grep -q "db.*Up"; then
    # Создаем директорию для бэкапа, если её нет
    mkdir -p "$BACKUP_DIR"
    
    # Создаем файл бэкапа с временной меткой
    BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"
    
    # Выполняем дамп базы данных
    docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_dumpall -U postgres > "$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
      echo "✅ Резервная копия базы данных успешно создана в $BACKUP_FILE"
      
      # Оставляем только последние 5 бэкапов
      find "$BACKUP_DIR" -name "*.sql" -type f | sort -r | tail -n +6 | xargs rm -f 2>/dev/null || true
    else
      echo "❌ Ошибка создания резервной копии базы данных"
      return 1
    fi
  else
    echo "❌ База данных не запущена, невозможно создать резервную копию"
    return 1
  fi
}

# Функция для остановки всего проекта
stop_all() {
  echo "🛑 Останавливаем все сервисы проекта..."
  cd "$PROJECT_DIR"
  docker-compose -f "$DOCKER_COMPOSE_FILE" down --remove-orphans
  echo "✅ Все сервисы проекта остановлены"
}

# Функция для запуска всего проекта
start_all() {
  echo "🚀 Запускаем все сервисы проекта..."
  cd "$PROJECT_DIR"
  docker-compose -f "$DOCKER_COMPOSE_FILE" up -d
  
  # Проверяем готовность базы данных
  echo "⏳ Проверяем готовность базы данных..."
  RETRY_COUNT=30
  for i in $(seq 1 $RETRY_COUNT); do
    if docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
      echo "✅ База данных готова!"
      break
    fi
    echo "⏳ Ожидаем готовность базы данных... Попытка $i/$RETRY_COUNT"
    sleep 2
  done
  
  if ! docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "❌ Не удалось запустить базу данных"
    return 1
  fi
  
  echo "✅ Все сервисы проекта запущены"
}

# Функция для перезапуска всего проекта
restart_all() {
  echo "🔄 Перезапускаем все сервисы проекта..."
  stop_all
  sleep 2
  start_all
  echo "✅ Все сервисы проекта перезапущены"
}

# Функция для остановки отдельного сервиса
stop_service() {
  SERVICE_NAME=$1
  
  if [ -z "$SERVICE_NAME" ]; then
    echo "❌ Необходимо указать имя сервиса"
    show_help
    return 1
  fi
  
  echo "🛑 Останавливаем сервис $SERVICE_NAME..."
  cd "$PROJECT_DIR"
  docker-compose -f "$DOCKER_COMPOSE_FILE" stop "$SERVICE_NAME"
  echo "✅ Сервис $SERVICE_NAME остановлен"
}

# Функция для запуска отдельного сервиса
start_service() {
  SERVICE_NAME=$1
  
  if [ -z "$SERVICE_NAME" ]; then
    echo "❌ Необходимо указать имя сервиса"
    show_help
    return 1
  fi
  
  echo "🚀 Запускаем сервис $SERVICE_NAME..."
  cd "$PROJECT_DIR"
  docker-compose -f "$DOCKER_COMPOSE_FILE" up -d "$SERVICE_NAME"
  
  if [ "$SERVICE_NAME" = "db" ]; then
    # Проверяем готовность базы данных
    echo "⏳ Проверяем готовность базы данных..."
    RETRY_COUNT=30
    for i in $(seq 1 $RETRY_COUNT); do
      if docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
        echo "✅ База данных готова!"
        break
      fi
      echo "⏳ Ожидаем готовность базы данных... Попытка $i/$RETRY_COUNT"
      sleep 2
    done
    
    if ! docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
      echo "❌ Не удалось запустить базу данных"
      return 1
    fi
  fi
  
  echo "✅ Сервис $SERVICE_NAME запущен"
}

# Функция для перезапуска отдельного сервиса
restart_service() {
  SERVICE_NAME=$1
  
  if [ -z "$SERVICE_NAME" ]; then
    echo "❌ Необходимо указать имя сервиса"
    show_help
    return 1
  fi
  
  echo "🔄 Перезапускаем сервис $SERVICE_NAME..."
  cd "$PROJECT_DIR"
  docker-compose -f "$DOCKER_COMPOSE_FILE" restart "$SERVICE_NAME"
  
  if [ "$SERVICE_NAME" = "db" ]; then
    # Проверяем готовность базы данных
    echo "⏳ Проверяем готовность базы данных..."
    RETRY_COUNT=30
    for i in $(seq 1 $RETRY_COUNT); do
      if docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
        echo "✅ База данных готова!"
        break
      fi
      echo "⏳ Ожидаем готовность базы данных... Попытка $i/$RETRY_COUNT"
      sleep 2
    done
    
    if ! docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
      echo "❌ Не удалось запустить базу данных"
      return 1
    fi
  fi
  
  echo "✅ Сервис $SERVICE_NAME перезапущен"
}

# Функция для получения статуса всех сервисов
status() {
  echo "📊 Статус сервисов проекта:"
  cd "$PROJECT_DIR"
  docker-compose -f "$DOCKER_COMPOSE_FILE" ps
}

# Функция для выполнения миграций
run_migrations() {
  echo "🔄 Запускаем миграции базы данных..."
  cd "$PROJECT_DIR"
  
  # Проверяем, запущена ли база данных
  if ! docker-compose -f "$DOCKER_COMPOSE_FILE" exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "❌ База данных не запущена, запускаем..."
    start_service "db"
  fi
  
  # Запускаем миграции
  docker run --rm --network hostel-booking-system_hostel_network \
    -v "$(pwd)/backend/migrations:/migrations" \
    migrate/migrate \
    -path=/migrations/ \
    -database="postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable" \
    up
  
  echo "✅ Миграции успешно выполнены"
}

# Функция для отображения справки
show_help() {
  echo "Управление хостельной системой бронирования"
  echo ""
  echo "Использование: $0 <команда> [аргументы]"
  echo ""
  echo "Команды:"
  echo "  start-all            - Запустить все сервисы"
  echo "  stop-all             - Остановить все сервисы"
  echo "  restart-all          - Перезапустить все сервисы"
  echo "  start <service>      - Запустить конкретный сервис"
  echo "  stop <service>       - Остановить конкретный сервис"
  echo "  restart <service>    - Перезапустить конкретный сервис"
  echo "  status               - Показать статус всех сервисов"
  echo "  backup               - Создать резервную копию базы данных"
  echo "  migrations           - Запустить миграции базы данных"
  echo "  help                 - Показать эту справку"
  echo ""
  echo "Доступные сервисы в проекте (примерный список):"
  echo "  db                   - База данных PostgreSQL"
  echo "  opensearch           - Поисковый движок OpenSearch"
  echo "  backend              - Бэкенд-сервис"
  echo "  frontend             - Фронтенд-сервис"
  echo "  nginx                - Веб-сервер и прокси"
  echo "  certbot              - Сервис для управления SSL-сертификатами"
  echo ""
  echo "Для получения полного списка сервисов используйте команду status"
}

# Проверяем наличие директории проекта
if [ ! -d "$PROJECT_DIR" ]; then
  echo "❌ Директория проекта $PROJECT_DIR не существует"
  exit 1
fi

# Проверяем наличие файла docker-compose
if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
  echo "❌ Файл docker-compose $DOCKER_COMPOSE_FILE не существует"
  exit 1
fi

# Обрабатываем команды
case "$1" in
  start-all)
    start_all
    ;;
  stop-all)
    stop_all
    ;;
  restart-all)
    restart_all
    ;;
  start)
    start_service "$2"
    ;;
  stop)
    stop_service "$2"
    ;;
  restart)
    restart_service "$2"
    ;;
  status)
    status
    ;;
  backup)
    backup_database
    ;;
  migrations)
    run_migrations
    ;;
  help|--help|-h)
    show_help
    ;;
  *)
    echo "❌ Неизвестная команда: $1"
    show_help
    exit 1
    ;;
esac