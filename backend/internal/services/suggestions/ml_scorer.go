package suggestions

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
)

// MLScorer provides machine learning based scoring for suggestions
type MLScorer struct {
	logger       *zap.Logger
	features     *FeatureExtractor
	weights      *ModelWeights
	clickHistory *ClickHistory
}

// FeatureExtractor extracts features from queries and context
type FeatureExtractor struct {
	tokenizer *Tokenizer
}

// ModelWeights contains learned weights for scoring
type ModelWeights struct {
	// TF-IDF weights
	TermFrequency  float64
	InverseDocFreq float64

	// Behavioral weights
	ClickThroughRate float64
	ConversionRate   float64
	DwellTime        float64

	// Contextual weights
	TimeOfDay   float64
	DayOfWeek   float64
	Seasonality float64

	// Personalization weights
	UserHistory      float64
	CategoryAffinity float64
	PriceRange       float64

	// Query characteristics
	QueryLength     float64
	QueryComplexity float64
	HasAttributes   float64
}

// ClickHistory tracks user interactions
type ClickHistory struct {
	Clicks      map[string]int64
	Impressions map[string]int64
	Conversions map[string]int64
	DwellTimes  map[string]time.Duration
}

// Tokenizer for text processing
type Tokenizer struct {
	stopWords map[string]bool
}

// NewMLScorer creates a new ML scorer
func NewMLScorer(logger *zap.Logger) *MLScorer {
	return &MLScorer{
		logger:       logger,
		features:     NewFeatureExtractor(),
		weights:      DefaultModelWeights(),
		clickHistory: NewClickHistory(),
	}
}

// NewFeatureExtractor creates a new feature extractor
func NewFeatureExtractor() *FeatureExtractor {
	return &FeatureExtractor{
		tokenizer: NewTokenizer(),
	}
}

// NewTokenizer creates a new tokenizer
func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		stopWords: initStopWords(),
	}
}

// NewClickHistory creates a new click history tracker
func NewClickHistory() *ClickHistory {
	return &ClickHistory{
		Clicks:      make(map[string]int64),
		Impressions: make(map[string]int64),
		Conversions: make(map[string]int64),
		DwellTimes:  make(map[string]time.Duration),
	}
}

// DefaultModelWeights returns default model weights
func DefaultModelWeights() *ModelWeights {
	return &ModelWeights{
		// TF-IDF weights
		TermFrequency:  0.25,
		InverseDocFreq: 0.20,

		// Behavioral weights (highest importance)
		ClickThroughRate: 0.35,
		ConversionRate:   0.40,
		DwellTime:        0.15,

		// Contextual weights
		TimeOfDay:   0.05,
		DayOfWeek:   0.05,
		Seasonality: 0.10,

		// Personalization weights
		UserHistory:      0.30,
		CategoryAffinity: 0.20,
		PriceRange:       0.10,

		// Query characteristics
		QueryLength:     0.05,
		QueryComplexity: 0.10,
		HasAttributes:   0.15,
	}
}

// ScoreSuggestion calculates ML score for a suggestion
func (s *MLScorer) ScoreSuggestion(ctx context.Context, suggestion *Suggestion, context ScoringContext) float64 {
	features := s.extractFeatures(suggestion, context)
	score := s.calculateScore(features)

	// Apply temporal boost
	score *= s.getTemporalBoost(context.Timestamp)

	// Apply personalization boost
	if context.UserID != nil {
		score *= s.getPersonalizationBoost(*context.UserID, suggestion)
	}

	// Normalize score to [0, 1]
	return math.Max(0, math.Min(1, score))
}

