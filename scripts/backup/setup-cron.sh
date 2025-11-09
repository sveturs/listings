#!/bin/bash
#
# setup-cron.sh - Setup automated backup cron jobs
#
# Features:
# - Configure daily/weekly backup schedules
# - Setup log rotation
# - Test first backup
# - Email notifications setup
# - Monitoring integration
#
# Usage:
#   sudo ./setup-cron.sh [--dry-run] [--uninstall]
#

set -euo pipefail

# ========================================
# Configuration
# ========================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_SCRIPT="${SCRIPT_DIR}/backup-db.sh"
VERIFY_SCRIPT="${SCRIPT_DIR}/verify-backup.sh"
MONITOR_SCRIPT="${SCRIPT_DIR}/monitor-backups.py"

# Backup schedule
DAILY_BACKUP_TIME="${DAILY_BACKUP_TIME:-02:00}"  # 2 AM
WEEKLY_BACKUP_DAY="${WEEKLY_BACKUP_DAY:-0}"      # Sunday
VERIFY_TIME="${VERIFY_TIME:-06:00}"              # 6 AM

# User to run backups as
BACKUP_USER="${BACKUP_USER:-listings}"

# Directories
LOG_DIR="${LOG_DIR:-/var/log/listings}"
BACKUP_DIR="${BACKUP_DIR:-/var/backups/listings}"

# Command line flags
DRY_RUN=false
UNINSTALL=false

# ========================================
# Parse command line arguments
# ========================================

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --uninstall)
            UNINSTALL=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--dry-run] [--uninstall]"
            exit 1
            ;;
    esac
done

# ========================================
# Helper functions
# ========================================

log_info() {
    echo "[INFO] $*"
}

log_error() {
    echo "[ERROR] $*" >&2
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

check_user_exists() {
    if ! id "$BACKUP_USER" >/dev/null 2>&1; then
        log_error "User '$BACKUP_USER' does not exist"
        log_info "Create user with: sudo useradd -r -s /bin/bash $BACKUP_USER"
        exit 1
    fi
}

create_directories() {
    log_info "Creating directories..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create directories"
        return 0
    fi

    mkdir -p "$LOG_DIR"
    mkdir -p "$BACKUP_DIR"/{daily,weekly,monthly,wal,pre-restore}
    mkdir -p "$LOG_DIR/reports"

    # Set permissions
    chown -R "$BACKUP_USER:$BACKUP_USER" "$LOG_DIR"
    chown -R "$BACKUP_USER:$BACKUP_USER" "$BACKUP_DIR"
    chmod -R 750 "$BACKUP_DIR"

    log_info "Directories created"
}

setup_logrotate() {
    log_info "Setting up log rotation..."

    local logrotate_conf="/etc/logrotate.d/listings-backup"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create logrotate config: $logrotate_conf"
        return 0
    fi

    cat > "$logrotate_conf" <<'EOF'
/var/log/listings/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 640 listings listings
    sharedscripts
    postrotate
        # Restart services if needed
    endscript
}
EOF

    log_info "Logrotate configured: $logrotate_conf"
}

