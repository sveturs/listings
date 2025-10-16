import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import RangeAttribute from '../RangeAttribute';

const mockAttribute = {
  id: 1,
  name: 'price',
  display_name: 'Price',
  unit: '$',
  min_value: 0,
  max_value: 1000,
  translations: {
    en: 'Price',
    ru: 'Цена',
    sr: 'Цена',
  },
  is_required: true,
};

describe('RangeAttribute', () => {
  const mockOnChange = jest.fn();

  beforeEach(() => {
    mockOnChange.mockClear();
  });

  it('renders with correct label and required indicator', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.getByText('Price')).toBeInTheDocument();
    expect(screen.getByText('*')).toBeInTheDocument();
  });

  it('displays correct translations based on locale', () => {
    const { rerender } = render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="ru"
      />
    );

    expect(screen.getByText('Цена')).toBeInTheDocument();

    rerender(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.getByText('Price')).toBeInTheDocument();
  });

  it('renders min and max input fields', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min');
    const maxInput = screen.getByPlaceholderText('create.max');

    expect(minInput).toBeInTheDocument();
    expect(maxInput).toBeInTheDocument();
    expect(minInput).toHaveAttribute('type', 'number');
    expect(maxInput).toHaveAttribute('type', 'number');
  });

  it('displays unit when provided', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.getByText('$')).toBeInTheDocument();
  });

  it('calls onChange when min value is entered', async () => {
    const user = userEvent.setup();
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min');
    await user.type(minInput, '100');

    expect(mockOnChange).toHaveBeenCalledWith({
      min: 100,
      max: undefined,
    });
  });

  it('calls onChange when max value is entered', async () => {
    const user = userEvent.setup();
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const maxInput = screen.getByPlaceholderText('create.max');
    await user.type(maxInput, '500');

    expect(mockOnChange).toHaveBeenCalledWith({
      min: undefined,
      max: 500,
    });
  });

  it('calls onChange with both values when both are entered', async () => {
    const user = userEvent.setup();
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min');
    const maxInput = screen.getByPlaceholderText('create.max');

    await user.type(minInput, '100');
    await user.type(maxInput, '500');

    expect(mockOnChange).toHaveBeenLastCalledWith({
      min: 100,
      max: 500,
    });
  });

  it('displays validation error when min is greater than max', async () => {
    const user = userEvent.setup();
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min');
    const maxInput = screen.getByPlaceholderText('create.max');

    await user.type(minInput, '600');
    await user.type(maxInput, '400');

    expect(screen.getByText('create.minGreaterThanMax')).toBeInTheDocument();
    expect(minInput).toHaveClass('input-error');
    expect(maxInput).toHaveClass('input-error');
  });

  it('displays error message when provided', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        error="Please enter a valid range"
        locale="en"
      />
    );

    expect(screen.getByText('Please enter a valid range')).toBeInTheDocument();
  });

  it('initializes with object value correctly', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={{ min: 100, max: 500 }}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min') as HTMLInputElement;
    const maxInput = screen.getByPlaceholderText('create.max') as HTMLInputElement;

    expect(minInput.value).toBe('100');
    expect(maxInput.value).toBe('500');
  });

  it('initializes with string value correctly', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={JSON.stringify({ min: 200, max: 800 })}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min') as HTMLInputElement;
    const maxInput = screen.getByPlaceholderText('create.max') as HTMLInputElement;

    expect(minInput.value).toBe('200');
    expect(maxInput.value).toBe('800');
  });

  it('displays allowed range when min and max values are defined', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.getByText(/create.allowedRange: 0 - 1000 \$/)).toBeInTheDocument();
  });

  it('displays only min allowed when only min_value is defined', () => {
    const minOnlyAttribute = {
      ...mockAttribute,
      max_value: undefined,
    };

    render(
      <RangeAttribute
        attribute={minOnlyAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.getByText(/create.minAllowed: 0 \$/)).toBeInTheDocument();
  });

  it('displays only max allowed when only max_value is defined', () => {
    const maxOnlyAttribute = {
      ...mockAttribute,
      min_value: undefined,
    };

    render(
      <RangeAttribute
        attribute={maxOnlyAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.getByText(/create.maxAllowed: 1000 \$/)).toBeInTheDocument();
  });

  it('sets min and max attributes on inputs', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min');
    const maxInput = screen.getByPlaceholderText('create.max');

    expect(minInput).toHaveAttribute('min', '0');
    expect(minInput).toHaveAttribute('max', '1000');
    expect(maxInput).toHaveAttribute('min', '0');
    expect(maxInput).toHaveAttribute('max', '1000');
  });

  it('handles clearing input values', async () => {
    const user = userEvent.setup();
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={{ min: 100, max: 500 }}
        locale="en"
      />
    );

    const minInput = screen.getByPlaceholderText('create.min');
    await user.clear(minInput);

    expect(mockOnChange).toHaveBeenCalledWith({
      min: undefined,
      max: 500,
    });
  });

  it('does not show validation error when error prop is provided but values are valid', () => {
    render(
      <RangeAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={{ min: 100, max: 500 }}
        error="Some error"
        locale="en"
      />
    );

    expect(screen.getByText('Some error')).toBeInTheDocument();
    expect(screen.queryByText('create.minGreaterThanMax')).not.toBeInTheDocument();
  });

  it('renders without unit when not provided', () => {
    const noUnitAttribute = {
      ...mockAttribute,
      unit: undefined,
    };

    render(
      <RangeAttribute
        attribute={noUnitAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.queryByText('$')).not.toBeInTheDocument();
  });

  it('uses display_name when translations are not available', () => {
    const noTranslationsAttribute = {
      ...mockAttribute,
      translations: undefined,
    };

    render(
      <RangeAttribute
        attribute={noTranslationsAttribute}
        onChange={mockOnChange}
        locale="ru"
      />
    );

    expect(screen.getByText('Price')).toBeInTheDocument();
  });
});
