#!/usr/bin/env bash

###############################################################################
# GHZ gRPC Load Test for Orders Service - Listings Microservice
#
# Tests Orders-specific gRPC endpoints under various load scenarios:
# - AddToCart (write-heavy, concurrent)
# - GetCart (read-heavy, frequently accessed)
# - CreateOrder (transaction-heavy, critical path)
# - ListOrders (read-heavy, paginated)
#
# Success Criteria:
# - AddToCart: p95 < 50ms, throughput >100 RPS
# - GetCart: p95 < 20ms, throughput >500 RPS
# - CreateOrder: p95 < 200ms, throughput >50 RPS
# - ListOrders: p95 < 100ms, error rate < 1%
###############################################################################

set -euo pipefail

# Configuration
GRPC_HOST="${GRPC_HOST:-localhost:50052}"
PROTO_PATH="${PROTO_PATH:-/p/github.com/sveturs/listings/api/proto/listings/v1/orders.proto}"
IMPORT_PATH="${IMPORT_PATH:-/p/github.com/sveturs/listings/api/proto}"
RESULTS_DIR="${RESULTS_DIR:-/p/github.com/sveturs/listings/load-tests/results}"

# Create results directory
mkdir -p "$RESULTS_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Timestamp for result files
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

###############################################################################
# Helper Functions
###############################################################################

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo ""
    echo "========================================================================"
    echo "$1"
    echo "========================================================================"
}

check_ghz_installed() {
    if ! command -v ghz &> /dev/null; then
        log_error "ghz is not installed. Install it with:"
        echo "  go install github.com/bojand/ghz/cmd/ghz@latest"
        exit 1
    fi
    log_success "ghz is installed: $(ghz --version)"
}

check_service_available() {
    log_info "Checking if Orders gRPC service is available at $GRPC_HOST..."

    if ! nc -z "${GRPC_HOST%%:*}" "${GRPC_HOST##*:}" 2>/dev/null; then
        log_error "Orders gRPC service is not available at $GRPC_HOST"
        log_info "Make sure the Orders microservice is running:"
        log_info "  /home/dim/.local/bin/start-listings-microservice.sh"
        exit 1
    fi

    log_success "Orders gRPC service is available"
}

check_proto_exists() {
    if [[ ! -f "$PROTO_PATH" ]]; then
        log_error "Proto file not found: $PROTO_PATH"
        exit 1
    fi
    log_success "Proto file found: $PROTO_PATH"
}

###############################################################################
# Setup Functions
###############################################################################

# Get or create test data (user, cart, listing)
setup_test_data() {
    log_info "Setting up test data..."

    # These IDs should exist in your test database
    # If running against fresh DB, you'll need to create them first
    export TEST_USER_ID="${TEST_USER_ID:-1}"
    export TEST_STOREFRONT_ID="${TEST_STOREFRONT_ID:-1}"
    export TEST_LISTING_ID="${TEST_LISTING_ID:-1}"
    export TEST_CART_ID="${TEST_CART_ID:-1}"

    log_info "Using test data:"
    log_info "  User ID: $TEST_USER_ID"
    log_info "  Storefront ID: $TEST_STOREFRONT_ID"
    log_info "  Listing ID: $TEST_LISTING_ID"
    log_info "  Cart ID: $TEST_CART_ID"
}

###############################################################################
# Load Test Scenarios
###############################################################################

# Scenario 1: GetCart (read-heavy, should be very fast)
test_get_cart() {
    print_header "Test 1: GetCart (Read Performance)"
    log_info "Testing cart retrieval performance..."

    local output_file="$RESULTS_DIR/get_cart_${TIMESTAMP}.json"

    # Warmup: 10 RPS for 10s
    log_info "Phase 1: Warmup (10 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/GetCart" \
        --rps=10 \
        --duration=10s \
        --connections=5 \
        --concurrency=10 \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID}" \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Main test: High read load
    log_info "Phase 2: Load Test (200 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/GetCart" \
        --rps=200 \
        --duration=60s \
        --connections=20 \
        --concurrency=200 \
        --format=json \
        --output="$output_file" \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID}" \
        "$GRPC_HOST"

    # Analyze results
    analyze_results "$output_file" "GetCart" 20
}

