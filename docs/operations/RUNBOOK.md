# Listings Microservice Operations Runbook

**Last Updated:** 2025-11-05
**Version:** 1.0.0
**Maintainer:** Platform Team

## Table of Contents

- [Quick Reference](#quick-reference)
- [Common Incidents](#common-incidents)
  - [1. High Error Rate (>1%)](#1-high-error-rate-1)
  - [2. High Latency (P99 >2s)](#2-high-latency-p99-2s)
  - [3. Service Down](#3-service-down)
  - [4. Database Connection Pool Exhausted](#4-database-connection-pool-exhausted)
  - [5. Redis Connection Issues](#5-redis-connection-issues)
  - [6. OpenSearch Cluster Red](#6-opensearch-cluster-red)
  - [7. Memory Leak / OOM Killed](#7-memory-leak--oom-killed)
  - [8. Disk Space Critical](#8-disk-space-critical)
  - [9. Rate Limit Abuse](#9-rate-limit-abuse)
  - [10. Slow Queries](#10-slow-queries)
- [Recovery Procedures](#recovery-procedures)
- [Escalation](#escalation)

---

## Quick Reference

### Critical Commands

```bash
# Service Status
sudo systemctl status listings-service
sudo systemctl restart listings-service
sudo systemctl stop listings-service
sudo systemctl start listings-service

# Logs (last 1 hour, follow mode)
sudo journalctl -u listings-service --since "1 hour ago" -f

# Check service health
curl http://localhost:8086/health
curl http://localhost:8086/ready

# Check metrics
curl http://localhost:8086/metrics | grep listings_

# Database connections
psql "postgres://listings_user:listings_password@localhost:35433/listings_db?sslmode=disable" \
  -c "SELECT count(*) FROM pg_stat_activity WHERE datname='listings_db';"

# Redis connectivity
redis-cli -h localhost -p 36380 -a redis_password ping

# gRPC health check
grpcurl -plaintext localhost:50053 grpc.health.v1.Health/Check

# Process info
ps aux | grep listings-service
top -p $(pgrep listings-service)
```

### Service Endpoints

| Protocol | Port | Endpoint | Purpose |
|----------|------|----------|---------|
| gRPC | 50053 | - | Internal API |
| HTTP | 8086 | /health | Health check |
| HTTP | 8086 | /ready | Readiness check |
| HTTP | 8086 | /api/v1/* | REST API |
| HTTP | 9093 | /metrics | Prometheus metrics |

### SLO Targets

- **Availability:** 99.9% (43.2 minutes downtime/month)
- **Error Rate:** < 1% (excluding 4xx)
- **Latency P99:** < 2 seconds
- **Latency P95:** < 1 second

### Emergency Contacts

| Role | Contact | When to Escalate |
|------|---------|------------------|
| On-Call Engineer | PagerDuty: listings-oncall | All alerts |
| Database SRE | PagerDuty: db-team | DB connection issues, slow queries |
| Platform Team Lead | Slack: #listings-incidents | SLO breach, data loss |
| Security Team | security@svetu.rs | Security incidents, DDoS |

---

## Common Incidents

### 1. High Error Rate (>1%)

#### Symptoms

- Grafana dashboard shows error rate > 1%
- Alert: `ListingsHighErrorRate`
- Logs contain frequent ERROR level messages
- Metrics: `listings_errors_total` increasing rapidly

#### Investigation Steps

**Step 1: Identify error types**
```bash
# Check error breakdown by component
curl -s http://localhost:8086/metrics | grep listings_errors_total

# Check logs for error patterns
sudo journalctl -u listings-service --since "10 minutes ago" | grep '"level":"error"' | jq .
```

**Step 2: Check dependencies**
```bash
# PostgreSQL
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"

# Redis
redis-cli -h localhost -p 36380 -a redis_password ping

# OpenSearch
curl -u admin:admin http://localhost:9200/_cluster/health
```

**Step 3: Check service resources**
```bash
# Memory usage
ps aux | grep listings-service | awk '{print $6/1024 " MB"}'

# CPU usage
top -b -n 1 -p $(pgrep listings-service) | tail -1

# File descriptors
ls -l /proc/$(pgrep listings-service)/fd | wc -l
```

#### Resolution Steps

**Scenario A: Database connection errors**
```bash
# Check connection pool exhaustion
curl -s http://localhost:8086/metrics | grep listings_db_connections

# If pool exhausted, restart service
sudo systemctl restart listings-service

# Monitor recovery
watch -n 1 'curl -s http://localhost:8086/metrics | grep listings_db_connections'
```

**Scenario B: Redis connection errors**
```bash
# Check Redis status
redis-cli -h localhost -p 36380 -a redis_password info | grep "connected_clients"

# Clear rate limit keys if Redis memory issue
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'rate_limit:*' | xargs redis-cli DEL

# Service will fail-open (allow requests) if Redis is down
```

**Scenario C: OpenSearch errors**
```bash
# Check cluster health
curl -u admin:admin http://localhost:9200/_cluster/health?pretty

# Service can operate without OpenSearch (search will fail, CRUD continues)
# No immediate action needed unless search is critical
```

#### Prevention

- Set up proper monitoring alerts for dependencies
- Configure connection pool sizes: `DB_MAX_OPEN_CONNS=25`
- Enable fail-open mode for non-critical dependencies (Redis, OpenSearch)
- Regular load testing to identify bottlenecks

#### Escalation Path

- **If error rate > 5%:** Page Database SRE team
- **If error rate > 10%:** Page Platform Team Lead
- **If SLO breach imminent:** Execute disaster recovery plan

---

### 2. High Latency (P99 >2s)

#### Symptoms

- Alert: `ListingsHighLatency`
- Grafana shows P99 > 2 seconds
- Metrics: `listings_grpc_request_duration_seconds` high
- Users reporting slow responses

#### Investigation Steps

**Step 1: Identify slow endpoints**
```bash
# Check latency by method
curl -s http://localhost:8086/metrics | grep listings_grpc_request_duration_seconds_bucket | sort -t= -k2 -n | tail -20

# Check timeout metrics
curl -s http://localhost:8086/metrics | grep listings_near_timeouts_total
```

**Step 2: Check database performance**
```bash
# Active queries
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pid, now() - pg_stat_activity.query_start AS duration, query
   FROM pg_stat_activity
   WHERE state = 'active' AND datname = 'listings_db'
   ORDER BY duration DESC;"

# Slow queries log
sudo tail -100 /var/log/postgresql/postgresql-15-main.log | grep "duration:"
```

**Step 3: Check system resources**
```bash
# CPU load
uptime

# Memory pressure
free -h

# Disk I/O
iostat -x 1 5

# Network latency to dependencies
ping -c 5 localhost  # PostgreSQL
```

#### Resolution Steps

**Scenario A: Slow database queries**
```bash
# Kill long-running query
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pg_terminate_backend(pid) FROM pg_stat_activity
   WHERE datname = 'listings_db' AND now() - query_start > interval '30 seconds';"

# Restart service to reset connections
sudo systemctl restart listings-service
```

**Scenario B: High CPU usage**
```bash
# Generate CPU profile (30 seconds)
curl http://localhost:8086/debug/pprof/profile?seconds=30 > /tmp/cpu.prof

# Analyze top functions
go tool pprof -top /tmp/cpu.prof

# If CPU > 80%, consider scaling horizontally
```

**Scenario C: Memory pressure**
```bash
# Generate heap profile
curl http://localhost:8086/debug/pprof/heap > /tmp/heap.prof

# Analyze memory usage
go tool pprof -top /tmp/heap.prof

# If memory leak suspected, restart service
sudo systemctl restart listings-service
```

#### Prevention

- Index optimization: Review query plans monthly
- Connection pool tuning: Adjust `DB_MAX_OPEN_CONNS` based on load
- Enable query timeout enforcement (already implemented)
- Regular performance testing under load

#### Escalation Path

- **If P99 > 5s:** Page Database SRE team
- **If persists > 15 minutes:** Page Platform Team Lead
- **If CPU > 90%:** Consider emergency scale-out

---

### 3. Service Down

#### Symptoms

- Alert: `ListingsServiceDown`
- Health check failing: `curl http://localhost:8086/health` returns error
- No metrics being exported
- Service not responding to requests

#### Investigation Steps

**Step 1: Check service status**
```bash
# Systemd status
sudo systemctl status listings-service

# Process check
ps aux | grep listings-service

# Port check
netstat -tlnp | grep -E '50053|8086|9093'
```

**Step 2: Check recent logs**
```bash
# Last 100 lines
sudo journalctl -u listings-service -n 100

# Look for panic/fatal errors
sudo journalctl -u listings-service | grep -E 'panic|fatal' | tail -20
```

**Step 3: Check system resources**
```bash
# Disk space
df -h

# Memory
free -h

# OOM killer logs
sudo dmesg | grep -i "killed process"
```

#### Resolution Steps

**Scenario A: Service crashed (not running)**
```bash
# Start service
sudo systemctl start listings-service

# Monitor startup
sudo journalctl -u listings-service -f

# Verify health
sleep 10
curl http://localhost:8086/health
```

**Scenario B: Service running but not responding**
```bash
# Check if ports are listening
sudo ss -tlnp | grep -E '50053|8086'

# Force restart
sudo systemctl restart listings-service

# If restart fails, check binary
ls -lh /opt/listings-dev/bin/listings-service
file /opt/listings-dev/bin/listings-service
```

**Scenario C: Dependency failure**
```bash
# Check PostgreSQL
sudo systemctl status postgresql
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"

# Check Redis
sudo systemctl status redis
redis-cli -h localhost -p 36380 -a redis_password ping

# If PostgreSQL down, service won't start
sudo systemctl start postgresql
sleep 5
sudo systemctl start listings-service
```

**Scenario D: Out of Memory (OOM killed)**
```bash
# Verify OOM kill
sudo dmesg | grep listings-service | grep -i "killed"

# Check systemd OOM settings
systemctl show listings-service | grep -i memory

# Restart with monitoring
sudo systemctl start listings-service
watch -n 2 'ps aux | grep listings-service'
```

#### Prevention

- Set up proper systemd restart policies (already configured: `Restart=on-failure`)
- Monitor memory usage and set alerts
- Implement graceful shutdown handling
- Regular health check monitoring (every 10 seconds)

#### Escalation Path

- **If restart fails 3 times:** Page Platform Team Lead
- **If database issue:** Page Database SRE team
- **If hardware issue:** Page Infrastructure team

---

### 4. Database Connection Pool Exhausted

#### Symptoms

- Alert: `ListingsDBPoolExhausted`
- Logs: "connection pool exhausted" or "too many connections"
- Metrics: `listings_db_connections_open` >= 25 (max configured)
- High request latency
- Timeouts on database operations

#### Investigation Steps

**Step 1: Check connection pool metrics**
```bash
# Current connections
curl -s http://localhost:8086/metrics | grep listings_db_connections

# PostgreSQL side
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT count(*), state FROM pg_stat_activity
   WHERE datname='listings_db' GROUP BY state;"
```

**Step 2: Identify long-running transactions**
```bash
# Active transactions
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pid, now() - xact_start AS duration, state, query
   FROM pg_stat_activity
   WHERE datname = 'listings_db' AND xact_start IS NOT NULL
   ORDER BY duration DESC LIMIT 20;"
```

**Step 3: Check for connection leaks**
```bash
# Idle connections
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT count(*) FROM pg_stat_activity
   WHERE datname='listings_db' AND state='idle'
   AND now() - state_change > interval '5 minutes';"
```

#### Resolution Steps

**Immediate Action:**
```bash
# Kill idle connections (>5 minutes)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pg_terminate_backend(pid)
   FROM pg_stat_activity
   WHERE datname = 'listings_db'
   AND state = 'idle'
   AND now() - state_change > interval '5 minutes';"

# Restart service to reset connection pool
sudo systemctl restart listings-service

# Monitor recovery
watch -n 1 'curl -s http://localhost:8086/metrics | grep listings_db_connections'
```

**Long-term Fix:**
```bash
# Adjust connection pool settings in .env
# Edit /opt/listings-dev/.env
sudo nano /opt/listings-dev/.env

# Increase pool size (if system can handle it)
SVETULISTINGS_DB_MAX_OPEN_CONNS=50
SVETULISTINGS_DB_MAX_IDLE_CONNS=20
SVETULISTINGS_DB_CONN_MAX_LIFETIME=3m

# Restart service
sudo systemctl restart listings-service
```

#### Prevention

- Set appropriate connection pool limits based on load testing
- Configure connection max lifetime: `DB_CONN_MAX_LIFETIME=5m`
- Set connection max idle time: `DB_CONN_MAX_IDLE_TIME=10m`
- Monitor connection pool usage with alerts
- Regular connection pool analysis

#### Escalation Path

- **If issue persists after restart:** Page Database SRE team
- **If caused by slow queries:** Coordinate with Database SRE for query optimization
- **If caused by traffic spike:** Consider horizontal scaling

---

### 5. Redis Connection Issues

#### Symptoms

- Alert: `ListingsRedisDown`
- Logs: "redis: connection refused" or "i/o timeout"
- Metrics: Rate limit metrics stop updating
- Service continues to operate (fail-open mode)

#### Investigation Steps

**Step 1: Check Redis status**
```bash
# Redis service status
sudo systemctl status redis

# Redis connectivity
redis-cli -h localhost -p 36380 -a redis_password ping

# Redis memory usage
redis-cli -h localhost -p 36380 -a redis_password info memory
```

**Step 2: Check Redis connections**
```bash
# Connected clients
redis-cli -h localhost -p 36380 -a redis_password info clients

# Check for connection errors
sudo journalctl -u redis --since "10 minutes ago" | grep -i error
```

**Step 3: Check network connectivity**
```bash
# Port check
telnet localhost 36380

# Network stats
netstat -an | grep 36380
```

#### Resolution Steps

**Scenario A: Redis service down**
```bash
# Start Redis
sudo systemctl start redis

# Verify connectivity
redis-cli -h localhost -p 36380 -a redis_password ping

# Check service logs (listings service will auto-reconnect)
sudo journalctl -u listings-service -f
```

**Scenario B: Redis out of memory**
```bash
# Check memory usage
redis-cli -h localhost -p 36380 -a redis_password info memory | grep used_memory_human

# Clear rate limit keys
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'rate_limit:*' | \
  xargs redis-cli -h localhost -p 36380 -a redis_password DEL

# Clear cache keys
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'cache:*' | \
  xargs redis-cli -h localhost -p 36380 -a redis_password DEL
```

**Scenario C: Redis connection pool exhausted**
```bash
# Check Redis max clients
redis-cli -h localhost -p 36380 -a redis_password config get maxclients

# Increase if needed
redis-cli -h localhost -p 36380 -a redis_password config set maxclients 10000

# Restart listings service
sudo systemctl restart listings-service
```

#### Prevention

- **Fail-open mode already implemented**: Service continues without Redis
- Monitor Redis memory usage: Alert at 80%
- Set maxmemory policy: `maxmemory-policy allkeys-lru`
- Regular Redis monitoring
- Configure Redis persistence: `appendonly yes`

#### Impact Assessment

- **Rate Limiting:** Disabled during Redis outage (fail-open)
- **Caching:** Disabled, all requests hit database (increased latency)
- **Service Availability:** Not affected
- **SLO Impact:** Latency may increase, error rate should remain normal

#### Escalation Path

- **If Redis down > 10 minutes:** Page Platform Team Lead
- **If memory issue persists:** Consider Redis instance upgrade
- **If impacting SLO:** Execute capacity scaling plan

---

### 6. OpenSearch Cluster Red

#### Symptoms

- Alert: `OpenSearchClusterRed`
- Search requests failing
- Logs: "opensearch: no available connection" or "cluster_block_exception"
- CRUD operations continue normally (OpenSearch only used for search)

#### Investigation Steps

**Step 1: Check cluster health**
```bash
# Cluster health
curl -u admin:admin http://localhost:9200/_cluster/health?pretty

# Node status
curl -u admin:admin http://localhost:9200/_cat/nodes?v

# Index status
curl -u admin:admin http://localhost:9200/_cat/indices?v | grep listings
```

**Step 2: Check shard allocation**
```bash
# Unassigned shards
curl -u admin:admin http://localhost:9200/_cat/shards?v | grep UNASSIGNED

# Cluster allocation explain
curl -u admin:admin -X GET "http://localhost:9200/_cluster/allocation/explain?pretty"
```

**Step 3: Check logs**
```bash
# OpenSearch logs
sudo journalctl -u opensearch --since "30 minutes ago" | grep -E "error|exception"

# Listings service logs (OpenSearch errors)
sudo journalctl -u listings-service --since "30 minutes ago" | grep opensearch
```

#### Resolution Steps

**Scenario A: Disk space issue**
```bash
# Check disk usage
curl -u admin:admin http://localhost:9200/_cat/allocation?v

# Clear old indices if needed
curl -u admin:admin -X DELETE "http://localhost:9200/old_index_name"

# Adjust cluster settings
curl -u admin:admin -X PUT "http://localhost:9200/_cluster/settings" \
  -H 'Content-Type: application/json' -d'
{
  "transient": {
    "cluster.routing.allocation.disk.watermark.low": "90%",
    "cluster.routing.allocation.disk.watermark.high": "95%"
  }
}'
```

**Scenario B: Shard allocation issues**
```bash
# Retry shard allocation
curl -u admin:admin -X POST "http://localhost:9200/_cluster/reroute?retry_failed=true"

# Force allocation of unassigned shards (DANGEROUS - use with caution)
# This may cause data loss
# curl -u admin:admin -X POST "http://localhost:9200/_cluster/reroute" \
#   -H 'Content-Type: application/json' -d'
# {
#   "commands": [
#     {
#       "allocate_empty_primary": {
#         "index": "listings_microservice",
#         "shard": 0,
#         "node": "node-name",
#         "accept_data_loss": true
#       }
#     }
#   ]
# }'
```

**Scenario C: Index corruption**
```bash
# Check index integrity
curl -u admin:admin -X POST "http://localhost:9200/listings_microservice/_flush"

# Reindex from PostgreSQL (safe recovery)
cd /p/github.com/sveturs/listings
python3 scripts/reindex_via_docker.py --target-password admin

# Validate reindexing
python3 scripts/validate_opensearch.py --target-password admin
```

#### Prevention

- Monitor disk space: Alert at 80%
- Regular index maintenance
- Snapshot and restore policy
- Monitor shard allocation
- Set up replica shards for high availability

#### Impact Assessment

- **Search Functionality:** Completely unavailable
- **CRUD Operations:** Not affected (direct database access)
- **Service Availability:** Not affected
- **SLO Impact:** Search SLO violated, overall service SLO maintained

#### Escalation Path

- **If cluster red > 5 minutes:** Begin reindexing from PostgreSQL
- **If cluster red > 15 minutes:** Page Platform Team Lead
- **If data loss suspected:** Page Database SRE team for backup recovery

---

### 7. Memory Leak / OOM Killed

#### Symptoms

- Alert: `ListingsHighMemoryUsage`
- Service process killed by OOM killer
- Logs: systemd shows "Main process exited, code=killed, status=9/KILL"
- Metrics: Memory usage continuously increasing
- Service restarts frequently

#### Investigation Steps

**Step 1: Verify OOM kill**
```bash
# Check dmesg for OOM killer
sudo dmesg -T | grep -E "listings-service|Out of memory"

# Check systemd logs
sudo journalctl -u listings-service | grep -i "killed"

# Check memory cgroup limits
systemctl show listings-service | grep Memory
```

**Step 2: Analyze memory usage patterns**
```bash
# Current memory usage
ps aux | grep listings-service | awk '{print "RSS: " $6/1024 " MB, VSZ: " $5/1024 " MB"}'

# Memory usage history (if available)
# Check Grafana dashboard: "Listings Service Memory Usage"

# Generate heap profile
curl http://localhost:8086/debug/pprof/heap > /tmp/heap-$(date +%s).prof

# Analyze heap
go tool pprof -top /tmp/heap-*.prof
go tool pprof -web /tmp/heap-*.prof  # Opens in browser
```

**Step 3: Check for goroutine leaks**
```bash
# Goroutine count
curl http://localhost:8086/debug/pprof/goroutine?debug=1 | grep "goroutine profile" | head -1

# Detailed goroutine dump
curl http://localhost:8086/debug/pprof/goroutine?debug=2 > /tmp/goroutines.txt

# Count by function
curl http://localhost:8086/debug/pprof/goroutine?debug=1 | grep -oP '# 0x[0-9a-f]+ \K.*' | sort | uniq -c | sort -rn | head -20
```

#### Resolution Steps

**Immediate Action:**
```bash
# Restart service
sudo systemctl restart listings-service

# Monitor memory growth
watch -n 5 'ps aux | grep listings-service | awk "{print \$6/1024 \" MB\"}"'

# Set temporary memory limit (16GB)
sudo systemctl set-property listings-service MemoryMax=16G
sudo systemctl restart listings-service
```

**Investigation Mode:**
```bash
# Enable memory profiling for analysis
# Add to service startup: GODEBUG=gctrace=1

# Collect heap profiles over time
for i in {1..10}; do
  curl http://localhost:8086/debug/pprof/heap > /tmp/heap-$i.prof
  echo "Captured profile $i"
  sleep 300  # 5 minutes between captures
done

# Compare profiles to identify leak
go tool pprof -base /tmp/heap-1.prof /tmp/heap-10.prof
```

**Common Leak Sources:**
```bash
# Check Redis connection pool
curl -s http://localhost:8086/metrics | grep listings_redis_pool

# Check database connection pool
curl -s http://localhost:8086/metrics | grep listings_db_connections

# Check HTTP connections
netstat -an | grep :8086 | wc -l

# Check OpenSearch connections
curl -u admin:admin http://localhost:9200/_nodes/stats/http | jq '.nodes[].http.current_open'
```

#### Prevention

- Implement memory limits in systemd: `MemoryMax=8G`
- Set `GOMEMLIMIT` environment variable
- Regular memory profiling in staging
- Connection pool limits strictly enforced
- Proper resource cleanup in code
- Monitor goroutine count

#### Known Memory Issues

1. **HTTP Response Body Not Closed**: Fixed in Sprint 4.2.1
2. **OpenSearch Client Connection Pooling**: Properly configured
3. **Redis Connection Leaks**: Connection pool implemented
4. **Large Result Set Buffering**: Pagination enforced

#### Escalation Path

- **If OOM occurs > 2 times/hour:** Page Platform Team Lead immediately
- **If memory > 80% consistently:** Begin investigation and profiling
- **If leak identified in code:** Create hotfix and deploy immediately

---

### 8. Disk Space Critical

#### Symptoms

- Alert: `ListingsDiskSpaceLow`
- Disk usage > 90%
- Logs: "no space left on device"
- Service unable to write logs
- Database operations failing

#### Investigation Steps

**Step 1: Check disk usage**
```bash
# Overall disk usage
df -h

# Service directory usage
du -sh /opt/listings-dev/*

# Log directory usage
du -sh /var/log/journal/*
du -sh /var/log/postgresql/*

# Find large files
find /opt/listings-dev -type f -size +100M -exec ls -lh {} \;
```

**Step 2: Identify disk consumers**
```bash
# Top 20 largest directories
du -h / --max-depth=3 2>/dev/null | sort -rh | head -20

# Recently created large files
find / -type f -size +100M -mtime -1 -exec ls -lh {} \; 2>/dev/null
```

**Step 3: Check logs and rotations**
```bash
# Journal logs size
journalctl --disk-usage

# PostgreSQL logs
ls -lh /var/log/postgresql/

# Check logrotate config
cat /etc/logrotate.d/listings-service
```

#### Resolution Steps

**Immediate Action (Free Up Space):**
```bash
# Clean systemd journal (keep 7 days)
sudo journalctl --vacuum-time=7d

# Clean old PostgreSQL logs
sudo find /var/log/postgresql/ -name "*.log" -mtime +7 -delete

# Clean old heap profiles
sudo find /tmp -name "*.prof" -mtime +1 -delete
sudo find /opt/listings-dev -name "*.prof" -mtime +1 -delete

# Clean package cache
sudo apt clean

# Verify space freed
df -h
```

**PostgreSQL Cleanup:**
```bash
# Vacuum database (reclaim space)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "VACUUM FULL ANALYZE;"

# Check table sizes
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
   FROM pg_tables
   WHERE schemaname = 'public'
   ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
   LIMIT 10;"
```

**Log Rotation Setup:**
```bash
# Configure journald limits
sudo nano /etc/systemd/journald.conf
# Set: SystemMaxUse=2G

# Restart journald
sudo systemctl restart systemd-journald

# Configure PostgreSQL log rotation
sudo nano /etc/logrotate.d/postgresql-common
# Ensure: rotate 7, maxage 7
```

#### Prevention

- Set up disk space monitoring: Alert at 80%
- Configure log rotation for all services
- Set journald limits: `SystemMaxUse=2G`
- Regular cleanup cron jobs
- Monitor database growth trends
- Implement data retention policies

#### Escalation Path

- **If disk > 95%:** Execute emergency cleanup immediately
- **If database affected:** Page Database SRE team
- **If requires disk expansion:** Page Infrastructure team

---

### 9. Rate Limit Abuse

#### Symptoms

- Alert: `ListingsHighRateLimitRejections`
- Metrics: `listings_rate_limit_rejected_total` increasing rapidly
- High rejection rate from specific IPs
- Logs: Frequent "rate limit exceeded" warnings

#### Investigation Steps

**Step 1: Identify abusive sources**
```bash
# Check rate limit metrics
curl -s http://localhost:8086/metrics | grep listings_rate_limit

# Check Redis for top rate-limited keys
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'rate_limit:*' | \
  while read key; do
    echo "$key: $(redis-cli -h localhost -p 36380 -a redis_password GET $key)"
  done | sort -t: -k2 -n

# Extract IPs from logs (last 10 minutes)
sudo journalctl -u listings-service --since "10 minutes ago" | \
  grep "rate limit exceeded" | \
  grep -oP 'identifier":\K"[^"]*"' | \
  sort | uniq -c | sort -rn | head -20
```

**Step 2: Analyze request patterns**
```bash
# Check request rate per endpoint
curl -s http://localhost:8086/metrics | \
  grep listings_grpc_requests_total | \
  grep -v "#" | \
  sort -t= -k2 -rn | head -20

# Check top methods being rate-limited
curl -s http://localhost:8086/metrics | \
  grep listings_rate_limit_rejected_total | \
  grep -v "#"
```

**Step 3: Verify legitimate vs malicious traffic**
```bash
# Check if requests are authenticated
sudo journalctl -u listings-service --since "5 minutes ago" | \
  grep "rate limit exceeded" | \
  jq -r 'select(.user_id != null) | .identifier' | \
  sort | uniq -c

# Geographic distribution (if available)
# Check nginx/reverse proxy logs for X-Forwarded-For
```

#### Resolution Steps

**Immediate Action (Block Abusive IP):**
```bash
# Option 1: Temporary block in iptables
sudo iptables -A INPUT -s 192.168.1.100 -j DROP

# Option 2: Block in nginx (if using reverse proxy)
sudo nano /etc/nginx/sites-available/listings.conf
# Add to http block:
#   deny 192.168.1.100;
sudo nginx -t && sudo systemctl reload nginx

# Option 3: Add to Redis blocklist
redis-cli -h localhost -p 36380 -a redis_password SET "blocklist:192.168.1.100" "1" EX 3600
```

**Adjust Rate Limits (If Legitimate Traffic):**
```bash
# Edit rate limit config
sudo nano /opt/listings-dev/internal/ratelimit/config.go

# Example: Increase GetListing limit
# Change:  Limit: 200
# To:      Limit: 500

# Rebuild and restart
cd /opt/listings-dev
make build
sudo systemctl restart listings-service
```

**Clear Rate Limit Counters:**
```bash
# Clear all rate limits (allow retry)
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'rate_limit:*' | \
  xargs redis-cli -h localhost -p 36380 -a redis_password DEL

# Clear specific IP
redis-cli -h localhost -p 36380 -a redis_password DEL "rate_limit:/listings.v1.ListingsService/GetListing:192.168.1.100"
```

#### Prevention

- Implement IP allowlist/blocklist
- Use CDN with DDoS protection (Cloudflare)
- Progressive rate limiting (stricter for unauthenticated users)
- CAPTCHA for suspicious patterns
- Geo-blocking if applicable
- Monitor rate limit metrics continuously

#### Known Patterns

1. **Scraper Bots**: High frequency GetListing calls
2. **DDoS Attacks**: Distributed sources, random endpoints
3. **Legitimate Spikes**: Marketing campaigns, mobile app releases
4. **Misconfigured Clients**: Retry loops, missing backoff

#### Escalation Path

- **If rejection rate > 20%:** Review and adjust rate limits
- **If DDoS suspected:** Page Security team immediately
- **If legitimate traffic impacted:** Page Platform Team Lead

---

### 10. Slow Queries

#### Symptoms

- Alert: `ListingsSlowQueries`
- High database CPU usage
- Increased request latency
- Metrics: `listings_db_query_duration_seconds` high
- PostgreSQL logs showing slow queries

#### Investigation Steps

**Step 1: Identify slow queries**
```bash
# Check slow query log
sudo tail -100 /var/log/postgresql/postgresql-15-main.log | grep "duration:" | sort -t: -k2 -rn

# Active slow queries
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pid, now() - pg_stat_activity.query_start AS duration, query
   FROM pg_stat_activity
   WHERE state = 'active' AND datname = 'listings_db'
   AND now() - pg_stat_activity.query_start > interval '1 second'
   ORDER BY duration DESC;"

# Query statistics (requires pg_stat_statements)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT query, calls, total_exec_time, mean_exec_time, max_exec_time
   FROM pg_stat_statements
   ORDER BY mean_exec_time DESC
   LIMIT 20;"
```

**Step 2: Analyze query plans**
```bash
# Get query plan for slow query
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "EXPLAIN ANALYZE SELECT * FROM listings WHERE category_id = 1301 LIMIT 10;"

# Check for missing indexes
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT schemaname, tablename, indexname, idx_scan
   FROM pg_stat_user_indexes
   WHERE idx_scan = 0
   ORDER BY schemaname, tablename;"
```

**Step 3: Check database health**
```bash
# Bloat and fragmentation
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
   FROM pg_tables
   WHERE schemaname = 'public'
   ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"

# Lock contention
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT blocked_locks.pid AS blocked_pid,
          blocked_activity.usename AS blocked_user,
          blocking_locks.pid AS blocking_pid,
          blocking_activity.usename AS blocking_user,
          blocked_activity.query AS blocked_statement,
          blocking_activity.query AS current_statement_in_blocking_process
   FROM pg_catalog.pg_locks blocked_locks
   JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
   JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
   AND blocking_locks.pid != blocked_locks.pid
   JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
   WHERE NOT blocked_locks.granted;"
```

#### Resolution Steps

**Immediate Action:**
```bash
# Kill blocking query (if identified)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pg_terminate_backend(12345);"  # Replace with actual PID

# Run ANALYZE to update statistics
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "ANALYZE;"
```

**Index Creation:**
```bash
# Create missing index (example)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "CREATE INDEX CONCURRENTLY idx_listings_category_status
   ON listings(category_id, status)
   WHERE deleted_at IS NULL;"

# Verify index usage
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "EXPLAIN ANALYZE SELECT * FROM listings
   WHERE category_id = 1301 AND status = 'active' AND deleted_at IS NULL
   LIMIT 10;"
```

**Database Maintenance:**
```bash
# Vacuum (reclaim space and update statistics)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "VACUUM ANALYZE;"

# Reindex (if fragmentation detected)
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "REINDEX DATABASE listings_db;"
```

#### Prevention

- Regular VACUUM ANALYZE (weekly)
- Monitor query performance continuously
- Review query plans during code review
- Add indexes based on query patterns
- Set slow query log threshold: `log_min_duration_statement = 1000ms`
- Enable pg_stat_statements extension

#### Escalation Path

- **If slow queries persist > 5 minutes:** Page Database SRE team
- **If requires schema changes:** Coordinate with Platform Team Lead
- **If impacting SLO:** Consider read replica or caching

---

## Recovery Procedures

### Service Restart Checklist

```bash
# 1. Backup current state
curl http://localhost:8086/metrics > /tmp/metrics-pre-restart.txt

# 2. Graceful shutdown
sudo systemctl stop listings-service

# 3. Verify clean shutdown
ps aux | grep listings-service

# 4. Check dependencies
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"
redis-cli -h localhost -p 36380 -a redis_password ping

# 5. Start service
sudo systemctl start listings-service

# 6. Monitor startup
sudo journalctl -u listings-service -f

# 7. Wait for readiness
sleep 10

# 8. Verify health
curl http://localhost:8086/health
curl http://localhost:8086/ready

# 9. Test endpoints
grpcurl -plaintext -d '{"id": 1}' localhost:50053 listings.v1.ListingsService/GetListing

# 10. Monitor metrics
watch -n 2 'curl -s http://localhost:8086/metrics | grep -E "listings_grpc_requests_total|listings_errors_total"'
```

### Database Recovery

See [DISASTER_RECOVERY.md](./DISASTER_RECOVERY.md) for complete procedures.

### Cache Invalidation

```bash
# Clear all Redis cache
redis-cli -h localhost -p 36380 -a redis_password FLUSHALL

# Clear specific cache patterns
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'cache:listing:*' | \
  xargs redis-cli -h localhost -p 36380 -a redis_password DEL
```

---

## Escalation

### When to Escalate

| Situation | Escalate To | When |
|-----------|-------------|------|
| Service down > 5 minutes | Platform Team Lead | Immediately |
| Database issue | Database SRE | Immediately |
| Data loss suspected | Platform Team Lead + Database SRE | Immediately |
| Security incident | Security Team | Immediately |
| SLO breach | Platform Team Lead | Within 15 minutes |
| Performance degradation | Platform Team Lead | After 30 minutes |
| Rate limit abuse (DDoS) | Security Team | Immediately |

### Escalation Contacts

```bash
# PagerDuty (primary)
listings-oncall     # On-call engineer (Level 1)
db-team             # Database SRE (Level 2)
platform-team-lead  # Platform Team Lead (Level 3)
security-team       # Security Team (parallel)

# Slack (secondary)
#listings-incidents  # Incident coordination
#listings-alerts     # Alert notifications
#platform-team       # Platform team channel

# Email (tertiary)
oncall@svetu.rs     # On-call rotation
platform@svetu.rs   # Platform team
security@svetu.rs   # Security team
```

### Incident Communication Template

```
INCIDENT: [SEVERITY] Listings Microservice [ISSUE]

STATUS: [INVESTIGATING | IDENTIFIED | RESOLVING | RESOLVED]
START TIME: [TIMESTAMP]
IMPACT: [DESCRIPTION]
AFFECTED: [USERS/ENDPOINTS]

TIMELINE:
- HH:MM - [EVENT]
- HH:MM - [ACTION]

CURRENT ACTIONS:
- [ACTION 1]
- [ACTION 2]

NEXT STEPS:
- [STEP 1]
- [STEP 2]

ESCALATED TO: [TEAM/PERSON]
ETA FOR RESOLUTION: [TIME]
```

---

## Post-Incident Procedures

1. **Document incident** in incident management system
2. **Restore monitoring** and verify all metrics normal
3. **Schedule postmortem** within 48 hours
4. **Create action items** for prevention
5. **Update runbook** with lessons learned

---

**Document Version:** 1.0.0
**Last Reviewed:** 2025-11-05
**Next Review:** 2025-12-05
**Owner:** Platform Team
**Contributors:** SRE Team, Database Team
