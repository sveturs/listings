# Спецификация: Определение и выбор адреса при AI создании продукта витрины

**Дата создания:** 2025-10-03
**Версия:** 1.0
**Статус:** ✅ Реализовано

---

## 📋 Оглавление

1. [Обзор](#обзор)
2. [Архитектура решения](#архитектура-решения)
3. [Этап 1: AI Обработка (ProcessView)](#этап-1-ai-обработка-processview)
4. [Этап 2: Улучшение (EnhanceView)](#этап-2-улучшение-enhanceview)
5. [Управление состоянием (Context)](#управление-состоянием-context)
6. [Структуры данных](#структуры-данных)
7. [Интернационализация](#интернационализация)
8. [Технические детали](#технические-детали)
9. [Диаграмма потока данных](#диаграмма-потока-данных)

---

## Обзор

### Цель функционала

Предоставить пользователям возможность автоматического определения адреса продукта из метаданных фотографий (EXIF GPS) и гибкого выбора адреса при создании продукта витрины с помощью AI.

### Ключевые возможности

1. ✅ **Автоматическое определение адреса из EXIF** - извлечение GPS координат из метаданных фотографий
2. ✅ **Геокодирование** - преобразование координат в человеко-читаемый адрес
3. ✅ **Три варианта выбора адреса**:
   - Использовать адрес витрины (по умолчанию)
   - Использовать адрес из EXIF (если определён)
   - Указать индивидуальный адрес вручную на карте
4. ✅ **Настройки приватности адреса** - уровень детализации и показ на карте
5. ✅ **Прозрачный процесс** - пользователь видит статус каждого шага обработки

### Путь пользователя

```
Загрузка фото → AI Обработка → Улучшение → Варианты → Публикация
                      ↓              ↓
              Извлечение EXIF   Выбор адреса
              GPS координат     + Приватность
```

---

## Архитектура решения

### Компоненты системы

```
┌─────────────────────────────────────────────────────────────┐
│                      AI Product Wizard                       │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ UploadView   │→ │ ProcessView  │→ │ EnhanceView  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                           ↓                   ↓               │
│                    EXIF Extraction      Location Picker      │
│                    Geocoding            Privacy Settings     │
│                                                               │
└─────────────────────────────────────────────────────────────┘
                            ↓
              ┌─────────────────────────────┐
              │  CreateAIProductContext     │
              │  - location data            │
              │  - privacy settings         │
              │  - useStorefrontLocation    │
              └─────────────────────────────┘
```

### Стек технологий

- **Frontend Framework:** React 19, Next.js 15
- **State Management:** React Context API
- **EXIF Processing:** exifr library
- **Geocoding:** Nominatim OpenStreetMap API
- **Map Components:** Leaflet + React-Leaflet
- **i18n:** next-intl
- **Типизация:** TypeScript

---

## Этап 1: AI Обработка (ProcessView)

### Файл: `src/components/storefronts/ai/ProcessView.tsx`

### Добавленный шаг: Определение адреса из EXIF

#### Позиция в процессе

Шаг 2 из 5 (между "Анализ изображений" и "Определение категории"):

```typescript
const [steps, setSteps] = useState<ProcessStep[]>([
  { id: 'analyze', label: t('analyzingImages'), status: 'pending' },
  {
    id: 'exif',
    label: t('detectingLocationFromExif'),
    status: 'pending'
  }, // ← НОВЫЙ ШАГ
  { id: 'category', label: t('detectingCategory'), status: 'pending' },
  { id: 'titles', label: t('generatingTitleVariants'), status: 'pending' },
  { id: 'translations', label: t('creatingTranslations'), status: 'pending' },
]);
```

#### Алгоритм обработки EXIF

```typescript
// Step 2: Extract location from EXIF
updateStepStatus('exif', 'processing');

let locationData = null;
try {
  // 1. Извлечение GPS координат из всех загруженных изображений
  const exifLocation = await extractLocationFromImages(state.imageFiles);

  if (exifLocation) {
    console.log('[ProcessView] EXIF location found:', exifLocation);

    // 2. Геокодирование координат в адрес (используя текущую локаль)
    const geocodedAddress = await geocoding.reverseGeocode(
      exifLocation.latitude,
      exifLocation.longitude
    );

    console.log('[ProcessView] Geocoded address:', geocodedAddress);

    if (geocodedAddress) {
      // 3. Формирование объекта locationData
      locationData = {
        latitude: exifLocation.latitude,
        longitude: exifLocation.longitude,
        address: geocodedAddress.display_name || '',
        city: geocodedAddress.address?.city || '',
        region: geocodedAddress.address?.state || '',
        source: 'exif' as const, // Важно: маркер источника
      };

      // 4. Обновление статуса с информативным сообщением
      updateStepStatus(
        'exif',
        'completed',
        `${geocodedAddress.address?.city || 'Location'} detected`
      );
    } else {
      updateStepStatus('exif', 'completed', 'No address found');
    }
  } else {
    updateStepStatus('exif', 'completed', 'No GPS data in photos');
  }
} catch (exifError) {
  console.error('[ProcessView] EXIF extraction error:', exifError);
  updateStepStatus('exif', 'completed', 'No location data available');
}
```

#### Интеграция с результатами AI

```typescript
// Сохранение locationData в контекст вместе с другими AI данными
setAIData({
  title: analysisResult.title,
  titleVariants: analysisResult.titleVariants,
  // ... другие поля ...
  location: locationData || analysisResult.location || null, // Приоритет EXIF
});
```

#### Возможные статусы шага

| Статус | Иконка | Сообщение | Описание |
|--------|--------|-----------|----------|
| `pending` | ⚪ Пустой круг | - | Ожидание выполнения |
| `processing` | 🔄 Спиннер | "Обработка..." | Идёт извлечение GPS и геокодирование |
| `completed` ✅ | ✓ Галочка | "Новый Сад detected" | Адрес успешно определён |
| `completed` ℹ️ | ✓ Галочка | "No GPS data in photos" | GPS данные отсутствуют в EXIF |
| `completed` ℹ️ | ✓ Галочка | "No address found" | Координаты есть, но геокодирование не удалось |
| `error` | ❌ Крестик | "No location data available" | Ошибка при обработке |

#### Зависимости

```typescript
import { extractLocationFromImages } from '@/utils/exifUtils';
import { useAddressGeocoding } from '@/hooks/useAddressGeocoding';
```

**Хук геокодирования:**

```typescript
const geocoding = useAddressGeocoding({
  country: ['rs'], // Сербия
  language: locale, // Текущий язык интерфейса (ru/en/sr)
});
```

---

## Этап 2: Улучшение (EnhanceView)

### Файл: `src/components/storefronts/ai/EnhanceView.tsx`

### Добавленная секция: Выбор адреса

#### UI компоненты

```
┌─────────────────────────────────────────────────────────┐
│  📍 Местоположение товара                               │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ⦿ Использовать адрес витрины                          │
│    Адрес будет взят из настроек вашей витрины          │
│                                                          │
│  ○ Использовать адрес из метаданных фото               │  ← Показывается только
│    ул. Boulevard of Liberation 12, Новый Сад           │    если EXIF определён
│    📍 Новый Сад, Vojvodina                             │
│                                                          │
│  ○ Указать индивидуальный адрес                        │
│    Выбрать адрес на карте вручную                      │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

#### Логика отображения вариантов

```typescript
// 1. Адрес витрины - всегда доступен
<label className="card">
  <input
    type="radio"
    checked={localUseStorefrontLocation}
    onChange={() => handleLocationTypeChange('storefront')}
  />
  <div>Использовать адрес витрины</div>
</label>

// 2. Адрес из EXIF - условный рендеринг
{state.aiData.location?.source === 'exif' && (
  <label className="card">
    <input
      type="radio"
      checked={useExifLocation}
      onChange={() => handleLocationTypeChange('exif')}
    />
    <div>
      Использовать адрес из метаданных фото
      <p>{state.aiData.location?.address}</p>
      <div className="text-info">
        📍 {state.aiData.location?.city}, {state.aiData.location?.region}
      </div>
    </div>
  </label>
)}

// 3. Индивидуальный адрес - всегда доступен
<label className="card">
  <input
    type="radio"
    checked={!localUseStorefrontLocation && !useExifLocation}
    onChange={() => handleLocationTypeChange('manual')}
  />
  <div>Указать индивидуальный адрес</div>
</label>
```

#### LocationPicker (для индивидуального адреса)

```typescript
{!localUseStorefrontLocation && !useExifLocation && (
  <div className="mt-4 p-4 bg-base-200 rounded-lg">
    <LocationPicker
      value={individualLocation}
      onChange={handleLocationChange}
      placeholder={t('locationPlaceholder')}
      height="300px"
      showCurrentLocation={true}
      defaultCountry="Србија"
    />
  </div>
)}
```

#### Настройки приватности

```
┌─────────────────────────────────────────────────────────┐
│  🛡️ Настройки приватности                  [▼]         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Уровень приватности:                                   │
│  ⦿ Точный адрес      (ул. Пушкина, д. 12)              │
│  ○ Улица             (ул. Пушкина)                      │
│  ○ Район             (Центральный район)                │
│  ○ Город             (Новый Сад)                        │
│                                                          │
│  ☑ Показывать на карте                                 │
│    Разрешить показывать товар на карте маркетплейса    │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

**Компонент приватности:**

```typescript
<LocationPrivacySettings
  selectedLevel={localPrivacyLevel}
  onLevelChange={setLocalPrivacyLevel}
  location={{
    lat: individualLocation.latitude,
    lng: individualLocation.longitude,
  }}
  showPreview={true}
/>
```

#### Управление состоянием

```typescript
// Локальное состояние компонента
const [localUseStorefrontLocation, setLocalUseStorefrontLocation] =
  useState(state.useStorefrontLocation ?? true);
const [useExifLocation, setUseExifLocation] =
  useState(state.aiData.location?.source === 'exif');
const [individualLocation, setIndividualLocation] = useState<LocationData | null>(
  state.aiData.location && state.aiData.location.source !== 'storefront'
    ? { /* ... данные из context ... */ }
    : null
);
const [localPrivacyLevel, setLocalPrivacyLevel] =
  useState<'exact' | 'street' | 'district' | 'city'>(
    state.locationPrivacyLevel || 'exact'
  );
const [localShowOnMap, setLocalShowOnMap] =
  useState(state.showOnMap ?? true);
```

#### Обработчик смены типа адреса

```typescript
const handleLocationTypeChange = (type: 'storefront' | 'exif' | 'manual') => {
  if (type === 'storefront') {
    setLocalUseStorefrontLocation(true);
    setUseStorefrontLocation(true); // → Context
    setUseExifLocation(false);
  } else if (type === 'exif') {
    setLocalUseStorefrontLocation(false);
    setUseStorefrontLocation(false); // → Context
    setUseExifLocation(true);
  } else {
    setLocalUseStorefrontLocation(false);
    setUseStorefrontLocation(false); // → Context
    setUseExifLocation(false);
    // Открывается LocationPicker
  }
};
```

#### Сохранение данных

```typescript
const handleSave = () => {
  // 1. Основные данные продукта
  setAIData({
    title: editedData.title,
    description: editedData.description,
    price: editedData.price,
    stockQuantity: editedData.stockQuantity,
    categoryId: editedData.categoryId,
    category: editedData.category,

    // 2. Обновление location в зависимости от выбора
    location: localUseStorefrontLocation
      ? null // Витрина использует свой адрес
      : useExifLocation && state.aiData.location
        ? { ...state.aiData.location, source: 'exif' as const }
        : individualLocation
          ? { ...individualLocation, source: 'manual' as const }
          : state.aiData.location,
  });

  // 3. Настройки приватности
  setLocationPrivacyLevel(localPrivacyLevel);
  setShowOnMap(localShowOnMap);

  setView('variants'); // → Следующий шаг
};
```

---

## Управление состоянием (Context)

### Файл: `src/contexts/CreateAIProductContext.tsx`

### Расширение интерфейса AIProductState

```typescript
interface AIProductState {
  // ... существующие поля ...

  aiData: {
    // ... другие поля ...

    location: {
      address: string;
      latitude: number;
      longitude: number;
      city: string;
      region: string;
      source: 'exif' | 'storefront' | 'manual'; // ← Источник адреса
    } | null;
  };

  // Настройки локации
  useStorefrontLocation: boolean; // ← Использовать адрес витрины
  locationPrivacyLevel: 'exact' | 'street' | 'district' | 'city'; // ← Расширен
  showOnMap: boolean; // ← Показывать на карте
  previewLanguage: string;
}
```

### Новые экшены

```typescript
type AIProductAction =
  | { type: 'SET_VIEW'; payload: View }
  | { type: 'SET_PROCESSING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_IMAGES'; payload: { urls: string[]; files: File[] } }
  | { type: 'SET_AI_DATA'; payload: Partial<AIProductState['aiData']> }
  | { type: 'SELECT_TITLE_VARIANT'; payload: number }
  | { type: 'SET_USE_STOREFRONT_LOCATION'; payload: boolean } // ← Новый
  | {
      type: 'SET_LOCATION_PRIVACY_LEVEL'; // ← Новый
      payload: 'exact' | 'street' | 'district' | 'city';
    }
  | { type: 'SET_SHOW_ON_MAP'; payload: boolean } // ← Новый
  | { type: 'SET_HAS_VARIANTS'; payload: boolean }
  // ... остальные экшены ...
```

### Новые методы контекста

```typescript
interface CreateAIProductContextType {
  state: AIProductState;
  dispatch: React.Dispatch<AIProductAction>;

  // ... существующие методы ...

  // Location methods ← НОВЫЕ
  setUseStorefrontLocation: (use: boolean) => void;
  setLocationPrivacyLevel: (
    level: 'exact' | 'street' | 'district' | 'city'
  ) => void;
  setShowOnMap: (show: boolean) => void;

  // ... остальные методы ...
}
```

### Реализация методов

```typescript
// Location methods
const setUseStorefrontLocation = (use: boolean) =>
  dispatch({ type: 'SET_USE_STOREFRONT_LOCATION', payload: use });

const setLocationPrivacyLevel = (
  level: 'exact' | 'street' | 'district' | 'city'
) => dispatch({ type: 'SET_LOCATION_PRIVACY_LEVEL', payload: level });

const setShowOnMap = (show: boolean) =>
  dispatch({ type: 'SET_SHOW_ON_MAP', payload: show });
```

### Reducer handlers

```typescript
case 'SET_USE_STOREFRONT_LOCATION':
  return { ...state, useStorefrontLocation: action.payload };

case 'SET_LOCATION_PRIVACY_LEVEL':
  return { ...state, locationPrivacyLevel: action.payload };

case 'SET_SHOW_ON_MAP':
  return { ...state, showOnMap: action.payload };
```

### Начальное состояние

```typescript
const initialState: AIProductState = {
  currentView: 'upload',
  isProcessing: false,
  error: null,
  images: [],
  imageFiles: [],
  aiData: {
    // ... другие поля ...
    location: null,
  },
  useStorefrontLocation: true, // ← По умолчанию используем адрес витрины
  locationPrivacyLevel: 'exact', // ← По умолчанию точный адрес
  showOnMap: true, // ← По умолчанию показываем на карте
  previewLanguage: 'ru',
};
```

---

## Структуры данных

### LocationData (EXIF source)

```typescript
interface LocationDataFromEXIF {
  latitude: number;        // GPS широта из EXIF
  longitude: number;       // GPS долгота из EXIF
  address: string;         // "ул. Освобождения 12, Новый Сад, Сербия"
  city: string;            // "Новый Сад"
  region: string;          // "Vojvodina"
  source: 'exif';          // Маркер источника
}
```

**Пример данных:**

```json
{
  "latitude": 45.2671,
  "longitude": 19.8335,
  "address": "Boulevard of Liberation 12, Novi Sad, Serbia",
  "city": "Novi Sad",
  "region": "Vojvodina",
  "source": "exif"
}
```

### LocationData (Manual source)

```typescript
interface LocationDataFromPicker {
  latitude: number;        // Координаты из LocationPicker
  longitude: number;
  address: string;         // Полный адрес из геокодирования
  city: string;
  region: string;
  country: string;         // "Србија"
  confidence: number;      // 0.9 (точность определения)
  source: 'manual';        // Маркер источника
}
```

### ProcessStep (AI обработка)

```typescript
interface ProcessStep {
  id: string;              // 'analyze' | 'exif' | 'category' | 'titles' | 'translations'
  label: string;           // Локализованное название
  status: 'pending' | 'processing' | 'completed' | 'error';
  message?: string;        // Дополнительное сообщение о статусе
}
```

### LocationPrivacyLevel

```typescript
type LocationPrivacyLevel = 'exact' | 'street' | 'district' | 'city';
```

**Примеры преобразования адреса:**

| Уровень | Исходный адрес | Отображаемый адрес |
|---------|----------------|-------------------|
| `exact` | ул. Пушкина, д. 12, Новый Сад | ул. Пушкина, д. 12, Новый Сад |
| `street` | ул. Пушкина, д. 12, Новый Сад | ул. Пушкина, Новый Сад |
| `district` | ул. Пушкина, д. 12, Новый Сад | Центральный район, Новый Сад |
| `city` | ул. Пушкина, д. 12, Новый Сад | Новый Сад |

---

## Интернационализация

### Добавленные ключи переводов

#### `storefronts.json`

##### Шаг AI обработки

| Ключ | Русский | English | Српски |
|------|---------|---------|--------|
| `detectingLocationFromExif` | Определение адреса из EXIF данных фотографии | Detecting location from EXIF photo metadata | Detekcija lokacije iz EXIF metapodataka fotografije |

##### Выбор типа адреса

| Ключ | Русский | English | Српски |
|------|---------|---------|--------|
| `useStorefrontLocation` | Использовать адрес витрины | Use storefront address | Koristi adresu prodavnice |
| `useExifLocation` | Использовать адрес из метаданных фото | Use address from photo metadata | Koristi adresu iz metapodataka fotografije |
| `useIndividualLocation` | Указать индивидуальный адрес | Specify individual address | Navedi posebnu adresu |

##### Настройки приватности

| Ключ | Русский | English | Српски |
|------|---------|---------|--------|
| `productLocation` | Местоположение товара | Product Location | Lokacija proizvoda |
| `privacySettings` | Настройки приватности | Privacy Settings | Podešavanja privatnosti |
| `showOnMap` | Показывать на карте | Show on map | Prikaži na mapi |
| `showOnMapDescription` | Товар будет отображаться на карте маркетплейса | Product will be displayed on marketplace map | Proizvod će biti prikazan na mapi tržišta |

### Файлы переводов

```
frontend/svetu/src/messages/
├── ru/storefronts.json  ← Русский
├── en/storefronts.json  ← Английский
└── sr/storefronts.json  ← Сербский
```

---

## Технические детали

### Утилиты и хуки

#### `extractLocationFromImages` (EXIF утилита)

**Файл:** `src/utils/exifUtils.ts`

```typescript
export async function extractLocationFromImages(
  files: File[]
): Promise<{ latitude: number; longitude: number } | null> {
  // Библиотека: exifr
  // 1. Перебирает все файлы
  // 2. Извлекает GPS теги (GPSLatitude, GPSLongitude)
  // 3. Конвертирует координаты в десятичный формат
  // 4. Возвращает координаты первого изображения с GPS данными
  // 5. Возвращает null если GPS данных нет
}
```

**Поддерживаемые EXIF теги:**

- `GPSLatitude`
- `GPSLongitude`
- `GPSLatitudeRef` (N/S)
- `GPSLongitudeRef` (E/W)

#### `useAddressGeocoding` (геокодирование)

**Файл:** `src/hooks/useAddressGeocoding.ts`

```typescript
export function useAddressGeocoding(options: {
  country?: string[];
  language?: string;
}) {
  return {
    // Прямое геокодирование (адрес → координаты)
    geocode: (address: string) => Promise<GeocodingResult[]>,

    // Обратное геокодирование (координаты → адрес)
    reverseGeocode: (lat: number, lon: number) => Promise<GeocodingResult>,
  };
}
```

**API:** Nominatim OpenStreetMap
**Endpoint:** `https://nominatim.openstreetmap.org/reverse`

**Пример запроса:**

```
GET https://nominatim.openstreetmap.org/reverse
  ?lat=45.2671
  &lon=19.8335
  &format=json
  &addressdetails=1
  &accept-language=ru
```

**Пример ответа:**

```json
{
  "display_name": "Boulevard of Liberation 12, Novi Sad, Vojvodina, Serbia",
  "address": {
    "road": "Boulevard of Liberation",
    "house_number": "12",
    "city": "Novi Sad",
    "state": "Vojvodina",
    "country": "Serbia",
    "country_code": "rs"
  },
  "lat": "45.2671",
  "lon": "19.8335"
}
```

### Компоненты карты

#### LocationPicker

**Файл:** `src/components/GIS/LocationPicker.tsx`

```typescript
interface LocationPickerProps {
  value: LocationData | null;
  onChange: (location: LocationData) => void;
  placeholder?: string;
  height?: string;
  showCurrentLocation?: boolean;
  defaultCountry?: string;
}
```

**Функционал:**

- Интерактивная карта (Leaflet)
- Поиск адреса (автодополнение)
- Клик на карте для выбора точки
- Определение текущего местоположения пользователя
- Автоматическое геокодирование выбранной точки

#### LocationPrivacySettings

**Файл:** `src/components/GIS/LocationPrivacySettings.tsx`

```typescript
interface LocationPrivacySettingsProps {
  selectedLevel: 'exact' | 'street' | 'district' | 'city';
  onLevelChange: (level: LocationPrivacyLevel) => void;
  location: { lat: number; lng: number };
  showPreview?: boolean;
}
```

**Функционал:**

- Радио-кнопки для выбора уровня приватности
- Превью адреса для каждого уровня
- Визуализация зоны на мини-карте (опционально)

### Библиотеки

```json
{
  "dependencies": {
    "exifr": "^7.1.3",           // EXIF данные
    "leaflet": "^1.9.4",         // Карты
    "react-leaflet": "^4.2.1",   // React обёртка для Leaflet
    "next-intl": "^3.x",         // Интернационализация
    "react": "^19.x",            // React
    "next": "^15.x"              // Next.js
  }
}
```

---

## Диаграмма потока данных

### Полный цикл определения и выбора адреса

```
┌─────────────────────────────────────────────────────────────────┐
│                     1. UPLOAD VIEW                               │
│  Пользователь загружает фотографии                              │
│  state.imageFiles = [File, File, ...]                           │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     2. PROCESS VIEW                              │
├─────────────────────────────────────────────────────────────────┤
│  Шаг 1: Анализ изображений                                      │
│    ✓ AI анализирует первое фото                                 │
│                                                                  │
│  Шаг 2: Определение адреса из EXIF ← НОВЫЙ ШАГ                 │
│    ┌────────────────────────────────────────┐                  │
│    │ extractLocationFromImages(imageFiles)  │                  │
│    │   ↓                                    │                  │
│    │ { lat: 45.2671, lon: 19.8335 }        │                  │
│    │   ↓                                    │                  │
│    │ geocoding.reverseGeocode(lat, lon)    │                  │
│    │   ↓                                    │                  │
│    │ {                                      │                  │
│    │   address: "Blvd Liberation 12...",   │                  │
│    │   city: "Novi Sad",                   │                  │
│    │   region: "Vojvodina",                │                  │
│    │   source: 'exif'                      │                  │
│    │ }                                      │                  │
│    └────────────────────────────────────────┘                  │
│    ✓ "Novi Sad detected"                                        │
│                                                                  │
│  Шаг 3: Определение категории                                  │
│  Шаг 4: Генерация вариантов названия                           │
│  Шаг 5: Создание переводов                                     │
│                                                                  │
│  → setAIData({ ..., location: locationData })                  │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     3. ENHANCE VIEW                              │
├─────────────────────────────────────────────────────────────────┤
│  Редактирование: название, описание, цена, категория...         │
│                                                                  │
│  📍 Местоположение товара ← НОВАЯ СЕКЦИЯ                       │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ Выбор типа адреса:                                        │ │
│  │                                                            │ │
│  │ ⦿ Использовать адрес витрины                             │ │
│  │   (useStorefrontLocation = true)                         │ │
│  │                                                            │ │
│  │ ○ Использовать адрес из метаданных фото                  │ │
│  │   (если state.aiData.location?.source === 'exif')        │ │
│  │   Показывает: address, city, region                      │ │
│  │                                                            │ │
│  │ ○ Указать индивидуальный адрес                           │ │
│  │   Открывается LocationPicker                             │ │
│  │   ┌──────────────────────────────────┐                   │ │
│  │   │ 🗺️ Интерактивная карта          │                   │ │
│  │   │ 🔍 Поиск адреса                 │                   │ │
│  │   │ 📍 Клик на карте                │                   │ │
│  │   └──────────────────────────────────┘                   │ │
│  │                                                            │ │
│  │ 🛡️ Настройки приватности:                               │ │
│  │   ⦿ Точный адрес / ○ Улица / ○ Район / ○ Город         │ │
│  │   ☑ Показывать на карте                                  │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  handleSave():                                                  │
│    → setAIData({ location: ... })                              │
│    → setLocationPrivacyLevel(level)                            │
│    → setShowOnMap(show)                                        │
│    → setView('variants')                                       │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     4. VARIANTS VIEW                             │
│  Управление вариантами товара                                   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     5. PUBLISH VIEW                              │
│  Финальная проверка и публикация                                │
│                                                                  │
│  Отправка на backend:                                           │
│  {                                                               │
│    title, description, price, category, ...                     │
│    location: {                                                  │
│      latitude, longitude, address, city, region,               │
│      source: 'exif' | 'storefront' | 'manual'                  │
│    },                                                            │
│    locationPrivacyLevel: 'exact' | 'street' | ...              │
│    showOnMap: true | false                                      │
│  }                                                               │
└─────────────────────────────────────────────────────────────────┘
```

---

## Примеры использования

### Сценарий 1: Пользователь загрузил фото с GPS данными

1. **Upload View:** Загружены 3 фотографии
2. **Process View → Шаг 2 (EXIF):**
   - Status: `processing` → "Обработка..."
   - Извлечены координаты: `45.2671, 19.8335`
   - Геокодирование: "Boulevard of Liberation 12, Novi Sad, Serbia"
   - Status: `completed` → "Novi Sad detected" ✅
3. **Enhance View:**
   - Отображаются 3 варианта адреса
   - Второй вариант (EXIF) содержит: "Boulevard of Liberation 12..."
   - Пользователь выбирает EXIF адрес
   - Устанавливает приватность: "Улица" (скрыть номер дома)
   - Включает "Показывать на карте"
4. **Результат:** Товар публикуется с адресом "Boulevard of Liberation, Novi Sad"

### Сценарий 2: Пользователь загрузил фото без GPS данных

1. **Upload View:** Загружены фотографии со смартфона (GPS отключён)
2. **Process View → Шаг 2 (EXIF):**
   - Status: `processing` → "Обработка..."
   - GPS данные не найдены
   - Status: `completed` → "No GPS data in photos" ℹ️
3. **Enhance View:**
   - Отображаются только 2 варианта:
     - Адрес витрины (выбрано по умолчанию)
     - Индивидуальный адрес
   - Вариант EXIF не отображается (т.к. `location?.source !== 'exif'`)
4. **Результат:** Товар использует адрес витрины

### Сценарий 3: Пользователь хочет указать точное местоположение товара

1. **Enhance View:**
   - Выбирает "Указать индивидуальный адрес"
   - Открывается LocationPicker
   - Вводит в поиск: "Нови Сад, Змай Јовина 24"
   - Кликает на карте для уточнения
   - Координаты определены: `45.2551, 19.8451`
   - Геокодирование: "Zmaj Jovina 24, Novi Sad"
2. **Настройки приватности:**
   - Выбирает "Точный адрес"
   - Включает "Показывать на карте"
3. **Результат:** Товар публикуется с точным адресом "Zmaj Jovina 24, Novi Sad"

---

## Тестирование

### Контрольные точки

#### ProcessView (AI обработка)

- [ ] EXIF шаг отображается между "Анализ изображений" и "Определение категории"
- [ ] Статус меняется: pending → processing → completed
- [ ] При наличии GPS: отображается название города
- [ ] При отсутствии GPS: отображается "No GPS data in photos"
- [ ] При ошибке геокодирования: отображается "No address found"
- [ ] locationData сохраняется в context с source: 'exif'

#### EnhanceView (Выбор адреса)

- [ ] Вариант "Адрес витрины" всегда доступен
- [ ] Вариант "EXIF адрес" отображается только если `location?.source === 'exif'`
- [ ] Вариант "Индивидуальный адрес" всегда доступен
- [ ] LocationPicker открывается при выборе индивидуального адреса
- [ ] Настройки приватности отображаются для EXIF и индивидуального адреса
- [ ] Данные корректно сохраняются при переходе на следующий шаг

#### Context (Состояние)

- [ ] setUseStorefrontLocation корректно обновляет state
- [ ] setLocationPrivacyLevel корректно обновляет state
- [ ] setShowOnMap корректно обновляет state
- [ ] location с source: 'exif' корректно сохраняется
- [ ] location с source: 'manual' корректно сохраняется

---

## Известные ограничения

1. **EXIF GPS данные**
   - Не все камеры и смартфоны записывают GPS в EXIF
   - Пользователь может отключить геотегирование в настройках камеры
   - Некоторые приложения удаляют EXIF при обработке фото

2. **Геокодирование**
   - Nominatim имеет лимит запросов (1 запрос/сек для бесплатного API)
   - Не все координаты имеют точный адрес (например, в поле или море)
   - Качество адресов зависит от полноты данных OpenStreetMap

3. **Приватность**
   - Преобразование адреса в приватный формат выполняется на frontend
   - Backend должен уважать настройки приватности при отображении

---

## Дальнейшие улучшения

### Возможные расширения

1. **Кэширование геокодирования**
   - Сохранять результаты геокодирования в локальное хранилище
   - Избежать повторных запросов для одинаковых координат

2. **Пакетная обработка EXIF**
   - Извлекать GPS из всех фотографий
   - Определять наиболее частое местоположение
   - Предупреждать если фото сделаны в разных местах

3. **Визуализация на карте**
   - Показывать предпросмотр местоположения с текущим уровнем приватности
   - Отображать радиус "размытия" для районов/городов

4. **Автодополнение адреса**
   - Использовать Nominatim для autocomplete
   - Предлагать популярные адреса в городе

5. **Улучшение UX**
   - Добавить tooltips с объяснением каждого варианта
   - Анимации при переключении между вариантами
   - Индикатор точности геокодирования

---

## Заключение

Реализованный функционал обеспечивает:

✅ **Автоматизацию** - определение адреса из EXIF без участия пользователя
✅ **Гибкость** - три варианта выбора адреса на выбор
✅ **Приватность** - контроль над детализацией адреса
✅ **Прозрачность** - визуализация каждого шага процесса
✅ **Локализация** - поддержка трёх языков (ru, en, sr)
✅ **Типобезопасность** - полная типизация на TypeScript

Функционал полностью интегрирован в существующий процесс AI создания продукта витрины и готов к использованию.

---

**Автор:** Claude Code
**Дата последнего обновления:** 2025-10-03
**Статус:** Production Ready
