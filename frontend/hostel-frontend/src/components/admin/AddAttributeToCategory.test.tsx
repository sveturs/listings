import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import axios from '../../api/axios';
import AddAttributeToCategory from './AddAttributeToCategory';

// Мокаем модули
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      if (key === 'admin.categoryAttributes.addAttribute') return 'Add Attribute';
      if (key === 'admin.categoryAttributes.selectAttribute') return 'Select Attribute';
      if (key === 'admin.attributes.required') return 'Required';
      if (key === 'admin.attributes.enabled') return 'Enabled';
      if (key === 'admin.attributes.sortOrder') return 'Sort Order';
      if (key === 'admin.categoryAttributes.addAttributeButton') return 'Add';
      if (key === 'admin.categoryAttributes.fetchAttributesError') return 'Failed to fetch attributes';
      if (key === 'admin.categoryAttributes.addAttributeError') return 'Failed to add attribute';
      if (key === 'admin.categoryAttributes.noAttributesAvailable') return 'No attributes available';
      if (key === 'admin.categoryAttributes.allAttributesMapped') return 'All attributes are already mapped';
      if (key.startsWith('admin.attributeTypes.')) return key.split('.').pop() || '';
      return key;
    }
  })
}));

jest.mock('../../api/axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('AddAttributeToCategory', () => {
  const mockCategoryId = 1;
  const mockOnAttributeAdded = jest.fn();
  const mockOnError = jest.fn();
  
  const mockAllAttributes = [
    { id: 1, name: 'color', display_name: 'Color', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 1, created_at: '2023-01-01' },
    { id: 2, name: 'size', display_name: 'Size', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 2, created_at: '2023-01-02' },
    { id: 3, name: 'brand', display_name: 'Brand', attribute_type: 'text', is_searchable: true, is_filterable: false, is_required: false, sort_order: 3, created_at: '2023-01-03' }
  ];
  
  const mockExistingMappings = [
    { category_id: 1, attribute_id: 1, is_required: true, is_enabled: true, sort_order: 10 }
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    
    // Устанавливаем успешные ответы API по умолчанию
    mockedAxios.get.mockImplementation((url) => {
      if (url === '/api/admin/attributes') {
        return Promise.resolve({ data: mockAllAttributes });
      }
      if (url === '/api/admin/categories/1/attributes') {
        return Promise.resolve({ data: mockExistingMappings });
      }
      return Promise.reject(new Error('not found'));
    });
    
    mockedAxios.post.mockResolvedValue({ data: { success: true } });
  });

  test('renders and loads data correctly', async () => {
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Проверяем, что заголовок компонента отображается
    expect(screen.getByText('Add Attribute')).toBeInTheDocument();
    
    // Проверяем, что API был вызван для загрузки атрибутов и маппингов
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/attributes');
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories/1/attributes');
    });
    
    // Проверяем, что поле выбора атрибута отображается
    expect(screen.getByRole('combobox')).toBeInTheDocument();
    expect(screen.getByLabelText('Select Attribute')).toBeInTheDocument();
    
    // Проверяем, что остальные элементы формы отображаются
    expect(screen.getByLabelText('Required')).toBeInTheDocument();
    expect(screen.getByLabelText('Enabled')).toBeInTheDocument();
    expect(screen.getByLabelText('Sort Order')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Add' })).toBeInTheDocument();
  });

  test('filters out already mapped attributes', async () => {
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Ждем загрузки данных
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/attributes');
    });

    // Открываем выпадающий список атрибутов
    const autocomplete = screen.getByLabelText('Select Attribute');
    fireEvent.mouseDown(autocomplete);
    
    await waitFor(() => {
      // В выпадающем списке не должно быть атрибута 'Color', так как он уже привязан
      const options = screen.getAllByRole('option');
      const optionTexts = options.map(option => option.textContent);
      expect(optionTexts).not.toContain('Color');
      
      // Но должны быть атрибуты 'Size' и 'Brand', которые еще не привязаны
      expect(optionTexts).toContain('Sizeselect');
      expect(optionTexts).toContain('Brandtext');
    });
  });

  test('displays message when all attributes are mapped', async () => {
    // Имитируем ситуацию, когда все атрибуты уже привязаны
    mockedAxios.get.mockImplementation((url) => {
      if (url === '/api/admin/attributes') {
        return Promise.resolve({ data: [mockAllAttributes[0]] }); // только один атрибут
      }
      if (url === '/api/admin/categories/1/attributes') {
        return Promise.resolve({ data: mockExistingMappings }); // и он уже привязан
      }
      return Promise.reject(new Error('not found'));
    });
    
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Ждем загрузки данных
    await waitFor(() => {
      // Должно отображаться сообщение о том, что все атрибуты уже привязаны
      expect(screen.getByText('All attributes are already mapped')).toBeInTheDocument();
    });
  });

  test('adds an attribute to the category', async () => {
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Ждем загрузки данных
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/attributes');
    });
    
    // Открываем выпадающий список атрибутов
    const autocomplete = screen.getByLabelText('Select Attribute');
    fireEvent.mouseDown(autocomplete);
    
    // Выбираем атрибут 'Size'
    await waitFor(() => {
      const sizeOption = screen.getByText('Sizeselect');
      fireEvent.click(sizeOption);
    });
    
    // Включаем флаг "обязательный"
    const requiredCheckbox = screen.getByLabelText('Required');
    fireEvent.click(requiredCheckbox);
    
    // Изменяем порядок сортировки
    const sortOrderInput = screen.getByLabelText('Sort Order');
    fireEvent.change(sortOrderInput, { target: { value: '25' } });
    
    // Нажимаем кнопку добавления
    const addButton = screen.getByRole('button', { name: 'Add' });
    fireEvent.click(addButton);
    
    // Проверяем, что был вызван API для добавления атрибута
    await waitFor(() => {
      expect(mockedAxios.post).toHaveBeenCalledWith(
        '/api/admin/categories/1/attributes',
        {
          attribute_id: 2, // ID атрибута 'Size'
          is_required: true,
          is_enabled: true, // По умолчанию включено
          sort_order: 25
        }
      );
    });
    
    // Проверяем, что был вызван колбэк onAttributeAdded
    expect(mockOnAttributeAdded).toHaveBeenCalled();
    
    // Проверяем, что данные были перезагружены
    expect(mockedAxios.get).toHaveBeenCalledTimes(4); // 2 вызова при инициализации + 2 после добавления
  });

  test('handles API error when loading attributes', async () => {
    // Имитируем ошибку при загрузке атрибутов
    mockedAxios.get.mockRejectedValueOnce(new Error('API error'));
    
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Ждем загрузки данных и проверяем, что отображается сообщение об ошибке
    await waitFor(() => {
      expect(screen.getByText('Failed to fetch attributes')).toBeInTheDocument();
    });
    
    // Проверяем, что был вызван колбэк onError
    expect(mockOnError).toHaveBeenCalledWith('Failed to fetch attributes');
  });

  test('handles API error when adding an attribute', async () => {
    // Имитируем ошибку при добавлении атрибута
    mockedAxios.post.mockRejectedValueOnce(new Error('API error'));
    
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Ждем загрузки данных
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/attributes');
    });
    
    // Открываем выпадающий список атрибутов и выбираем атрибут
    const autocomplete = screen.getByLabelText('Select Attribute');
    fireEvent.mouseDown(autocomplete);
    
    await waitFor(() => {
      const sizeOption = screen.getByText('Sizeselect');
      fireEvent.click(sizeOption);
    });
    
    // Нажимаем кнопку добавления
    const addButton = screen.getByRole('button', { name: 'Add' });
    fireEvent.click(addButton);
    
    // Проверяем, что отображается сообщение об ошибке
    await waitFor(() => {
      expect(screen.getByText('Failed to add attribute')).toBeInTheDocument();
    });
    
    // Проверяем, что был вызван колбэк onError
    expect(mockOnError).toHaveBeenCalledWith('Failed to add attribute');
  });

  test('calculates next sort order correctly', async () => {
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Ждем загрузки данных
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories/1/attributes');
    });
    
    // Проверяем, что значение сортировки установлено на 10 больше, чем максимальное существующее
    // В нашем случае существующее значение = 10, поэтому должно быть 20
    const sortOrderInput = screen.getByLabelText('Sort Order') as HTMLInputElement;
    expect(sortOrderInput.value).toBe('20');
    
    // Теперь имитируем ситуацию, когда нет существующих маппингов
    mockedAxios.get.mockImplementation((url) => {
      if (url === '/api/admin/attributes') {
        return Promise.resolve({ data: mockAllAttributes });
      }
      if (url === '/api/admin/categories/1/attributes') {
        return Promise.resolve({ data: [] }); // Пустой массив маппингов
      }
      return Promise.reject(new Error('not found'));
    });
    
    // Перерендерим компонент
    render(
      <AddAttributeToCategory 
        categoryId={mockCategoryId}
        onAttributeAdded={mockOnAttributeAdded}
        onError={mockOnError}
      />
    );
    
    // Ждем загрузки данных
    await waitFor(() => {
      // Берем второй элемент, так как у нас теперь два инпута от двух рендеров
      const sortOrderInputs = screen.getAllByLabelText('Sort Order') as HTMLInputElement[];
      // Проверяем, что значение сортировки установлено на 10, так как нет существующих маппингов
      expect(sortOrderInputs[1].value).toBe('10');
    });
  });
});