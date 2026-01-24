#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Load environment variables from .env
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

echo ""
echo -e "${BLUE}=========================================="
echo -e "  Listings Service - Environment Status"
echo -e "==========================================${NC}"
echo ""

# Function to check if port is in use
check_port() {
    local port=$1
    local name=$2
    local protocol=$3

    if lsof -ti:$port > /dev/null 2>&1; then
        echo -e "  ${GREEN}✓${NC} $name"
        echo -e "    Port: ${GREEN}$port${NC} ($protocol)"
    else
        echo -e "  ${RED}✗${NC} $name"
        echo -e "    Port: ${RED}$port${NC} ($protocol) - ${RED}NOT RUNNING${NC}"
    fi
}

# Function to check docker container
check_docker() {
    local container=$1
    local name=$2
    local port=$3

    if docker-compose ps --format json 2>/dev/null | grep -q "\"$container\".*running"; then
        echo -e "  ${GREEN}✓${NC} $name"
        echo -e "    Port: ${GREEN}$port${NC}"
    elif docker-compose ps 2>/dev/null | grep -q "$container.*Up"; then
        echo -e "  ${GREEN}✓${NC} $name"
        echo -e "    Port: ${GREEN}$port${NC}"
    else
        echo -e "  ${RED}✗${NC} $name"
        echo -e "    Port: ${RED}$port${NC} - ${RED}NOT RUNNING${NC}"
    fi
}

echo -e "${YELLOW}Docker Dependencies:${NC}"
check_docker "postgres" "PostgreSQL" "${VONDILISTINGS_DB_PORT:-35434}"
check_docker "redis" "Redis" "${VONDILISTINGS_REDIS_PORT:-36380}"
echo ""

echo -e "${YELLOW}Service:${NC}"
check_port "${VONDILISTINGS_HTTP_PORT:-8086}" "Listings Service (HTTP)" "HTTP"
check_port "${VONDILISTINGS_GRPC_PORT:-50053}" "Listings Service (gRPC)" "gRPC"
echo ""

echo -e "${BLUE}==========================================${NC}"
echo ""
echo -e "${YELLOW}Quick commands:${NC}"
echo "  make start      - Start service"
echo "  make stop       - Stop service"
echo "  make reset-all  - Full reset (stop + clean DB + migrate + start)"
echo ""
