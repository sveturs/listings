# Phase 28: Redis Caching Strategy for Search Analytics

## Overview

This document defines the Redis caching strategy for Search Analytics queries to achieve:
- **Target:** 80-85% cache hit rate
- **Performance:** < 10ms for cached queries
- **Memory:** < 100 MB Redis usage
- **Reliability:** Graceful degradation if Redis unavailable

## Cache Key Design

### Key Naming Convention

```
{service}:{entity}:{identifier}[:{filter}]
```

**Examples:**
```
search:trending:category:1301:days:7
search:popular:category:1301
search:history:user:123
search:ctr:query:iphone:days:30
```

### Key Hierarchy

```
search:
  ├── trending:
  │   ├── category:{id}:days:{n}    # Category-specific trending
  │   └── global:days:{n}            # Global trending
  ├── popular:
  │   ├── category:{id}              # Category-specific popular
  │   └── global                     # Global popular
  ├── history:
  │   ├── user:{id}                  # User search history
  │   └── session:{id}               # Session search history
  └── ctr:
      └── query:{text}:days:{n}      # CTR analysis per query
```

## TTL Strategy

### TTL Values and Rationale

| Cache Type | TTL | Update Frequency | Rationale |
|------------|-----|------------------|-----------|
| **Trending Queries** | 15 min | High (every search) | Real-time feel, rapidly changing |
| **Popular Queries** | 1 hour | Medium (daily) | Stable over time, expensive query |
| **User History** | 5 min | High (every search) | Personalized, frequent updates |
| **Session History** | 5 min | High (every search) | Personalized, frequent updates |
| **CTR Analysis** | 30 min | Low (on click) | Admin dashboard, not time-critical |

### TTL Implementation

```go
const (
    TrendingQueriesTTL  = 15 * time.Minute
    PopularQueriesTTL   = 60 * time.Minute
    UserHistoryTTL      = 5 * time.Minute
    SessionHistoryTTL   = 5 * time.Minute
    CTRAnalysisTTL      = 30 * time.Minute
)

// SetWithTTL sets a cache value with appropriate TTL
func (c *SearchCache) SetTrending(ctx context.Context, key string, value interface{}) error {
    return c.redis.Set(ctx, key, value, TrendingQueriesTTL).Err()
}
```

## Cache Value Format

### JSON Serialization

All cached values are stored as JSON strings for:
- Human readability (Redis CLI debugging)
- Cross-language compatibility
- Schema flexibility

**Example: Trending Queries Cache**
```json
{
  "searches": [
    {
      "query_text": "iphone",
      "search_count": 1203,
      "last_searched": "2025-11-19T15:30:00Z",
      "category_id": 1001
    },
    {
      "query_text": "samsung",
      "search_count": 891,
      "last_searched": "2025-11-19T15:28:00Z",
      "category_id": 1001
    }
  ],
  "took_ms": 350,
  "cached_at": "2025-11-19T15:30:15Z"
}
```

**Size Estimate:** ~2-5 KB per cached result

## Cache Invalidation Strategy

### 1. Time-Based Invalidation (TTL)

**Primary Strategy:** Use Redis TTL for automatic expiration

**Pros:**
- Simple, no manual invalidation needed
- Prevents stale data automatically
- Works even if app crashes

**Cons:**
- May serve slightly stale data
- No immediate updates

### 2. Event-Based Invalidation (Optional)

**For critical real-time updates:**

```go
// When a new search is logged, invalidate trending cache
func (s *Service) logSearchQuery(ctx context.Context, req *SearchRequest, resp *SearchResponse) {
    // ... log to database ...

    // Invalidate trending cache for this category (optional)
    if req.CategoryID != nil {
        cacheKey := fmt.Sprintf("search:trending:category:%d:days:7", *req.CategoryID)
        s.cache.Del(ctx, cacheKey) // Force refresh on next request
    }
}
```

**When to Use:**
- High-traffic categories (> 100 searches/min)
- Real-time leaderboards
- A/B testing scenarios

**When NOT to Use:**
- Low-traffic categories (< 10 searches/min)
- Background analytics
- Historical data

