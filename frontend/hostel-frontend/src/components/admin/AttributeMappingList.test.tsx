import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import axios from '../../api/axios';
import AttributeMappingList from './AttributeMappingList';

// Мокаем модули
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      if (key === 'admin.attributes.name') return 'Name';
      if (key === 'admin.attributes.type') return 'Type';
      if (key === 'admin.attributes.required') return 'Required';
      if (key === 'admin.attributes.enabled') return 'Enabled';
      if (key === 'admin.attributes.sortOrder') return 'Sort Order';
      if (key === 'admin.actions') return 'Actions';
      if (key === 'admin.edit') return 'Edit';
      if (key === 'admin.remove') return 'Remove';
      if (key === 'admin.categoryAttributes.noAttributesMapped') return 'No attributes mapped to this category';
      if (key === 'admin.categoryAttributes.advancedEdit') return 'Advanced Edit';
      if (key === 'admin.categoryAttributes.confirmRemove') return 'Are you sure you want to remove this attribute?';
      if (key === 'admin.categoryAttributes.fetchMappingsError') return 'Failed to fetch attribute mappings';
      if (key === 'admin.categoryAttributes.updateAttributeError') return 'Failed to update attribute';
      if (key === 'admin.categoryAttributes.removeAttributeError') return 'Failed to remove attribute';
      return key;
    }
  })
}));

jest.mock('../../api/axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

// Мокаем диалог расширенного редактирования
jest.mock('./EditAttributeInCategory', () => {
  return ({ 
    open, 
    categoryId, 
    attributeId, 
    onClose, 
    onUpdate, 
    onError 
  }: any) => {
    return (
      <div data-testid="edit-attribute-dialog">
        {open && (
          <div>
            <div>Category ID: {categoryId}</div>
            <div>Attribute ID: {attributeId}</div>
            <button onClick={onClose}>Close</button>
            <button onClick={onUpdate}>Update</button>
          </div>
        )}
      </div>
    );
  };
});

// Мокируем confirm браузера
window.confirm = jest.fn();

describe('AttributeMappingList', () => {
  const mockCategoryId = 1;
  const mockOnError = jest.fn();
  
  const mockAttributeMappings = [
    {
      category_id: 1,
      attribute_id: 10,
      is_required: true,
      is_enabled: true,
      sort_order: 1,
      attribute: {
        id: 10,
        name: 'color',
        display_name: 'Color',
        attribute_type: 'select',
        is_searchable: true,
        is_filterable: true,
        is_required: true,
        sort_order: 1,
        created_at: '2023-01-01'
      }
    },
    {
      category_id: 1,
      attribute_id: 20,
      is_required: false,
      is_enabled: true,
      sort_order: 2,
      attribute: {
        id: 20,
        name: 'size',
        display_name: 'Size',
        attribute_type: 'select',
        is_searchable: true,
        is_filterable: true,
        is_required: false,
        sort_order: 2,
        created_at: '2023-01-02'
      }
    },
    {
      category_id: 1,
      attribute_id: 30,
      is_required: false,
      is_enabled: false,
      sort_order: 3,
      attribute: {
        id: 30,
        name: 'brand',
        display_name: 'Brand',
        attribute_type: 'text',
        is_searchable: true,
        is_filterable: false,
        is_required: false,
        sort_order: 3,
        created_at: '2023-01-03'
      }
    }
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    
    // Устанавливаем успешный ответ API по умолчанию
    mockedAxios.get.mockResolvedValue({ data: mockAttributeMappings });
    mockedAxios.put.mockResolvedValue({ data: { success: true } });
    mockedAxios.delete.mockResolvedValue({ data: { success: true } });
    
    // Сбрасываем мок для window.confirm и заставляем его возвращать true по умолчанию
    (window.confirm as jest.Mock).mockReturnValue(true);
  });

  test('renders loading state initially', async () => {
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    // Должен отображаться индикатор загрузки
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
    
    // После загрузки данных таблица должна появиться
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
  });

  test('renders attribute mappings correctly', async () => {
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      // Проверяем, что в таблице отображаются заголовки
      expect(screen.getByText('Name')).toBeInTheDocument();
      expect(screen.getByText('Type')).toBeInTheDocument();
      expect(screen.getByText('Required')).toBeInTheDocument();
      expect(screen.getByText('Enabled')).toBeInTheDocument();
      expect(screen.getByText('Sort Order')).toBeInTheDocument();
      
      // Проверяем, что отображаются атрибуты
      expect(screen.getByText('Color')).toBeInTheDocument();
      expect(screen.getByText('Size')).toBeInTheDocument();
      expect(screen.getByText('Brand')).toBeInTheDocument();
      
      // Проверяем, что отображаются типы атрибутов
      expect(screen.getAllByText('select').length).toBe(2);
      expect(screen.getByText('text')).toBeInTheDocument();
    });
  });

  test('displays message when there are no attribute mappings', async () => {
    // Устанавливаем пустой список атрибутов
    mockedAxios.get.mockResolvedValueOnce({ data: [] });
    
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      expect(screen.getByText('No attributes mapped to this category')).toBeInTheDocument();
    });
  });

  test('handles API error', async () => {
    // Имитируем ошибку API
    mockedAxios.get.mockRejectedValueOnce(new Error('API error'));
    
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      // Должно отображаться сообщение об ошибке
      expect(screen.getByText('Failed to fetch attribute mappings')).toBeInTheDocument();
      
      // Должен быть вызван колбэк onError
      expect(mockOnError).toHaveBeenCalledWith('Failed to fetch attribute mappings');
    });
  });

  test('allows editing an attribute mapping', async () => {
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Находим кнопку редактирования и нажимаем на неё
    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);
    
    // Проверяем, что включился режим редактирования
    // (поиск поля ввода для sort_order)
    expect(screen.getByRole('spinbutton')).toBeInTheDocument();
    
    // Изменяем значение sort_order
    fireEvent.change(screen.getByRole('spinbutton'), { target: { value: '5' } });
    
    // Нажимаем кнопку сохранения
    const saveButton = screen.getByRole('button', { name: '' });  // У кнопки нет текстового содержимого, только иконка
    fireEvent.click(saveButton);
    
    // Проверяем, что был вызван API для обновления
    await waitFor(() => {
      expect(mockedAxios.put).toHaveBeenCalledWith(
        '/api/admin/categories/1/attributes/10',
        {
          is_required: true,
          is_enabled: true,
          sort_order: 5
        }
      );
    });
    
    // Проверяем, что после сохранения снова загружается список атрибутов
    expect(mockedAxios.get).toHaveBeenCalledTimes(2);
    expect(mockedAxios.get).toHaveBeenLastCalledWith('/api/admin/categories/1/attributes');
  });

  test('can cancel editing', async () => {
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Находим кнопку редактирования и нажимаем на неё
    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);
    
    // Проверяем, что включился режим редактирования
    expect(screen.getByRole('spinbutton')).toBeInTheDocument();
    
    // Находим кнопку отмены и нажимаем на неё
    const cancelButtons = screen.getAllByRole('button');
    const cancelButton = cancelButtons.find(btn => btn.querySelector('svg[data-testid="CancelIcon"]'));
    if (cancelButton) {
      fireEvent.click(cancelButton);
    }
    
    // Проверяем, что режим редактирования выключился
    await waitFor(() => {
      expect(screen.queryByRole('spinbutton')).not.toBeInTheDocument();
    });
    
    // Проверяем, что API не вызывался для обновления
    expect(mockedAxios.put).not.toHaveBeenCalled();
  });

  test('opens advanced edit dialog', async () => {
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Находим кнопку расширенного редактирования и нажимаем на неё
    const advancedEditButtons = screen.getAllByTitle('Advanced Edit');
    fireEvent.click(advancedEditButtons[0]);
    
    // Проверяем, что диалог открылся и передаются правильные параметры
    expect(screen.getByTestId('edit-attribute-dialog')).toBeInTheDocument();
    expect(screen.getByText('Category ID: 1')).toBeInTheDocument();
    expect(screen.getByText('Attribute ID: 10')).toBeInTheDocument();
    
    // Закрываем диалог
    fireEvent.click(screen.getByText('Close'));
    
    // После закрытия диалог должен исчезнуть
    await waitFor(() => {
      expect(screen.queryByText('Category ID: 1')).not.toBeInTheDocument();
    });
  });

  test('confirms and removes an attribute', async () => {
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Находим кнопку удаления и нажимаем на неё
    const removeButtons = screen.getAllByTitle('Remove');
    fireEvent.click(removeButtons[0]);
    
    // Проверяем, что был вызван confirm
    expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to remove this attribute?');
    
    // Проверяем, что был вызван API для удаления
    await waitFor(() => {
      expect(mockedAxios.delete).toHaveBeenCalledWith('/api/admin/categories/1/attributes/10');
    });
    
    // Проверяем, что после удаления снова загружается список атрибутов
    expect(mockedAxios.get).toHaveBeenCalledTimes(2);
    expect(mockedAxios.get).toHaveBeenLastCalledWith('/api/admin/categories/1/attributes');
  });

  test('does not remove attribute when confirmation is canceled', async () => {
    // Устанавливаем, что пользователь нажал Отмена в диалоге подтверждения
    (window.confirm as jest.Mock).mockReturnValue(false);
    
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Находим кнопку удаления и нажимаем на неё
    const removeButtons = screen.getAllByTitle('Remove');
    fireEvent.click(removeButtons[0]);
    
    // Проверяем, что был вызван confirm
    expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to remove this attribute?');
    
    // Проверяем, что API для удаления НЕ вызывался
    expect(mockedAxios.delete).not.toHaveBeenCalled();
  });

  test('handles attribute removal API error', async () => {
    // Имитируем ошибку API при удалении
    mockedAxios.delete.mockRejectedValueOnce(new Error('API deletion error'));
    
    render(<AttributeMappingList categoryId={mockCategoryId} onError={mockOnError} />);
    
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Находим кнопку удаления и нажимаем на неё
    const removeButtons = screen.getAllByTitle('Remove');
    fireEvent.click(removeButtons[0]);
    
    // Проверяем, что был вызван API для удаления
    await waitFor(() => {
      expect(mockedAxios.delete).toHaveBeenCalledWith('/api/admin/categories/1/attributes/10');
    });
    
    // Должно отображаться сообщение об ошибке
    await waitFor(() => {
      expect(screen.getByText('Failed to remove attribute')).toBeInTheDocument();
    });
    
    // Должен быть вызван колбэк onError
    expect(mockOnError).toHaveBeenCalledWith('Failed to remove attribute');
  });
});