'use client';

import { useState } from 'react';
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
  MenuItem
} from '@mui/material';
import { useTranslations } from 'next-intl';

export default function ListingFilters() {
  const t = useTranslations('marketplace.filters');
  const [priceRange, setPriceRange] = useState<number[]>([0, 1000]);
  const [category, setCategory] = useState('');

  const handlePriceChange = (event: Event, newValue: number | number[]) => {
    setPriceRange(newValue as number[]);
  };

  const handleReset = () => {
    setPriceRange([0, 1000]);
    setCategory('');
  };

  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h6" gutterBottom>
        {t('title')}
      </Typography>

      <Box sx={{ mt: 3 }}>
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
          </Select>
        </FormControl>

        <Box sx={{ mb: 3 }}>
          <Typography gutterBottom>{t('priceRange')}</Typography>
          <Slider
            value={priceRange}
            onChange={handlePriceChange}
            valueLabelDisplay="auto"
            min={0}
            max={5000}
            sx={{ mt: 4 }}
          />
          <Box display="flex" justifyContent="space-between">
            <Typography variant="body2">${priceRange[0]}</Typography>
            <Typography variant="body2">${priceRange[1]}</Typography>
          </Box>
        </Box>

        <TextField
          fullWidth
          label={t('location')}
          variant="outlined"
          sx={{ mb: 3 }}
        />

        <Box display="flex" gap={1}>
          <Button 
            variant="contained" 
            fullWidth
            color="primary"
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