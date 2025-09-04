package suggestions

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// PersonalizationService provides personalized recommendations
type PersonalizationService struct {
	redis  *redis.Client
	logger *zap.Logger
	engine *RecommendationEngine
}

// RecommendationEngine generates recommendations
type RecommendationEngine struct {
	collaborative *CollaborativeFilter
	contentBased  *ContentBasedFilter
	hybrid        *HybridRecommender
}

// CollaborativeFilter for user-based collaborative filtering
type CollaborativeFilter struct {
	userSimilarity map[int]map[int]float64    // user -> user -> similarity
	itemMatrix     map[int]map[string]float64 // user -> item -> rating
}

// ContentBasedFilter for content-based filtering
type ContentBasedFilter struct {
	itemProfiles map[string]*ItemProfile
	userProfiles map[int]*UserProfile
}

// HybridRecommender combines multiple recommendation strategies
type HybridRecommender struct {
	weights RecommenderWeights
}

// UserProfile represents user preferences
type UserProfile struct {
	UserID           int                `json:"user_id"`
	PreferredQueries []string           `json:"preferred_queries"`
	Categories       map[int]float64    `json:"categories"` // category_id -> affinity
	Attributes       map[string]float64 `json:"attributes"` // attribute -> frequency
	PriceRange       *PriceRange        `json:"price_range"`
	SearchHistory    []SearchRecord     `json:"search_history"`
	ClickHistory     []ClickRecord      `json:"click_history"`
	LastUpdated      time.Time          `json:"last_updated"`
}

// ItemProfile represents item characteristics
type ItemProfile struct {
	Query      string             `json:"query"`
	Category   int                `json:"category"`
	Attributes []string           `json:"attributes"`
	Popularity float64            `json:"popularity"`
	Features   map[string]float64 `json:"features"`
}

// SearchRecord represents a search action
type SearchRecord struct {
	Query     string    `json:"query"`
	Timestamp time.Time `json:"timestamp"`
	Results   int       `json:"results"`
	Clicked   bool      `json:"clicked"`
}

// ClickRecord represents a click action
type ClickRecord struct {
	ItemID    string        `json:"item_id"`
	Query     string        `json:"query"`
	Position  int           `json:"position"`
	Timestamp time.Time     `json:"timestamp"`
	DwellTime time.Duration `json:"dwell_time"`
	Converted bool          `json:"converted"`
}

// RecommenderWeights for hybrid approach
type RecommenderWeights struct {
	Collaborative float64
	ContentBased  float64
	Popularity    float64
	Trending      float64
	Random        float64
}

// NewPersonalizationService creates a new personalization service
func NewPersonalizationService(redis *redis.Client, logger *zap.Logger) *PersonalizationService {
	return &PersonalizationService{
		redis:  redis,
		logger: logger,
		engine: NewRecommendationEngine(),
	}
}

// NewRecommendationEngine creates a new recommendation engine
func NewRecommendationEngine() *RecommendationEngine {
	return &RecommendationEngine{
		collaborative: NewCollaborativeFilter(),
		contentBased:  NewContentBasedFilter(),
		hybrid:        NewHybridRecommender(),
	}
}

// NewCollaborativeFilter creates a new collaborative filter
func NewCollaborativeFilter() *CollaborativeFilter {
	return &CollaborativeFilter{
		userSimilarity: make(map[int]map[int]float64),
		itemMatrix:     make(map[int]map[string]float64),
	}
}

// NewContentBasedFilter creates a new content-based filter
func NewContentBasedFilter() *ContentBasedFilter {
	return &ContentBasedFilter{
		itemProfiles: make(map[string]*ItemProfile),
		userProfiles: make(map[int]*UserProfile),
	}
}

// NewHybridRecommender creates a new hybrid recommender
func NewHybridRecommender() *HybridRecommender {
	return &HybridRecommender{
		weights: RecommenderWeights{
			Collaborative: 0.35,
			ContentBased:  0.30,
			Popularity:    0.20,
			Trending:      0.10,
			Random:        0.05,
		},
	}
}

