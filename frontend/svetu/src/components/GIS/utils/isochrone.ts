import buffer from '@turf/buffer';
import { point } from '@turf/helpers';
import type { Feature, Polygon } from 'geojson';

/**
 * Генерирует приблизительную изохрону (зону пешей доступности)
 * @param center - Центральная точка [lng, lat]
 * @param walkingTimeMinutes - Время ходьбы в минутах
 * @param walkingSpeedKmh - Скорость ходьбы в км/ч (по умолчанию 5 км/ч)
 * @returns GeoJSON полигон зоны доступности
 */
export function generateIsochrone(
  center: [number, number],
  walkingTimeMinutes: number,
  walkingSpeedKmh: number = 5
): Feature<Polygon> {
  // Рассчитываем расстояние в километрах
  const distanceKm = (walkingTimeMinutes / 60) * walkingSpeedKmh;

  // Создаем точку из центра
  const centerPoint = point(center);

  // Создаем буфер вокруг точки
  // Используем коэффициент 0.7 для более реалистичной формы
  // (учитывает, что в городе путь не прямой)
  const isochrone = buffer(centerPoint, distanceKm * 0.7, {
    units: 'kilometers',
    steps: 64,
  });

  return isochrone as Feature<Polygon>;
}

/**
 * Генерирует стилизованную изохрону с волнообразными краями
 * для более естественного вида
 */
export function generateStylizedIsochrone(
  center: [number, number],
  walkingTimeMinutes: number,
  walkingSpeedKmh: number = 5
): Feature<Polygon> {
  const distanceKm = (walkingTimeMinutes / 60) * walkingSpeedKmh;
  const baseRadius = distanceKm * 0.7; // км

  // Преобразуем в градусы (приблизительно)
  const degreeRadius = baseRadius / 111.32;

  const points: [number, number][] = [];
  const numPoints = 64;

  for (let i = 0; i < numPoints; i++) {
    const angle = (i / numPoints) * 2 * Math.PI;

    // Добавляем небольшие вариации для более естественного вида
    const variation =
      1 + Math.sin(angle * 4) * 0.1 + Math.sin(angle * 7) * 0.05;
    const radius = degreeRadius * variation;

    const lat = center[1] + radius * Math.sin(angle);
    const lng =
      center[0] +
      (radius * Math.cos(angle)) / Math.cos((center[1] * Math.PI) / 180);

    points.push([lng, lat]);
  }

  // Замыкаем полигон
  points.push(points[0]);

  return {
    type: 'Feature',
    geometry: {
      type: 'Polygon',
      coordinates: [points],
    },
    properties: {
      walkingTime: walkingTimeMinutes,
      walkingSpeed: walkingSpeedKmh,
    },
  };
}

/**
 * Преобразует радиус в минуты ходьбы
 */
export function radiusToWalkingTime(
  radiusMeters: number,
  walkingSpeedKmh: number = 5
): number {
  const radiusKm = radiusMeters / 1000;
  const hours = radiusKm / walkingSpeedKmh;
  return Math.round(hours * 60);
}

/**
 * Преобразует минуты ходьбы в радиус
 */
export function walkingTimeToRadius(
  minutes: number,
  walkingSpeedKmh: number = 5
): number {
  const hours = minutes / 60;
  const radiusKm = hours * walkingSpeedKmh;
  return Math.round(radiusKm * 1000);
}
