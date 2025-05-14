import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Divider,
  Paper,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  CircularProgress,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material';
import { ExpandMore as ExpandMoreIcon } from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import { AttributeComponentFactory, CategoryComponentFactory, ComponentRegistry } from './registry';

interface FilterFactoryProps {
  categoryId: number;
  values: Record<string, any>;
  onChange: (values: Record<string, any>) => void;
  onFilter?: () => void;
  expandedByDefault?: boolean;
}

interface CategoryData {
  id: number;
  name: string;
  has_custom_ui: boolean;
  custom_ui_component?: string;
}

/**
 * Фабрика компонентов для фильтров
 * Динамически создает компоненты фильтрации на основе атрибутов категории
 */
const FilterFactory: React.FC<FilterFactoryProps> = ({
  categoryId,
  values,
  onChange,
  onFilter,
  expandedByDefault = true
}) => {
  const { t } = useTranslation();
  const [attributes, setAttributes] = useState<any[]>([]);
  const [category, setCategory] = useState<CategoryData | null>(null);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState(expandedByDefault);

  // Загружаем данные о категории
  useEffect(() => {
    const fetchCategoryData = async () => {
      try {
        setLoading(true);
        // Получаем информацию о категории
        const categoryResponse = await axios.get(`/api/v1/marketplace/categories/${categoryId}`);
        if (categoryResponse.data.data) {
          setCategory(categoryResponse.data.data);
        }

        // Получаем атрибуты категории
        const attributesResponse = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attributes`);
        if (attributesResponse.data.data) {
          // Сортируем атрибуты по порядку
          const sortedAttributes = [...attributesResponse.data.data].sort(
            (a, b) => a.sort_order - b.sort_order
          );
          setAttributes(
            // Фильтруем только атрибуты с флагом is_filterable
            sortedAttributes.filter(attr => attr.is_filterable)
          );
        }
      } catch (error) {
        console.error('Error fetching filter data:', error);
      } finally {
        setLoading(false);
      }
    };

    if (categoryId) {
      fetchCategoryData();
    }
  }, [categoryId]);

  // Обработчик изменения значения атрибута
  const handleAttributeChange = (attributeName: string, value: any) => {
    onChange({
      ...values,
      [attributeName]: value
    });
  };

  // Обработчик очистки фильтров
  const handleClearFilters = () => {
    onChange({});
  };

  const handleToggleExpanded = () => {
    setExpanded(!expanded);
  };

  // Если идет загрузка, показываем индикатор
  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
        <CircularProgress size={24} />
      </Box>
    );
  }

  // Если категория имеет кастомный UI компонент, используем его
  if (category?.has_custom_ui && category?.custom_ui_component) {
    return (
      <Paper sx={{ p: 2, mb: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="subtitle1" fontWeight="medium">
            {t('marketplace.filters.title')}
          </Typography>
          {Object.keys(values).length > 0 && (
            <Button size="small" onClick={handleClearFilters}>
              {t('marketplace.filters.clearAll')}
            </Button>
          )}
        </Box>
        
        <Divider sx={{ mb: 2 }} />
        
        <CategoryComponentFactory
          categoryId={categoryId}
          componentName={category.custom_ui_component}
          values={values}
          onChange={onChange}
        />
        
        {onFilter && (
          <Box sx={{ mt: 2, display: 'flex', justifyContent: 'flex-end' }}>
            <Button 
              variant="contained" 
              color="primary" 
              onClick={onFilter}
              disabled={loading}
            >
              {t('marketplace.filters.apply')}
            </Button>
          </Box>
        )}
      </Paper>
    );
  }

  // Если нет кастомного компонента, генерируем стандартные фильтры
  return (
    <Paper sx={{ mb: 2 }}>
      <Accordion expanded={expanded} onChange={handleToggleExpanded}>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography variant="subtitle1" fontWeight="medium">
            {t('marketplace.filters.title')}
          </Typography>
        </AccordionSummary>
        <AccordionDetails>
          {Object.keys(values).length > 0 && (
            <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 2 }}>
              <Button size="small" onClick={handleClearFilters}>
                {t('marketplace.filters.clearAll')}
              </Button>
            </Box>
          )}
          
          <Grid container spacing={2}>
            {attributes.length > 0 ? (
              attributes.map(attribute => (
                <Grid item xs={12} sm={6} key={attribute.id}>
                  <AttributeComponentFactory
                    attribute={attribute}
                    value={values[attribute.name]}
                    onChange={(value) => handleAttributeChange(attribute.name, value)}
                  />
                </Grid>
              ))
            ) : (
              <Grid item xs={12}>
                <Typography color="textSecondary" align="center" sx={{ py: 2 }}>
                  {t('marketplace.filters.noFilters')}
                </Typography>
              </Grid>
            )}
          </Grid>
          
          {onFilter && attributes.length > 0 && (
            <Box sx={{ mt: 2, display: 'flex', justifyContent: 'flex-end' }}>
              <Button 
                variant="contained" 
                color="primary" 
                onClick={onFilter}
                disabled={loading}
              >
                {t('marketplace.filters.apply')}
              </Button>
            </Box>
          )}
        </AccordionDetails>
      </Accordion>
    </Paper>
  );
};

export default FilterFactory;