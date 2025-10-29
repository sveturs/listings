/**
 * Unit Tests for Delivery Service
 *
 * Tests all deliveryService methods with mocked apiClient
 */

import { deliveryService } from './delivery';
import { apiClient } from './api-client';
import type {
  CalculateRateRequest,
  CalculateRateResponse,
  CreateShipmentRequest,
  CreateShipmentResponse,
  TrackingInfo,
  DeliveryProvider,
  Settlement,
  Street,
  ParcelLocker,
} from './delivery';

// Mock apiClient
jest.mock('./api-client', () => ({
  apiClient: {
    get: jest.fn(),
    post: jest.fn(),
    put: jest.fn(),
    delete: jest.fn(),
    patch: jest.fn(),
  },
}));

describe('deliveryService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // ========================================
  // calculateRate
  // ========================================

  describe('calculateRate', () => {
    it('should calculate delivery rate successfully', async () => {
      // Arrange
      const request: CalculateRateRequest = {
        provider: 'post_express',
        from_city: 'Belgrade',
        to_city: 'Novi Sad',
        weight: 2.5,
        length: 30,
        width: 20,
        height: 10,
      };

      const mockResponse: CalculateRateResponse = {
        base_price: 200,
        insurance: 0,
        cod_fee: 0,
        weight_fee: 50,
        distance_fee: 0,
        total_cost: 250,
        estimated_delivery: '2025-10-31T12:00:00Z',
        currency: 'RSD',
      };

      (apiClient.post as jest.Mock).mockResolvedValue({
        data: mockResponse,
        error: null,
      });

      // Act
      const result = await deliveryService.calculateRate(request);

      // Assert
      expect(apiClient.post).toHaveBeenCalledWith('/delivery/calculate-rate', request);
      expect(result.data).toEqual(mockResponse);
      expect(result.data?.total_cost).toBe(250);
      expect(result.error).toBeNull();
    });

    it('should handle COD and insurance fees', async () => {
      // Arrange
      const request: CalculateRateRequest = {
        provider: 'post_express',
        from_city: 'Belgrade',
        to_city: 'Nis',
        weight: 1.5,
        cash_on_delivery: true,
        cod_amount: 5000,
        insurance: true,
      };

      const mockResponse: CalculateRateResponse = {
        base_price: 200,
        insurance: 50,
        cod_fee: 100,
        weight_fee: 0,
        distance_fee: 100,
        total_cost: 450,
        estimated_delivery: '2025-11-01T12:00:00Z',
        currency: 'RSD',
      };

      (apiClient.post as jest.Mock).mockResolvedValue({
        data: mockResponse,
      });

      // Act
      const result = await deliveryService.calculateRate(request);

      // Assert
      expect(result.data?.cod_fee).toBe(100);
      expect(result.data?.insurance).toBe(50);
      expect(result.data?.total_cost).toBe(450);
    });

    it('should handle API errors', async () => {
      // Arrange
      const request: CalculateRateRequest = {
        provider: 'post_express',
        from_city: 'Belgrade',
        to_city: 'Novi Sad',
        weight: 2.5,
      };

      (apiClient.post as jest.Mock).mockResolvedValue({
        data: null,
        error: { message: 'Invalid provider', statusCode: 400 },
      });

      // Act
      const result = await deliveryService.calculateRate(request);

      // Assert
      expect(result.error).toBeDefined();
      expect(result.error?.message).toBe('Invalid provider');
      expect(result.data).toBeNull();
    });
  });

  // ========================================
  // calculateCart
  // ========================================

  describe('calculateCart', () => {
    it('should calculate cart delivery cost', async () => {
      // Arrange
      const storefrontId = 123;
      const mockResponse: CalculateRateResponse = {
        base_price: 150,
        insurance: 0,
        cod_fee: 0,
        weight_fee: 30,
        distance_fee: 20,
        total_cost: 200,
        estimated_delivery: '2025-10-31T12:00:00Z',
        currency: 'RSD',
      };

      (apiClient.post as jest.Mock).mockResolvedValue({
        data: mockResponse,
      });

      // Act
      const result = await deliveryService.calculateCart(storefrontId);

      // Assert
      expect(apiClient.post).toHaveBeenCalledWith(`/delivery/calculate-cart/${storefrontId}`);
      expect(result.data?.total_cost).toBe(200);
    });
  });

  // ========================================
  // getProviders
  // ========================================

  describe('getProviders', () => {
    it('should fetch all delivery providers', async () => {
      // Arrange
      const mockProviders: DeliveryProvider[] = [
        {
          code: 'post_express',
          name: 'Post Express',
          description: 'Official Serbian Post',
          enabled: true,
          logo_url: '/logos/post_express.png',
        },
        {
          code: 'bex',
          name: 'BEX Express',
          description: 'Fast delivery',
          enabled: true,
        },
      ];

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockProviders,
      });

      // Act
      const result = await deliveryService.getProviders();

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith('/delivery/providers');
      expect(result.data).toHaveLength(2);
      expect(result.data?.[0].code).toBe('post_express');
      expect(result.data?.[1].code).toBe('bex');
    });

    it('should handle empty provider list', async () => {
      // Arrange
      (apiClient.get as jest.Mock).mockResolvedValue({
        data: [],
      });

      // Act
      const result = await deliveryService.getProviders();

      // Assert
      expect(result.data).toEqual([]);
    });
  });

  // ========================================
  // createShipment
  // ========================================

  describe('createShipment', () => {
    it('should create shipment successfully', async () => {
      // Arrange
      const request: CreateShipmentRequest = {
        order_id: 123,
        provider_code: 'post_express',
        from_address: {
          name: 'Store Name',
          phone: '+381601234567',
          street: 'Main Street 1',
          city: 'Belgrade',
          postalCode: '11000',
          country: 'RS',
        },
        to_address: {
          name: 'Customer Name',
          phone: '+381607654321',
          street: 'Customer Street 5',
          city: 'Novi Sad',
          postalCode: '21000',
          country: 'RS',
        },
        packages: [
          {
            weight: 2.5,
            length: 30,
            width: 20,
            height: 10,
            description: 'Order #123',
            declaredValue: 5000,
          },
        ],
      };

      const mockResponse: CreateShipmentResponse = {
        shipment_id: 456,
        tracking_number: 'PE123456789',
        provider_code: 'post_express',
        status: 'pending',
        label_url: 'https://example.com/label.pdf',
      };

      (apiClient.post as jest.Mock).mockResolvedValue({
        data: mockResponse,
      });

      // Act
      const result = await deliveryService.createShipment(request);

      // Assert
      expect(apiClient.post).toHaveBeenCalledWith('/delivery/shipments', request);
      expect(result.data?.tracking_number).toBe('PE123456789');
      expect(result.data?.provider_code).toBe('post_express');
      expect(result.data?.label_url).toBeDefined();
    });

    it('should handle COD shipment', async () => {
      // Arrange
      const request: CreateShipmentRequest = {
        order_id: 123,
        provider_code: 'post_express',
        from_address: {
          street: 'Main St 1',
          city: 'Belgrade',
          postalCode: '11000',
          country: 'RS',
        },
        to_address: {
          street: 'Customer St 5',
          city: 'Novi Sad',
          postalCode: '21000',
          country: 'RS',
        },
        packages: [
          {
            weight: 1.5,
            cashOnDelivery: true,
            codAmount: 5000,
          },
        ],
      };

      const mockResponse: CreateShipmentResponse = {
        shipment_id: 789,
        tracking_number: 'PE987654321',
        provider_code: 'post_express',
        status: 'pending',
      };

      (apiClient.post as jest.Mock).mockResolvedValue({
        data: mockResponse,
      });

      // Act
      const result = await deliveryService.createShipment(request);

      // Assert
      expect(result.data?.tracking_number).toBe('PE987654321');
    });
  });

  // ========================================
  // trackShipment
  // ========================================

  describe('trackShipment', () => {
    it('should track shipment successfully', async () => {
      // Arrange
      const trackingToken = 'PE123456789';
      const mockResponse: TrackingInfo = {
        shipment_id: 456,
        tracking_number: trackingToken,
        status: 'in_transit',
        current_location: 'Belgrade Sorting Center',
        estimated_delivery: '2025-10-31T12:00:00Z',
        events: [
          {
            timestamp: '2025-10-29T10:00:00Z',
            location: 'Belgrade',
            status: 'accepted',
            description: 'Shipment accepted',
          },
          {
            timestamp: '2025-10-29T15:00:00Z',
            location: 'Belgrade Sorting Center',
            status: 'in_transit',
            description: 'In transit to Novi Sad',
          },
        ],
      };

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockResponse,
      });

      // Act
      const result = await deliveryService.trackShipment(trackingToken);

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith(`/delivery/track/${trackingToken}`);
      expect(result.data?.status).toBe('in_transit');
      expect(result.data?.events).toHaveLength(2);
      expect(result.data?.current_location).toBe('Belgrade Sorting Center');
    });

    it('should handle delivered shipment', async () => {
      // Arrange
      const trackingToken = 'PE123456789';
      const mockResponse: TrackingInfo = {
        shipment_id: 456,
        tracking_number: trackingToken,
        status: 'delivered',
        current_location: 'Novi Sad',
        actual_delivery: '2025-10-30T14:30:00Z',
        events: [
          {
            timestamp: '2025-10-30T14:30:00Z',
            location: 'Novi Sad',
            status: 'delivered',
            description: 'Package delivered to recipient',
          },
        ],
      };

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockResponse,
      });

      // Act
      const result = await deliveryService.trackShipment(trackingToken);

      // Assert
      expect(result.data?.status).toBe('delivered');
      expect(result.data?.actual_delivery).toBeDefined();
    });

    it('should handle invalid tracking number', async () => {
      // Arrange
      const trackingToken = 'INVALID123';

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: null,
        error: { message: 'Tracking number not found', statusCode: 404 },
      });

      // Act
      const result = await deliveryService.trackShipment(trackingToken);

      // Assert
      expect(result.error).toBeDefined();
      expect(result.error?.statusCode).toBe(404);
    });
  });

  // ========================================
  // cancelShipment
  // ========================================

  describe('cancelShipment', () => {
    it('should cancel shipment successfully', async () => {
      // Arrange
      const shipmentId = 456;

      (apiClient.delete as jest.Mock).mockResolvedValue({
        data: null,
        error: null,
      });

      // Act
      const result = await deliveryService.cancelShipment(shipmentId);

      // Assert
      expect(apiClient.delete).toHaveBeenCalledWith(`/delivery/shipments/${shipmentId}`);
      expect(result.error).toBeNull();
    });

    it('should handle cancellation error', async () => {
      // Arrange
      const shipmentId = 456;

      (apiClient.delete as jest.Mock).mockResolvedValue({
        data: null,
        error: { message: 'Cannot cancel shipped package', statusCode: 400 },
      });

      // Act
      const result = await deliveryService.cancelShipment(shipmentId);

      // Assert
      expect(result.error).toBeDefined();
      expect(result.error?.message).toContain('cancel');
    });
  });

  // ========================================
  // POST EXPRESS SPECIFIC
  // ========================================

  describe('getSettlements', () => {
    it('should fetch all settlements', async () => {
      // Arrange
      const mockSettlements: Settlement[] = [
        {
          id: '1',
          name: 'Beograd',
          postal_code: '11000',
          region: 'Central Serbia',
        },
        {
          id: '2',
          name: 'Novi Sad',
          postal_code: '21000',
          region: 'Vojvodina',
        },
      ];

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockSettlements,
      });

      // Act
      const result = await deliveryService.getSettlements();

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith('/delivery/settlements');
      expect(result.data).toHaveLength(2);
    });

    it('should search settlements by term', async () => {
      // Arrange
      const searchTerm = 'Beograd';
      const mockSettlements: Settlement[] = [
        {
          id: '1',
          name: 'Beograd',
          postal_code: '11000',
        },
      ];

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockSettlements,
      });

      // Act
      const result = await deliveryService.getSettlements(searchTerm);

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith('/delivery/settlements?search=Beograd');
      expect(result.data).toHaveLength(1);
      expect(result.data?.[0].name).toBe('Beograd');
    });
  });

  describe('getStreets', () => {
    it('should fetch streets for settlement', async () => {
      // Arrange
      const settlementId = '1';
      const mockStreets: Street[] = [
        {
          id: '101',
          name: 'Knez Mihailova',
          settlement_id: settlementId,
        },
        {
          id: '102',
          name: 'Terazije',
          settlement_id: settlementId,
        },
      ];

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockStreets,
      });

      // Act
      const result = await deliveryService.getStreets(settlementId);

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith(`/delivery/settlements/${settlementId}/streets`);
      expect(result.data).toHaveLength(2);
    });

    it('should search streets by term', async () => {
      // Arrange
      const settlementId = '1';
      const searchTerm = 'Knez';
      const mockStreets: Street[] = [
        {
          id: '101',
          name: 'Knez Mihailova',
          settlement_id: settlementId,
        },
      ];

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockStreets,
      });

      // Act
      const result = await deliveryService.getStreets(settlementId, searchTerm);

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith(
        `/delivery/settlements/${settlementId}/streets?search=Knez`
      );
      expect(result.data?.[0].name).toContain('Knez');
    });
  });

  describe('getParcelLockers', () => {
    it('should fetch all parcel lockers', async () => {
      // Arrange
      const mockLockers: ParcelLocker[] = [
        {
          id: 'PL001',
          name: 'Belgrade Center',
          address: 'Knez Mihailova 1',
          city: 'Belgrade',
          latitude: 44.817,
          longitude: 20.457,
          working_hours: '24/7',
          available: true,
        },
        {
          id: 'PL002',
          name: 'Novi Sad Center',
          address: 'Zmaj Jovina 1',
          city: 'Novi Sad',
          available: true,
        },
      ];

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockLockers,
      });

      // Act
      const result = await deliveryService.getParcelLockers();

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith('/delivery/parcel-lockers');
      expect(result.data).toHaveLength(2);
      expect(result.data?.[0].working_hours).toBe('24/7');
    });

    it('should filter parcel lockers by city', async () => {
      // Arrange
      const cityId = '1';
      const mockLockers: ParcelLocker[] = [
        {
          id: 'PL001',
          name: 'Belgrade Center',
          address: 'Knez Mihailova 1',
          city: 'Belgrade',
          available: true,
        },
      ];

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockLockers,
      });

      // Act
      const result = await deliveryService.getParcelLockers(cityId);

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith(`/delivery/parcel-lockers?city_id=${cityId}`);
      expect(result.data).toHaveLength(1);
    });
  });

  // ========================================
  // Product & Category Defaults
  // ========================================

  describe('getProductAttributes', () => {
    it('should fetch product delivery attributes', async () => {
      // Arrange
      const productId = '123';
      const type = 'c2c';
      const mockAttributes = {
        weight: 2.5,
        dimensions: { length: 30, width: 20, height: 10 },
        fragile: false,
      };

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockAttributes,
      });

      // Act
      const result = await deliveryService.getProductAttributes(productId, type);

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith(
        `/products/${productId}/delivery-attributes?type=${type}`
      );
      expect(result.data?.weight).toBe(2.5);
    });
  });

  describe('getCategoryDefaults', () => {
    it('should fetch category delivery defaults', async () => {
      // Arrange
      const categoryId = '1301';
      const mockDefaults = {
        default_weight: 1.0,
        requires_signature: false,
      };

      (apiClient.get as jest.Mock).mockResolvedValue({
        data: mockDefaults,
      });

      // Act
      const result = await deliveryService.getCategoryDefaults(categoryId);

      // Assert
      expect(apiClient.get).toHaveBeenCalledWith(`/categories/${categoryId}/delivery-defaults`);
      expect(result.data?.default_weight).toBe(1.0);
    });
  });
});
