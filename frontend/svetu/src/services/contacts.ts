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

  constructor() {
    this.baseUrl = `${configManager.getApiUrl()}/api/v1/contacts`;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    // Добавляем JWT токен в заголовки через tokenManager
    const { tokenManager } = await import('@/utils/tokenManager');
    const token = await tokenManager.getAccessToken();

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options?.headers as Record<string, string>),
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    } else {
      console.warn('[ContactsService] No access token available');
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
    return this.request<void>(`/${contactUserID}`, {
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
}

export const contactsService = new ContactsService();
