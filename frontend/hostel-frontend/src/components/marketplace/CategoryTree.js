//frontend/hostel-frontend/src/components/marketplace/CategoryTree.js
import React, { useState } from 'react';
import {
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Collapse,
  Typography,
  Box,
  useTheme
} from '@mui/material';
import { ChevronRight, ChevronDown } from 'lucide-react';

const CategoryTreeItem = ({ 
  category, 
  selectedId, 
  onSelect,
  level = 0 
}) => {
  // Меняем начальное состояние на false
  const [isOpen, setIsOpen] = useState(false);
  const theme = useTheme();

  // Добавляем функцию для подсчета всех объявлений в категории и подкатегориях
  const getTotalListings = (cat) => {
    let total = cat.listing_count || 0;
    if (cat.children) {
      total += cat.children.reduce((sum, child) => sum + getTotalListings(child), 0);
    }
    return total;
  };

  const handleClick = () => {
    if (hasChildren) {
      setIsOpen(!isOpen);
      onSelect(category.id); // Выбираем категорию при клике
    } else {
      onSelect(category.id);
    }
  };

  const totalListings = getTotalListings(category);
  const hasChildren = category.children && category.children.length > 0;
  const isSelected = selectedId === category.id;

  return (
    <>
      <ListItemButton
        dense
        onClick={() => {
          if (hasChildren) {
            setIsOpen(!isOpen);
          }
          onSelect(category.id);
        }}
        selected={isSelected}
        sx={{
          pl: level * 1.5 + 1,
          py: 0.5,
          minHeight: 32,
          color: isSelected ? 'primary.main' : 'text.primary',
          '&.Mui-selected': {
            backgroundColor: theme.palette.primary.main + '08',
            '&:hover': {
              backgroundColor: theme.palette.primary.main + '12',
            }
          }
        }}
      >
        <ListItemIcon sx={{ minWidth: 24, color: 'inherit' }}>
          {hasChildren && (
            isOpen ? <ChevronDown size={16} /> : <ChevronRight size={16} />
          )}
        </ListItemIcon>

        <ListItemText
          primary={
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <Typography
                variant="body2"
                sx={{
                  fontSize: '0.8125rem',
                  fontWeight: isSelected ? 500 : 400,
                }}
              >
                {category.name}
              </Typography>
              {totalListings > 0 && (
                <Typography
                  variant="caption"
                  color={isSelected ? 'primary' : 'text.secondary'}
                  sx={{ fontSize: '0.75rem' }}
                >
                  {totalListings}
                </Typography>
              )}
            </Box>
          }
        />
      </ListItemButton>

      {hasChildren && (
        <Collapse in={isOpen} timeout="auto">
          <List disablePadding>
            {category.children.map((child) => (
              <CategoryTreeItem
                key={child.id}
                category={child}
                selectedId={selectedId}
                onSelect={onSelect}
                level={level + 1}
              />
            ))}
          </List>
        </Collapse>
      )}
    </>
  );
};

const CompactCategoryTree = ({ categories, selectedId, onSelectCategory }) => {
  if (!categories?.length) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ p: 2 }}>
        Категории не найдены
      </Typography>
    );
  }

  return (
    <List 
      component="nav" 
      disablePadding
      sx={{
        '& .MuiListItemButton-root': {
          borderRadius: 0,
        }
      }}
    >
      {categories.map((category) => (
        <CategoryTreeItem
          key={category.id}
          category={category}
          selectedId={selectedId}
          onSelect={onSelectCategory}
        />
      ))}
    </List>
  );
};

export default CompactCategoryTree;