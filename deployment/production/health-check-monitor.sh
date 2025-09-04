#!/bin/bash

# Automated Health Check Monitor for Unified Attributes System
# Continuously monitors system health during canary deployment

set -e

# Configuration
API_BASE_URL="${API_BASE_URL:-https://api.svetu.rs}"
GRAFANA_URL="${GRAFANA_URL:-https://grafana.svetu.rs}"
SLACK_WEBHOOK="${SLACK_WEBHOOK:-}"
CHECK_INTERVAL="${CHECK_INTERVAL:-30}"  # seconds
ALERT_THRESHOLD_ERROR_RATE="0.001"      # 0.1%
ALERT_THRESHOLD_LATENCY="100"           # ms
ALERT_THRESHOLD_CACHE_HIT="0.6"         # 60%
LOG_FILE="/var/log/health-check-monitor.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Health status tracking
HEALTH_STATUS="healthy"
ERROR_COUNT=0
WARNING_COUNT=0
CHECK_COUNT=0

# Logging function
log() {
    local level=$1
    shift
    echo -e "${level}[$(date +'%Y-%m-%d %H:%M:%S')] $*${NC}" | tee -a "$LOG_FILE"
}

# Send alert
send_alert() {
    local severity=$1
    local message=$2
    
    log "$RED" "ALERT [$severity]: $message"
    
    # Send to Slack
    if [ -n "$SLACK_WEBHOOK" ]; then
        local emoji=""
        case $severity in
            "CRITICAL") emoji=":rotating_light:" ;;
            "WARNING") emoji=":warning:" ;;
            "INFO") emoji=":information_source:" ;;
        esac
        
        curl -X POST "$SLACK_WEBHOOK" \
            -H 'Content-Type: application/json' \
            -d "{\"text\": \"$emoji *Unified Attributes Health Check*\n*Severity:* $severity\n*Message:* $message\"}" \
            -s > /dev/null
    fi
    
    # Log to file for alertmanager
    echo "$(date -Iseconds),${severity},${message}" >> /var/log/health-alerts.csv
}

# Check endpoint health
check_endpoint() {
    local endpoint=$1
    local expected_status=$2
    local description=$3
    
    local response=$(curl -s -o /dev/null -w "%{http_code}" "${API_BASE_URL}${endpoint}")
    
    if [ "$response" -eq "$expected_status" ]; then
        log "$GREEN" "✓ $description: OK ($response)"
        return 0
    else
        log "$RED" "✗ $description: FAILED (Expected: $expected_status, Got: $response)"
        return 1
    fi
}

# Check metrics from Prometheus
check_metrics() {
    local metric_query=$1
    local threshold=$2
    local comparison=$3
    local description=$4
    
    local value=$(curl -s -G "${GRAFANA_URL}/api/datasources/proxy/1/api/v1/query" \
        --data-urlencode "query=${metric_query}" \
        -H "Authorization: Bearer ${GRAFANA_API_KEY:-}" \
        | jq -r '.data.result[0].value[1] // 0')
    
    local result=""
    case $comparison in
        "lt")
            if (( $(echo "$value < $threshold" | bc -l) )); then
                result="OK"
            else
                result="FAILED"
            fi
            ;;
        "gt")
            if (( $(echo "$value > $threshold" | bc -l) )); then
                result="OK"
            else
                result="FAILED"
            fi
            ;;
    esac
    
    if [ "$result" = "OK" ]; then
        log "$GREEN" "✓ $description: $value (threshold: $comparison $threshold)"
        return 0
    else
        log "$RED" "✗ $description: $value (threshold: $comparison $threshold)"
        return 1
    fi
}

