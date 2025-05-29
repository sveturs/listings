import { apiClient } from '@/lib/api-client';
import { Listing, ListingFilters, ListingSort, PaginatedResponse } from '@/types/listing';
import { getMockListings, getMockListingById } from '@/lib/mock-data';

const USE_MOCK_DATA = false; // Always use real API

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

    // Set default values
    const searchParams = {
      sort_by: 'date_desc',
      page: 1,
      size: 20, // Changed from 25 to match typical pagination
      ...params
    };

    // Add all search parameters
    Object.entries(searchParams).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        queryParams.append(key, value.toString());
      }
    });

    try {
      const response = await apiClient.get<{
        data?: Listing[];
        meta?: {
          total?: number;
          total_pages?: number;
          page?: number;
          size?: number;
          has_more?: boolean;
        };
      }>(`/api/v1/marketplace/search?${queryParams.toString()}`);
      
      console.log('API Response:', {
        url: `/api/v1/marketplace/search?${queryParams.toString()}`,
        responseData: response.data,
        meta: response.data?.meta
      });
      
      // Handle the response structure from the backend
      if (response.data) {
        const listings = response.data.data || [];
        const meta = response.data.meta || {};
        const total = meta.total || listings.length;
        const currentPage = meta.page || searchParams.page;
        const pageSize = meta.size || searchParams.size;
        
        // Calculate totalPages correctly
        // If backend provides total_pages, use it; otherwise calculate it
        let totalPages = meta.total_pages;
        if (!totalPages && total > 0) {
          totalPages = Math.ceil(total / pageSize);
        }
        
        const result = {
          items: listings,
          total,
          page: currentPage,
          pageSize: pageSize,
          totalPages: totalPages || 0
        };
        
        console.log('Processed result:', {
          ...result,
          calculatedTotalPages: Math.ceil(total / pageSize),
          metaTotalPages: meta.total_pages,
          hasMore: meta.has_more
        });
        return result;
      }
      
      // Fallback to empty response if structure is unexpected
      return {
        items: [],
        total: 0,
        page: searchParams.page,
        pageSize: searchParams.size,
        totalPages: 0
      };
    } catch (error) {
      console.error('Error searching listings:', error);
      // Return empty results on error
      return {
        items: [],
        total: 0,
        page: searchParams.page,
        pageSize: searchParams.size,
        totalPages: 0
      };
    }
  }
}

export const listingService = new ListingService();