'use client';

import { useState, useCallback, useEffect } from 'react';
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
import { useTranslations, useLocale } from 'next-intl';
import { useParams } from 'next/navigation';
import { categoryService, Category } from '@/services/category.service';

interface FilterProps {
  onFiltersChange: (filters: Record<string, unknown>) => void;
}

export default function ListingFilters({ onFiltersChange }: FilterProps) {
  const t = useTranslations('marketplace.filters');
  const nextIntlLocale = useLocale();
  const params = useParams();
  
  // Get locale from URL params, fallback to next-intl locale
  const locale = (params?.locale as string) || nextIntlLocale || 'en';
  
  // Manual translations object since useTranslations is not working with proper locale
  const manualTranslations = {
    en: {
      title: 'Filters',
      category: 'Category',
      allCategories: 'All Categories',
      price: 'Price',
      minPrice: 'Min',
      maxPrice: 'Max',
      distance: 'Distance',
      anyDistance: 'Any distance',
      location: 'Location',
      condition: 'Condition',
      any: 'Any',
      new: 'New',
      used: 'Used',
      apply: 'Apply Filters',
      reset: 'Reset',
      map: 'Map'
    },
    ru: {
      title: 'Фильтры',
      category: 'Категория',
      allCategories: 'Все категории',
      price: 'Цена',
      minPrice: 'Мин',
      maxPrice: 'Макс',
      distance: 'Расстояние',
      anyDistance: 'Любое расстояние',
      location: 'Местоположение',
      condition: 'Состояние',
      any: 'Любое',
      new: 'Новое',
      used: 'Б/у',
      apply: 'Применить фильтры',
      reset: 'Сбросить',
      map: 'Карта'
    },
    rs: {
      title: 'Filteri',
      category: 'Kategorija',
      allCategories: 'Sve kategorije',
      price: 'Cena',
      minPrice: 'Min',
      maxPrice: 'Maks',
      distance: 'Udaljenost',
      anyDistance: 'Bilo koja udaljenost',
      location: 'Lokacija',
      condition: 'Stanje',
      any: 'Bilo koje',
      new: 'Novo',
      used: 'Polovino',
      apply: 'Primeni filtere',
      reset: 'Resetuj',
      map: 'Mapa'
    }
  };
  
  const translations = manualTranslations[locale as keyof typeof manualTranslations] || manualTranslations.en;
  const [priceRange, setPriceRange] = useState<number[]>([0, 5000]);
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [category, setCategory] = useState('');
  const [condition, setCondition] = useState('');
  const [distance, setDistance] = useState('');
  const [location, setLocation] = useState('');
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadCategories = async () => {
      setLoading(true);
      try {
        const categoryTree = await categoryService.getCategoryTree();
        setCategories(categoryService.flattenCategories(categoryTree));
      } catch (error) {
        console.error('Error loading categories:', error);
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, []);

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
        {translations.title}
      </Typography>

      <Box sx={{ mt: 3 }}>
        {/* Category Filter */}
        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel>{translations.category}</InputLabel>
          <Select
            value={category}
            label={translations.category}
            onChange={(e) => setCategory(e.target.value)}
            disabled={loading}
          >
            <MenuItem value="">{translations.allCategories}</MenuItem>
            {categories.map((cat, index) => (
              <MenuItem key={`cat-${cat.id}-${index}`} value={cat.id.toString()}>
                {'—'.repeat(Math.max(0, cat.level - 1))} {categoryService.getCategoryName(cat, locale)}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <Divider sx={{ mb: 3 }} />

        {/* Price Filter */}
        <Box sx={{ mb: 3 }}>
          <Typography gutterBottom>{translations.price}</Typography>
          <Box display="flex" gap={2} mb={2}>
            <TextField
              label={translations.minPrice}
              type="number"
              value={minPrice}
              onChange={(e) => handleMinPriceChange(e.target.value)}
              size="small"
              fullWidth
            />
            <TextField
              label={translations.maxPrice}
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
          <InputLabel>{translations.distance}</InputLabel>
          <Select
            value={distance}
            label={translations.distance}
            onChange={(e) => setDistance(e.target.value)}
          >
            <MenuItem value="">{translations.anyDistance}</MenuItem>
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
          label={translations.location}
          variant="outlined"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          sx={{ mb: 3 }}
        />

        <Divider sx={{ mb: 3 }} />

        {/* Condition Filter */}
        <Box sx={{ mb: 3 }}>
          <Typography gutterBottom>{translations.condition}</Typography>
          <RadioGroup
            value={condition}
            onChange={(e) => setCondition(e.target.value)}
          >
            <FormControlLabel value="" control={<Radio />} label={translations.any} />
            <FormControlLabel value="new" control={<Radio />} label={translations.new} />
            <FormControlLabel value="used" control={<Radio />} label={translations.used} />
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
            {translations.apply}
          </Button>
          <Button 
            variant="outlined" 
            onClick={handleReset}
          >
            {translations.reset}
          </Button>
        </Box>
      </Box>
    </Paper>
  );
}