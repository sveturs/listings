import { unifiedAttributeService } from '../unifiedAttributeService';
import { apiClient } from '../api-client';
import type { components } from '@/types/generated/api';

type UnifiedAttribute = components['schemas']['models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['models.UnifiedAttributeValue'];

// Mock apiClient
jest.mock('../api-client', () => ({
  apiClient: {
    get: jest.fn(),
    post: jest.fn(),
    put: jest.fn(),
    delete: jest.fn(),
  },
}));

describe('UnifiedAttributeService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getCategoryAttributes', () => {
    it('should fetch category attributes successfully', async () => {
      const mockAttributes: UnifiedAttribute[] = [
        { id: 1, name: 'brand', display_name: 'Brand', attribute_type: 'text' },
        { id: 2, name: 'model', display_name: 'Model', attribute_type: 'text' },
      ];

      (apiClient.get as jest.Mock).mockResolvedValueOnce({
        data: { data: mockAttributes },
        status: 200,
      });

      const result = await unifiedAttributeService.getCategoryAttributes(123);

      expect(apiClient.get).toHaveBeenCalledWith('/v2/attributes/category/123');

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockAttributes);
    });

    it('should handle fetch error', async () => {
      (apiClient.get as jest.Mock).mockResolvedValueOnce({
        error: { message: 'Network error' },
        status: 500,
      });

      const result = await unifiedAttributeService.getCategoryAttributes(123);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Network error');
    });

    it('should handle API error response', async () => {
      (apiClient.get as jest.Mock).mockResolvedValueOnce({
        error: { message: 'Category not found' },
        status: 404,
      });

      const result = await unifiedAttributeService.getCategoryAttributes(123);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Category not found');
    });
  });

  describe('getAttributes', () => {
    it('should fetch attributes with filters', async () => {
      const mockAttributes: UnifiedAttribute[] = [
        { id: 1, name: 'test', is_required: true },
      ];

      (apiClient.get as jest.Mock).mockResolvedValueOnce({
        data: { data: mockAttributes },
        status: 200,
      });

      const result = await unifiedAttributeService.getAttributes({
        category_id: 123,
        required: true,
        searchable: false,
      });

      expect(apiClient.get).toHaveBeenCalledWith(
        '/v2/attributes?category_id=123&required=true&searchable=false'
      );

      expect(result.success).toBe(true);
    });

    it('should handle empty params', async () => {
      (apiClient.get as jest.Mock).mockResolvedValueOnce({
        data: { data: [] },
        status: 200,
      });

      await unifiedAttributeService.getAttributes();

      expect(apiClient.get).toHaveBeenCalledWith('/v2/attributes');
    });
  });

  describe('createAttributeValue', () => {
    it('should create attribute value for listing', async () => {
      const mockValue: UnifiedAttributeValue = {
        id: 1,
        attribute_id: 10,
        text_value: 'test value',
      };

      (apiClient.post as jest.Mock).mockResolvedValueOnce({
        data: { data: mockValue },
        status: 201,
      });

      const result = await unifiedAttributeService.createAttributeValue({
        listing_id: 100,
        attribute_id: 10,
        text_value: 'test value',
      });

      expect(apiClient.post).toHaveBeenCalledWith(
        '/v2/listings/100/attributes',
        {
          attribute_id: 10,
          text_value: 'test value',
          numeric_value: undefined,
          boolean_value: undefined,
          date_value: undefined,
          json_value: undefined,
        }
      );

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockValue);
    });

    it('should create attribute value for product', async () => {
      const mockValue: UnifiedAttributeValue = {
        id: 1,
        attribute_id: 10,
        numeric_value: 42,
      };

      (apiClient.post as jest.Mock).mockResolvedValueOnce({
        data: { data: mockValue },
        status: 201,
      });

      const result = await unifiedAttributeService.createAttributeValue({
        product_id: 200,
        attribute_id: 10,
        numeric_value: 42,
      });

      expect(apiClient.post).toHaveBeenCalledWith(
        '/v2/products/200/attributes',
        expect.anything()
      );

      expect(result.success).toBe(true);
    });
  });

  describe('updateAttributeValue', () => {
    it('should update attribute value', async () => {
      const mockValue: UnifiedAttributeValue = {
        id: 1,
        attribute_id: 10,
        text_value: 'updated value',
      };

      (apiClient.put as jest.Mock).mockResolvedValueOnce({
        data: { data: mockValue },
        status: 200,
      });

      const result = await unifiedAttributeService.updateAttributeValue({
        id: 1,
        listing_id: 100,
        attribute_id: 10,
        text_value: 'updated value',
      });

      expect(apiClient.put).toHaveBeenCalledWith(
        '/v2/listings/100/attributes/1',
        {
          text_value: 'updated value',
          numeric_value: undefined,
          boolean_value: undefined,
          date_value: undefined,
          json_value: undefined,
        }
      );

      expect(result.success).toBe(true);
    });
  });

  describe('deleteAttributeValue', () => {
    it('should delete attribute value', async () => {
      (apiClient.delete as jest.Mock).mockResolvedValueOnce({
        data: { message: 'Deleted' },
        status: 200,
      });

      const result = await unifiedAttributeService.deleteAttributeValue(1, 100);

      expect(apiClient.delete).toHaveBeenCalledWith(
        '/v2/listings/100/attributes/1'
      );

      expect(result.success).toBe(true);
    });
  });

  describe('saveListingAttributes', () => {
    it('should batch save listing attributes', async () => {
      const attributeValues: Partial<UnifiedAttributeValue>[] = [
        { attribute_id: 1, text_value: 'value1' },
        { attribute_id: 2, numeric_value: 42 },
      ];

      (apiClient.post as jest.Mock).mockResolvedValueOnce({
        data: { data: attributeValues },
        status: 200,
      });

      const result = await unifiedAttributeService.saveListingAttributes(
        100,
        attributeValues
      );

      expect(apiClient.post).toHaveBeenCalledWith(
        '/v2/listings/100/attributes/batch',
        { attributes: attributeValues }
      );

      expect(result.success).toBe(true);
    });
  });

  describe('validateAttributeValue', () => {
    it('should validate required field', () => {
      const attribute: UnifiedAttribute = {
        id: 1,
        name: 'test',
        is_required: true,
        attribute_type: 'text',
      };

      const value: UnifiedAttributeValue = {
        attribute_id: 1,
      };

      const errors = unifiedAttributeService.validateAttributeValue(
        attribute,
        value
      );
      expect(errors).toContain('This field is required');
    });

    it('should pass validation for filled required field', () => {
      const attribute: UnifiedAttribute = {
        id: 1,
        name: 'test',
        is_required: true,
        attribute_type: 'text',
      };

      const value: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: 'some value',
      };

      const errors = unifiedAttributeService.validateAttributeValue(
        attribute,
        value
      );
      expect(errors).toHaveLength(0);
    });

    it('should validate email format', () => {
      const attribute: UnifiedAttribute = {
        id: 1,
        name: 'email',
        attribute_type: 'email',
      };

      const invalidValue: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: 'invalid-email',
      };

      const errors = unifiedAttributeService.validateAttributeValue(
        attribute,
        invalidValue
      );
      expect(errors).toContain('Invalid email address');

      const validValue: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: 'test@example.com',
      };

      const noErrors = unifiedAttributeService.validateAttributeValue(
        attribute,
        validValue
      );
      expect(noErrors).toHaveLength(0);
    });

    it('should validate URL format', () => {
      const attribute: UnifiedAttribute = {
        id: 1,
        name: 'website',
        attribute_type: 'url',
      };

      const invalidValue: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: 'not a url',
      };

      const errors = unifiedAttributeService.validateAttributeValue(
        attribute,
        invalidValue
      );
      expect(errors).toContain('Invalid URL');

      const validValue: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: 'https://example.com',
      };

      const noErrors = unifiedAttributeService.validateAttributeValue(
        attribute,
        validValue
      );
      expect(noErrors).toHaveLength(0);
    });

    it('should validate phone number', () => {
      const attribute: UnifiedAttribute = {
        id: 1,
        name: 'phone',
        attribute_type: 'tel',
      };

      const invalidValue: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: 'abc123',
      };

      const errors = unifiedAttributeService.validateAttributeValue(
        attribute,
        invalidValue
      );
      expect(errors).toContain('Invalid phone number');

      const validValue: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: '+1234567890',
      };

      const noErrors = unifiedAttributeService.validateAttributeValue(
        attribute,
        validValue
      );
      expect(noErrors).toHaveLength(0);
    });
  });

  describe('checkV2ApiAvailability', () => {
    it('should check v2 API availability', async () => {
      (apiClient.get as jest.Mock).mockResolvedValueOnce({
        data: { success: true },
        status: 200,
      });

      const isAvailable =
        await unifiedAttributeService.checkV2ApiAvailability();

      expect(apiClient.get).toHaveBeenCalledWith('/v2/attributes/health');

      expect(isAvailable).toBe(true);
    });

    it('should return false when v2 API is not available', async () => {
      (apiClient.get as jest.Mock).mockResolvedValueOnce({
        error: { message: 'Not available' },
        status: 503,
      });

      const isAvailable =
        await unifiedAttributeService.checkV2ApiAvailability();
      expect(isAvailable).toBe(false);
    });
  });
});
