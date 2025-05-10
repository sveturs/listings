// frontend/hostel-frontend/src/api/axios.ts
import axios, { AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import i18n from '../i18n/config';

// Используем уже существующий тип Window.ENV из global.d.ts
const getBaseUrl = (): string => {
  if (window.ENV && window.ENV.REACT_APP_BACKEND_URL !== undefined) {
    return window.ENV.REACT_APP_BACKEND_URL;
  }
  return process.env.REACT_APP_BACKEND_URL || '';
};

const instance = axios.create({
  baseURL: getBaseUrl(),
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },
  validateStatus: function (status: number) {
    return status >= 200 && status < 500;
  }
});

// Add a request interceptor to include language in every request
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    // Get current language from i18n
    const currentLanguage = i18n.language;
    
    // Add language parameter to GET requests
    if (config.method === 'get') {
      config.params = config.params || {};
      // Only add if not already set
      if (!config.params.language) {
        config.params.language = currentLanguage;
      }
    }
    
    // Add Accept-Language header for all requests
    if (config.headers) {
      config.headers['Accept-Language'] = currentLanguage;
    }
    
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

instance.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Можно добавить редирект на страницу логина
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default instance;