// GetPersonalizedSuggestions returns personalized suggestions for a user
func (ps *PersonalizationService) GetPersonalizedSuggestions(ctx context.Context, userID int, input string, limit int) ([]*Suggestion, error) {
	// Load user profile
	profile, err := ps.loadUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get recommendations from different strategies
	collaborativeRecs := ps.engine.collaborative.Recommend(userID, input, limit)
	contentRecs := ps.engine.contentBased.Recommend(profile, input, limit)

	// Combine recommendations using hybrid approach
	suggestions := ps.engine.hybrid.Combine(
		collaborativeRecs,
		contentRecs,
		limit,
	)

	// Apply personalization boost based on user profile
	ps.applyPersonalizationBoost(suggestions, profile)

	// Sort by final score
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions, nil
}

// UpdateUserProfile updates user profile based on actions
func (ps *PersonalizationService) UpdateUserProfile(ctx context.Context, userID int, action UserAction) error {
	profile, err := ps.loadUserProfile(ctx, userID)
	if err != nil {
		profile = ps.createNewProfile(userID)
	}

	// Update profile based on action type
	switch action.Type {
	case "search":
		ps.updateSearchHistory(profile, action)
	case "click":
		ps.updateClickHistory(profile, action)
	case "conversion":
		ps.updateConversionData(profile, action)
	}

	// Update category and attribute preferences
	ps.updatePreferences(profile)

	// Save updated profile
	return ps.saveUserProfile(ctx, profile)
}

