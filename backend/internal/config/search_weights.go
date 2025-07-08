package config

// SearchWeights содержит все веса для поисковой системы
type SearchWeights struct {
	// OpenSearchBoosts содержит boost веса для OpenSearch запросов
	OpenSearchBoosts OpenSearchBoostWeights `yaml:"opensearch_boosts"`

	// CategoryAttributeWeights содержит веса атрибутов для разных категорий
	CategoryAttributeWeights map[int]map[string]float64 `yaml:"category_attribute_weights"`

	// DefaultAttributeWeights веса атрибутов по умолчанию для категорий без специфических весов
	DefaultAttributeWeights map[string]float64 `yaml:"default_attribute_weights"`

	// SimilarityScoring веса для расчета общей похожести объявлений
	SimilarityScoring SimilarityScoringWeights `yaml:"similarity_scoring"`

	// UnifiedSearchWeights веса для unified search
	UnifiedSearchWeights UnifiedSearchWeights `yaml:"unified_search_weights"`

	// FuzzySearchThreshold минимальный порог похожести для fuzzy search
	FuzzySearchThreshold float64 `yaml:"fuzzy_search_threshold"`
}

// OpenSearchBoostWeights содержит boost веса для различных полей в OpenSearch
type OpenSearchBoostWeights struct {
	// Основные поля
	Title            float64 `yaml:"title"`
	TitleNgram       float64 `yaml:"title_ngram"`
	Description      float64 `yaml:"description"`
	TranslationTitle float64 `yaml:"translation_title"`
	TranslationDesc  float64 `yaml:"translation_description"`

	// Атрибуты
	AttributeTextValue        float64 `yaml:"attribute_text_value"`
	AttributeDisplayValue     float64 `yaml:"attribute_display_value"`
	AttributeTextValueKeyword float64 `yaml:"attribute_text_value_keyword"`
	AttributeGeneralBoost     float64 `yaml:"attribute_general_boost"`

	// Специализированные поля
	RealEstateAttributesCombined float64 `yaml:"real_estate_attributes_combined"`
	PropertyType                 float64 `yaml:"property_type"`
	RoomsText                    float64 `yaml:"rooms_text"`
	CarMake                      float64 `yaml:"car_make"`
	CarModel                     float64 `yaml:"car_model"`
	CarKeywords                  float64 `yaml:"car_keywords"`

	// Веса для поиска по словам
	PerWordTitle                float64 `yaml:"per_word_title"`
	PerWordDescription          float64 `yaml:"per_word_description"`
	PerWordAllAttributes        float64 `yaml:"per_word_all_attributes"`
	PerWordRealEstateAttributes float64 `yaml:"per_word_real_estate_attributes"`
	PerWordRoomsText            float64 `yaml:"per_word_rooms_text"`

	// Прочие
	AutomotiveAttributePriority float64 `yaml:"automotive_attribute_priority"`
	SynonymBoost                float64 `yaml:"synonym_boost"`
}

// SimilarityScoringWeights содержит веса для расчета похожести
type SimilarityScoringWeights struct {
	// Веса когда категории совпадают (CategoryScore = 1.0)
	SameCategoryWeights struct {
		Category   float64 `yaml:"category"`
		Attributes float64 `yaml:"attributes"`
		Text       float64 `yaml:"text"`
		Price      float64 `yaml:"price"`
		Location   float64 `yaml:"location"`
	} `yaml:"same_category"`

	// Веса когда категории из одной группы (CategoryScore >= 0.6)
	SimilarCategoryWeights struct {
		Category   float64 `yaml:"category"`
		Attributes float64 `yaml:"attributes"`
		Text       float64 `yaml:"text"`
		Price      float64 `yaml:"price"`
		Location   float64 `yaml:"location"`
	} `yaml:"similar_category"`

	// Веса когда категории разные (CategoryScore < 0.6)
	DifferentCategoryWeights struct {
		Category   float64 `yaml:"category"`
		Attributes float64 `yaml:"attributes"`
		Text       float64 `yaml:"text"`
		Price      float64 `yaml:"price"`
		Location   float64 `yaml:"location"`
	} `yaml:"different_category"`
}

// UnifiedSearchWeights веса для unified search
type UnifiedSearchWeights struct {
	ExactTitleMatch    float64 `yaml:"exact_title_match"`
	PartialTitleMatch  float64 `yaml:"partial_title_match"`
	DescriptionMatch   float64 `yaml:"description_match"`
	PopularityMaxBoost float64 `yaml:"popularity_max_boost"`
	FreshnessWeek      float64 `yaml:"freshness_week"`
	FreshnessMonth     float64 `yaml:"freshness_month"`
	FreshnessQuarter   float64 `yaml:"freshness_quarter"`
}

