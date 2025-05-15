import React, { useState, useEffect } from 'react';
import {
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Collapse,
  Typography,
  Box,
  TextField,
  InputAdornment,
  Paper,
  Divider,
  Grid,
  Badge,
  Card,
  CardContent,
  Chip,
  Fade,
  useTheme,
  alpha
} from '@mui/material';
import {
  ExpandLess,
  ExpandMore,
  FolderOutlined,
  FolderOpen,
  Search as SearchIcon,
  Category as CategoryIcon,
  Info as InfoIcon,
  CodeOutlined,
  ListAlt
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';

interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number | null;
  icon?: string;
  translations?: Record<string, string>;
  listing_count: number;
  has_custom_ui: boolean;
  custom_ui_component?: string;
}

interface CategoryTreeNode extends Category {
  children: CategoryTreeNode[];
}

interface CategorySelectorProps {
  categories: Category[];
  onCategorySelect: (category: Category) => void;
  selectedCategoryId?: number;
}

const CategorySelector: React.FC<CategorySelectorProps> = ({
  categories,
  onCategorySelect,
  selectedCategoryId
}) => {
  const { t } = useTranslation();
  const theme = useTheme();
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());
  const [categoryTree, setCategoryTree] = useState<CategoryTreeNode[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);

  // Конвертация плоского списка категорий в иерархическое дерево
  useEffect(() => {
    console.log('CategorySelector - categories prop:', categories);
    
    // Даже если категорий нет, мы всё равно создаем пустое дерево, чтобы компонент отрендерился
    const buildCategoryTree = (parentId: number | null = null): CategoryTreeNode[] => {
      if (!categories || categories.length === 0) {
        return [];
      }
      
      return categories
        .filter(cat => cat.parent_id === parentId)
        .map(cat => ({
          ...cat,
          children: buildCategoryTree(cat.id)
        }))
        .sort((a, b) => a.name.localeCompare(b.name));
    };

    const tree = buildCategoryTree(null);
    console.log('CategorySelector - built tree:', tree);
    setCategoryTree(tree);

    // Если есть выбранная категория, расширяем её родителей и сохраняем данные категории
    if (selectedCategoryId) {
      const selectedCat = categories.find(cat => cat.id === selectedCategoryId);
      if (selectedCat) {
        setSelectedCategory(selectedCat);
      }
      
      const expandParents = (categoryId: number) => {
        const category = categories.find(cat => cat.id === categoryId);
        if (category && category.parent_id) {
          setExpandedIds(prev => {
            const newSet = new Set(prev);
            newSet.add(category.parent_id as number);
            return newSet;
          });
          expandParents(category.parent_id);
        }
      };
      expandParents(selectedCategoryId);
    }
  }, [categories, selectedCategoryId]);

  // Обработчик переключения раскрытия/свертывания категории
  const toggleExpand = (categoryId: number, e?: React.MouseEvent) => {
    if (e) {
      e.stopPropagation();
    }
    
    setExpandedIds(prev => {
      const newSet = new Set(prev);
      if (newSet.has(categoryId)) {
        newSet.delete(categoryId);
      } else {
        newSet.add(categoryId);
      }
      return newSet;
    });
  };

  // Обработчик выбора категории
  const handleCategorySelect = (category: Category) => {
    setSelectedCategory(category);
    onCategorySelect(category);
  };

  // Фильтрация категорий по поисковому запросу
  const filterCategories = (node: CategoryTreeNode): boolean => {
    if (!searchQuery) return true;
    
    const nameMatches = node.name.toLowerCase().includes(searchQuery.toLowerCase());
    const slugMatches = node.slug.toLowerCase().includes(searchQuery.toLowerCase());
    const childrenMatch = node.children.some(filterCategories);
    
    return nameMatches || slugMatches || childrenMatch;
  };

  // Найти родительскую цепочку для категории
  const getCategoryPath = (categoryId: number): string => {
    const path: string[] = [];
    
    const findPath = (id: number) => {
      const category = categories.find(cat => cat.id === id);
      if (!category) return;
      
      path.unshift(category.name);
      
      if (category.parent_id) {
        findPath(category.parent_id);
      }
    };
    
    findPath(categoryId);
    return path.join(' > ');
  };

  // Рекурсивный рендеринг дерева категорий
  const renderCategoryTree = (nodes: CategoryTreeNode[], level = 0) => {
    return nodes
      .filter(filterCategories)
      .map(node => {
        const isExpanded = expandedIds.has(node.id);
        const isSelected = selectedCategoryId === node.id;
        const hasChildren = node.children && node.children.length > 0;
        
        return (
          <React.Fragment key={node.id}>
            <ListItem 
              disablePadding 
              sx={{ 
                pl: level * 1.5, 
                borderRadius: 1,
                transition: theme.transitions.create(['background-color'], {
                  duration: theme.transitions.duration.shorter,
                }),
                ...(isSelected && {
                  bgcolor: alpha(theme.palette.primary.main, 0.1),
                  boxShadow: `inset 3px 0 0 ${theme.palette.primary.main}`
                })
              }}
            >
              <ListItemButton 
                dense
                onClick={() => handleCategorySelect(node)}
                selected={isSelected}
                sx={{
                  borderRadius: 1,
                  py: 0.5
                }}
              >
                <ListItemIcon sx={{ minWidth: 36 }}>
                  {hasChildren ? (
                    <Badge 
                      badgeContent={node.children.length} 
                      color="primary"
                      sx={{ '& .MuiBadge-badge': { fontSize: '0.6rem', height: 16, minWidth: 16 } }}
                    >
                      {isExpanded ? 
                        <FolderOpen color={isSelected ? "primary" : "action"} fontSize="small" /> : 
                        <FolderOutlined color={isSelected ? "primary" : "action"} fontSize="small" />
                      }
                    </Badge>
                  ) : (
                    <CategoryIcon color={isSelected ? "primary" : "action"} fontSize="small" />
                  )}
                </ListItemIcon>
                <ListItemText 
                  primary={
                    <Typography 
                      variant="body2" 
                      noWrap 
                      sx={{ 
                        fontWeight: isSelected ? 'bold' : 'normal',
                        color: isSelected ? 'primary.main' : 'text.primary'
                      }}
                    >
                      {node.name}
                    </Typography>
                  }
                  secondary={
                    <Typography variant="caption" color="text.secondary" noWrap>
                      {node.listing_count > 0 ? 
                        t('admin.categories.listingCount', { count: node.listing_count }) :
                        node.slug
                      }
                    </Typography>
                  }
                />
                {hasChildren && (
                  <Box onClick={(e) => toggleExpand(node.id, e)}>
                    {isExpanded ? <ExpandLess /> : <ExpandMore />}
                  </Box>
                )}
              </ListItemButton>
            </ListItem>
            {hasChildren && (
              <Collapse in={isExpanded} timeout="auto" unmountOnExit>
                <List disablePadding>
                  {renderCategoryTree(node.children, level + 1)}
                </List>
              </Collapse>
            )}
          </React.Fragment>
        );
      });
  };

  // Функция рендеринга панели с детальной информацией о выбранной категории
  const renderCategoryDetails = () => {
    if (!selectedCategory) return null;

    return (
      <Fade in={!!selectedCategory}>
        <Card variant="outlined">
          <CardContent>
            <Typography variant="h6" color="primary" gutterBottom>
              {selectedCategory.name}
            </Typography>
            
            {selectedCategory.parent_id && (
              <Typography variant="caption" display="block" gutterBottom>
                <strong>{t('admin.categories.path')}:</strong> {getCategoryPath(selectedCategory.id)}
              </Typography>
            )}
            
            <Divider sx={{ my: 1 }} />
            
            <Grid container spacing={1} sx={{ mt: 1 }}>
              <Grid item xs={6}>
                <Typography variant="body2">
                  <strong>ID:</strong> {selectedCategory.id}
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="body2">
                  <strong>{t('admin.categories.slug')}:</strong> {selectedCategory.slug}
                </Typography>
              </Grid>
              <Grid item xs={12}>
                <Box sx={{ mt: 1, mb: 1 }}>
                  <Typography variant="body2">
                    <strong>{t('admin.categories.parent')}:</strong> {
                      selectedCategory.parent_id ? 
                        categories.find(c => c.id === selectedCategory.parent_id)?.name || selectedCategory.parent_id : 
                        t('admin.categories.noParent')
                    }
                  </Typography>
                </Box>
              </Grid>
              <Grid item xs={12}>
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 1 }}>
                  <Chip 
                    icon={<ListAlt />} 
                    label={t('admin.categories.listingCount', { count: selectedCategory.listing_count })} 
                    size="small" 
                    color={selectedCategory.listing_count > 0 ? "primary" : "default"}
                    variant={selectedCategory.listing_count > 0 ? "filled" : "outlined"}
                  />
                  
                  {selectedCategory.has_custom_ui && (
                    <Chip 
                      icon={<CodeOutlined />} 
                      label={selectedCategory.custom_ui_component || t('admin.categories.customUi')} 
                      size="small" 
                      color="secondary"
                    />
                  )}
                </Box>
              </Grid>
            </Grid>
            
            {selectedCategory.translations && Object.keys(selectedCategory.translations).length > 0 && (
              <>
                <Divider sx={{ my: 1 }} />
                <Typography variant="subtitle2" gutterBottom>
                  {t('admin.categories.translations')}:
                </Typography>
                <Box sx={{ ml: 1 }}>
                  {Object.entries(selectedCategory.translations).map(([lang, text]) => (
                    <Typography key={lang} variant="caption" display="block">
                      <strong>{lang.toUpperCase()}:</strong> {text}
                    </Typography>
                  ))}
                </Box>
              </>
            )}
          </CardContent>
        </Card>
      </Fade>
    );
  };

  return (
    <Grid container spacing={2}>
      <Grid item xs={12} md={selectedCategory ? 7 : 12}>
        <Paper variant="outlined" sx={{ maxHeight: '70vh', display: 'flex', flexDirection: 'column' }}>
          <Box sx={{ p: 1.5, bgcolor: 'background.default' }}>
            <TextField
              fullWidth
              size="small"
              placeholder={t('admin.categories.searchPlaceholder')}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon fontSize="small" />
                  </InputAdornment>
                ),
              }}
            />
            {searchQuery && categories.length > 0 && (
              <Box sx={{ mt: 1 }}>
                <Typography variant="caption" color="text.secondary">
                  {categoryTree.flatMap(cat => [cat, ...cat.children]).filter(filterCategories).length} {t('admin.categories.searchResults')}
                </Typography>
              </Box>
            )}
          </Box>
          <Divider />
          <Box sx={{ overflow: 'auto', flexGrow: 1, p: 1 }}>
            {!categories || categories.length === 0 ? (
              <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', p: 3 }}>
                <CategoryIcon color="disabled" sx={{ fontSize: 48, mb: 2, opacity: 0.5 }} />
                <Typography variant="body1" color="text.secondary" align="center" gutterBottom>
                  {t('admin.categories.noCategories')}
                </Typography>
                <Typography variant="caption" color="text.secondary" align="center">
                  {t('admin.categories.addCategoriesFirst')}
                </Typography>
              </Box>
            ) : categoryTree.length > 0 ? (
              <List dense disablePadding>
                {renderCategoryTree(categoryTree)}
              </List>
            ) : (
              <Typography variant="body2" sx={{ p: 2, textAlign: 'center', color: 'text.secondary' }}>
                {t('admin.categories.loadingCategories')}
              </Typography>
            )}
          </Box>
        </Paper>
      </Grid>
      
      {selectedCategory && (
        <Grid item xs={12} md={5}>
          {renderCategoryDetails()}
        </Grid>
      )}
    </Grid>
  );
};

export default CategorySelector;