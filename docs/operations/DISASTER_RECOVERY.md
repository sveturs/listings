# Listings Microservice Disaster Recovery Plan

**Last Updated:** 2025-11-05
**Version:** 1.0.0
**Maintainer:** Platform Team
**Classification:** CONFIDENTIAL

## Table of Contents

- [Overview](#overview)
- [Recovery Objectives](#recovery-objectives)
- [Disaster Scenarios](#disaster-scenarios)
  - [1. Complete Database Loss](#1-complete-database-loss)
  - [2. Region/Datacenter Failure](#2-regiondatacenter-failure)
  - [3. Data Corruption](#3-data-corruption)
  - [4. Security Breach](#4-security-breach)
  - [5. Complete Service Loss](#5-complete-service-loss)
- [Backup Strategy](#backup-strategy)
- [Communication Protocol](#communication-protocol)
- [Post-Recovery Procedures](#post-recovery-procedures)

---

## Overview

This document outlines disaster recovery procedures for the Listings Microservice. Follow these procedures in the event of catastrophic failures that cannot be resolved through standard operational procedures.

### When to Use This Document

Execute disaster recovery procedures when:
- Standard recovery procedures (RUNBOOK.md) have failed
- Data loss or corruption is suspected
- Multiple system components have failed simultaneously
- Security breach requiring immediate containment
- Service downtime exceeds 30 minutes

### Incident Commander

The **Platform Team Lead** serves as Incident Commander for all disaster recovery scenarios. Contact via PagerDuty: `platform-team-lead`.

---

## Recovery Objectives

### RTO (Recovery Time Objective)

| Scenario | RTO Target | Maximum Acceptable |
|----------|------------|-------------------|
| Database Loss | 30 minutes | 1 hour |
| Region Failure | 15 minutes | 30 minutes |
| Data Corruption | 1 hour | 2 hours |
| Security Breach | Immediate | 15 minutes |
| Complete Service Loss | 15 minutes | 30 minutes |

### RPO (Recovery Point Objective)

| Data Type | RPO Target | Backup Frequency |
|-----------|------------|------------------|
| Database | 5 minutes | Continuous WAL archiving |
| OpenSearch Index | 15 minutes | Reindexable from DB |
| Redis Cache | Acceptable loss | Not backed up (ephemeral) |
| Configuration | 0 (version controlled) | Git repository |

### SLO Impact Thresholds

- **Critical:** Downtime > 30 minutes (breaches monthly SLO)
- **Major:** Downtime 15-30 minutes (warning threshold)
- **Minor:** Downtime < 15 minutes (within SLO budget)

---

## Disaster Scenarios

### 1. Complete Database Loss

**Scenario:** PostgreSQL database completely unavailable or corrupted beyond repair.

**Impact:**
- ✗ All listing CRUD operations fail
- ✗ Service completely unavailable
- ✗ Search unavailable (OpenSearch can't sync)
- ✓ Read-only mode possible from OpenSearch (future enhancement)

#### Detection

**Symptoms:**
```bash
# Database unreachable
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"
# psql: error: connection to server at "localhost", port 35433 failed

# All database operations failing
curl -s http://localhost:8086/metrics | grep listings_errors_total
# Spike in database errors

# Service logs showing database errors
sudo journalctl -u listings-service | grep -i "database connection"
```

**Alerts:**
- `ListingsServiceDown`
- `ListingsDBUnreachable`
- `ListingsHighErrorRate`

#### Recovery Procedure

**Step 1: Assess Damage (5 minutes)**

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check PostgreSQL logs
sudo tail -100 /var/log/postgresql/postgresql-15-main.log

# Attempt connection
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "\l"

# Check disk space
df -h /var/lib/postgresql

# Check for corruption indicators
sudo -u postgres pg_controldata /var/lib/postgresql/15/main
```

**Step 2: Stop Listings Service (1 minute)**

```bash
# Prevent writes to corrupted database
sudo systemctl stop listings-service

# Verify stopped
ps aux | grep listings-service
```

**Step 3: Attempt PostgreSQL Restart (2 minutes)**

```bash
# Try restart
sudo systemctl restart postgresql

# Check status
sudo systemctl status postgresql

# If successful, verify data integrity
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT COUNT(*) FROM listings;"

# If data intact, proceed to Step 8
# If data missing/corrupted, continue to Step 4
```

**Step 4: Locate Latest Backup (2 minutes)**

```bash
# List available backups (example location)
ls -lh /opt/backups/postgresql/listings_db/

# Check backup age
ls -lt /opt/backups/postgresql/listings_db/ | head -5

# Verify backup integrity
pg_restore -l /opt/backups/postgresql/listings_db/backup_latest.dump | head -20
```

**Example backup structure:**
```
/opt/backups/postgresql/listings_db/
├── backup_2025-11-05_00-00.dump
├── backup_2025-11-05_06-00.dump
├── backup_2025-11-05_12-00.dump
├── backup_latest.dump -> backup_2025-11-05_12-00.dump
└── wal_archive/
    ├── 000000010000000000000001
    ├── 000000010000000000000002
    └── ...
```

**Step 5: Stop PostgreSQL (1 minute)**

```bash
# Stop PostgreSQL
sudo systemctl stop postgresql

# Verify stopped
ps aux | grep postgres
```

**Step 6: Restore Database (10-15 minutes)**

```bash
# Backup current (corrupted) data directory
sudo mv /var/lib/postgresql/15/main /var/lib/postgresql/15/main.corrupted

# Create fresh data directory
sudo -u postgres mkdir /var/lib/postgresql/15/main
sudo -u postgres chmod 700 /var/lib/postgresql/15/main

# Initialize new cluster
sudo -u postgres /usr/lib/postgresql/15/bin/initdb -D /var/lib/postgresql/15/main

# Start PostgreSQL
sudo systemctl start postgresql

# Create database and user
sudo -u postgres psql -c "CREATE DATABASE listings_db;"
sudo -u postgres psql -c "CREATE USER listings_user WITH PASSWORD 'listings_password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE listings_db TO listings_user;"

# Restore from backup
pg_restore -U postgres -d listings_db /opt/backups/postgresql/listings_db/backup_latest.dump

# Apply WAL archives (if available)
# This recovers transactions since last backup
sudo -u postgres pg_waldump /opt/backups/postgresql/listings_db/wal_archive/000000010000000000000001
# ... (advanced procedure, consult Database SRE)
```

**Step 7: Verify Database Integrity (5 minutes)**

```bash
# Check table counts
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT 'listings', COUNT(*) FROM listings
   UNION ALL
   SELECT 'listing_images', COUNT(*) FROM listing_images
   UNION ALL
   SELECT 'listing_attributes', COUNT(*) FROM listing_attributes;"

# Check data consistency
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT id, title, created_at FROM listings ORDER BY id LIMIT 5;"

# Verify indexes
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT schemaname, tablename, indexname FROM pg_indexes
   WHERE schemaname = 'public' ORDER BY tablename;"

# Run ANALYZE
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "ANALYZE;"
```

**Step 8: Start Listings Service (2 minutes)**

```bash
# Start service
sudo systemctl start listings-service

# Monitor startup
sudo journalctl -u listings-service -f

# Wait for healthy status
sleep 10

# Verify health
curl http://localhost:8086/health
curl http://localhost:8086/ready
```

**Step 9: Verify Service Functionality (3 minutes)**

```bash
# Test gRPC endpoint
grpcurl -plaintext -d '{"id": 1}' localhost:50053 listings.v1.ListingsService/GetListing

# Test HTTP endpoint
curl http://localhost:8086/api/v1/listings?limit=5

# Check metrics
curl -s http://localhost:8086/metrics | grep -E 'listings_grpc_requests_total|listings_errors_total'

# Monitor for 2 minutes
watch -n 5 'curl -s http://localhost:8086/metrics | grep listings_errors_total'
```

**Step 10: Reindex OpenSearch (5 minutes)**

```bash
# Trigger reindexing
cd /p/github.com/sveturs/listings
python3 scripts/reindex_via_docker.py --target-password admin

# Monitor progress
curl -u admin:admin http://localhost:9200/listings_microservice/_count

# Validate
python3 scripts/validate_opensearch.py --target-password admin
```

#### Validation Checklist

- [ ] PostgreSQL running and accepting connections
- [ ] Database contains expected number of records
- [ ] Listings service health check passing
- [ ] No database errors in service logs
- [ ] gRPC endpoints responding
- [ ] HTTP endpoints responding
- [ ] Error rate < 1%
- [ ] OpenSearch reindexed successfully
- [ ] Search queries returning results
- [ ] Metrics being collected

#### Rollback Plan

If restoration fails:
```bash
# Restore original corrupted database for forensics
sudo systemctl stop postgresql
sudo rm -rf /var/lib/postgresql/15/main
sudo mv /var/lib/postgresql/15/main.corrupted /var/lib/postgresql/15/main
sudo systemctl start postgresql

# Escalate to Database SRE team immediately
```

#### Communication Template

```
DISASTER RECOVERY: Complete Database Loss - Listings Service

STATUS: RESTORING
START TIME: 2025-11-05 14:30 UTC
INCIDENT COMMANDER: [Name]

IMPACT:
- Listings service completely unavailable
- All CRUD operations failing
- Estimated users affected: ALL
- Estimated recovery time: 30 minutes

ACTIONS TAKEN:
- 14:30 - Database failure detected
- 14:32 - Listings service stopped
- 14:35 - Backup located and validated
- 14:40 - Database restoration started
- 14:55 - Database restored, service starting

CURRENT STATUS:
- Database restored from backup at 12:00 UTC (2.5 hours RPO)
- Service health checks passing
- Functionality tests in progress

NEXT STEPS:
- Complete validation
- Monitor for 30 minutes
- Schedule postmortem

ETA FOR FULL RECOVERY: 15:05 UTC
```

---

### 2. Region/Datacenter Failure

**Scenario:** Complete loss of primary datacenter/region where Listings service is deployed.

**Impact:**
- ✗ All services in region unavailable
- ✗ Network connectivity lost
- ✓ DR site can take over (if configured)

**Note:** Current deployment is single-region. Multi-region DR is a future enhancement.

#### Current Mitigation

```bash
# Single-region deployment
# Primary: /opt/listings-dev on svetu.rs server

# Future: Multi-region architecture
# Primary: us-east-1
# DR: eu-west-1
# RPO: < 1 minute (continuous replication)
# RTO: < 15 minutes (automated failover)
```

#### Manual Failover (Future)

**Prerequisites:**
- Secondary region configured
- Database replication enabled (streaming replication or logical replication)
- DNS failover configured (Route53 health checks)
- Container registry accessible from both regions

**Procedure (When Implemented):**
```bash
# 1. Verify primary region failure
# 2. Promote secondary database to primary
# 3. Update DNS to point to secondary region
# 4. Start services in secondary region
# 5. Verify functionality
```

#### Current Workaround

If hosting provider has total outage:
1. Provision new VM at different provider
2. Restore database from backup
3. Deploy service from git repository
4. Update DNS

Estimated recovery time: **2-4 hours** (manual process)

---

### 3. Data Corruption

**Scenario:** Listings data corrupted due to software bug, operator error, or partial failure.

**Impact:**
- Some listings showing incorrect data
- Relationships broken (images, attributes)
- Possible data loss

#### Detection

**Symptoms:**
```bash
# Data inconsistencies
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT COUNT(*) FROM listings WHERE deleted_at IS NOT NULL;"
# Unexpected number of deleted records

# Orphaned records
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT COUNT(*) FROM listing_images li
   LEFT JOIN listings l ON li.listing_id = l.id
   WHERE l.id IS NULL;"
# Should be 0

# Invalid data
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT COUNT(*) FROM listings WHERE price < 0 OR price IS NULL;"
# Should be 0 (or match business rules)
```

**User Reports:**
- "Listing images missing"
- "Wrong price displayed"
- "Listing details incorrect"

#### Recovery Procedure

**Step 1: Identify Scope of Corruption (10 minutes)**

```bash
# Determine affected time range
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT MIN(updated_at), MAX(updated_at)
   FROM listings
   WHERE /* corrupted condition */;"

# Count affected records
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT COUNT(*) FROM listings WHERE /* corrupted condition */;"

# Identify affected users
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT DISTINCT user_id FROM listings WHERE /* corrupted condition */;"
```

**Step 2: Stop Further Corruption (2 minutes)**

```bash
# If caused by buggy code, stop service immediately
sudo systemctl stop listings-service

# If caused by bad data import, identify source
sudo journalctl -u listings-service --since "1 hour ago" | grep -i import
```

**Step 3: Point-in-Time Recovery (PITR) (30 minutes)**

**Option A: Restore specific tables from backup**
```bash
# Extract affected table from backup
pg_restore -t listings -U postgres /opt/backups/postgresql/listings_db/backup_2025-11-05_12-00.dump > /tmp/listings_restore.sql

# Review SQL before applying
head -100 /tmp/listings_restore.sql

# Apply to temporary table
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "CREATE TABLE listings_backup AS SELECT * FROM listings WHERE id IN (/* affected IDs */);"

# Restore good data
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" < /tmp/listings_restore.sql

# Verify restoration
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT COUNT(*) FROM listings WHERE /* corruption check */;"
```

**Option B: Restore entire database to specific point**
```bash
# Create new database for recovery
sudo -u postgres psql -c "CREATE DATABASE listings_db_recovery;"

# Restore backup
pg_restore -U postgres -d listings_db_recovery /opt/backups/postgresql/listings_db/backup_2025-11-05_00-00.dump

# Apply WAL logs up to specific time
sudo -u postgres pg_waldump /opt/backups/postgresql/listings_db/wal_archive/* | \
  sudo -u postgres pg_replay -d listings_db_recovery

# Extract good data
psql "postgres://listings_user:listings_password@localhost:35433/listings_db_recovery" -c \
  "COPY (SELECT * FROM listings WHERE id IN (/* affected IDs */)) TO '/tmp/good_data.csv' WITH CSV HEADER;"

# Import to production
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "COPY listings FROM '/tmp/good_data.csv' WITH CSV HEADER ON CONFLICT (id) DO UPDATE SET ...;"
```

**Step 4: Validate Data Integrity (10 minutes)**

```bash
# Run data integrity checks
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -f /p/github.com/sveturs/listings/scripts/validate_data_integrity.sql

# Verify specific records
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT id, title, price, status FROM listings WHERE id IN (/* sample IDs */);"

# Check relationships
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT l.id, l.title, COUNT(li.id) AS image_count
   FROM listings l
   LEFT JOIN listing_images li ON li.listing_id = l.id
   WHERE l.id IN (/* affected IDs */)
   GROUP BY l.id, l.title;"
```

**Step 5: Restart Service and Monitor (10 minutes)**

```bash
# Start service
sudo systemctl start listings-service

# Monitor for errors
sudo journalctl -u listings-service -f

# Test affected endpoints
for id in 1 2 3 4 5; do
  grpcurl -plaintext -d "{\"id\": $id}" localhost:50053 listings.v1.ListingsService/GetListing
done

# Reindex OpenSearch
python3 /p/github.com/sveturs/listings/scripts/reindex_via_docker.py --target-password admin
```

#### Validation Checklist

- [ ] Corrupted data identified and isolated
- [ ] Good data restored from backup
- [ ] Data integrity checks passing
- [ ] Service restarted successfully
- [ ] Affected endpoints tested
- [ ] OpenSearch reindexed
- [ ] User-reported issues verified fixed

#### Prevention

```bash
# Add data validation constraints
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "ALTER TABLE listings ADD CONSTRAINT check_price_positive CHECK (price >= 0);"

# Add foreign key constraints (with cascades)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "ALTER TABLE listing_images
   ADD CONSTRAINT fk_listing_images_listing_id
   FOREIGN KEY (listing_id) REFERENCES listings(id) ON DELETE CASCADE;"

# Regular data integrity audits (weekly cron)
0 0 * * 0 psql "postgres://listings_user:listings_password@localhost:35433/listings_db" \
  -f /p/github.com/sveturs/listings/scripts/weekly_data_audit.sql | \
  mail -s "Weekly Data Integrity Report" platform@svetu.rs
```

---

### 4. Security Breach

**Scenario:** Unauthorized access detected, credentials compromised, or malicious activity.

**Impact:**
- Data confidentiality compromised
- Potential data loss or corruption
- Service integrity at risk
- Legal and compliance implications

**CRITICAL:** Execute immediately upon confirmation of breach.

#### Detection

**Indicators:**
- Unusual database queries in logs
- Unexpected authentication attempts
- Unauthorized data access
- Modified configuration files
- Unusual network traffic
- Alerts from security monitoring tools

#### Immediate Response (5 minutes)

**Step 1: Contain Breach**

```bash
# STOP THE SERVICE IMMEDIATELY
sudo systemctl stop listings-service

# Block all external access at firewall
sudo iptables -A INPUT -p tcp --dport 50053 -j DROP
sudo iptables -A INPUT -p tcp --dport 8086 -j DROP
sudo iptables -A INPUT -p tcp --dport 9093 -j DROP

# Verify blocked
netstat -tlnp | grep -E '50053|8086|9093'

# If behind load balancer/nginx
sudo systemctl stop nginx
```

**Step 2: Notify Security Team**

```bash
# Page security team immediately
# PagerDuty: security-team
# Email: security@svetu.rs
# Slack: #security-incidents

# Do not discuss publicly
```

**Step 3: Preserve Evidence**

```bash
# Capture current state
mkdir -p /tmp/security-incident-$(date +%Y%m%d-%H%M%S)
cd /tmp/security-incident-*

# Copy logs
sudo cp /var/log/journal/listings-service.* .
sudo journalctl -u listings-service --since "24 hours ago" > service-logs.txt

# Copy configuration
sudo cp -r /opt/listings-dev/.env .
sudo cp -r /opt/listings-dev/.env.backup .

# Network connections
netstat -an > network-connections.txt
sudo ss -tulpn > socket-stats.txt

# Process information
ps auxf > processes.txt

# Database connections
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT * FROM pg_stat_activity;" > db-connections.txt

# File integrity
find /opt/listings-dev -type f -exec sha256sum {} \; > file-checksums.txt

# Redis keys
redis-cli -h localhost -p 36380 -a redis_password KEYS '*' > redis-keys.txt

# Compress evidence
tar czf /tmp/security-incident-$(date +%Y%m%d-%H%M%S).tar.gz .
```

#### Investigation (Security Team)

**Step 4: Identify Attack Vector**

```bash
# Analyze access logs
sudo journalctl -u listings-service --since "7 days ago" | \
  jq -r 'select(.level == "warn" or .level == "error")'

# Check for suspicious queries
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT * FROM pg_stat_statements ORDER BY total_exec_time DESC LIMIT 100;"

# Check authentication failures
sudo journalctl | grep -i "authentication failed"

# Check for modified files
find /opt/listings-dev -type f -mtime -1 -ls
```

**Step 5: Assess Data Exposure**

```bash
# Check data access patterns
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT schemaname, tablename, n_tup_ins, n_tup_upd, n_tup_del
   FROM pg_stat_user_tables
   WHERE schemaname = 'public'
   ORDER BY n_tup_upd + n_tup_del DESC;"

# Identify accessed records (if audit log enabled)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT * FROM audit_log WHERE created_at > NOW() - INTERVAL '24 hours'
   ORDER BY created_at DESC;"

# Check for data exfiltration
sudo journalctl -u listings-service --since "24 hours ago" | \
  jq -r 'select(.response_size > 10000000)' # Responses > 10MB
```

#### Recovery (30-60 minutes)

**Step 6: Rotate All Credentials**

```bash
# Generate new database password
NEW_DB_PASSWORD=$(openssl rand -base64 32)

# Update database password
psql "postgres://postgres:current_password@localhost:35433/postgres" -c \
  "ALTER USER listings_user WITH PASSWORD '$NEW_DB_PASSWORD';"

# Update .env file
sudo sed -i "s/SVETULISTINGS_DB_PASSWORD=.*/SVETULISTINGS_DB_PASSWORD=$NEW_DB_PASSWORD/" \
  /opt/listings-dev/.env

# Generate new Redis password
NEW_REDIS_PASSWORD=$(openssl rand -base64 32)

# Update Redis password
redis-cli -h localhost -p 36380 -a old_password CONFIG SET requirepass "$NEW_REDIS_PASSWORD"

# Update .env file
sudo sed -i "s/SVETULISTINGS_REDIS_PASSWORD=.*/SVETULISTINGS_REDIS_PASSWORD=$NEW_REDIS_PASSWORD/" \
  /opt/listings-dev/.env

# Rotate OpenSearch credentials (contact OpenSearch admin)

# Rotate TLS certificates (if applicable)
# Contact infrastructure team
```

**Step 7: Patch Vulnerability**

```bash
# If caused by code vulnerability, deploy hotfix
cd /p/github.com/sveturs/listings
git pull origin security-hotfix

# Rebuild
make build

# Deploy new binary
sudo cp bin/listings-service /opt/listings-dev/bin/

# Update configuration with security hardening
sudo nano /opt/listings-dev/.env
# Add:
# SVETULISTINGS_SECURITY_ENHANCED=true
# SVETULISTINGS_AUDIT_LOGGING=true
```

**Step 8: Restore from Clean Backup (if compromised)**

```bash
# If data integrity compromised, restore from pre-breach backup
# Follow "Complete Database Loss" procedure
# Use backup from before breach was detected

# Identify breach timestamp
BREACH_TIME="2025-11-05 14:00:00"

# Find backup before breach
ls -lt /opt/backups/postgresql/listings_db/ | \
  awk -v breach="$BREACH_TIME" '{if ($6" "$7" "$8 < breach) print $0}' | \
  head -1

# Restore that backup
# (follow Step 6 from "Complete Database Loss")
```

**Step 9: Validate Security Posture**

```bash
# Run security audit
cd /p/github.com/sveturs/listings
make security-audit

# Check for backdoors
sudo rkhunter --check
sudo chkrootkit

# Verify file integrity
sudo aide --check

# Check open ports
sudo nmap -sV localhost

# Verify firewall rules
sudo iptables -L -n -v

# Check cron jobs for persistence
sudo crontab -l
sudo cat /etc/crontab
```

**Step 10: Controlled Service Restart**

```bash
# Remove iptables blocks
sudo iptables -D INPUT -p tcp --dport 50053 -j DROP
sudo iptables -D INPUT -p tcp --dport 8086 -j DROP
sudo iptables -D INPUT -p tcp --dport 9093 -j DROP

# Start service with enhanced logging
sudo systemctl start listings-service

# Monitor closely for 1 hour
sudo journalctl -u listings-service -f

# Enable additional security features
# - Rate limiting (stricter)
# - IP allowlisting (if applicable)
# - WAF rules (if using WAF)
# - Enhanced audit logging
```

#### Validation Checklist

- [ ] Service stopped and isolated immediately
- [ ] Security team notified
- [ ] Evidence preserved
- [ ] Attack vector identified
- [ ] Data exposure assessed
- [ ] All credentials rotated
- [ ] Vulnerability patched
- [ ] System restored from clean backup (if needed)
- [ ] Security posture validated
- [ ] Service restarted with monitoring

#### Post-Breach Actions

```bash
# 1. Full security audit
# 2. Penetration testing
# 3. Code review for vulnerabilities
# 4. Update security policies
# 5. Staff security training
# 6. Customer notification (if PII exposed)
# 7. Legal compliance (GDPR, etc.)
# 8. Insurance claim (if applicable)
```

#### Communication Template

```
SECURITY INCIDENT: Listings Service

CLASSIFICATION: CONFIDENTIAL
SEVERITY: HIGH
STATUS: CONTAINED

TIMELINE:
- 14:30 - Breach detected
- 14:31 - Service stopped and isolated
- 14:32 - Security team notified
- 14:35 - Evidence collection started
- 15:00 - Investigation in progress

IMPACT ASSESSMENT:
- Service unavailable (planned)
- Data exposure: Under investigation
- Credentials: Being rotated
- Affected users: TBD

ACTIONS TAKEN:
- Service stopped and network isolated
- Evidence preserved
- Security team investigating
- Credentials rotation in progress

NEXT STEPS:
- Complete investigation
- Patch vulnerability
- Restore service with enhanced security
- Customer notification (if required)

ETA FOR SERVICE RESTORATION: 16:30 (controlled)

DO NOT SHARE OUTSIDE SECURITY TEAM
```

---

### 5. Complete Service Loss

**Scenario:** Service binary deleted, VM destroyed, or complete infrastructure failure.

**Impact:**
- Service completely unavailable
- All operations failing
- No access to running instance

#### Recovery Procedure (15-30 minutes)

**Step 1: Verify Infrastructure (2 minutes)**

```bash
# Check VM status
# If cloud provider: Check console/dashboard
# If bare metal: Physical access

# Check SSH access
ssh svetu@svetu.rs

# If accessible, check service
sudo systemctl status listings-service
ps aux | grep listings-service
```

**Step 2: Verify Dependencies (3 minutes)**

```bash
# PostgreSQL
sudo systemctl status postgresql
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"

# Redis
sudo systemctl status redis
redis-cli -h localhost -p 36380 -a redis_password ping

# OpenSearch
curl -u admin:admin http://localhost:9200/_cluster/health
```

**Step 3: Deploy from Source (10 minutes)**

```bash
# Clone repository (if not present)
cd /opt
sudo git clone https://github.com/sveturs/listings.git listings-dev
cd listings-dev

# Checkout correct branch
sudo git fetch origin
sudo git checkout main  # or specific release tag

# Copy environment file
sudo cp /opt/backups/config/listings/.env .
# OR restore from version control

# Install dependencies
make deps

# Build service
make build

# Verify binary
ls -lh bin/listings-service
```

**Step 4: Run Database Migrations (3 minutes)**

```bash
# Check current migration version
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT * FROM schema_migrations;"

# Apply migrations if needed
make migrate-up

# Verify
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 5;"
```

**Step 5: Configure Systemd Service (2 minutes)**

```bash
# Copy service file
sudo cp deployment/listings-service.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable listings-service
```

**Step 6: Start Service (2 minutes)**

```bash
# Start service
sudo systemctl start listings-service

# Check status
sudo systemctl status listings-service

# Monitor logs
sudo journalctl -u listings-service -f
```

**Step 7: Verify Functionality (5 minutes)**

```bash
# Health checks
curl http://localhost:8086/health
curl http://localhost:8086/ready

# gRPC test
grpcurl -plaintext localhost:50053 grpc.health.v1.Health/Check

# API test
curl http://localhost:8086/api/v1/listings?limit=5

# Metrics
curl -s http://localhost:8086/metrics | grep listings_grpc_requests_total

# Monitor for 5 minutes
watch -n 5 'curl -s http://localhost:8086/metrics | grep -E "listings_grpc_requests_total|listings_errors_total"'
```

**Step 8: Restore Monitoring (3 minutes)**

```bash
# Verify Prometheus scraping
curl http://prometheus-server:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.job == "listings")'

# Verify Grafana dashboards
# Access: http://grafana-server:3000
# Dashboard: "Listings Service Overview"

# Check alerts
curl http://alertmanager-server:9093/api/v2/alerts
```

#### Validation Checklist

- [ ] Infrastructure accessible
- [ ] Dependencies (DB, Redis, OpenSearch) running
- [ ] Service deployed from source
- [ ] Database migrations current
- [ ] Service running and healthy
- [ ] gRPC endpoints responding
- [ ] HTTP endpoints responding
- [ ] Error rate < 1%
- [ ] Monitoring restored
- [ ] Alerts configured

---

## Backup Strategy

### Automated Backups

#### PostgreSQL Backups

**Schedule:**
```bash
# Daily full backup (midnight UTC)
0 0 * * * /opt/scripts/backup_listings_db.sh

# Continuous WAL archiving (real-time)
# Configured in postgresql.conf:
# archive_mode = on
# archive_command = 'cp %p /opt/backups/postgresql/listings_db/wal_archive/%f'
```

**Backup Script Example:**
```bash
#!/bin/bash
# /opt/scripts/backup_listings_db.sh

BACKUP_DIR="/opt/backups/postgresql/listings_db"
TIMESTAMP=$(date +%Y-%m-%d_%H-%M)
BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.dump"

# Create backup
pg_dump -U postgres -d listings_db -F c -f "$BACKUP_FILE"

# Verify backup
pg_restore -l "$BACKUP_FILE" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Backup successful: $BACKUP_FILE"

    # Create latest symlink
    ln -sf "$BACKUP_FILE" "$BACKUP_DIR/backup_latest.dump"

    # Compress old backups (>7 days)
    find "$BACKUP_DIR" -name "backup_*.dump" -mtime +7 -exec gzip {} \;

    # Delete very old backups (>30 days)
    find "$BACKUP_DIR" -name "backup_*.dump.gz" -mtime +30 -delete
else
    echo "Backup verification failed!"
    exit 1
fi
```

**Retention Policy:**
- **Daily backups:** 30 days
- **WAL archives:** 7 days
- **Monthly backups:** 12 months (first of month)

#### Configuration Backups

```bash
# Backup configuration (daily)
0 1 * * * tar czf /opt/backups/config/listings/config_$(date +\%Y\%m\%d).tar.gz \
  /opt/listings-dev/.env \
  /opt/listings-dev/deployment/ \
  /etc/systemd/system/listings-service.service

# Retention: 90 days
```

#### Backup Verification

```bash
# Weekly backup verification (Sunday 2 AM)
0 2 * * 0 /opt/scripts/verify_listings_backup.sh

# Script tests:
# 1. Backup file integrity
# 2. Restore to test database
# 3. Data consistency checks
# 4. Report generation
```

### Backup Locations

| Backup Type | Primary Location | Secondary Location | Offsite Location |
|-------------|------------------|-------------------|------------------|
| Database Dumps | `/opt/backups/postgresql/listings_db/` | NFS mount: `/mnt/backup-nas/` | S3: `s3://svetu-backups/listings/` |
| WAL Archives | `/opt/backups/postgresql/listings_db/wal_archive/` | NFS mount: `/mnt/backup-nas/wal/` | S3: `s3://svetu-backups/listings-wal/` |
| Configuration | `/opt/backups/config/listings/` | Git: `github.com/sveturs/listings` | S3: `s3://svetu-backups/config/` |

### Restore Testing

**Quarterly Restore Drills:**
```bash
# Q1, Q2, Q3, Q4 - First Monday of quarter
# Procedure:
# 1. Create test VM
# 2. Restore latest backup
# 3. Verify data integrity
# 4. Document any issues
# 5. Update procedures if needed
```

---

## Communication Protocol

### Severity Levels

| Level | Criteria | Response | Notification |
|-------|----------|----------|--------------|
| **P1 - Critical** | Complete service loss, data breach | Immediate | All stakeholders |
| **P2 - High** | Major functionality impaired | Within 15 min | Engineering teams |
| **P3 - Medium** | Degraded performance | Within 1 hour | On-call engineer |
| **P4 - Low** | Minor issues, no user impact | Next business day | Internal only |

### Notification Channels

**P1 (Critical):**
- PagerDuty: `listings-oncall`, `platform-team-lead`, `security-team`
- Slack: `#listings-incidents` (public), `#platform-team` (private)
- Email: `engineering-all@svetu.rs`, `exec-team@svetu.rs`
- Status Page: Update immediately

**P2 (High):**
- PagerDuty: `listings-oncall`, `platform-team-lead`
- Slack: `#listings-incidents`, `#platform-team`
- Email: `engineering@svetu.rs`

**P3 (Medium):**
- PagerDuty: `listings-oncall`
- Slack: `#listings-alerts`

**P4 (Low):**
- Jira ticket
- Slack: `#listings-team`

### Status Updates

**Frequency:**
- **P1:** Every 30 minutes until resolved
- **P2:** Every hour
- **P3:** Every 4 hours
- **P4:** Daily standup

**Update Template:**
```
INCIDENT UPDATE #[N] - Listings [ISSUE]

STATUS: [INVESTIGATING|IDENTIFIED|RESOLVING|RESOLVED]
TIME: [TIMESTAMP]
INCIDENT COMMANDER: [NAME]

SUMMARY:
[Brief description of current situation]

ACTIONS TAKEN SINCE LAST UPDATE:
- [Action 1]
- [Action 2]

CURRENT IMPACT:
- Service availability: [%]
- Affected users: [COUNT or ALL]
- Functionality: [LIST]

NEXT STEPS:
- [Step 1] (ETA: [TIME])
- [Step 2]

ETA FOR RESOLUTION: [TIME]

NEXT UPDATE: [TIME]
```

### Stakeholder Matrix

| Role | P1 | P2 | P3 | P4 |
|------|----|----|----|----|
| CEO | Notify | Inform | - | - |
| CTO | Notify | Notify | Inform | - |
| VP Engineering | Notify | Notify | Inform | - |
| Platform Team Lead | Notify | Notify | Notify | Inform |
| On-Call Engineer | Notify | Notify | Notify | Assign |
| Database SRE | Notify (if DB) | Notify (if DB) | Inform | - |
| Security Team | Notify (if security) | Notify (if security) | - | - |
| Customer Support | Notify | Inform | - | - |
| Marketing | Inform | - | - | - |

---

## Post-Recovery Procedures

### Immediate (0-4 hours)

**1. Stabilization Monitoring**
```bash
# Monitor service metrics closely
watch -n 10 'curl -s http://localhost:8086/metrics | grep -E "listings_errors_total|listings_grpc_requests_total"'

# Check for anomalies
sudo journalctl -u listings-service -f --since "30 minutes ago"

# Verify all functionality
bash /p/github.com/sveturs/listings/scripts/smoke_test.sh
```

**2. Update Status Page**
```
RESOLVED: Listings Service Issue

The Listings service has been fully restored and is operating normally.

Incident Summary:
- Start: [TIME]
- End: [TIME]
- Duration: [HH:MM]
- Root Cause: [BRIEF DESCRIPTION]

Impact:
- [DESCRIPTION OF IMPACT]

Resolution:
- [BRIEF DESCRIPTION OF FIX]

Next Steps:
- We will conduct a full postmortem within 48 hours
- Additional monitoring has been implemented

We apologize for any inconvenience.
```

**3. Notify Stakeholders**
```bash
# Send resolution notification
# Use same channels as incident notifications
```

### Short-term (4-48 hours)

**4. Create Incident Timeline**
```markdown
# Incident Timeline

## Detection (14:30 UTC)
- Alert: ListingsServiceDown triggered
- On-call engineer paged

## Investigation (14:31-14:45 UTC)
- Service status checked
- Database confirmed unavailable
- Backup recovery initiated

## Resolution (14:45-15:15 UTC)
- Database restored from backup
- Service restarted
- Validation tests passed

## Monitoring (15:15-16:00 UTC)
- Close monitoring for anomalies
- No further issues detected
```

**5. Schedule Postmortem**
```bash
# Within 48 hours of incident resolution
# Invite:
# - Incident Commander
# - On-call engineer(s)
# - Platform Team Lead
# - Relevant stakeholders
# - Database SRE (if DB involved)

# Agenda:
# 1. Timeline review
# 2. Root cause analysis
# 3. What went well
# 4. What could be improved
# 5. Action items
```

**6. Document Lessons Learned**
```markdown
# Postmortem: [INCIDENT NAME]

## Incident Summary
- Date: [DATE]
- Duration: [HH:MM]
- Severity: [P1/P2/P3/P4]
- Incident Commander: [NAME]

## Impact
- Users affected: [COUNT]
- Revenue impact: [$ if applicable]
- SLO breach: [YES/NO]

## Root Cause
[Detailed analysis of what caused the incident]

## Timeline
[Detailed timeline from detection to resolution]

## What Went Well
- [Item 1]
- [Item 2]

## What Could Be Improved
- [Item 1]
- [Item 2]

## Action Items
| Action | Owner | Due Date | Status |
|--------|-------|----------|--------|
| [Action 1] | [Name] | [Date] | Open |
| [Action 2] | [Name] | [Date] | Open |

## Follow-up
- Next review: [DATE]
```

### Long-term (1 week - 1 month)

**7. Implement Improvements**
```bash
# Create Jira tickets for action items
# Priority based on risk reduction

# Examples:
# - Enhanced monitoring
# - Automated recovery procedures
# - Improved documentation
# - Additional testing
# - Infrastructure improvements
```

**8. Update Procedures**
```bash
# Update this document with lessons learned
cd /p/github.com/sveturs/listings/docs/operations
git checkout -b update-dr-procedures

# Make improvements to:
# - DISASTER_RECOVERY.md (this document)
# - RUNBOOK.md
# - TROUBLESHOOTING.md
# - ON_CALL_GUIDE.md

git add .
git commit -m "Update DR procedures based on incident [INCIDENT-ID]"
git push origin update-dr-procedures
# Create PR for review
```

**9. Training and Drills**
```bash
# Schedule disaster recovery drill
# Simulate failure scenario
# Test recovery procedures
# Document drill results
# Improve procedures based on learnings
```

**10. Review and Report**
```bash
# Monthly: Review all incidents
# Quarterly: Disaster recovery drill
# Annually: Full DR plan review and update
```

---

## Emergency Contacts

### Primary Contacts

| Role | Name | PagerDuty | Phone | Email |
|------|------|-----------|-------|-------|
| Incident Commander | Platform Team Lead | `platform-team-lead` | +XXX | platform@svetu.rs |
| On-Call Engineer | Rotation | `listings-oncall` | Via PagerDuty | oncall@svetu.rs |
| Database SRE | DB Team | `db-team` | Via PagerDuty | db-sre@svetu.rs |
| Security Lead | Security Team | `security-team` | +XXX | security@svetu.rs |

### Escalation Path

```
Level 1: On-Call Engineer
    ↓ (If cannot resolve in 15 minutes)
Level 2: Platform Team Lead
    ↓ (If P1 or cannot resolve in 30 minutes)
Level 3: VP Engineering + CTO
    ↓ (If data breach or legal implications)
Level 4: CEO + Legal Team
```

### External Contacts

| Service | Contact | Purpose |
|---------|---------|---------|
| AWS Support | Enterprise Support | Infrastructure issues |
| Database Vendor | PostgreSQL Support | Database-specific issues |
| Security Firm | [Contact] | Security incident response |
| Legal Counsel | [Contact] | Legal/compliance issues |

---

## Appendix

### A. Backup Restoration Scripts

Located in: `/opt/scripts/disaster-recovery/`

```bash
# restore_database.sh - Full database restoration
# restore_pitr.sh - Point-in-time recovery
# validate_backup.sh - Backup verification
# emergency_deploy.sh - Emergency service deployment
```

### B. Security Incident Response Checklist

Located in: `/p/github.com/sveturs/listings/docs/security/INCIDENT_RESPONSE.md`

### C. Data Integrity Validation Queries

Located in: `/p/github.com/sveturs/listings/scripts/validate_data_integrity.sql`

### D. Compliance Requirements

- **GDPR:** Data breach notification within 72 hours
- **PCI DSS:** (If applicable) Follow PCI incident response procedures
- **Local Laws:** Consult legal team for jurisdiction-specific requirements

---

**Document Version:** 1.0.0
**Last Reviewed:** 2025-11-05
**Next Review:** 2025-12-05
**Classification:** CONFIDENTIAL
**Owner:** Platform Team
**Approved By:** CTO

**Distribution:**
- Platform Team
- On-Call Engineers
- Database SRE Team
- Security Team
- VP Engineering
- CTO
