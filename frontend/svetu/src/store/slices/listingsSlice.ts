import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { Listing, ListingFilters } from '@/types/listing';
import listingsApi from '@/services/api/listings';

interface ListingsState {
  listings: Listing[];
  currentListing: Listing | null;
  similarListings: Listing[];
  favoriteListings: Listing[];
  totalElements: number;
  totalPages: number;
  currentPage: number;
  pageSize: number;
  loading: boolean;
  error: string | null;
  filters: ListingFilters;
  viewMode: 'grid' | 'list';
}

const initialState: ListingsState = {
  listings: [],
  currentListing: null,
  similarListings: [],
  favoriteListings: [],
  totalElements: 0,
  totalPages: 0,
  currentPage: 0,
  pageSize: 20,
  loading: false,
  error: null,
  filters: {},
  viewMode: 'grid'
};

// Async thunks
export const fetchListings = createAsyncThunk(
  'listings/fetchListings',
  async ({ filters, pagination }: { filters?: ListingFilters; pagination?: { page?: number; size?: number } }) => {
    const response = await listingsApi.getListings(filters || {}, pagination || { page: 0, size: 20 });
    return response;
  }
);

export const fetchListing = createAsyncThunk(
  'listings/fetchListing',
  async (id: string | number) => {
    const listing = await listingsApi.getListing(id);
    return listing;
  }
);

export const fetchSimilarListings = createAsyncThunk(
  'listings/fetchSimilarListings',
  async ({ id, limit }: { id: string | number; limit?: number }) => {
    const listings = await listingsApi.getSimilarListings(id, limit);
    return listings;
  }
);

export const fetchFavoriteListings = createAsyncThunk(
  'listings/fetchFavoriteListings',
  async (pagination?: { page?: number; size?: number }) => {
    const response = await listingsApi.getFavoriteListings(pagination || { page: 0, size: 20 });
    return response;
  }
);

export const toggleFavorite = createAsyncThunk(
  'listings/toggleFavorite',
  async (id: string | number) => {
    const isFavorite = await listingsApi.toggleFavorite(id);
    return { id, isFavorite };
  }
);

export const searchListings = createAsyncThunk(
  'listings/searchListings',
  async ({ query, filters, pagination }: { query: string; filters?: ListingFilters; pagination?: { page?: number; size?: number } }) => {
    const response = await listingsApi.searchListings(query, filters || {}, pagination || { page: 0, size: 20 });
    return response;
  }
);

const listingsSlice = createSlice({
  name: 'listings',
  initialState,
  reducers: {
    setFilters: (state, action: PayloadAction<ListingFilters>) => {
      state.filters = action.payload;
    },
    updateFilter: (state, action: PayloadAction<{ key: string; value: unknown }>) => {
      state.filters[action.payload.key] = action.payload.value;
    },
    clearFilters: (state) => {
      state.filters = {};
    },
    setViewMode: (state, action: PayloadAction<'grid' | 'list'>) => {
      state.viewMode = action.payload;
    },
    setCurrentPage: (state, action: PayloadAction<number>) => {
      state.currentPage = action.payload;
    },
    setPageSize: (state, action: PayloadAction<number>) => {
      state.pageSize = action.payload;
    },
    clearError: (state) => {
      state.error = null;
    }
  },
  extraReducers: (builder) => {
    // Fetch listings
    builder
      .addCase(fetchListings.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchListings.fulfilled, (state, action) => {
        state.loading = false;
        state.listings = action.payload.items;
        state.totalElements = action.payload.total;
        state.totalPages = action.payload.totalPages;
        state.currentPage = action.payload.page;
        state.pageSize = action.payload.pageSize;
      })
      .addCase(fetchListings.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch listings';
      });

    // Fetch single listing
    builder
      .addCase(fetchListing.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchListing.fulfilled, (state, action) => {
        state.loading = false;
        state.currentListing = action.payload;
      })
      .addCase(fetchListing.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch listing';
      });

    // Fetch similar listings
    builder
      .addCase(fetchSimilarListings.fulfilled, (state, action) => {
        state.similarListings = action.payload;
      });

    // Fetch favorite listings
    builder
      .addCase(fetchFavoriteListings.fulfilled, (state, action) => {
        state.favoriteListings = action.payload.items;
      });

    // Toggle favorite
    builder
      .addCase(toggleFavorite.fulfilled, (state, action) => {
        const { id, isFavorite } = action.payload;
        
        // Update the listing in the main list
        const listing = state.listings.find(l => l.id === id);
        if (listing) {
          listing.is_favorite = isFavorite;
        }
        
        // Update current listing if it's the same
        if (state.currentListing && state.currentListing.id === id) {
          state.currentListing.is_favorite = isFavorite;
        }
        
        // Update favorites list
        if (isFavorite) {
          const listingToAdd = state.listings.find(l => l.id === id) || state.currentListing;
          if (listingToAdd && !state.favoriteListings.find(l => l.id === id)) {
            state.favoriteListings.push(listingToAdd);
          }
        } else {
          state.favoriteListings = state.favoriteListings.filter(l => l.id !== id);
        }
      });

    // Search listings
    builder
      .addCase(searchListings.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(searchListings.fulfilled, (state, action) => {
        state.loading = false;
        state.listings = action.payload.items;
        state.totalElements = action.payload.total;
        state.totalPages = action.payload.totalPages;
        state.currentPage = action.payload.page;
        state.pageSize = action.payload.pageSize;
      })
      .addCase(searchListings.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to search listings';
      });
  }
});

export const {
  setFilters,
  updateFilter,
  clearFilters,
  setViewMode,
  setCurrentPage,
  setPageSize,
  clearError
} = listingsSlice.actions;

export default listingsSlice.reducer;