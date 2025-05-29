// API service for listings
import axios from 'axios';
import { Listing, ListingFilters, PaginatedResponse } from '@/types/listing';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Add request interceptor to include language
api.interceptors.request.use(
  (config) => {
    // Get current language from next-intl or localStorage
    const currentLanguage = localStorage.getItem('locale') || 'en';
    
    // Add language parameter to GET requests
    if (config.method === 'get') {
      config.params = config.params || {};
      config.params.language = currentLanguage;
    }
    
    // Add Accept-Language header
    config.headers['Accept-Language'] = currentLanguage;
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Pagination params interface
interface PaginationParams {
  page?: number;
  size?: number;
  sort?: string;
  sort_by?: string;
}

// Listing API endpoints
export const listingsApi = {
  // Get listings with filters and pagination
  getListings: async (filters: ListingFilters = {}, pagination: PaginationParams = { page: 0, size: 20 }): Promise<PaginatedResponse<Listing>> => {
    const params = {
      ...filters,
      page: pagination.page,
      size: pagination.size,
      sort_by: pagination.sort_by || 'created_at,desc'
    };

    // Remove empty values
    const cleanParams: Record<string, unknown> = {};
    Object.entries(params).forEach(([key, value]) => {
      if (value !== '' && value !== null && value !== undefined) {
        cleanParams[key] = value;
      }
    });

    const response = await api.get<{ data: PaginatedResponse<Listing> }>('/api/v1/marketplace/listings', { params: cleanParams });
    return response.data.data;
  },

  // Get single listing by ID
  getListing: async (id: string | number): Promise<Listing> => {
    const response = await api.get<{ data: Listing }>(`/api/v1/marketplace/listings/${id}`);
    return response.data.data;
  },

  // Get similar listings
  getSimilarListings: async (id: string | number, limit: number = 8): Promise<Listing[]> => {
    const response = await api.get<{ data: Listing[] }>(`/api/v1/marketplace/listings/${id}/similar`, {
      params: { limit }
    });
    return response.data.data;
  },

  // Get featured listings
  getFeaturedListings: async (limit: number = 12): Promise<Listing[]> => {
    const response = await api.get<{ data: { content: Listing[] } }>('/api/v1/marketplace/listings/featured', {
      params: { size: limit }
    });
    return response.data.data.content;
  },

  // Get listings by category
  getListingsByCategory: async (categoryId: string | number, pagination: PaginationParams = { page: 0, size: 20 }): Promise<PaginatedResponse<Listing>> => {
    return listingsApi.getListings({ category_id: categoryId }, pagination);
  },

  // Search listings
  searchListings: async (query: string, filters: ListingFilters = {}, pagination: PaginationParams = { page: 0, size: 20 }): Promise<PaginatedResponse<Listing>> => {
    return listingsApi.getListings({ ...filters, query }, pagination);
  },

  // Get user's listings
  getUserListings: async (userId: string | number, pagination: PaginationParams = { page: 0, size: 20 }): Promise<PaginatedResponse<Listing>> => {
    const response = await api.get<{ data: PaginatedResponse<Listing> }>(`/api/v1/users/${userId}/listings`, {
      params: {
        page: pagination.page,
        size: pagination.size
      }
    });
    return response.data.data;
  },

  // Get my listings (authenticated user)
  getMyListings: async (pagination: PaginationParams = { page: 0, size: 20 }): Promise<PaginatedResponse<Listing>> => {
    const response = await api.get<{ data: PaginatedResponse<Listing> }>('/api/v1/marketplace/my-listings', {
      params: {
        page: pagination.page,
        size: pagination.size
      }
    });
    return response.data.data;
  },

  // Get favorite listings
  getFavoriteListings: async (pagination: PaginationParams = { page: 0, size: 20 }): Promise<PaginatedResponse<Listing>> => {
    const response = await api.get<{ data: PaginatedResponse<Listing> }>('/api/v1/marketplace/favorites', {
      params: {
        page: pagination.page,
        size: pagination.size
      }
    });
    return response.data.data;
  },

  // Create new listing
  createListing: async (data: Partial<Listing>): Promise<Listing> => {
    const response = await api.post<{ data: Listing }>('/api/v1/marketplace/listings', data);
    return response.data.data;
  },

  // Update listing
  updateListing: async (id: string | number, data: Partial<Listing>): Promise<Listing> => {
    const response = await api.put<{ data: Listing }>(`/api/v1/marketplace/listings/${id}`, data);
    return response.data.data;
  },

  // Delete listing
  deleteListing: async (id: string | number): Promise<void> => {
    await api.delete(`/api/v1/marketplace/listings/${id}`);
  },

  // Toggle favorite
  toggleFavorite: async (id: string | number): Promise<{ is_favorite: boolean }> => {
    const response = await api.post<{ data: { is_favorite: boolean } }>(`/api/v1/marketplace/listings/${id}/favorite`);
    return response.data.data;
  },

  // Upload images
  uploadImages: async (id: string | number, files: File[]): Promise<string[]> => {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append('images[]', file);
    });

    const response = await api.post<{ data: { urls: string[] } }>(
      `/api/v1/marketplace/listings/${id}/images`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      }
    );
    return response.data.data.urls;
  },

  // Delete image
  deleteImage: async (listingId: string | number, imageId: string | number): Promise<void> => {
    await api.delete(`/api/v1/marketplace/listings/${listingId}/images/${imageId}`);
  }
};

export default listingsApi;