import React, { useState, useEffect } from 'react';
import { Grid, Typography, FormControl, InputLabel, Select, MenuItem, Box, Divider } from '@mui/material';
import { CategoryUiComponentProps } from '../ComponentRegistry';
import AttributeComponentFactory from '../AttributeComponentFactory';
import axios from '../../../../api/axios';

/**
 * Кастомный компонент для категории "Автомобили"
 * Реализует специфичную форму для автомобильных атрибутов
 */
const AutoCategoryComponent: React.FC<CategoryUiComponentProps> = ({ categoryId, values, onChange }) => {
  const [attributes, setAttributes] = useState<any[]>([]);
  const [makes, setMakes] = useState<string[]>([]);
  const [models, setModels] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  // Получаем атрибуты для категории
  useEffect(() => {
    const fetchAttributes = async () => {
      try {
        const response = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attributes`);
        setAttributes(response.data.data || []);
        
        // Выделяем список марок автомобилей
        const makeAttr = response.data.data.find((attr: any) => attr.name === 'make');
        if (makeAttr && makeAttr.options && makeAttr.options.values) {
          setMakes(makeAttr.options.values);
        }
        
        setLoading(false);
      } catch (error) {
        console.error('Error fetching category attributes:', error);
        setLoading(false);
      }
    };

    fetchAttributes();
  }, [categoryId]);

  // Обновляем список моделей при выборе марки
  useEffect(() => {
    const selectedMake = values.make;
    if (selectedMake) {
      // В реальном приложении здесь был бы запрос к API
      // Сейчас просто имитируем получение моделей
      const getMakeModels = async () => {
        try {
          // Тут можно сделать запрос к API для получения моделей для выбранной марки
          // Например: const response = await axios.get(`/api/v1/car-models?make=${selectedMake}`);
          
          // Для демонстрации генерируем случайные модели
          const demoModels = ['Модель 1', 'Модель 2', 'Модель 3', 'Модель 4', 'Модель 5'];
          setModels(demoModels);
        } catch (error) {
          console.error('Error fetching models:', error);
          setModels([]);
        }
      };

      getMakeModels();
    } else {
      setModels([]);
    }
  }, [values.make]);

  const handleAttributeChange = (name: string, value: any) => {
    onChange({
      ...values,
      [name]: value
    });
  };

  if (loading) {
    return <Typography>Загрузка...</Typography>;
  }

  return (
    <Box>
      <Typography variant="subtitle1" sx={{ mb: 2 }}>
        Информация об автомобиле
      </Typography>
      
      <Divider sx={{ mb: 3 }} />
      
      <Grid container spacing={2}>
        {/* Специальные поля для автомобилей */}
        <Grid item xs={12} md={6}>
          <FormControl fullWidth size="small">
            <InputLabel>Марка</InputLabel>
            <Select
              value={values.make || ''}
              onChange={(e) => handleAttributeChange('make', e.target.value)}
              label="Марка"
            >
              <MenuItem value="">
                <em>Не выбрано</em>
              </MenuItem>
              {makes.map((make) => (
                <MenuItem key={make} value={make}>
                  {make}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <FormControl fullWidth size="small" disabled={!values.make}>
            <InputLabel>Модель</InputLabel>
            <Select
              value={values.model || ''}
              onChange={(e) => handleAttributeChange('model', e.target.value)}
              label="Модель"
            >
              <MenuItem value="">
                <em>Не выбрано</em>
              </MenuItem>
              {models.map((model) => (
                <MenuItem key={model} value={model}>
                  {model}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        
        {/* Остальные атрибуты выводятся динамически */}
        {attributes
          .filter(attr => !['make', 'model'].includes(attr.name))
          .map(attribute => (
            <Grid item xs={12} md={6} key={attribute.id}>
              <AttributeComponentFactory
                attribute={attribute}
                value={values[attribute.name]}
                onChange={(value) => handleAttributeChange(attribute.name, value)}
              />
            </Grid>
          ))
        }
      </Grid>
    </Box>
  );
};

export default AutoCategoryComponent;