//frontend/hostel-frontend/src/components/marketplace/ImprovedCategoryTree.js
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
import {
  ChevronRight,
  ChevronDown,
} from 'lucide-react';

const CategoryTreeItem = ({ 
  category, 
  selectedId, 
  onSelect,
  level = 0 
}) => {
  const [isOpen, setIsOpen] = useState(level === 0);
  const theme = useTheme();

  const getTotalListings = (cat) => {
    let total = cat.listing_count || 0;
    if (cat.children) {
      total += cat.children.reduce((sum, child) => sum + getTotalListings(child), 0);
    }
    return total;
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
          pl: level * 2 + 1,
          py: 0.5,
          minHeight: 32,
          color: isSelected ? 'primary.main' : 'text.primary',
          '&:hover': {
            color: 'primary.main',
          },
          '&.Mui-selected': {
            backgroundColor: theme.palette.primary.main + '08',
            '&:hover': {
              backgroundColor: theme.palette.primary.main + '12',
            }
          }
        }}
      >
        <ListItemIcon 
          sx={{ 
            minWidth: 24,
            color: 'inherit'
          }}
        >
          {hasChildren && (
            isOpen ? <ChevronDown size={18} /> : <ChevronRight size={18} />
          )}
        </ListItemIcon>

        <ListItemText
          primary={
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <Typography
                variant="body2"
                sx={{
                  fontWeight: isSelected ? 500 : 400,
                  fontSize: '0.875rem',
                }}
              >
                {category.name}
              </Typography>
              {totalListings > 0 && (
                <Typography
                  variant="caption"
                  sx={{
                    ml: 0.5,
                    color: isSelected ? 'primary.main' : 'text.secondary',
                  }}
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
          <List component="div" disablePadding>
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