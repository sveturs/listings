# Задание на изучение системы категорий и поиска

## Цель изучения
Понять текущую архитектуру категорий, поиска и индексации для разработки системы умного определения категорий на основе семантического анализа.

## 1. Изучение структуры БД

### 1.1 Таблицы категорий
- [ ] Изучить структуру `marketplace_categories`
  - Поля и их назначение
  - Иерархическая структура (parent_id)
  - Связи с другими таблицами
  
- [ ] Изучить `category_attributes`
  - Типы атрибутов
  - Связь с категориями
  - Опциональные/обязательные атрибуты

- [ ] Изучить `translations`
  - Структура мультиязычных переводов
  - entity_type для категорий
  - Поддерживаемые языки

### 1.2 Поисковые индексы
- [ ] Проанализировать структуру индексов в PostgreSQL
- [ ] Найти существующие полнотекстовые индексы

## 2. Изучение OpenSearch

### 2.1 Структура индексов
- [ ] Изучить mapping для `marketplace` индекса
  - Какие поля индексируются
  - Типы полей и анализаторы
  - Вложенные объекты

### 2.2 Поисковые запросы
- [ ] Проанализировать текущие поисковые запросы в:
  - `backend/internal/proj/marketplace/storage/opensearch/`
  - Особенно файлы поиска и фильтрации

### 2.3 Индексация категорий
- [ ] Как категории попадают в OpenSearch
- [ ] Какие данные о категориях хранятся
- [ ] Связь категорий с листингами

## 3. Backend анализ

### 3.1 API endpoints
- [ ] GET /api/v1/marketplace/categories
- [ ] GET /api/v1/marketplace/categories/:id/attributes
- [ ] POST /api/v1/marketplace/search

### 3.2 Сервисы и хендлеры
- [ ] marketplace/handler/categories.go
- [ ] marketplace/service/categories.go
- [ ] marketplace/storage/postgres/categories.go

### 3.3 Поисковая логика
- [ ] Изучить `SearchService` 
- [ ] Анализ `UnifiedSearchHandler`
- [ ] Веса и скоринг в поиске

## 4. Frontend анализ

### 4.1 Компоненты категорий
- [ ] CategorySelector компонент
- [ ] Как загружаются категории
- [ ] Кэширование на клиенте

### 4.2 AI интеграция
- [ ] create-listing-ai/page.tsx
- [ ] Функция getCategoryData
- [ ] Маппинг AI результатов на категории

### 4.3 Поисковые компоненты
- [ ] Как работает поиск по категориям
- [ ] Фильтрация по атрибутам

## Вопросы для исследования:

1. **Есть ли уже механизмы весов в системе?**
2. **Как хранится статистика использования категорий?**
3. **Есть ли аналитика по поисковым запросам?**
4. **Поддерживает ли OpenSearch fuzzy matching?**
5. **Есть ли Redis кэширование для категорий?**

## Ожидаемые результаты изучения:

1. Полная схема данных категорий
2. Понимание поисковых алгоритмов
3. Список точек интеграции для новой функциональности
4. Оценка сложности реализации
5. Возможные риски и ограничения

## Файлы для изучения:

```bash
# Backend
backend/internal/domain/models/category.go
backend/internal/proj/marketplace/handler/
backend/internal/proj/marketplace/service/
backend/internal/proj/marketplace/storage/
backend/internal/storage/opensearch/

# Frontend  
frontend/svetu/src/services/api/categories.ts
frontend/svetu/src/components/marketplace/CategorySelector.tsx
frontend/svetu/src/app/[locale]/create-listing-ai/

# Migrations
backend/migrations/*categories*.sql
backend/migrations/*attributes*.sql
```

---

После изучения создать документ с результатами в:
`/memory-bank/analysis/category-search-system-analysis.md`