// frontend/hostel-frontend/src/api/axios.js
import axios from 'axios';
import i18n from '../i18n/config';

const instance = axios.create({
  baseURL: process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000',
  withCredentials: true,
  headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
  },
  validateStatus: function (status) {
      return status >= 200 && status < 500;
  }
});

// Add a request interceptor to include language in every request
instance.interceptors.request.use(
  config => {
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
    config.headers['Accept-Language'] = currentLanguage;
    
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

instance.interceptors.response.use(
  response => response,
  error => {
      if (error.response?.status === 401) {
          // Можно добавить редирект на страницу логина
          window.location.href = '/login';
      }
      return Promise.reject(error);
  }
);

export default instance;