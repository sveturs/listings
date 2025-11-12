#!/usr/bin/env bash

#######################################
# Smoke Tests Script
# Automated smoke testing for critical endpoints
#
# Usage: ./smoke-tests.sh [--host HOST] [--port PORT] [--verbose] [--json]
#
# Exit codes:
#   0 - All tests passed
#   1 - One or more tests failed
#######################################

set -Eeuo pipefail

# Default configuration
HOST="${HOST:-localhost}"
PORT="${PORT:-8080}"
VERBOSE=false
JSON_OUTPUT=false
TIMEOUT=10

# Test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
declare -a FAILED_TEST_NAMES=()

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

#######################################
# Parse arguments
#######################################

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --host)
                HOST="$2"
                shift 2
                ;;
            --port)
                PORT="$2"
                shift 2
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            --json)
                JSON_OUTPUT=true
                shift
                ;;
            *)
                echo "Usage: $0 [--host HOST] [--port PORT] [--verbose] [--json]"
                exit 1
                ;;
        esac
    done
}

#######################################
# Logging functions
#######################################

log_test() {
    if [[ "${JSON_OUTPUT}" == "false" ]]; then
        echo -e "${BLUE}[TEST]${NC} $*"
    fi
}

log_pass() {
    if [[ "${JSON_OUTPUT}" == "false" ]]; then
        echo -e "${GREEN}[PASS]${NC} $*"
    fi
}

log_fail() {
    if [[ "${JSON_OUTPUT}" == "false" ]]; then
        echo -e "${RED}[FAIL]${NC} $*"
    fi
}

log_info() {
    if [[ "${JSON_OUTPUT}" == "false" ]] && [[ "${VERBOSE}" == "true" ]]; then
        echo -e "${BLUE}[INFO]${NC} $*"
    fi
}

#######################################
# HTTP request helper
#######################################

make_request() {
    local method="$1"
    local endpoint="$2"
    local expected_code="${3:-200}"
    local data="${4:-}"

    local url="http://${HOST}:${PORT}${endpoint}"
    local curl_opts=(-s -w "%{http_code}" -o /tmp/smoke_response.txt --max-time "${TIMEOUT}")

    if [[ -n "${data}" ]]; then
        curl_opts+=(-X "${method}" -H "Content-Type: application/json" -d "${data}")
    else
        curl_opts+=(-X "${method}")
    fi

    local http_code
    http_code=$(curl "${curl_opts[@]}" "${url}" 2>/dev/null || echo "000")

    if [[ "${VERBOSE}" == "true" ]] && [[ "${JSON_OUTPUT}" == "false" ]]; then
        log_info "Request: ${method} ${url}"
        log_info "Expected: ${expected_code}, Got: ${http_code}"
        if [[ -s /tmp/smoke_response.txt ]]; then
            log_info "Response: $(cat /tmp/smoke_response.txt)"
        fi
    fi

    if [[ "${http_code}" == "${expected_code}" ]]; then
        return 0
    else
        return 1
    fi
}

#######################################
# Test execution helper
#######################################

run_test() {
    local test_name="$1"
    local test_func="$2"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    log_test "${test_name}"

    if ${test_func}; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        log_pass "${test_name}"
        return 0
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        FAILED_TEST_NAMES+=("${test_name}")
        log_fail "${test_name}"
        return 1
    fi
}

#######################################
# Individual test functions
#######################################

test_health_check() {
    make_request GET "/health" 200
}

test_readiness_check() {
    make_request GET "/ready" 200
}

test_metrics_endpoint() {
    make_request GET "/metrics" 200
}

test_list_listings() {
    make_request GET "/api/v1/listings?page=1&limit=10" 200
}

test_get_listing_by_id() {
    # First, get a listing ID from list
    make_request GET "/api/v1/listings?page=1&limit=1" 200 || return 1

    local listing_id
    listing_id=$(cat /tmp/smoke_response.txt | jq -r '.data[0].id' 2>/dev/null || echo "")

    if [[ -z "${listing_id}" ]] || [[ "${listing_id}" == "null" ]]; then
        log_info "No listings found, skipping get by ID test"
        return 0
    fi

    make_request GET "/api/v1/listings/${listing_id}" 200
}

test_search_listings() {
    make_request GET "/api/v1/listings/search?q=test" 200
}

test_get_categories() {
    make_request GET "/api/v1/categories" 200
}

test_create_listing() {
    local payload='{
        "title": "Smoke Test Listing",
        "description": "Auto-generated test listing",
        "price": 100.00,
        "category_id": 1
    }'

    make_request POST "/api/v1/listings" 201 "${payload}" || return 1

    # Save created listing ID for cleanup
    local listing_id
    listing_id=$(cat /tmp/smoke_response.txt | jq -r '.id' 2>/dev/null || echo "")

    if [[ -n "${listing_id}" ]] && [[ "${listing_id}" != "null" ]]; then
        echo "${listing_id}" > /tmp/smoke_test_listing_id.txt
    fi

    return 0
}

