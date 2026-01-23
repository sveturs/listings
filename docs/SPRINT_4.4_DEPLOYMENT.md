# Sprint 4.4 - dev.vondi.rs Deployment Setup

**Status**: ✅ Complete
**Phase**: 4 - Deployment Infrastructure
**Sprint**: 4.4 - dev.vondi.rs Deployment Setup
**Duration**: 8 hours
**Date**: 2025-10-31

## Overview

This document describes the complete deployment infrastructure for listings-service on dev.vondi.rs server.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        dev.vondi.rs Server                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Nginx (Port 443/80)                                      │  │
│  │  - listings.dev.vondi.rs → localhost:8086 (HTTP REST)    │  │
│  └────────────────────┬─────────────────────────────────────┘  │
│                       │                                          │
│  ┌────────────────────▼─────────────────────────────────────┐  │
│  │  Listings Service (systemd)                              │  │
│  │  - HTTP REST API: 8086                                   │  │
│  │  - gRPC API: 50053 (internal)                            │  │
│  │  - Metrics: 9093 (internal)                              │  │
│  │  - Binary: /opt/listings-dev/bin/listings-service       │  │
│  └──────────┬──────────────┬──────────────┬─────────────────┘  │
│             │              │              │                     │
│  ┌──────────▼────┐  ┌──────▼─────┐  ┌────▼──────────┐         │
│  │  PostgreSQL   │  │   Redis     │  │  Auth Service │         │
│  │  Port: 35433  │  │  Port:36380 │  │  Port: 28086  │         │
│  │  (Docker)     │  │  (Docker)   │  │  (preprod)    │         │
│  └───────────────┘  └─────────────┘  └───────────────┘         │
│                                                                  │
│  ┌──────────────────┐  ┌─────────────────────────────────────┐│
│  │   OpenSearch     │  │         MinIO                        ││
│  │   Port: 9200     │  │         Port: 9000                   ││
│  │   (shared)       │  │         (shared, S3-compatible)      ││
│  └──────────────────┘  └─────────────────────────────────────┘│
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## File Structure

```
listings/
├── scripts/
│   └── deploy-to-dev.sh              # Automated deployment script
├── deployment/
│   ├── listings-service.service      # systemd service file
│   └── nginx-listings.conf           # Nginx reverse proxy config
├── .env.prod.example                 # Production environment template
└── docs/
    └── SPRINT_4.4_DEPLOYMENT.md      # This file
```

## Deployment Components

### 1. Deploy Script (`scripts/deploy-to-dev.sh`)

Automated deployment script that handles:

- ✅ Git commit & push
- ✅ Local binary build (`make build`)
- ✅ Upload binary to server
- ✅ Upload configuration files
- ✅ Restart dependencies (PostgreSQL, Redis)
- ✅ Run database migrations
- ✅ Restart systemd service
- ✅ Health check validation

**Usage**:

```bash
cd /p/github.com/sveturs/listings
./scripts/deploy-to-dev.sh
```

**Features**:

- Color-coded logging (green/yellow/red/blue)
- Error handling with rollback instructions
- Health checks with retries (6 attempts, 10s interval)
- Service status monitoring
- Automatic service restart

### 2. Systemd Service (`deployment/listings-service.service`)

Production-ready systemd service configuration:

**Key Features**:

- ✅ Automatic restart on failure (10s delay)
- ✅ Proper dependencies (PostgreSQL, Redis, network)
- ✅ Resource limits (65536 files, 4096 processes)
- ✅ Security hardening (NoNewPrivileges, PrivateTmp, ProtectSystem)
- ✅ Graceful shutdown (30s timeout)
- ✅ Journal logging with syslog identifier

**Installation**:

```bash
# On server (done automatically by deploy script)
sudo cp deployment/listings-service.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable listings-service
sudo systemctl start listings-service
```

**Management**:

```bash
# Check status
sudo systemctl status listings-service

# View logs
sudo journalctl -u listings-service -f

# Restart service
sudo systemctl restart listings-service

# Stop service
sudo systemctl stop listings-service
```

### 3. Nginx Configuration (`deployment/nginx-listings.conf`)

Reverse proxy configuration for HTTP REST API:

**Features**:

- ✅ HTTP to HTTPS redirect
- ✅ SSL/TLS with certbot support
- ✅ Security headers (HSTS, X-Frame-Options, CSP)
- ✅ Proper proxy headers (Host, X-Real-IP, X-Forwarded-*)
- ✅ Health check endpoint (no caching)
- ✅ Client limits (50MB max body size)
- ✅ Timeouts (60s connect/send/read)

**Installation**:

```bash
# On server
sudo cp deployment/nginx-listings.conf /etc/nginx/sites-available/listings-dev
sudo ln -s /etc/nginx/sites-available/listings-dev /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Setup SSL certificate
sudo certbot --nginx -d listings.dev.vondi.rs
```

**Important Notes**:

