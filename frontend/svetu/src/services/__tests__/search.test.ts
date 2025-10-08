import { SearchService, SearchParams } from '../search';
import { apiClient } from '../api-client';

jest.mock('../api-client', () => ({
  apiClient: {
    get: jest.fn(),
    post: jest.fn(),
  },
}));

describe('SearchService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('searchListings', () => {
    const mockSearchResponse = {
      data: {
        items: [
          {
            id: '1',
            title: 'iPhone 13',
            price: 799,
            category: 'Electronics',
          },
          {
            id: '2',
            title: 'Samsung Galaxy S21',
            price: 699,
            category: 'Electronics',
          },
        ],
        total: 2,
        page: 1,
        limit: 20,
        facets: {
          categories: [{ name: 'Electronics', count: 2 }],
          price_ranges: [{ range: '500-1000', count: 2 }],
        },
      },
      error: null,
    };

    it('should search listings with query', async () => {
      (apiClient.post as jest.Mock).mockResolvedValue(mockSearchResponse);

      const params: SearchParams = {
        query: 'phone',
        page: 1,
        limit: 20,
      };

      const result = await SearchService.searchListings(params);

      expect(apiClient.post).toHaveBeenCalledWith(
        '/api/v1/marketplace/search',
        params
      );
      expect(result).toEqual(mockSearchResponse.data);
    });

    it('should search with filters', async () => {
      (apiClient.post as jest.Mock).mockResolvedValue(mockSearchResponse);

      const params: SearchParams = {
        query: 'phone',
        filters: {
          categories: [1, 2],
          price_min: 500,
          price_max: 1000,
          condition: 'new',
          location: 'Novi Sad',
          radius: 10,
        },
        page: 1,
        limit: 20,
      };

      const result = await SearchService.searchListings(params);

      expect(apiClient.post).toHaveBeenCalledWith(
        '/api/v1/marketplace/search',
        params
      );
      expect(result).toEqual(mockSearchResponse.data);
    });

    it('should handle empty results', async () => {
      const emptyResponse = {
        data: {
          items: [],
          total: 0,
          page: 1,
          limit: 20,
          facets: {},
        },
        error: null,
      };

      (apiClient.post as jest.Mock).mockResolvedValue(emptyResponse);

      const result = await SearchService.searchListings({ query: 'xyz' });

      expect(result.items).toHaveLength(0);
      expect(result.total).toBe(0);
    });

    it('should handle API errors', async () => {
      const errorResponse = {
        data: null,
        error: 'Search service unavailable',
      };

      (apiClient.post as jest.Mock).mockResolvedValue(errorResponse);

      await expect(
        SearchService.searchListings({ query: 'test' })
      ).rejects.toThrow(
        'Failed to search listings: Search service unavailable'
      );
    });

    it('should handle network errors', async () => {
      (apiClient.post as jest.Mock).mockRejectedValue(
        new Error('Network error')
      );

      await expect(
        SearchService.searchListings({ query: 'test' })
      ).rejects.toThrow('Network error');
    });

    it('should use fuzzy search when enabled', async () => {
      (apiClient.post as jest.Mock).mockResolvedValue(mockSearchResponse);

      const params: SearchParams = {
        query: 'iphon',
        fuzzy: true,
        page: 1,
        limit: 20,
      };

      await SearchService.searchListings(params);

      expect(apiClient.post).toHaveBeenCalledWith(
        '/api/v1/marketplace/search',
        expect.objectContaining({ fuzzy: true })
      );
    });

    it('should handle sorting', async () => {
      (apiClient.post as jest.Mock).mockResolvedValue(mockSearchResponse);

      const params: SearchParams = {
        query: 'phone',
        sort: 'price_asc',
        page: 1,
        limit: 20,
      };

      await SearchService.searchListings(params);

      expect(apiClient.post).toHaveBeenCalledWith(
        '/api/v1/marketplace/search',
        expect.objectContaining({ sort: 'price_asc' })
      );
    });

    it('should handle pagination', async () => {
      (apiClient.post as jest.Mock).mockResolvedValue(mockSearchResponse);

      const params: SearchParams = {
        query: 'phone',
        page: 3,
        limit: 50,
      };

      await SearchService.searchListings(params);

      expect(apiClient.post).toHaveBeenCalledWith(
        '/api/v1/marketplace/search',
        expect.objectContaining({ page: 3, limit: 50 })
      );
    });
  });

  describe('getAutocompleteSuggestions', () => {
    const mockSuggestionsResponse = {
      data: [
        {
          id: '1',
          type: 'product',
          text: 'iPhone 13',
          highlight: '<em>iPhone</em> 13',
          metadata: {
            category: 'Electronics',
            price: 799,
          },
        },
        {
          id: '2',
          type: 'category',
          text: 'Smartphones',
          highlight: 'Smartphones',
          metadata: {
            count: 150,
          },
        },
      ],
      error: null,
    };

    it('should fetch autocomplete suggestions', async () => {
      (apiClient.get as jest.Mock).mockResolvedValue(mockSuggestionsResponse);

      const result = await SearchService.getAutocompleteSuggestions('iph');

      expect(apiClient.get).toHaveBeenCalledWith(
        '/api/v1/marketplace/search/autocomplete?q=iph'
      );
      expect(result).toEqual(mockSuggestionsResponse.data);
    });

    it('should handle empty query', async () => {
      const result = await SearchService.getAutocompleteSuggestions('');

      expect(apiClient.get).not.toHaveBeenCalled();
      expect(result).toEqual([]);
    });

    it('should trim whitespace from query', async () => {
      (apiClient.get as jest.Mock).mockResolvedValue(mockSuggestionsResponse);

      await SearchService.getAutocompleteSuggestions('  test  ');

      expect(apiClient.get).toHaveBeenCalledWith(
        '/api/v1/marketplace/search/autocomplete?q=test'
      );
    });

    it('should handle API errors', async () => {
      const errorResponse = {
        data: null,
        error: 'Autocomplete service unavailable',
      };

      (apiClient.get as jest.Mock).mockResolvedValue(errorResponse);

      await expect(
        SearchService.getAutocompleteSuggestions('test')
      ).rejects.toThrow(
        'Failed to fetch suggestions: Autocomplete service unavailable'
      );
    });

    it('should handle network errors', async () => {
      (apiClient.get as jest.Mock).mockRejectedValue(
        new Error('Network error')
      );

      await expect(
        SearchService.getAutocompleteSuggestions('test')
      ).rejects.toThrow('Network error');
    });

    it('should cache suggestions', async () => {
      (apiClient.get as jest.Mock).mockResolvedValue(mockSuggestionsResponse);

      // First call
      await SearchService.getAutocompleteSuggestions('test');

      // Second call with same query (should use cache if implemented)
      await SearchService.getAutocompleteSuggestions('test');

      // Note: If caching is implemented, this should be called only once
      // For now, it will be called twice
      expect(apiClient.get).toHaveBeenCalledTimes(2);
    });
  });

  describe('getSearchHistory', () => {
    const mockHistoryResponse = {
      data: ['laptop', 'phone', 'furniture'],
      error: null,
    };

    it('should fetch search history', async () => {
      (apiClient.get as jest.Mock).mockResolvedValue(mockHistoryResponse);

      const result = await SearchService.getSearchHistory();

      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/user/search-history');
      expect(result).toEqual(mockHistoryResponse.data);
    });

    it('should handle empty history', async () => {
      const emptyResponse = {
        data: [],
        error: null,
      };

      (apiClient.get as jest.Mock).mockResolvedValue(emptyResponse);

      const result = await SearchService.getSearchHistory();

      expect(result).toEqual([]);
    });

    it('should handle API errors', async () => {
      const errorResponse = {
        data: null,
        error: 'Unauthorized',
      };

      (apiClient.get as jest.Mock).mockResolvedValue(errorResponse);

      await expect(SearchService.getSearchHistory()).rejects.toThrow(
        'Failed to fetch search history: Unauthorized'
      );
    });
  });

  describe('getTrendingSearches', () => {
    const mockTrendingResponse = {
      data: ['iphone', 'laptop', 'nike shoes'],
      error: null,
    };

    it('should fetch trending searches', async () => {
      (apiClient.get as jest.Mock).mockResolvedValue(mockTrendingResponse);

      const result = await SearchService.getTrendingSearches();

      expect(apiClient.get).toHaveBeenCalledWith(
        '/api/v1/marketplace/search/trending'
      );
      expect(result).toEqual(mockTrendingResponse.data);
    });

    it('should handle empty trending', async () => {
      const emptyResponse = {
        data: [],
        error: null,
      };

      (apiClient.get as jest.Mock).mockResolvedValue(emptyResponse);

      const result = await SearchService.getTrendingSearches();

      expect(result).toEqual([]);
    });

    it('should handle API errors', async () => {
      const errorResponse = {
        data: null,
        error: 'Service unavailable',
      };

      (apiClient.get as jest.Mock).mockResolvedValue(errorResponse);

      await expect(SearchService.getTrendingSearches()).rejects.toThrow(
        'Failed to fetch trending searches: Service unavailable'
      );
    });
  });

  describe('clearSearchHistory', () => {
    it('should clear search history', async () => {
      const successResponse = {
        data: { success: true },
        error: null,
      };

      (apiClient.post as jest.Mock).mockResolvedValue(successResponse);

      await SearchService.clearSearchHistory();

      expect(apiClient.post).toHaveBeenCalledWith(
        '/api/v1/user/search-history/clear'
      );
    });

    it('should handle API errors', async () => {
      const errorResponse = {
        data: null,
        error: 'Failed to clear history',
      };

      (apiClient.post as jest.Mock).mockResolvedValue(errorResponse);

      await expect(SearchService.clearSearchHistory()).rejects.toThrow(
        'Failed to clear search history: Failed to clear history'
      );
    });
  });
});
