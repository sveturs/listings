package service

import (
	"context"
	"backend/internal/proj/search_optimization/storage"
	"time"
)

// OptimizationParams параметры для оптимизации весов
type OptimizationParams struct {
	FieldNames      []string `json:"field_names,omitempty"` // Конкретные поля для оптимизации (если пусто - все)
	ItemType        string   `json:"item_type"`             // Тип элементов: marketplace, storefront, global
	CategoryID      *int     `json:"category_id,omitempty"` // ID категории (опционально)
	MinSampleSize   int      `json:"min_sample_size"`       // Минимальное количество данных для анализа
	ConfidenceLevel float64  `json:"confidence_level"`      // Требуемый уровень уверенности (0.8-0.99)
	LearningRate    float64  `json:"learning_rate"`         // Скорость обучения для градиентного спуска
	MaxIterations   int      `json:"max_iterations"`        // Максимальное количество итераций
	AnalysisPeriod  int      `json:"analysis_period_days"`  // Период анализа данных в днях
	AutoApply       bool     `json:"auto_apply"`            // Автоматически применить оптимизированные веса
}

// OptimizationConfig конфигурация по умолчанию для оптимизации
type OptimizationConfig struct {
	DefaultMinSampleSize   int     `json:"default_min_sample_size"`
	DefaultConfidenceLevel float64 `json:"default_confidence_level"`
	DefaultLearningRate    float64 `json:"default_learning_rate"`
	DefaultMaxIterations   int     `json:"default_max_iterations"`
	DefaultAnalysisPeriod  int     `json:"default_analysis_period_days"`
	MaxWeightChange        float64 `json:"max_weight_change"` // Максимальное изменение веса за одну оптимизацию
	MinWeight              float64 `json:"min_weight"`        // Минимальное значение веса
	MaxWeight              float64 `json:"max_weight"`        // Максимальное значение веса
}

// WeightValidationError ошибка валидации весов
type WeightValidationError struct {
	Field   string  `json:"field"`
	Weight  float64 `json:"weight"`
	Message string  `json:"message"`
}

func (e *WeightValidationError) Error() string {
	return e.Message
}

// ABTestConfig конфигурация A/B тестирования
type ABTestConfig struct {
	TestDuration      time.Duration `json:"test_duration"`      // Длительность теста
	TrafficSplit      float64       `json:"traffic_split"`      // Доля трафика для тестирования (0.1 = 10%)
	MinSampleSize     int           `json:"min_sample_size"`    // Минимальный размер выборки
	SignificanceLevel float64       `json:"significance_level"` // Уровень значимости для статистических тестов
}

// SearchOptimizationService интерфейс сервиса оптимизации поиска
type SearchOptimizationService interface {
	// Основные методы оптимизации
	StartOptimization(ctx context.Context, params *OptimizationParams, adminID int) (int64, error)
	GetOptimizationStatus(ctx context.Context, sessionID int64) (*storage.OptimizationSession, error)
	CancelOptimization(ctx context.Context, sessionID int64, adminID int) error

	// Анализ и предложения
	AnalyzeCurrentWeights(ctx context.Context, itemType string, categoryID *int, fromDate, toDate time.Time) ([]*storage.WeightOptimizationResult, error)
	GenerateWeightRecommendations(ctx context.Context, behaviorData []*storage.BehaviorAnalysisData) ([]*storage.WeightOptimizationResult, error)
	ValidateWeights(ctx context.Context, weights []*storage.SearchWeight) ([]*WeightValidationError, error)

	// Управление весами
	ApplyOptimizedWeights(ctx context.Context, sessionID int64, adminID int, selectedResults []int64) error
	RollbackWeights(ctx context.Context, weightIDs []int64, adminID int) error
	CreateWeightBackup(ctx context.Context, itemType string, categoryID *int, adminID int) error

	// A/B тестирование
	StartABTest(ctx context.Context, config *ABTestConfig, newWeights []*storage.SearchWeight, adminID int) (int64, error)
	GetABTestResults(ctx context.Context, testID int64) (*ABTestResult, error)

	// Статистика и метрики
	GetWeightPerformanceMetrics(ctx context.Context, fieldName string, fromDate, toDate time.Time) (*WeightPerformanceMetrics, error)
	GetOptimizationHistory(ctx context.Context, limit int) ([]*storage.OptimizationSession, error)
	GetFieldCorrelationMatrix(ctx context.Context, fromDate, toDate time.Time) (map[string]map[string]float64, error)

	// Машинное обучение
	TrainWeightOptimizationModel(ctx context.Context, behaviorData []*storage.BehaviorAnalysisData) (*MLModel, error)
	PredictOptimalWeights(ctx context.Context, model *MLModel, currentWeights []*storage.SearchWeight) ([]*storage.WeightOptimizationResult, error)

	// Конфигурация
	GetOptimizationConfig(ctx context.Context) (*OptimizationConfig, error)
	UpdateOptimizationConfig(ctx context.Context, config *OptimizationConfig, adminID int) error
}

// ABTestResult результаты A/B тестирования
type ABTestResult struct {
	TestID                  int64                   `json:"test_id"`
	Status                  string                  `json:"status"`
	StartDate               time.Time               `json:"start_date"`
	EndDate                 *time.Time              `json:"end_date,omitempty"`
	ControlGroupMetrics     *ABTestMetrics          `json:"control_group_metrics"`
	ExperimentGroupMetrics  *ABTestMetrics          `json:"experiment_group_metrics"`
	StatisticalSignificance float64                 `json:"statistical_significance"`
	RecommendedAction       string                  `json:"recommended_action"` // adopt, reject, continue
	TestedWeights           []*storage.SearchWeight `json:"tested_weights"`
}

// ABTestMetrics метрики для группы A/B теста
type ABTestMetrics struct {
	TotalSearches   int     `json:"total_searches"`
	TotalClicks     int     `json:"total_clicks"`
	CTR             float64 `json:"ctr"`
	ConversionRate  float64 `json:"conversion_rate"`
	AvgPosition     float64 `json:"avg_position"`
	BounceRate      float64 `json:"bounce_rate"`
	SessionDuration float64 `json:"session_duration"`
}

// WeightPerformanceMetrics метрики производительности веса
type WeightPerformanceMetrics struct {
	FieldName       string                        `json:"field_name"`
	CurrentWeight   float64                       `json:"current_weight"`
	PerformanceData []*WeightPerformanceDataPoint `json:"performance_data"`
	Correlation     float64                       `json:"correlation"`     // Корреляция веса с CTR
	TrendDirection  string                        `json:"trend_direction"` // increasing, decreasing, stable
	Recommendations []string                      `json:"recommendations"`
}

// WeightPerformanceDataPoint точка данных производительности веса
type WeightPerformanceDataPoint struct {
	Date        time.Time `json:"date"`
	Weight      float64   `json:"weight"`
	CTR         float64   `json:"ctr"`
	Searches    int       `json:"searches"`
	Clicks      int       `json:"clicks"`
	Conversions int       `json:"conversions"`
}

// MLModel модель машинного обучения для оптимизации весов
type MLModel struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`          // gradient_descent, random_forest, neural_network
	Parameters   map[string]interface{} `json:"parameters"`    // Гиперпараметры модели
	TrainingData interface{}            `json:"training_data"` // Данные для обучения
	Accuracy     float64                `json:"accuracy"`      // Точность модели
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
