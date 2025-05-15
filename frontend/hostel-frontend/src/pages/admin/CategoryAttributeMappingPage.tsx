import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  Divider,
  CircularProgress,
  Alert,
  Card,
  Stack,
  Button,
  IconButton,
  Tooltip,
  alpha,
  useTheme
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import AttributeMappingList from '../../components/admin/AttributeMappingList';
import AddAttributeToCategory from '../../components/admin/AddAttributeToCategory';
import CategorySelector from '../../components/admin/CategorySelector';
import CustomUIManager from '../../components/admin/CustomUIManager';
import CategoryAttributeExporter from '../../components/admin/CategoryAttributeExporter';
import RefreshIcon from '@mui/icons-material/Refresh';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import LabelImportantIcon from '@mui/icons-material/LabelImportant';

// Типы данных для атрибутов и категорий
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

interface Attribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: string;
  options?: any;
  is_searchable: boolean;
  is_filterable: boolean;
  is_required: boolean;
  sort_order: number;
  custom_component?: string;
  created_at: string;
}

interface CategoryAttributeMapping {
  category_id: number;
  attribute_id: number;
  is_required: boolean;
  is_enabled: boolean;
  sort_order: number;
}

// Основной компонент страницы
const CategoryAttributeMappingPage: React.FC = () => {
  const { t } = useTranslation();
  const theme = useTheme();
  const [categories, setCategories] = useState<Category[]>([]);
  const [attributes, setAttributes] = useState<Attribute[]>([]);
  const [categoryAttributeMappings, setCategoryAttributeMappings] = useState<CategoryAttributeMapping[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState<number>(0);

  // Загрузка данных при монтировании компонента
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);
        
        // Добавим тестовые категории для отладки UI
        const mockCategories = [
          { id: 1, name: 'Электроника', slug: 'electronics', parent_id: null, listing_count: 15, has_custom_ui: true, custom_ui_component: 'ElectronicsUI' },
          { id: 2, name: 'Смартфоны', slug: 'smartphones', parent_id: 1, listing_count: 8, has_custom_ui: false },
          { id: 3, name: 'Ноутбуки', slug: 'laptops', parent_id: 1, listing_count: 5, has_custom_ui: true, custom_ui_component: 'LaptopsUI' },
          { id: 4, name: 'Одежда', slug: 'clothing', parent_id: null, listing_count: 20, has_custom_ui: false },
          { id: 5, name: 'Мужская', slug: 'mens', parent_id: 4, listing_count: 12, has_custom_ui: false },
          { id: 6, name: 'Женская', slug: 'womens', parent_id: 4, listing_count: 8, has_custom_ui: false },
          { id: 7, name: 'Детская', slug: 'kids', parent_id: 4, listing_count: 0, has_custom_ui: false },
        ];
        
        const mockAttributes: Attribute[] = [
          { id: 1, name: 'brand', display_name: 'Бренд', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: true, sort_order: 1, created_at: new Date().toISOString() },
          { id: 2, name: 'price', display_name: 'Цена', attribute_type: 'number', is_searchable: true, is_filterable: true, is_required: true, sort_order: 2, created_at: new Date().toISOString() },
          { id: 3, name: 'color', display_name: 'Цвет', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 3, created_at: new Date().toISOString() },
          { id: 4, name: 'condition', display_name: 'Состояние', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: true, sort_order: 4, created_at: new Date().toISOString() },
          { id: 5, name: 'size', display_name: 'Размер', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 5, created_at: new Date().toISOString() },
        ];
        
        // Используем мок-данные для отладки UI
        console.log('Используем мок-данные для категорий и атрибутов');
        setCategories(mockCategories);
        setAttributes(mockAttributes);
        
      } catch (err) {
        console.error('Ошибка загрузки данных:', err);
        setError(`Ошибка загрузки данных: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [t, refreshKey]);

  // Загрузка привязок атрибутов при выборе категории
  const loadCategoryAttributes = async (categoryId: number) => {
    try {
      setLoading(true);
      setError(null);
      
      // Мок-данные для привязок атрибутов
      const mockMappings = [
        // Электроника (id=1)
        ...(categoryId === 1 ? [
          { 
            category_id: 1, 
            attribute_id: 1, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 1,
            attribute: attributes.find(a => a.id === 1) 
          },
          { 
            category_id: 1, 
            attribute_id: 2, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 2,
            attribute: attributes.find(a => a.id === 2)
          },
          { 
            category_id: 1, 
            attribute_id: 4, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 3,
            attribute: attributes.find(a => a.id === 4)
          },
        ] : []),
        
        // Смартфоны (id=2)
        ...(categoryId === 2 ? [
          { 
            category_id: 2, 
            attribute_id: 1, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 1,
            attribute: attributes.find(a => a.id === 1)
          },
          { 
            category_id: 2, 
            attribute_id: 2, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 2,
            attribute: attributes.find(a => a.id === 2)
          },
          { 
            category_id: 2, 
            attribute_id: 3, 
            is_required: false, 
            is_enabled: true, 
            sort_order: 3,
            attribute: attributes.find(a => a.id === 3)
          },
          { 
            category_id: 2, 
            attribute_id: 4, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 4,
            attribute: attributes.find(a => a.id === 4)
          },
        ] : []),
        
        // Одежда (id=4)
        ...(categoryId === 4 ? [
          { 
            category_id: 4, 
            attribute_id: 1, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 1,
            attribute: attributes.find(a => a.id === 1)
          },
          { 
            category_id: 4, 
            attribute_id: 2, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 2,
            attribute: attributes.find(a => a.id === 2)
          },
          { 
            category_id: 4, 
            attribute_id: 3, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 3,
            attribute: attributes.find(a => a.id === 3)
          },
          { 
            category_id: 4, 
            attribute_id: 5, 
            is_required: true, 
            is_enabled: true, 
            sort_order: 4,
            attribute: attributes.find(a => a.id === 5)
          },
        ] : []),
        
        // Для остальных категорий - пустые массивы по умолчанию
      ];
      
      // Используем мок-данные
      console.log(`Загрузка привязок атрибутов для категории ID=${categoryId}:`, mockMappings);
      setCategoryAttributeMappings(mockMappings);
    } catch (err) {
      console.error('Error fetching category attributes:', err);
      setError(t('admin.categoryAttributes.fetchMappingsError'));
    } finally {
      setLoading(false);
    }
  };

  // Обработчик выбора категории
  const handleCategorySelect = (category: Category) => {
    console.log('Selected category:', category);
    setSelectedCategory(category);
    loadCategoryAttributes(category.id);
  };

  // Обработчик удаления атрибута из категории
  const handleRemoveAttributeFromCategory = async (attributeId: number) => {
    if (!selectedCategory) return;
    
    try {
      setLoading(true);
      setError(null);
      
      await axios.delete(`/api/admin/categories/${selectedCategory.id}/attributes/${attributeId}`);
      
      // Перезагружаем привязки после удаления
      await loadCategoryAttributes(selectedCategory.id);
      
    } catch (err) {
      console.error('Error removing attribute from category:', err);
      setError(t('admin.categoryAttributes.removeAttributeError'));
    } finally {
      setLoading(false);
    }
  };

  // Обработчик обновления параметров атрибута в категории
  const handleUpdateAttributeCategory = async (
    attributeId: number, 
    isRequired: boolean, 
    isEnabled: boolean,
    sortOrder: number
  ) => {
    if (!selectedCategory) return;
    
    try {
      setLoading(true);
      setError(null);
      
      await axios.put(`/api/admin/categories/${selectedCategory.id}/attributes/${attributeId}`, {
        is_required: isRequired,
        is_enabled: isEnabled,
        sort_order: sortOrder
      });
      
      // Перезагружаем привязки после обновления
      await loadCategoryAttributes(selectedCategory.id);
      
    } catch (err) {
      console.error('Error updating attribute category settings:', err);
      setError(t('admin.categoryAttributes.updateAttributeError'));
    } finally {
      setLoading(false);
    }
  };

  // Функция обновления данных
  const handleRefresh = () => {
    setRefreshKey(prev => prev + 1);
    if (selectedCategory) {
      loadCategoryAttributes(selectedCategory.id);
    }
  };

  // Рендеринг содержимого страницы
  return (
    <Box sx={{ p: 3 }}>
      <Paper 
        sx={{ 
          p: 2, 
          borderRadius: 2,
          boxShadow: theme.shadows[2],
          bgcolor: 'background.paper'
        }}
      >
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Box>
            <Typography variant="h5" component="h1" fontWeight="bold" color="primary.main">
              {t('admin.categoryAttributes.title')}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {t('admin.categoryAttributes.description')}
            </Typography>
          </Box>
          <Box>
            {loading && (
              <CircularProgress size={24} sx={{ mr: 2 }} />
            )}
            <Tooltip title={t('admin.common.refresh')}>
              <IconButton 
                onClick={handleRefresh} 
                size="small" 
                color="primary"
                disabled={loading}
              >
                <RefreshIcon />
              </IconButton>
            </Tooltip>
          </Box>
        </Box>

        <Divider sx={{ mb: 3 }} />

        {error && (
          <Alert 
            severity="error" 
            sx={{ 
              mb: 2,
              borderRadius: 1
            }}
            onClose={() => setError(null)}
          >
            {error}
          </Alert>
        )}

        <Grid container spacing={3}>
            {/* Левая панель - Выбор категории */}
            <Grid item xs={12} md={4}>
              <Box>
                <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                  <LabelImportantIcon fontSize="small" color="primary" />
                  {t('admin.categoryAttributes.selectCategory')}
                </Typography>
                <Box sx={{ mt: 2, position: 'relative' }}>
                  {loading && categories.length === 0 ? (
                    <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
                      <CircularProgress />
                    </Box>
                  ) : (
                    <CategorySelector
                      categories={categories}
                      onCategorySelect={handleCategorySelect}
                      selectedCategoryId={selectedCategory?.id}
                    />
                  )}
                </Box>
              </Box>
            </Grid>

            {/* Правая панель - Управление атрибутами */}
            <Grid item xs={12} md={8}>
              {selectedCategory ? (
                <Box>
                  <Paper 
                    variant="outlined" 
                    sx={{ 
                      p: 2, 
                      mb: 3,
                      borderRadius: 1,
                      borderColor: alpha(theme.palette.primary.main, 0.3),
                      borderWidth: '1px',
                      position: 'relative',
                      overflow: 'hidden',
                      '&:before': {
                        content: '""',
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        width: '4px',
                        height: '100%',
                        bgcolor: 'primary.main'
                      }
                    }}
                  >
                    <Typography 
                      variant="h6" 
                      gutterBottom 
                      sx={{ 
                        display: 'flex', 
                        alignItems: 'center', 
                        gap: 0.5,
                        color: 'primary.main',
                        pl: 1
                      }}
                    >
                      <InfoOutlinedIcon fontSize="small" />
                      {t('admin.categoryAttributes.categoryAttributes', { category: selectedCategory.name })}
                    </Typography>
                    
                    <Box sx={{ mt: 2, position: 'relative' }}>
                      {loading ? (
                        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                          <CircularProgress size={30} />
                        </Box>
                      ) : (
                        <AttributeMappingList 
                          categoryId={selectedCategory.id}
                          onError={setError}
                          mappings={categoryAttributeMappings}
                          onUpdateAttribute={handleUpdateAttributeCategory}
                          onRemoveAttribute={handleRemoveAttributeFromCategory}
                        />
                      )}
                    </Box>
                  </Paper>
                  
                  {/* Добавляем компонент управления пользовательским UI */}
                  <CustomUIManager
                    categoryId={selectedCategory.id}
                    onCategoryUpdate={handleRefresh}
                  />
                  
                  <Paper 
                    variant="outlined" 
                    sx={{ 
                      p: 2,
                      borderRadius: 1
                    }}
                  >
                    <Typography variant="h6" gutterBottom>
                      {t('admin.categoryAttributes.addAttribute')}
                    </Typography>
                    
                    <AddAttributeToCategory
                      categoryId={selectedCategory.id}
                      onAttributeAdded={() => loadCategoryAttributes(selectedCategory.id)}
                      onError={setError}
                    />
                  </Paper>
                  
                  {/* Компонент экспорта/импорта/копирования атрибутов */}
                  <CategoryAttributeExporter
                    categoryId={selectedCategory.id}
                    onSuccess={(message) => {
                      setError(null);
                      loadCategoryAttributes(selectedCategory.id);
                    }}
                    onError={setError}
                  />
                </Box>
              ) : (
                <Paper 
                  variant="outlined" 
                  sx={{ 
                    p: 4, 
                    textAlign: 'center',
                    borderRadius: 1,
                    height: '100%',
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'center',
                    alignItems: 'center',
                    bgcolor: alpha(theme.palette.background.default, 0.5)
                  }}
                >
                  <LabelImportantIcon sx={{ fontSize: 40, color: 'text.secondary', mb: 2, opacity: 0.5 }} />
                  <Typography variant="body1" color="text.secondary">
                    {t('admin.categoryAttributes.selectCategoryPrompt')}
                  </Typography>
                </Paper>
              )}
            </Grid>
          </Grid>
      </Paper>
    </Box>
  );
};

export default CategoryAttributeMappingPage;