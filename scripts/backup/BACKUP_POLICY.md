# Listings Microservice Backup Policy

## Overview

This document describes the backup and disaster recovery strategy for the Listings microservice database.

## Backup Strategy

### Backup Types

| Type | Frequency | Retention | Description |
|------|-----------|-----------|-------------|
| **Daily** | Every day at 2:00 AM | 7 days | Full database dump with compression |
| **Weekly** | Sunday at 2:00 AM | 4 weeks | Full database dump for weekly archival |
| **Monthly** | 1st of month at 2:00 AM | 12 months | Full database dump for long-term archival |
| **WAL Archive** | Continuous | 7 days | Write-Ahead Log files for Point-in-Time Recovery (PITR) |

### Backup Contents

All backups include:
- Full database schema (tables, indexes, constraints, sequences)
- All data from critical tables:
  - `listings` - Main listings data
  - `products` - Product catalog
  - `inventory` - Stock levels
  - `inventory_movements` - Stock history
  - `listing_images` - Image metadata
  - `listing_attributes` - Custom attributes
- Database users and permissions (excluding passwords)
- Extension configurations

### Storage Locations

#### Local Storage
- **Primary**: `/var/backups/listings/`
  - `/daily/` - Daily backups
  - `/weekly/` - Weekly backups
  - `/monthly/` - Monthly backups
  - `/wal/` - WAL archive files
  - `/pre-restore/` - Pre-restore safety backups

#### Remote Storage (Optional)
- **S3/MinIO**: `s3://listings-backups/`
  - Organized by date: `backups/YYYY/MM/DD/`
  - Versioning enabled
  - Lifecycle policy: 30 days

### Backup Format

- **Format**: Plain SQL (pg_dump --format=plain)
- **Compression**: gzip -9 (maximum compression)
- **Checksums**: SHA256 for each backup file
- **Metadata**: JSON file with backup information

Example backup structure:
```
/var/backups/listings/daily/
├── listings_dev_db_20241105_020000.sql.gz       # Compressed backup
├── listings_dev_db_20241105_020000.sql.gz.sha256 # Checksum
└── listings_dev_db_20241105_020000.sql.gz.meta   # Metadata
```

## Recovery Objectives

### Recovery Time Objective (RTO)
- **Target**: < 1 hour
- **Maximum**: 4 hours

**Breakdown**:
1. Incident detection: 5 minutes
2. Backup retrieval: 10 minutes
3. Database restore: 30 minutes (depends on size)
4. Verification: 10 minutes
5. Service restart: 5 minutes

### Recovery Point Objective (RPO)
- **Target**: < 1 hour
- **Maximum**: 24 hours

With daily backups + WAL archiving, we can recover to:
- Any point in time within the last 24 hours (using PITR)
- Any daily backup within the last 7 days
- Any weekly backup within the last 4 weeks
- Any monthly backup within the last 12 months

## Backup Process

### Automated Backups

Backups are automated via cron jobs:

```bash
# Daily backup at 2:00 AM
0 2 * * * listings /path/to/backup-db.sh

# Weekly verification on Sunday at 6:00 AM
0 6 * * 0 listings /path/to/verify-backup.sh --verify-all --quick

# Hourly monitoring
0 * * * * listings /path/to/monitor-backups.py
```

### Backup Steps

1. **Lock acquisition** - Prevent concurrent backups
2. **Disk space check** - Ensure sufficient space (2x database size)
3. **Database dump** - pg_dump with compression
4. **Checksum generation** - SHA256 for integrity
5. **Metadata creation** - Record backup details
6. **S3 upload** (optional) - Upload to remote storage
7. **Retention cleanup** - Remove old backups
8. **Notification** - Send success/failure alerts
9. **Lock release** - Allow next backup

### Backup Verification

Weekly automated verification:
1. File integrity check (gzip -t, checksum)
2. Test restore to temporary database
3. Data integrity checks (row counts, critical tables)
4. Generate verification report
5. Cleanup temporary resources

Manual verification recommended:
- After major database changes
- Before production deployments
- Quarterly full restore tests

## Restore Process

### Standard Restore

```bash
# 1. Stop services using the database
sudo systemctl stop listings

# 2. Run restore script
sudo -u listings /path/to/restore-db.sh \
    /var/backups/listings/daily/backup.sql.gz

# 3. Verify restore
sudo -u listings /path/to/verify-backup.sh --quick

# 4. Restart services
sudo systemctl start listings
```

### Point-in-Time Recovery (PITR)

```bash
# Restore to specific timestamp
sudo -u listings /path/to/restore-db.sh \
    /var/backups/listings/daily/backup.sql.gz \
    --pitr-target "2024-11-05 12:00:00"
```

### Emergency Restore

In case of critical failure:

1. **Assess situation**
   - Identify cause of failure
   - Determine required recovery point

2. **Select backup**
   ```bash
   # List available backups
   ls -lh /var/backups/listings/daily/
   ```

3. **Restore database**
   ```bash
   # Use --no-backup flag to skip pre-restore backup
   sudo -u listings /path/to/restore-db.sh \
       /var/backups/listings/daily/latest.sql.gz \
       --no-backup
   ```

