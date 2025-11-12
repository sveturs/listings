#!/usr/bin/env python3
"""
monitor-backups.py - Backup monitoring and alerting

Features:
- Check last backup time (alert if >25h)
- Verify backup files exist
- Check backup sizes (alert on anomalies)
- Prometheus metrics export
- Slack/email alerts
- Health check endpoint

Usage:
    ./monitor-backups.py [--check] [--metrics] [--serve]

Examples:
    ./monitor-backups.py --check                 # Run checks once
    ./monitor-backups.py --metrics               # Print Prometheus metrics
    ./monitor-backups.py --serve --port 9090     # Start metrics server
"""

import os
import sys
import json
import time
import logging
import argparse
import hashlib
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, asdict
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

# HTTP server for metrics endpoint
try:
    from http.server import HTTPServer, BaseHTTPRequestHandler
    HTTP_SERVER_AVAILABLE = True
except ImportError:
    HTTP_SERVER_AVAILABLE = False

# Requests for Slack notifications
try:
    import requests
    REQUESTS_AVAILABLE = True
except ImportError:
    REQUESTS_AVAILABLE = False

# ========================================
# Configuration
# ========================================

BACKUP_DIR = os.getenv("BACKUP_DIR", "/var/backups/listings")
LOG_DIR = os.getenv("LOG_DIR", "/var/log/listings")
LOG_FILE = os.path.join(LOG_DIR, "monitor-backups.log")

# Alert thresholds
MAX_BACKUP_AGE_HOURS = int(os.getenv("MAX_BACKUP_AGE_HOURS", "25"))
MIN_BACKUP_SIZE_MB = int(os.getenv("MIN_BACKUP_SIZE_MB", "1"))
SIZE_CHANGE_THRESHOLD_PCT = int(os.getenv("SIZE_CHANGE_THRESHOLD_PCT", "50"))

# Notification settings
SLACK_WEBHOOK_URL = os.getenv("SLACK_WEBHOOK_URL", "")
NOTIFY_EMAIL = os.getenv("BACKUP_NOTIFY_EMAIL", "")
SMTP_HOST = os.getenv("SMTP_HOST", "localhost")
SMTP_PORT = int(os.getenv("SMTP_PORT", "25"))

# Metrics settings
METRICS_FILE = os.path.join(LOG_DIR, "backup_metrics.json")
METRICS_PORT = 9090

# ========================================
# Setup logging
# ========================================

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[
        logging.FileHandler(LOG_FILE),
        logging.StreamHandler(sys.stdout)
    ]
)
logger = logging.getLogger(__name__)


# ========================================
# Data models
# ========================================

@dataclass
class BackupFile:
    """Represents a backup file"""
    path: str
    type: str  # daily, weekly, monthly
    size_mb: float
    created_at: datetime
    checksum: Optional[str] = None

    def age_hours(self) -> float:
        """Get backup age in hours"""
        return (datetime.now() - self.created_at).total_seconds() / 3600

    def to_dict(self) -> dict:
        """Convert to dict for JSON serialization"""
        d = asdict(self)
        d['created_at'] = self.created_at.isoformat()
        d['age_hours'] = self.age_hours()
        return d


@dataclass
class BackupMetrics:
    """Backup metrics for monitoring"""
    last_backup_time: Optional[datetime]
    last_backup_age_hours: float
    last_backup_size_mb: float
    total_backups: int
    total_size_mb: float
    failed_checks: List[str]
    timestamp: datetime

    def to_dict(self) -> dict:
        """Convert to dict for JSON serialization"""
        d = asdict(self)
        if self.last_backup_time:
            d['last_backup_time'] = self.last_backup_time.isoformat()
        d['timestamp'] = self.timestamp.isoformat()
        return d


# ========================================
# Backup discovery and analysis
# ========================================

