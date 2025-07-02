# Паспорт компонента: SearchPage

## 📋 Метаданные
- **Название**: SearchPage
- **Путь**: `frontend/svetu/src/app/[locale]/search/SearchPage.tsx`
- **Роль**: Страница результатов поиска
- **Уровень**: Страница приложения

## 🎯 Назначение
Страница отображения результатов поиска с фильтрацией, сортировкой и бесконечной прокруткой. Объединяет результаты из marketplace и storefront товаров.

## 🔧 Технические детали

### Props
```typescript
interface SearchPageProps {
  searchParams: {
    q?: string;           // Поисковый запрос
    category?: string;    // Фильтр по категории
    minPrice?: string;    // Минимальная цена
    maxPrice?: string;    // Максимальная цена
    sort?: 'relevance' | 'price_asc' | 'price_desc' | 'date';
    type?: 'all' | 'marketplace' | 'storefront';
  };
  locale: string;
}
```

### Состояние
```typescript
{
  results: SearchResult[];      // Результаты поиска
  isLoading: boolean;          // Статус загрузки
  hasMore: boolean;            // Есть ли еще результаты
  page: number;                // Текущая страница
  totalResults: number;        // Общее количество
  facets: SearchFacets;        // Фасеты для фильтров
  error: string | null;        // Ошибка поиска
}

interface SearchResult {
  id: string;
  type: 'marketplace' | 'storefront';
  title: string;
  description: string;
  price: number;
  currency: string;
  imageUrl?: string;
  category: string;
  location?: string;
  rating?: number;
  reviewsCount?: number;
  createdAt: string;
  seller: {
    id: number;
    name: string;
    avatarUrl?: string;
  };
}

interface SearchFacets {
  categories: Array<{
    id: number;
    name: string;
    count: number;
  }>;
  priceRange: {
    min: number;
    max: number;
  };
  types: Array<{
    type: string;
    count: number;
  }>;
}
```

### Основная логика

#### 1. Загрузка результатов
```typescript
const loadResults = async (append = false) => {
  setIsLoading(true);
  setError(null);
  
  try {
    const response = await unifiedSearchService.search({
      query: searchParams.q || '',
      filters: {
        category: searchParams.category,
        minPrice: searchParams.minPrice ? parseFloat(searchParams.minPrice) : undefined,
        maxPrice: searchParams.maxPrice ? parseFloat(searchParams.maxPrice) : undefined,
        type: searchParams.type,
      },
      sort: searchParams.sort || 'relevance',
      page: append ? page + 1 : 1,
      limit: 20,
    });
    
    if (append) {
      setResults(prev => [...prev, ...response.results]);
    } else {
      setResults(response.results);
    }
    
    setFacets(response.facets);
    setTotalResults(response.total);
    setHasMore(response.hasMore);
    setPage(response.page);
  } catch (error) {
    setError(t('search.error.failed'));
    console.error('Search failed:', error);
  } finally {
    setIsLoading(false);
  }
};
```

#### 2. Бесконечная прокрутка
```typescript
const { ref, inView } = useIntersectionObserver({
  threshold: 0.1,
  rootMargin: '100px',
});

useEffect(() => {
  if (inView && hasMore && !isLoading) {
    loadResults(true);
  }
}, [inView, hasMore, isLoading]);
```

#### 3. Обработка фильтров
```typescript
const handleFilterChange = (filters: Partial<typeof searchParams>) => {
  const newParams = new URLSearchParams();
  
  // Сохраняем поисковый запрос
  if (searchParams.q) {
    newParams.set('q', searchParams.q);
  }
  
  // Применяем новые фильтры
  Object.entries(filters).forEach(([key, value]) => {
    if (value) {
      newParams.set(key, value.toString());
    }
  });
  
  router.push(`/${locale}/search?${newParams.toString()}`);
};
```

