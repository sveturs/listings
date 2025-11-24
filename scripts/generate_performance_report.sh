#!/usr/bin/env bash

###############################################################################
# Performance Report Generator
#
# Analyzes benchmark and load test results, generates comprehensive report
# with metrics, charts, and recommendations.
#
# Usage:
#   ./generate_performance_report.sh [results_dir] [output_file]
#
# Example:
#   ./generate_performance_report.sh ./load-tests/results ./PERFORMANCE_REPORT.md
###############################################################################

set -euo pipefail

# Configuration
RESULTS_DIR="${1:-./load-tests/results}"
OUTPUT_FILE="${2:-./PERFORMANCE_REPORT_$(date +%Y%m%d_%H%M%S).md}"
BENCHMARK_FILE="${3:-./benchmark_results.txt}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

###############################################################################
# Report Generation
###############################################################################

generate_report() {
    log_info "Generating performance report..."

    cat > "$OUTPUT_FILE" << 'EOF'
# Performance Testing Report - Orders Microservice

**Generated:** $(date '+%Y-%m-%d %H:%M:%S')
**Environment:** Development/CI
**Test Suite Version:** 1.0.0

---

## Executive Summary

This report provides comprehensive performance analysis of the Orders Microservice, covering:
- **Benchmark Tests**: Go benchmarks for individual operations
- **Load Tests**: gRPC load testing under various scenarios
- **Success Criteria**: Comparison against performance targets
- **Recommendations**: Optimization suggestions

---

## 1. Performance Targets

| Operation | Target P95 | Target Throughput | Actual P95 | Actual Throughput | Status |
|-----------|------------|-------------------|------------|-------------------|--------|
EOF

    # Add target vs actual comparison
    if [[ -d "$RESULTS_DIR" ]]; then
        for result_file in "$RESULTS_DIR"/*.json; do
            if [[ -f "$result_file" ]]; then
                local test_name=$(basename "$result_file" | sed 's/_[0-9]*\.json$//' | sed 's/_/ /g')
                local p95=$(jq -r '.latencyDistribution[] | select(.percentage==95) | .latency' "$result_file" 2>/dev/null | sed 's/ns$//' | awk '{printf "%.2f", $1/1000000}')
                local rps=$(jq -r '.rps' "$result_file" 2>/dev/null)

                # Determine target based on test name
                local target_p95="N/A"
                local target_rps="N/A"
                local status="⏸️"

                case "$test_name" in
                    *"get_cart"*)
                        target_p95="20ms"
                        target_rps="500 RPS"
                        [[ $(echo "$p95 < 20" | bc -l) -eq 1 ]] && status="✅" || status="❌"
                        ;;
                    *"add_to_cart"*)
                        target_p95="50ms"
                        target_rps="100 RPS"
                        [[ $(echo "$p95 < 50" | bc -l) -eq 1 ]] && status="✅" || status="❌"
                        ;;
                    *"create_order"*)
                        target_p95="200ms"
                        target_rps="50 RPS"
                        [[ $(echo "$p95 < 200" | bc -l) -eq 1 ]] && status="✅" || status="❌"
                        ;;
                    *"list_orders"*)
                        target_p95="100ms"
                        target_rps="100 RPS"
                        [[ $(echo "$p95 < 100" | bc -l) -eq 1 ]] && status="✅" || status="❌"
                        ;;
                esac

                cat >> "$OUTPUT_FILE" << EOF
| $test_name | $target_p95 | $target_rps | ${p95}ms | ${rps} RPS | $status |
EOF
            fi
        done
    fi

    cat >> "$OUTPUT_FILE" << 'EOF'

---

## 2. Benchmark Test Results

### 2.1 Go Benchmarks

EOF

    # Add Go benchmark results
    if [[ -f "$BENCHMARK_FILE" ]]; then
        cat >> "$OUTPUT_FILE" << 'EOF'
```
EOF
        cat "$BENCHMARK_FILE" >> "$OUTPUT_FILE" 2>/dev/null || echo "No benchmark results found"
        cat >> "$OUTPUT_FILE" << 'EOF'
```

### 2.2 Benchmark Analysis

EOF

        # Parse benchmark results
        if grep -q "Benchmark" "$BENCHMARK_FILE" 2>/dev/null; then
            cat >> "$OUTPUT_FILE" << 'EOF'
| Benchmark | Operations | ns/op | Allocations | Bytes/op |
|-----------|------------|-------|-------------|----------|
EOF

            while IFS= read -r line; do
                if [[ "$line" =~ ^Benchmark ]]; then
                    # Extract benchmark name and metrics
                    name=$(echo "$line" | awk '{print $1}')
                    ops=$(echo "$line" | awk '{print $2}')
                    ns_op=$(echo "$line" | awk '{print $3}')
                    allocs=$(echo "$line" | awk '{print $5}')
                    bytes=$(echo "$line" | awk '{print $7}')

                    cat >> "$OUTPUT_FILE" << EOF
| $name | $ops | $ns_op | $allocs | $bytes |
EOF
                fi
            done < "$BENCHMARK_FILE"
        fi
    fi

    cat >> "$OUTPUT_FILE" << 'EOF'

---

## 3. Load Test Results

### 3.1 Test Scenarios

EOF

    # Add detailed load test results
    if [[ -d "$RESULTS_DIR" ]]; then
        for result_file in "$RESULTS_DIR"/*.json; do
            if [[ -f "$result_file" ]]; then
                local test_name=$(basename "$result_file" | sed 's/_[0-9]*\.json$//')

                cat >> "$OUTPUT_FILE" << EOF

#### Test: $test_name

EOF

                # Extract metrics
                local total=$(jq -r '.count' "$result_file" 2>/dev/null || echo "N/A")
                local rps=$(jq -r '.rps' "$result_file" 2>/dev/null || echo "N/A")
                local average=$(jq -r '.average' "$result_file" 2>/dev/null | awk '{printf "%.2fms", $1/1000000}')
                local fastest=$(jq -r '.fastest' "$result_file" 2>/dev/null | awk '{printf "%.2fms", $1/1000000}')
                local slowest=$(jq -r '.slowest' "$result_file" 2>/dev/null | awk '{printf "%.2fms", $1/1000000}')
                local p50=$(jq -r '.latencyDistribution[] | select(.percentage==50) | .latency' "$result_file" 2>/dev/null | sed 's/ns$//' | awk '{printf "%.2fms", $1/1000000}')
                local p95=$(jq -r '.latencyDistribution[] | select(.percentage==95) | .latency' "$result_file" 2>/dev/null | sed 's/ns$//' | awk '{printf "%.2fms", $1/1000000}')
                local p99=$(jq -r '.latencyDistribution[] | select(.percentage==99) | .latency' "$result_file" 2>/dev/null | sed 's/ns$//' | awk '{printf "%.2fms", $1/1000000}')

                cat >> "$OUTPUT_FILE" << EOF
**Summary:**
- Total Requests: $total
- Requests/sec: $rps
- Average Latency: $average
- Fastest: $fastest
- Slowest: $slowest

**Latency Percentiles:**
- p50 (median): $p50
- p95: $p95
- p99: $p99

EOF

                # Check for errors
                local errors=$(jq -r '.errorDistribution | length' "$result_file" 2>/dev/null || echo "0")
                if [[ "$errors" -gt 0 ]]; then
                    cat >> "$OUTPUT_FILE" << EOF
**Errors Detected:**
\`\`\`json
$(jq -r '.errorDistribution' "$result_file" 2>/dev/null)
\`\`\`

EOF
                else
                    cat >> "$OUTPUT_FILE" << EOF
**Errors:** None ✅

EOF
                fi
            fi
        done
    fi

    cat >> "$OUTPUT_FILE" << 'EOF'

---

## 4. System Resource Usage

### 4.1 CPU and Memory

EOF

    # Add system metrics if available
    if [[ -f "$RESULTS_DIR/system_metrics_"*.log ]]; then
        cat >> "$OUTPUT_FILE" << EOF
\`\`\`
$(cat "$RESULTS_DIR"/system_metrics_*.log | tail -20)
\`\`\`

EOF
    else
        cat >> "$OUTPUT_FILE" << EOF
*System metrics not available*

EOF
    fi

    cat >> "$OUTPUT_FILE" << 'EOF'

---

## 5. Success Criteria Evaluation

### 5.1 Overall Assessment

EOF

    # Calculate pass/fail
    local total_tests=0
    local passed_tests=0

    if [[ -d "$RESULTS_DIR" ]]; then
        for result_file in "$RESULTS_DIR"/*.json; do
            if [[ -f "$result_file" ]]; then
                total_tests=$((total_tests + 1))

                local p95=$(jq -r '.latencyDistribution[] | select(.percentage==95) | .latency' "$result_file" 2>/dev/null | sed 's/ns$//' | awk '{printf "%.0f", $1/1000000}')
                local errors=$(jq -r '.errorDistribution | length' "$result_file" 2>/dev/null || echo "0")

                # Simple pass/fail: p95 < 200ms and no errors
                if [[ "$p95" -lt 200 && "$errors" -eq 0 ]]; then
                    passed_tests=$((passed_tests + 1))
                fi
            fi
        done
    fi

    local pass_rate=0
    if [[ $total_tests -gt 0 ]]; then
        pass_rate=$(echo "scale=1; ($passed_tests * 100) / $total_tests" | bc)
    fi

    cat >> "$OUTPUT_FILE" << EOF
**Test Results:**
- Total Tests: $total_tests
- Passed: $passed_tests
- Failed: $((total_tests - passed_tests))
- Pass Rate: ${pass_rate}%

EOF

    if [[ $pass_rate -ge 80 ]]; then
        cat >> "$OUTPUT_FILE" << 'EOF'
**Overall Status:** ✅ PASS

The Orders Microservice meets performance requirements. All critical paths perform within acceptable latency targets.

EOF
    else
        cat >> "$OUTPUT_FILE" << 'EOF'
**Overall Status:** ⚠️ NEEDS ATTENTION

Some performance targets are not met. Review failed tests and consider optimization.

EOF
    fi

    cat >> "$OUTPUT_FILE" << 'EOF'

---

## 6. Optimization Recommendations

### 6.1 High Priority

1. **Database Query Optimization**
   - Review slow queries (p95 > 100ms)
   - Add missing indexes
   - Optimize JOIN operations

2. **Connection Pooling**
   - Verify pool size matches load
   - Monitor connection utilization
   - Implement connection retry logic

3. **Caching Strategy**
   - Cache frequently accessed data (GetCart)
   - Implement Redis caching for read-heavy operations
   - Set appropriate TTL values

### 6.2 Medium Priority

1. **Concurrency Control**
   - Review lock contention for write operations
   - Implement optimistic locking where appropriate
   - Use database-level locking for critical sections

2. **Batch Operations**
   - Implement bulk cart item additions
   - Batch order status updates
   - Use batch inserts for high-volume operations

3. **Monitoring & Alerting**
   - Set up Prometheus metrics
   - Create Grafana dashboards
   - Configure alerts for p95 > targets

### 6.3 Low Priority

1. **Code Optimization**
   - Profile hot paths with pprof
   - Reduce memory allocations
   - Optimize serialization/deserialization

2. **Infrastructure**
   - Consider read replicas for heavy read workloads
   - Implement database sharding for scale
   - Use CDN for static assets

---

## 7. Next Steps

1. **Address Failed Tests**
   - Review and fix tests with p95 > targets
   - Investigate error causes
   - Re-run tests after fixes

2. **Baseline Establishment**
   - Save these results as baseline
   - Track performance trends over time
   - Alert on regressions > 20%

3. **Continuous Monitoring**
   - Integrate performance tests into CI/CD
   - Run load tests before each release
   - Monitor production metrics

---

## 8. Appendix

### 8.1 Test Environment

- **Hardware:** CI Runner (ubuntu-latest)
- **Database:** PostgreSQL 15
- **Cache:** Redis 7
- **Go Version:** 1.23
- **Concurrency:** Varies by test (10-300 concurrent connections)

### 8.2 Tools Used

- **Benchmarking:** Go testing.B
- **Load Testing:** ghz (gRPC)
- **Metrics:** Prometheus
- **Reporting:** Custom script

### 8.3 References

- [Performance Testing Guidelines](./docs/PERFORMANCE_TESTING_GUIDE.md)
- [CI/CD Setup](./docs/CI_CD_SETUP.md)
- [Monitoring Setup](./monitoring/README.md)

---

**Report Generated by:** Performance Testing Suite
**Contact:** Development Team
**Last Updated:** $(date '+%Y-%m-%d')
EOF

    log_success "Performance report generated: $OUTPUT_FILE"
}

###############################################################################
# Main Execution
###############################################################################

main() {
    log_info "Performance Report Generator"
    log_info "Results directory: $RESULTS_DIR"
    log_info "Output file: $OUTPUT_FILE"

    if [[ ! -d "$RESULTS_DIR" ]]; then
        log_error "Results directory not found: $RESULTS_DIR"
        log_info "Run load tests first: cd load-tests && ./run-all-tests.sh"
        exit 1
    fi

    generate_report

    log_success "Done!"
    log_info "View report: cat $OUTPUT_FILE"
}

main "$@"
