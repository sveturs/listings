import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

function getCurrentLocale(): string {
  // Проверяем URL для получения локали (Next.js i18n)
  if (typeof window !== 'undefined') {
    const pathSegments = window.location.pathname.split('/');
    const locale = pathSegments[1];
    if (['en', 'ru', 'sr'].includes(locale)) {
      return locale;
    }
  }
  // По умолчанию возвращаем сербский
  return 'sr';
}

// УДАЛЕНО: getAuthHeaders больше не нужен
// JWT токены автоматически добавляются BFF прокси /api/v2 из httpOnly cookies

// Types
export interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number;
  parent_name?: string;
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
  listing_count?: number;
  translations?: {
    [key: string]: string;
  };
}

export interface Attribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type:
    | 'text'
    | 'number'
    | 'select'
    | 'multiselect'
    | 'boolean'
    | 'date'
    | 'range'
    | 'location'
    | 'file'
    | 'gallery';
  icon?: string;
  options?: string[];
  validation_rules?: Record<string, unknown>;
  custom_component?: string;
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
  is_variant_compatible?: boolean;
  affects_stock?: boolean;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
}

export interface AttributeGroup {
  id: number;
  code: string;
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
  group_id?: number;
}

// Admin API Service
export const adminApi = {
  // Categories
  categories: {
    async getAll(): Promise<Category[]> {
      try {
        const currentLocale = getCurrentLocale();
        const response = await apiClient.get(
          `/admin/categories-all?lang=${currentLocale}`
        );

        // Backend возвращает данные в поле data
        const result = response.data?.data || response.data || [];
        return Array.isArray(result) ? result : [];
      } catch (error) {
        console.error('Failed to load categories:', error);
        return [];
      }
    },

    async getById(id: number): Promise<Category> {
      const response = await apiClient.get(`/admin/categories/${id}`);

      if (response.error) {
        throw new Error(response.error.message);
      }

      // Backend returns wrapped response with success and data fields
      return response.data?.data || response.data;
    },

    async create(
      category: Partial<Category>
    ): Promise<{ id: number; message: string }> {
      const response = await apiClient.request('/admin/categories', {
        method: 'POST',
        body: JSON.stringify(category),
      });

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data as any;
    },

    async update(id: number, category: Partial<Category>): Promise<any> {
      const response = await apiClient.request(`/admin/categories/${id}`, {
        method: 'PUT',
        body: JSON.stringify(category),
      });

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data;
    },

    async delete(id: number): Promise<any> {
      const response = await apiClient.request(`/admin/categories/${id}`, {
        method: 'DELETE',
      });

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data;
    },

    async reorder(orderedIds: number[]): Promise<any> {
      const response = await apiClient.post('/admin/categories/reorder', {
        ordered_ids: orderedIds,
      });
      return response.data as any as { message: string };
    },

    async move(id: number, newParentId: number): Promise<any> {
      const response = await apiClient.put(`/admin/categories/${id}/move`, {
        new_parent_id: newParentId,
      });
      return response.data as any as { message: string };
    },

    // Variant attributes for category
    async getVariantAttributes(categoryId: number): Promise<any[]> {
      const response = await apiClient.get(
        `/admin/categories/${categoryId}/variant-attributes`
      );
      // API возвращает { success: true, data: [...] }
      return response.data?.data || [];
    },

    async updateVariantAttributes(
      categoryId: number,
      data: { variant_attributes: any[] }
    ): Promise<void> {
      await apiClient.put(
        `/admin/categories/${categoryId}/variant-attributes`,
        data
      );
    },

    // Category attributes
    async getAttributes(categoryId: number): Promise<Attribute[]> {
      const currentLocale = getCurrentLocale();
      const response = await apiClient.get(
        `/c2c/categories/${categoryId}/attributes?lang=${currentLocale}`
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      // Backend returns attributes in data field
      return response.data?.data || response.data || [];
    },

    async addAttribute(
      categoryId: number,
      attributeId: number,
      isRequired: boolean = false
    ): Promise<any> {
      const response = await apiClient.post(
        `/admin/categories/${categoryId}/attributes`,
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
        `/admin/categories/${categoryId}/attributes/${attributeId}`
      );
      return response.data as any;
    },

    async updateAttributeSettings(
      categoryId: number,
      attributeId: number,
      settings: Partial<CategoryAttributeMapping>
    ): Promise<any> {
      const response = await apiClient.put(
        `/admin/categories/${categoryId}/attributes/${attributeId}`,
        settings
      );
      return response.data as any;
    },

    // Category groups
    async getGroups(categoryId: number): Promise<AttributeGroup[]> {
      const response = await apiClient.get(
        `/admin/categories/${categoryId}/groups`
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data?.data || response.data || [];
    },

    async attachGroup(
      categoryId: number,
      groupId: number,
      sortOrder: number = 0
    ): Promise<any> {
      const response = await apiClient.post(
        `/admin/categories/${categoryId}/groups`,
        {
          group_id: groupId,
          sort_order: sortOrder,
        }
      );
      return response.data as any;
    },

    async detachGroup(categoryId: number, groupId: number): Promise<any> {
      const response = await apiClient.delete(
        `/admin/categories/${categoryId}/groups/${groupId}`
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

        const result = await apiClient.get(
          `/admin/attributes?${params.toString()}`
        );

        if (result.error) {
          console.error('API error:', result.error);
          return {
            data: [],
            page: 1,
            page_size: pageSize,
            total: 0,
            total_pages: 0,
          };
        }

        // Backend возвращает обёрнутый ответ с полем data
        // apiClient уже распаковал внешний wrapper, проверяем структуру
        const responseData = result.data;

        // Проверяем структуру ответа - может быть двойная вложенность
        if (
          responseData &&
          responseData.data &&
          typeof responseData.data.page !== 'undefined'
        ) {
          // Backend вернул {success: true, data: {data: [...], page: 1, ...}}
          const paginatedData = responseData.data;
          return {
            data: paginatedData.data || [],
            page: paginatedData.page || page,
            page_size: paginatedData.page_size || pageSize,
            total: paginatedData.total || 0,
            total_pages: paginatedData.total_pages || 0,
          };
        } else if (
          responseData &&
          typeof responseData.page !== 'undefined' &&
          Array.isArray(responseData.data)
        ) {
          // Ответ уже содержит структуру пагинации напрямую
          return {
            data: responseData.data || [],
            page: responseData.page || page,
            page_size: responseData.page_size || pageSize,
            total: responseData.total || 0,
            total_pages: responseData.total_pages || 0,
          };
        } else if (Array.isArray(responseData)) {
          // Старый формат ответа - простой массив
          return {
            data: responseData,
            page: 1,
            page_size: responseData.length,
            total: responseData.length,
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
      const response = await apiClient.get(`/admin/attributes/${id}`);

      if (response.error) {
        throw new Error(response.error.message);
      }

      // Backend returns wrapped response with success and data fields
      return response.data?.data || response.data;
    },

    async create(
      attribute: Partial<Attribute>
    ): Promise<{ id: number; message: string }> {
      const response = await apiClient.request('/admin/attributes', {
        method: 'POST',
        body: JSON.stringify(attribute),
      });

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data as any;
    },

    async update(id: number, attribute: Partial<Attribute>): Promise<any> {
      const response = await apiClient.request(`/admin/attributes/${id}`, {
        method: 'PUT',
        body: JSON.stringify(attribute),
      });

      if (response.error) {
        console.error(
          '[adminApi.attributes.update] Error response:',
          response.error
        );
        throw new Error(response.error.message);
      }
      return response.data;
    },

    async delete(id: number): Promise<any> {
      const response = await apiClient.request(`/admin/attributes/${id}`, {
        method: 'DELETE',
      });

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data;
    },

    async bulkUpdate(attributes: Attribute[]): Promise<{
      success_count: number;
      total_count: number;
      errors?: string[];
    }> {
      const response = await apiClient.post('/admin/attributes/bulk-update', {
        attributes,
      });
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
      const response = await apiClient.get('/admin/attribute-groups');
      return response.data?.data?.groups || [];
    },

    async getById(id: number): Promise<AttributeGroup> {
      const response = await apiClient.get(`/admin/attribute-groups/${id}`);
      return response.data?.data?.group;
    },

    async create(
      group: Partial<AttributeGroup>
    ): Promise<{ id: number; message: string }> {
      const response = await apiClient.post('/admin/attribute-groups', group);
      return response.data?.data;
    },

    async update(
      id: number,
      group: Partial<AttributeGroup>
    ): Promise<{ message: string }> {
      const response = await apiClient.put(
        `/admin/attribute-groups/${id}`,
        group
      );
      return response.data?.data;
    },

    async delete(id: number): Promise<{ message: string }> {
      const response = await apiClient.delete(`/admin/attribute-groups/${id}`);
      return response.data?.data;
    },

    async getWithItems(
      id: number
    ): Promise<{ group: AttributeGroup; items: AttributeGroupItem[] }> {
      const response = await apiClient.get(
        `/admin/attribute-groups/${id}/items`
      );
      return (response.data as any).data;
    },

    async addItem(
      groupId: number,
      item: Partial<AttributeGroupItem>
    ): Promise<{ id: number; success: boolean }> {
      const response = await apiClient.post(
        `/admin/attribute-groups/${groupId}/items`,
        item
      );
      return response.data as any;
    },

    async removeItem(
      groupId: number,
      attributeId: number
    ): Promise<{ success: boolean; message: string }> {
      const response = await apiClient.delete(
        `/admin/attribute-groups/${groupId}/items/${attributeId}`
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
        const response = await apiClient.post('/c2c/translations/translate', {
          text,
          source_lang: 'ru', // Предполагаем, что исходный текст на русском
          target_lang: targetLang,
          provider: 'google',
        });

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
      `/admin/attributes/${attributeId}/translate`,
      {
        source_language: sourceLanguage,
        target_languages: targetLanguages,
      }
    );
    return response.data as any;
  },

  // Get Translation Status
  async getTranslationStatus(
    entityType: 'category' | 'attribute',
    entityIds?: number[]
  ): Promise<components['schemas']['handler.TranslationStatusItem'][]> {
    const params = new URLSearchParams();
    params.append('entity_type', entityType);
    if (entityIds && entityIds.length > 0) {
      params.append('entity_ids', entityIds.join(','));
    }

    const response = await apiClient.get(
      `/admin/c2c-translations/status?${params.toString()}`
    );
    return (response.data as any).data || [];
  },

  // Batch Translate Categories
  async batchTranslateCategories(
    categoryIds: number[],
    sourceLanguage: string = 'ru',
    targetLanguages: string[] = ['en', 'sr']
  ): Promise<{
    translations: Record<string, any>;
    errors?: string[];
  }> {
    const response = await apiClient.post(
      '/admin/translations/batch-categories',
      {
        category_ids: categoryIds,
        source_language: sourceLanguage,
        target_languages: targetLanguages,
      }
    );
    return response.data as any;
  },

  // Batch Translate Attributes
  async batchTranslateAttributes(
    attributeIds: number[],
    sourceLanguage: string = 'ru',
    targetLanguages: string[] = ['en', 'sr']
  ): Promise<{
    translations: Record<string, any>;
    errors?: string[];
  }> {
    const response = await apiClient.post(
      '/admin/translations/batch-attributes',
      {
        attribute_ids: attributeIds,
        source_language: sourceLanguage,
        target_languages: targetLanguages,
      }
    );
    return response.data as any;
  },

  // Translate Category (single)
  async translateCategory(
    categoryId: number,
    sourceLanguage: string = 'ru',
    targetLanguages: string[] = ['en', 'sr']
  ): Promise<{
    category_id: number;
    translations: Record<string, any>;
    errors?: string[];
  }> {
    const response = await apiClient.post(
      `/admin/categories/${categoryId}/translate`,
      {
        source_language: sourceLanguage,
        target_languages: targetLanguages,
      }
    );
    return response.data as any;
  },

  // Update Field Translation
  async updateFieldTranslation(
    entityType: 'category' | 'attribute',
    entityId: number,
    fieldName: string,
    translations: Record<string, string>,
    provider: string = 'manual'
  ): Promise<
    Record<string, components['schemas']['handler.TranslationFieldStatus']>
  > {
    const response = await apiClient.put(
      `/admin/translations/${entityType}/${entityId}/${fieldName}`,
      {
        translations,
        provider,
      }
    );
    return (response.data as any).data || {};
  },

  // Variant Attributes API
  variantAttributes: {
    async getAll(
      page: number = 1,
      limit: number = 20,
      search?: string,
      type?: string
    ): Promise<{
      data: VariantAttribute[];
      total: number;
      page: number;
      total_pages: number;
    }> {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
      });

      if (search) params.append('search', search);
      if (type) params.append('type', type);

      const response = await apiClient.request(
        `/admin/variant-attributes?${params.toString()}`,
        { method: 'GET' }
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data as any;
    },

    async getById(id: number): Promise<VariantAttribute> {
      const response = await apiClient.request(
        `/admin/variant-attributes/${id}`,
        { method: 'GET' }
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return (response.data as any)?.data || response.data;
    },

    async create(
      attribute: Partial<VariantAttribute>
    ): Promise<{ id: number; message: string }> {
      const response = await apiClient.request('/admin/variant-attributes', {
        method: 'POST',
        body: JSON.stringify(attribute),
      });

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data as any;
    },

    async update(
      id: number,
      attribute: Partial<VariantAttribute>
    ): Promise<any> {
      const response = await apiClient.request(
        `/admin/variant-attributes/${id}`,
        {
          method: 'PUT',
          body: JSON.stringify(attribute),
        }
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data;
    },

    async delete(id: number): Promise<{ message: string }> {
      const response = await apiClient.request(
        `/admin/variant-attributes/${id}`,
        { method: 'DELETE' }
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data as any;
    },

    // API методы для управления связями между атрибутами
    async getVariantAttributeMappings(
      variantAttributeId: number
    ): Promise<any[]> {
      const response = await apiClient.request(
        `/admin/variant-attributes/${variantAttributeId}/mappings`,
        { method: 'GET' }
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return (response.data as any)?.data || [];
    },

    async updateVariantAttributeMappings(
      variantAttributeId: number,
      categoryAttributeIds: number[]
    ): Promise<any> {
      const response = await apiClient.request(
        `/admin/variant-attributes/${variantAttributeId}/mappings`,
        {
          method: 'PUT',
          body: JSON.stringify({
            category_attribute_ids: categoryAttributeIds,
          }),
        }
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data;
    },
  },
};

// Additional types for variant attributes
export interface VariantAttribute {
  id: number;
  name: string;
  display_name: string;
  type:
    | 'text'
    | 'number'
    | 'select'
    | 'multiselect'
    | 'boolean'
    | 'date'
    | 'range';
  is_required: boolean;
  sort_order: number;
  affects_stock: boolean;
  created_at: string;
  updated_at: string;
}
