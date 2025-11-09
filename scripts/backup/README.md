# Listings Microservice Backup System

Production-ready backup and restore system for the Listings microservice PostgreSQL database.

## Features

✅ **Automated Backups**
- Daily, weekly, and monthly full backups
- Write-Ahead Log (WAL) archiving for Point-in-Time Recovery
- Configurable retention policies
- Automatic cleanup of old backups

✅ **Reliable Restore**
- Restore from any backup with verification
- Point-in-Time Recovery (PITR) support
- Pre-restore safety backup
- Automatic rollback on failure

✅ **Verification**
- Test restore to temporary database
- Data integrity checks
- Checksum verification
- Automated verification reports

✅ **Remote Storage**
- S3/MinIO upload support
- Encryption at rest and in transit
- Multi-part upload for large files
- Automatic retention cleanup

✅ **Monitoring**
- Real-time health checks
- Prometheus metrics export
- Slack/email notifications
- Anomaly detection

✅ **Security**
- Encrypted credentials
- Access control (file permissions)
- Lock files to prevent concurrent operations
- Audit logging

## Quick Start

### 1. Installation

```bash
# Clone or copy scripts to server
cd /p/github.com/sveturs/listings/scripts/backup/

# Make scripts executable
chmod +x *.sh *.py

# Install dependencies
pip3 install requests  # For Slack notifications (optional)

# Setup automated backups
sudo ./setup-cron.sh
```

### 2. Configuration

Edit the environment file:
```bash
sudo nano /etc/listings-backup.env
```

Set required variables:
```bash
BACKUP_DB_PASSWORD=your_secure_password
BACKUP_NOTIFY_EMAIL=your-email@example.com
```

### 3. Test Backup

Run a test backup:
```bash
sudo -u listings ./backup-db.sh --dry-run
sudo -u listings ./backup-db.sh
```

### 4. Verify

Check logs:
```bash
tail -f /var/log/listings/backup.log
```

View backup files:
```bash
ls -lh /var/backups/listings/daily/
```

## Usage

### Create Backup

```bash
# Manual backup (uses current date to determine type)
./backup-db.sh

# Dry run (test without creating backup)
./backup-db.sh --dry-run

# Full backup without WAL archiving
./backup-db.sh --full-only

# Backup without S3 upload
./backup-db.sh --no-upload
```

### Restore Database

```bash
# Standard restore
./restore-db.sh /var/backups/listings/daily/backup.sql.gz

# Dry run (test without restoring)
./restore-db.sh backup.sql.gz --dry-run

# Skip pre-restore backup (faster, but risky)
./restore-db.sh backup.sql.gz --no-backup

# Point-in-Time Recovery
./restore-db.sh backup.sql.gz --pitr-target "2024-11-05 12:00:00"
```

### Verify Backup

```bash
# Verify single backup
./verify-backup.sh /var/backups/listings/daily/backup.sql.gz

# Quick check (faster, less thorough)
./verify-backup.sh backup.sql.gz --quick

# Verify all backups
./verify-backup.sh --verify-all

# Parallel verification (faster)
./verify-backup.sh --verify-all --parallel
```

### Upload to S3/MinIO

```bash
# Upload backup to S3
./backup-s3.sh /var/backups/listings/daily/backup.sql.gz

# Dry run
./backup-s3.sh backup.sql.gz --dry-run
```

### Monitor Backups

```bash
# Run health checks once
./monitor-backups.py --check

# Print Prometheus metrics
./monitor-backups.py --metrics

# Start metrics HTTP server (port 9090)
./monitor-backups.py --serve

# Custom port
./monitor-backups.py --serve --port 8080
```

### Integration Test

```bash
# Run full integration test
./test-backup-restore.sh

# Keep backup after test
./test-backup-restore.sh --keep-backup

# Skip cleanup (for debugging)
./test-backup-restore.sh --skip-cleanup
```

## Configuration

### Environment Variables

All scripts use these environment variables (can be set in `/etc/listings-backup.env`):

#### Database Configuration
```bash
BACKUP_DB_HOST=localhost          # Database host
BACKUP_DB_PORT=35434              # Database port
BACKUP_DB_NAME=listings_dev_db    # Database name
BACKUP_DB_USER=listings_user      # Database user
BACKUP_DB_PASSWORD=secret         # Database password (required)
BACKUP_DB_CONTAINER=listings_postgres  # Docker container name
```

