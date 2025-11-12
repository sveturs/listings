#!/usr/bin/env bash

#######################################
# Production Deployment Script (Blue-Green Strategy)
# Zero-downtime deployment with automated rollback
#
# Usage: ./deploy-to-prod.sh [--dry-run] [--verbose] [--skip-tests]
#
# Prerequisites:
#   - .env.deploy file configured
#   - SSH access to production server
#   - Validated pre-deployment checks
#
# Exit codes:
#   0 - Success
#   1 - Pre-deployment validation failed
#   2 - Build failed
#   3 - Upload failed
#   4 - Smoke tests failed
#   5 - Traffic switch failed
#   6 - Rollback executed
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
    echo "ERROR: .env.deploy not found. Copy .env.deploy.example and configure it."
    exit 1
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Flags
DRY_RUN=false
VERBOSE=false
SKIP_TESTS=false

# Lock file
LOCK_FILE="${PROJECT_ROOT}/.deployment.lock"

# Deployment metadata
DEPLOYMENT_ID="deploy-$(date +%Y%m%d-%H%M%S)"
DEPLOYMENT_LOG="${PROJECT_ROOT}/logs/deployments/${DEPLOYMENT_ID}.log"
START_TIME=$(date +%s)

# Cleanup function (called on any exit)
cleanup() {
    local exit_code=$?

    log_info "Cleaning up..."

    # Remove lock file
    if [[ -f "${LOCK_FILE}" ]]; then
        rm -f "${LOCK_FILE}"
    fi

    # Calculate deployment time
    local end_time=$(date +%s)
    local duration=$((end_time - START_TIME))

    if [[ $exit_code -eq 0 ]]; then
        log_success "Deployment completed successfully in ${duration}s"
        send_notification "✅ Deployment ${DEPLOYMENT_ID} succeeded" "Duration: ${duration}s"
    else
        log_error "Deployment failed with exit code ${exit_code} after ${duration}s"
        send_notification "❌ Deployment ${DEPLOYMENT_ID} failed" "Exit code: ${exit_code}, Duration: ${duration}s"
    fi
}

trap cleanup EXIT

#######################################
# Logging functions
#######################################

log_info() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [INFO] $*"
    echo -e "${BLUE}${msg}${NC}"
    echo "${msg}" >> "${DEPLOYMENT_LOG}"
}

log_success() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [SUCCESS] $*"
    echo -e "${GREEN}${msg}${NC}"
    echo "${msg}" >> "${DEPLOYMENT_LOG}"
}

log_warning() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [WARNING] $*"
    echo -e "${YELLOW}${msg}${NC}"
    echo "${msg}" >> "${DEPLOYMENT_LOG}"
}

log_error() {
    local msg="[$(date +'%Y-%m-%d %H:%M:%S')] [ERROR] $*"
    echo -e "${RED}${msg}${NC}" >&2
    echo "${msg}" >> "${DEPLOYMENT_LOG}"
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
    "text": "${title}",
    "blocks": [
        {
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": "*${title}*\n${message}\n*Deployment ID:* ${DEPLOYMENT_ID}\n*Server:* ${PROD_HOST}"
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
}

#######################################
# Parse command line arguments
#######################################

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --dry-run)
                DRY_RUN=true
                log_warning "DRY RUN MODE - No actual changes will be made"
                shift
                ;;
            --verbose)
                VERBOSE=true
                set -x
                shift
                ;;
            --skip-tests)
                SKIP_TESTS=true
                log_warning "Skipping smoke tests (not recommended for production)"
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                echo "Usage: $0 [--dry-run] [--verbose] [--skip-tests]"
                exit 1
                ;;
        esac
    done
}

#######################################
# Check lock file (prevent concurrent deployments)
#######################################

check_lock() {
    if [[ -f "${LOCK_FILE}" ]]; then
        local lock_info=$(cat "${LOCK_FILE}")
        log_error "Another deployment is in progress: ${lock_info}"
        log_error "If this is a stale lock, remove ${LOCK_FILE} manually"
        exit 1
    fi

    # Create lock file
    echo "${DEPLOYMENT_ID} started at $(date)" > "${LOCK_FILE}"
}

#######################################
# Pre-deployment validation
#######################################

validate_environment() {
    log_info "Running pre-deployment validation..."

    if ! "${SCRIPT_DIR}/validate-deployment.sh"; then
        log_error "Pre-deployment validation failed"
        exit 1
    fi

    log_success "Pre-deployment validation passed"
}

#######################################
# Create backup
#######################################

create_backup() {
    log_info "Creating database backup..."

    local backup_name="backup-before-${DEPLOYMENT_ID}.sql"

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would create backup: ${backup_name}"
        return 0
    fi

    ssh "${PROD_USER}@${PROD_HOST}" "cd ${PROD_DIR} && docker exec listings_postgres pg_dump -U ${DB_USER} ${DB_NAME} | gzip > backups/${backup_name}.gz"

    if [[ $? -eq 0 ]]; then
        log_success "Backup created: ${backup_name}.gz"
    else
        log_error "Backup creation failed"
        exit 1
    fi
}

