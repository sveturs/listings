# Фаза 2.5: Полная система управления поиском (Search Management System)

## 📍 Статус: ПЛАНИРУЕТСЯ
**Предполагаемая длительность**: 7-10 дней (расширено из-за полного охвата)
**Приоритет**: КРИТИЧЕСКИЙ (перед Фазой 3)

## 🎯 Цель
Создать комплексную систему управления всеми аспектами поиска для полного контроля, мониторинга и оптимизации.

## 📋 Архитектура админки

### 1. 🎯 Dashboard - Главная панель мониторинга (1 день)
- [ ] **Метрики реального времени**
  - QPS (Queries per second)
  - Среднее время ответа
  - Активные пользователи
  - Текущая нагрузка на OpenSearch/PostgreSQL

- [ ] **Статистика и графики**
  - График объёма поисков (час/день/неделя)
  - График времени ответа
  - Топ-10 запросов в реальном времени
  - Процент "нулевых" результатов

### 2. 📊 Analytics - Глубокая аналитика (2 дня)
- [ ] **Анализ запросов**
  - Топ-100 запросов с трендами
  - Запросы без результатов с предложениями
  - Медленные запросы (P95, P99)
  - Поисковые паттерны (query refinement)
  - Сезонность и тренды

- [ ] **Анализ поведения**
  - CTR по позициям (heatmap)
  - Средняя длина запроса
  - Использование фильтров
  - Успешность поисковых сессий
  - Процент уточнения запросов

- [ ] **Сегментация**
  - По языкам (ru/sr/en)
  - По категориям
  - По устройствам
  - По времени суток

### 3. ⚙️ Configuration - Управление всеми параметрами (2 дня)
- [ ] **Веса и ранжирование**
  - Управление весами категорий (30%)
  - Веса атрибутов по категориям
  - Текстовая релевантность (20%)
  - Ценовой скоринг (15%)
  - Географический boost (5%)
  - Копирование весов между категориями

- [ ] **Синонимы и языки**
  - CRUD для синонимов по языкам
  - Bulk import/export (CSV)
  - Статистика использования
  - Автоматические предложения

- [ ] **Boost правила**
  - Создание условных правил
  - Boost для конкретных listings
  - Boost по атрибутам
  - Временные акции
  - A/B тестирование правил

### 4. 🔍 Search Playground - Песочница (1 день)
- [ ] **Тестирование запросов**
  - Выполнение тестовых запросов
  - Детальное объяснение scoring
  - Визуализация весов
  - Сравнение конфигураций

- [ ] **Debugging tools**
  - Query analyzer
  - Объяснение почему результат на позиции N
  - Сравнение с/без fuzzy search
  - Проверка синонимов

### 5. 📈 Performance - Мониторинг производительности (1 день)
- [ ] **Метрики систем**
  - OpenSearch: размер индекса, cache hit rate
  - PostgreSQL: fuzzy search performance
  - API: latency по endpoints
  - Очереди и задержки

- [ ] **Алерты и рекомендации**
  - Автоматические алерты
  - Рекомендации по оптимизации
  - Предупреждения о проблемах
  - Health checks

### 6. 🎛️ Feature Management (1 день)
- [ ] **Feature flags**
  - Fuzzy search on/off
  - Синонимы on/off
  - Storefront integration
  - Персонализация
  - Geo-boost

- [ ] **A/B эксперименты**
  - Создание экспериментов
  - Распределение трафика
  - Метрики успеха
  - Статистическая значимость

### 7. 🔧 Maintenance - Обслуживание (1 день)
- [ ] **Управление индексами**
  - Ручная переиндексация
  - Расписание автоматической
  - Статус и прогресс
  - Rollback механизм

- [ ] **Очистка данных**
  - Старые логи запросов
  - Неиспользуемые синонимы
  - Оптимизация таблиц

- [ ] **Backup/Restore**
  - Сохранение конфигураций
  - История изменений
  - Откат на дату

## 🛠 Технические детали

### База данных - новые таблицы
```sql
-- Логи всех поисковых запросов
CREATE TABLE search_logs (
    id SERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    normalized_query TEXT,
    results_count INT NOT NULL,
    response_time_ms INT NOT NULL,
    filters JSONB,
    user_id INT,
    session_id VARCHAR(255),
    language VARCHAR(10),
    device_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Конфигурация весов по категориям
CREATE TABLE category_search_weights (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES marketplace_categories(id),
    attribute_weights JSONB NOT NULL, -- {"brand": 0.8, "model": 0.7, ...}
    text_weight FLOAT DEFAULT 0.2,
    price_weight FLOAT DEFAULT 0.15,
    location_weight FLOAT DEFAULT 0.05,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INT REFERENCES users(id)
);

-- Правила boost
CREATE TABLE search_boost_rules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    conditions JSONB NOT NULL, -- {"query": "iphone", "category": 123}
    boosts JSONB NOT NULL, -- {"listings": [1,2,3], "attributes": {"new": 1.5}}
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    valid_from TIMESTAMP,
    valid_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- История изменений конфигурации
CREATE TABLE search_config_history (
    id SERIAL PRIMARY KEY,
    config_type VARCHAR(50) NOT NULL, -- 'weights', 'synonyms', 'boost_rule'
    entity_id INT,
    old_value JSONB,
    new_value JSONB,
    changed_by INT REFERENCES users(id),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reason TEXT
);

-- A/B эксперименты
CREATE TABLE search_experiments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    variants JSONB NOT NULL, -- [{"name": "control", "config": {...}, "traffic": 50}]
    metrics JSONB, -- {"ctr": {...}, "conversion": {...}}
    status VARCHAR(50) DEFAULT 'draft', -- draft, running, completed
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    created_by INT REFERENCES users(id)
);
```

