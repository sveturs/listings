import configManager from '@/config';

// Types for AI Storefront Service
export interface AnalyzeProductImageRequest {
  imageData: string; // base64
  language?: string;
}

export interface ProductAnalysisResult {
  title: string;
  titleVariants: string[];
  description: string;
  category: string;
  categoryHints: {
    domain: string;
    productType: string;
    keywords: string[];
  };
  price: number;
  priceRange: { min: number; max: number };
  attributes: Record<string, any>;
  suggestedVariants?: Array<{
    attributes: Record<string, string>;
    sku?: string;
    price?: number;
  }>;
  stockEstimate: number;
  condition: 'new' | 'used' | 'refurbished';
  keywords: string[];
  location?: {
    address: string;
    latitude: number;
    longitude: number;
    city: string;
    region: string;
    source: 'exif' | 'manual';
  };
}

export interface DetectCategoryRequest {
  title: string;
  description: string;
  aiHints?: {
    domain: string;
    productType: string;
    keywords: string[];
  };
  language?: string;
}

export interface CategoryDetectionResult {
  categoryId: number;
  categoryName: string;
  categoryPath: string;
  confidenceScore: number;
  algorithm: string;
  processingTimeMs: number;
  keywords: string[];
  aiHints?: {
    domain: string;
    productType: string;
    keywords: string[];
  };
  alternativeIds?: number[];
}

export interface ABTestTitlesRequest {
  titleVariants: string[];
}

export interface ABTestResult {
  bestVariant: string;
  bestVariantIndex: number;
  scores: Array<{
    title: string;
    score: number;
    reasoning?: string;
  }>;
}

export interface TranslateContentRequest {
  content: {
    title: string;
    description: string;
  };
  targetLanguages: string[];
  sourceLanguage?: string;
}

export interface TranslationResult {
  translations: Record<
    string,
    {
      title: string;
      description: string;
    }
  >;
}

export interface AIMetrics {
  entity_type: string;
  total_detections: number;
  accuracy: number;
  cache_hit_rate: number;
  avg_processing_ms: number;
}

/**
 * StorefrontAIService - AI-powered сервис для работы с товарами витрин
 * Переиспользует marketplace AI infrastructure для products
 */
class StorefrontAIService {
  private baseUrl: string;
  private authToken: string | null = null;

  constructor() {
    this.baseUrl = configManager.getApiUrl();
  }

  /**
   * Установить токен авторизации
   */
  setAuthToken(token: string) {
    this.authToken = token;
  }

  /**
   * Получить заголовки для запроса
   */
  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (this.authToken) {
      headers['Authorization'] = `Bearer ${this.authToken}`;
    }

