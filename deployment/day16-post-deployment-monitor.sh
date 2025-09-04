#!/bin/bash
# Post-Deployment Monitoring Script - Day 16
# Continuous monitoring and performance baseline collection
# Version: 1.0.0
# Date: 04.09.2025

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
MONITORING_DURATION_HOURS=24
CHECK_INTERVAL_SECONDS=300  # 5 minutes
API_BASE_URL="https://api.svetu.rs"
PROMETHEUS_URL="https://prometheus.svetu.rs/api/v1"
GRAFANA_URL="https://grafana.svetu.rs"
ELASTICSEARCH_URL="https://elasticsearch.svetu.rs"
REPORT_DIR="/var/reports/unified-attributes"
BASELINE_FILE="$REPORT_DIR/performance-baseline-$(date +%Y%m%d).json"
ALERT_THRESHOLD_ERROR_RATE=0.001  # 0.1%
ALERT_THRESHOLD_P95_LATENCY=100   # milliseconds
ALERT_THRESHOLD_MEMORY_USAGE=90   # percentage
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL:-}"

# Metrics storage
declare -A current_metrics
declare -A baseline_metrics
declare -A anomalies

# Initialize directories
mkdir -p "$REPORT_DIR"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date +"%Y-%m-%d %H:%M:%S") - $1" | tee -a "$REPORT_DIR/monitoring.log"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date +"%Y-%m-%d %H:%M:%S") - $1" | tee -a "$REPORT_DIR/monitoring.log"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $(date +"%Y-%m-%d %H:%M:%S") - $1" | tee -a "$REPORT_DIR/monitoring.log"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date +"%Y-%m-%d %H:%M:%S") - $1" | tee -a "$REPORT_DIR/monitoring.log"
}

log_metric() {
    echo -e "${CYAN}[METRIC]${NC} $1: $2" | tee -a "$REPORT_DIR/metrics.log"
}

# Slack notification
send_alert() {
    local severity="$1"
    local message="$2"
    local color="#36a64f"
    
    case "$severity" in
        "critical")
            color="#ff0000"
            ;;
        "warning")
            color="#ff9900"
            ;;
        "info")
            color="#439FE0"
            ;;
    esac
    
    if [[ -n "$SLACK_WEBHOOK_URL" ]]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{
                \"attachments\": [{
                    \"color\": \"${color}\",
                    \"title\": \"Post-Deployment Monitoring Alert\",
                    \"text\": \"${message}\",
                    \"footer\": \"Unified Attributes System\",
                    \"ts\": $(date +%s)
                }]
            }" \
            "$SLACK_WEBHOOK_URL" 2>/dev/null || true
    fi
}

