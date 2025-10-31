#!/bin/bash

# Deploy script for listings-service to dev.svetu.rs
# Deploys to /opt/listings-dev on svetu@svetu.rs server

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

# Configuration
SERVER="svetu@svetu.rs"
DEPLOY_DIR="/opt/listings-dev"
SERVICE_NAME="listings-service"
BINARY_NAME="listings-service"
HTTP_PORT="8086"
GRPC_PORT="50053"
METRICS_PORT="9093"
HEALTH_CHECK_RETRIES=6

log "🚀 Starting deployment of listings-service to dev.svetu.rs"

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
    git commit -m "Deploy listings-service to dev server" || warn "Nothing to commit"
fi

# 3. Push changes
log "⬆️  Pushing to origin/$CURRENT_BRANCH..."
if ! git push origin "$CURRENT_BRANCH"; then
    error "Failed to push changes. Aborting deployment."
    exit 1
fi

# 4. Build binary locally
log "🔨 Building binary locally..."
if ! make build; then
    error "Failed to build binary"
    exit 1
fi

# Verify binary exists
if [ ! -f "bin/$BINARY_NAME" ]; then
    error "Binary not found at bin/$BINARY_NAME"
    exit 1
fi

BINARY_SIZE=$(du -h "bin/$BINARY_NAME" | cut -f1)
log "✅ Binary built successfully (size: $BINARY_SIZE)"

# 5. Upload files to server
log "📤 Uploading files to server..."

# Upload binary
if ! scp "bin/$BINARY_NAME" "$SERVER:$DEPLOY_DIR/bin/"; then
    error "Failed to upload binary"
    exit 1
fi
log "✅ Binary uploaded"

# Upload docker-compose.yml
if ! scp docker-compose.yml "$SERVER:$DEPLOY_DIR/"; then
    error "Failed to upload docker-compose.yml"
    exit 1
fi
log "✅ docker-compose.yml uploaded"

# Upload .env.prod if exists
if [ -f ".env.prod" ]; then
    if ! scp .env.prod "$SERVER:$DEPLOY_DIR/.env"; then
        warn "Failed to upload .env.prod (will use existing .env on server)"
    else
        log "✅ .env.prod uploaded"
    fi
else
    warn ".env.prod not found, using existing .env on server"
fi

# Upload systemd service file
if [ -f "deployment/listings-service.service" ]; then
    if ! scp deployment/listings-service.service "$SERVER:/tmp/"; then
        warn "Failed to upload systemd service file"
    else
        log "✅ systemd service file uploaded to /tmp/"
    fi
fi

# 6. Deploy on server
log "🔄 Deploying on server..."

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

log "📂 Switching to deployment directory..."
cd "$DEPLOY_DIR"

# Fetch latest changes (if git repo)
if [ -d ".git" ]; then
    log "📥 Fetching latest changes from git..."
    git fetch origin

    TARGET_BRANCH="$CURRENT_BRANCH"
    log "🔀 Updating to branch: \$TARGET_BRANCH"

    if ! git reset --hard origin/\$TARGET_BRANCH; then
        error "Failed to reset to origin/\$TARGET_BRANCH"
        exit 1
    fi

    NEW_COMMIT=\$(git rev-parse HEAD)
    log "✅ Updated to commit: \${NEW_COMMIT:0:8}"
else
    warn "Not a git repository, skipping git update"
fi

# Check if systemd service exists
if [ -f "/tmp/listings-service.service" ]; then
    log "🔧 Installing systemd service..."
    sudo cp /tmp/listings-service.service /etc/systemd/system/
    sudo systemctl daemon-reload
    log "✅ systemd service installed"
fi

# Start/restart dependencies (PostgreSQL, Redis)
log "🔄 Starting dependencies (Docker Compose)..."
if ! docker-compose up -d postgres redis; then
    error "Failed to start dependencies"
    exit 1
fi

# Wait for dependencies to be healthy
log "⏳ Waiting for dependencies to be healthy..."
sleep 5

