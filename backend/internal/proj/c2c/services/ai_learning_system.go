package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"backend/internal/proj/c2c/repository"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// AILearningSystem implements self-improving AI categorization
type AILearningSystem struct {
	logger       *zap.Logger
	redisClient  *redis.Client
	keywordRepo  *repository.KeywordRepository
	validator    *AICategoryValidator
	generator    *AIKeywordGenerator
	feedbackRepo FeedbackRepository // Interface for feedback storage
}

// FeedbackRepository interface for dependency injection
type FeedbackRepository interface {
	GetIncorrectDetections(ctx context.Context, limit int) ([]FeedbackRecord, error)
	GetRecentFeedbacks(ctx context.Context, hours int) ([]FeedbackRecord, error)
	LogIncorrectDetection(ctx context.Context, input AIDetectionInput, suggested ValidationResult) error
	UpdateKeywordSuccessRate(ctx context.Context, keyword string, categoryID int32, successful bool) error
}

// FeedbackRecord represents user feedback on category detection
type FeedbackRecord struct {
	ID                  int64                  `json:"id"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	DetectedCategoryID  int32                  `json:"detectedCategoryId"`
	SuggestedCategoryID int32                  `json:"suggestedCategoryId"`
	SuggestedCategory   string                 `json:"suggestedCategory"`
	SuggestedKeywords   []string               `json:"suggestedKeywords"`
	AIHints             map[string]interface{} `json:"aiHints"`
	Confidence          float64                `json:"confidence"`
	UserConfirmed       bool                   `json:"userConfirmed"`
	CreatedAt           time.Time              `json:"createdAt"`
}

// LearningMetrics represents learning system performance
type LearningMetrics struct {
	TotalFeedbacks      int      `json:"totalFeedbacks"`
	IncorrectDetections int      `json:"incorrectDetections"`
	ImprovementsApplied int      `json:"improvementsApplied"`
	KeywordsLearned     int      `json:"keywordsLearned"`
	AccuracyImprovement float64  `json:"accuracyImprovement"`
	LastLearningSession string   `json:"lastLearningSession"`
	AvgProcessingTimeMs int64    `json:"avgProcessingTimeMs"`
	RecommendedActions  []string `json:"recommendedActions"`
}

// NewAILearningSystem creates a new self-learning AI system
func NewAILearningSystem(
	logger *zap.Logger,
	redisClient *redis.Client,
	keywordRepo *repository.KeywordRepository,
	validator *AICategoryValidator,
	generator *AIKeywordGenerator,
	feedbackRepo FeedbackRepository,
) *AILearningSystem {
	return &AILearningSystem{
		logger:       logger,
		redisClient:  redisClient,
		keywordRepo:  keywordRepo,
		validator:    validator,
		generator:    generator,
		feedbackRepo: feedbackRepo,
	}
}

// LearnFromValidationFeedback processes AI validation results to improve the system
func (als *AILearningSystem) LearnFromValidationFeedback(ctx context.Context) (*LearningMetrics, error) {
	startTime := time.Now()
	metrics := &LearningMetrics{
		LastLearningSession: startTime.Format(time.RFC3339),
	}

	als.logger.Info("Starting learning session from validation feedback")

	// Get recent incorrect detections (last 24 hours)
	incorrectCases, err := als.feedbackRepo.GetIncorrectDetections(ctx, 100)
	if err != nil {
		als.logger.Error("Failed to get incorrect detections", zap.Error(err))
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	metrics.IncorrectDetections = len(incorrectCases)
	metrics.TotalFeedbacks = len(incorrectCases)

	if len(incorrectCases) == 0 {
		als.logger.Info("No incorrect detections found for learning")
		return metrics, nil
	}

	// Analyze patterns and learn
	improvementsApplied := 0
	keywordsLearned := 0

	for _, feedback := range incorrectCases {
		improvements, keywords := als.processFeedback(ctx, feedback)
		improvementsApplied += improvements
		keywordsLearned += keywords
	}

	metrics.ImprovementsApplied = improvementsApplied
	metrics.KeywordsLearned = keywordsLearned

	// Update learning statistics
	als.updateLearningStats(ctx, metrics)

	// Generate recommendations
	metrics.RecommendedActions = als.generateRecommendations(ctx, incorrectCases)

	metrics.AvgProcessingTimeMs = time.Since(startTime).Milliseconds()

	als.logger.Info("Learning session completed",
		zap.Int("totalFeedbacks", metrics.TotalFeedbacks),
		zap.Int("improvements", metrics.ImprovementsApplied),
		zap.Int("keywordsLearned", metrics.KeywordsLearned),
		zap.Int64("processingTime", metrics.AvgProcessingTimeMs))

	return metrics, nil
}

// processFeedback processes a single feedback item
func (als *AILearningSystem) processFeedback(ctx context.Context, feedback FeedbackRecord) (int, int) {
	improvements := 0
	keywordsLearned := 0

	// 1. Learn new keywords from AI suggestions
	if len(feedback.SuggestedKeywords) > 0 && feedback.SuggestedCategoryID > 0 {
		keywords := make([]repository.GeneratedKeyword, 0)
		for _, keyword := range feedback.SuggestedKeywords {
			keywords = append(keywords, repository.GeneratedKeyword{
				Keyword:     strings.ToLower(strings.TrimSpace(keyword)),
				Type:        "context", // AI-suggested keywords are context type
				Weight:      1.2,       // Medium-high weight for AI suggestions
				Confidence:  0.8,       // High confidence in AI suggestions
				Description: "Learned from AI validation feedback",
			})
		}

		err := als.keywordRepo.BulkInsertKeywords(ctx, feedback.SuggestedCategoryID, keywords, "ai_learned")
		if err != nil {
			als.logger.Error("Failed to save learned keywords",
				zap.Int32("categoryId", feedback.SuggestedCategoryID),
				zap.Error(err))
		} else {
			keywordsLearned += len(keywords)
			improvements++
		}
	}

	// 2. Extract patterns from title and description
	if feedback.SuggestedCategoryID > 0 {
		extractedKeywords := als.extractKeywordsFromText(feedback.Title, feedback.Description)
		if len(extractedKeywords) > 0 {
			keywords := make([]repository.GeneratedKeyword, 0)
			for _, keyword := range extractedKeywords {
				keywords = append(keywords, repository.GeneratedKeyword{
					Keyword:     keyword,
					Type:        "pattern",
					Weight:      0.8, // Lower weight for extracted patterns
					Confidence:  0.6,
					Description: "Extracted from feedback analysis",
				})
			}

			err := als.keywordRepo.BulkInsertKeywords(ctx, feedback.SuggestedCategoryID, keywords, "pattern_extracted")
			if err == nil {
				keywordsLearned += len(keywords)
				improvements++
			}
		}
	}

	// 3. Update weight of keywords that failed
	if feedback.DetectedCategoryID > 0 {
		// Get keywords that matched but led to wrong category
		existingKeywords, err := als.keywordRepo.GetKeywordsByCategory(ctx, feedback.DetectedCategoryID)
		if err == nil {
			titleWords := strings.Fields(strings.ToLower(feedback.Title))
			for _, kw := range existingKeywords {
				for _, word := range titleWords {
					if strings.Contains(word, kw.Keyword) || strings.Contains(kw.Keyword, word) {
						// Reduce weight of keyword that led to wrong classification
						newWeight := kw.Weight * 0.9 // Reduce by 10%
						if newWeight < 0.1 {
							newWeight = 0.1
						}
						if err := als.keywordRepo.UpdateKeywordWeight(ctx, kw.ID, newWeight); err != nil {
							als.logger.Warn("Failed to update keyword weight",
								zap.Int32("keywordId", kw.ID),
								zap.Float64("newWeight", newWeight),
								zap.Error(err))
						}
						improvements++
						break
					}
				}
			}
		}
	}

	return improvements, keywordsLearned
}

// extractKeywordsFromText extracts meaningful keywords from text
func (als *AILearningSystem) extractKeywordsFromText(title, description string) []string {
	text := strings.ToLower(title + " " + description)
	words := strings.Fields(text)

	// Filter out common words and extract meaningful terms
	meaningfulWords := make(map[string]bool)
	stopWords := map[string]bool{
		"и": true, "в": true, "на": true, "с": true, "для": true, "из": true,
		"по": true, "от": true, "до": true, "за": true, "к": true, "о": true,
		"the": true, "and": true, "or": true, "but": true, "in": true, "on": true,
		"at": true, "to": true, "for": true, "of": true, "with": true, "by": true,
	}

	for _, word := range words {
		// Clean word
		word = strings.Trim(word, ".,!?()[]{}\"'")

		// Filter criteria
		if len(word) < 3 || len(word) > 30 {
			continue
		}
		if stopWords[word] {
			continue
		}
		if strings.Contains(word, "http") || strings.Contains(word, "@") {
			continue
		}

		meaningfulWords[word] = true
	}

	// Convert to slice
	keywords := make([]string, 0, len(meaningfulWords))
	for word := range meaningfulWords {
		keywords = append(keywords, word)
	}

	// Limit to most relevant
	if len(keywords) > 10 {
		keywords = keywords[:10]
	}

	return keywords
}

// AutoImproveKeywords automatically improves keyword coverage for poor-performing categories
func (als *AILearningSystem) AutoImproveKeywords(ctx context.Context) error {
	als.logger.Info("Starting automatic keyword improvement")

	// Find categories with poor keyword coverage (less than 30 keywords)
	categories, err := als.keywordRepo.GetCategoriesNeedingKeywords(ctx, 30)
	if err != nil {
		return fmt.Errorf("failed to get categories needing improvement: %w", err)
	}

	if len(categories) == 0 {
		als.logger.Info("All categories have sufficient keyword coverage")
		return nil
	}

	als.logger.Info("Found categories needing keyword improvement",
		zap.Int("categoriesCount", len(categories)))

	// Process up to 5 categories per session to avoid overwhelming
	maxCategories := 5
	if len(categories) > maxCategories {
		categories = categories[:maxCategories]
	}

	for _, category := range categories {
		req := KeywordGenerationRequest{
			CategoryID:   category.ID,
			CategoryName: category.Name,
			Language:     "ru",
			MinKeywords:  50,
		}

		result, err := als.generator.GenerateKeywordsForCategory(ctx, req)
		if err != nil {
			als.logger.Error("Failed to auto-generate keywords",
				zap.Int32("categoryId", category.ID),
				zap.String("categoryName", category.Name),
				zap.Error(err))
			continue
		}

		// Save generated keywords
		if len(result.Keywords) > 0 {
			// Convert GeneratedKeyword to repository.GeneratedKeyword
			repoKeywords := make([]repository.GeneratedKeyword, len(result.Keywords))
			for i, kw := range result.Keywords {
				repoKeywords[i] = repository.GeneratedKeyword{
					Keyword:     kw.Keyword,
					Type:        kw.Type,
					Weight:      kw.Weight,
					Confidence:  kw.Confidence,
					Description: kw.Description,
				}
			}
			err = als.keywordRepo.BulkInsertKeywords(ctx, category.ID, repoKeywords, "auto_improved")
			if err != nil {
				als.logger.Error("Failed to save auto-generated keywords",
					zap.Int32("categoryId", category.ID),
					zap.Error(err))
			} else {
				als.logger.Info("Auto-improved keywords for category",
					zap.Int32("categoryId", category.ID),
					zap.String("categoryName", category.Name),
					zap.Int("keywordsAdded", len(result.Keywords)))
			}
		}

		// Small delay between categories
		time.Sleep(2 * time.Second)
	}

	return nil
}

// generateRecommendations generates actionable recommendations based on feedback
func (als *AILearningSystem) generateRecommendations(ctx context.Context, feedbacks []FeedbackRecord) []string {
	recommendations := []string{}

	if len(feedbacks) == 0 {
		return []string{"System is performing well, no immediate improvements needed"}
	}

	// Analyze common failure patterns
	categoryErrors := make(map[int32]int)
	domainErrors := make(map[string]int)

	for _, feedback := range feedbacks {
		categoryErrors[feedback.DetectedCategoryID]++

		if aiHints, ok := feedback.AIHints["domain"].(string); ok {
			domainErrors[aiHints]++
		}
	}

	// Generate recommendations based on patterns
	if len(categoryErrors) > 0 {
		// Find most problematic category
		maxErrors := 0
		problematicCategory := int32(0)
		for categoryID, errors := range categoryErrors {
			if errors > maxErrors {
				maxErrors = errors
				problematicCategory = categoryID
			}
		}

		if maxErrors > 3 {
			recommendations = append(recommendations,
				fmt.Sprintf("Category ID %d has %d incorrect detections - review and expand keywords",
					problematicCategory, maxErrors))
		}
	}

	// Check if bulk keyword generation is needed
	if len(feedbacks) > 10 {
		recommendations = append(recommendations,
			"High number of incorrect detections detected - consider running bulk keyword generation")
	}

	// Performance recommendations
	if len(feedbacks) > 5 {
		recommendations = append(recommendations,
			"System accuracy below optimal - run comprehensive test to identify specific issues")
	}

	// Domain-specific recommendations
	for domain, errors := range domainErrors {
		if errors > 2 {
			recommendations = append(recommendations,
				fmt.Sprintf("Domain '%s' has %d detection errors - review AI mapping rules", domain, errors))
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "System is learning well from feedback")
	}

	return recommendations
}

// updateLearningStats updates Redis cache with learning statistics
func (als *AILearningSystem) updateLearningStats(ctx context.Context, metrics *LearningMetrics) {
	statsJSON, err := json.Marshal(metrics)
	if err != nil {
		als.logger.Error("Failed to marshal learning stats", zap.Error(err))
		return
	}

	key := "ai_learning_stats:latest"
	err = als.redisClient.Set(ctx, key, statsJSON, 24*time.Hour).Err()
	if err != nil {
		als.logger.Error("Failed to cache learning stats", zap.Error(err))
	}

	// Also store historical data
	historyKey := fmt.Sprintf("ai_learning_stats:%s", time.Now().Format("2006-01-02"))
	als.redisClient.Set(ctx, historyKey, statsJSON, 7*24*time.Hour) // Keep for 7 days
}

// GetLearningStats retrieves current learning statistics
func (als *AILearningSystem) GetLearningStats(ctx context.Context) (*LearningMetrics, error) {
	key := "ai_learning_stats:latest"
	data, err := als.redisClient.Get(ctx, key).Result()
	if err != nil {
		// Return empty metrics if no cache
		return &LearningMetrics{
			LastLearningSession: "Never",
			RecommendedActions:  []string{"Run initial learning session"},
		}, nil
	}

	var metrics LearningMetrics
	err = json.Unmarshal([]byte(data), &metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to parse learning stats: %w", err)
	}

	return &metrics, nil
}

// ScheduledLearning performs all scheduled learning tasks
func (als *AILearningSystem) ScheduledLearning(ctx context.Context) error {
	als.logger.Info("Starting scheduled learning session")

	// 1. Learn from validation feedback
	metrics, err := als.LearnFromValidationFeedback(ctx)
	if err != nil {
		als.logger.Error("Failed to learn from validation feedback", zap.Error(err))
		return err
	}

	// 2. Auto-improve keywords for poor-performing categories
	err = als.AutoImproveKeywords(ctx)
	if err != nil {
		als.logger.Error("Failed to auto-improve keywords", zap.Error(err))
		// Don't return error, continue with other tasks
	}

	als.logger.Info("Scheduled learning session completed",
		zap.Int("improvements", metrics.ImprovementsApplied),
		zap.Int("keywordsLearned", metrics.KeywordsLearned))

	return nil
}
