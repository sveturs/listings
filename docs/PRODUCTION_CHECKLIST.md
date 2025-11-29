# Production Deployment Checklist

**Service:** Listings Microservice
**Version:** 0.1.0
**Target Environment:** dev.vondi.rs (Production)
**Date:** 2025-11-05
**Checklist Owner:** Platform Team

---

## Overview

This checklist ensures all production requirements are met before deploying the Listings microservice. Each item must be verified and signed off by the responsible team member.

**Deployment Criteria:** ALL items must be checked ✅ before deployment.

**Sign-off Required By:**
- [ ] Tech Lead
- [ ] DevOps Engineer
- [ ] QA Engineer
- [ ] Security Engineer

---

## 1. Infrastructure Readiness (15 checks)

### 1.1 Server Resources

- [ ] **Server capacity verified**
  - **Check:** `free -h` shows ≥4GB available RAM
  - **Command:** `free -h`
  - **Expected:** Available memory ≥ 4GB
  - **Owner:** DevOps

- [ ] **Disk space sufficient**
  - **Check:** `/opt` partition has ≥50GB free
  - **Command:** `df -h /opt`
  - **Expected:** Available space ≥ 50GB
  - **Owner:** DevOps

- [ ] **CPU capacity available**
  - **Check:** Server not under heavy load (load average < 4)
  - **Command:** `uptime`
  - **Expected:** Load average < 4.0
  - **Owner:** DevOps

### 1.2 Network Configuration

- [ ] **DNS records configured**
  - **Check:** `listings.dev.vondi.rs` resolves to correct IP
  - **Command:** `nslookup listings.dev.vondi.rs`
  - **Expected:** Resolves to dev.vondi.rs IP
  - **Owner:** DevOps

- [ ] **Nginx reverse proxy configured**
  - **Check:** Nginx config includes listings service proxy
  - **Command:** `sudo nginx -t && sudo cat /etc/nginx/sites-enabled/listings.conf`
  - **Expected:** Config valid, proxy_pass to 127.0.0.1:8086
  - **Owner:** DevOps

- [ ] **SSL/TLS certificates valid**
  - **Check:** Certificate for listings.dev.vondi.rs not expired
  - **Command:** `sudo certbot certificates | grep listings.dev.vondi.rs`
  - **Expected:** Valid until date > 30 days from now
  - **Owner:** DevOps

- [ ] **Firewall rules configured**
  - **Check:** Required ports open internally, external access blocked
  - **Command:** `sudo ufw status`
  - **Expected:** 8086 accessible from localhost only, 443 open externally
  - **Owner:** Security

### 1.3 Port Availability

- [ ] **Port 8086 available (HTTP)**
  - **Check:** No process listening on 8086
  - **Command:** `sudo netstat -tulpn | grep 8086`
  - **Expected:** No output (port free)
  - **Owner:** DevOps

- [ ] **Port 50053 available (gRPC)**
  - **Check:** No process listening on 50053
  - **Command:** `sudo netstat -tulpn | grep 50053`
  - **Expected:** No output (port free)
  - **Owner:** DevOps

- [ ] **Port 9093 available (Metrics)**
  - **Check:** No process listening on 9093
  - **Command:** `sudo netstat -tulpn | grep 9093`
  - **Expected:** No output (port free)
  - **Owner:** DevOps

### 1.4 Service Directories

- [ ] **Application directory exists**
  - **Check:** `/opt/listings-service` directory created with correct permissions
  - **Command:** `ls -la /opt/ | grep listings-service`
  - **Expected:** drwxr-xr-x owned by listings user
  - **Owner:** DevOps

- [ ] **Log directory exists**
  - **Check:** `/var/log/listings` directory created
  - **Command:** `ls -la /var/log/ | grep listings`
  - **Expected:** drwxr-xr-x owned by listings user
  - **Owner:** DevOps

- [ ] **Binary deployed**
  - **Check:** Latest binary exists and is executable
  - **Command:** `/opt/listings-service/listings-server --version`
  - **Expected:** Version 0.1.0 (or current release)
  - **Owner:** DevOps

- [ ] **Systemd service configured**
  - **Check:** Service unit file exists and is valid
  - **Command:** `systemctl cat listings-service && systemctl is-enabled listings-service`
  - **Expected:** Unit file exists, service enabled
  - **Owner:** DevOps

- [ ] **Systemd service tested**
  - **Check:** Service can start and stop cleanly
  - **Command:** `sudo systemctl start listings-service && sleep 5 && sudo systemctl status listings-service && sudo systemctl stop listings-service`
  - **Expected:** Service starts successfully, stops cleanly
  - **Owner:** DevOps

