import React, { useState, useEffect, ChangeEvent } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Drawer,
  Box,
  Typography,
  Divider,
  IconButton,
  Button,
  Slider,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  TextField,
  MenuItem,
  InputAdornment,
  Stack
} from '@mui/material';
import {
  Close as CloseIcon,
  FilterList as FilterIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';

interface PriceRange {
  min: number;
  max: number;
}

interface GISFilters {
  price: PriceRange;
  condition: 'any' | 'new' | 'used';
  radius: number;
  sort: 'newest' | 'priceAsc' | 'priceDesc' | 'distanceAsc' | 'popularityDesc' | 'id';
  sort_direction?: 'asc' | 'desc';
  [key: string]: any;
}

interface GISFilterPanelProps {
  open: boolean;
  onClose: () => void;
  filters: GISFilters;
  onFiltersChange?: (filters: GISFilters) => void;
  onApplyFilters?: (filters: GISFilters) => void;
  onResetFilters?: (filters: GISFilters) => void;
}

interface FilterDrawerProps {
  anchor?: 'left' | 'right' | 'top' | 'bottom';
  open?: boolean;
  onClose?: () => void;
  variant?: 'permanent' | 'persistent' | 'temporary';
  children?: React.ReactNode;
}

const FilterDrawer = styled(Drawer)<FilterDrawerProps>(({ theme }) => ({
  '& .MuiDrawer-paper': {
    width: 320,
    boxSizing: 'border-box',
    [theme.breakpoints.down('sm')]: {
      width: '100%'
    }
  }
})) as React.ComponentType<FilterDrawerProps>;

interface FilterSectionProps {
  children?: React.ReactNode;
}

const FilterSection = styled(Box)<FilterSectionProps>(({ theme }) => ({
  padding: theme.spacing(2),
  '&:not(:last-child)': {
    borderBottom: `1px solid ${theme.palette.divider}`
  }
})) as React.ComponentType<FilterSectionProps>;

interface StyledSliderProps {
  value?: number | number[];
  onChange?: (event: Event, newValue: number | number[]) => void;
  valueLabelDisplay?: 'auto' | 'on' | 'off';
  min?: number;
  max?: number;
  step?: number;
  marks?: boolean | { value: number; label?: React.ReactNode }[];
  'aria-labelledby'?: string;
  valueLabelFormat?: string | ((value: number, index: number) => React.ReactNode);
}

const StyledSlider = styled(Slider)<StyledSliderProps>(({ theme }) => ({
  '& .MuiSlider-valueLabel': {
    fontSize: 12,
    fontWeight: 'normal',
    top: -10,
    backgroundColor: 'unset',
    color: theme.palette.text.primary,
    '&:before': {
      display: 'none',
    },
  },
})) as React.ComponentType<StyledSliderProps>;

interface PriceInputProps {
  label?: React.ReactNode;
  value?: number;
  onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
  type?: string;
  InputLabelProps?: object;
  InputProps?: object;
}

const PriceInput = styled(TextField)<PriceInputProps>(({ theme }) => ({
  width: 110,
  marginTop: theme.spacing(2),
  '& input': {
    textAlign: 'right'
  }
})) as React.ComponentType<PriceInputProps>;

const GISFilterPanel: React.FC<GISFilterPanelProps> = ({ 
  open, 
  onClose, 
  filters, 
  onFiltersChange,
  onApplyFilters,
  onResetFilters
}) => {
  const { t } = useTranslation('gis');
  const [localFilters, setLocalFilters] = useState<GISFilters>({ ...filters });

  // Reset local filters when props change
  useEffect(() => {
    setLocalFilters({ ...filters });
  }, [filters]);

  const handleChange = (key: string, value: any): void => {
    setLocalFilters(prev => ({
      ...prev,
      [key]: value
    }));
  };

  const handlePriceChange = (_event: Event, newValue: number | number[]): void => {
    if (Array.isArray(newValue)) {
      setLocalFilters(prev => ({
        ...prev,
        price: {
          ...prev.price,
          min: newValue[0],
          max: newValue[1]
        }
      }));
    }
  };

  const handleInputPriceChange = (type: 'min' | 'max') => (event: ChangeEvent<HTMLInputElement>): void => {
    const value = event.target.value === '' ? 0 : Number(event.target.value);
    setLocalFilters(prev => ({
      ...prev,
      price: {
        ...prev.price,
        [type]: value
      }
    }));
  };

  const handleApply = (): void => {
    if (onFiltersChange) {
      onFiltersChange(localFilters);
    }
    if (onApplyFilters) {
      onApplyFilters(localFilters);
    }
    onClose();
  };

  const handleReset = (): void => {
    const resetFilters: GISFilters = {
      price: { min: 0, max: 100000 },
      condition: 'any',
      radius: 0,
      sort: 'id',          // Используем id вместо newest для совместимости с OpenSearch
      sort_direction: 'desc' // Новые объявления имеют бóльшие id
    };
    
    setLocalFilters(resetFilters);
    
    if (onResetFilters) {
      onResetFilters(resetFilters);
    }
    
    if (onFiltersChange) {
      onFiltersChange(resetFilters);
    }
  };

  return (
    <FilterDrawer
      anchor="right"
      open={open}
      onClose={onClose}
      variant="temporary"
    >
      <Box display="flex" alignItems="center" justifyContent="space-between" p={2}>
        <Box display="flex" alignItems="center">
          <FilterIcon sx={{ mr: 1 }} />
          <Typography variant="h6">{t('filters.title')}</Typography>
        </Box>
        <IconButton onClick={onClose} size="large">
          <CloseIcon />
        </IconButton>
      </Box>
      
      <Divider />
      
      <Box sx={{ overflow: 'auto', flexGrow: 1 }}>
        {/* Price Range Filter */}
        <FilterSection>
          <Typography id="price-range-slider" gutterBottom>
            {t('filters.price.title')}
          </Typography>
          
          <StyledSlider
            value={[localFilters.price.min, localFilters.price.max]}
            onChange={handlePriceChange}
            valueLabelDisplay="auto"
            min={0}
            max={100000}
            step={1000}
            aria-labelledby="price-range-slider"
          />
          
          <Box display="flex" justifyContent="space-between" mt={2}>
            <PriceInput
              label={t('filters.price.min')}
              value={localFilters.price.min}
              onChange={handleInputPriceChange('min')}
              type="number"
              InputLabelProps={{
                shrink: true,
              }}
              InputProps={{
                endAdornment: <InputAdornment position="end">RSD</InputAdornment>,
              }}
            />
            <PriceInput
              label={t('filters.price.max')}
              value={localFilters.price.max}
              onChange={handleInputPriceChange('max')}
              type="number"
              InputLabelProps={{
                shrink: true,
              }}
              InputProps={{
                endAdornment: <InputAdornment position="end">RSD</InputAdornment>,
              }}
            />
          </Box>
        </FilterSection>
        
        {/* Condition Filter */}
        <FilterSection>
          <FormControl component="fieldset">
            <FormLabel component="legend">
              {t('filters.condition.title')}
            </FormLabel>
            <RadioGroup
              value={localFilters.condition}
              onChange={(e) => handleChange('condition', e.target.value)}
            >
              <FormControlLabel 
                value="any" 
                control={<Radio />} 
                label={t('filters.condition.any')} 
              />
              <FormControlLabel 
                value="new" 
                control={<Radio />} 
                label={t('filters.condition.new')} 
              />
              <FormControlLabel 
                value="used" 
                control={<Radio />} 
                label={t('filters.condition.used')} 
              />
            </RadioGroup>
          </FormControl>
        </FilterSection>
        
        {/* Search Radius Filter */}
        <FilterSection>
          <Typography gutterBottom>
            {t('filters.radius.title')}
          </Typography>
          <Box px={1}>
            <StyledSlider
              value={localFilters.radius}
              onChange={(e, newValue) => handleChange('radius', newValue)}
              valueLabelDisplay="auto"
              min={0}
              max={50}
              marks={[
                { value: 0, label: t('filters.radius.unlimited') },
                { value: 5, label: '5' },
                { value: 10, label: '10' },
                { value: 25, label: '25' },
                { value: 50, label: '50' }
              ]}
              valueLabelFormat={(value) => value > 0 ? `${value} ${t('filters.radius.unit')}` : t('filters.radius.unlimited')}
            />
          </Box>
        </FilterSection>
        
        {/* Sort By Filter */}
        <FilterSection>
          <FormControl fullWidth>
            <FormLabel>{t('filters.sort.title')}</FormLabel>
            <TextField
              select
              value={localFilters.sort}
              onChange={(e) => handleChange('sort', e.target.value)}
              variant="outlined"
              fullWidth
              margin="normal"
            >
              <MenuItem value="newest">{t('filters.sort.newest')}</MenuItem>
              <MenuItem value="priceAsc">{t('filters.sort.priceAsc')}</MenuItem>
              <MenuItem value="priceDesc">{t('filters.sort.priceDesc')}</MenuItem>
              <MenuItem value="distanceAsc">{t('filters.sort.distanceAsc')}</MenuItem>
              <MenuItem value="popularityDesc">{t('filters.sort.popularityDesc')}</MenuItem>
            </TextField>
          </FormControl>
        </FilterSection>
      </Box>
      
      <Divider />
      
      <Box p={2}>
        <Stack direction="row" spacing={2}>
          <Button 
            variant="outlined" 
            onClick={handleReset}
            fullWidth
          >
            {t('filters.resetFilters')}
          </Button>
          <Button 
            variant="contained" 
            color="primary" 
            onClick={handleApply}
            fullWidth
          >
            {t('filters.applyFilters')}
          </Button>
        </Stack>
      </Box>
    </FilterDrawer>
  );
};

export default GISFilterPanel;