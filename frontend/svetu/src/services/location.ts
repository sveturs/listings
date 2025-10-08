import { apiClient } from './api-client';

export interface City {
  id: string;
  name: string;
  country: string;
  lat: number;
  lng: number;
  region?: string;
  population?: number;
}

export const LocationService = {
  async searchCities(query: string): Promise<City[]> {
    if (!query.trim()) {
      return [];
    }

    const response = await apiClient.get<City[]>(
      `/api/v1/locations/cities/search?q=${encodeURIComponent(query.trim())}`
    );

    if (response.error) {
      throw new Error(`Failed to search cities: ${response.error}`);
    }

    return response.data || [];
  },

  async getCityByCoordinates(lat: number, lng: number): Promise<City | null> {
    const response = await apiClient.get<City>(
      `/api/v1/locations/cities/reverse?lat=${lat}&lng=${lng}`
    );

    if (response.error) {
      throw new Error(`Failed to get city by coordinates: ${response.error}`);
    }

    return response.data || null;
  },

  async getNearbyListings(
    lat: number,
    lng: number,
    radius: number = 10
  ): Promise<any[]> {
    const response = await apiClient.get<any[]>(
      `/api/v1/c2c/listings/nearby?lat=${lat}&lng=${lng}&radius=${radius}`
    );

    if (response.error) {
      throw new Error(`Failed to get nearby listings: ${response.error}`);
    }

    return response.data || [];
  },
};