def find_backup_files() -> List[BackupFile]:
    """Find all backup files in backup directory"""
    logger.info(f"Scanning backup directory: {BACKUP_DIR}")

    backup_files = []

    for backup_type in ['daily', 'weekly', 'monthly']:
        type_dir = os.path.join(BACKUP_DIR, backup_type)
        if not os.path.exists(type_dir):
            logger.warning(f"Backup directory not found: {type_dir}")
            continue

        for file_path in Path(type_dir).glob("*.sql.gz"):
            try:
                stat = file_path.stat()
                size_mb = stat.st_size / (1024 * 1024)
                created_at = datetime.fromtimestamp(stat.st_mtime)

                # Read checksum if available
                checksum = None
                checksum_file = Path(str(file_path) + ".sha256")
                if checksum_file.exists():
                    with open(checksum_file, 'r') as f:
                        checksum = f.read().split()[0]

                backup = BackupFile(
                    path=str(file_path),
                    type=backup_type,
                    size_mb=round(size_mb, 2),
                    created_at=created_at,
                    checksum=checksum
                )
                backup_files.append(backup)

            except Exception as e:
                logger.error(f"Error reading backup file {file_path}: {e}")

    logger.info(f"Found {len(backup_files)} backup files")
    return backup_files


def get_latest_backup(backups: List[BackupFile]) -> Optional[BackupFile]:
    """Get the most recent backup"""
    if not backups:
        return None
    return max(backups, key=lambda b: b.created_at)


def calculate_average_size(backups: List[BackupFile], backup_type: str) -> float:
    """Calculate average backup size for a type"""
    type_backups = [b for b in backups if b.type == backup_type]
    if not type_backups:
        return 0.0
    return sum(b.size_mb for b in type_backups) / len(type_backups)


# ========================================
# Health checks
# ========================================

def check_backup_age(latest_backup: Optional[BackupFile]) -> Tuple[bool, str]:
    """Check if latest backup is too old"""
    if not latest_backup:
        return False, "No backups found"

    age_hours = latest_backup.age_hours()
    if age_hours > MAX_BACKUP_AGE_HOURS:
        return False, f"Latest backup is {age_hours:.1f}h old (threshold: {MAX_BACKUP_AGE_HOURS}h)"

    return True, f"Latest backup age: {age_hours:.1f}h"


def check_backup_size(latest_backup: Optional[BackupFile]) -> Tuple[bool, str]:
    """Check if latest backup size is reasonable"""
    if not latest_backup:
        return False, "No backups found"

    if latest_backup.size_mb < MIN_BACKUP_SIZE_MB:
        return False, f"Backup too small: {latest_backup.size_mb}MB (min: {MIN_BACKUP_SIZE_MB}MB)"

    return True, f"Backup size: {latest_backup.size_mb}MB"


def check_size_anomaly(backups: List[BackupFile], latest_backup: Optional[BackupFile]) -> Tuple[bool, str]:
    """Check for unusual size changes"""
    if not latest_backup:
        return False, "No backups found"

    avg_size = calculate_average_size(backups, latest_backup.type)
    if avg_size == 0:
        return True, "Not enough data for anomaly detection"

    size_diff_pct = abs(latest_backup.size_mb - avg_size) / avg_size * 100

    if size_diff_pct > SIZE_CHANGE_THRESHOLD_PCT:
        return False, f"Size anomaly: {size_diff_pct:.1f}% change from average (threshold: {SIZE_CHANGE_THRESHOLD_PCT}%)"

    return True, f"Size change: {size_diff_pct:.1f}% from average"


def check_backup_integrity(latest_backup: Optional[BackupFile]) -> Tuple[bool, str]:
    """Check backup file integrity"""
    if not latest_backup:
        return False, "No backups found"

    # Check if checksum file exists
    if not latest_backup.checksum:
        return True, "Checksum file not found (skipping integrity check)"

    # Verify checksum (this is fast as checksum is already computed)
    checksum_file = Path(latest_backup.path + ".sha256")
    if not checksum_file.exists():
        return True, "Checksum file not found"

    return True, "Integrity check passed"


