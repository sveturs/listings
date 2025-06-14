import {
  ApiResponse,
  RequestOptions,
  QueryParams,
  buildQueryString,
} from '@/types/api';
import configManager from '@/config';

/**
 * Базовый класс для всех API сервисов
 * Обеспечивает единообразную обработку запросов и ошибок
 */
export abstract class ApiClient {
  protected baseURL: string;

  constructor(baseURL?: string) {
    // В development режиме и в браузере используем пустую строку для proxy
    if (
      !baseURL &&
      typeof window !== 'undefined' &&
      process.env.NODE_ENV === 'development'
    ) {
      this.baseURL = '';
    } else {
      this.baseURL = baseURL || configManager.getConfig().api.url;
    }
  }

  /**
   * Выполняет GET запрос
   */
  protected async get<T>(
    endpoint: string,
    params?: QueryParams,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    const url = this.baseURL
      ? `${this.baseURL}${endpoint}${buildQueryString(params || {})}`
      : `${endpoint}${buildQueryString(params || {})}`;

    return this.request<T>(url, {
      method: 'GET',
      ...options,
    });
  }

  /**
   * Выполняет POST запрос
   */
  protected async post<T, D = unknown>(
    endpoint: string,
    data?: D,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    const url = this.baseURL ? `${this.baseURL}${endpoint}` : endpoint;
    return this.request<T>(url, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    });
  }

  /**
   * Выполняет PUT запрос
   */
  protected async put<T, D = unknown>(
    endpoint: string,
    data?: D,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    const url = this.baseURL ? `${this.baseURL}${endpoint}` : endpoint;
    return this.request<T>(url, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    });
  }

  /**
   * Выполняет DELETE запрос
   */
  protected async delete<T>(
    endpoint: string,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    const url = this.baseURL ? `${this.baseURL}${endpoint}` : endpoint;
    return this.request<T>(url, {
      method: 'DELETE',
      ...options,
    });
  }

  /**
   * Загружает файлы
   */
  protected async upload<T>(
    endpoint: string,
    formData: FormData,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    const url = this.baseURL ? `${this.baseURL}${endpoint}` : endpoint;
    return this.request<T>(url, {
      method: 'POST',
      body: formData,
      headers: {
        // Не устанавливаем Content-Type, браузер сам установит с boundary
        ...options?.headers,
      },
    });
  }

  /**
   * Получает JWT токен через tokenManager
   */
  private async getAuthToken(): Promise<string | null> {
    if (typeof window === 'undefined') return null;
    try {
      const { tokenManager } = await import('@/utils/tokenManager');
      return await tokenManager.getAccessToken();
    } catch {
      return null;
    }
  }

  /**
   * Получает CSRF токен через AuthService
   */
  private async getCsrfToken(): Promise<string | null> {
    if (typeof window === 'undefined') return null;
    try {
      const { AuthService } = await import('./auth');
      return await AuthService.getCsrfToken();
    } catch {
      return null;
    }
  }

  /**
   * Базовый метод для выполнения запросов
   */
  protected async request<T>(
    url: string,
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    try {
      // Добавляем JWT токен в заголовки если есть
      const headers: Record<string, string> = {
        ...(options?.headers as Record<string, string>),
      };

      // Устанавливаем Content-Type только если это не FormData
      if (!(options?.body instanceof FormData)) {
        headers['Content-Type'] = 'application/json';
      }

      const token = await this.getAuthToken();
      if (token) {
        (headers as any)['Authorization'] = `Bearer ${token}`;
      }

      // Добавляем CSRF токен для небезопасных методов
      const method = options?.method?.toUpperCase();
      if (method && ['POST', 'PUT', 'DELETE', 'PATCH'].includes(method)) {
        try {
          const csrfToken = await this.getCsrfToken();
          if (csrfToken) {
            (headers as any)['X-CSRF-Token'] = csrfToken;
          }
        } catch (error) {
          // Если не удалось получить CSRF токен, продолжаем без него
          console.warn('Failed to get CSRF token:', error);
        }
      }

      const response = await fetch(url, {
        credentials: 'include',
        ...options,
        headers,
      });

      // Обработка пустых ответов (204 No Content)
      if (response.status === 204) {
        return {
          success: true,
          data: null as T,
        };
      }

      const contentType = response.headers.get('content-type');
      let data: unknown;

      if (contentType?.includes('application/json')) {
        data = await response.json();
      } else {
        // Если ответ не JSON, возвращаем как есть
        data = await response.text();
      }

      if (!response.ok) {
        // Обработка ошибок
        const errorData = data as { error?: string; message?: string };
        const error =
          errorData?.error ||
          errorData?.message ||
          `HTTP error! status: ${response.status}`;
        console.error(`API Error [${response.status}]:`, error);

        // Если получили 401, удаляем невалидный токен
        if (response.status === 401) {
          try {
            const { tokenManager } = await import('@/utils/tokenManager');
            tokenManager.clearTokens();
          } catch {
            // Игнорируем ошибки
          }
        }

        // Обработка rate limiting
        if (response.status === 429) {
          const retryAfter = response.headers.get('Retry-After');
          console.warn(
            `API rate limited (429), retry after: ${retryAfter || 'unknown'} seconds`
          );
        }

        return {
          success: false,
          data: null as T,
          error,
          message: errorData?.message,
        };
      }

      // Нормализация успешного ответа
      if (
        data &&
        typeof data === 'object' &&
        'success' in data &&
        'data' in data
      ) {
        return data as ApiResponse<T>;
      }

      // Если сервер возвращает данные напрямую
      return {
        success: true,
        data: data as T,
      };
    } catch (error) {
      // Обработка сетевых ошибок
      if (error instanceof Error) {
        if (error.name === 'AbortError') {
          console.log('Request aborted');
          throw error; // Пробрасываем для обработки в компонентах
        }

        console.error('Network error:', error);
        return {
          success: false,
          data: null as T,
          error: error.message,
        };
      }

      console.error('Unknown error:', error);
      return {
        success: false,
        data: null as T,
        error: 'Unknown error occurred',
      };
    }
  }

  /**
   * Проверяет успешность ответа и возвращает данные или выбрасывает ошибку
   */
  protected unwrap<T>(response: ApiResponse<T>): T {
    if (response.success && response.data !== null) {
      return response.data;
    }

    throw new Error(response.error || 'Request failed');
  }
}

