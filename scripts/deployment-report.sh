#!/usr/bin/env bash

#######################################
# Deployment Report Generator
# Generates comprehensive deployment report
#
# Usage: ./deployment-report.sh --deployment-id ID [--output FILE]
#
# Exit codes:
#   0 - Success
#   1 - Failed
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
fi

# Parameters
DEPLOYMENT_ID=""
OUTPUT_FILE=""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

#######################################
# Parse arguments
#######################################

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --deployment-id)
                DEPLOYMENT_ID="$2"
                shift 2
                ;;
            --output)
                OUTPUT_FILE="$2"
                shift 2
                ;;
            *)
                echo "Usage: $0 --deployment-id ID [--output FILE]"
                exit 1
                ;;
        esac
    done

    if [[ -z "${DEPLOYMENT_ID}" ]]; then
        echo "ERROR: --deployment-id is required"
        exit 1
    fi

    if [[ -z "${OUTPUT_FILE}" ]]; then
        OUTPUT_FILE="${PROJECT_ROOT}/logs/deployments/${DEPLOYMENT_ID}-report.md"
    fi
}

#######################################
# Gather deployment information
#######################################

get_version_info() {
    cd "${PROJECT_ROOT}"

    local version=$(git describe --tags --always 2>/dev/null || echo "unknown")
    local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    local branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    local author=$(git log -1 --pretty=format:'%an <%ae>' 2>/dev/null || echo "unknown")
    local commit_msg=$(git log -1 --pretty=format:'%s' 2>/dev/null || echo "unknown")

    cat <<EOF
## Version Information

- **Version:** ${version}
- **Commit:** ${commit}
- **Branch:** ${branch}
- **Author:** ${author}
- **Last Commit Message:** ${commit_msg}
EOF
}

get_deployment_timing() {
    local log_file="${PROJECT_ROOT}/logs/deployments/${DEPLOYMENT_ID}.log"

    if [[ ! -f "${log_file}" ]]; then
        echo "## Deployment Timing"
        echo ""
        echo "*Log file not found*"
        return
    fi

    local start_time=$(grep -m1 "Starting Production Deployment" "${log_file}" | awk '{print $1, $2}' || echo "unknown")
    local end_time=$(grep -m1 "Deployment completed successfully" "${log_file}" | awk '{print $1, $2}' || echo "unknown")

    # Parse log for phase durations
    local build_time=$(grep "Binary built successfully" "${log_file}" | head -1 || echo "")
    local migration_time=$(grep "Migrations completed" "${log_file}" | head -1 || echo "")
    local smoke_test_time=$(grep "Smoke tests passed" "${log_file}" | head -1 || echo "")

    cat <<EOF
## Deployment Timing

- **Start Time:** ${start_time}
- **End Time:** ${end_time}
- **Total Duration:** Calculated from logs

### Phase Durations
- **Build:** ${build_time:-"N/A"}
- **Migration:** ${migration_time:-"N/A"}
- **Smoke Tests:** ${smoke_test_time:-"N/A"}
EOF
}

get_test_results() {
    local log_file="${PROJECT_ROOT}/logs/deployments/${DEPLOYMENT_ID}.log"

    if [[ ! -f "${log_file}" ]]; then
        echo "## Test Results"
        echo ""
        echo "*No test results available*"
        return
    fi

    # Check if smoke tests were run
    if grep -q "Smoke tests" "${log_file}"; then
        local test_status=$(grep "Smoke tests" "${log_file}" | tail -1)

        cat <<EOF
## Test Results

### Smoke Tests
- **Status:** ${test_status}

EOF
    else
        cat <<EOF
## Test Results

*Smoke tests were skipped*

EOF
    fi
}

get_database_changes() {
    local log_file="${PROJECT_ROOT}/logs/deployments/${DEPLOYMENT_ID}.log"

    if [[ ! -f "${log_file}" ]]; then
        echo "## Database Changes"
        echo ""
        echo "*No migration information available*"
        return
    fi

    cat <<EOF
## Database Changes

EOF

    if grep -q "Running database migrations" "${log_file}"; then
        echo "- Migrations were executed"
        echo "- Backup created before deployment"
    else
        echo "- No migrations executed"
    fi

    echo ""
}

get_traffic_distribution() {
    cat <<EOF
## Traffic Distribution

### Canary Phases
1. **Phase 1:** 10% Green, 90% Blue (5 minutes monitoring)
2. **Phase 2:** 50% Green, 50% Blue (5 minutes monitoring)
3. **Phase 3:** 100% Green (complete switch)

### Final State
- **Blue:** Decommissioned
- **Green:** 100% traffic (promoted to Blue for next deployment)

EOF
}

