# LocationStep.tsx

## Метаданные
- **Путь**: `frontend/svetu/src/components/create-listing/steps/LocationStep.tsx`
- **Роль**: UI компонент геолокации
- **Тип**: Step компонент мастера создания объявления
- **Размер**: 384 строки

## Назначение
Пятый шаг мастера создания объявления. Управляет геолокацией товара с поддержкой региональной специфики, точного позиционирования и рекомендаций безопасных мест встреч.

## Props структура
```typescript
interface LocationStepProps {
  onNext: () => void;   // Переход к следующему шагу
  onBack: () => void;   // Возврат к предыдущему шагу
}
```

## Зависимости
### Внешние зависимости
- `react` - состояние и эффекты
- `next-intl` - интернационализация

### Внутренние зависимости
- `@/contexts/CreateListingContext` - контекст создания объявления

## Управление состоянием
### Локальное состояние
```typescript
const [formData, setFormData] = useState({
  address: string;                    // Адрес (необязательно)
  city: string;                      // Город (обязательно)
  region: string;                    // Регион (обязательно)
  country: string;                   // Страна (по умолчанию 'Србија')
  exactLocation: boolean;            // Точная геолокация
  approximateArea: string;           // Приблизительная область
  safeMeetingPlaces: string[];       // Безопасные места встреч
});

const [currentLocation, setCurrentLocation] = useState<{lat, lng} | null>();
const [loadingLocation, setLoadingLocation] = useState(false);
```

### Региональные данные
```typescript
// TODO: Загружать список городов из API
// ВРЕМЕННОЕ РЕШЕНИЕ: Хардкодные города Сербии
const serbianCities = ['Београд', 'Нови Сад', 'Ниш', ...];
const serbianRegions = ['Београд (главни град)', 'Војводина', ...];
```

## Бизнес-логика
### Геолокация
- **Browser Geolocation API**: получение текущих координат
- **High accuracy mode**: enableHighAccuracy для точности
- **Timeout handling**: 10 секунд на получение координат
- **Error handling**: graceful fallback при отказе в доступе

### Валидация формы
```typescript
const canProceed = formData.city && formData.region;
```

### Региональная адаптация
- **Локальные страны**: Сербия, Хорватия, БиГ, Черногория, Словения, Македония
- **Флаги стран**: эмодзи для визуального представления
- **Локальные города**: предзаполненный список для Сербии

## UI структура
### Форма локации
1. **Страна**: Select с флагами стран региона
2. **Регион**: Select с административными областями
3. **Город**: Select + custom input для "Остало"
4. **Адрес**: Опциональное текстовое поле
5. **Точная локация**: Toggle для GPS координат

### Геолокация блок
- **Conditional rendering**: показывается при включении exactLocation
- **GPS кнопка**: с loading состоянием
- **Status indicator**: подтверждение получения координат
- **Visual feedback**: success состояние при получении

### Безопасные места встреч
- **Multi-select buttons**: toggle выбор мест
- **Regional options**: тржни центар, главни трг, etc.
- **Visual states**: primary/outline для выбранных/не выбранных

## Примеры использования
```tsx
// В мастере создания объявления
<LocationStep
  onNext={() => setCurrentStep(5)}
  onBack={() => setCurrentStep(3)}
/>

// Получение данных локации
const { state } = useCreateListing();
console.log(state.location); // { latitude, longitude, address, city, region, country }
```

## Известные особенности
### Позитивные
- ✅ Региональная адаптация для Балканских стран
- ✅ Browser Geolocation API integration
- ✅ Безопасные места встреч для доверия пользователей
- ✅ Гибкая адресация (точная или приблизительная)
- ✅ Валидация обязательных полей
- ✅ Автосохранение в глобальный контекст

### Технический долг
- ⚠️ Hardcoded списки городов и регионов (строки 33-62)
- ⚠️ Отсутствует API интеграция для городов
- ⚠️ Нет поддержки autocomplete для адресов
- ⚠️ Флаги как эмодзи вместо SVG
- ⚠️ Отсутствует карта для visual feedback

### Потенциальные улучшения
- Интегрировать API для динамической загрузки городов
- Добавить Google Maps/OpenStreetMap интеграцию
- Реализовать address autocomplete
- Добавить поддержку координат без GPS (клик по карте)
- Улучшить UX для выбора региона (nested selects)
- Добавить валидацию адресов через геосервисы
- Реализовать сохранение frequently used locations