# Comprehensive health check
perform_health_check() {
    local check_passed=true
    
    log "$BLUE" "=== Health Check #$CHECK_COUNT ==="
    
    # 1. API Endpoints Health
    log "$YELLOW" "Checking API endpoints..."
    
    check_endpoint "/health/live" 200 "Liveness probe" || check_passed=false
    check_endpoint "/health/ready" 200 "Readiness probe" || check_passed=false
    check_endpoint "/api/v2/categories/1/attributes" 200 "Get attributes endpoint" || check_passed=false
    
    # 2. Database Health
    log "$YELLOW" "Checking database health..."
    
    local db_check=$(curl -s "${API_BASE_URL}/health/db" | jq -r '.status // "unknown"')
    if [ "$db_check" = "healthy" ]; then
        log "$GREEN" "✓ Database connection: healthy"
    else
        log "$RED" "✗ Database connection: $db_check"
        check_passed=false
    fi
    
    # 3. Performance Metrics
    log "$YELLOW" "Checking performance metrics..."
    
    # Error rate
    check_metrics \
        'sum(rate(http_requests_total{status=~"5.."}[1m])) / sum(rate(http_requests_total[1m]))' \
        "$ALERT_THRESHOLD_ERROR_RATE" \
        "lt" \
        "Error rate" || check_passed=false
    
    # P95 Latency
    check_metrics \
        'histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[1m])) by (le)) * 1000' \
        "$ALERT_THRESHOLD_LATENCY" \
        "lt" \
        "P95 latency (ms)" || check_passed=false
    
    # Cache hit rate
    check_metrics \
        'sum(rate(cache_hits_total[1m])) / sum(rate(cache_requests_total[1m]))' \
        "$ALERT_THRESHOLD_CACHE_HIT" \
        "gt" \
        "Cache hit rate" || check_passed=false
    
    # 4. Dual-write validation
    log "$YELLOW" "Checking dual-write status..."
    
    check_metrics \
        'rate(dual_write_success_total[1m]) / rate(dual_write_attempts_total[1m])' \
        "0.99" \
        "gt" \
        "Dual-write success rate" || check_passed=false
    
    # 5. Canary traffic distribution
    log "$YELLOW" "Checking canary traffic..."
    
    local canary_percent=$(curl -s "${API_BASE_URL}/internal/canary/status" | jq -r '.percent // 0')
    log "$BLUE" "Current canary traffic: ${canary_percent}%"
    
    # 6. Memory and CPU usage
    log "$YELLOW" "Checking resource usage..."
    
    check_metrics \
        'avg(container_memory_usage_bytes{pod=~"backend-.*"}) / avg(container_spec_memory_limit_bytes{pod=~"backend-.*"})' \
        "0.9" \
        "lt" \
        "Memory usage" || check_passed=false
    
    check_metrics \
        'avg(rate(container_cpu_usage_seconds_total{pod=~"backend-.*"}[1m]))' \
        "0.8" \
        "lt" \
        "CPU usage" || check_passed=false
    
    # Update health status
    if [ "$check_passed" = true ]; then
        HEALTH_STATUS="healthy"
        ERROR_COUNT=0
        log "$GREEN" "=== Health Check PASSED ==="
    else
        ERROR_COUNT=$((ERROR_COUNT + 1))
        if [ $ERROR_COUNT -ge 3 ]; then
            HEALTH_STATUS="unhealthy"
            send_alert "CRITICAL" "System unhealthy for 3 consecutive checks"
        else
            HEALTH_STATUS="degraded"
            send_alert "WARNING" "Health check failed ($ERROR_COUNT/3)"
        fi
        log "$RED" "=== Health Check FAILED ==="
    fi
    
    echo
    return $([ "$check_passed" = true ] && echo 0 || echo 1)
}