// GetDefaultSearchWeights возвращает дефолтные веса поиска
func GetDefaultSearchWeights() *SearchWeights {
	return &SearchWeights{
		OpenSearchBoosts: OpenSearchBoostWeights{
			Title:            5.0,
			TitleNgram:       2.0,
			Description:      2.0,
			TranslationTitle: 4.0,
			TranslationDesc:  1.5,

			AttributeTextValue:        5.0,
			AttributeDisplayValue:     4.0,
			AttributeTextValueKeyword: 5.0,
			AttributeGeneralBoost:     4.0,

			RealEstateAttributesCombined: 5.0,
			PropertyType:                 4.0,
			RoomsText:                    4.0,
			CarMake:                      5.0,
			CarModel:                     4.0,
			CarKeywords:                  5.0,

			PerWordTitle:                2.0,
			PerWordDescription:          1.0,
			PerWordAllAttributes:        2.0,
			PerWordRealEstateAttributes: 3.0,
			PerWordRoomsText:            2.5,

			AutomotiveAttributePriority: 2.0,
			SynonymBoost:                0.5,
		},

		CategoryAttributeWeights: map[int]map[string]float64{
			// Недвижимость - Квартиры
			1100: {
				"rooms":         0.9,
				"area":          0.85,
				"floor":         0.7,
				"property_type": 0.8,
				"location":      0.75,
				"condition":     0.6,
				"heating":       0.5,
				"parking":       0.4,
				"balcony":       0.3,
				"elevator":      0.25,
			},
			// Автомобили
			2000: {
				"make":         0.9,
				"model":        0.85,
				"year":         0.8,
				"body_type":    0.75,
				"fuel_type":    0.7,
				"transmission": 0.65,
				"engine":       0.6,
				"mileage":      0.7,
				"condition":    0.6,
				"color":        0.3,
			},
			// Электроника
			3000: {
				"brand":        0.9,
				"model":        0.85,
				"type":         0.8,
				"condition":    0.7,
				"warranty":     0.5,
				"storage":      0.6,
				"display_size": 0.5,
				"color":        0.3,
			},
			// Мебель
			4000: {
				"type":      0.9,
				"material":  0.8,
				"style":     0.7,
				"size":      0.75,
				"condition": 0.65,
				"color":     0.6,
				"brand":     0.5,
			},
			// Одежда
			5000: {
				"type":      0.9,
				"size":      0.85,
				"brand":     0.8,
				"gender":    0.8,
				"material":  0.7,
				"season":    0.65,
				"color":     0.6,
				"condition": 0.6,
			},
		},

		DefaultAttributeWeights: map[string]float64{
			"type":      0.8,
			"brand":     0.7,
			"model":     0.65,
			"condition": 0.6,
			"size":      0.5,
			"material":  0.5,
			"color":     0.4,
		},

		SimilarityScoring: SimilarityScoringWeights{
			SameCategoryWeights: struct {
				Category   float64 `yaml:"category"`
				Attributes float64 `yaml:"attributes"`
				Text       float64 `yaml:"text"`
				Price      float64 `yaml:"price"`
				Location   float64 `yaml:"location"`
			}{
				Category:   0.3,
				Attributes: 0.3,
				Text:       0.2,
				Price:      0.15,
				Location:   0.05,
			},
			SimilarCategoryWeights: struct {
				Category   float64 `yaml:"category"`
				Attributes float64 `yaml:"attributes"`
				Text       float64 `yaml:"text"`
				Price      float64 `yaml:"price"`
				Location   float64 `yaml:"location"`
			}{
				Category:   0.2,
				Attributes: 0.2,
				Text:       0.25,
				Price:      0.25,
				Location:   0.1,
			},
			DifferentCategoryWeights: struct {
				Category   float64 `yaml:"category"`
				Attributes float64 `yaml:"attributes"`
				Text       float64 `yaml:"text"`
				Price      float64 `yaml:"price"`
				Location   float64 `yaml:"location"`
			}{
				Category:   0.1,
				Attributes: 0.15,
				Text:       0.2,
				Price:      0.35,
				Location:   0.2,
			},
		},

		UnifiedSearchWeights: UnifiedSearchWeights{
			ExactTitleMatch:    5.0,
			PartialTitleMatch:  3.0,
			DescriptionMatch:   2.0,
			PopularityMaxBoost: 1.0,
			FreshnessWeek:      0.5,
			FreshnessMonth:     0.3,
			FreshnessQuarter:   0.1,
		},

		FuzzySearchThreshold: 0.3,
	}
}
