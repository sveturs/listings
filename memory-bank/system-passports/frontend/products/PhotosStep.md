# Паспорт компонента: PhotosStep

## Метаданные
- **Путь**: `frontend/svetu/src/components/products/steps/PhotosStep.tsx`
- **Роль**: Шаг загрузки и управления фотографиями товара
- **Тип**: Upload Component
- **Размер**: 352 строки

## Назначение
Четвертый шаг в wizard'е создания товара. Позволяет:
- Загружать до 10 фотографий товара
- Просматривать превью загруженных изображений
- Изменять порядок фотографий (первая - главная)
- Удалять ненужные изображения

## Props структура
```typescript
interface PhotosStepProps {
  onNext: () => void;  // Переход к следующему шагу
  onBack: () => void;  // Возврат к предыдущему шагу
}
```

## Зависимости
- **Библиотеки**: React, next-intl
- **Контекст**: `CreateProductContext` - глобальное состояние
- **Утилиты**: `toast` - уведомления об ошибках
- **Хуки**: useState, useEffect

## Управление состоянием

### Локальное состояние
```typescript
const [images, setImagesState] = useState<File[]>(state.images || []);
const [previews, setPreviews] = useState<string[]>([]);
```

### Генерация превью
```typescript
useEffect(() => {
  const newPreviews: string[] = [];
  images.forEach((file) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      newPreviews.push(reader.result as string);
      if (newPreviews.length === images.length) {
        setPreviews([...newPreviews]);
      }
    };
    reader.readAsDataURL(file);
  });
}, [images]);
```

## Бизнес-логика

### Валидация файлов
```typescript
// Максимум 10 изображений
if (files.length + images.length > 10) {
  toast.error(t('storefronts.products.maxImagesError'));
  return;
}

// Проверка типа и размера
const invalidFiles = files.filter((file) => {
  const isValidType = file.type.startsWith('image/');
  const isValidSize = file.size <= 10 * 1024 * 1024; // 10MB
  return !isValidType || !isValidSize;
});
```

### Управление порядком
```typescript
const moveImage = (fromIndex: number, toIndex: number) => {
  const newImages = [...images];
  const newPreviews = [...previews];
  
  // Меняем местами
  [newImages[fromIndex], newImages[toIndex]] = 
    [newImages[toIndex], newImages[fromIndex]];
  [newPreviews[fromIndex], newPreviews[toIndex]] = 
    [newPreviews[toIndex], newPreviews[fromIndex]];
    
  setImagesState(newImages);
  setImages(newImages);
  setPreviews(newPreviews);
};
```

## UI структура

### Статистика
```
┌──────────────────┬──────────────────┐
│ 📸 Загружено: 3  │ ⚡ Рекомендуется │
│ Максимум: 10     │ 1200px, 1:1      │
└──────────────────┴──────────────────┘
```

### Дроп-зона
- Drag & Drop область
- Кнопка выбора файлов
- Поддерживаемые форматы: JPEG, PNG, GIF, WebP
- Визуальная обратная связь при наведении

### Сетка превью
```
┌───┬───┬───┬───┐
│ 1 │ 2 │ 3 │ 4 │  Номера фотографий
├───┼───┼───┼───┤
│ ← │ → │ × │   │  Кнопки управления (при hover)
└───┴───┴───┴───┘
```

### Особенности UI
- Первое фото помечено как "Главное фото"
- Hover эффект показывает кнопки управления
- Номера фотографий для ориентации
- Подсказки по оптимизации фото

## Примеры использования

### В ProductWizard
```tsx
case 3:
  return <PhotosStep onNext={nextStep} onBack={prevStep} />;
```

### Обработка загрузки
```tsx
const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
  const files = Array.from(e.target.files || []);
  
  // Валидация
  if (!validateFiles(files)) return;
  
  // Добавление к существующим
  const newImages = [...images, ...files];
  setImagesState(newImages);
  setImages(newImages); // В контекст
  clearError('images');
};
```

## Известные особенности

### Позитивные
- ✅ Интуитивный drag & drop интерфейс
- ✅ Визуальное управление порядком фото
- ✅ Превью в реальном времени
- ✅ Валидация типов и размеров файлов
- ✅ Информативная статистика

### Технический долг
- ⚠️ Нет прогресс-бара для загрузки
- ⚠️ Отсутствует сжатие изображений на клиенте
- ⚠️ Нет обрезки/поворота изображений
- ⚠️ FileReader может быть медленным для больших файлов

### Возможные улучшения
- 💡 Добавить drag & drop для изменения порядка
- 💡 Оптимизация изображений перед загрузкой
- 💡 Встроенный редактор (обрезка, поворот)
- 💡 Загрузка по URL
- 💡 Интеграция с камерой устройства