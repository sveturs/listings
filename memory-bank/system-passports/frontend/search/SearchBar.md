# Паспорт компонента: SearchBar

## 📋 Метаданные
- **Название**: SearchBar
- **Путь**: `frontend/svetu/src/components/SearchBar/SearchBar.tsx`
- **Роль**: Универсальная строка поиска с автодополнением
- **Уровень**: Компонент пользовательского интерфейса

## 🎯 Назначение
Универсальный компонент поиска с автодополнением, предложениями, историей и интеграцией с OpenSearch для поиска по товарам marketplace и storefront.

## 🔧 Технические детали

### Props
```typescript
interface SearchBarProps {
  className?: string;
  placeholder?: string;
  autoFocus?: boolean;
  onSearch?: (query: string) => void;
  showSuggestions?: boolean;
  maxSuggestions?: number;
  searchDebounceMs?: number;
}
```

### Состояние
```typescript
{
  query: string;                    // Текущий поисковый запрос
  suggestions: SearchSuggestion[];  // Предложения автодополнения
  isLoading: boolean;              // Статус загрузки
  isFocused: boolean;              // Фокус на поле ввода
  selectedIndex: number;           // Выбранное предложение (клавиатура)
  searchHistory: string[];         // История поиска из localStorage
}

interface SearchSuggestion {
  text: string;          // Текст предложения
  type: 'product' | 'category' | 'query';  // Тип предложения
  count?: number;        // Количество результатов
  categoryId?: number;   // ID категории
}
```

### Основная логика

#### 1. Автодополнение
```typescript
const fetchSuggestions = useDebounce(async (query: string) => {
  if (query.length < 2) {
    setSuggestions([]);
    return;
  }
  
  setIsLoading(true);
  try {
    const response = await unifiedSearchService.getSuggestions(query);
    setSuggestions(response.suggestions);
  } catch (error) {
    console.error('Failed to fetch suggestions:', error);
  } finally {
    setIsLoading(false);
  }
}, searchDebounceMs);
```

#### 2. История поиска
```typescript
const saveToHistory = (query: string) => {
  const history = getSearchHistory();
  const filtered = history.filter(item => item !== query);
  const updated = [query, ...filtered].slice(0, 10);
  localStorage.setItem('searchHistory', JSON.stringify(updated));
};

const getSearchHistory = (): string[] => {
  try {
    return JSON.parse(localStorage.getItem('searchHistory') || '[]');
  } catch {
    return [];
  }
};
```

#### 3. Навигация клавиатурой
```typescript
const handleKeyDown = (e: KeyboardEvent) => {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault();
      setSelectedIndex(prev => 
        Math.min(prev + 1, suggestions.length - 1)
      );
      break;
    case 'ArrowUp':
      e.preventDefault();
      setSelectedIndex(prev => Math.max(prev - 1, -1));
      break;
    case 'Enter':
      e.preventDefault();
      if (selectedIndex >= 0) {
        handleSelectSuggestion(suggestions[selectedIndex]);
      } else {
        handleSearch();
      }
      break;
    case 'Escape':
      setIsFocused(false);
      break;
  }
};
```

