#!/bin/bash
set -e # Останавливаем выполнение при ошибках
echo "Начинаем деплой..."
cd /opt/hostel-booking-system

# Настраиваем git pull strategy
git config pull.rebase false

# Создаем необходимые директории
mkdir -p backend/uploads
mkdir -p frontend/hostel-frontend/build
mkdir -p certbot/conf
mkdir -p certbot/www
mkdir -p /tmp/hostel-backup/db

# Сохраняем важные файлы
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true
cp -f frontend/hostel-frontend/.env /tmp/hostel-backup/frontend.env 2>/dev/null || true

# Сохраняем SSL сертификаты
if [ -d "certbot/conf" ]; then
  cp -r certbot/conf /tmp/hostel-backup/ 2>/dev/null || true
fi

# Сохраняем загруженные изображения
cp -r backend/uploads /tmp/hostel-backup/ 2>/dev/null || true

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

# Обеспечиваем чистое состояние git
git fetch origin
git reset --hard origin/main
git clean -fdx -e "*.env*" -e "uploads/" -e "certbot/"

# Восстанавливаем файлы
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true
if [ -d "/tmp/hostel-backup/conf" ]; then
  rm -rf certbot/conf
  cp -r /tmp/hostel-backup/conf certbot/ 2>/dev/null || true
fi

# Удаляем старые образы
docker image prune -f

# Очищаем сети и осиротевшие контейнеры
echo "Очищаем старые контейнеры и сети..."
docker-compose -f docker-compose.prod.yml down -v --remove-orphans || true
docker network prune -f || true

# Собираем фронтенд
echo "Собираем фронтенд..."
cd frontend/hostel-frontend

# Создаем кастомный скрипт для сборки фронтенда с установкой дополнительных пакетов
cat > build_frontend.sh << 'EOL'
#!/bin/sh
set -e
echo "Устанавливаем зависимости..."
npm cache clean --force
npm install --legacy-peer-deps

# Устанавливаем react-scripts
echo "Устанавливаем React Scripts и другие основные зависимости..."
npm install react-scripts@5.0.1 --save --legacy-peer-deps
npm install ajv@6.12.6 ajv-keywords@3.5.2 schema-utils@3.1.1 --legacy-peer-deps

# Устанавливаем пакеты, которых может не хватать
echo "Устанавливаем дополнительные зависимости, которые могут потребоваться..."
npm install react-window --legacy-peer-deps
npm install react-virtualized --legacy-peer-deps
npm install @tanstack/react-query --legacy-peer-deps
npm install @tanstack/react-table --legacy-peer-deps
npm install @mui/x-date-pickers --legacy-peer-deps
npm install @faker-js/faker --legacy-peer-deps
npm install framer-motion --legacy-peer-deps

# Проверяем и устанавливаем пакеты из сообщений об ошибках
if grep -q "Can't resolve" npm-debug.log 2>/dev/null; then
  echo "Устанавливаем пакеты из сообщений об ошибках..."
  for pkg in $(grep -o "'[^']*'" npm-debug.log | grep -v node_modules | sed "s/'//g"); do
    echo "Установка $pkg..."
    npm install $pkg --legacy-peer-deps || true
  done
fi

echo "Запускаем сборку..."
npm run build
EOL

chmod +x build_frontend.sh
NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "./build_frontend.sh"

# Проверяем, успешно ли прошла сборка
if [ ! -d "build" ] || [ -z "$(ls -A build)" ]; then
  echo "Сборка фронтенда не удалась. Проверьте логи выше на наличие ошибок."
  
  # Если сборка не удалась, попробуем установить все недостающие пакеты
  echo "Пытаемся исправить ошибки сборки..."
  
  cat > fix_build.sh << 'EOL'
#!/bin/sh
set -e

# Сначала сохраним лог с ошибками
npm run build > build_errors.log 2>&1 || true

# Находим все сообщения об отсутствующих пакетах
if grep -q "Can't resolve" build_errors.log; then
  echo "Найдены отсутствующие пакеты, устанавливаем их..."
  for pkg in $(grep -o "Can't resolve '[^']*'" build_errors.log | sed "s/Can't resolve '//g" | sed "s/'//g"); do
    echo "Устанавливаем $pkg..."
    npm install $pkg --legacy-peer-deps || true
  done
  
  # Пробуем собрать снова
  echo "Пробуем собрать снова..."
  npm run build
fi
EOL
  
  chmod +x fix_build.sh
  NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "./fix_build.sh"
  
  # Проверяем еще раз
  if [ ! -d "build" ] || [ -z "$(ls -A build)" ]; then
    echo "Не удалось исправить ошибки сборки. Деплой прерван."
    exit 1
  fi
fi

cd ../..

# Запускаем только базу данных
echo "Запускаем базу данных..."
docker-compose -f docker-compose.prod.yml up --build -d db

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

# Восстанавливаем данные из бэкапа, если он есть
if [ -n "$(ls -t /tmp/hostel-backup/db/*.sql 2>/dev/null | head -1)" ]; then
  LATEST_BACKUP=$(ls -t /tmp/hostel-backup/db/*.sql | head -1)
  echo "Восстанавливаем базу данных из $LATEST_BACKUP..."
  cat "$LATEST_BACKUP" | docker-compose -f docker-compose.prod.yml exec -T db psql -U postgres
  if [ $? -eq 0 ]; then
    echo "База данных успешно восстановлена"
  else
    echo "Ошибка восстановления базы данных, но продолжаем деплой"
  fi
else
  echo "Бэкап базы данных не найден, пропускаем восстановление"
fi

# Запускаем остальные сервисы
echo "Запускаем остальные сервисы..."
docker-compose -f docker-compose.prod.yml up --build -d

# Сохраняем последние 5 бэкапов и удаляем более старые
find /tmp/hostel-backup/db -name "*.sql" -type f | sort -r | tail -n +6 | xargs rm -f 2>/dev/null || true

echo "Деплой завершен!"