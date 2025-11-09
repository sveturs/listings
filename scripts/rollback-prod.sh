#!/usr/bin/env bash

#######################################
# Production Rollback Script
# Emergency rollback to previous stable version
#
# Usage: ./rollback-prod.sh [--reason "error description"] [--restore-db] [--dry-run]
#
# Exit codes:
#   0 - Success
#   1 - Rollback failed
#######################################

set -Eeuo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Import configuration
if [[ -f "${SCRIPT_DIR}/.env.deploy" ]]; then
    set -a
    source "${SCRIPT_DIR}/.env.deploy"
    set +a
else
    echo "ERROR: .env.deploy not found"
    exit 1
fi

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Flags
DRY_RUN=false
RESTORE_DB=false
ROLLBACK_REASON="manual-rollback"

# Rollback metadata
ROLLBACK_ID="rollback-$(date +%Y%m%d-%H%M%S)"
ROLLBACK_LOG="${PROJECT_ROOT}/logs/rollbacks/${ROLLBACK_ID}.log"
START_TIME=$(date +%s)

# Cleanup function
cleanup() {
    local exit_code=$?
    local end_time=$(date +%s)
    local duration=$((end_time - START_TIME))

    if [[ $exit_code -eq 0 ]]; then
        log_success "Rollback completed successfully in ${duration}s"
        send_notification "✅ Rollback ${ROLLBACK_ID} succeeded" "Reason: ${ROLLBACK_REASON}, Duration: ${duration}s"
    else
        log_error "Rollback failed with exit code ${exit_code} after ${duration}s"
        send_notification "❌ Rollback ${ROLLBACK_ID} failed" "Exit code: ${exit_code}"
    fi
}

trap cleanup EXIT

#######################################
# Logging functions
#######################################

log_info() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [INFO] $*"
    echo -e "${BLUE}${msg}${NC}"
    echo "${msg}" >> "${ROLLBACK_LOG}"
}

log_success() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [SUCCESS] $*"
    echo -e "${GREEN}${msg}${NC}"
    echo "${msg}" >> "${ROLLBACK_LOG}"
}

log_warning() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [WARNING] $*"
    echo -e "${YELLOW}${msg}${NC}"
    echo "${msg}" >> "${ROLLBACK_LOG}"
}

log_error() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [ERROR] $*"
    echo -e "${RED}${msg}${NC}" >&2
    echo "${msg}" >> "${ROLLBACK_LOG}"
}

#######################################
# Notification function
#######################################

send_notification() {
    local title="$1"
    local message="$2"

    if [[ -n "${SLACK_WEBHOOK_URL:-}" ]]; then
        local payload=$(cat <<EOF
{
    "text": "🚨 ${title}",
    "blocks": [
        {
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": "*${title}*\n${message}\n*Rollback ID:* ${ROLLBACK_ID}\n*Server:* ${PROD_HOST}"
            }
        }
    ]
}
EOF
)
        curl -s -X POST -H 'Content-type: application/json' \
            --data "${payload}" \
            "${SLACK_WEBHOOK_URL}" > /dev/null 2>&1 || true
    fi

    # Also send email if configured
    if [[ -n "${ALERT_EMAIL:-}" ]]; then
        echo -e "Subject: [URGENT] Production Rollback - ${ROLLBACK_ID}\n\n${title}\n${message}" | \
            sendmail "${ALERT_EMAIL}" 2>/dev/null || true
    fi
}

#######################################
# Parse arguments
#######################################

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --reason)
                ROLLBACK_REASON="$2"
                shift 2
                ;;
            --restore-db)
                RESTORE_DB=true
                log_warning "Database will be restored from backup"
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                log_warning "DRY RUN MODE - No actual changes will be made"
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                echo "Usage: $0 [--reason \"error description\"] [--restore-db] [--dry-run]"
                exit 1
                ;;
        esac
    done
}

#######################################
# Confirm rollback
#######################################

