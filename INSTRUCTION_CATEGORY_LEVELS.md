# Инструкция по отображению всех уровней вложенности категорий

## Проблема
В системе используется иерархическая структура категорий с неограниченным количеством уровней вложенности. Из-за кеширования иногда новые категории (особенно глубоких уровней) могут не отображаться.

## Решение

### 1. Очистка кеша при добавлении/изменении категорий

При любых изменениях в категориях необходимо очищать кеш:

```bash
# Очистить весь кеш Redis
docker exec hostel_redis redis-cli FLUSHALL

# Или очистить только кеш категорий (если реализован паттерн)
docker exec hostel_redis redis-cli --scan --pattern "category*" | xargs docker exec hostel_redis redis-cli DEL
```

### 2. Проверка корректности отображения

#### Backend API должен возвращать все уровни:
```bash
# Проверить количество категорий в API маркетплейса
curl -s "http://localhost:3000/api/v1/marketplace/categories" | jq '.data | length'

# Проверить количество категорий в API админки (требует авторизации)
curl -s -H "Authorization: Bearer YOUR_TOKEN" "http://localhost:3000/api/v1/admin/categories/all?lang=ru" | jq '.data | length'
```

Оба API должны возвращать одинаковое количество категорий.

### 3. Места в системе, где используются категории

#### a) Создание объявлений (`/ru/create-listing`)
**Файл**: `/frontend/svetu/src/app/[locale]/create-listing/components/CategorySelector.tsx`

Компонент должен использовать рекурсивное отображение:
```typescript
const renderCategoryTree = (categories: Category[], parentId: number | null = null, level = 0) => {
  return categories
    .filter(cat => cat.parent_id === parentId)
    .map(category => (
      <div key={category.id} style={{ marginLeft: `${level * 20}px` }}>
        <CategoryItem category={category} />
        {renderCategoryTree(categories, category.id, level + 1)}
      </div>
    ));
};
```

#### b) Создание товаров витрин
**Файлы**: 
- `/frontend/svetu/src/app/[locale]/storefront/[id]/products/create/page.tsx`
- Компоненты выбора категорий в витринах

#### c) Отображение товаров в витринах
**Файлы**:
- `/frontend/svetu/src/app/[locale]/storefront/[id]/page.tsx`
- Фильтры категорий в витринах

#### d) Карта с фильтром категорий (`/ru/map`)
**Файл**: `/frontend/svetu/src/app/[locale]/map/components/CategoryFilter.tsx`

### 4. Универсальный компонент для отображения категорий

Создайте переиспользуемый компонент:

```typescript
// src/components/categories/CategoryTreeSelector.tsx
import React, { useState } from 'react';

interface Category {
  id: number;
  name: string;
  parent_id: number | null;
  is_active: boolean;
  translations?: Record<string, string>;
}

interface CategoryTreeSelectorProps {
  categories: Category[];
  selectedId?: number | null;
  onSelect: (categoryId: number) => void;
  showInactive?: boolean;
  locale?: string;
}

export const CategoryTreeSelector: React.FC<CategoryTreeSelectorProps> = ({
  categories,
  selectedId,
  onSelect,
  showInactive = false,
  locale = 'ru'
}) => {
  const [expanded, setExpanded] = useState<Set<number>>(new Set());

  // Фильтруем категории по активности
  const filteredCategories = showInactive 
    ? categories 
    : categories.filter(c => c.is_active);

  // Рекурсивная функция для построения дерева
  const renderTree = (parentId: number | null = null, level = 0): JSX.Element[] => {
    return filteredCategories
      .filter(cat => cat.parent_id === parentId)
      .map(category => {
        const hasChildren = filteredCategories.some(c => c.parent_id === category.id);
        const isExpanded = expanded.has(category.id);
        const isSelected = selectedId === category.id;

        // Получаем локализованное название
        const displayName = category.translations?.[locale] || category.name;

        return (
          <div key={category.id} className="category-tree-item">
            <div 
              className={`
                flex items-center gap-2 p-2 hover:bg-base-200 rounded cursor-pointer
                ${isSelected ? 'bg-primary/10 font-semibold' : ''}
                ${!category.is_active ? 'opacity-50' : ''}
              `}
              style={{ paddingLeft: `${level * 20 + 8}px` }}
              onClick={() => onSelect(category.id)}
            >
              {hasChildren && (
                <button
                  className="btn btn-ghost btn-xs p-0 min-h-0 h-4 w-4"
                  onClick={(e) => {
                    e.stopPropagation();
                    setExpanded(prev => {
                      const next = new Set(prev);
                      if (next.has(category.id)) {
                        next.delete(category.id);
                      } else {
                        next.add(category.id);
                      }
                      return next;
                    });
                  }}
                >
                  <svg
                    className={`h-3 w-3 transition-transform ${isExpanded ? 'rotate-90' : ''}`}
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                </button>
              )}
              {!hasChildren && <div className="w-4" />}
              
              <span className="flex-1">{displayName}</span>
              
              {!category.is_active && (
                <span className="badge badge-ghost badge-xs">неактивна</span>
              )}
            </div>
            
            {hasChildren && isExpanded && (
              <div className="category-children">
                {renderTree(category.id, level + 1)}
              </div>
            )}
          </div>
        );
      });
  };

  return (
    <div className="category-tree-selector">
      {renderTree()}
    </div>
  );
};
```

### 5. Использование компонента

```typescript
// В любом месте где нужен выбор категории
import { CategoryTreeSelector } from '@/components/categories/CategoryTreeSelector';

function MyComponent() {
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);

  // Загрузка категорий
  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      const response = await fetch('/api/v1/marketplace/categories');
      const data = await response.json();
      setCategories(data.data);
    } catch (error) {
      console.error('Failed to load categories:', error);
    }
  };

  return (
    <CategoryTreeSelector
      categories={categories}
      selectedId={selectedCategory}
      onSelect={setSelectedCategory}
      showInactive={false}
      locale="ru"
    />
  );
}
```

### 6. Проверочный чеклист

При реализации отображения категорий проверьте:

- [ ] Компонент поддерживает рекурсивное отображение всех уровней
- [ ] Нет жестко заданных ограничений на количество уровней
- [ ] При загрузке используется актуальный API endpoint
- [ ] Обрабатывается локализация названий категорий
- [ ] Корректно отображаются неактивные категории (если нужно)
- [ ] При изменении категорий очищается кеш

### 7. Отладка проблем

Если категории не отображаются:

1. Проверьте консоль браузера на ошибки
2. Проверьте Network вкладку - какие данные приходят с API
3. Убедитесь что в ответе API есть все уровни категорий
4. Очистите кеш Redis: `docker exec hostel_redis redis-cli FLUSHALL`
5. Проверьте логи backend: `docker logs hostel_backend`

### 8. Автоматическая инвалидация кеша

Для автоматической очистки кеша при изменении категорий, убедитесь что в backend handlers вызывается:

```go
// В handlers при создании/обновлении/удалении категорий
h.InvalidateCategoryCache()
```

Это должно быть реализовано в:
- `CreateCategory` 
- `UpdateCategory`
- `DeleteCategory`
- `ReorderCategories`
- `MoveCategory`