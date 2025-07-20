import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import MultiSelectAttribute from '../MultiSelectAttribute';

const mockAttribute = {
  id: 1,
  name: 'colors',
  display_name: 'Colors',
  options: ['Red', 'Green', 'Blue', 'Yellow'],
  translations: {
    en: 'Colors',
    ru: 'Цвета',
    sr: 'Боје',
  },
  option_translations: {
    Red: { en: 'Red', ru: 'Красный', sr: 'Црвена' },
    Green: { en: 'Green', ru: 'Зеленый', sr: 'Зелена' },
    Blue: { en: 'Blue', ru: 'Синий', sr: 'Плава' },
    Yellow: { en: 'Yellow', ru: 'Желтый', sr: 'Жута' },
  },
  is_required: true,
};

describe('MultiSelectAttribute', () => {
  const mockOnChange = jest.fn();

  beforeEach(() => {
    mockOnChange.mockClear();
  });

  it('renders with correct label and required indicator', () => {
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    expect(screen.getByText('Colors')).toBeInTheDocument();
    expect(screen.getByText('*')).toBeInTheDocument();
  });

  it('displays correct translations based on locale', () => {
    const { rerender } = render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="ru"
      />
    );

    expect(screen.getByText('Цвета')).toBeInTheDocument();

    rerender(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="sr"
      />
    );

    expect(screen.getByText('Боје')).toBeInTheDocument();
  });

  it('opens dropdown when button is clicked', async () => {
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const button = screen.getByRole('button', { name: /selectOptions/i });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText('Red')).toBeInTheDocument();
      expect(screen.getByText('Green')).toBeInTheDocument();
      expect(screen.getByText('Blue')).toBeInTheDocument();
      expect(screen.getByText('Yellow')).toBeInTheDocument();
    });
  });

  it('selects and deselects options correctly', async () => {
    const user = userEvent.setup();
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    // Open dropdown
    const button = screen.getByRole('button', { name: /selectOptions/i });
    await user.click(button);

    // Select Red
    const redCheckbox = screen.getByRole('checkbox', { name: /Red/i });
    await user.click(redCheckbox);

    expect(mockOnChange).toHaveBeenCalledWith(['Red']);
    expect(redCheckbox).toBeChecked();

    // Select Green
    const greenCheckbox = screen.getByRole('checkbox', { name: /Green/i });
    await user.click(greenCheckbox);

    expect(mockOnChange).toHaveBeenCalledWith(['Red', 'Green']);
    expect(greenCheckbox).toBeChecked();

    // Deselect Red
    await user.click(redCheckbox);
    expect(mockOnChange).toHaveBeenCalledWith(['Green']);
    expect(redCheckbox).not.toBeChecked();
  });

  it('displays selected count correctly', async () => {
    const _user = userEvent.setup();
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={['Red', 'Green']}
        locale="en"
      />
    );

    const button = screen.getByRole('button');
    expect(button).toHaveTextContent('selected: 2');
  });

  it('displays selected values as badges', () => {
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={['Red', 'Green']}
        locale="en"
      />
    );

    expect(screen.getByText('Red')).toBeInTheDocument();
    expect(screen.getByText('Green')).toBeInTheDocument();
  });

  it('removes selected value when badge X is clicked', async () => {
    const user = userEvent.setup();
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={['Red', 'Green']}
        locale="en"
      />
    );

    // Find and click the X button for Red
    const redBadge = screen.getByText('Red').closest('.badge');
    const removeButton = redBadge?.querySelector('button');

    if (removeButton) {
      await user.click(removeButton);
      expect(mockOnChange).toHaveBeenCalledWith(['Green']);
    }
  });

  it('closes dropdown when clicking outside', async () => {
    const user = userEvent.setup();
    render(
      <div>
        <div data-testid="outside">Outside element</div>
        <MultiSelectAttribute
          attribute={mockAttribute}
          onChange={mockOnChange}
          locale="en"
        />
      </div>
    );

    // Open dropdown
    const button = screen.getByRole('button', { name: /selectOptions/i });
    await user.click(button);

    // Verify dropdown is open
    expect(screen.getByText('Red')).toBeInTheDocument();

    // Click outside
    const outsideElement = screen.getByTestId('outside');
    await user.click(outsideElement);

    // Verify dropdown is closed
    await waitFor(() => {
      expect(screen.queryByRole('checkbox')).not.toBeInTheDocument();
    });
  });

  it('displays error message when provided', () => {
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        error="Please select at least one option"
        locale="en"
      />
    );

    expect(
      screen.getByText('Please select at least one option')
    ).toBeInTheDocument();
  });

  it('handles string options format', () => {
    const stringOptionsAttribute = {
      ...mockAttribute,
      options: JSON.stringify(['Red', 'Green', 'Blue']),
    };

    render(
      <MultiSelectAttribute
        attribute={stringOptionsAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const button = screen.getByRole('button', { name: /selectOptions/i });
    fireEvent.click(button);

    expect(screen.getByText('Red')).toBeInTheDocument();
    expect(screen.getByText('Green')).toBeInTheDocument();
    expect(screen.getByText('Blue')).toBeInTheDocument();
  });

  it('handles object options format', async () => {
    const objectOptionsAttribute = {
      ...mockAttribute,
      options: [
        { value: 'red', label: 'Red Color' },
        { value: 'green', label: 'Green Color' },
      ],
    };

    render(
      <MultiSelectAttribute
        attribute={objectOptionsAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const button = screen.getByRole('button', { name: /selectOptions/i });
    fireEvent.click(button);

    expect(screen.getByText('Red Color')).toBeInTheDocument();
    expect(screen.getByText('Green Color')).toBeInTheDocument();
  });

  it('handles empty options gracefully', () => {
    const emptyOptionsAttribute = {
      ...mockAttribute,
      options: [],
    };

    render(
      <MultiSelectAttribute
        attribute={emptyOptionsAttribute}
        onChange={mockOnChange}
        locale="en"
      />
    );

    const button = screen.getByRole('button', { name: /selectOptions/i });
    fireEvent.click(button);

    expect(screen.getByText('noOptions')).toBeInTheDocument();
  });

  it('initializes with string value correctly', () => {
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        value={JSON.stringify(['Red', 'Blue'])}
        locale="en"
      />
    );

    expect(screen.getByText('Red')).toBeInTheDocument();
    expect(screen.getByText('Blue')).toBeInTheDocument();
  });

  it('uses option translations when available', async () => {
    render(
      <MultiSelectAttribute
        attribute={mockAttribute}
        onChange={mockOnChange}
        locale="ru"
      />
    );

    const button = screen.getByRole('button', { name: /selectOptions/i });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText('Красный')).toBeInTheDocument();
      expect(screen.getByText('Зеленый')).toBeInTheDocument();
      expect(screen.getByText('Синий')).toBeInTheDocument();
      expect(screen.getByText('Желтый')).toBeInTheDocument();
    });
  });
});
