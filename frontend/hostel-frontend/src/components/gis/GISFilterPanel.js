import React, { useState } from 'react';
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

const FilterDrawer = styled(Drawer)(({ theme }) => ({
  '& .MuiDrawer-paper': {
    width: 320,
    boxSizing: 'border-box',
    [theme.breakpoints.down('sm')]: {
      width: '100%'
    }
  }
}));

const FilterSection = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
  '&:not(:last-child)': {
    borderBottom: `1px solid ${theme.palette.divider}`
  }
}));

const StyledSlider = styled(Slider)(({ theme }) => ({
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
}));

const PriceInput = styled(TextField)(({ theme }) => ({
  width: 110,
  marginTop: theme.spacing(2),
  '& input': {
    textAlign: 'right'
  }
}));

const GISFilterPanel = ({ 
  open, 
  onClose, 
  filters, 
  onFiltersChange,
  onApplyFilters,
  onResetFilters
}) => {
  const { t } = useTranslation('gis');
  const [localFilters, setLocalFilters] = useState({ ...filters });

  // Reset local filters when props change
  React.useEffect(() => {
    setLocalFilters({ ...filters });
  }, [filters]);

  const handleChange = (key, value) => {
    setLocalFilters(prev => ({
      ...prev,
      [key]: value
    }));
  };

  const handlePriceChange = (event, newValue) => {
    setLocalFilters(prev => ({
      ...prev,
      price: {
        ...prev.price,
        min: newValue[0],
        max: newValue[1]
      }
    }));
  };

  const handleInputPriceChange = (type) => (event) => {
    const value = event.target.value === '' ? 0 : Number(event.target.value);
    setLocalFilters(prev => ({
      ...prev,
      price: {
        ...prev.price,
        [type]: value
      }
    }));
  };

  const handleApply = () => {
    if (onFiltersChange) {
      onFiltersChange(localFilters);
    }
    if (onApplyFilters) {
      onApplyFilters(localFilters);
    }
    onClose();
  };

  const handleReset = () => {
    const resetFilters = {
      price: { min: 0, max: 100000 },
      condition: 'any',
      radius: 0,
      sort: 'newest'
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