---

## 2. Security Configuration (10 checks)

### 2.1 Credentials and Secrets

- [ ] **All CHANGE_ME passwords replaced**
  - **Check:** No "CHANGE_ME" strings in .env.prod
  - **Command:** `grep -i "CHANGE_ME" /opt/listings-service/.env.prod`
  - **Expected:** No matches found
  - **Owner:** Security

- [ ] **Database password strength verified**
  - **Check:** Password ≥32 characters, random
  - **Command:** `echo $VONDILISTINGS_DB_PASSWORD | wc -c`
  - **Expected:** ≥33 (32 chars + newline)
  - **Owner:** Security

- [ ] **Redis password strength verified**
  - **Check:** Password ≥32 characters, random
  - **Command:** `echo $VONDILISTINGS_REDIS_PASSWORD | wc -c`
  - **Expected:** ≥33
  - **Owner:** Security

- [ ] **Passwords stored in vault**
  - **Check:** All production passwords backed up in 1Password/Vault
  - **Command:** Manual verification in password manager
  - **Expected:** Entry exists for "Listings Service Production"
  - **Owner:** Security

### 2.2 Access Control

- [ ] **Service user created with minimal permissions**
  - **Check:** Non-root user "listings" exists with no sudo
  - **Command:** `id listings && sudo -l -U listings`
  - **Expected:** User exists, "not allowed to run sudo"
  - **Owner:** Security

- [ ] **File permissions restricted**
  - **Check:** .env.prod readable only by listings user
  - **Command:** `ls -l /opt/listings-service/.env.prod`
  - **Expected:** -rw------- (600) owned by listings:listings
  - **Owner:** Security

- [ ] **Auth public key accessible**
  - **Check:** JWT public key exists and readable by service
  - **Command:** `sudo -u listings cat $VONDILISTINGS_AUTH_PUBLIC_KEY_PATH | head -1`
  - **Expected:** "-----BEGIN PUBLIC KEY-----"
  - **Owner:** Security

### 2.3 Network Security

- [ ] **CORS origins restricted**
  - **Check:** No wildcards (*) in CORS configuration
  - **Command:** `grep CORS_ALLOWED_ORIGINS /opt/listings-service/.env.prod`
  - **Expected:** Specific origins only (dev.vondi.rs, devapi.vondi.rs)
  - **Owner:** Security

- [ ] **Rate limiting enabled**
  - **Check:** Rate limit enabled in production config
  - **Command:** `grep RATE_LIMIT_ENABLED /opt/listings-service/.env.prod`
  - **Expected:** VONDILISTINGS_RATE_LIMIT_ENABLED=true
  - **Owner:** Security

- [ ] **Security headers enabled**
  - **Check:** Security headers enabled in config
  - **Command:** `grep SECURITY_HEADERS_ENABLED /opt/listings-service/.env.prod`
  - **Expected:** VONDILISTINGS_SECURITY_HEADERS_ENABLED=true
  - **Owner:** Security

---

## 3. Database Configuration (6 checks)

### 3.1 Database Setup

- [ ] **PostgreSQL instance running**
  - **Check:** PostgreSQL accessible on port 35433
  - **Command:** `sudo systemctl status postgresql && sudo netstat -tulpn | grep 35433`
  - **Expected:** Service active, port listening
  - **Owner:** DevOps

- [ ] **Database exists**
  - **Check:** Database "listings_dev_db" created
  - **Command:** `psql "postgres://listings_user:PASSWORD@localhost:35433/postgres?sslmode=disable" -c "\l" | grep listings_dev_db`
  - **Expected:** Database listed
  - **Owner:** DevOps

- [ ] **Database user has correct permissions**
  - **Check:** User can read/write to listings tables
  - **Command:** `psql "postgres://listings_user:PASSWORD@localhost:35433/listings_dev_db?sslmode=disable" -c "\dp listings"`
  - **Expected:** User has SELECT, INSERT, UPDATE, DELETE
  - **Owner:** DevOps

- [ ] **Migrations applied successfully**
  - **Check:** All migrations up to date
  - **Command:** `cd /opt/listings-service && ./listings-server migrate status`
  - **Expected:** All migrations applied, no pending
  - **Owner:** DevOps

- [ ] **Database connection pool tuned**
  - **Check:** Pool settings match expected load
  - **Command:** `grep DB_MAX_OPEN_CONNS /opt/listings-service/.env.prod`
  - **Expected:** MAX_OPEN_CONNS=50, MAX_IDLE_CONNS=25
  - **Owner:** DevOps

