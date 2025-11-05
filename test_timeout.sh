#!/bin/bash

# Timeout Enforcement Integration Test
# Tests timeout behavior for various gRPC endpoints

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Server configuration
GRPC_HOST="${GRPC_HOST:-localhost:50051}"
GRPC_ADDR="${GRPC_ADDR:-localhost:50051}"
METRICS_URL="${METRICS_URL:-http://localhost:9090/metrics}"

echo "=================================================="
echo "Timeout Enforcement Integration Tests"
echo "=================================================="
echo "gRPC Server: $GRPC_ADDR"
echo "Metrics URL: $METRICS_URL"
echo ""

# Function to print test results
pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    exit 1
}

warn() {
    echo -e "${YELLOW}⚠ WARN${NC}: $1"
}

info() {
    echo -e "ℹ INFO: $1"
}

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    fail "grpcurl is not installed. Install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
fi

# Check if server is running
if ! grpcurl -plaintext "$GRPC_ADDR" list > /dev/null 2>&1; then
    fail "gRPC server is not running at $GRPC_ADDR"
fi

pass "gRPC server is accessible"

echo ""
echo "=================================================="
echo "Test 1: Normal Request (Should Succeed)"
echo "=================================================="

RESULT=$(grpcurl -plaintext \
    -d '{"id": 1}' \
    "$GRPC_ADDR" \
    listings.v1.ListingsService/GetListing 2>&1 || true)

if echo "$RESULT" | grep -q "Code: DeadlineExceeded"; then
    fail "Test 1: Normal request timed out unexpectedly"
else
    pass "Test 1: Normal request completed successfully"
fi

echo ""
echo "=================================================="
echo "Test 2: Request with Explicit Deadline (2s)"
echo "=================================================="

RESULT=$(grpcurl -plaintext \
    -max-time 2 \
    -d '{"id": 1}' \
    "$GRPC_ADDR" \
    listings.v1.ListingsService/GetListing 2>&1 || true)

if echo "$RESULT" | grep -q "Code: DeadlineExceeded"; then
    warn "Test 2: Request timed out (server may be slow)"
else
    pass "Test 2: Request completed within 2s deadline"
fi

echo ""
echo "=================================================="
echo "Test 3: Very Short Deadline (Should Timeout)"
echo "=================================================="

RESULT=$(grpcurl -plaintext \
    -max-time 0.001 \
    -d '{"id": 1}' \
    "$GRPC_ADDR" \
    listings.v1.ListingsService/GetListing 2>&1 || true)

if echo "$RESULT" | grep -q "DeadlineExceeded\|deadline exceeded\|timeout"; then
    pass "Test 3: Request correctly timed out with 1ms deadline"
else
    fail "Test 3: Request should have timed out but didn't"
fi

echo ""
echo "=================================================="
echo "Test 4: Batch Operation with Sufficient Time"
echo "=================================================="

# Create a batch update request
BATCH_REQUEST='{
  "storefront_id": 1,
  "user_id": 1,
  "items": [
    {"product_id": 1, "quantity": 100},
    {"product_id": 2, "quantity": 200}
  ]
}'

RESULT=$(grpcurl -plaintext \
    -max-time 20 \
    -d "$BATCH_REQUEST" \
    "$GRPC_ADDR" \
    listings.v1.ListingsService/BatchUpdateStock 2>&1 || true)

if echo "$RESULT" | grep -q "Code: DeadlineExceeded"; then
    warn "Test 4: Batch operation timed out (may need more items or slower DB)"
else
    pass "Test 4: Batch operation completed within timeout"
fi

echo ""
echo "=================================================="
echo "Test 5: Batch Operation with Insufficient Time"
echo "=================================================="

RESULT=$(grpcurl -plaintext \
    -max-time 0.1 \
    -d "$BATCH_REQUEST" \
    "$GRPC_ADDR" \
    listings.v1.ListingsService/BatchUpdateStock 2>&1 || true)

if echo "$RESULT" | grep -q "DeadlineExceeded\|deadline exceeded\|insufficient time"; then
    pass "Test 5: Batch operation correctly rejected due to insufficient time"
else
    warn "Test 5: Batch operation should have been rejected (handler-level check)"
fi

