// frontend/hostel-frontend/src/components/marketplace/AutoPropertiesForm.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  TextField, 
  Grid, 
  MenuItem, 
  FormControl, 
  InputLabel, 
  Select,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Box,
  Divider
} from '@mui/material';
import { ExpandMore } from '@mui/icons-material';
import axios from '../../api/axios';

const AutoPropertiesForm = ({ values, onChange, errors = {} }) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const [constants, setConstants] = useState({
    brands: [],
    fuel_types: [],
    transmissions: [],
    body_types: [],
    drive_types: []
  });
  const [models, setModels] = useState([]);
  const [loading, setLoading] = useState(true);

  // Загрузка констант для автомобильных свойств
// Загрузка констант для автомобильных свойств
useEffect(() => {
    const fetchConstants = async () => {
        try {
            const response = await axios.get('/api/v1/auto/constants');
            if (response.data && response.data.data) {
                setConstants(response.data.data);
            }
            setLoading(false);
        } catch (err) {
            console.error('Ошибка загрузки констант для автомобилей:', err);
            setLoading(false);
        }
    };
    
    fetchConstants();
}, []);

  // Загрузка моделей при изменении марки
  useEffect(() => {
    const fetchModels = async () => {
      if (!values.brand) {
        setModels([]);
        return;
      }
      
      try {
        const response = await axios.get(`/api/v1/auto/models?brand=${encodeURIComponent(values.brand)}`);
        if (response.data && response.data.data && response.data.data.models) {
          setModels(response.data.data.models);
        }
      } catch (err) {
        console.error(`Ошибка загрузки моделей для марки ${values.brand}:`, err);
      }
    };
    
    fetchModels();
  }, [values.brand]);

  const handleChange = (field) => (event) => {
    const newValue = event.target.value;
    
    // Если изменилась марка, сбрасываем модель
    if (field === 'brand') {
      onChange({
        ...values,
        brand: newValue,
        model: ''
      });
    } else {
      onChange({
        ...values,
        [field]: newValue
      });
    }
  };

  const handleNumberChange = (field) => (event) => {
    const value = event.target.value;
    const numberValue = value === '' ? '' : Number(value);
    
    onChange({
      ...values,
      [field]: numberValue
    });
  };

  return (
    <Accordion defaultExpanded>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <Typography variant="h6">
          {t('auto.properties.title', { defaultValue: 'Параметры автомобиля' })}
        </Typography>
      </AccordionSummary>
      <AccordionDetails>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Divider />
            <Box sx={{ my: 2 }}>
              <Typography variant="subtitle1" gutterBottom>
                {t('auto.properties.basic', { defaultValue: 'Основные характеристики' })}
              </Typography>
            </Box>
          </Grid>
          
          {/* Марка */}
          <Grid item xs={12} md={6}>
            <FormControl fullWidth required error={!!errors.brand}>
              <InputLabel>{t('auto.properties.brand', { defaultValue: 'Марка' })}</InputLabel>
              <Select
                value={values.brand || ''}
                onChange={handleChange('brand')}
                label={t('auto.properties.brand', { defaultValue: 'Марка' })}
              >
                <MenuItem value="">
                  <em>{t('auto.properties.select.brand', { defaultValue: 'Выберите марку' })}</em>
                </MenuItem>
                {constants.brands.map((brand) => (
                  <MenuItem key={brand} value={brand}>
                    {brand}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Модель */}
          <Grid item xs={12} md={6}>
            <FormControl fullWidth required error={!!errors.model} disabled={!values.brand}>
              <InputLabel>{t('auto.properties.model', { defaultValue: 'Модель' })}</InputLabel>
              <Select
                value={values.model || ''}
                onChange={handleChange('model')}
                label={t('auto.properties.model', { defaultValue: 'Модель' })}
              >
                <MenuItem value="">
                  <em>{t('auto.properties.select.model', { defaultValue: 'Выберите модель' })}</em>
                </MenuItem>
                {models.map((model) => (
                  <MenuItem key={model} value={model}>
                    {model}
                  </MenuItem>
                ))}
                <MenuItem value="other">
                  {t('auto.properties.model.other', { defaultValue: 'Другая модель' })}
                </MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          {/* Год выпуска */}
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              required
              type="number"
              label={t('auto.properties.year', { defaultValue: 'Год выпуска' })}
              value={values.year || ''}
              onChange={handleNumberChange('year')}
              error={!!errors.year}
              helperText={errors.year}
              InputProps={{ inputProps: { min: 1900, max: new Date().getFullYear() } }}
            />
          </Grid>
          
          {/* Пробег */}
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              type="number"
              label={t('auto.properties.mileage', { defaultValue: 'Пробег (км)' })}
              value={values.mileage || ''}
              onChange={handleNumberChange('mileage')}
              error={!!errors.mileage}
              helperText={errors.mileage}
              InputProps={{ inputProps: { min: 0 } }}
            />
          </Grid>
          
          {/* Цвет */}
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              label={t('auto.properties.color', { defaultValue: 'Цвет' })}
              value={values.color || ''}
              onChange={handleChange('color')}
              error={!!errors.color}
              helperText={errors.color}
            />
          </Grid>
          
          <Grid item xs={12}>
            <Box sx={{ my: 2 }}>
              <Typography variant="subtitle1" gutterBottom>
                {t('auto.properties.engine', { defaultValue: 'Двигатель и трансмиссия' })}
              </Typography>
              <Divider />
            </Box>
          </Grid>
          
          {/* Тип топлива */}
          <Grid item xs={12} md={4}>
            <FormControl fullWidth>
              <InputLabel>{t('auto.properties.fuel_type', { defaultValue: 'Тип топлива' })}</InputLabel>
              <Select
                value={values.fuel_type || ''}
                onChange={handleChange('fuel_type')}
                label={t('auto.properties.fuel_type', { defaultValue: 'Тип топлива' })}
              >
                <MenuItem value="">
                  <em>{t('auto.properties.select.fuel_type', { defaultValue: 'Выберите тип топлива' })}</em>
                </MenuItem>
                {constants.fuel_types.map((type) => (
                  <MenuItem key={type} value={type}>
                    {t(`auto.properties.fuel_types.${type.toLowerCase()}`, { defaultValue: type })}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Трансмиссия */}
          <Grid item xs={12} md={4}>
            <FormControl fullWidth>
              <InputLabel>{t('auto.properties.transmission', { defaultValue: 'Трансмиссия' })}</InputLabel>
              <Select
                value={values.transmission || ''}
                onChange={handleChange('transmission')}
                label={t('auto.properties.transmission', { defaultValue: 'Трансмиссия' })}
              >
                <MenuItem value="">
                  <em>{t('auto.properties.select.transmission', { defaultValue: 'Выберите тип трансмиссии' })}</em>
                </MenuItem>
                {constants.transmissions.map((type) => (
                  <MenuItem key={type} value={type}>
                    {t(`auto.properties.transmissions.${type.toLowerCase()}`, { defaultValue: type })}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Объем двигателя */}
          <Grid item xs={12} md={2}>
            <TextField
              fullWidth
              type="number"
              label={t('auto.properties.engine_capacity', { defaultValue: 'Объем (л)' })}
              value={values.engine_capacity || ''}
              onChange={(e) => {
                const value = e.target.value;
                onChange({
                  ...values,
                  engine_capacity: value === '' ? '' : parseFloat(value)
                });
              }}
              error={!!errors.engine_capacity}
              helperText={errors.engine_capacity}
              inputProps={{ step: 0.1, min: 0 }}
            />
          </Grid>
          
          {/* Мощность */}
          <Grid item xs={12} md={2}>
            <TextField
              fullWidth
              type="number"
              label={t('auto.properties.power', { defaultValue: 'Мощность (л.с.)' })}
              value={values.power || ''}
              onChange={handleNumberChange('power')}
              error={!!errors.power}
              helperText={errors.power}
              InputProps={{ inputProps: { min: 0 } }}
            />
          </Grid>
          
          <Grid item xs={12}>
            <Box sx={{ my: 2 }}>
              <Typography variant="subtitle1" gutterBottom>
                {t('auto.properties.body', { defaultValue: 'Кузов и комплектация' })}
              </Typography>
              <Divider />
            </Box>
          </Grid>
          
          {/* Тип кузова */}
          <Grid item xs={12} md={4}>
            <FormControl fullWidth>
              <InputLabel>{t('auto.properties.body_type', { defaultValue: 'Тип кузова' })}</InputLabel>
              <Select
                value={values.body_type || ''}
                onChange={handleChange('body_type')}
                label={t('auto.properties.body_type', { defaultValue: 'Тип кузова' })}
              >
                <MenuItem value="">
                  <em>{t('auto.properties.select.body_type', { defaultValue: 'Выберите тип кузова' })}</em>
                </MenuItem>
                {constants.body_types.map((type) => (
                  <MenuItem key={type} value={type}>
                    {t(`auto.properties.body_types.${type.toLowerCase()}`, { defaultValue: type })}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Привод */}
          <Grid item xs={12} md={4}>
            <FormControl fullWidth>
              <InputLabel>{t('auto.properties.drive_type', { defaultValue: 'Привод' })}</InputLabel>
              <Select
                value={values.drive_type || ''}
                onChange={handleChange('drive_type')}
                label={t('auto.properties.drive_type', { defaultValue: 'Привод' })}
              >
                <MenuItem value="">
                  <em>{t('auto.properties.select.drive_type', { defaultValue: 'Выберите тип привода' })}</em>
                </MenuItem>
                {constants.drive_types.map((type) => (
                  <MenuItem key={type} value={type}>
                    {t(`auto.properties.drive_types.${type.toLowerCase()}`, { defaultValue: type })}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Количество дверей */}
          <Grid item xs={6} md={2}>
            <TextField
              fullWidth
              type="number"
              label={t('auto.properties.doors', { defaultValue: 'Двери' })}
              value={values.number_of_doors || ''}
              onChange={handleNumberChange('number_of_doors')}
              error={!!errors.number_of_doors}
              helperText={errors.number_of_doors}
              InputProps={{ inputProps: { min: 1, max: 6 } }}
            />
          </Grid>
          
          {/* Количество мест */}
          <Grid item xs={6} md={2}>
            <TextField
              fullWidth
              type="number"
              label={t('auto.properties.seats', { defaultValue: 'Места' })}
              value={values.number_of_seats || ''}
              onChange={handleNumberChange('number_of_seats')}
              error={!!errors.number_of_seats}
              helperText={errors.number_of_seats}
              InputProps={{ inputProps: { min: 1, max: 60 } }}
            />
          </Grid>
          
          {/* Дополнительные особенности */}
          <Grid item xs={12}>
            <TextField
              fullWidth
              multiline
              rows={2}
              label={t('auto.properties.features', { defaultValue: 'Дополнительные особенности и комплектация' })}
              value={values.additional_features || ''}
              onChange={handleChange('additional_features')}
              error={!!errors.additional_features}
              helperText={errors.additional_features}
              placeholder={t('auto.properties.features.placeholder', { defaultValue: 'Климат-контроль, кожаный салон, подогрев сидений, панорамная крыша и т.д.' })}
            />
          </Grid>
        </Grid>
      </AccordionDetails>
    </Accordion>
  );
};

export default AutoPropertiesForm;