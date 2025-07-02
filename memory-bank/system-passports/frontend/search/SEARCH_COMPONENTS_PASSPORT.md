# Паспорта компонентов поиска системы Sve Tu Platform

## 1. SearchBar Component

### Название и назначение
**Компонент**: `SearchBar`  
**Путь**: `/frontend/svetu/src/components/SearchBar/SearchBar.tsx`  
**Назначение**: Универсальная строка поиска с автодополнением, историей поиска и поддержкой трендовых запросов

### Основные props и их типы
```typescript
interface SearchBarProps {
  placeholder?: string;           // Текст плейсхолдера
  initialQuery?: string;         // Начальное значение поиска
  onSearch?: (query: string) => void; // Callback при поиске
  className?: string;            // Дополнительные CSS классы
  variant?: 'default' | 'hero' | 'minimal'; // Варианты отображения
  showTrending?: boolean;        // Показывать ли трендовые запросы
}
```

### Зависимости (импорты)
```typescript
// React и хуки
import { useState, useEffect, useRef, useCallback } from 'react';

// Интернационализация и навигация
import { useTranslations, useLocale } from 'next-intl';
import { useRouter } from 'next/navigation';

// Кастомные хуки
import { useDebounce } from '@/hooks/useDebounce';

// Сервисы и типы
import { UnifiedSearchService, SearchSuggestion } from '@/services/unifiedSearch';

// Компоненты иконок
import { SearchIcon, CloseIcon, TrendingIcon } from './icons';
```

### Ключевая логика работы

1. **Состояния компонента**:
   - `query` - текущее значение поиска
   - `suggestions` - массив предложений автодополнения
   - `isLoadingSuggestions` - флаг загрузки предложений
   - `showSuggestions` - флаг отображения выпадающего списка
   - `selectedIndex` - индекс выбранного элемента (навигация клавиатурой)
   - `searchHistory` - история поисковых запросов
   - `trendingSearches` - популярные поисковые запросы

2. **Debouncing**: Использует хук `useDebounce` с задержкой 300мс для оптимизации запросов

3. **Автодополнение**:
   - Загружает предложения при вводе минимум 2 символов
   - Показывает историю и тренды при пустом поле
   - Обрабатывает клики вне компонента для закрытия

4. **Навигация клавиатурой**:
   - `ArrowUp/ArrowDown` - перемещение по списку
   - `Enter` - выбор элемента или поиск
   - `Escape` - закрытие списка

5. **Типы предложений**:
   - `text` - обычный текстовый поиск
   - `category` - переход в категорию
   - `product` - переход к конкретному товару

### UI структура

```
<div className="relative">
  {/* Поле ввода */}
  <input
    type="text"
    placeholder={placeholder || t('placeholder')}
    className={getInputClasses()}
    value={query}
    onChange={(e) => setQuery(e.target.value)}
    onKeyDown={handleKeyDown}
  />
  
  {/* Кнопки очистки и поиска */}
  <div className="absolute right-3 top-1/2">
    {query && <button onClick={handleClear}><CloseIcon /></button>}
    <button onClick={() => handleSearch()}>
      {isLoadingSuggestions ? <span className="loading" /> : <SearchIcon />}
    </button>
  </div>
  
  {/* Выпадающий список предложений */}
  {showSuggestions && (
    <div className="absolute top-full">
      {/* Результаты поиска */}
      {suggestions.map((suggestion) => (
        <div onClick={() => handleSuggestionClick(suggestion)}>
          {renderSuggestionIcon(suggestion.type)}
          <div>{suggestion.text}</div>
          {suggestion.type === 'category' && <badge>{t('category')}</badge>}
        </div>
      ))}
      
      {/* История поиска */}
      {searchHistory.map((historyItem) => (
        <div onClick={() => handleSearch(historyItem)}>
          <svg>{/* История иконка */}</svg>
          <span>{historyItem}</span>
        </div>
      ))}
      
      {/* Популярные запросы */}
      {trendingSearches.map((trending) => (
        <div onClick={() => handleSearch(trending)}>
          <TrendingIcon />
          <span>{trending}</span>
        </div>
      ))}
    </div>
  )}
</div>
```

### Взаимодействие с API/Redux
- **API**: Использует `UnifiedSearchService` для:
  - `getSuggestions()` - получение предложений автодополнения
  - `getSearchHistory()` - загрузка истории из localStorage
  - `getTrendingSearches()` - получение популярных запросов
  - `saveToHistory()` - сохранение в историю
  - `clearHistory()` - очистка истории
- **Redux**: Не используется напрямую
- **Навигация**: Использует Next.js Router для перехода на страницу поиска

### Особенности реализации
1. **Три варианта отображения**:
   - `default` - стандартный размер
   - `hero` - увеличенный с тенью для главной страницы
   - `minimal` - компактный для header