- [ ] **Database backup configured**
  - **Check:** Automated backup script scheduled
  - **Command:** `crontab -l | grep listings-backup`
  - **Expected:** Daily backup job exists
  - **Owner:** DevOps

---

## 4. Dependencies Health (9 checks)

### 4.1 Redis

- [ ] **Redis instance running**
  - **Check:** Redis accessible on port 36380
  - **Command:** `redis-cli -p 36380 -a PASSWORD PING`
  - **Expected:** PONG
  - **Owner:** DevOps

- [ ] **Redis password configured**
  - **Check:** Redis requires authentication
  - **Command:** `redis-cli -p 36380 PING`
  - **Expected:** NOAUTH Authentication required
  - **Owner:** Security

- [ ] **Redis persistence enabled**
  - **Check:** RDB or AOF persistence configured
  - **Command:** `redis-cli -p 36380 -a PASSWORD CONFIG GET save`
  - **Expected:** Persistence configured (not empty)
  - **Owner:** DevOps

### 4.2 OpenSearch

- [ ] **OpenSearch accessible**
  - **Check:** OpenSearch responding on port 9200
  - **Command:** `curl -s http://localhost:9200/_cluster/health | jq .status`
  - **Expected:** "green" or "yellow"
  - **Owner:** DevOps

- [ ] **OpenSearch index exists**
  - **Check:** marketplace_listings index created
  - **Command:** `curl -s http://localhost:9200/_cat/indices | grep marketplace_listings`
  - **Expected:** Index exists with documents
  - **Owner:** DevOps

- [ ] **OpenSearch credentials work**
  - **Check:** Can authenticate with configured credentials
  - **Command:** `curl -u admin:PASSWORD -s http://localhost:9200/_cluster/health | jq .status`
  - **Expected:** "green" or "yellow"
  - **Owner:** Security

### 4.3 MinIO

- [ ] **MinIO accessible**
  - **Check:** MinIO responding on port 9000
  - **Command:** `curl -s http://localhost:9000/minio/health/live`
  - **Expected:** HTTP 200
  - **Owner:** DevOps

- [ ] **MinIO bucket exists**
  - **Check:** listings-images bucket created
  - **Command:** `mc ls minio/listings-images/`
  - **Expected:** Bucket exists (may be empty)
  - **Owner:** DevOps

### 4.4 Auth Service

- [ ] **Auth service accessible**
  - **Check:** Auth service responding
  - **Command:** `curl -s http://localhost:28086/health | jq .status`
  - **Expected:** "ok" or "healthy"
  - **Owner:** Platform Team

---

## 5. Monitoring Setup (8 checks)

### 5.1 Prometheus

- [ ] **Prometheus scraping service metrics**
  - **Check:** Listings target in Prometheus
  - **Command:** `curl -s http://prometheus.vondi.rs:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.job == "listings")'`
  - **Expected:** Target exists, state="up"
  - **Owner:** DevOps

- [ ] **Metrics endpoint accessible**
  - **Check:** Can fetch metrics from service
  - **Command:** `curl -s http://localhost:9093/metrics | grep listings_grpc_requests_total`
  - **Expected:** Metrics present
  - **Owner:** DevOps

- [ ] **Recording rules loaded**
  - **Check:** Prometheus has listings recording rules
  - **Command:** `curl -s http://prometheus.vondi.rs:9090/api/v1/rules | jq '.data.groups[] | select(.name == "listings_recording_rules")'`
  - **Expected:** Group exists with rules
  - **Owner:** DevOps

### 5.2 Alerting

- [ ] **Alert rules loaded**
  - **Check:** Prometheus has listings alert rules
  - **Command:** `curl -s http://prometheus.vondi.rs:9090/api/v1/rules | jq '.data.groups[] | select(.name == "listings_alerts")'`
  - **Expected:** Group exists with alerts
  - **Owner:** DevOps

- [ ] **AlertManager configured**
  - **Check:** AlertManager has routing for listings alerts
  - **Command:** `curl -s http://alertmanager.vondi.rs:9093/api/v2/status | jq .config.route`
  - **Expected:** Route exists for service="listings"
  - **Owner:** DevOps

- [ ] **PagerDuty integration working**
  - **Check:** Test alert sent to PagerDuty
  - **Command:** `cd /p/github.com/sveturs/listings/deployment/prometheus && ./test-alerts.sh --test-pagerduty`
  - **Expected:** Alert received in PagerDuty
  - **Owner:** DevOps

