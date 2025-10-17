import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { CarFilters } from '../CarFilters';
import { CarsService } from '@/services/cars';

// Mock the CarsService
jest.mock('@/services/cars');

// Mock next-intl
jest.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
  useLocale: () => 'en',
}));

describe('CarFilters', () => {
  const mockOnFiltersChange = jest.fn();
  const mockCarMakes = [
    { id: 1, name: 'BMW', slug: 'bmw' },
    { id: 2, name: 'Mercedes', slug: 'mercedes' },
  ];
  const mockCarModels = [
    { id: 1, make_id: 1, name: '3 Series', slug: '3-series' },
    { id: 2, make_id: 1, name: '5 Series', slug: '5-series' },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    (CarsService.getMakes as jest.Mock).mockResolvedValue({
      success: true,
      data: mockCarMakes,
    });
    (CarsService.getModelsByMake as jest.Mock).mockResolvedValue({
      success: true,
      data: mockCarModels,
    });
  });

  it('renders without crashing', () => {
    render(<CarFilters onFiltersChange={mockOnFiltersChange} />);
    expect(screen.getByText('filters.make')).toBeInTheDocument();
  });

  it('loads car makes on mount', async () => {
    render(<CarFilters onFiltersChange={mockOnFiltersChange} />);

    await waitFor(() => {
      expect(CarsService.getMakes).toHaveBeenCalled();
    });
  });

  it('loads models when make is selected', async () => {
    render(<CarFilters onFiltersChange={mockOnFiltersChange} />);

    await waitFor(() => {
      expect(screen.getByText('filters.make')).toBeInTheDocument();
    });

    // Select BMW
    const makeSelects = screen.getAllByRole('combobox');
    const makeSelect = makeSelects[0]; // First select is make
    fireEvent.change(makeSelect, { target: { value: 'bmw' } });

    await waitFor(() => {
      expect(CarsService.getModelsByMake).toHaveBeenCalledWith('bmw');
    });
  });

  it('calls onFiltersChange when filters change', async () => {
    render(<CarFilters onFiltersChange={mockOnFiltersChange} />);

    // Change price from
    const priceFromInput = screen.getByPlaceholderText('filters.from');
    fireEvent.change(priceFromInput, { target: { value: '10000' } });

    await waitFor(() => {
      expect(mockOnFiltersChange).toHaveBeenCalled();
    });
  });

  it('toggles advanced filters', () => {
    render(<CarFilters onFiltersChange={mockOnFiltersChange} />);

    const advancedButton = screen.getByText('filters.advancedFilters');
    fireEvent.click(advancedButton);

    expect(screen.getByText('filters.bodyType')).toBeInTheDocument();
  });

  it('handles body type selection', async () => {
    render(<CarFilters onFiltersChange={mockOnFiltersChange} />);

    // Show advanced filters
    const advancedButton = screen.getByText('filters.advancedFilters');
    fireEvent.click(advancedButton);

    // Select sedan body type
    const sedanCheckbox = screen.getByLabelText('filters.bodyTypes.sedan');
    fireEvent.click(sedanCheckbox);

    await waitFor(() => {
      expect(mockOnFiltersChange).toHaveBeenCalledWith(
        expect.objectContaining({
          car_body_type: ['sedan'],
        })
      );
    });
  });

  it('resets filters when reset button is clicked', async () => {
    render(<CarFilters onFiltersChange={mockOnFiltersChange} />);

    // Set some filters
    const priceFromInput = screen.getByPlaceholderText('filters.from');
    fireEvent.change(priceFromInput, { target: { value: '10000' } });

    // Reset filters
    const resetButton = screen.getByText('filters.reset');
    fireEvent.click(resetButton);

    await waitFor(() => {
      expect(priceFromInput).toHaveValue(null);
    });
  });
});
