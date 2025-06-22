import {
  createSlice,
  createAsyncThunk,
  createSelector,
  PayloadAction,
} from '@reduxjs/toolkit';
import { RootState } from '../index';
import { tokenManager } from '@/utils/tokenManager';
import type { components } from '@/types/generated/api';
import type { PaymentMethodType } from '@/types/storefront';

// Типы из сгенерированного API
type Storefront = components['schemas']['models.Storefront'];
type StorefrontCreateDTO = components['schemas']['models.StorefrontCreateDTO'];
type StorefrontUpdateDTO = components['schemas']['models.StorefrontUpdateDTO'];
type StorefrontAnalytics = components['schemas']['models.StorefrontAnalytics'];
type StorefrontRatingSummary =
  components['schemas']['models.StorefrontRatingSummary'];
type StorefrontDeliveryOption =
  components['schemas']['models.StorefrontDeliveryOption'];
type StorefrontPaymentMethod =
  components['schemas']['models.StorefrontPaymentMethod'];
type StorefrontHours = components['schemas']['models.StorefrontHours'];
type StorefrontsListResponse =
  components['schemas']['handler.StorefrontsListResponse'];

// Интерфейсы для фильтров и настроек
export interface StorefrontFilters {
  search?: string;
  city?: string;
  isActive?: boolean;
  isVerified?: boolean;
  minRating?: number;
  paymentMethods?: PaymentMethodType[];
  hasDelivery?: boolean;
  hasSelfPickup?: boolean;
  isOpenNow?: boolean;
  latitude?: number;
  longitude?: number;
  radiusKm?: number;
}

export interface PaginationParams {
  limit: number;
  offset: number;
}

// Состояние slice
export interface StorefrontState {
  // Основные данные
  storefronts: Storefront[];
  currentStorefront: Storefront | null;
  myStorefronts: Storefront[];

  // UI состояние
  isLoading: boolean;
  error: string | null;

  // Операции с отдельными действиями
  isCreating: boolean;
  isUpdating: boolean;
  isDeleting: boolean;

  // Пагинация и фильтры
  filters: StorefrontFilters;
  pagination: PaginationParams;
  totalCount: number;
  hasMore: boolean;

  // Дополнительные данные
  analytics: StorefrontAnalytics[] | null;
  ratingSummary: StorefrontRatingSummary | null;

  // Настройки витрины
  deliveryOptions: StorefrontDeliveryOption[];
  paymentMethods: StorefrontPaymentMethod[];
  workingHours: StorefrontHours[];

  // Состояние загрузки дополнительных данных
  isLoadingAnalytics: boolean;
  isLoadingSettings: boolean;
}

// Начальное состояние
const initialState: StorefrontState = {
  storefronts: [],
  currentStorefront: null,
  myStorefronts: [],

  isLoading: false,
  error: null,

  isCreating: false,
  isUpdating: false,
  isDeleting: false,

  filters: {
    search: '',
    city: '',
    isActive: true,
    isVerified: undefined,
    minRating: undefined,
    paymentMethods: [],
    hasDelivery: undefined,
    hasSelfPickup: undefined,
    isOpenNow: undefined,
  },

  pagination: {
    limit: 20,
    offset: 0,
  },

  totalCount: 0,
  hasMore: true,

  analytics: null,
  ratingSummary: null,

  deliveryOptions: [],
  paymentMethods: [],
  workingHours: [],

  isLoadingAnalytics: false,
  isLoadingSettings: false,
};

// Async thunks

// Получение списка витрин с фильтрами
export const fetchStorefronts = createAsyncThunk<
  StorefrontsListResponse,
  { filters?: StorefrontFilters; pagination?: PaginationParams },
  { rejectValue: string }
