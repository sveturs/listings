import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['models.UnifiedAttribute'];
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
  private baseUrl: string;
  private apiVersion: string = 'v2'; // Используем v2 API для унифицированных атрибутов

  constructor() {
    this.baseUrl =
      process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3000/api';
  }

  /**
   * Получить список атрибутов для категории
   */
  async getCategoryAttributes(
    categoryId: number
  ): Promise<ApiResponse<UnifiedAttribute[]>> {
    try {
      const response = await fetch(
        `${this.baseUrl}/${this.apiVersion}/attributes/category/${categoryId}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        }
      );

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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

      const url = `${this.baseUrl}/${this.apiVersion}/attributes${
        queryParams.toString() ? `?${queryParams.toString()}` : ''
      }`;

      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
      const response = await fetch(
        `${this.baseUrl}/${this.apiVersion}/attributes/${attributeId}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        }
      );

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
      const response = await fetch(
        `${this.baseUrl}/${this.apiVersion}/listings/${listingId}/attributes`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        }
      );

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
      const response = await fetch(
        `${this.baseUrl}/${this.apiVersion}/products/${productId}/attributes`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        }
      );

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
        ? `${this.baseUrl}/${this.apiVersion}/listings/${params.listing_id}/attributes`
        : `${this.baseUrl}/${this.apiVersion}/products/${params.product_id}/attributes`;

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          attribute_id: params.attribute_id,
          text_value: params.text_value,
          numeric_value: params.numeric_value,
          boolean_value: params.boolean_value,
          date_value: params.date_value,
          json_value: params.json_value,
        }),
      });

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
        ? `${this.baseUrl}/${this.apiVersion}/listings/${params.listing_id}/attributes/${params.id}`
        : `${this.baseUrl}/${this.apiVersion}/products/${params.product_id}/attributes/${params.id}`;

      const response = await fetch(endpoint, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          text_value: params.text_value,
          numeric_value: params.numeric_value,
          boolean_value: params.boolean_value,
          date_value: params.date_value,
          json_value: params.json_value,
        }),
      });

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
        ? `${this.baseUrl}/${this.apiVersion}/listings/${listingId}/attributes/${valueId}`
        : `${this.baseUrl}/${this.apiVersion}/products/${productId}/attributes/${valueId}`;

      const response = await fetch(endpoint, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      const data = await response.json();
      return {
        success: response.ok,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
      const response = await fetch(
        `${this.baseUrl}/${this.apiVersion}/listings/${listingId}/attributes/batch`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({ attributes: attributeValues }),
        }
      );

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
      const response = await fetch(
        `${this.baseUrl}/${this.apiVersion}/products/${productId}/attributes/batch`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({ attributes: attributeValues }),
        }
      );

      const data = await response.json();
      return {
        success: response.ok,
        data: response.ok ? data.data : undefined,
        message: data.message,
        error: !response.ok ? data.error : undefined,
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
      const response = await fetch(
        `${this.baseUrl}/${this.apiVersion}/attributes/health`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        }
      );
      return response.ok;
    } catch {
      console.warn('V2 API not available, falling back to v1');
      return false;
    }
  }

  /**
   * Переключиться на v1 API (fallback)
   */
  useV1Api() {
    this.apiVersion = 'v1';
  }

  /**
   * Переключиться на v2 API
   */
  useV2Api() {
    this.apiVersion = 'v2';
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
