package translation_admin

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// CostTracker отслеживает расходы на AI провайдеров
type CostTracker struct {
	redis    *redis.Client
	inMemory map[string]*ProviderCosts
	mu       sync.RWMutex
	useRedis bool
}

// ProviderCosts расходы по провайдеру
type ProviderCosts struct {
	Provider      string             `json:"provider"`
	TotalCost     float64            `json:"total_cost"`
	TotalTokens   int64              `json:"total_tokens"`
	TotalRequests int64              `json:"total_requests"`
	LastUpdated   time.Time          `json:"last_updated"`
	DailyCosts    map[string]float64 `json:"daily_costs"`  // map[date]cost
	HourlyCosts   map[string]float64 `json:"hourly_costs"` // map[hour]cost
}

// CostConfig конфигурация стоимости для провайдеров
type CostConfig struct {
	OpenAI struct {
		GPT4InputCost   float64 `json:"gpt4_input_cost"`   // per 1K tokens
		GPT4OutputCost  float64 `json:"gpt4_output_cost"`  // per 1K tokens
		GPT35InputCost  float64 `json:"gpt35_input_cost"`  // per 1K tokens
		GPT35OutputCost float64 `json:"gpt35_output_cost"` // per 1K tokens
	} `json:"openai"`

	Google struct {
		CostPer1MChars float64 `json:"cost_per_1m_chars"` // $20 per 1M chars
	} `json:"google"`

	DeepL struct {
		CostPer1MChars float64 `json:"cost_per_1m_chars"` // €20 per 1M chars
	} `json:"deepl"`

	Claude struct {
		InputCost  float64 `json:"input_cost"`  // per 1K tokens
		OutputCost float64 `json:"output_cost"` // per 1K tokens
	} `json:"claude"`
}

// DefaultCostConfig возвращает конфигурацию стоимости по умолчанию
func DefaultCostConfig() *CostConfig {
	cfg := &CostConfig{}

	// OpenAI pricing (as of 2024)
	cfg.OpenAI.GPT4InputCost = 0.03    // $0.03 per 1K input tokens
	cfg.OpenAI.GPT4OutputCost = 0.06   // $0.06 per 1K output tokens
	cfg.OpenAI.GPT35InputCost = 0.001  // $0.001 per 1K input tokens
	cfg.OpenAI.GPT35OutputCost = 0.002 // $0.002 per 1K output tokens

	// Google Translate pricing
	cfg.Google.CostPer1MChars = 20.0 // $20 per 1M characters

	// DeepL pricing
	cfg.DeepL.CostPer1MChars = 20.0 // €20 per 1M characters (примерно)

	// Claude pricing
	cfg.Claude.InputCost = 0.015  // $0.015 per 1K input tokens
	cfg.Claude.OutputCost = 0.075 // $0.075 per 1K output tokens

	return cfg
}

// NewCostTracker создает новый трекер расходов
func NewCostTracker(ctx context.Context, redisClient *redis.Client) *CostTracker {
	tracker := &CostTracker{
		redis:    redisClient,
		inMemory: make(map[string]*ProviderCosts),
		useRedis: redisClient != nil,
	}

	// Инициализируем провайдеров
	providers := []string{"openai", "google", "deepl", "claude"}
	for _, provider := range providers {
		tracker.inMemory[provider] = &ProviderCosts{
			Provider:    provider,
			DailyCosts:  make(map[string]float64),
			HourlyCosts: make(map[string]float64),
		}
	}

	// Загружаем данные из Redis если доступно
	if tracker.useRedis {
		tracker.loadFromRedis(ctx)
	}

	return tracker
}

// TrackOpenAIUsage отслеживает использование OpenAI
func (t *CostTracker) TrackOpenAIUsage(ctx context.Context, model string, inputTokens, outputTokens int) error {
	cfg := DefaultCostConfig()

	var cost float64
	if model == "gpt-4" || model == "gpt-4-turbo" {
		cost = (float64(inputTokens)/1000.0)*cfg.OpenAI.GPT4InputCost +
			(float64(outputTokens)/1000.0)*cfg.OpenAI.GPT4OutputCost
	} else {
		cost = (float64(inputTokens)/1000.0)*cfg.OpenAI.GPT35InputCost +
			(float64(outputTokens)/1000.0)*cfg.OpenAI.GPT35OutputCost
	}

	return t.trackUsage(ctx, "openai", cost, int64(inputTokens+outputTokens))
}

// TrackGoogleUsage отслеживает использование Google Translate
func (t *CostTracker) TrackGoogleUsage(ctx context.Context, characters int) error {
	cfg := DefaultCostConfig()
	cost := (float64(characters) / 1000000.0) * cfg.Google.CostPer1MChars
	return t.trackUsage(ctx, "google", cost, int64(characters))
}

