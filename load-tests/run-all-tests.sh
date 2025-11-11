#!/usr/bin/env bash

###############################################################################
# Run All Load Tests - Listings Microservice
#
# Orchestrates execution of both HTTP and gRPC load tests with proper
# monitoring, result collection, and reporting.
#
# Usage:
#   ./run-all-tests.sh [OPTIONS]
#
# Options:
#   --http-only      Run only HTTP tests
#   --grpc-only      Run only gRPC tests
#   --skip-checks    Skip pre-flight checks
#   --no-monitor     Don't collect system metrics during tests
#   --help           Show this help message
###############################################################################

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="${RESULTS_DIR:-$SCRIPT_DIR/results}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Service endpoints
HTTP_BASE_URL="${HTTP_BASE_URL:-http://localhost:8086}"
GRPC_HOST="${GRPC_HOST:-localhost:50051}"

# Test options
RUN_HTTP=true
RUN_GRPC=true
SKIP_CHECKS=false
MONITOR_SYSTEM=true

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

###############################################################################
# Logging Functions
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

log_section() {
    echo -e "\n${MAGENTA}===========================================================${NC}"
    echo -e "${MAGENTA}$1${NC}"
    echo -e "${MAGENTA}===========================================================${NC}\n"
}

###############################################################################
# Command Line Parsing
###############################################################################

show_help() {
    cat << EOF
Load Testing Suite for Listings Microservice

Usage: $0 [OPTIONS]

Options:
    --http-only      Run only HTTP tests
    --grpc-only      Run only gRPC tests
    --skip-checks    Skip pre-flight checks
    --no-monitor     Don't collect system metrics during tests
    --help           Show this help message

Environment Variables:
    HTTP_BASE_URL    HTTP endpoint (default: http://localhost:8086)
    GRPC_HOST        gRPC endpoint (default: localhost:50051)
    RESULTS_DIR      Results directory (default: ./results)

Examples:
    # Run all tests
    $0

    # Run only HTTP tests
    $0 --http-only

    # Run without monitoring
    $0 --no-monitor

EOF
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --http-only)
                RUN_GRPC=false
                shift
                ;;
            --grpc-only)
                RUN_HTTP=false
                shift
                ;;
            --skip-checks)
                SKIP_CHECKS=true
                shift
                ;;
            --no-monitor)
                MONITOR_SYSTEM=false
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

###############################################################################
# Pre-flight Checks
###############################################################################

check_dependencies() {
    log_section "Checking Dependencies"

    local missing_deps=()

    # Check k6 for HTTP tests
    if $RUN_HTTP && ! command -v k6 &> /dev/null; then
        log_warning "k6 not found (required for HTTP tests)"
        missing_deps+=("k6")
    else
        log_success "k6 found: $(k6 version 2>&1 | head -n1 || echo 'installed')"
    fi

    # Check ghz for gRPC tests
    if $RUN_GRPC && ! command -v ghz &> /dev/null; then
        log_warning "ghz not found (required for gRPC tests)"
        missing_deps+=("ghz")
    else
        log_success "ghz found: $(ghz --version 2>&1 || echo 'installed')"
    fi

    # Check jq for result parsing
    if ! command -v jq &> /dev/null; then
        log_warning "jq not found (required for result analysis)"
        missing_deps+=("jq")
    else
        log_success "jq found: $(jq --version)"
    fi

    # Check bc for calculations
    if ! command -v bc &> /dev/null; then
        log_warning "bc not found (used for calculations)"
        missing_deps+=("bc")
    fi

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        echo ""
        log_error "Missing dependencies: ${missing_deps[*]}"
        echo ""
        echo "Installation instructions:"
        echo ""
        if [[ " ${missing_deps[*]} " =~ " k6 " ]]; then
            echo "  k6: https://k6.io/docs/getting-started/installation/"
            echo "      snap install k6"
            echo ""
        fi
        if [[ " ${missing_deps[*]} " =~ " ghz " ]]; then
            echo "  ghz: go install github.com/bojand/ghz/cmd/ghz@latest"
            echo ""
        fi
        if [[ " ${missing_deps[*]} " =~ " jq " ]]; then
            echo "  jq: sudo apt-get install jq"
            echo ""
        fi
        if [[ " ${missing_deps[*]} " =~ " bc " ]]; then
            echo "  bc: sudo apt-get install bc"
            echo ""
        fi
        exit 1
    fi

    log_success "All dependencies satisfied"
}

