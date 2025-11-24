#!/bin/bash
#
# Attribute Indexer Test Runner
# Runs all tests for the attribute_indexer component
#
# Usage:
#   ./scripts/run_attribute_tests.sh [--unit|--integration|--bench|--all]
#

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test results
UNIT_RESULT=""
INTEGRATION_RESULT=""
BENCH_RESULT=""

# Functions
print_header() {
    echo ""
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW} $1${NC}"
    echo -e "${YELLOW}========================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

run_unit_tests() {
    print_header "Running Unit Tests"

    if go test ./internal/indexer -v -count=1 -short; then
        UNIT_RESULT="PASS"
        print_success "Unit tests passed"
    else
        UNIT_RESULT="FAIL"
        print_error "Unit tests failed"
        return 1
    fi
}

run_integration_tests() {
    print_header "Running Integration Tests"

    # Check database connectivity
    if ! docker exec listings_postgres psql -U listings_user -d listings_dev_db -c "SELECT 1" > /dev/null 2>&1; then
        print_error "Database not accessible. Is Docker container running?"
        INTEGRATION_RESULT="SKIP"
        return 1
    fi

    if go test ./internal/indexer -v -tags=integration -count=1; then
        INTEGRATION_RESULT="PASS"
        print_success "Integration tests passed"
    else
        INTEGRATION_RESULT="FAIL"
        print_error "Integration tests failed"
        return 1
    fi
}

run_benchmarks() {
    print_header "Running Performance Benchmarks"

    # Check database connectivity
    if ! docker exec listings_postgres psql -U listings_user -d listings_dev_db -c "SELECT 1" > /dev/null 2>&1; then
        print_error "Database not accessible. Is Docker container running?"
        BENCH_RESULT="SKIP"
        return 1
    fi

    if go test ./internal/indexer -bench=. -benchmem -tags=integration -run=^$; then
        BENCH_RESULT="PASS"
        print_success "Benchmarks completed"
    else
        BENCH_RESULT="FAIL"
        print_error "Benchmarks failed"
        return 1
    fi
}

verify_cache() {
    print_header "Verifying Cache Health"

    docker exec listings_postgres psql -U listings_user -d listings_dev_db <<EOF
SELECT
    COUNT(*) as cache_entries,
    pg_size_pretty(pg_total_relation_size('attribute_search_cache')) as total_size
FROM attribute_search_cache;
EOF

    print_success "Cache verification complete"
}

print_summary() {
    print_header "Test Summary"

    echo "Unit Tests:        ${UNIT_RESULT:-SKIP}"
    echo "Integration Tests: ${INTEGRATION_RESULT:-SKIP}"
    echo "Benchmarks:        ${BENCH_RESULT:-SKIP}"
    echo ""

    if [[ "$UNIT_RESULT" == "PASS" && "$INTEGRATION_RESULT" == "PASS" ]]; then
        print_success "All tests PASSED! ✅"
        echo ""
        echo "Production readiness: APPROVED"
        return 0
    elif [[ "$UNIT_RESULT" == "FAIL" || "$INTEGRATION_RESULT" == "FAIL" ]]; then
        print_error "Some tests FAILED! ❌"
        echo ""
        echo "Production readiness: NOT APPROVED"
        return 1
    else
        echo "Some tests were skipped"
        return 0
    fi
}

# Main logic
cd "$(dirname "$0")/.."

case "${1:-all}" in
    --unit)
        run_unit_tests
        ;;
    --integration)
        run_integration_tests
        verify_cache
        ;;
    --bench)
        run_benchmarks
        ;;
    --all)
        run_unit_tests || true
        run_integration_tests || true
        run_benchmarks || true
        verify_cache || true
        print_summary
        ;;
    *)
        echo "Usage: $0 [--unit|--integration|--bench|--all]"
        echo ""
        echo "Options:"
        echo "  --unit         Run unit tests only"
        echo "  --integration  Run integration tests + cache verification"
        echo "  --bench        Run performance benchmarks"
        echo "  --all          Run all tests and show summary (default)"
        exit 1
        ;;
esac
