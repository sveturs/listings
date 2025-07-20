import { config } from '@/config';

export abstract class BaseApiService {
  protected baseUrl: string;

  constructor() {
    this.baseUrl = config.api.url;
  }

  protected createUrl(endpoint: string, params?: Record<string, any>): string {
    const url = new URL(`${this.baseUrl}${endpoint}`);

    // Автоматически добавляем текущий язык из URL или localStorage
    const currentLocale = this.getCurrentLocale();
    url.searchParams.append('lang', currentLocale);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value));
        }
      });
    }

    return url.toString();
  }

  protected getCurrentLocale(): string {
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

  protected async request<T>(url: string, options?: RequestInit): Promise<T> {
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          // Добавляем язык в заголовок Accept-Language
          'Accept-Language': this.getCurrentLocale(),
          ...options?.headers,
        },
        ...options,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          error: `HTTP error! status: ${response.status}`,
        }));
        throw new Error(
          errorData.error || errorData.message || 'Network error'
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Unknown error occurred');
    }
  }
}