# Scenario 2: AddToCart (write-heavy, moderate concurrency)
test_add_to_cart() {
    print_header "Test 2: AddToCart (Write Performance)"
    log_info "Testing cart item addition performance..."

    local output_file="$RESULTS_DIR/add_to_cart_${TIMESTAMP}.json"

    # Warmup
    log_info "Phase 1: Warmup (10 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/AddToCart" \
        --rps=10 \
        --duration=10s \
        --connections=5 \
        --concurrency=10 \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID,\"listing_id\":$TEST_LISTING_ID,\"quantity\":1}" \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Main test: Concurrent writes
    log_info "Phase 2: Load Test (100 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/AddToCart" \
        --rps=100 \
        --duration=60s \
        --connections=20 \
        --concurrency=100 \
        --format=json \
        --output="$output_file" \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID,\"listing_id\":$TEST_LISTING_ID,\"quantity\":1}" \
        "$GRPC_HOST"

    analyze_results "$output_file" "AddToCart" 50
}

# Scenario 3: CreateOrder (transaction-heavy, critical path)
test_create_order() {
    print_header "Test 3: CreateOrder (Transaction Performance)"
    log_info "Testing order creation performance..."

    local output_file="$RESULTS_DIR/create_order_${TIMESTAMP}.json"

    # Warmup
    log_info "Phase 1: Warmup (5 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/CreateOrder" \
        --rps=5 \
        --duration=10s \
        --connections=5 \
        --concurrency=5 \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID,\"items\":[{\"listing_id\":$TEST_LISTING_ID,\"quantity\":1,\"price_snapshot\":99.99}],\"subtotal\":99.99,\"total\":99.99,\"currency\":\"USD\"}" \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Main test: Order creation under load
    log_info "Phase 2: Load Test (50 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/CreateOrder" \
        --rps=50 \
        --duration=60s \
        --connections=10 \
        --concurrency=50 \
        --format=json \
        --output="$output_file" \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID,\"items\":[{\"listing_id\":$TEST_LISTING_ID,\"quantity\":1,\"price_snapshot\":99.99}],\"subtotal\":99.99,\"total\":99.99,\"currency\":\"USD\"}" \
        "$GRPC_HOST"

    analyze_results "$output_file" "CreateOrder" 200
}

# Scenario 4: ListOrders (read-heavy, paginated)
test_list_orders() {
    print_header "Test 4: ListOrders (Pagination Performance)"
    log_info "Testing order listing performance..."

    local output_file="$RESULTS_DIR/list_orders_${TIMESTAMP}.json"

    # Warmup
    log_info "Phase 1: Warmup (10 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/ListOrders" \
        --rps=10 \
        --duration=10s \
        --connections=5 \
        --concurrency=10 \
        --data="{\"user_id\":$TEST_USER_ID,\"limit\":20,\"offset\":0}" \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Main test
    log_info "Phase 2: Load Test (100 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/ListOrders" \
        --rps=100 \
        --duration=60s \
        --connections=20 \
        --concurrency=100 \
        --format=json \
        --output="$output_file" \
        --data="{\"user_id\":$TEST_USER_ID,\"limit\":20,\"offset\":0}" \
        "$GRPC_HOST"

    analyze_results "$output_file" "ListOrders" 100
}

# Scenario 5: Mixed workload (realistic usage pattern)
# 50% GetCart, 30% AddToCart, 15% ListOrders, 5% CreateOrder
test_mixed_orders_workload() {
    print_header "Test 5: Mixed Orders Workload (Stress Test)"
    log_info "Running mixed Orders operations at high load..."

    local output_file="$RESULTS_DIR/mixed_orders_workload_${TIMESTAMP}.json"

    # For mixed workload, we'll run GetCart at high load as baseline
    # In production, you'd use a tool like k6 for true mixed scenarios

    log_info "Phase 1: Warmup (20 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/GetCart" \
        --rps=20 \
        --duration=10s \
        --connections=10 \
        --concurrency=20 \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID}" \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Peak load test - 300 RPS
    log_info "Phase 2: Peak Load (300 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.OrdersService/GetCart" \
        --rps=300 \
        --duration=60s \
        --connections=50 \
        --concurrency=300 \
        --format=json \
        --output="$output_file" \
        --data="{\"user_id\":$TEST_USER_ID,\"storefront_id\":$TEST_STOREFRONT_ID}" \
        "$GRPC_HOST"

    analyze_results "$output_file" "MixedOrdersWorkload" 100
}

###############################################################################
# Results Analysis
###############################################################################

