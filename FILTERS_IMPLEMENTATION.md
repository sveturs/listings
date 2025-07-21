# Инструкция по реализации системы фильтров

## Дата создания: 21.07.2025

## 📋 Обзор системы фильтров

Система фильтров на главной странице маркетплейса состоит из двух типов:
1. **Постоянные фильтры** - всегда доступны независимо от категории
2. **Динамические фильтры** - загружаются на основе атрибутов выбранной категории

## 🏗️ Архитектура

### Компоненты системы:

```
HomePage.tsx
├── CategorySidebar.tsx (выбор категории)
├── ListingFilters.tsx (управление фильтрами)
│   ├── Постоянные фильтры (цена, состояние, продавец, доставка)
│   └── SmartFilters.tsx (динамические атрибуты категории)
└── MarketplaceList.tsx (отображение результатов с фильтрами)
```

### Поток данных:

1. **Выбор категории** → CategorySidebar → HomePage
2. **Загрузка атрибутов** → ListingFilters → API `/categories/{id}/attributes`
3. **Применение фильтров** → ListingFilters → MarketplaceList
4. **Поиск с фильтрами** → MarketplaceList → UnifiedSearchService

## 🔧 Реализация компонентов

### 1. ListingFilters Component

```typescript
export interface FilterValues {
  // Постоянные фильтры
  priceMin?: number;
  priceMax?: number;
  condition?: string;
  sellerType?: 'private' | 'company';
  hasDelivery?: boolean;
  
  // Динамические фильтры атрибутов
  attributeFilters?: Record<string, any>;
}

interface ListingFiltersProps {
  selectedCategoryId?: number | null;
  filters: FilterValues;
  onFiltersChange: (filters: FilterValues) => void;
  className?: string;
}
```

**Основные функции:**
- Отображение постоянных фильтров
- Загрузка атрибутов категории с `is_filterable = true`
- Управление состоянием фильтров
- Передача изменений в родительский компонент

### 2. SmartFilters Component

Динамически генерирует фильтры на основе атрибутов категории:

```typescript
interface SmartFiltersProps {
  categoryId: number;
  onChange: (filters: Record<string, any>) => void;
  lang: string;
  className?: string;
}
```

**Поддерживаемые типы фильтров:**
- `select` - выпадающий список
- `multiselect` - множественный выбор
- `number` - диапазон чисел
- `boolean` - чекбокс
- `text` - текстовый поиск

### 3. Интеграция с MarketplaceList

```typescript
// Конвертация фильтров для API
function convertFiltersToSearchParams(filters: FilterValues) {
  const searchParams: Record<string, any> = {};
  
  // Постоянные фильтры
  if (filters.priceMin) searchParams.priceMin = filters.priceMin;
  if (filters.priceMax) searchParams.priceMax = filters.priceMax;
  if (filters.condition) searchParams.condition = filters.condition;
  if (filters.sellerType) {
    searchParams.storefrontID = 
      filters.sellerType === 'company' ? 'not_null' : null;
  }
  
  // Атрибуты категории
  if (filters.attributeFilters) {
    searchParams.attributeFilters = filters.attributeFilters;
  }
  
  return searchParams;
}
```

## 📊 Постоянные фильтры

### 1. Диапазон цен (Price Range)
```tsx
<div className="flex gap-2">
  <input
    type="number"
    placeholder={t('filters.priceFrom')}
    value={filters.priceMin || ''}
    onChange={(e) => handleFilterChange('priceMin', parseFloat(e.target.value))}
  />
  <input
    type="number"
    placeholder={t('filters.priceTo')}
    value={filters.priceMax || ''}
    onChange={(e) => handleFilterChange('priceMax', parseFloat(e.target.value))}
  />
</div>
```

### 2. Состояние товара (Condition)
```tsx
<select value={filters.condition || ''} onChange={(e) => handleFilterChange('condition', e.target.value)}>
  <option value="">{t('filters.anyCondition')}</option>
  <option value="new">{t('condition.new')}</option>
  <option value="used">{t('condition.used')}</option>
  <option value="refurbished">{t('condition.refurbished')}</option>
  <option value="damaged">{t('condition.damaged')}</option>
</select>
```

### 3. Тип продавца (Seller Type)
```tsx
<select value={filters.sellerType || ''} onChange={(e) => handleFilterChange('sellerType', e.target.value)}>
  <option value="">{t('filters.anySeller')}</option>
  <option value="private">{t('filters.privateSeller')}</option>
  <option value="company">{t('filters.companySeller')}</option>
</select>
```

