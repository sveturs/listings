import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { UnifiedAttributeField } from '../UnifiedAttributeField';
import type { components } from '@/types/generated/api';

// Mock next-intl
jest.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
  useLocale: () => 'en',
}));

type UnifiedAttribute = components['schemas']['models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['models.UnifiedAttributeValue'];

describe('UnifiedAttributeField', () => {
  const mockOnChange = jest.fn();

  const createMockAttribute = (
    overrides?: Partial<UnifiedAttribute>
  ): UnifiedAttribute => ({
    id: 1,
    name: 'test_attribute',
    display_name: 'Test Attribute',
    attribute_type: 'text',
    is_required: false,
    is_active: true,
    ...overrides,
  });

  beforeEach(() => {
    mockOnChange.mockClear();
  });

  describe('Text field', () => {
    it('should render text input', () => {
      const attribute = createMockAttribute({ attribute_type: 'text' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const input = screen.getByRole('textbox');
      expect(input).toBeInTheDocument();
    });

    it('should handle text input change', () => {
      const attribute = createMockAttribute({ attribute_type: 'text' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const input = screen.getByRole('textbox');
      fireEvent.change(input, { target: { value: 'test value' } });

      expect(mockOnChange).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute_id: 1,
          text_value: 'test value',
        })
      );
    });

    it('should display initial value', () => {
      const attribute = createMockAttribute({ attribute_type: 'text' });
      const value: UnifiedAttributeValue = {
        attribute_id: 1,
        text_value: 'initial value',
      };

      render(
        <UnifiedAttributeField
          attribute={attribute}
          value={value}
          onChange={mockOnChange}
        />
      );

      const input = screen.getByRole('textbox') as HTMLInputElement;
      expect(input.value).toBe('initial value');
    });
  });

  describe('Number field', () => {
    it('should render number input', () => {
      const attribute = createMockAttribute({ attribute_type: 'number' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const input = screen.getByRole('spinbutton');
      expect(input).toBeInTheDocument();
    });

    it('should handle number input change', () => {
      const attribute = createMockAttribute({ attribute_type: 'number' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const input = screen.getByRole('spinbutton');
      fireEvent.change(input, { target: { value: '42' } });

      expect(mockOnChange).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute_id: 1,
          numeric_value: 42,
        })
      );
    });
  });

  describe('Select field', () => {
    it('should render select with options', () => {
      const attribute = createMockAttribute({
        attribute_type: 'select',
        options: ['Option 1', 'Option 2', 'Option 3'] as any,
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const select = screen.getByRole('combobox');
      expect(select).toBeInTheDocument();

      // Check for default option
      expect(screen.getByText('select_option')).toBeInTheDocument();
    });

    it('should handle select change', () => {
      const attribute = createMockAttribute({
        attribute_type: 'select',
        options: ['Option 1', 'Option 2'] as any,
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const select = screen.getByRole('combobox');
      fireEvent.change(select, { target: { value: 'Option 2' } });

      expect(mockOnChange).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute_id: 1,
          text_value: 'Option 2',
        })
      );
    });
  });

  describe('Boolean field', () => {
    it('should render checkbox', () => {
      const attribute = createMockAttribute({ attribute_type: 'boolean' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toBeInTheDocument();
    });

    it('should handle checkbox change', () => {
      const attribute = createMockAttribute({ attribute_type: 'boolean' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const checkbox = screen.getByRole('checkbox');
      fireEvent.click(checkbox);

      expect(mockOnChange).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute_id: 1,
          boolean_value: true,
        })
      );
    });
  });

  describe('Date field', () => {
    it('should render date input', () => {
      const attribute = createMockAttribute({ attribute_type: 'date' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const input = screen.getByDisplayValue('');
      expect(input).toHaveAttribute('type', 'date');
    });

    it('should handle date change', () => {
      const attribute = createMockAttribute({ attribute_type: 'date' });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const input = screen.getByDisplayValue('');
      fireEvent.change(input, { target: { value: '2025-09-03' } });

      expect(mockOnChange).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute_id: 1,
          date_value: '2025-09-03',
        })
      );
    });
  });

  describe('Multiselect field', () => {
    it('should render checkboxes for multiselect', () => {
      const attribute = createMockAttribute({
        attribute_type: 'multiselect',
        options: ['Option 1', 'Option 2', 'Option 3'] as any,
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const checkboxes = screen.getAllByRole('checkbox');
      expect(checkboxes).toHaveLength(3);
    });

    it('should handle multiselect changes', () => {
      const attribute = createMockAttribute({
        attribute_type: 'multiselect',
        options: ['Option 1', 'Option 2'] as any,
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const checkboxes = screen.getAllByRole('checkbox');
      fireEvent.click(checkboxes[0]);

      expect(mockOnChange).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute_id: 1,
          text_value: JSON.stringify(['Option 1']),
        })
      );
    });
  });

  describe('Validation', () => {
    it('should show required indicator', () => {
      const attribute = createMockAttribute({
        attribute_type: 'text',
        display_name: 'Required Field',
      });

      render(
        <UnifiedAttributeField
          attribute={attribute}
          required={true}
          onChange={mockOnChange}
        />
      );

      expect(screen.getByText('*')).toBeInTheDocument();
    });

    it('should show error message', () => {
      const attribute = createMockAttribute({ attribute_type: 'text' });

      render(
        <UnifiedAttributeField
          attribute={attribute}
          error="This field is required"
          onChange={mockOnChange}
        />
      );

      expect(screen.getByText('This field is required')).toBeInTheDocument();
    });

    it('should disable input when disabled prop is true', () => {
      const attribute = createMockAttribute({ attribute_type: 'text' });

      render(
        <UnifiedAttributeField
          attribute={attribute}
          disabled={true}
          onChange={mockOnChange}
        />
      );

      const input = screen.getByRole('textbox');
      expect(input).toBeDisabled();
    });
  });

  describe('Localization', () => {
    it('should use localized display name if available', () => {
      const attribute = createMockAttribute({
        attribute_type: 'text',
        display_name: 'Default Name',
        translations: {
          en: 'English Name',
          ru: 'Русское название',
        } as any,
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      expect(screen.getByText('English Name')).toBeInTheDocument();
    });

    it('should fall back to display_name if translation not available', () => {
      const attribute = createMockAttribute({
        attribute_type: 'text',
        display_name: 'Default Name',
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      expect(screen.getByText('Default Name')).toBeInTheDocument();
    });
  });

  describe('Edge cases', () => {
    it('should handle unknown attribute type', () => {
      const attribute = createMockAttribute({
        attribute_type: 'unknown_type' as any,
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      // Should render text input as fallback
      const input = screen.getByRole('textbox');
      expect(input).toBeInTheDocument();
    });

    it('should handle empty options array', () => {
      const attribute = createMockAttribute({
        attribute_type: 'select',
        options: [] as any,
      });

      render(
        <UnifiedAttributeField attribute={attribute} onChange={mockOnChange} />
      );

      const select = screen.getByRole('combobox');
      expect(select).toBeInTheDocument();
      // Should only have default option
      expect(select.children).toHaveLength(1);
    });
  });
});
