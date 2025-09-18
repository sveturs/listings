#!/bin/bash

# Deploy script for dev.svetu.rs
set -e

echo "🚀 Starting deployment to dev.svetu.rs"

# 1. Commit and push current changes
echo "📝 Committing current changes..."
CURRENT_BRANCH=$(git branch --show-current)
git add -A
git commit -m "Deploy to dev server" || echo "Nothing to commit"
git push origin $CURRENT_BRANCH

echo "Current branch: $CURRENT_BRANCH"

# 2. Create database dump
echo "💾 Creating database dump..."
DUMP_FILE="svetubd_dump_$(date +%Y%m%d_%H%M%S).sql"
PGPASSWORD=mX3g1XGhMRUZEX3l pg_dump -h localhost -U postgres -d svetubd --no-owner --no-acl --column-inserts --inserts -f /tmp/$DUMP_FILE
echo "Database dumped to /tmp/$DUMP_FILE"

# 3. Upload dump to server
echo "📤 Uploading database dump to server..."
scp /tmp/$DUMP_FILE svetu@svetu.rs:/tmp/

# 4. Deploy on server
echo "🔄 Deploying on server..."
ssh svetu@svetu.rs << 'ENDSSH'
set -e

echo "Switching to deployment directory..."
cd /opt/svetu-dev

# Get current branch from local
BRANCH=$(git branch --show-current)
echo "Current branch on server: $BRANCH"

# Pull latest changes
echo "Pulling latest changes..."
git fetch origin
git pull origin $BRANCH

# Restore database
echo "Restoring database..."
DUMP_FILE=$(ls -t /tmp/svetubd_dump_*.sql | head -1)
echo "Using dump: $DUMP_FILE"

# Clear and restore database in docker
cd /opt/svetu-dev
docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db < $DUMP_FILE

# Update schema_migrations to prevent migration issues
docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -c "UPDATE schema_migrations SET dirty = false WHERE dirty = true;" 2>/dev/null || true

# Restart backend
echo "Restarting backend..."
cd /opt/svetu-dev/backend
make dev-restart || {
    echo "make dev-restart not found, using custom restart"
    pkill -f 'go run ./cmd/api/main.go' || true
    lsof -ti:3002 | xargs kill -9 2>/dev/null || true
    screen -S backend-dev -X quit 2>/dev/null || true
    screen -dmS backend-dev bash -c 'go run ./cmd/api/main.go 2>&1 | tee /tmp/backend-dev.log'
}

# Restart frontend
echo "Restarting frontend..."
cd /opt/svetu-dev/frontend/svetu
make dev-restart || {
    echo "make dev-restart not found, using custom restart"
    pkill -f 'yarn dev.*3003' || true
    lsof -ti:3003 | xargs kill -9 2>/dev/null || true
    screen -S frontend-dev -X quit 2>/dev/null || true
    screen -dmS frontend-dev bash -c 'yarn dev -p 3003 2>&1 | tee /tmp/frontend-dev.log'
}

echo "✅ Services restarted"

# Clean up old dumps (keep last 3)
ls -t /tmp/svetubd_dump_*.sql | tail -n +4 | xargs rm -f 2>/dev/null || true

# Check services
sleep 5
echo "Checking services..."
curl -s -o /dev/null -w "Backend: %{http_code}\n" http://localhost:3002/health || echo "Backend: NOT READY"
curl -s -o /dev/null -w "Frontend: %{http_code}\n" http://localhost:3003 || echo "Frontend: NOT READY"

ENDSSH

# 5. Clean up local dump
rm -f /tmp/$DUMP_FILE

echo "✅ Deployment complete!"
echo "📍 Site: https://dev.svetu.rs"
echo "📍 API: https://devapi.svetu.rs"
echo "📍 Logs: ssh svetu@svetu.rs 'tail -f /tmp/backend-dev.log'"