#!/bin/bash

# Deploy script for dev.svetu.rs
# Improved version with proper error handling and environment checks

set -euo pipefail  # Exit on error, undefined vars, pipe failures

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" >&2
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO:${NC} $1"
}

# Load DB password from env or use default
DB_PASSWORD="${PGPASSWORD:-mX3g1XGhMRUZEX3l}"

# Configuration
SERVER="svetu@svetu.rs"
DEPLOY_DIR="/opt/svetu-dev"
BACKEND_PORT="3002"
FRONTEND_PORT="3003"
HEALTH_CHECK_RETRIES=6
REQUIRED_GO_VERSION="1.25"

log "🚀 Starting deployment to dev.svetu.rs"

# 1. Get current branch
CURRENT_BRANCH=$(git branch --show-current)
if [ -z "$CURRENT_BRANCH" ]; then
    error "Failed to get current branch"
    exit 1
fi
log "📌 Current branch: $CURRENT_BRANCH"

# 2. Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    log "📝 Committing current changes..."
    git add -A
    git commit -m "Deploy to dev server" || warn "Nothing to commit"
fi

# 3. Push changes
log "⬆️  Pushing to origin/$CURRENT_BRANCH..."
if ! git push origin "$CURRENT_BRANCH"; then
    error "Failed to push changes. Aborting deployment."
    exit 1
fi

# 4. Create database dump
log "💾 Creating database dump..."
DUMP_FILE="svetubd_dump_$(date +%Y%m%d_%H%M%S).sql"
DUMP_PATH="/tmp/$DUMP_FILE"

if ! PGPASSWORD="$DB_PASSWORD" pg_dump -h localhost -U postgres -d svetubd \
    --no-owner --no-acl --column-inserts --inserts -f "$DUMP_PATH"; then
    error "Failed to create database dump"
    exit 1
fi
log "✅ Database dumped to $DUMP_PATH ($(du -h "$DUMP_PATH" | cut -f1))"

# 5. Get Mapbox token from local env
MAPBOX_TOKEN=""
if [ -f "/data/hostel-booking-system/frontend/svetu/.env.local" ]; then
    MAPBOX_TOKEN=$(grep "^NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=" /data/hostel-booking-system/frontend/svetu/.env.local 2>/dev/null | cut -d'=' -f2 || true)
    if [ -n "$MAPBOX_TOKEN" ]; then
        log "🗺️  Mapbox token found (will sync to server)"
    fi
fi

# 6. Upload dump to server
log "📤 Uploading database dump to server..."
if ! scp "$DUMP_PATH" "$SERVER:/tmp/"; then
    error "Failed to upload dump to server"
    rm -f "$DUMP_PATH"
    exit 1
fi

# 7. Deploy on server
log "🔄 Deploying on server..."

# Create heredoc with proper variable substitution
ssh "$SERVER" /bin/bash <<ENDSSH
set -euo pipefail

# Enable verbose error tracking
trap 'echo "❌ Error on line \$LINENO. Exit code: \$?" >&2' ERR

# Colors for remote logging
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "\${GREEN}[Server \$(date +'%H:%M:%S')]\${NC} \$1"; }
error() { echo -e "\${RED}[Server \$(date +'%H:%M:%S')] ERROR:\${NC} \$1" >&2; }
warn() { echo -e "\${YELLOW}[Server \$(date +'%H:%M:%S')] WARNING:\${NC} \$1"; }
info() { echo -e "\${BLUE}[Server \$(date +'%H:%M:%S')] INFO:\${NC} \$1"; }
debug() { echo -e "\${BLUE}[Server \$(date +'%H:%M:%S')] DEBUG:\${NC} \$1"; }

# Check Go version
log "🔍 Checking Go version..."
CURRENT_GO_VERSION=\$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "0.0")
REQUIRED_VERSION="$REQUIRED_GO_VERSION"

if [ "\$CURRENT_GO_VERSION" != "\$REQUIRED_VERSION" ]; then
    warn "Go version mismatch: found \$CURRENT_GO_VERSION, required \$REQUIRED_VERSION"
    log "📥 Installing Go \$REQUIRED_VERSION..."

    cd /tmp
    wget -q https://go.dev/dl/go\${REQUIRED_VERSION}.0.linux-amd64.tar.gz

    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go\${REQUIRED_VERSION}.0.linux-amd64.tar.gz

    sudo rm -f /usr/bin/go
    sudo ln -s /usr/local/go/bin/go /usr/bin/go

    log "✅ Go \$REQUIRED_VERSION installed"
    go version
