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
    this.baseURL = baseURL || configManager.getConfig().api.url;
  }

  /**
   * Выполняет GET запрос
   */
  protected async get<T>(
    endpoint: string,
    params?: QueryParams,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}${buildQueryString(params || {})}`;

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
    return this.request<T>(`${this.baseURL}${endpoint}`, {
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
    return this.request<T>(`${this.baseURL}${endpoint}`, {
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
    return this.request<T>(`${this.baseURL}${endpoint}`, {
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
    return this.request<T>(`${this.baseURL}${endpoint}`, {
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
   * Базовый метод для выполнения запросов
   */
  protected async request<T>(
    url: string,
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    try {
      // Добавляем JWT токен в заголовки если есть
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options?.headers as Record<string, string>),
      };

      const token = await this.getAuthToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
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
  // Добавляем axios-подобные методы для совместимости с tokenManager
  interceptors = {
    request: {
      use: () => {
        // Пустая реализация для совместимости
        console.log('[ApiClient] Request interceptor registered (no-op)');
      },
    },
    response: {
      use: () => {
        // Пустая реализация для совместимости
        console.log('[ApiClient] Response interceptor registered (no-op)');
      },
    },
  };

  // Дополнительный метод для axios-подобного интерфейса
  async axiosRequest(config: {
    url?: string;
    method?: string;
    headers?: HeadersInit;
    data?: unknown;
  }) {
    const url = config.url || '';
    const response = await fetch(url, {
      method: config.method || 'GET',
      headers: config.headers,
      body: config.data ? JSON.stringify(config.data) : undefined,
      credentials: 'include',
    });

    return {
      data: await response.json().catch(() => ({})),
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
      config,
    };
  }

  // Добавляем свойство request как объект для совместимости
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  request = this.axiosRequest.bind(this) as any;
}

export const apiClient = new DefaultApiClient();
