# Паспорт сервиса: UnifiedSearchService

## 📋 Метаданные
- **Название**: UnifiedSearchService
- **Путь**: `frontend/svetu/src/services/search/unifiedSearchService.ts`
- **Роль**: Сервис унифицированного поиска
- **Уровень**: Сервисный слой

## 🎯 Назначение
Централизованный сервис для работы с поисковой функциональностью, объединяющий поиск по marketplace и storefront товарам через OpenSearch backend.

## 🔧 Технические детали

### Интерфейсы
```typescript
interface SearchRequest {
  query: string;
  filters?: {
    category?: string;
    minPrice?: number;
    maxPrice?: number;
    type?: 'all' | 'marketplace' | 'storefront';
    location?: {
      lat: number;
      lng: number;
      radius: number;
    };
  };
  sort?: 'relevance' | 'price_asc' | 'price_desc' | 'date';
  page?: number;
  limit?: number;
  locale?: string;
}

interface SearchResponse {
  results: SearchResult[];
  facets: SearchFacets;
  total: number;
  page: number;
  hasMore: boolean;
  query: {
    original: string;
    corrected?: string;
  };
}

interface SuggestionsRequest {
  query: string;
  limit?: number;
  types?: Array<'product' | 'category' | 'query'>;
}

interface SuggestionsResponse {
  suggestions: SearchSuggestion[];
  trending?: string[];
}
```

### Основные методы

#### 1. Поиск товаров
```typescript
async search(request: SearchRequest): Promise<SearchResponse> {
  const params = new URLSearchParams({
    q: request.query,
    page: (request.page || 1).toString(),
    limit: (request.limit || 20).toString(),
    sort: request.sort || 'relevance',
    locale: request.locale || 'ru',
  });
  
  // Добавляем фильтры
  if (request.filters) {
    Object.entries(request.filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (key === 'location' && typeof value === 'object') {
          params.append('lat', value.lat.toString());
          params.append('lng', value.lng.toString());
          params.append('radius', value.radius.toString());
        } else {
          params.append(key, value.toString());
        }
      }
    });
  }
  
  const response = await fetch(`/api/search?${params}`, {
    headers: {
      'Accept': 'application/json',
      'X-Request-ID': generateRequestId(),
    },
  });
  
  if (!response.ok) {
    throw new SearchError('Search request failed', response.status);
  }
  
  const data = await response.json();
  return this.transformSearchResponse(data);
}
```

#### 2. Получение предложений
```typescript
async getSuggestions(request: SuggestionsRequest): Promise<SuggestionsResponse> {
  // Проверяем кэш
  const cacheKey = `suggestions:${request.query}`;
  const cached = this.cache.get(cacheKey);
  if (cached) {
    return cached;
  }
  
  const params = new URLSearchParams({
    q: request.query,
    limit: (request.limit || 10).toString(),
  });
  
  if (request.types) {
    params.append('types', request.types.join(','));
  }
  
  const response = await fetch(`/api/search/suggestions?${params}`, {
    signal: AbortSignal.timeout(3000), // 3 секунды таймаут
  });
  
  if (!response.ok) {
    console.error('Suggestions request failed');
    return { suggestions: [] };
  }
  
  const data = await response.json();
  const result = this.transformSuggestionsResponse(data);
  
  // Кэшируем на 5 минут
  this.cache.set(cacheKey, result, 5 * 60 * 1000);
  
  return result;
}
```

#### 3. Поиск по координатам
```typescript
async searchByLocation(
  lat: number, 
  lng: number, 
  radius: number = 5000
): Promise<SearchResponse> {
  return this.search({
    query: '',
    filters: {
      type: 'storefront',
      location: { lat, lng, radius },
    },
    sort: 'distance',
  });
}
```

