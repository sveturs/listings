#!/bin/bash
#
# verify-backup.sh - Backup integrity verification with test restore
#
# Features:
# - Test restore to temporary database
# - Validate data integrity (checksums, row counts)
# - Generate verification report
# - Cleanup temporary resources
# - Parallel verification for multiple backups
#
# Usage:
#   ./verify-backup.sh <backup_file> [--dry-run] [--quick]
#   ./verify-backup.sh --verify-all [--parallel]
#
# Examples:
#   ./verify-backup.sh /var/backups/listings/daily/backup.sql.gz
#   ./verify-backup.sh --verify-all
#   ./verify-backup.sh --verify-all --parallel
#

set -euo pipefail

# ========================================
# Configuration
# ========================================

# Database configuration
DB_HOST="${VERIFY_DB_HOST:-localhost}"
DB_PORT="${VERIFY_DB_PORT:-35434}"
DB_USER="${VERIFY_DB_USER:-listings_user}"
DB_PASSWORD="${VERIFY_DB_PASSWORD:-}"
DB_CONTAINER="${VERIFY_DB_CONTAINER:-listings_postgres}"

# Verification configuration
BACKUP_DIR="${BACKUP_DIR:-/var/backups/listings}"
LOG_DIR="${LOG_DIR:-/var/log/listings}"
REPORT_DIR="${REPORT_DIR:-/var/log/listings/reports}"
TEMP_DB_PREFIX="listings_verify_"

# Command line arguments
BACKUP_FILE=""
DRY_RUN=false
QUICK_CHECK=false
VERIFY_ALL=false
PARALLEL=false

# ========================================
# Parse command line arguments
# ========================================

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --quick)
            QUICK_CHECK=true
            shift
            ;;
        --verify-all)
            VERIFY_ALL=true
            shift
            ;;
        --parallel)
            PARALLEL=true
            shift
            ;;
        -*)
            echo "Unknown option: $1"
            echo "Usage: $0 <backup_file> [--dry-run] [--quick]"
            echo "       $0 --verify-all [--parallel]"
            exit 1
            ;;
        *)
            BACKUP_FILE="$1"
            shift
            ;;
    esac
done

# ========================================
# Helper functions
# ========================================

log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_DIR/verify-backup.log"
}

log_info() {
    log "INFO" "$@"
}

log_error() {
    log "ERROR" "$@"
}

log_warning() {
    log "WARNING" "$@"
}

cleanup_temp_database() {
    local temp_db="$1"

    log_info "Cleaning up temporary database: $temp_db"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would drop database: $temp_db"
        return 0
    fi

    # Terminate connections
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "
        SELECT pg_terminate_backend(pid)
        FROM pg_stat_activity
        WHERE datname = '$temp_db'
          AND pid <> pg_backend_pid();
    " >/dev/null 2>&1 || true

    sleep 1

    # Drop database
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c \
        "DROP DATABASE IF EXISTS ${temp_db};" >/dev/null 2>&1 || true
}

verify_file_integrity() {
    local file="$1"

    log_info "Verifying file integrity: $(basename "$file")"

    # Check file exists
    if [[ ! -f "$file" ]]; then
        log_error "File not found: $file"
        return 1
    fi

    # Check file size
    local size
    size=$(stat -c%s "$file")
    if [[ "$size" -lt 1000 ]]; then
        log_error "File too small (${size} bytes), possibly corrupted"
        return 1
    fi

    # Verify compression
    if [[ "$file" == *.gz ]]; then
        if ! gzip -t "$file" 2>/dev/null; then
            log_error "Gzip test failed - file corrupted"
            return 1
        fi
    fi

    # Verify checksum
    if [[ -f "${file}.sha256" ]]; then
        if ! sha256sum -c "${file}.sha256" >/dev/null 2>&1; then
            log_error "Checksum verification failed"
            return 1
        fi
        log_info "Checksum verified"
    else
        log_warning "Checksum file not found"
    fi

    log_info "File integrity: OK"
    return 0
}