#### Backup Configuration
```bash
BACKUP_DIR=/var/backups/listings  # Backup directory
LOG_DIR=/var/log/listings         # Log directory
BACKUP_RETENTION_DAYS=7           # Daily backups to keep
BACKUP_RETENTION_WEEKS=4          # Weekly backups to keep
BACKUP_RETENTION_MONTHS=12        # Monthly backups to keep
```

#### S3/MinIO Configuration
```bash
BACKUP_ENABLE_S3=false            # Enable S3 upload
S3_ENDPOINT=localhost:9000        # S3/MinIO endpoint
S3_ACCESS_KEY=your_key            # S3 access key
S3_SECRET_KEY=your_secret         # S3 secret key
S3_BUCKET=listings-backups        # S3 bucket name
S3_REGION=us-east-1               # S3 region
S3_USE_SSL=false                  # Use SSL for S3
S3_RETENTION_DAYS=30              # S3 retention days
```

#### Notification Configuration
```bash
BACKUP_NOTIFY_EMAIL=email@example.com  # Email for notifications
SMTP_HOST=localhost                    # SMTP server
SMTP_PORT=25                           # SMTP port
SLACK_WEBHOOK_URL=https://hooks.slack.com/...  # Slack webhook
```

#### Monitoring Configuration
```bash
MAX_BACKUP_AGE_HOURS=25           # Alert if backup older than
MIN_BACKUP_SIZE_MB=1              # Alert if backup smaller than
SIZE_CHANGE_THRESHOLD_PCT=50      # Alert if size change > %
```

### Cron Schedule

Default schedule (configured by `setup-cron.sh`):

```bash
# Daily backup at 2:00 AM
0 2 * * * listings /path/to/backup-db.sh

# Weekly verification on Sunday at 6:00 AM
0 6 * * 0 listings /path/to/verify-backup.sh --verify-all --quick

# Hourly monitoring
0 * * * * listings /path/to/monitor-backups.py
```

Customize schedule:
```bash
# Edit cron file
sudo nano /etc/cron.d/listings-backup

# Or use systemd timers
sudo systemctl enable --now listings-backup.timer
sudo systemctl enable --now listings-backup-verify.timer
```

## Architecture

### Directory Structure

```
/var/backups/listings/
├── daily/              # Daily backups (7 days retention)
├── weekly/             # Weekly backups (4 weeks retention)
├── monthly/            # Monthly backups (12 months retention)
├── wal/                # WAL archive for PITR
└── pre-restore/        # Pre-restore safety backups

/var/log/listings/
├── backup.log          # Backup operation logs
├── restore.log         # Restore operation logs
├── verify-backup.log   # Verification logs
├── monitor-backups.log # Monitoring logs
└── reports/            # Verification reports
```

### Backup File Format

Each backup consists of:
- `*.sql.gz` - Compressed SQL dump
- `*.sql.gz.sha256` - SHA256 checksum
- `*.sql.gz.meta` - Metadata (JSON)

Example:
```
listings_dev_db_20241105_020000.sql.gz
listings_dev_db_20241105_020000.sql.gz.sha256
listings_dev_db_20241105_020000.sql.gz.meta
```

Metadata file:
```json
{
  "backup_date": "2024-11-05 02:00:00",
  "backup_type": "daily",
  "database_name": "listings_dev_db",
  "database_size": "150MB",
  "backup_size": "45MB",
  "duration_seconds": 120,
  "host": "localhost",
  "port": "35434"
}
```

## Monitoring

### Prometheus Metrics

Metrics exposed on `http://localhost:9090/metrics`:

```
# Backup age in hours
listings_backup_age_hours 2.5

# Backup size in MB
listings_backup_size_mb 45.2

# Total number of backups
listings_backup_total 21

# Total size of all backups in MB
listings_backup_total_size_mb 950.4

# Backup health (1=healthy, 0=unhealthy)
listings_backup_health 1
```

### Health Check

Health check endpoint on `http://localhost:9090/health`:

```json
{
  "status": "healthy",
  "timestamp": "2024-11-05T02:00:00",
  "last_backup_age_hours": 2.5
}
```

### Alerts

Automatic alerts sent when:
- Last backup is older than 25 hours
- Backup size is unusually small (< 1MB)
- Backup size changes by more than 50%
- Backup creation fails
- Restore operation fails

Notification channels:
- Slack (if `SLACK_WEBHOOK_URL` configured)
- Email (if `BACKUP_NOTIFY_EMAIL` configured)
- Logs (always)

## Troubleshooting

### Backup fails with "insufficient disk space"