- ✅ HTTP REST API exposed via Nginx: `listings.dev.vondi.rs`
- ❌ gRPC (port 50053) - INTERNAL ONLY (no Nginx exposure)
- ❌ Metrics (port 9093) - INTERNAL ONLY (security risk if exposed)

### 4. Production Environment (`.env.prod.example`)

Template for production configuration:

**Key Sections**:

1. **Application Settings**: ENV, log level, log format
2. **Server Ports**: gRPC (50053), HTTP (8086), Metrics (9093)
3. **Database**: Separate PostgreSQL (listings_dev_db, port 35433)
4. **Redis**: Separate instance (port 36380)
5. **OpenSearch**: Shared instance (marketplace_listings index)
6. **MinIO**: Shared S3 storage (listings-images bucket)
7. **Auth Service**: Preprod instance (http://localhost:28086)
8. **Worker**: Async indexing configuration
9. **CORS**: Dev frontend origins
10. **Feature Flags**: Async indexing, image optimization, caching

**Setup**:

```bash
# Copy template
cp .env.prod.example .env.prod

# Edit with production values
vim .env.prod

# IMPORTANT: Never commit .env.prod to git!
```

## Server Setup

### Prerequisites

1. **Server Access**:

```bash
ssh svetu@vondi.rs
```

2. **Create Deployment Directory**:

```bash
sudo mkdir -p /opt/listings-dev
sudo chown svetu:svetu /opt/listings-dev
cd /opt/listings-dev

# Clone repository (if not exists)
git clone git@github.com:sveturs/listings.git .
```

3. **Setup Environment**:

```bash
# Copy production env
cp .env.prod.example .env.prod
vim .env.prod  # Configure with production values
```

4. **Create Database**:

```bash
# Connect to PostgreSQL (from Docker)
docker exec -it svetu-dev_db_1 psql -U svetu_dev_user -d postgres

# Create listings database
CREATE DATABASE listings_dev_db;
CREATE USER listings_user WITH PASSWORD 'STRONG_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE listings_dev_db TO listings_user;
\q
```

5. **Setup Redis**:

```bash
# Start Redis (via docker-compose)
cd /opt/listings-dev
docker-compose up -d redis
```

6. **Install Nginx Config**:

```bash
sudo cp deployment/nginx-listings.conf /etc/nginx/sites-available/listings-dev
sudo ln -s /etc/nginx/sites-available/listings-dev /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Setup SSL
sudo certbot --nginx -d listings.dev.vondi.rs
```

## Deployment Process

### Automated Deployment (Recommended)

```bash
# From local machine
cd /p/github.com/sveturs/listings
./scripts/deploy-to-dev.sh
```

The script will:

1. Commit and push changes
2. Build binary locally
3. Upload to server
4. Restart dependencies
5. Run migrations
6. Restart service
7. Validate health

### Manual Deployment

If you need to deploy manually:

```bash
# 1. Build locally
make build

# 2. Upload to server
scp bin/listings-service svetu@vondi.rs:/opt/listings-dev/bin/
scp docker-compose.yml svetu@vondi.rs:/opt/listings-dev/
scp .env.prod svetu@vondi.rs:/opt/listings-dev/.env

# 3. On server
ssh svetu@vondi.rs
cd /opt/listings-dev

# 4. Restart dependencies
docker-compose up -d postgres redis

# 5. Run migrations
make migrate-up

# 6. Restart service
sudo systemctl restart listings-service

# 7. Check health
curl http://localhost:8086/health
curl http://localhost:9093/metrics
```

## Verification

### 1. Service Status

```bash
# On server
sudo systemctl status listings-service
```

Expected output:

```
● listings-service.service - Listings Microservice (dev)
   Loaded: loaded (/etc/systemd/system/listings-service.service; enabled)
   Active: active (running) since ...
```

### 2. Health Checks

```bash
# HTTP REST API
curl http://localhost:8086/health
# Expected: {"status":"ok"}

# Metrics
curl http://localhost:9093/metrics
# Expected: Prometheus metrics output

# Public HTTPS
curl https://listings.dev.vondi.rs/health
# Expected: {"status":"ok"}
```

### 3. Logs

```bash
# Service logs
sudo journalctl -u listings-service -f

# Nginx logs
sudo tail -f /var/log/nginx/listings-dev-access.log
sudo tail -f /var/log/nginx/listings-dev-error.log
```

### 4. Process Verification

```bash
# Check process
ps aux | grep listings-service

# Check ports
sudo netstat -tlnp | grep -E "8086|50053|9093"
```

Expected:

```
tcp        0      0 0.0.0.0:8086            0.0.0.0:*               LISTEN      12345/listings-service
tcp        0      0 0.0.0.0:50053           0.0.0.0:*               LISTEN      12345/listings-service
tcp        0      0 0.0.0.0:9093            0.0.0.0:*               LISTEN      12345/listings-service
```

## Troubleshooting

### Service Won't Start

**Check logs**:

```bash
sudo journalctl -u listings-service -n 100 --no-pager
```

**Common issues**:

1. **Database connection failed**:

```bash
# Check PostgreSQL
docker ps | grep listings_postgres
docker exec listings_postgres pg_isready -U listings_user
```

2. **Redis connection failed**:

```bash
# Check Redis
docker ps | grep listings_redis
docker exec listings_redis redis-cli ping
```

3. **Port already in use**:

```bash
# Check what's using the port
sudo lsof -i :8086
sudo kill -9 <PID>
```

### Health Checks Failing

**Check service is running**:

```bash
sudo systemctl status listings-service
curl http://localhost:8086/health
```

**Check Nginx**:

```bash
sudo nginx -t
sudo systemctl status nginx
curl -I http://localhost:8086
```

### Deployment Script Fails

**Check SSH access**:

```bash
ssh svetu@vondi.rs echo "OK"
```

**Check build**:

```bash
make clean
make build
./bin/listings-service --version
```

**Check server disk space**:

```bash
ssh svetu@vondi.rs "df -h /opt/listings-dev"
```

## Rollback Procedure

If deployment fails or causes issues:

### 1. Quick Rollback (Service Only)

```bash
# On server
cd /opt/listings-dev

# Restore previous binary (if backup exists)
cp bin/listings-service.backup bin/listings-service

# Restart service
sudo systemctl restart listings-service
```

### 2. Full Rollback (Git + Service)

```bash
# On server
cd /opt/listings-dev

# Find previous commit
git log --oneline -10

# Reset to previous commit
git reset --hard <PREVIOUS_COMMIT>

# Rebuild
make build

# Restart service
sudo systemctl restart listings-service
```

### 3. Database Rollback

```bash
# Rollback last migration
make migrate-down

# Or reset to specific version
migrate -path migrations -database "$DATABASE_URL" force <VERSION>
```

## Monitoring

### Service Metrics

Prometheus metrics available at: `http://localhost:9093/metrics`

**Key metrics**:

- `listings_http_requests_total` - HTTP request count
- `listings_http_request_duration_seconds` - Request latency
- `listings_grpc_requests_total` - gRPC request count
- `listings_db_queries_total` - Database query count
- `listings_cache_hits_total` - Cache hit rate

### Alerting

Set up Prometheus alerts for:

- Service down (health check fails)
- High error rate (>5%)
- High latency (p95 > 1s)
- Database connection errors
- Redis connection errors

## Security

### Hardening Checklist

- ✅ systemd security features enabled
- ✅ Non-root user (svetu)
- ✅ Firewall rules (UFW)
- ✅ SSL/TLS with certbot
- ✅ Security headers in Nginx
- ✅ gRPC and Metrics NOT exposed publicly
- ✅ Strong passwords in .env.prod
- ✅ Read-only file system where possible

### Firewall Rules

```bash
# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Block direct access to application ports (only allow localhost)
sudo ufw deny 8086/tcp
sudo ufw deny 50053/tcp
sudo ufw deny 9093/tcp
```

## Performance

### Expected Performance

- **HTTP REST**: 1000+ RPS
- **gRPC**: 5000+ RPS
- **Latency**: p95 < 100ms
- **Memory**: ~200-500MB
- **CPU**: ~10-30% (under load)

### Optimization Tips

1. **Increase connection pools** (if needed):

```bash
VONDILISTINGS_DB_MAX_OPEN_CONNS=100
VONDILISTINGS_REDIS_POOL_SIZE=50
```

2. **Enable caching**:

```bash
VONDILISTINGS_FEATURE_CACHE_ENABLED=true
VONDILISTINGS_CACHE_LISTING_TTL=15m
```

3. **Increase worker concurrency**:

```bash
VONDILISTINGS_WORKER_CONCURRENCY=20
```

## Future Improvements

### Phase 5 (Next Steps)

- [ ] Production deployment (prod.vondi.rs)
- [ ] Multi-instance deployment (load balancing)
- [ ] Blue-green deployment
- [ ] Canary releases
- [ ] Automated rollback on failure
- [ ] Enhanced monitoring (Grafana dashboards)
- [ ] Log aggregation (ELK/Loki)
- [ ] Distributed tracing (Jaeger)

## Conclusion

Sprint 4.4 provides complete deployment infrastructure for listings-service on dev.vondi.rs:

✅ **Automated deployment** via `deploy-to-dev.sh`
✅ **systemd service** with proper dependencies and restart policy
✅ **Nginx reverse proxy** with SSL/TLS and security headers
✅ **Production environment** template with comprehensive settings
✅ **Complete documentation** with troubleshooting and rollback procedures

The deployment is production-ready and can be used as a template for future prod deployment.

---

**Sprint Status**: ✅ **COMPLETE**
**Next Sprint**: Phase 5 - Production Hardening & Monitoring
**Deployment URL**: https://listings.dev.vondi.rs