test_restore() {
    local backup_file="$1"
    local temp_db="$2"

    log_info "Testing restore to temporary database: $temp_db"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would test restore"
        return 0
    fi

    # Create temporary database
    log_info "Creating temporary database..."
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c \
        "CREATE DATABASE ${temp_db} OWNER ${DB_USER};" >/dev/null 2>&1

    # Decompress if needed
    local restore_file="$backup_file"
    local cleanup_restore_file=false

    if [[ "$backup_file" == *.gz ]]; then
        restore_file="/tmp/verify_restore_$$.sql"
        log_info "Decompressing backup..."
        gunzip -c "$backup_file" > "$restore_file"
        cleanup_restore_file=true
    fi

    # Restore to temp database
    log_info "Restoring backup..."
    local start_time
    start_time=$(date +%s)

    PGPASSWORD="$DB_PASSWORD" psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$temp_db" \
        -f "$restore_file" \
        >> "$LOG_DIR/verify-backup.log" 2>&1

    local restore_result=$?
    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))

    # Cleanup decompressed file
    if [[ "$cleanup_restore_file" == "true" ]]; then
        rm -f "$restore_file"
    fi

    if [[ $restore_result -ne 0 ]]; then
        log_error "Test restore failed (exit code: $restore_result)"
        cleanup_temp_database "$temp_db"
        return 1
    fi

    log_info "Test restore completed in ${duration}s"
    return 0
}

verify_data_integrity() {
    local temp_db="$1"
    local report_file="$2"

    log_info "Verifying data integrity..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would verify data integrity"
        return 0
    fi

    # Check critical tables
    local critical_tables=("listings" "products" "inventory" "inventory_movements")
    local all_tables_ok=true

    echo "=== Data Integrity Report ===" >> "$report_file"
    echo "Database: $temp_db" >> "$report_file"
    echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')" >> "$report_file"
    echo "" >> "$report_file"

    for table in "${critical_tables[@]}"; do
        log_info "Checking table: $table"

        # Check if table exists
        local table_exists
        table_exists=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$temp_db" -t -c \
            "SELECT COUNT(*) FROM information_schema.tables WHERE table_name='$table';" | xargs)

        if [[ "$table_exists" != "1" ]]; then
            log_error "Table '$table' not found"
            echo "Table '$table': MISSING" >> "$report_file"
            all_tables_ok=false
            continue
        fi

        # Get row count
        local row_count
        row_count=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$temp_db" -t -c \
            "SELECT COUNT(*) FROM $table;" | xargs)

        log_info "Table '$table': $row_count rows"
        echo "Table '$table': $row_count rows" >> "$report_file"

        # Quick data validation (sample some rows)
        if [[ "$QUICK_CHECK" != "true" ]]; then
            local sample_query="SELECT * FROM $table LIMIT 10;"
            PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$temp_db" -c \
                "$sample_query" >> "$report_file" 2>&1 || true
        fi
    done

    # Get database statistics
    echo "" >> "$report_file"
    echo "=== Database Statistics ===" >> "$report_file"

    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$temp_db" -c "
        SELECT
            schemaname,
            COUNT(*) as table_count,
            SUM(n_live_tup) as total_rows
        FROM pg_stat_user_tables
        GROUP BY schemaname;
    " >> "$report_file" 2>&1 || true

    # Check for obvious data corruption
    echo "" >> "$report_file"
    echo "=== Corruption Check ===" >> "$report_file"

    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$temp_db" -c \
        "SELECT * FROM pg_stat_database WHERE datname='$temp_db';" \
        >> "$report_file" 2>&1 || true

    if [[ "$all_tables_ok" == "true" ]]; then
        log_info "Data integrity: OK"
        echo "" >> "$report_file"
        echo "Result: PASSED" >> "$report_file"
        return 0
    else
        log_error "Data integrity check failed"
        echo "" >> "$report_file"
        echo "Result: FAILED" >> "$report_file"
        return 1
    fi
}

verify_backup() {
    local backup_file="$1"

    log_info "========================================="
    log_info "Verifying backup: $(basename "$backup_file")"
    log_info "========================================="

    local temp_db="${TEMP_DB_PREFIX}$(date +%s)"
    local report_file="${REPORT_DIR}/verify_$(basename "$backup_file")_$(date +%Y%m%d_%H%M%S).txt"

    mkdir -p "$REPORT_DIR"

    # Step 1: File integrity
    if ! verify_file_integrity "$backup_file"; then
        log_error "File integrity check failed"
        echo "File integrity: FAILED" > "$report_file"
        return 1
    fi

    # Step 2: Test restore
    if ! test_restore "$backup_file" "$temp_db"; then
        log_error "Test restore failed"
        echo "Test restore: FAILED" >> "$report_file"
        cleanup_temp_database "$temp_db"
        return 1
    fi

    # Step 3: Data integrity
    local data_ok=true
    if ! verify_data_integrity "$temp_db" "$report_file"; then
        log_error "Data integrity check failed"
        data_ok=false
    fi

    # Cleanup
    cleanup_temp_database "$temp_db"

    # Generate summary
    log_info "Verification report: $report_file"
    log_info "========================================="

    if [[ "$data_ok" == "true" ]]; then
        log_info "Backup verification: PASSED"
        return 0
    else
        log_error "Backup verification: FAILED"
        return 1
    fi
}

verify_all_backups() {
    log_info "Verifying all backups in: $BACKUP_DIR"

    local backup_files=()

    # Find all backup files
    while IFS= read -r -d '' file; do
        backup_files+=("$file")
    done < <(find "$BACKUP_DIR" -name "*.sql.gz" -type f -print0)

    local total_backups="${#backup_files[@]}"
    log_info "Found $total_backups backup files"

    if [[ "$total_backups" -eq 0 ]]; then
        log_warning "No backup files found"
        return 0
    fi

    local failed_backups=0
    local passed_backups=0

    # Verify backups
    if [[ "$PARALLEL" == "true" ]]; then
        log_info "Running parallel verification (max 4 concurrent)"

        # Export functions for parallel execution
        export -f verify_backup verify_file_integrity test_restore verify_data_integrity cleanup_temp_database log log_info log_error log_warning
        export DB_HOST DB_PORT DB_USER DB_PASSWORD DB_CONTAINER
        export BACKUP_DIR LOG_DIR REPORT_DIR TEMP_DB_PREFIX
        export DRY_RUN QUICK_CHECK

        # Run in parallel with xargs
        printf '%s\n' "${backup_files[@]}" | xargs -P 4 -I {} bash -c 'verify_backup "{}"' || true

    else
        log_info "Running sequential verification"

        for backup_file in "${backup_files[@]}"; do
            if verify_backup "$backup_file"; then
                passed_backups=$((passed_backups + 1))
            else
                failed_backups=$((failed_backups + 1))
            fi
        done

        log_info "========================================="
        log_info "Verification Summary"
        log_info "Total: $total_backups"
        log_info "Passed: $passed_backups"
        log_info "Failed: $failed_backups"
        log_info "========================================="
    fi
}

# ========================================
# Main verification logic
# ========================================

main() {
    log_info "Starting backup verification"

    # Validate configuration
    if [[ -z "$DB_PASSWORD" ]]; then
        log_error "VERIFY_DB_PASSWORD is required"
        exit 1
    fi

    mkdir -p "$LOG_DIR"
    mkdir -p "$REPORT_DIR"

    # Verify all or single backup
    if [[ "$VERIFY_ALL" == "true" ]]; then
        verify_all_backups
    else
        if [[ -z "$BACKUP_FILE" ]]; then
            log_error "Backup file not specified"
            echo "Usage: $0 <backup_file> [--dry-run] [--quick]"
            echo "       $0 --verify-all [--parallel]"
            exit 1
        fi

        verify_backup "$BACKUP_FILE"
    fi
}

# ========================================
# Entry point
# ========================================

main "$@"
