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

const CategoryDrawer = styled(Drawer)(({ theme }) => ({
  '& .MuiDrawer-paper': {
    width: 320,
    boxSizing: 'border-box',
    [theme.breakpoints.down('sm')]: {
      width: '100%'
    }
  }
}));

const SearchInput = styled(TextField)(({ theme }) => ({
  margin: theme.spacing(2),
  '& .MuiOutlinedInput-root': {
    borderRadius: 20
  }
}));

const CategoryListItem = styled(ListItem)(({ theme, depth = 0 }) => ({
  paddingLeft: theme.spacing(2 + (depth * 2)),
  cursor: 'pointer',
  '&:hover': {
    backgroundColor: theme.palette.action.hover
  }
}));

const GISCategoryPanel = ({ open, onClose, onCategorySelect }) => {
  const { t } = useTranslation('gis');
  const [categories, setCategories] = useState([]);
  const [expandedCategories, setExpandedCategories] = useState({});
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [filteredCategories, setFilteredCategories] = useState([]);

  useEffect(() => {
    const fetchCategories = async () => {
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
    const filterCategoriesRecursive = (items) => {
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

  const handleToggle = (categoryId) => {
    setExpandedCategories(prev => ({
      ...prev,
      [categoryId]: !prev[categoryId]
    }));
  };

  const handleCategorySelect = (category) => {
    if (onCategorySelect) {
      onCategorySelect(category);
    }
    if (window.innerWidth < 600) {
      onClose(); // Close the drawer on mobile after selection
    }
  };

  const renderCategoryItem = (category, depth = 0) => {
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
      
      <Box overflow="auto" flexGrow={1}>
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
              {searchQuery ? t('search.noResults') : t('categories.all')}
            </Typography>
          </Box>
        )}
      </Box>
    </CategoryDrawer>
  );
};

export default GISCategoryPanel;