else
    log "✅ Go version is correct: \$CURRENT_GO_VERSION"
fi

# Check Git configuration for private repos
log "🔍 Checking Git configuration..."
if ! git config --global --get url."git@github.com:".insteadOf &>/dev/null; then
    warn "Git not configured for SSH, fixing..."
    git config --global url."git@github.com:".insteadOf "https://github.com/"
    log "✅ Git configured to use SSH for GitHub"
else
    log "✅ Git already configured for SSH"
fi

log "📂 Switching to deployment directory..."
cd "$DEPLOY_DIR"

# Save current commit for potential rollback
PREVIOUS_COMMIT=\$(git rev-parse HEAD)
log "💾 Current commit (for rollback): \${PREVIOUS_COMMIT:0:8}"

# Fetch and reset to target branch
log "📥 Fetching latest changes..."
git fetch origin

TARGET_BRANCH="$CURRENT_BRANCH"
log "🔀 Deploying branch: \$TARGET_BRANCH"

if ! git reset --hard origin/\$TARGET_BRANCH; then
    error "Failed to reset to origin/\$TARGET_BRANCH"
    exit 1
fi

NEW_COMMIT=\$(git rev-parse HEAD)
log "✅ Updated to commit: \${NEW_COMMIT:0:8}"

# Database restore
log "💾 Restoring database..."
DUMP_FILE="/tmp/$DUMP_FILE"

if [ ! -f "\$DUMP_FILE" ]; then
    error "Dump file not found: \$DUMP_FILE"
    exit 1
fi

log "🗄️  Clearing database schema..."
if ! docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db \
    -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" &>/tmp/db_clear.log; then
    error "Failed to clear database schema"
    cat /tmp/db_clear.log
    exit 1
fi

log "📥 Loading database dump..."
if ! docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db \
    < "\$DUMP_FILE" &>/tmp/db_load.log; then
    error "Failed to load database dump"
    tail -20 /tmp/db_load.log
    exit 1
fi

log "✅ Database restored successfully"
tail -5 /tmp/db_load.log | sed 's/^/  /'

# Fix dirty migrations
docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db \
    -c "UPDATE schema_migrations SET dirty = false WHERE dirty = true;" >/dev/null 2>&1 || true

# Sync Mapbox token if provided
if [ -n "$MAPBOX_TOKEN" ]; then
    log "🗺️  Syncing Mapbox token..."
    ENV_FILE="$DEPLOY_DIR/frontend/svetu/.env.local"
    if [ -f "\$ENV_FILE" ]; then
        # Update or append token
        if grep -q "^NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=" "\$ENV_FILE"; then
            sed -i "s|^NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=.*|NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=$MAPBOX_TOKEN|" "\$ENV_FILE"
        else
            echo "NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=$MAPBOX_TOKEN" >> "\$ENV_FILE"
        fi
        log "✅ Mapbox token updated"
    else
        warn "Env file not found: \$ENV_FILE"
    fi
fi

# Kill old backend processes before restart
cd "$DEPLOY_DIR/backend" || { error "Failed to cd to backend dir"; exit 1; }
log "🔪 Killing old backend processes..."
pkill -9 -f "bin/api_dev" 2>/dev/null || true
sleep 2

# Restart backend
log "🔄 Restarting backend..."
if ! timeout 120 make dev-restart &>/tmp/backend_restart.log; then
    error "Failed to restart backend (timeout or error)"
    tail -50 /tmp/backend_restart.log

    # Check if it's just a "port already in use" issue
    if grep -q "bind: address already in use" /tmp/backend_restart.log; then
        warn "Port already in use - backend might be already running from previous deploy"
        info "Checking if backend is actually running..."

        if pgrep -f "bin/api_dev" > /dev/null; then
            log "✅ Backend process is running (build succeeded, restart skipped)"
        else
            error "Backend not running despite port being in use"
            exit 1
        fi
    else
        exit 1
    fi
