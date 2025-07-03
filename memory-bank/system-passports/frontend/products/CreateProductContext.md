# Паспорт компонента: CreateProductContext

## Метаданные
- **Путь**: `frontend/svetu/src/contexts/CreateProductContext.tsx`
- **Роль**: Контекст управления состоянием создания товара
- **Тип**: Context Provider
- **Размер**: 278 строк

## Назначение
Глобальное управление состоянием процесса создания товара. Обеспечивает централизованное хранение данных, навигацию между шагами, валидацию и управление ошибками для всего wizard'а создания товара.

## Структура состояния
```typescript
interface ProductState {
  // Навигация
  currentStep: number;              // Текущий шаг (0-4)
  completedSteps: Set<number>;      // Завершенные шаги
  
  // Данные товара
  category?: MarketplaceCategory;   // Выбранная категория
  productData: CreateProductRequest; // Основные данные товара
  attributes: Record<number, any>;   // Атрибуты категории
  images: File[];                   // Загруженные изображения
  
  // Метаданные
  isDraft: boolean;                 // Черновик
  isValid: boolean;                 // Валидность формы
  errors: Record<string, string>;   // Ошибки по полям
  loading: boolean;                 // Состояние загрузки
}
```

## Actions (Redux pattern)
```typescript
type ProductAction =
  | { type: 'SET_STEP'; payload: number }
  | { type: 'COMPLETE_STEP'; payload: number }
  | { type: 'SET_CATEGORY'; payload: MarketplaceCategory }
  | { type: 'SET_PRODUCT_DATA'; payload: Partial<CreateProductRequest> }
  | { type: 'SET_ATTRIBUTE'; payload: { id: number; value: any } }
  | { type: 'SET_IMAGES'; payload: File[] }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: { field: string; message: string } }
  | { type: 'CLEAR_ERROR'; payload: string }
  | { type: 'CLEAR_ALL_ERRORS' }
  | { type: 'RESET' }
  | { type: 'LOAD_DRAFT'; payload: ProductState };
```

## Экспортируемый интерфейс
```typescript
interface CreateProductContextType {
  state: ProductState;
  dispatch: React.Dispatch<ProductAction>;
  
  // Навигация
  goToStep: (step: number) => void;
  nextStep: () => void;
  prevStep: () => void;
  completeStep: (step: number) => void;
  
  // Управление данными
  setCategory: (category: MarketplaceCategory) => void;
  setProductData: (data: Partial<CreateProductRequest>) => void;
  setAttribute: (id: number, value: any) => void;
  setImages: (images: File[]) => void;
  
  // Управление ошибками
  setError: (field: string, message: string) => void;
  clearError: (field: string) => void;
  clearAllErrors: () => void;
  
  // Валидация
  validateStep: (step: number) => boolean;
  canGoToStep: (step: number) => boolean;
}
```

## Бизнес-логика

### Валидация шагов
```typescript
const validateStep = (step: number): boolean => {
  switch (step) {
    case 0: // Категория
      return !!state.category;
    case 1: // Основная информация
      return !!(
        state.productData.name &&
        state.productData.description &&
        state.productData.price > 0
      );
    case 2: // Атрибуты
      return true; // Опциональны
    case 3: // Фотографии
      return state.images.length > 0;
    case 4: // Превью
      return true;
  }
};
```

### Управление атрибутами
- Атрибуты сохраняются в двух местах:
  - `attributes` - для UI компонентов
  - `productData.attributes` - для API запроса

### Управление навигацией
- Можно перейти только на завершенные шаги или следующий
- Валидация предыдущих шагов при попытке перехода вперед

## Примеры использования

### Провайдер в корне wizard'а
```tsx
<CreateProductProvider>
  <ProductWizard storefrontSlug="my-store" />
</CreateProductProvider>
```

### Использование в компоненте шага
```tsx
function BasicInfoStep() {
  const { 
    state, 
    setProductData, 
    setError, 
    clearError 
  } = useCreateProduct();
  
  const handleNameChange = (e) => {
    setProductData({ name: e.target.value });
    if (e.target.value.length >= 3) {
      clearError('name');
    }
  };
  
  // ...
}
```

### Валидация и переход
```tsx
function CategoryStep({ onNext }) {
  const { state, setCategory, validateStep } = useCreateProduct();
  
  const handleNext = () => {
    if (validateStep(0)) {
      onNext();
    } else {
      setError('category', 'Выберите категорию');
    }
  };
}
```

## Известные особенности

### Позитивные
- ✅ Централизованное управление состоянием
- ✅ Type-safe actions через TypeScript
- ✅ Встроенная валидация шагов
- ✅ Гибкое управление ошибками по полям
- ✅ Поддержка черновиков через LOAD_DRAFT

### Технический долг
- ⚠️ Нет персистентности состояния
- ⚠️ Отсутствует отмена/повтор действий
- ⚠️ Жестко закодированная валюта RSD
- ⚠️ Нет оптимистичных обновлений

### Возможные улучшения
- 💡 Сохранение в localStorage для восстановления
- 💡 История изменений с undo/redo
- 💡 Async валидация (проверка уникальности SKU)
- 💡 Интеграция с React Query для кеширования
- 💡 DevTools для отладки состояния