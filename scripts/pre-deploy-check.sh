#!/bin/bash
set -e

echo "🔍 Running pre-deployment checks..."
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track if any check fails
HAS_ERRORS=0

# Function to print success
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to print error
error() {
    echo -e "${RED}❌ $1${NC}"
    HAS_ERRORS=1
}

# Function to print warning
warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# 1. Check if we're in the backend directory
echo ""
echo "1. Checking directory structure..."
if [ ! -f "go.mod" ]; then
    error "go.mod not found. Are you in the backend directory?"
else
    success "Backend directory structure valid"
fi

# 2. Check Go installation
echo ""
echo "2. Checking Go installation..."
if ! command -v go &> /dev/null; then
    error "Go is not installed"
else
    GO_VERSION=$(go version | awk '{print $3}')
    success "Go installed: $GO_VERSION"
fi

# 3. Check dependencies
echo ""
echo "3. Checking Go dependencies..."
if go mod verify &> /dev/null; then
    success "Go modules verified"
else
    error "Go module verification failed. Run 'go mod tidy'"
fi

# 4. Run tests
echo ""
echo "4. Running unit tests..."
if make test &> /tmp/test-output.log; then
    success "All unit tests passed"
else
    error "Unit tests failed. Check /tmp/test-output.log"
    cat /tmp/test-output.log
fi

# 5. Run integration tests
echo ""
echo "5. Running integration tests..."
if make test-integration &> /tmp/integration-test-output.log; then
    success "All integration tests passed"
else
    warning "Integration tests failed. Check /tmp/integration-test-output.log"
    # Not a blocker, just a warning
fi

# 6. Check linting
echo ""
echo "6. Running linter..."
if make lint &> /tmp/lint-output.log; then
    success "No lint errors"
else
    error "Linting failed. Run 'make lint' to see errors"
    cat /tmp/lint-output.log
fi

# 7. Check formatting
echo ""
echo "7. Checking code formatting..."
UNFORMATTED=$(gofmt -l . 2>/dev/null | grep -v vendor || true)
if [ -z "$UNFORMATTED" ]; then
    success "Code is properly formatted"
else
    error "Code formatting issues found. Run 'make format'"
    echo "$UNFORMATTED"
fi

# 8. Check migrations status
echo ""
echo "8. Checking database migrations..."
if [ -x "./migrator" ]; then
    MIGRATION_STATUS=$(./migrator status 2>&1 || true)
    if echo "$MIGRATION_STATUS" | grep -q "All migrations applied"; then
        success "All migrations applied"
    else
        warning "Pending migrations detected"
        echo "$MIGRATION_STATUS"
    fi
else
    warning "Migrator not found. Skipping migration check"
fi

# 9. Check Docker image build
echo ""
echo "9. Testing Docker image build..."
if [ -f "Dockerfile" ]; then
    if docker build -t svetu-backend:pre-deploy-test . &> /tmp/docker-build.log; then
        success "Docker image builds successfully"
        # Clean up test image
        docker rmi svetu-backend:pre-deploy-test &> /dev/null || true
    else
        error "Docker build failed. Check /tmp/docker-build.log"
        tail -n 50 /tmp/docker-build.log
    fi
else
    warning "Dockerfile not found. Skipping Docker build check"
fi

# 10. Check Prometheus config (if exists)
echo ""
echo "10. Checking Prometheus configuration..."
if [ -f "../monitoring/prometheus.yml" ]; then
    if command -v promtool &> /dev/null; then
        if promtool check config ../monitoring/prometheus.yml &> /dev/null; then
            success "Prometheus config is valid"
        else
            error "Invalid Prometheus configuration"
        fi
    else
        warning "promtool not installed. Skipping Prometheus config check"
    fi
else
    warning "Prometheus config not found. Skipping check"
fi

# 11. Check environment variables
echo ""
echo "11. Checking required environment variables..."
REQUIRED_VARS=(
    "POSTGRES_DSN"
    "REDIS_URL"
    "JWT_SECRET"
    "MINIO_ENDPOINT"
    "OPENSEARCH_URL"
)

for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        warning "Environment variable $var is not set"
    else
        success "Environment variable $var is set"
    fi
done

# 12. Check service connectivity (if in deployment environment)
echo ""
echo "12. Checking service connectivity..."

# PostgreSQL
if [ -n "$POSTGRES_DSN" ]; then
    if psql "$POSTGRES_DSN" -c "SELECT 1" &> /dev/null; then
        success "PostgreSQL connection OK"
    else
        error "Cannot connect to PostgreSQL"
    fi
fi

# Redis
if [ -n "$REDIS_URL" ]; then
    if redis-cli -u "$REDIS_URL" PING &> /dev/null; then
        success "Redis connection OK"
    else
        error "Cannot connect to Redis"
    fi
fi

# 13. Check disk space
echo ""
echo "13. Checking disk space..."
DISK_USAGE=$(df -h . | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -lt 80 ]; then
    success "Disk space OK ($DISK_USAGE% used)"
else
    warning "Disk space high: $DISK_USAGE% used"
fi

# 14. Check memory
echo ""
echo "14. Checking available memory..."
if command -v free &> /dev/null; then
    MEMORY_AVAILABLE=$(free -m | awk 'NR==2 {print $7}')
    if [ "$MEMORY_AVAILABLE" -gt 500 ]; then
        success "Memory available: ${MEMORY_AVAILABLE}MB"
    else
        warning "Low memory: ${MEMORY_AVAILABLE}MB available"
    fi
fi

# 15. Check if required ports are available
echo ""
echo "15. Checking port availability..."
REQUIRED_PORTS=(3000 3001)
for port in "${REQUIRED_PORTS[@]}"; do
    if ! netstat -tuln 2>/dev/null | grep -q ":$port "; then
        success "Port $port is available"
    else
        warning "Port $port is already in use"
    fi
done

# Summary
echo ""
echo "========================================"
if [ $HAS_ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ All pre-deployment checks passed!${NC}"
    echo ""
    echo "Ready to deploy 🚀"
    exit 0
else
    echo -e "${RED}❌ Some checks failed!${NC}"
    echo ""
    echo "Please fix the errors above before deploying."
    exit 1
fi
