# Phase 28: Search Analytics Performance Analysis

## Overview

This document provides detailed performance analysis for the Search Analytics infrastructure, including EXPLAIN PLAN analysis, index optimization recommendations, and Redis caching strategy.

## Database Schema Summary

```sql
CREATE TABLE search_queries (
    id BIGSERIAL PRIMARY KEY,
    query_text VARCHAR(500) NOT NULL,
    category_id BIGINT NULL,
    user_id BIGINT NULL,
    session_id VARCHAR(255) NULL,
    results_count INTEGER NOT NULL DEFAULT 0,
    clicked_listing_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
1. `idx_search_queries_trending` - (category_id, created_at DESC, query_text) WHERE results_count > 0
2. `idx_search_queries_user_history` - (user_id, created_at DESC) WHERE user_id IS NOT NULL
3. `idx_search_queries_session_history` - (session_id, created_at DESC) WHERE session_id IS NOT NULL
4. `idx_search_queries_created_at` - (created_at DESC)
5. `idx_search_queries_category_agg` - (category_id, query_text) WHERE results_count > 0
6. `idx_search_queries_query_text_fts` - GIN (to_tsvector('english', query_text))
7. `idx_search_queries_ctr_analysis` - (clicked_listing_id, created_at DESC) WHERE clicked_listing_id IS NOT NULL

## EXPLAIN PLAN Analysis

### 1. Trending Queries (Primary Use Case)

**Query:**
```sql
SELECT
    query_text,
    COUNT(*) as search_count,
    MAX(created_at) as last_searched,
    category_id
FROM search_queries
WHERE created_at > NOW() - INTERVAL '7 days'
  AND category_id = 1301
  AND results_count > 0
GROUP BY query_text, category_id
ORDER BY search_count DESC, last_searched DESC
LIMIT 10;
```

**Expected EXPLAIN PLAN (with 1M rows):**
```
QUERY PLAN
--------------------------------------------------------------------------------
Limit  (cost=12543.21..12543.24 rows=10 width=520)
  ->  Sort  (cost=12543.21..12643.45 rows=40098 width=520)
        Sort Key: (count(*)) DESC, (max(created_at)) DESC
        ->  GroupAggregate  (cost=0.42..11234.89 rows=40098 width=520)
              Group Key: query_text, category_id
              ->  Index Scan using idx_search_queries_trending on search_queries
                  (cost=0.42..10234.45 rows=40098 width=520)
                    Index Cond: ((category_id = 1301) AND
                                 (created_at > (now() - '7 days'::interval)))
                    Filter: (results_count > 0)
```

**Performance Characteristics:**
- **Index Used:** `idx_search_queries_trending` (optimal)
- **Scan Type:** Index Scan (no sequential scan)
- **Estimated Cost:** ~12,543 (reasonable for aggregation)
- **Expected Duration:** < 500ms for 1M rows
- **Rows Scanned:** ~40,098 (7 days of data)

**Optimization Notes:**
- Partial index `WHERE results_count > 0` reduces index size by ~10-15%
- Composite index (category_id, created_at, query_text) allows index-only scan
- No need for additional sorting if query_text is in index

### 2. User Search History

**Query:**
```sql
SELECT
    id, query_text, category_id, user_id, session_id,
    results_count, clicked_listing_id, created_at
FROM search_queries
WHERE user_id = 123
ORDER BY created_at DESC
LIMIT 20;
```

**Expected EXPLAIN PLAN:**
```
QUERY PLAN
--------------------------------------------------------------------------------
Limit  (cost=0.42..12.45 rows=20 width=100)
  ->  Index Scan using idx_search_queries_user_history on search_queries
      (cost=0.42..234.56 rows=389 width=100)
        Index Cond: (user_id = 123)
```

**Performance Characteristics:**
- **Index Used:** `idx_search_queries_user_history` (optimal)
- **Scan Type:** Index Scan (covering index)
- **Estimated Cost:** ~12.45 (very low)
- **Expected Duration:** < 50ms
- **Rows Scanned:** Only user's searches (~389 rows per user on average)

**Optimization Notes:**
- Partial index `WHERE user_id IS NOT NULL` saves space for anonymous searches
- Index already sorted by `created_at DESC` - no extra sort needed
- LIMIT 20 ensures fast response even for heavy users

### 3. Session Search History (Anonymous Users)

**Query:**
```sql
SELECT
    id, query_text, category_id, user_id, session_id,
    results_count, clicked_listing_id, created_at
FROM search_queries
WHERE session_id = 'abc123-uuid'
ORDER BY created_at DESC
LIMIT 20;
```

**Expected EXPLAIN PLAN:**
```
QUERY PLAN
--------------------------------------------------------------------------------
Limit  (cost=0.42..8.54 rows=20 width=100)
  ->  Index Scan using idx_search_queries_session_history on search_queries
      (cost=0.42..123.45 rows=304 width=100)
        Index Cond: (session_id = 'abc123-uuid'::text)
