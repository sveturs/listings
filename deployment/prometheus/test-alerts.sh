#!/bin/bash
# Prometheus Alert Testing Script
# Tests alert rules by simulating various scenarios

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
ALERTMANAGER_URL="${ALERTMANAGER_URL:-http://localhost:9093}"

echo "=========================================="
echo "Prometheus Alert Testing Suite"
echo "=========================================="
echo ""
echo "Prometheus:   $PROMETHEUS_URL"
echo "Alertmanager: $ALERTMANAGER_URL"
echo ""

# ===================================================================
# Function: Check if Prometheus is running
# ===================================================================
check_prometheus() {
    echo -n "Checking Prometheus health... "

    if curl -sf "$PROMETHEUS_URL/-/healthy" &> /dev/null; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        echo "Error: Prometheus is not running or not healthy"
        echo "Start with: docker compose up -d prometheus"
        return 1
    fi
}

# ===================================================================
# Function: Check if Alertmanager is running
# ===================================================================
check_alertmanager() {
    echo -n "Checking Alertmanager health... "

    if curl -sf "$ALERTMANAGER_URL/-/healthy" &> /dev/null; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        echo "Warning: Alertmanager is not running"
        echo "Start with: docker compose up -d alertmanager"
        return 1
    fi
}

# ===================================================================
# Function: Reload Prometheus configuration
# ===================================================================
reload_prometheus() {
    echo -n "Reloading Prometheus configuration... "

    if curl -sf -X POST "$PROMETHEUS_URL/-/reload" &> /dev/null; then
        echo -e "${GREEN}✓${NC}"
        sleep 2 # Wait for reload to complete
        return 0
    else
        echo -e "${RED}✗${NC}"
        echo "Error: Failed to reload Prometheus"
        echo "Ensure Prometheus was started with --web.enable-lifecycle"
        return 1
    fi
}

# ===================================================================
# Function: List all alert rules
# ===================================================================
list_alert_rules() {
    echo ""
    echo "=========================================="
    echo "Alert Rules Inventory"
    echo "=========================================="

    local response=$(curl -sf "$PROMETHEUS_URL/api/v1/rules")

    if [[ -z "$response" ]]; then
        echo -e "${RED}✗ Failed to fetch rules${NC}"
        return 1
    fi

    # Parse and display rules
    echo "$response" | jq -r '
        .data.groups[] |
        "Group: \(.name) (interval: \(.interval)s)",
        (.rules[] |
            "  - \(.name) [\(.labels.severity // "info")]"
        )
    ' 2>/dev/null || {
        echo -e "${YELLOW}⊘ jq not installed, showing raw response${NC}"
        echo "$response"
    }

    # Count rules
    local total_rules=$(echo "$response" | jq '.data.groups[].rules | length' | awk '{s+=$1} END {print s}' 2>/dev/null || echo "?")
    echo ""
    echo "Total rules loaded: $total_rules"
}

# ===================================================================
# Function: Show active alerts
# ===================================================================
show_active_alerts() {
    echo ""
    echo "=========================================="
    echo "Active Alerts"
    echo "=========================================="

    local response=$(curl -sf "$PROMETHEUS_URL/api/v1/alerts")

    if [[ -z "$response" ]]; then
        echo -e "${RED}✗ Failed to fetch alerts${NC}"
        return 1
    fi

    # Parse and display alerts
    local alert_count=$(echo "$response" | jq '.data.alerts | length' 2>/dev/null || echo "0")

    if [[ "$alert_count" -eq 0 ]]; then
        echo -e "${GREEN}✓ No active alerts${NC}"
    else
        echo -e "${YELLOW}⚠ $alert_count active alerts:${NC}"
        echo ""

        echo "$response" | jq -r '
            .data.alerts[] |
            "Alert: \(.labels.alertname)",
            "  Severity: \(.labels.severity)",
            "  State: \(.state)",
            "  Summary: \(.annotations.summary // "N/A")",
            "  Since: \(.activeAt)",
            ""
        ' 2>/dev/null || echo "$response"
    fi
}

