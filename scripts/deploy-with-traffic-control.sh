#!/bin/bash
set -e

# Traffic-controlled deployment script with monitoring and auto-rollback
# Usage: ./deploy-with-traffic-control.sh [TRAFFIC_PERCENT] [ENVIRONMENT]
#
# Example:
#   ./deploy-with-traffic-control.sh 10 staging   # Deploy with 10% traffic to staging
#   ./deploy-with-traffic-control.sh 50 production # Deploy with 50% traffic to production

TRAFFIC_PERCENT=${1:-0}
ENVIRONMENT=${2:-staging}
MONITORING_DURATION=${3:-300}  # Default: 5 minutes

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration based on environment
case "$ENVIRONMENT" in
    staging)
        SERVER_HOST="staging.svetu.rs"
        PROMETHEUS_URL="http://prometheus.svetu.rs"
        PROJECT_DIR="/opt/svetu-staging"
        ;;
    production)
        SERVER_HOST="svetu.rs"
        PROMETHEUS_URL="http://prometheus.svetu.rs"
        PROJECT_DIR="/opt/svetu-prod"
        ;;
    *)
        echo -e "${RED}❌ Invalid environment: $ENVIRONMENT${NC}"
        echo "Valid options: staging, production"
        exit 1
        ;;
esac

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🚀 Traffic-Controlled Deployment${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Environment:       $ENVIRONMENT"
echo "Server:            $SERVER_HOST"
echo "Traffic Percent:   $TRAFFIC_PERCENT%"
echo "Monitoring:        ${MONITORING_DURATION}s"
echo ""

# Validate traffic percentage
if ! [[ "$TRAFFIC_PERCENT" =~ ^[0-9]+$ ]] || [ "$TRAFFIC_PERCENT" -lt 0 ] || [ "$TRAFFIC_PERCENT" -gt 100 ]; then
    echo -e "${RED}❌ Invalid traffic percentage: $TRAFFIC_PERCENT${NC}"
    echo "Must be a number between 0 and 100"
    exit 1
fi

# Function to check service health
check_health() {
    local url="https://$SERVER_HOST/health"
    local max_attempts=30
    local attempt=0

    echo ""
    echo -e "${YELLOW}🏥 Performing health check...${NC}"

    while [ $attempt -lt $max_attempts ]; do
        attempt=$((attempt + 1))

        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$url" || echo "000")

        if [ "$HTTP_CODE" == "200" ]; then
            echo -e "${GREEN}✅ Health check passed (attempt $attempt/$max_attempts)${NC}"
            return 0
        fi

        echo "Attempt $attempt/$max_attempts: HTTP $HTTP_CODE - retrying in 2s..."
        sleep 2
    done

    echo -e "${RED}❌ Health check failed after $max_attempts attempts${NC}"
    return 1
}

# Function to query Prometheus metrics
query_prometheus() {
    local query=$1
    local result=$(curl -s "$PROMETHEUS_URL/api/v1/query?query=$query" | jq -r '.data.result[0].value[1] // "0"')
    echo "$result"
}

# Function to monitor metrics
monitor_metrics() {
    local duration=$1
    local checks=$((duration / 5))  # Check every 5 seconds
    local check=0

    echo ""
    echo -e "${YELLOW}📊 Monitoring deployment metrics for ${duration}s...${NC}"
    echo ""

    while [ $check -lt $checks ]; do
        check=$((check + 1))
        local elapsed=$((check * 5))

        echo "[$elapsed/${duration}s] Collecting metrics..."

        # Query metrics from Prometheus
        local error_rate=$(query_prometheus "rate(traffic_router_errors_total[1m])")
        local response_time_p99=$(query_prometheus "histogram_quantile(0.99, rate(traffic_router_request_duration_seconds_bucket[1m]))")
        local microservice_up=$(query_prometheus "up{job=\"marketplace-microservice\"}")
        local monolith_up=$(query_prometheus "up{job=\"marketplace-monolith\"}")

        echo "  Error rate:         $error_rate"
        echo "  Response time P99:  ${response_time_p99}s"
        echo "  Microservice up:    $microservice_up"
        echo "  Monolith up:        $monolith_up"

        # Check error rate threshold (> 1%)
        if (( $(echo "$error_rate > 0.01" | bc -l) )); then
            echo -e "${RED}🔴 ALERT: Error rate too high: $error_rate (threshold: 0.01)${NC}"
            return 1
        fi

        # Check response time threshold (> 2s)
        if (( $(echo "$response_time_p99 > 2.0" | bc -l) )); then
            echo -e "${RED}🔴 ALERT: Response time too slow: ${response_time_p99}s (threshold: 2.0s)${NC}"
            return 1
        fi

        # Check if microservice is down (when traffic > 0)
        if [ "$TRAFFIC_PERCENT" -gt 0 ] && (( $(echo "$microservice_up < 1" | bc -l) )); then
            echo -e "${RED}🔴 ALERT: Microservice is down${NC}"
            return 1
        fi

        # Check if monolith is down
        if (( $(echo "$monolith_up < 1" | bc -l) )); then
            echo -e "${RED}🔴 ALERT: Monolith is down${NC}"
            return 1
        fi

        echo ""
        sleep 5
    done

    echo -e "${GREEN}✅ Monitoring completed successfully${NC}"
    return 0
}

# Function to rollback deployment
rollback() {
    echo ""
    echo -e "${RED}🔄 Initiating rollback...${NC}"

    ssh svetu@$SERVER_HOST "cd $PROJECT_DIR && \
        export TRAFFIC_ROUTER_PERCENT=0 && \
        docker-compose restart backend"

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Rollback completed - all traffic routed to monolith${NC}"
    else
        echo -e "${RED}❌ Rollback failed! Manual intervention required${NC}"
        exit 1
    fi
}

# Main deployment flow
echo -e "${YELLOW}📦 Starting deployment...${NC}"

# Step 1: Deploy with traffic control
echo ""
echo -e "${YELLOW}1️⃣  Deploying backend with traffic router...${NC}"

ssh svetu@$SERVER_HOST "cd $PROJECT_DIR && \
    export TRAFFIC_ROUTER_ENABLED=true && \
    export TRAFFIC_ROUTER_PERCENT=$TRAFFIC_PERCENT && \
    export MICROSERVICE_TIMEOUT=500ms && \
    export CIRCUIT_BREAKER_ENABLED=true && \
    docker-compose pull backend && \
    docker-compose up -d backend"

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Deployment failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Backend deployed successfully${NC}"

# Step 2: Wait for service to be ready
echo ""
echo -e "${YELLOW}2️⃣  Waiting for service to be ready...${NC}"
sleep 10

# Step 3: Health check
if ! check_health; then
    echo -e "${RED}❌ Health check failed${NC}"
    rollback
    exit 1
fi

# Step 4: Run smoke tests
echo ""
echo -e "${YELLOW}3️⃣  Running smoke tests...${NC}"

ssh svetu@$SERVER_HOST "cd $PROJECT_DIR/backend/tests/e2e && \
    BACKEND_URL=https://$SERVER_HOST go test -v -tags=smoke -timeout=5m" || {
    echo -e "${RED}❌ Smoke tests failed${NC}"
    rollback
    exit 1
}

echo -e "${GREEN}✅ Smoke tests passed${NC}"

# Step 5: Monitor metrics (only if traffic > 0)
if [ "$TRAFFIC_PERCENT" -gt 0 ]; then
    if ! monitor_metrics "$MONITORING_DURATION"; then
        echo -e "${RED}❌ Monitoring detected issues${NC}"
        rollback
        exit 1
    fi
else
    echo ""
    echo -e "${YELLOW}⏭️  Skipping monitoring (traffic is 0%)${NC}"
fi

# Step 6: Verify deployment version
echo ""
echo -e "${YELLOW}4️⃣  Verifying deployment version...${NC}"

DEPLOYED_VERSION=$(curl -s "https://$SERVER_HOST/" | grep -oP 'Svetu API \K[0-9.]+' || echo "unknown")
echo "Deployed version: $DEPLOYED_VERSION"

# Success!
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Environment:       $ENVIRONMENT"
echo "Server:            https://$SERVER_HOST"
echo "Version:           $DEPLOYED_VERSION"
echo "Traffic:           $TRAFFIC_PERCENT% to microservice"
echo "Status:            🟢 Healthy"
echo ""

# Next steps suggestion
if [ "$TRAFFIC_PERCENT" -lt 100 ]; then
    NEXT_PERCENT=$((TRAFFIC_PERCENT + 10))
    if [ "$NEXT_PERCENT" -gt 100 ]; then
        NEXT_PERCENT=100
    fi

    echo -e "${YELLOW}💡 Next steps:${NC}"
    echo "   Monitor metrics in Grafana: http://grafana.svetu.rs"
    echo "   Increase traffic: ./deploy-with-traffic-control.sh $NEXT_PERCENT $ENVIRONMENT"
fi

echo ""
exit 0