#######################################
# Build binary
#######################################

build_binary() {
    log_info "Building production binary..."

    cd "${PROJECT_ROOT}"

    # Get version
    local version=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    local build_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    log_info "Version: ${version}, Commit: ${commit}, Build time: ${build_time}"

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would build binary with CGO_ENABLED=0"
        return 0
    fi

    # Build with optimizations
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-s -w -X main.Version=${version} -X main.Commit=${commit} -X main.BuildTime=${build_time}" \
        -o bin/listings-prod \
        ./cmd/server/main.go

    if [[ $? -eq 0 ]]; then
        local size=$(du -h bin/listings-prod | cut -f1)
        log_success "Binary built successfully (${size})"
    else
        log_error "Build failed"
        exit 2
    fi
}

#######################################
# Upload binary to server
#######################################

upload_binary() {
    log_info "Uploading binary to production server..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would upload binary to ${PROD_HOST}:${PROD_DIR}/bin/listings-green"
        return 0
    fi

    # Upload to green environment
    scp -C "${PROJECT_ROOT}/bin/listings-prod" \
        "${PROD_USER}@${PROD_HOST}:${PROD_DIR}/bin/listings-green"

    if [[ $? -eq 0 ]]; then
        log_success "Binary uploaded successfully"
    else
        log_error "Upload failed"
        exit 3
    fi
}

#######################################
# Upload .env.prod file
#######################################

upload_env_file() {
    log_info "Uploading .env.prod configuration..."

    if [[ ! -f "${PROJECT_ROOT}/.env.prod" ]]; then
        log_error ".env.prod file not found"
        exit 1
    fi

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would upload .env.prod to server"
        return 0
    fi

    scp "${PROJECT_ROOT}/.env.prod" \
        "${PROD_USER}@${PROD_HOST}:${PROD_DIR}/.env.green"

    if [[ $? -eq 0 ]]; then
        log_success "Configuration uploaded successfully"
    else
        log_error "Configuration upload failed"
        exit 3
    fi
}

#######################################
# Run database migrations
#######################################

run_migrations() {
    log_info "Running database migrations..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would run migrations on Green environment"
        return 0
    fi

    # Run migrations using Green config (Blue still serving traffic)
    ssh "${PROD_USER}@${PROD_HOST}" "cd ${PROD_DIR} && ./bin/listings-green migrate up --env-file=.env.green"

    if [[ $? -eq 0 ]]; then
        log_success "Migrations completed successfully"
    else
        log_error "Migrations failed"
        log_info "Consider manual rollback: ./rollback-prod.sh"
        exit 1
    fi
}

#######################################
# Start Green instance
#######################################

start_green() {
    log_info "Starting Green instance..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would start Green instance on port ${GREEN_PORT}"
        return 0
    fi

    # Start Green instance with separate config
    ssh "${PROD_USER}@${PROD_HOST}" << 'EOF'
cd ${PROD_DIR}

# Stop existing Green instance if running
if pgrep -f "listings-green" > /dev/null; then
    pkill -f "listings-green"
    sleep 2
fi

# Start new Green instance
nohup ./bin/listings-green --port=${GREEN_PORT} --env-file=.env.green > logs/green.log 2>&1 &

echo $! > green.pid

sleep 5

# Verify it started
if pgrep -f "listings-green" > /dev/null; then
    echo "Green instance started successfully"
    exit 0
else
    echo "Failed to start Green instance"
    exit 1
fi
EOF

    if [[ $? -eq 0 ]]; then
        log_success "Green instance started on port ${GREEN_PORT}"
    else
        log_error "Failed to start Green instance"
        exit 1
    fi
}

#######################################
# Run smoke tests on Green
#######################################

run_smoke_tests() {
    if [[ "${SKIP_TESTS}" == "true" ]]; then
        log_warning "Skipping smoke tests (--skip-tests flag)"
        return 0
    fi

    log_info "Running smoke tests on Green instance..."

    if [[ "${DRY_RUN}" == "true" ]]; then
        log_info "[DRY RUN] Would run smoke tests against port ${GREEN_PORT}"
        return 0
    fi

    # Run smoke tests against Green port
    if "${SCRIPT_DIR}/smoke-tests.sh" --host "${PROD_HOST}" --port "${GREEN_PORT}"; then
        log_success "Smoke tests passed"
    else
        log_error "Smoke tests failed on Green instance"
        log_info "Rolling back..."

        # Stop Green instance
        ssh "${PROD_USER}@${PROD_HOST}" "pkill -f 'listings-green'"

        exit 4
    fi
}

#######################################
# Canary deployment (gradual traffic shift)
#######################################

