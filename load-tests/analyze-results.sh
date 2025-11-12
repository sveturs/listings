#!/usr/bin/env bash

###############################################################################
# Result Analysis Script - Load Test Results
#
# Analyzes load test results and generates comparison reports.
#
# Usage:
#   ./analyze-results.sh [TIMESTAMP]
#
# If TIMESTAMP is not provided, analyzes the most recent test run.
###############################################################################

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

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

get_latest_timestamp() {
    ls -t "$RESULTS_DIR"/summary_*.txt 2>/dev/null | head -n1 | sed 's/.*summary_\(.*\)\.txt/\1/' || echo ""
}

###############################################################################
# Analysis Functions
###############################################################################

analyze_http_results() {
    local timestamp="$1"
    local file="$RESULTS_DIR/k6_results_${timestamp}.json"

    if [[ ! -f "$file" ]]; then
        log_warning "HTTP results not found for timestamp: $timestamp"
        return
    fi

    echo ""
    echo "=========================================="
    echo "HTTP Load Test Results (k6)"
    echo "=========================================="

    # Extract key metrics
    local total_requests=$(jq -r '.metrics.http_reqs.values.count // 0' "$file")
    local failed_requests=$(jq -r '.metrics.http_req_failed.values.passes // 0' "$file")
    local avg_duration=$(jq -r '.metrics.http_req_duration.values.avg // 0' "$file")
    local p95_duration=$(jq -r '.metrics.http_req_duration.values["p(95)"] // 0' "$file")
    local p99_duration=$(jq -r '.metrics.http_req_duration.values["p(99)"] // 0' "$file")
    local rps=$(jq -r '.metrics.http_reqs.values.rate // 0' "$file")

    # Calculate error rate
    local error_rate=0
    if [[ "$total_requests" -gt 0 ]]; then
        error_rate=$(echo "scale=2; ($failed_requests / $total_requests) * 100" | bc)
    fi

    echo ""
    echo "📊 Request Statistics:"
    echo "  Total Requests:    $total_requests"
    echo "  Failed Requests:   $failed_requests"
    echo "  Error Rate:        ${error_rate}%"
    echo "  Requests/sec:      $(printf '%.2f' "$rps")"

    echo ""
    echo "⏱️  Response Time:"
    echo "  Average:           $(printf '%.2f' "$avg_duration")ms"
    echo "  p95:               $(printf '%.2f' "$p95_duration")ms"
    echo "  p99:               $(printf '%.2f' "$p99_duration")ms"

    # Evaluate success criteria
    echo ""
    echo "✅ Success Criteria:"

    if (( $(echo "$p95_duration < 100" | bc -l) )); then
        log_success "  ✓ p95 latency < 100ms: $(printf '%.2f' "$p95_duration")ms"
    else
        log_error "  ✗ p95 latency >= 100ms: $(printf '%.2f' "$p95_duration")ms"
    fi

    if (( $(echo "$error_rate < 1" | bc -l) )); then
        log_success "  ✓ Error rate < 1%: ${error_rate}%"
    else
        log_error "  ✗ Error rate >= 1%: ${error_rate}%"
    fi

    if (( $(echo "$rps >= 100" | bc -l) )); then
        log_success "  ✓ Throughput >= 100 RPS: $(printf '%.2f' "$rps") RPS"
    else
        log_warning "  ⚠ Throughput < 100 RPS: $(printf '%.2f' "$rps") RPS"
    fi

    # Group statistics
    echo ""
    echo "📋 Endpoint Breakdown:"

    jq -r '.root_group.groups | to_entries[] | "\(.key): \(.value.checks | values | length) checks"' "$file" 2>/dev/null || echo "  No group data available"
}