# Metric collection functions
collect_response_time_metrics() {
    local endpoint="$1"
    local label="$2"
    
    # Make 10 requests and calculate statistics
    local times=()
    for i in {1..10}; do
        local response_time=$(curl -w "%{time_total}" -o /dev/null -sf "$endpoint" 2>/dev/null || echo "999")
        times+=("$response_time")
        sleep 0.5
    done
    
    # Calculate p50, p95, p99
    IFS=$'\n' sorted=($(sort -n <<<"${times[*]}"))
    unset IFS
    
    local p50_index=$((${#sorted[@]} * 50 / 100))
    local p95_index=$((${#sorted[@]} * 95 / 100))
    local p99_index=$((${#sorted[@]} * 99 / 100))
    
    current_metrics["${label}_p50"]=$(echo "${sorted[$p50_index]} * 1000" | bc)
    current_metrics["${label}_p95"]=$(echo "${sorted[$p95_index]} * 1000" | bc)
    current_metrics["${label}_p99"]=$(echo "${sorted[$p99_index]} * 1000" | bc)
    
    log_metric "${label}_p50" "${current_metrics[${label}_p50]}ms"
    log_metric "${label}_p95" "${current_metrics[${label}_p95]}ms"
    log_metric "${label}_p99" "${current_metrics[${label}_p99]}ms"
}

collect_error_rate() {
    # Query Prometheus for error rate
    local query="rate(http_requests_total{status=~\"5..\"}[5m])"
    local error_rate=$(curl -sf "${PROMETHEUS_URL}/query?query=${query}" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    current_metrics["error_rate"]="$error_rate"
    log_metric "error_rate" "${error_rate}"
    
    # Check threshold
    if (( $(echo "$error_rate > $ALERT_THRESHOLD_ERROR_RATE" | bc -l) )); then
        log_error "Error rate exceeds threshold: ${error_rate} > ${ALERT_THRESHOLD_ERROR_RATE}"
        send_alert "critical" "⚠️ Error rate critical: ${error_rate} (threshold: ${ALERT_THRESHOLD_ERROR_RATE})"
        anomalies["high_error_rate"]="$error_rate"
    fi
}

collect_throughput_metrics() {
    # Query for request rate
    local query="rate(http_requests_total[5m])"
    local request_rate=$(curl -sf "${PROMETHEUS_URL}/query?query=${query}" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    current_metrics["throughput"]="$request_rate"
    log_metric "throughput" "${request_rate} req/s"
}

collect_cache_metrics() {
    # Query cache hit rate
    local cache_hits=$(curl -sf "${PROMETHEUS_URL}/query?query=unified_attributes_cache_hits_total" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    local cache_misses=$(curl -sf "${PROMETHEUS_URL}/query?query=unified_attributes_cache_misses_total" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    local total=$((cache_hits + cache_misses))
    local hit_rate=0
    if [[ $total -gt 0 ]]; then
        hit_rate=$(echo "scale=2; $cache_hits * 100 / $total" | bc)
    fi
    
    current_metrics["cache_hit_rate"]="$hit_rate"
    log_metric "cache_hit_rate" "${hit_rate}%"
}

collect_resource_metrics() {
    # CPU usage
    local cpu_usage=$(curl -sf "${PROMETHEUS_URL}/query?query=container_cpu_usage_seconds_total" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    current_metrics["cpu_usage"]="$cpu_usage"
    log_metric "cpu_usage" "${cpu_usage}%"
    
    # Memory usage
    local memory_usage=$(curl -sf "${PROMETHEUS_URL}/query?query=container_memory_usage_bytes" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    local memory_usage_mb=$(echo "$memory_usage / 1024 / 1024" | bc)
    current_metrics["memory_usage_mb"]="$memory_usage_mb"
    log_metric "memory_usage" "${memory_usage_mb}MB"
    
    # Check memory threshold
    local memory_percent=$(curl -sf "${PROMETHEUS_URL}/query?query=container_memory_usage_bytes/container_spec_memory_limit_bytes*100" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    if (( $(echo "$memory_percent > $ALERT_THRESHOLD_MEMORY_USAGE" | bc -l) )); then
        log_warning "Memory usage high: ${memory_percent}%"
        send_alert "warning" "⚠️ Memory usage high: ${memory_percent}% (threshold: ${ALERT_THRESHOLD_MEMORY_USAGE}%)"
        anomalies["high_memory"]="$memory_percent"
    fi
}

collect_database_metrics() {
    # Active connections
    local db_connections=$(curl -sf "${PROMETHEUS_URL}/query?query=pg_stat_database_numbackends" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    current_metrics["db_connections"]="$db_connections"
    log_metric "db_connections" "$db_connections"
    
    # Query execution time
    local query_time=$(curl -sf "${PROMETHEUS_URL}/query?query=pg_stat_statements_mean_seconds" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    local query_time_ms=$(echo "$query_time * 1000" | bc)
    current_metrics["db_query_time_ms"]="$query_time_ms"
    log_metric "db_query_time" "${query_time_ms}ms"
    
    # Replication lag
    local replication_lag=$(curl -sf "${PROMETHEUS_URL}/query?query=pg_replication_lag_seconds" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    current_metrics["replication_lag"]="$replication_lag"
    log_metric "replication_lag" "${replication_lag}s"
}

collect_business_metrics() {
    # Listings created
    local listings_created=$(curl -sf "${PROMETHEUS_URL}/query?query=increase(listings_created_total[1h])" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    current_metrics["listings_created_hourly"]="$listings_created"
    log_metric "listings_created_hourly" "$listings_created"
    
    # Searches performed
    local searches=$(curl -sf "${PROMETHEUS_URL}/query?query=increase(searches_total[1h])" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    current_metrics["searches_hourly"]="$searches"
    log_metric "searches_hourly" "$searches"
    
    # Attribute usage
    local attribute_usage=$(curl -sf "${PROMETHEUS_URL}/query?query=increase(unified_attributes_usage_total[1h])" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    current_metrics["attribute_usage_hourly"]="$attribute_usage"
    log_metric "attribute_usage_hourly" "$attribute_usage"
}

analyze_logs() {
    log_info "Analyzing application logs for errors..."
    
    # Get error logs from last hour
    local error_count=$(curl -sf "${ELASTICSEARCH_URL}/_search" \
        -H "Content-Type: application/json" \
        -d '{
            "query": {
                "bool": {
                    "must": [
                        {"match": {"level": "ERROR"}},
                        {"range": {"timestamp": {"gte": "now-1h"}}}
                    ]
                }
            },
            "size": 0
        }' | jq -r '.hits.total.value' 2>/dev/null || echo "0")
    
    current_metrics["error_logs_hourly"]="$error_count"
    log_metric "error_logs_hourly" "$error_count"
    
    if [[ $error_count -gt 100 ]]; then
        log_warning "High number of error logs: $error_count in last hour"
        anomalies["high_error_logs"]="$error_count"
    fi
}

compare_with_baseline() {
    if [[ -f "$BASELINE_FILE" ]]; then
        log_info "Comparing with baseline metrics..."
        
        # Load baseline
        while IFS="=" read -r key value; do
            baseline_metrics["$key"]="$value"
        done < <(jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' "$BASELINE_FILE")
        
        # Compare key metrics
        for metric in "error_rate" "cache_hit_rate" "db_query_time_ms"; do
            if [[ -n "${baseline_metrics[$metric]:-}" ]] && [[ -n "${current_metrics[$metric]:-}" ]]; then
                local baseline="${baseline_metrics[$metric]}"
                local current="${current_metrics[$metric]}"
                local diff=$(echo "scale=2; (($current - $baseline) / $baseline) * 100" | bc)
                
                if (( $(echo "$diff > 20" | bc -l) )); then
                    log_warning "Metric $metric degraded by ${diff}% compared to baseline"
                    anomalies["${metric}_degradation"]="$diff"
                elif (( $(echo "$diff < -20" | bc -l) )); then
                    log_success "Metric $metric improved by ${diff#-}% compared to baseline"
                fi
            fi
        done
    fi
}

generate_hourly_report() {
    local report_file="$REPORT_DIR/hourly-report-$(date +%Y%m%d-%H).json"
    
    cat > "$report_file" <<EOF
{
  "timestamp": "$(date -Iseconds)",
  "environment": "production",
  "system": "unified_attributes",
  "metrics": $(echo '{}' | jq --argjson metrics "$(declare -p current_metrics | sed 's/declare -A current_metrics=//' | sed "s/'/\"/g")" '$metrics'),
  "anomalies": $(echo '{}' | jq --argjson anomalies "$(declare -p anomalies | sed 's/declare -A anomalies=//' | sed "s/'/\"/g")" '$anomalies'),
  "health_status": "$([ ${#anomalies[@]} -eq 0 ] && echo 'healthy' || echo 'degraded')"
}
EOF
    
    log_success "Hourly report generated: $report_file"
}

update_baseline() {
    # Update baseline with current metrics if healthy
    if [[ ${#anomalies[@]} -eq 0 ]]; then
        echo "{" > "$BASELINE_FILE"
        for key in "${!current_metrics[@]}"; do
            echo "  \"$key\": ${current_metrics[$key]}," >> "$BASELINE_FILE"
        done
        # Remove last comma and close JSON
        sed -i '$ s/,$//' "$BASELINE_FILE"
        echo "}" >> "$BASELINE_FILE"
        
        log_success "Performance baseline updated"
    else
        log_warning "Baseline not updated due to anomalies"
    fi
}

continuous_monitoring() {
    local start_time=$(date +%s)
    local end_time=$((start_time + (MONITORING_DURATION_HOURS * 3600)))
    local check_count=0
    
    log_info "Starting continuous monitoring for ${MONITORING_DURATION_HOURS} hours"
    send_alert "info" "📊 Post-deployment monitoring started for ${MONITORING_DURATION_HOURS} hours"
    
    while [[ $(date +%s) -lt $end_time ]]; do
        check_count=$((check_count + 1))
        log_info "=== Monitoring Check #${check_count} ==="
        
        # Clear previous metrics and anomalies
        current_metrics=()
        anomalies=()
        
        # Collect all metrics
        log_info "Collecting performance metrics..."
        collect_response_time_metrics "${API_BASE_URL}/api/v1/unified-attributes" "unified_attributes"
        collect_response_time_metrics "${API_BASE_URL}/api/v1/marketplace/search" "search"
        
        log_info "Collecting system metrics..."
        collect_error_rate
        collect_throughput_metrics
        collect_cache_metrics
        collect_resource_metrics
        
        log_info "Collecting database metrics..."
        collect_database_metrics
        
        log_info "Collecting business metrics..."
        collect_business_metrics
        
        log_info "Analyzing logs..."
        analyze_logs
        
        # Analysis
        compare_with_baseline
        
        # Generate reports
        if [[ $((check_count % 12)) -eq 0 ]]; then
            # Generate hourly report
            generate_hourly_report
            update_baseline
        fi
        
        # Alert on anomalies
        if [[ ${#anomalies[@]} -gt 0 ]]; then
            log_warning "Anomalies detected: ${!anomalies[*]}"
            send_alert "warning" "⚠️ Anomalies detected: ${!anomalies[*]}"
        else
            log_success "All metrics within normal ranges"
        fi
        
        # Sleep until next check
        log_info "Next check in ${CHECK_INTERVAL_SECONDS} seconds..."
        sleep "$CHECK_INTERVAL_SECONDS"
    done
    
    log_success "Monitoring period completed"
}

generate_final_report() {
    local final_report="$REPORT_DIR/post-deployment-summary-$(date +%Y%m%d).md"
    
    cat > "$final_report" <<EOF
# Post-Deployment Monitoring Summary
## Unified Attributes System

**Date:** $(date +"%Y-%m-%d")
**Duration:** ${MONITORING_DURATION_HOURS} hours
**Status:** ${#anomalies[@]} anomalies detected

## Key Metrics Summary

### Performance
- Response Time (p95): ${current_metrics[unified_attributes_p95]:-N/A}ms
- Throughput: ${current_metrics[throughput]:-N/A} req/s
- Error Rate: ${current_metrics[error_rate]:-N/A}

### Cache Performance
- Hit Rate: ${current_metrics[cache_hit_rate]:-N/A}%

### Resource Usage
- CPU: ${current_metrics[cpu_usage]:-N/A}%
- Memory: ${current_metrics[memory_usage_mb]:-N/A}MB

### Database
- Active Connections: ${current_metrics[db_connections]:-N/A}
- Query Time: ${current_metrics[db_query_time_ms]:-N/A}ms
- Replication Lag: ${current_metrics[replication_lag]:-N/A}s

### Business Metrics
- Listings Created: ${current_metrics[listings_created_hourly]:-N/A}/hour
- Searches: ${current_metrics[searches_hourly]:-N/A}/hour
- Attribute Usage: ${current_metrics[attribute_usage_hourly]:-N/A}/hour

## Anomalies Detected

$(if [[ ${#anomalies[@]} -eq 0 ]]; then
    echo "No anomalies detected during monitoring period."
else
    for key in "${!anomalies[@]}"; do
        echo "- $key: ${anomalies[$key]}"
    done
fi)

## Recommendations

$(if [[ ${#anomalies[@]} -eq 0 ]]; then
    echo "- System is performing well, continue monitoring"
    echo "- Consider increasing cache TTL for better performance"
    echo "- Review and optimize slow database queries"
else
    echo "- Investigate and resolve detected anomalies"
    echo "- Consider scaling resources if usage is high"
    echo "- Review error logs for root cause analysis"
fi)

## Next Steps

1. Continue monitoring for next 24 hours
2. Analyze trends and patterns
3. Implement performance optimizations
4. Update documentation with findings

---
*Generated: $(date)*
EOF
    
    log_success "Final report generated: $final_report"
    send_alert "info" "✅ Post-deployment monitoring complete. Report: $final_report"
}

# Main execution
main() {
    log_info "=========================================="
    log_info "Post-Deployment Monitoring System"
    log_info "Version: 1.0.0"
    log_info "Duration: ${MONITORING_DURATION_HOURS} hours"
    log_info "=========================================="
    
    # Initial health check
    if ! curl -sf "${API_BASE_URL}/health" >/dev/null 2>&1; then
        log_error "API health check failed - system may not be operational"
        exit 1
    fi
    
    # Run continuous monitoring
    continuous_monitoring
    
    # Generate final report
    generate_final_report
    
    log_success "Post-deployment monitoring completed successfully"
}

# Run if not sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi