#!/bin/bash

# ============================================================================
# Foreign Keys Integration Tests Runner
# ============================================================================
# This script runs all FK migration tests and generates a coverage report
# ============================================================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Database connection (PostgreSQL на порту 5433 согласно CLAUDE.md)
DB_URL="postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"

# Test directory
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo ""
echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}  Foreign Keys Integration Tests${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""

# Check if PostgreSQL is accessible
echo -e "${YELLOW}Checking database connection...${NC}"
if ! psql "$DB_URL" -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}❌ Cannot connect to database${NC}"
    echo -e "${RED}   Connection string: $DB_URL${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Database connection OK${NC}"
echo ""

# Function to run a test file
run_test() {
    local test_file=$1
    local test_name=$2

    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Running: ${test_name}${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    if psql "$DB_URL" -f "$test_file" 2>&1; then
        echo ""
        echo -e "${GREEN}✅ ${test_name} PASSED${NC}"
        return 0
    else
        echo ""
        echo -e "${RED}❌ ${test_name} FAILED${NC}"
        return 1
    fi
}

# Track test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Run CASCADE DELETE tests
TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "$TEST_DIR/test_foreign_keys_cascade.sql" "CASCADE DELETE Tests"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
echo ""

# Run RESTRICT tests
TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "$TEST_DIR/test_foreign_keys_restrict.sql" "RESTRICT Tests"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
echo ""

# Generate Coverage Report
echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}  Generating Coverage Report${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""

# Query FK constraints from database
COVERAGE_REPORT=$(psql "$DB_URL" -t -A -F'|' <<'EOSQL'
WITH fk_stats AS (
    SELECT
        tc.table_name,
        tc.constraint_name,
        kcu.column_name,
        ccu.table_name AS foreign_table_name,
        ccu.column_name AS foreign_column_name,
        rc.delete_rule,
        CASE
            WHEN tc.table_name IN ('c2c_images', 'c2c_attributes', 'c2c_favorites',
                                   'b2c_product_images', 'b2c_product_variants')
                 AND rc.delete_rule = 'CASCADE'
            THEN 'CASCADE-tested'
            WHEN tc.table_name IN ('c2c_categories', 'b2c_categories', 'storefronts')
                 AND rc.delete_rule IN ('RESTRICT', 'NO ACTION')
            THEN 'RESTRICT-tested'
            ELSE 'not-tested'
        END AS test_status
    FROM information_schema.table_constraints tc
    JOIN information_schema.key_column_usage kcu
        ON tc.constraint_name = kcu.constraint_name
        AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage ccu
        ON ccu.constraint_name = tc.constraint_name
        AND ccu.table_schema = tc.table_schema
    JOIN information_schema.referential_constraints rc
        ON rc.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'FOREIGN KEY'
    AND tc.table_schema = 'public'
    AND tc.table_name IN (
        'c2c_listings', 'c2c_images', 'c2c_attributes', 'c2c_favorites',
        'b2c_products', 'b2c_product_images', 'b2c_product_variants',
        'storefronts', 'c2c_categories', 'b2c_categories'
    )
    ORDER BY tc.table_name, tc.constraint_name
)
SELECT
    '📊 Foreign Keys Coverage Report:' AS report
UNION ALL
SELECT '   Total FK constraints: ' || COUNT(*)::TEXT FROM fk_stats
UNION ALL
SELECT '   CASCADE DELETE: ' || COUNT(*)::TEXT FROM fk_stats WHERE delete_rule = 'CASCADE'
UNION ALL
SELECT '   RESTRICT/NO ACTION: ' || COUNT(*)::TEXT FROM fk_stats WHERE delete_rule IN ('RESTRICT', 'NO ACTION')
UNION ALL
SELECT '   Tested (CASCADE): ' || COUNT(*)::TEXT FROM fk_stats WHERE test_status = 'CASCADE-tested'
UNION ALL
SELECT '   Tested (RESTRICT): ' || COUNT(*)::TEXT FROM fk_stats WHERE test_status = 'RESTRICT-tested'
UNION ALL
SELECT ''
UNION ALL
SELECT '📋 FK Constraints by Table:' AS report
UNION ALL
SELECT '   ' || table_name || ': ' || COUNT(*)::TEXT || ' FKs (' ||
       STRING_AGG(delete_rule, ', ') || ')'
FROM fk_stats
GROUP BY table_name
ORDER BY table_name;
EOSQL
)

echo "$COVERAGE_REPORT" | while IFS= read -r line; do
    if [[ $line == *"📊"* ]] || [[ $line == *"📋"* ]]; then
        echo -e "${BLUE}$line${NC}"
    elif [[ $line == *"Total"* ]] || [[ $line == *"CASCADE"* ]] || [[ $line == *"RESTRICT"* ]] || [[ $line == *"Tested"* ]]; then
        echo -e "${YELLOW}$line${NC}"
    else
        echo "$line"
    fi
done

echo ""

# Generate detailed FK list
echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}  Detailed FK Constraints List${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""

psql "$DB_URL" -c "
SELECT
    tc.table_name AS \"Table\",
    tc.constraint_name AS \"Constraint\",
    kcu.column_name AS \"Column\",
    ccu.table_name AS \"References\",
    rc.delete_rule AS \"On Delete\"
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
JOIN information_schema.referential_constraints rc
    ON rc.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
AND tc.table_schema = 'public'
AND tc.table_name IN (
    'c2c_listings', 'c2c_images', 'c2c_attributes', 'c2c_favorites',
    'b2c_products', 'b2c_product_images', 'b2c_product_variants',
    'storefronts', 'c2c_categories', 'b2c_categories'
)
ORDER BY tc.table_name, tc.constraint_name;
"

echo ""

# Final Summary
echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}  Test Summary${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""
echo -e "${YELLOW}Total test suites: ${TOTAL_TESTS}${NC}"
echo -e "${GREEN}Passed: ${PASSED_TESTS}${NC}"
echo -e "${RED}Failed: ${FAILED_TESTS}${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}============================================================================${NC}"
    echo -e "${GREEN}  ✅ ALL TESTS PASSED SUCCESSFULLY!${NC}"
    echo -e "${GREEN}============================================================================${NC}"
    exit 0
else
    echo -e "${RED}============================================================================${NC}"
    echo -e "${RED}  ❌ SOME TESTS FAILED${NC}"
    echo -e "${RED}============================================================================${NC}"
    exit 1
fi
