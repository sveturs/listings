import { BaseApiService } from './baseApi';

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

export interface C2CListingsResponse {
  data: ListingCardData[];
  pagination: {
    total: number;
    limit: number;
    offset: number;
  };
}

export interface C2CFilters {
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

class MarketplaceApiService extends BaseApiService {
  /**
   * Get marketplace listings with filters
   */
  async getListings(filters?: C2CFilters): Promise<C2CListingsResponse> {
    const url = this.createUrl('/api/v1/c2c/listings', filters);
    return this.request<C2CListingsResponse>(url);
  }

  /**
   * Get listing by ID
   */
  async getListingById(id: number): Promise<ListingCardData> {
    const url = this.createUrl(`/api/v1/c2c/listings/${id}`);
    return this.request<ListingCardData>(url);
  }
}

export const marketplaceApi = new MarketplaceApiService();
export default marketplaceApi;
