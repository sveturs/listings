#!/bin/bash

# AI Learning Cron Script
# Автоматически запускает процесс обучения AI системы категоризации

set -e

# Configuration
BASE_URL="http://localhost:3000"
API_URL="$BASE_URL/api/v1/marketplace/ai"
LOG_FILE="/var/log/ai_learning_cron.log"
LOCK_FILE="/tmp/ai_learning.lock"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Function to log with timestamp
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') [AI-LEARNING] $1" | tee -a "$LOG_FILE"
}

# Function to check if another instance is running
check_lock() {
    if [ -f "$LOCK_FILE" ]; then
        local pid=$(cat "$LOCK_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            log "Another AI learning process is already running (PID: $pid)"
            exit 1
        else
            log "Removing stale lock file"
            rm -f "$LOCK_FILE"
        fi
    fi
}

# Function to create lock file
create_lock() {
    echo $$ > "$LOCK_FILE"
}

# Function to remove lock file
remove_lock() {
    rm -f "$LOCK_FILE"
}

# Function to check if backend is running
check_backend() {
    local response=$(curl -s -f "$BASE_URL/health" 2>/dev/null || echo "")
    if [ -z "$response" ]; then
        log "❌ Backend is not running at $BASE_URL"
        return 1
    fi
    return 0
}

# Function to trigger learning from feedback
trigger_learning() {
    log "🧠 Triggering AI learning from feedback..."

    local response=$(curl -s -X POST "$API_URL/learn" \
        -H "Content-Type: application/json" \
        -w "%{http_code}" \
        -o /tmp/learning_response.json)

    if [ "$response" = "200" ]; then
        log "✅ Learning from feedback completed successfully"

        # Parse response for metrics
        if [ -f "/tmp/learning_response.json" ]; then
            local improvements=$(jq -r '.data.improvementsApplied // 0' /tmp/learning_response.json 2>/dev/null || echo "0")
            local keywords=$(jq -r '.data.keywordsLearned // 0' /tmp/learning_response.json 2>/dev/null || echo "0")

            if [ "$improvements" != "0" ] || [ "$keywords" != "0" ]; then
                log "📊 Metrics: $improvements improvements applied, $keywords keywords learned"
            fi
        fi

        rm -f /tmp/learning_response.json
        return 0
    else
        log "❌ Learning trigger failed with HTTP code: $response"
        if [ -f "/tmp/learning_response.json" ]; then
            log "Response: $(cat /tmp/learning_response.json)"
            rm -f /tmp/learning_response.json
        fi
        return 1
    fi
}

# Function to trigger bulk keyword generation for categories that need it
trigger_bulk_keywords() {
    log "🔤 Checking categories needing keyword expansion..."

    local response=$(curl -s -X POST "$API_URL/generate-keywords-all?minKeywords=40" \
        -H "Content-Type: application/json" \
        -w "%{http_code}" \
        -o /tmp/keywords_response.json)

    if [ "$response" = "200" ]; then
        if [ -f "/tmp/keywords_response.json" ]; then
            local categories_found=$(jq -r '.data.categoriesFound // 0' /tmp/keywords_response.json 2>/dev/null || echo "0")
            local message=$(jq -r '.data.message // "Unknown"' /tmp/keywords_response.json 2>/dev/null || echo "Unknown")

            log "✅ Bulk keyword generation: $message"
            if [ "$categories_found" != "0" ]; then
                log "📊 Processing $categories_found categories in background"
            fi
        fi
        rm -f /tmp/keywords_response.json
        return 0
    else
        log "❌ Bulk keyword generation failed with HTTP code: $response"
        if [ -f "/tmp/keywords_response.json" ]; then
            log "Response: $(cat /tmp/keywords_response.json)"
            rm -f /tmp/keywords_response.json
        fi
        return 1
    fi
}

# Function to get learning metrics
get_learning_metrics() {
    log "📊 Retrieving learning metrics..."

    local response=$(curl -s "$API_URL/metrics?days=1" \
        -w "%{http_code}" \
        -o /tmp/metrics_response.json)

    if [ "$response" = "200" ]; then
        if [ -f "/tmp/metrics_response.json" ]; then
            local accuracy=$(jq -r '.data.accuracy // "unknown"' /tmp/metrics_response.json 2>/dev/null || echo "unknown")
            local total_detections=$(jq -r '.data.totalDetections // 0' /tmp/metrics_response.json 2>/dev/null || echo "0")

            if [ "$accuracy" != "unknown" ] && [ "$total_detections" != "0" ]; then
                log "📈 Current accuracy: $accuracy% (based on $total_detections detections)"

                # Alert if accuracy is below threshold
                if [ "$accuracy" != "unknown" ]; then
                    local accuracy_int=$(echo "$accuracy" | cut -d'.' -f1)
                    if [ "$accuracy_int" -lt 95 ]; then
                        log "⚠️  WARNING: Accuracy below 95% - consider additional training"
                    fi
                fi
            fi
        fi
        rm -f /tmp/metrics_response.json
    else
        log "❌ Failed to retrieve metrics with HTTP code: $response"
    fi
}

# Function to cleanup old logs
cleanup_logs() {
    # Keep only last 30 days of logs
    if [ -f "$LOG_FILE" ]; then
        local temp_log="/tmp/ai_learning_temp.log"
        tail -10000 "$LOG_FILE" > "$temp_log" 2>/dev/null || true
        mv "$temp_log" "$LOG_FILE" 2>/dev/null || true
    fi
}

# Main execution
main() {
    log "🚀 Starting AI learning cron job"

    # Check if another instance is running
    check_lock
    create_lock

    # Set trap to cleanup on exit
    trap remove_lock EXIT

    # Check if backend is running
    if ! check_backend; then
        log "❌ Backend health check failed - skipping learning session"
        exit 1
    fi

    local success_count=0
    local total_tasks=3

    # 1. Trigger learning from feedback
    if trigger_learning; then
        success_count=$((success_count + 1))
    fi

    # Small delay between operations
    sleep 2

    # 2. Trigger bulk keyword generation for categories that need it
    if trigger_bulk_keywords; then
        success_count=$((success_count + 1))
    fi

    # Small delay
    sleep 2

    # 3. Get and log current metrics
    get_learning_metrics
    success_count=$((success_count + 1))  # Always count metrics as success

    # Cleanup old logs
    cleanup_logs

    # Final status
    if [ $success_count -eq $total_tasks ]; then
        log "✅ AI learning cron job completed successfully ($success_count/$total_tasks tasks)"
    else
        log "⚠️  AI learning cron job completed with issues ($success_count/$total_tasks tasks successful)"
    fi

    log "🏁 Learning session finished"
}

# Execution based on parameters
case "${1:-run}" in
    "run")
        main
        ;;
    "test")
        echo "🧪 Testing AI learning system connectivity..."
        if check_backend; then
            echo "✅ Backend is accessible"
            echo "🔗 API URL: $API_URL"
            echo "📝 Log file: $LOG_FILE"
            echo "🔒 Lock file: $LOCK_FILE"
        else
            echo "❌ Backend is not accessible"
            exit 1
        fi
        ;;
    "install")
        echo "📅 Installing AI learning cron job..."

        # Add to crontab (every 6 hours)
        local cron_entry="0 */6 * * * $0 run >/dev/null 2>&1"

        # Check if already exists
        if crontab -l 2>/dev/null | grep -q "$0"; then
            echo "⚠️  Cron job already exists"
        else
            (crontab -l 2>/dev/null; echo "$cron_entry") | crontab -
            echo "✅ Cron job installed: every 6 hours"
        fi

        # Create log directory if needed
        mkdir -p "$(dirname "$LOG_FILE")"
        touch "$LOG_FILE"

        echo "📝 Log file: $LOG_FILE"
        echo "🧪 Test with: $0 test"
        echo "🔄 Manual run: $0 run"
        ;;
    "uninstall")
        echo "🗑️  Removing AI learning cron job..."
        crontab -l 2>/dev/null | grep -v "$0" | crontab -
        echo "✅ Cron job removed"
        ;;
    "status")
        echo "📊 AI Learning System Status"
        echo "=========================="

        if check_backend; then
            echo "✅ Backend: Running"
        else
            echo "❌ Backend: Not accessible"
        fi

        if [ -f "$LOCK_FILE" ]; then
            local pid=$(cat "$LOCK_FILE")
            if kill -0 "$pid" 2>/dev/null; then
                echo "🔄 Learning: Running (PID: $pid)"
            else
                echo "⚠️  Learning: Stale lock file"
            fi
        else
            echo "⏹️  Learning: Not running"
        fi

        if [ -f "$LOG_FILE" ]; then
            local log_size=$(du -h "$LOG_FILE" | cut -f1)
            local last_run=$(tail -1 "$LOG_FILE" 2>/dev/null | grep -o '^[0-9-]* [0-9:]*' || echo "Never")
            echo "📝 Log file: $log_size, last run: $last_run"
        else
            echo "📝 Log file: Not found"
        fi

        # Check crontab
        if crontab -l 2>/dev/null | grep -q "$0"; then
            echo "📅 Cron job: Installed"
        else
            echo "📅 Cron job: Not installed"
        fi
        ;;
    "help"|*)
        echo "AI Learning System Cron Manager"
        echo ""
        echo "Commands:"
        echo "  run        - Execute learning session (default)"
        echo "  test       - Test connectivity and configuration"
        echo "  install    - Install cron job (every 6 hours)"
        echo "  uninstall  - Remove cron job"
        echo "  status     - Show system status"
        echo "  help       - Show this help"
        echo ""
        echo "Examples:"
        echo "  $0 run       # Manual execution"
        echo "  $0 install   # Setup automatic execution"
        echo "  $0 status    # Check current status"
        ;;
esac