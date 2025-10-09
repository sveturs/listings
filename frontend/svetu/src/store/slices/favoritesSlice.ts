import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { api } from '@/lib/api';
import { components } from '@/types/generated/api';

type C2CListing = components['schemas']['models.MarketplaceListing'];

interface FavoritesState {
  items: C2CListing[];
  itemIds: Set<number>;
  loading: boolean;
  error: string | null;
  count: number;
}

const initialState: FavoritesState = {
  items: [],
  itemIds: new Set(),
  loading: false,
  error: null,
  count: 0,
};

// Async thunks
export const fetchFavorites = createAsyncThunk(
  'favorites/fetchFavorites',
  async () => {
    const response = await api.get('/c2c/favorites');
    return response.data;
  },
  {
    // Предотвращаем повторные запросы если уже загружается
    condition: (_, { getState }) => {
      const state = getState() as { favorites: FavoritesState };
      return !state.favorites.loading;
    },
  }
);

export const fetchFavoritesCount = createAsyncThunk(
  'favorites/fetchCount',
  async () => {
    const response = await api.get('/c2c/favorites/count');
    return response.data;
  }
);

export const addToFavorites = createAsyncThunk(
  'favorites/add',
  async ({ id, type = 'marketplace' }: { id: number; type?: string }) => {
    const url = `/c2c/favorites/${id}${type === 'storefront' ? '?type=storefront' : ''}`;
    const response = await api.post(url);
    return { id, type, data: response.data };
  }
);

export const removeFromFavorites = createAsyncThunk(
  'favorites/remove',
  async ({ id, type = 'marketplace' }: { id: number; type?: string }) => {
    const url = `/c2c/favorites/${id}${type === 'storefront' ? '?type=storefront' : ''}`;
    const response = await api.delete(url);
    return { id, type, data: response.data };
  }
);

export const checkIfInFavorites = createAsyncThunk(
  'favorites/check',
  async (id: number) => {
    const response = await api.get(`/c2c/favorites/${id}/check`);
    return { id, isInFavorites: response.data.is_in_favorites };
  }
);

const favoritesSlice = createSlice({
  name: 'favorites',
  initialState,
  reducers: {
    clearFavorites: (state) => {
      state.items = [];
      state.itemIds = new Set();
      state.count = 0;
      state.error = null;
    },
    setFavorites: (state, action: PayloadAction<C2CListing[]>) => {
      state.items = action.payload;
      state.itemIds = new Set(
        action.payload
          .map((item) => item.id)
          .filter((id): id is number => id !== undefined)
      );
      state.count = action.payload.length;
      state.error = null;
    },
    toggleFavoriteOptimistic: (state, action: PayloadAction<number>) => {
      const id = action.payload;
      if (state.itemIds.has(id)) {
        state.itemIds.delete(id);
        state.items = state.items.filter(
          (item) => item.id !== undefined && item.id !== id
        );
        state.count = Math.max(0, state.count - 1);
      } else {
        state.itemIds.add(id);
        state.count += 1;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch favorites
      .addCase(fetchFavorites.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchFavorites.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload || [];
        state.itemIds = new Set(
          (action.payload || []).map((item: C2CListing) => item.id)
        );
        state.count = (action.payload || []).length;
      })
      .addCase(fetchFavorites.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch favorites';
      })
      // Fetch count
      .addCase(fetchFavoritesCount.fulfilled, (state, action) => {
        state.count = action.payload.count || 0;
      })
      // Add to favorites
      .addCase(addToFavorites.pending, (state, action) => {
        // Optimistically add
        const id = action.meta.arg.id;
        state.itemIds.add(id);
        state.count += 1;
      })
      .addCase(addToFavorites.fulfilled, (state, action) => {
        // Confirm addition
        const id = action.payload.id;
        state.itemIds.add(id);
      })
      .addCase(addToFavorites.rejected, (state, action) => {
        // Rollback on error
        const id = action.meta.arg.id;
        state.itemIds.delete(id);
        state.count = Math.max(0, state.count - 1);
        state.error = action.error.message || 'Failed to add to favorites';
      })
      // Remove from favorites
      .addCase(removeFromFavorites.pending, (state, action) => {
        // Optimistically remove
        const id = action.meta.arg.id;
        state.itemIds.delete(id);
        state.items = state.items.filter((item) => item.id !== id);
        state.count = Math.max(0, state.count - 1);
      })
      .addCase(removeFromFavorites.fulfilled, (state, action) => {
        // Confirm removal
        const id = action.payload.id;
        state.itemIds.delete(id);
        state.items = state.items.filter((item) => item.id !== id);
      })
      .addCase(removeFromFavorites.rejected, (state, action) => {
        // Rollback on error - would need to refetch
        state.error = action.error.message || 'Failed to remove from favorites';
      })
      // Check if in favorites
      .addCase(checkIfInFavorites.fulfilled, (state, action) => {
        if (action.payload.isInFavorites) {
          state.itemIds.add(action.payload.id);
        } else {
          state.itemIds.delete(action.payload.id);
        }
      });
  },
});

export const { clearFavorites, setFavorites, toggleFavoriteOptimistic } =
  favoritesSlice.actions;
export default favoritesSlice.reducer;
