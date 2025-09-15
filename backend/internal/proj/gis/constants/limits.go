package constants

// Лимиты для API согласно плану MAP_IMPLEMENTATION_PLAN.md
const (
	// DEFAULT_LIMIT - стандартный лимит для запросов
	DEFAULT_LIMIT = 1000 // увеличен с 50 до 1000

	// MAX_LIMIT - максимальный лимит для запросов
	MAX_LIMIT = 5000 // увеличен с 200 до 5000

	// DEFAULT_NEARBY_LIMIT - стандартный лимит для nearby запросов
	DEFAULT_NEARBY_LIMIT = 200 // увеличен с 20 до 200

	// DEFAULT_RADIUS_KM - стандартный радиус поиска в километрах
	DEFAULT_RADIUS_KM = 5.0

	// DEFAULT_RADIUS_METERS - стандартный радиус поиска в метрах
	DEFAULT_RADIUS_METERS = 5000
)
