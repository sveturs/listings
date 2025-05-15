import React from 'react';
import { render, screen, waitFor, fireEvent, act } from '@testing-library/react';
import '@testing-library/jest-dom';
import axios from '../../api/axios';
import CategoryAttributeMappingPage from './CategoryAttributeMappingPage';

// Мокаем модули
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, params?: any) => {
      // Простая имитация функции перевода
      if (key === 'admin.categoryAttributes.title') return 'Category Attribute Mapping';
      if (key === 'admin.categoryAttributes.description') return 'Manage attributes for each category';
      if (key === 'admin.categoryAttributes.selectCategory') return 'Select Category';
      if (key === 'admin.categoryAttributes.fetchError') return 'Error fetching data';
      if (key === 'admin.common.refresh') return 'Refresh';
      if (key === 'admin.categoryAttributes.selectCategoryPrompt') return 'Please select a category';
      
      // Если передан параметр category
      if (key === 'admin.categoryAttributes.categoryAttributes' && params?.category) {
        return `Attributes for ${params.category}`;
      }
      
      return key;
    }
  })
}));

jest.mock('../../api/axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

// Мокаем компоненты
jest.mock('../../components/admin/AttributeMappingList', () => {
  return ({ categoryId, onError }: { categoryId: number, onError: (error: string) => void }) => (
    <div data-testid="attribute-mapping-list" data-category-id={categoryId}>
      Attribute Mapping List
    </div>
  );
});

jest.mock('../../components/admin/AddAttributeToCategory', () => {
  return ({ categoryId, onAttributeAdded, onError }: { 
    categoryId: number, 
    onAttributeAdded: () => void, 
    onError: (error: string) => void 
  }) => (
    <div data-testid="add-attribute-to-category" data-category-id={categoryId}>
      Add Attribute To Category
    </div>
  );
});

jest.mock('../../components/admin/CategorySelector', () => {
  return ({ categories, onCategorySelect, selectedCategoryId }: { 
    categories: any[], 
    onCategorySelect: (category: any) => void, 
    selectedCategoryId?: number 
  }) => (
    <div data-testid="category-selector">
      <button 
        onClick={() => onCategorySelect({ id: 1, name: 'Electronics', slug: 'electronics', listing_count: 5 })}
        data-testid="select-category-button"
      >
        Select Category
      </button>
    </div>
  );
});

jest.mock('../../components/admin/CustomUIManager', () => {
  return ({ categoryId, onCategoryUpdate }: { categoryId: number, onCategoryUpdate: () => void }) => (
    <div data-testid="custom-ui-manager" data-category-id={categoryId}>
      Custom UI Manager
    </div>
  );
});

jest.mock('../../components/admin/CategoryAttributeExporter', () => {
  return ({ categoryId, onSuccess, onError }: { 
    categoryId: number, 
    onSuccess: (message: string) => void,
    onError: (error: string) => void 
  }) => (
    <div data-testid="category-attribute-exporter" data-category-id={categoryId}>
      Category Attribute Exporter
    </div>
  );
});

describe('CategoryAttributeMappingPage', () => {
  const mockCategories = [
    { id: 1, name: 'Electronics', slug: 'electronics', listing_count: 5 },
    { id: 2, name: 'Clothing', slug: 'clothing', listing_count: 10 }
  ];
  
  const mockAttributes = [
    { id: 1, name: 'color', display_name: 'Color', attribute_type: 'select' },
    { id: 2, name: 'size', display_name: 'Size', attribute_type: 'select' }
  ];
  
  const mockCategoryAttributes = [
    { category_id: 1, attribute_id: 1, is_required: true, is_enabled: true, sort_order: 1 },
    { category_id: 1, attribute_id: 2, is_required: false, is_enabled: true, sort_order: 2 }
  ];

  beforeEach(() => {
    // Сбрасываем моки перед каждым тестом
    jest.clearAllMocks();
    
    // Настраиваем успешные ответы API
    mockedAxios.get.mockImplementation((url) => {
      if (url === '/api/admin/categories') {
        return Promise.resolve({ data: mockCategories });
      }
      if (url === '/api/admin/attributes') {
        return Promise.resolve({ data: mockAttributes });
      }
      if (url === '/api/admin/categories/1/attributes') {
        return Promise.resolve({ data: mockCategoryAttributes });
      }
      return Promise.reject(new Error('not found'));
    });
  });

  test('renders initial loading state', async () => {
    render(<CategoryAttributeMappingPage />);
    
    // Проверяем, что отображается заголовок и описание
    expect(screen.getByText('Category Attribute Mapping')).toBeInTheDocument();
    expect(screen.getByText('Manage attributes for each category')).toBeInTheDocument();
    
    // Проверяем наличие индикатора загрузки
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
    
    // После загрузки данных индикатор должен исчезнуть
    await waitFor(() => {
      expect(screen.queryByRole('progressbar')).not.toBeInTheDocument();
    });
  });

  test('loads categories and attributes on mount', async () => {
    render(<CategoryAttributeMappingPage />);
    
    await waitFor(() => {
      // Проверяем, что API был вызван для загрузки категорий и атрибутов
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories');
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/attributes');
    });
    
    // Ожидаем, что будет отображаться CategorySelector
    expect(screen.getByTestId('category-selector')).toBeInTheDocument();
    
    // Ожидаем, что будет отображаться подсказка выбрать категорию
    expect(screen.getByText('Please select a category')).toBeInTheDocument();
  });

  test('selects a category and loads its attributes', async () => {
    render(<CategoryAttributeMappingPage />);
    
    // Ждем, пока загрузятся категории
    await waitFor(() => {
      expect(screen.getByTestId('category-selector')).toBeInTheDocument();
    });
    
    // Нажимаем на кнопку выбора категории
    fireEvent.click(screen.getByTestId('select-category-button'));
    
    // Проверяем, что был сделан запрос для загрузки атрибутов категории
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories/1/attributes');
    });
    
    // Проверяем, что отображаются компоненты для управления атрибутами
    expect(screen.getByTestId('attribute-mapping-list')).toBeInTheDocument();
    expect(screen.getByTestId('add-attribute-to-category')).toBeInTheDocument();
    expect(screen.getByTestId('custom-ui-manager')).toBeInTheDocument();
    expect(screen.getByTestId('category-attribute-exporter')).toBeInTheDocument();
    
    // Проверяем, что компоненты получают правильный categoryId
    expect(screen.getByTestId('attribute-mapping-list').getAttribute('data-category-id')).toBe('1');
    expect(screen.getByTestId('add-attribute-to-category').getAttribute('data-category-id')).toBe('1');
    expect(screen.getByTestId('custom-ui-manager').getAttribute('data-category-id')).toBe('1');
    expect(screen.getByTestId('category-attribute-exporter').getAttribute('data-category-id')).toBe('1');
  });

  test('handles API error', async () => {
    // Имитируем ошибку при загрузке категорий
    mockedAxios.get.mockRejectedValueOnce(new Error('Network error'));
    
    render(<CategoryAttributeMappingPage />);
    
    // Ожидаем появления сообщения об ошибке
    await waitFor(() => {
      expect(screen.getByText('Error fetching data')).toBeInTheDocument();
    });
  });

  test('refreshes data when refresh button is clicked', async () => {
    render(<CategoryAttributeMappingPage />);
    
    // Ждем загрузки данных
    await waitFor(() => {
      expect(screen.getByTestId('category-selector')).toBeInTheDocument();
    });
    
    // Находим и нажимаем кнопку обновления
    const refreshButton = screen.getByTitle('Refresh');
    fireEvent.click(refreshButton);
    
    // Проверяем, что API был вызван повторно
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories');
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/attributes');
      // Первый вызов при монтировании, второй при обновлении
      expect(mockedAxios.get).toHaveBeenCalledTimes(4);
    });
  });
});