// TrackDeepLUsage отслеживает использование DeepL
func (t *CostTracker) TrackDeepLUsage(ctx context.Context, characters int) error {
	cfg := DefaultCostConfig()
	cost := (float64(characters) / 1000000.0) * cfg.DeepL.CostPer1MChars
	return t.trackUsage(ctx, "deepl", cost, int64(characters))
}

// TrackClaudeUsage отслеживает использование Claude
func (t *CostTracker) TrackClaudeUsage(ctx context.Context, inputTokens, outputTokens int) error {
	cfg := DefaultCostConfig()
	cost := (float64(inputTokens)/1000.0)*cfg.Claude.InputCost +
		(float64(outputTokens)/1000.0)*cfg.Claude.OutputCost
	return t.trackUsage(ctx, "claude", cost, int64(inputTokens+outputTokens))
}

// trackUsage внутренний метод для отслеживания использования
func (t *CostTracker) trackUsage(ctx context.Context, provider string, cost float64, tokens int64) error {
	now := time.Now()
	dateKey := now.Format("2006-01-02")
	hourKey := now.Format("2006-01-02T15")

	log.Info().
		Str("provider", provider).
		Float64("cost", cost).
		Int64("tokens", tokens).
		Str("date", dateKey).
		Str("hour", hourKey).
		Msg("trackUsage called")

	t.mu.Lock()
	defer t.mu.Unlock()

	// Обновляем in-memory данные
	if costs, ok := t.inMemory[provider]; ok {
		costs.TotalCost += cost
		costs.TotalTokens += tokens
		costs.TotalRequests++
		costs.LastUpdated = now
		costs.DailyCosts[dateKey] += cost
		costs.HourlyCosts[hourKey] += cost

		log.Info().
			Str("provider", provider).
			Float64("total_cost", costs.TotalCost).
			Int64("total_tokens", costs.TotalTokens).
			Int64("total_requests", costs.TotalRequests).
			Msg("Updated in-memory costs")

		// Очищаем старые данные (храним только последние 30 дней и 24 часа)
		t.cleanupOldData(costs)
	}

	// Сохраняем в Redis если доступно
	if t.useRedis {
		log.Info().Str("provider", provider).Msg("Saving to Redis")
		err := t.saveToRedis(ctx, provider)
		if err != nil {
			log.Error().Err(err).Str("provider", provider).Msg("Failed to save to Redis")
			return err
		}
		log.Info().Str("provider", provider).Msg("Successfully saved to Redis")
	} else {
		log.Info().Msg("Redis not available, using in-memory storage only")
	}

	return nil
}

// GetProviderCosts возвращает расходы по провайдеру
func (t *CostTracker) GetProviderCosts(ctx context.Context, provider string) (*ProviderCosts, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if costs, ok := t.inMemory[provider]; ok {
		// Возвращаем копию чтобы избежать race conditions
		costsCopy := *costs
		costsCopy.DailyCosts = make(map[string]float64)
		costsCopy.HourlyCosts = make(map[string]float64)

		for k, v := range costs.DailyCosts {
			costsCopy.DailyCosts[k] = v
		}
		for k, v := range costs.HourlyCosts {
			costsCopy.HourlyCosts[k] = v
		}

		return &costsCopy, nil
	}

	return nil, fmt.Errorf("provider %s not found", provider)
}

// GetAllProvidersCosts возвращает расходы по всем провайдерам
func (t *CostTracker) GetAllProvidersCosts(ctx context.Context) (map[string]*ProviderCosts, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[string]*ProviderCosts)

	for provider, costs := range t.inMemory {
		// Создаем копию
		costsCopy := *costs
		costsCopy.DailyCosts = make(map[string]float64)
		costsCopy.HourlyCosts = make(map[string]float64)

		for k, v := range costs.DailyCosts {
			costsCopy.DailyCosts[k] = v
		}
		for k, v := range costs.HourlyCosts {
			costsCopy.HourlyCosts[k] = v
		}

		result[provider] = &costsCopy
	}

	return result, nil
}

// GetDailyCosts возвращает расходы за день
func (t *CostTracker) GetDailyCosts(ctx context.Context, date string) (map[string]float64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[string]float64)

	for provider, costs := range t.inMemory {
		if dailyCost, ok := costs.DailyCosts[date]; ok {
			result[provider] = dailyCost
		}
	}

	return result, nil
}

// GetMonthlyCosts возвращает расходы за месяц
func (t *CostTracker) GetMonthlyCosts(ctx context.Context, year, month int) (map[string]float64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[string]float64)
	monthPrefix := fmt.Sprintf("%04d-%02d", year, month)

	for provider, costs := range t.inMemory {
		monthTotal := 0.0
		for date, cost := range costs.DailyCosts {
			if len(date) >= 7 && date[:7] == monthPrefix {
				monthTotal += cost
			}
		}
		if monthTotal > 0 {
			result[provider] = monthTotal
		}
	}

	return result, nil
}

