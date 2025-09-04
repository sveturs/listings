#!/bin/bash
# Final Validation Script - Day 15
# Post-deployment validation for Unified Attributes System
# Version: 1.0.0
# Date: 03.09.2025

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
API_BASE_URL="https://api.svetu.rs/api/v1"
HEALTH_URL="https://api.svetu.rs/health"
METRICS_URL="https://prometheus.svetu.rs/api/v1"
DB_HOST="prod-db-master.svetu.rs"
DB_NAME="svetubd"
DB_USER="postgres"
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=require"

# Validation results
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
WARNINGS=0

# Result tracking
declare -A validation_results

# Logging functions
log_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
}

log_pass() {
    echo -e "${GREEN}[✓]${NC} $1"
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
    validation_results["$1"]="PASSED"
}

log_fail() {
    echo -e "${RED}[✗]${NC} $1"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
    validation_results["$1"]="FAILED"
}

log_warn() {
    echo -e "${YELLOW}[!]${NC} $1"
    WARNINGS=$((WARNINGS + 1))
}

# API validation functions
validate_api_health() {
    log_test "API Health Check"
    
    if curl -sf "$HEALTH_URL" >/dev/null 2>&1; then
        log_pass "API is healthy"
    else
        log_fail "API health check failed"
        return 1
    fi
}

validate_unified_attributes_endpoint() {
    log_test "Unified Attributes Endpoint"
    
    response=$(curl -sf "${API_BASE_URL}/unified-attributes" 2>/dev/null || echo "")
    
    if [[ -n "$response" ]]; then
        log_pass "Unified attributes endpoint responding"
    else
        log_fail "Unified attributes endpoint not responding"
        return 1
    fi
}

validate_attribute_search() {
    log_test "Attribute Search Functionality"
    
    response=$(curl -sf "${API_BASE_URL}/unified-attributes/search?category_id=1" 2>/dev/null || echo "")
    
    if [[ -n "$response" ]] && [[ "$response" != *"error"* ]]; then
        log_pass "Attribute search working"
    else
        log_fail "Attribute search not working"
        return 1
    fi
}

validate_listing_creation_with_attributes() {
    log_test "Listing Creation with Unified Attributes"
    
    # Create test listing with attributes
    response=$(curl -sf -X POST "${API_BASE_URL}/marketplace/listings" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TEST_AUTH_TOKEN:-test}" \
        -d '{
            "title": "Test Listing",
            "description": "Validation test",
            "category_id": 1,
            "price": 100,
            "unified_attributes": {
                "brand": "TestBrand",
                "model": "TestModel",
                "year": 2024
            }
        }' 2>/dev/null || echo "")
    
    if [[ "$response" == *"id"* ]]; then
        log_pass "Listing creation with unified attributes successful"
        
        # Extract listing ID for cleanup
        listing_id=$(echo "$response" | jq -r '.data.id' 2>/dev/null || echo "")
        if [[ -n "$listing_id" ]]; then
            # Delete test listing
            curl -sf -X DELETE "${API_BASE_URL}/marketplace/listings/${listing_id}" \
                -H "Authorization: Bearer ${TEST_AUTH_TOKEN:-test}" >/dev/null 2>&1
        fi
    else
        log_fail "Listing creation with unified attributes failed"
        return 1
    fi
}

# Database validation functions
validate_database_schema() {
    log_test "Database Schema Validation"
    
    # Check if unified_attributes table exists
    table_exists=$(psql "$DATABASE_URL" -tAc "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'unified_attributes')" 2>/dev/null || echo "f")
    
    if [[ "$table_exists" == "t" ]]; then
        log_pass "Unified attributes table exists"
    else
        log_fail "Unified attributes table not found"
        return 1
    fi
    
    # Check required columns
    columns=$(psql "$DATABASE_URL" -tAc "SELECT column_name FROM information_schema.columns WHERE table_name = 'unified_attributes'" 2>/dev/null || echo "")
    
    required_columns=("id" "category_id" "attribute_key" "attribute_type" "validation_rules" "display_order")
    for col in "${required_columns[@]}"; do
        if echo "$columns" | grep -q "$col"; then
            log_pass "Column '$col' exists"
        else
            log_fail "Column '$col' missing"
        fi
    done
}

