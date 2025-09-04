# 🚀 Production Deployment Plan: Unified Attributes System
## Zero-Downtime Migration Strategy

*Version: 1.0.0*
*Date: 03.09.2025*
*Environment: Production (svetu.rs)*

---

## 📋 Executive Summary

Production deployment план для миграции на unified attributes систему с использованием blue-green deployment и canary release стратегий для обеспечения zero-downtime и минимизации рисков.

---

## 🎯 Deployment Objectives

1. **Zero Downtime** - никакого прерывания сервиса
2. **Gradual Rollout** - поэтапный переход (10% → 25% → 50% → 100%)
3. **Quick Rollback** - откат за < 5 минут
4. **Data Integrity** - 100% сохранность данных
5. **Performance** - без деградации производительности

---

## 🏗️ Infrastructure Setup

### 1. Database Configuration
```yaml
Production Database:
  Master: prod-db-master.svetu.rs
  Replica: prod-db-replica.svetu.rs
  Backup: prod-db-backup.svetu.rs
  
Migrations:
  - 000034_unified_attributes.up.sql
  - 000035_migrate_attributes_data.up.sql
  
Rollback Scripts:
  - 000035_migrate_attributes_data.down.sql
  - 000034_unified_attributes.down.sql
```

### 2. Load Balancer Configuration
```nginx
# Blue-Green Traffic Distribution
upstream backend_blue {
    server backend-blue-1.svetu.rs:3000;
    server backend-blue-2.svetu.rs:3000;
    server backend-blue-3.svetu.rs:3000;
}

upstream backend_green {
    server backend-green-1.svetu.rs:3000;
    server backend-green-2.svetu.rs:3000;
    server backend-green-3.svetu.rs:3000;
}

# Canary Release Rules
map $cookie_canary $backend_pool {
    "1"     backend_green;
    default backend_blue;
}
```

### 3. Application Configuration
```yaml
# Environment Variables
ENVIRONMENT: production
USE_UNIFIED_ATTRIBUTES: true
UNIFIED_ATTRIBUTES_FALLBACK: true
DUAL_WRITE_ATTRIBUTES: true
UNIFIED_ATTRIBUTES_PERCENT: 0  # Start with 0%

# Feature Flags
feature_flags:
  unified_attributes_enabled: false
  unified_attributes_percentage: 0
  dual_write_enabled: true
  fallback_enabled: true
```

---

## 📊 Deployment Phases

### Phase 1: Pre-Deployment (Day 13)
```yaml
Duration: 4 hours
Tasks:
  - Database backup (full)
  - Deploy migrations to staging
  - Run integration tests
  - Update monitoring dashboards
  - Prepare rollback scripts
  
Validation:
  - [ ] All tests pass in staging
  - [ ] Backup verified
  - [ ] Monitoring ready
  - [ ] Team briefed
```

### Phase 2: Database Migration (Day 13)
```yaml
Duration: 2 hours
Tasks:
  - Apply migration 000034 (schema)
  - Apply migration 000035 (data)
  - Verify data integrity
  - Enable read replica sync
  
Validation:
  - [ ] No data loss
  - [ ] Replica in sync
  - [ ] Query performance normal
```

### Phase 3: Blue Environment Update (Day 14)
```yaml
Duration: 3 hours
Tasks:
  - Deploy new code to blue environment
  - Enable dual-write mode
  - Test with internal traffic
  - Monitor error rates
  
Validation:
  - [ ] Blue environment healthy
  - [ ] Dual-write working
  - [ ] No errors in logs
```

### Phase 4: Performance Optimizations Rollout (Day 21)
```yaml
Duration: 3 hours
Optimizations Deployment:
  Step 1 (30min): Deploy 6 new database indexes
  Step 2 (30min): Implement Redis caching strategy
  Step 3 (60min): Deploy frontend performance optimizations
  Step 4 (60min): A/B testing with performance validation
  
Success Criteria Based on Day 20 Testing:
  - API response time < 5ms for 95% requests (was 2.9ms in testing)
  - Cache hit rate > 75% (achieved 78% in testing)
  - Categories endpoint > 2000 req/sec (achieved 3411 req/sec)
  - Search endpoint > 1000 req/sec (achieved 1236 req/sec)
  - Zero data loss or errors during deployment
```

### Phase 5: Green Environment Deployment (Day 15)
```yaml
Duration: 2 hours
Tasks:
  - Deploy to green environment
  - Switch primary traffic to green
  - Keep blue as fallback
  - Monitor all metrics
  
Validation:
  - [ ] Green environment stable
  - [ ] Performance metrics good
  - [ ] Blue ready for rollback
```

### Phase 6: Finalization (Day 15)
```yaml
Duration: 2 hours
Tasks:
  - Disable dual-write mode
  - Remove fallback flags
  - Decommission blue environment
  - Archive old tables
  
Validation:
  - [ ] System fully on unified attributes
  - [ ] All metrics normal
  - [ ] Documentation updated
```

---

## 🔄 Rollback Procedures

### Instant Rollback (< 1 minute)
```bash
# Switch traffic back to blue
kubectl set image deployment/backend backend=backend:v1.0.0
kubectl rollout status deployment/backend

# Or via load balancer
curl -X POST https://lb-control.svetu.rs/switch-to-blue
```

### Database Rollback (< 5 minutes)
```bash
# Apply down migrations
migrate -path migrations -database $DATABASE_URL down 2

# Restore from backup if needed
pg_restore -d svetubd backup_before_migration.dump

# Restart applications
kubectl rollout restart deployment/backend
```

