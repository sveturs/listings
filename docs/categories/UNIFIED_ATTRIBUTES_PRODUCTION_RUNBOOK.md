# Unified Attributes System - Production Runbook
*Версия: 1.0*
*Дата: 03.09.2025*
*Статус: Production Ready*

## 📋 Обзор системы

### Unified Attributes Architecture
- **Backend**: Go-based API на порту 3002
- **Database**: PostgreSQL в Docker контейнере (порт 5433)
- **Cache**: Redis на порту 6380
- **Search**: OpenSearch на порту 9201
- **Monitoring**: Prometheus (9090) + Grafana (3001)

### Production Environment
- **Server**: svetu.rs (Ubuntu, 15GB RAM, 125GB доступно)
- **Uptime**: 80+ дней
- **Load**: 0.19 average (low load, stable)

## 🚀 Daily Operations

### Cache Management

#### Cache Warmup (ежедневно)
```bash
cd /opt/cache-strategy
REDIS_ADDR=localhost:6380 go run day20-cache-strategy.go warmup
```

#### Cache Statistics Check
```bash
cd /opt/cache-strategy
REDIS_ADDR=localhost:6380 go run day20-cache-strategy.go stats
```

### Performance Monitoring

#### API Performance Check
```bash
# Categories endpoint
ab -n 100 -c 5 http://localhost:3002/api/v1/marketplace/categories

# Search endpoint  
ab -n 50 -c 3 'http://localhost:3002/api/v1/marketplace/search?q=car'
```

#### Expected Performance Benchmarks
- **Categories API**: >800 req/sec, <10ms avg response
- **Search API**: >300 req/sec, <15ms avg response
- **Cache Hit Rate**: >70% после warmup
- **Error Rate**: <1%

## 🔧 Troubleshooting Guide

### Common Issues

#### 1. High API Response Time (>20ms)
**Solutions**:
1. Warmup cache: `REDIS_ADDR=localhost:6380 go run day20-cache-strategy.go warmup`
2. Restart backend if unhealthy: `docker-compose restart backend`
3. Check database performance

#### 2. Cache Hit Rate <50%
**Solutions**:
1. Clear and re-warmup cache
2. Check Redis memory limits
3. Verify cache key patterns

### Emergency Procedures

#### Backend Service Restart
```bash
cd /app && docker-compose restart backend
curl http://localhost:3002/api/v1/marketplace/categories
```

## 📊 Current Production Metrics (Day 21)
- **Categories API**: 884 req/sec (baseline)
- **Search API**: 380 req/sec, 13ms avg response time
- **Cache Keys**: 7 in unified_attrs namespace
- **System Load**: 0.19 average

---

**Document Status**: Production Ready
**Last Updated**: 03.09.2025