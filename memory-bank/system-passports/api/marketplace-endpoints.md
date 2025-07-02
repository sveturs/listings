# Паспорт API Endpoints: Marketplace (Маркетплейс)

## 📋 Метаданные
- **Группа API**: Marketplace
- **Базовый путь**: `/api/v1/marketplace`
- **Handler**: `backend/internal/proj/marketplace/handler/handler.go`
- **Количество endpoints**: 41 (13 публичных, 28 защищенных)
- **Интеграции**: PostgreSQL, OpenSearch, MinIO

## 🎯 Назначение
Ядро платформы маркетплейса, включающее:
- CRUD операции с объявлениями
- Поиск и фильтрация через OpenSearch
- Система категорий с динамическими атрибутами
- Управление изображениями через MinIO
- Избранное и рекомендации
- Геопространственный поиск

## 📡 Endpoints

### 🌐 Публичные (читайте без авторизации)

#### Основные листинги
```typescript
GET /api/v1/marketplace/listings
// Получение списка объявлений с пагинацией и фильтрами
// Handler: h.Listings.GetListings
// Query: page, limit, category_id, price_min, price_max, location, status
// Response: ListingsResponse с пагинацией
```

```typescript
GET /api/v1/marketplace/listings/:id
// Получение детальной информации об объявлении
// Handler: h.Listings.GetListing
// Response: DetailedListing с изображениями и метаданными
```

#### Система категорий
```typescript
GET /api/v1/marketplace/categories
// Плоский список всех категорий
// Handler: h.Categories.GetCategories
// Response: Category[] с базовой информацией
```

```typescript
GET /api/v1/marketplace/category-tree  
// Иерархическое дерево категорий
// Handler: h.Categories.GetCategoryTree
// Response: CategoryTree с nested children
```

```typescript
GET /api/v1/marketplace/categories/:id/attributes
// Атрибуты конкретной категории для фильтрации
// Handler: h.Categories.GetCategoryAttributes
// Response: CategoryAttribute[] с типами и ограничениями
```

```typescript
GET /api/v1/marketplace/categories/:id/attribute-ranges
// Диапазоны значений атрибутов (min/max для числовых)
// Handler: h.Categories.GetAttributeRanges
// Response: AttributeRange[] для построения слайдеров
```

#### Поиск и автодополнение
```typescript
GET /api/v1/marketplace/search
// Расширенный поиск через OpenSearch
// Handler: h.Search.SearchListingsAdvanced
// Query: q, filters, sort, geo_filter
// Response: SearchResults с подсветкой и фасетами
```

```typescript
GET /api/v1/marketplace/suggestions
// Автодополнение для поисковой строки
// Handler: h.Search.GetSuggestions  
// Query: q, limit
// Response: SearchSuggestion[] с релевантностью
```

```typescript
GET /api/v1/marketplace/category-suggestions
// Предложения категорий для поискового запроса
// Handler: h.Search.GetCategorySuggestions
// Response: CategorySuggestion[] с весами
```

```typescript
GET /api/v1/marketplace/enhanced-suggestions
// Улучшенные предложения (товары + категории + тренды)
// Handler: h.Search.GetEnhancedSuggestions
// Response: EnhancedSuggestions с группировкой
```

#### Аналитика и рекомендации
```typescript
GET /api/v1/marketplace/listings/:id/price-history
// История изменения цены объявления
// Handler: h.Listings.GetPriceHistory
// Response: PriceHistoryPoint[] для графиков
```

```typescript
GET /api/v1/marketplace/listings/:id/similar
// Похожие объявления на основе ML
// Handler: h.Search.GetSimilarListings
// Response: SimilarListing[] с score релевантности
```

#### Геопространственный поиск
```typescript
GET /api/v1/marketplace/map/bounds
// Объявления в границах карты
// Handler: h.GetListingsInBounds
// Query: ne_lat, ne_lng, sw_lat, sw_lng, zoom
// Response: BoundedListings для отображения на карте
```

```typescript
GET /api/v1/marketplace/map/clusters
// Кластерные данные для карты (для производительности)
// Handler: h.GetMapClusters  
// Query: bounds, zoom_level
// Response: MapCluster[] с агрегированными данными
```

### 🔒 Защищенные (требуют авторизации)

#### Управление объявлениями
```typescript
POST /api/v1/marketplace/listings
// Создание нового объявления
// Handler: h.Listings.CreateListing
// Body: CreateListingRequest с полной информацией
// Response: CreatedListing с ID и статусом
```

```typescript
PUT /api/v1/marketplace/listings/:id
// Обновление существующего объявления
// Handler: h.Listings.UpdateListing
// Security: Только владелец или админ
// Body: UpdateListingRequest (partial)
```

```typescript
DELETE /api/v1/marketplace/listings/:id
// Удаление объявления (soft delete)
// Handler: h.Listings.DeleteListing
// Security: Только владелец или админ
// Effect: Статус меняется на 'deleted'
```

#### Управление изображениями
```typescript
POST /api/v1/marketplace/listings/:id/images
// Загрузка изображений к объявлению
// Handler: h.Images.UploadImages
// Content-Type: multipart/form-data
// Files: Поддержка до 10 изображений
// Integration: MinIO для хранения
```

#### Система избранного
```typescript
POST /api/v1/marketplace/listings/:id/favorite
// Добавление в избранное
// Handler: h.Favorites.AddToFavorites
// Effect: Создается запись в marketplace_favorites
```

```typescript
DELETE /api/v1/marketplace/listings/:id/favorite
// Удаление из избранного
// Handler: h.Favorites.RemoveFromFavorites
// Effect: Удаляется запись из marketplace_favorites
```