#### 4. Популярные запросы
```typescript
async getTrendingSearches(limit: number = 10): Promise<string[]> {
  const cacheKey = 'trending';
  const cached = this.cache.get(cacheKey);
  if (cached) {
    return cached;
  }
  
  try {
    const response = await fetch(`/api/search/trending?limit=${limit}`);
    const data = await response.json();
    
    // Кэшируем на 1 час
    this.cache.set(cacheKey, data.queries, 60 * 60 * 1000);
    
    return data.queries;
  } catch (error) {
    console.error('Failed to fetch trending searches');
    return [];
  }
}
```

### Вспомогательные методы

#### Трансформация данных
```typescript
private transformSearchResponse(data: any): SearchResponse {
  return {
    results: data.hits.map(this.transformSearchResult),
    facets: this.transformFacets(data.aggregations),
    total: data.total.value,
    page: data.page,
    hasMore: data.page * data.limit < data.total.value,
    query: {
      original: data.query,
      corrected: data.corrected_query,
    },
  };
}

private transformSearchResult(hit: any): SearchResult {
  const source = hit._source;
  return {
    id: hit._id,
    type: source.type,
    title: this.highlightText(source.title, hit.highlight?.title),
    description: this.highlightText(
      source.description, 
      hit.highlight?.description
    ),
    price: source.price,
    currency: source.currency || 'RSD',
    imageUrl: source.images?.[0]?.url,
    category: source.category_name,
    location: source.location?.city,
    rating: source.rating,
    reviewsCount: source.reviews_count,
    createdAt: source.created_at,
    seller: {
      id: source.seller_id,
      name: source.seller_name,
      avatarUrl: source.seller_avatar,
    },
  };
}
```

#### Кэширование
```typescript
class SimpleCache {
  private cache = new Map<string, { data: any; expires: number }>();
  
  get(key: string): any | null {
    const item = this.cache.get(key);
    if (!item) return null;
    
    if (Date.now() > item.expires) {
      this.cache.delete(key);
      return null;
    }
    
    return item.data;
  }
  
  set(key: string, data: any, ttl: number): void {
    this.cache.set(key, {
      data,
      expires: Date.now() + ttl,
    });
  }
  
  clear(): void {
    this.cache.clear();
  }
}
```

## 🔗 Зависимости

### Внешние
- `fetch`: HTTP запросы
- `AbortController`: Отмена запросов

### Внутренние
- Нет прямых зависимостей от других сервисов

## 📊 Управление состоянием
- **Memory Cache**: Кэширование предложений и трендов
- **Request Deduplication**: Предотвращение дублирующих запросов
- **Error Handling**: Graceful degradation при ошибках

## ⚡ Оптимизации
1. **Кэширование**: In-memory cache для частых запросов
2. **Request Batching**: Группировка похожих запросов
3. **Abort Signals**: Отмена устаревших запросов
4. **Compression**: Gzip для больших ответов

## 🎯 Примеры использования

### Базовый поиск
```typescript
const results = await unifiedSearchService.search({
  query: 'велосипед',
  page: 1,
  limit: 20,
});
```

### Поиск с фильтрами
```typescript
const results = await unifiedSearchService.search({
  query: 'ноутбук',
  filters: {
    category: 'electronics',
    minPrice: 500,
    maxPrice: 1500,
    type: 'marketplace',
  },
  sort: 'price_asc',
});
```

### Получение предложений
```typescript
const suggestions = await unifiedSearchService.getSuggestions({
  query: 'вел',
  limit: 5,
  types: ['product', 'category'],
});
```

### Поиск по геолокации
```typescript
const nearbyStores = await unifiedSearchService.searchByLocation(
  44.786568,  // lat
  20.448922,  // lng
  3000        // radius in meters
);
```

## 🐛 Известные проблемы
1. **TODO**: Реализовать offline mode
2. **TODO**: Добавить метрики производительности
3. **Missing**: Персонализация результатов
4. **Hardcoded**: Таймауты и лимиты

## 🔒 Безопасность
- Санитизация входных данных
- Request ID для трейсинга
- Rate limiting headers
- CORS политики

## 📈 Метрики
- Время ответа поиска
- Cache hit rate
- Популярные запросы
- Конверсия поиска в покупку