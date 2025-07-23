# LocationPicker - Многоязычное геокодирование

## Описание изменений

Компонент `LocationPicker` теперь поддерживает многоязычное геокодирование. Адреса при обратном геокодировании (клик на карте) и поиске через SmartAddressInput отображаются на языке, выбранном пользователем в интерфейсе.

## Внесенные изменения

### 1. Импорт хука локали
```typescript
import { useLocale } from 'next-intl';
```

### 2. Получение текущей локали
```typescript
const locale = useLocale();
```

### 3. Использование локали в API запросах

#### Обратное геокодирование (reverseGeocode)
```typescript
const response = await fetch(
  `https://api.mapbox.com/geocoding/v5/mapbox.places/${lng},${lat}.json?` +
    `access_token=${mapboxToken}&language=${locale}&types=address,place`
);
```

#### SmartAddressInput
```typescript
<SmartAddressInput
  // ...другие пропсы
  language={locale}
/>
```

## Поддерживаемые языки

- `en` - Английский
- `ru` - Русский  
- `sr` - Сербский

## Использование

Компонент автоматически определяет текущий язык интерфейса и использует его для:
- Обратного геокодирования при клике на карту
- Поиска адресов через SmartAddressInput
- Отображения результатов геокодирования

## Файлы, где используется LocationPicker

1. `/components/create-listing/steps/LocationStep.tsx` - создание объявлений
2. `/components/storefronts/create/steps/LocationStep.tsx` - создание витрин
3. `/components/products/steps/LocationStep.tsx` - создание товаров
4. `/components/products/steps/EditLocationStep.tsx` - редактирование товаров

Все эти компоненты теперь автоматически получают многоязычную поддержку геокодирования.