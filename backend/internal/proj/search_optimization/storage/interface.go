package storage

import (
	"context"
	"time"
)

// SearchWeight представляет вес поля поиска
type SearchWeight struct {
	ID          int64     `json:"id" db:"id"`
	FieldName   string    `json:"field_name" db:"field_name"`
	Weight      float64   `json:"weight" db:"weight"`
	SearchType  string    `json:"search_type" db:"search_type"`
	ItemType    string    `json:"item_type" db:"item_type"`
	CategoryID  *int      `json:"category_id,omitempty" db:"category_id"`
	Description *string   `json:"description,omitempty" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Version     int       `json:"version" db:"version"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy   *int      `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy   *int      `json:"updated_by,omitempty" db:"updated_by"`
}

// SearchWeightHistory представляет историю изменений весов
type SearchWeightHistory struct {
	ID             int64     `json:"id" db:"id"`
	WeightID       int64     `json:"weight_id" db:"weight_id"`
	OldWeight      float64   `json:"old_weight" db:"old_weight"`
	NewWeight      float64   `json:"new_weight" db:"new_weight"`
	ChangeReason   string    `json:"change_reason" db:"change_reason"`
	ChangeMetadata *string   `json:"change_metadata,omitempty" db:"change_metadata"`
	ChangedBy      *int      `json:"changed_by,omitempty" db:"changed_by"`
	ChangedAt      time.Time `json:"changed_at" db:"changed_at"`
}

// WeightOptimizationResult представляет результат оптимизации весов
type WeightOptimizationResult struct {
	FieldName           string  `json:"field_name"`
	CurrentWeight       float64 `json:"current_weight"`
	OptimizedWeight     float64 `json:"optimized_weight"`
	ImprovementScore    float64 `json:"improvement_score"` // Ожидаемое улучшение CTR
	ConfidenceLevel     float64 `json:"confidence_level"`  // Уровень уверенности в рекомендации
	SampleSize          int     `json:"sample_size"`       // Количество данных для анализа
	CurrentCTR          float64 `json:"current_ctr"`
	PredictedCTR        float64 `json:"predicted_ctr"`
	StatisticalSigLevel float64 `json:"statistical_significance_level"`
}

// OptimizationSession представляет сессию оптимизации весов
type OptimizationSession struct {
	ID              int64                       `json:"id"`
	Status          string                      `json:"status"` // running, completed, failed
	StartTime       time.Time                   `json:"start_time"`
	EndTime         *time.Time                  `json:"end_time,omitempty"`
	TotalFields     int                         `json:"total_fields"`
	ProcessedFields int                         `json:"processed_fields"`
	Results         []*WeightOptimizationResult `json:"results,omitempty"`
	ErrorMessage    *string                     `json:"error_message,omitempty"`
	CreatedBy       int                         `json:"created_by"`
}

// BehaviorAnalysisData представляет данные для анализа поведения
type BehaviorAnalysisData struct {
	SearchQuery    string  `json:"search_query"`
	FieldName      string  `json:"field_name"`
	TotalSearches  int     `json:"total_searches"`
	TotalClicks    int     `json:"total_clicks"`
	CTR            float64 `json:"ctr"`
	AvgPosition    float64 `json:"avg_position"`
	Conversions    int     `json:"conversions"`
	ConversionRate float64 `json:"conversion_rate"`
}

// SynonymData представляет синоним для поиска
type SynonymData struct {
	ID        int64     `json:"id" db:"id"`
	Term      string    `json:"term" db:"term"`
	Synonym   string    `json:"synonym" db:"synonym"`
	Language  string    `json:"language" db:"language"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SearchOptimizationRepository интерфейс для работы с оптимизацией поиска
type SearchOptimizationRepository interface {
	// Управление весами
	GetSearchWeights(ctx context.Context, itemType string, categoryID *int) ([]*SearchWeight, error)
	GetSearchWeightByField(ctx context.Context, fieldName, searchType, itemType string, categoryID *int) (*SearchWeight, error)
	UpdateSearchWeight(ctx context.Context, id int64, weight float64, updatedBy int) error
	CreateSearchWeight(ctx context.Context, weight *SearchWeight) error
	BulkUpdateSearchWeights(ctx context.Context, weights []*SearchWeight, updatedBy int) error

	// История изменений
	GetWeightHistory(ctx context.Context, weightID int64, limit int) ([]*SearchWeightHistory, error)
	CreateWeightHistoryEntry(ctx context.Context, entry *SearchWeightHistory) error

	// Анализ поведенческих данных
	GetBehaviorAnalysisData(ctx context.Context, fromDate, toDate time.Time, fieldNames []string) ([]*BehaviorAnalysisData, error)
	GetSearchQueryMetrics(ctx context.Context, query string, fromDate, toDate time.Time) (*BehaviorAnalysisData, error)
	GetFieldPerformanceMetrics(ctx context.Context, fieldName string, fromDate, toDate time.Time) ([]*BehaviorAnalysisData, error)

	// Корреляционный анализ
	GetCTRByPosition(ctx context.Context, fieldName string, fromDate, toDate time.Time) (map[int]float64, error)
	GetWeightVsCTRCorrelation(ctx context.Context, fieldName string, fromDate, toDate time.Time) ([]struct {
		Weight float64 `json:"weight"`
		CTR    float64 `json:"ctr"`
		Count  int     `json:"count"`
	}, error)

	// Сессии оптимизации
	CreateOptimizationSession(ctx context.Context, session *OptimizationSession) error
	UpdateOptimizationSession(ctx context.Context, sessionID int64, status string, results []*WeightOptimizationResult, errorMsg *string) error
	GetOptimizationSession(ctx context.Context, sessionID int64) (*OptimizationSession, error)
	GetRecentOptimizationSessions(ctx context.Context, limit int) ([]*OptimizationSession, error)

	// Управление синонимами
	GetSynonyms(ctx context.Context, language, search string, offset, limit int) ([]*SynonymData, int, error)
	CreateSynonym(ctx context.Context, synonym *SynonymData) (int64, error)
	UpdateSynonym(ctx context.Context, synonymID int64, synonym *SynonymData) error
	DeleteSynonym(ctx context.Context, synonymID int64) error
	CheckSynonymExists(ctx context.Context, term, synonym, language string) (bool, error)
	CheckSynonymExistsByID(ctx context.Context, synonymID int64) (bool, error)
}
