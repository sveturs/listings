import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { CarListingCard } from '../CarListingCard';
import { useRouter } from 'next/navigation';

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
}));

// Mock next-intl
jest.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

// Mock next/image
jest.mock('next/image', () => ({
  __esModule: true,
  default: (props: any) => {
    return <img {...props} />;
  },
}));

describe('CarListingCard', () => {
  const mockPush = jest.fn();

  const mockListing = {
    id: 1,
    title: 'BMW 3 Series 2020',
    description: 'Excellent condition BMW',
    price: 25000,
    category_id: 10100,
    user_id: 1,
    status: 'active' as const,
    condition: 'excellent' as const,
    created_at: '2025-01-01T00:00:00Z',
    updated_at: '2025-01-01T00:00:00Z',
    views_count: 100,
    images: [{ id: 1, url: '/test-image.jpg', is_main: true }],
    attributes: [
      { attribute_name: 'car_make', string_value: 'BMW', display_name: 'Make' },
      {
        attribute_name: 'car_model',
        string_value: '3 Series',
        display_name: 'Model',
      },
      { attribute_name: 'year', string_value: '2020', display_name: 'Year' },
      {
        attribute_name: 'mileage',
        string_value: '50000',
        display_name: 'Mileage',
      },
      {
        attribute_name: 'fuel_type',
        string_value: 'Gasoline',
        display_name: 'Fuel Type',
      },
      {
        attribute_name: 'transmission',
        string_value: 'Automatic',
        display_name: 'Transmission',
      },
    ],
    city: 'Belgrade',
    country: 'Serbia',
  };

  beforeEach(() => {
    jest.clearAllMocks();
    (useRouter as jest.Mock).mockReturnValue({
      push: mockPush,
    });
  });

  it('renders listing information correctly', () => {
    render(<CarListingCard listing={mockListing} locale="en" />);

    expect(screen.getByText('BMW 3 Series 2020')).toBeInTheDocument();
    expect(screen.getByText('€25,000')).toBeInTheDocument();
    expect(screen.getByText('2020')).toBeInTheDocument();
    expect(screen.getByText('50000 km')).toBeInTheDocument();
    expect(screen.getByText('Gasoline')).toBeInTheDocument();
    expect(screen.getByText('Automatic')).toBeInTheDocument();
  });

  it('renders in grid view mode', () => {
    render(
      <CarListingCard listing={mockListing} locale="en" viewMode="grid" />
    );

    const card = screen.getByRole('article');
    expect(card).toHaveClass('card');
  });

  it('renders in list view mode', () => {
    render(
      <CarListingCard listing={mockListing} locale="en" viewMode="list" />
    );

    const card = screen.getByRole('article');
    expect(card).toHaveClass('flex-row');
  });

  it('navigates to listing detail on click', () => {
    render(<CarListingCard listing={mockListing} locale="en" />);

    const card = screen.getByRole('article');
    fireEvent.click(card);

    expect(mockPush).toHaveBeenCalledWith('/en/listing/1');
  });

  it('displays location if available', () => {
    render(<CarListingCard listing={mockListing} locale="en" />);

    expect(screen.getByText(/Belgrade, Serbia/)).toBeInTheDocument();
  });

  it('displays views count', () => {
    render(<CarListingCard listing={mockListing} locale="en" />);

    expect(screen.getByText(/100/)).toBeInTheDocument();
  });

  it('displays condition badge', () => {
    render(<CarListingCard listing={mockListing} locale="en" />);

    expect(screen.getByText('conditions.excellent')).toBeInTheDocument();
  });

  it('handles listing without images gracefully', () => {
    const listingWithoutImages = {
      ...mockListing,
      images: [],
    };

    render(<CarListingCard listing={listingWithoutImages} locale="en" />);

    // Should render placeholder or default image
    const image = screen.getByRole('img');
    expect(image).toBeInTheDocument();
  });

  it('handles listing without attributes gracefully', () => {
    const listingWithoutAttributes = {
      ...mockListing,
      attributes: [],
    };

    render(<CarListingCard listing={listingWithoutAttributes} locale="en" />);

    // Should still render the card
    expect(screen.getByText('BMW 3 Series 2020')).toBeInTheDocument();
  });

  it('formats mileage correctly', () => {
    const listingWithHighMileage = {
      ...mockListing,
      attributes: [
        ...mockListing.attributes.filter((a) => a.attribute_name !== 'mileage'),
        {
          attribute_name: 'mileage',
          string_value: '150000',
          display_name: 'Mileage',
        },
      ],
    };

    render(<CarListingCard listing={listingWithHighMileage} locale="en" />);

    expect(screen.getByText('150000 km')).toBeInTheDocument();
  });
});
