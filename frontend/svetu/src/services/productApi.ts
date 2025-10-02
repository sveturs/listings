import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

// Типы из сгенерированного API
type BulkCreateProductsRequest =
  components['schemas']['backend_internal_domain_models.BulkCreateProductsRequest'];
type BulkCreateProductsResponse =
  components['schemas']['backend_internal_domain_models.BulkCreateProductsResponse'];
type BulkUpdateProductsRequest =
  components['schemas']['backend_internal_domain_models.BulkUpdateProductsRequest'];
type BulkUpdateProductsResponse =
  components['schemas']['backend_internal_domain_models.BulkUpdateProductsResponse'];
type BulkDeleteProductsRequest =
  components['schemas']['backend_internal_domain_models.BulkDeleteProductsRequest'];
type BulkDeleteProductsResponse =
  components['schemas']['backend_internal_domain_models.BulkDeleteProductsResponse'];
type BulkUpdateStatusRequest =
  components['schemas']['backend_internal_domain_models.BulkUpdateStatusRequest'];
type BulkUpdateStatusResponse =
  components['schemas']['backend_internal_domain_models.BulkUpdateStatusResponse'];

export const productApi = {
  /**
   * Массовое создание товаров
   */
  bulkCreate: async (
    storefrontSlug: string,
    request: BulkCreateProductsRequest
  ) => {
    const response = await apiClient.post<BulkCreateProductsResponse>(
      `/api/v1/storefronts/slug/${storefrontSlug}/products/bulk/create`,
      request
    );
    return response.data;
  },

  /**
   * Массовое обновление товаров
   */
  bulkUpdate: async (
    storefrontSlug: string,
    request: BulkUpdateProductsRequest
  ) => {
    const response = await apiClient.put<BulkUpdateProductsResponse>(
      `/api/v1/storefronts/slug/${storefrontSlug}/products/bulk/update`,
      request
    );
    return response.data;
  },

  /**
   * Массовое удаление товаров
   */
  bulkDelete: async (storefrontSlug: string, productIds: number[]) => {
    const request: BulkDeleteProductsRequest = { product_ids: productIds };
    const response = await apiClient.delete<BulkDeleteProductsResponse>(
      `/api/v1/storefronts/slug/${storefrontSlug}/products/bulk/delete`,
      { data: request }
    );
    return response.data;
  },

  /**
   * Массовое обновление статуса товаров
   */
  bulkUpdateStatus: async (
    storefrontSlug: string,
    productIds: number[],
    isActive: boolean
  ) => {
    const request: BulkUpdateStatusRequest = {
      product_ids: productIds,
      is_active: isActive,
    };
    const response = await apiClient.put<BulkUpdateStatusResponse>(
      `/api/v1/storefronts/slug/${storefrontSlug}/products/bulk/status`,
      request
    );
    return response.data;
  },

  /**
   * Экспорт товаров в CSV
   */
  exportToCSV: async (storefrontSlug: string, productIds?: number[]) => {
    const params = productIds?.length ? `?ids=${productIds.join(',')}` : '';
    // Используем BFF proxy - автоматически добавит cookies для авторизации
    const response = await fetch(
      `/api/v2/storefronts/slug/${storefrontSlug}/products/export/csv${params}`,
      {
        credentials: 'include', // Включаем cookies
      }
    );

    if (!response.ok) {
      throw new Error('Export failed');
    }

    const blob = await response.blob();

    // Создаем ссылку для скачивания
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute(
      'download',
      `products_${storefrontSlug}_${Date.now()}.csv`
    );
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  },

  /**
   * Экспорт товаров в XML
   */
  exportToXML: async (storefrontSlug: string, productIds?: number[]) => {
    const params = productIds?.length ? `?ids=${productIds.join(',')}` : '';
    // Используем BFF proxy - автоматически добавит cookies для авторизации
    const response = await fetch(
      `/api/v2/storefronts/slug/${storefrontSlug}/products/export/xml${params}`,
      {
        credentials: 'include', // Включаем cookies
      }
    );

    if (!response.ok) {
      throw new Error('Export failed');
    }

    const blob = await response.blob();

    // Создаем ссылку для скачивания
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute(
      'download',
      `products_${storefrontSlug}_${Date.now()}.xml`
    );
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  },
};
