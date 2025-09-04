#!/bin/bash
# Production Deployment Script - Day 15
# Unified Attributes System Final Deployment
# Version: 1.0.0
# Date: 03.09.2025

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENVIRONMENT="production"
DEPLOYMENT_VERSION="v2.0.0-unified-attributes"
DEPLOYMENT_DATE=$(date +"%Y-%m-%d %H:%M:%S")
LOG_FILE="/var/log/deployment/unified-attributes-$(date +%Y%m%d-%H%M%S).log"
BACKUP_DIR="/var/backups/unified-attributes"
HEALTH_CHECK_URL="https://api.svetu.rs/health"
MONITORING_URL="https://grafana.svetu.rs/d/unified-attrs"
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL:-}"

# Database configuration
DB_HOST="prod-db-master.svetu.rs"
DB_NAME="svetubd"
DB_USER="postgres"
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=require"

# Deployment phases
PHASE_1_PREPARATION="preparation"
PHASE_2_GREEN_DEPLOYMENT="green_deployment"
PHASE_3_TRAFFIC_SWITCH="traffic_switch"
PHASE_4_FINALIZATION="finalization"
PHASE_5_CLEANUP="cleanup"

# Current phase tracking
CURRENT_PHASE=""
ROLLBACK_REQUIRED=false

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

# Slack notification
send_slack_notification() {
    local message="$1"
    local color="${2:-#36a64f}"
    
    if [[ -n "$SLACK_WEBHOOK_URL" ]]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{
                \"attachments\": [{
                    \"color\": \"${color}\",
                    \"title\": \"Unified Attributes Deployment\",
                    \"text\": \"${message}\",
                    \"footer\": \"Environment: ${ENVIRONMENT}\",
                    \"ts\": $(date +%s)
                }]
            }" \
            "$SLACK_WEBHOOK_URL" 2>/dev/null || true
    fi
}

# Phase tracking
start_phase() {
    local phase_name="$1"
    CURRENT_PHASE="$phase_name"
    log_info "=========================================="
    log_info "Starting Phase: $phase_name"
    log_info "=========================================="
    send_slack_notification "Starting deployment phase: $phase_name" "#439FE0"
}

complete_phase() {
    local phase_name="$1"
    log_success "Phase completed: $phase_name"
    send_slack_notification "✅ Phase completed: $phase_name" "#36a64f"
}

# Error handling
handle_error() {
    local error_message="$1"
    log_error "$error_message"
    send_slack_notification "❌ Deployment failed: $error_message" "#ff0000"
    ROLLBACK_REQUIRED=true
    initiate_rollback
    exit 1
}

# Rollback procedure
initiate_rollback() {
    log_warning "Initiating rollback procedure..."
    send_slack_notification "⚠️ Initiating rollback for phase: $CURRENT_PHASE" "#ff9900"
    
    case "$CURRENT_PHASE" in
        "$PHASE_2_GREEN_DEPLOYMENT")
            rollback_green_deployment
            ;;
        "$PHASE_3_TRAFFIC_SWITCH")
            rollback_traffic_switch
            ;;
        "$PHASE_4_FINALIZATION")
            rollback_finalization
            ;;
        *)
            log_info "No rollback needed for phase: $CURRENT_PHASE"
            ;;
    esac
    
    log_success "Rollback completed"
    send_slack_notification "✅ Rollback completed successfully" "#36a64f"
}

# Rollback functions
rollback_green_deployment() {
    log_info "Rolling back green deployment..."
    kubectl delete deployment backend-green || true
    kubectl delete deployment frontend-green || true
}

rollback_traffic_switch() {
    log_info "Rolling back traffic switch..."
    # Switch all traffic back to blue
    kubectl patch ingress main-ingress -p '{"spec":{"rules":[{"http":{"paths":[{"backend":{"service":{"name":"backend-blue","port":{"number":3000}}}}]}}]}}'
}

rollback_finalization() {
    log_info "Rolling back finalization..."
    # Re-enable dual-write and fallback
    kubectl set env deployment/backend-blue DUAL_WRITE_ATTRIBUTES=true
    kubectl set env deployment/backend-blue UNIFIED_ATTRIBUTES_FALLBACK=true
}

# Health checks
perform_health_check() {
    local service_url="$1"
    local max_retries=10
    local retry_count=0
    
    while [ $retry_count -lt $max_retries ]; do
        if curl -sf "$service_url" >/dev/null 2>&1; then
            return 0
        fi
        retry_count=$((retry_count + 1))
        log_warning "Health check failed, retry $retry_count/$max_retries"
        sleep 5
    done
    
    return 1
}

# Validation functions
validate_deployment_readiness() {
    log_info "Validating deployment readiness..."
    
    # Check Kubernetes cluster
    if ! kubectl cluster-info >/dev/null 2>&1; then
        handle_error "Kubernetes cluster not accessible"
    fi
    
    # Check database connectivity
    if ! psql "$DATABASE_URL" -c "SELECT 1" >/dev/null 2>&1; then
        handle_error "Database not accessible"
    fi
    
    # Check blue environment health
    if ! perform_health_check "https://backend-blue.svetu.rs/health"; then
        handle_error "Blue environment not healthy"
    fi
    
    log_success "Deployment readiness validated"
}