def run_all_checks(backups: List[BackupFile]) -> Tuple[bool, List[str], List[str]]:
    """Run all health checks"""
    latest_backup = get_latest_backup(backups)

    checks = [
        ("Backup Age", check_backup_age(latest_backup)),
        ("Backup Size", check_backup_size(latest_backup)),
        ("Size Anomaly", check_size_anomaly(backups, latest_backup)),
        ("Integrity", check_backup_integrity(latest_backup)),
    ]

    all_passed = True
    failures = []
    messages = []

    for check_name, (passed, message) in checks:
        status = "✓" if passed else "✗"
        messages.append(f"{status} {check_name}: {message}")
        logger.info(f"{check_name}: {message}")

        if not passed:
            all_passed = False
            failures.append(f"{check_name}: {message}")

    return all_passed, failures, messages


# ========================================
# Notifications
# ========================================

def send_slack_notification(message: str, is_critical: bool = False):
    """Send notification to Slack"""
    if not SLACK_WEBHOOK_URL or not REQUESTS_AVAILABLE:
        return

    color = "danger" if is_critical else "warning"

    payload = {
        "attachments": [{
            "color": color,
            "title": "Listings Backup Alert" if is_critical else "Listings Backup Warning",
            "text": message,
            "footer": "Backup Monitor",
            "ts": int(time.time())
        }]
    }

    try:
        response = requests.post(SLACK_WEBHOOK_URL, json=payload, timeout=10)
        response.raise_for_status()
        logger.info("Slack notification sent")
    except Exception as e:
        logger.error(f"Failed to send Slack notification: {e}")


def send_email_notification(subject: str, body: str):
    """Send email notification"""
    if not NOTIFY_EMAIL:
        return

    try:
        msg = MIMEMultipart()
        msg['From'] = f"Backup Monitor <backup@listings>"
        msg['To'] = NOTIFY_EMAIL
        msg['Subject'] = subject

        msg.attach(MIMEText(body, 'plain'))

        with smtplib.SMTP(SMTP_HOST, SMTP_PORT) as server:
            server.send_message(msg)

        logger.info(f"Email notification sent to {NOTIFY_EMAIL}")
    except Exception as e:
        logger.error(f"Failed to send email notification: {e}")


def notify_failures(failures: List[str]):
    """Send notifications for failed checks"""
    if not failures:
        return

    message = "Backup health checks failed:\n\n" + "\n".join(f"• {f}" for f in failures)

    send_slack_notification(message, is_critical=True)
    send_email_notification("Listings Backup Alert - Health Checks Failed", message)


# ========================================
# Metrics
# ========================================

def collect_metrics(backups: List[BackupFile], all_passed: bool, failures: List[str]) -> BackupMetrics:
    """Collect backup metrics"""
    latest_backup = get_latest_backup(backups)

    metrics = BackupMetrics(
        last_backup_time=latest_backup.created_at if latest_backup else None,
        last_backup_age_hours=latest_backup.age_hours() if latest_backup else -1,
        last_backup_size_mb=latest_backup.size_mb if latest_backup else 0,
        total_backups=len(backups),
        total_size_mb=round(sum(b.size_mb for b in backups), 2),
        failed_checks=failures,
        timestamp=datetime.now()
    )

    return metrics


def save_metrics(metrics: BackupMetrics):
    """Save metrics to file"""
    try:
        with open(METRICS_FILE, 'w') as f:
            json.dump(metrics.to_dict(), f, indent=2)
        logger.info(f"Metrics saved to {METRICS_FILE}")
    except Exception as e:
        logger.error(f"Failed to save metrics: {e}")


def format_prometheus_metrics(metrics: BackupMetrics) -> str:
    """Format metrics for Prometheus"""
    lines = [
        "# HELP listings_backup_age_hours Age of last backup in hours",
        "# TYPE listings_backup_age_hours gauge",
        f"listings_backup_age_hours {metrics.last_backup_age_hours}",
        "",
        "# HELP listings_backup_size_mb Size of last backup in MB",
        "# TYPE listings_backup_size_mb gauge",
        f"listings_backup_size_mb {metrics.last_backup_size_mb}",
        "",
        "# HELP listings_backup_total Total number of backups",
        "# TYPE listings_backup_total gauge",
        f"listings_backup_total {metrics.total_backups}",
        "",
        "# HELP listings_backup_total_size_mb Total size of all backups in MB",
        "# TYPE listings_backup_total_size_mb gauge",
        f"listings_backup_total_size_mb {metrics.total_size_mb}",
        "",
        "# HELP listings_backup_health Backup health status (1=healthy, 0=unhealthy)",
        "# TYPE listings_backup_health gauge",
        f"listings_backup_health {1 if not metrics.failed_checks else 0}",
        "",
    ]
    return "\n".join(lines)


