#!/bin/bash
#
# Performance Baseline Measurement Script
# Tests critical gRPC endpoints and measures P50/P95/P99 latencies
#
# Usage:
#   ./baseline.sh [options]
#
# Options:
#   -h, --help              Show this help message
#   -d, --duration SECONDS  Duration for each test (default: 30)
#   -c, --concurrency NUM   Number of concurrent connections (default: 10)
#   -r, --rate NUM          Requests per second (default: 100)
#   -o, --output FILE       Output results to file (default: baseline_results.json)
#   --grpc-addr ADDR        gRPC server address (default: localhost:8086)
#   --metrics-url URL       Prometheus metrics URL (default: http://localhost:8086/metrics)
#

set -euo pipefail

# Default configuration
DURATION=${DURATION:-30}
CONCURRENCY=${CONCURRENCY:-10}
RATE=${RATE:-100}
OUTPUT_FILE=${OUTPUT_FILE:-baseline_results.json}
GRPC_ADDR=${GRPC_ADDR:-localhost:8086}
METRICS_URL=${METRICS_URL:-http://localhost:8086/metrics}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
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

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -d|--duration)
                DURATION="$2"
                shift 2
                ;;
            -c|--concurrency)
                CONCURRENCY="$2"
                shift 2
                ;;
            -r|--rate)
                RATE="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_FILE="$2"
                shift 2
                ;;
            --grpc-addr)
                GRPC_ADDR="$2"
                shift 2
                ;;
            --metrics-url)
                METRICS_URL="$2"
                shift 2
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

show_help() {
    head -n 20 "$0" | grep "^#" | sed 's/^# //' | sed 's/^#//'
}

# Check if required tools are installed
check_prerequisites() {
    log_info "Checking prerequisites..."

    local missing_tools=()

    if ! command -v grpcurl &> /dev/null; then
        missing_tools+=("grpcurl")
    fi

    if ! command -v ghz &> /dev/null; then
        missing_tools+=("ghz")
    fi

    if ! command -v jq &> /dev/null; then
        missing_tools+=("jq")
    fi

    if ! command -v bc &> /dev/null; then
        missing_tools+=("bc")
    fi

    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        echo ""
        echo "Install instructions:"
        echo "  - grpcurl: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
        echo "  - ghz:     go install github.com/bojand/ghz/cmd/ghz@latest"
        echo "  - jq:      apt-get install jq (Ubuntu/Debian) or brew install jq (macOS)"
        echo "  - bc:      apt-get install bc (Ubuntu/Debian) or brew install bc (macOS)"
        exit 1
    fi

    log_success "All prerequisites satisfied"
}

# Check if gRPC server is accessible
check_server() {
    log_info "Checking gRPC server connectivity at $GRPC_ADDR..."

    if ! grpcurl -plaintext "$GRPC_ADDR" list &> /dev/null; then
        log_error "Cannot connect to gRPC server at $GRPC_ADDR"
        log_error "Make sure the listings service is running"
        exit 1
    fi

    log_success "gRPC server is accessible"
}

# Fetch current Prometheus metrics
fetch_metrics() {
    log_info "Fetching current Prometheus metrics..."

    if ! curl -s "$METRICS_URL" > /dev/null; then
        log_warning "Cannot fetch metrics from $METRICS_URL"
        log_warning "Metrics collection will be limited"
        return 1
    fi

    log_success "Metrics endpoint is accessible"
    return 0
}

