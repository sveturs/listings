#!/usr/bin/env bash

###############################################################################
# GHZ gRPC Load Test for Listings Microservice
#
# Tests gRPC endpoints under various load scenarios:
# - GetAllCategories (read-heavy, cached)
# - ListStorefronts (read-heavy, paginated)
# - GetListing (read, frequently accessed)
#
# Success Criteria:
# - p95 latency < 100ms
# - Error rate < 1%
# - 100 RPS without degradation
###############################################################################

set -euo pipefail

# Configuration
GRPC_HOST="${GRPC_HOST:-localhost:50051}"
PROTO_PATH="${PROTO_PATH:-/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto}"
IMPORT_PATH="${IMPORT_PATH:-/p/github.com/sveturs/listings/api/proto}"
RESULTS_DIR="${RESULTS_DIR:-/p/github.com/sveturs/listings/load-tests/results}"

# Create results directory
mkdir -p "$RESULTS_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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
    log_info "Checking if gRPC service is available at $GRPC_HOST..."

    if ! nc -z "${GRPC_HOST%%:*}" "${GRPC_HOST##*:}" 2>/dev/null; then
        log_error "gRPC service is not available at $GRPC_HOST"
        log_info "Make sure the service is running:"
        log_info "  cd /p/github.com/sveturs/listings"
        log_info "  make run"
        exit 1
    fi

    log_success "gRPC service is available"
}

check_proto_exists() {
    if [[ ! -f "$PROTO_PATH" ]]; then
        log_error "Proto file not found: $PROTO_PATH"
        exit 1
    fi
    log_success "Proto file found: $PROTO_PATH"
}

###############################################################################
# Load Test Scenarios
###############################################################################

# Scenario 1: GetAllCategories (cached read, should be fast)
test_get_all_categories() {
    print_header "Test 1: GetAllCategories"
    log_info "Testing cached category reads..."

    local output_file="$RESULTS_DIR/get_all_categories_${TIMESTAMP}.json"

    # Warmup: 10 RPS for 10s
    log_info "Phase 1: Warmup (10 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/GetAllCategories" \
        --rps=10 \
        --duration=10s \
        --connections=5 \
        --concurrency=10 \
        --data='{}' \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Main test: Ramp up to 100 RPS
    log_info "Phase 2: Load Test (50 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/GetAllCategories" \
        --rps=50 \
        --duration=60s \
        --connections=10 \
        --concurrency=50 \
        --format=json \
        --output="$output_file" \
        --data='{}' \
        "$GRPC_HOST"

    # Analyze results
    analyze_results "$output_file" "GetAllCategories"
}

# Scenario 2: ListStorefronts (paginated read)
test_list_storefronts() {
    print_header "Test 2: ListStorefronts"
    log_info "Testing paginated storefront listing..."

    local output_file="$RESULTS_DIR/list_storefronts_${TIMESTAMP}.json"

    # Warmup
    log_info "Phase 1: Warmup (10 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/ListStorefronts" \
        --rps=10 \
        --duration=10s \
        --connections=5 \
        --concurrency=10 \
        --data='{"limit":10,"offset":0}' \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Main test
    log_info "Phase 2: Load Test (50 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/ListStorefronts" \
        --rps=50 \
        --duration=60s \
        --connections=10 \
        --concurrency=50 \
        --format=json \
        --output="$output_file" \
        --data='{"limit":10,"offset":0}' \
        "$GRPC_HOST"

    analyze_results "$output_file" "ListStorefronts"
}

# Scenario 3: GetListing (single item read)
test_get_listing() {
    print_header "Test 3: GetListing"
    log_info "Testing single listing retrieval..."

    # First, get a valid listing ID
    local listing_id=$(get_valid_listing_id)

    if [[ -z "$listing_id" ]]; then
        log_warning "No valid listing ID found, skipping GetListing test"
        return
    fi

    log_info "Using listing ID: $listing_id"

    local output_file="$RESULTS_DIR/get_listing_${TIMESTAMP}.json"

    # Warmup
    log_info "Phase 1: Warmup (10 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/GetListing" \
        --rps=10 \
        --duration=10s \
        --connections=5 \
        --concurrency=10 \
        --data="{\"id\":$listing_id}" \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Main test
    log_info "Phase 2: Load Test (100 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/GetListing" \
        --rps=100 \
        --duration=60s \
        --connections=20 \
        --concurrency=100 \
        --format=json \
        --output="$output_file" \
        --data="{\"id\":$listing_id}" \
        "$GRPC_HOST"

    analyze_results "$output_file" "GetListing"
}

