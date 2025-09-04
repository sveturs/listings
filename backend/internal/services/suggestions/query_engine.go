package suggestions

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// QueryEngine provides intelligent query suggestions
type QueryEngine struct {
	redis     *redis.Client
	logger    *zap.Logger
	mu        sync.RWMutex
	patterns  map[string]*QueryPattern
	stopWords map[string]bool
}

// QueryPattern represents a search pattern with metadata
type QueryPattern struct {
	Query       string    `json:"query"`
	Frequency   int64     `json:"frequency"`
	LastUsed    time.Time `json:"last_used"`
	Attributes  []string  `json:"attributes"`
	CategoryID  *int      `json:"category_id,omitempty"`
	Conversions int64     `json:"conversions"`
	CTR         float64   `json:"ctr"` // Click-through rate
}

// Suggestion represents a query suggestion
type Suggestion struct {
	Query      string                 `json:"query"`
	Type       string                 `json:"type"` // "recent", "popular", "trending", "personalized"
	Score      float64                `json:"score"`
	Attributes []string               `json:"attributes,omitempty"`
	Highlight  string                 `json:"highlight,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// NewQueryEngine creates a new query suggestions engine
func NewQueryEngine(redis *redis.Client, logger *zap.Logger) *QueryEngine {
	return &QueryEngine{
		redis:     redis,
		logger:    logger,
		patterns:  make(map[string]*QueryPattern),
		stopWords: initStopWords(),
	}
}

// GetSuggestions returns query suggestions based on input
func (qe *QueryEngine) GetSuggestions(ctx context.Context, input string, options SuggestionOptions) ([]*Suggestion, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if len(input) < 2 {
		return nil, nil
	}

	// Collect suggestions from different sources
	suggestions := make([]*Suggestion, 0)

	// 1. Get prefix matches
	prefixSuggestions := qe.getPrefixSuggestions(ctx, input, options)
	suggestions = append(suggestions, prefixSuggestions...)

	// 2. Get fuzzy matches for typos
	if options.EnableFuzzy {
		fuzzySuggestions := qe.getFuzzySuggestions(ctx, input, options)
		suggestions = append(suggestions, fuzzySuggestions...)
	}

	// 3. Get trending suggestions
	if options.IncludeTrending {
		trendingSuggestions := qe.getTrendingSuggestions(ctx, options.CategoryID)
		suggestions = append(suggestions, trendingSuggestions...)
	}

	// 4. Get personalized suggestions
	if options.UserID != nil {
		personalSuggestions := qe.getPersonalizedSuggestions(ctx, *options.UserID, input)
		suggestions = append(suggestions, personalSuggestions...)
	}

	// Score and rank suggestions
	qe.scoreSuggestions(suggestions, input, options)

	// Sort by score
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	// Limit results
	if len(suggestions) > options.Limit {
		suggestions = suggestions[:options.Limit]
	}

	// Add highlighting
	for _, s := range suggestions {
		s.Highlight = qe.highlightMatch(s.Query, input)
	}

	return suggestions, nil
}

// RecordQuery records a query for learning
func (qe *QueryEngine) RecordQuery(ctx context.Context, query string, metadata QueryMetadata) error {
	query = strings.TrimSpace(strings.ToLower(query))

	// Update pattern in memory
	qe.mu.Lock()
	pattern, exists := qe.patterns[query]
	if !exists {
		pattern = &QueryPattern{
			Query:      query,
			Frequency:  0,
			Attributes: metadata.Attributes,
			CategoryID: metadata.CategoryID,
		}
		qe.patterns[query] = pattern
	}
	pattern.Frequency++
	pattern.LastUsed = time.Now()
	if metadata.Converted {
		pattern.Conversions++
	}
	pattern.CTR = float64(pattern.Conversions) / float64(pattern.Frequency)
	qe.mu.Unlock()

	// Store in Redis for persistence
	key := fmt.Sprintf("query_pattern:%s", query)
	data, _ := json.Marshal(pattern)

	pipe := qe.redis.Pipeline()
	pipe.Set(ctx, key, data, 30*24*time.Hour) // 30 days TTL
	pipe.ZIncrBy(ctx, "query_frequencies", 1, query)

	// Update trending scores
	trendingKey := fmt.Sprintf("trending:%s", time.Now().Format("2006-01-02"))
	pipe.ZIncrBy(ctx, trendingKey, 1, query)
	pipe.Expire(ctx, trendingKey, 7*24*time.Hour)

	// Update user history if user ID provided
	if metadata.UserID != nil {
		userKey := fmt.Sprintf("user_queries:%d", *metadata.UserID)
		pipe.ZAdd(ctx, userKey, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: query,
		})
		pipe.Expire(ctx, userKey, 30*24*time.Hour)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// getPrefixSuggestions returns suggestions matching the prefix
func (qe *QueryEngine) getPrefixSuggestions(ctx context.Context, prefix string, options SuggestionOptions) []*Suggestion {
	suggestions := make([]*Suggestion, 0)

	// Get from Redis sorted set
	results, err := qe.redis.ZRevRangeByLex(ctx, "query_frequencies",
		&redis.ZRangeBy{
			Min: "[" + prefix,
			Max: "[" + prefix + "\xff",
		}).Result()
	if err != nil {
		qe.logger.Error("Failed to get prefix suggestions", zap.Error(err))
		return suggestions
	}

	for _, result := range results {
		// Load pattern details
		pattern := qe.loadPattern(ctx, result)
		if pattern != nil {
			suggestions = append(suggestions, &Suggestion{
				Query:      pattern.Query,
				Type:       "popular",
				Attributes: pattern.Attributes,
				Metadata: map[string]interface{}{
					"frequency":   pattern.Frequency,
					"ctr":         pattern.CTR,
					"conversions": pattern.Conversions,
				},
			})
		}
	}

	return suggestions
}

// getFuzzySuggestions returns fuzzy matches for typo correction
func (qe *QueryEngine) getFuzzySuggestions(ctx context.Context, input string, options SuggestionOptions) []*Suggestion {
	suggestions := make([]*Suggestion, 0)
	maxDistance := 2 // Maximum Levenshtein distance

	qe.mu.RLock()
	defer qe.mu.RUnlock()

	for query, pattern := range qe.patterns {
		distance := levenshteinDistance(input, query)
		if distance > 0 && distance <= maxDistance {
			suggestions = append(suggestions, &Suggestion{
				Query:      pattern.Query,
				Type:       "fuzzy",
				Attributes: pattern.Attributes,
				Metadata: map[string]interface{}{
					"distance":    distance,
					"frequency":   pattern.Frequency,
					"conversions": pattern.Conversions,
				},
			})
		}
	}

	return suggestions
}

// getTrendingSuggestions returns currently trending queries
func (qe *QueryEngine) getTrendingSuggestions(ctx context.Context, categoryID *int) []*Suggestion {
	suggestions := make([]*Suggestion, 0)

	// Get trending queries from last 24 hours
	trendingKey := fmt.Sprintf("trending:%s", time.Now().Format("2006-01-02"))
	results, err := qe.redis.ZRevRangeWithScores(ctx, trendingKey, 0, 9).Result()
	if err != nil {
		qe.logger.Error("Failed to get trending suggestions", zap.Error(err))
		return suggestions
	}

	for _, result := range results {
		query := result.Member.(string)
		pattern := qe.loadPattern(ctx, query)

		if pattern != nil {
			// Filter by category if specified
			if categoryID != nil && pattern.CategoryID != nil && *pattern.CategoryID != *categoryID {
				continue
			}

			suggestions = append(suggestions, &Suggestion{
				Query:      pattern.Query,
				Type:       "trending",
				Attributes: pattern.Attributes,
				Metadata: map[string]interface{}{
					"trend_score": result.Score,
					"frequency":   pattern.Frequency,
				},
			})
		}
	}

	return suggestions
}

// getPersonalizedSuggestions returns personalized suggestions for user
func (qe *QueryEngine) getPersonalizedSuggestions(ctx context.Context, userID int, input string) []*Suggestion {
	suggestions := make([]*Suggestion, 0)

	// Get user's recent queries
	userKey := fmt.Sprintf("user_queries:%d", userID)
	results, err := qe.redis.ZRevRange(ctx, userKey, 0, 19).Result()
	if err != nil {
		qe.logger.Error("Failed to get user queries", zap.Error(err))
		return suggestions
	}

	for _, query := range results {
		// Skip if doesn't match prefix
		if !strings.HasPrefix(query, input) {
			continue
		}

		pattern := qe.loadPattern(ctx, query)
		if pattern != nil {
			suggestions = append(suggestions, &Suggestion{
				Query:      pattern.Query,
				Type:       "personalized",
				Attributes: pattern.Attributes,
				Metadata: map[string]interface{}{
					"user_frequency": pattern.Frequency,
					"last_used":      pattern.LastUsed,
				},
			})
		}
	}

	return suggestions
}

// scoreSuggestions calculates scores for suggestions
func (qe *QueryEngine) scoreSuggestions(suggestions []*Suggestion, input string, options SuggestionOptions) {
	now := time.Now()

	for _, s := range suggestions {
		score := 0.0

		// Base score by type
		switch s.Type {
		case "personalized":
			score = 1.0
		case "trending":
			score = 0.8
		case "popular":
			score = 0.6
		case "fuzzy":
			score = 0.4
		case "recent":
			score = 0.5
		}

		// Boost for exact prefix match
		if strings.HasPrefix(s.Query, input) {
			score += 0.3
		}

		// Boost based on frequency (logarithmic)
		if freq, ok := s.Metadata["frequency"].(int64); ok && freq > 0 {
			score += math.Log10(float64(freq)) * 0.1
		}

		// Boost based on CTR
		if ctr, ok := s.Metadata["ctr"].(float64); ok {
			score += ctr * 0.2
		}

		// Boost for conversions
		if conv, ok := s.Metadata["conversions"].(int64); ok && conv > 0 {
			score += math.Log10(float64(conv)) * 0.15
		}

		// Recency boost for personalized suggestions
		if lastUsed, ok := s.Metadata["last_used"].(time.Time); ok {
			daysSince := now.Sub(lastUsed).Hours() / 24
			if daysSince < 7 {
				score += (7 - daysSince) / 7 * 0.2
			}
		}

		// Penalty for fuzzy matches based on distance
		if distance, ok := s.Metadata["distance"].(int); ok {
			score -= float64(distance) * 0.2
		}

		// Apply length normalization
		lengthDiff := math.Abs(float64(len(s.Query) - len(input)))
		score -= lengthDiff * 0.01

		s.Score = math.Max(0, math.Min(1, score))
	}
}

// highlightMatch highlights the matching part of the query
func (qe *QueryEngine) highlightMatch(query, input string) string {
	lowerQuery := strings.ToLower(query)
	lowerInput := strings.ToLower(input)

	index := strings.Index(lowerQuery, lowerInput)
	if index == -1 {
		return query
	}

	// Return with <mark> tags for highlighting
	return fmt.Sprintf("%s<mark>%s</mark>%s",
		query[:index],
		query[index:index+len(input)],
		query[index+len(input):])
}

// loadPattern loads a query pattern from Redis
func (qe *QueryEngine) loadPattern(ctx context.Context, query string) *QueryPattern {
	key := fmt.Sprintf("query_pattern:%s", query)
	data, err := qe.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil
	}

	var pattern QueryPattern
	if err := json.Unmarshal(data, &pattern); err != nil {
		return nil
	}

	return &pattern
}

// SuggestionOptions contains options for getting suggestions
type SuggestionOptions struct {
	Limit           int
	CategoryID      *int
	UserID          *int
	EnableFuzzy     bool
	IncludeTrending bool
	MinScore        float64
}

// QueryMetadata contains metadata about a query
type QueryMetadata struct {
	UserID     *int
	CategoryID *int
	Attributes []string
	Converted  bool
	Duration   time.Duration
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	if len(s1) == 0 {
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first column and row
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// initStopWords initializes common stop words
func initStopWords() map[string]bool {
	return map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"as": true, "is": true, "was": true, "are": true, "were": true,
	}
}

func min(values ...int) int {
	minVal := values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}
