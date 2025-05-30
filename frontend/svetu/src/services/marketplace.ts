import {
  MarketplaceSearchParams,
  MarketplaceSearchResponse,
} from '@/types/marketplace';
import configManager from '@/config';

export class MarketplaceService {
  static async search(
    params: MarketplaceSearchParams = {}
  ): Promise<MarketplaceSearchResponse> {
    const searchParams = new URLSearchParams();

    if (params.sort_by) searchParams.append('sort_by', params.sort_by);
    if (params.page !== undefined)
      searchParams.append('page', params.page.toString());
    if (params.size !== undefined)
      searchParams.append('size', params.size.toString());
    if (params.category) searchParams.append('category', params.category);
    if (params.minPrice !== undefined)
      searchParams.append('minPrice', params.minPrice.toString());
    if (params.maxPrice !== undefined)
      searchParams.append('maxPrice', params.maxPrice.toString());
    if (params.search) searchParams.append('search', params.search);

    const url = `${configManager.getApiUrl()}/api/v1/marketplace/search?${searchParams.toString()}`;

    try {
      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(
          `Failed to fetch marketplace items: ${response.status}`
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching marketplace items:', error);
      throw error;
    }
  }
}
