import { apiClient, ApiResponse } from '../api-client';
import type { components } from '@/types/generated/api';

// Типы из сгенерированной схемы
type UserProfile = components['schemas']['models.UserProfile'];
type C2CListing = components['schemas']['models.MarketplaceListing'];
type ChatMessage = components['schemas']['models.MarketplaceMessage'];
type ListingsResponse = components['schemas']['handler.ListingsResponse'];

/**
 * Базовый класс для API endpoints
 * Все запросы идут через BFF proxy /api/v2
 */
export class ApiEndpoint {
  constructor(protected basePath: string) {}
}

/**
 * User API endpoints
 */
export class UserApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/users');
  }

  async getCurrentUser(): Promise<ApiResponse<UserProfile>> {
    return apiClient.get<UserProfile>(`${this.basePath}/me`);
  }

  async updateProfile(
    data: Partial<UserProfile>
  ): Promise<ApiResponse<UserProfile>> {
    return apiClient.patch<UserProfile>(`${this.basePath}/me`, data);
  }
}

/**
 * Marketplace API endpoints (unified c2c/b2c)
 */
export class MarketplaceApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/marketplace');
  }

  async getListings(params?: {
    page?: number;
    limit?: number;
    category?: string;
  }): Promise<ApiResponse<ListingsResponse>> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.set('page', params.page.toString());
    if (params?.limit) queryParams.set('limit', params.limit.toString());
    if (params?.category) queryParams.set('category', params.category);

    return apiClient.get<ListingsResponse>(
      `${this.basePath}/listings?${queryParams}`
    );
  }

  async getListingById(id: string): Promise<ApiResponse<C2CListing>> {
    return apiClient.get<C2CListing>(`${this.basePath}/listings/${id}`);
  }

  async createListing(data: FormData): Promise<ApiResponse<C2CListing>> {
    return apiClient.upload<C2CListing>(`${this.basePath}/listings`, data);
  }

  async addToFavorites(listingId: string): Promise<ApiResponse<void>> {
    return apiClient.post(`${this.basePath}/favorites/${listingId}`);
  }
}

/**
 * Chat API endpoints
 */
export class ChatApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/chat');
  }

  async getMessages(chatId: string): Promise<ApiResponse<ChatMessage[]>> {
    return apiClient.get<ChatMessage[]>(`${this.basePath}/${chatId}/messages`);
  }

  async sendMessage(
    chatId: string,
    message: string
  ): Promise<ApiResponse<ChatMessage>> {
    return apiClient.post<ChatMessage>(`${this.basePath}/${chatId}/messages`, {
      content: message,
    });
  }
}

// Экспортируем singleton instances
export const userApi = new UserApi();
export const marketplaceApi = new MarketplaceApi();
export const chatApi = new ChatApi();