- [ ] **Slack notifications working**
  - **Check:** Test alert sent to Slack
  - **Command:** `cd /p/github.com/sveturs/listings/deployment/prometheus && ./test-alerts.sh --test-slack`
  - **Expected:** Message in #alerts-listings channel
  - **Owner:** DevOps

### 5.3 Grafana

- [ ] **Grafana dashboards imported**
  - **Check:** All 4 dashboards exist
  - **Command:** Manual check in Grafana UI → Dashboards → Listings
  - **Expected:** Overview, Details, Database, Redis, SLO dashboards visible
  - **Owner:** DevOps

---

## 6. Backup System (6 checks)

### 6.1 Database Backups

- [ ] **Backup script exists**
  - **Check:** Backup script in place and executable
  - **Command:** `ls -l /opt/listings-service/scripts/backup/backup-db.sh`
  - **Expected:** -rwxr-xr-x (executable)
  - **Owner:** DevOps

- [ ] **Backup script tested**
  - **Check:** Manual backup completes successfully
  - **Command:** `cd /opt/listings-service && ./scripts/backup/backup-db.sh`
  - **Expected:** Backup file created in /backup/listings/
  - **Owner:** DevOps

- [ ] **Backup cron job scheduled**
  - **Check:** Cron job for daily backups
  - **Command:** `crontab -u listings -l | grep backup-db.sh`
  - **Expected:** Daily backup at 02:00
  - **Owner:** DevOps

### 6.2 Restore Procedures

- [ ] **Restore script exists**
  - **Check:** Restore script in place and executable
  - **Command:** `ls -l /opt/listings-service/scripts/backup/restore-db.sh`
  - **Expected:** -rwxr-xr-x (executable)
  - **Owner:** DevOps

- [ ] **Restore tested successfully**
  - **Check:** Restore from backup verified
  - **Command:** `cd /opt/listings-service && ./scripts/backup/test-backup-restore.sh`
  - **Expected:** Test passes, data restored correctly
  - **Owner:** DevOps

- [ ] **Backup retention configured**
  - **Check:** Old backups cleaned up automatically
  - **Command:** `grep RETENTION /opt/listings-service/.env.prod`
  - **Expected:** BACKUP_RETENTION_DAYS=30
  - **Owner:** DevOps

---

## 7. Documentation (5 checks)

### 7.1 Operations Documentation

- [ ] **Runbook reviewed and accessible**
  - **Check:** Runbook exists and is up to date
  - **Command:** `ls -l /opt/listings-service/docs/operations/RUNBOOK.md`
  - **Expected:** File exists, modified within last 7 days
  - **Owner:** Platform Team

- [ ] **Troubleshooting guide accessible**
  - **Check:** Troubleshooting guide exists
  - **Command:** `ls -l /opt/listings-service/docs/operations/TROUBLESHOOTING.md`
  - **Expected:** File exists
  - **Owner:** Platform Team

- [ ] **SLO documentation complete**
  - **Check:** SLO targets and error budgets documented
  - **Command:** `grep "99.9%" /opt/listings-service/docs/operations/SLO_GUIDE.md`
  - **Expected:** Availability SLO documented
  - **Owner:** Platform Team

### 7.2 Team Readiness

- [ ] **On-call rotation configured**
  - **Check:** PagerDuty schedule includes listings service
  - **Command:** Manual verification in PagerDuty
  - **Expected:** Schedule exists with ≥2 engineers
  - **Owner:** Team Lead

- [ ] **Team trained on dashboards**
  - **Check:** All team members completed dashboard walkthrough
  - **Command:** Manual verification (training log)
  - **Expected:** 100% of on-call engineers trained
  - **Owner:** Team Lead

---

## 8. Deployment Validation (10 checks)

### 8.1 Health Checks

- [ ] **Health endpoint returns 200**
  - **Check:** Basic health check passes
  - **Command:** `curl -s -o /dev/null -w "%{http_code}" http://localhost:8086/health`
  - **Expected:** 200
  - **Owner:** QA

- [ ] **Readiness endpoint returns 200**
  - **Check:** Readiness check passes (all dependencies healthy)
  - **Command:** `curl -s http://localhost:8086/ready | jq .status`
  - **Expected:** "ok" with all dependencies "healthy"
  - **Owner:** QA

- [ ] **Metrics endpoint returns data**
  - **Check:** Prometheus metrics exposed
  - **Command:** `curl -s http://localhost:9093/metrics | head -20`
  - **Expected:** Metrics data returned
  - **Owner:** QA

### 8.2 Smoke Tests

