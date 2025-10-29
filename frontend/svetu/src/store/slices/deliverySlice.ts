/**
 * Delivery Redux Slice
 *
 * Централизованное управление состоянием delivery системы:
 * - Кэширование провайдеров
 * - Кэширование расчетов доставки с TTL
 * - Управление выбранными quotes
 * - Tracking информация
 */

import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';
import { deliveryService } from '@/services/delivery';
import type {
  DeliveryQuote,
  CalculationRequest,
  CalculationResponse,
} from '@/types/delivery';
import type {
  DeliveryProvider as ServiceDeliveryProvider,
  TrackingInfo,
} from '@/services/delivery';
import type { RootState } from '@/store';

// Используем тип провайдера из services (более простой, соответствует API)
type DeliveryProvider = ServiceDeliveryProvider;

// ========================================
// TYPES
// ========================================

interface CachedCalculation {
  data: CalculationResponse;
  timestamp: number;
  params: CalculationRequest;
}

interface DeliveryState {
  // Providers
  providers: DeliveryProvider[];
  providersLoading: boolean;
  providersError: string | null;

  // Selected quotes by storefront
  selectedQuotes: Record<string, DeliveryQuote>;

  // Calculation cache with TTL
  calculations: Record<string, CachedCalculation>;
  calculationsLoading: Record<string, boolean>;
  calculationsError: Record<string, string>;

  // Tracking info
  tracking: Record<string, TrackingInfo>;
  trackingLoading: Record<string, boolean>;
  trackingError: Record<string, string>;
}

// Cache TTL: 5 minutes
const CACHE_TTL = 5 * 60 * 1000;

// ========================================
// INITIAL STATE
// ========================================

const initialState: DeliveryState = {
  providers: [],
  providersLoading: false,
  providersError: null,
  selectedQuotes: {},
  calculations: {},
  calculationsLoading: {},
  calculationsError: {},
  tracking: {},
  trackingLoading: {},
  trackingError: {},
};

// ========================================
// HELPERS
// ========================================

/**
 * Проверяет валидность кэша
 */
export function isCacheValid(
  cached: CachedCalculation,
  newParams: CalculationRequest
): boolean {
  // Check TTL
  if (Date.now() - cached.timestamp > CACHE_TTL) {
    return false;
  }

  // Compare params (deep comparison)
  return JSON.stringify(cached.params) === JSON.stringify(newParams);
}

/**
 * Генерирует ключ кэша для calculation request
 */
function getCacheKey(request: CalculationRequest): string {
  return JSON.stringify({
    from: request.from_location,
    to: request.to_location,
    items: request.items,
    provider_id: request.provider_id,
    insurance_value: request.insurance_value,
    cod_amount: request.cod_amount,
  });
}

// ========================================
// ASYNC THUNKS
// ========================================

/**
 * Загрузить список доступных провайдеров
 */
export const fetchProviders = createAsyncThunk(
  'delivery/fetchProviders',
  async (_, { rejectWithValue }) => {
    try {
      const response = await deliveryService.getProviders();

      if (response.error) {
        return rejectWithValue(response.error.message);
      }

      return response.data || [];
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch providers');
    }
  }
);

/**
 * Рассчитать стоимость доставки (unified endpoint)
 */
export const calculateRate = createAsyncThunk(
  'delivery/calculateRate',
  async (
    {
      request,
      bypassCache = false,
    }: { request: CalculationRequest; bypassCache?: boolean },
    { getState, rejectWithValue }
  ) => {
    try {
      const state = getState() as RootState;
      const cacheKey = getCacheKey(request);

      // Check cache (if not bypassed)
      if (!bypassCache) {
        const cached = state.delivery.calculations[cacheKey];
        if (cached && isCacheValid(cached, request)) {
          console.log('[deliverySlice] Using cached calculation');
          return { cacheKey, data: cached.data, fromCache: true };
        }
      }

      // Make API call
      const response = await fetch('/api/v2/delivery/calculate-universal', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return rejectWithValue(
          errorData.error?.message || 'Failed to calculate rate'
        );
      }

      const data = await response.json();

      return {
        cacheKey,
        data: {
          success: true,
          data,
        } as CalculationResponse,
        fromCache: false,
      };
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to calculate rate');
    }
  }
);

/**
 * Отследить отправление
 */
export const trackShipment = createAsyncThunk(
  'delivery/trackShipment',
  async (trackingToken: string, { rejectWithValue }) => {
    try {
      const response = await deliveryService.trackShipment(trackingToken);

      if (response.error) {
        return rejectWithValue(response.error.message);
      }

      return { trackingToken, data: response.data! };
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to track shipment');
    }
  }
);

// ========================================
// SLICE
// ========================================

