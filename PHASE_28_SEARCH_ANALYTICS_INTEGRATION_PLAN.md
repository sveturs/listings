# Phase 28: Search Analytics Integration Plan

## Overview

This document describes how to integrate **async search query logging** into the existing `SearchListings` handler without blocking search responses.

## Architecture

```
User → SearchListings gRPC → Search Service
                                  ↓
                           Execute Search (OpenSearch)
                                  ↓
                           Return Results (synchronous)
                                  ↓
                           Log Search Query (async goroutine) → PostgreSQL
```

## Implementation Strategy

### 1. Extract User Context

**From authenticated users:**
- Use `middleware.GetUserID(ctx)` to extract user_id from JWT context
- Available via existing auth interceptor

**From anonymous users:**
- Extract `x-session-id` header from gRPC metadata
- Frontend must send UUID in metadata for session tracking
- If missing → generate UUID in service (for backward compatibility)

### 2. Async Logging Pattern

**Key Requirements:**
- Must NOT block search response
- Must NOT fail search if logging fails
- Use goroutine with separate context (5s timeout)
- Log errors but don't propagate to user

**Code Pattern:**
```go
// In SearchListings handler, AFTER returning results to user
go func() {
    // Create separate context with timeout
    logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Extract user/session from original context
    userID, _ := middleware.GetUserID(ctx)
    sessionID := extractSessionID(ctx)

    // Log search query
    input := &domain.CreateSearchQueryInput{
        QueryText:    req.Query,
        CategoryID:   req.CategoryID,
        UserID:       userID,
        SessionID:    sessionID,
        ResultsCount: int32(len(results.Listings)),
    }

    if _, err := searchQueriesRepo.CreateSearchQuery(logCtx, input); err != nil {
        logger.Warn().Err(err).Msg("failed to log search query (non-blocking)")
    }
}()
```

### 3. Session ID Extraction

Create helper function in `internal/middleware/context.go`:

```go
// SessionIDKey is context key for session ID
type SessionIDKey struct{}

// GetSessionID extracts session ID from context (from x-session-id header)
func GetSessionID(ctx context.Context) (*string, bool) {
    sessionID, ok := ctx.Value(SessionIDKey{}).(string)
    if !ok || sessionID == "" {
        return nil, false
    }
    return &sessionID, true
}

// extractSessionIDFromMetadata extracts session ID from gRPC metadata
func extractSessionIDFromMetadata(ctx context.Context) *string {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil
    }

    sessionIDs := md.Get("x-session-id")
    if len(sessionIDs) == 0 {
        return nil
    }

    sessionID := strings.TrimSpace(sessionIDs[0])
    if sessionID == "" {
        return nil
    }

    return &sessionID
}
```

### 4. Modify Search Service

Update `internal/service/search/service.go`:

**Add repository dependency:**
```go
type Service struct {
    searchClient       *opensearch.SearchClient
    cache              *cache.SearchCache
    searchQueriesRepo  repository.SearchQueriesRepository  // NEW
    logger             zerolog.Logger
}
```

**Modify SearchListings method:**
```go
func (s *Service) SearchListings(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    // ... existing search logic ...

    response := &SearchResponse{
        Listings: listings,
        Total:    searchResp.Hits.Total.Value,
        TookMs:   int32(searchResp.Took),
        Cached:   false,
    }

    // Async logging (non-blocking)
    go s.logSearchQuery(ctx, req, response)

    return response, nil
}

// logSearchQuery logs search query asynchronously
func (s *Service) logSearchQuery(originalCtx context.Context, req *SearchRequest, resp *SearchResponse) {
    // Create separate context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Extract user/session from original context
    var userID *int64
    if uid, ok := middleware.GetUserID(originalCtx); ok {
        userID = &uid
    }

    sessionID := middleware.GetSessionID(originalCtx)

    // Skip logging if no user and no session
    if userID == nil && sessionID == nil {
        s.logger.Debug().Msg("skipping search query logging: no user or session")
        return
    }

    // Create input
    input := &domain.CreateSearchQueryInput{
        QueryText:    req.Query,
        CategoryID:   req.CategoryID,
        UserID:       userID,
        SessionID:    sessionID,
        ResultsCount: int32(len(resp.Listings)),
    }

    // Log to database
    if _, err := s.searchQueriesRepo.CreateSearchQuery(ctx, input); err != nil {
        s.logger.Warn().
            Err(err).
            Str("query", req.Query).
            Msg("failed to log search query (non-blocking)")
    } else {
        s.logger.Debug().
            Str("query", req.Query).
            Int32("results", input.ResultsCount).
            Msg("search query logged")
    }
}
```

### 5. Update GetPopularSearches Implementation

Replace mock implementation in `service.go:625-682`:

```go
func (s *Service) GetPopularSearches(ctx context.Context, req *PopularSearchesRequest) (*PopularSearchesResponse, error) {
    start := time.Now()

    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }

    s.logger.Debug().
        Interface("category_id", req.CategoryID).
        Int32("limit", req.Limit).
        Str("time_range", req.TimeRange).
        Msg("fetching popular searches")

    // Generate cache key (15 min TTL - trending data)
    cacheKey := ""
    if s.cache != nil {
        cacheKey = s.cache.GeneratePopularKey(req.CategoryID, req.TimeRange)

        // Check cache
        if cached, err := s.cache.GetPopular(ctx, cacheKey); err == nil && cached != nil {
            s.logger.Debug().Msg("popular searches cache hit")
            return s.convertCachedPopularSearches(cached), nil
        }
    }

    // Convert time_range to days_ago
    daysAgo := s.timeRangeToDays(req.TimeRange)

    // Query PostgreSQL for trending searches
    filter := &domain.GetTrendingQueriesFilter{
        CategoryID:         req.CategoryID,
        Limit:              req.Limit,
        DaysAgo:            daysAgo,
        MinResultsCount:    1,
        IncludeZeroResults: false,
    }

    trending, err := s.searchQueriesRepo.GetTrendingQueries(ctx, filter)
    if err != nil {
        s.logger.Error().
            Err(err).
            Msg("failed to fetch trending queries")
        return nil, fmt.Errorf("failed to fetch trending queries: %w", err)
    }

    // Convert to PopularSearch format
    searches := s.convertTrendingToPopular(trending)

    response := &PopularSearchesResponse{
        Searches: searches,
        TookMs:   int32(time.Since(start).Milliseconds()),
    }

    // Cache results
    if s.cache != nil && cacheKey != "" {
        go func() {
            cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()

            popularMap := s.convertPopularSearchesForCache(response)
            if err := s.cache.SetPopular(cacheCtx, cacheKey, popularMap); err != nil {
                s.logger.Warn().Err(err).Msg("failed to cache popular searches")
            }
        }()
    }

    s.logger.Info().
        Dur("duration", time.Since(start)).
        Int32("took_ms", response.TookMs).
        Int("searches", len(response.Searches)).
        Msg("popular searches fetched from database")

    return response, nil
}

// timeRangeToDays converts time_range string to days_ago
func (s *Service) timeRangeToDays(timeRange string) int32 {
    switch timeRange {
    case "24h":
        return 1
    case "7d":
        return 7
    case "30d":
        return 30
    default:
        return 7 // default to 7 days
    }
}

// convertTrendingToPopular converts TrendingSearch to PopularSearch
func (s *Service) convertTrendingToPopular(trending []domain.TrendingSearch) []PopularSearch {
    searches := make([]PopularSearch, 0, len(trending))
    for _, t := range trending {
        searches = append(searches, PopularSearch{
            Query:       t.QueryText,
            SearchCount: t.SearchCount,
            TrendScore:  0.0, // TODO: Calculate trend score based on growth rate
        })
    }
    return searches
}
```

## Error Handling Strategy

1. **Logging Failures:** Log warnings but NEVER fail the search
2. **Database Unavailable:** Skip logging, don't block user
3. **Invalid Input:** Log error details for debugging
4. **Context Timeout:** Use separate context with 5s timeout

## Frontend Requirements

Frontend must send `x-session-id` header for anonymous users:

```typescript
// Generate session ID once per browser session
const sessionId = localStorage.getItem('session-id') || uuidv4();
localStorage.setItem('session-id', sessionId);

// Include in gRPC metadata
const metadata = new grpc.Metadata();
metadata.set('x-session-id', sessionId);
```

## Testing Plan

1. **Unit Tests:**
   - Test CreateSearchQuery with valid/invalid input
   - Test GetTrendingQueries with various filters
   - Test async logging doesn't block search

2. **Integration Tests:**
   - Test search query logging for authenticated users
   - Test search query logging for anonymous users
   - Test trending queries aggregation
   - Test CTR tracking

3. **Performance Tests:**
   - Verify async logging adds < 10ms overhead
   - Verify trending queries < 500ms (1M rows)
   - Verify user history < 50ms

## Rollout Plan

1. **Phase 1:** Deploy migrations (000032)
2. **Phase 2:** Deploy microservice with logging enabled
3. **Phase 3:** Monitor logging performance (CloudWatch/Prometheus)
4. **Phase 4:** Enable frontend session-id header
5. **Phase 5:** Replace mock GetPopularSearches with real queries

## Monitoring

**Key Metrics:**
- Search query logging success rate (target: > 99%)
- Async logging duration (target: < 50ms p95)
- Trending queries duration (target: < 500ms p95)
- Database connection pool usage

**Alerts:**
- Search query logging failures > 1%
- Trending queries duration > 1s
- Database connection pool exhaustion

## Rollback Plan

If issues occur:
1. Disable async logging via feature flag
2. Revert to mock GetPopularSearches
3. Rollback migration 000032 if needed

## Success Criteria

- ✅ Migrations applied without errors
- ✅ Search queries logged asynchronously
- ✅ Search response time unchanged (< +10ms)
- ✅ Trending queries return real data
- ✅ User history working for authenticated users
- ✅ No database connection pool exhaustion
