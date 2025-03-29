// frontend/hostel-frontend/src/components/marketplace/HierarchicalCategorySelect.jsx
// Исправляем компонент, чтобы сохранить правильную иерархию категорий

import React, { useState, useEffect } from 'react';
import { 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  Box,
  Typography,
  Divider
} from '@mui/material';
import { ChevronRight } from 'lucide-react';

/**
 * Компонент для иерархического выбора категорий, позволяющий выбирать категории любого уровня
 */
const HierarchicalCategorySelect = ({ 
  categories, 
  value, 
  onChange, 
  placeholder,
  label,
  error,
  size = "medium"
}) => {
  // Рекурсивная функция для рендеринга категорий с правильной иерархией
  const renderCategories = (categories, level = 0) => {
    return categories.map(category => (
      <React.Fragment key={category.id}>
        {/* Сама категория */}
        <MenuItem 
          value={category.id}
          sx={{ 
            pl: 2 + level * 2,
            borderLeft: level > 0 ? '1px dashed rgba(0,0,0,0.1)' : 'none',
            display: 'flex',
            alignItems: 'center'
          }}
        >
          {level > 0 && (
            <ChevronRight size={14 + level} style={{ marginRight: 8, opacity: 0.4 + (0.1 * level) }} />
          )}
          <Box>
            <Typography variant="body2">
              {category.name}
            </Typography>
            {category.parent_name && (
              <Typography variant="caption" color="text.secondary">
                {category.parent_name}
              </Typography>
            )}
          </Box>
        </MenuItem>
        
        {/* Дочерние категории, если они есть */}
        {category.children && category.children.length > 0 && 
          renderCategories(category.children.map(child => ({
            ...child,
            parent_name: category.name
          })), level + 1)
        }
      </React.Fragment>
    ));
  };
  
  // Функция для подготовки древовидной структуры категорий, если она не существует
  const prepareCategories = (cats) => {
    // Если категории уже имеют иерархическую структуру (с children)
    if (cats.some(c => c.children && c.children.length > 0)) {
      return cats;
    }
    
    // Иначе преобразуем из плоского списка в дерево
    const categoryMap = {};
    cats.forEach(cat => {
      categoryMap[cat.id] = { ...cat, children: [] };
    });
    
    const rootCategories = [];
    cats.forEach(cat => {
      if (cat.parent_id) {
        if (categoryMap[cat.parent_id]) {
          categoryMap[cat.parent_id].children.push(categoryMap[cat.id]);
        } else {
          rootCategories.push(categoryMap[cat.id]);
        }
      } else {
        rootCategories.push(categoryMap[cat.id]);
      }
    });
    
    return rootCategories;
  };
  
  // Подготавливаем категории
  const preparedCategories = prepareCategories(categories);
  
  return (
    <FormControl fullWidth error={error} size={size}>
      <InputLabel>{label || placeholder || 'Выберите категорию'}</InputLabel>
      <Select
        value={value || ''}
        onChange={(e) => onChange(e.target.value)}
        label={label || placeholder || 'Выберите категорию'}
        MenuProps={{
          PaperProps: {
            style: {
              maxHeight: 400,
            },
          },
        }}
      >
        {/* Специальная категория "Прочее" */}
        <MenuItem value={9999}>
          Прочее (ID: 9999)
        </MenuItem>
        
        <Divider />
        
        {/* Отображаем категории в иерархическом виде */}
        {renderCategories(preparedCategories)}
      </Select>
    </FormControl>
  );
};

export default HierarchicalCategorySelect;