    return headers;
  }

  /**
   * Анализ изображения товара через AI
   * @param imageBase64 - base64 encoded image
   * @param language - язык для анализа (ru/en/sr)
   */
  async analyzeProductImage(
    imageBase64: string,
    language: string = 'ru'
  ): Promise<ProductAnalysisResult> {
    const response = await fetch(
      `${this.baseUrl}/api/v1/storefronts/ai/analyze-product-image`,
      {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({
          imageData: imageBase64,
          language,
        }),
      }
    );

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ message: response.statusText }));
      throw new Error(
        errorData.message || `AI analysis failed: ${response.statusText}`
      );
    }

    const data = await response.json();
    return data.data || data;
  }

  /**
   * Определение категории для товара витрины
   * @param title - заголовок товара
   * @param description - описание товара
   * @param categoryHints - подсказки для определения категории
   * @param language - язык контента
   */
  async detectCategory(
    title: string,
    description: string,
    categoryHints?: {
      domain: string;
      productType: string;
      keywords: string[];
    },
    language: string = 'ru'
  ): Promise<CategoryDetectionResult> {
    const response = await fetch(
      `${this.baseUrl}/api/v1/storefronts/ai/detect-category`,
      {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({
          title,
          description,
          aiHints: categoryHints,
          language,
        }),
      }
    );

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ message: response.statusText }));
      throw new Error(
        errorData.message || `Category detection failed: ${response.statusText}`
      );
    }

    const data = await response.json();
    return data.data || data;
  }

  /**
   * A/B тестирование заголовков для товара
   * @param titleVariants - варианты заголовков для тестирования
   */
  async abTestTitles(titleVariants: string[]): Promise<ABTestResult> {
    if (titleVariants.length < 2) {
      throw new Error('At least 2 title variants are required for A/B testing');
    }

    const response = await fetch(
      `${this.baseUrl}/api/v1/storefronts/ai/ab-test-titles`,
      {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({
          titleVariants,
        }),
      }
    );

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ message: response.statusText }));
      throw new Error(
        errorData.message || `A/B testing failed: ${response.statusText}`
      );
    }

    const data = await response.json();
    return data.data || data;
  }

  /**
   * Перевод контента товара на другие языки
   * @param content - контент для перевода (title, description)
   * @param targetLanguages - целевые языки (en, ru, sr)
   * @param sourceLanguage - исходный язык (default: ru)
   */
  async translateProductContent(
    content: { title: string; description: string },
    targetLanguages: string[],
    sourceLanguage: string = 'ru'
  ): Promise<TranslationResult> {
    if (targetLanguages.length === 0) {
      throw new Error('At least one target language is required');
    }

    const response = await fetch(
      `${this.baseUrl}/api/v1/storefronts/ai/translate-content`,
      {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({
          content,
          targetLanguages,
          sourceLanguage,
        }),
      }
    );

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ message: response.statusText }));
      throw new Error(
        errorData.message || `Translation failed: ${response.statusText}`
      );
    }

    const data = await response.json();
    return data.data || data;
  }

  /**
   * Получить метрики AI для товаров витрин
   */
  async getMetrics(): Promise<AIMetrics> {
    const response = await fetch(
      `${this.baseUrl}/api/v1/storefronts/ai/metrics`,
      {
        method: 'GET',
        headers: this.getHeaders(),
      }
    );

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ message: response.statusText }));
      throw new Error(
        errorData.message || `Failed to get metrics: ${response.statusText}`
      );
    }

    const data = await response.json();
    return data.data || data;
  }

  /**
   * Извлечь GPS координаты из EXIF данных изображения
   * Это вспомогательная функция на клиенте
   * @param file - файл изображения
   */
  async extractGPSFromImage(_file: File): Promise<{
    latitude: number;
    longitude: number;
  } | null> {
    // TODO: Реализовать извлечение EXIF GPS данных на клиенте
    // Используем библиотеку exif-js или piexifjs
    return null;
  }

  /**
   * Конвертировать изображение в base64
   * @param file - файл изображения
   */
  async imageToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        const result = reader.result as string;
        // Удаляем префикс data:image/...;base64,
        const base64 = result.split(',')[1];
        resolve(base64);
      };
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  }

  /**
   * Compress изображение перед отправкой на сервер
   * @param file - файл изображения
   * @param maxWidth - максимальная ширина
   * @param maxHeight - максимальная высота
   * @param quality - качество сжатия (0-1)
   */
  async compressImage(
    file: File,
    maxWidth: number = 1920,
    maxHeight: number = 1080,
    quality: number = 0.85
  ): Promise<string> {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.onload = () => {
        const canvas = document.createElement('canvas');
        let { width, height } = img;

        // Вычисляем новые размеры сохраняя пропорции
        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height);
          width *= ratio;
          height *= ratio;
        }

        canvas.width = width;
        canvas.height = height;

        const ctx = canvas.getContext('2d');
        if (!ctx) {
          reject(new Error('Failed to get canvas context'));
          return;
        }

        ctx.drawImage(img, 0, 0, width, height);

        // Конвертируем в base64
        const base64 = canvas.toDataURL('image/jpeg', quality).split(',')[1];
        resolve(base64);
      };

      img.onerror = () => reject(new Error('Failed to load image'));

      const reader = new FileReader();
      reader.onload = (e) => {
        img.src = e.target?.result as string;
      };
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  }
}

// Экспортируем singleton instance
export const storefrontAI = new StorefrontAIService();
