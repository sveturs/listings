# AttributesStep.tsx

## Метаданные
- **Путь**: `frontend/svetu/src/components/create-listing/steps/AttributesStep.tsx`
- **Роль**: UI компонент динамических атрибутов
- **Тип**: Step компонент мастера создания объявления
- **Размер**: 724 строки

## Назначение
Третий шаг мастера создания объявления. Динамически загружает и отображает атрибуты выбранной категории с группировкой по типам и валидацией обязательных полей.

## Props структура
```typescript
interface AttributesStepProps {
  onNext: () => void;   // Переход к следующему шагу
  onBack: () => void;   // Возврат к предыдущему шагу
}
```

## Зависимости
### Внешние зависимости
- `react` - состояние, эффекты, callbacks
- `next-intl` - интернационализация с поддержкой локали

### Внутренние зависимости
- `@/contexts/CreateListingContext` - контекст создания объявления
- `@/services/marketplace` - API для загрузки атрибутов категории
- `@/utils/translatedAttribute` - утилита для локализации атрибутов

## Управление состоянием
### Локальное состояние
```typescript
const [attributes, setAttributes] = useState<CategoryAttributeMapping[]>([]);
const [formData, setFormData] = useState<Record<number, AttributeFormData>>();
const [loading, setLoading] = useState(true);
const [expandedGroups, setExpandedGroups] = useState<Set<string>>();
```

### Типы атрибутов
```typescript
interface AttributeFormData {
  attribute_id: number;      // ID атрибута
  attribute_name: string;    // Системное имя
  display_name: string;      // Отображаемое имя
  attribute_type: string;    // text|number|select|boolean|multiselect
  text_value?: string;       // Текстовое значение
  numeric_value?: number;    // Числовое значение
  boolean_value?: boolean;   // Булево значение
  json_value?: any;          // JSON значение (массивы)
  display_value?: string;    // Форматированное значение для отображения
  unit?: string;            // Единица измерения
}
```

## Бизнес-логика
### Загрузка атрибутов
1. **API запрос** на основе выбранной категории
2. **Преобразование данных** в CategoryAttributeMapping формат
3. **Дедупликация** по attribute_id
4. **Сортировка** по sort_order

### Группировка атрибутов
Автоматическая группировка по логическим категориям:
- **basic**: brand, model, type, category, name, title
- **technical**: year, engine, processor, ram, display, etc.
- **condition**: condition, warranty, used, new
- **accessories**: accessories, included, box, charger
- **dimensions**: width, height, length, weight, size
- **other**: все остальные атрибуты

### Валидация форм
- **Обязательные поля**: проверка заполненности required атрибутов
- **Типизированные значения**: проверка типов данных
- **Автосохранение**: синхронизация с глобальным контекстом

## UI структура
### Группы атрибутов
- **Collapsible cards** с иконками и счетчиками
- **Автораскрытие** групп с обязательными полями
- **Индикаторы прогресса** заполнения обязательных полей
- **Adaptive expand/collapse** с анимацией

### Типы полей
1. **text**: Обычное текстовое поле
2. **number**: Числовое поле с указанием шага
3. **select**: Выпадающий список с переводами опций
4. **boolean**: Checkbox с лейблом "Да"
5. **multiselect**: Группа checkbox для множественного выбора

### Интерактивные элементы
- **Группы атрибутов**: clickable headers для expand/collapse
- **Progress badges**: показывают статус заполнения групп
- **Required indicators**: звездочки для обязательных полей

## Примеры использования
```tsx
// В мастере создания объявления
<AttributesStep
  onNext={() => setCurrentStep(3)}
  onBack={() => setCurrentStep(1)}
/>

// Получение атрибутов
const { state } = useCreateListing();
console.log(state.attributes); // Record<number, AttributeFormData>
```

## Известные особенности
### Позитивные
- ✅ Динамическая загрузка атрибутов на основе категории
- ✅ Автоматическая группировка по логическим типам
- ✅ Поддержка всех типов полей (text, number, select, boolean, multiselect)
- ✅ Полная локализация атрибутов и опций
- ✅ Валидация обязательных полей в реальном времени
- ✅ UX оптимизация: автораскрытие важных групп

### Технический долг
- ⚠️ Extensive console.log в production коде (строки 88-98)
- ⚠️ Hardcoded группировка атрибутов по именам (строки 184-227)
- ⚠️ Хардкодная локализация "Да"/"Не" для boolean (строка 314)
- ⚠️ Отсутствует обработка ошибок загрузки атрибутов
- ⚠️ Нет поддержки условной видимости атрибутов

### Потенциальные улучшения
- Вынести логику группировки в отдельную утилиту
- Добавить поддержку условных атрибутов (зависимости)
- Реализовать кэширование атрибутов категорий
- Добавить валидацию диапазонов для числовых полей
- Улучшить UX для long lists атрибутов (поиск/фильтрация)
- Добавить подсказки и описания для сложных атрибутов