### 4. С доставкой (With Delivery)
```tsx
<input
  type="checkbox"
  checked={filters.hasDelivery || false}
  onChange={(e) => handleFilterChange('hasDelivery', e.target.checked)}
/>
```

## 🔄 Динамические фильтры

### Загрузка атрибутов категории:

```typescript
useEffect(() => {
  const loadAttributes = async () => {
    if (!selectedCategoryId) {
      setFilterableAttributes([]);
      return;
    }
    
    const response = await MarketplaceService.getCategoryAttributes(
      selectedCategoryId,
      locale
    );
    
    // Фильтруем только атрибуты с is_filterable = true
    const filterable = response.data.filter(attr => attr.is_filterable);
    setFilterableAttributes(filterable);
  };
  
  loadAttributes();
}, [selectedCategoryId, locale]);
```

### Типы динамических фильтров:

1. **Select атрибуты**
```typescript
if (attribute.attribute_type === 'select') {
  return (
    <select onChange={(e) => updateFilter(attribute.name, e.target.value)}>
      {attribute.options.values.map(value => (
        <option key={value} value={value}>{value}</option>
      ))}
    </select>
  );
}
```

2. **Number атрибуты (range)**
```typescript
if (attribute.attribute_type === 'number') {
  return (
    <div className="flex gap-2">
      <input
        type="number"
        placeholder="Min"
        onChange={(e) => updateFilter(`${attribute.name}_min`, e.target.value)}
      />
      <input
        type="number"
        placeholder="Max"
        onChange={(e) => updateFilter(`${attribute.name}_max`, e.target.value)}
      />
    </div>
  );
}
```

## 🌐 Локализация

### Ключи переводов для фильтров:

```json
{
  "marketplace": {
    "filters": {
      "title": "Filters",
      "clearAll": "Clear all",
      "priceRange": "Price range",
      "priceFrom": "From",
      "priceTo": "To",
      "condition": "Condition",
      "anyCondition": "Any condition",
      "sellerType": "Seller type",
      "anySeller": "Any seller",
      "privateSeller": "Private sellers",
      "companySeller": "Company sellers",
      "withDelivery": "With delivery",
      "categoryFilters": "Category filters",
      "noAttributeFilters": "No additional filters for this category"
    }
  }
}
```

## 🔌 Backend интеграция

### API Endpoints:

1. **Получение атрибутов категории:**
```
GET /api/v1/marketplace/categories/{id}/attributes?lang={locale}
```

2. **Поиск с фильтрами:**
```
POST /api/v1/unified-search
{
  "query": "",
  "product_types": ["marketplace", "storefront"],
  "category_id": "1001",
  "priceMin": 100,
  "priceMax": 1000,
  "condition": "new",
  "attributeFilters": {
    "brand": "Apple",
    "storage": "128GB"
  }
}
```

### Обработка фильтров в Backend:

```go
// Постоянные фильтры
if priceMin > 0 {
    query = query.Where("price >= ?", priceMin)
}
if priceMax > 0 {
    query = query.Where("price <= ?", priceMax)
}
if condition != "" {
    query = query.Where("condition = ?", condition)
}

// Динамические атрибуты
if attributeFilters != nil {
    for key, value := range attributeFilters {
        query = query.Where("attributes->? @> ?", key, value)
    }
}
```

## ✅ Чеклист реализации

- [x] Создать компонент ListingFilters
- [x] Реализовать постоянные фильтры
- [x] Интегрировать SmartFilters для динамических атрибутов
- [x] Добавить загрузку атрибутов при выборе категории
- [x] Реализовать передачу фильтров в MarketplaceList
- [x] Обновить UnifiedSearchService для поддержки фильтров
- [x] Добавить переводы для всех элементов фильтров
- [x] Протестировать работу фильтров
- [x] Добавить индикацию активных фильтров
- [x] Реализовать кнопку "Clear all"

## 🚀 Оптимизации и улучшения

### Реализовано:
1. ✅ Счетчик активных фильтров в заголовке
2. ✅ Сворачивание/разворачивание блока фильтров
3. ✅ Кнопка очистки всех фильтров
4. ✅ Визуальная индикация при наличии активных фильтров

### Планируется:
1. 🔄 Сохранение состояния фильтров в URL
2. 🔄 Debounce для числовых полей
3. 🔄 Показ количества товаров для каждого значения
4. 🔄 Виртуализация для больших списков опций
5. 🔄 Предзагрузка популярных категорий

## 📝 Примеры использования