>(
  'storefronts/fetchStorefronts',
  async (
    { filters = {}, pagination = { limit: 20, offset: 0 } },
    { rejectWithValue }
  ) => {
    try {
      const queryParams = new URLSearchParams();

      // Добавляем пагинацию
      queryParams.append('limit', pagination.limit.toString());
      queryParams.append('offset', pagination.offset.toString());

      // Добавляем фильтры
      if (filters.search) queryParams.append('search', filters.search);
      if (filters.city) queryParams.append('city', filters.city);
      if (filters.isActive !== undefined)
        queryParams.append('is_active', filters.isActive.toString());
      if (filters.isVerified !== undefined)
        queryParams.append('is_verified', filters.isVerified.toString());
      if (filters.minRating)
        queryParams.append('min_rating', filters.minRating.toString());
      if (filters.hasDelivery !== undefined)
        queryParams.append('has_delivery', filters.hasDelivery.toString());
      if (filters.hasSelfPickup !== undefined)
        queryParams.append('has_self_pickup', filters.hasSelfPickup.toString());
      if (filters.isOpenNow !== undefined)
        queryParams.append('is_open_now', filters.isOpenNow.toString());

      // Геолокация
      if (filters.latitude && filters.longitude) {
        queryParams.append('latitude', filters.latitude.toString());
        queryParams.append('longitude', filters.longitude.toString());
        if (filters.radiusKm) {
          queryParams.append('radius_km', filters.radiusKm.toString());
        }
      }

      // Методы оплаты
      if (filters.paymentMethods && filters.paymentMethods.length > 0) {
        filters.paymentMethods.forEach((method) => {
          queryParams.append('payment_methods', method);
        });
      }

      const response = await fetch(
        `/api/v1/storefronts?${queryParams.toString()}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        const errorData = await response
          .json()
          .catch(() => ({ error: 'Network error' }));
        return rejectWithValue(
          errorData.error || `HTTP error! status: ${response.status}`
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      return rejectWithValue(
        error instanceof Error ? error.message : 'Unknown error'
      );
    }
  }
);

// Получение витрины по ID
export const fetchStorefrontById = createAsyncThunk<
  Storefront,
  number,
  { rejectValue: string }
>('storefronts/fetchStorefrontById', async (id, { rejectWithValue }) => {
  try {
    const response = await fetch(`/api/v1/storefronts/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ error: 'Network error' }));
      return rejectWithValue(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    const data = await response.json();
    return data;
  } catch (error) {
    return rejectWithValue(
      error instanceof Error ? error.message : 'Unknown error'
    );
  }
});

// Получение витрины по slug
export const fetchStorefrontBySlug = createAsyncThunk<
  Storefront,
  string,
  { rejectValue: string }
>('storefronts/fetchStorefrontBySlug', async (slug, { rejectWithValue }) => {
  try {
    const response = await fetch(`/api/v1/storefronts/slug/${slug}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ error: 'Network error' }));
      return rejectWithValue(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    const data = await response.json();
    return data;
  } catch (error) {
    return rejectWithValue(
      error instanceof Error ? error.message : 'Unknown error'
    );
  }
});

// Получение моих витрин
export const fetchMyStorefronts = createAsyncThunk<
  StorefrontsListResponse,
  void,
  { rejectValue: string }
>('storefronts/fetchMyStorefronts', async (_, { rejectWithValue }) => {
  try {
    const accessToken = tokenManager.getAccessToken();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    const response = await fetch('/api/v1/storefronts/my', {
      method: 'GET',
      headers,
      credentials: 'include', // Важно для отправки cookies
    });

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ error: 'Network error' }));
      return rejectWithValue(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    const data = await response.json();
    // API возвращает массив витрин, преобразуем в ожидаемый формат
    if (Array.isArray(data)) {
      return {
        storefronts: data,
        total: data.length,
        limit: data.length,
        offset: 0,
      };
    }
    return data;
  } catch (error) {
    return rejectWithValue(
      error instanceof Error ? error.message : 'Unknown error'
    );
  }
});

// Создание витрины
export const createStorefront = createAsyncThunk<
  Storefront,
  StorefrontCreateDTO,
  { rejectValue: string }
>(
  'storefronts/createStorefront',
  async (storefrontData, { rejectWithValue }) => {
    try {
      const response = await fetch('/api/v1/storefronts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(storefrontData),
      });

      if (!response.ok) {
        const errorData = await response
          .json()
          .catch(() => ({ error: 'Network error' }));
        return rejectWithValue(
          errorData.error || `HTTP error! status: ${response.status}`
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      return rejectWithValue(
        error instanceof Error ? error.message : 'Unknown error'
      );
    }
  }
);

// Обновление витрины
export const updateStorefront = createAsyncThunk<
  Storefront,
  { id: number; data: StorefrontUpdateDTO },
  { rejectValue: string }
>('storefronts/updateStorefront', async ({ id, data }, { rejectWithValue }) => {
  try {
    const response = await fetch(`/api/v1/storefronts/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ error: 'Network error' }));
      return rejectWithValue(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    const responseData = await response.json();
    return responseData;
  } catch (error) {
    return rejectWithValue(
      error instanceof Error ? error.message : 'Unknown error'
    );
  }
});

// Удаление витрины
export const deleteStorefront = createAsyncThunk<
  number,
  number,
  { rejectValue: string }
