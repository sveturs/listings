import type { Feature, Polygon } from 'geojson';
import booleanPointInPolygon from '@turf/boolean-point-in-polygon';
import { point } from '@turf/helpers';
import { generateLocalIsochrone } from './isochrone';

const MAPBOX_ACCESS_TOKEN = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;
const MAPBOX_ISOCHRONE_API = 'https://api.mapbox.com/isochrone/v1/mapbox';

interface MapboxIsochroneOptions {
  coordinates: [number, number];
  minutes: number;
  profile?: 'walking' | 'cycling' | 'driving';
}

interface MapboxIsochroneResponse {
  features: Feature<Polygon>[];
  type: 'FeatureCollection';
}

/**
 * Получает изохрону из Mapbox Isochrone API
 * @param options - параметры запроса
 * @returns GeoJSON Feature с полигоном изохроны
 */
export async function getMapboxIsochrone(
  options: MapboxIsochroneOptions
): Promise<Feature<Polygon>> {
  const { coordinates, minutes, profile = 'walking' } = options;

  try {
    // Формируем URL для запроса
    const [lng, lat] = coordinates;
    const url = new URL(`${MAPBOX_ISOCHRONE_API}/${profile}/${lng},${lat}`);

    // Проверяем наличие токена
    if (!MAPBOX_ACCESS_TOKEN) {
      throw new Error('MAPBOX_ACCESS_TOKEN is not configured');
    }

    // Добавляем параметры
    url.searchParams.append('contours_minutes', minutes.toString());
    url.searchParams.append('polygons', 'true');
    url.searchParams.append('access_token', MAPBOX_ACCESS_TOKEN);

    // Дополнительные параметры для улучшения качества
    url.searchParams.append('denoise', '1'); // Сглаживание границ
    url.searchParams.append('generalize', '50'); // Упрощение геометрии (в метрах)

    // Выполняем запрос
    const response = await fetch(url.toString());

    if (!response.ok) {
      throw new Error(
        `Mapbox API error: ${response.status} ${response.statusText}`
      );
    }

    const data: MapboxIsochroneResponse = await response.json();

    // Проверяем наличие результатов
    if (!data.features || data.features.length === 0) {
      throw new Error('No isochrone data received from Mapbox');
    }

    // Возвращаем первый полигон (самый внешний контур)
    return data.features[0];
  } catch (error) {
    console.error('Mapbox Isochrone API error:', error);

    // Fallback на локальную генерацию
    console.log('Falling back to local isochrone generation');
    return generateLocalIsochrone(coordinates, minutes);
  }
}

/**
 * Получает несколько изохрон для разных временных интервалов
 * @param coordinates - координаты центральной точки
 * @param minutesArray - массив временных интервалов
 * @param profile - профиль передвижения
 * @returns Массив GeoJSON Features с полигонами изохрон
 */
export async function getMultipleMapboxIsochrones(
  coordinates: [number, number],
  minutesArray: number[],
  profile: 'walking' | 'cycling' | 'driving' = 'walking'
): Promise<Feature<Polygon>[]> {
  try {
    // Проверяем наличие токена
    if (!MAPBOX_ACCESS_TOKEN) {
      throw new Error('MAPBOX_ACCESS_TOKEN is not configured');
    }

    // Mapbox API поддерживает запрос нескольких контуров за раз
    const url = new URL(
      `${MAPBOX_ISOCHRONE_API}/${profile}/${coordinates[0]},${coordinates[1]}`
    );

    // Добавляем все временные интервалы
    url.searchParams.append('contours_minutes', minutesArray.join(','));
    url.searchParams.append('polygons', 'true');
    url.searchParams.append('access_token', MAPBOX_ACCESS_TOKEN);
    url.searchParams.append('denoise', '1');
    url.searchParams.append('generalize', '50');

    const response = await fetch(url.toString());

    if (!response.ok) {
      throw new Error(
        `Mapbox API error: ${response.status} ${response.statusText}`
      );
    }

    const data: MapboxIsochroneResponse = await response.json();

    if (!data.features || data.features.length === 0) {
      throw new Error('No isochrone data received from Mapbox');
    }

    // Сортируем по времени (от меньшего к большему)
    return data.features.sort((a, b) => {
      const timeA = a.properties?.contour || 0;
      const timeB = b.properties?.contour || 0;
      return timeA - timeB;
    });
  } catch (error) {
    console.error('Mapbox Isochrone API error:', error);

    // Fallback на локальную генерацию для каждого интервала
    console.log('Falling back to local isochrone generation');
    return Promise.all(
      minutesArray.map((minutes) =>
        generateLocalIsochrone(coordinates, minutes)
      )
    );
  }
}

/**
 * Проверяет доступность Mapbox Isochrone API
 * @returns true если API доступен
 */
export async function checkMapboxIsochroneAvailability(): Promise<boolean> {
  try {
    // Проверяем наличие токена
    if (!MAPBOX_ACCESS_TOKEN) {
      return false;
    }

    // Тестовый запрос с минимальными параметрами
    const testCoordinates: [number, number] = [0, 0];
    const url = new URL(
      `${MAPBOX_ISOCHRONE_API}/walking/${testCoordinates[0]},${testCoordinates[1]}`
    );
    url.searchParams.append('contours_minutes', '1');
    url.searchParams.append('access_token', MAPBOX_ACCESS_TOKEN);

    const response = await fetch(url.toString(), { method: 'HEAD' });
    return response.ok;
  } catch {
    return false;
  }
}

/**
 * Проверяет, находится ли точка внутри изохрона
 * @param pointCoordinates - координаты точки [lng, lat]
 * @param isochrone - полигон изохрона (GeoJSON Feature)
 * @returns true если точка находится внутри изохрона
 */
export function isPointInIsochrone(
  pointCoordinates: [number, number],
  isochrone: Feature<Polygon>
): boolean {
  try {
    // Создаем точку из координат
    const testPoint = point(pointCoordinates);

    // Проверяем, находится ли точка внутри полигона
    return booleanPointInPolygon(testPoint, isochrone);
  } catch (error) {
    console.error('Error checking point in isochrone:', error);
    return false;
  }
}
