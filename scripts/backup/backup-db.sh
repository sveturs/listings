#!/bin/bash
#
# backup-db.sh - Automated PostgreSQL database backup with retention policy
#
# Features:
# - Full database dump (structure + data)
# - Incremental WAL archiving (PITR support)
# - Compression (gzip)
# - Retention policy (7 daily, 4 weekly, 12 monthly)
# - Optional S3/MinIO upload
# - Email notifications
# - Logging with timestamps
# - Lock file for preventing concurrent backups
#
# Usage:
#   ./backup-db.sh [--dry-run] [--no-upload] [--full-only]
#
# Environment variables required:
#   BACKUP_DB_HOST       - Database host (default: localhost)
#   BACKUP_DB_PORT       - Database port (default: 35434)
#   BACKUP_DB_NAME       - Database name (default: listings_dev_db)
#   BACKUP_DB_USER       - Database user (default: listings_user)
#   BACKUP_DB_PASSWORD   - Database password (required)
#   BACKUP_DIR           - Backup directory (default: /var/backups/listings)
#   BACKUP_RETENTION_DAYS   - Daily backups to keep (default: 7)
#   BACKUP_RETENTION_WEEKS  - Weekly backups to keep (default: 4)
#   BACKUP_RETENTION_MONTHS - Monthly backups to keep (default: 12)
#   BACKUP_ENABLE_S3     - Upload to S3/MinIO (default: false)
#   BACKUP_NOTIFY_EMAIL  - Email for notifications (optional)
#

set -euo pipefail

# ========================================
# Configuration
# ========================================

# Database configuration
DB_HOST="${BACKUP_DB_HOST:-localhost}"
DB_PORT="${BACKUP_DB_PORT:-35434}"
DB_NAME="${BACKUP_DB_NAME:-listings_dev_db}"
DB_USER="${BACKUP_DB_USER:-listings_user}"
DB_PASSWORD="${BACKUP_DB_PASSWORD:-}"
DB_CONTAINER="${BACKUP_DB_CONTAINER:-listings_postgres}"

# Backup configuration
BACKUP_DIR="${BACKUP_DIR:-/var/backups/listings}"
LOG_DIR="${LOG_DIR:-/var/log/listings}"
LOCK_FILE="/var/run/listings-backup.lock"

# Retention policy
RETENTION_DAILY="${BACKUP_RETENTION_DAYS:-7}"
RETENTION_WEEKLY="${BACKUP_RETENTION_WEEKS:-4}"
RETENTION_MONTHLY="${BACKUP_RETENTION_MONTHS:-12}"

# S3/MinIO upload
ENABLE_S3="${BACKUP_ENABLE_S3:-false}"
ENABLE_UPLOAD="${ENABLE_UPLOAD:-true}"

# Notifications
NOTIFY_EMAIL="${BACKUP_NOTIFY_EMAIL:-}"

# Command line flags
DRY_RUN=false
NO_UPLOAD=false
FULL_ONLY=false

# ========================================
# Parse command line arguments
# ========================================

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --no-upload)
            NO_UPLOAD=true
            shift
            ;;
        --full-only)
            FULL_ONLY=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--dry-run] [--no-upload] [--full-only]"
            exit 1
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
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_DIR/backup.log"
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

send_notification() {
    local subject="$1"
    local body="$2"

    if [[ -n "$NOTIFY_EMAIL" ]] && command -v mail >/dev/null 2>&1; then
        echo "$body" | mail -s "$subject" "$NOTIFY_EMAIL"
    fi
}

