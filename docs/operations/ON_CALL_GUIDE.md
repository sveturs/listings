# Listings Microservice On-Call Guide

**Last Updated:** 2025-11-05
**Version:** 1.0.0
**Maintainer:** Platform Team

## Table of Contents

- [On-Call Overview](#on-call-overview)
- [Responsibilities](#responsibilities)
- [Getting Started](#getting-started)
- [Communication Channels](#communication-channels)
- [Escalation Matrix](#escalation-matrix)
- [Alert Response](#alert-response)
- [Common Alerts](#common-alerts)
- [After-Hours Procedures](#after-hours-procedures)
- [Incident Severity Levels](#incident-severity-levels)
- [Handoff Procedures](#handoff-procedures)
- [Post-Incident Procedures](#post-incident-procedures)
- [On-Call Toolkit](#on-call-toolkit)
- [Tips and Best Practices](#tips-and-best-practices)

---

## On-Call Overview

### What is On-Call?

As the on-call engineer for the Listings microservice, you are the first responder for all alerts and incidents. Your goal is to:
1. **Respond quickly** to alerts (within 5 minutes)
2. **Assess severity** and impact
3. **Resolve** or escalate appropriately
4. **Communicate** status to stakeholders
5. **Document** actions taken

### On-Call Rotation

- **Schedule:** Weekly rotation, Monday 9:00 AM to Monday 9:00 AM (local time)
- **Managed via:** PagerDuty
- **Rotation size:** 4-6 engineers
- **Backup:** Secondary on-call engineer
- **View schedule:** https://svetu.pagerduty.com/schedules

### Expected Response Times

| Alert Severity | Acknowledgment | Initial Assessment | Regular Updates |
|----------------|----------------|-------------------|-----------------|
| **P1 (Critical)** | 5 minutes | 10 minutes | Every 30 minutes |
| **P2 (High)** | 15 minutes | 30 minutes | Every hour |
| **P3 (Medium)** | 1 hour | 2 hours | Every 4 hours |
| **P4 (Low)** | Next business day | N/A | N/A |

### Compensation

- **Weekday on-call:** [Company policy]
- **Weekend on-call:** [Company policy]
- **Incident response:** [Company policy]

---

## Responsibilities

### Primary Responsibilities

1. **Monitor alerts** from PagerDuty and Slack
2. **Respond to pages** within SLA timeframes
3. **Triage incidents** and determine severity
4. **Take action** to resolve or mitigate
5. **Escalate** when necessary
6. **Communicate** with stakeholders
7. **Document** all actions and decisions
8. **Hand off** unresolved incidents properly

### Secondary Responsibilities

1. Review and improve runbooks
2. Conduct postmortems for incidents
3. Identify and file bugs
4. Suggest monitoring improvements
5. Participate in on-call training

### Out of Scope

❌ **NOT your responsibility:**
- Feature development
- Code reviews (unless blocking critical fix)
- Attending regular meetings during on-call
- Working on non-urgent tickets
- Supporting other services (unless explicitly assigned)

✅ **Focus on:**
- Service availability
- Incident response
- Critical bug fixes
- Emergency deployments

---

## Getting Started

### Before Your Shift

**1 week before:**
```bash
# Review the schedule
# Confirm no conflicts (vacation, important meetings)
# Arrange swap if needed via PagerDuty
```

**1 day before:**
```bash
# Review recent incidents
https://svetu.pagerduty.com/incidents

# Check current service health
# Grafana: https://grafana.svetu.rs/d/listings-overview
# Prometheus: http://prometheus.svetu.rs:9090

# Read this guide
# Review RUNBOOK.md and TROUBLESHOOTING.md
```

**Start of shift:**
```bash
# Test PagerDuty
# Mobile app: Trigger test alert
# Verify: Email, SMS, phone call working

# Test SSH access
ssh svetu@svetu.rs

# Test database access
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"

# Test Redis access
redis-cli -h localhost -p 36380 -a redis_password ping

# Check service health
curl http://localhost:8086/health

# Review open incidents
https://svetu.pagerduty.com/incidents

# Introduce yourself in Slack
# #listings-team: "On-call for Listings this week. Feel free to reach out!"
```

### Access Requirements

**Before starting on-call, ensure you have:**

- [ ] PagerDuty account configured (mobile app installed)
- [ ] SSH access to production servers
- [ ] VPN credentials (if required)
- [ ] Database credentials (stored in password manager)
- [ ] Grafana access (monitoring dashboards)
- [ ] Prometheus access (metrics)
- [ ] Slack access (all relevant channels)
- [ ] GitHub access (repository for emergency fixes)
- [ ] AWS Console access (if applicable)
- [ ] sudo privileges on production servers

**Test all access before your shift!**

### Development Environment

Set up a local development environment for emergency fixes:

```bash
# Clone repository
git clone https://github.com/sveturs/listings.git
cd listings

# Set up development environment
cp .env.example .env
# Edit .env with development values

# Install dependencies
make deps

# Run tests to verify setup
make test

# Build to verify toolchain
make build
```

---

## Communication Channels

### PagerDuty

**Primary alert channel**
- Mobile app: iOS/Android
- SMS alerts
- Phone calls (escalation)
- Email (low priority)

**Actions:**
- **Acknowledge:** "I'm looking into this"
- **Escalate:** Pass to next level
- **Resolve:** Issue fixed
- **Snooze:** Temporary suppression (use carefully!)

### Slack Channels

| Channel | Purpose | When to Use |
|---------|---------|-------------|
| **#listings-incidents** | Active incident coordination | All P1/P2 incidents |
| **#listings-alerts** | Alert notifications | Monitoring alerts |
| **#listings-team** | Team communication | Questions, non-urgent updates |
| **#platform-team** | Cross-team coordination | Escalation, help needed |
| **#security-incidents** | Security issues | Security breaches, suspicious activity |
| **#engineering-oncall** | All on-call engineers | General on-call support |

**Slack Slash Commands:**
```
/incident create severity:high summary:"Listings service down"
/pagerduty trigger
/statuspage update
```

### Email

**Use for:**
- Stakeholder updates (non-urgent)
- Postmortem distribution
- Weekly on-call summaries

**Mailing Lists:**
- `oncall@svetu.rs` - On-call engineers
- `platform@svetu.rs` - Platform team
- `engineering@svetu.rs` - All engineering
- `incidents@svetu.rs` - Incident notifications

### Phone

**Reserved for:**
- P1 incidents (if no PagerDuty response)
- Escalation to management
- Emergency communication

**Phone Tree:**
1. On-Call Engineer (PagerDuty)
2. Backup On-Call (PagerDuty)
3. Platform Team Lead (direct)
4. VP Engineering (direct)

---

## Escalation Matrix

### When to Escalate

**Escalate immediately if:**
- You don't have necessary access
- Issue is outside your expertise
- Standard procedures aren't working
- Severity is P1 and no resolution in 15 minutes
- Data loss or corruption suspected
- Security breach confirmed
- Multiple systems affected

**DO NOT hesitate to escalate!** Better to escalate early than miss SLO.

### Escalation Levels

#### Level 1: On-Call Engineer (You)
**Handles:** Most alerts, standard incidents, known issues
**Contact:** Via PagerDuty
**Expected resolution:** 80% of incidents

#### Level 2: Backup On-Call Engineer
**Handles:** Complex issues, second opinion, coverage
**Contact:** Via PagerDuty (escalate button)
**When:** After 15 minutes with no progress, or need help

#### Level 3: Platform Team Lead
**Handles:** Cross-service issues, architectural decisions, P1 incidents
**Contact:** PagerDuty `platform-team-lead` or direct phone
**When:** P1 incidents, escalation from Level 2, policy decisions

#### Level 4: VP Engineering
**Handles:** Major outages, business impact decisions, executive communication
**Contact:** Direct phone (after informing Platform Team Lead)
**When:** Complete service loss, data breach, media attention

### Specialist Escalations

**Database SRE Team:**
- PostgreSQL issues
- Slow queries, connection problems
- Data corruption
**Contact:** PagerDuty `db-team`

**Security Team:**
- Security breaches
- Suspicious activity
- Credential compromise
**Contact:** PagerDuty `security-team` (any time)

**Infrastructure Team:**
- Server/VM issues
- Network problems
- DNS issues
**Contact:** PagerDuty `infra-team`

**OpenSearch Team:**
- Search cluster issues
- Index problems
**Contact:** Slack `#opensearch-support`

### Escalation Template

When escalating, provide:
```
ESCALATION: Listings Service [ISSUE]

SEVERITY: P1/P2/P3
DURATION: [HH:MM]
CURRENT STATUS: [Brief description]

ACTIONS TAKEN:
- [Action 1]
- [Action 2]

WHY ESCALATING:
- [Reason 1]
- [Reason 2]

ASSISTANCE NEEDED:
- [Specific help required]

INCIDENT LINK: [PagerDuty URL]
SLACK THREAD: [Link]
```

---

## Alert Response

### Step-by-Step Response Process

#### Step 1: Acknowledge (1 minute)

```bash
# In PagerDuty app:
# 1. Tap notification
# 2. Tap "Acknowledge"
# 3. Optionally add note: "Investigating"

# In Slack #listings-incidents:
# Post: "Acknowledged alert: [ALERT_NAME]. Investigating."
```

#### Step 2: Assess Severity (2-5 minutes)

**Check service status:**
```bash
# SSH to server
ssh svetu@svetu.rs

# Check service
sudo systemctl status listings-service

# Check health endpoints
curl http://localhost:8086/health
curl http://localhost:8086/ready

# Check metrics
curl -s http://localhost:8086/metrics | grep -E 'listings_errors_total|listings_grpc_requests_total'

# Check recent logs
sudo journalctl -u listings-service --since "10 minutes ago" | tail -50
```

**Check dependencies:**
```bash
# PostgreSQL
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"

# Redis
redis-cli -h localhost -p 36380 -a redis_password ping

# OpenSearch
curl -u admin:admin http://localhost:9200/_cluster/health
```

**Determine severity:**
- **P1:** Service down, all users affected, data loss
- **P2:** Major functionality impaired, many users affected
- **P3:** Degraded performance, some users affected
- **P4:** Minor issue, no user impact

#### Step 3: Communicate (2 minutes)

**For P1/P2:**
```bash
# Create incident in Slack
/incident create severity:high summary:"Listings [ISSUE]"

# Post initial update
"INVESTIGATING: Listings service [issue]. Checking [components].
Will update in 15 minutes."

# Update PagerDuty incident
# Add note with initial findings
```

**For P3/P4:**
```bash
# Post in #listings-alerts
"Investigating [alert]. Will update when resolved."
```

#### Step 4: Investigate and Resolve (varies)

**Use RUNBOOK.md:**
- Find matching incident
- Follow resolution steps
- Document actions taken

**If not in runbook:**
- Use TROUBLESHOOTING.md decision tree
- Check logs for error patterns
- Review recent changes (deployments, config)
- Check Grafana for metric anomalies

**Common quick wins:**
```bash
# Restart service (if safe)
sudo systemctl restart listings-service

# Clear cache
redis-cli -h localhost -p 36380 -a redis_password FLUSHALL

# Kill long-running queries
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pg_terminate_backend(pid) FROM pg_stat_activity
   WHERE state='active' AND query_start < NOW() - INTERVAL '30 seconds';"
```

#### Step 5: Verify Resolution (5 minutes)

```bash
# Health checks
curl http://localhost:8086/health
curl http://localhost:8086/ready

# Test functionality
grpcurl -plaintext -d '{"id": 1}' localhost:50053 listings.v1.ListingsService/GetListing

# Monitor metrics for 5 minutes
watch -n 10 'curl -s http://localhost:8086/metrics | grep listings_errors_total'

# Check alert status
# Should stop firing if resolved
```

#### Step 6: Resolve and Document (5 minutes)

**In PagerDuty:**
```
Resolution Note:
- Issue: [Brief description]
- Root cause: [Cause]
- Resolution: [Action taken]
- Duration: [HH:MM]
- Impact: [Description]
```

**In Slack:**
```
RESOLVED: Listings [issue] resolved.
- Duration: [HH:MM]
- Root cause: [Brief description]
- Resolution: [Action taken]
- Follow-up: [Ticket number if needed]
```

**Create follow-up ticket:**
- If root cause needs fixing
- If runbook needs updating
- If monitoring needs improvement

---

## Common Alerts

### ListingsServiceDown

**Meaning:** Service health check failing
**Severity:** P1
**Response time:** 5 minutes

**Quick diagnosis:**
```bash
sudo systemctl status listings-service
curl http://localhost:8086/health
sudo journalctl -u listings-service -n 50
```

**Common causes:**
- Service crashed (check logs for panic)
- Database connection lost
- Configuration error

**Resolution:** See RUNBOOK.md → Service Down

---

### ListingsHighErrorRate

**Meaning:** Error rate > 1%
**Severity:** P2
**Response time:** 15 minutes

**Quick diagnosis:**
```bash
curl -s http://localhost:8086/metrics | grep listings_errors_total
sudo journalctl -u listings-service --since "10 minutes ago" | grep '"level":"error"' | jq .
```

**Common causes:**
- Database connection pool exhausted
- External service down (OpenSearch)
- Bad deployment

**Resolution:** See RUNBOOK.md → High Error Rate

---

### ListingsHighLatency

**Meaning:** P99 latency > 2 seconds
**Severity:** P2
**Response time:** 15 minutes

**Quick diagnosis:**
```bash
curl -s http://localhost:8086/metrics | grep listings_grpc_request_duration_seconds
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pid, now() - query_start AS duration, query FROM pg_stat_activity
   WHERE state='active' ORDER BY duration DESC LIMIT 10;"
```

**Common causes:**
- Slow database queries
- Database connection pool full
- High CPU usage

**Resolution:** See RUNBOOK.md → High Latency

---

### ListingsDBPoolExhausted

**Meaning:** Database connection pool at capacity
**Severity:** P2
**Response time:** 15 minutes

**Quick diagnosis:**
```bash
curl -s http://localhost:8086/metrics | grep listings_db_connections
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT count(*), state FROM pg_stat_activity WHERE datname='listings_db' GROUP BY state;"
```

**Resolution:** See RUNBOOK.md → Database Connection Pool Exhausted

---

### ListingsRedisDown

**Meaning:** Redis connection failing
**Severity:** P3 (service continues with fail-open)
**Response time:** 1 hour

**Quick diagnosis:**
```bash
redis-cli -h localhost -p 36380 -a redis_password ping
sudo systemctl status redis
```

**Impact:**
- Rate limiting disabled (fail-open)
- Caching disabled (higher DB load)

**Resolution:** See RUNBOOK.md → Redis Connection Issues

---

### OpenSearchClusterRed

**Meaning:** OpenSearch cluster in red state
**Severity:** P3 (search unavailable, CRUD continues)
**Response time:** 1 hour

**Quick diagnosis:**
```bash
curl -u admin:admin http://localhost:9200/_cluster/health
curl -u admin:admin http://localhost:9200/_cat/indices | grep listings
```

**Impact:**
- Search functionality unavailable
- CRUD operations continue normally

**Resolution:** See RUNBOOK.md → OpenSearch Cluster Red

---

### ListingsHighMemoryUsage

**Meaning:** Memory usage > 80%
**Severity:** P3 (warning, may lead to OOM)
**Response time:** 1 hour

**Quick diagnosis:**
```bash
ps aux | grep listings-service | awk '{print $6/1024 " MB"}'
curl http://localhost:8086/debug/pprof/heap > /tmp/heap.prof
go tool pprof -top /tmp/heap.prof
```

**Resolution:** See RUNBOOK.md → Memory Leak / OOM Killed

---

### ListingsDiskSpaceLow

**Meaning:** Disk usage > 80%
**Severity:** P2 (may lead to service failure)
**Response time:** 30 minutes

**Quick diagnosis:**
```bash
df -h
du -sh /opt/listings-dev/*
du -sh /var/log/journal/*
```

**Resolution:** See RUNBOOK.md → Disk Space Critical

---

### ListingsRateLimitAbuse

**Meaning:** High rate limit rejections
**Severity:** P3 (may indicate DDoS or legitimate traffic spike)
**Response time:** 1 hour

**Quick diagnosis:**
```bash
curl -s http://localhost:8086/metrics | grep listings_rate_limit_rejected_total
sudo journalctl -u listings-service --since "10 minutes ago" | \
  grep "rate limit exceeded" | jq -r '.identifier' | sort | uniq -c | sort -rn
```

**Resolution:** See RUNBOOK.md → Rate Limit Abuse

---

### ListingsSlowQueries

**Meaning:** Database queries taking > 1 second
**Severity:** P3
**Response time:** 1 hour

**Quick diagnosis:**
```bash
sudo tail -100 /var/log/postgresql/postgresql-15-main.log | grep "duration:"
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pid, now() - query_start, query FROM pg_stat_activity
   WHERE state='active' AND now() - query_start > interval '1 second';"
```

**Resolution:** See RUNBOOK.md → Slow Queries

---

## After-Hours Procedures

### Responding After Hours

**Upon receiving page:**
1. **Acknowledge within 5 minutes** (even if you need time to get to computer)
2. **Post in Slack** that you're investigating
3. **Assess severity** quickly
4. **If P1:** Begin immediate response
5. **If P2/P3:** Assess if it can wait until morning (document decision)
6. **If P4:** Snooze until morning (should be rare after hours)

### Decision: Respond Now vs Morning?

**Respond immediately if:**
- ✅ P1 (service down, data loss)
- ✅ P2 with high user impact
- ✅ Security incident
- ✅ Ongoing data corruption
- ✅ SLO breach imminent

**Can wait until morning if:**
- ✅ P3 with no user-facing impact
- ✅ Alert is informational (monitoring suggestion)
- ✅ Issue is stable (not getting worse)
- ✅ Workaround is in place

**Document your decision:**
```
PagerDuty note: "Assessed as P3, no immediate user impact.
Monitoring continues. Will address first thing in morning."

Slack: "Reviewed alert, currently stable. Monitoring overnight.
Will investigate in detail tomorrow morning."
```

### Working from Home Setup

**Before your on-call week:**
```bash
# Test VPN from home
# Test SSH from home network
# Test laptop has all necessary tools installed
# Test PagerDuty on mobile AND laptop
# Have phone charger accessible
# Have laptop charger accessible
# Inform family/roommates you're on-call
```

**If you can't access from home:**
- Acknowledge alert immediately
- Escalate to backup on-call or Platform Team Lead
- Provide context: "Unable to access from home, escalating"

### Weekend Coverage

**Weekend on-call expectations:**
- Check service health once per day (morning or evening)
- Respond to pages with same SLA
- Can request backup coverage for specific times (with advance notice)

**Friday handoff:**
- Review any ongoing issues
- Check for scheduled maintenance
- Ensure backup on-call is aware of current state

**Monday handoff:**
- Document weekend incidents
- Brief incoming on-call
- Create tickets for follow-up work

---

## Incident Severity Levels

### P1 - Critical

**Criteria:**
- Complete service outage
- Data loss or corruption
- Security breach
- All users affected
- Revenue-impacting

**Response:**
- Acknowledge: **5 minutes**
- Initial assessment: **10 minutes**
- Updates: **Every 30 minutes**
- All hands on deck if needed
- Escalate to Platform Team Lead immediately

**Communication:**
- Create Slack incident channel
- Page backup on-call
- Update status page
- Notify stakeholders every 30 min

**Examples:**
- Service completely down
- Database loss
- Security breach

---

### P2 - High

**Criteria:**
- Major functionality impaired
- Many users affected
- Degraded performance (P99 >5s)
- Error rate >5%
- Partial outage

**Response:**
- Acknowledge: **15 minutes**
- Initial assessment: **30 minutes**
- Updates: **Every hour**
- Escalate if no progress in 1 hour

**Communication:**
- Post in #listings-incidents
- Update PagerDuty incident
- Notify team lead

**Examples:**
- High error rate
- Database connection issues
- Significant latency increase

---

### P3 - Medium

**Criteria:**
- Degraded performance
- Some users affected
- Non-critical functionality impaired
- Warning threshold reached

**Response:**
- Acknowledge: **1 hour**
- Initial assessment: **2 hours**
- Updates: **Every 4 hours**

**Communication:**
- Post in #listings-alerts
- Update PagerDuty

**Examples:**
- Redis down (fail-open active)
- OpenSearch yellow
- Elevated latency
- Rate limit warnings

---

### P4 - Low

**Criteria:**
- Informational alerts
- No user impact
- Monitoring suggestions

**Response:**
- Next business day
- Create ticket for follow-up

**Communication:**
- Internal only

**Examples:**
- Disk space warning (60% used)
- Maintenance notifications
- Monitoring config suggestions

---

## Handoff Procedures

### Daily Handoff (During Shifts)

**At end of your day (even if incidents ongoing):**

```bash
# Post in #listings-team:
"Daily on-call handoff:

ACTIVE INCIDENTS:
- [Incident 1]: [Status]
- [Incident 2]: [Status]

MONITORING:
- [Metric] is elevated, watching closely
- [Alert] keeps firing, investigating

FOLLOW-UP NEEDED:
- [Ticket] for [issue]

Backup on-call tonight: @backup-engineer"
```

### Weekly Handoff

**End of your on-call week:**

1. **Schedule 15-minute sync** with next on-call
2. **Review open incidents** together
3. **Discuss any trends** you noticed
4. **Share lessons learned**
5. **Provide context** on ongoing investigations

**Handoff template:**
```markdown
# On-Call Handoff: Week of [DATE]

## Summary
- Total alerts: [COUNT]
- Incidents: [COUNT] (P1: X, P2: Y, P3: Z)
- Escalations: [COUNT]
- Service uptime: [%]

## Active Issues
1. [Issue 1]
   - Status: [In progress/Monitoring]
   - Actions taken: [List]
   - Next steps: [List]

2. [Issue 2]
   - ...

## Trends Noticed
- [Observation 1]
- [Observation 2]

## Improvements Made
- [Runbook update]
- [Monitoring added]

## Open Tickets
- [TICKET-1]: [Description]
- [TICKET-2]: [Description]

## Notes for Next On-Call
- [Important info]
- [Watch for...]

## Questions?
Feel free to ping me on Slack!
```

---

## Post-Incident Procedures

### Immediately After Resolving

**1. Update systems (5 minutes):**
```bash
# Resolve PagerDuty incident with note
# Post resolution in Slack
# Update status page (if used)
# Close Slack incident channel (if P1/P2)
```

**2. Create follow-up ticket (5 minutes):**
```markdown
Title: [COMPONENT] Follow-up from incident [INCIDENT-ID]

Description:
Incident occurred on [DATE] at [TIME].

Summary:
- [Brief description]

Root cause:
- [If known]

Follow-up needed:
- [ ] Investigate root cause (if unknown)
- [ ] Fix underlying issue
- [ ] Update runbook
- [ ] Add/improve monitoring
- [ ] Add automated remediation

Links:
- PagerDuty: [URL]
- Slack thread: [URL]
- Grafana: [URL]
```

### Within 24 Hours

**3. Write incident summary:**
```markdown
## Incident Summary: [INCIDENT-ID]

**Date:** [DATE]
**Duration:** [HH:MM]
**Severity:** [P1/P2/P3]
**On-Call:** [Your name]

**Impact:**
- Users affected: [COUNT or description]
- Functionality impacted: [Description]

**Timeline:**
- HH:MM - Alert fired
- HH:MM - Acknowledged
- HH:MM - Root cause identified
- HH:MM - Fix applied
- HH:MM - Verified resolved

**Root Cause:**
[Description]

**Resolution:**
[What was done to fix it]

**Prevention:**
[What can be done to prevent recurrence]

**Action Items:**
- [TICKET-1]: [Description]
```

### Within 48 Hours (for P1/P2)

**4. Schedule postmortem:**
- Invite: On-call engineer, team lead, relevant stakeholders
- Duration: 1 hour
- Agenda: Timeline review, root cause, improvements

**5. Update documentation:**
- Add to RUNBOOK.md if new scenario
- Update TROUBLESHOOTING.md with new findings
- Improve alert descriptions based on experience

---

## On-Call Toolkit

### Required Tools

**Local machine:**
```bash
# SSH client
ssh -V

# psql client
psql --version

# redis-cli
redis-cli --version

# curl
curl --version

# jq (JSON parsing)
jq --version

# grpcurl (gRPC testing)
grpcurl --version

# Git
git --version

# Go toolchain (for emergency builds)
go version
```

**Mobile apps:**
- PagerDuty
- Slack
- SSH client (Termius, Blink, etc.)
- VPN client

### Useful Commands Cheatsheet

**Save to `/opt/listings-on-call-commands.sh`:**
```bash
#!/bin/bash
# Listings On-Call Quick Commands

# Service status
alias ls-status='sudo systemctl status listings-service'
alias ls-restart='sudo systemctl restart listings-service'
alias ls-logs='sudo journalctl -u listings-service -f'

# Health checks
alias ls-health='curl http://localhost:8086/health'
alias ls-ready='curl http://localhost:8086/ready'
alias ls-metrics='curl -s http://localhost:8086/metrics | grep listings_'

# Database
alias ls-db='psql "postgres://listings_user:listings_password@localhost:35433/listings_db"'
alias ls-db-conns='psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT count(*), state FROM pg_stat_activity WHERE datname=\"listings_db\" GROUP BY state;"'

# Redis
alias ls-redis='redis-cli -h localhost -p 36380 -a redis_password'
alias ls-redis-ping='redis-cli -h localhost -p 36380 -a redis_password ping'

# OpenSearch
alias ls-os-health='curl -u admin:admin http://localhost:9200/_cluster/health'
alias ls-os-indices='curl -u admin:admin http://localhost:9200/_cat/indices | grep listings'

# Quick diagnostics
alias ls-errors='sudo journalctl -u listings-service --since "10 minutes ago" | grep '"'"'"level":"error"'"'"' | jq .'
alias ls-error-rate='curl -s http://localhost:8086/metrics | grep listings_errors_total'
alias ls-db-slow='psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT pid, now() - query_start, query FROM pg_stat_activity WHERE state=\"active\" AND now() - query_start > interval \"1 second\";"'
```

### Bookmarks

**Essential URLs:**
- PagerDuty: https://svetu.pagerduty.com
- Grafana: https://grafana.svetu.rs/d/listings-overview
- Prometheus: http://prometheus.svetu.rs:9090
- GitHub: https://github.com/sveturs/listings
- RUNBOOK: `file:///p/github.com/sveturs/listings/docs/operations/RUNBOOK.md`
- Slack: https://svetu.slack.com/messages/listings-incidents

### Emergency Contacts

**Save in phone:**
```
Listings Backup On-Call: [Phone via PagerDuty]
Platform Team Lead: [Phone]
Database SRE: [Phone via PagerDuty]
VP Engineering: [Phone]
```

---

## Tips and Best Practices

### Do's ✅

1. **Acknowledge quickly** - Even if you need time to investigate
2. **Communicate often** - Stakeholders prefer over-communication
3. **Document everything** - Future you (and others) will thank you
4. **Ask for help early** - Don't wait until you're stuck
5. **Use runbooks** - They're tested and reliable
6. **Test before applying** - Especially for fixes that modify data
7. **Monitor after fixing** - Ensure issue doesn't recur
8. **Take breaks** - During long incidents, take 5-minute breaks
9. **Stay calm** - Panic doesn't help anyone
10. **Learn and improve** - Use incidents as learning opportunities

### Don'ts ❌

1. **Don't ignore alerts** - They exist for a reason
2. **Don't apply fixes without understanding** - You might make it worse
3. **Don't skip communication** - Silent incident response is scary
4. **Don't be afraid to escalate** - It's better to escalate early
5. **Don't restart services blindly** - Understand the problem first
6. **Don't delete data without backups** - Even if you think it's safe
7. **Don't work on incidents while impaired** - If you're unable (sick, intoxicated), escalate
8. **Don't skip documentation** - It helps future responders
9. **Don't take it personally** - Incidents happen, it's not your fault
10. **Don't burn out** - Ask for relief if needed

### Mental Model for Incidents

**Think like a doctor:**
1. **Triage:** How severe? Who's affected?
2. **Diagnose:** What's the root cause?
3. **Stabilize:** Stop the bleeding (mitigation)
4. **Treat:** Fix the root cause
5. **Monitor:** Ensure recovery is complete
6. **Prevent:** How do we avoid this in future?

**Remember:**
- Service availability > perfect fixes
- Communication > silent investigation
- Mitigation > complete resolution (initially)
- Documented guess > undocumented certainty

### Stress Management

**During high-stress incidents:**
- Take deep breaths
- Ask for backup if overwhelmed
- Focus on one step at a time
- Don't rush, even under pressure
- Remember: This will end

**After stressful incidents:**
- Take a break
- Talk to team lead or colleagues
- Document lessons learned
- Don't dwell on mistakes
- Celebrate wins (resolution!)

### Work-Life Balance

**Protecting your time:**
- Set clear boundaries with team
- Swap shifts if you have important events
- Take comp time after major incidents
- Don't check PagerDuty when not on-call
- Turn off work Slack when off-duty

**If on-call is impacting your wellbeing:**
- Talk to your manager
- Consider rotation adjustments
- Seek support from team

---

## Frequently Asked Questions

**Q: What if I'm unsure if something is an incident?**
A: When in doubt, treat it as P3 and investigate. Better safe than sorry.

**Q: Can I make changes to production during on-call?**
A: Yes, for incident response. For non-urgent changes, follow change management process.

**Q: What if I'm not available during my on-call shift?**
A: Arrange a swap via PagerDuty **before** your shift. Last-minute unavailability should be rare and coordinated with backup.

**Q: What if an alert fires but everything looks fine?**
A: Check if it's a false positive, adjust alert if needed, document, and resolve.

**Q: Should I wake up the Platform Team Lead for P2 at 3 AM?**
A: Only if you can't resolve it or need help. P1 always escalate immediately.

**Q: What if I make a mistake during incident response?**
A: Communicate it immediately, document it, work to fix it. Mistakes happen. We learn from them.

**Q: Can I ignore alerts if I'm about to end my shift?**
A: No. Respond and hand off properly to next on-call if needed.

**Q: How do I know if I should escalate?**
A: If you've been investigating for 30 minutes (P1) or 1 hour (P2) without progress, escalate.

---

## On-Call Checklist

### Start of Week
- [ ] Test PagerDuty (mobile app, SMS, phone)
- [ ] Test SSH access to production
- [ ] Test database access
- [ ] Review recent incidents
- [ ] Check current service health
- [ ] Review RUNBOOK.md updates
- [ ] Post in #listings-team that you're on-call

### Daily
- [ ] Check service health in morning
- [ ] Review overnight alerts (if any)
- [ ] Monitor #listings-alerts channel
- [ ] Keep laptop and phone charged
- [ ] Stay reachable

### After Incident
- [ ] Resolve PagerDuty incident
- [ ] Post resolution in Slack
- [ ] Create follow-up ticket
- [ ] Write incident summary
- [ ] Update documentation (if needed)

### End of Week
- [ ] Prepare handoff document
- [ ] Schedule handoff meeting with next on-call
- [ ] File any outstanding tickets
- [ ] Update runbooks based on week's learnings
- [ ] Post weekly summary in #listings-team

---

## Resources

- **RUNBOOK.md** - Common incidents and resolutions
- **TROUBLESHOOTING.md** - Debugging guide
- **DISASTER_RECOVERY.md** - Emergency procedures
- **MONITORING_GUIDE.md** - Grafana dashboard tour
- **SLO_GUIDE.md** - SLO tracking and error budgets

**External Resources:**
- PagerDuty Best Practices: https://www.pagerduty.com/resources/learn/oncall-best-practices/
- Google SRE Book (On-Call): https://sre.google/sre-book/being-on-call/
- Incident Response: https://response.pagerduty.com/

---

## Feedback

Have suggestions for improving the on-call experience?
- Create ticket: [JIRA project]
- Post in: #listings-team
- Email: platform@svetu.rs

---

**Document Version:** 1.0.0
**Last Reviewed:** 2025-11-05
**Next Review:** 2025-12-05
**Owner:** Platform Team

**Good luck on your on-call rotation! Remember: You've got this, and you're not alone. The team is here to help.**
