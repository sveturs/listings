package service

import (
	"context"
	"testing"

	"backend/internal/config"
	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestSimilarityCalculator_CalculateSimilarity(t *testing.T) {
	searchWeights := config.GetDefaultSearchWeights()
	calculator := NewSimilarityCalculator(searchWeights)
	ctx := context.Background()

	// Исходное объявление
	source := &models.MarketplaceListing{
		ID:         1,
		CategoryID: 1100,
		Title:      "3-комнатная квартира в центре",
		Price:      200000,
		City:       "Белград",
		Attributes: []models.ListingAttributeValue{
			{AttributeName: "rooms", DisplayValue: "3"},
			{AttributeName: "area", DisplayValue: "85"},
		},
	}

	// Очень похожее объявление
	similar := &models.MarketplaceListing{
		ID:         2,
		CategoryID: 1100,
		Title:      "3-комнатная квартира люкс",
		Price:      220000,
		City:       "Белград",
		Attributes: []models.ListingAttributeValue{
			{AttributeName: "rooms", DisplayValue: "3"},
			{AttributeName: "area", DisplayValue: "90"},
		},
	}

	// Непохожее объявление
	different := &models.MarketplaceListing{
		ID:         3,
		CategoryID: 2000,
		Title:      "Автомобиль BMW",
		Price:      15000,
		City:       "Нови-Сад",
	}

	// Тест для похожего объявления
	score, err := calculator.CalculateSimilarity(ctx, source, similar)
	assert.NoError(t, err)
	assert.NotNil(t, score)
	assert.Greater(t, score.TotalScore, 0.5) // Должен быть высокий скор
	assert.Contains(t, score.MatchReasons, "Та же категория")
	assert.Contains(t, score.MatchReasons, "Похожая цена")
	assert.Contains(t, score.MatchReasons, "То же место")

	// Тест для непохожего объявления
	score, err = calculator.CalculateSimilarity(ctx, source, different)
	assert.NoError(t, err)
	assert.NotNil(t, score)
	assert.Less(t, score.TotalScore, 0.3) // Должен быть низкий скор
}

func TestSimilarityCalculator_BasicFunctionality(t *testing.T) {
	searchWeights := config.GetDefaultSearchWeights()
	calculator := NewSimilarityCalculator(searchWeights)
	assert.NotNil(t, calculator)

	// Тест что конструктор создает валидный объект
	assert.NotNil(t, calculator)
}