```typescript
GET /api/v1/marketplace/favorites
// Получение списка избранных объявлений
// Handler: h.Favorites.GetFavorites
// Response: FavoriteListing[] с метаданными
```

## 🏗️ Структуры данных

### Основные модели
```typescript
interface Listing {
  id: string;
  title: string;
  description: string;
  price: number;
  currency: "RSD" | "EUR";
  category_id: string;
  user_id: string;
  status: "draft" | "active" | "sold" | "deleted";
  location: {
    city: string;
    address?: string;
    coordinates?: [number, number];
  };
  attributes: Record<string, any>;
  images: ListingImage[];
  created_at: string;
  updated_at: string;
  expires_at: string;
}

interface ListingImage {
  id: string;
  url: string;
  thumbnail_url: string;
  order: number;
  alt_text?: string;
}

interface Category {
  id: string;
  name: Record<string, string>; // локализованные названия
  slug: string;
  parent_id?: string;
  icon?: string;
  attributes: CategoryAttribute[];
  children?: Category[];
}

interface CategoryAttribute {
  id: string;
  name: Record<string, string>;
  type: "string" | "number" | "boolean" | "select" | "multiselect";
  required: boolean;
  options?: string[];
  min_value?: number;
  max_value?: number;
  unit?: string;
}
```

### Поисковые структуры
```typescript
interface SearchFilters {
  category_ids?: string[];
  price_range?: [number, number];
  location?: {
    city?: string;
    radius?: number; // км
    coordinates?: [number, number];
  };
  attributes?: Record<string, any>;
  date_range?: [string, string];
  status?: ("active" | "sold")[];
}

interface SearchResults {
  listings: SearchListing[];
  total: number;
  facets: SearchFacets;
  suggestions?: string[];
  took: number; // время выполнения в мс
}

interface SearchFacets {
  categories: FacetBucket[];
  price_ranges: FacetBucket[];
  locations: FacetBucket[];
  attributes: Record<string, FacetBucket[]>;
}
```

## 🔄 Интеграции

### OpenSearch Integration
- **Index**: `marketplace_listings`
- **Mapping**: Поддержка мультиязычности (ru/en)
- **Analyzers**: Кастомные для кириллицы и латиницы
- **Geospatial**: Поддержка geo_point для карты

### MinIO Integration
- **Bucket**: `listings`
- **Path**: `/listings/{listing_id}/{image_id}.{ext}`
- **Thumbnails**: Автогенерация при загрузке
- **CDN**: Раздача через nginx

### Database Schema
```sql
-- Основная таблица объявлений
marketplace_listings (
  id, title, description, price, currency,
  category_id, user_id, status, location_data,
  attributes_json, created_at, updated_at
);

-- Изображения
marketplace_images (
  id, listing_id, file_path, thumbnail_path,
  order_index, alt_text
);

-- Избранное
marketplace_favorites (
  user_id, listing_id, created_at
);

-- История цен
price_history (
  listing_id, price, currency, created_at
);
```

## 🎛️ Бизнес-логика

### Статусы объявлений
- **draft**: Черновик, видим только автору
- **active**: Активное, индексируется и показывается
- **sold**: Продано, только для чтения
- **deleted**: Удалено, скрыто от пользователей

### Система цен
- Поддержка RSD и EUR
- Автоматическое отслеживание изменений в price_history
- Валидация разумных диапазонов по категориям

### Геолокация
- Безопасное хранение (только город по умолчанию)
- Точные координаты только для доверенных пользователей
- Радиусный поиск для карты

## 🔍 OpenSearch Queries

### Основной поиск
```json
{
  "query": {
    "bool": {
      "must": [
        {"multi_match": {"query": "текст", "fields": ["title^2", "description"]}}
      ],
      "filter": [
        {"term": {"status": "active"}},
        {"range": {"price": {"gte": 100, "lte": 5000}}},
        {"geo_distance": {"distance": "10km", "location": {"lat": 44.8, "lon": 20.5}}}
      ]
    }
  },
  "aggs": {
    "categories": {"terms": {"field": "category_id"}},
    "price_ranges": {"histogram": {"field": "price", "interval": 1000}}
  }
}
```

## ⚠️ Известные особенности

### Performance
- Пагинация: default 20, max 100 items
- OpenSearch timeout: 5 секунд
- Image upload: max 10MB per file, 10 files max
- Rate limiting: 100 requests/minute per user

### Security
- Файлы проверяются на MIME type
- Изображения сканируются на вирусы
- XSS защита в user-generated content
- CSRF защита для всех POST/PUT/DELETE

### Локализация
- Категории и атрибуты поддерживают ru/en
- Поиск работает с обеими раскладками
- Геолокация использует локальные справочники

### Caching
- Categories кешируются на 1 час
- Search suggestions на 15 минут
- Listing details на 5 минут для не-владельцев

## 🧪 Примеры использования

### Поиск с фильтрами
```bash
curl "/api/v1/marketplace/search?q=telefon&category_ids=electronics&price_max=50000&location=Belgrade"
```

### Создание объявления
```bash
curl -X POST /api/v1/marketplace/listings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "iPhone 15 Pro",
    "description": "Excellent condition",
    "price": 120000,
    "currency": "RSD",
    "category_id": "smartphones",
    "location": {"city": "Belgrade"}
  }'
```

### Загрузка изображений
```bash
curl -X POST /api/v1/marketplace/listings/123/images \
  -H "Authorization: Bearer <token>" \
  -F "images=@photo1.jpg" \
  -F "images=@photo2.jpg"
```