# Run benchmark for a specific RPC method
benchmark_rpc() {
    local method=$1
    local data=$2
    local description=$3

    log_info "Benchmarking: $method - $description"

    # Run ghz benchmark
    local result_file="/tmp/ghz_${method//\//_}_$$.json"

    ghz --insecure \
        --proto /p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto \
        --import-paths /p/github.com/sveturs/listings/api/proto \
        --call "$method" \
        --data "$data" \
        --duration "${DURATION}s" \
        --concurrency "$CONCURRENCY" \
        --rps "$RATE" \
        --connections "$CONCURRENCY" \
        --format json \
        --output "$result_file" \
        "$GRPC_ADDR" 2>&1 | grep -v "^$" || true

    if [ ! -f "$result_file" ]; then
        log_error "Benchmark failed for $method"
        return 1
    fi

    # Extract key metrics
    local count=$(jq -r '.count' "$result_file")
    local total_time=$(jq -r '.total' "$result_file")
    local average=$(jq -r '.average' "$result_file")
    local fastest=$(jq -r '.fastest' "$result_file")
    local slowest=$(jq -r '.slowest' "$result_file")
    local rps=$(jq -r '.rps' "$result_file")

    # Extract histogram (P50, P95, P99)
    local p50=$(jq -r '.histogram[] | select(.mark == 0.5) | .latency' "$result_file" 2>/dev/null || echo "0")
    local p95=$(jq -r '.histogram[] | select(.mark == 0.95) | .latency' "$result_file" 2>/dev/null || echo "0")
    local p99=$(jq -r '.histogram[] | select(.mark == 0.99) | .latency' "$result_file" 2>/dev/null || echo "0")

    # Extract error rate
    local error_dist=$(jq -r '.errorDistribution' "$result_file")
    local total_errors=0
    if [ "$error_dist" != "null" ] && [ -n "$error_dist" ]; then
        total_errors=$(echo "$error_dist" | jq -r 'to_entries | map(.value) | add // 0')
    fi

    local error_rate=0
    if [ "$count" -gt 0 ]; then
        error_rate=$(echo "scale=4; ($total_errors / $count) * 100" | bc)
    fi

    # Convert nanoseconds to milliseconds for readability
    convert_ns_to_ms() {
        local ns=$1
        echo "scale=3; $ns / 1000000" | bc
    }

    local avg_ms=$(convert_ns_to_ms "$average")
    local p50_ms=$(convert_ns_to_ms "$p50")
    local p95_ms=$(convert_ns_to_ms "$p95")
    local p99_ms=$(convert_ns_to_ms "$p99")
    local fastest_ms=$(convert_ns_to_ms "$fastest")
    local slowest_ms=$(convert_ns_to_ms "$slowest")

    # Print results
    echo ""
    echo "  Total Requests: $count"
    echo "  RPS:            $rps"
    echo "  Error Rate:     ${error_rate}%"
    echo "  Latency:"
    echo "    Average:  ${avg_ms} ms"
    echo "    Fastest:  ${fastest_ms} ms"
    echo "    Slowest:  ${slowest_ms} ms"
    echo "    P50:      ${p50_ms} ms"
    echo "    P95:      ${p95_ms} ms"
    echo "    P99:      ${p99_ms} ms"
    echo ""

    # Add to results JSON
    jq -n \
        --arg method "$method" \
        --arg description "$description" \
        --argjson count "$count" \
        --argjson rps "$rps" \
        --argjson error_rate "$error_rate" \
        --argjson average "$avg_ms" \
        --argjson p50 "$p50_ms" \
        --argjson p95 "$p95_ms" \
        --argjson p99 "$p99_ms" \
        --argjson fastest "$fastest_ms" \
        --argjson slowest "$slowest_ms" \
        '{
            method: $method,
            description: $description,
            total_requests: $count,
            rps: $rps,
            error_rate_pct: $error_rate,
            latency_ms: {
                average: $average,
                p50: $p50,
                p95: $p95,
                p99: $p99,
                fastest: $fastest,
                slowest: $slowest
            }
        }' >> /tmp/baseline_results_$$.jsonl

    # Cleanup
    rm -f "$result_file"

    log_success "Completed: $method"
}

