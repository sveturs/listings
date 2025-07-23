package service

// AttributeWeight представляет вес атрибута для конкретной категории
type AttributeWeight struct {
	AttributeName string  `json:"attribute_name"`
	Weight        float64 `json:"weight"`
	Category      string  `json:"category"`
}

// CategoryWeights содержит веса для всех атрибутов категории
type CategoryWeights struct {
	CategoryID   int                `json:"category_id"`
	CategoryName string             `json:"category_name"`
	Weights      map[string]float64 `json:"weights"`
	ParentID     *int               `json:"parent_id,omitempty"`
}

// WeightManager управляет весами атрибутов
type WeightManager struct {
	weights        map[int]*CategoryWeights
	defaultWeights map[string]float64
}

func NewWeightManager(categoryWeights map[int]map[string]float64, defaultWeights map[string]float64) *WeightManager {
	wm := &WeightManager{
		weights:        make(map[int]*CategoryWeights),
		defaultWeights: defaultWeights,
	}

	// Преобразуем веса из конфигурации в формат CategoryWeights
	for categoryID, weights := range categoryWeights {
		wm.weights[categoryID] = &CategoryWeights{
			CategoryID: categoryID,
			Weights:    weights,
		}
	}

	return wm
}

// InitializeDefaultWeights теперь не нужен - веса загружаются из конфигурации
// Метод оставлен для обратной совместимости
func (wm *WeightManager) InitializeDefaultWeights() {
	// Ничего не делаем - веса уже загружены из конфигурации
}

func (wm *WeightManager) GetCategoryWeights(categoryID int) map[string]float64 {
	if weights, exists := wm.weights[categoryID]; exists {
		return weights.Weights
	}

	// Попытка найти веса родительской категории
	parentWeights := wm.findParentWeights(categoryID)
	if parentWeights != nil {
		return parentWeights
	}

	// Возвращаем базовые веса
	return wm.getDefaultWeights()
}

func (wm *WeightManager) getDefaultWeights() map[string]float64 {
	if len(wm.defaultWeights) > 0 {
		return wm.defaultWeights
	}

	// Fallback на случай если defaultWeights не заданы
	return map[string]float64{
		"brand":     0.7,
		"model":     0.65,
		"type":      0.8,
		"condition": 0.6,
		"color":     0.4,
		"size":      0.5,
		"material":  0.5,
	}
}

func (wm *WeightManager) findParentWeights(categoryID int) map[string]float64 {
	// Здесь должна быть логика поиска родительской категории
	// Пока возвращаем nil для упрощения
	return nil
}
