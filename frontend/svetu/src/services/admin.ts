import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

async function getAuthHeaders(): Promise<Record<string, string>> {
  // Helper function to get headers with auth and CSRF tokens
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (typeof window !== 'undefined') {
    // Get auth token
    try {
      const { tokenManager } = await import('@/utils/tokenManager');
      const token = await tokenManager.getAccessToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
        console.log('[getAuthHeaders] Auth token added');
      }
    } catch {
      console.log('No auth token available');
    }
  }

  console.log('[getAuthHeaders] Final headers:', headers);
  return headers;
}

// Types
export interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number;
  icon?: string;
  description?: string;
  is_active?: boolean;
  seo_title?: string;
  seo_description?: string;
  seo_keywords?: string;
  created_at?: string;
  updated_at?: string;
  children?: Category[];
  items_count?: number;
}

export interface Attribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type:
    | 'text'
    | 'number'
    | 'select'
    | 'boolean'
    | 'date'
    | 'range'
    | 'location'
    | 'file'
    | 'gallery';
  icon?: string;
  options?: string[];
  validation_rules?: Record<string, unknown>;
  is_searchable?: boolean;
  is_filterable?: boolean;
  is_required?: boolean;
  sort_order?: number;
  created_at?: string;
  unit?: string;
  min_value?: number;
  max_value?: number;
  min_length?: number;
  max_length?: number;
  pattern?: string;
  default_value?: string | number | boolean;
  show_in_card?: boolean;
  show_in_list?: boolean;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
}

export interface AttributeGroup {
  id: number;
  name: string;
  display_name: string;
  description?: string;
  icon?: string;
  sort_order?: number;
  is_active?: boolean;
  is_system?: boolean;
  created_at?: string;
  updated_at?: string;
  attributes?: AttributeGroupItem[];
}

export interface AttributeGroupItem {
  id: number;
  group_id: number;
  attribute_id: number;
  icon?: string;
  sort_order?: number;
  custom_display_name?: string;
  visibility_condition?: {
    attribute_id: number;
    operator: string;
    value: string | number | boolean;
  };
  attribute?: Attribute;
}

// Use generated PaginatedResponse type
type PaginatedResponse<T> =
  components['schemas']['models.PaginatedResponse'] & {
    data: T[];
  };

export interface CategoryAttributeMapping {
  category_id: number;
  attribute_id: number;
  is_enabled?: boolean;
  is_required?: boolean;
  sort_order?: number;
  custom_component?: string;
}