2. **Оптимизация производительности**:
   - Debouncing для минимизации запросов
   - useCallback для handleSearch
   - Локальное хранение истории в localStorage

3. **Адаптивность**:
   - Responsive классы для разных экранов
   - Адаптивные размеры кнопок и input

4. **Доступность**:
   - Полная поддержка навигации клавиатурой
   - ARIA атрибуты для screen readers
   - Визуальные индикаторы загрузки

---

## 2. SearchBar Icons

### Название и назначение
**Компоненты**: `SearchIcon`, `CloseIcon`, `TrendingIcon`  
**Путь**: `/frontend/svetu/src/components/SearchBar/icons.tsx`  
**Назначение**: SVG иконки для компонента поиска

### Props
```typescript
interface IconProps {
  className?: string; // CSS классы (по умолчанию 'w-5 h-5')
}
```

### Компоненты
1. **SearchIcon** - иконка лупы для кнопки поиска
2. **CloseIcon** - крестик для очистки поля
3. **TrendingIcon** - график тренда для популярных запросов

---

## 3. SearchPage Component

### Название и назначение
**Компонент**: `SearchPage`  
**Путь**: `/frontend/svetu/src/app/[locale]/search/SearchPage.tsx`  
**Назначение**: Страница результатов поиска с фильтрами, сортировкой и бесконечной прокруткой

### Состояния и типы
```typescript
interface SearchFilters {
  category_id?: string;
  price_min?: number;
  price_max?: number;
  product_types?: string[];
  sort_by?: string;
  sort_order?: string;
  city?: string;
}

// Основные состояния
const [query, setQuery] = useState(initialQuery);
const [results, setResults] = useState<UnifiedSearchResult | null>(null);
const [loading, setLoading] = useState(false);
const [error, setError] = useState<string | null>(null);
const [filters, setFilters] = useState<SearchFilters>({...});
const [page, setPage] = useState(1);
const [allItems, setAllItems] = useState<any[]>([]);
const [showFilters, setShowFilters] = useState(false);
const [viewMode, setViewMode] = useViewPreference('grid');
```

### Зависимости
```typescript
// React хуки
import { useState, useEffect } from 'react';

// Интернационализация и навигация
import { useTranslations, useLocale } from 'next-intl';
import { useSearchParams, useRouter } from 'next/navigation';

// Компоненты
import { SearchBar } from '@/components/SearchBar';
import MarketplaceCard from '@/components/MarketplaceCard';
import ViewToggle from '@/components/common/ViewToggle';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';

// Хуки
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';

// Сервисы и типы
import { UnifiedSearchService, UnifiedSearchResult, UnifiedSearchParams } from '@/services/unifiedSearch';
import { MarketplaceItem } from '@/types/marketplace';
```

### Ключевая логика работы

1. **Управление поиском**:
   - Инициализация из URL параметра `q`
   - Автоматический поиск при изменении URL
   - Обновление URL при новом поиске
   - Пагинация с бесконечной прокруткой

2. **Фильтрация и сортировка**:
   - Типы товаров (marketplace/storefront)
   - Диапазон цен
   - Сортировка по релевантности/цене/дате/популярности
   - Город (для локального поиска)

3. **Обработка результатов**:
   - Конвертация данных в `MarketplaceItem` формат
   - Накопление результатов при пагинации
   - Fallback на пустые результаты при ошибке

4. **UI состояния**:
   - Загрузка
   - Ошибки (скрытые от пользователя)
   - Пустые результаты с рекомендациями
   - Статистика поиска

### UI структура

```
<div className="min-h-screen">
  {/* Sticky header с поиском */}
  <div className="sticky top-0 z-50">
    <SearchBar variant="minimal" />
    
    {/* Быстрые фильтры */}
    <div className="flex gap-2">
      <button>Объявления</button>
      <button>Магазины</button>
      <button>По цене</button>
      <button>По дате</button>
    </div>
  </div>
  
  {/* Основной контент */}
  <div className="container">
    {/* Статистика поиска - 4 карточки */}
    <div className="grid grid-cols-1 md:grid-cols-4">
      <div className="stat">Найдено: {results.total}</div>
      <div className="stat">Скорость: {results.took_ms}мс</div>
      <div className="stat">Время: {results.took_ms}мс</div>
      <div className="stat">Фильтры: {activeFiltersCount()}</div>
    </div>
    
    <div className="flex">
      {/* Боковая панель фильтров */}
      <aside className="lg:w-80">
        <div className="card">
          {/* Тип товаров - карточки выбора */}
          {/* Диапазон цен - input группы */}
          {/* Сортировка - select */}
        </div>
      </aside>
      
      {/* Результаты */}
      <main className="flex-1">
        {loading && <span className="loading" />}
        
        {allItems.length === 0 && (
          <div className="card">
            {/* Красивое сообщение "нет результатов" */}
            {/* Рекомендации категорий */}
            {/* Кнопки действий */}
          </div>
        )}
        
        {allItems.length > 0 && (
          <>
            <ViewToggle />
            <div className={viewMode === 'grid' ? 'grid' : 'space-y-4'}>
              {allItems.map(item => <MarketplaceCard />)}
            </div>
            <InfiniteScrollTrigger />
          </>
        )}
      </main>
    </div>
  </div>
</div>
```

