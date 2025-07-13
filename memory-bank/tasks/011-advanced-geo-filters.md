# Задача 011: Создание расширенных геофильтров для поиска объявлений

**Дата выполнения**: 2025-01-13
**Статус**: ✅ Завершено
**Ветка**: feature/maps-phase2

## Описание задачи
Реализовать расширенные геофильтры для улучшения поиска объявлений, включающие фильтрацию по времени пути, близости к точкам интереса и плотности объявлений.

## Реализованный функционал

### 1. Backend компоненты

#### Новые типы и структуры (`/backend/internal/proj/gis/types/advanced_filters.go`):
- `AdvancedGeoFilters` - основная структура для всех геофильтров
- `TravelTimeFilter` - фильтр по времени пути через изохроны
- `POIFilter` - фильтр по близости к точкам интереса
- `DensityFilter` - фильтр по плотности объявлений
- Типы POI: школы, больницы, метро, супермаркеты, парки, банки, аптеки, автобусные остановки

#### Новые сервисы:
- `IsochroneService` (`/backend/internal/proj/gis/service/isochrone_service.go`)
  - Интеграция с MapBox Isochrone API
  - Кеширование изохрон для оптимизации
  - Фильтрация объявлений по времени пути
  
- `POIService` (`/backend/internal/proj/gis/service/poi_service.go`)
  - Поиск точек интереса через MapBox Geocoding API
  - Кеширование найденных POI
  - Фильтрация объявлений по близости к POI
  
- `DensityService` (`/backend/internal/proj/gis/service/density_service.go`)
  - Анализ плотности объявлений по сетке
  - Материализованное представление для быстрого доступа
  - Генерация данных для тепловой карты

#### Новый обработчик (`/backend/internal/proj/gis/handler/advanced_filters.go`):
- `GetIsochrone` - получение изохроны для визуализации
- `SearchPOI` - поиск точек интереса
- `AnalyzeDensity` - анализ плотности в области
- `GetDensityHeatmap` - данные для тепловой карты
- `ApplyAdvancedFilters` - применение всех фильтров к списку объявлений

#### Интеграция с marketplace:
- Обновлен `SearchListingsAdvanced` для поддержки расширенных фильтров
- Добавлен метод `applyAdvancedGeoFilters` для вызова GIS сервиса
- Поддержка POST запросов для передачи сложных фильтров

### 2. Frontend компоненты

#### Новый компонент фильтров (`/frontend/svetu/src/components/GIS/AdvancedGeoFilters.tsx`):
- Интерфейс для всех типов фильтров
- Slider для времени пути (5-60 минут)
- Выбор типа транспорта
- Поиск POI с live результатами
- Контроль плотности объявлений
- Сохранение настроек в localStorage

#### Хук для управления фильтрами (`/frontend/svetu/src/hooks/useAdvancedGeoFilters.ts`):
- Синхронизация с URL параметрами
- Управление состоянием фильтров
- События для обновления результатов

#### Обновленный HomePage (`/frontend/svetu/src/components/marketplace/HomePageWithFilters.tsx`):
- Боковая панель с фильтрами (desktop)
- Мобильная версия фильтров
- Индикатор активных фильтров
- Интеграция с системой поиска

### 3. База данных

#### Новые таблицы (миграция 036):
- `gis_isochrone_cache` - кеширование изохрон
- `gis_poi_cache` - кеширование точек интереса
- `gis_filter_analytics` - аналитика использования
- `gis_listing_density_grid` - материализованное представление плотности

### 4. API endpoints

#### GIS модуль (`/api/v1/gis/advanced/`):
- `POST /isochrone` - получение изохроны
- `GET /poi/search` - поиск точек интереса
- `POST /density/analyze` - анализ плотности
- `GET /density/heatmap` - данные для тепловой карты
- `POST /apply-filters` - применение всех фильтров

#### Обновленный marketplace:
- `POST /api/v1/marketplace/search` - поддержка расширенных фильтров в теле запроса

## Технические особенности

### Оптимизации:
1. Кеширование результатов MapBox API
2. Материализованное представление для плотности
3. Пакетная обработка фильтров
4. Использование PostGIS для геозапросов

### Безопасность:
1. Валидация всех входных параметров
2. Ограничения на размер области поиска
3. Rate limiting для внешних API
4. Санитизация пользовательского ввода

### Производительность:
1. Индексы на всех геополях
2. Оптимизированные SQL запросы
3. Параллельное применение фильтров
4. Прогрессивная загрузка результатов

## Конфигурация

### Переменные окружения:
- `MAPBOX_ACCESS_TOKEN` - токен для MapBox API (обязательно)

### Feature flags:
- Фильтры работают независимо и могут быть отключены по отдельности

## Использование

### Пример запроса с расширенными фильтрами:
```javascript
const result = await UnifiedSearchService.search({
  query: '',
  page: 1,
  limit: 20,
  advanced_geo_filters: {
    travel_time: {
      center_lat: 44.8176,
      center_lng: 20.4649,
      max_minutes: 30,
      transport_mode: 'walking'
    },
    poi_filter: {
      poi_type: 'metro',
      max_distance: 1000,
      min_count: 1
    },
    density_filter: {
      avoid_crowded: true
    }
  }
});
```

### URL параметры:
- `travel_time=44.8176,20.4649,30,walking`
- `poi=metro,1000,1`
- `density=avoid`

## Будущие улучшения

1. **Кастомные изохроны**:
   - Учет расписания общественного транспорта
   - Время суток и пробки
   - Мультимодальные маршруты

2. **Расширенные POI**:
   - Больше категорий точек интереса
   - Рейтинги и отзывы POI
   - Кастомные пользовательские POI

3. **Умная плотность**:
   - Прогнозирование изменения плотности
   - Рекомендации оптимальных районов
   - Исторические данные плотности

4. **Визуализация**:
   - 3D тепловые карты
   - Анимированные изохроны
   - Интерактивные слои фильтров

## Файлы проекта

### Backend:
- `/backend/internal/proj/gis/types/advanced_filters.go`
- `/backend/internal/proj/gis/service/isochrone_service.go`
- `/backend/internal/proj/gis/service/poi_service.go`
- `/backend/internal/proj/gis/service/density_service.go`
- `/backend/internal/proj/gis/handler/advanced_filters.go`
- `/backend/internal/proj/marketplace/service/marketplace.go` (обновлен)
- `/backend/internal/domain/search/types.go` (обновлен)
- `/backend/migrations/036_add_advanced_geo_filters.up.sql`
- `/backend/migrations/036_add_advanced_geo_filters.down.sql`

### Frontend:
- `/frontend/svetu/src/components/GIS/AdvancedGeoFilters.tsx`
- `/frontend/svetu/src/components/marketplace/HomePageWithFilters.tsx`
- `/frontend/svetu/src/hooks/useAdvancedGeoFilters.ts`
- `/frontend/svetu/src/services/unifiedSearch.ts` (обновлен)
- `/frontend/svetu/src/messages/ru.json` (обновлен)
- `/frontend/svetu/src/messages/en.json` (обновлен)