# Synthetic transaction test
synthetic_test() {
    log "$BLUE" "Running synthetic transaction test..."
    
    # 1. Create a test listing with attributes
    local test_listing=$(cat <<EOF
{
  "title": "Health Check Test Listing",
  "category_id": 1,
  "attributes": {
    "brand": "TestBrand",
    "model": "TestModel",
    "year": 2024
  }
}
EOF
)
    
    local create_response=$(curl -s -X POST \
        "${API_BASE_URL}/api/v2/listings" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TEST_TOKEN:-}" \
        -d "$test_listing")
    
    local listing_id=$(echo "$create_response" | jq -r '.data.id // ""')
    
    if [ -n "$listing_id" ]; then
        log "$GREEN" "✓ Created test listing: $listing_id"
        
        # 2. Verify attributes were saved
        local get_response=$(curl -s \
            "${API_BASE_URL}/api/v2/listings/${listing_id}" \
            -H "Authorization: Bearer ${TEST_TOKEN:-}")
        
        local saved_brand=$(echo "$get_response" | jq -r '.data.attributes.brand // ""')
        
        if [ "$saved_brand" = "TestBrand" ]; then
            log "$GREEN" "✓ Attributes correctly saved and retrieved"
        else
            log "$RED" "✗ Attributes not correctly saved"
            send_alert "WARNING" "Synthetic test: Attributes save/retrieve failed"
        fi
        
        # 3. Clean up test data
        curl -s -X DELETE \
            "${API_BASE_URL}/api/v2/listings/${listing_id}" \
            -H "Authorization: Bearer ${TEST_TOKEN:-}" > /dev/null
        
        log "$GREEN" "✓ Test listing cleaned up"
    else
        log "$RED" "✗ Failed to create test listing"
        send_alert "CRITICAL" "Synthetic test: Cannot create listing"
    fi
}

# Monitor dashboard status
check_dashboard() {
    log "$BLUE" "Checking monitoring dashboard..."
    
    local dashboard_status=$(curl -s \
        "${GRAFANA_URL}/api/dashboards/uid/unified-attrs-canary" \
        -H "Authorization: Bearer ${GRAFANA_API_KEY:-}" \
        | jq -r '.dashboard.title // "not found"')
    
    if [ "$dashboard_status" != "not found" ]; then
        log "$GREEN" "✓ Grafana dashboard accessible: $dashboard_status"
        
        # Check for active alerts
        local active_alerts=$(curl -s \
            "${GRAFANA_URL}/api/alerts" \
            -H "Authorization: Bearer ${GRAFANA_API_KEY:-}" \
            | jq '[.[] | select(.state == "alerting")] | length')
        
        if [ "$active_alerts" -gt 0 ]; then
            log "$YELLOW" "⚠ Active alerts in Grafana: $active_alerts"
            send_alert "WARNING" "Grafana has $active_alerts active alerts"
        else
            log "$GREEN" "✓ No active alerts in Grafana"
        fi
    else
        log "$YELLOW" "⚠ Cannot access Grafana dashboard"
    fi
}

# Main monitoring loop
main() {
    log "$BLUE" "========================================="
    log "$BLUE" "  UNIFIED ATTRIBUTES HEALTH MONITOR     "
    log "$BLUE" "========================================="
    log "$BLUE" "API: $API_BASE_URL"
    log "$BLUE" "Check Interval: ${CHECK_INTERVAL}s"
    log "$BLUE" "========================================="
    echo
    
    # Trap signals for clean shutdown
    trap 'log "$YELLOW" "Health monitor stopped"; exit 0' INT TERM
    
    # Initial synthetic test
    synthetic_test
    
    # Continuous monitoring loop
    while true; do
        CHECK_COUNT=$((CHECK_COUNT + 1))
        
        # Perform health check
        if perform_health_check; then
            # Run synthetic test every 10 checks
            if [ $((CHECK_COUNT % 10)) -eq 0 ]; then
                synthetic_test
            fi
            
            # Check dashboard every 5 checks
            if [ $((CHECK_COUNT % 5)) -eq 0 ]; then
                check_dashboard
            fi
        else
            # If health check failed, increase frequency
            log "$YELLOW" "Increasing check frequency due to failures..."
            sleep 10
            continue
        fi
        
        # Wait for next check
        sleep "$CHECK_INTERVAL"
    done
}

# Export health status for other scripts
export_status() {
    cat <<EOF > /tmp/unified-attrs-health.json
{
  "status": "$HEALTH_STATUS",
  "last_check": "$(date -Iseconds)",
  "error_count": $ERROR_COUNT,
  "warning_count": $WARNING_COUNT,
  "check_count": $CHECK_COUNT
}
EOF
}

# Run main function
main "$@"