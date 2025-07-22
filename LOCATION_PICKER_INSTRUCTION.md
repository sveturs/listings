# Инструкция по реализации выбора адреса с помощью LocationPicker

## Обзор

Компонент `LocationPicker` предоставляет унифицированный интерфейс для выбора местоположения во всей системе. Он поддерживает два режима работы:
1. **Поиск по адресу** - ввод адреса с автодополнением
2. **Выбор на карте** - клик по карте для установки маркера

## Компонент LocationPicker

### Расположение
```
/frontend/svetu/src/components/GIS/LocationPicker.tsx
```

### Основные возможности
- 🔍 Поиск адреса с автодополнением (использует SmartAddressInput)
- 🗺️ Интерактивная карта на базе Mapbox
- 📍 Визуальный маркер местоположения
- 🔄 Обратное геокодирование (получение адреса по координатам)
- 📱 Адаптивный дизайн
- 🌍 Поддержка международных адресов

### Интерфейс компонента

```typescript
interface LocationData {
  latitude: number;
  longitude: number;
  address: string;
  city: string;
  region: string;
  country: string;
  confidence: number;
}

interface LocationPickerProps {
  value?: LocationData;
  onChange: (location: LocationData) => void;
  placeholder?: string;
  height?: string;
  showCurrentLocation?: boolean;
  defaultCountry?: string;
  mapboxToken?: string;
}
```

## Примеры использования

### 1. Создание обычных объявлений (marketplace listings)

```typescript
// src/components/create-listing/steps/LocationStep.tsx

import LocationPicker from '@/components/GIS/LocationPicker';

export default function LocationStep() {
  const [location, setLocation] = useState<LocationData | undefined>();

  return (
    <LocationPicker
      value={location}
      onChange={setLocation}
      placeholder="Начните вводить адрес (например: Београд, Кнез Михаилова)"
      height="500px"
      showCurrentLocation={true}
      defaultCountry="Србија"
    />
  );
}
```

### 2. Создание витрин (storefronts)

```typescript
// src/components/storefronts/create/steps/LocationStep.tsx

import LocationPicker from '@/components/GIS/LocationPicker';

export default function LocationStep() {
  const [location, setLocation] = useState<LocationData | undefined>();

  return (
    <LocationPicker
      value={location}
      onChange={handleLocationChange}
      placeholder="Введите адрес вашей витрины или выберите точку на карте"
      height="400px"
      showCurrentLocation={false}
      defaultCountry="Србија"
    />
  );
}
```

### 3. Создание товаров для витрин (storefront products)

Для добавления выбора адреса при создании товара витрины, нужно:

1. **Добавить шаг LocationStep в wizard создания товара**:

```typescript
// src/components/products/ProductWizard.tsx

import LocationStep from './steps/LocationStep';

const steps = [
  'category',
  'basic',
  'attributes',
  'photos',
  'location', // Новый шаг
  'preview'
];
```

2. **Создать компонент LocationStep**:

```typescript
// src/components/products/steps/LocationStep.tsx

'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import LocationPicker from '@/components/GIS/LocationPicker';

interface LocationStepProps {
  onNext: () => void;
  onBack: () => void;
  formData: any;
  updateFormData: (data: any) => void;
}

interface LocationData {
  latitude: number;
  longitude: number;
  address: string;
  city: string;
  region: string;
  country: string;
  confidence: number;
}

export default function LocationStep({ 
  onNext, 
  onBack, 
  formData, 
  updateFormData 
}: LocationStepProps) {
  const t = useTranslations();
  const [location, setLocation] = useState<LocationData | undefined>(
    formData.location
  );
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (location) {
      updateFormData({ location });
    }
  }, [location, updateFormData]);

  const handleLocationChange = (locationData: LocationData) => {
    setLocation(locationData);
    setErrors({}); // Очищаем ошибки при выборе
  };

  const validate = () => {
    const newErrors: Record<string, string> = {};
    
    if (!location) {
      newErrors.location = 'Необходимо выбрать местоположение товара';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleNext = () => {
    if (validate()) {
      onNext();
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold mb-2">
          📍 Местоположение товара
        </h2>
        <p className="text-base-content/70">
          Укажите, где находится товар для удобства покупателей
        </p>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <LocationPicker
            value={location}
            onChange={handleLocationChange}
            placeholder="Введите адрес местонахождения товара"
            height="400px"
            showCurrentLocation={true}
            defaultCountry="Србија"
          />
          
          {errors.location && (
            <p className="text-error text-sm mt-2">{errors.location}</p>
          )}

          {location && (
            <div className="mt-4 p-4 bg-success/10 border border-success/20 rounded-lg">
              <h4 className="font-medium text-success-content mb-2">
                ✅ Местоположение выбрано
              </h4>
              <div className="text-sm text-success-content/80">
                <p><strong>Адрес:</strong> {location.address}</p>
                <p><strong>Город:</strong> {location.city}</p>
                <p><strong>Координаты:</strong> {location.latitude.toFixed(6)}, {location.longitude.toFixed(6)}</p>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="flex justify-between">
        <button className="btn btn-outline" onClick={onBack}>
          ← {t('common.back')}
        </button>
        <button 
          className={`btn btn-primary ${!location ? 'btn-disabled' : ''}`}
          onClick={handleNext}
          disabled={!location}
        >
          {t('common.next')} →
        </button>
      </div>
    </div>
  );
}
```

3. **Обновить контекст для хранения данных о местоположении**:

```typescript
// src/contexts/CreateProductContext.tsx

interface ProductFormData {
  // ... существующие поля
  location?: {
    latitude: number;
    longitude: number;
    address: string;
    city: string;
    region: string;
    country: string;
  };
}
```

## Ключевые особенности реализации

### 1. Переключение режимов
Компонент предоставляет два режима с плавным переключением между ними:
- Поиск по адресу (по умолчанию)
- Выбор на карте

### 2. Обратное геокодирование
При клике на карте автоматически происходит:
1. Установка маркера для мгновенной обратной связи
2. Запрос к Mapbox API для получения адреса
3. Обновление полей адреса с возможностью редактирования

### 3. Валидация
- Проверка обязательности выбора местоположения
- Возможность редактирования адреса и города после выбора
- Сохранение уровня достоверности (confidence)

### 4. Визуальная обратная связь
- Красный маркер на карте
- Индикатор загрузки при геокодировании
- Блок с подтверждением выбранного местоположения

## Настройка окружения

### Mapbox Token
Для работы карты необходим токен Mapbox в `.env`:
```
NEXT_PUBLIC_MAPBOX_TOKEN=pk.eyJ1I...
```

### Стили карты
По умолчанию используется `mapbox://styles/mapbox/streets-v12`

## Полезные компоненты

### SmartAddressInput
Используется внутри LocationPicker для поиска адресов:
```
/frontend/svetu/src/components/GIS/SmartAddressInput.tsx
```

### LocationPrivacySettings
Дополнительный компонент для настройки приватности местоположения:
```
/frontend/svetu/src/components/GIS/LocationPrivacySettings.tsx
```

## Проверка работы

### Визуальная проверка
1. Переключение между режимами работает корректно
2. Маркер появляется при клике на карте
3. Адрес автоматически заполняется
4. Поиск адреса показывает подсказки

### Функциональная проверка
1. Данные сохраняются в правильном формате
2. Валидация работает корректно
3. Координаты точные (6 знаков после запятой)
4. Город и адрес можно редактировать

## Типичные проблемы и решения

### Маркер не отображается
1. Проверьте импорт стилей: `import 'mapbox-gl/dist/mapbox-gl.css'`
2. Убедитесь, что используется правильный anchor: `anchor="bottom"`
3. Проверьте, что selectedLocation не null

### Карта не загружается
1. Проверьте токен Mapbox в `.env`
2. Убедитесь, что компонент обернут в `'use client'`
3. Проверьте сетевое соединение

### Адрес не определяется
1. Проверьте работу Mapbox Geocoding API
2. Убедитесь, что координаты передаются в правильном порядке (lng, lat)
3. Проверьте language параметр в запросе

## Заключение

LocationPicker обеспечивает единообразный и удобный способ выбора местоположения во всей системе. Компонент легко интегрируется в любые формы создания контента и обеспечивает отличный пользовательский опыт.