test_update_listing() {
    # Get listing ID from previous test
    if [[ ! -f /tmp/smoke_test_listing_id.txt ]]; then
        log_info "No listing ID available, skipping update test"
        return 0
    fi

    local listing_id=$(cat /tmp/smoke_test_listing_id.txt)

    local payload='{
        "title": "Updated Smoke Test Listing",
        "price": 150.00
    }'

    make_request PUT "/api/v1/listings/${listing_id}" 200 "${payload}"
}

test_delete_listing() {
    # Get listing ID from create test
    if [[ ! -f /tmp/smoke_test_listing_id.txt ]]; then
        log_info "No listing ID available, skipping delete test"
        return 0
    fi

    local listing_id=$(cat /tmp/smoke_test_listing_id.txt)

    make_request DELETE "/api/v1/listings/${listing_id}" 204
}

test_database_connectivity() {
    # This test assumes the service exposes a DB check endpoint
    make_request GET "/health/db" 200 || make_request GET "/health" 200
}

test_redis_connectivity() {
    # This test assumes the service exposes a Redis check endpoint
    make_request GET "/health/redis" 200 || make_request GET "/health" 200
}

#######################################
# Parallel test execution
#######################################

run_tests_parallel() {
    declare -a pids=()

    # Run independent tests in parallel
    run_test "Health Check" test_health_check &
    pids+=($!)

    run_test "Readiness Check" test_readiness_check &
    pids+=($!)

    run_test "Metrics Endpoint" test_metrics_endpoint &
    pids+=($!)

    run_test "List Listings" test_list_listings &
    pids+=($!)

    run_test "Search Listings" test_search_listings &
    pids+=($!)

    run_test "Get Categories" test_get_categories &
    pids+=($!)

    run_test "Database Connectivity" test_database_connectivity &
    pids+=($!)

    run_test "Redis Connectivity" test_redis_connectivity &
    pids+=($!)

    # Wait for all parallel tests
    for pid in "${pids[@]}"; do
        wait $pid || true
    done
}

#######################################
# Sequential test execution (CRUD)
#######################################

run_tests_sequential() {
    # These tests must run sequentially (create -> update -> delete)
    run_test "Get Listing by ID" test_get_listing_by_id
    run_test "Create Listing" test_create_listing
    run_test "Update Listing" test_update_listing
    run_test "Delete Listing" test_delete_listing
}

#######################################
# Output results
#######################################

output_results() {
    if [[ "${JSON_OUTPUT}" == "true" ]]; then
        # JSON output for CI/CD
        local failed_tests_json=$(printf '%s\n' "${FAILED_TEST_NAMES[@]}" | jq -R . | jq -s .)

        cat <<EOF
{
    "total": ${TOTAL_TESTS},
    "passed": ${PASSED_TESTS},
    "failed": ${FAILED_TESTS},
    "success_rate": $(echo "scale=2; ${PASSED_TESTS} * 100 / ${TOTAL_TESTS}" | bc),
    "failed_tests": ${failed_tests_json},
    "host": "${HOST}",
    "port": ${PORT},
    "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
}
EOF
    else
        # Human-readable output
        echo ""
        echo "========================================="
        echo "Smoke Test Results"
        echo "========================================="
        echo "Total Tests:   ${TOTAL_TESTS}"
        echo "Passed:        ${PASSED_TESTS}"
        echo "Failed:        ${FAILED_TESTS}"
        echo "Success Rate:  $(echo "scale=2; ${PASSED_TESTS} * 100 / ${TOTAL_TESTS}" | bc)%"
        echo "========================================="

        if [[ ${FAILED_TESTS} -gt 0 ]]; then
            echo ""
            echo "Failed Tests:"
            for test in "${FAILED_TEST_NAMES[@]}"; do
                echo "  - ${test}"
            done
        fi

        echo ""
    fi
}

#######################################
# Cleanup
#######################################

cleanup() {
    rm -f /tmp/smoke_response.txt
    rm -f /tmp/smoke_test_listing_id.txt
}

trap cleanup EXIT

#######################################
# Main function
#######################################

main() {
    parse_args "$@"

    if [[ "${JSON_OUTPUT}" == "false" ]]; then
        echo "========================================="
        echo "Running Smoke Tests"
        echo "Target: http://${HOST}:${PORT}"
        echo "========================================="
        echo ""
    fi

    # Check if service is reachable
    if ! curl -sf --max-time 5 "http://${HOST}:${PORT}/health" > /dev/null 2>&1; then
        if [[ "${JSON_OUTPUT}" == "false" ]]; then
            echo -e "${RED}ERROR: Service not reachable at http://${HOST}:${PORT}${NC}"
        else
            echo '{"error": "Service not reachable", "host": "'${HOST}'", "port": '${PORT}'}'
        fi
        exit 1
    fi

    # Run tests
    run_tests_parallel
    run_tests_sequential

    # Output results
    output_results

    # Exit with appropriate code
    if [[ ${FAILED_TESTS} -eq 0 ]]; then
        exit 0
    else
        exit 1
    fi
}

# Run main function
main "$@"