// extractFeatures extracts features from suggestion and context
func (s *MLScorer) extractFeatures(suggestion *Suggestion, context ScoringContext) map[string]float64 {
	features := make(map[string]float64)

	// Text features
	features["tf_idf"] = s.calculateTFIDF(suggestion.Query, context.Corpus)
	features["query_length"] = float64(len(strings.Fields(suggestion.Query)))
	features["has_attributes"] = boolToFloat(len(suggestion.Attributes) > 0)

	// Behavioral features
	ctr := s.clickHistory.GetCTR(suggestion.Query)
	features["ctr"] = ctr
	features["conversion_rate"] = s.clickHistory.GetConversionRate(suggestion.Query)
	features["avg_dwell_time"] = s.clickHistory.GetAvgDwellTime(suggestion.Query).Seconds()

	// Temporal features
	now := context.Timestamp
	features["hour_of_day"] = float64(now.Hour()) / 24.0
	features["day_of_week"] = float64(now.Weekday()) / 7.0
	features["day_of_month"] = float64(now.Day()) / 31.0
	features["month"] = float64(now.Month()) / 12.0

	// Query complexity
	features["complexity"] = s.calculateQueryComplexity(suggestion.Query)

	// Category affinity
	if context.CategoryID != nil && suggestion.Metadata != nil {
		if catID, ok := suggestion.Metadata["category_id"].(int); ok {
			features["category_match"] = boolToFloat(catID == *context.CategoryID)
		}
	}

	// Price range affinity
	if context.PriceRange != nil && suggestion.Metadata != nil {
		if price, ok := suggestion.Metadata["avg_price"].(float64); ok {
			if price >= context.PriceRange.Min && price <= context.PriceRange.Max {
				features["price_in_range"] = 1.0
			}
		}
	}

	return features
}

// calculateScore applies model weights to features
func (s *MLScorer) calculateScore(features map[string]float64) float64 {
	score := 0.0

	// TF-IDF component
	score += features["tf_idf"] * s.weights.TermFrequency

	// Behavioral component
	score += features["ctr"] * s.weights.ClickThroughRate
	score += features["conversion_rate"] * s.weights.ConversionRate
	score += math.Min(features["avg_dwell_time"]/60, 1.0) * s.weights.DwellTime // Normalize to minutes

	// Contextual component
	score += features["hour_of_day"] * s.weights.TimeOfDay
	score += features["day_of_week"] * s.weights.DayOfWeek

	// Query characteristics
	queryLengthNorm := math.Min(features["query_length"]/10, 1.0) // Normalize to 10 words max
	score += queryLengthNorm * s.weights.QueryLength
	score += features["complexity"] * s.weights.QueryComplexity
	score += features["has_attributes"] * s.weights.HasAttributes

	// Category and price affinity
	score += features["category_match"] * s.weights.CategoryAffinity
	score += features["price_in_range"] * s.weights.PriceRange

	return score
}

// calculateTFIDF calculates TF-IDF score
func (s *MLScorer) calculateTFIDF(query string, corpus []string) float64 {
	tokens := s.features.tokenizer.Tokenize(query)
	if len(tokens) == 0 {
		return 0
	}

	// Calculate term frequency
	tf := make(map[string]float64)
	for _, token := range tokens {
		tf[token]++
	}
	for token := range tf {
		tf[token] /= float64(len(tokens))
	}

	// Calculate inverse document frequency
	idf := make(map[string]float64)
	corpusSize := float64(len(corpus))

	for token := range tf {
		docCount := 0.0
		for _, doc := range corpus {
			if strings.Contains(strings.ToLower(doc), token) {
				docCount++
			}
		}
		if docCount > 0 {
			idf[token] = math.Log(corpusSize / docCount)
		}
	}

	// Calculate TF-IDF score
	score := 0.0
	for token := range tf {
		score += tf[token] * idf[token]
	}

	return score / float64(len(tokens)) // Normalize by token count
}

// calculateQueryComplexity estimates query complexity
func (s *MLScorer) calculateQueryComplexity(query string) float64 {
	tokens := s.features.tokenizer.Tokenize(query)

	// Factors for complexity
	numTokens := float64(len(tokens))
	uniqueTokens := float64(len(uniqueTokenSet(tokens)))
	avgTokenLength := averageTokenLength(tokens)
	hasSpecialChars := containsSpecialChars(query)

	// Calculate complexity score
	complexity := 0.0
	complexity += math.Min(numTokens/10, 1.0) * 0.3      // Number of tokens
	complexity += (uniqueTokens / numTokens) * 0.3       // Token diversity
	complexity += math.Min(avgTokenLength/10, 1.0) * 0.2 // Average token length
	complexity += boolToFloat(hasSpecialChars) * 0.2     // Special characters

	return complexity
}