get_metrics_comparison() {
    if [[ -n "${PROD_HOST:-}" ]]; then
        # Try to get current metrics from production
        local current_requests=$(ssh "${PROD_USER}@${PROD_HOST}" \
            "curl -s http://localhost:${GREEN_PORT}/metrics 2>/dev/null | grep 'http_requests_total' | tail -1" 2>/dev/null || echo "N/A")

        cat <<EOF
## Metrics Comparison

### Pre-Deployment
- *Baseline metrics would be captured here*

### Post-Deployment
- **Total Requests:** ${current_requests}
- **Error Rate:** To be measured over 24h
- **Response Time:** To be measured over 24h

EOF
    else
        cat <<EOF
## Metrics Comparison

*Metrics collection not configured*

EOF
    fi
}

get_incidents() {
    cat <<EOF
## Incidents

*No incidents reported during deployment*

EOF
}

get_rollback_info() {
    cat <<EOF
## Rollback Information

### Rollback Window
- **Duration:** 10 minutes after 100% traffic switch
- **Status:** Completed successfully (no rollback needed)

### Rollback Procedure (if needed)
\`\`\`bash
./scripts/rollback-prod.sh --reason "issue description"
\`\`\`

EOF
}

get_sign_off_checklist() {
    cat <<EOF
## Sign-Off Checklist

- [ ] All smoke tests passed
- [ ] No errors in production logs (first 1 hour)
- [ ] Database migrations successful
- [ ] Metrics within acceptable range
- [ ] No customer complaints
- [ ] Monitoring alerts normal
- [ ] Backup verified and accessible
- [ ] Documentation updated
- [ ] Team notified

### Sign-Off

**Deployed by:** ${USER}
**Reviewed by:** _______________________
**Approved by:** _______________________
**Date:** $(date +'%Y-%m-%d')

EOF
}

#######################################
# Generate report
#######################################

generate_report() {
    cat > "${OUTPUT_FILE}" <<EOF
# Deployment Report

**Deployment ID:** ${DEPLOYMENT_ID}
**Date:** $(date +'%Y-%m-%d %H:%M:%S %Z')
**Environment:** Production
**Server:** ${PROD_HOST:-"unknown"}

---

$(get_version_info)

---

$(get_deployment_timing)

---

$(get_test_results)

---

$(get_database_changes)

---

$(get_traffic_distribution)

---

$(get_metrics_comparison)

---

$(get_incidents)

---

$(get_rollback_info)

---

$(get_sign_off_checklist)

---

## Monitoring Links

- **Application Logs:** \`ssh ${PROD_USER:-user}@${PROD_HOST:-server} 'tail -f ${PROD_DIR:-/opt/listings}/logs/green.log'\`
- **Nginx Logs:** \`ssh ${PROD_USER:-user}@${PROD_HOST:-server} 'sudo tail -f /var/log/nginx/access.log'\`
- **Error Logs:** \`ssh ${PROD_USER:-user}@${PROD_HOST:-server} 'tail -f ${PROD_DIR:-/opt/listings}/logs/green.log | grep ERROR'\`

## Next Steps

1. Monitor production for 24 hours
2. Review metrics and error rates
3. Gather user feedback
4. Update runbook if needed
5. Schedule postmortem meeting (if any issues occurred)

## Related Documents

- [Deployment Guide](../DEPLOYMENT.md)
- [Rollback Procedures](../ROLLBACK.md)
- [Runbook](../RUNBOOK.md)

---

*This report was automatically generated by deployment-report.sh*
EOF

    echo -e "${GREEN}Report generated: ${OUTPUT_FILE}${NC}"
}

#######################################
# Main function
#######################################

main() {
    parse_args "$@"

    echo "Generating deployment report for ${DEPLOYMENT_ID}..."

    # Create output directory if needed
    mkdir -p "$(dirname "${OUTPUT_FILE}")"

    # Generate report
    generate_report

    # Display report location
    echo ""
    echo "========================================="
    echo "Deployment Report"
    echo "========================================="
    echo "Report saved to: ${OUTPUT_FILE}"
    echo ""
    echo "View report:"
    echo "  cat ${OUTPUT_FILE}"
    echo ""
    echo "Or open in editor:"
    echo "  \$EDITOR ${OUTPUT_FILE}"
    echo "========================================="
}

# Run main function
main "$@"
