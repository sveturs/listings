/**
 * API клиент для взаимодействия с backend через BFF proxy
 *
 * Все запросы идут через /api/v2/* → backend /api/v1/*
 * BFF proxy автоматически добавляет JWT токен из httpOnly cookies
 */

// Используем BFF proxy (/api/v2/*) вместо прямого обращения к backend
const USE_BFF_PROXY = process.env.NEXT_PUBLIC_USE_API_V2 !== 'false'; // по умолчанию true
const API_BASE_URL = USE_BFF_PROXY ? '' : (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000');

interface ApiClientOptions {
  headers?: Record<string, string>;
}

// CSRF токен кэш
let csrfToken: string | null = null;

/**
 * Получает CSRF токен для мутирующих запросов
 */
async function getCSRFToken(): Promise<string> {
  if (csrfToken) {
    return csrfToken;
  }

  try {
    const url = USE_BFF_PROXY ? '/api/v2/csrf-token' : `${API_BASE_URL}/api/v1/csrf-token`;
    const response = await fetch(url, {
      method: 'GET',
      credentials: USE_BFF_PROXY ? 'include' : 'omit',
    });

    if (!response.ok) {
      throw new Error(`Failed to get CSRF token: ${response.status}`);
    }

    const data = await response.json();
    // Backend может вернуть либо {csrf_token: "..."} либо {data: {csrf_token: "..."}}
    csrfToken = data.csrf_token || data.data?.csrf_token;

    return csrfToken!;
  } catch (error) {
    console.error('Error getting CSRF token:', error);
    throw error;
  }
}

/**
 * Преобразует путь для BFF proxy
 * /api/v1/... → /api/v2/... (убирает v1, BFF proxy сам добавит v1 при проксировании)
 */
function transformPath(path: string): string {
  if (!USE_BFF_PROXY) {
    return path;
  }

  // Если путь начинается с /api/v1/, заменяем на /api/v2/
  if (path.startsWith('/api/v1/')) {
    return path.replace('/api/v1/', '/api/v2/');
  }

  // Если путь не содержит /api/, добавляем префикс /api/v2
  if (!path.startsWith('/api/')) {
    return `/api/v2${path}`;
  }

  return path;
}

export const apiClient = {
  async get(path: string, options?: ApiClientOptions) {
    const url = USE_BFF_PROXY ? transformPath(path) : `${API_BASE_URL}${path}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      credentials: USE_BFF_PROXY ? 'include' : 'omit', // для BFF proxy нужны cookies
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  },

  async post(path: string, data: any, options?: ApiClientOptions) {
    const url = USE_BFF_PROXY ? transformPath(path) : `${API_BASE_URL}${path}`;
    const token = await getCSRFToken();
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': token,
        ...options?.headers,
      },
      body: JSON.stringify(data),
      credentials: USE_BFF_PROXY ? 'include' : 'omit',
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  },

  async put(path: string, data: any, options?: ApiClientOptions) {
    const url = USE_BFF_PROXY ? transformPath(path) : `${API_BASE_URL}${path}`;
    const token = await getCSRFToken();
    const response = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': token,
        ...options?.headers,
      },
      body: JSON.stringify(data),
      credentials: USE_BFF_PROXY ? 'include' : 'omit',
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  },

  async delete(path: string, options?: ApiClientOptions) {
    const url = USE_BFF_PROXY ? transformPath(path) : `${API_BASE_URL}${path}`;
    const token = await getCSRFToken();
    console.log('[apiClient DELETE] CSRF token obtained:', token ? 'present' : 'MISSING');
    console.log('[apiClient DELETE] URL:', url);
    const response = await fetch(url, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': token,
        ...options?.headers,
      },
      credentials: USE_BFF_PROXY ? 'include' : 'omit',
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  },
};