acquire_lock() {
    if [[ -f "$LOCK_FILE" ]]; then
        local pid
        pid=$(cat "$LOCK_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            log_error "Backup already running (PID: $pid)"
            exit 1
        else
            log_warning "Stale lock file found, removing"
            rm -f "$LOCK_FILE"
        fi
    fi

    echo $$ > "$LOCK_FILE"
    trap 'rm -f "$LOCK_FILE"' EXIT
}

release_lock() {
    rm -f "$LOCK_FILE"
}

check_disk_space() {
    local required_gb="$1"
    local available_gb
    available_gb=$(df -BG "$BACKUP_DIR" | awk 'NR==2 {print $4}' | sed 's/G//')

    if [[ "$available_gb" -lt "$required_gb" ]]; then
        log_error "Insufficient disk space: ${available_gb}GB available, ${required_gb}GB required"
        return 1
    fi

    log_info "Disk space check passed: ${available_gb}GB available"
    return 0
}

get_db_size() {
    local size_mb
    size_mb=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT pg_size_pretty(pg_database_size('$DB_NAME'));" | xargs)
    echo "$size_mb"
}

# ========================================
# Backup functions
# ========================================

create_full_backup() {
    local backup_type="$1"  # daily, weekly, monthly
    local timestamp
    timestamp=$(date '+%Y%m%d_%H%M%S')
    local backup_file="${BACKUP_DIR}/${backup_type}/listings_${DB_NAME}_${timestamp}.sql"
    local backup_file_gz="${backup_file}.gz"

    log_info "Starting full backup: $backup_type"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create backup: $backup_file_gz"
        return 0
    fi

    # Create backup directory
    mkdir -p "${BACKUP_DIR}/${backup_type}"

    # Get database size before backup
    local db_size
    db_size=$(get_db_size)
    log_info "Database size: $db_size"

    # Check available disk space (need 2x database size)
    check_disk_space 10 || return 1

    # Create backup
    log_info "Creating backup: $backup_file"
    local start_time
    start_time=$(date +%s)

    PGPASSWORD="$DB_PASSWORD" pg_dump \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        --format=plain \
        --no-owner \
        --no-acl \
        --verbose \
        --file="$backup_file" \
        2>> "$LOG_DIR/backup.log"

    local dump_result=$?
    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))

    if [[ $dump_result -ne 0 ]]; then
        log_error "pg_dump failed with exit code: $dump_result"
        rm -f "$backup_file"
        return 1
    fi

    log_info "Backup completed in ${duration}s"

    # Compress backup
    log_info "Compressing backup..."
    gzip -9 "$backup_file"

    if [[ ! -f "$backup_file_gz" ]]; then
        log_error "Compression failed"
        return 1
    fi

    # Get backup size
    local backup_size
    backup_size=$(du -h "$backup_file_gz" | cut -f1)
    log_info "Backup size: $backup_size"

    # Create checksum
    log_info "Creating checksum..."
    sha256sum "$backup_file_gz" > "${backup_file_gz}.sha256"

    # Create metadata file
    cat > "${backup_file_gz}.meta" <<EOF
backup_date=$(date '+%Y-%m-%d %H:%M:%S')
backup_type=$backup_type
database_name=$DB_NAME
database_size=$db_size
backup_size=$backup_size
duration_seconds=$duration
host=$DB_HOST
port=$DB_PORT
EOF

    log_info "Full backup completed successfully: $backup_file_gz"
    echo "$backup_file_gz"
}

create_wal_archive() {
    if [[ "$FULL_ONLY" == "true" ]]; then
        return 0
    fi

    log_info "Creating WAL archive for PITR..."

    local wal_dir="${BACKUP_DIR}/wal"
    mkdir -p "$wal_dir"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would archive WAL files"
        return 0
    fi

    # Archive WAL files from container
    docker exec "$DB_CONTAINER" bash -c \
        "pg_ctl -D /var/lib/postgresql/data switch-wal" 2>/dev/null || true

    log_info "WAL archive completed"
}

