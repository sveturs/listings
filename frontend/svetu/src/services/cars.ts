import configManager from '@/config';
import type { CarMake, CarModel, CarGeneration } from '@/types/cars';

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export class CarsService {
  private static get baseUrl() {
    return `${configManager.getApiUrl()}/api/v1/cars`;
  }

  /**
   * Получить все марки автомобилей
   */
  static async getMakes(): Promise<ApiResponse<CarMake[]>> {
    try {
      const response = await fetch(`${this.baseUrl}/makes`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data.data || data,
      };
    } catch (error) {
      console.error('Error fetching car makes:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  /**
   * Получить модели для конкретной марки
   */
  static async getModelsByMake(
    makeSlug: string
  ): Promise<ApiResponse<CarModel[]>> {
    try {
      const response = await fetch(`${this.baseUrl}/makes/${makeSlug}/models`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data.data || data,
      };
    } catch (error) {
      console.error('Error fetching car models:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  /**
   * Получить поколения для конкретной модели
   */
  static async getGenerationsByModel(
    modelId: number
  ): Promise<ApiResponse<CarGeneration[]>> {
    try {
      const response = await fetch(
        `${this.baseUrl}/models/${modelId}/generations`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data.data || data,
      };
    } catch (error) {
      console.error('Error fetching car generations:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  /**
   * Поиск марок по названию
   */
  static async searchMakes(query: string): Promise<ApiResponse<CarMake[]>> {
    try {
      const response = await fetch(
        `${this.baseUrl}/makes/search?q=${encodeURIComponent(query)}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data.data || data,
      };
    } catch (error) {
      console.error('Error searching car makes:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }
}
