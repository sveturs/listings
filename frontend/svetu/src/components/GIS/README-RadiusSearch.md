# Компоненты радиусного поиска

Данная документация описывает компоненты для радиусного поиска, созданные для маркетплейса.

## Компоненты

### 1. RadiusSearchControl.tsx

Главный компонент управления радиусным поиском с адаптивным интерфейсом.

**Возможности:**
- Слайдер для выбора радиуса поиска (0.1км - 50км)
- Поле ввода адреса с автодополнением
- Кнопка "Использовать мое местоположение"
- Переключатель показа области поиска на карте
- Адаптивный дизайн для мобильных устройств

**Пример использования:**
```tsx
<RadiusSearchControl
  config={{
    minRadius: 0.1,
    maxRadius: 50,
    defaultRadius: 5,
    showMyLocation: true,
    showAddressInput: true,
  }}
  onSearch={(results, center, radius) => {
    // Обработка результатов поиска
  }}
  onRadiusChange={(radius) => {
    // Обработка изменения радиуса
  }}
/>
```

### 2. useRadiusSearch.ts

Хук для работы с радиусным поиском и управления состоянием.

**API функции:**
- `search(params)` - основной поиск по параметрам
- `searchByAddress(address, radius, category?)` - поиск по адресу
- `searchByCurrentLocation(radius, category?)` - поиск по текущему местоположению
- `clearResults()` - очистка результатов

**Утилитные функции:**
- `formatRadius(radius)` - форматирование радиуса
- `validateRadius(radius, min, max)` - валидация радиуса
- `normalizeRadius(radius, min, max)` - нормализация радиуса
- `createRadiusCircle(center, radius)` - создание геоданных круга

### 3. MarketplaceMapWithRadiusSearch.tsx

Интегрированный компонент карты с радиусным поиском для маркетплейса.

**Возможности:**
- Интерактивная карта с маркерами
- Интегрированный контроль радиусного поиска
- Отображение результатов поиска
- Обработка кликов по маркерам
- Индикаторы загрузки и ошибок

### 4. HomePage.tsx

Обновленная главная страница с переключением между списком и картой.

**Новые возможности:**
- Переключатель "Список/Карта"
- Интеграция радиусного поиска
- Отображение выбранного объявления
- Адаптивный дизайн

## Типы данных

### RadiusSearchParams
```typescript
interface RadiusSearchParams {
  latitude: number;
  longitude: number;
  radius: number; // в километрах
  category?: string;
  limit?: number;
  query?: string;
}
```

### RadiusSearchResult
```typescript
interface RadiusSearchResult {
  id: string;
  title: string;
  description?: string;
  latitude: number;
  longitude: number;
  distance: number; // расстояние в км
  category?: string;
  price?: number;
  currency?: string;
  imageUrl?: string;
  metadata?: Record<string, any>;
}
```

### RadiusSearchControlConfig
```typescript
interface RadiusSearchControlConfig {
  minRadius: number; // в км
  maxRadius: number; // в км
  defaultRadius: number; // в км
  step: number; // шаг изменения радиуса
  showMyLocation: boolean;
  showAddressInput: boolean;
  showRadiusCircle: boolean;
  enableGeolocation: boolean;
}
```

## API Endpoints

Компоненты работают с API endpoint:
- `POST /api/v1/gis/search/radius` - радиусный поиск

Ожидаемый формат ответа:
```json
{
  "data": {
    "items": [RadiusSearchResult[]],
    "total": number,
    "center": {
      "latitude": number,
      "longitude": number
    },
    "radius": number
  }
}
```

## Переводы

Добавлены переводы в `src/messages/`:
- `gis.radius_search.title` - "Радиусный поиск"
- `gis.radius_search.radius` - "Радиус поиска"
- `gis.radius_search.address` - "Адрес"
- `gis.radius_search.use_my_location` - "Использовать мое местоположение"
- И другие ключи для полной локализации

## Стили

Компоненты используют существующие стили DaisyUI:
- `btn`, `btn-primary`, `btn-outline` для кнопок
- `input`, `input-bordered` для полей ввода
- `range`, `range-primary` для слайдеров
- `alert`, `alert-error` для уведомлений
- Адаптивные классы для мобильных устройств

## Мобильная адаптация

Все компоненты адаптированы для мобильных устройств:
- Компактные контролы на маленьких экранах
- Адаптивная сетка и отступы
- Оптимизированные размеры кнопок для touch-интерфейса
- Скрываемые/раскрываемые панели для экономии места

## Интеграция

Компоненты интегрированы в:
1. **Главную страницу** (`/`) - переключатель список/карта
2. **Систему типов** - новые TypeScript интерфейсы
3. **Систему переводов** - многоязычная поддержка
4. **Систему экспортов** - добавлены в `GIS/index.ts`

## Требования

- React 19+
- TypeScript
- TailwindCSS + DaisyUI
- next-intl для переводов
- react-map-gl для карт
- Переменная окружения `NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN`