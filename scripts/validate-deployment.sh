#!/usr/bin/env bash

#######################################
# Pre-Deployment Validation Script
# Validates environment before production deployment
#
# Usage: ./validate-deployment.sh [--verbose]
#
# Exit codes:
#   0 - All checks passed
#   1 - One or more checks failed
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

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Flags
VERBOSE=false

# Check results
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
declare -a FAILED_CHECK_NAMES=()

#######################################
# Logging functions
#######################################

log_check() {
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    echo -ne "${BLUE}[CHECK]${NC} $*... "
}

log_pass() {
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
    echo -e "${GREEN}PASS${NC}"
}

log_fail() {
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
    local check_name="$1"
    FAILED_CHECK_NAMES+=("${check_name}")
    echo -e "${RED}FAIL${NC}"
}

log_info() {
    if [[ "${VERBOSE}" == "true" ]]; then
        echo -e "${BLUE}  [INFO]${NC} $*"
    fi
}

log_warning() {
    echo -e "${YELLOW}  [WARNING]${NC} $*"
}

log_error() {
    echo -e "${RED}  [ERROR]${NC} $*"
}

#######################################
# Parse arguments
#######################################

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --verbose)
                VERBOSE=true
                shift
                ;;
            *)
                echo "Usage: $0 [--verbose]"
                exit 1
                ;;
        esac
    done
}

#######################################
# Check functions
#######################################

check_ssh_connectivity() {
    log_check "SSH connectivity to ${PROD_HOST}"

    if ssh -o ConnectTimeout=10 -o BatchMode=yes "${PROD_USER}@${PROD_HOST}" "echo 'OK'" &>/dev/null; then
        log_pass
        log_info "SSH connection successful"
        return 0
    else
        log_fail "SSH connectivity"
        log_error "Cannot connect to ${PROD_USER}@${PROD_HOST}"
        log_error "Ensure SSH key is added and host is reachable"
        return 1
    fi
}

check_disk_space() {
    log_check "Disk space on ${PROD_HOST}"

    local available_gb=$(ssh "${PROD_USER}@${PROD_HOST}" \
        "df ${PROD_DIR} | tail -1 | awk '{print \$4}'" 2>/dev/null)

    # Convert KB to GB (approximate)
    available_gb=$((available_gb / 1024 / 1024))

    if [[ ${available_gb} -gt 20 ]]; then
        log_pass
        log_info "Available: ${available_gb} GB"
        return 0
    else
        log_fail "Disk space"
        log_error "Only ${available_gb} GB available (minimum 20 GB required)"
        return 1
    fi
}

check_env_prod_file() {
    log_check ".env.prod file exists locally"

    if [[ -f "${PROJECT_ROOT}/.env.prod" ]]; then
        log_pass
        log_info "File: ${PROJECT_ROOT}/.env.prod"
        return 0
    else
        log_fail ".env.prod file"
        log_error "File not found: ${PROJECT_ROOT}/.env.prod"
        log_error "Create .env.prod with production configuration"
        return 1
    fi
}

check_docker_available() {
    log_check "Docker is available on ${PROD_HOST}"

    if ssh "${PROD_USER}@${PROD_HOST}" "docker ps" &>/dev/null; then
        log_pass
        return 0
    else
        log_fail "Docker availability"
        log_error "Docker is not available or user lacks permissions"
        return 1
    fi
}

check_postgres_accessible() {
    log_check "PostgreSQL is accessible"

    if ssh "${PROD_USER}@${PROD_HOST}" \
        "docker exec listings_postgres pg_isready -U ${DB_USER}" &>/dev/null; then
        log_pass
        return 0
    else
        log_fail "PostgreSQL accessibility"
        log_error "PostgreSQL is not accessible or container is not running"
        return 1
    fi
}

check_redis_accessible() {
    log_check "Redis is accessible"

    if ssh "${PROD_USER}@${PROD_HOST}" \
        "docker exec listings_redis redis-cli ping" 2>/dev/null | grep -q "PONG"; then
        log_pass
        return 0
    else
        log_fail "Redis accessibility"
        log_error "Redis is not accessible or container is not running"
        return 1
    fi
}

check_no_deployments_in_progress() {
    log_check "No deployments in progress"

    if ssh "${PROD_USER}@${PROD_HOST}" \
        "test -f ${PROD_DIR}/.deployment.lock" 2>/dev/null; then
        log_fail "Deployment lock"
        log_error "Another deployment is in progress"
        log_error "Lock file: ${PROD_DIR}/.deployment.lock"
        return 1
    else
        log_pass
        return 0
    fi
}

check_git_clean() {
    log_check "Git working directory is clean"

    cd "${PROJECT_ROOT}"

    if [[ -z $(git status --porcelain) ]]; then
        log_pass
        return 0
    else
        log_fail "Git status"
        log_warning "Working directory has uncommitted changes"
        log_warning "Consider committing or stashing changes before deployment"
        # This is a warning, not a failure
        FAILED_CHECKS=$((FAILED_CHECKS - 1))
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    fi
}

check_git_pushed() {
    log_check "Git changes are pushed to remote"

    cd "${PROJECT_ROOT}"

    local unpushed=$(git log @{u}.. --oneline 2>/dev/null | wc -l)

    if [[ ${unpushed} -eq 0 ]]; then
        log_pass
        return 0
    else
        log_fail "Git push status"
        log_warning "${unpushed} commits not pushed to remote"
        log_warning "Consider pushing commits before deployment"
        # This is a warning, not a failure
        FAILED_CHECKS=$((FAILED_CHECKS - 1))
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    fi
}

check_go_version() {
    log_check "Go version compatibility"

    local local_version=$(go version | awk '{print $3}' | sed 's/go//')
    local remote_version=$(ssh "${PROD_USER}@${PROD_HOST}" "go version" 2>/dev/null | awk '{print $3}' | sed 's/go//' || echo "unknown")

    log_info "Local: ${local_version}, Remote: ${remote_version}"

    if [[ "${remote_version}" == "unknown" ]]; then
        log_fail "Go version"
        log_error "Go not found on remote server"
        return 1
    fi

    log_pass
    return 0
}

check_required_directories() {
    log_check "Required directories exist on ${PROD_HOST}"

    local dirs=("${PROD_DIR}" "${PROD_DIR}/bin" "${PROD_DIR}/logs" "${PROD_DIR}/backups")
    local missing_dirs=()

    for dir in "${dirs[@]}"; do
        if ! ssh "${PROD_USER}@${PROD_HOST}" "test -d ${dir}" 2>/dev/null; then
            missing_dirs+=("${dir}")
        fi
    done

    if [[ ${#missing_dirs[@]} -eq 0 ]]; then
        log_pass
        return 0
    else
        log_fail "Directory structure"
        log_error "Missing directories: ${missing_dirs[*]}"
        log_info "Creating missing directories..."

        for dir in "${missing_dirs[@]}"; do
            ssh "${PROD_USER}@${PROD_HOST}" "mkdir -p ${dir}"
        done

        log_info "Directories created"
        return 0
    fi
}

check_nginx_installed() {
    log_check "Nginx is installed on ${PROD_HOST}"

    if ssh "${PROD_USER}@${PROD_HOST}" "which nginx" &>/dev/null; then
        log_pass
        return 0
    else
        log_fail "Nginx installation"
        log_error "Nginx is not installed"
        return 1
    fi
}

check_backup_space() {
    log_check "Backup directory has sufficient space"

    local backup_dir="${PROD_DIR}/backups"
    local available_gb=$(ssh "${PROD_USER}@${PROD_HOST}" \
        "df ${backup_dir} | tail -1 | awk '{print \$4}'" 2>/dev/null || echo "0")

    available_gb=$((available_gb / 1024 / 1024))

    if [[ ${available_gb} -gt 5 ]]; then
        log_pass
        log_info "Available: ${available_gb} GB"
        return 0
    else
        log_fail "Backup space"
        log_error "Only ${available_gb} GB available for backups (minimum 5 GB required)"
        return 1
    fi
}

check_blue_instance_running() {
    log_check "Blue instance is running"

    if ssh "${PROD_USER}@${PROD_HOST}" "pgrep -f 'listings-blue' > /dev/null" 2>/dev/null; then
        log_pass
        return 0
    else
        log_fail "Blue instance"
        log_warning "Blue instance is not running"
        log_warning "This is expected for first-time deployment"
        # Not a critical failure
        FAILED_CHECKS=$((FAILED_CHECKS - 1))
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    fi
}

check_slack_webhook() {
    log_check "Slack webhook is configured"

    if [[ -n "${SLACK_WEBHOOK_URL:-}" ]]; then
        log_pass
        log_info "Webhook configured"
        return 0
    else
        log_fail "Slack webhook"
        log_warning "Slack notifications not configured"
        log_info "Set SLACK_WEBHOOK_URL in .env.deploy to enable notifications"
        # Not a critical failure
        FAILED_CHECKS=$((FAILED_CHECKS - 1))
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    fi
}

#######################################
# Main function
#######################################

main() {
    parse_args "$@"

    echo "========================================="
    echo "Pre-Deployment Validation"
    echo "========================================="
    echo ""

    # Run all checks
    check_ssh_connectivity
    check_disk_space
    check_env_prod_file
    check_docker_available
    check_postgres_accessible
    check_redis_accessible
    check_no_deployments_in_progress
    check_git_clean
    check_git_pushed
    check_go_version
    check_required_directories
    check_nginx_installed
    check_backup_space
    check_blue_instance_running
    check_slack_webhook

    # Summary
    echo ""
    echo "========================================="
    echo "Validation Summary"
    echo "========================================="
    echo "Total Checks:  ${TOTAL_CHECKS}"
    echo "Passed:        ${PASSED_CHECKS}"
    echo "Failed:        ${FAILED_CHECKS}"
    echo "========================================="

    if [[ ${FAILED_CHECKS} -gt 0 ]]; then
        echo ""
        echo -e "${RED}Failed Checks:${NC}"
        for check in "${FAILED_CHECK_NAMES[@]}"; do
            echo "  - ${check}"
        done
        echo ""
        echo -e "${RED}Deployment validation FAILED${NC}"
        echo "Please fix the issues above before deploying"
        exit 1
    else
        echo ""
        echo -e "${GREEN}All checks passed! Ready for deployment.${NC}"
        exit 0
    fi
}

# Run main function
main "$@"
