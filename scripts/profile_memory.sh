#!/bin/bash

set -euo pipefail

# Configuration
SERVICE_HOST="localhost"
PPROF_PORT="6060"
DURATION="60s"  # Profile for 1 minute
OUTPUT_DIR="/tmp/memory_profiles_$(date +%Y%m%d_%H%M%S)"

mkdir -p "$OUTPUT_DIR"

echo "=== Memory Leak Detection Tool ==="
echo "Output directory: $OUTPUT_DIR"
echo ""

# Check if service is running
if ! nc -z $SERVICE_HOST $PPROF_PORT 2>/dev/null; then
    echo "ERROR: pprof server not accessible on $SERVICE_HOST:$PPROF_PORT"
    echo "Make sure the listings service is running with pprof enabled"
    exit 1
fi

# 1. Baseline heap profile
echo "[1/7] Capturing baseline heap profile..."
curl -s "http://$SERVICE_HOST:$PPROF_PORT/debug/pprof/heap" > "$OUTPUT_DIR/heap_baseline.pprof"

# 2. Baseline goroutine profile
echo "[2/7] Capturing baseline goroutine profile..."
curl -s "http://$SERVICE_HOST:$PPROF_PORT/debug/pprof/goroutine" > "$OUTPUT_DIR/goroutine_baseline.pprof"

# 3. Force GC
echo "[3/7] Forcing garbage collection..."
curl -s "http://$SERVICE_HOST:$PPROF_PORT/debug/pprof/heap?gc=1" > /dev/null

# 4. Wait and generate load
echo "[4/7] Generating load for $DURATION..."

# Check if ghz is available
if ! command -v ghz &> /dev/null; then
    echo "WARNING: ghz not found, using curl-based load instead"
    # Fallback to curl-based load
    for i in {1..1000}; do
        curl -s http://localhost:8080/health > /dev/null &
    done
    wait
    sleep 60
else
    # Use ghz for gRPC load testing
    ghz --insecure \
        --proto="/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto" \
        --call=listings.v1.ListingsService.GetListing \
        -d '{"id": 328}' \
        -c 50 \
        -z $DURATION \
        --rps 1000 \
        localhost:50051 > /dev/null 2>&1 || true
fi

echo "[4/7] Load generation complete, waiting for stabilization..."
sleep 5

# 5. Post-load heap profile
echo "[5/7] Capturing post-load heap profile..."
curl -s "http://$SERVICE_HOST:$PPROF_PORT/debug/pprof/heap" > "$OUTPUT_DIR/heap_postload.pprof"

# 6. Post-load goroutine profile
echo "[6/7] Capturing post-load goroutine profile..."
curl -s "http://$SERVICE_HOST:$PPROF_PORT/debug/pprof/goroutine" > "$OUTPUT_DIR/goroutine_postload.pprof"

# 7. Analyze differences
echo "[7/7] Analyzing memory differences..."

go tool pprof -base="$OUTPUT_DIR/heap_baseline.pprof" \
    -top "$OUTPUT_DIR/heap_postload.pprof" > "$OUTPUT_DIR/heap_diff.txt" 2>/dev/null || true

go tool pprof -base="$OUTPUT_DIR/goroutine_baseline.pprof" \
    -top "$OUTPUT_DIR/goroutine_postload.pprof" > "$OUTPUT_DIR/goroutine_diff.txt" 2>/dev/null || true

echo ""
echo "=== Analysis Complete ==="
echo ""

# Display top memory allocations
echo "Top Memory Allocations:"
if [ -f "$OUTPUT_DIR/heap_diff.txt" ]; then
    head -20 "$OUTPUT_DIR/heap_diff.txt"
else
    echo "(No significant heap changes detected)"
fi

echo ""
echo "Goroutine Growth:"
if [ -f "$OUTPUT_DIR/goroutine_diff.txt" ]; then
    head -20 "$OUTPUT_DIR/goroutine_diff.txt"
else
    echo "(No significant goroutine changes detected)"
fi

echo ""
echo "Full reports available in: $OUTPUT_DIR"
echo ""

# Generate visualizations (if go-torch is available)
if command -v go-torch &> /dev/null; then
    echo "Generating flamegraphs..."
    go-torch --file="$OUTPUT_DIR/heap_flamegraph.svg" \
        --url="http://$SERVICE_HOST:$PPROF_PORT" \
        -t heap 2>/dev/null || echo "Flamegraph generation skipped"
    if [ -f "$OUTPUT_DIR/heap_flamegraph.svg" ]; then
        echo "Flamegraph: $OUTPUT_DIR/heap_flamegraph.svg"
    fi
fi

echo ""
echo "✅ Memory profiling complete!"
echo ""
echo "Next steps:"
echo "  1. Review heap diff: cat $OUTPUT_DIR/heap_diff.txt"
echo "  2. Review goroutine diff: cat $OUTPUT_DIR/goroutine_diff.txt"
echo "  3. Interactive analysis: go tool pprof $OUTPUT_DIR/heap_postload.pprof"
echo "  4. Compare profiles: go tool pprof -base $OUTPUT_DIR/heap_baseline.pprof $OUTPUT_DIR/heap_postload.pprof"