else
    log "✅ Backend restarted"
fi

# Kill old frontend processes before restart
log "🔪 Killing old frontend processes..."
# Убиваем все возможные варианты процессов Next.js
pkill -9 -f "yarn dev.*3003" 2>/dev/null || true
pkill -9 -f "yarn start.*3003" 2>/dev/null || true
pkill -9 -f "next dev.*3003" 2>/dev/null || true
pkill -9 -f "next start.*3003" 2>/dev/null || true
pkill -9 -f "next-server.*3003" 2>/dev/null || true
pkill -9 -f "node.*next.*3003" 2>/dev/null || true
sleep 3

# Verify port 3003 is free (более надежная проверка)
log "🔍 Checking if port 3003 is free..."
PORT_CHECK_ATTEMPTS=0
MAX_PORT_ATTEMPTS=5

while netstat -tlnp 2>/dev/null | grep -q ":3003 " && [ \$PORT_CHECK_ATTEMPTS -lt \$MAX_PORT_ATTEMPTS ]; do
    warn "Port 3003 still occupied (attempt \$((PORT_CHECK_ATTEMPTS + 1))/\$MAX_PORT_ATTEMPTS), forcing cleanup..."
    fuser -k 3003/tcp 2>/dev/null || true
    sleep 2
    PORT_CHECK_ATTEMPTS=\$((PORT_CHECK_ATTEMPTS + 1))
done

if netstat -tlnp 2>/dev/null | grep -q ":3003 "; then
    error "Failed to free port 3003 after \$MAX_PORT_ATTEMPTS attempts"
    warn "Processes still using port 3003:"
    fuser -v 3003/tcp 2>&1 || true
    warn "You may need to manually kill the process or reboot"
    exit 1
fi

log "✅ Port 3003 is free"

# Restart frontend with production build
log "🔄 Restarting frontend (production build)..."
cd "$DEPLOY_DIR/frontend/svetu" || { error "Failed to cd to frontend dir"; exit 1; }

# КРИТИЧНО: Удаляем старый .next чтобы не использовать недельный билд!
log "🧹 Removing old .next build directory..."
rm -rf .next
log "✅ Old build removed"

# Билд с увеличенным таймаутом (10 минут вместо 5)
log "🏗️  Building fresh production version (timeout: 10 min)..."
if ! timeout 600 yarn build &>/tmp/frontend_build.log; then
    error "Failed to build frontend (timeout or error)"
    tail -100 /tmp/frontend_build.log
    error "BUILD IS MANDATORY - deployment aborted!"
    error "Old .next was deleted, cannot fallback to old build"
    exit 1
fi
log "✅ Frontend built successfully"

# Проверяем свежесть .next (должна быть не старше 2 минут)
NEXT_AGE=\$(find .next -maxdepth 0 -mmin -2 2>/dev/null | wc -l)
if [ "\$NEXT_AGE" -eq 0 ]; then
    error ".next directory is too old or missing!"
    error "Build might have failed silently"
    exit 1
fi
log "✅ .next is fresh (created within last 2 minutes)"

# Останавливаем старый процесс
log "🔪 Stopping old frontend process..."
lsof -ti:3003 | xargs -r kill 2>/dev/null || true
fuser -k 3003/tcp 2>/dev/null || true
sleep 2

# Запускаем production сервер
log "🚀 Starting production server on port 3003..."
if ! nohup yarn start -p 3003 > frontend-dev.log 2>&1 & then
    error "Failed to start frontend server"
    tail -50 frontend-dev.log
    exit 1
fi
sleep 3
log "✅ Frontend restarted (production mode with FRESH build)"

# Проверяем что frontend действительно запустился
sleep 5
if ! pgrep -f "next.*3003" > /dev/null; then
    error "Frontend process not found after restart"
    tail -50 /tmp/frontend_restart.log
    exit 1
fi
log "✅ Frontend process is running"

# Clean up old dumps (keep last 3)
log "🧹 Cleaning old dumps..."
ls -t /tmp/svetubd_dump_*.sql 2>/dev/null | tail -n +4 | xargs rm -f 2>/dev/null || true

# Wait for services to fully start
log "⏳ Waiting for services to initialize..."
sleep 5

