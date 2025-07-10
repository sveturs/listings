import { apiClient, ApiResponse } from '../api-client';
import type { components } from '@/types/generated/api';

// Типы из сгенерированной схемы
type UserProfile =
  components['schemas']['backend_internal_domain_models.UserProfile'];
type MarketplaceListing =
  components['schemas']['backend_internal_domain_models.MarketplaceListing'];
type ChatMessage =
  components['schemas']['backend_internal_domain_models.MarketplaceMessage'];
type ListingsResponse =
  components['schemas']['internal_proj_marketplace_handler.ListingsResponse'];

/**
 * Базовый класс для API endpoints
 */
export class ApiEndpoint {
  constructor(protected basePath: string) {}

  /**
   * Определяет, должен ли запрос использовать внутренний URL
   * Переопределите в наследниках для специфичной логики
   */
  protected shouldUseInternalUrl(): boolean {
    // По умолчанию используем внутренний URL только для SSR
    return typeof window === 'undefined';
  }
}

/**
 * User API endpoints
 */
export class UserApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/users');
  }

  async getCurrentUser(): Promise<ApiResponse<UserProfile>> {
    return apiClient.get<UserProfile>(`${this.basePath}/me`, {
      internal: this.shouldUseInternalUrl(),
    });
  }

  async updateProfile(
    data: Partial<UserProfile>
  ): Promise<ApiResponse<UserProfile>> {
    return apiClient.patch<UserProfile>(`${this.basePath}/me`, data, {
      internal: this.shouldUseInternalUrl(),
    });
  }
}

/**
 * Marketplace API endpoints
 */
export class MarketplaceApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/marketplace');
  }

  // Для публичных данных можем использовать внутренний URL при SSR
  protected shouldUseInternalUrl(): boolean {
    return typeof window === 'undefined';
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
      `${this.basePath}/listings?${queryParams}`,
      {
        internal: this.shouldUseInternalUrl(),
      }
    );
  }

  async getListingById(id: string): Promise<ApiResponse<MarketplaceListing>> {
    return apiClient.get<MarketplaceListing>(
      `${this.basePath}/listings/${id}`,
      {
        internal: this.shouldUseInternalUrl(),
      }
    );
  }

  async createListing(
    data: FormData
  ): Promise<ApiResponse<MarketplaceListing>> {
    return apiClient.upload<MarketplaceListing>(
      `${this.basePath}/listings`,
      data
    );
  }

  async addToFavorites(listingId: string): Promise<ApiResponse<void>> {
    return apiClient.post(`${this.basePath}/listings/${listingId}/favorite`);
  }
}

/**
 * Chat API endpoints - всегда использует публичный URL для WebSocket
 */
export class ChatApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/chat');
  }

  // Для чата всегда используем публичный URL из-за WebSocket
  protected shouldUseInternalUrl(): boolean {
    return false;
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
