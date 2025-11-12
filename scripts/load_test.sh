#!/bin/bash

set -euo pipefail

# Configuration
HOST="localhost:50051"
PROTO_PATH="/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto"
GHZ="/home/dim/go/bin/ghz"
DURATION="10m"
TARGET_RPS="10000"
CONNECTIONS="100"
WORKERS="50"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Listings Microservice Load Test"
echo "=========================================="
echo "Target: $TARGET_RPS RPS"
echo "Duration: $DURATION"
echo "Host: $HOST"
echo ""

# Check if service is running
if ! nc -z localhost 50051; then
    echo -e "${RED}ERROR: Service not running on port 50051${NC}"
    exit 1
fi

# Check if ghz is installed
if [ ! -f "$GHZ" ]; then
    echo -e "${RED}ERROR: ghz not found at $GHZ${NC}"
    exit 1
fi

# Create results directory
RESULTS_DIR="/tmp/load_test_results_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$RESULTS_DIR"

echo "Results will be saved to: $RESULTS_DIR"
echo ""

# Test 1: GetProductStats (Read-Heavy - 40% of traffic)
echo -e "${GREEN}[1/6] Testing GetProductStats (Read-Heavy)${NC}"
echo "Expected load: 4,000 RPS"
$GHZ --insecure \
    --proto="$PROTO_PATH" \
    --call=listings.v1.ListingsService.GetProductStats \
    -d '{"storefront_id": 1}' \
    -c $CONNECTIONS \
    -z 1m \
    -t 5s \
    --rps 4000 \
    -o "$RESULTS_DIR/get_product_stats.json" \
    "$HOST"

echo ""

# Test 2: IncrementProductViews (Write-Light - 5% of traffic)
echo -e "${GREEN}[2/6] Testing IncrementProductViews (Write-Light)${NC}"
echo "Expected load: 500 RPS"
$GHZ --insecure \
    --proto="$PROTO_PATH" \
    --call=listings.v1.ListingsService.IncrementProductViews \
    -d '{"product_id": 1}' \
    -c $CONNECTIONS \
    -z 1m \
    -t 3s \
    --rps 500 \
    -o "$RESULTS_DIR/increment_views.json" \
    "$HOST"

echo ""

# Test 3: RecordInventoryMovement (Write-Heavy - 5% of traffic)
echo -e "${GREEN}[3/6] Testing RecordInventoryMovement (Write-Heavy)${NC}"
echo "Expected load: 500 RPS"
$GHZ --insecure \
    --proto="$PROTO_PATH" \
    --call=listings.v1.ListingsService.RecordInventoryMovement \
    -d '{"storefront_id": 1, "product_id": 1, "movement_type": "adjustment", "quantity": 10, "reason": "manual_adjustment", "user_id": 1}' \
    -c $CONNECTIONS \
    -z 1m \
    -t 8s \
    --rps 500 \
    -o "$RESULTS_DIR/record_movement.json" \
    "$HOST"

echo ""

# Test 4: GetProduct (Read-Heavy - 20% of traffic)
echo -e "${GREEN}[4/6] Testing GetProduct (Read-Heavy)${NC}"
echo "Expected load: 2,000 RPS"
$GHZ --insecure \
    --proto="$PROTO_PATH" \
    --call=listings.v1.ListingsService.GetProduct \
    -d '{"product_id": 1, "storefront_id": 1}' \
    -c $CONNECTIONS \
    -z 1m \
    -t 5s \
    --rps 2000 \
    -o "$RESULTS_DIR/get_product.json" \
    "$HOST"

echo ""

# Test 5: CheckStockAvailability (Read-Heavy - 10% of traffic)
echo -e "${GREEN}[5/6] Testing CheckStockAvailability (Critical Path)${NC}"
echo "Expected load: 1,000 RPS"
$GHZ --insecure \
    --proto="$PROTO_PATH" \
    --call=listings.v1.ListingsService.CheckStockAvailability \
    -d '{"items": [{"product_id": 1, "quantity": 5}]}' \
    -c $CONNECTIONS \
    -z 1m \
    -t 3s \
    --rps 1000 \
    -o "$RESULTS_DIR/check_stock.json" \
    "$HOST"