```

**Performance Characteristics:**
- **Index Used:** `idx_search_queries_session_history` (optimal)
- **Scan Type:** Index Scan
- **Estimated Cost:** ~8.54 (very low)
- **Expected Duration:** < 50ms
- **Rows Scanned:** Session's searches (~304 rows per session on average)

### 4. All-Time Popular Queries

**Query:**
```sql
SELECT
    query_text,
    COUNT(*) as search_count,
    MAX(created_at) as last_searched,
    category_id
FROM search_queries
WHERE results_count > 0
  AND category_id = 1301
GROUP BY query_text, category_id
ORDER BY search_count DESC
LIMIT 10;
```

**Expected EXPLAIN PLAN:**
```
QUERY PLAN
--------------------------------------------------------------------------------
Limit  (cost=45678.21..45678.24 rows=10 width=520)
  ->  Sort  (cost=45678.21..46234.56 rows=222542 width=520)
        Sort Key: (count(*)) DESC
        ->  GroupAggregate  (cost=0.42..42345.67 rows=222542 width=520)
              Group Key: query_text, category_id
              ->  Index Scan using idx_search_queries_category_agg on search_queries
                  (cost=0.42..38234.12 rows=222542 width=520)
                    Index Cond: (category_id = 1301)
                    Filter: (results_count > 0)
```

**Performance Characteristics:**
- **Index Used:** `idx_search_queries_category_agg` (optimal)
- **Scan Type:** Index Scan
- **Estimated Cost:** ~45,678 (higher due to all-time aggregation)
- **Expected Duration:** < 800ms for 1M rows (ALL TIME)
- **Rows Scanned:** ~222,542 (all successful searches in category)

**Optimization Notes:**
- This query should be **heavily cached** (1 hour TTL)
- Consider materialized view for very large datasets (> 10M rows)
- Partial index `WHERE results_count > 0` critical for performance

### 5. CTR Analysis

**Query:**
```sql
SELECT
    query_text,
    COUNT(*) as total_searches,
    COUNT(clicked_listing_id) as total_clicks,
    ROUND(100.0 * COUNT(clicked_listing_id) / NULLIF(COUNT(*), 0), 2) as ctr_percent,
    category_id
FROM search_queries
WHERE created_at > NOW() - INTERVAL '30 days'
  AND query_text = 'iphone'
GROUP BY query_text, category_id;
```

**Expected EXPLAIN PLAN:**
```
QUERY PLAN
--------------------------------------------------------------------------------
GroupAggregate  (cost=234.56..456.78 rows=1 width=520)
  Group Key: query_text, category_id
  ->  Index Scan using idx_search_queries_trending on search_queries
      (cost=0.42..234.12 rows=234 width=520)
        Index Cond: (created_at > (now() - '30 days'::interval))
        Filter: (query_text = 'iphone'::text)
```

**Performance Characteristics:**
- **Index Used:** `idx_search_queries_trending` (suboptimal but acceptable)
- **Scan Type:** Index Scan with Filter
- **Estimated Cost:** ~456.78
- **Expected Duration:** < 300ms
- **Rows Scanned:** ~234 (searches for "iphone" in 30 days)

**Optimization Notes:**
- Could add separate index on (query_text, created_at) if CTR queries are frequent
- Current indexes are sufficient for admin dashboard use case
- Consider caching CTR results per query (1 hour TTL)

## Index Size Estimates

**Assumptions:**
- 1M total rows
- Average query_text length: 30 characters
- 30% authenticated users, 70% anonymous
- 20% of searches result in clicks

**Index Sizes (approximate):**

| Index Name | Size (MB) | % of Table |
|------------|-----------|------------|
| idx_search_queries_trending | 45 | 35% |
| idx_search_queries_user_history | 15 | 12% |
| idx_search_queries_session_history | 25 | 19% |
| idx_search_queries_created_at | 8 | 6% |
| idx_search_queries_category_agg | 40 | 31% |
| idx_search_queries_query_text_fts | 60 | 46% |
| idx_search_queries_ctr_analysis | 5 | 4% |
| **Total Index Size** | **198 MB** | **152%** |
| **Table Size** | **130 MB** | 100% |

**Notes:**
- Total index size is larger than table (expected for analytics tables)
- Partial indexes reduce size by ~10-15%
- GIN full-text index is largest (60 MB)

## Redis Caching Strategy

### Cache Keys Design

```
# Trending Queries
trending:category:{category_id}:days:{days}  → TTL: 15 minutes
trending:global:days:{days}                  → TTL: 15 minutes

# Popular Queries (All-Time)
popular:category:{category_id}               → TTL: 1 hour
popular:global                               → TTL: 1 hour

# User History
history:user:{user_id}                       → TTL: 5 minutes
history:session:{session_id}                 → TTL: 5 minutes

