# Production Rollback Procedures

Emergency procedures for rolling back failed deployments.

## Table of Contents

- [When to Rollback](#when-to-rollback)
- [Rollback Types](#rollback-types)
- [Quick Rollback](#quick-rollback)
- [Database Rollback](#database-rollback)
- [Manual Rollback](#manual-rollback)
- [Post-Rollback](#post-rollback)
- [Incident Response](#incident-response)

---

## When to Rollback

### Automatic Rollback Triggers

The deployment system automatically rolls back if:
- Smoke tests fail on Green instance
- High error rate detected during canary phase (>10 errors in 1000 log lines)
- Green instance health checks fail

### Manual Rollback Indicators

Consider manual rollback if:
- **High Error Rate:** >5% errors in production logs
- **Performance Degradation:** Response times >2x baseline
- **Database Issues:** Data corruption or integrity problems
- **Customer Impact:** Multiple user complaints
- **Critical Bug:** Security vulnerability or data loss
- **Service Unavailable:** Repeated 500/503 errors

### Decision Matrix

| Severity | Impact | Action |
|----------|--------|--------|
| P0 - Critical | Service down | Immediate rollback |
| P1 - High | Data loss risk | Rollback within 5 minutes |
| P2 - Medium | Degraded performance | Rollback within 15 minutes |
| P3 - Low | Minor issues | Fix forward or schedule rollback |

---

## Rollback Types

### 1. Traffic-Only Rollback

**When:** Application issues, no database changes
**Duration:** < 1 minute
**Risk:** Low

Switches traffic back to Blue without database changes.

### 2. Full Rollback

**When:** Application issues + database migrations executed
**Duration:** 5-10 minutes
**Risk:** Medium (database restore involved)

Switches traffic and restores database from backup.

### 3. Partial Rollback

**When:** Specific feature causing issues
**Duration:** Varies
**Risk:** Low to Medium

Feature flag disable or configuration change.

---

## Quick Rollback

### Emergency Rollback (Traffic Only)

```bash
cd /p/github.com/sveturs/listings

# Execute immediate rollback
./scripts/rollback-prod.sh --reason "high error rate"
```

**What happens:**
1. ✅ Traffic switches to Blue (100%) - **Instant**
2. ✅ Green instance stopped
3. ✅ Blue instance verified healthy
4. ✅ Incident logs captured
5. ✅ Incident report generated
6. ✅ Team notified

**Expected duration:** 1-2 minutes

### Verify Rollback

```bash
# Check traffic distribution
curl -I https://listings.example.com/ | grep 'X-.*-Weight'
# Should show: X-Blue-Weight: 100, X-Green-Weight: 0

# Check Blue is serving
curl https://listings.example.com/health

# Check logs for errors
ssh user@production-server 'tail -100 /opt/listings/logs/blue.log | grep ERROR'
```

---

## Database Rollback

### When to Restore Database

Restore database if:
- Migration caused data corruption
- Schema change broke functionality
- Data loss occurred
- Referential integrity violated

### Full Rollback with Database Restore

```bash
# Rollback with database restoration
./scripts/rollback-prod.sh --reason "migration failed" --restore-db
```

**⚠️ WARNING:** This will restore database to pre-deployment state. All data changes since deployment will be lost!

**Confirmation required:**
```
This will OVERWRITE the current database. Type 'RESTORE' to confirm:
```

**What happens:**
1. ✅ Traffic switches to Blue (100%)
2. ✅ Green instance stopped
3. ✅ Blue instance verified
4. ⚠️ Database backup restored (requires confirmation)
5. ✅ Incident logs captured
6. ✅ Team notified

**Expected duration:** 5-10 minutes (depends on database size)

### Verify Database Restore

```bash
# Check latest backup
ssh user@production-server 'ls -lth /opt/listings/backups/*.sql.gz | head -1'

# Verify data integrity
ssh user@production-server 'docker exec listings_postgres psql -U listings -d listings_prod -c "SELECT COUNT(*) FROM listings;"'

# Check application connectivity
curl https://listings.example.com/health/db
```

---

## Manual Rollback

### Scenario 1: Rollback Script Failed

If automated rollback fails:

```bash
# 1. SSH to production server
ssh user@production-server

# 2. Switch Nginx to Blue manually
sudo nano /etc/nginx/sites-available/listings

# Change upstream weights:
# Blue: weight=100
# Green: weight=0

# 3. Test config
sudo nginx -t

# 4. Reload Nginx
sudo systemctl reload nginx

# 5. Stop Green instance
pkill -f listings-green

# 6. Verify Blue is serving
curl localhost:8080/health
```

### Scenario 2: Both Blue and Green Failed

If both instances are down:

```bash
# 1. SSH to production server
ssh user@production-server

# 2. Check what's running
ps aux | grep listings
docker ps

# 3. Restart Blue instance
cd /opt/listings
nohup ./bin/listings-blue --port=8080 --env-file=.env.blue > logs/blue.log 2>&1 &

# 4. Wait for startup
sleep 10

# 5. Verify health
curl localhost:8080/health

# 6. Update Nginx to point to Blue only
# (follow Scenario 1 steps 2-4)
```

### Scenario 3: Database Connection Lost

If application can't connect to database:

```bash
# 1. Check PostgreSQL container
ssh user@production-server
docker ps | grep postgres

# 2. If not running, start it
docker start listings_postgres

# 3. Wait for PostgreSQL to be ready
docker exec listings_postgres pg_isready -U listings

# 4. Restart application
pkill -f listings-blue
cd /opt/listings
nohup ./bin/listings-blue --port=8080 --env-file=.env.blue > logs/blue.log 2>&1 &

# 5. Verify connectivity
curl localhost:8080/health/db
```

### Scenario 4: Configuration Issues

If configuration is wrong:

```bash
# 1. SSH to production server
ssh user@production-server

# 2. Check current config
cat /opt/listings/.env.blue

# 3. Compare with backup
cat /opt/listings/.env.blue.old

# 4. Restore backup config
cp /opt/listings/.env.blue.old /opt/listings/.env.blue

# 5. Restart Blue instance
pkill -f listings-blue
cd /opt/listings
nohup ./bin/listings-blue --port=8080 --env-file=.env.blue > logs/blue.log 2>&1 &

# 6. Verify
curl localhost:8080/health
```

---

## Post-Rollback

### Immediate Actions

**Within 5 minutes:**

1. **Verify Service**
   ```bash
   # Health check
   curl https://listings.example.com/health

   # Test critical endpoints
   ./scripts/smoke-tests.sh --host production.example.com --port 80
   ```

2. **Monitor Logs**
   ```bash
   # Watch for errors
   ssh user@production-server 'tail -f /opt/listings/logs/blue.log | grep -i error'
   ```

3. **Check Metrics**
   - Error rate should drop to normal
   - Response times should improve
   - CPU/Memory usage stable

### Incident Documentation

**Within 30 minutes:**

1. **Review Incident Report**
   ```bash
   # Generated automatically by rollback script
   cat incidents/rollback-YYYYMMDD-HHMMSS/INCIDENT_REPORT.md
   ```

2. **Capture Logs**
   ```bash
   # Logs are automatically archived in:
   ls -l incidents/rollback-YYYYMMDD-HHMMSS/

   # Files included:
   # - green-rollback-*.log.gz (application logs)
   # - nginx-error.log (nginx errors)
   # - systemd.log (system logs)
   ```

3. **Notify Stakeholders**
   - Update status page
   - Notify affected customers
   - Inform team via Slack/email

### Root Cause Analysis

**Within 24 hours:**

1. **Analyze Logs**
   ```bash
   # Green instance logs
   gunzip -c incidents/rollback-YYYYMMDD-HHMMSS/green-rollback-*.log.gz | less

   # Look for:
   # - Panic/fatal errors
   # - Database query failures
   # - Timeout errors
   # - Memory issues
   ```

2. **Identify Root Cause**
   - Code bug?
   - Configuration issue?
   - Database migration problem?
   - Infrastructure limitation?

3. **Create Bug Ticket**
   - Title: Clear description of issue
   - Priority: Based on severity
   - Assignee: Developer who can fix
   - Labels: `bug`, `production`, `rollback`

### Postmortem Meeting

**Within 1 week:**

Agenda:
1. **Timeline of events** - What happened and when
2. **Root cause** - Why did it happen
3. **Impact assessment** - Who/what was affected
4. **Action items** - How to prevent recurrence
5. **Documentation updates** - Improve runbooks

Template:
```markdown
# Postmortem: [Issue Title]

## Summary
Brief description of the incident

## Timeline
- HH:MM - Deployment started
- HH:MM - Issue detected
- HH:MM - Rollback initiated
- HH:MM - Service restored

## Root Cause
Detailed explanation of what went wrong

## Impact
- Downtime: X minutes
- Users affected: X
- Revenue impact: $X

## Resolution
How the issue was resolved

## Action Items
- [ ] Fix code bug (Assignee, Due date)
- [ ] Add test coverage (Assignee, Due date)
- [ ] Update documentation (Assignee, Due date)
- [ ] Improve monitoring (Assignee, Due date)

## Lessons Learned
What we learned and how to improve
```

---

## Incident Response

### Communication Protocol

#### Internal Communication

**Slack Channels:**
- `#incidents` - Real-time incident updates
- `#deployments` - Deployment notifications
- `#engineering` - Technical discussion

**Update Template:**
```
🚨 INCIDENT: [Title]
Status: Investigating / Mitigating / Resolved
Impact: [High/Medium/Low]
Started: HH:MM
ETA: HH:MM
Actions: [What we're doing]
Owner: @username
```

#### External Communication

**Status Page Updates:**
- Immediate: "Investigating issues with Listings service"
- Update every 15 minutes
- Resolution: "Service restored, monitoring for stability"

**Customer Notifications:**
- Email affected users (if data loss)
- Provide timeline and resolution
- Offer compensation if applicable

### Escalation Path

1. **Level 1:** Engineer on duty (immediate)
2. **Level 2:** Team Lead (if unresolved in 15 min)
3. **Level 3:** Engineering Manager (if unresolved in 30 min)
4. **Level 4:** CTO (if customer-facing outage >1 hour)

### On-Call Responsibilities

**During Rollback:**
- Monitor service health
- Respond to alerts
- Communicate status
- Execute rollback procedures
- Document incident

**After Rollback:**
- Verify service stability
- Analyze root cause
- Create tickets
- Update documentation
- Conduct postmortem

---

## Testing Rollback Procedures

### Staging Environment Test

**Quarterly:** Test rollback in staging

```bash
# 1. Deploy to staging
./scripts/deploy-to-staging.sh

# 2. Wait for deployment to complete

# 3. Test rollback
./scripts/rollback-staging.sh --reason "testing"

# 4. Verify rollback worked

# 5. Document any issues
```

### Dry-Run Test

```bash
# Test rollback without making changes
./scripts/rollback-prod.sh --dry-run --reason "testing"
```

### Checklist

- [ ] Rollback script executes successfully
- [ ] Traffic switches to Blue
- [ ] Green stops cleanly
- [ ] Blue serves traffic correctly
- [ ] Logs captured
- [ ] Report generated
- [ ] Notifications sent

---

## Rollback Checklist

### Pre-Rollback

- [ ] Severity assessed (P0-P3)
- [ ] Rollback type determined (traffic/full/partial)
- [ ] Team notified in Slack
- [ ] Incident commander assigned
- [ ] Status page updated

### During Rollback

- [ ] Rollback script executed
- [ ] Traffic switched to Blue
- [ ] Green instance stopped
- [ ] Blue verified healthy
- [ ] Database restored (if needed)
- [ ] Logs captured

### Post-Rollback

- [ ] Service verified functional
- [ ] Smoke tests passed
- [ ] Metrics normal
- [ ] Customers notified (if affected)
- [ ] Incident report reviewed
- [ ] Root cause investigation started
- [ ] Bug ticket created
- [ ] Postmortem scheduled

---

## Common Rollback Scenarios

### Scenario 1: Deployment Panic

**Symptoms:** Application crashes immediately

**Rollback:**
```bash
./scripts/rollback-prod.sh --reason "application panic on startup"
```

**Investigation:**
- Check Green logs for panic stacktrace
- Review recent code changes
- Test locally with production config

### Scenario 2: Performance Regression

**Symptoms:** 3x slower response times

**Rollback:**
```bash
./scripts/rollback-prod.sh --reason "performance regression"
```

**Investigation:**
- Review database query performance
- Check for N+1 queries
- Profile code with production data

### Scenario 3: Data Corruption

**Symptoms:** Invalid data in database

**Rollback:**
```bash
./scripts/rollback-prod.sh --reason "data corruption" --restore-db
```

**Investigation:**
- Review migration scripts
- Check for race conditions
- Verify data validation logic

### Scenario 4: Integration Failure

**Symptoms:** Third-party API errors

**Rollback:**
```bash
# May not need rollback, try config change first
ssh user@production-server
nano /opt/listings/.env.green
# Update API endpoint/key
pkill -f listings-green && nohup ./bin/listings-green --port=8081 --env-file=.env.green > logs/green.log 2>&1 &
```

**Investigation:**
- Verify third-party API status
- Check credentials/permissions
- Review API version compatibility

---

## Prevention

### Before Deployment

- [ ] Code reviewed
- [ ] Tests pass (unit + integration)
- [ ] Staging deployment successful
- [ ] Rollback plan documented
- [ ] Team available for support

### During Deployment

- [ ] Monitor logs in real-time
- [ ] Watch metrics dashboards
- [ ] Have rollback command ready
- [ ] Keep communication channels open

### After Deployment

- [ ] Monitor for 24 hours
- [ ] Review metrics trends
- [ ] Check error rates
- [ ] Gather user feedback

---

## Related Documents

- [Deployment Guide](./DEPLOYMENT.md)
- [Runbook](./RUNBOOK.md)
- [Incident Response](./INCIDENT_RESPONSE.md)

---

## Emergency Contacts

**On-Call Engineer:** See `.env.deploy` for current rotation

**Escalation:**
- Team Lead: [contact]
- Engineering Manager: [contact]
- CTO: [contact]

**External:**
- Database Admin: [contact]
- DevOps Team: [contact]
- Security Team: [contact]

---

*Last updated: 2025-11-05*
*Maintained by: DevOps Team*
