import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import axios from '../../api/axios';
import EditAttributeInCategory from './EditAttributeInCategory';

// Мокаем модули
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      // Простая имитация функции перевода
      if (key === 'admin.categoryAttributes.editAttributeInCategory') return 'Edit Attribute in Category';
      if (key === 'admin.categoryAttributes.basicTab') return 'Basic';
      if (key === 'admin.categoryAttributes.descriptionsTab') return 'Descriptions';
      if (key === 'admin.categoryAttributes.optionsTab') return 'Options';
      if (key === 'admin.categoryAttributes.unitsTab') return 'Units';
      if (key === 'admin.attributes.required') return 'Required';
      if (key === 'admin.attributes.enabled') return 'Enabled';
      if (key === 'admin.attributes.sortOrder') return 'Sort Order';
      if (key === 'admin.categoryAttributes.hint') return 'Hint';
      if (key === 'admin.categoryAttributes.description') return 'Description';
      if (key === 'admin.categoryAttributes.defaultLanguage') return 'Default Language';
      if (key === 'admin.categoryAttributes.translations') return 'Translations';
      if (key === 'admin.categoryAttributes.hintHelp') return 'Short hint text for this attribute';
      if (key === 'admin.categoryAttributes.descriptionHelp') return 'Detailed description for this attribute';
      if (key === 'admin.categoryAttributes.fetchDetailsError') return 'Failed to fetch attribute details';
      if (key === 'admin.categoryAttributes.updateAttributeError') return 'Failed to update attribute';
      if (key === 'admin.categoryAttributes.customizeOptions') return 'Customize Options';
      if (key === 'admin.attributes.addOption') return 'Add Option';
      if (key === 'admin.attributes.optionValue') return 'Option Value';
      if (key === 'admin.categoryAttributes.noCustomOptions') return 'No custom options';
      if (key === 'admin.remove') return 'Remove';
      if (key === 'admin.categoryAttributes.optionTranslations') return 'Option Translations';
      if (key === 'admin.categoryAttributes.unit') return 'Unit';
      if (key === 'admin.categoryAttributes.unitHelp') return 'Unit for numeric values';
      if (key === 'admin.categoryAttributes.unitTranslations') return 'Unit Translations';
      if (key === 'common.cancel') return 'Cancel';
      if (key === 'common.save') return 'Save';
      if (key === 'languages.en') return 'English';
      if (key === 'languages.ru') return 'Russian';
      if (key === 'languages.sr') return 'Serbian';
      if (key.startsWith('admin.attributeTypes.')) {
        const type = key.split('.').pop();
        if (type === 'text') return 'Text';
        if (type === 'number') return 'Number';
        if (type === 'select') return 'Select';
        if (type === 'multiselect') return 'Multi-select';
        if (type === 'boolean') return 'Boolean';
        if (type === 'range') return 'Range';
        if (type === 'date') return 'Date';
        return type || '';
      }
      return key;
    },
    i18n: {
      language: 'en'
    }
  })
}));

