import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import CategorySelector from './CategorySelector';

// Мокаем react-i18next
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, params?: any) => {
      if (key === 'admin.categories.searchPlaceholder') return 'Search category...';
      if (key === 'admin.categories.noCategories') return 'No categories found';
      if (key === 'admin.categories.listingCount') {
        return params?.count === 1 ? '1 listing' : `${params?.count} listings`;
      }
      if (key === 'admin.categories.searchResults') return 'search results';
      if (key === 'admin.categories.path') return 'Path';
      if (key === 'admin.categories.slug') return 'Slug';
      if (key === 'admin.categories.parent') return 'Parent';
      if (key === 'admin.categories.noParent') return 'No parent';
      if (key === 'admin.categories.translations') return 'Translations';
      if (key === 'admin.categories.customUi') return 'Custom UI';
      return key;
    }
  })
}));

describe('CategorySelector', () => {
  const mockCategories = [
    {
      id: 1,
      name: 'Electronics',
      slug: 'electronics',
      parent_id: null,
      listing_count: 10,
      has_custom_ui: false
    },
    {
      id: 2,
      name: 'Smartphones',
      slug: 'smartphones',
      parent_id: 1,
      listing_count: 5,
      has_custom_ui: false
    },
    {
      id: 3,
      name: 'Laptops',
      slug: 'laptops',
      parent_id: 1,
      listing_count: 3,
      has_custom_ui: true,
      custom_ui_component: 'LaptopForm'
    },
    {
      id: 4,
      name: 'Clothing',
      slug: 'clothing',
      parent_id: null,
      listing_count: 7,
      has_custom_ui: false,
      translations: {
        ru: 'Одежда',
        sr: 'Odeća'
      }
    }
  ];

  const mockOnCategorySelect = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders empty state when no categories are available', () => {
    render(
      <CategorySelector
        categories={[]}
        onCategorySelect={mockOnCategorySelect}
      />
    );

    expect(screen.getByText('No categories found')).toBeInTheDocument();
  });

  test('renders category tree correctly', () => {
    render(
      <CategorySelector
        categories={mockCategories}
        onCategorySelect={mockOnCategorySelect}
      />
    );

    // Проверка, что родительские категории отображаются
    expect(screen.getByText('Electronics')).toBeInTheDocument();
    expect(screen.getByText('Clothing')).toBeInTheDocument();
    
    // Дочерние категории могут быть скрыты до раскрытия
    // Нажимаем, чтобы раскрыть категорию Electronics
    const electronicsExpandButtons = screen.getAllByRole('button').filter(
      button => button.textContent?.includes('Electronics')
    );
    fireEvent.click(electronicsExpandButtons[0]);
    
    // Теперь должны быть видны дочерние категории
    expect(screen.getByText('Smartphones')).toBeInTheDocument();
    expect(screen.getByText('Laptops')).toBeInTheDocument();
  });

  test('calls onCategorySelect when a category is clicked', () => {
    render(
      <CategorySelector
        categories={mockCategories}
        onCategorySelect={mockOnCategorySelect}
      />
    );

    // Нажимаем на категорию
    const electronicsButtons = screen.getAllByRole('button').filter(
      button => button.textContent?.includes('Electronics')
    );
    fireEvent.click(electronicsButtons[0]);
    
    // Проверяем, что функция обратного вызова была вызвана с правильной категорией
    expect(mockOnCategorySelect).toHaveBeenCalledWith({
      id: 1,
      name: 'Electronics',
      slug: 'electronics',
      parent_id: null,
      listing_count: 10,
      has_custom_ui: false
    });
  });

  test('filters categories based on search query', () => {
    render(
      <CategorySelector
        categories={mockCategories}
        onCategorySelect={mockOnCategorySelect}
      />
    );

    // Вводим поисковый запрос
    const searchInput = screen.getByPlaceholderText('Search category...');
    fireEvent.change(searchInput, { target: { value: 'smart' } });
    
    // Должна отображаться только категория Smartphones
    expect(screen.getByText('Smartphones')).toBeInTheDocument();
    
    // Electronics также должна быть видна, т.к. это родитель Smartphones
    expect(screen.getByText('Electronics')).toBeInTheDocument();
    
    // Clothing не должна быть видна
    expect(screen.queryByText('Clothing')).not.toBeInTheDocument();
  });

  test('shows category details when a category is selected', () => {
    render(
      <CategorySelector
        categories={mockCategories}
        onCategorySelect={mockOnCategorySelect}
        selectedCategoryId={4} // Clothing, которая имеет переводы
      />
    );

    // Проверяем, что детали категории отображаются
    expect(screen.getByText('Clothing')).toBeInTheDocument(); // Название категории
    expect(screen.getByText('Slug:')).toBeInTheDocument();
    expect(screen.getByText('clothing')).toBeInTheDocument(); // slug категории
    expect(screen.getByText('Parent:')).toBeInTheDocument();
    expect(screen.getByText('No parent')).toBeInTheDocument(); // У категории нет родителя
    
    // Проверяем, что отображается информация о кол-ве объявлений
    expect(screen.getByText('7 listings')).toBeInTheDocument();
    
    // Проверяем, что отображаются переводы
    expect(screen.getByText('Translations:')).toBeInTheDocument();
    expect(screen.getByText('RU:')).toBeInTheDocument();
    expect(screen.getByText('Одежда')).toBeInTheDocument();
    expect(screen.getByText('SR:')).toBeInTheDocument();
    expect(screen.getByText('Odeća')).toBeInTheDocument();
  });

  test('expands parent categories when a child category is selected', () => {
    render(
      <CategorySelector
        categories={mockCategories}
        onCategorySelect={mockOnCategorySelect}
        selectedCategoryId={2} // Smartphones, дочерняя категория для Electronics
      />
    );

    // Проверяем, что родительская категория Electronics раскрыта
    // и дочерняя категория Smartphones видна
    expect(screen.getByText('Smartphones')).toBeInTheDocument();
  });

  test('toggles category expansion when expand button is clicked', () => {
    render(
      <CategorySelector
        categories={mockCategories}
        onCategorySelect={mockOnCategorySelect}
      />
    );

    // Проверяем, что дочерние категории изначально не видны
    expect(screen.queryByText('Smartphones')).not.toBeInTheDocument();
    
    // Находим кнопку раскрытия для Electronics
    const expandButtons = screen.getAllByRole('button').filter(
      button => button.textContent?.includes('Electronics')
    );
    
    // Нажимаем на кнопку раскрытия
    fireEvent.click(expandButtons[0]);
    
    // Теперь дочерние категории должны быть видны
    expect(screen.getByText('Smartphones')).toBeInTheDocument();
    expect(screen.getByText('Laptops')).toBeInTheDocument();
    
    // Нажимаем еще раз, чтобы свернуть
    fireEvent.click(expandButtons[0]);
    
    // Дочерние категории должны снова быть скрыты
    expect(screen.queryByText('Smartphones')).not.toBeInTheDocument();
  });

  test('shows custom UI information for categories that have it', () => {
    render(
      <CategorySelector
        categories={mockCategories}
        onCategorySelect={mockOnCategorySelect}
        selectedCategoryId={3} // Laptops с custom_ui_component
      />
    );

    // Раскрываем родительскую категорию, чтобы увидеть Laptops
    const expandButtons = screen.getAllByRole('button').filter(
      button => button.textContent?.includes('Electronics')
    );
    fireEvent.click(expandButtons[0]);
    
    // Нажимаем на категорию Laptops
    const laptopsButtons = screen.getAllByRole('button').filter(
      button => button.textContent?.includes('Laptops')
    );
    fireEvent.click(laptopsButtons[0]);
    
    // Проверяем, что информация о пользовательском UI отображается
    expect(screen.getByText('LaptopForm')).toBeInTheDocument();
  });
});