export const deliverySlice = createSlice({
  name: 'delivery',
  initialState,
  reducers: {
    /**
     * Выбрать quote для конкретного storefront
     */
    selectQuote: (
      state,
      action: PayloadAction<{ storefrontId: string; quote: DeliveryQuote }>
    ) => {
      state.selectedQuotes[action.payload.storefrontId] = action.payload.quote;
    },

    /**
     * Очистить quote для storefront
     */
    clearQuote: (state, action: PayloadAction<string>) => {
      delete state.selectedQuotes[action.payload];
    },

    /**
     * Очистить все выбранные quotes
     */
    clearAllQuotes: (state) => {
      state.selectedQuotes = {};
    },

    /**
     * Очистить кэш расчетов
     */
    clearCalculationsCache: (state) => {
      state.calculations = {};
      state.calculationsError = {};
    },

    /**
     * Очистить tracking кэш
     */
    clearTrackingCache: (state) => {
      state.tracking = {};
      state.trackingError = {};
    },
  },

  extraReducers: (builder) => {
    // ========================================
    // FETCH PROVIDERS
    // ========================================
    builder.addCase(fetchProviders.pending, (state) => {
      state.providersLoading = true;
      state.providersError = null;
    });

    builder.addCase(fetchProviders.fulfilled, (state, action) => {
      state.providers = action.payload;
      state.providersLoading = false;
      state.providersError = null;
    });

    builder.addCase(fetchProviders.rejected, (state, action) => {
      state.providersLoading = false;
      state.providersError =
        (action.payload as string) || 'Failed to fetch providers';
    });

    // ========================================
    // CALCULATE RATE
    // ========================================
    builder.addCase(calculateRate.pending, (state, action) => {
      const cacheKey = getCacheKey(action.meta.arg.request);
      state.calculationsLoading[cacheKey] = true;
      delete state.calculationsError[cacheKey];
    });

    builder.addCase(calculateRate.fulfilled, (state, action) => {
      const { cacheKey, data } = action.payload;

      // Store in cache with timestamp
      state.calculations[cacheKey] = {
        data,
        timestamp: Date.now(),
        params: action.meta.arg.request,
      };

      state.calculationsLoading[cacheKey] = false;
      delete state.calculationsError[cacheKey];
    });

    builder.addCase(calculateRate.rejected, (state, action) => {
      const cacheKey = getCacheKey(action.meta.arg.request);
      state.calculationsLoading[cacheKey] = false;
      state.calculationsError[cacheKey] =
        (action.payload as string) || 'Failed to calculate rate';
    });

    // ========================================
    // TRACK SHIPMENT
    // ========================================
    builder.addCase(trackShipment.pending, (state, action) => {
      state.trackingLoading[action.meta.arg] = true;
      delete state.trackingError[action.meta.arg];
    });

    builder.addCase(trackShipment.fulfilled, (state, action) => {
      const { trackingToken, data } = action.payload;
      state.tracking[trackingToken] = data;
      state.trackingLoading[trackingToken] = false;
      delete state.trackingError[trackingToken];
    });

    builder.addCase(trackShipment.rejected, (state, action) => {
      const trackingToken = action.meta.arg;
      state.trackingLoading[trackingToken] = false;
      state.trackingError[trackingToken] =
        (action.payload as string) || 'Failed to track shipment';
    });
  },
});

// ========================================
// ACTIONS
// ========================================

export const {
  selectQuote,
  clearQuote,
  clearAllQuotes,
  clearCalculationsCache,
  clearTrackingCache,
} = deliverySlice.actions;

// ========================================
// SELECTORS
// ========================================

/**
 * Получить список провайдеров
 */
export const selectProviders = (state: RootState) => state.delivery.providers;

/**
 * Получить loading state провайдеров
 */
export const selectProvidersLoading = (state: RootState) =>
  state.delivery.providersLoading;

/**
 * Получить выбранный quote для storefront
 */
export const selectSelectedQuote =
  (storefrontId: string) => (state: RootState) =>
    state.delivery.selectedQuotes[storefrontId];

/**
 * Получить все выбранные quotes
 */
export const selectAllSelectedQuotes = (state: RootState) =>
  state.delivery.selectedQuotes;

/**
 * Получить calculation для request (с проверкой кэша)
 */
export const selectCalculation =
  (request: CalculationRequest) => (state: RootState) => {
    const cacheKey = getCacheKey(request);
    const cached = state.delivery.calculations[cacheKey];

    if (!cached) return null;

    // Check if cache is still valid
    if (!isCacheValid(cached, request)) {
      return null;
    }

    return cached.data;
  };

/**
 * Получить loading state для calculation
 */
export const selectCalculationLoading =
  (request: CalculationRequest) => (state: RootState) => {
    const cacheKey = getCacheKey(request);
    return state.delivery.calculationsLoading[cacheKey] || false;
  };

/**
 * Получить error для calculation
 */
export const selectCalculationError =
  (request: CalculationRequest) => (state: RootState) => {
    const cacheKey = getCacheKey(request);
    return state.delivery.calculationsError[cacheKey];
  };

/**
 * Получить tracking info
 */
export const selectTracking = (trackingToken: string) => (state: RootState) =>
  state.delivery.tracking[trackingToken];

/**
 * Получить loading state для tracking
 */
export const selectTrackingLoading =
  (trackingToken: string) => (state: RootState) =>
    state.delivery.trackingLoading[trackingToken] || false;

/**
 * Получить error для tracking
 */
export const selectTrackingError =
  (trackingToken: string) => (state: RootState) =>
    state.delivery.trackingError[trackingToken];

// ========================================
// EXPORT REDUCER
// ========================================

export default deliverySlice.reducer;
