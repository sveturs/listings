import configManager from '@/config';

export interface UserContact {
  id: number;
  user_id: number;
  contact_user_id: number;
  status: 'pending' | 'accepted' | 'blocked';
  notes?: string;
  added_from_chat_id?: number;
  created_at: string;
  updated_at: string;
  contact_user?: {
    id: number;
    name: string;
    email: string;
    picture_url?: string;
  };
  user?: {
    id: number;
    name: string;
    email: string;
    picture_url?: string;
  };
}

export interface ContactsResponse {
  contacts: UserContact[];
  total_count: number;
  page: number;
  limit: number;
}

export interface PrivacySettings {
  user_id: number;
  allow_contact_requests: boolean;
  allow_messages_from_contacts_only: boolean;
  created_at: string;
  updated_at: string;
}

export interface AddContactRequest {
  contact_user_id: number;
  notes?: string;
  added_from_chat_id?: number;
}

export interface UpdateContactRequest {
  status: 'accepted' | 'blocked';
  notes?: string;
}

export interface UpdatePrivacySettingsRequest {
  allow_contact_requests?: boolean;
  allow_messages_from_contacts_only?: boolean;
}

class ContactsService {
  private baseUrl: string;
  private cache = new Map<string, { data: any; timestamp: number }>();
  private cacheTimeout = 5000; // 5 seconds cache
  private pendingRequests = new Map<string, Promise<any>>();

  constructor() {
    this.baseUrl = `${configManager.getApiUrl()}/api/v1/contacts`;
  }

  private getCacheKey(endpoint: string, options?: RequestInit): string {
    return `${options?.method || 'GET'}:${endpoint}`;
  }

  private getFromCache(key: string): any | null {
    const cached = this.cache.get(key);
    if (!cached) return null;

    if (Date.now() - cached.timestamp > this.cacheTimeout) {
      this.cache.delete(key);
      return null;
    }

    return cached.data;
  }

  private setCache(key: string, data: any): void {
    this.cache.set(key, { data, timestamp: Date.now() });
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const method = options?.method?.toUpperCase() || 'GET';

    // Проверяем кэш только для GET запросов
    if (method === 'GET') {
      const cacheKey = this.getCacheKey(endpoint, options);
      const cached = this.getFromCache(cacheKey);
      if (cached) {
        return cached as T;
      }

      // Проверяем, есть ли уже запрос в процессе
      const pendingRequest = this.pendingRequests.get(cacheKey);
      if (pendingRequest) {
        return pendingRequest as Promise<T>;
      }
    }

    // Создаём promise для запроса
    const requestPromise = (async () => {
      // Добавляем JWT токен в заголовки через tokenManager
      const { tokenManager } = await import('@/utils/tokenManager');
      const token = await tokenManager.getAccessToken();

      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options?.headers as Record<string, string>),
      };

      if (token) {
        (headers as any)['Authorization'] = `Bearer ${token}`;
      } else {
        console.warn('[ContactsService] No access token available');
      }

      // Добавляем CSRF токен для небезопасных методов
      if (['POST', 'PUT', 'DELETE', 'PATCH'].includes(method)) {
        try {
          const { AuthService } = await import('./auth');
          const csrfToken = await AuthService.getCsrfToken();
          if (csrfToken) {
            (headers as any)['X-CSRF-Token'] = csrfToken;
          }
        } catch (error) {
          console.warn('[ContactsService] Failed to get CSRF token:', error);
        }
      }

      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        credentials: 'include',
        headers,
        ...options,
      });

      if (!response.ok) {
        let errorMessage = `HTTP ${response.status}`;
        try {
          const errorData = await response.json();
          errorMessage = errorData.error || errorData.message || errorMessage;
        } catch {
          const errorText = await response.text();
          errorMessage = errorText || errorMessage;
        }
        if (response.status === 401) {
          errorMessage = 'Unauthorized';
        }
        throw new Error(errorMessage);
      }

      const data = await response.json();
      return data.data || data;
    })();

    // Сохраняем promise для GET запросов
    if (method === 'GET') {
      const cacheKey = this.getCacheKey(endpoint, options);
      this.pendingRequests.set(cacheKey, requestPromise);
    }

    try {
      const result = await requestPromise;

      // Кэшируем результат для GET запросов
      if (method === 'GET') {
        const cacheKey = this.getCacheKey(endpoint, options);
        this.setCache(cacheKey, result);
        this.pendingRequests.delete(cacheKey);
      }

      // Инвалидируем кэш при изменениях
      if (['POST', 'PUT', 'DELETE', 'PATCH'].includes(method)) {
        this.cache.clear();
      }

      return result;
    } catch (error) {
      // Удаляем из pending requests при ошибке
      if (method === 'GET') {
        const cacheKey = this.getCacheKey(endpoint, options);
        this.pendingRequests.delete(cacheKey);
      }
      throw error;
    }
  }

  async getContacts(
    status?: string,
    page: number = 1,
    limit: number = 20
  ): Promise<ContactsResponse> {
    const params = new URLSearchParams();
    if (status) params.append('status', status);
    params.append('page', page.toString());
    params.append('limit', limit.toString());

    return this.request<ContactsResponse>(`?${params.toString()}`);
  }

  async addContact(request: AddContactRequest): Promise<UserContact> {
    return this.request<UserContact>('', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async updateContactStatus(
    contactUserID: number,
    request: UpdateContactRequest
  ): Promise<void> {
    return this.request<void>(`/${contactUserID}/status`, {
      method: 'PUT',
      body: JSON.stringify(request),
    });
  }

  async removeContact(contactUserID: number): Promise<void> {
    return this.request<void>(`/${contactUserID}`, {
      method: 'DELETE',
    });
  }

  async getPrivacySettings(): Promise<PrivacySettings> {
    return this.request<PrivacySettings>('/privacy');
  }

  async updatePrivacySettings(
    request: UpdatePrivacySettingsRequest
  ): Promise<PrivacySettings> {
    return this.request<PrivacySettings>('/privacy', {
      method: 'PUT',
      body: JSON.stringify(request),
    });
  }

  async getContactStatus(
    contactUserID: number
  ): Promise<{ are_contacts: boolean; user_id: number; contact_id: number }> {
    return this.request<{
      are_contacts: boolean;
      user_id: number;
      contact_id: number;
    }>(`/status/${contactUserID}`);
  }

  async getIncomingRequests(
    params: {
      page?: number;
      limit?: number;
    } = {}
  ): Promise<ContactsResponse> {
    const queryParams = new URLSearchParams();
    if (params.page) {
      queryParams.append('page', params.page.toString());
    }
    if (params.limit) {
      queryParams.append('limit', params.limit.toString());
    }

    return this.request<ContactsResponse>(
      `/incoming?${queryParams.toString()}`
    );
  }
}

export const contactsService = new ContactsService();