# Check PostgreSQL health
if ! docker exec listings_postgres pg_isready -U listings_user &>/dev/null; then
    warn "PostgreSQL not ready, waiting..."
    sleep 5
    if ! docker exec listings_postgres pg_isready -U listings_user &>/dev/null; then
        error "PostgreSQL failed to start"
        exit 1
    fi
fi
log "✅ PostgreSQL is healthy"

# Check Redis health
if ! docker exec listings_redis redis-cli ping &>/dev/null; then
    warn "Redis not ready, waiting..."
    sleep 5
    if ! docker exec listings_redis redis-cli ping &>/dev/null; then
        error "Redis failed to start"
        exit 1
    fi
fi
log "✅ Redis is healthy"

# Run database migrations
log "🗄️  Running database migrations..."
if ! make migrate-up; then
    error "Failed to run migrations"
    exit 1
fi
log "✅ Migrations applied"

# Stop old service
log "🛑 Stopping old service..."
if systemctl is-active --quiet $SERVICE_NAME; then
    sudo systemctl stop $SERVICE_NAME
    log "✅ Old service stopped"
else
    info "Service was not running"
fi

# Make binary executable
chmod +x bin/$BINARY_NAME
log "✅ Binary permissions set"

# Start service
log "🚀 Starting service..."
if ! sudo systemctl start $SERVICE_NAME; then
    error "Failed to start service"
    sudo journalctl -u $SERVICE_NAME -n 50 --no-pager
    exit 1
fi

# Enable service to start on boot
sudo systemctl enable $SERVICE_NAME &>/dev/null || true

log "✅ Service started"

# Wait for service to initialize
log "⏳ Waiting for service to initialize..."
sleep 5

# Health checks with retries
log "🏥 Checking service health..."
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

# Check HTTP REST endpoint
check_service "HTTP API" "http://localhost:$HTTP_PORT/health" || HEALTH_OK=false

# Check metrics endpoint
check_service "Metrics" "http://localhost:$METRICS_PORT/metrics" || HEALTH_OK=false

if [ "\$HEALTH_OK" = "false" ]; then
    error "Health checks failed!"
    warn "Check service logs:"
    echo "  sudo journalctl -u $SERVICE_NAME -n 100 --no-pager"
    echo ""
    warn "Check service status:"
    echo "  sudo systemctl status $SERVICE_NAME"
    exit 1
fi

# Show service status
log "📊 Service status:"
sudo systemctl status $SERVICE_NAME --no-pager -l | head -15

# Show process info
info "  HTTP Port: $HTTP_PORT"
info "  gRPC Port: $GRPC_PORT"
info "  Metrics Port: $METRICS_PORT"

log "🎉 Deployment completed successfully!"
ENDSSH

DEPLOY_EXIT_CODE=$?

# 7. Final status
if [ $DEPLOY_EXIT_CODE -eq 0 ]; then
    log "✅ Deployment complete!"
    echo ""
    log "📍 Service URLs:"
    log "  HTTP API: https://listings.dev.svetu.rs"
    log "  Metrics: http://svetu.rs:$METRICS_PORT/metrics (internal only)"
    log "  gRPC: svetu.rs:$GRPC_PORT (internal only)"
    echo ""
    log "📋 Useful commands:"
    log "  Status: ssh $SERVER 'sudo systemctl status $SERVICE_NAME'"
    log "  Logs: ssh $SERVER 'sudo journalctl -u $SERVICE_NAME -f'"
    log "  Stop: ssh $SERVER 'sudo systemctl stop $SERVICE_NAME'"
    log "  Restart: ssh $SERVER 'sudo systemctl restart $SERVICE_NAME'"
else
    error "Deployment failed with exit code $DEPLOY_EXIT_CODE"
    error "Check server logs for details:"
    echo "  ssh $SERVER 'sudo journalctl -u $SERVICE_NAME -n 100 --no-pager'"
    exit $DEPLOY_EXIT_CODE
fi
