import axios from 'axios';
import configManager from '@/config';

const api = axios.create({
  baseURL: configManager.getApiUrl(),
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Можно добавить токен авторизации если нужно
    const token =
      typeof window !== 'undefined' ? localStorage.getItem('token') : null;
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
  (error) => {
    if (error.response?.status === 401) {
      // Обработка ошибки авторизации
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        // Можно добавить редирект на страницу логина
      }
    }
    return Promise.reject(error);
  }
);

export default api;