### На главной странице:
```tsx
const [filters, setFilters] = useState<FilterValues>({});
const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null);

<CategorySidebar
  onCategorySelect={setSelectedCategoryId}
  selectedCategoryId={selectedCategoryId}
/>

<ListingFilters
  selectedCategoryId={selectedCategoryId}
  filters={filters}
  onFiltersChange={setFilters}
/>

<MarketplaceList
  filters={filters}
  selectedCategoryId={selectedCategoryId}
/>
```

### В других компонентах:
```tsx
// Можно использовать только постоянные фильтры
<ListingFilters
  filters={filters}
  onFiltersChange={setFilters}
  selectedCategoryId={null} // Без динамических атрибутов
/>
```

## 🐛 Отладка проблем

### Фильтры не применяются:
1. Проверить передачу filters в MarketplaceList
2. Убедиться что convertFiltersToSearchParams вызывается
3. Проверить Network запросы на наличие параметров фильтров

### Атрибуты не загружаются:
1. Проверить selectedCategoryId не null
2. Проверить ответ API `/categories/{id}/attributes`
3. Убедиться что атрибуты имеют `is_filterable = true`

### Проблемы с локализацией:
1. Проверить наличие переводов в messages файлах
2. Убедиться что используется правильный namespace (marketplace.filters)
3. Проверить передачу locale в API запросы

## 🔴 Известные проблемы и их решения

### 1. Бесконечный цикл при выборе категории (РЕШЕНО - 21.07.2025)
**Проблема**: При выборе категории на главной странице происходит бесконечный цикл перезагрузки данных.

**Причины**: 
1. Объект `filters` пересоздается при каждом рендере, что вызывает срабатывание useEffect в MarketplaceList.
2. В SmartFilters был вызов `onChange({})` в useEffect с зависимостью от `onChange`.

**Решения**: 
1. В MarketplaceList использована мемоизация через `useMemo` для создания стабильного `filtersKey`:
```typescript
const filtersKey = useMemo(() => {
  return JSON.stringify({
    priceMin: filters.priceMin,
    priceMax: filters.priceMax,
    condition: filters.condition,
    sellerType: filters.sellerType,
    hasDelivery: filters.hasDelivery,
    attributeFilters: filters.attributeFilters,
  });
}, [filters]);

// Использовать filtersKey вместо filters в зависимостях
useEffect(() => {
  // загрузка данных
}, [productTypes, selectedCategoryId, filtersKey, initialized, t]);
```

### 2. Переводы для condition фильтров (РЕШЕНО - 21.07.2025)
**Проблема**: Отсутствовали переводы для `marketplace.condition.*` во всех языках.

**Решение**: 
1. Добавлены переводы во все языковые файлы:
   - `ru.json`: Новое, Б/у, Восстановленное, Повреждённое
   - `en.json`: New, Used, Refurbished, Damaged
   - `sr.json`: Novo, Korišćeno, Obnovljeno, Oštećeno

2. Удалены дубликаты ключей в ru.json.

3. Обновлен компонент `ListingFilters.tsx` для использования правильного namespace:
```typescript
// useTranslations('marketplace')
// t('condition.new') - без двойного marketplace
```

### 3. Новый бесконечный цикл в SmartFilters (РЕШЕНО - 21.07.2025)
**Проблема**: При выборе категории с атрибутами возникала ошибка "Maximum update depth exceeded".

**Причина**: В SmartFilters.tsx был вызов `onChange({})` в useEffect с зависимостью от `onChange`, что создавало бесконечный цикл.

**Решение**: Убрали `onChange` из зависимостей useEffect и добавили проверку:
```typescript
useEffect(() => {
  setFilterValues({});
  // Вызываем onChange только если действительно были фильтры
  if (Object.keys(filterValues).length > 0) {
    onChange({});
  }
}, [categoryId]); // Убираем onChange из зависимостей
```

### 4. Недостающий перевод marketplace.clearFilter (РЕШЕНО - 21.07.2025)
**Проблема**: Отсутствовал перевод для `marketplace.clearFilter`.

**Решение**: Добавлены переводы:
- ru.json: "Очистить фильтр"
- en.json: "Clear filter"
- sr.json: "Обриши филтер"

### 5. Неправильный namespace в SmartFilters (РЕШЕНО - 21.07.2025)
**Проблема**: SmartFilters использовал namespace 'map' вместо 'marketplace'.

**Решение**: Исправлено:
```typescript
// Было: const t = useTranslations('map');
// Стало: const t = useTranslations('marketplace');
```