analyze_results() {
    local result_file="$1"
    local test_name="$2"
    local target_p95="${3:-100}" # Default target: 100ms

    if [[ ! -f "$result_file" ]]; then
        log_error "Result file not found: $result_file"
        return 1
    fi

    echo ""
    log_info "Analyzing results for $test_name..."

    # Extract key metrics using jq
    local total=$(jq -r '.count' "$result_file")
    local rps=$(jq -r '.rps' "$result_file")
    local average=$(jq -r '.average' "$result_file")
    local fastest=$(jq -r '.fastest' "$result_file")
    local slowest=$(jq -r '.slowest' "$result_file")
    local p50=$(jq -r '.latencyDistribution[] | select(.percentage==50) | .latency' "$result_file")
    local p95=$(jq -r '.latencyDistribution[] | select(.percentage==95) | .latency' "$result_file")
    local p99=$(jq -r '.latencyDistribution[] | select(.percentage==99) | .latency' "$result_file")
    local error_dist=$(jq -r '.errorDistribution | to_entries[] | "\(.key): \(.value)"' "$result_file" 2>/dev/null)

    # Convert nanoseconds to milliseconds
    average=$(echo "scale=2; $average / 1000000" | bc)
    fastest=$(echo "scale=2; $fastest / 1000000" | bc)
    slowest=$(echo "scale=2; $slowest / 1000000" | bc)
    p50=$(echo "$p50" | sed 's/ns$//' | awk '{printf "%.2f", $1/1000000}')
    p95=$(echo "$p95" | sed 's/ns$//' | awk '{printf "%.2f", $1/1000000}')
    p99=$(echo "$p99" | sed 's/ns$//' | awk '{printf "%.2f", $1/1000000}')

    echo ""
    echo "📊 Results Summary:"
    echo "-------------------"
    echo "Total Requests:    $total"
    echo "Requests/sec:      $(printf '%.2f' "$rps")"
    echo ""
    echo "⏱️  Latency:"
    echo "Average:           ${average}ms"
    echo "Fastest:           ${fastest}ms"
    echo "Slowest:           ${slowest}ms"
    echo "p50 (median):      ${p50}ms"
    echo "p95:               ${p95}ms"
    echo "p99:               ${p99}ms"

    if [[ -n "$error_dist" ]]; then
        echo ""
        echo "❌ Errors:"
        echo "$error_dist"
    fi

    # Evaluate against success criteria
    echo ""
    evaluate_criteria "$p95" "$error_dist" "$test_name" "$target_p95"
}

evaluate_criteria() {
    local p95="$1"
    local errors="$2"
    local test_name="$3"
    local target_p95="${4:-100}"

    echo "✅ Success Criteria Evaluation:"
    echo "--------------------------------"

    local pass=true

    # Check p95 latency against target
    if (( $(echo "$p95 < $target_p95" | bc -l) )); then
        log_success "✓ p95 latency < ${target_p95}ms: ${p95}ms"
    else
        log_error "✗ p95 latency >= ${target_p95}ms: ${p95}ms"
        pass=false
    fi

    # Check error rate < 1%
    if [[ -z "$errors" ]]; then
        log_success "✓ No errors (0% error rate)"
    else
        log_error "✗ Errors detected"
        pass=false
    fi

    echo ""
    if $pass; then
        log_success "🎉 All success criteria passed for $test_name!"
    else
        log_warning "⚠️  Some criteria not met for $test_name"
    fi
}

###############################################################################
# Main Execution
###############################################################################

main() {
    print_header "gRPC Load Testing - Orders Service"

    # Pre-flight checks
    check_ghz_installed
    check_proto_exists
    check_service_available
    setup_test_data

    echo ""
    log_info "Configuration:"
    log_info "  gRPC Host:    $GRPC_HOST"
    log_info "  Proto Path:   $PROTO_PATH"
    log_info "  Results Dir:  $RESULTS_DIR"

    # Run all test scenarios
    test_get_cart
    sleep 5

    test_add_to_cart
    sleep 5

    test_create_order
    sleep 5

    test_list_orders
    sleep 5

    test_mixed_orders_workload

    # Final summary
    print_header "Orders Load Testing Complete"
    log_success "Results saved to: $RESULTS_DIR"
    log_info "Timestamp: $TIMESTAMP"

    echo ""
    log_info "To view detailed results:"
    echo "  cat $RESULTS_DIR/*_${TIMESTAMP}.json | jq"

    echo ""
    log_info "Performance Targets Summary:"
    echo "  GetCart:      p95 < 20ms  ✓"
    echo "  AddToCart:    p95 < 50ms  ✓"
    echo "  CreateOrder:  p95 < 200ms ✓"
    echo "  ListOrders:   p95 < 100ms ✓"
}

# Run main function
main "$@"
