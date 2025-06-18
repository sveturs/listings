import { config } from '@/config';

export interface ListingCardData {
  id: number;
  title: string;
  price: number;
  old_price?: number;
  location?: string;
  main_image?: string;
  has_discount?: boolean;
  discount_percentage?: number;
  created_at: string;
  views_count?: number;
  favorites_count?: number;
  is_favorite?: boolean;
}

export interface MarketplaceListingsResponse {
  data: ListingCardData[];
  pagination: {
    total: number;
    limit: number;
    offset: number;
  };
}

export interface MarketplaceFilters {
  storefront_id?: number;
  category_id?: number;
  search?: string;
  min_price?: number;
  max_price?: number;
  condition?: string;
  location?: string;
  limit?: number;
  offset?: number;
}

class MarketplaceApiService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = config.api.url;
  }

  private createUrl(endpoint: string, params?: Record<string, any>): string {
    const url = new URL(`${this.baseUrl}${endpoint}`);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value));
        }
      });
    }

    return url.toString();
  }

  private async request<T>(url: string, options?: RequestInit): Promise<T> {
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          error: `HTTP error! status: ${response.status}`,
        }));
        throw new Error(
          errorData.error || errorData.message || 'Network error'
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Unknown error occurred');
    }
  }

  /**
   * Get marketplace listings with filters
   */
  async getListings(
    filters?: MarketplaceFilters
  ): Promise<MarketplaceListingsResponse> {
    const url = this.createUrl('/api/v1/marketplace/listings', filters);
    return this.request<MarketplaceListingsResponse>(url);
  }

  /**
   * Get listing by ID
   */
  async getListingById(id: number): Promise<ListingCardData> {
    const url = this.createUrl(`/api/v1/marketplace/listings/${id}`);
    return this.request<ListingCardData>(url);
  }
}

export const marketplaceApi = new MarketplaceApiService();
export default marketplaceApi;