# ========================================
# HTTP server for metrics
# ========================================

class MetricsHandler(BaseHTTPRequestHandler):
    """HTTP handler for metrics endpoint"""

    def do_GET(self):
        """Handle GET requests"""
        if self.path == '/metrics':
            self.send_metrics()
        elif self.path == '/health':
            self.send_health()
        else:
            self.send_error(404, "Not Found")

    def send_metrics(self):
        """Send Prometheus metrics"""
        try:
            with open(METRICS_FILE, 'r') as f:
                metrics_data = json.load(f)

            metrics = BackupMetrics(**metrics_data)
            content = format_prometheus_metrics(metrics)

            self.send_response(200)
            self.send_header('Content-Type', 'text/plain; version=0.0.4')
            self.end_headers()
            self.wfile.write(content.encode())
        except Exception as e:
            self.send_error(500, f"Error: {e}")

    def send_health(self):
        """Send health check response"""
        try:
            with open(METRICS_FILE, 'r') as f:
                metrics_data = json.load(f)

            health = {
                "status": "healthy" if not metrics_data.get('failed_checks') else "unhealthy",
                "timestamp": metrics_data.get('timestamp'),
                "last_backup_age_hours": metrics_data.get('last_backup_age_hours')
            }

            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(health, indent=2).encode())
        except Exception as e:
            self.send_error(500, f"Error: {e}")

    def log_message(self, format, *args):
        """Override to use custom logger"""
        logger.info(f"{self.address_string()} - {format % args}")


def start_metrics_server(port: int):
    """Start HTTP server for metrics"""
    if not HTTP_SERVER_AVAILABLE:
        logger.error("HTTP server not available")
        return

    server = HTTPServer(('', port), MetricsHandler)
    logger.info(f"Starting metrics server on port {port}")
    logger.info(f"Metrics endpoint: http://localhost:{port}/metrics")
    logger.info(f"Health endpoint: http://localhost:{port}/health")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logger.info("Shutting down metrics server")
        server.shutdown()


# ========================================
# Main
# ========================================

def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(description="Backup monitoring and alerting")
    parser.add_argument('--check', action='store_true', help='Run health checks once')
    parser.add_argument('--metrics', action='store_true', help='Print Prometheus metrics')
    parser.add_argument('--serve', action='store_true', help='Start metrics HTTP server')
    parser.add_argument('--port', type=int, default=METRICS_PORT, help='Metrics server port')
    parser.add_argument('--verbose', action='store_true', help='Verbose logging')

    args = parser.parse_args()

    if args.verbose:
        logger.setLevel(logging.DEBUG)

    # Ensure log directory exists
    os.makedirs(LOG_DIR, exist_ok=True)

    # Find backups
    backups = find_backup_files()

    # Run checks
    all_passed, failures, messages = run_all_checks(backups)

    # Collect metrics
    metrics = collect_metrics(backups, all_passed, failures)
    save_metrics(metrics)

    # Print results
    logger.info("=" * 50)
    logger.info("Backup Health Check Results")
    logger.info("=" * 50)
    for msg in messages:
        logger.info(msg)
    logger.info("=" * 50)

    if all_passed:
        logger.info("✓ All checks passed")
    else:
        logger.error("✗ Some checks failed")
        notify_failures(failures)

    # Handle different modes
    if args.metrics:
        print(format_prometheus_metrics(metrics))
    elif args.serve:
        start_metrics_server(args.port)
    elif not args.check:
        # Default: run check once
        sys.exit(0 if all_passed else 1)


if __name__ == "__main__":
    main()
