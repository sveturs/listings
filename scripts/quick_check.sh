#!/bin/bash

# Quick Memory Health Check
# Performs a rapid check of service memory health without load testing

set -euo pipefail

PPROF_URL="http://localhost:6060"
METRICS_URL="http://localhost:8080/metrics"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Quick Memory Health Check ==="
echo ""

# Check if pprof is accessible
if ! curl -s -f "$PPROF_URL/debug/pprof/" > /dev/null; then
    echo -e "${RED}✗ pprof server not accessible at $PPROF_URL${NC}"
    echo "  Make sure the service is running with pprof enabled"
    exit 1
fi

echo -e "${GREEN}✓ pprof server accessible${NC}"

# Check if metrics are accessible
if ! curl -s -f "$METRICS_URL" > /dev/null; then
    echo -e "${YELLOW}⚠ Metrics endpoint not accessible at $METRICS_URL${NC}"
    echo "  Continuing without metrics..."
    METRICS_AVAILABLE=false
else
    echo -e "${GREEN}✓ Metrics endpoint accessible${NC}"
    METRICS_AVAILABLE=true
fi

echo ""
echo "=== Current Memory Stats ==="

# Get heap stats
HEAP_STATS=$(curl -s "$PPROF_URL/debug/pprof/heap?debug=1" | head -50)

HEAP_ALLOC=$(echo "$HEAP_STATS" | grep "HeapAlloc = " | sed 's/.*= //' | awk '{printf "%.1f", $1/1024/1024}')
HEAP_SYS=$(echo "$HEAP_STATS" | grep "HeapSys = " | sed 's/.*= //' | awk '{printf "%.1f", $1/1024/1024}')
HEAP_IDLE=$(echo "$HEAP_STATS" | grep "HeapIdle = " | sed 's/.*= //' | awk '{printf "%.1f", $1/1024/1024}')
HEAP_INUSE=$(echo "$HEAP_STATS" | grep "HeapInuse = " | sed 's/.*= //' | awk '{printf "%.1f", $1/1024/1024}')
NUM_GC=$(echo "$HEAP_STATS" | grep "NumGC = " | sed 's/.*= //' | awk '{print $1}')

echo "Heap Allocation: ${HEAP_ALLOC} MB"
echo "Heap System:     ${HEAP_SYS} MB"
echo "Heap In-Use:     ${HEAP_INUSE} MB"
echo "Heap Idle:       ${HEAP_IDLE} MB"
echo "Total GC Runs:   ${NUM_GC}"

# Assess heap health
if (( $(echo "$HEAP_ALLOC > 200" | bc -l) )); then
    echo -e "${RED}⚠ WARNING: High heap allocation (> 200 MB)${NC}"
elif (( $(echo "$HEAP_ALLOC > 100" | bc -l) )); then
    echo -e "${YELLOW}⚠ NOTICE: Moderate heap allocation (> 100 MB)${NC}"
else
    echo -e "${GREEN}✓ Heap allocation normal${NC}"
fi

echo ""
echo "=== Goroutines ==="

# Get goroutine count
GOROUTINE_COUNT=$(curl -s "$PPROF_URL/debug/pprof/goroutine?debug=1" | grep "^goroutine profile:" | awk '{print $3}')

echo "Active Goroutines: ${GOROUTINE_COUNT}"

# Assess goroutine count
if [ "$GOROUTINE_COUNT" -gt 200 ]; then
    echo -e "${RED}⚠ WARNING: High goroutine count (> 200)${NC}"
elif [ "$GOROUTINE_COUNT" -gt 100 ]; then
    echo -e "${YELLOW}⚠ NOTICE: Moderate goroutine count (> 100)${NC}"
else
    echo -e "${GREEN}✓ Goroutine count normal${NC}"
fi

# Show top goroutine stacks
echo ""
echo "Top 5 Goroutine States:"
curl -s "$PPROF_URL/debug/pprof/goroutine?debug=1" | grep -A 1 "^goroutine " | head -10

echo ""
echo "=== Database Connections ==="

