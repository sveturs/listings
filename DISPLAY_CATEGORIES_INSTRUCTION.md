# Инструкция по правильному отображению категорий

## Проблема
Категории должны отображаться рекурсивно с неограниченной вложенностью и правильно реагировать на смену языка интерфейса.

## Ключевые принципы отображения

### 1. Рекурсивное отображение без ограничений
```tsx
// ПРАВИЛЬНО - рекурсивная функция
const renderCategoryTree = (categories: Category[], level = 0) => {
  return categories.map(category => (
    <div key={category.id} style={{ paddingLeft: `${level * 20}px` }}>
      <div>{category.name}</div>
      {category.children && category.children.length > 0 && 
        renderCategoryTree(category.children, level + 1)
      }
    </div>
  ));
};

// НЕПРАВИЛЬНО - ограничение уровней
if (level < 3) { // НЕ ДЕЛАТЬ ТАК!
  // render children
}
```

### 2. Визуальная иерархия
- Использовать отступы для показа вложенности
- Иконки для раскрытия/сворачивания веток
- Разные стили для родительских и листовых категорий

### 3. Показ полного пути
```tsx
// Пример: "Транспорт > Автомобили > Седаны"
const getCategoryPath = (categoryId: number, categories: Category[]): string => {
  const path: string[] = [];
  
  const findPath = (cats: Category[], targetId: number): boolean => {
    for (const cat of cats) {
      if (cat.id === targetId) {
        path.unshift(cat.name);
        return true;
      }
      if (cat.children && findPath(cat.children, targetId)) {
        path.unshift(cat.name);
        return true;
      }
    }
    return false;
  };
  
  findPath(categories, categoryId);
  return path.join(' > ');
};
```

## Места где нужно правильное отображение

### 1. Выпадающий список категорий (Select/Dropdown)
```tsx
// Компонент должен показывать дерево с отступами
<CategoryTreeSelect>
  📁 Транспорт
    📁 Автомобили  
      📄 Седаны
      📄 Внедорожники
    📁 Мотоциклы
  📁 Недвижимость
    📁 Квартиры
      📄 Студии
      📄 Однокомнатные
</CategoryTreeSelect>
```

### 2. Фильтры на главной странице
- Чекбоксы с древовидной структурой
- Счетчики количества объявлений
- Возможность выбрать родительскую категорию (выберет все дочерние)

### 3. Breadcrumbs (хлебные крошки)
```tsx
<Breadcrumbs>
  <Link href="/">Главная</Link> / 
  <Link href="/transport">Транспорт</Link> / 
  <Link href="/transport/cars">Автомобили</Link> / 
  Седаны
</Breadcrumbs>
```

### 4. Сайдбар категорий
- Аккордеон с раскрывающимися ветками
- Подсветка активной категории
- Плавная анимация раскрытия

## Поддержка многоязычности

### Backend обработка языка:
Backend автоматически обрабатывает язык через LocaleMiddleware, который:
1. Проверяет query параметр `lang` или `locale`
2. Проверяет заголовок `Accept-Language`
3. Проверяет cookie `locale`
4. По умолчанию использует сербский язык (`sr`)

```tsx
// Frontend должен передавать язык в URL
const response = await fetch(`/api/v1/marketplace/categories?lang=${locale}`);
```

### Frontend интеграция с next-intl:
```tsx
import { useLocale } from 'next-intl';

const CategorySidebar = () => {
  const locale = useLocale(); // Получаем текущий язык
  
  useEffect(() => {
    const fetchCategories = async () => {
      // Передаем язык в API запрос
      const response = await MarketplaceService.getCategories(locale);
      // ...
    };
    
    fetchCategories();
  }, [locale]); // Перезапрашиваем при смене языка
};
```

### Важно для мобильной версии:
Убедитесь, что мобильные компоненты также используют текущую локаль:
```tsx
// В мобильном drawer
const MobileCategoryDrawer = () => {
  const locale = useLocale();
  // Используем тот же CategorySidebar компонент
  return <CategorySidebar locale={locale} />;
};
```

