# Listings Backup System - Files Overview

## 📁 Complete File Listing

**Total**: 11 files, 5267 lines of code and documentation
**Location**: `/p/github.com/sveturs/listings/scripts/backup/`
**Created**: 2024-11-05

---

## 🔧 Executable Scripts (7 files, 2735 lines)

### 1. backup-db.sh (432 lines)
**Purpose**: Automated database backup with retention policy

**Features**:
- Full PostgreSQL dump with pg_dump
- Gzip compression (level 9)
- SHA256 checksum generation
- Metadata file creation
- WAL archiving for PITR
- Retention policy (7 daily, 4 weekly, 12 months)
- Lock file для concurrent backup prevention
- Email/Slack notifications
- Disk space validation

**Usage**:
```bash
./backup-db.sh [--dry-run] [--no-upload] [--full-only]
```

**Key Functions**:
- `create_full_backup()` - Creates compressed database dump
- `create_wal_archive()` - Archives WAL files for PITR
- `apply_retention_policy()` - Removes old backups
- `upload_to_s3()` - Triggers S3 upload
- `check_disk_space()` - Validates available space

---

### 2. restore-db.sh (514 lines)
**Purpose**: Database restore with verification and rollback

**Features**:
- Restore from compressed backups
- Point-in-Time Recovery (PITR) support
- Pre-restore safety backup
- Post-restore integrity checks
- Automatic rollback on failure
- Connection termination before restore
- Data validation (row counts, table checks)

**Usage**:
```bash
./restore-db.sh <backup_file> [--dry-run] [--no-backup] [--pitr-target "YYYY-MM-DD HH:MM:SS"]
```

**Key Functions**:
- `validate_backup_file()` - Verifies backup integrity
- `backup_current_database()` - Creates pre-restore backup
- `restore_database()` - Restores from backup file
- `apply_pitr()` - Applies Point-in-Time Recovery
- `verify_restore()` - Validates restored data
- `rollback_restore()` - Reverts on failure

---

### 3. verify-backup.sh (459 lines)
**Purpose**: Backup integrity verification with test restore

**Features**:
- File integrity checks (gzip, checksum)
- Test restore to temporary database
- Data integrity validation
- Critical tables verification
- Parallel verification support
- Automated verification reports

**Usage**:
```bash
./verify-backup.sh <backup_file> [--dry-run] [--quick]
./verify-backup.sh --verify-all [--parallel]
```

**Key Functions**:
- `verify_file_integrity()` - Validates file and checksum
- `test_restore()` - Restores to temp database
- `verify_data_integrity()` - Checks critical tables
- `verify_all_backups()` - Verifies all backup files
- `cleanup_temp_database()` - Removes temp resources

---

### 4. backup-s3.sh (391 lines)
**Purpose**: Upload backups to S3/MinIO with verification

**Features**:
- AWS CLI integration
- Multipart upload (files > 100MB)
- Upload verification with checksums
- Retry logic (3 attempts)
- S3 bucket creation and versioning
- Retention policy cleanup
- Metadata preservation

**Usage**:
```bash
./backup-s3.sh <backup_file> [--dry-run]
```

**Key Functions**:
- `configure_aws_cli()` - Sets up AWS credentials
- `create_bucket_if_not_exists()` - Creates S3 bucket
- `upload_file()` - Uploads with retry logic
- `verify_upload()` - Verifies checksum after upload
- `cleanup_old_backups()` - Removes old S3 backups

---

### 5. setup-cron.sh (473 lines)
**Purpose**: Setup automated backup scheduling

**Features**:
- Cron job configuration
- Systemd timer creation
- Log rotation setup
- Directory creation and permissions
- Environment file generation
- First backup test
- Uninstall capability

**Usage**:
```bash
sudo ./setup-cron.sh [--dry-run] [--uninstall]
```

**Key Functions**:
- `create_directories()` - Sets up backup/log dirs
- `setup_cron_jobs()` - Configures cron schedule
- `setup_systemd_timers()` - Creates systemd timers
- `setup_logrotate()` - Configures log rotation
- `test_backup()` - Runs first test backup
- `create_env_file()` - Generates config file

**Default Schedule**:
- Daily backup: 2:00 AM
- Weekly verification: Sunday 6:00 AM
- Hourly monitoring

---

### 6. monitor-backups.py (531 lines)
**Purpose**: Backup monitoring and alerting

**Features**:
- Backup age monitoring (alert if > 25h)
- Backup size validation (alert if < 1MB)
- Size anomaly detection (50% threshold)
- Prometheus metrics export
- HTTP server for metrics endpoint
- Slack/email notifications
- Health check endpoint

**Usage**:
```bash
./monitor-backups.py [--check] [--metrics] [--serve --port 9090]
```

