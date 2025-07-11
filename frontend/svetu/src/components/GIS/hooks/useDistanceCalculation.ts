import { useMemo } from 'react';
import distance from '@turf/distance';
import { point } from '@turf/helpers';
import booleanPointInPolygon from '@turf/boolean-point-in-polygon';
import type { MapMarkerData } from '../types/gis';

interface UseDistanceCalculationProps {
  userLocation: { latitude: number; longitude: number };
  markers: MapMarkerData[];
  searchRadius: number; // в метрах
  radiusCircleGeoJSON?: any;
}

interface MarkerWithDistance extends MapMarkerData {
  distance: number; // в километрах
  isWithinRadius: boolean;
}

export function useDistanceCalculation({
  userLocation,
  markers,
  searchRadius,
  radiusCircleGeoJSON,
}: UseDistanceCalculationProps) {
  const markersWithDistance = useMemo(() => {
    if (!userLocation || !markers) return [];

    const userPoint = point([userLocation.longitude, userLocation.latitude]);

    return markers.map((marker): MarkerWithDistance => {
      const markerPoint = point([marker.longitude, marker.latitude]);

      // Рассчитываем расстояние в километрах
      const distanceInKm = distance(userPoint, markerPoint, {
        units: 'kilometers',
      });

      // Проверяем, находится ли маркер в радиусе поиска
      let isWithinRadius = false;
      if (radiusCircleGeoJSON) {
        isWithinRadius = booleanPointInPolygon(
          markerPoint,
          radiusCircleGeoJSON
        );
      } else {
        // Если нет GeoJSON круга, используем простое сравнение расстояний
        isWithinRadius = distanceInKm <= searchRadius / 1000;
      }

      return {
        ...marker,
        distance: distanceInKm,
        isWithinRadius,
      };
    });
  }, [userLocation, markers, searchRadius, radiusCircleGeoJSON]);

  // Фильтруем маркеры в радиусе и сортируем по расстоянию
  const markersInRadius = useMemo(() => {
    return markersWithDistance
      .filter((marker) => marker.isWithinRadius)
      .sort((a, b) => a.distance - b.distance);
  }, [markersWithDistance]);

  return {
    markersWithDistance,
    markersInRadius,
  };
}
