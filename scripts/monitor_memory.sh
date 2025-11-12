#!/bin/bash

# Monitor memory usage over time

INTERVAL=10  # seconds
OUTPUT="/tmp/memory_monitoring_$(date +%Y%m%d_%H%M%S).csv"
PPROF_PORT=6060
METRICS_PORT=8080

echo "timestamp,heap_alloc_mb,heap_sys_mb,num_gc,goroutines,db_connections_open,db_connections_in_use" > "$OUTPUT"

echo "=== Memory Monitoring Tool ==="
echo "Monitoring memory every ${INTERVAL}s (press Ctrl+C to stop)"
echo "Output: $OUTPUT"
echo ""

# Trap Ctrl+C to show final summary
trap 'echo ""; echo "Monitoring stopped. Results saved to: $OUTPUT"; exit 0' INT

while true; do
    timestamp=$(date +%s)

    # Get memory stats from pprof
    stats=$(curl -s http://localhost:$PPROF_PORT/debug/pprof/heap?debug=1 2>/dev/null | head -50)

    heap_alloc=$(echo "$stats" | grep "HeapAlloc = " | sed 's/.*= //' | sed 's/ bytes.*//' | awk '{print $1/1024/1024}')
    heap_sys=$(echo "$stats" | grep "HeapSys = " | sed 's/.*= //' | sed 's/ bytes.*//' | awk '{print $1/1024/1024}')
    num_gc=$(echo "$stats" | grep "NumGC = " | sed 's/.*= //' | awk '{print $1}')

    # Default values if parsing fails
    [ -z "$heap_alloc" ] && heap_alloc=0
    [ -z "$heap_sys" ] && heap_sys=0
    [ -z "$num_gc" ] && num_gc=0

    # Get goroutine count from pprof
    goroutines=$(curl -s http://localhost:$PPROF_PORT/debug/pprof/goroutine?debug=1 2>/dev/null | grep "^goroutine profile:" | awk '{print $3}')
    [ -z "$goroutines" ] && goroutines=0

    # Get DB connection stats from metrics endpoint
    db_conn_open=$(curl -s http://localhost:$METRICS_PORT/metrics 2>/dev/null | grep "^listings_db_connections_open " | awk '{print $2}')
    db_conn_in_use=$(curl -s http://localhost:$METRICS_PORT/metrics 2>/dev/null | grep "^listings_db_connections_in_use " | awk '{print $2}')

    [ -z "$db_conn_open" ] && db_conn_open=0
    [ -z "$db_conn_in_use" ] && db_conn_in_use=0

    # Save to CSV
    echo "$timestamp,$heap_alloc,$heap_sys,$num_gc,$goroutines,$db_conn_open,$db_conn_in_use" >> "$OUTPUT"

    # Display current values
    printf "\r[%s] Heap: %6.1f MB | Sys: %6.1f MB | Goroutines: %4d | DB: %2d/%2d | GC: %d  " \
        "$(date +%H:%M:%S)" "$heap_alloc" "$heap_sys" "$goroutines" "$db_conn_in_use" "$db_conn_open" "$num_gc"

    sleep $INTERVAL
done