check_services() {
    log_section "Checking Service Availability"

    local services_ok=true

    # Check HTTP service
    if $RUN_HTTP; then
        log_info "Checking HTTP service at $HTTP_BASE_URL..."
        if curl -sf "$HTTP_BASE_URL/health" > /dev/null 2>&1; then
            log_success "HTTP service is available"
        else
            log_error "HTTP service is not available at $HTTP_BASE_URL"
            services_ok=false
        fi
    fi

    # Check gRPC service
    if $RUN_GRPC; then
        log_info "Checking gRPC service at $GRPC_HOST..."
        local host="${GRPC_HOST%%:*}"
        local port="${GRPC_HOST##*:}"
        if nc -z "$host" "$port" 2>/dev/null; then
            log_success "gRPC service is available"
        else
            log_error "gRPC service is not available at $GRPC_HOST"
            services_ok=false
        fi
    fi

    if ! $services_ok; then
        echo ""
        log_error "One or more services are unavailable"
        log_info "Start the service with: cd /p/github.com/sveturs/listings && make run"
        exit 1
    fi
}

###############################################################################
# System Monitoring
###############################################################################

start_monitoring() {
    if ! $MONITOR_SYSTEM; then
        return
    fi

    log_info "Starting system monitoring..."

    local monitor_file="$RESULTS_DIR/system_metrics_${TIMESTAMP}.log"

    # Start background monitoring
    (
        echo "timestamp,cpu_percent,mem_percent,mem_used_mb,load_avg"
        while true; do
            local timestamp=$(date +%s)
            local cpu=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}')
            local mem=$(free | grep Mem | awk '{printf "%.1f,%.0f", ($3/$2) * 100, $3/1024}')
            local load=$(uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | tr -d ',')

            echo "$timestamp,$cpu,$mem,$load"
            sleep 5
        done
    ) > "$monitor_file" 2>&1 &

    MONITOR_PID=$!
    log_success "Monitoring started (PID: $MONITOR_PID)"
}

stop_monitoring() {
    if ! $MONITOR_SYSTEM || [[ -z "${MONITOR_PID:-}" ]]; then
        return
    fi

    log_info "Stopping system monitoring..."
    kill "$MONITOR_PID" 2>/dev/null || true
    wait "$MONITOR_PID" 2>/dev/null || true
    log_success "Monitoring stopped"
}

###############################################################################
# Test Execution
###############################################################################

run_http_tests() {
    log_section "Running HTTP Load Tests"

    local output_file="$RESULTS_DIR/k6_results_${TIMESTAMP}.json"

    log_info "Starting k6 HTTP load test..."
    log_info "Duration: ~5 minutes"
    log_info "Max RPS: 200"
    echo ""

    if k6 run \
        --out json="$output_file" \
        --env BASE_URL="$HTTP_BASE_URL" \
        "$SCRIPT_DIR/k6-http.js"; then
        log_success "HTTP load test completed successfully"
        return 0
    else
        log_error "HTTP load test failed"
        return 1
    fi
}

run_grpc_tests() {
    log_section "Running gRPC Load Tests"

    log_info "Starting ghz gRPC load tests..."
    log_info "Duration: ~6 minutes (4 scenarios)"
    echo ""

    if GRPC_HOST="$GRPC_HOST" \
       RESULTS_DIR="$RESULTS_DIR" \
       "$SCRIPT_DIR/ghz-grpc.sh"; then
        log_success "gRPC load tests completed successfully"
        return 0
    else
        log_error "gRPC load tests failed"
        return 1
    fi
}