# CTR Analysis
ctr:query:{query_text}:days:{days}           → TTL: 30 minutes
```

### TTL Rationale

| Cache Type | TTL | Reason |
|------------|-----|--------|
| Trending Queries | 15 min | Rapidly changing, real-time feel |
| Popular Queries | 1 hour | Stable over time, expensive query |
| User History | 5 min | Frequently updated, low cost query |
| CTR Analysis | 30 min | Admin dashboard, not time-critical |

### Cache Hit Rate Expectations

**With proper TTLs:**
- Trending Queries: 85-90% (high traffic, stable results)
- Popular Queries: 95%+ (very stable, long TTL)
- User History: 60-70% (personalized, frequent updates)
- CTR Analysis: 80-85% (admin use, moderate traffic)

**Overall Cache Hit Rate:** ~80-85%

### Memory Usage Estimates

**Assumptions:**
- 100 active categories
- 10,000 active users
- 50,000 active sessions
- Average cached result size: 2 KB

**Redis Memory Usage:**

| Cache Type | Keys | Size per Key | Total Size |
|------------|------|--------------|------------|
| Trending | 200 | 2 KB | 400 KB |
| Popular | 100 | 5 KB | 500 KB |
| User History | 10,000 | 1 KB | 10 MB |
| Session History | 50,000 | 1 KB | 50 MB |
| CTR Analysis | 500 | 2 KB | 1 MB |
| **Total** | **60,800** | - | **~62 MB** |

**Conclusion:** Redis memory usage is negligible (< 100 MB)

## Performance Benchmarks (Expected)

### Query Performance (1M rows, cold cache)

| Query Type | Target | Expected | Worst Case |
|------------|--------|----------|------------|
| Trending (7 days, category) | < 500ms | 350ms | 600ms |
| Trending (7 days, global) | < 800ms | 550ms | 1000ms |
| User History (20 results) | < 50ms | 25ms | 80ms |
| Session History (20 results) | < 50ms | 30ms | 80ms |
| Popular (all-time, category) | < 800ms | 600ms | 1200ms |
| CTR Analysis (30 days) | < 300ms | 200ms | 500ms |

### Query Performance (with Redis cache)

| Query Type | Target | Expected | p95 |
|------------|--------|----------|-----|
| Trending (cached) | < 10ms | 5ms | 15ms |
| Popular (cached) | < 10ms | 5ms | 15ms |
| User History (cached) | < 10ms | 5ms | 15ms |

### Async Logging Performance

| Metric | Target | Expected |
|--------|--------|----------|
| Logging duration | < 50ms | 30ms |
| Search overhead | < 10ms | 5ms |
| Success rate | > 99% | 99.5% |

## Scaling Considerations

### When to Partition?

**Consider table partitioning when:**
- Total rows > 10M
- Query performance degrades > 20%
- Index size > 1 GB

**Partition Strategy:**
```sql
-- Partition by created_at (monthly)
CREATE TABLE search_queries_2025_11 PARTITION OF search_queries
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

**Benefits:**
- Faster queries (scan only relevant partitions)
- Easier cleanup (DROP partition instead of DELETE)
- Better index maintenance

### When to Add Read Replicas?

**Consider read replicas when:**
- Trending queries > 1000 RPS
- Database CPU > 70%
- Primary instance impacted by analytics queries

**Strategy:**
- Route trending/popular queries to replica
- Keep write operations (CreateSearchQuery) on primary
- Use connection pooling with read/write splitting

## Optimization Recommendations

### Immediate (Phase 28)

1. ✅ Create all 7 indexes from migration 000032
2. ✅ Implement Redis caching with TTLs
3. ✅ Use async logging pattern (non-blocking)
4. ✅ Monitor query performance (CloudWatch)

### Short-term (1-2 months)

1. Add materialized view for popular queries if > 5M rows
2. Implement query result caching in application layer
3. Add database query monitoring (pg_stat_statements)
4. Set up alerts for slow queries (> 1s)

### Long-term (6+ months)

1. Consider table partitioning if > 10M rows
2. Evaluate read replicas for analytics queries
3. Implement query result pre-warming (cache refresh job)
4. Add database connection pooling optimization

## Monitoring Queries

### Index Usage Statistics

```sql
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE tablename = 'search_queries'
ORDER BY idx_scan DESC;
```

### Slow Queries Detection

```sql
SELECT
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
WHERE query LIKE '%search_queries%'
ORDER BY mean_time DESC
LIMIT 10;
```

### Table Size Monitoring

```sql
SELECT
    pg_size_pretty(pg_total_relation_size('search_queries')) as total_size,
    pg_size_pretty(pg_relation_size('search_queries')) as table_size,
    pg_size_pretty(pg_indexes_size('search_queries')) as indexes_size;
```

## Conclusion

The Search Analytics infrastructure is designed for:

- **High Performance:** < 500ms for trending queries (1M rows)
- **Scalability:** Handles 10M+ rows with partitioning
- **Reliability:** 99.5%+ logging success rate
- **Efficiency:** 80-85% cache hit rate
- **Cost-Effective:** < 100 MB Redis memory

**Performance targets are realistic and achievable with proper indexing and caching strategy.**