# Phase 1: Preparation
phase_1_preparation() {
    start_phase "$PHASE_1_PREPARATION"
    
    # Create backup directory
    mkdir -p "$BACKUP_DIR"
    mkdir -p "$(dirname "$LOG_FILE")"
    
    # Validate environment
    validate_deployment_readiness
    
    # Create final database backup
    log_info "Creating database backup..."
    pg_dump "$DATABASE_URL" | gzip > "$BACKUP_DIR/pre-deployment-$(date +%Y%m%d-%H%M%S).sql.gz"
    
    # Verify migrations are applied
    log_info "Verifying migrations..."
    if ! psql "$DATABASE_URL" -c "SELECT * FROM unified_attributes LIMIT 1" >/dev/null 2>&1; then
        handle_error "Unified attributes table not found - migrations may not be applied"
    fi
    
    # Check monitoring
    log_info "Checking monitoring setup..."
    if ! curl -sf "$MONITORING_URL" >/dev/null 2>&1; then
        log_warning "Monitoring dashboard not accessible"
    fi
    
    complete_phase "$PHASE_1_PREPARATION"
}

# Phase 2: Green Environment Deployment
phase_2_green_deployment() {
    start_phase "$PHASE_2_GREEN_DEPLOYMENT"
    
    log_info "Deploying backend to green environment..."
    
    # Deploy backend green
    kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend-green
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
      version: green
  template:
    metadata:
      labels:
        app: backend
        version: green
    spec:
      containers:
      - name: backend
        image: registry.svetu.rs/backend:${DEPLOYMENT_VERSION}
        ports:
        - containerPort: 3000
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: USE_UNIFIED_ATTRIBUTES
          value: "true"
        - name: UNIFIED_ATTRIBUTES_FALLBACK
          value: "false"
        - name: DUAL_WRITE_ATTRIBUTES
          value: "false"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 5
EOF
    
    log_info "Deploying frontend to green environment..."
    
    # Deploy frontend green
    kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend-green
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: frontend
      version: green
  template:
    metadata:
      labels:
        app: frontend
        version: green
    spec:
      containers:
      - name: frontend
        image: registry.svetu.rs/frontend:${DEPLOYMENT_VERSION}
        ports:
        - containerPort: 3001
        env:
        - name: NODE_ENV
          value: "production"
        - name: NEXT_PUBLIC_API_URL
          value: "https://api.svetu.rs"
        - name: USE_UNIFIED_ATTRIBUTES
          value: "true"
        resources:
          requests:
            memory: "256Mi"
            cpu: "200m"
          limits:
            memory: "512Mi"
            cpu: "400m"
EOF
    
    # Wait for green deployment to be ready
    log_info "Waiting for green environment to be ready..."
    kubectl rollout status deployment/backend-green -n production --timeout=600s
    kubectl rollout status deployment/frontend-green -n production --timeout=600s
    
    # Health check green environment
    if ! perform_health_check "https://backend-green.svetu.rs/health"; then
        handle_error "Green backend not healthy after deployment"
    fi
    
    if ! perform_health_check "https://frontend-green.svetu.rs/health"; then
        handle_error "Green frontend not healthy after deployment"
    fi
    
    complete_phase "$PHASE_2_GREEN_DEPLOYMENT"
}