###############################################################################
# Result Analysis
###############################################################################

generate_summary() {
    log_section "Generating Test Summary"

    local summary_file="$RESULTS_DIR/summary_${TIMESTAMP}.txt"

    {
        echo "=========================================="
        echo "Load Test Summary"
        echo "=========================================="
        echo ""
        echo "Timestamp: $(date -d @${TIMESTAMP:0:8} +%Y-%m-%d) ${TIMESTAMP:9:2}:${TIMESTAMP:11:2}:${TIMESTAMP:13:2}"
        echo "Results Directory: $RESULTS_DIR"
        echo ""

        if $RUN_HTTP; then
            echo "HTTP Tests:"
            echo "  Endpoint: $HTTP_BASE_URL"
            echo "  Duration: ~5 minutes"
            echo "  Max RPS: 200"
            echo ""
        fi

        if $RUN_GRPC; then
            echo "gRPC Tests:"
            echo "  Endpoint: $GRPC_HOST"
            echo "  Duration: ~6 minutes"
            echo "  Scenarios: 4"
            echo ""
        fi

        echo "Files Generated:"
        ls -lh "$RESULTS_DIR"/*_${TIMESTAMP}.* 2>/dev/null || echo "  (no files found)"
        echo ""

        echo "=========================================="
        echo "Next Steps"
        echo "=========================================="
        echo ""
        echo "1. View HTTP results:"
        echo "   cat $RESULTS_DIR/k6_results_${TIMESTAMP}.json | jq"
        echo ""
        echo "2. View gRPC results:"
        echo "   cat $RESULTS_DIR/get_all_categories_${TIMESTAMP}.json | jq"
        echo ""
        echo "3. View system metrics:"
        echo "   cat $RESULTS_DIR/system_metrics_${TIMESTAMP}.log"
        echo ""
        echo "4. Compare with previous runs:"
        echo "   ls -lh $RESULTS_DIR/"
        echo ""

    } | tee "$summary_file"

    log_success "Summary saved to: $summary_file"
}

###############################################################################
# Cleanup
###############################################################################

cleanup() {
    local exit_code=$?

    echo ""
    log_info "Cleaning up..."

    stop_monitoring

    if [[ $exit_code -ne 0 ]]; then
        log_error "Tests failed with exit code: $exit_code"
    fi

    exit $exit_code
}

###############################################################################
# Main Execution
###############################################################################

main() {
    # Parse command line arguments
    parse_args "$@"

    # Setup cleanup trap
    trap cleanup EXIT INT TERM

    # Create results directory
    mkdir -p "$RESULTS_DIR"

    # Show configuration
    log_section "Load Testing Suite - Listings Microservice"

    echo "Configuration:"
    echo "  HTTP Endpoint: $HTTP_BASE_URL"
    echo "  gRPC Endpoint: $GRPC_HOST"
    echo "  Results Dir:   $RESULTS_DIR"
    echo "  Run HTTP:      $RUN_HTTP"
    echo "  Run gRPC:      $RUN_GRPC"
    echo "  Monitor:       $MONITOR_SYSTEM"
    echo ""

    # Pre-flight checks
    if ! $SKIP_CHECKS; then
        check_dependencies
        check_services
    else
        log_warning "Skipping pre-flight checks"
    fi

    # Start monitoring
    start_monitoring

    # Run tests
    local test_failed=false

    if $RUN_HTTP; then
        if ! run_http_tests; then
            test_failed=true
        fi
    fi

    if $RUN_GRPC; then
        if ! run_grpc_tests; then
            test_failed=true
        fi
    fi

    # Stop monitoring
    stop_monitoring

    # Generate summary
    generate_summary

    # Final status
    log_section "Test Execution Complete"

    if $test_failed; then
        log_error "Some tests failed"
        return 1
    else
        log_success "All tests passed!"
        return 0
    fi
}

# Run main function
main "$@"