# Health checks with retries
log "🏥 Checking services health..."
check_service() {
    local name=\$1
    local url=\$2
    local retries=$HEALTH_CHECK_RETRIES
    local wait=10

    for i in \$(seq 1 \$retries); do
        HTTP_CODE=\$(curl -s -o /dev/null -w "%{http_code}" "\$url" 2>/dev/null || echo "000")

        # Accept 200 (OK), 307 (redirect), 404 (route not found but server running)
        if echo "\$HTTP_CODE" | grep -qE "^(200|307|404)$"; then
            log "✅ \$name is healthy (HTTP \$HTTP_CODE)"
            return 0
        fi

        if [ \$i -lt \$retries ]; then
            warn "\$name not ready yet (HTTP \$HTTP_CODE, attempt \$i/\$retries)..."
            sleep \$wait
        else
            error "\$name failed health check after \$retries attempts (last HTTP: \$HTTP_CODE)"
            return 1
        fi
    done
}

HEALTH_OK=true
check_service "Backend" "http://localhost:$BACKEND_PORT/" || HEALTH_OK=false
check_service "Frontend" "http://localhost:$FRONTEND_PORT" || HEALTH_OK=false

if [ "\$HEALTH_OK" = "false" ]; then
    error "Health checks failed!"
    warn "Check logs for details:"
    echo "  ssh $SERVER 'tail -100 /tmp/backend-dev.log'"
    echo "  ssh $SERVER 'tail -100 /tmp/frontend-dev.log'"
    echo "  ssh $SERVER 'cd /opt/svetu-dev/backend && tail -50 api_dev.log'"
    echo "  ssh $SERVER 'cd /opt/svetu-dev/frontend/svetu && tail -50 frontend-dev.log'"
    echo ""
    warn "Manual restart commands:"
    echo "  ssh $SERVER 'cd /opt/svetu-dev/backend && make dev-restart'"
    echo "  ssh $SERVER 'cd /opt/svetu-dev/frontend/svetu && make dev-restart'"
    echo ""
    warn "If needed, rollback with:"
    echo "  ssh $SERVER 'cd $DEPLOY_DIR && git reset --hard \$PREVIOUS_COMMIT'"
    exit 1
fi

# Show deployed version
BACKEND_VERSION=\$(curl -s http://localhost:$BACKEND_PORT/ 2>/dev/null | head -1 || echo "unknown")
log "🎯 Deployed backend version: \$BACKEND_VERSION"
log "🎯 Deployed commit: \${NEW_COMMIT:0:8}"

# Show process info
log "📊 Process status:"
info "  Backend PID: \$(pgrep -f 'bin/api_dev' || echo 'not found')"
info "  Frontend PID: \$(pgrep -f 'next.*3003' || echo 'not found')"

log "🎉 Deployment completed successfully!"
ENDSSH

DEPLOY_EXIT_CODE=$?

# 8. Clean up local dump
rm -f "$DUMP_PATH"
log "🧹 Local dump cleaned up"

# 9. Final status
if [ $DEPLOY_EXIT_CODE -eq 0 ]; then
    log "✅ Deployment complete!"
    echo ""
    log "📍 Site: https://dev.svetu.rs"
    log "📍 API: https://devapi.svetu.rs"
    echo ""
    log "📊 Deployed:"
    log "  Branch: $CURRENT_BRANCH"
    log "  Commit: $(git rev-parse --short HEAD)"
    echo ""
    log "📋 Useful commands:"
    log "  Logs: ssh $SERVER 'tail -f /tmp/backend-dev.log'"
    log "  Backend: curl https://devapi.svetu.rs/"
    log "  Frontend: curl -I https://dev.svetu.rs"
else
    error "Deployment failed with exit code $DEPLOY_EXIT_CODE"
    error "Check server logs for details:"
    echo "  ssh $SERVER 'tail -100 /tmp/backend-dev.log'"
    echo "  ssh $SERVER 'tail -100 /tmp/frontend-dev.log'"
    echo "  ssh $SERVER 'tail -50 /tmp/backend_restart.log'"
    echo "  ssh $SERVER 'tail -50 /tmp/frontend_restart.log'"
    exit $DEPLOY_EXIT_CODE
fi
