# Listings Backup System - Quick Installation Guide

## Prerequisites

- PostgreSQL database running in Docker (container: `listings_postgres`)
- Database: `listings_dev_db` on port `35434`
- User with backup privileges: `listings_user`
- Python 3.6+ (for monitoring script)
- `psql`, `pg_dump`, `gzip`, `sha256sum` utilities

## Quick Installation (5 minutes)

### 1. Create Backup User (if not exists)

```bash
# Create system user for running backups
sudo useradd -r -s /bin/bash listings

# Create home directory
sudo mkdir -p /home/listings
sudo chown listings:listings /home/listings
```

### 2. Create Required Directories

```bash
# Create backup and log directories
sudo mkdir -p /var/backups/listings/{daily,weekly,monthly,wal,pre-restore}
sudo mkdir -p /var/log/listings/reports

# Set ownership
sudo chown -R listings:listings /var/backups/listings
sudo chown -R listings:listings /var/log/listings

# Set permissions
sudo chmod 750 /var/backups/listings
sudo chmod 750 /var/log/listings
```

### 3. Configure Environment

```bash
# Create configuration file
sudo tee /etc/listings-backup.env > /dev/null <<'EOF'
# Database Configuration
BACKUP_DB_HOST=localhost
BACKUP_DB_PORT=35434
BACKUP_DB_NAME=listings_dev_db
BACKUP_DB_USER=listings_user
BACKUP_DB_PASSWORD=YOUR_PASSWORD_HERE
BACKUP_DB_CONTAINER=listings_postgres

# Backup Configuration
BACKUP_DIR=/var/backups/listings
LOG_DIR=/var/log/listings

# Retention Policy
BACKUP_RETENTION_DAYS=7
BACKUP_RETENTION_WEEKS=4
BACKUP_RETENTION_MONTHS=12

# S3/MinIO (Optional - set BACKUP_ENABLE_S3=true to enable)
BACKUP_ENABLE_S3=false
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=listings-backups
S3_USE_SSL=false

# Notifications (Optional)
BACKUP_NOTIFY_EMAIL=
SLACK_WEBHOOK_URL=
EOF

# Set secure permissions
sudo chmod 600 /etc/listings-backup.env
sudo chown listings:listings /etc/listings-backup.env
```

**⚠️ IMPORTANT**: Edit `/etc/listings-backup.env` and set `BACKUP_DB_PASSWORD`!

```bash
sudo nano /etc/listings-backup.env
```

### 4. Install Scripts

Scripts are already in the project at:
```
/p/github.com/sveturs/listings/scripts/backup/
```

Make them accessible system-wide (optional):
```bash
# Create symbolic links
sudo ln -s /p/github.com/sveturs/listings/scripts/backup/backup-db.sh /usr/local/bin/listings-backup
sudo ln -s /p/github.com/sveturs/listings/scripts/backup/restore-db.sh /usr/local/bin/listings-restore
sudo ln -s /p/github.com/sveturs/listings/scripts/backup/verify-backup.sh /usr/local/bin/listings-verify
```

### 5. Setup Automated Backups

```bash
cd /p/github.com/sveturs/listings/scripts/backup/

# Setup cron jobs (automated)
sudo ./setup-cron.sh

# Or manual cron setup
sudo tee /etc/cron.d/listings-backup > /dev/null <<'EOF'
# Listings Database Backup Cron Jobs

SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

# Daily backup at 2:00 AM
0 2 * * * listings /p/github.com/sveturs/listings/scripts/backup/backup-db.sh >> /var/log/listings/cron.log 2>&1

# Weekly verification on Sunday at 6:00 AM
0 6 * * 0 listings /p/github.com/sveturs/listings/scripts/backup/verify-backup.sh --verify-all --quick >> /var/log/listings/cron.log 2>&1

# Hourly monitoring
0 * * * * listings /p/github.com/sveturs/listings/scripts/backup/monitor-backups.py >> /var/log/listings/monitor.log 2>&1
EOF

sudo chmod 644 /etc/cron.d/listings-backup
```

### 6. Test Installation

```bash
# Source environment
set -a
source /etc/listings-backup.env
set +a

# Test backup (dry run)
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/backup-db.sh --dry-run

# If dry run succeeds, create first backup
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/backup-db.sh

# Check backup was created
ls -lh /var/backups/listings/daily/

# Verify backup integrity
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/verify-backup.sh \
    /var/backups/listings/daily/*.sql.gz --quick

# Run integration test (optional)
cd /p/github.com/sveturs/listings/scripts/backup/
sudo -u listings TEST_DB_PASSWORD="$BACKUP_DB_PASSWORD" ./test-backup-restore.sh
```

### 7. Setup Monitoring (Optional)

#### Option A: Standalone Metrics Server

```bash
# Start metrics server
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/monitor-backups.py --serve --port 9090 &

# Test metrics endpoint
curl http://localhost:9090/metrics
curl http://localhost:9090/health
```

#### Option B: Systemd Service

