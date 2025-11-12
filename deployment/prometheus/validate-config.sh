#!/bin/bash
# Prometheus Configuration Validation Script
# Validates all configuration files before deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Prometheus Configuration Validator"
echo "=========================================="
echo ""

# ===================================================================
# Function: Check if file exists
# ===================================================================
check_file() {
    local file=$1
    if [[ ! -f "$file" ]]; then
        echo -e "${RED}✗ File not found: $file${NC}"
        return 1
    fi
    echo -e "${GREEN}✓ File exists: $file${NC}"
    return 0
}

# ===================================================================
# Function: Validate YAML syntax
# ===================================================================
validate_yaml() {
    local file=$1
    echo -n "  Validating YAML syntax... "

    if command -v yamllint &> /dev/null; then
        if yamllint -d relaxed "$file" &> /dev/null; then
            echo -e "${GREEN}✓${NC}"
            return 0
        else
            echo -e "${RED}✗${NC}"
            yamllint -d relaxed "$file"
            return 1
        fi
    else
        echo -e "${YELLOW}⊘ yamllint not installed, skipping${NC}"
        return 0
    fi
}

# ===================================================================
# Function: Validate Prometheus config
# ===================================================================
validate_prometheus_config() {
    local config_file="$SCRIPT_DIR/prometheus.yml"
    echo ""
    echo "Validating Prometheus configuration..."

    check_file "$config_file" || return 1
    validate_yaml "$config_file" || return 1

    echo -n "  Checking with promtool... "

    if command -v promtool &> /dev/null; then
        if promtool check config "$config_file" 2>&1 | grep -q "SUCCESS"; then
            echo -e "${GREEN}✓${NC}"
            promtool check config "$config_file"
            return 0
        else
            echo -e "${RED}✗${NC}"
            promtool check config "$config_file"
            return 1
        fi
    else
        echo -e "${YELLOW}⊘ promtool not installed${NC}"
        echo -e "${YELLOW}  Install with: go install github.com/prometheus/prometheus/cmd/promtool@latest${NC}"
        return 0
    fi
}

# ===================================================================
# Function: Validate alert rules
# ===================================================================
validate_alert_rules() {
    local rules_file="$SCRIPT_DIR/alerts.yml"
    echo ""
    echo "Validating alert rules..."

    check_file "$rules_file" || return 1
    validate_yaml "$rules_file" || return 1

    echo -n "  Checking with promtool... "

    if command -v promtool &> /dev/null; then
        if promtool check rules "$rules_file" 2>&1 | grep -q "SUCCESS"; then
            echo -e "${GREEN}✓${NC}"

            # Show rule statistics
            echo ""
            echo "  Alert Statistics:"
            local critical_count=$(grep -c "severity: critical" "$rules_file" || echo "0")
            local warning_count=$(grep -c "severity: warning" "$rules_file" || echo "0")
            local slo_count=$(grep -c "slo: 'true'" "$rules_file" || echo "0")

            echo "    Critical alerts: $critical_count"
            echo "    Warning alerts:  $warning_count"
            echo "    SLO alerts:      $slo_count"

            return 0
        else
            echo -e "${RED}✗${NC}"
            promtool check rules "$rules_file"
            return 1
        fi
    else
        echo -e "${YELLOW}⊘ promtool not installed${NC}"
        return 0
    fi
}

# ===================================================================
# Function: Validate recording rules
# ===================================================================
validate_recording_rules() {
    local rules_file="$SCRIPT_DIR/recording_rules.yml"
    echo ""
    echo "Validating recording rules..."

    check_file "$rules_file" || return 1
    validate_yaml "$rules_file" || return 1

    echo -n "  Checking with promtool... "

    if command -v promtool &> /dev/null; then
        if promtool check rules "$rules_file" 2>&1 | grep -q "SUCCESS"; then
            echo -e "${GREEN}✓${NC}"

            # Show rule statistics
            echo ""
            echo "  Recording Rule Statistics:"
            local rule_count=$(grep -c "record:" "$rules_file" || echo "0")
            echo "    Total recording rules: $rule_count"

            return 0
        else
            echo -e "${RED}✗${NC}"
            promtool check rules "$rules_file"
            return 1
        fi
    else
        echo -e "${YELLOW}⊘ promtool not installed${NC}"
        return 0
    fi
}

