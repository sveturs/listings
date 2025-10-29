/**
 * Unit Tests for Delivery Redux Slice
 *
 * Tests reducers and selectors (simplified without redux-mock-store)
 */

import {
  deliverySlice,
  selectQuote,
  clearQuote,
  clearAllQuotes,
  clearCalculationsCache,
  clearTrackingCache,
  isCacheValid,
} from './deliverySlice';
import type { DeliveryQuote, CalculationRequest, CalculationResponse } from '@/types/delivery';

describe('deliverySlice', () => {
  // ========================================
  // INITIAL STATE
  // ========================================

  describe('initial state', () => {
    it('should have correct initial state', () => {
      const state = deliverySlice.reducer(undefined, { type: 'unknown' });

      expect(state.providers).toEqual([]);
      expect(state.providersLoading).toBe(false);
      expect(state.providersError).toBeNull();
      expect(state.selectedQuotes).toEqual({});
      expect(state.calculations).toEqual({});
      expect(state.tracking).toEqual({});
    });
  });

  // ========================================
  // REDUCERS: selectQuote, clearQuote
  // ========================================

  describe('selectQuote reducer', () => {
    it('should select quote for storefront', () => {
      // Arrange
      const initialState = deliverySlice.reducer(undefined, { type: 'unknown' });
      const quote: DeliveryQuote = {
        provider_code: 'post_express',
        provider_name: 'Post Express',
        total_cost: 250,
        base_price: 200,
        estimated_delivery: '2025-10-31T12:00:00Z',
        currency: 'RSD',
      };

      // Act
      const state = deliverySlice.reducer(
        initialState,
        selectQuote({ storefrontId: '123', quote })
      );

      // Assert
      expect(state.selectedQuotes['123']).toEqual(quote);
    });

    it('should overwrite existing quote', () => {
      // Arrange
      const initialState = {
        ...deliverySlice.reducer(undefined, { type: 'unknown' }),
        selectedQuotes: {
          '123': {
            provider_code: 'old_provider',
            provider_name: 'Old Provider',
            total_cost: 100,
            base_price: 100,
            estimated_delivery: '2025-10-30T12:00:00Z',
            currency: 'RSD',
          },
        },
      };

      const newQuote: DeliveryQuote = {
        provider_code: 'post_express',
        provider_name: 'Post Express',
        total_cost: 250,
        base_price: 200,
        estimated_delivery: '2025-10-31T12:00:00Z',
        currency: 'RSD',
      };

      // Act
      const state = deliverySlice.reducer(
        initialState,
        selectQuote({ storefrontId: '123', quote: newQuote })
      );

      // Assert
      expect(state.selectedQuotes['123']).toEqual(newQuote);
    });
  });

  describe('clearQuote reducer', () => {
    it('should clear quote for specific storefront', () => {
      // Arrange
      const initialState = {
        ...deliverySlice.reducer(undefined, { type: 'unknown' }),
        selectedQuotes: {
          '123': {
            provider_code: 'post_express',
            provider_name: 'Post Express',
            total_cost: 250,
            base_price: 200,
            estimated_delivery: '2025-10-31T12:00:00Z',
            currency: 'RSD',
          },
          '456': {
            provider_code: 'bex',
            provider_name: 'BEX',
            total_cost: 300,
            base_price: 250,
            estimated_delivery: '2025-10-31T12:00:00Z',
            currency: 'RSD',
          },
        },
      };

      // Act
      const state = deliverySlice.reducer(initialState, clearQuote('123'));

      // Assert
      expect(state.selectedQuotes['123']).toBeUndefined();
      expect(state.selectedQuotes['456']).toBeDefined();
    });
  });

  describe('clearAllQuotes reducer', () => {
    it('should clear all quotes', () => {
      // Arrange
      const initialState = {
        ...deliverySlice.reducer(undefined, { type: 'unknown' }),
        selectedQuotes: {
          '123': {
            provider_code: 'post_express',
            provider_name: 'Post Express',
            total_cost: 250,
            base_price: 200,
            estimated_delivery: '2025-10-31T12:00:00Z',
            currency: 'RSD',
          },
          '456': {
            provider_code: 'bex',
            provider_name: 'BEX',
            total_cost: 300,
            base_price: 250,
            estimated_delivery: '2025-10-31T12:00:00Z',
            currency: 'RSD',
          },
        },
      };

      // Act
      const state = deliverySlice.reducer(initialState, clearAllQuotes());

      // Assert
      expect(state.selectedQuotes).toEqual({});
    });
  });

  // ========================================
  // REDUCERS: clearCalculationsCache, clearTrackingCache
  // ========================================

  describe('clearCalculationsCache reducer', () => {
    it('should clear calculations cache', () => {
      // Arrange
      const initialState = {
        ...deliverySlice.reducer(undefined, { type: 'unknown' }),
        calculations: {
          'cache-key-1': {
            data: { success: true } as CalculationResponse,
            timestamp: Date.now(),
            params: {} as CalculationRequest,
          },
        },
        calculationsError: { 'cache-key-1': 'some error' },
      };

      // Act
      const state = deliverySlice.reducer(initialState, clearCalculationsCache());

      // Assert
      expect(state.calculations).toEqual({});
      expect(state.calculationsError).toEqual({});
    });
  });

  describe('clearTrackingCache reducer', () => {
    it('should clear tracking cache', () => {
      // Arrange
      const initialState = {
        ...deliverySlice.reducer(undefined, { type: 'unknown' }),
        tracking: {
          PE123: {
            shipment_id: 1,
            tracking_number: 'PE123',
            status: 'delivered',
            events: [],
          },
        },
        trackingError: { PE123: 'some error' },
      };

      // Act
      const state = deliverySlice.reducer(initialState, clearTrackingCache());

      // Assert
      expect(state.tracking).toEqual({});
      expect(state.trackingError).toEqual({});
    });
  });

  // ========================================
  // CACHE VALIDATION
  // ========================================

  describe('isCacheValid helper', () => {
    it('should return true for valid cache', () => {
      // Arrange
      const request: CalculationRequest = {
        from_location: { city: 'Belgrade', postal_code: '11000', country: 'RS' },
        to_location: { city: 'Novi Sad', postal_code: '21000', country: 'RS' },
        items: [{ weight: 2.5 }],
      };

      const cached = {
        data: { success: true } as CalculationResponse,
        timestamp: Date.now(),
        params: request,
      };

      // Act
      const result = isCacheValid(cached, request);

      // Assert
      expect(result).toBe(true);
    });

    it('should return false for expired cache', () => {
      // Arrange
      const request: CalculationRequest = {
        from_location: { city: 'Belgrade', postal_code: '11000', country: 'RS' },
        to_location: { city: 'Novi Sad', postal_code: '21000', country: 'RS' },
        items: [{ weight: 2.5 }],
      };

      const cached = {
        data: { success: true } as CalculationResponse,
        timestamp: Date.now() - 10 * 60 * 1000, // 10 minutes ago (TTL = 5 minutes)
        params: request,
      };

      // Act
      const result = isCacheValid(cached, request);

      // Assert
      expect(result).toBe(false);
    });

    it('should return false for different params', () => {
      // Arrange
      const cachedRequest: CalculationRequest = {
        from_location: { city: 'Belgrade', postal_code: '11000', country: 'RS' },
        to_location: { city: 'Novi Sad', postal_code: '21000', country: 'RS' },
        items: [{ weight: 2.5 }],
      };

      const newRequest: CalculationRequest = {
        from_location: { city: 'Belgrade', postal_code: '11000', country: 'RS' },
        to_location: { city: 'Nis', postal_code: '18000', country: 'RS' }, // Different city
        items: [{ weight: 2.5 }],
      };

      const cached = {
        data: { success: true } as CalculationResponse,
        timestamp: Date.now(),
        params: cachedRequest,
      };

      // Act
      const result = isCacheValid(cached, newRequest);

      // Assert
      expect(result).toBe(false);
    });
  });

  // ========================================
  // STATE SHAPE VALIDATION
  // ========================================

  describe('state shape validation', () => {
    it('should maintain correct state shape after multiple operations', () => {
      // Arrange
      let state = deliverySlice.reducer(undefined, { type: 'unknown' });

      // Act - perform multiple operations
      const quote: DeliveryQuote = {
        provider_code: 'post_express',
        provider_name: 'Post Express',
        total_cost: 250,
        base_price: 200,
        estimated_delivery: '2025-10-31T12:00:00Z',
        currency: 'RSD',
      };

      state = deliverySlice.reducer(state, selectQuote({ storefrontId: '123', quote }));
      state = deliverySlice.reducer(state, selectQuote({ storefrontId: '456', quote }));
      state = deliverySlice.reducer(state, clearQuote('123'));

      // Assert
      expect(Object.keys(state.selectedQuotes)).toHaveLength(1);
      expect(state.selectedQuotes['456']).toBeDefined();
      expect(state.selectedQuotes['123']).toBeUndefined();
    });

    it('should handle clearing empty caches gracefully', () => {
      // Arrange
      const initialState = deliverySlice.reducer(undefined, { type: 'unknown' });

      // Act
      const state = deliverySlice.reducer(initialState, clearCalculationsCache());

      // Assert
      expect(state.calculations).toEqual({});
      expect(state.calculationsError).toEqual({});
    });
  });
});