// getTemporalBoost returns temporal boost factor
func (s *MLScorer) getTemporalBoost(timestamp time.Time) float64 {
	hour := timestamp.Hour()

	// Peak hours boost (9-11 AM, 7-9 PM)
	if (hour >= 9 && hour <= 11) || (hour >= 19 && hour <= 21) {
		return 1.2
	}

	// Normal hours
	if hour >= 8 && hour <= 22 {
		return 1.0
	}

	// Off-peak hours
	return 0.8
}

// getPersonalizationBoost returns personalization boost
func (s *MLScorer) getPersonalizationBoost(userID int, suggestion *Suggestion) float64 {
	// This would typically query user preferences from database
	// For now, return a simple boost based on suggestion type
	if suggestion.Type == "personalized" {
		return 1.3
	}
	return 1.0
}

// Tokenize splits query into tokens
func (t *Tokenizer) Tokenize(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	tokens := make([]string, 0, len(words))

	for _, word := range words {
		// Remove punctuation
		word = strings.Trim(word, ".,!?;:'\"")

		// Skip stop words
		if !t.stopWords[word] && len(word) > 0 {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

// GetCTR returns click-through rate for a query
func (ch *ClickHistory) GetCTR(query string) float64 {
	clicks := ch.Clicks[query]
	impressions := ch.Impressions[query]

	if impressions == 0 {
		return 0
	}

	return float64(clicks) / float64(impressions)
}

// GetConversionRate returns conversion rate for a query
func (ch *ClickHistory) GetConversionRate(query string) float64 {
	conversions := ch.Conversions[query]
	clicks := ch.Clicks[query]

	if clicks == 0 {
		return 0
	}

	return float64(conversions) / float64(clicks)
}

// GetAvgDwellTime returns average dwell time for a query
func (ch *ClickHistory) GetAvgDwellTime(query string) time.Duration {
	totalTime := ch.DwellTimes[query]
	clicks := ch.Clicks[query]

	if clicks == 0 {
		return 0
	}

	return totalTime / time.Duration(clicks)
}

// RecordImpression records an impression
func (ch *ClickHistory) RecordImpression(query string) {
	ch.Impressions[query]++
}

// RecordClick records a click
func (ch *ClickHistory) RecordClick(query string, dwellTime time.Duration) {
	ch.Clicks[query]++
	ch.DwellTimes[query] += dwellTime
}

// RecordConversion records a conversion
func (ch *ClickHistory) RecordConversion(query string) {
	ch.Conversions[query]++
}

// ScoringContext contains context for scoring
type ScoringContext struct {
	UserID     *int
	CategoryID *int
	PriceRange *PriceRange
	Timestamp  time.Time
	Corpus     []string // Document corpus for TF-IDF
	Location   *Location
}

// PriceRange represents a price range
type PriceRange struct {
	Min float64
	Max float64
}

// Location represents geographic location
type Location struct {
	Lat float64
	Lng float64
}

// RankSuggestions ranks suggestions using ML scoring
func (s *MLScorer) RankSuggestions(ctx context.Context, suggestions []*Suggestion, scoringContext ScoringContext) []*Suggestion {
	// Score each suggestion
	for _, suggestion := range suggestions {
		suggestion.Score = s.ScoreSuggestion(ctx, suggestion, scoringContext)
	}

	// Sort by score
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	return suggestions
}

// Helper functions

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

func uniqueTokenSet(tokens []string) map[string]bool {
	set := make(map[string]bool)
	for _, token := range tokens {
		set[token] = true
	}
	return set
}

func averageTokenLength(tokens []string) float64 {
	if len(tokens) == 0 {
		return 0
	}

	totalLength := 0
	for _, token := range tokens {
		totalLength += len(token)
	}

	return float64(totalLength) / float64(len(tokens))
}

func containsSpecialChars(text string) bool {
	for _, char := range text {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == ' ') {
			return true
		}
	}
	return false
}