jest.mock('../../api/axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('EditAttributeInCategory', () => {
  const mockCategoryId = 1;
  const mockAttributeId = 10;
  const mockOnClose = jest.fn();
  const mockOnUpdate = jest.fn();
  const mockOnError = jest.fn();
  
  // Mock data for text attribute
  const mockTextAttributeMapping = {
    category_id: 1,
    attribute_id: 10,
    is_required: true,
    is_enabled: true,
    sort_order: 1,
    hint: 'Enter model number',
    description: 'The model number of the product',
    translations: {
      ru: {
        hint: 'Введите номер модели',
        description: 'Номер модели продукта'
      }
    },
    attribute: {
      id: 10,
      name: 'model',
      display_name: 'Model',
      attribute_type: 'text',
      is_searchable: true,
      is_filterable: true,
      is_required: true,
      sort_order: 1,
      created_at: '2023-01-01'
    }
  };
  
  // Mock data for select attribute
  const mockSelectAttributeMapping = {
    category_id: 1,
    attribute_id: 20,
    is_required: false,
    is_enabled: true,
    sort_order: 2,
    attribute: {
      id: 20,
      name: 'color',
      display_name: 'Color',
      attribute_type: 'select',
      options: {
        values: ['Red', 'Blue', 'Green']
      },
      is_searchable: true,
      is_filterable: true,
      is_required: false,
      sort_order: 2,
      created_at: '2023-01-02'
    },
    options: {
      values: ['Red', 'Blue', 'Green', 'Yellow'],
      translations: {
        ru: {
          'Red': 'Красный',
          'Blue': 'Синий',
          'Green': 'Зеленый'
        }
      }
    }
  };
  
  // Mock data for number attribute
  const mockNumberAttributeMapping = {
    category_id: 1,
    attribute_id: 30,
    is_required: false,
    is_enabled: true,
    sort_order: 3,
    unit: 'kg',
    unit_translations: {
      ru: 'кг'
    },
    attribute: {
      id: 30,
      name: 'weight',
      display_name: 'Weight',
      attribute_type: 'number',
      is_searchable: true,
      is_filterable: true,
      is_required: false,
      sort_order: 3,
      created_at: '2023-01-03'
    }
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders loading state initially', async () => {
    // Задерживаем ответ API, чтобы протестировать состояние загрузки
    mockedAxios.get.mockImplementationOnce(() => new Promise(resolve => {
      setTimeout(() => {
        resolve({ data: mockTextAttributeMapping });
      }, 100);
    }));
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={mockAttributeId}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Должен отображаться заголовок диалога
    expect(screen.getByText('Edit Attribute in Category')).toBeInTheDocument();
    
    // Должен отображаться индикатор загрузки
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Model')).toBeInTheDocument();
    });
  });

  test('fetches and displays text attribute data', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockTextAttributeMapping });
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={mockAttributeId}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Проверяем, что был вызван API для загрузки данных
    expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories/1/attributes/10/details');
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      // Проверяем, что отображается атрибут
      expect(screen.getByText('Model')).toBeInTheDocument();
      expect(screen.getByText('model (Text)')).toBeInTheDocument();
      
      // Проверяем, что отображаются вкладки
      expect(screen.getByText('Basic')).toBeInTheDocument();
      expect(screen.getByText('Descriptions')).toBeInTheDocument();
      
      // Проверяем, что отображаются базовые настройки
      expect(screen.getByLabelText('Required')).toBeInTheDocument();
      expect(screen.getByLabelText('Enabled')).toBeInTheDocument();
      expect(screen.getByLabelText('Sort Order')).toBeInTheDocument();
    });
    
    // Проверяем значения полей
    const requiredCheckbox = screen.getByLabelText('Required') as HTMLInputElement;
    expect(requiredCheckbox.checked).toBe(true);
    
    const enabledCheckbox = screen.getByLabelText('Enabled') as HTMLInputElement;
    expect(enabledCheckbox.checked).toBe(true);
    
    const sortOrderInput = screen.getByLabelText('Sort Order') as HTMLInputElement;
    expect(sortOrderInput.value).toBe('1');
  });

  test('displays descriptions tab for text attribute', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockTextAttributeMapping });
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={mockAttributeId}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Model')).toBeInTheDocument();
    });
    
    // Переключаемся на вкладку Descriptions
    fireEvent.click(screen.getByText('Descriptions'));
    
    // Проверяем, что отображаются поля для описаний
    expect(screen.getByLabelText('Hint')).toBeInTheDocument();
    expect(screen.getByLabelText('Description')).toBeInTheDocument();
    
    // Проверяем значения полей
    const hintInput = screen.getByLabelText('Hint') as HTMLInputElement;
    expect(hintInput.value).toBe('Enter model number');
    
    const descriptionInput = screen.getByLabelText('Description') as HTMLInputElement;
    expect(descriptionInput.value).toBe('The model number of the product');
    
    // Проверяем, что отображается аккордеон с переводами
    expect(screen.getByText('Translations')).toBeInTheDocument();
  });

  test('displays options tab for select attribute', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockSelectAttributeMapping });
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={20} // ID select атрибута
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Проверяем, что есть вкладка Options
    expect(screen.getByText('Options')).toBeInTheDocument();
    
    // Переключаемся на вкладку Options
    fireEvent.click(screen.getByText('Options'));
    
    // Проверяем, что отображаются опции
    expect(screen.getByText('Customize Options')).toBeInTheDocument();
    expect(screen.getByText('Add Option')).toBeInTheDocument();
    
    // Проверяем, что отображаются значения опций
    const optionInputs = screen.getAllByRole('textbox');
    const optionValues = optionInputs.map(input => (input as HTMLInputElement).value);
    
    // Проверяем, что все опции отображаются
    expect(optionValues).toContain('Red');
    expect(optionValues).toContain('Blue');
    expect(optionValues).toContain('Green');
    expect(optionValues).toContain('Yellow');
    
    // Проверяем, что есть кнопки удаления опций
    expect(screen.getAllByTitle('Remove').length).toBe(4); // 4 опции
  });

  test('displays units tab for number attribute', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockNumberAttributeMapping });
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={30} // ID number атрибута
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Weight')).toBeInTheDocument();
    });
    
    // Проверяем, что есть вкладка Units
    expect(screen.getByText('Units')).toBeInTheDocument();
    
    // Переключаемся на вкладку Units
    fireEvent.click(screen.getByText('Units'));
    
    // Проверяем, что отображается поле для единицы измерения
    expect(screen.getByLabelText('Unit')).toBeInTheDocument();
    
    // Проверяем значение поля
    const unitInput = screen.getByLabelText('Unit') as HTMLInputElement;
    expect(unitInput.value).toBe('kg');
    
    // Проверяем, что отображается аккордеон с переводами единиц измерения
    expect(screen.getByText('Unit Translations')).toBeInTheDocument();
  });

  test('handles form input changes', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockTextAttributeMapping });
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={mockAttributeId}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Model')).toBeInTheDocument();
    });
    
    // Изменяем значения в базовых настройках
    const requiredCheckbox = screen.getByLabelText('Required');
    fireEvent.click(requiredCheckbox); // Должно стать false
    
    const sortOrderInput = screen.getByLabelText('Sort Order');
    fireEvent.change(sortOrderInput, { target: { value: '5' } });
    
    // Переключаемся на вкладку Descriptions
    fireEvent.click(screen.getByText('Descriptions'));
    
    // Изменяем значения в описаниях
    const hintInput = screen.getByLabelText('Hint');
    fireEvent.change(hintInput, { target: { value: 'Updated hint' } });
    
    const descriptionInput = screen.getByLabelText('Description');
    fireEvent.change(descriptionInput, { target: { value: 'Updated description' } });
    
    // Сохраняем изменения
    const saveButton = screen.getByRole('button', { name: 'Save' });
    fireEvent.click(saveButton);
    
    // Проверяем, что был вызван API для обновления
    await waitFor(() => {
      expect(mockedAxios.put).toHaveBeenCalledWith(
        '/api/admin/categories/1/attributes/10',
        expect.objectContaining({
          is_required: false, // Изменено с true на false
          is_enabled: true,
          sort_order: 5, // Изменено с 1 на 5
          hint: 'Updated hint', // Изменено
          description: 'Updated description', // Изменено
        })
      );
    });
    
    // Проверяем, что были вызваны колбэки
    expect(mockOnUpdate).toHaveBeenCalled();
    expect(mockOnClose).toHaveBeenCalled();
  });

  test('handles option changes for select attribute', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockSelectAttributeMapping });
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={20} // ID select атрибута
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Color')).toBeInTheDocument();
    });
    
    // Переключаемся на вкладку Options
    fireEvent.click(screen.getByText('Options'));
    
    // Добавляем новую опцию
    const addButton = screen.getByText('Add Option');
    fireEvent.click(addButton);
    
    // Должен появиться новый инпут
    const optionInputs = screen.getAllByRole('textbox');
    expect(optionInputs.length).toBe(5); // 4 существующих + 1 новый
    
    // Изменяем значение последней (новой) опции
    fireEvent.change(optionInputs[optionInputs.length - 1], { target: { value: 'Purple' } });
    
    // Удаляем одну из опций
    const removeButtons = screen.getAllByTitle('Remove');
    fireEvent.click(removeButtons[0]); // Удаляем первую опцию (Red)
    
    // Сохраняем изменения
    const saveButton = screen.getByRole('button', { name: 'Save' });
    fireEvent.click(saveButton);
    
    // Проверяем, что был вызван API для обновления
    await waitFor(() => {
      expect(mockedAxios.put).toHaveBeenCalledWith(
        '/api/admin/categories/1/attributes/20',
        expect.objectContaining({
          options: expect.objectContaining({
            values: ['Blue', 'Green', 'Yellow', 'Purple'] // Red удален, Purple добавлен
          })
        })
      );
    });
  });

  test('handles API error when loading data', async () => {
    // Имитируем ошибку API при загрузке данных
    mockedAxios.get.mockRejectedValueOnce(new Error('API error'));
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={mockAttributeId}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Должно отображаться сообщение об ошибке
    await waitFor(() => {
      expect(screen.getByText('Failed to fetch attribute details')).toBeInTheDocument();
    });
    
    // Должен быть вызван колбэк onError
    expect(mockOnError).toHaveBeenCalledWith('Failed to fetch attribute details');
  });

  test('handles API error when saving', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockTextAttributeMapping });
    
    // Имитируем ошибку API при сохранении
    mockedAxios.put.mockRejectedValueOnce(new Error('API save error'));
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={mockAttributeId}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Model')).toBeInTheDocument();
    });
    
    // Нажимаем кнопку сохранения
    const saveButton = screen.getByRole('button', { name: 'Save' });
    fireEvent.click(saveButton);
    
    // Должно отображаться сообщение об ошибке
    await waitFor(() => {
      expect(screen.getByText('Failed to update attribute')).toBeInTheDocument();
    });
    
    // Должен быть вызван колбэк onError
    expect(mockOnError).toHaveBeenCalledWith('Failed to update attribute');
    
    // Колбэки onUpdate и onClose не должны быть вызваны
    expect(mockOnUpdate).not.toHaveBeenCalled();
    expect(mockOnClose).not.toHaveBeenCalled();
  });

  test('closes dialog when cancel button is clicked', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockTextAttributeMapping });
    
    render(
      <EditAttributeInCategory
        open={true}
        categoryId={mockCategoryId}
        attributeId={mockAttributeId}
        onClose={mockOnClose}
        onUpdate={mockOnUpdate}
        onError={mockOnError}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByText('Model')).toBeInTheDocument();
    });
    
    // Нажимаем кнопку отмены
    const cancelButton = screen.getByRole('button', { name: 'Cancel' });
    fireEvent.click(cancelButton);
    
    // Должен быть вызван колбэк onClose
    expect(mockOnClose).toHaveBeenCalled();
    
    // Колбэк onUpdate не должен быть вызван
    expect(mockOnUpdate).not.toHaveBeenCalled();
  });
});