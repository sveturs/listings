import React from 'react';
import ComponentRegistry, { CategoryUiComponentProps } from './ComponentRegistry';
import { Box, Typography } from '@mui/material';

/**
 * Компонент, который выводится, если кастомного компонента для категории не существует
 */
const DefaultCategoryComponent: React.FC<CategoryUiComponentProps> = ({ categoryId, values, onChange }) => {
  return (
    <Box sx={{ p: 2, border: '1px dashed #ccc', borderRadius: 1 }}>
      <Typography variant="body2" color="text.secondary">
        Стандартная форма атрибутов для категории #{categoryId}
      </Typography>
    </Box>
  );
};

/**
 * Фабрика компонентов для категорий
 * Выбирает подходящий компонент в зависимости от ID категории или компонента
 */
const CategoryComponentFactory: React.FC<{
  categoryId: number;
  componentName?: string;
  values: Record<string, any>;
  onChange: (values: Record<string, any>) => void;
}> = ({ categoryId, componentName, values, onChange }) => {
  // Если указано имя компонента, используем его
  if (componentName) {
    const CustomComponent = ComponentRegistry.getCategoryComponent(componentName);
    if (CustomComponent) {
      return <CustomComponent categoryId={categoryId} values={values} onChange={onChange} />;
    }
  }

  // Если компонент не найден, возвращаем стандартный компонент
  return <DefaultCategoryComponent categoryId={categoryId} values={values} onChange={onChange} />;
};

export default CategoryComponentFactory;

// Регистрируем стандартный компонент
ComponentRegistry.registerCategoryComponent('DefaultCategoryComponent', DefaultCategoryComponent);