# Listings Microservice SLO Management Guide

**Last Updated:** 2025-11-05
**Version:** 1.0.0
**Maintainer:** Platform Team

## Table of Contents

- [Overview](#overview)
- [SLO Definitions](#slo-definitions)
- [Error Budget](#error-budget)
- [Tracking and Reporting](#tracking-and-reporting)
- [Incident Impact Calculation](#incident-impact-calculation)
- [Monthly Reviews](#monthly-reviews)
- [SLO Breach Response](#slo-breach-response)
- [Improving SLOs](#improving-slos)

---

## Overview

### What are SLOs?

**Service Level Objectives (SLOs)** are measurable targets for service reliability. They define the level of service we promise to deliver to our users.

**Key Concepts:**
- **SLO:** Target reliability (e.g., 99.9% uptime)
- **SLI (Service Level Indicator):** Actual measurement (e.g., 99.92% uptime)
- **Error Budget:** Allowed unreliability (100% - SLO = 0.1% downtime)
- **SLA (Service Level Agreement):** External contract with consequences

### Why SLOs Matter

1. **User Expectations:** Define what users can expect
2. **Risk Management:** Quantify reliability risks
3. **Decision Making:** Balance velocity vs reliability
4. **Prioritization:** Focus engineering effort
5. **Communication:** Common language for reliability

### SLO Philosophy

```
Perfect reliability (100%) is:
- Impossible to achieve
- Too expensive
- Slows down innovation

SLOs provide:
- Realistic targets
- Room for failure
- Balance between reliability and velocity
```

---

## SLO Definitions

### Primary SLOs

#### 1. Availability SLO

**Definition:** Percentage of time the service is successfully responding to requests.

```
SLO Target: 99.9% availability (monthly)
Measurement Period: 30 days (rolling window)
Error Budget: 43.2 minutes downtime per month
```

**Success Criteria:**
```
Availability = (Successful Requests / Total Requests) × 100

Successful = HTTP 2xx, 3xx, 4xx (except 429)
Failed = HTTP 5xx, timeouts, connection errors
```

**Measurement:**
```promql
# PromQL Query
sum(rate(listings_grpc_requests_total{status!~"5.."}[30d]))
/
sum(rate(listings_grpc_requests_total[30d]))
* 100
```

**Exclusions:**
- Planned maintenance (with advance notice)
- User errors (4xx except 429 rate limits)
- Client-side failures
- DDoS attacks (mitigated)

---

#### 2. Latency SLO

**Definition:** Percentage of requests completing within latency threshold.

```
SLO Target: 95% of requests < 1 second (P95)
SLO Target: 99% of requests < 2 seconds (P99)
Measurement Period: 30 days (rolling window)
```

**Measurement:**
```promql
# P95 Latency
histogram_quantile(0.95,
  rate(listings_grpc_request_duration_seconds_bucket[30d])
)

# P99 Latency
histogram_quantile(0.99,
  rate(listings_grpc_request_duration_seconds_bucket[30d])
)
```

**Per-Endpoint Targets:**

| Endpoint | P95 Target | P99 Target |
|----------|-----------|-----------|
| GetListing | < 500ms | < 1s |
| ListListings | < 800ms | < 1.5s |
| CreateListing | < 1s | < 2s |
| UpdateListing | < 1s | < 2s |
| SearchListings | < 1s | < 2s |
| BatchUpdateStock | < 5s | < 10s |

---

#### 3. Error Rate SLO

**Definition:** Percentage of requests that fail with errors.

```
SLO Target: < 1% error rate
Measurement Period: 30 days (rolling window)
Error Budget: 1% of all requests
```

**Error Definition:**
```
Errors include:
- HTTP 5xx responses
- gRPC Internal, Unavailable, DataLoss errors
- Timeouts (DeadlineExceeded)
- Database connection failures

NOT errors:
- 4xx client errors (except 429)
- Rate limit rejections (429)
- Invalid input (400)
```

**Measurement:**
```promql
# Error Rate
sum(rate(listings_grpc_requests_total{status=~"5.."}[30d]))
/
sum(rate(listings_grpc_requests_total[30d]))
* 100
```

---

### Secondary SLOs

#### 4. Data Durability SLO

**Definition:** No data loss under normal operations.

```
SLO Target: 99.99% durability (annual)
Maximum Acceptable Loss: 1 in 10,000 writes
```

**Measurement:**
- Database transaction failures
- Replication lag monitoring
- Backup verification

---

#### 5. Data Freshness SLO

**Definition:** Search index reflects database state within threshold.

```
SLO Target: < 1 minute lag between DB and OpenSearch
Measurement: Reindexing latency monitoring
```

---

## Error Budget

### Understanding Error Budget

**Error Budget = Allowed Unreliability**

```
If SLO = 99.9% availability
Then Error Budget = 0.1% unavailability
     = 43.2 minutes per month
     = 0.72 minutes per day
     = 43 seconds per hour
```

### Error Budget Calculation

**Monthly Budget (30 days):**
```
Total time = 30 days × 24 hours × 60 minutes = 43,200 minutes
Error budget = 43,200 × (100% - 99.9%) = 43.2 minutes
```

**Requests-based Budget:**
```
If average 1000 requests/second:
Total requests = 1000 × 60 × 60 × 24 × 30 = 2,592,000,000
Error budget = 2,592,000,000 × 0.1% = 2,592,000 failed requests
```

### Error Budget Policy

**When Error Budget is Healthy (>50% remaining):**
- ✅ Deploy new features freely
- ✅ Conduct experiments
- ✅ Scheduled maintenance allowed
- ✅ Normal release velocity

**When Error Budget is Low (<50% remaining):**
- ⚠️ Increase caution on deployments
- ⚠️ Focus on stability improvements
- ⚠️ Defer non-critical features
- ⚠️ Enhanced testing required

**When Error Budget is Depleted (0% remaining):**
- 🛑 **FEATURE FREEZE**
- 🛑 Only reliability improvements
- 🛑 Only critical bug fixes
- 🛑 Incident postmortems mandatory
- 🛑 Weekly review until recovered

### Error Budget Reset

- **Automatic:** Beginning of each month
- **Emergency Reset:** CTO approval only (documented)

---

## Tracking and Reporting

### Real-Time Monitoring

**Grafana Dashboard:** "Listings SLO Overview"
- Current SLO compliance
- Error budget remaining
- Trend graphs (7d, 30d)
- Per-endpoint breakdown

**Access:** https://grafana.svetu.rs/d/listings-slo

### Daily Check

```bash
# Quick SLO status
curl -s http://localhost:9090/api/v1/query --data-urlencode 'query=
  (sum(rate(listings_grpc_requests_total{status!~"5.."}[24h])) /
   sum(rate(listings_grpc_requests_total[24h]))) * 100
' | jq '.data.result[0].value[1]'

# Expected output: "99.95" (above 99.9% SLO)
```

### Weekly Report

**Automated email every Monday:**
```
Subject: Listings SLO Weekly Report - Week ending [DATE]

Availability: 99.92% ✅ (Target: 99.9%)
P95 Latency: 850ms ✅ (Target: <1s)
P99 Latency: 1.8s ✅ (Target: <2s)
Error Rate: 0.5% ✅ (Target: <1%)

Error Budget:
- Remaining: 35.2 minutes (81.5% of monthly budget) ✅
- Consumed: 8.0 minutes (18.5%)

Notable Incidents:
- [2025-11-02] High latency (5 minutes) - Database slow query
- [2025-11-04] Elevated errors (2 minutes) - Redis connection issue

Action Items:
- [TICKET-123] Optimize slow query from 2025-11-02
- [TICKET-124] Improve Redis connection handling

Next Review: 2025-11-12
```

### Monthly Report

**Comprehensive analysis:**
```markdown
# Listings SLO Monthly Report - [MONTH YEAR]

## Executive Summary
- Overall SLO Compliance: [PASS/FAIL]
- Error Budget Status: [% remaining]
- Major Incidents: [COUNT]
- Improvement Trend: [UP/DOWN/STABLE]

## SLO Performance

### Availability
- Target: 99.9%
- Actual: 99.91% ✅
- Error Budget: 38.4 minutes available, 4.8 minutes used
- Incidents contributing to downtime:
  - [2025-11-05] Service restart (1.5 min)
  - [2025-11-12] Database failover (2.0 min)
  - [2025-11-20] Deployment rollback (1.3 min)

### Latency
- P95 Target: < 1s, Actual: 920ms ✅
- P99 Target: < 2s, Actual: 1.85s ✅
- Latency spikes: 3 incidents totaling 12 minutes

### Error Rate
- Target: < 1%, Actual: 0.6% ✅
- Total errors: 15,552 out of 2,592,000 requests

## Trends
[Graphs and analysis]

## Top Contributors to SLO Violations
1. Database connection pool exhaustion (40% of budget)
2. Slow queries (30% of budget)
3. Deployment issues (20% of budget)
4. External service failures (10% of budget)

## Improvements Implemented
- [Improvement 1]
- [Improvement 2]

## Action Items for Next Month
- [Action 1]
- [Action 2]

## Recommendations
[Strategic recommendations]

Prepared by: [Name]
Date: [Date]
```

---

## Incident Impact Calculation

### How to Calculate SLO Impact

**For Availability:**
```
Impact (minutes) = Duration of incident (minutes)

Example:
- Incident duration: 10 minutes
- Impact on SLO: 10 minutes
- Remaining budget: 43.2 - 10 = 33.2 minutes
```

**For Partial Outage:**
```
Impact = Duration × Error Rate

Example:
- Incident duration: 30 minutes
- Error rate during incident: 20%
- Impact on SLO: 30 × 0.20 = 6 minutes
- Remaining budget: 43.2 - 6 = 37.2 minutes
```

**For Latency Degradation:**
```
If P99 exceeds threshold:
Count requests above threshold as "failed" for error budget

Example:
- 1000 requests had P99 > 2s
- Total requests: 1,000,000
- Impact on error budget: 0.1%
```

### Impact Classification

| Impact | Error Budget Consumed | Severity |
|--------|---------------------|----------|
| **Minor** | < 10% (< 4.3 min) | P3 |
| **Moderate** | 10-25% (4.3-10.8 min) | P2 |
| **Major** | 25-50% (10.8-21.6 min) | P2 |
| **Critical** | > 50% (> 21.6 min) | P1 |

### Real-Time Impact Tracking

**During incidents, track impact:**
```bash
# Start time
START_TIME=$(date +%s)

# After resolution
END_TIME=$(date +%s)
DURATION=$(( (END_TIME - START_TIME) / 60 ))

# Check error rate during incident
ERROR_RATE=$(curl -s "http://localhost:9090/api/v1/query" --data-urlencode "query=
  sum(rate(listings_grpc_requests_total{status=~\"5..\"}[${DURATION}m])) /
  sum(rate(listings_grpc_requests_total[${DURATION}m]))
" | jq -r '.data.result[0].value[1]')

# Calculate impact
IMPACT=$(echo "$DURATION * $ERROR_RATE" | bc)
echo "SLO Impact: $IMPACT minutes"

# Remaining budget
MONTHLY_BUDGET=43.2
REMAINING=$(echo "$MONTHLY_BUDGET - $IMPACT" | bc)
echo "Remaining error budget: $REMAINING minutes"
```

---

## Monthly Reviews

### Review Schedule

**First Wednesday of each month at 2:00 PM**
- Duration: 1 hour
- Attendees: Platform team, Engineering leads, Product manager
- Location: Conference room / Zoom

### Review Agenda

**1. SLO Performance Review (15 minutes)**
- Present monthly SLO dashboard
- Discuss compliance vs targets
- Review error budget status

**2. Incident Analysis (20 minutes)**
- Review all incidents from past month
- Identify patterns and trends
- Discuss root causes

**3. Action Items Review (15 minutes)**
- Status of previous month's action items
- Effectiveness of improvements

**4. Forward Planning (10 minutes)**
- Planned work that may impact SLO
- Risk assessment for next month
- Preventive measures

**5. SLO Adjustment Discussion (optional, 10 minutes)**
- Should SLOs be tightened or relaxed?
- New SLOs to add?
- Measurement improvements needed?

### Review Template

```markdown
# Listings SLO Monthly Review - [Month Year]

Date: [Date]
Attendees: [Names]
Facilitator: [Name]

## 1. SLO Performance

| Metric | Target | Actual | Status | Notes |
|--------|--------|--------|--------|-------|
| Availability | 99.9% | 99.91% | ✅ PASS | Excellent |
| P95 Latency | <1s | 920ms | ✅ PASS | Good |
| P99 Latency | <2s | 1.85s | ✅ PASS | Good |
| Error Rate | <1% | 0.6% | ✅ PASS | Excellent |

**Overall: PASS** ✅

Error Budget: 38.4 min remaining (88.9%) - Healthy

## 2. Incidents Summary

Total Incidents: 5 (P1: 0, P2: 2, P3: 3)

### P2 Incidents
1. **Database Connection Pool Exhausted** (2025-11-12, 10 min)
   - Impact: 10 minutes downtime
   - Root cause: Traffic spike + small pool
   - Action: Increased pool size to 50

2. **High Latency** (2025-11-18, 5 min)
   - Impact: 5 minutes degraded performance
   - Root cause: Slow query without index
   - Action: Added index, query optimized

### P3 Incidents
[List...]

## 3. Trends

- ⬆️ **Traffic increased 15%** vs previous month
- ⬇️ **Error rate decreased** from 0.8% to 0.6%
- ➡️ **Latency stable**
- ⬆️ **Incidents increased** from 3 to 5 (traffic-related)

## 4. Previous Action Items

| Action | Status | Outcome |
|--------|--------|---------|
| Increase DB pool size | ✅ DONE | Resolved pool exhaustion |
| Add slow query monitoring | ✅ DONE | Caught 2 slow queries |
| Improve deployment process | 🟡 IN PROGRESS | 60% complete |

## 5. New Action Items

| Action | Owner | Due Date | Priority |
|--------|-------|----------|----------|
| Implement circuit breaker for OpenSearch | @alice | 2025-12-15 | High |
| Add automated canary deployments | @bob | 2025-12-31 | Medium |
| Review and optimize top 10 queries | @charlie | 2025-12-10 | High |

## 6. Risks for Next Month

- **Black Friday traffic spike** - Expected 3x traffic
  - Mitigation: Load testing scheduled, capacity increased
- **Planned database upgrade** - 10-minute maintenance window
  - Mitigation: Scheduled during low-traffic period

## 7. SLO Adjustments

**Proposal:** No changes recommended
- Current SLOs appropriate for service maturity
- Error budget policy working well
- Re-evaluate in 3 months

## 8. Conclusion

Overall: **Excellent month** ✅
- All SLOs met comfortably
- Error budget healthy
- Improvements showing positive impact
- Team handling incidents effectively

Focus for next month:
- Prepare for traffic spike
- Continue query optimization
- Complete deployment automation

---

Next Review: 2025-12-04 at 2:00 PM
```

---

## SLO Breach Response

### When SLO is Breached

**Definition:** SLO breach occurs when SLI falls below SLO target at end of measurement period.

**Example:**
```
Monthly availability: 99.85%
Target SLO: 99.9%
Status: ❌ SLO BREACHED
```

### Immediate Actions (Day 1)

**1. Acknowledge Breach**
```bash
# Send notification
# To: Engineering leadership, Product, Customer success
# Subject: Listings SLO Breach - [Month]

"The Listings service did not meet its 99.9% availability SLO
for [Month], achieving 99.85% availability.

Root cause analysis and corrective action plan in progress.
Full report will be shared within 3 business days."
```

**2. Initiate Root Cause Analysis**
- Review all incidents from the period
- Identify primary contributors to breach
- Analyze patterns and systemic issues

**3. Emergency Team Meeting**
- Schedule within 24 hours
- Attendees: Platform team, engineering leads, CTO
- Focus: Understanding breach and immediate mitigations

### Corrective Action Plan (Week 1)

**1. Document Root Causes**
```markdown
# SLO Breach Root Cause Analysis

## Breach Details
- SLO Target: 99.9%
- Actual: 99.85%
- Gap: 0.05% (21.6 minutes over budget)
- Measurement Period: [Start] to [End]

## Contributing Incidents
1. Database failover (15 min downtime) - 35% of breach
2. Deployment rollback (10 min downtime) - 23% of breach
3. High latency incidents (18 min degraded) - 42% of breach

## Root Causes
1. **Insufficient database failover testing**
   - Failover took 15 minutes instead of expected 5 minutes
   - Monitoring didn't detect stale connection issue

2. **Deployment process lacking safeguards**
   - Bad deployment reached production
   - Rollback took 10 minutes (manual process)

3. **Query performance regression**
   - New feature introduced N+1 query pattern
   - Code review didn't catch database implications

## Systemic Issues
- Lack of automated canary deployments
- Insufficient load testing before releases
- Query performance not part of CI/CD
```

**2. Create Corrective Actions**
```markdown
## Corrective Actions

### Immediate (Week 1)
- [ ] Implement automated database failover testing (weekly)
- [ ] Add rollback automation to deployment pipeline
- [ ] Conduct emergency query optimization sprint

### Short-term (Month 1)
- [ ] Implement automated canary deployments
- [ ] Add query performance regression tests to CI
- [ ] Increase database connection pool timeout
- [ ] Enhance monitoring for database failover

### Long-term (Quarter 1)
- [ ] Implement read replicas for database
- [ ] Add chaos engineering practices
- [ ] Conduct quarterly disaster recovery drills
```

**3. Communicate Plan**
```bash
# Send to stakeholders
# Subject: Listings SLO Breach - Corrective Action Plan

"Following the SLO breach in [Month], we have completed
our root cause analysis and developed a corrective action plan.

Key findings:
- Database failover took longer than expected
- Deployment process lacked automated rollback
- Query performance regression went undetected

Actions being taken:
- [List top 3-5 actions]

Timeline:
- Immediate actions: Week 1
- Short-term actions: Month 1
- Long-term improvements: Quarter 1

We are confident these measures will prevent recurrence
and improve overall service reliability.

Full report attached."
```

### Feature Freeze Implementation

**When error budget is depleted:**

**1. Announce Feature Freeze**
```
Subject: FEATURE FREEZE: Listings Service

Due to error budget depletion, Listings service is now
in FEATURE FREEZE until error budget recovers.

ALLOWED:
- Reliability improvements
- Bug fixes (P1, P2)
- Security patches
- Documentation updates

NOT ALLOWED:
- New features
- Non-critical refactoring
- Experiments
- Performance optimizations (unless P1/P2)

Duration: Until error budget recovers to >25%
Estimated: [Date]

Questions? Contact Platform Team Lead.
```

**2. Review Process**
```bash
# All changes must be approved by Platform Team Lead
# Include in PR description:
# - [ ] This change improves reliability
# - [ ] This is a critical bug fix (P1/P2)
# - [ ] This has been thoroughly tested
# - [ ] This has minimal risk
```

**3. Monitor Recovery**
```bash
# Daily error budget check
curl -s "http://localhost:9090/api/v1/query" --data-urlencode 'query=
  (43.2 - (43200 * (1 - (
    sum(rate(listings_grpc_requests_total{status!~"5.."}[30d])) /
    sum(rate(listings_grpc_requests_total[30d]))
  )))) / 43.2 * 100
' | jq -r '.data.result[0].value[1]'

# If > 25%, consider lifting freeze
# If > 50%, lift freeze
```

---

## Improving SLOs

### When to Tighten SLOs

**Consider tightening when:**
- Consistently exceeding SLO by >10% for 3+ months
- Error budget always >80% remaining
- Users expect better reliability
- Competitive pressure

**Example:**
```
Current: 99.9% → Proposed: 99.95%
Impact: Error budget reduces from 43.2 min to 21.6 min
```

**Process:**
1. Analyze 6 months of historical data
2. Model impact of proposed SLO
3. Assess engineering effort required
4. Business value vs cost analysis
5. Phased rollout (if approved)

### When to Relax SLOs

**Consider relaxing when:**
- Consistently missing SLO despite best efforts
- SLO doesn't reflect user experience
- Cost of achieving SLO is prohibitive
- Business priorities have changed

**⚠️ Warning:** Relaxing SLOs should be rare and well-justified.

**Process:**
1. Document why current SLO can't be met
2. Analyze user impact of proposed change
3. Get executive approval
4. Communicate to all stakeholders
5. Implement with 30-day notice period

### Adding New SLOs

**Good candidates:**
- Data durability
- Data freshness
- API versioning compatibility
- Security patch latency

**Process:**
1. Define SLI (measurement)
2. Baseline current performance (3 months)
3. Set realistic target based on baseline
4. Implement measurement and alerting
5. Shadow mode for 1 month
6. Official adoption

---

## Tools and Automation

### Prometheus Queries

**Availability:**
```promql
100 - (
  sum(rate(listings_grpc_requests_total{status=~"5.."}[30d])) /
  sum(rate(listings_grpc_requests_total[30d])) * 100
)
```

**Error Budget Remaining:**
```promql
(43.2 - (43200 * (1 - (
  sum(rate(listings_grpc_requests_total{status!~"5.."}[30d])) /
  sum(rate(listings_grpc_requests_total[30d]))
)))) / 43.2 * 100
```

**P95 Latency:**
```promql
histogram_quantile(0.95,
  rate(listings_grpc_request_duration_seconds_bucket[30d])
)
```

**P99 Latency:**
```promql
histogram_quantile(0.99,
  rate(listings_grpc_request_duration_seconds_bucket[30d])
)
```

### Alerting Rules

**Located:** `/p/github.com/sveturs/listings/deployment/prometheus/alerts.yml`

**Key alerts:**
- `SLOAvailabilityAtRisk` - Error budget < 25%
- `SLOAvailabilityBreach` - Availability < 99.9%
- `SLOLatencyAtRisk` - P99 > 2s for 30 minutes
- `SLOErrorRateHigh` - Error rate > 1% for 30 minutes

---

## Appendix

### Further Reading

- **Google SRE Book - Chapter 4:** Service Level Objectives
  https://sre.google/sre-book/service-level-objectives/

- **Implementing SLOs:** https://sre.google/workbook/implementing-slos/

- **Error Budget Policy:** https://sre.google/workbook/error-budget-policy/

### SLO Calculator

Online tool: https://sre.google/sre-book/service-level-objectives/#calculator

### Review Checklist

**Monthly Review Checklist:**
- [ ] Gather SLO data for past month
- [ ] Calculate actual SLIs
- [ ] Determine SLO compliance (pass/fail)
- [ ] Review all incidents
- [ ] Calculate error budget remaining
- [ ] Identify trends
- [ ] Review previous action items
- [ ] Create new action items
- [ ] Schedule next review
- [ ] Distribute report to stakeholders

---

**Document Version:** 1.0.0
**Last Reviewed:** 2025-11-05
**Next Review:** 2025-12-05
**Owner:** Platform Team

**Remember: SLOs are a tool, not a goal. The goal is happy users and sustainable engineering.**