4. **Verify data**
   - Check critical tables
   - Verify data consistency
   - Test application functionality

5. **Resume operations**
   - Restart services
   - Monitor for issues
   - Notify stakeholders

## Monitoring and Alerting

### Automated Monitoring

The `monitor-backups.py` script runs hourly and checks:

1. **Backup age** - Alert if last backup > 25 hours old
2. **Backup size** - Alert if size < 1MB or anomalous change > 50%
3. **File integrity** - Verify checksums
4. **Disk space** - Alert if < 10GB available

### Alerts

Notifications sent via:
- **Slack** - Real-time alerts (if configured)
- **Email** - Detailed reports (if configured)
- **Logs** - Always logged to `/var/log/listings/`

Alert levels:
- **Critical**: Backup failed, backup too old (>25h)
- **Warning**: Size anomaly, missing checksum
- **Info**: Successful backup, verification passed

### Metrics

Prometheus metrics exposed on port 9090:
- `listings_backup_age_hours` - Age of last backup
- `listings_backup_size_mb` - Size of last backup
- `listings_backup_total` - Total number of backups
- `listings_backup_health` - Health status (1=healthy, 0=unhealthy)

## Security

### Access Control

- Backup files owned by `listings:listings` with permissions `640`
- Backup directory permissions: `750`
- Database credentials stored in environment file (`/etc/listings-backup.env`) with permissions `600`

### Encryption

- **At rest**:
  - Local backups on encrypted filesystem (LUKS)
  - S3 backups with server-side encryption (SSE-S3)
- **In transit**:
  - S3 uploads over HTTPS
  - Database connections use SSL (if configured)

### Credential Management

Database password stored in:
1. Environment file: `/etc/listings-backup.env`
2. Or environment variable: `BACKUP_DB_PASSWORD`
3. Never in scripts or version control

## Testing and Validation

### Regular Testing

1. **Weekly** - Automated verification of latest backups
2. **Monthly** - Manual restore test to staging environment
3. **Quarterly** - Full disaster recovery drill

### Test Checklist

- [ ] Backup creation successful
- [ ] File integrity verified (checksum)
- [ ] Test restore to temporary database
- [ ] Data integrity checks passed
- [ ] Critical tables have expected row counts
- [ ] S3 upload successful (if enabled)
- [ ] Monitoring alerts working
- [ ] Documentation up to date

### Disaster Recovery Drill

Quarterly full DR drill:
1. Simulate catastrophic failure
2. Restore from backup to fresh server
3. Verify all data and functionality
4. Measure RTO/RPO achieved
5. Document lessons learned
6. Update procedures as needed

## Maintenance

### Regular Tasks

**Daily**:
- Review backup logs for errors
- Check monitoring dashboard

**Weekly**:
- Verify backup verification reports
- Check disk space utilization

**Monthly**:
- Review retention policy
- Test restore procedure
- Update documentation

**Quarterly**:
- Full disaster recovery drill
- Review and update backup strategy
- Audit access controls

### Log Rotation

Logs rotated daily, kept for 30 days:
- `/var/log/listings/backup.log`
- `/var/log/listings/restore.log`
- `/var/log/listings/verify-backup.log`
- `/var/log/listings/monitor-backups.log`

## Troubleshooting

### Common Issues

#### Backup fails with "insufficient disk space"
```bash
# Check available space
df -h /var/backups/listings

# Clean up old backups manually
find /var/backups/listings/daily -name "*.sql.gz" -mtime +7 -delete

# Or increase retention
export BACKUP_RETENTION_DAYS=3
```

#### Restore fails with "database already exists"
```bash
# Drop database first (CAUTION!)
PGPASSWORD=$DB_PASSWORD psql -h localhost -p 35434 -U listings_user -d postgres \
    -c "DROP DATABASE listings_dev_db;"

# Or use restore script (it handles this)
./restore-db.sh backup.sql.gz
```

#### Checksum verification fails
```bash
# Re-generate checksum
sha256sum backup.sql.gz > backup.sql.gz.sha256

# Or skip verification
VERIFY_CHECKSUM=false ./verify-backup.sh backup.sql.gz
```

#### WAL archiving not working
```bash
# Check PostgreSQL configuration
docker exec listings_postgres cat /var/lib/postgresql/data/postgresql.conf | grep archive

# Enable WAL archiving
docker exec listings_postgres bash -c "
echo 'archive_mode = on' >> /var/lib/postgresql/data/postgresql.conf
echo 'archive_command = 'cp %p /var/lib/postgresql/wal_archive/%f'' >> /var/lib/postgresql/data/postgresql.conf
"

# Restart PostgreSQL
docker restart listings_postgres
```

## Scripts Reference

### backup-db.sh
Full database backup with retention policy.

**Usage**:
```bash
./backup-db.sh [--dry-run] [--no-upload] [--full-only]
```