### Решение проблемы с отсутствующими переводами:
Если категории отображаются на сербском языке независимо от выбранного языка, это означает, что в базе данных отсутствуют переводы. Backend возвращает название на языке по умолчанию (сербский).

Для добавления переводов:
1. Используйте админ-панель для перевода категорий
2. Или выполните SQL-запрос для добавления переводов:
```sql
-- Пример добавления перевода для категории
INSERT INTO translations (entity_type, entity_id, field_name, language, value)
VALUES 
  ('category', 1, 'name', 'ru', 'Электроника'),
  ('category', 1, 'name', 'en', 'Electronics');
```

### Проверка работы переключения языка:
1. Откройте Network вкладку в DevTools
2. Переключите язык на странице
3. Найдите запрос к `/api/v1/marketplace/categories`
4. Убедитесь, что передается правильный параметр `lang` (en/ru/sr)
5. Проверьте ответ - если переводов нет, все названия будут на сербском

## Работа с кешем при отображении

### Проверка актуальности данных:
```tsx
// В компоненте
useEffect(() => {
  // Если категории пустые или устарели - запросить заново
  if (!categories || categories.length === 0) {
    fetchCategories();
  }
}, []);

// Принудительное обновление
const handleForceRefresh = async () => {
  // Добавить параметр для обхода кеша
  const response = await fetch('/api/v1/marketplace/categories?refresh=true');
  // ...
};
```

## Стилизация с DaisyUI

### Древовидный список:
```tsx
<ul className="menu menu-compact">
  {categories.map(category => (
    <li key={category.id}>
      <details open={expandedCategories.includes(category.id)}>
        <summary>{category.name}</summary>
        {category.children && (
          <ul>
            {category.children.map(child => (
              <li key={child.id}>
                <a>{child.name}</a>
              </li>
            ))}
          </ul>
        )}
      </details>
    </li>
  ))}
</ul>
```

### Селект с группировкой:
```tsx
<select className="select select-bordered w-full">
  {renderOptions(categories)}
</select>

const renderOptions = (cats: Category[], prefix = ''): JSX.Element[] => {
  return cats.flatMap(cat => [
    <option key={cat.id} value={cat.id}>
      {prefix}{cat.name}
    </option>,
    ...(cat.children ? renderOptions(cat.children, prefix + '  ') : [])
  ]);
};
```

## Оптимизация производительности

### 1. Мемоизация дерева:
```tsx
const memoizedTree = useMemo(() => 
  buildCategoryTree(categories), 
  [categories]
);
```

### 2. Виртуализация для больших списков:
```tsx
import { FixedSizeTree } from 'react-vtree';

// Для списков с сотнями категорий
<FixedSizeTree
  treeWalker={treeWalker}
  itemSize={30}
  height={400}
  width="100%"
>
  {Node}
</FixedSizeTree>
```

### 3. Lazy loading веток:
```tsx
const [loadedBranches, setLoadedBranches] = useState<Set<number>>(new Set());

const handleExpandCategory = async (categoryId: number) => {
  if (!loadedBranches.has(categoryId)) {
    const children = await fetchCategoryChildren(categoryId);
    // Обновить дерево с новыми детьми
    setLoadedBranches(prev => new Set(prev).add(categoryId));
  }
};
```

## Проверка правильности отображения

### Чеклист:
- [ ] Все уровни вложенности отображаются (не обрезаются на 2-3 уровне)
- [ ] Визуально понятна иерархия (отступы, иконки, линии)
- [ ] Работает раскрытие/сворачивание веток
- [ ] Показывается полный путь где нужно
- [ ] Правильно работает выбор категории
- [ ] Отображается количество элементов в категории
- [ ] Поддерживается поиск по названию
- [ ] Корректно работает на мобильных устройствах

## Отладка проблем с отображением

### 1. Категории не показывают детей:
```bash
# Проверить структуру данных
console.log(JSON.stringify(categories, null, 2));

# Убедиться что children не null
categories.forEach(cat => {
  console.log(`${cat.name}: children = ${cat.children?.length || 0}`);
});
```

