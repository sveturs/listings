import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import axios from '../../api/axios';
import CustomUIManager from './CustomUIManager';

// Мокаем модули
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      if (key === 'admin.customUI.title') return 'Custom UI Manager';
      if (key === 'admin.customUI.description') return 'Configure custom UI for this category';
      if (key === 'admin.customUI.enableCustomUI') return 'Enable custom UI';
      if (key === 'admin.customUI.selectComponent') return 'Select UI component';
      if (key === 'admin.customUI.noComponent') return 'No component';
      if (key === 'admin.customUI.componentDescription') return 'Select a custom UI component for this category';
      if (key === 'admin.customUI.helpTooltip') return 'Custom UI components provide specialized forms and views for specific category types';
      if (key === 'admin.customUI.fetchError') return 'Failed to fetch category UI settings';
      if (key === 'admin.customUI.updateError') return 'Failed to update category UI settings';
      if (key === 'common.save') return 'Save';
      return key;
    }
  })
}));

jest.mock('../../api/axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('CustomUIManager', () => {
  const mockCategoryId = 1;
  const mockOnCategoryUpdate = jest.fn();
  
  // Мокируем данные категории
  const categoryWithoutCustomUI = {
    id: 1,
    name: 'Electronics',
    slug: 'electronics',
    has_custom_ui: false,
    custom_ui_component: ''
  };
  
  const categoryWithCustomUI = {
    id: 1,
    name: 'Automobiles',
    slug: 'automobiles',
    has_custom_ui: true,
    custom_ui_component: 'AutoCategoryUI'
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders and loads category data without custom UI', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: categoryWithoutCustomUI });
    
    render(
      <CustomUIManager
        categoryId={mockCategoryId}
        onCategoryUpdate={mockOnCategoryUpdate}
      />
    );
    
    // Проверяем, что заголовок компонента отображается
    expect(screen.getByText('Custom UI Manager')).toBeInTheDocument();
    
    // Проверяем, что API был вызван для загрузки данных категории
    expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories/1');
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      // Проверяем, что чекбокс отображается и не выбран
      const checkbox = screen.getByLabelText('Enable custom UI') as HTMLInputElement;
      expect(checkbox).toBeInTheDocument();
      expect(checkbox.checked).toBe(false);
    });
    
    // Проверяем, что нет выпадающего списка для выбора компонента
    expect(screen.queryByLabelText('Select UI component')).not.toBeInTheDocument();
  });

  test('renders and loads category data with custom UI', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: categoryWithCustomUI });
    
    render(
      <CustomUIManager
        categoryId={mockCategoryId}
        onCategoryUpdate={mockOnCategoryUpdate}
      />
    );
    
    // Проверяем, что API был вызван для загрузки данных категории
    expect(mockedAxios.get).toHaveBeenCalledWith('/api/admin/categories/1');
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      // Проверяем, что чекбокс отображается и выбран
      const checkbox = screen.getByLabelText('Enable custom UI') as HTMLInputElement;
      expect(checkbox).toBeInTheDocument();
      expect(checkbox.checked).toBe(true);
      
      // Проверяем, что выпадающий список отображается
      expect(screen.getByLabelText('Select UI component')).toBeInTheDocument();
      
      // Проверяем, что выбранное значение корректно отображается
      expect(screen.getByText('Автомобили')).toBeInTheDocument();
    });
  });

  test('handles enabling and disabling custom UI', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: categoryWithoutCustomUI });
    
    render(
      <CustomUIManager
        categoryId={mockCategoryId}
        onCategoryUpdate={mockOnCategoryUpdate}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByLabelText('Enable custom UI')).toBeInTheDocument();
    });
    
    // Включаем кастомный UI
    const checkbox = screen.getByLabelText('Enable custom UI');
    fireEvent.click(checkbox);
    
    // Проверяем, что выпадающий список появился
    expect(screen.getByLabelText('Select UI component')).toBeInTheDocument();
    
    // Выбираем компонент из выпадающего списка
    const select = screen.getByLabelText('Select UI component');
    fireEvent.mouseDown(select);
    
    await waitFor(() => {
      // Выбираем компонент из списка
      const option = screen.getByText('Электроника');
      fireEvent.click(option);
    });
    
    // Сохраняем изменения
    const saveButton = screen.getByRole('button', { name: 'Save' });
    fireEvent.click(saveButton);
    
    // Проверяем, что был вызван API для обновления
    await waitFor(() => {
      expect(mockedAxios.patch).toHaveBeenCalledWith(
        '/api/admin/categories/1/custom-ui',
        {
          has_custom_ui: true,
          custom_ui_component: 'ElectronicsCategoryUI'
        }
      );
    });
    
    // Проверяем, что был вызван колбэк onCategoryUpdate
    expect(mockOnCategoryUpdate).toHaveBeenCalled();
    
    // Отключаем кастомный UI
    fireEvent.click(checkbox);
    
    // Проверяем, что выпадающий список исчез
    expect(screen.queryByLabelText('Select UI component')).not.toBeInTheDocument();
    
    // Сохраняем изменения снова
    fireEvent.click(saveButton);
    
    // Проверяем, что был вызван API для обновления
    await waitFor(() => {
      expect(mockedAxios.patch).toHaveBeenCalledWith(
        '/api/admin/categories/1/custom-ui',
        {
          has_custom_ui: false,
          custom_ui_component: ''
        }
      );
    });
  });

  test('handles API error when loading data', async () => {
    // Имитируем ошибку API при загрузке данных
    mockedAxios.get.mockRejectedValueOnce(new Error('API error'));
    
    render(
      <CustomUIManager
        categoryId={mockCategoryId}
        onCategoryUpdate={mockOnCategoryUpdate}
      />
    );
    
    // Должно отображаться сообщение об ошибке
    await waitFor(() => {
      expect(screen.getByText('Failed to fetch category UI settings')).toBeInTheDocument();
    });
  });

  test('handles API error when saving', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: categoryWithoutCustomUI });
    
    // Имитируем ошибку API при сохранении
    mockedAxios.patch.mockRejectedValueOnce(new Error('API save error'));
    
    render(
      <CustomUIManager
        categoryId={mockCategoryId}
        onCategoryUpdate={mockOnCategoryUpdate}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByLabelText('Enable custom UI')).toBeInTheDocument();
    });
    
    // Включаем кастомный UI
    const checkbox = screen.getByLabelText('Enable custom UI');
    fireEvent.click(checkbox);
    
    // Сохраняем изменения
    const saveButton = screen.getByRole('button', { name: 'Save' });
    fireEvent.click(saveButton);
    
    // Должно отображаться сообщение об ошибке
    await waitFor(() => {
      expect(screen.getByText('Failed to update category UI settings')).toBeInTheDocument();
    });
    
    // Колбэк onCategoryUpdate не должен быть вызван
    expect(mockOnCategoryUpdate).not.toHaveBeenCalled();
  });

  test('disables save button when no changes are made', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: categoryWithoutCustomUI });
    
    render(
      <CustomUIManager
        categoryId={mockCategoryId}
        onCategoryUpdate={mockOnCategoryUpdate}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByLabelText('Enable custom UI')).toBeInTheDocument();
    });
    
    // Проверяем, что кнопка сохранения изначально отключена
    const saveButton = screen.getByRole('button', { name: 'Save' }) as HTMLButtonElement;
    expect(saveButton).toBeDisabled();
    
    // Делаем изменения
    const checkbox = screen.getByLabelText('Enable custom UI');
    fireEvent.click(checkbox);
    
    // Проверяем, что кнопка сохранения стала активной
    expect(saveButton).not.toBeDisabled();
    
    // Возвращаем исходное состояние
    fireEvent.click(checkbox);
    
    // Проверяем, что кнопка сохранения снова отключена
    expect(saveButton).toBeDisabled();
  });

  test('shows all available UI components in the dropdown', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: categoryWithCustomUI });
    
    render(
      <CustomUIManager
        categoryId={mockCategoryId}
        onCategoryUpdate={mockOnCategoryUpdate}
      />
    );
    
    // Ждем, пока данные загрузятся
    await waitFor(() => {
      expect(screen.getByLabelText('Select UI component')).toBeInTheDocument();
    });
    
    // Открываем выпадающий список
    const select = screen.getByLabelText('Select UI component');
    fireEvent.mouseDown(select);
    
    // Проверяем, что все доступные компоненты отображаются
    await waitFor(() => {
      expect(screen.getByText('No component')).toBeInTheDocument();
      expect(screen.getByText('Автомобили')).toBeInTheDocument();
      expect(screen.getByText('Недвижимость')).toBeInTheDocument();
      expect(screen.getByText('Электроника')).toBeInTheDocument();
      expect(screen.getByText('Мебель')).toBeInTheDocument();
      expect(screen.getByText('Одежда')).toBeInTheDocument();
      expect(screen.getByText('Услуги')).toBeInTheDocument();
    });
  });
});