### Feature Flag Rollback (< 30 seconds)
```bash
# Disable via API
curl -X POST https://api.svetu.rs/internal/feature-flags \
  -d '{"unified_attributes_enabled": false}'

# Or via environment variable
kubectl set env deployment/backend USE_UNIFIED_ATTRIBUTES=false
```

---

## 📈 Monitoring & Alerts

### Key Metrics to Monitor
```yaml
Application Metrics:
  - Request rate (req/s)
  - Error rate (%)
  - Response time p50, p95, p99
  - Active connections
  
Database Metrics:
  - Query execution time
  - Connection pool usage
  - Replication lag
  - Lock wait time
  
Business Metrics:
  - Listing creation rate
  - Search success rate
  - Attribute usage patterns
  - User activity
```

### Alert Thresholds
```yaml
Critical Alerts:
  - Error rate > 1%
  - p95 latency > 100ms
  - Database replication lag > 1s
  - Memory usage > 90%
  
Warning Alerts:
  - Error rate > 0.5%
  - p95 latency > 75ms
  - CPU usage > 80%
  - Disk usage > 85%
```

### Dashboards
```yaml
Grafana Dashboards:
  - Unified Attributes Overview
  - Migration Progress
  - Performance Comparison
  - Error Analysis
  - User Impact
  
URLs:
  - https://grafana.svetu.rs/d/unified-attrs
  - https://grafana.svetu.rs/d/migration-status
```

---

## 👥 Team Responsibilities

### DevOps Team
- Infrastructure setup
- Load balancer configuration
- Deployment execution
- Monitoring setup

### Backend Team
- Code deployment
- Migration execution
- Performance monitoring
- Bug fixes

### QA Team
- Staging validation
- Canary testing
- User acceptance
- Bug reporting

### Product Team
- User communication
- Success metrics
- Feedback collection
- Go/No-go decisions

---

## 📅 Timeline

### Day 21 (Today) - Performance Optimizations Deployment
- 14:00-15:00: Database performance optimizations deployment
- 15:00-15:30: Redis cache strategy implementation
- 15:30-16:00: Backend optimizations deployment
- 16:00-16:30: Performance validation and testing
- 16:30-17:00: A/B testing setup and monitoring

### Day 22
- 09:00-12:00: Monitor A/B testing results (50% traffic to optimized version)
- 12:00-14:00: Gradual rollout increase (50% → 75% → 100%)
- 14:00-16:00: Full production monitoring
- 16:00-17:00: Performance metrics analysis and reporting

### Day 23
- 09:00-11:00: Final optimization deployment and cleanup
- 11:00-13:00: Legacy system deprecation
- 13:00-15:00: Documentation updates and team briefing
- 15:00-17:00: Success validation and next steps planning

---

## ✅ Pre-Deployment Checklist

### Technical
- [ ] All tests passing in CI/CD
- [ ] Staging environment validated
- [ ] Database backups completed
- [ ] Rollback scripts tested
- [ ] Monitoring dashboards ready
- [ ] Alert rules configured

### Process
- [ ] Deployment approved by stakeholders
- [ ] Maintenance window scheduled
- [ ] Communication sent to users
- [ ] Support team briefed
- [ ] On-call schedule confirmed
- [ ] Runbook distributed

### Documentation
- [ ] Deployment guide finalized
- [ ] Rollback procedures documented
- [ ] Troubleshooting guide ready
- [ ] Post-mortem template prepared

---

## 🚨 Risk Mitigation

### High Risk Items
1. **Data Corruption**
   - Mitigation: Full backup, dual-write validation
   - Recovery: Restore from backup (15 min)

2. **Performance Degradation**
   - Mitigation: Canary release, real-time monitoring
   - Recovery: Feature flag disable (30 sec)

3. **Compatibility Issues**
   - Mitigation: Fallback mechanism, extensive testing
   - Recovery: Rollback deployment (5 min)

### Communication Plan
```yaml
Internal:
  - Slack: #deployment-unified-attrs
  - Email: tech-team@svetu.rs
  - War Room: meet.google.com/unified-deploy
  
External:
  - Status Page: status.svetu.rs
  - User Email: If issues > 30 min
  - Support Team: Tier 1 briefing
```

---

## 📊 Success Criteria

### Technical Metrics
- ✅ Zero data loss
- ✅ Error rate < 0.1%
- ✅ p95 latency < 50ms
- ✅ 100% feature parity
- ✅ All tests passing

### Business Metrics
- ✅ No drop in listing creation
- ✅ Search functionality maintained
- ✅ User complaints < 5
- ✅ Support tickets < 10

---

## 📝 Post-Deployment Tasks

### Immediate (Day 15)
- [ ] Remove feature flags
- [ ] Disable dual-write
- [ ] Update documentation
- [ ] Send success communication

### Week 1
- [ ] Performance analysis
- [ ] User feedback collection
- [ ] Bug fix deployment
- [ ] Optimization implementation

### Month 1
- [ ] Archive old tables
- [ ] Remove legacy code
- [ ] Final documentation
- [ ] Post-mortem review

---

## 🔗 Related Documents

- [Technical Specification](TZ_ATTRIBUTE_SYSTEM_UNIFICATION.md)
- [Testing Guide](UNIFIED_ATTRIBUTES_TESTING_GUIDE.md)
- [CI/CD Pipeline](./.github/workflows/unified-attributes-test.yml)
- [Monitoring Setup](./backend/monitoring/)

---

**Document Status:** READY FOR REVIEW
**Last Updated:** 03.09.2025
**Next Review:** Before deployment

---