**Environment**:
- `BACKUP_DB_PASSWORD` - Database password (required)
- `BACKUP_DIR` - Backup directory (default: /var/backups/listings)
- `BACKUP_RETENTION_DAYS` - Daily retention (default: 7)
- `BACKUP_ENABLE_S3` - Upload to S3 (default: false)

### restore-db.sh
Restore database from backup with verification.

**Usage**:
```bash
./restore-db.sh <backup_file> [--dry-run] [--no-backup] [--pitr-target "YYYY-MM-DD HH:MM:SS"]
```

**Features**:
- Pre-restore backup (can be disabled with `--no-backup`)
- Point-in-time recovery support
- Post-restore verification
- Rollback on failure

### verify-backup.sh
Verify backup integrity with test restore.

**Usage**:
```bash
./verify-backup.sh <backup_file> [--dry-run] [--quick]
./verify-backup.sh --verify-all [--parallel]
```

**Features**:
- File integrity checks
- Test restore to temporary database
- Data integrity validation
- Generate verification reports

### backup-s3.sh
Upload backup to S3/MinIO storage.

**Usage**:
```bash
./backup-s3.sh <backup_file> [--dry-run]
```

**Environment**:
- `S3_ACCESS_KEY` - S3 access key (required)
- `S3_SECRET_KEY` - S3 secret key (required)
- `S3_ENDPOINT` - S3 endpoint (default: localhost:9000)
- `S3_BUCKET` - S3 bucket (default: listings-backups)

### monitor-backups.py
Monitor backup health and generate alerts.

**Usage**:
```bash
./monitor-backups.py [--check] [--metrics] [--serve]
```

**Features**:
- Check backup age and size
- Detect anomalies
- Slack/email alerts
- Prometheus metrics

### setup-cron.sh
Setup automated backup cron jobs.

**Usage**:
```bash
sudo ./setup-cron.sh [--dry-run] [--uninstall]
```

**Features**:
- Configure daily/weekly schedules
- Setup log rotation
- Create systemd timers
- Test first backup

### test-backup-restore.sh
Integration test for backup and restore.

**Usage**:
```bash
./test-backup-restore.sh [--keep-backup] [--skip-cleanup]
```

**Tests**:
- Backup creation
- Backup verification
- Database restore
- Data integrity

## Contact and Escalation

### Backup Issues
1. Check logs: `/var/log/listings/backup.log`
2. Run manual backup: `./backup-db.sh`
3. Contact: DevOps team

### Restore Emergencies
1. **P0 (Critical)**: Production data loss
   - Escalate to: Senior DevOps Engineer
   - Contact: On-call engineer

2. **P1 (High)**: Development/Staging restore needed
   - Contact: DevOps team
   - SLA: 4 hours

## References

- PostgreSQL Backup Documentation: https://www.postgresql.org/docs/current/backup.html
- pg_dump Reference: https://www.postgresql.org/docs/current/app-pgdump.html
- Point-in-Time Recovery: https://www.postgresql.org/docs/current/continuous-archiving.html

## Change History

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2024-11-05 | 1.0 | Initial backup policy | DevOps |

## Appendix

### Environment File Template

```bash
# /etc/listings-backup.env

# Database Configuration
BACKUP_DB_HOST=localhost
BACKUP_DB_PORT=35434
BACKUP_DB_NAME=listings_dev_db
BACKUP_DB_USER=listings_user
BACKUP_DB_PASSWORD=your_secure_password_here
BACKUP_DB_CONTAINER=listings_postgres

# Backup Configuration
BACKUP_DIR=/var/backups/listings
LOG_DIR=/var/log/listings

# Retention Policy
BACKUP_RETENTION_DAYS=7
BACKUP_RETENTION_WEEKS=4
BACKUP_RETENTION_MONTHS=12

# S3/MinIO (Optional)
BACKUP_ENABLE_S3=false
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=your_s3_access_key
S3_SECRET_KEY=your_s3_secret_key
S3_BUCKET=listings-backups
S3_USE_SSL=false
S3_RETENTION_DAYS=30

# Notifications (Optional)
BACKUP_NOTIFY_EMAIL=devops@example.com
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
SMTP_HOST=localhost
SMTP_PORT=25
```

### Quick Reference Card

```
┌─────────────────────────────────────────────────┐
│ Listings Backup Quick Reference                 │
├─────────────────────────────────────────────────┤
│                                                 │
│ BACKUP:                                         │
│   sudo -u listings /path/to/backup-db.sh        │
│                                                 │
│ RESTORE:                                        │
│   sudo -u listings /path/to/restore-db.sh \     │
│       /var/backups/listings/daily/backup.sql.gz │
│                                                 │
│ VERIFY:                                         │
│   sudo -u listings /path/to/verify-backup.sh \  │
│       /var/backups/listings/daily/backup.sql.gz │
│                                                 │
│ MONITOR:                                        │
│   sudo -u listings /path/to/monitor-backups.py  │
│                                                 │
│ LOGS:                                           │
│   tail -f /var/log/listings/backup.log          │
│                                                 │
│ METRICS:                                        │
│   curl http://localhost:9090/metrics            │
│                                                 │
└─────────────────────────────────────────────────┘
```
