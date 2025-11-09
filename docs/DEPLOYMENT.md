# Production Deployment Guide

Complete guide for deploying the Listings microservice to production using Blue-Green deployment strategy.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Deployment Process](#deployment-process)
- [Post-Deployment](#post-deployment)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Overview

### Blue-Green Deployment Strategy

We use Blue-Green deployment to achieve zero-downtime deployments:

- **Blue:** Current production version (stable)
- **Green:** New version being deployed
- **Canary:** Gradual traffic shift (10% → 50% → 100%)

### Deployment Flow

```
1. Pre-deployment validation
2. Build & upload new binary (Green)
3. Run database migrations (Blue still serving)
4. Start Green instance
5. Smoke tests on Green
6. Canary deployment:
   - 10% traffic to Green (5 min monitoring)
   - 50% traffic to Green (5 min monitoring)
   - 100% traffic to Green
7. Decommission Blue (after 10 min rollback window)
8. Generate deployment report
```

### Key Features

- **Zero Downtime:** Traffic switches without interruption
- **Automated Tests:** Smoke tests before traffic switch
- **Gradual Rollout:** Canary phases catch issues early
- **Quick Rollback:** Instant switch back to Blue if needed
- **Audit Trail:** Comprehensive logs and reports

---

## Prerequisites

### Local Machine

1. **SSH Access**
   ```bash
   ssh-copy-id user@production-server
   ssh user@production-server "echo 'OK'"
   ```

2. **Required Tools**
   - Go 1.21+
   - Git
   - curl, jq
   - bash 4.0+

3. **Repository Access**
   ```bash
   cd /p/github.com/sveturs/listings
   git pull origin main
   ```

### Production Server

1. **System Requirements**
   - Ubuntu 20.04+ or similar
   - 20GB+ free disk space
   - 4GB+ RAM

2. **Services Running**
   - Docker (for PostgreSQL, Redis)
   - Nginx (load balancer)
   - PostgreSQL container
   - Redis container

3. **Directory Structure**
   ```
   /opt/listings/
   ├── bin/
   │   ├── listings-blue   (current stable)
   │   └── listings-green  (new version)
   ├── logs/
   │   ├── blue.log
   │   ├── green.log
   │   └── archived/
   ├── backups/
   │   └── backup-*.sql.gz
   ├── .env.blue
   └── .env.green
   ```

---

## Configuration

### Step 1: Create Deployment Config

```bash
cd /p/github.com/sveturs/listings/scripts
cp .env.deploy.example .env.deploy
```

### Step 2: Edit Configuration

```bash
nano .env.deploy
```

Required variables:

```bash
# Production server
PROD_HOST="production.example.com"
PROD_USER="deploy"
PROD_DIR="/opt/listings"
PROD_DOMAIN="listings.example.com"

# Blue/Green ports
BLUE_PORT=8080
GREEN_PORT=8081

# Database
DB_USER="listings"
DB_NAME="listings_prod"

# Nginx
NGINX_CONFIG_PATH="/etc/nginx/sites-available/listings"

# Notifications (optional)
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
ALERT_EMAIL="ops@example.com"

# On-call (optional)
ONCALL_ENGINEER="John Doe <john@example.com>"
INCIDENT_COMMANDER="Jane Smith <jane@example.com>"
```

### Step 3: Create Production Environment File

```bash
cp .env.example .env.prod
nano .env.prod
```

Update with production values:

```env
# Server
PORT=8080
ENV=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=listings_prod
DB_USER=listings
DB_PASSWORD=<secure-password>

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Monitoring
SENTRY_DSN=<your-sentry-dsn>
```

---

## Deployment Process

### Quick Start

```bash
cd /p/github.com/sveturs/listings

# Deploy to production
./scripts/deploy-to-prod.sh
```

### Detailed Steps

#### 1. Pre-Deployment Validation

```bash
# Validate environment
./scripts/validate-deployment.sh --verbose
```

This checks:
- SSH connectivity
- Disk space (>20GB)
- Services availability (PostgreSQL, Redis, Nginx)
- No deployments in progress
- Git status

#### 2. Dry Run (Recommended)

```bash
# Test deployment without making changes
./scripts/deploy-to-prod.sh --dry-run
```

This shows what would happen without actually deploying.

#### 3. Full Deployment

```bash
# Execute production deployment
./scripts/deploy-to-prod.sh
```

**What happens:**

1. **Validation** - Runs pre-flight checks
2. **Backup** - Creates database backup
3. **Build** - Compiles optimized binary
4. **Upload** - Transfers binary and config to server
5. **Migrations** - Runs database migrations (Blue still serving)
6. **Start Green** - Launches new version on separate port
7. **Smoke Tests** - Validates Green instance health
8. **Canary Phase 1** - Routes 10% traffic to Green, monitors for 5 minutes
9. **Canary Phase 2** - Routes 50% traffic to Green, monitors for 5 minutes
10. **Full Switch** - Routes 100% traffic to Green
11. **Rollback Window** - Waits 10 minutes for issues
12. **Decommission** - Stops Blue, promotes Green to Blue
13. **Report** - Generates deployment report

#### 4. Skip Tests (Not Recommended)

```bash
# Deploy without smoke tests (emergency only)
./scripts/deploy-to-prod.sh --skip-tests
```

#### 5. Verbose Mode

```bash
# Show detailed execution
./scripts/deploy-to-prod.sh --verbose
```

---

## Post-Deployment

### Monitor Production

#### Application Logs

```bash
# Real-time Green logs
ssh user@production-server 'tail -f /opt/listings/logs/green.log'

# Filter errors
ssh user@production-server 'tail -f /opt/listings/logs/green.log | grep -i error'
```

#### Nginx Logs

```bash
# Access logs
ssh user@production-server 'sudo tail -f /var/log/nginx/access.log'

# Error logs
ssh user@production-server 'sudo tail -f /var/log/nginx/error.log'
```

#### System Metrics

```bash
# CPU, Memory, Disk
ssh user@production-server 'htop'

# Disk space
ssh user@production-server 'df -h'

# Network connections
ssh user@production-server 'netstat -tlnp | grep listings'
```

### Health Checks

```bash
# Application health
curl -I https://listings.example.com/health

# Metrics
curl https://listings.example.com/metrics

# Database connectivity
curl https://listings.example.com/health/db

# Redis connectivity
curl https://listings.example.com/health/redis
```

### Smoke Tests (Post-Deploy)

```bash
# Run smoke tests against production
./scripts/smoke-tests.sh --host production.example.com --port 80
```

### Review Deployment Report

```bash
# Find latest deployment report
ls -lt logs/deployments/*-report.md | head -1

# View report
cat logs/deployments/deploy-20250105-143022-report.md
```

### Sign-Off Checklist

- [ ] All smoke tests passed
- [ ] No errors in production logs (first hour)
- [ ] Database migrations successful
- [ ] Metrics within acceptable range
- [ ] No customer complaints
- [ ] Monitoring alerts normal
- [ ] Backup verified and accessible
- [ ] Team notified

---

## Troubleshooting

### Deployment Failed at Build

**Error:** `Build failed`

**Solution:**
```bash
# Check Go version
go version

# Run tests locally
go test ./...

# Try building locally
go build -o bin/listings ./cmd/server/main.go
```

### Deployment Failed at Migrations

**Error:** `Migrations failed`

**Solution:**
```bash
# Check migration files
ls -la migrations/

# Test migrations locally
go run cmd/migrate/main.go up

# Rollback if needed
./scripts/rollback-prod.sh --restore-db
```

### Smoke Tests Failed

**Error:** `Smoke tests failed on Green instance`

**Solution:**
```bash
# Check Green logs
ssh user@production-server 'tail -100 /opt/listings/logs/green.log'

# Verify Green is running
ssh user@production-server 'pgrep -f listings-green'

# Check port binding
ssh user@production-server 'netstat -tlnp | grep 8081'

# Manually test endpoint
curl http://production-server:8081/health
```

### High Error Rate in Canary

**Error:** `High error rate detected in Phase 1`

**Solution:**
- Deployment automatically rolls back
- Review Green logs for errors
- Fix issues and redeploy

### Unable to Connect to Database

**Error:** `PostgreSQL is not accessible`

**Solution:**
```bash
# Check PostgreSQL container
ssh user@production-server 'docker ps | grep postgres'

# Restart PostgreSQL if needed
ssh user@production-server 'docker restart listings_postgres'

# Check connection
ssh user@production-server 'docker exec listings_postgres pg_isready'
```

### Disk Space Issues

**Error:** `Only 5 GB available (minimum 20 GB required)`

**Solution:**
```bash
# Check disk usage
ssh user@production-server 'du -sh /opt/listings/*'

# Clean old backups (keep last 10)
ssh user@production-server 'cd /opt/listings/backups && ls -t | tail -n +11 | xargs rm -f'

# Clean old logs
ssh user@production-server 'find /opt/listings/logs/archived -mtime +30 -delete'
```

### Nginx Configuration Issues

**Error:** `Nginx configuration test failed`

**Solution:**
```bash
# Test current config
ssh user@production-server 'sudo nginx -t'

# Restore backup
ssh user@production-server 'sudo cp /tmp/nginx_backup.conf /etc/nginx/sites-available/listings'

# Reload Nginx
ssh user@production-server 'sudo systemctl reload nginx'
```

---

## FAQ

### How long does a deployment take?

Typical deployment: 25-35 minutes
- Build: 2-3 minutes
- Upload: 1 minute
- Migrations: 1-5 minutes (depends on changes)
- Smoke tests: 2 minutes
- Canary Phase 1: 5 minutes
- Canary Phase 2: 5 minutes
- Rollback window: 10 minutes

### Can I deploy during business hours?

Yes! Blue-Green deployment with canary phases ensures zero downtime. However, for major changes, consider deploying during low-traffic periods.

### What if I need to rollback?

```bash
# Immediate rollback
./scripts/rollback-prod.sh --reason "issue description"

# Rollback with database restore
./scripts/rollback-prod.sh --reason "bad migration" --restore-db
```

See [ROLLBACK.md](./ROLLBACK.md) for details.

### How do I monitor the deployment?

1. Watch deployment logs in real-time
2. Monitor Slack notifications (if configured)
3. Check server logs via SSH
4. Review deployment report after completion

### Can I pause the deployment?

No automatic pause, but:
- Each canary phase has 5-minute monitoring
- If error rate is high, deployment auto-rolls back
- You can manually rollback anytime

### What happens to database during deployment?

1. Backup created before any changes
2. Migrations run while Blue serves traffic
3. Green uses updated schema
4. If rollback needed, database can be restored

### How are secrets managed?

- Secrets stored in `.env.prod` (NOT in git)
- File transferred securely via SCP
- Permissions: `600` (owner read/write only)
- Backed up separately from code

### What if deployment locks up?

```bash
# Remove stale lock file
ssh user@production-server 'rm /opt/listings/.deployment.lock'

# Verify no deployment is actually running
ssh user@production-server 'ps aux | grep deploy'
```

### How do I deploy a hotfix?

Same process, but consider:
```bash
# Create hotfix branch
git checkout -b hotfix/critical-fix

# Make changes, commit, push

# Deploy immediately (skip tests if critical)
./scripts/deploy-to-prod.sh --skip-tests
```

### Can I deploy multiple times per day?

Yes, but:
- Ensure previous deployment is complete
- Monitor production after each deployment
- Consider batching small changes

---

## Related Documents

- [Rollback Procedures](./ROLLBACK.md)
- [Runbook](./RUNBOOK.md)
- [Architecture](./ARCHITECTURE.md)

---

## Support

**On-Call Team:**
- Primary: Check `.env.deploy` for current on-call engineer
- Escalation: Incident Commander

**Deployment Issues:**
- Slack: `#deployments` channel
- PagerDuty: Alerts sent automatically for failures

**Questions:**
- Documentation: This file
- Team Wiki: [Internal Wiki Link]

---

*Last updated: 2025-11-05*
*Maintained by: DevOps Team*