```bash
sudo tee /etc/systemd/system/listings-backup-monitor.service > /dev/null <<'EOF'
[Unit]
Description=Listings Backup Monitoring Service
After=network.target

[Service]
Type=simple
User=listings
EnvironmentFile=/etc/listings-backup.env
ExecStart=/p/github.com/sveturs/listings/scripts/backup/monitor-backups.py --serve --port 9090
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable listings-backup-monitor
sudo systemctl start listings-backup-monitor

# Check status
sudo systemctl status listings-backup-monitor
```

### 8. Configure Notifications (Optional)

#### Email Notifications

```bash
# Install mailutils
sudo apt-get install mailutils

# Edit environment file
sudo nano /etc/listings-backup.env

# Add:
BACKUP_NOTIFY_EMAIL=your-email@example.com
SMTP_HOST=localhost
SMTP_PORT=25
```

#### Slack Notifications

```bash
# Get Slack webhook URL from: https://api.slack.com/messaging/webhooks

# Edit environment file
sudo nano /etc/listings-backup.env

# Add:
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

Install Python requests library:
```bash
pip3 install requests
```

## Verification Checklist

After installation, verify everything works:

- [ ] Backup user created: `id listings`
- [ ] Directories exist: `ls -ld /var/backups/listings /var/log/listings`
- [ ] Environment file configured: `sudo cat /etc/listings-backup.env`
- [ ] Scripts executable: `ls -lh /p/github.com/sveturs/listings/scripts/backup/*.sh`
- [ ] Cron jobs configured: `sudo cat /etc/cron.d/listings-backup`
- [ ] Test backup successful: `ls -lh /var/backups/listings/daily/`
- [ ] Verification passed: `grep "PASSED" /var/log/listings/verify-backup.log`
- [ ] Monitoring working: `curl http://localhost:9090/health`
- [ ] Logs readable: `tail /var/log/listings/backup.log`

## Quick Commands Reference

```bash
# Manual backup
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/backup-db.sh

# Restore from backup
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/restore-db.sh \
    /var/backups/listings/daily/backup.sql.gz

# Verify backup
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/verify-backup.sh \
    /var/backups/listings/daily/backup.sql.gz

# Check monitoring
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/monitor-backups.py --check

# View logs
tail -f /var/log/listings/backup.log
tail -f /var/log/listings/monitor-backups.log

# List backups
ls -lh /var/backups/listings/daily/
ls -lh /var/backups/listings/weekly/
ls -lh /var/backups/listings/monthly/

# Check cron status
sudo grep listings /var/log/syslog | tail -20
```

## Troubleshooting

### Permission Denied Errors

```bash
# Fix ownership
sudo chown -R listings:listings /var/backups/listings
sudo chown -R listings:listings /var/log/listings

# Fix permissions
sudo chmod 750 /var/backups/listings
sudo chmod 640 /var/log/listings/*.log
```

### Database Connection Failed

```bash
# Test database connection
PGPASSWORD=your_password psql -h localhost -p 35434 -U listings_user -d listings_dev_db -c "SELECT 1;"

# Check Docker container
docker ps | grep listings_postgres
docker logs listings_postgres | tail -20
```

### Backup Not Running Automatically

```bash
# Check cron service
sudo systemctl status cron

# Check cron logs
sudo grep CRON /var/log/syslog | grep listings

# Test cron job manually
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/backup-db.sh
```

### Disk Space Issues

```bash
# Check available space
df -h /var/backups/listings

# Clean old backups
find /var/backups/listings/daily -name "*.sql.gz" -mtime +7 -delete

# Adjust retention
sudo nano /etc/listings-backup.env
# Set BACKUP_RETENTION_DAYS=3
```

## Uninstallation

If you need to remove the backup system:

```bash
# Stop services
sudo systemctl stop listings-backup-monitor
sudo systemctl disable listings-backup-monitor

# Remove cron jobs
sudo rm /etc/cron.d/listings-backup

# Remove systemd services
sudo rm /etc/systemd/system/listings-backup*.{service,timer}
sudo systemctl daemon-reload

# Remove configuration
sudo rm /etc/listings-backup.env

# Remove backups (CAUTION!)
# sudo rm -rf /var/backups/listings

# Remove logs
# sudo rm -rf /var/log/listings

# Remove user (optional)
# sudo userdel listings
```

## Next Steps

1. **Configure S3 backup** (for off-site storage)
   - Edit `/etc/listings-backup.env`
   - Set `BACKUP_ENABLE_S3=true` and S3 credentials
   - Test: `sudo -u listings ./backup-s3.sh /var/backups/listings/daily/*.sql.gz`

2. **Setup monitoring dashboard**
   - Configure Prometheus to scrape `http://localhost:9090/metrics`
   - Create Grafana dashboard for backup metrics

3. **Document restore procedures**
   - Train team on restore process
   - Run quarterly disaster recovery drills

4. **Review backup policy**
   - Read [BACKUP_POLICY.md](BACKUP_POLICY.md)
   - Adjust retention based on requirements

## Support

- Documentation: [README.md](README.md)
- Policy: [BACKUP_POLICY.md](BACKUP_POLICY.md)
- Logs: `/var/log/listings/`
- Issues: Contact DevOps team

## Installation Complete! ✅

Your backup system is now installed and configured. The first automated backup will run at 2:00 AM tomorrow.

To create a backup immediately:
```bash
sudo -u listings /p/github.com/sveturs/listings/scripts/backup/backup-db.sh
```
