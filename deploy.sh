#!/bin/bash
set -e # Stop execution on errors
echo "Starting deployment..."
cd /opt/hostel-booking-system

# Setup git pull strategy
git config pull.rebase false

# Create necessary directories
mkdir -p backend/uploads
mkdir -p frontend/hostel-frontend/build
mkdir -p certbot/conf
mkdir -p certbot/www
mkdir -p /tmp/hostel-backup/db

# Save important files
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true
cp -f frontend/hostel-frontend/.env /tmp/hostel-backup/frontend.env 2>/dev/null || true

# Save SSL certificates
if [ -d "certbot/conf" ]; then
  cp -r certbot/conf /tmp/hostel-backup/ 2>/dev/null || true
fi

# Save uploaded images
cp -r backend/uploads /tmp/hostel-backup/ 2>/dev/null || true

# Backup database
echo "Creating database backup..."
BACKUP_FILE="/tmp/hostel-backup/db/backup_$(date +%Y%m%d_%H%M%S).sql"
docker-compose -f docker-compose.prod.yml exec -T db pg_dumpall -U postgres > "$BACKUP_FILE"
if [ $? -ne 0 ]; then
    echo "Database backup failed!"
    exit 1
fi
echo "Database backup created at $BACKUP_FILE"

# Clean git state
git fetch origin
git reset --hard origin/main
git clean -fdx -e "*.env*" -e "uploads/" -e "certbot/"

# Restore files
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true
if [ -d "/tmp/hostel-backup/conf" ]; then
  rm -rf certbot/conf
  cp -r /tmp/hostel-backup/conf certbot/ 2>/dev/null || true
fi

# Clean old images
docker image prune -f

# Clean networks and orphaned containers
echo "Cleaning up orphan containers and networks..."
docker-compose -f docker-compose.prod.yml down -v --remove-orphans || true
docker network prune -f || true

# Build frontend
echo "Building frontend..."
cd frontend/hostel-frontend
NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "\
  npm cache clean --force && \
  npm install --legacy-peer-deps && \
  npm install react-scripts@5.0.1 --save --legacy-peer-deps && \
  npm install ajv@6.12.6 ajv-keywords@3.5.2 schema-utils@3.1.1 --legacy-peer-deps && \
  npm run build"
cd ../..

# Start database
echo "Starting database..."
docker-compose -f docker-compose.prod.yml up --build -d db

# Check database
echo "Checking database readiness..."
RETRY_COUNT=30
for i in $(seq 1 $RETRY_COUNT); do
  if docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "Database is ready!"
    break
  fi
  echo "Waiting for database to be ready... Attempt $i/$RETRY_COUNT"
  sleep 2
done

if ! docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
  echo "Database failed to start"
  exit 1
fi

# Run migrations
echo "Running migrations..."
docker run --rm --network hostel-booking-system_hostel_network -v $(pwd)/backend/migrations:/migrations migrate/migrate -path=/migrations/ -database="postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable" up

# Restore database data
echo "Restoring database data..."
cat "$BACKUP_FILE" | docker-compose -f docker-compose.prod.yml exec -T db psql -U postgres
if [ $? -ne 0 ]; then
    echo "Database restore failed!"
    exit 1
fi
echo "Database restored successfully"

# Start remaining services
echo "Starting services..."
docker-compose -f docker-compose.prod.yml up --build -d

# Keep last 5 backups and remove older ones
cd /tmp/hostel-backup/db
ls -t *.sql | tail -n +6 | xargs rm -f 2>/dev/null || true

echo "Deployment completed!"