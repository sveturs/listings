import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

type StorefrontProduct =
  components['schemas']['backend_internal_domain_models.StorefrontProduct'];
type UpdateProductRequest =
  components['schemas']['backend_internal_domain_models.UpdateProductRequest'];
type CreateProductRequest =
  components['schemas']['backend_internal_domain_models.CreateProductRequest'];

// Расширенный тип для создания товара с вариантами
interface CreateProductWithVariantsRequest extends CreateProductRequest {
  has_variants?: boolean;
  variants?: Array<{
    sku?: string;
    barcode?: string;
    price?: number;
    compare_at_price?: number;
    cost_price?: number;
    stock_quantity: number;
    low_stock_threshold?: number;
    variant_attributes: Record<string, any>;
    weight?: number;
    dimensions?: Record<string, any>;
    is_default: boolean;
  }>;
  variant_settings?: {
    track_inventory: boolean;
    continue_selling: boolean;
    require_shipping: boolean;
    taxable_product: boolean;
    weight_unit?: string;
    selected_attributes: string[];
  };
}

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

    if (response.error) {
      throw new Error(response.error.message || 'Failed to fetch product');
    }

    if (!response.data) {
      throw new Error('Product not found');
    }

    return response.data;
  },

  // Создать товар (с поддержкой вариантов)
  async createProduct(
    storefrontSlug: string,
    productData: CreateProductRequest | CreateProductWithVariantsRequest
  ): Promise<StorefrontProduct> {
    const response = await apiClient.post(
      `/api/v1/storefronts/${storefrontSlug}/products`,
      productData
    );
    return response.data;
  },

  // Обновить товар
  async updateProduct(
    storefrontSlug: string,
    productId: number,
    productData: UpdateProductRequest
  ): Promise<StorefrontProduct> {
    const response = await apiClient.put(
      `/api/v1/storefronts/${storefrontSlug}/products/${productId}`,
      productData
    );
    return response.data;
  },

  // Удалить товар
  async deleteProduct(
    storefrontSlug: string,
    productId: number
  ): Promise<void> {
    await apiClient.delete(
      `/api/v1/storefronts/${storefrontSlug}/products/${productId}`
    );
  },

  // Загрузить изображения товара
  async uploadProductImages(
    storefrontSlug: string,
    productId: number,
    images: File[],
    mainImageIndex?: number
  ): Promise<{ uploaded: any[] }> {
    const uploaded: any[] = [];

    console.log('uploadProductImages called with:', {
      storefrontSlug,
      productId,
      imagesCount: images.length,
      mainImageIndex,
      images: images.map((f) => ({ name: f.name, size: f.size, type: f.type })),
    });

    // Загружаем изображения по одному
    for (let index = 0; index < images.length; index++) {
      const image = images[index];
      const formData = new FormData();
      formData.append('image', image); // Backend ожидает 'image', не 'images'

      // Устанавливаем главное изображение если это нужный индекс
      if (index === mainImageIndex) {
        formData.append('is_main', 'true');
      } else {
        formData.append('is_main', 'false');
      }

      // Устанавливаем порядок отображения
      formData.append('display_order', index.toString());

      console.log(`Uploading image ${index}:`, {
        fileName: image.name,
        fileSize: image.size,
        fileType: image.type,
        isMain: index === mainImageIndex,
        displayOrder: index,
      });

      try {
        // НЕ устанавливаем Content-Type вручную - браузер сделает это сам с правильным boundary
        const response = await apiClient.post(
          `/api/v1/storefronts/${storefrontSlug}/products/${productId}/images`,
          formData
        );
        console.log(`Image ${index} uploaded successfully:`, response.data);
        uploaded.push(response.data);
      } catch (error: any) {
        console.error(`Failed to upload image ${index}:`, error);
        console.error('Error response:', error.response?.data);
        throw error;
      }
    }

    return { uploaded };
  },

  // Удалить изображение товара
  async deleteProductImage(
    storefrontSlug: string,
    productId: number,
    imageId: number
  ): Promise<void> {
    await apiClient.delete(
      `/api/v1/storefronts/${storefrontSlug}/products/${productId}/images/${imageId}`
    );
  },

  // Установить главное изображение
  async setMainImage(
    storefrontSlug: string,
    productId: number,
    imageId: number
  ): Promise<void> {
    await apiClient.put(
      `/api/v1/storefronts/${storefrontSlug}/products/${productId}/images/${imageId}/main`,
      {}
    );
  },
};