>('storefronts/deleteStorefront', async (id, { rejectWithValue }) => {
  try {
    const response = await fetch(`/api/v1/storefronts/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData = await response
        .json()
        .catch(() => ({ error: 'Network error' }));
      return rejectWithValue(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    return id;
  } catch (error) {
    return rejectWithValue(
      error instanceof Error ? error.message : 'Unknown error'
    );
  }
});

// Получение аналитики витрины
export const fetchStorefrontAnalytics = createAsyncThunk<
  StorefrontAnalytics[],
  { id: number; from?: string; to?: string },
  { rejectValue: string }
>(
  'storefronts/fetchStorefrontAnalytics',
  async ({ id, from, to }, { rejectWithValue }) => {
    try {
      const queryParams = new URLSearchParams();
      if (from) queryParams.append('from', from);
      if (to) queryParams.append('to', to);

      const response = await fetch(
        `/api/v1/storefronts/${id}/analytics?${queryParams.toString()}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        const errorData = await response
          .json()
          .catch(() => ({ error: 'Network error' }));
        return rejectWithValue(
          errorData.error || `HTTP error! status: ${response.status}`
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      return rejectWithValue(
        error instanceof Error ? error.message : 'Unknown error'
      );
    }
  }
);

// Создание slice
const storefrontSlice = createSlice({
  name: 'storefronts',
  initialState,
  reducers: {
    // Синхронные действия
    setFilters: (state, action: PayloadAction<Partial<StorefrontFilters>>) => {
      state.filters = { ...state.filters, ...action.payload };
      // Сбрасываем пагинацию при изменении фильтров
      state.pagination.offset = 0;
    },

    clearFilters: (state) => {
      state.filters = initialState.filters;
      state.pagination.offset = 0;
    },

    setPagination: (
      state,
      action: PayloadAction<Partial<PaginationParams>>
    ) => {
      state.pagination = { ...state.pagination, ...action.payload };
    },

    setCurrentStorefront: (state, action: PayloadAction<Storefront | null>) => {
      state.currentStorefront = action.payload;
    },

    clearCurrentStorefront: (state) => {
      state.currentStorefront = null;
      state.analytics = null;
      state.ratingSummary = null;
      state.deliveryOptions = [];
      state.paymentMethods = [];
      state.workingHours = [];
    },

    clearError: (state) => {
      state.error = null;
    },

    // Обновление отдельной витрины в списках
    updateStorefrontInLists: (state, action: PayloadAction<Storefront>) => {
      const updatedStorefront = action.payload;

      // Обновляем в общем списке
      const storefrontIndex = state.storefronts.findIndex(
        (s) => s.id === updatedStorefront.id
      );
      if (storefrontIndex !== -1) {
        state.storefronts[storefrontIndex] = updatedStorefront;
      }

      // Обновляем в моих витринах
      const myStorefrontIndex = state.myStorefronts.findIndex(
        (s) => s.id === updatedStorefront.id
      );
      if (myStorefrontIndex !== -1) {
        state.myStorefronts[myStorefrontIndex] = updatedStorefront;
      }

      // Обновляем текущую витрину если это она
      if (state.currentStorefront?.id === updatedStorefront.id) {
        state.currentStorefront = updatedStorefront;
      }
    },
  },

  extraReducers: (builder) => {
    // Получение списка витрин
    builder
      .addCase(fetchStorefronts.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchStorefronts.fulfilled, (state, action) => {
        state.isLoading = false;

        const { storefronts = [], total = 0, offset = 0 } = action.payload;

        // Если это первая страница или новый поиск - заменяем, иначе добавляем
        if (offset === 0) {
          state.storefronts = storefronts;
        } else {
          state.storefronts = [...state.storefronts, ...storefronts];
        }

        state.totalCount = total;
        state.hasMore = state.storefronts.length < total;
      })
      .addCase(fetchStorefronts.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload || 'Failed to fetch storefronts';
      });

    // Получение витрины по ID
    builder
      .addCase(fetchStorefrontById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchStorefrontById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentStorefront = action.payload;
      })
      .addCase(fetchStorefrontById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload || 'Failed to fetch storefront';
      });

    // Получение витрины по slug
    builder
      .addCase(fetchStorefrontBySlug.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchStorefrontBySlug.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentStorefront = action.payload;
      })
      .addCase(fetchStorefrontBySlug.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload || 'Failed to fetch storefront by slug';
      });

    // Получение моих витрин
    builder
      .addCase(fetchMyStorefronts.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchMyStorefronts.fulfilled, (state, action) => {
        state.isLoading = false;
        state.myStorefronts = action.payload.storefronts || [];
      })
      .addCase(fetchMyStorefronts.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload || 'Failed to fetch my storefronts';
      });

    // Создание витрины
    builder
      .addCase(createStorefront.pending, (state) => {
        state.isCreating = true;
        state.error = null;
      })
      .addCase(createStorefront.fulfilled, (state, action) => {
        state.isCreating = false;
        // Добавляем новую витрину в начало списка моих витрин
        state.myStorefronts.unshift(action.payload);
        // Также добавляем в общий список если он загружен
        if (state.storefronts.length > 0) {
          state.storefronts.unshift(action.payload);
          state.totalCount += 1;
        }
      })
      .addCase(createStorefront.rejected, (state, action) => {
        state.isCreating = false;
        state.error = action.payload || 'Failed to create storefront';
      });

    // Обновление витрины
    builder
      .addCase(updateStorefront.pending, (state) => {
        state.isUpdating = true;
        state.error = null;
      })
      .addCase(updateStorefront.fulfilled, (state, action) => {
        state.isUpdating = false;
        // Используем reducer для обновления во всех списках
        storefrontSlice.caseReducers.updateStorefrontInLists(state, action);
      })
      .addCase(updateStorefront.rejected, (state, action) => {
        state.isUpdating = false;
        state.error = action.payload || 'Failed to update storefront';
      });

    // Удаление витрины
    builder
      .addCase(deleteStorefront.pending, (state) => {
        state.isDeleting = true;
        state.error = null;
      })
      .addCase(deleteStorefront.fulfilled, (state, action) => {
        state.isDeleting = false;
        const deletedId = action.payload;

        // Удаляем из всех списков
        state.storefronts = state.storefronts.filter((s) => s.id !== deletedId);
        state.myStorefronts = state.myStorefronts.filter(
          (s) => s.id !== deletedId
        );

        // Обновляем счетчик
        state.totalCount = Math.max(0, state.totalCount - 1);

        // Сбрасываем текущую витрину если она была удалена
        if (state.currentStorefront?.id === deletedId) {
          state.currentStorefront = null;
          state.analytics = null;
          state.ratingSummary = null;
          state.deliveryOptions = [];
          state.paymentMethods = [];
          state.workingHours = [];
        }
      })
      .addCase(deleteStorefront.rejected, (state, action) => {
        state.isDeleting = false;
        state.error = action.payload || 'Failed to delete storefront';
      });

    // Получение аналитики
    builder
      .addCase(fetchStorefrontAnalytics.pending, (state) => {
        state.isLoadingAnalytics = true;
      })
      .addCase(fetchStorefrontAnalytics.fulfilled, (state, action) => {
        state.isLoadingAnalytics = false;
        // API возвращает массив аналитики, берем последний элемент для совместимости
        state.analytics = Array.isArray(action.payload)
          ? action.payload
          : [action.payload];
      })
      .addCase(fetchStorefrontAnalytics.rejected, (state, action) => {
        state.isLoadingAnalytics = false;
        state.error = action.payload || 'Failed to fetch analytics';
      });
  },
});