confirm_rollback() {
    log_warning "========================================="
    log_warning "PRODUCTION ROLLBACK"
    log_warning "========================================="
    log_warning "This will switch traffic back to Blue (previous version)"
    log_warning "Reason: ${ROLLBACK_REASON}"
    log_warning "Restore DB: ${RESTORE_DB}"
    log_warning ""

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Skipping confirmation"
        return 0
    fi

    read -p "Are you sure you want to proceed? (type 'yes' to confirm): " confirmation
    if [[ "${confirmation}" != "yes" ]]; then
        log_info "Rollback cancelled by user"
        exit 0
    fi
}

#######################################
# Switch traffic to Blue
#######################################

switch_to_blue() {
    log_info "Switching all traffic to Blue (previous version)..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would switch 100% traffic to Blue"
        return 0
    fi

    # Immediate switch - no gradual rollback
    "${SCRIPT_DIR}/traffic-split.sh" --green-weight 0

    if [[ $? -eq 0 ]]; then
        log_success "Traffic switched to Blue"
    else
        log_error "Failed to switch traffic"
        exit 1
    fi

    # Wait a moment for connections to drain
    sleep 5
}

#######################################
# Stop Green instance
#######################################

stop_green() {
    log_info "Stopping Green instance..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would stop Green instance"
        return 0
    fi

    ssh "${PROD_USER}@${PROD_HOST}" << 'EOF'
cd ${PROD_DIR}

# Archive Green logs for postmortem
if [[ -f "logs/green.log" ]]; then
    gzip -c logs/green.log > logs/incidents/green-rollback-$(date +%Y%m%d-%H%M%S).log.gz
    echo "Green logs archived for postmortem"
fi

# Stop Green instance
if pgrep -f "listings-green" > /dev/null; then
    pkill -f "listings-green"
    echo "Green instance stopped"
else
    echo "Green instance was not running"
fi
EOF

    log_success "Green instance stopped"
}

#######################################
# Verify Blue is healthy
#######################################

verify_blue() {
    log_info "Verifying Blue instance is healthy..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would verify Blue health"
        return 0
    fi

    # Check if Blue is running
    ssh "${PROD_USER}@${PROD_HOST}" "pgrep -f 'listings-blue' > /dev/null"
    if [[ $? -ne 0 ]]; then
        log_error "Blue instance is not running!"
        log_info "Starting Blue instance..."

        ssh "${PROD_USER}@${PROD_HOST}" << 'EOF'
cd ${PROD_DIR}
nohup ./bin/listings-blue --port=${BLUE_PORT} --env-file=.env.blue > logs/blue.log 2>&1 &
echo $! > blue.pid
sleep 5
EOF
    fi

    # Health check
    local retries=5
    local count=0

    while [[ $count -lt $retries ]]; do
        if curl -sf "http://${PROD_HOST}:${BLUE_PORT}/health" > /dev/null 2>&1; then
            log_success "Blue instance is healthy"
            return 0
        fi

        count=$((count + 1))
        log_warning "Health check attempt ${count}/${retries} failed, retrying..."
        sleep 2
    done

    log_error "Blue instance health check failed after ${retries} attempts"
    exit 1
}

#######################################
# Restore database backup
#######################################

restore_database() {
    if [[ "${RESTORE_DB}" != "true" ]]; then
        log_info "Skipping database restore (--restore-db not specified)"
        return 0
    fi

    log_warning "Restoring database from backup..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would restore latest database backup"
        return 0
    fi

    # Find latest backup
    local latest_backup=$(ssh "${PROD_USER}@${PROD_HOST}" \
        "cd ${PROD_DIR}/backups && ls -t backup-before-*.sql.gz 2>/dev/null | head -n1")

    if [[ -z "${latest_backup}" ]]; then
        log_error "No backup file found"
        exit 1
    fi

    log_info "Restoring backup: ${latest_backup}"

    # Confirm database restore
    read -p "This will OVERWRITE the current database. Type 'RESTORE' to confirm: " db_confirm
    if [[ "${db_confirm}" != "RESTORE" ]]; then
        log_info "Database restore cancelled"
        return 0
    fi

    # Restore database
    ssh "${PROD_USER}@${PROD_HOST}" << EOF
cd ${PROD_DIR}/backups
gunzip -c ${latest_backup} | docker exec -i listings_postgres psql -U ${DB_USER} -d ${DB_NAME}
EOF

    if [[ $? -eq 0 ]]; then
        log_success "Database restored from ${latest_backup}"
    else
        log_error "Database restore failed"
        exit 1
    fi
}