### 3. Cache Warming Strategy

**Pre-populate cache before traffic spike:**

```go
// Warm cache for top 10 categories at midnight
func (s *Service) warmTrendingCache(ctx context.Context) error {
    topCategories := []int64{1001, 1002, 1301, /* ... */}

    for _, catID := range topCategories {
        filter := &domain.GetTrendingQueriesFilter{
            CategoryID: &catID,
            Limit:      10,
            DaysAgo:    7,
        }

        // Fetch from DB and cache
        _, err := s.GetTrendingQueries(ctx, filter)
        if err != nil {
            log.Warn().Err(err).Int64("category", catID).Msg("cache warming failed")
        }
    }

    return nil
}
```

**Schedule:** Daily at 00:00 UTC via cron job

## Cache Hit Rate Optimization

### Predicted Cache Hit Rates

**Without Optimization:**
- Trending Queries: 70-75%
- Popular Queries: 85-90%
- User History: 50-60%
- Overall: 65-70%

**With Optimization:**
- Trending Queries: 85-90%
- Popular Queries: 95%+
- User History: 60-70%
- Overall: **80-85%** ✅

### Optimization Techniques

#### 1. Query Normalization

**Problem:** Cache misses due to query variations
- "iPhone 13" vs "iphone 13" (case)
- "laptop " vs "laptop" (whitespace)

**Solution:** Normalize query_text before caching

```go
func normalizeQuery(query string) string {
    // Trim whitespace
    query = strings.TrimSpace(query)

    // Lowercase
    query = strings.ToLower(query)

    // Remove extra spaces
    query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

    return query
}

// Use normalized query for cache key
cacheKey := fmt.Sprintf("search:trending:category:%d:days:7:%s",
    categoryID, normalizeQuery(query))
```

#### 2. Category Grouping

**Problem:** Too many cache keys for subcategories

**Solution:** Cache at parent category level

```go
// Instead of caching subcategory 1301 (Cars)
// Cache parent category 1000 (Vehicles)
func getParentCategoryID(categoryID int64) int64 {
    // Fetch from categories table
    return parentID
}
```

#### 3. Wildcard TTL Extension

**For high-traffic queries, extend TTL dynamically:**

```go
// If query is popular, extend TTL to 30 min
if searchCount > 1000 {
    c.redis.Expire(ctx, key, 30*time.Minute)
}
```

## Memory Management

### Memory Limits

**Total Redis Memory Budget:** 1 GB
**Search Analytics Allocation:** 100 MB (10%)

**Per-Type Limits:**
- Trending Queries: 20 MB
- Popular Queries: 10 MB
- User History: 50 MB
- Session History: 15 MB
- CTR Analysis: 5 MB

### Eviction Policy

**Redis Configuration:**
```
maxmemory 1gb
maxmemory-policy allkeys-lru
```

**LRU (Least Recently Used):**
- Automatically evicts least accessed keys
- No manual cleanup needed
- Works well with TTL

### Memory Monitoring

**Key Metrics:**
```bash
# Total memory usage
redis-cli INFO memory | grep used_memory_human

# Key count by pattern
redis-cli --scan --pattern "search:trending:*" | wc -l

# Sample key size
redis-cli DEBUG OBJECT "search:trending:category:1301:days:7"
```

**Alerts:**
- Memory usage > 80 MB → Warning
- Memory usage > 95 MB → Critical
- Cache evictions/sec > 10 → Investigate

## Error Handling

### Graceful Degradation

**If Redis is unavailable:**

```go
func (c *SearchCache) Get(ctx context.Context, key string) (interface{}, error) {
    val, err := c.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        // Cache miss - normal
        return nil, ErrCacheMiss
    }
    if err != nil {
        // Redis error - log and continue without cache
        c.logger.Warn().Err(err).Str("key", key).Msg("redis error, bypassing cache")
        return nil, ErrCacheUnavailable
    }
    return val, nil
}

// In service layer
func (s *Service) GetTrendingQueries(ctx context.Context, filter *domain.GetTrendingQueriesFilter) ([]domain.TrendingSearch, error) {
    // Try cache first
    cached, err := s.cache.GetTrending(ctx, cacheKey)
    if err == nil {
        return cached, nil
    }

    // Cache miss or error - query database
    // Don't fail request if cache is down
    return s.searchQueriesRepo.GetTrendingQueries(ctx, filter)
}
```