// Создаем базовый экземпляр API клиента для использования в других модулях
class DefaultApiClient extends ApiClient {
  // Делаем методы публичными для использования в других сервисах
  public async get<T>(
    endpoint: string,
    params?: QueryParams,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return super.get<T>(endpoint, params, options);
  }

  public async post<T, D = unknown>(
    endpoint: string,
    data?: D,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return super.post<T, D>(endpoint, data, options);
  }

  public async put<T, D = unknown>(
    endpoint: string,
    data?: D,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return super.put<T, D>(endpoint, data, options);
  }

  public async delete<T>(
    endpoint: string,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return super.delete<T>(endpoint, options);
  }

  public async upload<T>(
    endpoint: string,
    formData: FormData,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return super.upload<T>(endpoint, formData, options);
  }

  // Добавляем axios-подобные методы для совместимости
  interceptors = {
    request: {
      use: (
        onFulfilled: (
          config: Record<string, unknown>
        ) => Record<string, unknown>,
        onRejected?: (error: unknown) => unknown
      ) => {
        // Простая реализация для совместимости с tokenManager
        this.requestInterceptor = onFulfilled;
        this.requestErrorInterceptor = onRejected;
      },
    },
    response: {
      use: (
        onFulfilled: (
          response: Record<string, unknown>
        ) => Record<string, unknown>,
        onRejected?: (error: unknown) => unknown
      ) => {
        // Простая реализация для совместимости с tokenManager
        this.responseInterceptor = onFulfilled;
        this.responseErrorInterceptor = onRejected;
      },
    },
  };

  private requestInterceptor?: (
    config: Record<string, unknown>
  ) => Record<string, unknown>;
  private requestErrorInterceptor?: (error: unknown) => unknown;
  private responseInterceptor?: (
    response: Record<string, unknown>
  ) => Record<string, unknown>;
  private responseErrorInterceptor?: (error: unknown) => unknown;
}

export const apiClient = new DefaultApiClient();
