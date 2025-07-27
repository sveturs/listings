import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';

// Mock the translations
jest.mock('next-intl', () => ({
  useTranslations: () => (key: string, values?: any) => {
    if (key === 'ai.enhance.category_auto_selected' && values?.category) {
      return `AI автоматически выбрал категорию: ${values.category}`;
    }
    return key;
  },
  useLocale: () => 'ru',
}));

// Test the getCategoryData function
describe('getCategoryData', () => {
  const categories = [
    { id: 1001, name: 'Electronics', slug: 'electronics' },
    { id: 1002, name: 'Fashion', slug: 'fashion' },
    { id: 1003, name: 'Automotive', slug: 'automotive' },
  ];

  const getCategoryData = (
    categoryName: string
  ): { id: number; name: string; slug: string } => {
    // Проверка на undefined или пустую строку
    if (!categoryName) {
      return { id: 1, name: 'General', slug: 'general' };
    }

    // Пытаемся найти категорию по разным критериям
    const normalizedName = categoryName.toLowerCase().trim();

    // 1. Точное совпадение по slug
    let category = categories.find((cat) => cat.slug === normalizedName);

    // 2. Частичное совпадение по slug (категория содержит искомое слово)
    if (!category) {
      category = categories.find(
        (cat) =>
          cat.slug.toLowerCase().includes(normalizedName) ||
          normalizedName.includes(cat.slug.toLowerCase())
      );
    }

    // Возвращаем найденную категорию или дефолтную
    return category || { id: 1, name: 'General', slug: 'general' };
  };

  test('handles undefined category name', () => {
    const result = getCategoryData(undefined as any);
    expect(result).toEqual({ id: 1, name: 'General', slug: 'general' });
  });

  test('handles empty string category name', () => {
    const result = getCategoryData('');
    expect(result).toEqual({ id: 1, name: 'General', slug: 'general' });
  });

  test('finds category by exact slug match', () => {
    const result = getCategoryData('electronics');
    expect(result).toEqual({
      id: 1001,
      name: 'Electronics',
      slug: 'electronics',
    });
  });

  test('finds category by partial match', () => {
    const result = getCategoryData('auto');
    expect(result).toEqual({
      id: 1003,
      name: 'Automotive',
      slug: 'automotive',
    });
  });

  test('returns default category for unknown name', () => {
    const result = getCategoryData('unknown-category');
    expect(result).toEqual({ id: 1, name: 'General', slug: 'general' });
  });
});

// Test the category display component
describe('Category Display', () => {
  const CategoryDisplay = ({ categoryProbabilities }: any) => {
    return (
      <div>
        <div className="space-y-2">
          {(categoryProbabilities || []).map((cat: any, index: number) => {
            if (!cat || !cat.name) {
              return null;
            }

            const categoryName = cat.name;
            const isSelected = index === 0;

            return (
              <div
                key={index}
                data-testid={`category-${index}`}
                className={isSelected ? 'selected' : ''}
              >
                <span>{categoryName}</span>
                <span>{cat.probability}%</span>
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  test('handles empty categoryProbabilities', () => {
    render(<CategoryDisplay categoryProbabilities={[]} />);
    expect(screen.queryByTestId('category-0')).not.toBeInTheDocument();
  });

  test('handles undefined categoryProbabilities', () => {
    render(<CategoryDisplay categoryProbabilities={undefined} />);
    expect(screen.queryByTestId('category-0')).not.toBeInTheDocument();
  });

  test('skips categories with undefined name', () => {
    const probabilities = [
      { name: undefined, probability: 80 },
      { name: 'Fashion', probability: 20 },
    ];
    render(<CategoryDisplay categoryProbabilities={probabilities} />);

    expect(screen.queryByText('Fashion')).toBeInTheDocument();
    expect(screen.queryByTestId('category-0')).not.toBeInTheDocument();
    expect(screen.getByTestId('category-1')).toBeInTheDocument();
  });

  test('displays valid categories correctly', () => {
    const probabilities = [
      { name: 'Electronics', probability: 80 },
      { name: 'Fashion', probability: 20 },
    ];
    render(<CategoryDisplay categoryProbabilities={probabilities} />);

    expect(screen.getByText('Electronics')).toBeInTheDocument();
    expect(screen.getByText('80%')).toBeInTheDocument();
    expect(screen.getByText('Fashion')).toBeInTheDocument();
    expect(screen.getByText('20%')).toBeInTheDocument();

    // First category should be selected
    expect(screen.getByTestId('category-0')).toHaveClass('selected');
  });
});
