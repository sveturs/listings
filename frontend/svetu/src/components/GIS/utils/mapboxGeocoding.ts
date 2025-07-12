import mbxGeocoding from '@mapbox/mapbox-sdk/services/geocoding';
import type {
  GeocodeRequest,
} from '@mapbox/mapbox-sdk/services/geocoding';

// Определяем тип для MapboxFeature
interface MapboxFeature {
  center: [number, number];
  place_name: string;
  text?: string;
  bbox?: [number, number, number, number];
  place_type?: string[];
  relevance?: number;
  properties?: Record<string, any>;
}

// Определяем интерфейсы для совместимости с существующим API
export interface GeocodingResult {
  display_name: string;
  lat: string;
  lon: string;
  boundingbox?: [string, string, string, string];
  class?: string;
  type?: string;
  importance?: number;
}

export interface NearbySearchResult {
  name: string;
  category: string;
  latitude: number;
  longitude: number;
  distance?: number;
}

class MapboxGeocodingService {
  private client: ReturnType<typeof mbxGeocoding>;

  constructor() {
    const token = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;
    if (!token) {
      throw new Error('NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN is required');
    }
    this.client = mbxGeocoding({
      accessToken: token,
    });
  }

  /**
   * Прямое геокодирование: адрес → координаты
   */
  async geocode(
    query: string,
    options?: {
      limit?: number;
      proximity?: [number, number];
      bbox?: [number, number, number, number];
      country?: string[];
      language?: string[];
      types?: string[];
    }
  ): Promise<GeocodingResult[]> {
    try {
      const request = {
        query,
        limit: options?.limit || 5,
        proximity: options?.proximity,
        bbox: options?.bbox,
        countries: options?.country,
        language: options?.language,
        types: options?.types,
      } as any;

      const response = await this.client.forwardGeocode(request).send();

      return this.convertMapboxFeaturesToResults(response.body.features as any || []);
    } catch (error) {
      console.error('MapBox Geocoding error:', error);
      throw new Error('Geocoding failed');
    }
  }

  /**
   * Обратное геокодирование: координаты → адрес
   */
  async reverseGeocode(
    latitude: number,
    longitude: number,
    options?: {
      types?: string[];
      language?: string[];
    }
  ): Promise<GeocodingResult[]> {
    try {
      const response = await this.client
        .reverseGeocode({
          query: [longitude, latitude],
          types: options?.types,
          language: options?.language,
        } as any)
        .send();

      return this.convertMapboxFeaturesToResults(response.body.features as any || []);
    } catch (error) {
      console.error('MapBox Reverse Geocoding error:', error);
      throw new Error('Reverse geocoding failed');
    }
  }

  /**
   * Поиск мест с категориями (замена searchNearby)
   */
  async searchNearby(
    latitude: number,
    longitude: number,
    query: string = '',
    options?: {
      radius?: number; // в метрах, но MapBox использует proximity
      limit?: number;
      types?: string[];
    }
  ): Promise<NearbySearchResult[]> {
    try {
      const searchQuery = query || 'restaurant shop hotel';

      const response = await this.client
        .forwardGeocode({
          query: searchQuery,
          proximity: [longitude, latitude],
          limit: options?.limit || 10,
          types: options?.types || ['poi'],
        } as any)
        .send();

      return this.convertToNearbyResults(response.body.features as any || [], [
        longitude,
        latitude,
      ]);
    } catch (error) {
      console.error('MapBox Nearby Search error:', error);
      throw new Error('Nearby search failed');
    }
  }

  /**
   * Автодополнение для городов
   */
  async getCitySuggestions(
    query: string,
    options?: {
      country?: string[];
      language?: string[];
      limit?: number;
    }
  ): Promise<GeocodingResult[]> {
    try {
      const response = await this.client
        .forwardGeocode({
          query,
          types: ['place', 'locality', 'administrative'],
          countries: options?.country,
          language: options?.language,
          limit: options?.limit || 5,
        } as any)
        .send();

      return this.convertMapboxFeaturesToResults(response.body.features as any || []);
    } catch (error) {
      console.error('MapBox City Suggestions error:', error);
      throw new Error('City suggestions failed');
    }
  }

  /**
   * Конвертация MapBox features в стандартный формат
   */
  private convertMapboxFeaturesToResults(
    features: MapboxFeature[]
  ): GeocodingResult[] {
    return features.map((feature) => {
      const [lon, lat] = feature.center;

      return {
        display_name: feature.place_name,
        lat: lat.toString(),
        lon: lon.toString(),
        boundingbox: feature.bbox
          ? [
              feature.bbox[1].toString(), // south
              feature.bbox[3].toString(), // north
              feature.bbox[0].toString(), // west
              feature.bbox[2].toString(), // east
            ]
          : undefined,
        class: feature.place_type?.[0] || 'unknown',
        type:
          feature.properties?.category || feature.place_type?.[0] || 'unknown',
        importance: feature.relevance,
      };
    });
  }

  /**
   * Конвертация в формат NearbySearchResult
   */
  private convertToNearbyResults(
    features: MapboxFeature[],
    center: [number, number]
  ): NearbySearchResult[] {
    return features.map((feature) => {
      const [lon, lat] = feature.center;

      // Вычисляем расстояние от центра
      const distance = this.calculateDistance(
        center[1],
        center[0], // lat, lon центра
        lat,
        lon // lat, lon точки
      );

      return {
        name: feature.text || feature.place_name,
        category:
          feature.properties?.category || feature.place_type?.[0] || 'poi',
        latitude: lat,
        longitude: lon,
        distance: Math.round(distance),
      };
    });
  }

  /**
   * Вычисление расстояния между двумя точками (в метрах)
   */
  private calculateDistance(
    lat1: number,
    lon1: number,
    lat2: number,
    lon2: number
  ): number {
    const R = 6371e3; // Радиус Земли в метрах
    const φ1 = (lat1 * Math.PI) / 180;
    const φ2 = (lat2 * Math.PI) / 180;
    const Δφ = ((lat2 - lat1) * Math.PI) / 180;
    const Δλ = ((lon2 - lon1) * Math.PI) / 180;

    const a =
      Math.sin(Δφ / 2) * Math.sin(Δφ / 2) +
      Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

    return R * c;
  }
}

// Создаем синглтон экземпляр
export const mapboxGeocoding = new MapboxGeocodingService();

// Экспортируем также класс для создания дополнительных экземпляров
export { MapboxGeocodingService };
