declare module '@mapbox/mapbox-sdk/services/geocoding' {
  export interface MapboxResponse {
    body: {
      features: Array<{
        id: string;
        type: string;
        place_type: string[];
        relevance: number;
        properties: Record<string, any>;
        text: string;
        place_name: string;
        bbox?: number[];
        center: [number, number];
        geometry: {
          type: string;
          coordinates: [number, number];
        };
        context?: Array<{
          id: string;
          short_code?: string;
          wikidata?: string;
          text: string;
        }>;
      }>;
    };
  }

  export interface MapboxGeocodingConfig {
    accessToken: string;
  }

  export interface GeocodeRequest {
    query: string;
    limit?: number;
    language?: string | string[];
    country?: string | string[];
    proximity?: [number, number];
    bbox?: [number, number, number, number];
    types?: string[];
  }

  export interface MapboxGeocodingClient {
    forwardGeocode(options: GeocodeRequest): {
      send(): Promise<MapboxResponse>;
    };
    reverseGeocode(options: {
      query: [number, number];
      limit?: number;
      language?: string | string[];
    }): {
      send(): Promise<MapboxResponse>;
    };
  }

  export default function mbxGeocoding(
    config: MapboxGeocodingConfig
  ): MapboxGeocodingClient;
}
