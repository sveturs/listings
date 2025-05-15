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
  custom_component?: string;
  attribute?: Attribute;
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
        
        // Параллельная загрузка категорий и атрибутов
        const [categoriesResponse, attributesResponse] = await Promise.all([
          axios.get('/api/admin/categories'),
          axios.get('/api/admin/attributes')
        ]);
        
        console.log('Loaded categories:', categoriesResponse.data);
        console.log('Loaded attributes:', attributesResponse.data);
        
        // Извлекаем данные из ответа API, если они в формате { data: [...] }
        const categoryData = categoriesResponse.data && categoriesResponse.data.data ? 
                            categoriesResponse.data.data : categoriesResponse.data;
        const attributeData = attributesResponse.data && attributesResponse.data.data ? 
                             attributesResponse.data.data : attributesResponse.data;
                             
        setCategories(categoryData);
        setAttributes(attributeData);
        
      } catch (err) {
        console.error('Ошибка загрузки данных:', err);
        setError(t('admin.common.fetchError'));
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
      
      // Запрос к API для получения привязок атрибутов категории
      // Явно указываем метод GET
      const response = await axios.request({
        method: 'GET',
        url: `/api/admin/categories/${categoryId}/attributes/export`,
      });
      console.log(`Загрузка привязок атрибутов для категории ID=${categoryId}:`, response.data);
      
      // Преобразуем данные в нужный формат, добавляя информацию об атрибутах
      // Проверяем, что полученные данные являются массивом
      const responseData = Array.isArray(response.data) ? response.data : [];
      const mappingsWithAttributes = responseData.map((mapping: CategoryAttributeMapping) => {
        // Проверяем наличие атрибута в mapping
        if (mapping.attribute) {
          return mapping;
        }
        // Иначе ищем атрибут в списке всех атрибутов
        return {
          ...mapping,
          attribute: attributes.find(a => a.id === mapping.attribute_id)
        };
      });
      
      setCategoryAttributeMappings(mappingsWithAttributes);
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
    <Box sx={{ p: 0, maxWidth: '100%' }}>
      <Paper 
        sx={{ 
          p: { xs: 1, sm: 2 }, 
          borderRadius: 0,
          boxShadow: 'none',
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
              <span>
                <IconButton 
                  onClick={handleRefresh} 
                  size="small" 
                  color="primary"
                  disabled={loading}
                >
                  <RefreshIcon />
                </IconButton>
              </span>
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

        <Grid container spacing={{ xs: 1, sm: 2, md: 3 }}>
            {/* Левая панель - Выбор категории */}
            <Grid item xs={12} md={3}>
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
            <Grid item xs={12} md={9}>
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