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
  location_privacy?: string;
  address_multilingual?: Record<string, string>;
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
    data: CreateListingState,
    locationPrivacy?: string
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
      // Используем сербскую версию как основную (мы в Сербии)
      location:
        data.location?.addressMultilingual?.sr || data.location?.address || '',
      latitude: data.location?.latitude,
      longitude: data.location?.longitude,
      // Извлекаем город из сербского адреса или используем то что есть
      city:
        data.location?.addressMultilingual?.sr?.split(',')[0]?.trim() ||
        data.location?.city ||
        '',
      country: 'Србија',
      show_on_map:
        data.location?.latitude && data.location?.longitude ? true : false,
      location_privacy: locationPrivacy,
      // Добавляем мультиязычные адреса
      address_multilingual: data.location?.addressMultilingual,

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

    // Добавляем мультиязычные адреса в переводы
    const addressTranslations: Record<string, Record<string, string>> = {};

    if (data.location?.addressMultilingual) {
      // Добавляем переводы адресов для каждого языка
      if (data.location.addressMultilingual.en) {
        addressTranslations.en = {
          location: data.location.addressMultilingual.en,
          city:
            data.location.addressMultilingual.en.split(',')[0]?.trim() ||
            'Belgrade',
          country: 'Serbia',
        };
      }

      if (data.location.addressMultilingual.ru) {
        addressTranslations.ru = {
          location: data.location.addressMultilingual.ru,
          city:
            data.location.addressMultilingual.ru.split(',')[0]?.trim() ||
            'Белград',
          country: 'Сербия',
        };
      }

      if (data.location.addressMultilingual.sr) {
        addressTranslations.sr = {
          location: data.location.addressMultilingual.sr,
          city:
            data.location.addressMultilingual.sr.split(',')[0]?.trim() ||
            'Београд',
          country: 'Србија',
        };
      }
    }

    // Добавляем переводы, если они есть
    let finalTranslations: Record<string, Record<string, string>> = {};

    // Сначала обрабатываем существующие переводы
    if (data.translations) {
      const translationsObj = data.translations as any;
      if (
        'data' in translationsObj &&
        typeof translationsObj.data === 'object'
      ) {
        // Это ответ от API с полями data, message, success
        finalTranslations = translationsObj.data as Record<
          string,
          Record<string, string>
        >;
      } else if (
        !('success' in translationsObj) &&
        !('message' in translationsObj)
      ) {
        // Если это уже чистый объект переводов (без полей success/message)
        finalTranslations = translationsObj as Record<
          string,
          Record<string, string>
        >;
      }
    }

    // Добавляем или перезаписываем переводы адресов
    if (Object.keys(addressTranslations).length > 0) {
      for (const lang of Object.keys(addressTranslations)) {
        finalTranslations[lang] = {
          ...finalTranslations[lang],
          ...addressTranslations[lang],
        };
      }
    }

    // Добавляем переводы в запрос, если есть
    if (Object.keys(finalTranslations).length > 0) {
      request.translations = finalTranslations;
      console.log(
        'Translations being sent (including addresses):',
        request.translations
      );
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

    files.forEach((file) => {
      formData.append('file', file);
    });

    if (mainImageIndex >= 0 && mainImageIndex < files.length) {
      formData.append('main_image_index', mainImageIndex.toString());
    }

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