### API Endpoints - полный список
```typescript
// Dashboard
GET /api/admin/search/dashboard/stats
GET /api/admin/search/dashboard/realtime

// Analytics
GET /api/admin/search/analytics/queries?period=7d&limit=100
GET /api/admin/search/analytics/zero-results
GET /api/admin/search/analytics/performance
GET /api/admin/search/analytics/behavior
GET /api/admin/search/analytics/patterns

// Configuration
GET /api/admin/search/config/weights
PUT /api/admin/search/config/weights/:categoryId
POST /api/admin/search/config/weights/copy

// Synonyms
GET /api/admin/search/synonyms?language=ru
POST /api/admin/search/synonyms
PUT /api/admin/search/synonyms/:id
DELETE /api/admin/search/synonyms/:id
POST /api/admin/search/synonyms/import
GET /api/admin/search/synonyms/export

// Boost Rules
GET /api/admin/search/boost-rules
POST /api/admin/search/boost-rules
PUT /api/admin/search/boost-rules/:id
DELETE /api/admin/search/boost-rules/:id
POST /api/admin/search/boost-rules/:id/toggle

// Testing & Debug
POST /api/admin/search/test
POST /api/admin/search/explain
POST /api/admin/search/compare
GET /api/admin/search/debug/:query

// Performance
GET /api/admin/search/performance/metrics
GET /api/admin/search/performance/opensearch
GET /api/admin/search/performance/postgres
GET /api/admin/search/performance/recommendations

// Feature Flags
GET /api/admin/search/features
PUT /api/admin/search/features/:feature

// Experiments
GET /api/admin/search/experiments
POST /api/admin/search/experiments
PUT /api/admin/search/experiments/:id
POST /api/admin/search/experiments/:id/start
POST /api/admin/search/experiments/:id/stop

// Maintenance
POST /api/admin/search/reindex
GET /api/admin/search/reindex/status
POST /api/admin/search/cache/clear
GET /api/admin/search/health
POST /api/admin/search/backup
POST /api/admin/search/restore/:backupId
```

### UI компоненты - полная структура
```
/admin/search/
├── layout/
│   ├── SearchAdminLayout.tsx
│   ├── SearchAdminSidebar.tsx
│   └── SearchAdminHeader.tsx
├── dashboard/
│   ├── SearchDashboard.tsx
│   ├── RealtimeMetrics.tsx
│   ├── SearchVolumeChart.tsx
│   └── TopQueriesWidget.tsx
├── analytics/
│   ├── QueryAnalytics.tsx
│   ├── ZeroResultsAnalysis.tsx
│   ├── BehaviorAnalytics.tsx
│   ├── SearchPatterns.tsx
│   └── ClickHeatmap.tsx
├── configuration/
│   ├── WeightsManager.tsx
│   ├── CategoryWeightEditor.tsx
│   ├── SynonymsManager.tsx
│   ├── BoostRulesManager.tsx
│   └── ConfigHistory.tsx
├── playground/
│   ├── SearchPlayground.tsx
│   ├── QueryExplainer.tsx
│   ├── ConfigComparison.tsx
│   └── ScoringVisualizer.tsx
├── performance/
│   ├── PerformanceMonitor.tsx
│   ├── SystemMetrics.tsx
│   ├── AlertsManager.tsx
│   └── Recommendations.tsx
├── experiments/
│   ├── ExperimentsManager.tsx
│   ├── ExperimentCreator.tsx
│   ├── ExperimentResults.tsx
│   └── TrafficSplitter.tsx
├── maintenance/
│   ├── MaintenancePanel.tsx
│   ├── ReindexManager.tsx
│   ├── BackupRestore.tsx
│   └── DataCleanup.tsx
└── shared/
    ├── SearchAdminTable.tsx
    ├── MetricCard.tsx
    ├── ConfigEditor.tsx
    └── ConfirmDialog.tsx
```

## 📊 Что это даст

1. **Полная прозрачность**: Видим ВСЁ, что происходит с поиском
2. **Тотальный контроль**: Можем настроить любой аспект
3. **Data-driven решения**: Все изменения на основе метрик
4. **Безопасность**: История изменений и быстрый откат
5. **Масштабируемость**: Готовы к росту нагрузки

## 🔄 Интеграция с Фазой 3

Имея полную систему управления:
- **Видим паттерны**: Какие запросы требуют оптимизации
- **Тестируем гипотезы**: A/B тесты перед автоматизацией
- **Измеряем эффект**: Точные метрики до/после
- **Контролируем риски**: Можем вмешаться в любой момент

## ⏱️ Timeline (расширенный)

**Неделя 1:**
1. **День 1-2**: База данных + начальное логирование
2. **День 3-4**: Backend API (все endpoints)
3. **День 5**: Dashboard + базовая аналитика

**Неделя 2:**
1. **День 6-7**: Configuration management UI
2. **День 8**: Search Playground
3. **День 9**: Performance monitoring
4. **День 10**: Testing + документация

## ✅ Критерии готовности

- [ ] Логируется 100% поисковых запросов
- [ ] Dashboard показывает метрики в реальном времени
- [ ] Можно изменить любой параметр через UI
- [ ] Работает explain для любого запроса
- [ ] A/B тесты можно запускать без кода
- [ ] История всех изменений сохраняется
- [ ] Документация для всех компонентов