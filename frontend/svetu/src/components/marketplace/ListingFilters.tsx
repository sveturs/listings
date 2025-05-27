'use client';

import { useState, useEffect, useCallback } from 'react';
import { 
  Paper, 
  Typography, 
  TextField, 
  Slider, 
  Button, 
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  RadioGroup,
  FormControlLabel,
  Radio,
  Divider
} from '@mui/material';
import { useTranslations } from 'next-intl';

interface FilterProps {
  onFiltersChange: (filters: any) => void;
}

export default function ListingFilters({ onFiltersChange }: FilterProps) {
  const t = useTranslations('marketplace.filters');
  const [priceRange, setPriceRange] = useState<number[]>([0, 5000]);
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [category, setCategory] = useState('');
  const [condition, setCondition] = useState('');
  const [distance, setDistance] = useState('');
  const [location, setLocation] = useState('');

  const handlePriceChange = (event: Event, newValue: number | number[]) => {
    setPriceRange(newValue as number[]);
    if (Array.isArray(newValue)) {
      setMinPrice(newValue[0].toString());
      setMaxPrice(newValue[1].toString());
    }
  };

  const handleMinPriceChange = (value: string) => {
    setMinPrice(value);
    const numValue = parseInt(value) || 0;
    setPriceRange([numValue, priceRange[1]]);
  };

  const handleMaxPriceChange = (value: string) => {
    setMaxPrice(value);
    const numValue = parseInt(value) || 5000;
    setPriceRange([priceRange[0], numValue]);
  };

  const handleApply = useCallback(() => {
    const filters = {
      category_id: category,
      min_price: minPrice,
      max_price: maxPrice,
      condition,
      distance,
      location
    };
    onFiltersChange(filters);
  }, [category, minPrice, maxPrice, condition, distance, location, onFiltersChange]);

  const handleReset = () => {
    setPriceRange([0, 5000]);
    setMinPrice('');
    setMaxPrice('');
    setCategory('');
    setCondition('');
    setDistance('');
    setLocation('');
    onFiltersChange({});
  };

  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h6" gutterBottom>
        {t('title')}
      </Typography>

      <Box sx={{ mt: 3 }}>
        {/* Category Filter */}
        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel>{t('category')}</InputLabel>
          <Select
            value={category}
            label={t('category')}
            onChange={(e) => setCategory(e.target.value)}
          >
            <MenuItem value="">All</MenuItem>
            <MenuItem value="electronics">Electronics</MenuItem>
            <MenuItem value="furniture">Furniture</MenuItem>
            <MenuItem value="clothing">Clothing</MenuItem>
            <MenuItem value="real_estate">Real Estate</MenuItem>
            <MenuItem value="other">Other</MenuItem>
          </Select>
        </FormControl>

        <Divider sx={{ mb: 3 }} />

        {/* Price Filter */}
        <Box sx={{ mb: 3 }}>
          <Typography gutterBottom>{t('price')}</Typography>
          <Box display="flex" gap={2} mb={2}>
            <TextField
              label={t('minPrice')}
              type="number"
              value={minPrice}
              onChange={(e) => handleMinPriceChange(e.target.value)}
              size="small"
              fullWidth
            />
            <TextField
              label={t('maxPrice')}
              type="number"
              value={maxPrice}
              onChange={(e) => handleMaxPriceChange(e.target.value)}
              size="small"
              fullWidth
            />
          </Box>
          <Slider
            value={priceRange}
            onChange={handlePriceChange}
            valueLabelDisplay="auto"
            min={0}
            max={5000}
          />
        </Box>

        <Divider sx={{ mb: 3 }} />

        {/* Distance Filter */}
        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel>{t('distance')}</InputLabel>
          <Select
            value={distance}
            label={t('distance')}
            onChange={(e) => setDistance(e.target.value)}
          >
            <MenuItem value="">{t('anyDistance')}</MenuItem>
            <MenuItem value="5">5 km</MenuItem>
            <MenuItem value="10">10 km</MenuItem>
            <MenuItem value="25">25 km</MenuItem>
            <MenuItem value="50">50 km</MenuItem>
            <MenuItem value="100">100 km</MenuItem>
          </Select>
        </FormControl>

        {/* Location */}
        <TextField
          fullWidth
          label={t('location')}
          variant="outlined"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          sx={{ mb: 3 }}
        />

        <Divider sx={{ mb: 3 }} />

        {/* Condition Filter */}
        <Box sx={{ mb: 3 }}>
          <Typography gutterBottom>{t('condition')}</Typography>
          <RadioGroup
            value={condition}
            onChange={(e) => setCondition(e.target.value)}
          >
            <FormControlLabel value="" control={<Radio />} label={t('any')} />
            <FormControlLabel value="new" control={<Radio />} label={t('new')} />
            <FormControlLabel value="used" control={<Radio />} label={t('used')} />
          </RadioGroup>
        </Box>

        {/* Action Buttons */}
        <Box display="flex" gap={1}>
          <Button 
            variant="contained" 
            fullWidth
            color="primary"
            onClick={handleApply}
          >
            {t('apply')}
          </Button>
          <Button 
            variant="outlined" 
            onClick={handleReset}
          >
            {t('reset')}
          </Button>
        </Box>
      </Box>
    </Paper>
  );
}