# ===================================================================
# Function: Test specific alert query
# ===================================================================
test_alert_query() {
    local alert_name=$1
    local query=$2

    echo ""
    echo "Testing alert: $alert_name"
    echo "Query: $query"
    echo ""

    # URL encode the query
    local encoded_query=$(echo "$query" | jq -sRr @uri)

    # Execute query
    local response=$(curl -sf "$PROMETHEUS_URL/api/v1/query?query=$encoded_query")

    if [[ -z "$response" ]]; then
        echo -e "${RED}✗ Query failed${NC}"
        return 1
    fi

    # Check if query returned results
    local result_count=$(echo "$response" | jq '.data.result | length' 2>/dev/null || echo "0")

    if [[ "$result_count" -gt 0 ]]; then
        echo -e "${YELLOW}⚠ Alert condition is TRUE (would fire)${NC}"
        echo "Result:"
        echo "$response" | jq -C '.data.result[]' 2>/dev/null || echo "$response"
    else
        echo -e "${GREEN}✓ Alert condition is FALSE (would not fire)${NC}"
    fi

    return 0
}

# ===================================================================
# Function: Test critical alerts
# ===================================================================
test_critical_alerts() {
    echo ""
    echo "=========================================="
    echo "Testing Critical Alerts"
    echo "=========================================="

    # Test ServiceDown
    test_alert_query "ServiceDown" \
        'up{job="listings-microservice"} == 0'

    # Test HighErrorRateCritical
    test_alert_query "HighErrorRateCritical" \
        '(sum(rate(http_requests_total{job="listings-microservice",status=~"5.."}[5m]))/sum(rate(http_requests_total{job="listings-microservice"}[5m])))*100 > 1'

    # Test DatabaseConnectionsHigh
    test_alert_query "DatabaseConnectionsHigh" \
        '(listings_db_connections_open / listings_db_connections_max) * 100 > 90'
}

# ===================================================================
# Function: Test SLO alerts
# ===================================================================
test_slo_alerts() {
    echo ""
    echo "=========================================="
    echo "Testing SLO Alerts"
    echo "=========================================="

    # Test Availability SLO
    test_alert_query "AvailabilitySLO" \
        '(sum(rate(http_requests_total{job="listings-microservice",status!~"5.."}[1h]))/sum(rate(http_requests_total{job="listings-microservice"}[1h])))*100 < 99.9'

    # Test Success Rate SLO
    test_alert_query "SuccessRateSLO" \
        '(sum(rate(http_requests_total{job="listings-microservice",status=~"2.."}[1h]))/sum(rate(http_requests_total{job="listings-microservice"}[1h])))*100 < 99.5'
}

# ===================================================================
# Function: Test recording rules
# ===================================================================
test_recording_rules() {
    echo ""
    echo "=========================================="
    echo "Testing Recording Rules"
    echo "=========================================="

    # Test if recording rules are producing data
    local recording_rules=(
        "job:http_requests:rate1m"
        "job:http_requests_error_rate:ratio5m"
        "job:http_request_duration:p95_5m"
        "service:availability:ratio1h"
        "job:cache_hit_ratio:rate5m"
    )

    for rule in "${recording_rules[@]}"; do
        echo -n "Testing $rule... "

        local encoded_query=$(echo "$rule" | jq -sRr @uri)
        local response=$(curl -sf "$PROMETHEUS_URL/api/v1/query?query=$encoded_query")
        local result_count=$(echo "$response" | jq '.data.result | length' 2>/dev/null || echo "0")

        if [[ "$result_count" -gt 0 ]]; then
            echo -e "${GREEN}✓ ($result_count series)${NC}"
        else
            echo -e "${YELLOW}⊘ (no data yet)${NC}"
        fi
    done
}

# ===================================================================
# Function: Show Prometheus targets
# ===================================================================
show_targets() {
    echo ""
    echo "=========================================="
    echo "Scrape Targets Status"
    echo "=========================================="

    local response=$(curl -sf "$PROMETHEUS_URL/api/v1/targets")

    if [[ -z "$response" ]]; then
        echo -e "${RED}✗ Failed to fetch targets${NC}"
        return 1
    fi

    # Parse and display targets
    echo "$response" | jq -r '
        .data.activeTargets[] |
        "Job: \(.labels.job)",
        "  Instance: \(.labels.instance)",
        "  Health: \(.health)",
        "  Last Scrape: \(.lastScrape)",
        "  Error: \(.lastError // "none")",
        ""
    ' 2>/dev/null || {
        echo -e "${YELLOW}⊘ jq not installed${NC}"
        echo "$response"
    }
}