## 🎨 UI структура
```jsx
<div className="relative w-full">
  {/* Основное поле поиска */}
  <div className="relative">
    <input
      type="text"
      value={query}
      onChange={handleChange}
      onKeyDown={handleKeyDown}
      onFocus={() => setIsFocused(true)}
      placeholder={placeholder || t('search.placeholder')}
      className="input input-bordered w-full pr-10"
    />
    
    {/* Кнопка поиска */}
    <button
      onClick={handleSearch}
      className="absolute right-2 top-1/2 -translate-y-1/2"
    >
      <SearchIcon className="w-5 h-5" />
    </button>
  </div>
  
  {/* Dropdown с предложениями */}
  {isFocused && (query.length > 0 || searchHistory.length > 0) && (
    <div className="absolute z-50 w-full mt-1 bg-base-100 rounded-lg shadow-lg">
      {/* Предложения из API */}
      {suggestions.length > 0 && (
        <div className="py-2">
          {suggestions.map((suggestion, index) => (
            <button
              key={index}
              onClick={() => handleSelectSuggestion(suggestion)}
              className={cn(
                "w-full px-4 py-2 text-left hover:bg-base-200",
                selectedIndex === index && "bg-base-200"
              )}
            >
              <div className="flex items-center justify-between">
                <span>{suggestion.text}</span>
                {suggestion.count && (
                  <span className="text-sm text-base-content/60">
                    {suggestion.count}
                  </span>
                )}
              </div>
              {suggestion.type === 'category' && (
                <span className="text-xs text-base-content/60">
                  {t('search.inCategory')}
                </span>
              )}
            </button>
          ))}
        </div>
      )}
      
      {/* История поиска */}
      {query.length === 0 && searchHistory.length > 0 && (
        <div className="py-2">
          <div className="px-4 py-1 text-sm text-base-content/60">
            {t('search.recent')}
          </div>
          {searchHistory.map((item, index) => (
            <button
              key={index}
              onClick={() => handleHistoryClick(item)}
              className="w-full px-4 py-2 text-left hover:bg-base-200"
            >
              <HistoryIcon className="inline w-4 h-4 mr-2" />
              {item}
            </button>
          ))}
        </div>
      )}
    </div>
  )}
</div>
```

## 🔗 Зависимости

### Внешние библиотеки
- `react`: Хуки состояния и эффектов
- `next/navigation`: Навигация между страницами
- `clsx`: Управление классами
- `next-intl`: Интернационализация

### Внутренние зависимости
- `@/services/search`: UnifiedSearchService для API
- `@/hooks/useDebounce`: Дебаунсинг запросов
- `@/utils/cn`: Утилита для классов
- `./icons`: SVG иконки поиска

### API интеграция
```typescript
// Получение предложений
unifiedSearchService.getSuggestions(query)
// Response: { suggestions: SearchSuggestion[] }

// Выполнение поиска
router.push(`/${locale}/search?q=${encodeURIComponent(query)}`);
```

## 📊 Управление состоянием
- **LocalStorage**: История поиска пользователя
- **Component State**: Временные данные (предложения, фокус)
- **URL State**: Поисковый запрос в query параметрах

## 🔄 Жизненный цикл
1. **Mount**: Загрузка истории из localStorage
2. **Input Change**: Дебаунс → API запрос → Обновление предложений
3. **Submit**: Сохранение в историю → Навигация
4. **Unmount**: Очистка таймеров дебаунса

## 📱 Адаптивность
- **Mobile**: Полная ширина, крупные элементы для тач
- **Desktop**: Фиксированная ширина в хедере
- **Keyboard**: Полная поддержка навигации

## ⚡ Оптимизации
1. **Debouncing**: Задержка API запросов (300ms)
2. **Memoization**: Кэширование предложений
3. **Lazy Loading**: Динамическая загрузка иконок
4. **LocalStorage**: Кэш истории поиска

## 🎯 Примеры использования

### Базовое использование
```jsx
<SearchBar 
  placeholder="Поиск товаров..."
  onSearch={(query) => console.log('Searching:', query)}
/>
```

### С настройками
```jsx
<SearchBar 
  className="max-w-lg"
  autoFocus
  showSuggestions={true}
  maxSuggestions={8}
  searchDebounceMs={500}
  onSearch={handleSearch}
/>
```

### В хедере
```jsx
<Header>
  <SearchBar className="flex-1 max-w-xl mx-4" />
</Header>
```

## 🐛 Известные проблемы
1. **TODO**: Кэширование предложений в памяти
2. **TODO**: Поддержка голосового ввода
3. **Hardcoded**: Минимальная длина запроса (2 символа)
4. **Missing**: Аналитика поисковых запросов

## 🔒 Безопасность
- XSS защита через санитизацию запросов
- Валидация длины запроса
- Rate limiting на уровне API
- Безопасное хранение истории в localStorage