# Listings Microservice Operations Documentation

**Last Updated:** 2025-11-05
**Version:** 1.0.0
**Maintainer:** Platform Team

---

## 📚 Overview

This directory contains comprehensive operational documentation for the Listings microservice. These documents are designed for on-call engineers, SREs, and anyone responsible for maintaining service reliability.

---

## 📖 Documentation Structure

### Core Documents

| Document | Purpose | Primary Audience | When to Use |
|----------|---------|-----------------|-------------|
| **[RUNBOOK.md](./RUNBOOK.md)** | Common incidents & resolutions | On-call engineers | During active incidents |
| **[TROUBLESHOOTING.md](./TROUBLESHOOTING.md)** | Systematic debugging guide | All engineers | When investigating issues |
| **[DISASTER_RECOVERY.md](./DISASTER_RECOVERY.md)** | Emergency recovery procedures | Senior engineers, Management | Major outages, data loss |
| **[ON_CALL_GUIDE.md](./ON_CALL_GUIDE.md)** | On-call handbook | On-call engineers | Before/during on-call shifts |
| **[SLO_GUIDE.md](./SLO_GUIDE.md)** | SLO tracking & management | Platform team, Leadership | Monthly reviews, planning |
| **[MONITORING_GUIDE.md](./MONITORING_GUIDE.md)** | Monitoring tools usage | All engineers | Dashboard usage, metrics |

---

## 🚨 Quick Start for On-Call Engineers

### First Time On-Call? Start Here:

1. **Read:** [ON_CALL_GUIDE.md](./ON_CALL_GUIDE.md) - Complete handbook
2. **Bookmark:** [RUNBOOK.md](./RUNBOOK.md) - Your first stop during incidents
3. **Familiarize:** [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - Decision trees and tools
4. **Access:** Ensure you have all required access (see ON_CALL_GUIDE.md)
5. **Test:** Verify SSH, database, and PagerDuty access

### During an Incident:

```
1. Acknowledge alert (5 min) → PagerDuty
2. Check RUNBOOK.md → Find matching incident
3. Follow resolution steps → Document actions
4. If stuck → Use TROUBLESHOOTING.md decision tree
5. Can't resolve? → Escalate per ON_CALL_GUIDE.md
6. After resolution → Document in post-incident
```

---

## 📋 Document Summaries

### RUNBOOK.md

**Top 10 Common Incidents with Resolution Procedures**

1. High Error Rate (>1%)
2. High Latency (P99 >2s)
3. Service Down
4. Database Connection Pool Exhausted
5. Redis Connection Issues
6. OpenSearch Cluster Red
7. Memory Leak / OOM Killed
8. Disk Space Critical
9. Rate Limit Abuse
10. Slow Queries

**Each incident includes:**
- ✅ Symptoms (what you see)
- ✅ Investigation steps (commands to run)
- ✅ Resolution steps (how to fix)
- ✅ Prevention measures
- ✅ Escalation path

**Quick Reference Commands:**
```bash
# Service status
sudo systemctl status listings-service

# Logs
sudo journalctl -u listings-service -f

# Health check
curl http://localhost:8086/health

# Metrics
curl http://localhost:8086/metrics | grep listings_
```

---

### TROUBLESHOOTING.md

**Systematic Debugging Guide**

**Key Sections:**
- **Decision Tree:** Flowchart for problem diagnosis
- **Debugging Tools:** Logs, metrics, profiling, tracing
- **Common Errors:** Database, Redis, OpenSearch, gRPC
- **Performance Investigation:** Slow requests, query analysis, memory/CPU profiling
- **Network Issues:** Connectivity testing, DNS, certificates

**Tools Covered:**
- journalctl (log analysis)
- Prometheus/Grafana (metrics)
- pprof (CPU/memory profiling)
- psql (database queries)
- redis-cli (cache debugging)

**When to Use:**
- Investigation doesn't match any RUNBOOK scenario
- Need systematic approach to unknown issue
- Performance problems requiring deep analysis
- Learning how to debug service issues

---

### DISASTER_RECOVERY.md

**Emergency Recovery Procedures**

**5 Major Disaster Scenarios:**

1. **Complete Database Loss**
   - RTO: 30 minutes
   - RPO: 5 minutes
   - Full restoration from backup

2. **Region/Datacenter Failure**
   - RTO: 15 minutes
   - Failover procedures (future: multi-region)

3. **Data Corruption**
   - Point-in-time recovery
   - Partial restoration
   - Data integrity validation

4. **Security Breach**
   - Immediate containment
   - Evidence preservation
   - Credential rotation
   - System restoration

5. **Complete Service Loss**
   - Rapid redeployment
   - From source code
   - Database validation

**⚠️ CRITICAL:** Execute immediately upon catastrophic failure

**Includes:**
- Step-by-step recovery procedures
- Communication templates
- Backup/restore strategies
- Post-recovery validation

---

### ON_CALL_GUIDE.md

**Complete On-Call Handbook**

**Contents:**
- On-call responsibilities and expectations
- Alert response procedures
- Communication protocols
- Escalation matrix
- Common alerts with quick diagnosis
- After-hours procedures
- Handoff procedures
- Post-incident documentation
- On-call toolkit and commands
- Tips and best practices
- Mental health and stress management

**Key Information:**
- **Response Times:** P1: 5 min, P2: 15 min, P3: 1 hour
- **Escalation:** When and how to escalate
- **Severity Levels:** P1 (Critical) → P4 (Low)
- **Communication Channels:** PagerDuty, Slack, Email
- **Emergency Contacts:** Platform Lead, Database SRE, Security

---

### SLO_GUIDE.md

**SLO Management and Tracking**

**SLO Definitions:**
- **Availability:** 99.9% (43.2 min downtime/month)
- **Latency:** P95 < 1s, P99 < 2s
- **Error Rate:** < 1%

**Error Budget Management:**
- How to calculate remaining budget
- Error budget policy (feature freeze when depleted)
- Monthly reset process

**Tracking & Reporting:**
- Daily, weekly, monthly checks
- Dashboard usage
- Incident impact calculation
- SLO breach response

**Monthly Review Process:**
- Review schedule and agenda
- Template for review meetings
- Action item tracking
- SLO adjustment criteria

---

### MONITORING_GUIDE.md

**Monitoring Tools and Usage**

**Grafana Dashboards:**
1. **Listings Overview** - High-level health, SLO tracking
2. **Listings Details** - Deep-dive into metrics
3. **Database Performance** - PostgreSQL monitoring
4. **Redis Performance** - Cache and rate limiter
5. **SLO Dashboard** - Error budget tracking

**Prometheus Metrics:**
- gRPC handler metrics (requests, latency, errors)
- Database metrics (connections, queries)
- Rate limiting metrics
- Timeout metrics
- Business metrics (listings, inventory)
- Cache metrics
- Worker metrics

**Alert Definitions:**
- P1 (Critical): ServiceDown, CriticalErrorRate
- P2 (High): HighErrorRate, HighLatency, DBPoolExhausted
- P3 (Medium): ElevatedLatency, RedisDown, HighMemory

**Log Analysis:**
- Structured JSON logging format
- Common log queries with examples
- Finding errors, slow queries, tracking requests

---

## 🔧 Tools and Access

### Required Access

- **SSH:** Production servers
- **PagerDuty:** Alert management
- **Grafana:** Dashboards (https://grafana.svetu.rs)
- **Prometheus:** Metrics (http://prometheus.svetu.rs:9090)
- **Slack:** Communication channels
- **GitHub:** Repository access
- **Database:** psql credentials
- **Redis:** redis-cli credentials

### Essential URLs

| Service | URL | Purpose |
|---------|-----|---------|
| Grafana | https://grafana.svetu.rs | Monitoring dashboards |
| Prometheus | http://prometheus.svetu.rs:9090 | Metrics and queries |
| PagerDuty | https://svetu.pagerduty.com | Alert management |
| Slack | https://svetu.slack.com | Communication |
| GitHub | https://github.com/sveturs/listings | Source code |

---

## 📊 Service SLOs

| Metric | Target | Error Budget (Monthly) |
|--------|--------|----------------------|
| **Availability** | 99.9% | 43.2 minutes downtime |
| **P95 Latency** | < 1 second | - |
| **P99 Latency** | < 2 seconds | - |
| **Error Rate** | < 1% | 1% of all requests |

**Current Status:** See Grafana → Listings → SLO Dashboard

---

## 🚦 Incident Severity Levels

| Severity | Criteria | Response Time | Example |
|----------|----------|---------------|---------|
| **P1 - Critical** | Complete outage, data loss, security breach | 5 minutes | Service completely down |
| **P2 - High** | Major functionality impaired, many users affected | 15 minutes | High error rate (>1%) |
| **P3 - Medium** | Degraded performance, some users affected | 1 hour | Redis down (fail-open) |
| **P4 - Low** | Informational, no user impact | Next business day | Disk space warning (60%) |

---

## 📞 Emergency Contacts

| Role | Contact Method | When to Use |
|------|----------------|-------------|
| **On-Call Engineer** | PagerDuty: `listings-oncall` | All alerts |
| **Backup On-Call** | PagerDuty: escalate button | Need assistance |
| **Platform Team Lead** | PagerDuty: `platform-team-lead` | P1 incidents, escalation |
| **Database SRE** | PagerDuty: `db-team` | Database issues |
| **Security Team** | PagerDuty: `security-team` | Security incidents |

**Escalation Path:**
```
Level 1: On-Call Engineer
    ↓ (15 min for P1, 1 hour for P2)
Level 2: Backup On-Call or Platform Team Lead
    ↓ (if P1 or no progress)
Level 3: VP Engineering + CTO
    ↓ (data breach, major outage)
Level 4: CEO + Legal (security/compliance)
```

---

## 📝 Document Maintenance

### Review Schedule

| Document | Review Frequency | Next Review | Owner |
|----------|-----------------|-------------|-------|
| RUNBOOK.md | After each incident | Continuous | Platform Team |
| TROUBLESHOOTING.md | Monthly | 2025-12-05 | Platform Team |
| DISASTER_RECOVERY.md | Quarterly | 2025-12-05 | Platform Team + Security |
| ON_CALL_GUIDE.md | Quarterly | 2025-12-05 | Platform Team |
| SLO_GUIDE.md | Monthly | 2025-12-05 | Platform Team Lead |
| MONITORING_GUIDE.md | Monthly | 2025-12-05 | Platform Team |

### Contributing

**How to update these documents:**

1. **Identify improvement needed**
   - After incident: What was missing?
   - After on-call: What would have helped?
   - New procedures or tools added

2. **Create update branch**
   ```bash
   git checkout -b update-ops-docs-[topic]
   ```

3. **Make changes**
   - Update relevant document(s)
   - Add examples, commands, screenshots if helpful
   - Update "Last Updated" date

4. **Create Pull Request**
   - Title: "Update ops docs: [brief description]"
   - Description: Why this update is needed
   - Tag: Platform Team for review

5. **After merge**
   - Announce changes in #listings-team
   - Update any related training materials

**Quick fixes:**
- Typos, broken links: Direct commit to main (notify team)
- Command updates: Direct commit + notify team
- Process changes: Requires PR + review

---

## 🎓 Training and Onboarding

### New Engineer Onboarding

**Week 1:**
- [ ] Read all operations documents
- [ ] Get all required access
- [ ] Shadow on-call engineer
- [ ] Review recent incidents in PagerDuty

**Week 2:**
- [ ] Practice using RUNBOOK procedures in test environment
- [ ] Set up local development environment
- [ ] Run through disaster recovery scenarios
- [ ] Review Grafana dashboards

**Week 3:**
- [ ] Be backup on-call (with experienced engineer primary)
- [ ] Practice incident response procedures
- [ ] Write postmortem for practice incident

**Week 4+:**
- [ ] Ready for on-call rotation!

### Ongoing Training

**Monthly:**
- Review recent incidents and learnings
- Practice fire drills (simulate incidents)
- Share tips and tricks in team meetings

**Quarterly:**
- Disaster recovery drill
- Update documentation based on learnings
- Review and update SLOs

**Annually:**
- Comprehensive disaster recovery test
- Full documentation review and update
- On-call process retrospective

---

## 📊 Metrics and Reporting

### Service Health Dashboard

**Daily Check:** Grafana → Listings → Overview
- Service status: UP/DOWN
- Error rate: < 1%
- Latency P99: < 2s
- SLO compliance: ON TRACK/AT RISK

### Weekly On-Call Report

**Sent every Monday:**
- Incidents count (by severity)
- Total downtime
- Error budget consumed
- Notable incidents
- Action items

### Monthly SLO Review

**First Wednesday of month:**
- SLO performance vs targets
- Error budget status
- Incident analysis
- Improvement plan
- Next month risks

---

## 🔗 Related Documentation

### Service Documentation
- [Main README](../../README.md) - Service overview and architecture
- [API Documentation](../../api/proto/listings/v1/) - gRPC definitions
- [Deployment Guide](../../deployment/) - Deployment procedures

### External Resources
- **Google SRE Book:** https://sre.google/sre-book/
- **Prometheus Docs:** https://prometheus.io/docs/
- **Grafana Docs:** https://grafana.com/docs/
- **PostgreSQL Docs:** https://www.postgresql.org/docs/

---

## ❓ FAQ

**Q: Which document should I read first?**
A: If you're on-call, start with ON_CALL_GUIDE.md. Otherwise, start with this README.

**Q: What if I can't find my issue in RUNBOOK?**
A: Use TROUBLESHOOTING.md decision tree, or escalate if time-critical.

**Q: How do I know when to escalate?**
A: See ON_CALL_GUIDE.md escalation matrix. When in doubt, escalate.

**Q: Can I make changes to production during an incident?**
A: Yes, as documented in these procedures. Always document what you change.

**Q: What if I'm not sure about an action?**
A: Ask for help! Use Slack #listings-incidents or escalate to backup on-call.

**Q: How do I update these documents?**
A: Create PR with changes, tag Platform Team for review. See "Contributing" above.

---

## 📞 Support

**Need help with operations documentation?**
- **Slack:** #listings-team or #platform-team
- **Email:** platform@svetu.rs
- **On-Call:** PagerDuty escalation

**Found an error or have suggestions?**
- Create GitHub issue: `ops-docs: [description]`
- Or directly create PR with fix

---

## 📜 Document History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0.0 | 2025-11-05 | Initial comprehensive operations documentation | Platform Team |

---

**Remember: These documents exist to help you. If something is unclear, confusing, or missing, please improve it for the next person!**

**Good luck, and may your incidents be few and quickly resolved! 🚀**
