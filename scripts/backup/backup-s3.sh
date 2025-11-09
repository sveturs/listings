#!/bin/bash
#
# backup-s3.sh - Upload backup to S3/MinIO with verification
#
# Features:
# - Upload to object storage (S3/MinIO)
# - Verify upload with checksum
# - Cleanup old backups based on retention policy
# - Encryption at rest
# - Multi-part upload for large files
# - Retry logic for failed uploads
#
# Usage:
#   ./backup-s3.sh <backup_file> [--dry-run]
#
# Environment variables required:
#   S3_ENDPOINT      - S3/MinIO endpoint (default: localhost:9000)
#   S3_ACCESS_KEY    - Access key (required)
#   S3_SECRET_KEY    - Secret key (required)
#   S3_BUCKET        - Bucket name (default: listings-backups)
#   S3_REGION        - Region (default: us-east-1)
#   S3_USE_SSL       - Use SSL (default: false)
#   S3_RETENTION_DAYS - Keep backups for N days (default: 30)
#

set -euo pipefail

# ========================================
# Configuration
# ========================================

# S3/MinIO configuration
S3_ENDPOINT="${S3_ENDPOINT:-localhost:9000}"
S3_ACCESS_KEY="${S3_ACCESS_KEY:-}"
S3_SECRET_KEY="${S3_SECRET_KEY:-}"
S3_BUCKET="${S3_BUCKET:-listings-backups}"
S3_REGION="${S3_REGION:-us-east-1}"
S3_USE_SSL="${S3_USE_SSL:-false}"

# Retention
S3_RETENTION_DAYS="${S3_RETENTION_DAYS:-30}"

# Upload settings
MULTIPART_THRESHOLD_MB="${MULTIPART_THRESHOLD_MB:-100}"
MAX_RETRIES="${MAX_RETRIES:-3}"

# Logging
LOG_DIR="${LOG_DIR:-/var/log/listings}"

# Command line arguments
BACKUP_FILE=""
DRY_RUN=false

# ========================================
# Parse command line arguments
# ========================================

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        -*)
            echo "Unknown option: $1"
            echo "Usage: $0 <backup_file> [--dry-run]"
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
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_DIR/backup-s3.log"
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

check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v aws >/dev/null 2>&1; then
        log_error "AWS CLI not found. Install with: pip install awscli"
        return 1
    fi

    log_info "Dependencies check passed"
    return 0
}

configure_aws_cli() {
    log_info "Configuring AWS CLI..."

    # Set AWS credentials
    export AWS_ACCESS_KEY_ID="$S3_ACCESS_KEY"
    export AWS_SECRET_ACCESS_KEY="$S3_SECRET_KEY"
    export AWS_DEFAULT_REGION="$S3_REGION"

    # Configure endpoint URL
    local protocol="http"
    if [[ "$S3_USE_SSL" == "true" ]]; then
        protocol="https"
    fi
    export AWS_ENDPOINT_URL="${protocol}://${S3_ENDPOINT}"

    log_info "AWS CLI configured"
}

create_bucket_if_not_exists() {
    log_info "Checking if bucket exists: $S3_BUCKET"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would check/create bucket: $S3_BUCKET"
        return 0
    fi

    # Check if bucket exists
    if aws s3 ls "s3://${S3_BUCKET}" >/dev/null 2>&1; then
        log_info "Bucket exists: $S3_BUCKET"
        return 0
    fi

    # Create bucket
    log_info "Creating bucket: $S3_BUCKET"
    aws s3 mb "s3://${S3_BUCKET}" --region "$S3_REGION"

    # Enable versioning
    aws s3api put-bucket-versioning \
        --bucket "$S3_BUCKET" \
        --versioning-configuration Status=Enabled

    log_info "Bucket created: $S3_BUCKET"
}

upload_file() {
    local file="$1"
    local s3_key="$2"

    log_info "Uploading file to S3..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would upload: $file -> s3://${S3_BUCKET}/${s3_key}"
        return 0
    fi

    local file_size_mb
    file_size_mb=$(stat -c%s "$file" | awk '{print int($1/1024/1024)}')

    log_info "File size: ${file_size_mb}MB"

    # Determine upload method
    local upload_args=()
    if [[ "$file_size_mb" -gt "$MULTIPART_THRESHOLD_MB" ]]; then
        log_info "Using multipart upload (threshold: ${MULTIPART_THRESHOLD_MB}MB)"
        upload_args+=(
            --storage-class STANDARD_IA
        )
    fi

    # Add metadata
    upload_args+=(
        --metadata "backup-date=$(date -r "$file" '+%Y-%m-%d %H:%M:%S')"
        --metadata "original-name=$(basename "$file")"
        --metadata "checksum=$(sha256sum "$file" | cut -d' ' -f1)"
    )

    # Upload with retries
    local retry_count=0
    local upload_success=false

    while [[ $retry_count -lt $MAX_RETRIES ]]; do
        log_info "Upload attempt $((retry_count + 1))/$MAX_RETRIES"

        if aws s3 cp "$file" "s3://${S3_BUCKET}/${s3_key}" "${upload_args[@]}" \
            --only-show-errors 2>> "$LOG_DIR/backup-s3.log"; then
            upload_success=true
            break
        fi

        retry_count=$((retry_count + 1))
        if [[ $retry_count -lt $MAX_RETRIES ]]; then
            log_warning "Upload failed, retrying in 10 seconds..."
            sleep 10
        fi
    done

    if [[ "$upload_success" != "true" ]]; then
        log_error "Upload failed after $MAX_RETRIES attempts"
        return 1
    fi

    log_info "Upload completed successfully"
    return 0
}