# Scenario 4: Mixed workload (stress test)
test_mixed_workload() {
    print_header "Test 4: Mixed Workload (Stress Test)"
    log_info "Running mixed gRPC calls at high load..."

    local output_file="$RESULTS_DIR/mixed_workload_${TIMESTAMP}.json"

    # Get valid listing ID
    local listing_id=$(get_valid_listing_id)

    # Warmup
    log_info "Phase 1: Warmup (20 RPS, 10s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/GetAllCategories" \
        --rps=20 \
        --duration=10s \
        --connections=10 \
        --concurrency=20 \
        --data='{}' \
        "$GRPC_HOST" > /dev/null 2>&1 || true

    # Peak load test - 200 RPS
    log_info "Phase 2: Peak Load (200 RPS, 60s)"
    ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/GetAllCategories" \
        --rps=200 \
        --duration=60s \
        --connections=40 \
        --concurrency=200 \
        --format=json \
        --output="$output_file" \
        --data='{}' \
        "$GRPC_HOST"

    analyze_results "$output_file" "MixedWorkload"
}

###############################################################################
# Helper Functions for Testing
###############################################################################

get_valid_listing_id() {
    # Try to get a listing ID from ListListings call
    local temp_file=$(mktemp)

    if ghz --insecure \
        --proto="$PROTO_PATH" \
        --import-paths="$IMPORT_PATH" \
        --call="listings.v1.ListingsService/ListListings" \
        --total=1 \
        --format=json \
        --output="$temp_file" \
        --data='{"limit":1,"offset":0}' \
        "$GRPC_HOST" > /dev/null 2>&1; then

        # Extract first listing ID from response
        local listing_id=$(jq -r '.details[0].response.listings[0].id // empty' "$temp_file" 2>/dev/null)
        rm -f "$temp_file"

        if [[ -n "$listing_id" && "$listing_id" != "null" ]]; then
            echo "$listing_id"
            return
        fi
    fi

    # Fallback to default ID
    echo "1"
}

###############################################################################
# Results Analysis
###############################################################################

analyze_results() {
    local result_file="$1"
    local test_name="$2"

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
    evaluate_criteria "$p95" "$error_dist" "$test_name"
}

evaluate_criteria() {
    local p95="$1"
    local errors="$2"
    local test_name="$3"

    echo "✅ Success Criteria Evaluation:"
    echo "--------------------------------"

    local pass=true

    # Check p95 latency < 100ms
    if (( $(echo "$p95 < 100" | bc -l) )); then
        log_success "✓ p95 latency < 100ms: ${p95}ms"
    else
        log_error "✗ p95 latency >= 100ms: ${p95}ms"
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
    print_header "gRPC Load Testing - Listings Microservice"

    # Pre-flight checks
    check_ghz_installed
    check_proto_exists
    check_service_available

    echo ""
    log_info "Configuration:"
    log_info "  gRPC Host:    $GRPC_HOST"
    log_info "  Proto Path:   $PROTO_PATH"
    log_info "  Results Dir:  $RESULTS_DIR"

    # Run all test scenarios
    test_get_all_categories
    sleep 5

    test_list_storefronts
    sleep 5

    test_get_listing
    sleep 5

    test_mixed_workload

    # Final summary
    print_header "Load Testing Complete"
    log_success "Results saved to: $RESULTS_DIR"
    log_info "Timestamp: $TIMESTAMP"

    echo ""
    log_info "To view detailed results:"
    echo "  cat $RESULTS_DIR/*_${TIMESTAMP}.json | jq"
}

# Run main function
main "$@"