### Взаимодействие с API/Redux
- **API**: 
  - `UnifiedSearchService.search()` - основной поиск
  - Поддержка всех параметров фильтрации
  - Обработка пагинации
- **Redux**: Не используется
- **LocalStorage**: Через `useViewPreference` для сохранения режима отображения

### Особенности реализации

1. **Оптимизация производительности**:
   - Бесконечная прокрутка вместо пагинации
   - Debounce в SearchBar компоненте
   - Накопление результатов без перерисовки

2. **UX улучшения**:
   - Сохранение режима отображения
   - Быстрые фильтры в header
   - Статистика поиска в реальном времени
   - Graceful error handling (без показа ошибок)

3. **Адаптивность**:
   - Скрытие фильтров на мобильных
   - Responsive grid для результатов
   - Адаптивные размеры кнопок

4. **Интеграция с OpenSearch**:
   - Полнотекстовый поиск
   - Фасетный поиск
   - Подсветка найденных фрагментов
   - Релевантность результатов

---

## 4. Search Page Wrapper

### Название и назначение
**Компонент**: `Page`  
**Путь**: `/frontend/svetu/src/app/[locale]/search/page.tsx`  
**Назначение**: Next.js страница-обертка для SearchPage с Suspense

### Особенности
- Использует `Suspense` для отложенной загрузки
- Показывает loading spinner во время загрузки
- Обеспечивает SSR совместимость

---

## 5. UnifiedSearchService

### Название и назначение
**Сервис**: `UnifiedSearchService`  
**Путь**: `/frontend/svetu/src/services/unifiedSearch.ts`  
**Назначение**: Сервис для работы с унифицированным поиском по всем типам товаров

### Основные методы

1. **search(params: UnifiedSearchParams)**
   - Выполняет поиск по всем типам товаров
   - Поддерживает фильтрацию, сортировку, пагинацию
   - Возвращает `UnifiedSearchResult`

2. **getSuggestions(prefix: string, size: number)**
   - Получает предложения для автодополнения
   - Endpoint: `/api/v1/marketplace/enhanced-suggestions`
   - Возвращает массив `SearchSuggestion`

3. **getTrendingSearches()**
   - Возвращает популярные запросы
   - TODO: Требуется реализация backend endpoint

4. **getSearchHistory() / saveToHistory() / clearHistory()**
   - Работа с историей поиска в localStorage
   - Хранит до 10 последних запросов

### Типы данных

```typescript
interface UnifiedSearchParams {
  query: string;
  product_types?: ('marketplace' | 'storefront')[];
  page?: number;
  limit?: number;
  category_id?: string;
  price_min?: number;
  price_max?: number;
  sort_by?: 'relevance' | 'price' | 'date' | 'popularity';
  sort_order?: 'asc' | 'desc';
  storefront_id?: number;
  city?: string;
  language?: string;
}

interface UnifiedSearchItem {
  id: string;
  product_type: 'marketplace' | 'storefront';
  product_id: number;
  name: string;
  description: string;
  price: number;
  images: Array<{...}>;
  category: {...};
  location?: {...};
  storefront?: {...};
  score: number;
  highlights?: Record<string, string[]>;
}

interface SearchSuggestion {
  text: string;
  type: 'text' | 'category' | 'product';
  category?: {...};
  product_id?: number;
}
```

### Особенности интеграции
- Использует `configManager` для получения API URL
- Поддерживает development proxy
- Обрабатывает ошибки gracefully
- Работает с localStorage для истории

## Общая архитектура поиска

### Поток данных
1. **Ввод пользователя** → SearchBar → debounce → UnifiedSearchService.getSuggestions()
2. **Выбор предложения** → Навигация на /search?q=... или прямой переход
3. **Страница поиска** → UnifiedSearchService.search() → Отображение результатов
4. **Фильтрация** → Обновление параметров → Новый запрос
5. **Infinite scroll** → Загрузка следующей страницы → Добавление к результатам

### Интеграция с OpenSearch
- Backend использует OpenSearch для полнотекстового поиска
- Поддержка фасетного поиска
- Scoring по релевантности
- Highlights для найденных фрагментов

### Ключевые особенности
1. **Унифицированный поиск** по marketplace и storefront товарам
2. **Умное автодополнение** с категориями и товарами
3. **История и тренды** для улучшения UX
4. **Фасетная фильтрация** для точного поиска
5. **Бесконечная прокрутка** для удобства просмотра
6. **Адаптивный дизайн** для всех устройств