// Admin API Service
export const adminApi = {
  // Categories
  categories: {
    async getAll(): Promise<Category[]> {
      try {
        console.log(
          'adminApi.categories.getAll: Making request to admin categories API'
        );

        // Получаем токен если есть
        const headers: Record<string, string> = {
          'Content-Type': 'application/json',
        };

        if (typeof window !== 'undefined') {
          try {
            const { tokenManager } = await import('@/utils/tokenManager');
            const token = await tokenManager.getAccessToken();
            if (token) {
              headers['Authorization'] = `Bearer ${token}`;
            }
          } catch {
            console.log('No token available, proceeding without auth');
          }
        }

        const response = await fetch('/api/v1/admin/categories', {
          method: 'GET',
          headers,
          credentials: 'include',
        });

        console.log('Response status:', response.status);

        if (!response.ok) {
          console.error(
            'Response not OK:',
            response.status,
            response.statusText
          );
          return [];
        }

        const data = await response.json();
        console.log('Response data:', data);

        // Backend возвращает данные в поле data
        const result = data.data || data || [];
        console.log('Final categories array:', result);
        console.log(
          'Categories count:',
          Array.isArray(result) ? result.length : 'not array'
        );
        return Array.isArray(result) ? result : [];
      } catch (error) {
        console.error('Failed to load categories:', error);
        return [];
      }
    },

    async getById(id: number): Promise<Category> {
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      };

      if (typeof window !== 'undefined') {
        try {
          const { tokenManager } = await import('@/utils/tokenManager');
          const token = await tokenManager.getAccessToken();
          if (token) {
            headers['Authorization'] = `Bearer ${token}`;
          }
        } catch {
          console.log('No token available, proceeding without auth');
        }
      }

      const response = await fetch(`/api/v1/admin/categories/${id}`, {
        method: 'GET',
        headers,
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      // Backend returns wrapped response with success and data fields
      return data.data;
    },

    async create(
      category: Partial<Category>
    ): Promise<{ id: number; message: string }> {
      const headers = await getAuthHeaders();

      const response = await fetch('/api/v1/admin/categories', {
        method: 'POST',
        headers,
        body: JSON.stringify(category),
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    },

    async update(id: number, category: Partial<Category>): Promise<any> {
      const headers = await getAuthHeaders();

      const response = await fetch(`/api/v1/admin/categories/${id}`, {
        method: 'PUT',
        headers,
        body: JSON.stringify(category),
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    },

    async delete(id: number): Promise<any> {
      const headers = await getAuthHeaders();

      const response = await fetch(`/api/v1/admin/categories/${id}`, {
        method: 'DELETE',
        headers,
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    },

    async reorder(orderedIds: number[]): Promise<any> {
      const response = await apiClient.post(
        '/api/v1/admin/categories/reorder',
        {
          ordered_ids: orderedIds,
        }
      );
      return response.data as any as { message: string };
    },

    async move(id: number, newParentId: number): Promise<any> {
      const response = await apiClient.put(
        `/api/v1/admin/categories/${id}/move`,
        {
          new_parent_id: newParentId,
        }
      );
      return response.data as any as { message: string };
    },

    // Category attributes
    async getAttributes(categoryId: number): Promise<Attribute[]> {
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      };

      if (typeof window !== 'undefined') {
        try {
          const { tokenManager } = await import('@/utils/tokenManager');
          const token = await tokenManager.getAccessToken();
          if (token) {
            headers['Authorization'] = `Bearer ${token}`;
          }
        } catch {
          console.log('No token available, proceeding without auth');
        }
      }

      const response = await fetch(
        `/api/v1/marketplace/categories/${categoryId}/attributes`,
        {
          method: 'GET',
          headers,
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      // Backend returns attributes in data field
      return data.data || [];
    },

    async addAttribute(
      categoryId: number,
      attributeId: number,
      isRequired: boolean = false
    ): Promise<any> {
      const response = await apiClient.post(
        `/api/v1/admin/categories/${categoryId}/attributes`,
        {
          attribute_id: attributeId,
          is_required: isRequired,
        }
      );
      return response.data as any;
    },

    async removeAttribute(
      categoryId: number,
      attributeId: number
    ): Promise<any> {
      const response = await apiClient.delete(
        `/api/v1/admin/categories/${categoryId}/attributes/${attributeId}`
      );
      return response.data as any;
    },

    async updateAttributeSettings(
      categoryId: number,
      attributeId: number,
      settings: Partial<CategoryAttributeMapping>
    ): Promise<any> {
      const response = await apiClient.put(
        `/api/v1/admin/categories/${categoryId}/attributes/${attributeId}`,
        settings
      );
      return response.data as any;
    },
  },

  // Attributes
  attributes: {
    async getAll(
      page: number = 1,
      pageSize: number = 20,
      search?: string,
      filterType?: string
    ): Promise<PaginatedResponse<Attribute>> {
      try {
        console.log(
          'adminApi.attributes.getAll: Making request to admin attributes API with pagination',
          { page, pageSize, search, filterType }
        );

        // Получаем токен если есть
        const headers: Record<string, string> = {
          'Content-Type': 'application/json',
        };

        if (typeof window !== 'undefined') {
          try {
            const { tokenManager } = await import('@/utils/tokenManager');
            const token = await tokenManager.getAccessToken();
            if (token) {
              headers['Authorization'] = `Bearer ${token}`;
            }
          } catch {
            console.log('No token available, proceeding without auth');
          }
        }

        const params = new URLSearchParams({
          page: page.toString(),
          page_size: pageSize.toString(),
        });

        // Добавляем параметры поиска и фильтрации
        if (search) {
          params.append('search', search);
        }
        if (filterType) {
          params.append('type', filterType);
        }

        const response = await fetch(
          `/api/v1/admin/attributes?${params.toString()}`,
          {
            method: 'GET',
            headers,
            credentials: 'include',
          }
        );

        console.log('Response status:', response.status);

        if (!response.ok) {
          console.error(
            'Response not OK:',
            response.status,
            response.statusText
          );
          return {
            data: [],
            page: 1,
            page_size: pageSize,
            total: 0,
            total_pages: 0,
          };
        }

        const result = await response.json();
        console.log('Response data:', result);

        // Backend возвращает обёрнутый ответ с полем data
        // Проверяем структуру ответа
        if (
          result.data &&
          result.data.data &&
          typeof result.data.page !== 'undefined'
        ) {
          // Ответ уже содержит структуру пагинации
          const paginatedData = result.data;
          return {
            data: paginatedData.data || [],
            page: paginatedData.page || page,
            page_size: paginatedData.page_size || pageSize,
            total: paginatedData.total || 0,
            total_pages: paginatedData.total_pages || 0,
          };
        } else if (Array.isArray(result.data)) {
          // Старый формат ответа - простой массив
          return {
            data: result.data,
            page: 1,
            page_size: result.data.length,
            total: result.data.length,
            total_pages: 1,
          };
        } else {
          // Неизвестный формат
          return {
            data: [],
            page: 1,
            page_size: pageSize,
            total: 0,
            total_pages: 0,
          };
        }
      } catch (error) {
        console.error('Failed to load attributes:', error);
        return {
          data: [],
          page: 1,
          page_size: pageSize,
          total: 0,
          total_pages: 0,
        };
      }
    },

    async getById(id: number): Promise<Attribute> {
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      };

      if (typeof window !== 'undefined') {
        try {
          const { tokenManager } = await import('@/utils/tokenManager');
          const token = await tokenManager.getAccessToken();
          if (token) {
            headers['Authorization'] = `Bearer ${token}`;
          }
        } catch {
          console.log('No token available, proceeding without auth');
        }
      }

      const response = await fetch(`/api/v1/admin/attributes/${id}`, {
        method: 'GET',
        headers,
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      // Backend returns wrapped response with success and data fields
      return data.data;
    },

    async create(
      attribute: Partial<Attribute>
    ): Promise<{ id: number; message: string }> {
      const headers = await getAuthHeaders();

      const response = await fetch('/api/v1/admin/attributes', {
        method: 'POST',
        headers,
        body: JSON.stringify(attribute),
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    },

    async update(id: number, attribute: Partial<Attribute>): Promise<any> {
      console.log(
        '[adminApi.attributes.update] Starting update for attribute:',
        id,
        attribute
      );
      const headers = await getAuthHeaders();

      const response = await fetch(`/api/v1/admin/attributes/${id}`, {
        method: 'PUT',
        headers,
        body: JSON.stringify(attribute),
        credentials: 'include',
      });

      console.log(
        '[adminApi.attributes.update] Response status:',
        response.status
      );

      if (!response.ok) {
        const errorText = await response.text();
        console.error(
          '[adminApi.attributes.update] Error response:',
          errorText
        );
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log('[adminApi.attributes.update] Success response:', data);
      return data;
    },

    async delete(id: number): Promise<any> {
      const headers = await getAuthHeaders();

      const response = await fetch(`/api/v1/admin/attributes/${id}`, {
        method: 'DELETE',
        headers,
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    },

    async bulkUpdate(attributes: Attribute[]): Promise<{
      success_count: number;
      total_count: number;
      errors?: string[];
    }> {
      const response = await apiClient.post(
        '/api/v1/admin/attributes/bulk-update',
        { attributes }
      );
      return response.data as any as {
        success_count: number;
        total_count: number;
        errors?: string[];
      };
    },
  },

  // Attribute Groups
  attributeGroups: {
    async getAll(): Promise<AttributeGroup[]> {
      const response = await apiClient.get('/api/v1/admin/attribute-groups');
      return (response.data as any as any).groups || [];
    },

    async getById(id: number): Promise<AttributeGroup> {
      const response = await apiClient.get(
        `/api/v1/admin/attribute-groups/${id}`
      );
      return (response.data as any as any).group;
    },

    async create(
      group: Partial<AttributeGroup>
    ): Promise<{ id: number; success: boolean }> {
      const response = await apiClient.post(
        '/api/v1/admin/attribute-groups',
        group
      );
      return response.data as any;
    },

    async update(
      id: number,
      group: Partial<AttributeGroup>
    ): Promise<{ success: boolean; message: string }> {
      const response = await apiClient.put(
        `/api/v1/admin/attribute-groups/${id}`,
        group
      );
      return response.data as any;
    },

    async delete(id: number): Promise<{ success: boolean; message: string }> {
      const response = await apiClient.delete(
        `/api/v1/admin/attribute-groups/${id}`
      );
      return response.data as any;
    },

    async getWithItems(
      id: number
    ): Promise<{ group: AttributeGroup; items: AttributeGroupItem[] }> {
      const response = await apiClient.get(
        `/api/v1/admin/attribute-groups/${id}/items`
      );
      return (response.data as any).data;
    },

    async addItem(
      groupId: number,
      item: Partial<AttributeGroupItem>
    ): Promise<{ id: number; success: boolean }> {
      const response = await apiClient.post(
        `/api/v1/admin/attribute-groups/${groupId}/items`,
        item
      );
      return response.data as any;
    },

    async removeItem(
      groupId: number,
      attributeId: number
    ): Promise<{ success: boolean; message: string }> {
      const response = await apiClient.delete(
        `/api/v1/admin/attribute-groups/${groupId}/items/${attributeId}`
      );
      return response.data as any;
    },
  },

  // Translation
  async translate(
    text: string,
    targetLanguages: string[] = ['ru', 'en', 'sr']
  ): Promise<Record<string, string>> {
    // Backend принимает только один язык за раз, поэтому делаем несколько запросов
    const translations: Record<string, string> = {};

    for (const targetLang of targetLanguages) {
      try {
        const response = await apiClient.post(
          '/api/v1/marketplace/translations/translate',
          {
            text,
            source_lang: 'ru', // Предполагаем, что исходный текст на русском
            target_lang: targetLang,
            provider: 'google',
          }
        );

        const data = response.data as any;
        if (data && data.data && data.data.translated_text) {
          translations[targetLang] = data.data.translated_text;
        }
      } catch (error) {
        console.error(`Failed to translate to ${targetLang}:`, error);
        // Продолжаем с другими языками даже если один не удался
      }
    }

    return translations;
  },

  // Translate Attribute
  async translateAttribute(
    attributeId: number,
    sourceLanguage: string = 'en',
    targetLanguages: string[] = ['ru', 'sr']
  ): Promise<{
    attribute_id: number;
    translations: Record<string, any>;
    errors?: string[];
  }> {
    const response = await apiClient.post(
      `/api/v1/admin/attributes/${attributeId}/translate`,
      {
        source_language: sourceLanguage,
        target_languages: targetLanguages,
      }
    );
    return response.data as any;
  },
};