#######################################
# Capture incident logs
#######################################

capture_logs() {
    log_info "Capturing logs for postmortem analysis..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would capture logs"
        return 0
    fi

    local incident_dir="incidents/${ROLLBACK_ID}"
    mkdir -p "${PROJECT_ROOT}/${incident_dir}"

    # Download Green logs
    scp "${PROD_USER}@${PROD_HOST}:${PROD_DIR}/logs/incidents/green-rollback-*.log.gz" \
        "${PROJECT_ROOT}/${incident_dir}/" 2>/dev/null || true

    # Download nginx error logs
    ssh "${PROD_USER}@${PROD_HOST}" \
        "sudo tail -n 1000 /var/log/nginx/error.log" > "${PROJECT_ROOT}/${incident_dir}/nginx-error.log" || true

    # Download system logs
    ssh "${PROD_USER}@${PROD_HOST}" \
        "journalctl -u listings -n 1000" > "${PROJECT_ROOT}/${incident_dir}/systemd.log" || true

    log_success "Logs captured in ${incident_dir}/"
}

#######################################
# Create incident report
#######################################

create_incident_report() {
    log_info "Creating incident report..."

    local report_file="${PROJECT_ROOT}/incidents/${ROLLBACK_ID}/INCIDENT_REPORT.md"

    cat > "${report_file}" << EOF
# Incident Report - ${ROLLBACK_ID}

## Incident Details
- **Rollback ID:** ${ROLLBACK_ID}
- **Timestamp:** $(date +'%Y-%m-%d %H:%M:%S %Z')
- **Reason:** ${ROLLBACK_REASON}
- **Database Restored:** ${RESTORE_DB}

## Actions Taken
1. Switched traffic from Green to Blue (100%)
2. Stopped Green instance
3. Verified Blue instance health
$(if [[ "${RESTORE_DB}" == "true" ]]; then echo "4. Restored database from backup"; fi)
5. Captured logs for analysis

## Next Steps
- [ ] Analyze Green logs in: incidents/${ROLLBACK_ID}/
- [ ] Identify root cause
- [ ] Create bug ticket
- [ ] Fix issue
- [ ] Plan redeployment

## Postmortem Checklist
- [ ] Timeline of events
- [ ] Root cause analysis
- [ ] Impact assessment
- [ ] Action items
- [ ] Prevention measures

## Monitoring
- Check error rates: \`tail -f ${PROD_DIR}/logs/blue.log | grep ERROR\`
- Check response times: \`curl -w "@curl-format.txt" http://${PROD_HOST}/health\`

## Contacts
- On-call Engineer: ${ONCALL_ENGINEER:-TBD}
- Incident Commander: ${INCIDENT_COMMANDER:-TBD}

---
*Generated automatically by rollback-prod.sh*
EOF

    log_success "Incident report created: ${report_file}"
}

#######################################
# Main rollback flow
#######################################

main() {
    log_warning "========================================="
    log_warning "EMERGENCY PRODUCTION ROLLBACK"
    log_warning "Rollback ID: ${ROLLBACK_ID}"
    log_warning "Target: ${PROD_HOST}"
    log_warning "========================================="

    # Create directories
    mkdir -p "${PROJECT_ROOT}/logs/rollbacks"
    mkdir -p "${PROJECT_ROOT}/incidents/${ROLLBACK_ID}"

    # Parse arguments
    parse_args "$@"

    # Confirm rollback
    confirm_rollback

    # Send initial notification
    send_notification "🚨 Rollback initiated" "Reason: ${ROLLBACK_REASON}"

    # Step 1: Switch traffic to Blue immediately
    switch_to_blue

    # Step 2: Verify Blue is healthy
    verify_blue

    # Step 3: Stop Green instance
    stop_green

    # Step 4: Restore database if requested
    restore_database

    # Step 5: Capture logs for postmortem
    capture_logs

    # Step 6: Create incident report
    create_incident_report

    log_success "========================================="
    log_success "Rollback completed successfully!"
    log_success "Review incident report: incidents/${ROLLBACK_ID}/INCIDENT_REPORT.md"
    log_success "========================================="
}

# Run main function
main "$@"