## 🎨 UI структура
```jsx
<div className="container mx-auto px-4 py-6">
  {/* Заголовок с количеством результатов */}
  <div className="mb-6">
    <h1 className="text-2xl font-bold">
      {searchParams.q ? (
        <>
          {t('search.resultsFor')} "{searchParams.q}"
          <span className="text-base-content/60 ml-2">
            ({totalResults} {t('search.found')})
          </span>
        </>
      ) : (
        t('search.allProducts')
      )}
    </h1>
  </div>
  
  <div className="flex gap-6">
    {/* Боковая панель фильтров */}
    <aside className="w-64 hidden lg:block">
      <SearchFilters
        facets={facets}
        currentFilters={searchParams}
        onChange={handleFilterChange}
      />
    </aside>
    
    {/* Основной контент */}
    <main className="flex-1">
      {/* Панель сортировки */}
      <div className="flex justify-between items-center mb-4">
        <button
          className="lg:hidden btn btn-sm"
          onClick={() => setShowMobileFilters(true)}
        >
          <FilterIcon className="w-4 h-4 mr-1" />
          {t('search.filters')}
        </button>
        
        <SearchSort
          value={searchParams.sort || 'relevance'}
          onChange={(sort) => handleFilterChange({ sort })}
        />
      </div>
      
      {/* Результаты поиска */}
      {error ? (
        <div className="alert alert-error">{error}</div>
      ) : results.length === 0 && !isLoading ? (
        <SearchEmpty query={searchParams.q} />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {results.map((result) => (
            <SearchResultCard
              key={`${result.type}-${result.id}`}
              result={result}
              locale={locale}
            />
          ))}
        </div>
      )}
      
      {/* Индикатор загрузки */}
      {isLoading && (
        <div className="flex justify-center py-8">
          <span className="loading loading-spinner loading-lg" />
        </div>
      )}
      
      {/* Триггер бесконечной прокрутки */}
      {hasMore && !isLoading && (
        <div ref={ref} className="h-20" />
      )}
    </main>
  </div>
  
  {/* Мобильные фильтры */}
  <MobileFiltersDrawer
    isOpen={showMobileFilters}
    onClose={() => setShowMobileFilters(false)}
    facets={facets}
    currentFilters={searchParams}
    onChange={handleFilterChange}
  />
</div>
```

## 🔗 Зависимости

### Внешние библиотеки
- `react`: Хуки и компоненты
- `next/navigation`: Роутинг и параметры
- `react-intersection-observer`: Бесконечная прокрутка
- `next-intl`: Локализация

### Внутренние компоненты
- `SearchFilters`: Панель фильтров
- `SearchSort`: Выбор сортировки
- `SearchResultCard`: Карточка результата
- `SearchEmpty`: Пустые результаты
- `MobileFiltersDrawer`: Мобильные фильтры

### API интеграция
```typescript
// Основной поиск
unifiedSearchService.search({
  query: string,
  filters: SearchFilters,
  sort: SortOption,
  page: number,
  limit: number
})
// Response: {
//   results: SearchResult[],
//   facets: SearchFacets,
//   total: number,
//   hasMore: boolean,
//   page: number
// }
```

## 📊 Управление состоянием
- **URL State**: Все параметры поиска в query string
- **Component State**: Результаты и состояние загрузки
- **Server State**: Кэширование на уровне API

## 🔄 Жизненный цикл
1. **Mount**: Загрузка результатов из URL параметров
2. **URL Change**: Перезагрузка с новыми параметрами
3. **Scroll**: Подгрузка следующей страницы
4. **Filter Change**: Сброс результатов и новый поиск

## 📱 Адаптивность
- **Mobile**: Скрытые фильтры в drawer, одна колонка
- **Tablet**: Две колонки результатов
- **Desktop**: Три колонки + боковая панель

## ⚡ Оптимизации
1. **Infinite Scroll**: Постепенная загрузка результатов
2. **URL State**: Сохранение состояния при навигации
3. **Debounced Filters**: Задержка применения фильтров
4. **Image Lazy Loading**: Отложенная загрузка изображений

## 🎯 Примеры использования

### Прямой поиск
```
/ru/search?q=велосипед
```

### С фильтрами
```
/ru/search?q=велосипед&category=sport&minPrice=100&maxPrice=500&sort=price_asc
```

### По типу товаров
```
/ru/search?type=storefront&category=electronics
```

## 🐛 Известные проблемы
1. **TODO**: Сохранение позиции скролла при навигации
2. **TODO**: Предзагрузка следующей страницы
3. **Missing**: История просмотров
4. **Hardcoded**: Количество результатов на страницу (20)

## 🔒 Безопасность
- Санитизация поискового запроса
- Валидация числовых параметров
- Защита от SQL инъекций на backend
- Rate limiting для API запросов

## 🌐 SEO оптимизация
- Мета-теги с поисковым запросом
- Структурированные данные для товаров
- Canonical URL для дубликатов
- Sitemap для популярных запросов