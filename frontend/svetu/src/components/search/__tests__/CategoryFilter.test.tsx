import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import CategoryFilter from '../CategoryFilter';
import { CategoryService } from '@/services/category';
import { useTranslations, useLocale } from 'next-intl';

jest.mock('next-intl', () => ({
  useTranslations: jest.fn(),
  useLocale: jest.fn(),
}));

jest.mock('@/services/category', () => ({
  CategoryService: {
    getCategories: jest.fn(),
  },
}));

describe('CategoryFilter', () => {
  const mockOnCategoryChange = jest.fn();
  const mockTranslations = {
    categories: 'Categories',
    expandAll: 'Expand all',
    collapseAll: 'Collapse all',
    clear: 'Clear',
    noCategories: 'No categories available',
  };

  const mockCategories = [
    {
      id: 1,
      name: 'Electronics',
      slug: 'electronics',
      parent_id: null,
      listing_count: 50,
      sort_order: 1,
      translations: { en: 'Electronics', ru: 'Электроника' },
    },
    {
      id: 2,
      name: 'Smartphones',
      slug: 'smartphones',
      parent_id: 1,
      listing_count: 30,
      sort_order: 1,
      translations: { en: 'Smartphones', ru: 'Смартфоны' },
    },
    {
      id: 3,
      name: 'Laptops',
      slug: 'laptops',
      parent_id: 1,
      listing_count: 20,
      sort_order: 2,
      translations: { en: 'Laptops', ru: 'Ноутбуки' },
    },
    {
      id: 4,
      name: 'Furniture',
      slug: 'furniture',
      parent_id: null,
      listing_count: 25,
      sort_order: 2,
      translations: { en: 'Furniture', ru: 'Мебель' },
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    (useTranslations as jest.Mock).mockReturnValue(
      (key: string) => mockTranslations[key] || key
    );
    (useLocale as jest.Mock).mockReturnValue('en');
    (CategoryService.getCategories as jest.Mock).mockResolvedValue(
      mockCategories
    );
  });

  it('should render loading skeleton initially', () => {
    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    const skeletons = screen.getAllByTestId('skeleton');
    expect(skeletons).toHaveLength(3);
  });

  it('should load and display categories in tree structure', async () => {
    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Electronics')).toBeInTheDocument();
      expect(screen.getByText('Furniture')).toBeInTheDocument();
    });

    expect(screen.queryByText('Smartphones')).not.toBeInTheDocument();
    expect(screen.queryByText('Laptops')).not.toBeInTheDocument();
  });

  it('should expand category to show children', async () => {
    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Electronics')).toBeInTheDocument();
    });

    const expandButton = screen.getByRole('button', {
      name: /expand electronics/i,
    });
    fireEvent.click(expandButton);

    await waitFor(() => {
      expect(screen.getByText('Smartphones')).toBeInTheDocument();
      expect(screen.getByText('Laptops')).toBeInTheDocument();
    });
  });

  it('should handle category selection', async () => {
    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Electronics')).toBeInTheDocument();
    });

    const electronicsCheckbox = screen.getByRole('checkbox', {
      name: /electronics/i,
    });
    fireEvent.click(electronicsCheckbox);

    expect(mockOnCategoryChange).toHaveBeenCalledWith([1]);
  });

  it('should handle multiple category selections', async () => {
    render(
      <CategoryFilter
        selectedCategories={[1]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Furniture')).toBeInTheDocument();
    });

    const furnitureCheckbox = screen.getByRole('checkbox', {
      name: /furniture/i,
    });
    fireEvent.click(furnitureCheckbox);

    expect(mockOnCategoryChange).toHaveBeenCalledWith([1, 4]);
  });

  it('should deselect category when clicked again', async () => {
    render(
      <CategoryFilter
        selectedCategories={[1, 4]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Electronics')).toBeInTheDocument();
    });

    const electronicsCheckbox = screen.getByRole('checkbox', {
      name: /electronics/i,
    });
    fireEvent.click(electronicsCheckbox);

    expect(mockOnCategoryChange).toHaveBeenCalledWith([4]);
  });

  it('should expand all categories', async () => {
    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Electronics')).toBeInTheDocument();
    });

    const expandAllButton = screen.getByTitle('Expand all');
    fireEvent.click(expandAllButton);

    await waitFor(() => {
      expect(screen.getByText('Smartphones')).toBeInTheDocument();
      expect(screen.getByText('Laptops')).toBeInTheDocument();
    });
  });

  it('should collapse all categories', async () => {
    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Electronics')).toBeInTheDocument();
    });

    const expandAllButton = screen.getByTitle('Expand all');
    fireEvent.click(expandAllButton);

    await waitFor(() => {
      expect(screen.getByText('Smartphones')).toBeInTheDocument();
    });

    const collapseAllButton = screen.getByTitle('Collapse all');
    fireEvent.click(collapseAllButton);

    await waitFor(() => {
      expect(screen.queryByText('Smartphones')).not.toBeInTheDocument();
      expect(screen.queryByText('Laptops')).not.toBeInTheDocument();
    });
  });

  it('should clear all selections', async () => {
    render(
      <CategoryFilter
        selectedCategories={[1, 2, 3]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Clear (3)')).toBeInTheDocument();
    });

    const clearButton = screen.getByText('Clear (3)');
    fireEvent.click(clearButton);

    expect(mockOnCategoryChange).toHaveBeenCalledWith([]);
  });

  it('should show selected children count badge', async () => {
    render(
      <CategoryFilter
        selectedCategories={[2, 3]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Electronics')).toBeInTheDocument();
    });

    const badge = screen.getByText('2');
    expect(badge).toHaveClass('badge', 'badge-primary', 'badge-xs');
  });

  it('should display listing counts', async () => {
    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('(50)')).toBeInTheDocument();
      expect(screen.getByText('(25)')).toBeInTheDocument();
    });
  });

  it('should handle API error gracefully', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    (CategoryService.getCategories as jest.Mock).mockRejectedValue(
      new Error('API Error')
    );

    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.queryByText('Electronics')).not.toBeInTheDocument();
    });

    expect(consoleErrorSpy).toHaveBeenCalledWith(
      'Failed to load categories:',
      expect.any(Error)
    );
    consoleErrorSpy.mockRestore();
  });

  it('should use locale for translations', async () => {
    (useLocale as jest.Mock).mockReturnValue('ru');

    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Электроника')).toBeInTheDocument();
      expect(screen.getByText('Мебель')).toBeInTheDocument();
    });
  });

  it('should display empty state when no categories', async () => {
    (CategoryService.getCategories as jest.Mock).mockResolvedValue([]);

    render(
      <CategoryFilter
        selectedCategories={[]}
        onCategoryChange={mockOnCategoryChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('No categories available')).toBeInTheDocument();
    });
  });
});
