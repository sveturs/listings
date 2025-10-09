import { apiClient } from './api-client';

export interface SearchFilters {
  categories?: number[];
  price_min?: number;
  price_max?: number;
  condition?: string;
  location?: string;
  radius?: number;
  [key: string]: any;
}

export interface SearchParams {
  query?: string;
  filters?: SearchFilters;
  sort?: string;
  page?: number;
  limit?: number;
  fuzzy?: boolean;
}

export interface SearchResult {
  items: any[];
  total: number;
  page: number;
  limit: number;
  facets?: {
    categories?: Array<{ name: string; count: number }>;
    price_ranges?: Array<{ range: string; count: number }>;
    conditions?: Array<{ value: string; count: number }>;
  };
}

export interface SearchSuggestion {
  id: string;
  type: 'product' | 'category' | 'query';
  text: string;
  highlight?: string;
  metadata?: Record<string, any>;
}

export const SearchService = {
  async searchListings(params: SearchParams): Promise<SearchResult> {
    const response = await apiClient.post<SearchResult>(
      '/api/v1/c2c/search',
      params
    );

    if (response.error) {
      throw new Error(`Failed to search listings: ${response.error}`);
    }

    return response.data!;
  },

  async getAutocompleteSuggestions(query: string): Promise<SearchSuggestion[]> {
    if (!query.trim()) {
      return [];
    }

    const response = await apiClient.get<SearchSuggestion[]>(
      `/api/v1/c2c/search/autocomplete?q=${encodeURIComponent(query.trim())}`
    );

    if (response.error) {
      throw new Error(`Failed to fetch suggestions: ${response.error}`);
    }

    return response.data || [];
  },

  async getSearchHistory(): Promise<string[]> {
    const response = await apiClient.get<string[]>(
      '/api/v1/user/search-history'
    );

    if (response.error) {
      throw new Error(`Failed to fetch search history: ${response.error}`);
    }

    return response.data || [];
  },

  async getTrendingSearches(): Promise<string[]> {
    const response = await apiClient.get<string[]>(
      '/api/v1/c2c/search/trending'
    );

    if (response.error) {
      throw new Error(`Failed to fetch trending searches: ${response.error}`);
    }

    return response.data || [];
  },

  async clearSearchHistory(): Promise<void> {
    const response = await apiClient.post('/api/v1/user/search-history/clear');

    if (response.error) {
      throw new Error(`Failed to clear search history: ${response.error}`);
    }
  },
};