canary_deployment() {
    log_info "Starting canary deployment (gradual traffic shift)..."

    # Phase 1: 10% to Green
    log_info "Phase 1: Routing 10% traffic to Green..."
    if [[ "${DRY_RUN}" == "false" ]]; then
        "${SCRIPT_DIR}/traffic-split.sh" --green-weight 10
        sleep 300 # Wait 5 minutes

        # Monitor for errors
        check_error_rate
        if [[ $? -ne 0 ]]; then
            log_error "High error rate detected in Phase 1, rolling back..."
            "${SCRIPT_DIR}/rollback-prod.sh" --reason "canary-phase1-errors"
            exit 5
        fi
    else
        log_info "[DRY RUN] Would route 10% to Green and wait 5 minutes"
    fi

    # Phase 2: 50% to Green
    log_info "Phase 2: Routing 50% traffic to Green..."
    if [[ "${DRY_RUN}" == "false" ]]; then
        "${SCRIPT_DIR}/traffic-split.sh" --green-weight 50
        sleep 300 # Wait 5 minutes

        check_error_rate
        if [[ $? -ne 0 ]]; then
            log_error "High error rate detected in Phase 2, rolling back..."
            "${SCRIPT_DIR}/rollback-prod.sh" --reason "canary-phase2-errors"
            exit 5
        fi
    else
        log_info "[DRY RUN] Would route 50% to Green and wait 5 minutes"
    fi

    # Phase 3: 100% to Green
    log_info "Phase 3: Routing 100% traffic to Green..."
    if [[ "${DRY_RUN}" == "false" ]]; then
        "${SCRIPT_DIR}/traffic-split.sh" --green-weight 100
        sleep 60 # Wait 1 minute
    else
        log_info "[DRY RUN] Would route 100% to Green"
    fi

    log_success "Canary deployment completed successfully"
}

#######################################
# Check error rate
#######################################

check_error_rate() {
    local error_count=$(ssh "${PROD_USER}@${PROD_HOST}" \
        "tail -n 1000 ${PROD_DIR}/logs/green.log | grep -i 'error\|panic\|fatal' | wc -l")

    if [[ ${error_count} -gt 10 ]]; then
        log_error "High error rate detected: ${error_count} errors in last 1000 lines"
        return 1
    fi

    return 0
}

#######################################
# Decommission Blue instance
#######################################

decommission_blue() {
    log_info "Waiting 10 minutes before decommissioning Blue (rollback window)..."

    if [[ "${DRY_RUN}" == "false" ]]; then
        sleep 600 # 10 minutes
    else
        log_info "[DRY RUN] Would wait 10 minutes"
    fi

    log_info "Decommissioning Blue instance..."

    if [[ "${DRY_RUN}" == "false" ]]; then
        ssh "${PROD_USER}@${PROD_HOST}" << 'EOF'
cd ${PROD_DIR}

# Archive Blue logs
if [[ -f "logs/blue.log" ]]; then
    gzip -c logs/blue.log > logs/archived/blue-$(date +%Y%m%d-%H%M%S).log.gz
fi

# Stop Blue instance
if pgrep -f "listings-blue" > /dev/null; then
    pkill -f "listings-blue"
    echo "Blue instance stopped"
fi

# Promote Green to Blue for next deployment
mv bin/listings-blue bin/listings-blue.old 2>/dev/null || true
cp bin/listings-green bin/listings-blue
mv .env.blue .env.blue.old 2>/dev/null || true
cp .env.green .env.blue

echo "Green promoted to Blue"
EOF

        log_success "Blue instance decommissioned and Green promoted"
    else
        log_info "[DRY RUN] Would stop Blue and promote Green to Blue"
    fi
}

#######################################
# Generate deployment report
#######################################

generate_report() {
    log_info "Generating deployment report..."

    if [[ "${DRY_RUN}" == "false" ]]; then
        "${SCRIPT_DIR}/deployment-report.sh" --deployment-id "${DEPLOYMENT_ID}"
    else
        log_info "[DRY RUN] Would generate deployment report"
    fi
}

#######################################
# Main deployment flow
#######################################

main() {
    log_info "========================================="
    log_info "Starting Production Deployment"
    log_info "Deployment ID: ${DEPLOYMENT_ID}"
    log_info "Target: ${PROD_HOST}"
    log_info "========================================="

    # Create logs directory
    mkdir -p "${PROJECT_ROOT}/logs/deployments"

    # Parse arguments
    parse_args "$@"

    # Check for concurrent deployments
    check_lock

    # Step 1: Pre-deployment checks
    validate_environment

    # Step 2: Create backup
    create_backup

    # Step 3: Build binary
    build_binary

    # Step 4: Upload binary and config
    upload_binary
    upload_env_file

    # Step 5: Run migrations (Blue still serving)
    run_migrations

    # Step 6: Start Green instance
    start_green

    # Step 7: Smoke tests on Green
    run_smoke_tests

    # Step 8: Canary deployment (gradual traffic shift)
    canary_deployment

    # Step 9: Decommission Blue
    decommission_blue

    # Step 10: Generate report
    generate_report

    log_success "========================================="
    log_success "Deployment completed successfully!"
    log_success "========================================="
}

# Run main function
main "$@"
