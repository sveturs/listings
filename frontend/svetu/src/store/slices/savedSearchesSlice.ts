import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { apiClient } from '@/services/api-client';

export interface SavedSearch {
  id: number;
  user_id: number;
  name: string;
  filters: Record<string, any>;
  search_type: string;
  notify_enabled: boolean;
  notify_frequency: string;
  results_count: number;
  last_notified_at?: string;
  created_at: string;
  updated_at: string;
}

interface SavedSearchesState {
  items: SavedSearch[];
  loading: boolean;
  error: string | null;
  executingId: number | null;
  executionResults: any | null;
}

const initialState: SavedSearchesState = {
  items: [],
  loading: false,
  error: null,
  executingId: null,
  executionResults: null,
};

// Load saved searches from localStorage on startup
const loadSavedSearchesFromStorage = (): SavedSearch[] => {
  if (typeof window === 'undefined') return [];
  try {
    const stored = localStorage.getItem('savedSearches');
    return stored ? JSON.parse(stored) : [];
  } catch {
    return [];
  }
};

// Save to localStorage
const saveSavedSearchesToStorage = (searches: SavedSearch[]) => {
  if (typeof window === 'undefined') return;
  try {
    localStorage.setItem('savedSearches', JSON.stringify(searches));
  } catch (e) {
    console.error('Failed to save searches to localStorage:', e);
  }
};

// Async thunks
export const fetchSavedSearches = createAsyncThunk(
  'savedSearches/fetch',
  async (searchType?: string) => {
    const params = searchType ? `?search_type=${searchType}` : '';
    const response = await apiClient.get(
      `/marketplace/saved-searches${params}`
    );
    return response.data;
  }
);

export const createSavedSearch = createAsyncThunk(
  'savedSearches/create',
  async (data: {
    name: string;
    filters: Record<string, any>;
    search_type?: string;
    notify_enabled?: boolean;
    notify_frequency?: string;
  }) => {
    const response = await apiClient.post('/marketplace/saved-searches', {
      ...data,
      search_type: data.search_type || 'cars',
      notify_enabled: data.notify_enabled || false,
      notify_frequency: data.notify_frequency || 'daily',
    });
    return response.data;
  }
);

export const updateSavedSearch = createAsyncThunk(
  'savedSearches/update',
  async ({
    id,
    data,
  }: {
    id: number;
    data: {
      name?: string;
      filters?: Record<string, any>;
      notify_enabled?: boolean;
      notify_frequency?: string;
    };
  }) => {
    const response = await apiClient.put(
      `/marketplace/saved-searches/${id}`,
      data
    );
    return response.data;
  }
);

export const deleteSavedSearch = createAsyncThunk(
  'savedSearches/delete',
  async (id: number) => {
    await apiClient.delete(`/marketplace/saved-searches/${id}`);
    return id;
  }
);

export const executeSavedSearch = createAsyncThunk(
  'savedSearches/execute',
  async (id: number) => {
    const response = await apiClient.get(
      `/marketplace/saved-searches/${id}/execute`
    );
    return { id, results: response.data };
  }
);

const savedSearchesSlice = createSlice({
  name: 'savedSearches',
  initialState: {
    ...initialState,
    items: loadSavedSearchesFromStorage(),
  },
  reducers: {
    clearSavedSearches: (state) => {
      state.items = [];
      state.error = null;
      state.executionResults = null;
      saveSavedSearchesToStorage([]);
    },
    clearExecutionResults: (state) => {
      state.executionResults = null;
      state.executingId = null;
    },
    addLocalSavedSearch: (
      state,
      action: PayloadAction<
        Omit<SavedSearch, 'id' | 'user_id' | 'created_at' | 'updated_at'>
      >
    ) => {
      // For non-authenticated users, save locally
      const newSearch: SavedSearch = {
        ...action.payload,
        id: Date.now(), // Use timestamp as temporary ID
        user_id: 0, // 0 for local searches
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
      state.items.push(newSearch);
      saveSavedSearchesToStorage(state.items);
    },
    removeLocalSavedSearch: (state, action: PayloadAction<number>) => {
      state.items = state.items.filter((item) => item.id !== action.payload);
      saveSavedSearchesToStorage(state.items);
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch saved searches
      .addCase(fetchSavedSearches.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchSavedSearches.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload || [];
        saveSavedSearchesToStorage(state.items);
      })
      .addCase(fetchSavedSearches.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch saved searches';
      })
      // Create saved search
      .addCase(createSavedSearch.pending, (state) => {
        state.error = null;
      })
      .addCase(createSavedSearch.fulfilled, (state, action) => {
        state.items.push(action.payload);
        saveSavedSearchesToStorage(state.items);
      })
      .addCase(createSavedSearch.rejected, (state, action) => {
        state.error = action.error.message || 'Failed to create saved search';
      })
      // Update saved search
      .addCase(updateSavedSearch.fulfilled, (state, action) => {
        const index = state.items.findIndex(
          (item) => item.id === action.payload.id
        );
        if (index !== -1) {
          state.items[index] = action.payload;
          saveSavedSearchesToStorage(state.items);
        }
      })
      .addCase(updateSavedSearch.rejected, (state, action) => {
        state.error = action.error.message || 'Failed to update saved search';
      })
      // Delete saved search
      .addCase(deleteSavedSearch.fulfilled, (state, action) => {
        state.items = state.items.filter((item) => item.id !== action.payload);
        saveSavedSearchesToStorage(state.items);
      })
      .addCase(deleteSavedSearch.rejected, (state, action) => {
        state.error = action.error.message || 'Failed to delete saved search';
      })
      // Execute saved search
      .addCase(executeSavedSearch.pending, (state, action) => {
        state.executingId = action.meta.arg;
        state.executionResults = null;
        state.error = null;
      })
      .addCase(executeSavedSearch.fulfilled, (state, action) => {
        state.executingId = null;
        state.executionResults = action.payload.results;
        // Update results count
        const index = state.items.findIndex(
          (item) => item.id === action.payload.id
        );
        if (index !== -1 && action.payload.results?.total) {
          state.items[index].results_count = action.payload.results.total;
          saveSavedSearchesToStorage(state.items);
        }
      })
      .addCase(executeSavedSearch.rejected, (state, action) => {
        state.executingId = null;
        state.error = action.error.message || 'Failed to execute saved search';
      });
  },
});

export const {
  clearSavedSearches,
  clearExecutionResults,
  addLocalSavedSearch,
  removeLocalSavedSearch,
} = savedSearchesSlice.actions;

export default savedSearchesSlice.reducer;