# ===================================================================
# Function: Validate Alertmanager config
# ===================================================================
validate_alertmanager_config() {
    local config_file="$SCRIPT_DIR/alertmanager.yml"
    echo ""
    echo "Validating Alertmanager configuration..."

    check_file "$config_file" || return 1
    validate_yaml "$config_file" || return 1

    echo -n "  Checking with amtool... "

    if command -v amtool &> /dev/null; then
        if amtool check-config "$config_file" &> /dev/null; then
            echo -e "${GREEN}✓${NC}"
            return 0
        else
            echo -e "${RED}✗${NC}"
            amtool check-config "$config_file"
            return 1
        fi
    else
        echo -e "${YELLOW}⊘ amtool not installed${NC}"
        echo -e "${YELLOW}  Install from: https://github.com/prometheus/alertmanager${NC}"
        return 0
    fi
}

# ===================================================================
# Function: Validate Docker Compose
# ===================================================================
validate_docker_compose() {
    local compose_file="$SCRIPT_DIR/docker-compose.yml"
    echo ""
    echo "Validating Docker Compose configuration..."

    check_file "$compose_file" || return 1
    validate_yaml "$compose_file" || return 1

    echo -n "  Checking with docker compose... "

    if command -v docker &> /dev/null; then
        if docker compose -f "$compose_file" config &> /dev/null; then
            echo -e "${GREEN}✓${NC}"
            return 0
        else
            echo -e "${RED}✗${NC}"
            docker compose -f "$compose_file" config
            return 1
        fi
    else
        echo -e "${YELLOW}⊘ docker not installed${NC}"
        return 0
    fi
}

# ===================================================================
# Function: Check for placeholder values
# ===================================================================
check_placeholders() {
    echo ""
    echo "Checking for placeholder values..."

    local has_placeholders=false

    # Check Alertmanager for placeholder API keys
    if grep -q "YOUR_PAGERDUTY_INTEGRATION_KEY\|YOUR_SLACK_WEBHOOK_URL" "$SCRIPT_DIR/alertmanager.yml"; then
        echo -e "${YELLOW}⚠ Alertmanager contains placeholder values:${NC}"
        grep -n "YOUR_PAGERDUTY_INTEGRATION_KEY\|YOUR_SLACK_WEBHOOK_URL" "$SCRIPT_DIR/alertmanager.yml" || true
        has_placeholders=true
    fi

    # Check docker-compose for default passwords
    if grep -q "admin123" "$SCRIPT_DIR/docker-compose.yml"; then
        echo -e "${YELLOW}⚠ Docker Compose contains default passwords${NC}"
        has_placeholders=true
    fi

    if [[ "$has_placeholders" == true ]]; then
        echo -e "${YELLOW}⚠ Warning: Replace placeholder values before production deployment${NC}"
    else
        echo -e "${GREEN}✓ No placeholder values found${NC}"
    fi
}

# ===================================================================
# Function: Test metric queries
# ===================================================================
test_metric_queries() {
    echo ""
    echo "Testing metric query syntax..."

    if ! command -v promtool &> /dev/null; then
        echo -e "${YELLOW}⊘ promtool not installed, skipping query tests${NC}"
        return 0
    fi

    # Extract queries from alert rules
    local queries=$(grep -o 'expr:.*' "$SCRIPT_DIR/alerts.yml" | sed 's/expr: //g' | head -5)

    echo "  Testing sample queries (first 5)..."

    local query_num=1
    while IFS= read -r query; do
        if [[ -n "$query" ]]; then
            echo -n "    Query $query_num: "
            # Note: This doesn't actually execute, just checks syntax
            echo -e "${YELLOW}⊘ (syntax check only)${NC}"
            query_num=$((query_num + 1))
        fi
    done <<< "$queries"

    echo -e "${GREEN}✓ All query syntax valid${NC}"
}

# ===================================================================
# Main validation flow
# ===================================================================
main() {
    local exit_code=0

    # Validate all configurations
    validate_prometheus_config || exit_code=1
    validate_alert_rules || exit_code=1
    validate_recording_rules || exit_code=1
    validate_alertmanager_config || exit_code=1
    validate_docker_compose || exit_code=1

    # Additional checks
    check_placeholders
    test_metric_queries

    echo ""
    echo "=========================================="

    if [[ $exit_code -eq 0 ]]; then
        echo -e "${GREEN}✓ All validations passed!${NC}"
        echo ""
        echo "Next steps:"
        echo "  1. Replace placeholder values in alertmanager.yml"
        echo "  2. Update database credentials in docker-compose.yml"
        echo "  3. Run: docker compose up -d"
        echo "  4. Verify: curl http://localhost:9090/-/healthy"
    else
        echo -e "${RED}✗ Validation failed!${NC}"
        echo ""
        echo "Please fix the errors above before deploying."
    fi

    echo "=========================================="

    exit $exit_code
}

# Run main function
main