if [ "$METRICS_AVAILABLE" = true ]; then
    DB_OPEN=$(curl -s "$METRICS_URL" | grep "^listings_db_connections_open " | awk '{print $2}')
    DB_INUSE=$(curl -s "$METRICS_URL" | grep "^listings_db_connections_in_use " | awk '{print $2}')
    DB_IDLE=$(curl -s "$METRICS_URL" | grep "^listings_db_connections_idle " | awk '{print $2}')
    DB_WAIT=$(curl -s "$METRICS_URL" | grep "^listings_db_connections_wait_count " | awk '{print $2}')

    echo "Open Connections:   ${DB_OPEN:-N/A}"
    echo "In-Use Connections: ${DB_INUSE:-N/A}"
    echo "Idle Connections:   ${DB_IDLE:-N/A}"
    echo "Wait Count:         ${DB_WAIT:-N/A}"

    # Assess DB connection health
    if [ -n "$DB_OPEN" ] && [ "$DB_OPEN" -gt 45 ]; then
        echo -e "${RED}⚠ WARNING: High DB connection count (> 45/50 max)${NC}"
    elif [ -n "$DB_OPEN" ] && [ "$DB_OPEN" -gt 30 ]; then
        echo -e "${YELLOW}⚠ NOTICE: Moderate DB connection count (> 30)${NC}"
    else
        echo -e "${GREEN}✓ DB connection count normal${NC}"
    fi
else
    echo -e "${YELLOW}Metrics not available${NC}"
fi

echo ""
echo "=== GC Performance ==="

# Get GC stats
GC_PAUSE=$(echo "$HEAP_STATS" | grep "PauseNs" | head -1 | sed 's/.*: //' | awk '{printf "%.2f", $1/1000000}')

echo "Last GC Pause: ${GC_PAUSE:-N/A} ms"

# Assess GC performance
if [ -n "$GC_PAUSE" ]; then
    if (( $(echo "$GC_PAUSE > 100" | bc -l) )); then
        echo -e "${RED}⚠ WARNING: High GC pause time (> 100ms)${NC}"
    elif (( $(echo "$GC_PAUSE > 50" | bc -l) )); then
        echo -e "${YELLOW}⚠ NOTICE: Moderate GC pause time (> 50ms)${NC}"
    else
        echo -e "${GREEN}✓ GC pause time normal${NC}"
    fi
fi

echo ""
echo "=== Runtime Info ==="

# Get runtime info from /debug/pprof/
CMDLINE=$(curl -s "$PPROF_URL/debug/pprof/cmdline")
echo "Command: $CMDLINE"

# Get uptime estimate from GC count (approximate)
if [ "$NUM_GC" -gt 0 ]; then
    # Assume GC runs every ~5 minutes on average
    UPTIME_HOURS=$(echo "$NUM_GC * 5 / 60" | bc -l | awk '{printf "%.1f", $1}')
    echo "Estimated Uptime: ~${UPTIME_HOURS} hours (based on GC runs)"
fi

echo ""
echo "=== Health Summary ==="

# Overall health assessment
WARNINGS=0

# Check heap
if (( $(echo "$HEAP_ALLOC > 200" | bc -l) )); then
    ((WARNINGS++))
fi

# Check goroutines
if [ "$GOROUTINE_COUNT" -gt 200 ]; then
    ((WARNINGS++))
fi

# Check DB connections
if [ "$METRICS_AVAILABLE" = true ] && [ -n "$DB_OPEN" ] && [ "$DB_OPEN" -gt 45 ]; then
    ((WARNINGS++))
fi

if [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ Service health: GOOD${NC}"
    echo "  No issues detected"
elif [ $WARNINGS -eq 1 ]; then
    echo -e "${YELLOW}⚠ Service health: FAIR${NC}"
    echo "  $WARNINGS issue detected - monitor closely"
else
    echo -e "${RED}⚠ Service health: POOR${NC}"
    echo "  $WARNINGS issues detected - investigate immediately"
    echo ""
    echo "Recommended actions:"
    echo "  1. Run: ./scripts/profile_memory.sh"
    echo "  2. Review logs for errors"
    echo "  3. Check for memory leaks"
fi

echo ""
echo "=== Quick Actions ==="
echo ""
echo "View detailed heap profile:"
echo "  go tool pprof -top http://localhost:6060/debug/pprof/heap"
echo ""
echo "View goroutine details:"
echo "  go tool pprof -top http://localhost:6060/debug/pprof/goroutine"
echo ""
echo "Force garbage collection:"
echo "  curl http://localhost:6060/debug/pprof/heap?gc=1"
echo ""
echo "Run full leak detection:"
echo "  ./scripts/profile_memory.sh"
echo ""
echo "Start continuous monitoring:"
echo "  ./scripts/monitor_memory.sh"
echo ""
