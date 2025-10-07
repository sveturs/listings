import { apiClient } from '@/services/api-client';
import type { components } from '@/types/generated/api';

type CarMake = components['schemas']['models.CarMake'];
type CarModel =
  components['schemas']['models.CarModel'];
type CarGeneration =
  components['schemas']['models.CarGeneration'];
type VINDecodeResult =
  components['schemas']['models.VINDecodeResult'];

export interface CarMakesResponse {
  success: boolean;
  data?: CarMake[];
  error?: string;
}

export interface CarModelsResponse {
  success: boolean;
  data?: CarModel[];
  error?: string;
}

export interface CarGenerationsResponse {
  success: boolean;
  data?: CarGeneration[];
  error?: string;
}

export interface VINDecodeResponse {
  success: boolean;
  data?: VINDecodeResult;
  error?: string;
}

export class CarsService {
  static async getMakes(params?: {
    country?: string;
    is_domestic?: boolean;
    active_only?: boolean;
  }): Promise<CarMakesResponse> {
    try {
      const queryParams = new URLSearchParams();
      if (params?.country) queryParams.append('country', params.country);
      if (params?.is_domestic !== undefined) {
        queryParams.append('is_domestic', params.is_domestic.toString());
      }
      if (params?.active_only !== undefined) {
        queryParams.append('active_only', params.active_only.toString());
      }

      const response = await apiClient.get<{
        success: boolean;
        data?: CarMake[];
      }>(
        `/api/v1/cars/makes${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      );

      if (!response.data) {
        throw new Error('No data received');
      }

      return {
        success: response.data.success,
        data: response.data.data,
      };
    } catch (error) {
      console.error('Failed to fetch car makes:', error);
      return {
        success: false,
        error: 'Failed to load car makes',
      };
    }
  }

  static async getModels(
    makeSlug: string,
    activeOnly = true
  ): Promise<CarModelsResponse> {
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('active_only', activeOnly.toString());

      const response = await apiClient.get<{
        success: boolean;
        data?: CarModel[];
      }>(`/api/v1/cars/makes/${makeSlug}/models?${queryParams.toString()}`);

      if (!response.data) {
        throw new Error('No data received');
      }

      return {
        success: response.data.success,
        data: response.data.data,
      };
    } catch (error) {
      console.error('Failed to fetch car models:', error);
      return {
        success: false,
        error: 'Failed to load car models',
      };
    }
  }

  static async getGenerations(
    modelId: number,
    activeOnly = true
  ): Promise<CarGenerationsResponse> {
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('active_only', activeOnly.toString());

      const response = await apiClient.get<{
        success: boolean;
        data?: CarGeneration[];
      }>(
        `/api/v1/cars/models/${modelId}/generations?${queryParams.toString()}`
      );

      if (!response.data) {
        throw new Error('No data received');
      }

      return {
        success: response.data.success,
        data: response.data.data,
      };
    } catch (error) {
      console.error('Failed to fetch car generations:', error);
      return {
        success: false,
        error: 'Failed to load car generations',
      };
    }
  }

  /**
   * Decode a Vehicle Identification Number (VIN)
   * @param vin The VIN to decode (must be 17 characters)
   * @returns VIN decode result with vehicle information
   */
  static async decodeVIN(vin: string): Promise<VINDecodeResponse> {
    try {
      // Validate VIN format (basic check)
      if (!vin || vin.length !== 17) {
        return {
          success: false,
          error: 'VIN must be exactly 17 characters',
        };
      }

      // Clean VIN (remove spaces and convert to uppercase)
      const cleanVIN = vin.trim().toUpperCase();

      const response = await apiClient.get<{
        success: boolean;
        data?: VINDecodeResult;
      }>(`/api/v1/cars/vin/${cleanVIN}/decode`);

      if (!response.data) {
        throw new Error('No data received');
      }

      return {
        success: response.data.success,
        data: response.data.data,
      };
    } catch (error: any) {
      console.error('Failed to decode VIN:', error);

      // Handle specific error cases
      if (error.response?.status === 404) {
        return {
          success: false,
          error: 'VIN not found or invalid',
        };
      }

      if (error.response?.status === 429) {
        return {
          success: false,
          error: 'Too many requests. Please try again later',
        };
      }

      return {
        success: false,
        error: error.message || 'Failed to decode VIN',
      };
    }
  }
}

export default CarsService;
