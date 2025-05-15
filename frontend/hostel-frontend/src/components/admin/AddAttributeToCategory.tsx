import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormControlLabel,
  Checkbox,
  Button,
  Paper,
  Grid,
  CircularProgress,
  Alert,
  Autocomplete
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';

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

interface AddAttributeToCategoryProps {
  categoryId: number;
  onAttributeAdded?: () => void;
  onError?: (message: string) => void;
}

const AddAttributeToCategory: React.FC<AddAttributeToCategoryProps> = ({
  categoryId,
  onAttributeAdded,
  onError
}) => {
  const { t } = useTranslation();
  const [availableAttributes, setAvailableAttributes] = useState<Attribute[]>([]);
  const [existingMappings, setExistingMappings] = useState<CategoryAttributeMapping[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedAttribute, setSelectedAttribute] = useState<Attribute | null>(null);
  const [isRequired, setIsRequired] = useState<boolean>(false);
  const [isEnabled, setIsEnabled] = useState<boolean>(true);
  const [sortOrder, setSortOrder] = useState<number>(0);
  const [searchTerm, setSearchTerm] = useState<string>('');

  // Fetch available attributes and existing mappings
  useEffect(() => {
    if (categoryId) {
      fetchData();
    }
  }, [categoryId, t]);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);

      try {
        // Загружаем все атрибуты и текущие привязки через API
        const [allAttributesResponse, categoryMappingsResponse] = await Promise.all([
          axios.get('/api/admin/attributes'),
          axios.get(`/api/admin/categories/${categoryId}/attributes/export`)
        ]);
        
        console.log('All attributes:', allAttributesResponse.data);
        console.log('Category mappings:', categoryMappingsResponse.data);
        
        const allAttributes = allAttributesResponse.data;
        const categoryMappings = categoryMappingsResponse.data;
        
        setExistingMappings(categoryMappings);
        
        // Filter out attributes that are already mapped to this category
        const mappedAttributeIds = categoryMappings.map((mapping: CategoryAttributeMapping) => mapping.attribute_id);
        const filteredAttributes = allAttributes.filter((attr: Attribute) => 
          !mappedAttributeIds.includes(attr.id)
        );
        
        setAvailableAttributes(filteredAttributes);
      } catch (apiError) {
        console.error('API error:', apiError);
        
        // В режиме разработки используем мок-данные при ошибке API
        if (process.env.NODE_ENV === 'development') {
          console.warn('Используем мок-данные в режиме разработки');
          
          // Мок-данные для всех атрибутов
          const mockAllAttributes: Attribute[] = [
            { id: 1, name: 'brand', display_name: 'Бренд', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 1, created_at: new Date().toISOString() },
            { id: 2, name: 'price', display_name: 'Цена', attribute_type: 'number', is_searchable: true, is_filterable: true, is_required: false, sort_order: 2, created_at: new Date().toISOString() },
            { id: 3, name: 'color', display_name: 'Цвет', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 3, created_at: new Date().toISOString() },
            { id: 4, name: 'condition', display_name: 'Состояние', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 4, created_at: new Date().toISOString() },
            { id: 5, name: 'size', display_name: 'Размер', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 5, created_at: new Date().toISOString() },
            { id: 6, name: 'material', display_name: 'Материал', attribute_type: 'select', is_searchable: true, is_filterable: true, is_required: false, sort_order: 6, created_at: new Date().toISOString() },
            { id: 7, name: 'description', display_name: 'Описание', attribute_type: 'text', is_searchable: true, is_filterable: false, is_required: false, sort_order: 7, created_at: new Date().toISOString() },
          ];
          
          // Мок-данные для существующих привязок атрибутов к категории
          let mockMappings: CategoryAttributeMapping[] = [];
          
          // Например, для категории "Электроника" (id=1)
          if (categoryId === 1) {
            mockMappings = [
              { category_id: 1, attribute_id: 1, is_required: true, is_enabled: true, sort_order: 1 },
              { category_id: 1, attribute_id: 2, is_required: true, is_enabled: true, sort_order: 2 },
              { category_id: 1, attribute_id: 4, is_required: true, is_enabled: true, sort_order: 3 },
            ];
          }
          // Для категории "Смартфоны" (id=2)
          else if (categoryId === 2) {
            mockMappings = [
              { category_id: 2, attribute_id: 1, is_required: true, is_enabled: true, sort_order: 1 },
              { category_id: 2, attribute_id: 2, is_required: true, is_enabled: true, sort_order: 2 },
              { category_id: 2, attribute_id: 3, is_required: false, is_enabled: true, sort_order: 3 },
              { category_id: 2, attribute_id: 4, is_required: true, is_enabled: true, sort_order: 4 },
            ];
          }
          // Для категории "Одежда" (id=4)
          else if (categoryId === 4) {
            mockMappings = [
              { category_id: 4, attribute_id: 1, is_required: true, is_enabled: true, sort_order: 1 },
              { category_id: 4, attribute_id: 2, is_required: true, is_enabled: true, sort_order: 2 },
              { category_id: 4, attribute_id: 3, is_required: true, is_enabled: true, sort_order: 3 },
              { category_id: 4, attribute_id: 5, is_required: true, is_enabled: true, sort_order: 4 },
            ];
          }
          
          setExistingMappings(mockMappings);
          
          // Filter out attributes that are already mapped to this category
          const mappedAttributeIds = mockMappings.map(mapping => mapping.attribute_id);
          const filteredAttributes = mockAllAttributes.filter(attr => 
            !mappedAttributeIds.includes(attr.id)
          );
          
          setAvailableAttributes(filteredAttributes);
        } else {
          // В продакшн-режиме пробрасываем ошибку
          throw apiError;
        }
      }
      
    } catch (err) {
      console.error('Error fetching attributes data:', err);
      const errorMessage = t('admin.categoryAttributes.fetchAttributesError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Filter attributes based on search term
  const filteredAttributes = availableAttributes.filter(attr => 
    (attr.display_name || attr.name).toLowerCase().includes(searchTerm.toLowerCase())
  );

  // Handle adding the attribute to the category
  const handleAddAttribute = async () => {
    if (!selectedAttribute) return;
    
    try {
      setLoading(true);
      setError(null);
      
      try {
        // Добавляем атрибут к категории через API
        await axios.post(`/api/admin/categories/${categoryId}/attributes/${selectedAttribute.id}`, {
          is_required: isRequired,
          is_enabled: isEnabled,
          sort_order: sortOrder
        });
        
        console.log(`Attribute ${selectedAttribute.id} added to category ${categoryId}`, {
          is_required: isRequired,
          is_enabled: isEnabled,
          sort_order: sortOrder
        });
      } catch (apiError) {
        console.error('API error while adding attribute:', apiError);
        
        // В режиме разработки симулируем успешное добавление
        if (process.env.NODE_ENV === 'development') {
          console.warn('Симуляция добавления атрибута в режиме разработки');
          console.log(`${t('admin.categoryAttributes.addAttributeSimulation')}:`, {
            categoryId,
            attributeId: selectedAttribute.id,
            isRequired,
            isEnabled,
            sortOrder
          });
        } else {
          // В продакшн-режиме пробрасываем ошибку
          throw apiError;
        }
      }
      
      // Reset form and refresh data
      setSelectedAttribute(null);
      setIsRequired(false);
      setIsEnabled(true);
      setSortOrder(0);
      
      // Simulate delay for a more realistic user experience
      await new Promise(resolve => setTimeout(resolve, 500));
      
      // Fetch updated data
      await fetchData();
      
      // Notify parent component
      if (onAttributeAdded) {
        onAttributeAdded();
      }
    } catch (err) {
      console.error('Error adding attribute to category:', err);
      const errorMessage = t('admin.categoryAttributes.addAttributeError');
      setError(errorMessage);
      if (onError) onError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Calculate next available sort order
  useEffect(() => {
    if (existingMappings.length > 0) {
      const maxSortOrder = Math.max(...existingMappings.map(m => m.sort_order));
      setSortOrder(maxSortOrder + 10);
    } else {
      setSortOrder(10);
    }
  }, [existingMappings]);

  return (
    <Paper variant="outlined" sx={{ p: { xs: 1, sm: 2 }, mt: 2 }}>
      <Typography variant="h6" gutterBottom>
        {t('admin.categoryAttributes.addAttribute')}
      </Typography>
      
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      <Grid container spacing={{ xs: 1, sm: 2 }}>
        <Grid item xs={12}>
          <Autocomplete
            value={selectedAttribute}
            onChange={(event, newValue) => {
              setSelectedAttribute(newValue);
            }}
            inputValue={searchTerm}
            onInputChange={(event, newInputValue) => {
              setSearchTerm(newInputValue);
            }}
            options={filteredAttributes}
            getOptionLabel={(option) => option.display_name || option.name}
            renderInput={(params) => (
              <TextField 
                {...params} 
                label={t('admin.categoryAttributes.selectAttribute')}
                variant="outlined"
                fullWidth
                size="small"
              />
            )}
            renderOption={(props, option) => (
              <li {...props}>
                <Box>
                  <Typography variant="body1">
                    {option.display_name || option.name}
                  </Typography>
                  <Typography variant="caption" color="textSecondary">
                    {t(`admin.attributeTypes.${option.attribute_type}`)}
                  </Typography>
                </Box>
              </li>
            )}
            loading={loading}
            disabled={loading || availableAttributes.length === 0}
            noOptionsText={t('admin.categoryAttributes.noAttributesAvailable')}
          />
        </Grid>
        
        <Grid item xs={12} sm={4}>
          <FormControlLabel
            control={
              <Checkbox
                checked={isRequired}
                onChange={(e) => setIsRequired(e.target.checked)}
                disabled={loading || !selectedAttribute}
              />
            }
            label={t('admin.attributes.required')}
          />
        </Grid>
        
        <Grid item xs={12} sm={4}>
          <FormControlLabel
            control={
              <Checkbox
                checked={isEnabled}
                onChange={(e) => setIsEnabled(e.target.checked)}
                disabled={loading || !selectedAttribute}
              />
            }
            label={t('admin.attributes.enabled')}
          />
        </Grid>
        
        <Grid item xs={12} sm={4}>
          <TextField
            label={t('admin.attributes.sortOrder')}
            type="number"
            value={sortOrder}
            onChange={(e) => setSortOrder(parseInt(e.target.value) || 0)}
            variant="outlined"
            size="small"
            fullWidth
            InputProps={{ inputProps: { min: 0 } }}
            disabled={loading || !selectedAttribute}
          />
        </Grid>
        
        <Grid item xs={12}>
          <Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Button
              variant="contained"
              color="primary"
              onClick={handleAddAttribute}
              disabled={loading || !selectedAttribute}
              startIcon={loading ? <CircularProgress size={20} /> : null}
            >
              {t('admin.categoryAttributes.addAttributeButton')}
            </Button>
          </Box>
        </Grid>
      </Grid>
      
      {availableAttributes.length === 0 && !loading && (
        <Typography variant="body2" color="textSecondary" sx={{ mt: 2, textAlign: 'center' }}>
          {t('admin.categoryAttributes.allAttributesMapped')}
        </Typography>
      )}
    </Paper>
  );
};

export default AddAttributeToCategory;