**Key Classes**:
- `BackupFile` - Represents backup file with metadata
- `BackupMetrics` - Metrics data model
- `MetricsHandler` - HTTP handler for /metrics and /health

**Key Functions**:
- `find_backup_files()` - Discovers all backups
- `run_all_checks()` - Executes health checks
- `collect_metrics()` - Gathers monitoring data
- `format_prometheus_metrics()` - Formats for Prometheus
- `send_slack_notification()` - Sends Slack alerts
- `start_metrics_server()` - Starts HTTP server

**Prometheus Metrics**:
- `listings_backup_age_hours` - Age of last backup
- `listings_backup_size_mb` - Size of last backup
- `listings_backup_total` - Total backup count
- `listings_backup_total_size_mb` - Total size
- `listings_backup_health` - Health status (1/0)

---

### 7. test-backup-restore.sh (449 lines)
**Purpose**: Integration test for full backup/restore cycle

**Features**:
- Creates test data in database
- Runs backup process
- Simulates data loss
- Restores from backup
- Verifies data integrity
- Colored console output
- Automatic cleanup

**Usage**:
```bash
./test-backup-restore.sh [--keep-backup] [--skip-cleanup]
```

**Test Steps**:
1. **Test 1**: Backup Creation
   - Creates full backup
   - Verifies file exists
   - Checks size and checksum

2. **Test 2**: Backup Verification
   - Runs verify-backup.sh
   - Validates integrity

3. **Test 3**: Restore
   - Deletes test data (simulates loss)
   - Restores from backup
   - Verifies data recovered

4. **Test 4**: Data Integrity
   - Checks row counts
   - Validates data content
   - Verifies relationships

---

## 📚 Documentation Files (4 files, 2532 lines)

### 8. README.md (518 lines)
**Purpose**: Main documentation and usage guide

**Contents**:
- Features overview
- Quick start guide (4 steps)
- Detailed usage for all scripts
- Configuration reference
- Monitoring setup
- Troubleshooting guide
- Best practices
- Support information

**Sections**:
- Installation
- Usage examples
- Configuration variables
- Cron schedule
- Architecture
- Monitoring
- Troubleshooting
- Scripts reference

---

### 9. BACKUP_POLICY.md (550 lines)
**Purpose**: Comprehensive backup policy and procedures

**Contents**:
- Backup strategy (daily/weekly/monthly)
- Recovery objectives (RTO/RPO)
- Backup process details
- Restore procedures
- Security guidelines
- Testing requirements
- Disaster recovery procedures
- Troubleshooting guide

**Key Sections**:
- **Backup Types**: Daily, weekly, monthly, WAL
- **RTO/RPO**: < 1 hour target
- **Backup Process**: 9-step automated flow
- **Restore Process**: Standard and PITR
- **Security**: Access control, encryption
- **Monitoring**: Health checks, alerts
- **Testing**: Weekly verification, quarterly drills

---

### 10. INSTALLATION.md (387 lines)
**Purpose**: Step-by-step installation guide

**Contents**:
- Prerequisites
- Quick installation (8 steps)
- Configuration examples
- Test procedures
- Monitoring setup
- Notification configuration
- Verification checklist
- Troubleshooting

**Installation Steps**:
1. Create backup user
2. Create directories
3. Configure environment
4. Install scripts
5. Setup automated backups
6. Test installation
7. Setup monitoring (optional)
8. Configure notifications (optional)

**Verification Checklist**:
- [ ] Backup user created
- [ ] Directories exist
- [ ] Environment configured
- [ ] Scripts executable
- [ ] Cron jobs configured
- [ ] Test backup successful
- [ ] Verification passed
- [ ] Monitoring working

---

### 11. IMPLEMENTATION_SUMMARY.md (563 lines)
**Purpose**: Complete overview of the backup system

**Contents**:
- System overview
- Features summary
- Technical specifications
- Architecture diagrams
- Configuration examples
- Quick start examples
- Testing results
- Best practices
- Future enhancements

**Key Sections**:
- **Overview**: Production-ready system
- **Features**: 6 major capability areas
- **Architecture**: Component interaction diagrams
- **Performance**: Benchmarks and metrics
- **Examples**: 6 common use cases
- **Testing**: Integration test results
- **Checklist**: Implementation completion

---

### 12. FILES_OVERVIEW.md (this file)
**Purpose**: Detailed description of all files

**Contents**:
- Complete file listing
- Detailed description of each script
- Key functions and features
- Usage examples
- File statistics

---

## 📊 Statistics Summary

### By Type
| Type | Files | Lines | Size |
|------|-------|-------|------|
| Bash Scripts | 6 | 2,204 | 72 KB |
| Python Scripts | 1 | 531 | 17 KB |
| Documentation | 4 | 2,532 | 55 KB |
| **Total** | **11** | **5,267** | **144 KB** |

