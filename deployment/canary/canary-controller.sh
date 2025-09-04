#!/bin/bash

# Canary Release Controller for Unified Attributes System
# Manages gradual traffic rollout from blue to green environment

set -e

# Configuration
CONTROL_API="https://control.svetu.rs"
METRICS_API="https://grafana.svetu.rs/api"
SLACK_WEBHOOK="${SLACK_WEBHOOK:-}"
LOG_FILE="/var/log/canary-release.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${2:-$GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >> "$LOG_FILE"
}

# Send notification to Slack
notify_slack() {
    if [ -n "$SLACK_WEBHOOK" ]; then
        curl -X POST "$SLACK_WEBHOOK" \
            -H 'Content-Type: application/json' \
            -d "{\"text\": \"$1\"}" \
            -s > /dev/null
    fi
}

# Check environment health
check_health() {
    local env=$1
    local response=$(curl -s "${CONTROL_API}/health/${env}")
    
    if [ "$response" == "200" ]; then
        echo "healthy"
    else
        echo "unhealthy"
    fi
}

# Get current metrics
get_metrics() {
    local query=$1
    local response=$(curl -s -G "${METRICS_API}/datasources/proxy/1/api/v1/query" \
        --data-urlencode "query=${query}")
    
    echo "$response" | jq -r '.data.result[0].value[1]'
}

# Set canary percentage
set_canary_percent() {
    local percent=$1
    log "Setting canary traffic to ${percent}%"
    
    # Update nginx configuration
    curl -X POST "${CONTROL_API}/canary/${percent}" -s
    
    # Update application feature flags
    kubectl set env deployment/backend-green \
        UNIFIED_ATTRIBUTES_PERCENT="${percent}" \
        --namespace production
    
    kubectl set env deployment/frontend-green \
        NEXT_PUBLIC_UNIFIED_ATTRIBUTES_PERCENT="${percent}" \
        --namespace production
    
    log "Canary set to ${percent}%" "$GREEN"
    notify_slack ":rocket: Canary release updated to ${percent}%"
}

# Monitor canary metrics
monitor_canary() {
    local duration=$1
    local error_threshold=0.001  # 0.1%
    local latency_threshold=50   # 50ms
    
    log "Monitoring canary for ${duration} seconds..." "$YELLOW"
    
    local end_time=$(($(date +%s) + duration))
    local errors_detected=false
    
    while [ $(date +%s) -lt $end_time ]; do
        # Check error rate
        local error_rate=$(get_metrics 'rate(http_requests_total{status=~"5.."}[1m])')
        if (( $(echo "$error_rate > $error_threshold" | bc -l) )); then
            log "ERROR: High error rate detected: ${error_rate}" "$RED"
            errors_detected=true
            break
        fi
        
        # Check latency
        local p95_latency=$(get_metrics 'histogram_quantile(0.95, http_request_duration_seconds_bucket)')
        if (( $(echo "$p95_latency > $latency_threshold" | bc -l) )); then
            log "WARNING: High latency detected: ${p95_latency}ms" "$YELLOW"
        fi
        
        # Check database metrics
        local db_errors=$(get_metrics 'rate(postgres_errors_total[1m])')
        if (( $(echo "$db_errors > 0" | bc -l) )); then
            log "WARNING: Database errors detected: ${db_errors}" "$YELLOW"
        fi
        
        echo -n "."
        sleep 10
    done
    
    echo ""
    
    if [ "$errors_detected" = true ]; then
        return 1
    fi
    
    log "Monitoring complete. No issues detected." "$GREEN"
    return 0
}

