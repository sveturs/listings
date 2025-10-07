import { unifiedAttributeService } from '../unifiedAttributeService';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['models.UnifiedAttributeValue'];

// Mock fetch globally
global.fetch = jest.fn();

describe('UnifiedAttributeService', () => {
  const mockBaseUrl = 'http://localhost:3000/api';

  beforeEach(() => {
    jest.clearAllMocks();
    // Reset to v2 API
    unifiedAttributeService.useV2Api();
  });

  describe('getCategoryAttributes', () => {
    it('should fetch category attributes successfully', async () => {
      const mockAttributes: UnifiedAttribute[] = [
        { id: 1, name: 'brand', display_name: 'Brand', attribute_type: 'text' },
        { id: 2, name: 'model', display_name: 'Model', attribute_type: 'text' },
      ];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: mockAttributes }),
      });

      const result = await unifiedAttributeService.getCategoryAttributes(123);

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/attributes/category/123`,
        expect.objectContaining({
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        })
      );

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockAttributes);
    });

    it('should handle fetch error', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      const result = await unifiedAttributeService.getCategoryAttributes(123);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Failed to fetch category attributes');
    });

    it('should handle API error response', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ success: false, error: 'Category not found' }),
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

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: mockAttributes }),
      });

      const result = await unifiedAttributeService.getAttributes({
        category_id: 123,
        required: true,
        searchable: false,
      });

      const expectedUrl = `${mockBaseUrl}/v2/attributes?category_id=123&required=true&searchable=false`;
      expect(global.fetch).toHaveBeenCalledWith(expectedUrl, expect.anything());

      expect(result.success).toBe(true);
    });

    it('should handle empty params', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: [] }),
      });

      await unifiedAttributeService.getAttributes();

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/attributes`,
        expect.anything()
      );
    });
  });

  describe('createAttributeValue', () => {
    it('should create attribute value for listing', async () => {
      const mockValue: UnifiedAttributeValue = {
        id: 1,
        attribute_id: 10,
        text_value: 'test value',
      };

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: mockValue }),
      });

      const result = await unifiedAttributeService.createAttributeValue({
        listing_id: 100,
        attribute_id: 10,
        text_value: 'test value',
      });

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/listings/100/attributes`,
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            attribute_id: 10,
            text_value: 'test value',
          }),
        })
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

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: mockValue }),
      });

      const result = await unifiedAttributeService.createAttributeValue({
        product_id: 200,
        attribute_id: 10,
        numeric_value: 42,
      });

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/products/200/attributes`,
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

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: mockValue }),
      });

      const result = await unifiedAttributeService.updateAttributeValue({
        id: 1,
        listing_id: 100,
        attribute_id: 10,
        text_value: 'updated value',
      });

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/listings/100/attributes/1`,
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({
            text_value: 'updated value',
          }),
        })
      );

      expect(result.success).toBe(true);
    });
  });

  describe('deleteAttributeValue', () => {
    it('should delete attribute value', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, message: 'Deleted' }),
      });

      const result = await unifiedAttributeService.deleteAttributeValue(1, 100);

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/listings/100/attributes/1`,
        expect.objectContaining({
          method: 'DELETE',
        })
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

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: attributeValues }),
      });

      const result = await unifiedAttributeService.saveListingAttributes(
        100,
        attributeValues
      );

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/listings/100/attributes/batch`,
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ attributes: attributeValues }),
        })
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

  describe('API version switching', () => {
    it('should switch to v1 API', async () => {
      unifiedAttributeService.useV1Api();

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true, data: [] }),
      });

      await unifiedAttributeService.getAttributes();

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v1/attributes`,
        expect.anything()
      );
    });

    it('should check v2 API availability', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true }),
      });

      const isAvailable =
        await unifiedAttributeService.checkV2ApiAvailability();

      expect(global.fetch).toHaveBeenCalledWith(
        `${mockBaseUrl}/v2/attributes/health`,
        expect.anything()
      );

      expect(isAvailable).toBe(true);
    });

    it('should return false when v2 API is not available', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ success: false }),
      });

      const isAvailable =
        await unifiedAttributeService.checkV2ApiAvailability();
      expect(isAvailable).toBe(false);
    });
  });
});