// loadUserProfile loads user profile from storage
func (ps *PersonalizationService) loadUserProfile(ctx context.Context, userID int) (*UserProfile, error) {
	key := fmt.Sprintf("user_profile:%d", userID)
	data, err := ps.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// saveUserProfile saves user profile to storage
func (ps *PersonalizationService) saveUserProfile(ctx context.Context, profile *UserProfile) error {
	profile.LastUpdated = time.Now()

	key := fmt.Sprintf("user_profile:%d", profile.UserID)
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return ps.redis.Set(ctx, key, data, 30*24*time.Hour).Err()
}

// createNewProfile creates a new user profile
func (ps *PersonalizationService) createNewProfile(userID int) *UserProfile {
	return &UserProfile{
		UserID:           userID,
		PreferredQueries: []string{},
		Categories:       make(map[int]float64),
		Attributes:       make(map[string]float64),
		SearchHistory:    []SearchRecord{},
		ClickHistory:     []ClickRecord{},
		LastUpdated:      time.Now(),
	}
}

// updateSearchHistory updates search history in profile
func (ps *PersonalizationService) updateSearchHistory(profile *UserProfile, action UserAction) {
	record := SearchRecord{
		Query:     action.Query,
		Timestamp: action.Timestamp,
		Results:   action.ResultCount,
		Clicked:   action.Clicked,
	}

	profile.SearchHistory = append(profile.SearchHistory, record)

	// Keep only last 100 searches
	if len(profile.SearchHistory) > 100 {
		profile.SearchHistory = profile.SearchHistory[len(profile.SearchHistory)-100:]
	}

	// Update preferred queries
	if action.Clicked {
		ps.addPreferredQuery(profile, action.Query)
	}
}

// updateClickHistory updates click history in profile
func (ps *PersonalizationService) updateClickHistory(profile *UserProfile, action UserAction) {
	record := ClickRecord{
		ItemID:    action.ItemID,
		Query:     action.Query,
		Position:  action.Position,
		Timestamp: action.Timestamp,
		DwellTime: action.DwellTime,
		Converted: action.Converted,
	}

	profile.ClickHistory = append(profile.ClickHistory, record)

	// Keep only last 200 clicks
	if len(profile.ClickHistory) > 200 {
		profile.ClickHistory = profile.ClickHistory[len(profile.ClickHistory)-200:]
	}
}

// updateConversionData updates conversion-related data
func (ps *PersonalizationService) updateConversionData(profile *UserProfile, action UserAction) {
	// Mark the click as converted
	for i := len(profile.ClickHistory) - 1; i >= 0; i-- {
		if profile.ClickHistory[i].ItemID == action.ItemID {
			profile.ClickHistory[i].Converted = true
			break
		}
	}
}

// updatePreferences updates category and attribute preferences
func (ps *PersonalizationService) updatePreferences(profile *UserProfile) {
	// Reset preferences
	profile.Categories = make(map[int]float64)
	profile.Attributes = make(map[string]float64)

	// Calculate from click history
	for _, click := range profile.ClickHistory {
		// Weight based on recency and engagement
		_ = ps.calculateWeight(click)

		// Update category preference (would need actual category from item)
		// profile.Categories[categoryID] += weight

		// Update attribute preferences (would need actual attributes)
		// for _, attr := range attributes {
		//     profile.Attributes[attr] += weight
		// }
	}

	// Normalize preferences
	ps.normalizePreferences(profile)
}

// calculateWeight calculates weight based on engagement
func (ps *PersonalizationService) calculateWeight(click ClickRecord) float64 {
	weight := 1.0

	// Recency factor (exponential decay)
	daysSince := time.Since(click.Timestamp).Hours() / 24
	weight *= math.Exp(-daysSince / 30) // 30-day half-life

	// Engagement factors
	if click.DwellTime > 30*time.Second {
		weight *= 1.5
	}
	if click.Converted {
		weight *= 2.0
	}

	// Position bias correction
	weight *= 1.0 / (1.0 + float64(click.Position)*0.1)

	return weight
}

// normalizePreferences normalizes preference scores
func (ps *PersonalizationService) normalizePreferences(profile *UserProfile) {
	// Normalize categories
	maxCat := 0.0
	for _, score := range profile.Categories {
		if score > maxCat {
			maxCat = score
		}
	}
	if maxCat > 0 {
		for cat := range profile.Categories {
			profile.Categories[cat] /= maxCat
		}
	}

	// Normalize attributes
	maxAttr := 0.0
	for _, score := range profile.Attributes {
		if score > maxAttr {
			maxAttr = score
		}
	}
	if maxAttr > 0 {
		for attr := range profile.Attributes {
			profile.Attributes[attr] /= maxAttr
		}
	}
}

// addPreferredQuery adds a query to preferred list
func (ps *PersonalizationService) addPreferredQuery(profile *UserProfile, query string) {
	// Check if already exists
	for _, q := range profile.PreferredQueries {
		if q == query {
			return
		}
	}

	profile.PreferredQueries = append(profile.PreferredQueries, query)

	// Keep only top 20 preferred queries
	if len(profile.PreferredQueries) > 20 {
		profile.PreferredQueries = profile.PreferredQueries[len(profile.PreferredQueries)-20:]
	}
}

// applyPersonalizationBoost applies boost based on user profile
func (ps *PersonalizationService) applyPersonalizationBoost(suggestions []*Suggestion, profile *UserProfile) {
	for _, suggestion := range suggestions {
		boost := 1.0

		// Boost if query is in preferred list
		for _, pq := range profile.PreferredQueries {
			if suggestion.Query == pq {
				boost *= 1.5
				break
			}
		}

		// Boost based on category affinity
		if catID, ok := suggestion.Metadata["category_id"].(int); ok {
			if affinity, exists := profile.Categories[catID]; exists {
				boost *= (1.0 + affinity*0.3)
			}
		}

		// Boost based on attribute match
		matchCount := 0
		for _, attr := range suggestion.Attributes {
			if score, exists := profile.Attributes[attr]; exists {
				matchCount++
				boost *= (1.0 + score*0.1)
			}
		}

		suggestion.Score *= boost
	}
}

// Recommend generates collaborative filtering recommendations
func (cf *CollaborativeFilter) Recommend(userID int, input string, limit int) []*Suggestion {
	suggestions := make([]*Suggestion, 0)

	// Find similar users
	similarUsers := cf.findSimilarUsers(userID, 10)

	// Get items liked by similar users
	recommendedItems := make(map[string]float64)
	for _, simUser := range similarUsers {
		if items, exists := cf.itemMatrix[simUser.UserID]; exists {
			for item, rating := range items {
				if strings.Contains(strings.ToLower(item), strings.ToLower(input)) {
					recommendedItems[item] += rating * simUser.Similarity
				}
			}
		}
	}

	// Convert to suggestions
	for query, score := range recommendedItems {
		suggestions = append(suggestions, &Suggestion{
			Query: query,
			Type:  "collaborative",
			Score: score,
			Metadata: map[string]interface{}{
				"recommendation_type": "collaborative",
			},
		})
	}

	// Sort by score
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions
}

// Recommend generates content-based recommendations
func (cb *ContentBasedFilter) Recommend(profile *UserProfile, input string, limit int) []*Suggestion {
	suggestions := make([]*Suggestion, 0)

	// Find items similar to user's preferred items
	for query, itemProfile := range cb.itemProfiles {
		if strings.Contains(strings.ToLower(query), strings.ToLower(input)) {
			similarity := cb.calculateSimilarity(profile, itemProfile)
			if similarity > 0 {
				suggestions = append(suggestions, &Suggestion{
					Query:      query,
					Type:       "content",
					Score:      similarity,
					Attributes: itemProfile.Attributes,
					Metadata: map[string]interface{}{
						"recommendation_type": "content",
						"category":            itemProfile.Category,
					},
				})
			}
		}
	}

	// Sort by similarity
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions
}

// calculateSimilarity calculates similarity between user and item profiles
func (cb *ContentBasedFilter) calculateSimilarity(userProfile *UserProfile, itemProfile *ItemProfile) float64 {
	similarity := 0.0

	// Category similarity
	if affinity, exists := userProfile.Categories[itemProfile.Category]; exists {
		similarity += affinity * 0.4
	}

	// Attribute similarity
	matchCount := 0
	for _, attr := range itemProfile.Attributes {
		if score, exists := userProfile.Attributes[attr]; exists {
			similarity += score * 0.1
			matchCount++
		}
	}

	// Popularity factor
	similarity += itemProfile.Popularity * 0.2

	return similarity
}

// Combine combines recommendations from different strategies
func (hr *HybridRecommender) Combine(collaborative, content []*Suggestion, limit int) []*Suggestion {
	combined := make([]*Suggestion, 0)
	seen := make(map[string]bool)

	// Add collaborative recommendations
	for _, s := range collaborative {
		if !seen[s.Query] {
			s.Score *= hr.weights.Collaborative
			combined = append(combined, s)
			seen[s.Query] = true
		}
	}

	// Add content-based recommendations
	for _, s := range content {
		if !seen[s.Query] {
			s.Score *= hr.weights.ContentBased
			combined = append(combined, s)
			seen[s.Query] = true
		}
	}

	return combined
}

// findSimilarUsers finds similar users for collaborative filtering
func (cf *CollaborativeFilter) findSimilarUsers(userID int, limit int) []SimilarUser {
	similar := make([]SimilarUser, 0)

	userItems := cf.itemMatrix[userID]
	if len(userItems) == 0 {
		return similar
	}

	// Calculate similarity with other users
	for otherID, otherItems := range cf.itemMatrix {
		if otherID == userID {
			continue
		}

		similarity := cf.cosineSimilarity(userItems, otherItems)
		if similarity > 0 {
			similar = append(similar, SimilarUser{
				UserID:     otherID,
				Similarity: similarity,
			})
		}
	}

	// Sort by similarity
	sort.Slice(similar, func(i, j int) bool {
		return similar[i].Similarity > similar[j].Similarity
	})

	if len(similar) > limit {
		similar = similar[:limit]
	}

	return similar
}

// cosineSimilarity calculates cosine similarity between two item vectors
func (cf *CollaborativeFilter) cosineSimilarity(items1, items2 map[string]float64) float64 {
	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for item, rating1 := range items1 {
		norm1 += rating1 * rating1
		if rating2, exists := items2[item]; exists {
			dotProduct += rating1 * rating2
		}
	}

	for _, rating2 := range items2 {
		norm2 += rating2 * rating2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

// UserAction represents a user action
type UserAction struct {
	Type        string // "search", "click", "conversion"
	UserID      int
	Query       string
	ItemID      string
	Position    int
	Timestamp   time.Time
	DwellTime   time.Duration
	ResultCount int
	Clicked     bool
	Converted   bool
	CategoryID  *int
	Attributes  []string
}

// SimilarUser represents a similar user
type SimilarUser struct {
	UserID     int
	Similarity float64
}