echo ""

# Test 6: Sustained Mixed Load (10 minutes at 10k RPS)
echo -e "${GREEN}[6/6] Sustained Load Test (10 minutes at 10k RPS)${NC}"
echo "This will take 10 minutes..."
echo ""

# We'll use GetProductStats as representative of mixed load
$GHZ --insecure \
    --proto="$PROTO_PATH" \
    --call=listings.v1.ListingsService.GetProductStats \
    -d '{"storefront_id": 1}' \
    -c $CONNECTIONS \
    -z 10m \
    -t 5s \
    --rps $TARGET_RPS \
    -o "$RESULTS_DIR/sustained_10min.json" \
    "$HOST"

# Generate summary report
echo ""
echo "=========================================="
echo "Load Test Summary"
echo "=========================================="
echo ""

# Use Python to parse JSON results
python3 << 'EOF'
import json
import glob
import sys
import os

results_dir = sys.argv[1]
files = sorted(glob.glob(f"{results_dir}/*.json"))

print("\n| Test | Total | Success | Failed | p50 | p95 | p99 | RPS | Error% |")
print("|------|-------|---------|--------|-----|-----|-----|-----|--------|")

failed_tests = []

for file in files:
    with open(file) as f:
        data = json.load(f)
        test_name = os.path.basename(file).replace('.json', '')
        total = data['count']
        status_dist = data.get('statusCodeDist', {})
        success = status_dist.get('OK', 0)
        failed = total - success
        error_rate = (failed / total * 100) if total > 0 else 0

        # Latency in ms
        p50 = data['latencyDistribution'][5]['latency'] / 1e6 if len(data['latencyDistribution']) > 5 else 0
        p95 = data['latencyDistribution'][9]['latency'] / 1e6 if len(data['latencyDistribution']) > 9 else 0
        p99 = data['latencyDistribution'][10]['latency'] / 1e6 if len(data['latencyDistribution']) > 10 else 0
        rps = data['rps']

        # SLA validation
        sla_pass = True
        if 'sustained' in test_name and rps < 10000:
            sla_pass = False
        if p95 >= 100:
            sla_pass = False
        if error_rate >= 0.1:
            sla_pass = False

        if not sla_pass:
            failed_tests.append(test_name)

        status_icon = "✅" if sla_pass else "❌"

        print(f"| {status_icon} {test_name} | {total:,} | {success:,} | {failed} | {p50:.1f}ms | {p95:.1f}ms | {p99:.1f}ms | {rps:.0f} | {error_rate:.3f}% |")

print("\n" + "="*80)
if failed_tests:
    print(f"❌ SLA FAILED for: {', '.join(failed_tests)}")
    sys.exit(1)
else:
    print("✅ All tests passed SLA requirements!")

print("\n📊 Full results available in:", results_dir)
EOF

SUMMARY_EXIT_CODE=$?

echo ""
echo "=========================================="
echo "Next Steps"
echo "=========================================="
echo ""
echo "1. Analyze detailed results:"
echo "   python3 scripts/analyze_results.py $RESULTS_DIR/sustained_10min.json"
echo ""
echo "2. Check service logs:"
echo "   journalctl -u listings.service --since '20 minutes ago' | grep -E '(ERROR|WARN)'"
echo ""
echo "3. View resource usage:"
echo "   cat /tmp/resource_monitoring.csv"
echo ""

if [ $SUMMARY_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ Load test complete - SLA requirements met!${NC}"
    exit 0
else
    echo -e "${RED}❌ Load test complete - SLA requirements NOT met!${NC}"
    exit 1
fi