# ===================================================================
# Function: Simulate alert (send test alert to Alertmanager)
# ===================================================================
simulate_alert() {
    echo ""
    echo "=========================================="
    echo "Simulating Test Alert"
    echo "=========================================="

    local alert_json=$(cat <<EOF
[
  {
    "labels": {
      "alertname": "TestAlert",
      "service": "listings",
      "severity": "warning",
      "instance": "test-instance"
    },
    "annotations": {
      "summary": "This is a test alert",
      "description": "Test alert generated by test-alerts.sh script"
    },
    "startsAt": "$(date -u +%Y-%m-%dT%H:%M:%S.000Z)",
    "endsAt": "$(date -u -d '+5 minutes' +%Y-%m-%dT%H:%M:%S.000Z)"
  }
]
EOF
)

    echo "Sending test alert to Alertmanager..."
    echo ""

    if curl -sf -X POST -H "Content-Type: application/json" \
        -d "$alert_json" \
        "$ALERTMANAGER_URL/api/v1/alerts" &> /dev/null; then
        echo -e "${GREEN}✓ Test alert sent successfully${NC}"
        echo ""
        echo "Check Alertmanager UI: $ALERTMANAGER_URL"
        echo "Alert should appear in ~30 seconds"
    else
        echo -e "${RED}✗ Failed to send test alert${NC}"
        return 1
    fi
}

# ===================================================================
# Function: Check metrics availability
# ===================================================================
check_metrics() {
    echo ""
    echo "=========================================="
    echo "Checking Metrics Availability"
    echo "=========================================="

    local metrics=(
        "http_requests_total"
        "http_request_duration_seconds_bucket"
        "listings_db_connections_open"
        "listings_cache_hits_total"
        "go_goroutines"
    )

    for metric in "${metrics[@]}"; do
        echo -n "Checking $metric... "

        local encoded_query=$(echo "$metric" | jq -sRr @uri)
        local response=$(curl -sf "$PROMETHEUS_URL/api/v1/query?query=$encoded_query")
        local result_count=$(echo "$response" | jq '.data.result | length' 2>/dev/null || echo "0")

        if [[ "$result_count" -gt 0 ]]; then
            echo -e "${GREEN}✓ ($result_count series)${NC}"
        else
            echo -e "${YELLOW}⊘ (not available)${NC}"
        fi
    done

    echo ""
    echo "Note: Some metrics may not be available if the service hasn't been running long enough"
}

# ===================================================================
# Main menu
# ===================================================================
show_menu() {
    echo ""
    echo "=========================================="
    echo "Select test to run:"
    echo "=========================================="
    echo "1. Full test suite (recommended)"
    echo "2. List alert rules"
    echo "3. Show active alerts"
    echo "4. Test critical alerts"
    echo "5. Test SLO alerts"
    echo "6. Test recording rules"
    echo "7. Show scrape targets"
    echo "8. Check metrics availability"
    echo "9. Simulate test alert"
    echo "10. Reload Prometheus config"
    echo "0. Exit"
    echo ""
    echo -n "Enter choice [0-10]: "
}

# ===================================================================
# Main function
# ===================================================================
main() {
    # Check prerequisites
    check_prometheus || exit 1
    check_alertmanager || true # Non-fatal

    # If argument provided, run that specific test
    case "${1:-}" in
        --full)
            list_alert_rules
            show_active_alerts
            test_critical_alerts
            test_slo_alerts
            test_recording_rules
            show_targets
            check_metrics
            exit 0
            ;;
        --reload)
            reload_prometheus
            exit 0
            ;;
        --simulate)
            simulate_alert
            exit 0
            ;;
    esac

    # Interactive menu
    while true; do
        show_menu
        read -r choice

        case $choice in
            1)
                list_alert_rules
                show_active_alerts
                test_critical_alerts
                test_slo_alerts
                test_recording_rules
                show_targets
                check_metrics
                ;;
            2) list_alert_rules ;;
            3) show_active_alerts ;;
            4) test_critical_alerts ;;
            5) test_slo_alerts ;;
            6) test_recording_rules ;;
            7) show_targets ;;
            8) check_metrics ;;
            9) simulate_alert ;;
            10) reload_prometheus ;;
            0)
                echo "Exiting..."
                exit 0
                ;;
            *)
                echo -e "${RED}Invalid choice${NC}"
                ;;
        esac
    done
}

# Run main function
main "$@"
