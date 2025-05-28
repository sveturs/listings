import { apiClient } from '@/lib/api-client';
import { Listing, ListingFilters, ListingSort, PaginatedResponse } from '@/types/listing';
import { getMockListings, getMockListingById } from '@/lib/mock-data';

const USE_MOCK_DATA = process.env.NODE_ENV === 'development';

interface GetListingsParams {
  page?: number;
  pageSize?: number;
  filters?: ListingFilters;
  sort?: ListingSort;
}

interface SearchListingsParams {
  q?: string;
  category_id?: string;
  min_price?: string;
  max_price?: string;
  condition?: string;
  distance?: string;
  location?: string;
  page?: number;
  size?: number;
  sort_by?: string;
}

class ListingService {
  async getListings(params: GetListingsParams = {}): Promise<PaginatedResponse<Listing>> {
    const { page = 1, pageSize = 12, filters = {}, sort } = params;
    
    if (USE_MOCK_DATA) {
      return getMockListings(page, pageSize);
    }
    
    const queryParams = new URLSearchParams({
      page: page.toString(),
      pageSize: pageSize.toString(),
    });

    // Add filters
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value.toString());
      }
    });

    // Add sorting
    if (sort) {
      queryParams.append('sortBy', sort.field);
      queryParams.append('sortOrder', sort.order);
    }

    const response = await apiClient.get<PaginatedResponse<Listing>>(
      `/marketplace/listings?${queryParams.toString()}`
    );
    
    return response.data;
  }

  async getListingById(id: string): Promise<Listing> {
    if (USE_MOCK_DATA) {
      const listing = getMockListingById(id);
      if (!listing) {
        throw new Error('Listing not found');
      }
      return listing;
    }
    
    const response = await apiClient.get<Listing>(`/marketplace/listings/${id}`);
    return response.data;
  }

  async createListing(data: Partial<Listing>): Promise<Listing> {
    const response = await apiClient.post<Listing>('/marketplace/listings', data);
    return response.data;
  }

  async updateListing(id: string, data: Partial<Listing>): Promise<Listing> {
    const response = await apiClient.put<Listing>(`/marketplace/listings/${id}`, data);
    return response.data;
  }

  async deleteListing(id: string): Promise<void> {
    await apiClient.delete(`/marketplace/listings/${id}`);
  }

  async toggleFavorite(id: string): Promise<{ isFavorite: boolean }> {
    const response = await apiClient.post<{ isFavorite: boolean }>(
      `/marketplace/listings/${id}/favorite`
    );
    return response.data;
  }

  async getSimilarListings(id: string, limit: number = 4): Promise<Listing[]> {
    const response = await apiClient.get<Listing[]>(
      `/marketplace/listings/${id}/similar?limit=${limit}`
    );
    return response.data;
  }

  async searchListings(params: SearchListingsParams = {}): Promise<PaginatedResponse<Listing>> {
    const queryParams = new URLSearchParams();

    // Add all search parameters
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        queryParams.append(key, value.toString());
      }
    });

    try {
      const response = await apiClient.get<{
        data?: {
          data?: Listing[];
          total?: number;
          totalPages?: number;
        }
      }>(`/api/v1/marketplace/search?${queryParams.toString()}`);
      
      // Handle the response structure from the backend
      if (response.data && response.data.data) {
        const listings = Array.isArray(response.data.data) ? response.data.data : response.data.data.data || [];
        const total = response.data.data.total || listings.length;
        const totalPages = response.data.data.totalPages || Math.ceil(total / (params.size || 20));
        
        return {
          items: listings,
          total,
          page: params.page || 1,
          pageSize: params.size || 20,
          totalPages
        };
      }
      
      // Fallback to empty response if structure is unexpected
      return {
        items: [],
        total: 0,
        page: params.page || 1,
        pageSize: params.size || 20,
        totalPages: 0
      };
    } catch (error) {
      console.error('Error searching listings:', error);
      // Return empty results on error
      return {
        items: [],
        total: 0,
        page: params.page || 1,
        pageSize: params.size || 20,
        totalPages: 0
      };
    }
  }
}

export const listingService = new ListingService();