# Rollback canary
rollback_canary() {
    log "ROLLING BACK CANARY RELEASE!" "$RED"
    notify_slack ":warning: Rolling back canary release due to errors!"
    
    # Set traffic to 0% (all blue)
    set_canary_percent 0
    
    # Disable unified attributes
    kubectl set env deployment/backend-green \
        USE_UNIFIED_ATTRIBUTES="false" \
        UNIFIED_ATTRIBUTES_FALLBACK="false" \
        --namespace production
    
    log "Rollback complete. All traffic on blue environment." "$GREEN"
    notify_slack ":white_check_mark: Rollback complete. System stable."
}

# Progressive canary rollout
progressive_rollout() {
    local stages=(10 25 50 100)
    local monitor_duration=600  # 10 minutes per stage
    
    log "Starting progressive canary rollout" "$BLUE"
    notify_slack ":rocket: Starting canary release for Unified Attributes System"
    
    # Pre-flight checks
    log "Running pre-flight checks..." "$YELLOW"
    
    # Check blue health
    if [ "$(check_health blue)" != "healthy" ]; then
        log "ERROR: Blue environment is not healthy!" "$RED"
        exit 1
    fi
    
    # Check green health
    if [ "$(check_health green)" != "healthy" ]; then
        log "ERROR: Green environment is not healthy!" "$RED"
        exit 1
    fi
    
    log "Pre-flight checks passed" "$GREEN"
    
    # Progressive rollout
    for percent in "${stages[@]}"; do
        log "=== Stage: ${percent}% ===" "$BLUE"
        
        # Set canary percentage
        set_canary_percent "$percent"
        
        # Wait for changes to propagate
        sleep 30
        
        # Monitor
        if monitor_canary "$monitor_duration"; then
            log "Stage ${percent}% completed successfully" "$GREEN"
            notify_slack ":white_check_mark: Canary at ${percent}% - No issues detected"
            
            # Ask for confirmation before proceeding to next stage
            if [ "$percent" != "100" ]; then
                read -p "Proceed to next stage? (y/n) " -n 1 -r
                echo
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    log "Canary release paused at ${percent}%" "$YELLOW"
                    exit 0
                fi
            fi
        else
            log "Issues detected at ${percent}%" "$RED"
            rollback_canary
            exit 1
        fi
    done
    
    log "CANARY RELEASE COMPLETED SUCCESSFULLY!" "$GREEN"
    notify_slack ":tada: Canary release completed! All traffic on green (unified attributes)"
}

# Manual control functions
manual_set() {
    local percent=$1
    if [ -z "$percent" ]; then
        echo "Usage: $0 set <percentage>"
        exit 1
    fi
    
    if [ "$percent" -lt 0 ] || [ "$percent" -gt 100 ]; then
        log "ERROR: Percentage must be between 0 and 100" "$RED"
        exit 1
    fi
    
    set_canary_percent "$percent"
}

# Status check
check_status() {
    local status=$(curl -s "${CONTROL_API}/status")
    echo -e "${BLUE}=== Canary Release Status ===${NC}"
    echo "$status" | jq .
}

# Main script logic
case "${1:-}" in
    start)
        progressive_rollout
        ;;
    set)
        manual_set "$2"
        ;;
    rollback)
        rollback_canary
        ;;
    status)
        check_status
        ;;
    monitor)
        monitor_canary "${2:-600}"
        ;;
    *)
        echo "Unified Attributes Canary Release Controller"
        echo ""
        echo "Usage: $0 {start|set|rollback|status|monitor} [options]"
        echo ""
        echo "Commands:"
        echo "  start              Start progressive rollout (10% -> 25% -> 50% -> 100%)"
        echo "  set <percent>      Manually set canary percentage"
        echo "  rollback          Rollback to blue environment (0%)"
        echo "  status            Check current deployment status"
        echo "  monitor [seconds] Monitor metrics (default: 600s)"
        echo ""
        echo "Examples:"
        echo "  $0 start          # Start automatic progressive rollout"
        echo "  $0 set 25        # Set canary to 25%"
        echo "  $0 rollback      # Emergency rollback"
        echo "  $0 status        # Check current status"
        exit 1
        ;;
esac