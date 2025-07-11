import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

type StorefrontProduct =
  components['schemas']['backend_internal_domain_models.StorefrontProduct'];

export const storefrontProductsService = {
  // Получить товары витрины
  async getProducts(
    storefrontSlug: string,
    params?: {
      limit?: number;
      offset?: number;
      category_id?: number;
      search?: string;
      in_stock_only?: boolean;
    }
  ): Promise<{ products: StorefrontProduct[]; total: number }> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.append('limit', params.limit.toString());
    if (params?.offset) searchParams.append('offset', params.offset.toString());
    if (params?.category_id)
      searchParams.append('category_id', params.category_id.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.in_stock_only) searchParams.append('in_stock_only', 'true');

    const response = await apiClient.get(
      `/api/v1/storefronts/slug/${storefrontSlug}/products?${searchParams}`
    );
    return {
      products: response.data || [],
      total: response.data?.length || 0,
    };
  },

  // Получить товар по ID
  async getProduct(
    storefrontSlug: string,
    productId: number
  ): Promise<StorefrontProduct> {
    const response = await apiClient.get(
      `/api/v1/storefronts/slug/${storefrontSlug}/products/${productId}`
    );
    return response.data;
  },
};