### Retry Strategy

**Don't retry cache operations:**
- Cache is performance optimization, not critical path
- Retry adds latency
- Database query is fallback

**Exception:** Cache warming jobs (can retry 3x)

## Implementation Checklist

### Phase 1: Basic Caching (Week 1)

- ✅ Implement cache key generation
- ✅ Add TTL constants
- ✅ Create cache wrapper functions
- ✅ Add cache hit/miss logging

### Phase 2: Optimization (Week 2)

- ✅ Implement query normalization
- ✅ Add cache warming job
- ✅ Optimize category grouping
- ✅ Add memory monitoring

### Phase 3: Advanced Features (Week 3-4)

- ✅ Implement event-based invalidation
- ✅ Add dynamic TTL extension
- ✅ Create cache metrics dashboard
- ✅ Set up alerts

## Testing Strategy

### Unit Tests

```go
func TestTrendingQueriesCache(t *testing.T) {
    // Test cache hit
    // Test cache miss
    // Test TTL expiration
    // Test query normalization
}

func TestCacheGracefulDegradation(t *testing.T) {
    // Test with Redis unavailable
    // Test with Redis timeout
    // Test with malformed cache data
}
```

### Integration Tests

```go
func TestTrendingQueriesEndToEnd(t *testing.T) {
    // 1. Query trending (cache miss) - measure time
    // 2. Query trending (cache hit) - measure time
    // 3. Assert cache hit is < 10ms
    // 4. Assert cache miss is < 500ms
}
```

### Performance Tests

```bash
# Cache hit rate benchmark
go test -bench=BenchmarkTrendingQueries -benchtime=10s

# Memory usage benchmark
go test -bench=BenchmarkCacheMemory -benchmem
```

## Monitoring Dashboards

### Key Metrics

1. **Cache Hit Rate**
   - Formula: `hits / (hits + misses) * 100`
   - Target: > 80%
   - Alert: < 60%

2. **Cache Response Time**
   - p50: < 5ms
   - p95: < 15ms
   - p99: < 30ms

3. **Redis Memory Usage**
   - Current: MB
   - Limit: 100 MB
   - Alert: > 95 MB

4. **Cache Eviction Rate**
   - Evictions/sec
   - Target: < 1/sec
   - Alert: > 10/sec

### CloudWatch/Prometheus Queries

```promql
# Cache hit rate
rate(search_cache_hits_total[5m]) /
(rate(search_cache_hits_total[5m]) + rate(search_cache_misses_total[5m])) * 100

# Cache response time (p95)
histogram_quantile(0.95, rate(search_cache_duration_seconds_bucket[5m]))

# Redis memory usage
redis_memory_used_bytes{service="search-analytics"} / 1024 / 1024
```

## Rollout Plan

### Stage 1: Canary (10% traffic)
- Deploy with caching enabled
- Monitor cache hit rate
- Verify no errors

### Stage 2: Gradual Rollout (50% traffic)
- Increase traffic gradually
- Monitor memory usage
- Tune TTLs if needed

### Stage 3: Full Rollout (100% traffic)
- Enable for all users
- Monitor performance improvements
- Document lessons learned

## Success Criteria

- ✅ Cache hit rate > 80%
- ✅ Cached query response time < 10ms (p95)
- ✅ Redis memory usage < 100 MB
- ✅ No cache-related errors
- ✅ Graceful degradation working

## Conclusion

The Redis caching strategy is designed to:

1. **Maximize Cache Hit Rate:** 80-85% through query normalization and warming
2. **Minimize Latency:** < 10ms for cached queries
3. **Efficient Memory:** < 100 MB with LRU eviction
4. **Reliability:** Graceful degradation if Redis unavailable

**Result:** 10-20x performance improvement for trending/popular queries.
