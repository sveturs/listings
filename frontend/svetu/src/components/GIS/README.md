# GIS Map Components

Набор компонентов для работы с интерактивными картами на базе Mapbox GL и react-map-gl.

## Установка зависимостей

```bash
yarn add mapbox-gl react-map-gl @turf/turf @types/mapbox-gl
```

## Настройка

1. Добавьте токен Mapbox в `.env.local`:

```env
NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=pk.eyJ1IjoieW91cnVzZXJuYW1lIiwiYSI6ImNsZjFjZmZiNDAzZGczZHFxcXFxcXFxcSJ9.example
```

2. Получите токен на [https://account.mapbox.com/access-tokens/](https://account.mapbox.com/access-tokens/)

## Основные компоненты

### InteractiveMap

Основной компонент интерактивной карты с поддержкой маркеров, всплывающих окон и элементов управления.

```tsx
import { InteractiveMap } from '@/components/GIS';

<InteractiveMap
  initialViewState={{
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 12,
  }}
  markers={markers}
  popup={popup}
  onMarkerClick={handleMarkerClick}
  onMapClick={handleMapClick}
  controlsConfig={{
    showZoom: true,
    showGeolocate: true,
    position: 'top-right',
  }}
/>;
```

### MapMarker

Компонент для отображения маркеров на карте с поддержкой различных типов и стилей.

```tsx
import { MapMarker } from '@/components/GIS';

<MapMarker
  marker={{
    id: '1',
    position: [20.4649, 44.8176],
    title: 'Хостел "Центр"',
    type: 'listing',
    data: { price: 25, currency: '€' },
  }}
  onClick={handleMarkerClick}
  selected={selectedMarkerId === '1'}
/>;
```

### MapPopup

Компонент всплывающих окон с настраиваемым содержимым.

```tsx
import { MapPopup, ListingPopup } from '@/components/GIS';

// Обычный popup
<MapPopup
  popup={{
    id: '1',
    position: [20.4649, 44.8176],
    title: 'Название',
    description: 'Описание',
    content: <CustomContent />
  }}
  onClose={handleClose}
/>

// Специализированный popup для объявлений
<ListingPopup
  listing={{
    id: '1',
    title: 'Хостел',
    price: 25,
    currency: '€',
    rating: 4.5,
    reviews: 123
  }}
  position={[20.4649, 44.8176]}
  onClose={handleClose}
  onViewDetails={handleViewDetails}
/>
```

### MapControls

Элементы управления картой включая поиск, переключение стилей и геолокацию.

```tsx
import { MapControls } from '@/components/GIS';

<MapControls
  config={{
    showZoom: true,
    showCompass: true,
    showFullscreen: true,
    showGeolocate: true,
    position: 'top-right',
  }}
  onStyleChange={handleStyleChange}
  onSearch={handleSearch}
/>;
```

### MapboxClusterLayer

Компонент для нативной кластеризации маркеров с использованием встроенных возможностей Mapbox GL.

```tsx
import { MapboxClusterLayer } from '@/components/GIS';

<Map {...viewState} mapboxAccessToken={accessToken}>
  <MapboxClusterLayer
    markers={markers}
    clusterRadius={50}
    clusterMaxZoom={14}
    onClusterClick={handleClusterClick}
    onMarkerClick={handleMarkerClick}
    showPrices={true}
    clusterStyles={{
      small: { color: '#3b82f6', size: 40 },
      medium: { color: '#059669', size: 55 },
      large: { color: '#dc2626', size: 70 },
    }}
  />
</Map>;
```

## Хуки

### useGeolocation

Хук для работы с геолокацией пользователя.

```tsx
import { useGeolocation } from '@/components/GIS';

const {
  location,
  error,
  loading,
  supported,
  getCurrentPosition,
  watchPosition,
  clearWatch,
} = useGeolocation({
  enableHighAccuracy: true,
  timeout: 15000,
  watch: false,
});

// Получить текущую позицию
const handleGetLocation = async () => {
  try {
    const position = await getCurrentPosition();
    console.log('Текущая позиция:', position);
  } catch (error) {
    console.error('Ошибка геолокации:', error);
  }
};
```

### useGeoSearch

Хук для поиска по адресам и координатам.

```tsx
import { useGeoSearch } from '@/components/GIS';

const { results, loading, error, search, searchNearby, clearResults } =
  useGeoSearch();

// Поиск по адресу
const handleSearch = async (query: string) => {
  const results = await search({
    query,
    limit: 10,
    language: 'ru',
    countrycodes: 'rs',
  });
  console.log('Результаты поиска:', results);
};

// Поиск поблизости
const handleNearbySearch = async () => {
  const nearby = await searchNearby({
    latitude: 44.8176,
    longitude: 20.4649,
    radius: 1000,
    category: 'accommodation',
  });
  console.log('Объекты поблизости:', nearby);
};
```

### useRouteCalculation

Хук для расчета маршрутов между точками.

```tsx
import { useRouteCalculation } from '@/components/GIS';

const { loading, error, calculateRoute } = useRouteCalculation();

const handleGetRoute = async () => {
  const route = await calculateRoute(
    [20.4649, 44.8176], // откуда
    [20.47, 44.82], // куда
    'driving' // тип маршрута
  );
  console.log('Маршрут:', route);
};
```

### useGeocoding

Хук для прямого и обратного геокодирования.

```tsx
import { useGeocoding } from '@/components/GIS';

const { loading, error, geocode, reverseGeocode } = useGeocoding();

// Прямое геокодирование (адрес → координаты)
const handleGeocode = async () => {
  const result = await geocode('Knez Mihailova 10, Belgrade');
  console.log('Координаты:', result);
};

// Обратное геокодирование (координаты → адрес)
const handleReverseGeocode = async () => {
  const result = await reverseGeocode(44.8176, 20.4649);
  console.log('Адрес:', result);
};
```

## Типы данных

### MapMarkerData

```tsx
interface MapMarkerData {
  id: string;
  position: [number, number]; // [longitude, latitude]
  title: string;
  description?: string;
  type: 'listing' | 'user' | 'poi';
  data?: any;
}
```

### MapViewState

```tsx
interface MapViewState {
  longitude: number;
  latitude: number;
  zoom: number;
  pitch?: number;
  bearing?: number;
}
```

### GeoLocation

```tsx
interface GeoLocation {
  latitude: number;
  longitude: number;
  accuracy?: number;
  altitude?: number;
  heading?: number;
  speed?: number;
}
```

## Стили карты

Доступные стили Mapbox:

- `streets` - Улицы (по умолчанию)
- `satellite` - Спутниковые снимки
- `outdoors` - Природа и туризм
- `light` - Светлая тема
- `dark` - Темная тема

## Утилитные функции

### Расчет расстояния

```tsx
import { calculateDistance, formatDistance } from '@/components/GIS';

const distance = calculateDistance(
  44.8176,
  20.4649, // точка A
  44.82,
  20.47 // точка B
);

const formatted = formatDistance(distance * 1000, 'm'); // "1.2 км"
```

### Расчет направления

```tsx
import { calculateBearing } from '@/components/GIS';

const bearing = calculateBearing(
  44.8176,
  20.4649, // откуда
  44.82,
  20.47 // куда
); // возвращает угол в градусах
```

## Примеры использования

### Простая карта с маркерами

```tsx
import React, { useState } from 'react';
import { InteractiveMap, MapMarkerData } from '@/components/GIS';

const SimpleMap: React.FC = () => {
  const [markers] = useState<MapMarkerData[]>([
    {
      id: '1',
      position: [20.4649, 44.8176],
      title: 'Белград',
      type: 'poi',
    },
  ]);

  return (
    <div className="w-full h-96">
      <InteractiveMap
        initialViewState={{
          longitude: 20.4649,
          latitude: 44.8176,
          zoom: 12,
        }}
        markers={markers}
      />
    </div>
  );
};
```

### Карта с поиском

```tsx
import React, { useState } from 'react';
import { InteractiveMap, useGeoSearch } from '@/components/GIS';

const SearchableMap: React.FC = () => {
  const { search } = useGeoSearch();
  const [viewState, setViewState] = useState({
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 12,
  });

  const handleSearch = async (query: string) => {
    const results = await search({ query, limit: 1 });
    if (results.length > 0) {
      const result = results[0];
      setViewState({
        longitude: parseFloat(result.lon),
        latitude: parseFloat(result.lat),
        zoom: 14,
      });
    }
  };

  return (
    <InteractiveMap
      initialViewState={viewState}
      controlsConfig={{
        showZoom: true,
        showGeolocate: true,
      }}
      onSearch={handleSearch}
    />
  );
};
```

## Интернационализация

Компоненты поддерживают переводы через next-intl. Ключи переводов находятся в разделе `gis` файлов `messages/ru.json` и `messages/en.json`.

```tsx
import { useTranslations } from 'next-intl';

const t = useTranslations('gis.map');
const searchPlaceholder = t('search_placeholder'); // "Поиск на карте..."
```

## Производительность

### Оптимизация маркеров

Для большого количества маркеров рекомендуется использовать кластеризацию:

```tsx
import { MapboxClusterLayer } from '@/components/GIS';

// Нативная кластеризация Mapbox (рекомендуется)
<MapboxClusterLayer
  markers={markers}
  clusterRadius={50}
  clusterMaxZoom={14}
  onClusterClick={handleClusterClick}
  onMarkerClick={handleMarkerClick}
/>

// Или серверная кластеризация для очень больших датасетов
<InteractiveMap
  loadClusters={loadClustersFromServer}
  // другие пропсы
/>
```

### Ленивая загрузка

Загружайте карту только когда она действительно нужна:

```tsx
import dynamic from 'next/dynamic';

const InteractiveMap = dynamic(
  () => import('@/components/GIS').then((mod) => mod.InteractiveMap),
  { ssr: false }
);
```

## Безопасность

- Токен Mapbox должен быть ограничен по доменам в настройках аккаунта
- Не включайте секретные токены в клиентский код
- Используйте переменные окружения для конфигурации

## Устранение неполадок

### Карта не загружается

1. Проверьте токен Mapbox в переменных окружения
2. Убедитесь что домен добавлен в настройки токена
3. Проверьте консоль браузера на ошибки

### Медленная работа

1. Уменьшите количество маркеров на карте
2. Используйте оптимизированные изображения
3. Настройте кэширование статических ресурсов

### Проблемы с геолокацией

1. Проверьте что сайт использует HTTPS
2. Убедитесь что пользователь дал разрешение
3. Проверьте поддержку геолокации в браузере
