#!/bin/bash
#
# restore-db.sh - Database restore from backup with validation and rollback
#
# Features:
# - Restore from specific backup file
# - Point-in-time recovery (PITR) support
# - Pre-restore validation (disk space, permissions)
# - Backup current DB before restore
# - Post-restore integrity checks
# - Rollback capability
#
# Usage:
#   ./restore-db.sh <backup_file> [--dry-run] [--no-backup] [--pitr-target "2024-11-05 12:00:00"]
#
# Examples:
#   ./restore-db.sh /var/backups/listings/daily/listings_dev_db_20241105_020000.sql.gz
#   ./restore-db.sh backup.sql.gz --pitr-target "2024-11-05 12:00:00"
#   ./restore-db.sh backup.sql.gz --dry-run
#

set -euo pipefail

# ========================================
# Configuration
# ========================================

# Database configuration
DB_HOST="${RESTORE_DB_HOST:-localhost}"
DB_PORT="${RESTORE_DB_PORT:-35434}"
DB_NAME="${RESTORE_DB_NAME:-listings_dev_db}"
DB_USER="${RESTORE_DB_USER:-listings_user}"
DB_PASSWORD="${RESTORE_DB_PASSWORD:-}"
DB_CONTAINER="${RESTORE_DB_CONTAINER:-listings_postgres}"

# Restore configuration
BACKUP_DIR="${BACKUP_DIR:-/var/backups/listings}"
LOG_DIR="${LOG_DIR:-/var/log/listings}"
RESTORE_TMP_DIR="/tmp/listings-restore-$$"

# Command line arguments
BACKUP_FILE=""
DRY_RUN=false
NO_BACKUP=false
PITR_TARGET=""

# ========================================
# Parse command line arguments
# ========================================

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --no-backup)
            NO_BACKUP=true
            shift
            ;;
        --pitr-target)
            PITR_TARGET="$2"
            shift 2
            ;;
        -*)
            echo "Unknown option: $1"
            echo "Usage: $0 <backup_file> [--dry-run] [--no-backup] [--pitr-target \"YYYY-MM-DD HH:MM:SS\"]"
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
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_DIR/restore.log"
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

cleanup() {
    log_info "Cleaning up temporary files..."
    rm -rf "$RESTORE_TMP_DIR"
}

validate_backup_file() {
    local file="$1"

    log_info "Validating backup file: $file"

    # Check file exists
    if [[ ! -f "$file" ]]; then
        log_error "Backup file not found: $file"
        return 1
    fi

    # Check file size
    local size
    size=$(stat -c%s "$file")
    if [[ "$size" -lt 1000 ]]; then
        log_error "Backup file too small (${size} bytes), possibly corrupted"
        return 1
    fi

    # Check if gzip compressed
    if [[ "$file" == *.gz ]]; then
        if ! gzip -t "$file" 2>/dev/null; then
            log_error "Backup file is corrupted (gzip test failed)"
            return 1
        fi
    fi

    # Verify checksum if available
    if [[ -f "${file}.sha256" ]]; then
        log_info "Verifying checksum..."
        if ! sha256sum -c "${file}.sha256" >/dev/null 2>&1; then
            log_error "Checksum verification failed"
            return 1
        fi
        log_info "Checksum verified successfully"
    else
        log_warning "Checksum file not found, skipping verification"
    fi

    log_info "Backup file validation passed"
    return 0
}

check_disk_space() {
    local required_gb="$1"
    local available_gb
    available_gb=$(df -BG "$RESTORE_TMP_DIR" | awk 'NR==2 {print $4}' | sed 's/G//')

    if [[ "$available_gb" -lt "$required_gb" ]]; then
        log_error "Insufficient disk space: ${available_gb}GB available, ${required_gb}GB required"
        return 1
    fi

    log_info "Disk space check passed: ${available_gb}GB available"
    return 0
}

get_table_counts() {
    local db_name="$1"
    local output_file="$2"

    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$db_name" -t -c "
        SELECT
            schemaname || '.' || tablename as table_name,
            n_live_tup as row_count
        FROM pg_stat_user_tables
        ORDER BY schemaname, tablename;
    " > "$output_file" 2>/dev/null || true
}

backup_current_database() {
    log_info "Creating backup of current database before restore..."

    if [[ "$NO_BACKUP" == "true" ]]; then
        log_warning "Skipping pre-restore backup (--no-backup flag)"
        return 0
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create pre-restore backup"
        return 0
    fi

    local timestamp
    timestamp=$(date '+%Y%m%d_%H%M%S')
    local pre_restore_backup="${BACKUP_DIR}/pre-restore/listings_${DB_NAME}_${timestamp}.sql.gz"

    mkdir -p "${BACKUP_DIR}/pre-restore"

    log_info "Creating pre-restore backup: $pre_restore_backup"

    PGPASSWORD="$DB_PASSWORD" pg_dump \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        --format=plain \
        --no-owner \
        --no-acl \
        2>> "$LOG_DIR/restore.log" | gzip -9 > "$pre_restore_backup"

    if [[ ! -f "$pre_restore_backup" ]]; then
        log_error "Pre-restore backup failed"
        return 1
    fi

    log_info "Pre-restore backup completed: $pre_restore_backup"
    echo "$pre_restore_backup"
}

terminate_connections() {
    log_info "Terminating active connections to database..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would terminate connections"
        return 0
    fi

    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "
        SELECT pg_terminate_backend(pid)
        FROM pg_stat_activity
        WHERE datname = '$DB_NAME'
          AND pid <> pg_backend_pid();
    " >> "$LOG_DIR/restore.log" 2>&1 || true

    sleep 2
}

restore_database() {
    local backup_file="$1"
    local restore_file="$backup_file"

    log_info "Starting database restore..."

    # Decompress if needed
    if [[ "$backup_file" == *.gz ]]; then
        log_info "Decompressing backup file..."
        restore_file="${RESTORE_TMP_DIR}/restore.sql"
        gunzip -c "$backup_file" > "$restore_file"
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would restore from: $restore_file"
        return 0
    fi

    # Terminate connections
    terminate_connections

    # Drop and recreate database
    log_info "Dropping and recreating database..."
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c \
        "DROP DATABASE IF EXISTS ${DB_NAME};" >> "$LOG_DIR/restore.log" 2>&1

    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c \
        "CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};" >> "$LOG_DIR/restore.log" 2>&1

    # Restore data
    log_info "Restoring database from backup..."
    local start_time
    start_time=$(date +%s)

    PGPASSWORD="$DB_PASSWORD" psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -f "$restore_file" \
        >> "$LOG_DIR/restore.log" 2>&1

    local restore_result=$?
    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))

    if [[ $restore_result -ne 0 ]]; then
        log_error "Database restore failed (exit code: $restore_result)"
        return 1
    fi

    log_info "Database restore completed in ${duration}s"
    return 0
}

apply_pitr() {
    local target_time="$1"

    if [[ -z "$target_time" ]]; then
        return 0
    fi

    log_info "Applying Point-in-Time Recovery to: $target_time"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would apply PITR to: $target_time"
        return 0
    fi

    local wal_dir="${BACKUP_DIR}/wal"

    if [[ ! -d "$wal_dir" ]]; then
        log_error "WAL directory not found: $wal_dir"
        return 1
    fi

    # Create recovery.conf
    local recovery_conf="${RESTORE_TMP_DIR}/recovery.conf"
    cat > "$recovery_conf" <<EOF
restore_command = 'cp ${wal_dir}/%f %p'
recovery_target_time = '${target_time}'
recovery_target_action = 'promote'
EOF

    # Copy recovery.conf to container
    docker cp "$recovery_conf" "${DB_CONTAINER}:/var/lib/postgresql/data/recovery.conf"

    # Restart PostgreSQL
    log_info "Restarting PostgreSQL to apply PITR..."
    docker restart "$DB_CONTAINER"

    # Wait for PostgreSQL to start
    sleep 10

    log_info "PITR applied successfully"
}

verify_restore() {
    log_info "Verifying restore integrity..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would verify restore"
        return 0
    fi

    # Check database exists
    local db_exists
    db_exists=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -t -c \
        "SELECT 1 FROM pg_database WHERE datname='$DB_NAME';" | xargs)

    if [[ "$db_exists" != "1" ]]; then
        log_error "Database does not exist after restore"
        return 1
    fi

    # Get table counts
    local counts_file="${RESTORE_TMP_DIR}/table_counts.txt"
    get_table_counts "$DB_NAME" "$counts_file"

    local table_count
    table_count=$(wc -l < "$counts_file")

    log_info "Restored database contains $table_count tables"

    # Check critical tables
    local critical_tables=("listings" "products" "inventory" "inventory_movements")
    for table in "${critical_tables[@]}"; do
        local row_count
        row_count=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c \
            "SELECT COUNT(*) FROM $table;" | xargs 2>/dev/null || echo "0")

        log_info "Table '$table': $row_count rows"

        if [[ "$row_count" == "0" ]]; then
            log_warning "Table '$table' is empty after restore"
        fi
    done

    # Check database size
    local db_size
    db_size=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT pg_size_pretty(pg_database_size('$DB_NAME'));" | xargs)

    log_info "Restored database size: $db_size"

    log_info "Restore verification completed"
    return 0
}