# Main baseline test suite
run_baseline_tests() {
    log_info "Starting baseline measurement suite"
    log_info "Configuration: duration=${DURATION}s, concurrency=${CONCURRENCY}, rate=${RATE} RPS"
    echo ""

    # Initialize results file
    rm -f /tmp/baseline_results_$$.jsonl

    # Critical RPC Methods - Organized by priority

    # === CRITICAL (Stock Management - Order Processing) ===
    benchmark_rpc \
        "listings.v1.ListingsService/CheckStockAvailability" \
        '{"items":[{"product_id":1,"quantity":5}]}' \
        "Check stock availability (single item)"

    benchmark_rpc \
        "listings.v1.ListingsService/DecrementStock" \
        '{"order_id":"test-order-001","items":[{"product_id":1,"quantity":2}]}' \
        "Decrement stock (single item)"

    # === HIGH PRIORITY (Core CRUD Operations) ===
    benchmark_rpc \
        "listings.v1.ListingsService/GetProduct" \
        '{"id":1}' \
        "Get single product by ID"

    benchmark_rpc \
        "listings.v1.ListingsService/ListProducts" \
        '{"storefront_id":1,"limit":20,"offset":0}' \
        "List products (paginated)"

    benchmark_rpc \
        "listings.v1.ListingsService/GetListing" \
        '{"id":1}' \
        "Get single listing by ID"

    benchmark_rpc \
        "listings.v1.ListingsService/SearchListings" \
        '{"query":"test","limit":20,"offset":0}' \
        "Search listings (basic query)"

    # === HIGH PRIORITY (Categories) ===
    benchmark_rpc \
        "listings.v1.ListingsService/GetAllCategories" \
        '{}' \
        "Get all categories"

    benchmark_rpc \
        "listings.v1.ListingsService/GetRootCategories" \
        '{}' \
        "Get root categories"

    # === HIGH PRIORITY (Storefronts) ===
    benchmark_rpc \
        "listings.v1.ListingsService/GetStorefront" \
        '{"id":1}' \
        "Get storefront by ID"

    benchmark_rpc \
        "listings.v1.ListingsService/ListStorefronts" \
        '{"limit":20,"offset":0}' \
        "List storefronts (paginated)"

    # === MEDIUM PRIORITY (Batch Operations) ===
    benchmark_rpc \
        "listings.v1.ListingsService/GetProductsByIDs" \
        '{"ids":[1,2,3,4,5]}' \
        "Get products by IDs (batch)"

    benchmark_rpc \
        "listings.v1.ListingsService/CheckStockAvailability" \
        '{"items":[{"product_id":1,"quantity":5},{"product_id":2,"quantity":3},{"product_id":3,"quantity":10}]}' \
        "Check stock availability (multiple items)"

    # === MEDIUM PRIORITY (Inventory) ===
    benchmark_rpc \
        "listings.v1.ListingsService/GetProductStats" \
        '{"product_id":1}' \
        "Get product statistics"

    # Combine all results into single JSON
    if [ -f /tmp/baseline_results_$$.jsonl ]; then
        jq -s \
            --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
            --arg duration "$DURATION" \
            --argjson concurrency "$CONCURRENCY" \
            --argjson rate "$RATE" \
            --arg grpc_addr "$GRPC_ADDR" \
            '{
                metadata: {
                    timestamp: $timestamp,
                    grpc_address: $grpc_addr,
                    test_config: {
                        duration_seconds: ($duration | tonumber),
                        concurrency: $concurrency,
                        target_rps: $rate
                    }
                },
                results: .
            }' /tmp/baseline_results_$$.jsonl > "$OUTPUT_FILE"

        rm -f /tmp/baseline_results_$$.jsonl

        log_success "Baseline results saved to: $OUTPUT_FILE"
    else
        log_error "No results generated"
        exit 1
    fi
}

