#!/bin/bash
#
# Validate Grafana dashboard JSON files
# Usage: ./validate-dashboards.sh
#

set -euo pipefail

DASHBOARD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/dashboards" && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=================================================="
echo "  Grafana Dashboard Validation"
echo "=================================================="
echo ""

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed${NC}"
    echo "Please install jq: sudo apt-get install jq"
    exit 1
fi

success_count=0
fail_count=0
total_panels=0

# Validate a single dashboard
validate_dashboard() {
    local file=$1
    local basename=$(basename "$file")

    echo -e "${BLUE}Validating ${basename}...${NC}"

    # Check if file is valid JSON
    if ! jq empty "$file" 2>/dev/null; then
        echo -e "  ${RED}✗${NC} Invalid JSON"
        return 1
    fi
    echo -e "  ${GREEN}✓${NC} Valid JSON"

    # Check required fields
    local required_fields=("dashboard" "dashboard.title" "dashboard.uid" "dashboard.panels")
    for field in "${required_fields[@]}"; do
        if ! jq -e ".$field" "$file" > /dev/null 2>&1; then
            echo -e "  ${RED}✗${NC} Missing required field: $field"
            return 1
        fi
    done
    echo -e "  ${GREEN}✓${NC} All required fields present"

    # Extract dashboard info
    local title=$(jq -r '.dashboard.title' "$file")
    local uid=$(jq -r '.dashboard.uid' "$file")
    local panel_count=$(jq '.dashboard.panels | length' "$file")
    local tags=$(jq -r '.dashboard.tags | join(", ")' "$file")

    echo -e "  ${GREEN}✓${NC} Title: ${title}"
    echo -e "  ${GREEN}✓${NC} UID: ${uid}"
    echo -e "  ${GREEN}✓${NC} Panels: ${panel_count}"
    echo -e "  ${GREEN}✓${NC} Tags: ${tags}"

    total_panels=$((total_panels + panel_count))

    # Validate panels
    local invalid_panels=0
    for i in $(seq 0 $((panel_count - 1))); do
        local panel_id=$(jq -r ".dashboard.panels[$i].id" "$file")
        local panel_title=$(jq -r ".dashboard.panels[$i].title" "$file")
        local panel_type=$(jq -r ".dashboard.panels[$i].type" "$file")

        if [ "$panel_id" == "null" ] || [ "$panel_title" == "null" ]; then
            echo -e "  ${YELLOW}⚠${NC} Panel $i has missing id or title"
            ((invalid_panels++))
        fi
    done

    if [ $invalid_panels -gt 0 ]; then
        echo -e "  ${YELLOW}⚠${NC} Found $invalid_panels panels with issues"
    else
        echo -e "  ${GREEN}✓${NC} All panels valid"
    fi

    # Check for template variables
    local var_count=$(jq '.dashboard.templating.list | length' "$file" 2>/dev/null || echo "0")
    if [ "$var_count" -gt 0 ]; then
        echo -e "  ${GREEN}✓${NC} Template variables: ${var_count}"
    fi

    # Check for annotations
    local annotation_count=$(jq '.dashboard.annotations.list | length' "$file" 2>/dev/null || echo "0")
    if [ "$annotation_count" -gt 0 ]; then
        echo -e "  ${GREEN}✓${NC} Annotations: ${annotation_count}"
    fi

    echo ""
    return 0
}

# Validate all dashboards
for dashboard in "${DASHBOARD_DIR}"/*.json; do
    if [ -f "$dashboard" ]; then
        if validate_dashboard "$dashboard"; then
            ((success_count++))
        else
            ((fail_count++))
        fi
    fi
done

# Summary
echo "=================================================="
echo "  Validation Summary"
echo "=================================================="
echo -e "Dashboards validated: ${success_count}"
echo -e "Total panels: ${total_panels}"
echo -e "Passed: ${GREEN}${success_count}${NC}"
echo -e "Failed: ${RED}${fail_count}${NC}"
echo ""

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}All dashboards are valid!${NC}"
    exit 0
else
    echo -e "${RED}Some dashboards have validation errors${NC}"
    exit 1
fi
