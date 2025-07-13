# MapboxClusterLayer

Компонент для нативной кластеризации маркеров на карте с использованием возможностей Mapbox GL JS.

## Особенности

- **Нативная кластеризация**: Использует встроенные возможности Mapbox GL для оптимальной производительности
- **Автоматическое группирование**: Маркеры автоматически группируются в кластеры на основе настроек
- **Интерактивность**: Поддержка кликов по кластерам и индивидуальным маркерам
- **Настраиваемые стили**: Гибкая настройка внешнего вида кластеров и маркеров
- **Отображение цен**: Опциональное отображение цен на маркерах объявлений

## Использование

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
  />
</Map>;
```

## Props

| Prop               | Type              | Default | Description                                        |
| ------------------ | ----------------- | ------- | -------------------------------------------------- |
| `markers`          | `MapMarkerData[]` | -       | Массив маркеров для кластеризации                  |
| `clusterRadius`    | `number`          | `50`    | Радиус кластеризации в пикселях                    |
| `clusterMaxZoom`   | `number`          | `14`    | Максимальный уровень зума для кластеризации        |
| `clusterMinPoints` | `number`          | `2`     | Минимальное количество точек для создания кластера |
| `onClusterClick`   | `function`        | -       | Обработчик клика по кластеру                       |
| `onMarkerClick`    | `function`        | -       | Обработчик клика по маркеру                        |
| `showPrices`       | `boolean`         | `false` | Показывать цены на маркерах объявлений             |
| `clusterStyles`    | `object`          | -       | Настройки стилей кластеров                         |

## Настройка стилей кластеров

```tsx
const clusterStyles = {
  small: {
    color: '#3b82f6',
    size: 40,
    textColor: '#ffffff',
  },
  medium: {
    color: '#059669',
    size: 55,
    textColor: '#ffffff',
  },
  large: {
    color: '#dc2626',
    size: 70,
    textColor: '#ffffff',
  },
};

<MapboxClusterLayer markers={markers} clusterStyles={clusterStyles} />;
```

## Обработчики событий

### onClusterClick

```tsx
const handleClusterClick = (
  clusterId: number,
  coordinates: [number, number]
) => {
  // Автоматический зум к кластеру
  map.flyTo({
    center: coordinates,
    zoom: map.getZoom() + 2,
  });
};
```

### onMarkerClick

```tsx
const handleMarkerClick = (marker: MapMarkerData) => {
  // Показать подробности маркера
  setSelectedMarker(marker);
  // Открыть popup
  setPopup({
    id: marker.id,
    position: marker.position,
    title: marker.title,
    content: <MarkerDetails marker={marker} />,
  });
};
```

## Типы маркеров

Компонент поддерживает три типа маркеров:

- **`listing`**: Объявления (синий цвет, размер 24px)
- **`user`**: Пользователи (оранжевый цвет, размер 20px)
- **`poi`**: Точки интереса (красный цвет, размер 18px)

## Производительность

- **Нативная кластеризация**: Использует встроенные возможности Mapbox GL
- **Автоматическое управление**: Кластеры создаются и обновляются автоматически
- **Оптимизация рендеринга**: Только видимые маркеры отрисовываются на карте

## Совместимость

- Требует Mapbox GL JS 2.0+
- Совместим с React 18+
- Поддерживает TypeScript

## Примеры

### Базовое использование

```tsx
import { MapboxClusterLayer } from '@/components/GIS';

const markers = [
  {
    id: '1',
    position: [20.4649, 44.8176],
    title: 'Объявление 1',
    type: 'listing',
    data: { price: 150 },
  },
  // ... другие маркеры
];

<MapboxClusterLayer markers={markers} onMarkerClick={handleMarkerClick} />;
```

### С настраиваемыми стилями

```tsx
const customStyles = {
  small: { color: '#10b981', size: 35 },
  medium: { color: '#f59e0b', size: 50 },
  large: { color: '#ef4444', size: 65 },
};

<MapboxClusterLayer
  markers={markers}
  clusterStyles={customStyles}
  showPrices={true}
/>;
```

### С обработчиками событий

```tsx
const handleClusterClick = (clusterId, coordinates) => {
  console.log('Cluster clicked:', clusterId);
  // Логика обработки клика по кластеру
};

const handleMarkerClick = (marker) => {
  console.log('Marker clicked:', marker);
  // Логика обработки клика по маркеру
};

<MapboxClusterLayer
  markers={markers}
  onClusterClick={handleClusterClick}
  onMarkerClick={handleMarkerClick}
/>;
```

## Связанные компоненты

- [`InteractiveMap`](./InteractiveMap.md) - Основной компонент карты
- [`MapMarker`](./MapMarker.md) - Индивидуальный маркер
- [`MapCluster`](./MapCluster.md) - Компонент кластера (для серверной кластеризации)

## Утилиты

- [`markersToGeoJSON`](../utils/geoJsonHelpers.md) - Преобразование маркеров в GeoJSON
- [`geoJsonToMarker`](../utils/geoJsonHelpers.md) - Преобразование GeoJSON в маркеры