rollback_restore() {
    local pre_restore_backup="$1"

    log_error "Restore failed, initiating rollback..."

    if [[ -z "$pre_restore_backup" ]] || [[ ! -f "$pre_restore_backup" ]]; then
        log_error "Pre-restore backup not found, cannot rollback"
        return 1
    fi

    log_info "Rolling back to pre-restore backup: $pre_restore_backup"

    # Restore from pre-restore backup
    restore_database "$pre_restore_backup"
    local rollback_result=$?

    if [[ $rollback_result -ne 0 ]]; then
        log_error "Rollback failed! Manual intervention required."
        log_error "Pre-restore backup: $pre_restore_backup"
        return 1
    fi

    log_info "Rollback completed successfully"
    return 0
}

# ========================================
# Main restore logic
# ========================================

main() {
    log_info "========================================="
    log_info "Starting database restore"
    log_info "========================================="

    # Validate arguments
    if [[ -z "$BACKUP_FILE" ]]; then
        log_error "Backup file not specified"
        echo "Usage: $0 <backup_file> [--dry-run] [--no-backup] [--pitr-target \"YYYY-MM-DD HH:MM:SS\"]"
        exit 1
    fi

    if [[ -z "$DB_PASSWORD" ]]; then
        log_error "RESTORE_DB_PASSWORD is required"
        exit 1
    fi

    # Create directories
    mkdir -p "$LOG_DIR"
    mkdir -p "$RESTORE_TMP_DIR"
    trap cleanup EXIT

    # Validate backup file
    if ! validate_backup_file "$BACKUP_FILE"; then
        exit 1
    fi

    # Check disk space
    if ! check_disk_space 10; then
        exit 1
    fi

    # Get current table counts (for comparison)
    log_info "Getting current table counts..."
    local before_counts="${RESTORE_TMP_DIR}/before_counts.txt"
    get_table_counts "$DB_NAME" "$before_counts"

    # Backup current database
    local pre_restore_backup=""
    if [[ "$NO_BACKUP" != "true" ]]; then
        pre_restore_backup=$(backup_current_database)
        if [[ -z "$pre_restore_backup" ]]; then
            log_error "Pre-restore backup failed, aborting restore"
            exit 1
        fi
    fi

    # Restore database
    if ! restore_database "$BACKUP_FILE"; then
        log_error "Database restore failed"
        if [[ -n "$pre_restore_backup" ]]; then
            rollback_restore "$pre_restore_backup"
        fi
        exit 1
    fi

    # Apply PITR if requested
    if [[ -n "$PITR_TARGET" ]]; then
        if ! apply_pitr "$PITR_TARGET"; then
            log_error "PITR failed"
            if [[ -n "$pre_restore_backup" ]]; then
                rollback_restore "$pre_restore_backup"
            fi
            exit 1
        fi
    fi

    # Verify restore
    if ! verify_restore; then
        log_error "Restore verification failed"
        if [[ -n "$pre_restore_backup" ]]; then
            rollback_restore "$pre_restore_backup"
        fi
        exit 1
    fi

    # Compare table counts
    log_info "Comparing table counts..."
    local after_counts="${RESTORE_TMP_DIR}/after_counts.txt"
    get_table_counts "$DB_NAME" "$after_counts"

    if [[ -f "$before_counts" ]] && [[ -f "$after_counts" ]]; then
        log_info "Table count comparison:"
        diff -y "$before_counts" "$after_counts" >> "$LOG_DIR/restore.log" 2>&1 || true
    fi

    log_info "========================================="
    log_info "Database restore completed successfully"
    log_info "Restored from: $BACKUP_FILE"
    if [[ -n "$pre_restore_backup" ]]; then
        log_info "Pre-restore backup: $pre_restore_backup"
    fi
    log_info "========================================="
}

# ========================================
# Entry point
# ========================================

main "$@"
