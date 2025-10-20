import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import SearchAutocomplete from '@/components/SearchBar/SearchAutocomplete';
import { useTranslations } from 'next-intl';

jest.mock('next-intl', () => ({
  useTranslations: jest.fn(),
}));

jest.mock('next/image', () => ({
  __esModule: true,
  default: (props: any) => {
    return <img alt="" {...props} />;
  },
}));

describe('SearchAutocomplete', () => {
  const mockOnSelect = jest.fn();
  const mockOnCategorySelect = jest.fn();
  const mockOnProductSelect = jest.fn();
  const mockTranslations = {
    suggestions: 'Suggestions',
    noSuggestions: 'No suggestions available',
    searchHistory: 'Search History',
    trending: 'Trending Searches',
    category: 'Category',
    product: 'Product',
    results: 'results',
    clear: 'Clear',
  };

  const mockSuggestions = [
    {
      type: 'product' as const,
      text: 'iPhone 13',
      product_id: 123,
      metadata: {
        category: 'Electronics',
        price: 799,
        image: '/images/iphone.jpg',
      },
    },
    {
      type: 'category' as const,
      text: 'Smartphones',
      category: { id: 10, name: 'Smartphones', slug: 'smartphones' },
      metadata: {
        count: 150,
      },
    },
    {
      type: 'text' as const,
      text: 'laptop deals',
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    (useTranslations as jest.Mock).mockReturnValue(
      (key: string) =>
        mockTranslations[key as keyof typeof mockTranslations] || key
    );
  });

  it('should not render when showSuggestions is false', () => {
    const { container } = render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={false}
        selectedIndex={-1}
        query="test"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    expect(container.firstChild).toBeNull();
  });

  it('should render suggestions when showSuggestions is true', () => {
    render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="iphone"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    expect(screen.getByText('Suggestions')).toBeInTheDocument();
    expect(
      screen.getByText((content, element) => {
        return element?.textContent === 'iPhone 13';
      })
    ).toBeInTheDocument();
  });

  it('should display product suggestions with metadata', () => {
    const { container } = render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="iphone"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    expect(container.textContent).toContain('iPhone 13');
    expect(screen.getByText('Product')).toBeInTheDocument();
    expect(container.textContent).toContain('799');
    expect(container.textContent).toContain('RSD');
  });

  it('should display category suggestions', () => {
    const { container } = render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="smart"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    expect(container.textContent).toContain('Smartphones');
    expect(screen.getByText('Category')).toBeInTheDocument();
    expect(screen.getByText('(150 results)')).toBeInTheDocument();
  });

  it('should call onProductSelect when product is clicked', () => {
    render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="iphone"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    const clickableElements = document.querySelectorAll('.cursor-pointer');
    const productElement = Array.from(clickableElements).find((el) =>
      el.textContent?.includes('iPhone 13')
    );

    if (productElement) {
      fireEvent.click(productElement);
    }

    expect(mockOnProductSelect).toHaveBeenCalledWith(123);
  });

  it('should call onCategorySelect when category is clicked', () => {
    render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="smart"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    const clickableElements = document.querySelectorAll('.cursor-pointer');
    const categoryElement = Array.from(clickableElements).find((el) =>
      el.textContent?.includes('Smartphones')
    );

    if (categoryElement) {
      fireEvent.click(categoryElement);
    }

    expect(mockOnCategorySelect).toHaveBeenCalledWith(10);
  });

  it('should call onSelect when text suggestion is clicked', () => {
    render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="laptop"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    const clickableElements = document.querySelectorAll('.cursor-pointer');
    const textElement = Array.from(clickableElements).find((el) =>
      el.textContent?.includes('laptop deals')
    );

    if (textElement) {
      fireEvent.click(textElement);
    }

    expect(mockOnSelect).toHaveBeenCalledWith('laptop deals');
  });

  it('should show loading state', () => {
    render(
      <SearchAutocomplete
        suggestions={[]}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="test"
        isLoading={true}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    const loadingSpinner = document.querySelector('.loading-spinner');
    expect(loadingSpinner).toBeInTheDocument();
  });

  it('should show no suggestions message', () => {
    render(
      <SearchAutocomplete
        suggestions={[]}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="xyz"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    expect(screen.getByText('No suggestions available')).toBeInTheDocument();
  });

  it('should show search history when provided', () => {
    const mockHistory = ['previous search 1', 'previous search 2'];

    render(
      <SearchAutocomplete
        suggestions={[]}
        searchHistory={mockHistory}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query=""
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    expect(screen.getByText('Search History')).toBeInTheDocument();
    expect(screen.getByText('previous search 1')).toBeInTheDocument();
    expect(screen.getByText('previous search 2')).toBeInTheDocument();
  });

  it('should show trending searches when provided', () => {
    const mockTrending = ['trending 1', 'trending 2'];

    render(
      <SearchAutocomplete
        suggestions={[]}
        searchHistory={[]}
        trendingSearches={mockTrending}
        showSuggestions={true}
        selectedIndex={-1}
        query=""
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    expect(screen.getByText('Trending Searches')).toBeInTheDocument();
    expect(screen.getByText('trending 1')).toBeInTheDocument();
    expect(screen.getByText('trending 2')).toBeInTheDocument();
  });

  it('should call onSelect when history item is clicked', () => {
    const mockHistory = ['previous search'];

    render(
      <SearchAutocomplete
        suggestions={[]}
        searchHistory={mockHistory}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query=""
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    fireEvent.click(screen.getByText('previous search'));

    expect(mockOnSelect).toHaveBeenCalledWith('previous search');
  });

  it('should highlight matching text in suggestions', () => {
    render(
      <SearchAutocomplete
        suggestions={mockSuggestions}
        searchHistory={[]}
        trendingSearches={[]}
        showSuggestions={true}
        selectedIndex={-1}
        query="phone"
        isLoading={false}
        onSelect={mockOnSelect}
        onCategorySelect={mockOnCategorySelect}
        onProductSelect={mockOnProductSelect}
      />
    );

    const highlightedText = document.querySelector('mark');
    expect(highlightedText).toBeTruthy();
    expect(highlightedText?.textContent?.toLowerCase()).toContain('phone');
  });
});
