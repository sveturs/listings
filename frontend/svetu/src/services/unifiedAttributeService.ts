import { apiClient } from '@/services/api-client';
import type { components } from '@/types/generated/api';

type UnifiedAttribute = components['schemas']['models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['models.UnifiedAttributeValue'];

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

interface GetAttributesParams {
  category_id?: number;
  group?: string;
  required?: boolean;
  searchable?: boolean;
  filterable?: boolean;
  page?: number;
  limit?: number;
}

interface CreateAttributeValueParams {
  listing_id?: number;
  product_id?: number;
  attribute_id: number;
  text_value?: string;
  numeric_value?: number;
  boolean_value?: boolean;
  date_value?: string;
  json_value?: any;
}

interface UpdateAttributeValueParams extends CreateAttributeValueParams {
  id: number;
}

class UnifiedAttributeService {
  /**
   * Получить список атрибутов для категории
   */
  async getCategoryAttributes(
    categoryId: number
  ): Promise<ApiResponse<UnifiedAttribute[]>> {
    try {
      const response = await apiClient.get<{ data: UnifiedAttribute[] }>(
        `/v2/attributes/category/${categoryId}`
      );

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error fetching category attributes:', error);
      return {
        success: false,
        error: 'Failed to fetch category attributes',
      };
    }
  }

  /**
   * Получить список всех атрибутов с фильтрами
   */
  async getAttributes(
    params?: GetAttributesParams
  ): Promise<ApiResponse<UnifiedAttribute[]>> {
    try {
      const queryParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null) {
            queryParams.append(key, value.toString());
          }
        });
      }

      const endpoint = `/v2/attributes${
        queryParams.toString() ? `?${queryParams.toString()}` : ''
      }`;

      const response = await apiClient.get<{ data: UnifiedAttribute[] }>(
        endpoint
      );

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error fetching attributes:', error);
      return {
        success: false,
        error: 'Failed to fetch attributes',
      };
    }
  }

  /**
   * Получить один атрибут по ID
   */
  async getAttribute(
    attributeId: number
  ): Promise<ApiResponse<UnifiedAttribute>> {
    try {
      const response = await apiClient.get<{ data: UnifiedAttribute }>(
        `/v2/attributes/${attributeId}`
      );

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error fetching attribute:', error);
      return {
        success: false,
        error: 'Failed to fetch attribute',
      };
    }
  }

  /**
   * Получить значения атрибутов для объявления
   */
  async getListingAttributeValues(
    listingId: number
  ): Promise<ApiResponse<UnifiedAttributeValue[]>> {
    try {
      const response = await apiClient.get<{
        data: UnifiedAttributeValue[];
      }>(`/v2/listings/${listingId}/attributes`);

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error fetching listing attribute values:', error);
      return {
        success: false,
        error: 'Failed to fetch listing attribute values',
      };
    }
  }

  /**
   * Получить значения атрибутов для продукта
   */
  async getProductAttributeValues(
    productId: number
  ): Promise<ApiResponse<UnifiedAttributeValue[]>> {
    try {
      const response = await apiClient.get<{
        data: UnifiedAttributeValue[];
      }>(`/v2/products/${productId}/attributes`);

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error fetching product attribute values:', error);
      return {
        success: false,
        error: 'Failed to fetch product attribute values',
      };
    }
  }

  /**
   * Создать значение атрибута
   */
  async createAttributeValue(
    params: CreateAttributeValueParams
  ): Promise<ApiResponse<UnifiedAttributeValue>> {
    try {
      const endpoint = params.listing_id
        ? `/v2/listings/${params.listing_id}/attributes`
        : `/v2/products/${params.product_id}/attributes`;

      const payload = {
        attribute_id: params.attribute_id,
        text_value: params.text_value,
        numeric_value: params.numeric_value,
        boolean_value: params.boolean_value,
        date_value: params.date_value,
        json_value: params.json_value,
      };

      const response = await apiClient.post<{ data: UnifiedAttributeValue }>(
        endpoint,
        payload
      );

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error creating attribute value:', error);
      return {
        success: false,
        error: 'Failed to create attribute value',
      };
    }
  }

  /**
   * Обновить значение атрибута
   */
  async updateAttributeValue(
    params: UpdateAttributeValueParams
  ): Promise<ApiResponse<UnifiedAttributeValue>> {
    try {
      const endpoint = params.listing_id
        ? `/v2/listings/${params.listing_id}/attributes/${params.id}`
        : `/v2/products/${params.product_id}/attributes/${params.id}`;

      const payload = {
        text_value: params.text_value,
        numeric_value: params.numeric_value,
        boolean_value: params.boolean_value,
        date_value: params.date_value,
        json_value: params.json_value,
      };

      const response = await apiClient.put<{ data: UnifiedAttributeValue }>(
        endpoint,
        payload
      );

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error updating attribute value:', error);
      return {
        success: false,
        error: 'Failed to update attribute value',
      };
    }
  }

  /**
   * Удалить значение атрибута
   */
  async deleteAttributeValue(
    valueId: number,
    listingId?: number,
    productId?: number
  ): Promise<ApiResponse<void>> {
    try {
      const endpoint = listingId
        ? `/v2/listings/${listingId}/attributes/${valueId}`
        : `/v2/products/${productId}/attributes/${valueId}`;

      const response = await apiClient.delete<{ message?: string }>(endpoint);

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        message: response.data?.message,
      };
    } catch (error) {
      console.error('Error deleting attribute value:', error);
      return {
        success: false,
        error: 'Failed to delete attribute value',
      };
    }
  }

  /**
   * Пакетное создание/обновление значений атрибутов для объявления
   */
  async saveListingAttributes(
    listingId: number,
    attributeValues: Partial<UnifiedAttributeValue>[]
  ): Promise<ApiResponse<UnifiedAttributeValue[]>> {
    try {
      const response = await apiClient.post<{
        data: UnifiedAttributeValue[];
      }>(`/v2/listings/${listingId}/attributes/batch`, {
        attributes: attributeValues,
      });

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error saving listing attributes:', error);
      return {
        success: false,
        error: 'Failed to save listing attributes',
      };
    }
  }

  /**
   * Пакетное создание/обновление значений атрибутов для продукта
   */
  async saveProductAttributes(
    productId: number,
    attributeValues: Partial<UnifiedAttributeValue>[]
  ): Promise<ApiResponse<UnifiedAttributeValue[]>> {
    try {
      const response = await apiClient.post<{
        data: UnifiedAttributeValue[];
      }>(`/v2/products/${productId}/attributes/batch`, {
        attributes: attributeValues,
      });

      if (response.error) {
        return {
          success: false,
          error: response.error.message,
        };
      }

      return {
        success: true,
        data: response.data?.data,
      };
    } catch (error) {
      console.error('Error saving product attributes:', error);
      return {
        success: false,
        error: 'Failed to save product attributes',
      };
    }
  }

  /**
   * Валидация значений атрибутов
   */
  validateAttributeValue(
    attribute: UnifiedAttribute,
    value: UnifiedAttributeValue
  ): string[] {
    const errors: string[] = [];

    // Проверка обязательности
    if (attribute.is_required) {
      const hasValue =
        value.text_value ||
        value.numeric_value !== undefined ||
        value.boolean_value !== undefined ||
        value.date_value ||
        (value.json_value && value.json_value.length > 0);

      if (!hasValue) {
        errors.push('This field is required');
      }
    }

    // Валидация числовых значений
    if (
      attribute.attribute_type === 'number' &&
      value.numeric_value !== undefined
    ) {
      // Validation rules stored in validation_rules array for now
      // TODO: Update when backend provides proper validation structure
    }

    // Валидация текстовых значений
    if (attribute.attribute_type === 'text' && value.text_value) {
      // Text validation - check length if needed
      // TODO: Update when backend provides proper validation structure
    }

    // Валидация select/multiselect
    if (
      (attribute.attribute_type === 'select' ||
        attribute.attribute_type === 'radio') &&
      value.text_value
    ) {
      // Options are stored as number[] in the current schema
      // TODO: Validate against actual options when backend provides them
    }

    if (attribute.attribute_type === 'multiselect' && value.json_value) {
      // Multiselect validation
      // TODO: Validate against actual options when backend provides them
    }

    // Валидация email
    if (attribute.attribute_type === 'email' && value.text_value) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(value.text_value)) {
        errors.push('Invalid email address');
      }
    }

    // Валидация URL
    if (attribute.attribute_type === 'url' && value.text_value) {
      try {
        new URL(value.text_value);
      } catch {
        errors.push('Invalid URL');
      }
    }

    // Валидация телефона
    if (attribute.attribute_type === 'tel' && value.text_value) {
      const phoneRegex = /^\+?[\d\s\-\(\)]+$/;
      if (!phoneRegex.test(value.text_value)) {
        errors.push('Invalid phone number');
      }
    }

    return errors;
  }

  /**
   * Проверить доступность v2 API (feature flag)
   */
  async checkV2ApiAvailability(): Promise<boolean> {
    try {
      const response = await apiClient.get('/v2/attributes/health');
      return !response.error;
    } catch {
      console.warn('V2 API not available');
      return false;
    }
  }
}

// Экспортируем синглтон
export const unifiedAttributeService = new UnifiedAttributeService();

// Экспортируем типы для использования в компонентах
export type {
  UnifiedAttribute,
  UnifiedAttributeValue,
  ApiResponse,
  GetAttributesParams,
  CreateAttributeValueParams,
  UpdateAttributeValueParams,
};