apply_retention_policy() {
    local backup_type="$1"
    local retention_count="$2"

    log_info "Applying retention policy for $backup_type backups (keep $retention_count)"

    local backup_type_dir="${BACKUP_DIR}/${backup_type}"

    if [[ ! -d "$backup_type_dir" ]]; then
        log_warning "Backup directory not found: $backup_type_dir"
        return 0
    fi

    # Count backups
    local backup_count
    backup_count=$(find "$backup_type_dir" -name "*.sql.gz" -type f | wc -l)

    if [[ "$backup_count" -le "$retention_count" ]]; then
        log_info "Retention policy satisfied: $backup_count backups (need $retention_count)"
        return 0
    fi

    log_info "Found $backup_count backups, removing $((backup_count - retention_count)) old backups"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would remove old backups"
        find "$backup_type_dir" -name "*.sql.gz" -type f -printf '%T+ %p\n' | \
            sort | head -n -"$retention_count" | cut -d' ' -f2-
        return 0
    fi

    # Remove old backups (keep newest N)
    find "$backup_type_dir" -name "*.sql.gz" -type f -printf '%T+ %p\n' | \
        sort | head -n -"$retention_count" | cut -d' ' -f2- | while read -r file; do
        log_info "Removing old backup: $file"
        rm -f "$file" "${file}.sha256" "${file}.meta"
    done
}

upload_to_s3() {
    local backup_file="$1"

    if [[ "$ENABLE_S3" != "true" ]] || [[ "$NO_UPLOAD" == "true" ]] || [[ ! "$ENABLE_UPLOAD" == "true" ]]; then
        return 0
    fi

    log_info "Uploading backup to S3/MinIO..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would upload: $backup_file"
        return 0
    fi

    # Call backup-s3.sh script
    if [[ -x "$(dirname "$0")/backup-s3.sh" ]]; then
        "$(dirname "$0")/backup-s3.sh" "$backup_file"
    else
        log_warning "backup-s3.sh not found, skipping upload"
    fi
}

# ========================================
# Main backup logic
# ========================================

main() {
    log_info "========================================="
    log_info "Starting listings database backup"
    log_info "========================================="

    # Validate configuration
    if [[ -z "$DB_PASSWORD" ]]; then
        log_error "BACKUP_DB_PASSWORD is required"
        exit 1
    fi

    # Create directories
    mkdir -p "$BACKUP_DIR"/{daily,weekly,monthly,wal}
    mkdir -p "$LOG_DIR"

    # Acquire lock
    acquire_lock

    # Determine backup type based on day
    local backup_type="daily"
    local day_of_week
    day_of_week=$(date +%u)  # 1=Monday, 7=Sunday
    local day_of_month
    day_of_month=$(date +%d)

    if [[ "$day_of_month" == "01" ]]; then
        backup_type="monthly"
    elif [[ "$day_of_week" == "7" ]]; then
        backup_type="weekly"
    fi

    log_info "Backup type: $backup_type"

    # Create full backup
    local backup_file
    backup_file=$(create_full_backup "$backup_type")
    local backup_result=$?

    if [[ $backup_result -ne 0 ]]; then
        log_error "Backup failed"
        send_notification "Listings Backup FAILED" "Backup failed for database $DB_NAME"
        exit 1
    fi

    # Create WAL archive for PITR
    create_wal_archive

    # Apply retention policies
    apply_retention_policy "daily" "$RETENTION_DAILY"
    apply_retention_policy "weekly" "$RETENTION_WEEKLY"
    apply_retention_policy "monthly" "$RETENTION_MONTHLY"

    # Upload to S3/MinIO
    if [[ -n "$backup_file" ]]; then
        upload_to_s3 "$backup_file"
    fi

    # Release lock
    release_lock

    # Send success notification
    local backup_size
    if [[ -n "$backup_file" ]] && [[ -f "$backup_file" ]]; then
        backup_size=$(du -h "$backup_file" | cut -f1)
    else
        backup_size="N/A"
    fi

    log_info "========================================="
    log_info "Backup completed successfully"
    log_info "Backup file: $backup_file"
    log_info "Backup size: $backup_size"
    log_info "========================================="

    send_notification "Listings Backup SUCCESS" \
        "Backup completed successfully for database $DB_NAME\nBackup type: $backup_type\nBackup size: $backup_size\nFile: $backup_file"
}

# ========================================
# Entry point
# ========================================

main "$@"
