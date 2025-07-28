import { apiClient } from './api-client';
import type { CreateListingState } from '@/contexts/CreateListingContext';
import type { components } from '@/types/generated/api';

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
  // Переводы полей
  translations?: Record<string, Record<string, string>>;

  // Язык оригинала
  original_language?: string;

  // ID статистики определения категории для обновления success_rate
  category_detection_stats_id?: number;
  
  // Ключевые слова, которые были использованы для определения категории
  detected_keywords?: string[];

  // Региональные поля
  payment_methods?: string[];
  delivery_options?: string[];
  meeting_types?: string[];
  phone_verified?: boolean;
  negotiable_price?: boolean;
  bundle_deals?: boolean;
}

// Используем сгенерированный тип из backend с обёрткой SuccessResponse
export type CreateListingResponse =
  components['schemas']['backend_pkg_utils.SuccessResponseSwag'] & {
    data?: components['schemas']['internal_proj_marketplace_handler.IDMessageResponse'];
  };

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
    // Отладочное логирование состояния формы
    console.log('CreateListing Data:', JSON.stringify(data, null, 2));

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

    // Добавляем переводы, если они есть
    if (data.translations && Object.keys(data.translations).length > 0) {
      request.translations = data.translations;
    }

    // Добавляем язык оригинала, если он указан
    if (data.originalLanguage) {
      request.original_language = data.originalLanguage;
    }

    // Добавляем ID статистики определения категории
    if (data.categoryDetectionStatsId) {
      request.category_detection_stats_id = data.categoryDetectionStatsId;
    }
    
    // Добавляем ключевые слова для определения категории
    if (data.detectedKeywords && data.detectedKeywords.length > 0) {
      request.detected_keywords = data.detectedKeywords;
    }

    // Отладочное логирование запроса к API
    console.log('API Request:', JSON.stringify(request, null, 2));

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

    if (response.error) {
      throw new Error(response.error.message);
    }

    if (!response.data) {
      throw new Error('No data received from server');
    }

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

    // Отладочное логирование
    console.log('Uploading images:', {
      listingId,
      filesCount: files.length,
      formDataEntries: Array.from(formData.entries()).map(([key, value]) => ({
        key,
        value: value instanceof File ? `File: ${value.name}` : value,
      })),
    });

    const response = await apiClient.upload<UploadImagesResponse>(
      `/api/v1/marketplace/listings/${listingId}/images`,
      formData
      // Не устанавливаем Content-Type для FormData - браузер сделает это автоматически с boundary
    );

    if (response.error) {
      throw new Error(response.error.message);
    }

    if (!response.data) {
      throw new Error('No data received from server');
    }

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