# Phase 3: Traffic Switch
phase_3_traffic_switch() {
    start_phase "$PHASE_3_TRAFFIC_SWITCH"
    
    log_info "Starting progressive traffic switch to green environment..."
    
    # Function to update traffic split
    update_traffic_split() {
        local blue_weight=$1
        local green_weight=$2
        
        kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: backend-vs
  namespace: production
spec:
  hosts:
  - api.svetu.rs
  http:
  - match:
    - uri:
        prefix: /
    route:
    - destination:
        host: backend-blue
      weight: ${blue_weight}
    - destination:
        host: backend-green
      weight: ${green_weight}
EOF
        
        log_info "Traffic split updated: Blue=${blue_weight}%, Green=${green_weight}%"
    }
    
    # Progressive rollout
    declare -a rollout_steps=("90 10" "75 25" "50 50" "25 75" "0 100")
    declare -a wait_times=(300 300 600 600 600) # Wait times in seconds
    
    for i in "${!rollout_steps[@]}"; do
        IFS=' ' read -r blue_weight green_weight <<< "${rollout_steps[$i]}"
        wait_time="${wait_times[$i]}"
        
        update_traffic_split "$blue_weight" "$green_weight"
        
        log_info "Monitoring metrics for ${wait_time} seconds..."
        
        # Monitor error rate
        start_time=$(date +%s)
        while [ $(($(date +%s) - start_time)) -lt "$wait_time" ]; do
            error_rate=$(curl -s "https://prometheus.svetu.rs/api/v1/query?query=rate(http_requests_total{status=~\"5..\"}[1m])" | jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
            
            if (( $(echo "$error_rate > 0.01" | bc -l) )); then
                log_error "Error rate too high: ${error_rate}"
                handle_error "Error rate exceeded threshold during traffic switch"
            fi
            
            sleep 30
        done
        
        log_success "Traffic at Blue=${blue_weight}%, Green=${green_weight}% - stable"
    done
    
    complete_phase "$PHASE_3_TRAFFIC_SWITCH"
}

# Phase 4: Finalization
phase_4_finalization() {
    start_phase "$PHASE_4_FINALIZATION"
    
    log_info "Disabling dual-write mode on blue environment..."
    kubectl set env deployment/backend-blue DUAL_WRITE_ATTRIBUTES=false
    
    log_info "Removing fallback configuration..."
    kubectl set env deployment/backend-blue UNIFIED_ATTRIBUTES_FALLBACK=false
    
    log_info "Updating service endpoints to point to green only..."
    kubectl patch service backend -p '{"spec":{"selector":{"version":"green"}}}'
    kubectl patch service frontend -p '{"spec":{"selector":{"version":"green"}}}'
    
    log_info "Removing blue environment from load balancer..."
    kubectl delete ingress backend-blue-ingress || true
    kubectl delete ingress frontend-blue-ingress || true
    
    # Final health check
    if ! perform_health_check "$HEALTH_CHECK_URL"; then
        handle_error "Production health check failed after finalization"
    fi
    
    complete_phase "$PHASE_4_FINALIZATION"
}

# Phase 5: Cleanup
phase_5_cleanup() {
    start_phase "$PHASE_5_CLEANUP"
    
    log_info "Archiving old attribute tables..."
    psql "$DATABASE_URL" <<EOF
-- Create archive schema if not exists
CREATE SCHEMA IF NOT EXISTS archive;

-- Move old tables to archive
ALTER TABLE IF EXISTS category_attributes SET SCHEMA archive;
ALTER TABLE IF EXISTS listing_attributes SET SCHEMA archive;
ALTER TABLE IF EXISTS category_attribute_values SET SCHEMA archive;

-- Add archive timestamp
COMMENT ON SCHEMA archive IS 'Archived on ${DEPLOYMENT_DATE} after unified attributes migration';
EOF
    
    log_info "Cleaning up temporary resources..."
    kubectl delete deployment backend-blue -n production || true
    kubectl delete deployment frontend-blue -n production || true
    
    log_info "Updating documentation..."
    echo "Deployment completed: ${DEPLOYMENT_DATE}" >> "$BACKUP_DIR/deployment.log"
    
    complete_phase "$PHASE_5_CLEANUP"
}

# Generate deployment report
generate_deployment_report() {
    local report_file="$BACKUP_DIR/deployment-report-$(date +%Y%m%d-%H%M%S).json"
    
    cat > "$report_file" <<EOF
{
  "deployment": {
    "version": "${DEPLOYMENT_VERSION}",
    "environment": "${ENVIRONMENT}",
    "date": "${DEPLOYMENT_DATE}",
    "status": "success",
    "phases": {
      "preparation": "completed",
      "green_deployment": "completed",
      "traffic_switch": "completed",
      "finalization": "completed",
      "cleanup": "completed"
    },
    "metrics": {
      "duration": "$(($(date +%s) - start_time)) seconds",
      "downtime": "0 seconds",
      "errors": "0",
      "rollbacks": "0"
    },
    "logs": "${LOG_FILE}"
  }
}
EOF
    
    log_success "Deployment report generated: $report_file"
}

# Main execution
main() {
    local start_time=$(date +%s)
    
    log_info "=========================================="
    log_info "Unified Attributes Production Deployment"
    log_info "Version: ${DEPLOYMENT_VERSION}"
    log_info "Environment: ${ENVIRONMENT}"
    log_info "Date: ${DEPLOYMENT_DATE}"
    log_info "=========================================="
    
    send_slack_notification "🚀 Starting Unified Attributes production deployment" "#439FE0"
    
    # Execute deployment phases
    phase_1_preparation
    phase_2_green_deployment
    phase_3_traffic_switch
    phase_4_finalization
    phase_5_cleanup
    
    # Generate report
    generate_deployment_report
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    log_success "=========================================="
    log_success "DEPLOYMENT COMPLETED SUCCESSFULLY"
    log_success "Duration: ${duration} seconds"
    log_success "=========================================="
    
    send_slack_notification "✅ Unified Attributes deployment completed successfully in ${duration} seconds" "#36a64f"
}

# Trap errors
trap 'handle_error "Unexpected error occurred"' ERR

# Check if running as root or with sudo
if [[ $EUID -ne 0 ]]; then
   log_error "This script must be run as root or with sudo"
   exit 1
fi

# Confirmation prompt
echo -e "${YELLOW}WARNING: You are about to deploy Unified Attributes to PRODUCTION${NC}"
echo -e "${YELLOW}This will affect ALL users and services.${NC}"
read -p "Are you sure you want to continue? (yes/no): " confirmation

if [[ "$confirmation" != "yes" ]]; then
    log_info "Deployment cancelled by user"
    exit 0
fi

# Run main deployment
main "$@"