// GetCostsSummary возвращает сводку расходов
func (t *CostTracker) GetCostsSummary(ctx context.Context) (map[string]interface{}, error) {
	allCosts, err := t.GetAllProvidersCosts(ctx)
	if err != nil {
		return nil, err
	}

	// Считаем общие расходы
	var totalCost float64
	var totalTokens int64
	var totalRequests int64

	for _, costs := range allCosts {
		totalCost += costs.TotalCost
		totalTokens += costs.TotalTokens
		totalRequests += costs.TotalRequests
	}

	// Расходы за сегодня
	today := time.Now().Format("2006-01-02")
	todayCosts, _ := t.GetDailyCosts(ctx, today)

	var todayTotal float64
	for _, cost := range todayCosts {
		todayTotal += cost
	}

	// Расходы за текущий месяц
	now := time.Now()
	monthCosts, _ := t.GetMonthlyCosts(ctx, now.Year(), int(now.Month()))

	var monthTotal float64
	for _, cost := range monthCosts {
		monthTotal += cost
	}

	return map[string]interface{}{
		"total_cost":        totalCost,
		"total_tokens":      totalTokens,
		"total_requests":    totalRequests,
		"today_cost":        todayTotal,
		"month_cost":        monthTotal,
		"by_provider":       allCosts,
		"today_by_provider": todayCosts,
		"month_by_provider": monthCosts,
	}, nil
}

// ResetProviderCosts сбрасывает счетчики для провайдера
func (t *CostTracker) ResetProviderCosts(ctx context.Context, provider string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if costs, ok := t.inMemory[provider]; ok {
		costs.TotalCost = 0
		costs.TotalTokens = 0
		costs.TotalRequests = 0
		costs.DailyCosts = make(map[string]float64)
		costs.HourlyCosts = make(map[string]float64)
		costs.LastUpdated = time.Now()

		if t.useRedis {
			key := fmt.Sprintf("translation:costs:%s", provider)
			return t.redis.Del(ctx, key).Err()
		}
	}

	return nil
}

// Вспомогательные методы

func (t *CostTracker) cleanupOldData(costs *ProviderCosts) {
	now := time.Now()

	// Удаляем данные старше 30 дней
	cutoffDate := now.AddDate(0, 0, -30).Format("2006-01-02")
	for date := range costs.DailyCosts {
		if date < cutoffDate {
			delete(costs.DailyCosts, date)
		}
	}

	// Удаляем данные старше 24 часов
	cutoffHour := now.Add(-24 * time.Hour).Format("2006-01-02T15")
	for hour := range costs.HourlyCosts {
		if hour < cutoffHour {
			delete(costs.HourlyCosts, hour)
		}
	}
}

func (t *CostTracker) saveToRedis(ctx context.Context, provider string) error {
	if costs, ok := t.inMemory[provider]; ok {
		data, err := json.Marshal(costs)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("translation:costs:%s", provider)
		return t.redis.Set(ctx, key, data, 0).Err()
	}

	return nil
}

func (t *CostTracker) loadFromRedis(ctx context.Context) {
	providers := []string{"openai", "google", "deepl", "claude"}

	for _, provider := range providers {
		key := fmt.Sprintf("translation:costs:%s", provider)
		data, err := t.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var costs ProviderCosts
		if err := json.Unmarshal([]byte(data), &costs); err != nil {
			log.Error().Err(err).Str("provider", provider).Msg("Failed to unmarshal costs from Redis")
			continue
		}

		t.inMemory[provider] = &costs
	}
}

// GetCostAlerts возвращает алерты о превышении бюджета
func (t *CostTracker) GetCostAlerts(ctx context.Context, dailyLimit, monthlyLimit float64) ([]string, error) {
	var alerts []string

	// Проверяем дневной лимит
	today := time.Now().Format("2006-01-02")
	todayCosts, _ := t.GetDailyCosts(ctx, today)

	var todayTotal float64
	for _, cost := range todayCosts {
		todayTotal += cost
	}

	if todayTotal > dailyLimit {
		alerts = append(alerts, fmt.Sprintf("Daily cost limit exceeded: $%.2f > $%.2f", todayTotal, dailyLimit))
	} else if todayTotal > dailyLimit*0.8 {
		alerts = append(alerts, fmt.Sprintf("Daily cost approaching limit: $%.2f (80%% of $%.2f)", todayTotal, dailyLimit))
	}

	// Проверяем месячный лимит
	now := time.Now()
	monthCosts, _ := t.GetMonthlyCosts(ctx, now.Year(), int(now.Month()))

	var monthTotal float64
	for _, cost := range monthCosts {
		monthTotal += cost
	}

	if monthTotal > monthlyLimit {
		alerts = append(alerts, fmt.Sprintf("Monthly cost limit exceeded: $%.2f > $%.2f", monthTotal, monthlyLimit))
	} else if monthTotal > monthlyLimit*0.8 {
		alerts = append(alerts, fmt.Sprintf("Monthly cost approaching limit: $%.2f (80%% of $%.2f)", monthTotal, monthlyLimit))
	}

	return alerts, nil
}
