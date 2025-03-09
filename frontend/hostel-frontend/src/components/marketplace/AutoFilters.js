// frontend/hostel-frontend/src/components/marketplace/AutoFilters.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Typography,
  Grid,
  TextField,
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  Button,
  Divider,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Slider
} from '@mui/material';
import { ExpandMore } from '@mui/icons-material';
import axios from '../../api/axios';

const AutoFilters = ({ filters, onFilterChange }) => {
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
  const [expanded, setExpanded] = useState(true);
  const [yearRange, setYearRange] = useState([
    filters.year_from || 1990,
    filters.year_to || new Date().getFullYear()
  ]);

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
      if (!filters.brand) {
        setModels([]);
        return;
      }
      
      try {
        const response = await axios.get(`/api/v1/auto/models?brand=${encodeURIComponent(filters.brand)}`);
        if (response.data && response.data.data && response.data.data.models) {
          setModels(response.data.data.models);
        }
      } catch (err) {
        console.error(`Ошибка загрузки моделей для марки ${filters.brand}:`, err);
      }
    };
    
    fetchModels();
  }, [filters.brand]);

  // Обработка изменения фильтра
  const handleFilterChange = (field, value) => {
    onFilterChange({
      ...filters,
      [field]: value
    });
  };

  // Обработка изменения диапазона года
  const handleYearRangeChange = (event, newValue) => {
    setYearRange(newValue);
  };

  // Применение диапазона года к фильтрам
  const applyYearRange = () => {
    onFilterChange({
      ...filters,
      year_from: yearRange[0],
      year_to: yearRange[1]
    });
  };

  // Сброс всех автомобильных фильтров
  const resetAutoFilters = () => {
    const resetFilters = {
      ...filters,
      brand: '',
      model: '',
      year_from: '',
      year_to: '',
      fuel_type: '',
      transmission: '',
      body_type: '',
      drive_type: '',
      mileage_from: '',
      mileage_to: ''
    };
    
    onFilterChange(resetFilters);
    
    // Сбрасываем локальное состояние слайдера
    setYearRange([1990, new Date().getFullYear()]);
  };

  return (
    <Accordion expanded={expanded} onChange={() => setExpanded(!expanded)}>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <Typography variant="subtitle1" fontWeight="medium">
          {t('auto.filters.title', { defaultValue: 'Параметры автомобиля' })}
        </Typography>
      </AccordionSummary>
      <AccordionDetails>
        <Grid container spacing={2}>
          {/* Марка автомобиля */}
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth size="small">
              <InputLabel>{t('auto.properties.brand', { defaultValue: 'Марка' })}</InputLabel>
              <Select
                value={filters.brand || ''}
                onChange={(e) => handleFilterChange('brand', e.target.value)}
                label={t('auto.properties.brand', { defaultValue: 'Марка' })}
              >
                <MenuItem value="">
                  <em>{t('auto.filters.all_brands', { defaultValue: 'Все марки' })}</em>
                </MenuItem>
                {constants.brands.map((brand) => (
                  <MenuItem key={brand} value={brand}>
                    {brand}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Модель автомобиля */}
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth size="small" disabled={!filters.brand}>
              <InputLabel>{t('auto.properties.model', { defaultValue: 'Модель' })}</InputLabel>
              <Select
                value={filters.model || ''}
                onChange={(e) => handleFilterChange('model', e.target.value)}
                label={t('auto.properties.model', { defaultValue: 'Модель' })}
              >
                <MenuItem value="">
                  <em>{t('auto.filters.all_models', { defaultValue: 'Все модели' })}</em>
                </MenuItem>
                {models.map((model) => (
                  <MenuItem key={model} value={model}>
                    {model}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Год выпуска (слайдер) */}
          <Grid item xs={12}>
            <Typography variant="body2" gutterBottom>
              {t('auto.properties.year', { defaultValue: 'Год выпуска' })}: {yearRange[0]} - {yearRange[1]}
            </Typography>
            <Box sx={{ px: 1 }}>
              <Slider
                value={yearRange}
                onChange={handleYearRangeChange}
                onChangeCommitted={applyYearRange}
                valueLabelDisplay="auto"
                min={1980}
                max={new Date().getFullYear()}
                step={1}
              />
            </Box>
          </Grid>
          
          {/* Пробег */}
          <Grid item xs={6}>
            <TextField
              fullWidth
              size="small"
              type="number"
              label={t('auto.filters.mileage_from', { defaultValue: 'Пробег от' })}
              value={filters.mileage_from || ''}
              onChange={(e) => handleFilterChange('mileage_from', e.target.value)}
              InputProps={{ inputProps: { min: 0 } }}
            />
          </Grid>
          
          <Grid item xs={6}>
            <TextField
              fullWidth
              size="small"
              type="number"
              label={t('auto.filters.mileage_to', { defaultValue: 'Пробег до' })}
              value={filters.mileage_to || ''}
              onChange={(e) => handleFilterChange('mileage_to', e.target.value)}
              InputProps={{ inputProps: { min: 0 } }}
            />
          </Grid>
          
          <Grid item xs={12}>
            <Divider sx={{ my: 1 }} />
          </Grid>
          
          {/* Тип топлива */}
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth size="small">
              <InputLabel>{t('auto.properties.fuel_type', { defaultValue: 'Тип топлива' })}</InputLabel>
              <Select
                value={filters.fuel_type || ''}
                onChange={(e) => handleFilterChange('fuel_type', e.target.value)}
                label={t('auto.properties.fuel_type', { defaultValue: 'Тип топлива' })}
              >
                <MenuItem value="">
                  <em>{t('auto.filters.all_fuel_types', { defaultValue: 'Все типы' })}</em>
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
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth size="small">
              <InputLabel>{t('auto.properties.transmission', { defaultValue: 'Трансмиссия' })}</InputLabel>
              <Select
                value={filters.transmission || ''}
                onChange={(e) => handleFilterChange('transmission', e.target.value)}
                label={t('auto.properties.transmission', { defaultValue: 'Трансмиссия' })}
              >
                <MenuItem value="">
                  <em>{t('auto.filters.all_transmissions', { defaultValue: 'Все типы' })}</em>
                </MenuItem>
                {constants.transmissions.map((type) => (
                  <MenuItem key={type} value={type}>
                    {t(`auto.properties.transmissions.${type.toLowerCase()}`, { defaultValue: type })}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Тип кузова */}
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth size="small">
              <InputLabel>{t('auto.properties.body_type', { defaultValue: 'Тип кузова' })}</InputLabel>
              <Select
                value={filters.body_type || ''}
                onChange={(e) => handleFilterChange('body_type', e.target.value)}
                label={t('auto.properties.body_type', { defaultValue: 'Тип кузова' })}
              >
                <MenuItem value="">
                  <em>{t('auto.filters.all_body_types', { defaultValue: 'Все типы' })}</em>
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
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth size="small">
              <InputLabel>{t('auto.properties.drive_type', { defaultValue: 'Привод' })}</InputLabel>
              <Select
                value={filters.drive_type || ''}
                onChange={(e) => handleFilterChange('drive_type', e.target.value)}
                label={t('auto.properties.drive_type', { defaultValue: 'Привод' })}
              >
                <MenuItem value="">
                  <em>{t('auto.filters.all_drive_types', { defaultValue: 'Все типы' })}</em>
                </MenuItem>
                {constants.drive_types.map((type) => (
                  <MenuItem key={type} value={type}>
                    {t(`auto.properties.drive_types.${type.toLowerCase()}`, { defaultValue: type })}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          {/* Кнопка сброса фильтров */}
          <Grid item xs={12}>
            <Box sx={{ mt: 2, display: 'flex', justifyContent: 'center' }}>
              <Button 
                variant="outlined" 
                color="primary" 
                onClick={resetAutoFilters}
                size="small"
              >
                {t('auto.filters.reset', { defaultValue: 'Сбросить фильтры' })}
              </Button>
            </Box>
          </Grid>
        </Grid>
      </AccordionDetails>
    </Accordion>
  );
};

export default AutoFilters;