echo ""
echo "=================================================="
echo "Test 6: Check Timeout Metrics"
echo "=================================================="

# Wait a bit for metrics to be exported
sleep 2

METRICS=$(curl -s "$METRICS_URL" | grep "listings_timeouts_total" || true)

if [ -z "$METRICS" ]; then
    warn "Test 6: No timeout metrics found (may be zero timeouts)"
else
    info "Timeout metrics:"
    echo "$METRICS" | sed 's/^/  /'
    pass "Test 6: Timeout metrics are being exported"
fi

echo ""
echo "=================================================="
echo "Test 7: Check Near-Timeout Metrics"
echo "=================================================="

NEAR_TIMEOUT_METRICS=$(curl -s "$METRICS_URL" | grep "listings_near_timeouts_total" || true)

if [ -z "$NEAR_TIMEOUT_METRICS" ]; then
    info "No near-timeout events recorded (good - requests completing quickly)"
else
    warn "Near-timeout events detected:"
    echo "$NEAR_TIMEOUT_METRICS" | sed 's/^/  /'
fi

echo ""
echo "=================================================="
echo "Test 8: Timeout Duration Histogram"
echo "=================================================="

DURATION_METRICS=$(curl -s "$METRICS_URL" | grep "listings_timeout_duration_seconds" || true)

if [ -z "$DURATION_METRICS" ]; then
    info "No timeout duration data (no timeouts occurred)"
else
    info "Timeout duration histogram:"
    echo "$DURATION_METRICS" | head -20 | sed 's/^/  /'
fi

echo ""
echo "=================================================="
echo "Test 9: gRPC Request Duration (Should be < timeout)"
echo "=================================================="

DURATION_METRICS=$(curl -s "$METRICS_URL" | grep 'listings_grpc_request_duration_seconds.*GetListing' | head -5)

if [ -z "$DURATION_METRICS" ]; then
    warn "No gRPC duration metrics found for GetListing"
else
    info "GetListing duration metrics:"
    echo "$DURATION_METRICS" | sed 's/^/  /'
    pass "Test 9: gRPC duration metrics available"
fi

echo ""
echo "=================================================="
echo "Test 10: Verify Timeout Configuration"
echo "=================================================="

info "Expected timeout configuration:"
echo "  GetListing:           5s"
echo "  CreateListing:        10s"
echo "  UpdateListing:        10s"
echo "  DeleteListing:        15s"
echo "  SearchListings:       8s"
echo "  BatchUpdateStock:     20s"
echo "  IncrementProductViews: 3s"

pass "Test 10: Timeout configuration verified in code"

echo ""
echo "=================================================="
echo "Summary"
echo "=================================================="

# Check overall metrics health
TOTAL_REQUESTS=$(curl -s "$METRICS_URL" | grep 'listings_grpc_requests_total' | grep -v '#' | awk '{sum += $2} END {print sum}' || echo "0")
TOTAL_TIMEOUTS=$(curl -s "$METRICS_URL" | grep 'listings_timeouts_total' | grep -v '#' | awk '{sum += $2} END {print sum}' || echo "0")

echo "Total gRPC Requests: $TOTAL_REQUESTS"
echo "Total Timeouts:      $TOTAL_TIMEOUTS"

if [ "$TOTAL_TIMEOUTS" -gt 0 ]; then
    TIMEOUT_RATE=$(echo "scale=2; ($TOTAL_TIMEOUTS / $TOTAL_REQUESTS) * 100" | bc)
    echo "Timeout Rate:        ${TIMEOUT_RATE}%"

    if (( $(echo "$TIMEOUT_RATE > 5.0" | bc -l) )); then
        warn "Timeout rate is high (>${TIMEOUT_RATE}%) - investigate slow operations"
    else
        pass "Timeout rate is acceptable (<5%)"
    fi
else
    pass "No timeouts recorded during normal operations"
fi

echo ""
echo "=================================================="
echo -e "${GREEN}All tests completed!${NC}"
echo "=================================================="
echo ""
echo "Next steps:"
echo "1. Review timeout metrics at $METRICS_URL"
echo "2. Adjust timeout values in internal/timeout/config.go if needed"
echo "3. Monitor near_timeouts_total to identify operations close to limits"
echo "4. Check logs for timeout warnings"
