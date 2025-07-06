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
	weights map[int]*CategoryWeights
}

func NewWeightManager() *WeightManager {
	return &WeightManager{
		weights: make(map[int]*CategoryWeights),
	}
}

func (wm *WeightManager) InitializeDefaultWeights() {
	// Недвижимость - Квартиры (ID: 1100)
	wm.weights[1100] = &CategoryWeights{
		CategoryID:   1100,
		CategoryName: "Квартиры",
		Weights: map[string]float64{
			"rooms":         0.9,  // Количество комнат - критично
			"area":          0.85, // Площадь
			"floor":         0.7,  // Этаж
			"property_type": 0.8,  // Тип недвижимости
			"location":      0.75, // Местоположение
			"condition":     0.6,  // Состояние
			"heating":       0.5,  // Отопление
			"parking":       0.4,  // Парковка
			"balcony":       0.3,  // Балкон
			"elevator":      0.25, // Лифт
		},
	}

	// Автомобили (ID: 2000)
	wm.weights[2000] = &CategoryWeights{
		CategoryID:   2000,
		CategoryName: "Автомобили",
		Weights: map[string]float64{
			"make":         0.9,  // Марка
			"model":        0.85, // Модель
			"year":         0.8,  // Год выпуска
			"body_type":    0.75, // Тип кузова
			"fuel_type":    0.7,  // Тип топлива
			"transmission": 0.65, // Коробка передач
			"engine":       0.6,  // Двигатель
			"color":        0.3,  // Цвет
			"mileage":      0.7,  // Пробег
			"condition":    0.6,  // Состояние
		},
	}

	// Электроника (ID: 3000)
	wm.weights[3000] = &CategoryWeights{
		CategoryID:   3000,
		CategoryName: "Электроника",
		Weights: map[string]float64{
			"brand":        0.9,  // Бренд
			"model":        0.85, // Модель
			"type":         0.8,  // Тип устройства
			"condition":    0.7,  // Состояние
			"warranty":     0.5,  // Гарантия
			"color":        0.3,  // Цвет
			"storage":      0.6,  // Объем памяти
			"display_size": 0.5,  // Размер экрана
		},
	}

	// Мебель (ID: 4000)
	wm.weights[4000] = &CategoryWeights{
		CategoryID:   4000,
		CategoryName: "Мебель",
		Weights: map[string]float64{
			"type":      0.9,  // Тип мебели
			"material":  0.8,  // Материал
			"style":     0.7,  // Стиль
			"size":      0.75, // Размер
			"color":     0.6,  // Цвет
			"condition": 0.65, // Состояние
			"brand":     0.5,  // Бренд
		},
	}

	// Одежда (ID: 5000)
	wm.weights[5000] = &CategoryWeights{
		CategoryID:   5000,
		CategoryName: "Одежда",
		Weights: map[string]float64{
			"type":      0.9,  // Тип одежды
			"brand":     0.8,  // Бренд
			"size":      0.85, // Размер
			"color":     0.6,  // Цвет
			"material":  0.7,  // Материал
			"season":    0.65, // Сезон
			"gender":    0.8,  // Пол
			"condition": 0.6,  // Состояние
		},
	}
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
