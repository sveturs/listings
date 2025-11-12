#!/bin/bash
#
# test-backup-restore.sh - Integration test for backup and restore cycle
#
# Tests the complete backup and restore workflow:
# 1. Create test data in database
# 2. Run backup
# 3. Modify database
# 4. Restore from backup
# 5. Verify data integrity
# 6. Cleanup
#
# Usage:
#   ./test-backup-restore.sh [--keep-backup] [--skip-cleanup]
#

set -euo pipefail

# ========================================
# Configuration
# ========================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_SCRIPT="${SCRIPT_DIR}/backup-db.sh"
RESTORE_SCRIPT="${SCRIPT_DIR}/restore-db.sh"
VERIFY_SCRIPT="${SCRIPT_DIR}/verify-backup.sh"

# Database configuration
DB_HOST="${TEST_DB_HOST:-localhost}"
DB_PORT="${TEST_DB_PORT:-35434}"
DB_NAME="${TEST_DB_NAME:-listings_dev_db}"
DB_USER="${TEST_DB_USER:-listings_user}"
DB_PASSWORD="${TEST_DB_PASSWORD:-}"

# Test configuration
TEST_BACKUP_DIR="/tmp/listings-backup-test-$$"
TEST_LOG_DIR="/tmp/listings-backup-test-logs-$$"

# Command line flags
KEEP_BACKUP=false
SKIP_CLEANUP=false

# ========================================
# Colors for output
# ========================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ========================================
# Parse command line arguments
# ========================================

while [[ $# -gt 0 ]]; do
    case "$1" in
        --keep-backup)
            KEEP_BACKUP=true
            shift
            ;;
        --skip-cleanup)
            SKIP_CLEANUP=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--keep-backup] [--skip-cleanup]"
            exit 1
            ;;
    esac
done

# ========================================
# Helper functions
# ========================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $*"
}

log_step() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$*${NC}"
    echo -e "${BLUE}========================================${NC}"
}

cleanup() {
    if [[ "$SKIP_CLEANUP" == "true" ]]; then
        log_warning "Skipping cleanup (--skip-cleanup flag)"
        return 0
    fi

    log_info "Cleaning up test artifacts..."

    # Remove test backup directory
    if [[ "$KEEP_BACKUP" != "true" ]]; then
        rm -rf "$TEST_BACKUP_DIR"
    else
        log_info "Keeping test backup directory: $TEST_BACKUP_DIR"
    fi

    rm -rf "$TEST_LOG_DIR"

    log_info "Cleanup completed"
}

trap cleanup EXIT

# ========================================
# Database helper functions
# ========================================

run_sql() {
    local sql="$1"
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "$sql" | xargs
}

get_table_count() {
    local table="$1"
    run_sql "SELECT COUNT(*) FROM $table;"
}

insert_test_data() {
    log_info "Inserting test data..."

    # Insert test listing
    run_sql "
        INSERT INTO listings (user_id, title, description, category_id, status)
        VALUES ('test-user-1', 'Test Listing - Backup Test', 'This is a test listing for backup testing', 1301, 'active')
        RETURNING id;
    " > /tmp/test_listing_id.txt

    local listing_id
    listing_id=$(cat /tmp/test_listing_id.txt)

    # Insert test product
    run_sql "
        INSERT INTO products (listing_id, name, price, currency, stock_quantity)
        VALUES ($listing_id, 'Test Product', 100.00, 'RSD', 50);
    "

    log_success "Test data inserted (listing_id: $listing_id)"
    echo "$listing_id"
}

verify_test_data() {
    local listing_id="$1"

    log_info "Verifying test data exists..."

    local count
    count=$(run_sql "SELECT COUNT(*) FROM listings WHERE id=$listing_id;")

    if [[ "$count" == "1" ]]; then
        log_success "Test data verified (listing_id: $listing_id found)"
        return 0
    else
        log_error "Test data not found (listing_id: $listing_id)"
        return 1
    fi
}

delete_test_data() {
    local listing_id="$1"

    log_info "Deleting test data..."

    run_sql "DELETE FROM listings WHERE id=$listing_id;" >/dev/null

    log_success "Test data deleted"
}