setup_cron_jobs() {
    log_info "Setting up cron jobs..."

    local cron_file="/etc/cron.d/listings-backup"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create cron file: $cron_file"
        cat <<EOF
# Listings Database Backup Cron Jobs

# Daily backup at ${DAILY_BACKUP_TIME}
$(echo "$DAILY_BACKUP_TIME" | awk -F: '{print $2, $1}') * * * $BACKUP_USER $BACKUP_SCRIPT >> $LOG_DIR/cron.log 2>&1

# Weekly verification on Sunday at ${VERIFY_TIME}
$(echo "$VERIFY_TIME" | awk -F: '{print $2, $1}') * * ${WEEKLY_BACKUP_DAY} $BACKUP_USER $VERIFY_SCRIPT --verify-all >> $LOG_DIR/cron.log 2>&1

# Monitoring check every hour
0 * * * * $BACKUP_USER $MONITOR_SCRIPT >> $LOG_DIR/monitor.log 2>&1
EOF
        return 0
    fi

    # Parse time (HH:MM -> MM HH)
    local daily_min daily_hour
    daily_min=$(echo "$DAILY_BACKUP_TIME" | cut -d: -f2)
    daily_hour=$(echo "$DAILY_BACKUP_TIME" | cut -d: -f1)

    local verify_min verify_hour
    verify_min=$(echo "$VERIFY_TIME" | cut -d: -f2)
    verify_hour=$(echo "$VERIFY_TIME" | cut -d: -f1)

    cat > "$cron_file" <<EOF
# Listings Database Backup Cron Jobs
# Generated: $(date '+%Y-%m-%d %H:%M:%S')

SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

# Daily backup at ${DAILY_BACKUP_TIME}
${daily_min} ${daily_hour} * * * $BACKUP_USER $BACKUP_SCRIPT >> $LOG_DIR/cron.log 2>&1

# Weekly verification on day ${WEEKLY_BACKUP_DAY} at ${VERIFY_TIME}
${verify_min} ${verify_hour} * * ${WEEKLY_BACKUP_DAY} $BACKUP_USER $VERIFY_SCRIPT --verify-all --quick >> $LOG_DIR/cron.log 2>&1

# Monitoring check every hour
0 * * * * $BACKUP_USER $MONITOR_SCRIPT >> $LOG_DIR/monitor.log 2>&1
EOF

    chmod 644 "$cron_file"
    log_info "Cron jobs configured: $cron_file"
}

setup_systemd_timers() {
    log_info "Setting up systemd timers (alternative to cron)..."

    local timer_dir="/etc/systemd/system"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create systemd timers"
        return 0
    fi

    # Backup service
    cat > "${timer_dir}/listings-backup.service" <<EOF
[Unit]
Description=Listings Database Backup
After=network.target

[Service]
Type=oneshot
User=$BACKUP_USER
ExecStart=$BACKUP_SCRIPT
StandardOutput=append:$LOG_DIR/backup.log
StandardError=append:$LOG_DIR/backup.log
EOF

    # Backup timer
    cat > "${timer_dir}/listings-backup.timer" <<EOF
[Unit]
Description=Daily Listings Database Backup
Requires=listings-backup.service

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
EOF

    # Verify service
    cat > "${timer_dir}/listings-backup-verify.service" <<EOF
[Unit]
Description=Listings Backup Verification
After=network.target

[Service]
Type=oneshot
User=$BACKUP_USER
ExecStart=$VERIFY_SCRIPT --verify-all --quick
StandardOutput=append:$LOG_DIR/verify.log
StandardError=append:$LOG_DIR/verify.log
EOF

    # Verify timer
    cat > "${timer_dir}/listings-backup-verify.timer" <<EOF
[Unit]
Description=Weekly Listings Backup Verification
Requires=listings-backup-verify.service

[Timer]
OnCalendar=weekly
Persistent=true

[Install]
WantedBy=timers.target
EOF

    # Reload systemd
    systemctl daemon-reload

    log_info "Systemd timers created (not enabled by default)"
    log_info "To enable: systemctl enable --now listings-backup.timer"
    log_info "To enable: systemctl enable --now listings-backup-verify.timer"
}

test_backup() {
    log_info "Running test backup..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would run test backup"
        return 0
    fi

    if [[ ! -x "$BACKUP_SCRIPT" ]]; then
        log_error "Backup script not executable: $BACKUP_SCRIPT"
        return 1
    fi

    log_info "Testing backup script..."
    if sudo -u "$BACKUP_USER" "$BACKUP_SCRIPT" --dry-run; then
        log_info "Test backup: PASSED"
        return 0
    else
        log_error "Test backup: FAILED"
        return 1
    fi
}

