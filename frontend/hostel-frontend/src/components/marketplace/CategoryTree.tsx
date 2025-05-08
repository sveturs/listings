// frontend/hostel-frontend/src/components/marketplace/CategoryTree.tsx
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
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
import { Category as BaseCategory } from './HierarchicalCategorySelect';

// Расширенный интерфейс категории для дерева категорий
interface Category extends BaseCategory {
  listing_count?: number;
  children?: Category[];
}

interface CategoryTreeItemProps {
  category: Category;
  selectedId: number | string | null;
  onSelect: (id: number | string) => void;
  level?: number;
}

interface CompactCategoryTreeProps {
  categories: Category[];
  selectedId: number | string | null;
  onSelectCategory: (id: number | string) => void;
}

const CategoryTreeItem: React.FC<CategoryTreeItemProps> = ({ 
  category, 
  selectedId, 
  onSelect,
  level = 0 
}) => {
  const { i18n } = useTranslation();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const theme = useTheme();

  const getTotalListings = (cat: Category): number => {
    let total = cat.listing_count || 0;
    if (cat.children) {
      total += cat.children.reduce((sum, child) => sum + getTotalListings(child), 0);
    }
    return total;
  };

  const getTranslatedName = (category: Category, currentLanguage: string): string => {
    if (category.translations && category.translations[currentLanguage]) {
      return category.translations[currentLanguage];
    }
    return category.name;
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
                {getTranslatedName(category, i18n.language)}
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

const CompactCategoryTree: React.FC<CompactCategoryTreeProps> = ({ 
  categories, 
  selectedId, 
  onSelectCategory 
}) => {
  const { t } = useTranslation('marketplace');

  if (!categories?.length) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ p: 2 }}>
        {t('listings.create.ChooseAcategory')}
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