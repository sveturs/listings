# PreviewPublishStep.tsx

## Метаданные
- **Путь**: `frontend/svetu/src/components/create-listing/steps/PreviewPublishStep.tsx`
- **Роль**: UI компонент предпросмотра и публикации
- **Тип**: Step компонент мастера создания объявления (финальный)
- **Размер**: 363 строки

## Назначение
Финальный шаг мастера создания объявления. Показывает preview итогового объявления, обрабатывает публикацию через API, загрузку изображений и навигацию к опубликованному объявлению.

## Props структура
```typescript
interface PreviewPublishStepProps {
  onBack: () => void;      // Возврат к предыдущему шагу
  onComplete: () => void;  // Завершение процесса (не используется)
}
```

## Зависимости
### Внешние зависимости
- `react` - управление состоянием
- `next-intl` - интернационализация
- `next/image` - оптимизированные изображения

### Внутренние зависимости
- `@/contexts/CreateListingContext` - получение данных объявления
- `@/services/listings` - API для создания объявления и загрузки изображений
- `@/utils/toast` - уведомления пользователю
- `@/i18n/routing` - навигация с поддержкой локализации

## Управление состоянием
### Локальное состояние
```typescript
const [isPublishing, setIsPublishing] = useState(false);      // Статус публикации
const [isSavingDraft, setIsSavingDraft] = useState(false);    // Статус сохранения черновика
const [uploadingImages, setUploadingImages] = useState(false); // Статус загрузки изображений
```

### Интеграция с контекстом
```typescript
const { state, saveDraft, publish } = useCreateListing();
// state содержит все данные собранные в мастере
```

## Бизнес-логика
### Процесс публикации
1. **Создание объявления** через ListingsService.createListing()
2. **Загрузка изображений** (конвертация base64 → File → MinIO)
3. **Удаление черновика** из localStorage
4. **Обновление контекста** и показ success message
5. **Навигация** к странице созданного объявления

### Обработка изображений
```typescript
// Преобразование base64 в File объекты
const files = await Promise.all(
  state.images.map(async (imageUrl, index) => {
    const res = await fetch(imageUrl);
    const blob = await res.blob();
    return new File([blob], `image_${index}.jpg`, { type: 'image/jpeg' });
  })
);

await ListingsService.uploadImages(response.id, files, state.mainImageIndex);
```

### Сохранение черновика
- **localStorage persistence** через ListingsService.saveDraft()
- **Context update** через saveDraft() action
- **Toast notification** для подтверждения

## UI структура
### Preview объявления
- **Card layout**: имитация реального объявления
- **Main image**: главное фото с avatar стилем
- **Basic info**: название, цена, состояние, категория
- **Description**: обрезанное описание с line-clamp
- **Image gallery**: миниатюры остальных фото
- **Location**: город и регион
- **Trust indicators**: верификация телефона, способы оплаты

### Региональная сводка
- **Script mode**: отображение выбранного скрипта
- **Trust status**: статус верификации
- **Payment options**: доступные способы оплаты
- **Meeting type**: предпочтительный тип встречи

### Правила публикации
- **Warning alert**: региональные правила торговли
- **Community guidelines**: честность, справедливость, безопасность

### Действия
- **Back button**: возврат к предыдущему шагу
- **Save draft**: сохранение в localStorage
- **Publish**: финальная публикация с loading состояниями

## Примеры использования
```tsx
// В мастере создания объявления
<PreviewPublishStep
  onBack={() => setCurrentStep(6)}
  onComplete={() => console.log('Published')}
/>

// После успешной публикации
router.push(`/marketplace/${response.id}`);
```

## Известные особенности
### Позитивные
- ✅ Полный preview объявления как будет выглядеть в маркетплейсе
- ✅ Обработка загрузки изображений с error handling
- ✅ Региональная сводка всех настроек
- ✅ Multiple loading states для UX feedback
- ✅ Автоматическая навигация после публикации
- ✅ Draft система для сохранения прогресса

### Технический долг
- ⚠️ Console.log в production коде (строки 31-34)
- ⚠️ Hardcoded currency symbols (строки 105-111)
- ⚠️ Hardcoded condition labels (строки 114-121)
- ⚠️ setTimeout для navigation вместо callback (строка 78)
- ⚠️ onComplete prop не используется

### Потенциальные улучшения
- Добавить validation перед публикацией (required fields check)
- Реализовать retry mechanism для failed uploads
- Добавить preview в разных размерах (mobile/desktop)
- Улучшить error handling для network failures
- Добавить progress indicator для upload процесса
- Реализовать auto-save draft во время предпросмотра