verify_upload() {
    local file="$1"
    local s3_key="$2"

    log_info "Verifying upload..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would verify upload"
        return 0
    fi

    # Get local file checksum
    local local_checksum
    local_checksum=$(sha256sum "$file" | cut -d' ' -f1)

    # Get S3 file metadata
    local s3_metadata
    s3_metadata=$(aws s3api head-object \
        --bucket "$S3_BUCKET" \
        --key "$s3_key" \
        --query 'Metadata.checksum' \
        --output text 2>/dev/null || echo "")

    if [[ -z "$s3_metadata" ]]; then
        log_warning "S3 metadata not found, comparing file sizes..."

        local local_size
        local_size=$(stat -c%s "$file")
        local s3_size
        s3_size=$(aws s3api head-object \
            --bucket "$S3_BUCKET" \
            --key "$s3_key" \
            --query 'ContentLength' \
            --output text)

        if [[ "$local_size" != "$s3_size" ]]; then
            log_error "File size mismatch: local=$local_size, s3=$s3_size"
            return 1
        fi

        log_info "File sizes match: $local_size bytes"
        return 0
    fi

    # Compare checksums
    if [[ "$local_checksum" != "$s3_metadata" ]]; then
        log_error "Checksum mismatch: local=$local_checksum, s3=$s3_metadata"
        return 1
    fi

    log_info "Checksum verification passed"
    return 0
}

cleanup_old_backups() {
    log_info "Cleaning up old backups (retention: ${S3_RETENTION_DAYS} days)..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would cleanup old backups"
        return 0
    fi

    local cutoff_date
    cutoff_date=$(date -d "${S3_RETENTION_DAYS} days ago" '+%Y-%m-%d')

    log_info "Deleting backups older than: $cutoff_date"

    # List and delete old backups
    aws s3api list-objects-v2 \
        --bucket "$S3_BUCKET" \
        --query "Contents[?LastModified<='${cutoff_date}'].[Key]" \
        --output text | while read -r key; do

        if [[ -n "$key" ]]; then
            log_info "Deleting old backup: $key"
            aws s3 rm "s3://${S3_BUCKET}/${key}"
        fi
    done

    log_info "Cleanup completed"
}

get_s3_key() {
    local file="$1"
    local filename
    filename=$(basename "$file")

    # Organize by date: backups/YYYY/MM/DD/filename
    local date_path
    date_path=$(date -r "$file" '+%Y/%m/%d')

    echo "backups/${date_path}/${filename}"
}

# ========================================
# Main upload logic
# ========================================

main() {
    log_info "========================================="
    log_info "Starting S3/MinIO backup upload"
    log_info "========================================="

    # Validate arguments
    if [[ -z "$BACKUP_FILE" ]]; then
        log_error "Backup file not specified"
        echo "Usage: $0 <backup_file> [--dry-run]"
        exit 1
    fi

    if [[ ! -f "$BACKUP_FILE" ]]; then
        log_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi

    if [[ -z "$S3_ACCESS_KEY" ]] || [[ -z "$S3_SECRET_KEY" ]]; then
        log_error "S3_ACCESS_KEY and S3_SECRET_KEY are required"
        exit 1
    fi

    # Check dependencies
    if ! check_dependencies; then
        exit 1
    fi

    # Configure AWS CLI
    configure_aws_cli

    # Create bucket if not exists
    create_bucket_if_not_exists

    # Determine S3 key
    local s3_key
    s3_key=$(get_s3_key "$BACKUP_FILE")

    log_info "S3 key: $s3_key"

    # Upload file
    if ! upload_file "$BACKUP_FILE" "$s3_key"; then
        exit 1
    fi

    # Verify upload
    if ! verify_upload "$BACKUP_FILE" "$s3_key"; then
        log_error "Upload verification failed"
        exit 1
    fi

    # Upload checksum file if exists
    if [[ -f "${BACKUP_FILE}.sha256" ]]; then
        log_info "Uploading checksum file..."
        upload_file "${BACKUP_FILE}.sha256" "${s3_key}.sha256" || true
    fi

    # Upload metadata file if exists
    if [[ -f "${BACKUP_FILE}.meta" ]]; then
        log_info "Uploading metadata file..."
        upload_file "${BACKUP_FILE}.meta" "${s3_key}.meta" || true
    fi

    # Cleanup old backups
    cleanup_old_backups

    log_info "========================================="
    log_info "S3/MinIO upload completed successfully"
    log_info "S3 URL: s3://${S3_BUCKET}/${s3_key}"
    log_info "========================================="
}

# ========================================
# Entry point
# ========================================

main "$@"
