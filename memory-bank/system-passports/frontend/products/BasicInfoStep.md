# Паспорт компонента: BasicInfoStep

## Метаданные
- **Путь**: `frontend/svetu/src/components/products/steps/BasicInfoStep.tsx`
- **Роль**: Шаг ввода основной информации о товаре
- **Тип**: Form Component
- **Размер**: 340 строк

## Назначение
Второй шаг в wizard'е создания товара. Собирает основную информацию:
- Название и описание товара
- Цена и валюта
- Количество на складе
- SKU и штрих-код
- Статус активности

## Props структура
```typescript
interface BasicInfoStepProps {
  onNext: () => void;  // Переход к следующему шагу
  onBack: () => void;  // Возврат к предыдущему шагу
}
```

## Зависимости
- **Библиотеки**: React, next-intl
- **Контекст**: `CreateProductContext` - глобальное состояние
- **Хуки**: useState, useEffect

## Управление состоянием

### Локальное состояние формы
```typescript
const [formData, setFormData] = useState({
  name: state.productData.name || '',
  description: state.productData.description || '',
  price: state.productData.price || 0,
  currency: state.productData.currency || 'RSD',
  stock_quantity: state.productData.stock_quantity || 0,
  sku: state.productData.sku || '',
  barcode: state.productData.barcode || '',
  is_active: state.productData.is_active !== undefined 
    ? state.productData.is_active 
    : true,
});
```

### Синхронизация с глобальным состоянием
```typescript
useEffect(() => {
  setProductData(formData);
}, [formData, setProductData]);
```

## Бизнес-логика

### Обработка изменений
```typescript
const handleChange = (e) => {
  const { name, value, type } = e.target;
  
  if (type === 'checkbox') {
    const checked = e.target.checked;
    setFormData(prev => ({ ...prev, [name]: checked }));
  } else if (name === 'price' || name === 'stock_quantity') {
    const numValue = parseFloat(value) || 0;
    setFormData(prev => ({ ...prev, [name]: numValue }));
  } else {
    setFormData(prev => ({ ...prev, [name]: value }));
  }
  
  clearError(name); // Очистка ошибки при изменении
};
```

### Валидация
```typescript
const validateForm = (): boolean => {
  let isValid = true;
  
  if (!formData.name || formData.name.length < 3) {
    setError('name', t('storefronts.products.nameRequired'));
    isValid = false;
  }
  
  if (!formData.description || formData.description.length < 10) {
    setError('description', t('storefronts.products.descriptionRequired'));
    isValid = false;
  }
  
  if (formData.price <= 0) {
    setError('price', t('storefronts.products.priceRequired'));
    isValid = false;
  }
  
  return isValid;
};
```

## UI структура

### Layout - две колонки
1. **Левая колонка - Основная информация**
   - Название товара (обязательное)
   - Описание (обязательное, textarea)

2. **Правая колонка - Цена и инвентарь**
   - Цена с выбором валюты
   - Количество на складе
   - SKU (опционально)
   - Штрих-код (опционально)
   - Toggle активности

### Визуальные элементы
- Карточки с иконками для группировки
- Input группы для цены с валютой
- Toggle для статуса активности
- Подсказки (label-text-alt) для полей

## Примеры использования

### В ProductWizard
```tsx
case 1:
  return <BasicInfoStep onNext={nextStep} onBack={prevStep} />;
```

### Валидация при переходе
```tsx
const handleNext = () => {
  if (validateForm()) {
    onNext();
  }
  // Ошибки отображаются автоматически
};
```

## Известные особенности

### Позитивные
- ✅ Двухколоночный адаптивный layout
- ✅ Визуальная группировка связанных полей
- ✅ Inline валидация с очисткой при изменении
- ✅ Поддержка трех валют (RSD, EUR, USD)
- ✅ Автоматическая синхронизация с контекстом

### Технический долг
- ⚠️ Жестко закодированный список валют
- ⚠️ Нет маски ввода для штрих-кода
- ⚠️ Отсутствует проверка уникальности SKU
- ⚠️ Минимальные длины строк захардкожены

### Возможные улучшения
- 💡 Динамический список валют из API
- 💡 Автогенерация SKU на основе категории
- 💡 Сканирование штрих-кода камерой
- 💡 Rich text editor для описания
- 💡 Предпросмотр карточки товара