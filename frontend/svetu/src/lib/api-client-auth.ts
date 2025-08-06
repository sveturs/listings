import { apiClient } from './api-client';
import { tokenManager } from '@/utils/tokenManager';

// Функция для получения токена
function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null;

  // Сначала пробуем получить токен из tokenManager
  const token = tokenManager.getAccessToken();
  if (token) {
    return token;
  }

  // Если tokenManager не инициализирован, пробуем sessionStorage напрямую
  return sessionStorage.getItem('svetu_access_token');
}

// Расширенный API клиент с автоматической авторизацией
export const apiClientAuth = {
  async get(path: string, options?: any) {
    const token = getAuthToken();
    const headers = {
      ...options?.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return apiClient.get(path, { ...options, headers });
  },

  async post(path: string, data: any, options?: any) {
    const token = getAuthToken();
    const headers = {
      ...options?.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return apiClient.post(path, data, { ...options, headers });
  },

  async put(path: string, data: any, options?: any) {
    const token = getAuthToken();
    const headers = {
      ...options?.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return apiClient.put(path, data, { ...options, headers });
  },

  async delete(path: string, options?: any) {
    const token = getAuthToken();
    const headers = {
      ...options?.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return apiClient.delete(path, { ...options, headers });
  },
};
