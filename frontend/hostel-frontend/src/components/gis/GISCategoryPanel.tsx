import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Drawer,
  Box,
  Typography,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Collapse,
  IconButton,
  Divider,
  TextField,
  InputAdornment,
  CircularProgress
} from '@mui/material';
import { 
  ExpandLess, 
  ExpandMore, 
  Search as SearchIcon,
  Close as CloseIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import axios from '../../api/axios';

interface Category {
  id: number | string;
  name: string;
  children?: Category[];
  [key: string]: any;
}

interface ExpandedCategories {
  [key: string]: boolean;
}

interface GISCategoryPanelProps {
  open: boolean;
  onClose: () => void;
  onCategorySelect: (category: Category) => void;
}

// Это объявление теперь не нужно, так как есть более полное ниже

// Правильное типизирование styled-компонентов
interface CategoryDrawerProps {
  anchor?: 'left' | 'right' | 'top' | 'bottom';
  open?: boolean;
  onClose?: () => void;
  variant?: 'permanent' | 'persistent' | 'temporary';
  children?: React.ReactNode;
}

const CategoryDrawer = styled(Drawer)<CategoryDrawerProps>(({ theme }) => ({
  '& .MuiDrawer-paper': {
    width: 320,
    boxSizing: 'border-box',
    [theme.breakpoints.down('sm')]: {
      width: '100%'
    }
  }
})) as React.ComponentType<CategoryDrawerProps>;

interface SearchInputProps {
  variant?: 'outlined' | 'standard' | 'filled';
  placeholder?: string;
  fullWidth?: boolean;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  InputProps?: any;
  children?: React.ReactNode;
}

const SearchInput = styled(TextField)<SearchInputProps>(({ theme }) => ({
  margin: theme.spacing(2),
  '& .MuiOutlinedInput-root': {
    borderRadius: 20
  }
})) as React.ComponentType<SearchInputProps>;

interface CategoryListItemProps {
  depth?: number;
  onClick?: () => void;
  children?: React.ReactNode;
}

// Create a typed styled component for CategoryListItem
const CategoryListItem = styled(ListItem, {
  shouldForwardProp: (prop) => prop !== 'depth'
})<CategoryListItemProps>(({ theme, depth = 0 }) => ({
  paddingLeft: theme.spacing(2 + (depth * 2)),
  cursor: 'pointer',
  '&:hover': {
    backgroundColor: theme.palette.action.hover
  }
})) as React.ComponentType<CategoryListItemProps>;

const GISCategoryPanel: React.FC<GISCategoryPanelProps> = ({ open, onClose, onCategorySelect }) => {
  const { t } = useTranslation('gis');
  const [categories, setCategories] = useState<Category[]>([]);
  const [expandedCategories, setExpandedCategories] = useState<ExpandedCategories>({});
  const [loading, setLoading] = useState<boolean>(false);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [filteredCategories, setFilteredCategories] = useState<Category[]>([]);

  useEffect(() => {
    const fetchCategories = async (): Promise<void> => {
      setLoading(true);
      try {
        // We'll fetch categories from the API
        const response = await axios.get('/api/categories');
        setCategories(response.data);
        setFilteredCategories(response.data);
      } catch (error) {
        console.error('Error fetching categories:', error);
      } finally {
        setLoading(false);
      }
    };

    if (open) {
      fetchCategories();
    }
  }, [open]);

  useEffect(() => {
    // Filter categories based on search query
    if (!searchQuery) {
      setFilteredCategories(categories);
      return;
    }

    const query = searchQuery.toLowerCase();
    const filterCategoriesRecursive = (items: Category[]): Category[] => {
      return items.filter(item => {
        const nameMatches = item.name.toLowerCase().includes(query);
        let childrenMatch = false;
        
        if (item.children && item.children.length > 0) {
          const matchingChildren = filterCategoriesRecursive(item.children);
          childrenMatch = matchingChildren.length > 0;
          // Replace children with filtered children
          if (childrenMatch) {
            item = { ...item, children: matchingChildren };
          }
        }
        
        return nameMatches || childrenMatch;
      });
    };

    const filtered = filterCategoriesRecursive([...categories]);
    setFilteredCategories(filtered);
  }, [searchQuery, categories]);

  const handleToggle = (categoryId: string | number): void => {
    setExpandedCategories(prev => ({
      ...prev,
      [categoryId.toString()]: !prev[categoryId.toString()]
    }));
  };

  const handleCategorySelect = (category: Category): void => {
    if (onCategorySelect) {
      onCategorySelect(category);
    }
    if (window.innerWidth < 600) {
      onClose(); // Close the drawer on mobile after selection
    }
  };

  const renderCategoryItem = (category: Category, depth = 0): React.ReactNode => {
    const hasChildren = category.children && category.children.length > 0;
    const isExpanded = expandedCategories[category.id];

    return (
      <React.Fragment key={category.id}>
        <CategoryListItem 
          depth={depth}
          onClick={() => {
            if (hasChildren) {
              handleToggle(category.id);
            }
            handleCategorySelect(category);
          }}
        >
          <ListItemText primary={category.name} />
          {hasChildren && (
            isExpanded ? <ExpandLess /> : <ExpandMore />
          )}
        </CategoryListItem>
        
        {hasChildren && (
          <Collapse in={isExpanded} timeout="auto" unmountOnExit>
            <List component="div" disablePadding>
              {category.children.map(child => renderCategoryItem(child, depth + 1))}
            </List>
          </Collapse>
        )}
      </React.Fragment>
    );
  };

  return (
    <CategoryDrawer
      anchor="left"
      open={open}
      onClose={onClose}
      variant="temporary"
    >
      <Box display="flex" alignItems="center" justifyContent="space-between" p={2}>
        <Typography variant="h6">{t('categories.title')}</Typography>
        <IconButton onClick={onClose} size="large">
          <CloseIcon />
        </IconButton>
      </Box>
      
      <Divider />
      
      <SearchInput
        variant="outlined"
        placeholder={t('search.placeholder')}
        fullWidth
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon />
            </InputAdornment>
          ),
          endAdornment: searchQuery ? (
            <InputAdornment position="end">
              <IconButton onClick={() => setSearchQuery('')} size="small">
                <CloseIcon fontSize="small" />
              </IconButton>
            </InputAdornment>
          ) : null
        }}
      />
      
      <Box sx={{ overflow: "auto", flexGrow: 1 }}>
        {loading ? (
          <Box display="flex" justifyContent="center" my={4}>
            <CircularProgress />
          </Box>
        ) : filteredCategories.length > 0 ? (
          <List>
            {filteredCategories.map(category => renderCategoryItem(category))}
          </List>
        ) : (
          <Box p={3} textAlign="center">
            <Typography color="textSecondary">
              {searchQuery ? t('search.noResults') : t('categories.title')}
            </Typography>
          </Box>
        )}
      </Box>
    </CategoryDrawer>
  );
};

export default GISCategoryPanel;