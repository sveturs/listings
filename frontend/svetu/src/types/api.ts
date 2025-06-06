// Стандартный формат ответа API
export interface ApiResponse<T = unknown> {
  success: boolean;
  data: T;
  error?: string;
  message?: string;
  metadata?: {
    page?: number;
    limit?: number;
    total?: number;
    hasMore?: boolean;
  };
}

// Пагинированный ответ
export interface PaginatedResponse<T> {
  items: T[];
  page: number;
  limit: number;
  total: number;
  hasMore: boolean;
}

// Ошибка API
export interface ApiError {
  error: string;
  message?: string;
  code?: string;
  details?: Record<string, unknown>;
}

// Параметры пагинации
export interface PaginationParams {
  page?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
}

// Параметры с отменой запроса
export interface RequestOptions {
  signal?: AbortSignal;
  headers?: HeadersInit;
}

// Тип для формирования URL параметров
export type QueryParams = Record<
  string,
  string | number | boolean | undefined | null
>;

// Хелпер для построения query string
export function buildQueryString(params: QueryParams): string {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      searchParams.append(key, String(value));
    }
  });

  const queryString = searchParams.toString();
  return queryString ? `?${queryString}` : '';
}

// Типизированный fetch wrapper
export async function apiFetch<T>(
  url: string,
  options?: RequestInit
): Promise<ApiResponse<T>> {
  try {
    const response = await fetch(url, {
      credentials: 'include',
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || `HTTP error! status: ${response.status}`);
    }

    // Нормализуем ответ к стандартному формату
    if ('success' in data) {
      return data as ApiResponse<T>;
    }

    // Если сервер возвращает данные напрямую, оборачиваем их
    return {
      success: true,
      data: data as T,
    };
  } catch (error) {
    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw error; // Пробрасываем AbortError дальше
      }

      return {
        success: false,
        data: null as T,
        error: error.message,
      };
    }

    return {
      success: false,
      data: null as T,
      error: 'Unknown error occurred',
    };
  }
}