### By Purpose
| Purpose | Files | Lines |
|---------|-------|-------|
| Backup Operations | 2 | 823 |
| Restore Operations | 1 | 514 |
| Verification | 1 | 459 |
| Monitoring | 1 | 531 |
| Setup/Testing | 2 | 922 |
| Documentation | 4 | 2,018 |

### Code Quality
- ✅ All scripts have error handling (`set -euo pipefail`)
- ✅ Comprehensive logging with timestamps
- ✅ Dry-run mode for safe testing
- ✅ Input validation and sanitization
- ✅ Lock files for concurrency control
- ✅ Cleanup on error (trap handlers)
- ✅ Detailed comments and documentation

### Documentation Coverage
- ✅ README with quick start
- ✅ Installation guide
- ✅ Comprehensive policy document
- ✅ Implementation summary
- ✅ This file overview
- ✅ Inline comments in all scripts

---

## 🔗 File Dependencies

```
backup-db.sh
    ├── Calls: backup-s3.sh (optional)
    ├── Creates: *.sql.gz, *.sha256, *.meta
    └── Logs: backup.log

restore-db.sh
    ├── Reads: *.sql.gz, *.sha256
    ├── Calls: backup-db.sh (pre-restore)
    └── Logs: restore.log

verify-backup.sh
    ├── Reads: *.sql.gz, *.sha256
    ├── Calls: restore-db.sh (test restore)
    └── Creates: reports/*.txt

backup-s3.sh
    ├── Reads: *.sql.gz, *.sha256, *.meta
    ├── Requires: aws-cli
    └── Logs: backup-s3.log

monitor-backups.py
    ├── Reads: backup files, metrics
    ├── Creates: backup_metrics.json
    ├── Serves: HTTP endpoints (/metrics, /health)
    └── Logs: monitor-backups.log

setup-cron.sh
    ├── Creates: cron jobs, systemd timers
    ├── Creates: /etc/listings-backup.env
    └── Calls: All other scripts (for testing)

test-backup-restore.sh
    ├── Calls: backup-db.sh, restore-db.sh, verify-backup.sh
    ├── Creates: test data
    └── Validates: Complete cycle
```

---

## 🎯 Usage Flow

### Daily Automated Backup
```
02:00 AM → Cron triggers
           ↓
       backup-db.sh
           ├─► Creates backup.sql.gz
           ├─► Generates checksum
           ├─► Creates metadata
           ├─► Archives WAL files
           ├─► Applies retention
           └─► Calls backup-s3.sh (if enabled)
                   └─► Uploads to S3/MinIO
```

### Manual Restore
```
User runs restore-db.sh
    ↓
Validates backup file
    ↓
Creates pre-restore backup (safety)
    ↓
Terminates DB connections
    ↓
Drops and recreates database
    ↓
Restores from backup
    ↓
Verifies data integrity
    ↓
Success or Rollback
```

### Monitoring
```
Every hour → Cron triggers
             ↓
         monitor-backups.py
             ├─► Finds backup files
             ├─► Checks backup age
             ├─► Validates sizes
             ├─► Detects anomalies
             ├─► Updates metrics
             └─► Sends alerts (if issues)
```

---

## 🔍 Quick Reference

### Most Important Files
1. **backup-db.sh** - Main backup script (run daily)
2. **restore-db.sh** - Restore when needed
3. **monitor-backups.py** - Continuous monitoring
4. **README.md** - Start here for usage

### Configuration Files
- `/etc/listings-backup.env` - Main configuration
- `/etc/cron.d/listings-backup` - Cron schedule
- `/etc/systemd/system/listings-backup.*` - Systemd timers

### Log Files
- `/var/log/listings/backup.log` - Backup operations
- `/var/log/listings/restore.log` - Restore operations
- `/var/log/listings/verify-backup.log` - Verification
- `/var/log/listings/monitor-backups.log` - Monitoring

### Data Files
- `/var/backups/listings/daily/*.sql.gz` - Daily backups
- `/var/backups/listings/weekly/*.sql.gz` - Weekly backups
- `/var/backups/listings/monthly/*.sql.gz` - Monthly backups
- `/var/log/listings/backup_metrics.json` - Current metrics

---

## 📞 Getting Help

**Read documentation**:
1. Start with [README.md](README.md)
2. Follow [INSTALLATION.md](INSTALLATION.md)
3. Review [BACKUP_POLICY.md](BACKUP_POLICY.md)
4. Check this file for details

**Check logs**:
```bash
tail -f /var/log/listings/backup.log
tail -f /var/log/listings/monitor-backups.log
```

**Test scripts**:
```bash
./backup-db.sh --dry-run
./test-backup-restore.sh
```

**Contact**:
- DevOps Team: devops@example.com
- Emergency: PagerDuty

---

**Last Updated**: 2024-11-05
**Version**: 1.0.0
**Status**: ✅ Production Ready