# Generate human-readable report
generate_report() {
    if [ ! -f "$OUTPUT_FILE" ]; then
        log_error "Results file not found: $OUTPUT_FILE"
        return 1
    fi

    log_info "Generating baseline report..."

    local report_file="${OUTPUT_FILE%.json}.txt"

    cat > "$report_file" <<EOF
========================================
LISTINGS MICROSERVICE - BASELINE REPORT
========================================

Test Date: $(jq -r '.metadata.timestamp' "$OUTPUT_FILE")
gRPC Server: $(jq -r '.metadata.grpc_address' "$OUTPUT_FILE")

Configuration:
  - Duration: $(jq -r '.metadata.test_config.duration_seconds' "$OUTPUT_FILE")s per test
  - Concurrency: $(jq -r '.metadata.test_config.concurrency' "$OUTPUT_FILE") connections
  - Target RPS: $(jq -r '.metadata.test_config.target_rps' "$OUTPUT_FILE")

========================================
PERFORMANCE BASELINE RESULTS
========================================

EOF

    # Add each result
    jq -r '.results[] |
        "Method: \(.method)\n" +
        "Description: \(.description)\n" +
        "Total Requests: \(.total_requests)\n" +
        "Throughput: \(.rps | tonumber | floor) RPS\n" +
        "Error Rate: \(.error_rate_pct)%\n" +
        "Latency (ms):\n" +
        "  Average: \(.latency_ms.average)\n" +
        "  P50:     \(.latency_ms.p50)\n" +
        "  P95:     \(.latency_ms.p95)\n" +
        "  P99:     \(.latency_ms.p99)\n" +
        "  Range:   \(.latency_ms.fastest) - \(.latency_ms.slowest)\n" +
        "----------------------------------------\n"
    ' "$OUTPUT_FILE" >> "$report_file"

    # Add summary statistics
    cat >> "$report_file" <<EOF

========================================
SUMMARY STATISTICS
========================================

EOF

    # Calculate aggregates
    local avg_p50=$(jq -r '[.results[].latency_ms.p50] | add / length' "$OUTPUT_FILE")
    local avg_p95=$(jq -r '[.results[].latency_ms.p95] | add / length' "$OUTPUT_FILE")
    local avg_p99=$(jq -r '[.results[].latency_ms.p99] | add / length' "$OUTPUT_FILE")
    local max_p99=$(jq -r '[.results[].latency_ms.p99] | max' "$OUTPUT_FILE")
    local avg_error_rate=$(jq -r '[.results[].error_rate_pct] | add / length' "$OUTPUT_FILE")
    local total_requests=$(jq -r '[.results[].total_requests] | add' "$OUTPUT_FILE")

    cat >> "$report_file" <<EOF
Average Latency Across All Methods:
  P50: ${avg_p50} ms
  P95: ${avg_p95} ms
  P99: ${avg_p99} ms

Worst Case P99: ${max_p99} ms
Average Error Rate: ${avg_error_rate}%
Total Requests Processed: ${total_requests}

========================================
BASELINE THRESHOLDS (for alerts)
========================================

Recommended Alert Thresholds:
  - P95 Latency Warning:  > $(echo "$avg_p95 * 1.5" | bc) ms
  - P95 Latency Critical: > $(echo "$avg_p95 * 2" | bc) ms
  - P99 Latency Warning:  > $(echo "$avg_p99 * 1.5" | bc) ms
  - P99 Latency Critical: > $(echo "$avg_p99 * 2" | bc) ms
  - Error Rate Warning:   > 0.5%
  - Error Rate Critical:  > 1%

========================================
EOF

    log_success "Report generated: $report_file"

    # Print summary to console
    echo ""
    log_info "BASELINE SUMMARY:"
    echo "  Average P50: ${avg_p50} ms"
    echo "  Average P95: ${avg_p95} ms"
    echo "  Average P99: ${avg_p99} ms"
    echo "  Max P99:     ${max_p99} ms"
    echo "  Error Rate:  ${avg_error_rate}%"
    echo ""
}

# Main execution
main() {
    parse_args "$@"

    echo "========================================"
    echo "  Listings Microservice - Baseline"
    echo "========================================"
    echo ""

    check_prerequisites
    check_server
    fetch_metrics || true

    run_baseline_tests
    generate_report

    log_success "Baseline measurement complete!"
    log_info "Results: $OUTPUT_FILE"
    log_info "Report:  ${OUTPUT_FILE%.json}.txt"
}

# Run main function
main "$@"