// Экспорт действий
export const {
  setFilters,
  clearFilters,
  setPagination,
  setCurrentStorefront,
  clearCurrentStorefront,
  clearError,
  updateStorefrontInLists,
} = storefrontSlice.actions;

// Селекторы
export const selectStorefronts = (state: RootState) =>
  state.storefronts.storefronts;
export const selectCurrentStorefront = (state: RootState) =>
  state.storefronts.currentStorefront;
export const selectMyStorefronts = (state: RootState) =>
  state.storefronts.myStorefronts;
export const selectIsLoading = (state: RootState) =>
  state.storefronts.isLoading;
export const selectIsCreating = (state: RootState) =>
  state.storefronts.isCreating;
export const selectIsUpdating = (state: RootState) =>
  state.storefronts.isUpdating;
export const selectIsDeleting = (state: RootState) =>
  state.storefronts.isDeleting;
export const selectError = (state: RootState) => state.storefronts.error;
export const selectFilters = (state: RootState) => state.storefronts.filters;
export const selectPagination = (state: RootState) =>
  state.storefronts.pagination;
export const selectTotalCount = (state: RootState) =>
  state.storefronts.totalCount;
export const selectHasMore = (state: RootState) => state.storefronts.hasMore;
export const selectAnalytics = (state: RootState) =>
  state.storefronts.analytics;
export const selectIsLoadingAnalytics = (state: RootState) =>
  state.storefronts.isLoadingAnalytics;

// Фильтрованные селекторы (мемоизированные)

export const selectVerifiedStorefronts = createSelector(
  [selectStorefronts],
  (storefronts) => storefronts.filter((storefront) => storefront.is_verified)
);

export const selectStorefrontsByCity = createSelector(
  [selectStorefronts, (state: RootState, city: string) => city],
  (storefronts, city) =>
    storefronts.filter((storefront) => storefront.city === city)
);

export const selectActiveStorefronts = createSelector(
  [selectStorefronts],
  (storefronts) => storefronts.filter((storefront) => storefront.is_active)
);

// Экспорт reducer
export default storefrontSlice.reducer;