get_database_stats() {
    log_info "Database statistics:"

    local stats
    stats=$(run_sql "
        SELECT
            COUNT(*) as table_count,
            SUM(n_live_tup) as total_rows
        FROM pg_stat_user_tables;
    ")

    echo "  Tables and rows: $stats"
}

# ========================================
# Test functions
# ========================================

test_backup_creation() {
    log_step "Test 1: Backup Creation"

    log_info "Creating backup..."

    # Set environment for backup script
    export BACKUP_DIR="$TEST_BACKUP_DIR"
    export LOG_DIR="$TEST_LOG_DIR"
    export BACKUP_DB_HOST="$DB_HOST"
    export BACKUP_DB_PORT="$DB_PORT"
    export BACKUP_DB_NAME="$DB_NAME"
    export BACKUP_DB_USER="$DB_USER"
    export BACKUP_DB_PASSWORD="$DB_PASSWORD"
    export BACKUP_ENABLE_S3="false"
    export ENABLE_UPLOAD="false"

    mkdir -p "$TEST_BACKUP_DIR"
    mkdir -p "$TEST_LOG_DIR"

    # Run backup script
    if ! "$BACKUP_SCRIPT" --full-only; then
        log_error "Backup creation failed"
        return 1
    fi

    # Check if backup file exists
    local backup_file
    backup_file=$(find "$TEST_BACKUP_DIR/daily" -name "*.sql.gz" -type f | head -n 1)

    if [[ -z "$backup_file" ]]; then
        log_error "Backup file not found"
        return 1
    fi

    log_success "Backup created: $backup_file"

    # Check backup size
    local size
    size=$(du -h "$backup_file" | cut -f1)
    log_info "Backup size: $size"

    # Store backup file path for later tests
    echo "$backup_file" > /tmp/test_backup_file.txt

    return 0
}

test_backup_verification() {
    log_step "Test 2: Backup Verification"

    local backup_file
    backup_file=$(cat /tmp/test_backup_file.txt)

    log_info "Verifying backup: $backup_file"

    # Set environment for verify script
    export VERIFY_DB_HOST="$DB_HOST"
    export VERIFY_DB_PORT="$DB_PORT"
    export VERIFY_DB_USER="$DB_USER"
    export VERIFY_DB_PASSWORD="$DB_PASSWORD"
    export REPORT_DIR="$TEST_LOG_DIR/reports"

    # Run verify script
    if ! "$VERIFY_SCRIPT" "$backup_file" --quick; then
        log_error "Backup verification failed"
        return 1
    fi

    log_success "Backup verification passed"
    return 0
}

test_restore() {
    log_step "Test 3: Restore from Backup"

    local backup_file
    backup_file=$(cat /tmp/test_backup_file.txt)

    local listing_id
    listing_id=$(cat /tmp/test_listing_id.txt)

    log_info "Deleting test data (simulating data loss)..."
    delete_test_data "$listing_id"

    # Verify data is deleted
    if verify_test_data "$listing_id"; then
        log_error "Test data still exists (should be deleted)"
        return 1
    fi

    log_info "Test data deleted successfully"

    log_info "Restoring from backup: $backup_file"

    # Set environment for restore script
    export RESTORE_DB_HOST="$DB_HOST"
    export RESTORE_DB_PORT="$DB_PORT"
    export RESTORE_DB_NAME="$DB_NAME"
    export RESTORE_DB_USER="$DB_USER"
    export RESTORE_DB_PASSWORD="$DB_PASSWORD"
    export BACKUP_DIR="$TEST_BACKUP_DIR"

    # Run restore script
    if ! "$RESTORE_SCRIPT" "$backup_file" --no-backup; then
        log_error "Restore failed"
        return 1
    fi

    log_success "Restore completed"

    # Verify test data is back
    if ! verify_test_data "$listing_id"; then
        log_error "Test data not found after restore"
        return 1
    fi

    log_success "Data successfully restored"
    return 0
}

test_data_integrity() {
    log_step "Test 4: Data Integrity Check"

    local listing_id
    listing_id=$(cat /tmp/test_listing_id.txt)

    log_info "Checking data integrity after restore..."

    # Verify test data details
    local title
    title=$(run_sql "SELECT title FROM listings WHERE id=$listing_id;")

    if [[ "$title" == "Test Listing - Backup Test" ]]; then
        log_success "Data integrity verified: title matches"
    else
        log_error "Data integrity failed: title mismatch (got: $title)"
        return 1
    fi

    # Check products
    local product_count
    product_count=$(run_sql "SELECT COUNT(*) FROM products WHERE listing_id=$listing_id;")

    if [[ "$product_count" == "1" ]]; then
        log_success "Data integrity verified: product count matches"
    else
        log_error "Data integrity failed: expected 1 product, got $product_count"
        return 1
    fi

    log_success "All data integrity checks passed"
    return 0
}

# ========================================
# Main test execution
# ========================================

main() {
    log_step "Listings Backup & Restore Integration Test"

    # Validate configuration
    if [[ -z "$DB_PASSWORD" ]]; then
        log_error "TEST_DB_PASSWORD is required"
        exit 1
    fi

    # Check required scripts exist
    for script in "$BACKUP_SCRIPT" "$RESTORE_SCRIPT" "$VERIFY_SCRIPT"; do
        if [[ ! -f "$script" ]]; then
            log_error "Required script not found: $script"
            exit 1
        fi
    done

    # Get initial database stats
    get_database_stats

    # Create test data
    local listing_id
    listing_id=$(insert_test_data)

    # Run tests
    local test_results=()
    local failed_tests=0

    # Test 1: Backup creation
    if test_backup_creation; then
        test_results+=("✓ Backup Creation")
    else
        test_results+=("✗ Backup Creation")
        failed_tests=$((failed_tests + 1))
    fi

    # Test 2: Backup verification
    if test_backup_verification; then
        test_results+=("✓ Backup Verification")
    else
        test_results+=("✗ Backup Verification")
        failed_tests=$((failed_tests + 1))
    fi

    # Test 3: Restore
    if test_restore; then
        test_results+=("✓ Restore")
    else
        test_results+=("✗ Restore")
        failed_tests=$((failed_tests + 1))
    fi

    # Test 4: Data integrity
    if test_data_integrity; then
        test_results+=("✓ Data Integrity")
    else
        test_results+=("✗ Data Integrity")
        failed_tests=$((failed_tests + 1))
    fi

    # Cleanup test data
    delete_test_data "$listing_id"

    # Print summary
    log_step "Test Summary"

    for result in "${test_results[@]}"; do
        echo "$result"
    done

    echo ""

    if [[ $failed_tests -eq 0 ]]; then
        log_success "All tests passed! ✓"
        exit 0
    else
        log_error "$failed_tests test(s) failed"
        exit 1
    fi
}

# ========================================
# Entry point
# ========================================

main "$@"
