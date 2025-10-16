import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import LocationFilter from '../LocationFilter';
import { useTranslations } from 'next-intl';

jest.mock('next-intl', () => ({
  useTranslations: jest.fn(),
}));

jest.mock('@/services/location', () => ({
  LocationService: {
    searchCities: jest.fn(),
  },
}));

import { LocationService } from '@/services/location';

describe('LocationFilter', () => {
  const mockOnLocationChange = jest.fn();
  const mockTranslations = {
    location: 'Location',
    enterCity: 'Enter city name',
    searchRadius: 'Search radius',
  };

  const mockCities = [
    {
      id: '1',
      name: 'Novi Sad',
      country: 'Serbia',
      lat: 45.2671,
      lng: 19.8335,
    },
    {
      id: '2',
      name: 'Belgrade',
      country: 'Serbia',
      lat: 44.7866,
      lng: 20.4489,
    },
    { id: '3', name: 'Niš', country: 'Serbia', lat: 43.3209, lng: 21.8958 },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    (useTranslations as jest.Mock).mockReturnValue(
      (key: string) => mockTranslations[key] || key
    );
    (LocationService.searchCities as jest.Mock).mockResolvedValue(mockCities);
  });

  it('should render location input and radius slider', () => {
    render(
      <LocationFilter
        location="Belgrade"
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    expect(screen.getByLabelText('Location')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter city name')).toBeInTheDocument();
    // Radius is shown only when location is set
    expect(screen.getByText(/Search radius/i)).toBeInTheDocument();
    expect(screen.getByText('10 km')).toBeInTheDocument();
  });

  it('should search for cities when typing', async () => {
    render(
      <LocationFilter
        location=""
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText('Enter city name');
    fireEvent.change(input, { target: { value: 'Nov' } });

    await waitFor(() => {
      expect(LocationService.searchCities).toHaveBeenCalledWith('Nov');
    });
  });

  it('should display city suggestions', async () => {
    render(
      <LocationFilter
        location=""
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText('Enter city name');
    fireEvent.change(input, { target: { value: 'Nov' } });

    await waitFor(() => {
      expect(screen.getByText('Novi Sad')).toBeInTheDocument();
      expect(screen.getByText('Belgrade')).toBeInTheDocument();
      expect(screen.getByText('Niš')).toBeInTheDocument();
    });
  });

  it('should select a city from suggestions', async () => {
    render(
      <LocationFilter
        location=""
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText('Enter city name');
    fireEvent.change(input, { target: { value: 'Nov' } });

    await waitFor(() => {
      expect(screen.getByText('Novi Sad')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Novi Sad'));

    expect(mockOnLocationChange).toHaveBeenCalledWith('Novi Sad', 10);
  });

  it('should update radius', () => {
    render(
      <LocationFilter
        location="Novi Sad"
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    const slider = screen.getByRole('slider');
    fireEvent.change(slider, { target: { value: '25' } });

    expect(mockOnLocationChange).toHaveBeenCalledWith('Novi Sad', 25);
  });

  it('should display current location', () => {
    render(
      <LocationFilter
        location="Belgrade"
        radius={15}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText(
      'Enter city name'
    ) as HTMLInputElement;
    expect(input.value).toBe('Belgrade');
  });

  it('should clear location when input is emptied', () => {
    render(
      <LocationFilter
        location="Belgrade"
        radius={15}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText('Enter city name');
    fireEvent.change(input, { target: { value: '' } });

    expect(mockOnLocationChange).toHaveBeenCalledWith('', 15);
  });

  it('should handle API error gracefully', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    (LocationService.searchCities as jest.Mock).mockRejectedValue(
      new Error('API Error')
    );

    render(
      <LocationFilter
        location=""
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText('Enter city name');
    fireEvent.change(input, { target: { value: 'Nov' } });

    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Failed to search cities:',
        expect.any(Error)
      );
    });

    consoleErrorSpy.mockRestore();
  });

  it('should not search with empty input', () => {
    render(
      <LocationFilter
        location=""
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText('Enter city name');
    fireEvent.change(input, { target: { value: '' } });

    expect(LocationService.searchCities).not.toHaveBeenCalled();
  });

  it('should debounce search requests', async () => {
    jest.useFakeTimers();

    render(
      <LocationFilter
        location=""
        radius={10}
        onLocationChange={mockOnLocationChange}
      />
    );

    const input = screen.getByPlaceholderText('Enter city name');

    fireEvent.change(input, { target: { value: 'N' } });
    fireEvent.change(input, { target: { value: 'No' } });
    fireEvent.change(input, { target: { value: 'Nov' } });

    expect(LocationService.searchCities).not.toHaveBeenCalled();

    jest.advanceTimersByTime(300);

    await waitFor(() => {
      expect(LocationService.searchCities).toHaveBeenCalledTimes(1);
      expect(LocationService.searchCities).toHaveBeenCalledWith('Nov');
    });

    jest.useRealTimers();
  });
});