### 2. Проверка кеша:
```bash
# В DevTools Network
# Проверить заголовки ответа
Cache-Control: public, max-age=3600
X-Cache: HIT/MISS
```

### 3. Принудительное обновление:
```tsx
// Добавить кнопку для отладки
<button onClick={() => window.location.reload(true)}>
  Обновить без кеша
</button>
```

## Пример правильной реализации

✅ **Реализовано в**: `frontend/svetu/src/components/common/CategoryTreeSelector.tsx`

```tsx
// CategoryTreeDisplay.tsx
import React, { useState, useMemo } from 'react';

// Встроенные SVG иконки вместо внешних библиотек
const ChevronRight = ({ size = 16 }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <polyline points="9 18 15 12 9 6"></polyline>
  </svg>
);

interface Category {
  id: number;
  name: string;
  parent_id: number | null;
  children?: Category[];
  count?: number;
}

export const CategoryTreeDisplay: React.FC<{
  categories: Category[];
  selectedId?: number;
  onSelect?: (id: number) => void;
}> = ({ categories, selectedId, onSelect }) => {
  const [expanded, setExpanded] = useState<Set<number>>(new Set());

  const toggleExpand = (id: number) => {
    setExpanded(prev => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const renderTree = (items: Category[], level = 0): JSX.Element[] => {
    return items.map(item => {
      const hasChildren = item.children && item.children.length > 0;
      const isExpanded = expanded.has(item.id);
      const isSelected = item.id === selectedId;

      return (
        <div key={item.id}>
          <div 
            className={`
              flex items-center gap-2 px-2 py-1 cursor-pointer
              hover:bg-base-200 rounded
              ${isSelected ? 'bg-primary/10 text-primary' : ''}
            `}
            style={{ paddingLeft: `${level * 20 + 8}px` }}
            onClick={() => onSelect?.(item.id)}
          >
            {hasChildren && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  toggleExpand(item.id);
                }}
                className="btn btn-xs btn-ghost p-0 min-h-0 h-5 w-5"
              >
                {isExpanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
              </button>
            )}
            {!hasChildren && <div className="w-5" />}
            
            {hasChildren ? <Folder size={16} /> : <File size={16} />}
            
            <span className="flex-1">{item.name}</span>
            
            {item.count !== undefined && (
              <span className="badge badge-sm">{item.count}</span>
            )}
          </div>
          
          {hasChildren && isExpanded && (
            <div>
              {renderTree(item.children!, level + 1)}
            </div>
          )}
        </div>
      );
    });
  };

  return (
    <div className="overflow-auto max-h-96">
      {renderTree(categories)}
    </div>
  );
};
```

Эта реализация обеспечивает:
- Неограниченную вложенность
- Визуальную иерархию
- Раскрытие/сворачивание
- Выбор категории
- Отображение счетчиков
- Правильную стилизацию

## Важные моменты реализации

### 1. Построение дерева из плоского списка
Backend возвращает плоский список категорий. Frontend самостоятельно строит дерево:
```typescript
// Первый проход: создаем Map всех категорий
const categoryMap = new Map<number, Category>();
response.data.forEach((cat) => {
  categoryMap.set(cat.id, { ...cat, children: [] });
});

// Второй проход: строим иерархию через parent_id
response.data.forEach((cat) => {
  if (cat.parent_id) {
    const parent = categoryMap.get(cat.parent_id);
    if (parent) parent.children.push(categoryMap.get(cat.id));
  } else {
    rootCategories.push(categoryMap.get(cat.id));
  }
});
```

### 2. UX без лишних кликов
Для страниц создания/редактирования объявлений дерево категорий отображается сразу, без dropdown. Пользователь приходит за выбором категории - не заставляйте его делать лишний клик.

### 3. Встроенные иконки
Используйте встроенные SVG компоненты вместо внешних библиотек иконок для уменьшения размера бандла.

### 4. Ключи переводов
Все ключи переводов для marketplace компонентов должны быть в секции `marketplace` в файлах локализации:
```json
{
  "marketplace": {
    "selectCategory": "Select category",
    "searchCategories": "Search categories...",
    // остальные ключи
  }
}