analyze_grpc_results() {
    local timestamp="$1"

    echo ""
    echo "=========================================="
    echo "gRPC Load Test Results (ghz)"
    echo "=========================================="

    # Find all gRPC result files for this timestamp
    local grpc_files=("$RESULTS_DIR"/*_${timestamp}.json)

    if [[ ${#grpc_files[@]} -eq 0 ]]; then
        log_warning "No gRPC results found for timestamp: $timestamp"
        return
    fi

    for file in "$RESULTS_DIR"/*_${timestamp}.json; do
        # Skip k6 results
        if [[ "$file" == *"k6_results"* || "$file" == *"system_metrics"* ]]; then
            continue
        fi

        if [[ ! -f "$file" ]]; then
            continue
        fi

        local test_name=$(basename "$file" | sed "s/_${timestamp}.json//" | tr '_' ' ' | sed 's/\b\(.\)/\u\1/g')

        echo ""
        echo "--- $test_name ---"

        # Extract metrics
        local total=$(jq -r '.count // 0' "$file")
        local rps=$(jq -r '.rps // 0' "$file")
        local average=$(jq -r '.average // 0' "$file")
        local p95=$(jq -r '.latencyDistribution[] | select(.percentage==95) | .latency' "$file" 2>/dev/null || echo "0")
        local p99=$(jq -r '.latencyDistribution[] | select(.percentage==99) | .latency' "$file" 2>/dev/null || echo "0")
        local errors=$(jq -r '.errorDistribution | length' "$file" 2>/dev/null || echo "0")

        # Convert nanoseconds to milliseconds
        average=$(echo "scale=2; $average / 1000000" | bc 2>/dev/null || echo "0")
        p95_ms=$(echo "$p95" | sed 's/ns$//' | awk '{printf "%.2f", $1/1000000}' 2>/dev/null || echo "0")
        p99_ms=$(echo "$p99" | sed 's/ns$//' | awk '{printf "%.2f", $1/1000000}' 2>/dev/null || echo "0")

        echo "  Total Requests:    $total"
        echo "  Requests/sec:      $(printf '%.2f' "$rps")"
        echo "  Avg Latency:       ${average}ms"
        echo "  p95 Latency:       ${p95_ms}ms"
        echo "  p99 Latency:       ${p99_ms}ms"
        echo "  Errors:            $errors"

        # Evaluate
        if (( $(echo "$p95_ms < 100" | bc -l) )); then
            log_success "  ✓ Passed (p95 < 100ms)"
        else
            log_error "  ✗ Failed (p95 >= 100ms)"
        fi
    done
}

analyze_system_metrics() {
    local timestamp="$1"
    local file="$RESULTS_DIR/system_metrics_${timestamp}.log"

    if [[ ! -f "$file" ]]; then
        log_warning "System metrics not found for timestamp: $timestamp"
        return
    fi

    echo ""
    echo "=========================================="
    echo "System Metrics"
    echo "=========================================="

    # Skip header and calculate statistics
    local cpu_avg=$(tail -n +2 "$file" | awk -F, '{sum+=$2; count++} END {printf "%.2f", sum/count}')
    local cpu_max=$(tail -n +2 "$file" | awk -F, '{max=0} {if($2>max) max=$2} END {printf "%.2f", max}')
    local mem_avg=$(tail -n +2 "$file" | awk -F, '{sum+=$3; count++} END {printf "%.2f", sum/count}')
    local mem_max=$(tail -n +2 "$file" | awk -F, '{max=0} {if($3>max) max=$3} END {printf "%.2f", max}')

    echo ""
    echo "🖥️  CPU Usage:"
    echo "  Average:           ${cpu_avg}%"
    echo "  Peak:              ${cpu_max}%"

    echo ""
    echo "💾 Memory Usage:"
    echo "  Average:           ${mem_avg}%"
    echo "  Peak:              ${mem_max}%"

    # Evaluate
    echo ""
    echo "✅ Resource Utilization:"

    if (( $(echo "$cpu_max < 80" | bc -l) )); then
        log_success "  ✓ CPU usage healthy (peak: ${cpu_max}%)"
    else
        log_warning "  ⚠ High CPU usage (peak: ${cpu_max}%)"
    fi

    if (( $(echo "$mem_max < 80" | bc -l) )); then
        log_success "  ✓ Memory usage healthy (peak: ${mem_max}%)"
    else
        log_warning "  ⚠ High memory usage (peak: ${mem_max}%)"
    fi
}

compare_with_baseline() {
    local timestamp="$1"

    echo ""
    echo "=========================================="
    echo "Comparison with Baseline"
    echo "=========================================="

    # TODO: Implement baseline comparison
    # This would compare current results with a saved baseline
    # to detect performance regressions

    log_info "Baseline comparison not yet implemented"
    echo "  Run: ./analyze-results.sh --set-baseline $timestamp"
}

generate_html_report() {
    local timestamp="$1"
    local output_file="$RESULTS_DIR/report_${timestamp}.html"

    log_info "Generating HTML report..."

    # TODO: Generate HTML report with charts
    # Could use tools like plotly, chart.js, or gnuplot

    log_info "HTML report generation not yet implemented"
}

###############################################################################
# Main
###############################################################################

main() {
    local timestamp="${1:-}"

    # If no timestamp provided, use latest
    if [[ -z "$timestamp" ]]; then
        timestamp=$(get_latest_timestamp)
        if [[ -z "$timestamp" ]]; then
            log_error "No test results found in $RESULTS_DIR"
            exit 1
        fi
        log_info "Using latest test results: $timestamp"
    fi

    # Verify results exist
    if [[ ! -f "$RESULTS_DIR/summary_${timestamp}.txt" ]]; then
        log_error "No results found for timestamp: $timestamp"
        exit 1
    fi

    # Show summary
    echo "=========================================="
    echo "Load Test Results Analysis"
    echo "=========================================="
    echo ""
    echo "Timestamp: $timestamp"
    echo "Directory: $RESULTS_DIR"
    echo ""

    # Analyze each component
    analyze_http_results "$timestamp"
    analyze_grpc_results "$timestamp"
    analyze_system_metrics "$timestamp"
    compare_with_baseline "$timestamp"

    # Final summary
    echo ""
    echo "=========================================="
    echo "Analysis Complete"
    echo "=========================================="
    echo ""
    log_success "Results saved to: $RESULTS_DIR"
}

# Parse arguments
if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
    echo "Usage: $0 [TIMESTAMP]"
    echo ""
    echo "Analyzes load test results for the given timestamp."
    echo "If no timestamp provided, analyzes the most recent results."
    echo ""
    echo "Examples:"
    echo "  $0                    # Analyze latest results"
    echo "  $0 20251110_143000    # Analyze specific run"
    exit 0
fi

main "$@"
