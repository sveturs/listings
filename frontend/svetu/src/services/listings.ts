import { apiClient } from './api-client';
import type { CreateListingState } from '@/contexts/CreateListingContext';

export interface CreateListingRequest {
  category_id: number;
  title: string;
  description: string;
  price: number;
  condition: string;
  location: string;
  latitude?: number;
  longitude?: number;
  city: string;
  country: string;
  show_on_map: boolean;
  attributes?: Array<{
    attribute_id: number;
    text_value?: string;
    numeric_value?: number;
    boolean_value?: boolean;
    json_value?: any;
  }>;
  // Региональные поля
  payment_methods?: string[];
  delivery_options?: string[];
  meeting_types?: string[];
  phone_verified?: boolean;
  negotiable_price?: boolean;
  bundle_deals?: boolean;
}

export interface CreateListingResponse {
  id: number;
  message: string;
}

export interface UploadImagesRequest {
  listing_id: number;
  images: File[];
  main_image_index?: number;
}

export interface UploadImagesResponse {
  images: Array<{
    id: number;
    file_path: string;
    public_url: string;
    is_main: boolean;
  }>;
}

export class ListingsService {
  static async createListing(
    data: CreateListingState
  ): Promise<CreateListingResponse> {
    // Преобразуем данные из формы в формат API
    const request: CreateListingRequest = {
      category_id: data.category?.id || 0,
      title: data.title,
      description: data.description,
      price: data.price,
      condition: data.condition,
      location: data.location?.address || '',
      latitude: data.location?.latitude,
      longitude: data.location?.longitude,
      city: data.location?.city || '',
      country: data.location?.country || 'Србија',
      show_on_map:
        data.location?.latitude && data.location?.longitude ? true : false,

      // Региональные поля сохраняем в атрибутах
      payment_methods: data.payment.methods,
      delivery_options: data.payment.deliveryOptions,
      meeting_types: data.trust.preferredMeetingType
        ? [data.trust.preferredMeetingType]
        : [],
      phone_verified: data.trust.phoneVerified,
      negotiable_price: data.payment.negotiablePrice,
      bundle_deals: data.payment.bundleDeals,
    };

    // Преобразуем атрибуты из формы
    if (data.attributes && Object.keys(data.attributes).length > 0) {
      request.attributes = Object.entries(data.attributes).map(
        ([_key, value]) => {
          // value здесь уже объект AttributeFormData
          const attributeData = value as any;
          return {
            attribute_id: attributeData.attribute_id,
            attribute_name: attributeData.attribute_name,
            display_name: attributeData.display_name,
            attribute_type: attributeData.attribute_type,
            text_value: attributeData.text_value,
            numeric_value: attributeData.numeric_value,
            boolean_value: attributeData.boolean_value,
            json_value: attributeData.json_value,
            unit: attributeData.unit,
          };
        }
      );
    }

    const response = await apiClient.post<CreateListingResponse>(
      '/api/v1/marketplace/listings',
      request
    );

    return response.data;
  }

  static async uploadImages(
    listingId: number,
    files: File[],
    mainImageIndex: number = 0
  ): Promise<UploadImagesResponse> {
    const formData = new FormData();

    files.forEach((file, index) => {
      formData.append('images', file);
      if (index === mainImageIndex) {
        formData.append('main_image_index', index.toString());
      }
    });

    const response = await apiClient.post<UploadImagesResponse>(
      `/api/v1/marketplace/listings/${listingId}/images`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );

    return response.data;
  }

  static async saveDraft(data: CreateListingState): Promise<void> {
    // TODO: Реализовать сохранение черновика
    // Можно использовать localStorage для временного хранения
    const draftKey = `listing_draft_${data.category?.id || 'new'}`;
    localStorage.setItem(draftKey, JSON.stringify(data));
  }

  static async loadDraft(
    categoryId?: number
  ): Promise<CreateListingState | null> {
    const draftKey = `listing_draft_${categoryId || 'new'}`;
    const draft = localStorage.getItem(draftKey);

    if (draft) {
      try {
        return JSON.parse(draft);
      } catch (error) {
        console.error('Failed to parse draft:', error);
        return null;
      }
    }

    return null;
  }

  static async deleteDraft(categoryId?: number): Promise<void> {
    const draftKey = `listing_draft_${categoryId || 'new'}`;
    localStorage.removeItem(draftKey);
  }
}