- [ ] **gRPC health check passes**
  - **Check:** gRPC service responding
  - **Command:** `grpcurl -plaintext localhost:50053 grpc.health.v1.Health/Check`
  - **Expected:** "status": "SERVING"
  - **Owner:** QA

- [ ] **Can create listing**
  - **Check:** POST /api/v1/listings works
  - **Command:** `curl -X POST -H "Content-Type: application/json" -d '{"title":"Test","price":100}' http://localhost:8086/api/v1/listings`
  - **Expected:** HTTP 201, listing created
  - **Owner:** QA

- [ ] **Can retrieve listing**
  - **Check:** GET /api/v1/listings/:id works
  - **Command:** `curl -s http://localhost:8086/api/v1/listings/1 | jq .id`
  - **Expected:** Listing data returned
  - **Owner:** QA

- [ ] **Can search listings**
  - **Check:** GET /api/v1/listings/search works
  - **Command:** `curl -s "http://localhost:8086/api/v1/listings/search?q=test" | jq .total`
  - **Expected:** Search results returned
  - **Owner:** QA

### 8.3 Load Testing

- [ ] **Load test passes**
  - **Check:** Service handles expected load
  - **Command:** `cd /opt/listings-service && ./scripts/load_test.sh`
  - **Expected:** All requests succeed, P95 < 1s
  - **Owner:** QA

- [ ] **No memory leaks detected**
  - **Check:** Memory usage stable under load
  - **Command:** `cd /opt/listings-service && ./scripts/profile_memory.sh`
  - **Expected:** Memory usage stable after 10 minutes
  - **Owner:** QA

### 8.4 Rollback Readiness

- [ ] **Rollback script tested**
  - **Check:** Can rollback to previous version
  - **Command:** `cd /opt/listings-service && ./scripts/rollback-prod.sh --dry-run`
  - **Expected:** Dry run succeeds, shows rollback steps
  - **Owner:** DevOps

---

## 9. Final Sign-off

### 9.1 Sign-off Matrix

| Role | Name | Date | Signature |
|------|------|------|-----------|
| **Tech Lead** | | | [ ] |
| **DevOps Engineer** | | | [ ] |
| **QA Engineer** | | | [ ] |
| **Security Engineer** | | | [ ] |

### 9.2 Deployment Decision

- [ ] **All checklist items verified** ✅
- [ ] **All sign-offs obtained** ✅
- [ ] **Deployment window scheduled** ✅
- [ ] **Rollback plan reviewed** ✅
- [ ] **Stakeholders notified** ✅

**Deployment Approved:** [ ] YES / [ ] NO

**Approved By:** _________________________ **Date:** _____________

**Deployment Date:** _____________
**Deployment Time:** _____________
**Expected Duration:** _____________ hours

---

## 10. Post-Deployment Verification (First 24 hours)

### Immediate (First 5 minutes)

- [ ] Service started successfully
- [ ] Health checks passing
- [ ] No error alerts firing
- [ ] Logs show normal startup

### Short-term (First hour)

- [ ] Request rate normal (>0 RPS)
- [ ] Error rate < 1%
- [ ] P95 latency < 1s
- [ ] Database connections stable
- [ ] Redis cache working

### Medium-term (First 24 hours)

- [ ] No unexpected alerts
- [ ] SLO compliance maintained
- [ ] Memory usage stable
- [ ] No customer complaints
- [ ] Dashboards showing expected metrics

---

## Troubleshooting

**If any check fails:**

1. **STOP** - Do not proceed to next section
2. **Investigate** - Determine root cause
3. **Fix** - Resolve the issue
4. **Retest** - Verify fix with check command
5. **Document** - Add notes to deployment log
6. **Continue** - Proceed to next check

**If multiple checks fail:**

1. Consider postponing deployment
2. Review with team lead
3. Update deployment timeline
4. Notify stakeholders

**Emergency rollback criteria:**

- Error rate > 5%
- Service down > 1 minute
- Data corruption detected
- Security vulnerability discovered

---

## Resources

- **Runbook:** `/opt/listings-service/docs/operations/RUNBOOK.md`
- **Troubleshooting:** `/opt/listings-service/docs/operations/TROUBLESHOOTING.md`
- **Monitoring:** `/opt/listings-service/docs/operations/MONITORING_GUIDE.md`
- **Rollback:** `/opt/listings-service/docs/ROLLBACK.md`
- **Architecture:** `/opt/listings-service/README.md`

---

**Checklist Version:** 1.0.0
**Last Updated:** 2025-11-05
**Next Review:** After first deployment
