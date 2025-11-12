#!/bin/bash
#
# Resource Monitoring Script for Load Testing
#
# Monitors system resources during load tests:
# - CPU usage (%)
# - Memory usage (MB)
# - Database connections
# - Goroutines count
# - gRPC requests/sec
#
# Usage: ./monitor_resources.sh [output_file] [interval_seconds]
#

set -euo pipefail

# Configuration
OUTPUT="${1:-/tmp/resource_monitoring.csv}"
INTERVAL="${2:-5}"  # seconds
SERVICE_NAME="listings"
METRICS_URL="http://localhost:8086/metrics"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Starting resource monitoring...${NC}"
echo "Output file: $OUTPUT"
echo "Interval: ${INTERVAL}s"
echo "Press Ctrl+C to stop"
echo ""

# Create CSV header
echo "timestamp,datetime,cpu_percent,mem_mb,db_connections,goroutines,grpc_rps,grpc_errors" > "$OUTPUT"

# Function to get process ID
get_pid() {
    pgrep -f "$SERVICE_NAME" || echo "0"
}

# Function to get metrics from Prometheus endpoint
get_metric() {
    local metric_name="$1"
    curl -s "$METRICS_URL" 2>/dev/null | grep "^${metric_name}" | grep -v "^#" | awk '{print $2}' | head -1 || echo "0"
}

# Function to calculate rate
declare -A prev_values
declare -A prev_times

get_rate() {
    local metric_name="$1"
    local current_value=$(get_metric "$metric_name")
    local current_time=$(date +%s)

    if [ -z "${prev_values[$metric_name]:-}" ]; then
        prev_values[$metric_name]=$current_value
        prev_times[$metric_name]=$current_time
        echo "0"
        return
    fi

    local prev_value=${prev_values[$metric_name]}
    local prev_time=${prev_times[$metric_name]}
    local time_diff=$((current_time - prev_time))

    if [ "$time_diff" -eq 0 ]; then
        echo "0"
        return
    fi

    local value_diff=$(echo "$current_value - $prev_value" | bc)
    local rate=$(echo "$value_diff / $time_diff" | bc)

    prev_values[$metric_name]=$current_value
    prev_times[$metric_name]=$current_time

    echo "$rate"
}

# Trap Ctrl+C to gracefully exit
trap 'echo -e "\n${GREEN}Monitoring stopped. Results saved to: $OUTPUT${NC}"; exit 0' INT

# Main monitoring loop
while true; do
    timestamp=$(date +%s)
    datetime=$(date '+%Y-%m-%d %H:%M:%S')

    # Get process ID
    pid=$(get_pid)

    if [ "$pid" == "0" ]; then
        echo "WARNING: Process '$SERVICE_NAME' not found!"
        cpu=0
        mem=0
    else
        # CPU usage (%)
        cpu=$(ps -p "$pid" -o %cpu | tail -1 | tr -d ' ')

        # Memory usage (MB)
        mem=$(ps -p "$pid" -o rss | tail -1 | awk '{print $1/1024}')
    fi

    # Database connections (from Prometheus metrics)
    db_conn=$(get_metric "listings_db_connections_open")

    # Goroutines count
    goroutines=$(get_metric "go_goroutines")

    # gRPC requests per second (rate)
    grpc_rps=$(get_rate "grpc_server_handled_total")

    # gRPC errors per second
    grpc_errors=$(get_metric "grpc_server_handled_total{grpc_code!=\"OK\"}" | awk '{sum+=$1} END {print sum}')
    grpc_errors=${grpc_errors:-0}

    # Write to CSV
    echo "$timestamp,$datetime,$cpu,$mem,$db_conn,$goroutines,$grpc_rps,$grpc_errors" >> "$OUTPUT"

    # Display current values
    printf "\r[%s] CPU: %5.1f%% | Mem: %7.1f MB | DB: %3d | Goroutines: %5d | RPS: %5d | Errors: %3d" \
        "$datetime" "$cpu" "$mem" "$db_conn" "$goroutines" "$grpc_rps" "$grpc_errors"

    sleep "$INTERVAL"
done
