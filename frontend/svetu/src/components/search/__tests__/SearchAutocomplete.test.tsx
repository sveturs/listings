import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import SearchAutocomplete from '@/components/SearchBar/SearchAutocomplete';
import { SearchService } from '@/services/search';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';

jest.mock('next-intl', () => ({
  useTranslations: jest.fn(),
}));

jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
}));

jest.mock('@/services/search', () => ({
  SearchService: {
    getAutocompleteSuggestions: jest.fn(),
  },
}));

jest.mock('next/image', () => ({
  __esModule: true,
  default: (props: any) => {
    return <img {...props} />;
  },
}));

describe('SearchAutocomplete', () => {
  const mockOnSearch = jest.fn();
  const mockOnSuggestionSelect = jest.fn();
  const mockRouter = { push: jest.fn() };
  const mockTranslations = {
    suggestions: 'Suggestions',
    noSuggestions: 'No suggestions available',
    history: 'Search History',
    trending: 'Trending Searches',
    inCategory: 'in category',
  };

  const mockSuggestions = [
    {
      id: '1',
      type: 'product' as const,
      text: 'iPhone 13',
      highlight: '<em>iPhone</em> 13',
      metadata: {
        category: 'Electronics',
        price: 799,
        image: '/images/iphone.jpg',
      },
    },
    {
      id: '2',
      type: 'category' as const,
      text: 'Smartphones',
      highlight: '<em>Smart</em>phones',
      metadata: {
        count: 150,
      },
    },
    {
      id: '3',
      type: 'query' as const,
      text: 'laptop deals',
      highlight: 'laptop <em>deals</em>',
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    (useTranslations as jest.Mock).mockReturnValue(
      (key: string) => mockTranslations[key] || key
    );
    (useRouter as jest.Mock).mockReturnValue(mockRouter);
    (SearchService.getAutocompleteSuggestions as jest.Mock).mockResolvedValue(
      mockSuggestions
    );
  });

  it('should render search input', () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    expect(screen.getByRole('combobox')).toBeInTheDocument();
  });

  it('should show suggestions when typing', async () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'iphone' } });

    await waitFor(() => {
      expect(SearchService.getAutocompleteSuggestions).toHaveBeenCalledWith(
        'iphone'
      );
      expect(screen.getByText('Suggestions')).toBeInTheDocument();
    });
  });

  it('should display product suggestions with images', async () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'iphone' } });

    await waitFor(() => {
      const productSuggestion = screen.getByText('iPhone 13');
      expect(productSuggestion).toBeInTheDocument();

      const image = screen.getByAltText('iPhone 13');
      expect(image).toHaveAttribute('src', '/images/iphone.jpg');
    });
  });

  it('should display category suggestions', async () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'smart' } });

    await waitFor(() => {
      expect(screen.getByText('Smartphones')).toBeInTheDocument();
      expect(screen.getByText('(150)')).toBeInTheDocument();
    });
  });

  it('should handle suggestion selection', async () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'iphone' } });

    await waitFor(() => {
      expect(screen.getByText('iPhone 13')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('iPhone 13'));

    expect(mockOnSuggestionSelect).toHaveBeenCalledWith(mockSuggestions[0]);
  });

  it('should navigate with keyboard', async () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'test' } });

    await waitFor(() => {
      expect(screen.getByText('iPhone 13')).toBeInTheDocument();
    });

    fireEvent.keyDown(input, { key: 'ArrowDown' });
    fireEvent.keyDown(input, { key: 'ArrowDown' });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(mockOnSuggestionSelect).toHaveBeenCalledWith(mockSuggestions[1]);
  });

  it('should close suggestions on escape', async () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'test' } });

    await waitFor(() => {
      expect(screen.getByText('Suggestions')).toBeInTheDocument();
    });

    fireEvent.keyDown(input, { key: 'Escape' });

    await waitFor(() => {
      expect(screen.queryByText('Suggestions')).not.toBeInTheDocument();
    });
  });

  it('should show search history when enabled', async () => {
    const mockHistory = ['previous search 1', 'previous search 2'];

    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
        showHistory={true}
        searchHistory={mockHistory}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.focus(input);

    await waitFor(() => {
      expect(screen.getByText('Search History')).toBeInTheDocument();
      expect(screen.getByText('previous search 1')).toBeInTheDocument();
      expect(screen.getByText('previous search 2')).toBeInTheDocument();
    });
  });

  it('should show trending searches', async () => {
    const mockTrending = ['trending 1', 'trending 2'];

    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
        trendingSearches={mockTrending}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.focus(input);

    await waitFor(() => {
      expect(screen.getByText('Trending Searches')).toBeInTheDocument();
      expect(screen.getByText('trending 1')).toBeInTheDocument();
      expect(screen.getByText('trending 2')).toBeInTheDocument();
    });
  });

  it('should handle API error gracefully', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    (SearchService.getAutocompleteSuggestions as jest.Mock).mockRejectedValue(
      new Error('API Error')
    );

    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'test' } });

    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Failed to fetch suggestions:',
        expect.any(Error)
      );
    });

    consoleErrorSpy.mockRestore();
  });

  it('should debounce API calls', async () => {
    jest.useFakeTimers();

    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');

    fireEvent.change(input, { target: { value: 'i' } });
    fireEvent.change(input, { target: { value: 'ip' } });
    fireEvent.change(input, { target: { value: 'iph' } });
    fireEvent.change(input, { target: { value: 'ipho' } });
    fireEvent.change(input, { target: { value: 'iphone' } });

    expect(SearchService.getAutocompleteSuggestions).not.toHaveBeenCalled();

    jest.advanceTimersByTime(300);

    await waitFor(() => {
      expect(SearchService.getAutocompleteSuggestions).toHaveBeenCalledTimes(1);
      expect(SearchService.getAutocompleteSuggestions).toHaveBeenCalledWith(
        'iphone'
      );
    });

    jest.useRealTimers();
  });

  it('should call onSearch when form is submitted', () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'search query' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(mockOnSearch).toHaveBeenCalledWith('search query');
  });

  it('should highlight matching text in suggestions', async () => {
    render(
      <SearchAutocomplete
        query=""
        onSearch={mockOnSearch}
        onSuggestionSelect={mockOnSuggestionSelect}
      />
    );

    const input = screen.getByRole('combobox');
    fireEvent.change(input, { target: { value: 'phone' } });

    await waitFor(() => {
      const highlightedText = document.querySelector('em');
      expect(highlightedText).toHaveTextContent('iPhone');
    });
  });
});
