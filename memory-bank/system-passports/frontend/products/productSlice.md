# Паспорт компонента: productSlice

## Метаданные
- **Путь**: `frontend/svetu/src/store/slices/productSlice.ts`
- **Роль**: Redux slice для управления состоянием товаров
- **Тип**: State Management (Redux Toolkit)
- **Размер**: 405 строк

## Назначение
Централизованное управление состоянием товаров в витрине:
- Хранение списка товаров с пагинацией
- Фильтрация и сортировка
- Массовый выбор и операции
- UI состояния (режим отображения, выбор)
- Асинхронные операции (удаление, обновление, экспорт)

## Структура состояния
```typescript
interface ProductState {
  // Данные
  products: StorefrontProduct[];      // Список товаров
  selectedIds: number[];              // Выбранные товары
  loading: boolean;                   // Общая загрузка
  error: string | null;              // Ошибка
  
  // Фильтры
  filters: {
    search: string;                   // Поисковый запрос
    categoryId: number | null;        // ID категории
    minPrice: number | null;          // Мин. цена
    maxPrice: number | null;          // Макс. цена
    stockStatus: 'all' | 'in_stock' | 'low_stock' | 'out_of_stock';
    isActive: boolean | null;         // Активность
  };
  
  // Пагинация
  pagination: {
    page: number;                     // Текущая страница
    limit: number;                    // Элементов на странице
    total: number;                    // Всего элементов
    hasMore: boolean;                 // Есть ли еще
  };
  
  // Массовые операции
  bulkOperation: {
    isProcessing: boolean;            // Идет обработка
    progress: number;                 // Прогресс (0-100)
    total: number;                    // Всего операций
    errors: BulkOperationError[];     // Ошибки
    successCount: number;             // Успешно обработано
    currentOperation: 'idle' | 'delete' | 'update' | 'status' | 'export';
  };
  
  // UI
  ui: {
    isSelectMode: boolean;            // Режим выбора
    viewMode: 'grid' | 'list' | 'table';  // Режим отображения
    sortBy: 'name' | 'price' | 'created_at' | 'stock_quantity';
    sortOrder: 'asc' | 'desc';
  };
}
```

## Async Thunks

### bulkDeleteProducts
```typescript
export const bulkDeleteProducts = createAsyncThunk(
  'products/bulkDelete',
  async ({ storefrontSlug, productIds }) => {
    const response = await productApi.bulkDelete(storefrontSlug, productIds);
    
    // Уведомления
    if (response?.deleted?.length > 0) {
      toast.success(`Удалено товаров: ${response.deleted.length}`);
    }
    if (response?.failed?.length > 0) {
      toast.error(`Не удалось удалить: ${response.failed.length}`);
    }
    
    return response;
  }
);
```

### bulkUpdateStatus
- Массовое изменение статуса активности
- Обновляет локальное состояние после успеха
- Показывает toast уведомления

### exportProducts
- Экспорт в CSV или XML
- Поддержка выбранных товаров или всех
- Автоматическая загрузка файла

## Синхронные actions

### Управление выбором
- `toggleProductSelection` - переключить выбор товара
- `selectAll` - выбрать все видимые
- `clearSelection` - сбросить выбор
- `selectByFilter` - выбрать по условию

### UI управление
- `toggleSelectMode` - режим массового выбора
- `setViewMode` - изменить режим отображения
- `setSortBy` - установить сортировку

### Фильтрация
- `setFilters` - обновить фильтры (сбрасывает пагинацию)
- `resetFilters` - сбросить все фильтры

### Управление товарами
- `setProducts` - заменить список
- `appendProducts` - добавить к списку
- `updateProduct` - обновить один товар
- `removeProducts` - удалить товары

## Extra Reducers

### Обработка bulkDelete
```typescript
.addCase(bulkDeleteProducts.fulfilled, (state, action) => {
  // Удаляем из списка
  const idsToRemove = new Set(action.payload.deleted);
  state.products = state.products.filter(
    p => p.id && !idsToRemove.has(p.id)
  );
  
  // Удаляем из выбранных
  state.selectedIds = state.selectedIds.filter(
    id => !idsToRemove.has(id)
  );
  
  // Сохраняем ошибки
  state.bulkOperation.errors = action.payload.failed || [];
});
```

## Примеры использования

### В компоненте списка
```tsx
const dispatch = useDispatch();
const { products, filters, ui } = useSelector(state => state.products);

// Изменение фильтров
dispatch(setFilters({ search: 'iPhone' }));

// Массовое удаление
dispatch(bulkDeleteProducts({ 
  storefrontSlug: 'my-store', 
  productIds: [1, 2, 3] 
}));
```

### Массовый выбор
```tsx
// Переключить режим
dispatch(toggleSelectMode());

// Выбрать товары дешевле 100
dispatch(selectByFilter(product => product.price < 100));

// Деактивировать выбранные
dispatch(bulkUpdateStatus({
  storefrontSlug,
  productIds: selectedIds,
  isActive: false
}));
```

## Известные особенности

### Позитивные
- ✅ Полноценное управление состоянием товаров
- ✅ Поддержка массовых операций с прогрессом
- ✅ Гибкая система фильтров
- ✅ Интеграция с toast уведомлениями
- ✅ Type-safe с TypeScript

### Технический долг
- ⚠️ Нет кеширования запросов
- ⚠️ Отсутствует нормализация данных
- ⚠️ Нет оптимистичных обновлений
- ⚠️ Прогресс операций не real-time

### Возможные улучшения
- 💡 Интеграция с RTK Query для кеширования
- 💡 WebSocket для real-time прогресса
- 💡 Undo/Redo для массовых операций
- 💡 Сохранение фильтров в localStorage
- 💡 Виртуализация для больших списков