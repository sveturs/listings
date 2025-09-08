import axios from 'axios';
import { tokenManager } from '@/utils/tokenManager';
import configManager from '@/config';

// Получаем базовый URL из конфигурации
const apiUrl = configManager.get('api.url');

const api = axios.create({
  baseURL: apiUrl, // Используем URL из NEXT_PUBLIC_API_URL (https://devapi.svetu.rs)
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Получаем токен из tokenManager
    const token = tokenManager.getAccessToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    if (error.response?.status === 401) {
      // Пытаемся обновить токен
      try {
        const newToken = await tokenManager.refreshAccessToken();
        if (newToken && error.config && !error.config._retry) {
          // Помечаем запрос как повторный чтобы избежать бесконечного цикла
          error.config._retry = true;
          // Обновляем заголовок авторизации с новым токеном
          error.config.headers.Authorization = `Bearer ${newToken}`;
          // Повторяем оригинальный запрос с новым токеном
          return api(error.config);
        }
      } catch {
        // Если обновление токена не удалось, очищаем токены
        tokenManager.clearTokens();
      }
    }
    return Promise.reject(error);
  }
);

export default api;