**Problem**: Not enough disk space for backup.

**Solution**:
```bash
# Check available space
df -h /var/backups/listings

# Clean up old backups
./backup-db.sh  # Will auto-cleanup based on retention

# Or manually
find /var/backups/listings/daily -name "*.sql.gz" -mtime +7 -delete
```

### Restore fails with "too many clients already"

**Problem**: Database connection limit reached.

**Solution**:
```bash
# Terminate existing connections
PGPASSWORD=password psql -h localhost -p 35434 -U listings_user -d postgres -c "
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'listings_dev_db';
"

# Restart PostgreSQL
docker restart listings_postgres
```

### Checksum verification fails

**Problem**: Backup file corrupted or checksum missing.

**Solution**:
```bash
# Re-generate checksum
sha256sum backup.sql.gz > backup.sql.gz.sha256

# Or verify manually
gzip -t backup.sql.gz  # Test gzip integrity
```

### Backup script hangs

**Problem**: Lock file from previous run or deadlock.

**Solution**:
```bash
# Check for stale lock file
ls -l /var/run/listings-backup.lock

# Remove if stale (check PID first)
cat /var/run/listings-backup.lock  # Shows PID
ps -p <PID>  # Check if process exists
rm /var/run/listings-backup.lock  # If stale
```

### S3 upload fails

**Problem**: Credentials or network issues.

**Solution**:
```bash
# Test S3 credentials
aws s3 ls s3://listings-backups/ \
    --endpoint-url http://localhost:9000

# Check logs
tail -f /var/log/listings/backup-s3.log

# Upload manually
./backup-s3.sh /var/backups/listings/daily/backup.sql.gz
```

## Best Practices

### Security
1. ✅ Store credentials in `/etc/listings-backup.env` with permissions `600`
2. ✅ Run backups as dedicated user (`listings`)
3. ✅ Use encrypted filesystems for backups
4. ✅ Enable S3 server-side encryption
5. ✅ Rotate access keys regularly
6. ✅ Audit backup access logs

### Reliability
1. ✅ Test restores monthly
2. ✅ Monitor backup health continuously
3. ✅ Keep multiple backup types (daily/weekly/monthly)
4. ✅ Store backups off-site (S3/MinIO)
5. ✅ Document restore procedures
6. ✅ Run quarterly disaster recovery drills

### Performance
1. ✅ Schedule backups during low-traffic hours (2 AM)
2. ✅ Use compression (gzip -9)
3. ✅ Use parallel verification when possible
4. ✅ Clean up old backups regularly
5. ✅ Monitor disk I/O during backups
6. ✅ Use incremental WAL archiving

## Scripts Reference

| Script | Purpose | Usage |
|--------|---------|-------|
| `backup-db.sh` | Create database backup | `./backup-db.sh [--dry-run]` |
| `restore-db.sh` | Restore from backup | `./restore-db.sh <file> [--pitr-target]` |
| `verify-backup.sh` | Verify backup integrity | `./verify-backup.sh <file> [--quick]` |
| `backup-s3.sh` | Upload to S3/MinIO | `./backup-s3.sh <file>` |
| `monitor-backups.py` | Monitor backup health | `./monitor-backups.py [--check]` |
| `setup-cron.sh` | Setup automation | `sudo ./setup-cron.sh` |
| `test-backup-restore.sh` | Integration test | `./test-backup-restore.sh` |

## Documentation

- **[BACKUP_POLICY.md](BACKUP_POLICY.md)** - Comprehensive backup policy and procedures
- **Logs**: `/var/log/listings/`
- **Configuration**: `/etc/listings-backup.env`
- **Cron**: `/etc/cron.d/listings-backup`
- **Systemd**: `/etc/systemd/system/listings-backup.{service,timer}`

## Support

### Getting Help

1. Check logs: `tail -f /var/log/listings/backup.log`
2. Review troubleshooting section above
3. Run with verbose logging: `./backup-db.sh --dry-run`
4. Test in isolation: `./test-backup-restore.sh`

### Reporting Issues

When reporting issues, include:
- Script version and command used
- Error message from logs
- Database and system information
- Steps to reproduce

### Contact

- DevOps Team: devops@example.com
- On-call Engineer: Use PagerDuty for emergencies

## License

Internal use only. See project LICENSE file.

## Changelog

### v1.0.0 (2024-11-05)
- Initial release
- Full backup and restore functionality
- S3/MinIO upload support
- Monitoring and alerting
- Automated scheduling
- Integration tests
- Comprehensive documentation
