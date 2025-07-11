import { MapMarkerData } from '../types/gis';

export interface GeoJSONFeature {
  type: 'Feature';
  geometry: {
    type: 'Point';
    coordinates: [number, number]; // [longitude, latitude]
  };
  properties: {
    id: string;
    title: string;
    description?: string;
    type: 'listing' | 'user' | 'poi';
    data?: any;
    // Дополнительные свойства для стилизации
    color?: string;
    size?: number;
    icon?: string;
  };
}

export interface GeoJSONFeatureCollection {
  type: 'FeatureCollection';
  features: GeoJSONFeature[];
}

/**
 * Преобразует маркеры в GeoJSON FeatureCollection для использования с MapBox кластеризацией
 * @param markers - Массив маркеров для преобразования
 * @returns GeoJSON FeatureCollection
 */
export function markersToGeoJSON(
  markers: MapMarkerData[]
): GeoJSONFeatureCollection {
  return {
    type: 'FeatureCollection',
    features: markers.map((marker) => ({
      type: 'Feature',
      geometry: {
        type: 'Point',
        coordinates: [marker.position[0], marker.position[1]],
      },
      properties: {
        id: marker.id,
        title: marker.title,
        description: marker.description,
        type: marker.type,
        data: marker.data,
        // Добавляем дополнительные свойства для стилизации
        color: getMarkerColor(marker.type),
        size: getMarkerSize(marker.type),
        icon: getMarkerIcon(marker.type),
      },
    })),
  };
}

/**
 * Получает цвет маркера по типу
 * @param type - Тип маркера
 * @returns Цвет в формате hex
 */
export function getMarkerColor(type: 'listing' | 'user' | 'poi'): string {
  switch (type) {
    case 'listing':
      return '#3b82f6'; // blue-500
    case 'user':
      return '#f97316'; // orange-500
    case 'poi':
      return '#ef4444'; // red-500
    default:
      return '#6b7280'; // gray-500
  }
}

/**
 * Получает размер маркера по типу
 * @param type - Тип маркера
 * @returns Размер маркера
 */
export function getMarkerSize(type: 'listing' | 'user' | 'poi'): number {
  switch (type) {
    case 'listing':
      return 24;
    case 'user':
      return 20;
    case 'poi':
      return 18;
    default:
      return 16;
  }
}

/**
 * Получает иконку маркера по типу
 * @param type - Тип маркера
 * @returns Эмодзи иконка
 */
export function getMarkerIcon(type: 'listing' | 'user' | 'poi'): string {
  switch (type) {
    case 'listing':
      return '🏠';
    case 'user':
      return '👤';
    case 'poi':
      return '📍';
    default:
      return '📍';
  }
}

/**
 * Преобразует GeoJSON Feature обратно в MapMarkerData
 * @param feature - GeoJSON Feature
 * @returns MapMarkerData
 */
export function geoJsonToMarker(feature: GeoJSONFeature): MapMarkerData {
  return {
    id: feature.properties.id,
    position: [
      feature.geometry.coordinates[0],
      feature.geometry.coordinates[1],
    ],
    title: feature.properties.title,
    description: feature.properties.description,
    type: feature.properties.type,
    data: feature.properties.data,
  };
}

/**
 * Создает пустую GeoJSON FeatureCollection
 * @returns Пустая GeoJSON FeatureCollection
 */
export function createEmptyGeoJSON(): GeoJSONFeatureCollection {
  return {
    type: 'FeatureCollection',
    features: [],
  };
}

/**
 * Объединяет несколько GeoJSON FeatureCollection в одну
 * @param collections - Массив GeoJSON FeatureCollection
 * @returns Объединенная GeoJSON FeatureCollection
 */
export function mergeGeoJSONCollections(
  collections: GeoJSONFeatureCollection[]
): GeoJSONFeatureCollection {
  return {
    type: 'FeatureCollection',
    features: collections.flatMap((collection) => collection.features),
  };
}
