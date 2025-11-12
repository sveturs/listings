#!/bin/bash
#
# Quick Performance Check Script
# Performs a quick 10-second test on critical endpoints
#
# Usage: ./quick-check.sh [grpc-address]
#

set -euo pipefail

GRPC_ADDR=${1:-localhost:8086}

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "=========================================="
echo "  Quick Performance Check"
echo "=========================================="
echo "gRPC Server: $GRPC_ADDR"
echo "Duration: 10 seconds per endpoint"
echo ""

# Check prerequisites
if ! command -v ghz &> /dev/null; then
    echo -e "${RED}ERROR: ghz not installed${NC}"
    echo "Install: go install github.com/bojand/ghz/cmd/ghz@latest"
    exit 1
fi

# Test function
quick_test() {
    local method=$1
    local data=$2
    local name=$3

    echo -n "Testing $name... "

    result=$(ghz --insecure \
        --proto /p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto \
        --import-paths /p/github.com/sveturs/listings/api/proto \
        --call "$method" \
        --data "$data" \
        --duration 10s \
        --concurrency 5 \
        --rps 50 \
        --format json \
        "$GRPC_ADDR" 2>&1)

    if echo "$result" | jq -e '.error' > /dev/null 2>&1; then
        echo -e "${RED}FAILED${NC}"
        echo "$result" | jq -r '.error'
        return 1
    fi

    local p95=$(echo "$result" | jq -r '.histogram[] | select(.mark == 0.95) | .latency' 2>/dev/null || echo "0")
    local p95_ms=$(echo "scale=2; $p95 / 1000000" | bc)
    local rps=$(echo "$result" | jq -r '.rps' 2>/dev/null || echo "0")

    if (( $(echo "$p95_ms > 100" | bc -l) )); then
        echo -e "${RED}SLOW${NC} (P95: ${p95_ms}ms, RPS: $rps)"
    elif (( $(echo "$p95_ms > 50" | bc -l) )); then
        echo -e "${YELLOW}OK${NC} (P95: ${p95_ms}ms, RPS: $rps)"
    else
        echo -e "${GREEN}FAST${NC} (P95: ${p95_ms}ms, RPS: $rps)"
    fi
}

# Run quick tests
quick_test "listings.v1.ListingsService/GetProduct" '{"id":1}' "GetProduct"
quick_test "listings.v1.ListingsService/CheckStockAvailability" '{"items":[{"product_id":1,"quantity":5}]}' "CheckStock"
quick_test "listings.v1.ListingsService/GetAllCategories" '{}' "GetCategories"
quick_test "listings.v1.ListingsService/ListProducts" '{"storefront_id":1,"limit":20}' "ListProducts"

echo ""
echo "Quick check complete!"