validate_data_migration() {
    log_test "Data Migration Validation"
    
    # Count records in unified_attributes
    count=$(psql "$DATABASE_URL" -tAc "SELECT COUNT(*) FROM unified_attributes" 2>/dev/null || echo "0")
    
    if [[ "$count" -gt 0 ]]; then
        log_pass "Unified attributes table has $count records"
    else
        log_fail "Unified attributes table is empty"
        return 1
    fi
    
    # Check for data integrity
    integrity_check=$(psql "$DATABASE_URL" -tAc "
        SELECT COUNT(*) 
        FROM unified_attributes 
        WHERE category_id IS NULL 
           OR attribute_key IS NULL 
           OR attribute_type IS NULL
    " 2>/dev/null || echo "0")
    
    if [[ "$integrity_check" == "0" ]]; then
        log_pass "Data integrity check passed"
    else
        log_warn "Found $integrity_check records with null values"
    fi
}

validate_old_tables_archived() {
    log_test "Old Tables Archive Status"
    
    # Check if old tables are in archive schema
    archived=$(psql "$DATABASE_URL" -tAc "
        SELECT COUNT(*) 
        FROM information_schema.tables 
        WHERE table_schema = 'archive' 
          AND table_name IN ('category_attributes', 'listing_attributes', 'category_attribute_values')
    " 2>/dev/null || echo "0")
    
    if [[ "$archived" == "3" ]]; then
        log_pass "Old tables successfully archived"
    else
        log_warn "Not all old tables are archived (found $archived/3)"
    fi
}

# Performance validation
validate_performance_metrics() {
    log_test "Performance Metrics"
    
    # Check response time
    response_time=$(curl -w "%{time_total}" -o /dev/null -sf "${API_BASE_URL}/unified-attributes" 2>/dev/null || echo "999")
    response_time_ms=$(echo "$response_time * 1000" | bc)
    
    if (( $(echo "$response_time_ms < 100" | bc -l) )); then
        log_pass "Response time: ${response_time_ms}ms (< 100ms)"
    else
        log_warn "Response time: ${response_time_ms}ms (expected < 100ms)"
    fi
    
    # Check error rate from Prometheus
    error_rate=$(curl -s "${METRICS_URL}/query?query=rate(http_requests_total{status=~\"5..\"}[5m])" | \
        jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    if (( $(echo "$error_rate < 0.001" | bc -l) )); then
        log_pass "Error rate: ${error_rate} (< 0.1%)"
    else
        log_fail "Error rate: ${error_rate} (> 0.1%)"
    fi
}

validate_cache_functionality() {
    log_test "Cache Functionality"
    
    # Make first request
    time1=$(curl -w "%{time_total}" -o /dev/null -sf "${API_BASE_URL}/unified-attributes?category_id=1" 2>/dev/null || echo "999")
    
    # Make second request (should be cached)
    time2=$(curl -w "%{time_total}" -o /dev/null -sf "${API_BASE_URL}/unified-attributes?category_id=1" 2>/dev/null || echo "999")
    
    # Second request should be significantly faster
    if (( $(echo "$time2 < $time1 * 0.5" | bc -l) )); then
        log_pass "Cache is working (${time2}s < ${time1}s)"
    else
        log_warn "Cache may not be working effectively"
    fi
}

# Feature validation
validate_dual_write_disabled() {
    log_test "Dual-Write Mode Status"
    
    # Check if dual-write is disabled in configuration
    dual_write_status=$(kubectl get deployment backend-green -o json | \
        jq -r '.spec.template.spec.containers[0].env[] | select(.name=="DUAL_WRITE_ATTRIBUTES") | .value' 2>/dev/null || echo "unknown")
    
    if [[ "$dual_write_status" == "false" ]]; then
        log_pass "Dual-write mode is disabled"
    else
        log_fail "Dual-write mode is still enabled: $dual_write_status"
    fi
}

validate_fallback_disabled() {
    log_test "Fallback Mode Status"
    
    # Check if fallback is disabled
    fallback_status=$(kubectl get deployment backend-green -o json | \
        jq -r '.spec.template.spec.containers[0].env[] | select(.name=="UNIFIED_ATTRIBUTES_FALLBACK") | .value' 2>/dev/null || echo "unknown")
    
    if [[ "$fallback_status" == "false" ]]; then
        log_pass "Fallback mode is disabled"
    else
        log_fail "Fallback mode is still enabled: $fallback_status"
    fi
}

# Integration validation
validate_frontend_integration() {
    log_test "Frontend Integration"
    
    # Check if frontend is using unified attributes
    frontend_response=$(curl -sf "https://svetu.rs" 2>/dev/null || echo "")
    
    if [[ "$frontend_response" == *"unified-attributes"* ]] || [[ "$frontend_response" == *"unifiedAttributes"* ]]; then
        log_pass "Frontend is using unified attributes"
    else
        log_warn "Could not verify frontend integration"
    fi
}

validate_search_with_attributes() {
    log_test "Search with Unified Attributes"
    
    # Test search with attribute filters
    search_response=$(curl -sf "${API_BASE_URL}/marketplace/search?attributes[brand]=TestBrand" 2>/dev/null || echo "")
    
    if [[ -n "$search_response" ]] && [[ "$search_response" != *"error"* ]]; then
        log_pass "Search with attribute filters working"
    else
        log_warn "Search with attribute filters may have issues"
    fi
}

# Monitoring validation
validate_monitoring_setup() {
    log_test "Monitoring Dashboard"
    
    # Check if Grafana dashboard is accessible
    grafana_status=$(curl -sf -o /dev/null -w "%{http_code}" "https://grafana.svetu.rs/d/unified-attrs" 2>/dev/null || echo "000")
    
    if [[ "$grafana_status" == "200" ]]; then
        log_pass "Grafana dashboard accessible"
    else
        log_warn "Grafana dashboard returned status: $grafana_status"
    fi
}

validate_alerts_configured() {
    log_test "Alert Configuration"
    
    # Check if alerts are configured in Prometheus
    alerts=$(curl -sf "${METRICS_URL}/rules" | jq -r '.data.groups[].rules[] | select(.labels.component=="unified-attributes")' 2>/dev/null || echo "")
    
    if [[ -n "$alerts" ]]; then
        log_pass "Alerts configured for unified attributes"
    else
        log_warn "No specific alerts found for unified attributes"
    fi
}

# Generate validation report
generate_validation_report() {
    local report_file="/var/log/deployment/validation-report-$(date +%Y%m%d-%H%M%S).json"
    
    # Calculate success rate
    local success_rate=0
    if [[ $TOTAL_CHECKS -gt 0 ]]; then
        success_rate=$(echo "scale=2; $PASSED_CHECKS * 100 / $TOTAL_CHECKS" | bc)
    fi
    
    cat > "$report_file" <<EOF
{
  "validation": {
    "timestamp": "$(date -Iseconds)",
    "environment": "production",
    "summary": {
      "total_checks": $TOTAL_CHECKS,
      "passed": $PASSED_CHECKS,
      "failed": $FAILED_CHECKS,
      "warnings": $WARNINGS,
      "success_rate": "${success_rate}%"
    },
    "results": $(echo '{}' | jq --argjson results "$(declare -p validation_results | sed 's/declare -A validation_results=//' | sed "s/'/\"/g")" '$results'),
    "status": $([ $FAILED_CHECKS -eq 0 ] && echo '"SUCCESS"' || echo '"FAILURE"')
  }
}
EOF
    
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}Validation Report Summary${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "Total Checks: $TOTAL_CHECKS"
    echo -e "${GREEN}Passed: $PASSED_CHECKS${NC}"
    echo -e "${RED}Failed: $FAILED_CHECKS${NC}"
    echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
    echo -e "Success Rate: ${success_rate}%"
    echo -e "Report saved to: $report_file"
    
    if [[ $FAILED_CHECKS -eq 0 ]]; then
        echo -e "\n${GREEN}✅ VALIDATION SUCCESSFUL - System is fully operational${NC}"
        return 0
    else
        echo -e "\n${RED}❌ VALIDATION FAILED - Review failed checks above${NC}"
        return 1
    fi
}

# Main execution
main() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Unified Attributes System - Final Validation${NC}"
    echo -e "${BLUE}Date: $(date)${NC}"
    echo -e "${BLUE}========================================${NC}\n"
    
    # API Validation
    echo -e "${BLUE}[1/6] API Validation${NC}"
    validate_api_health
    validate_unified_attributes_endpoint
    validate_attribute_search
    validate_listing_creation_with_attributes
    
    # Database Validation
    echo -e "\n${BLUE}[2/6] Database Validation${NC}"
    validate_database_schema
    validate_data_migration
    validate_old_tables_archived
    
    # Performance Validation
    echo -e "\n${BLUE}[3/6] Performance Validation${NC}"
    validate_performance_metrics
    validate_cache_functionality
    
    # Feature Validation
    echo -e "\n${BLUE}[4/6] Feature Validation${NC}"
    validate_dual_write_disabled
    validate_fallback_disabled
    
    # Integration Validation
    echo -e "\n${BLUE}[5/6] Integration Validation${NC}"
    validate_frontend_integration
    validate_search_with_attributes
    
    # Monitoring Validation
    echo -e "\n${BLUE}[6/6] Monitoring Validation${NC}"
    validate_monitoring_setup
    validate_alerts_configured
    
    # Generate report
    generate_validation_report
}

# Run validation
main "$@"