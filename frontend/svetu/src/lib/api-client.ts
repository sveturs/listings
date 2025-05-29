import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import { ApiError } from '@/types/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || process.env.NEXT_PUBLIC_BACKEND_URL || 'https://svetu.rs';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      withCredentials: true,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      validateStatus: function (status: number) {
        return status >= 200 && status < 500;
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // Get token from localStorage or Redux store
        const token = this.getAuthToken();
        
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        // Add language to requests
        const currentLanguage = this.getCurrentLanguage();
        
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

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError<ApiError>) => {
        if (error.response) {
          // Handle 401 Unauthorized
          if (error.response.status === 401) {
            // Clear auth data and redirect to login
            this.clearAuthData();
            window.location.href = '/login';
          }

          // Create standardized error
          const apiError: ApiError = {
            message: error.response.data?.message || 'An error occurred',
            statusCode: error.response.status,
            errors: error.response.data?.errors,
          };

          return Promise.reject(apiError);
        }

        // Network error or timeout
        const networkError: ApiError = {
          message: 'Network error. Please check your connection.',
          statusCode: 0,
        };

        return Promise.reject(networkError);
      }
    );
  }

  private getAuthToken(): string | null {
    // Try to get from localStorage first
    if (typeof window !== 'undefined') {
      return localStorage.getItem('authToken');
    }
    return null;
  }

  private clearAuthData(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('authToken');
      localStorage.removeItem('user');
    }
  }

  // Store token when login successful
  public setAuthToken(token: string | null): void {
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem('authToken', token);
      } else {
        localStorage.removeItem('authToken');
      }
    }
  }

  private getCurrentLanguage(): string {
    if (typeof window !== 'undefined') {
      // Try to get from localStorage first
      const storedLang = localStorage.getItem('locale');
      if (storedLang) return storedLang;
      
      // Try to get from URL path
      const pathLang = window.location.pathname.split('/')[1];
      if (['en', 'rs', 'ru'].includes(pathLang)) {
        return pathLang;
      }
    }
    return 'rs'; // Default language
  }

  // Get axios instance for direct use
  public getInstance(): AxiosInstance {
    return this.client;
  }

  // Convenience methods
  public get<T>(url: string, config?: Record<string, unknown>) {
    return this.client.get<T>(url, config);
  }

  public post<T>(url: string, data?: unknown, config?: Record<string, unknown>) {
    return this.client.post<T>(url, data, config);
  }

  public put<T>(url: string, data?: unknown, config?: Record<string, unknown>) {
    return this.client.put<T>(url, data, config);
  }

  public patch<T>(url: string, data?: unknown, config?: Record<string, unknown>) {
    return this.client.patch<T>(url, data, config);
  }

  public delete<T>(url: string, config?: Record<string, unknown>) {
    return this.client.delete<T>(url, config);
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

// Export for type usage
export default ApiClient;