create_env_file() {
    log_info "Creating environment file for backup scripts..."

    local env_file="/etc/listings-backup.env"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create env file: $env_file"
        return 0
    fi

    if [[ -f "$env_file" ]]; then
        log_info "Environment file already exists: $env_file"
        return 0
    fi

    cat > "$env_file" <<'EOF'
# Listings Backup Environment Configuration
# Edit this file to configure backup settings

# Database Configuration
BACKUP_DB_HOST=localhost
BACKUP_DB_PORT=35434
BACKUP_DB_NAME=listings_dev_db
BACKUP_DB_USER=listings_user
BACKUP_DB_PASSWORD=listings_secret
BACKUP_DB_CONTAINER=listings_postgres

# Backup Directories
BACKUP_DIR=/var/backups/listings
LOG_DIR=/var/log/listings

# Retention Policy
BACKUP_RETENTION_DAYS=7
BACKUP_RETENTION_WEEKS=4
BACKUP_RETENTION_MONTHS=12

# S3/MinIO Upload (optional)
BACKUP_ENABLE_S3=false
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=listings-backups
S3_USE_SSL=false

# Email Notifications (optional)
BACKUP_NOTIFY_EMAIL=

# Monitoring (optional)
SLACK_WEBHOOK_URL=
EOF

    chmod 600 "$env_file"
    chown "$BACKUP_USER:$BACKUP_USER" "$env_file"

    log_info "Environment file created: $env_file"
    log_error "IMPORTANT: Edit $env_file and set DB_PASSWORD and other credentials!"
}

make_scripts_executable() {
    log_info "Making scripts executable..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would make scripts executable"
        return 0
    fi

    chmod +x "$BACKUP_SCRIPT"
    chmod +x "$VERIFY_SCRIPT"
    chmod +x "${SCRIPT_DIR}/backup-s3.sh"
    chmod +x "${SCRIPT_DIR}/restore-db.sh"
    chmod +x "$MONITOR_SCRIPT" 2>/dev/null || true

    log_info "Scripts are now executable"
}

uninstall_cron() {
    log_info "Uninstalling cron jobs..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would remove cron jobs"
        return 0
    fi

    local cron_file="/etc/cron.d/listings-backup"

    if [[ -f "$cron_file" ]]; then
        rm -f "$cron_file"
        log_info "Removed cron file: $cron_file"
    fi

    # Remove systemd timers
    systemctl stop listings-backup.timer 2>/dev/null || true
    systemctl stop listings-backup-verify.timer 2>/dev/null || true
    systemctl disable listings-backup.timer 2>/dev/null || true
    systemctl disable listings-backup-verify.timer 2>/dev/null || true

    rm -f /etc/systemd/system/listings-backup.{service,timer}
    rm -f /etc/systemd/system/listings-backup-verify.{service,timer}

    systemctl daemon-reload

    log_info "Cron jobs and timers uninstalled"
}

print_summary() {
    cat <<EOF

========================================
Backup System Setup Complete
========================================

Scripts Location: $SCRIPT_DIR
Backup Directory: $BACKUP_DIR
Log Directory: $LOG_DIR

Cron Schedule:
  - Daily backup: ${DAILY_BACKUP_TIME}
  - Weekly verify: Day ${WEEKLY_BACKUP_DAY} at ${VERIFY_TIME}
  - Hourly monitoring

Configuration:
  - Edit: /etc/listings-backup.env
  - Set BACKUP_DB_PASSWORD and other credentials

Next Steps:
  1. Edit /etc/listings-backup.env with your configuration
  2. Test backup: sudo -u $BACKUP_USER $BACKUP_SCRIPT --dry-run
  3. Monitor logs: tail -f $LOG_DIR/backup.log

Manual Commands:
  - Backup:  sudo -u $BACKUP_USER $BACKUP_SCRIPT
  - Restore: sudo -u $BACKUP_USER $BACKUP_SCRIPT/restore-db.sh <file>
  - Verify:  sudo -u $BACKUP_USER $VERIFY_SCRIPT <file>
  - Monitor: sudo -u $BACKUP_USER $MONITOR_SCRIPT

========================================
EOF
}

# ========================================
# Main setup logic
# ========================================

main() {
    log_info "========================================="
    log_info "Listings Backup System Setup"
    log_info "========================================="

    # Check prerequisites
    check_root

    # Uninstall mode
    if [[ "$UNINSTALL" == "true" ]]; then
        uninstall_cron
        log_info "Uninstall completed"
        exit 0
    fi

    # Check user exists
    check_user_exists

    # Setup
    create_directories
    make_scripts_executable
    create_env_file
    setup_logrotate
    setup_cron_jobs
    setup_systemd_timers

    # Test
    if ! test_backup; then
        log_error "Test backup failed. Please check configuration."
        exit 1
    fi

    # Summary
    print_summary
}

# ========================================
# Entry point
# ========================================

main "$@"
