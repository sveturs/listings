// frontend/hostel-frontend/src/components/marketplace/HierarchicalCategorySelect.jsx
import React, { useState, useEffect } from 'react';
import { 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  Box,
  Typography,
  Divider,
  IconButton,
  ListSubheader
} from '@mui/material';
import { ChevronRight, ChevronDown } from 'lucide-react';

/**
 * Компонент для иерархического выбора категорий, позволяющий выбирать категории любого уровня
 * с возможностью сворачивания/разворачивания категорий и сохранением состояния
 */
const HierarchicalCategorySelect = ({ 
  categories, 
  value, 
  onChange, 
  placeholder,
  label,
  error,
  size = "medium",
  expansionState = {},
  onExpansionChange
}) => {
  // Состояние для хранения развернутых категорий
  const [expandedCategories, setExpandedCategories] = useState(expansionState);

  // Эффект для синхронизации с внешним состоянием, если оно передано
  useEffect(() => {
    if (Object.keys(expansionState).length > 0) {
      setExpandedCategories(expansionState);
    }
  }, [expansionState]);

  // Обновляем внешнее состояние при изменении локального
  useEffect(() => {
    if (onExpansionChange) {
      onExpansionChange(expandedCategories);
    }
  }, [expandedCategories, onExpansionChange]);

  // Функция для переключения состояния категории (свернута/развернута)
  const toggleCategory = (categoryId, event) => {
    event.stopPropagation(); // Предотвращаем всплытие события, чтобы не выбрать категорию при клике на иконку
    setExpandedCategories(prev => ({
      ...prev,
      [categoryId]: !prev[categoryId]
    }));
  };

  // Функция для переключения всех категорий верхнего уровня
  const toggleAllTopCategories = (expandAll, event) => {
    if (event) {
      event.stopPropagation();
    }
    
    const topLevelIds = prepareCategories(categories)
      .filter(cat => cat.children && cat.children.length > 0)
      .map(cat => cat.id);
    
    const newExpandedState = {};
    topLevelIds.forEach(id => {
      newExpandedState[id] = expandAll;
    });
    
    setExpandedCategories(prev => ({
      ...prev,
      ...newExpandedState
    }));
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
  
  // Рекурсивно формируем плоский список MenuItems для Select
  // Вместо рендеринга с React.Fragment, создаем плоский массив MenuItem
  const generateMenuItems = (categories, level = 0) => {
    let items = [];
    
    categories.forEach(category => {
      // Добавляем саму категорию
      items.push(
        <MenuItem 
          key={`cat-${category.id}`}
          value={category.id}
          sx={{ 
            pl: 2 + level * 2,
            borderLeft: level > 0 ? '1px dashed rgba(0,0,0,0.1)' : 'none',
            display: 'flex',
            alignItems: 'center',
            position: 'relative'
          }}
        >
          {/* Иконка разворачивания/сворачивания только для категорий с дочерними элементами */}
          {category.children && category.children.length > 0 && (
            <Box 
              sx={{ 
                position: 'absolute', 
                left: level === 0 ? 0 : level * 1.5,
                display: 'flex',
                alignItems: 'center',
                cursor: 'pointer',
                mr: 1,
                zIndex: 1 // Убедимся, что иконка кликабельна
              }}
              onClick={(e) => toggleCategory(category.id, e)}
            >
              {expandedCategories[category.id] ? 
                <ChevronDown size={14} /> : 
                <ChevronRight size={14} />
              }
            </Box>
          )}
          
          {level > 0 && (
            <Box sx={{ width: 14, height: 14, mr: 1 }} /> // Пустое место для отступа у подкатегорий
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
      );
      
      // Если у этой категории есть дочерние элементы и она развернута,
      // рекурсивно добавляем и их
      if (
        category.children && 
        category.children.length > 0 && 
        expandedCategories[category.id]
      ) {
        // Добавляем дочерние элементы с информацией о родителе
        const childItems = generateMenuItems(
          category.children.map(child => ({
            ...child,
            parent_name: category.name
          })),
          level + 1
        );
        
        // Добавляем в общий массив
        items = [...items, ...childItems];
      }
    });
    
    return items;
  };

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
        
        {/* Заголовок с кнопками управления */}
        <ListSubheader sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="caption" color="text.secondary">
            Категории
          </Typography>
          <Box>
            <IconButton 
              size="small" 
              onClick={(e) => toggleAllTopCategories(true, e)}
              aria-label="Развернуть все"
            >
              <ChevronDown size={14} />
            </IconButton>
            <IconButton 
              size="small" 
              onClick={(e) => toggleAllTopCategories(false, e)}
              aria-label="Свернуть все"
            >
              <ChevronRight size={14} />
            </IconButton>
          </Box>
        </ListSubheader>
        
        {/* Генерируем плоский список элементов меню вместо использования React.Fragment */}
        {generateMenuItems(preparedCategories)}
      </Select>
    </FormControl>
  );
};

export default HierarchicalCategorySelect;