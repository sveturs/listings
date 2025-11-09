#!/bin/bash
#
# Import Grafana dashboards via API
# Usage: ./import-dashboards.sh [grafana-url] [api-key]
#
# Example:
#   ./import-dashboards.sh http://localhost:3000 eyJrIjoiXXXXXX
#

set -euo pipefail

GRAFANA_URL="${1:-http://localhost:3000}"
API_KEY="${2:-}"
DASHBOARD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/dashboards" && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=================================================="
echo "  Grafana Dashboard Import Script"
echo "=================================================="
echo ""
echo "Grafana URL: ${GRAFANA_URL}"
echo "Dashboard directory: ${DASHBOARD_DIR}"
echo ""

# Check if API key is provided
if [ -z "${API_KEY}" ]; then
    echo -e "${YELLOW}Warning: No API key provided${NC}"
    echo "You can create an API key in Grafana: Configuration -> API Keys"
    echo ""
    read -p "Continue without authentication? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
    AUTH_HEADER=""
else
    AUTH_HEADER="Authorization: Bearer ${API_KEY}"
fi

# Function to import a single dashboard
import_dashboard() {
    local file=$1
    local basename=$(basename "$file")

    echo -n "Importing ${basename}... "

    # Read the dashboard JSON and wrap it in the import format
    local dashboard_json=$(cat "$file")
    local import_json=$(jq -n \
        --argjson dashboard "$dashboard_json" \
        '{dashboard: $dashboard.dashboard, overwrite: true}')

    # Make the API call
    local response
    if [ -n "${AUTH_HEADER}" ]; then
        response=$(curl -s -X POST \
            -H "Content-Type: application/json" \
            -H "${AUTH_HEADER}" \
            -d "${import_json}" \
            "${GRAFANA_URL}/api/dashboards/db")
    else
        response=$(curl -s -X POST \
            -H "Content-Type: application/json" \
            -d "${import_json}" \
            "${GRAFANA_URL}/api/dashboards/db")
    fi

    # Check response
    if echo "$response" | jq -e '.status == "success"' > /dev/null 2>&1; then
        local url=$(echo "$response" | jq -r '.url')
        echo -e "${GREEN}✓${NC} Imported successfully"
        echo "   URL: ${GRAFANA_URL}${url}"
    else
        echo -e "${RED}✗${NC} Failed"
        echo "   Response: $response"
        return 1
    fi
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed${NC}"
    echo "Please install jq: sudo apt-get install jq"
    exit 1
fi

# Check if Grafana is reachable
echo "Checking Grafana connectivity..."
if ! curl -s -f "${GRAFANA_URL}/api/health" > /dev/null; then
    echo -e "${RED}Error: Cannot connect to Grafana at ${GRAFANA_URL}${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} Grafana is reachable"
echo ""

# Import all dashboards
success_count=0
fail_count=0

for dashboard in "${DASHBOARD_DIR}"/*.json; do
    if [ -f "$dashboard" ]; then
        if import_dashboard "$dashboard"; then
            ((success_count++))
        else
            ((fail_count++))
        fi
        echo ""
    fi
done

# Summary
echo "=================================================="
echo "  Import Summary"
echo "=================================================="
echo -e "Successfully imported: ${GREEN}${success_count}${NC}"
echo -e "Failed: ${RED}${fail_count}${NC}"
echo ""

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}All dashboards imported successfully!${NC}"
    exit 0
else
    echo -e "${YELLOW}Some dashboards failed to import${NC}"
    exit 1
fi
