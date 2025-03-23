#!/bin/bash
set -e # Останавливаем выполнение при ошибках
echo "Начинаем деплой..."
cd /opt/hostel-booking-system

# Создаем директории для хранения данных, если их еще нет
mkdir -p /opt/hostel-data/uploads
mkdir -p /opt/hostel-data/db
mkdir -p /opt/hostel-data/opensearch
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

# Удаляем старые образы (но не удаляем контейнеры и данные)
docker image prune -f

# Собираем фронтенд
echo "Собираем фронтенд..."
cd frontend/hostel-frontend

# Создаем измененный package.json с дополнительными зависимостями
echo "Добавляем необходимые зависимости в package.json..."
if [ -f "package.json" ]; then
  # Сохраняем оригинальный package.json
  cp package.json package.json.orig
  
  # Добавляем react-query в зависимости
  sed -i 's/"dependencies": {/"dependencies": {\n    "react-query": "^3.39.3",\n    "react-window": "^1.8.9",/g' package.json
fi

# Создаем кастомный скрипт для сборки фронтенда
cat > build_frontend.sh << 'EOL'
#!/bin/sh
set -e
echo "Устанавливаем зависимости..."
npm cache clean --force
npm install --legacy-peer-deps

# Проверяем, установился ли react-query
if ! npm list react-query >/dev/null 2>&1; then
  echo "Принудительная установка react-query..."
  npm install react-query@3.39.3 --save --legacy-peer-deps
fi

# Проверяем, установился ли react-window
if ! npm list react-window >/dev/null 2>&1; then
  echo "Принудительная установка react-window..."
  npm install react-window@1.8.9 --save --legacy-peer-deps
fi

# Устанавливаем react-scripts
echo "Устанавливаем React Scripts и другие основные зависимости..."
npm install react-scripts@5.0.1 --save --legacy-peer-deps
npm install ajv@6.12.6 ajv-keywords@3.5.2 schema-utils@3.1.1 --legacy-peer-deps

echo "Пробуем сборку проекта..."
npm run build || {
  echo "Сборка не удалась, пробуем установить дополнительные зависимости..."

  # Анализируем ошибки сборки
  npm run build 2>&1 | tee build_error.log
  
  # Находим все упоминания о недостающих пакетах
  if grep -q "Can't resolve" build_error.log; then
    packages=$(grep -o "Can't resolve '[^']*'" build_error.log | sed "s/Can't resolve '//g" | sed "s/'//g")
    echo "Найдены отсутствующие пакеты: $packages"
    
    for pkg in $packages; do
      echo "Устанавливаем $pkg..."
      npm install "$pkg" --save --legacy-peer-deps
    done
    
    # Пробуем собрать снова
    echo "Пробуем собрать проект снова..."
    npm run build || {
      echo "Сборка не удалась повторно. Пробуем последний вариант..."
      
      # Особые случаи для конкретных пакетов
      if grep -q "Can't resolve 'react-query'" build_error.log; then
        echo "Особый случай для react-query..."
        npm uninstall react-query
        npm install react-query@3.39.0 --legacy-peer-deps
        npm install @tanstack/react-query --legacy-peer-deps
      fi
      
      # Последняя попытка
      npm run build
    }
  fi
}
EOL

chmod +x build_frontend.sh
NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "./build_frontend.sh"

# Проверяем, успешно ли прошла сборка
if [ ! -d "build" ] || [ -z "$(ls -A build)" ]; then
  echo "Сборка фронтенда не удалась. Деплой прерван."
  exit 1
fi

cd ../..

# Останавливаем только контейнеры, не удаляя тома
echo "Останавливаем сервисы (без удаления томов)..."
docker-compose -f docker-compose.prod.yml down

# Запускаем с обновленными образами
echo "Запускаем сервисы с обновленными образами..."
docker-compose -f docker-compose.prod.yml up --build -d

# Проверяем базу данных
echo "Проверяем готовность базы данных..."
RETRY_COUNT=30
for i in $(seq 1 $RETRY_COUNT); do
  if docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "База данных готова!"
    break
  fi
  echo "Ожидаем готовность базы данных... Попытка $i/$RETRY_COUNT"
  sleep 2
done

if ! docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
  echo "Не удалось запустить базу данных"
  exit 1
fi

# Запускаем миграции
echo "Запускаем миграции..."
docker run --rm --network hostel-booking-system_hostel_network -v $(pwd)/backend/migrations:/migrations migrate/migrate -path=/migrations/ -database="postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable" up

# Сохраняем последние 5 бэкапов и удаляем более старые
find /tmp/hostel-backup/db -name "*.sql" -type f | sort -r | tail -n +6 | xargs rm -f 2>/dev/null || true

echo "Деплой завершен!"