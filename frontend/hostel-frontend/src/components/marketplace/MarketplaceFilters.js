// frontend/hostel-frontend/src/components/marketplace/MarketplaceFilters.js
import React from 'react';
import { useTranslation } from 'react-i18next';

import {
    Paper,
    Box,
    TextField,
    Select,
    MenuItem,
    InputAdornment,
    IconButton,
    Typography,
    Divider,
    Slider,
    Stack,
} from '@mui/material';
import { Search, X } from 'lucide-react';
import CompactCategoryTree from './CategoryTree';
import CategoryFilters from './CategoryFilters';
 
const CompactMarketplaceFilters = ({ filters, onFilterChange, categories, selectedCategoryId }) => {
    const { t } = useTranslation('marketplace', 'common');
    const selectedCategory = categories.find(c => c.id === selectedCategoryId);

    return (
        <Paper variant="elevation" elevation={3} sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            {/* Поиск */}
            <Box sx={{ p: 2, backgroundColor: 'background.default', boxShadow: '0px 1px 2px rgba(0, 0, 0, 0.1)', zIndex: 1 }}>
                <TextField
                    fullWidth
                    size="small"
                    placeholder= {t('buttons.search', { ns: 'common' })}
                    
                    value={filters.query || ''}
                    onChange={(e) => onFilterChange({ query: e.target.value })}
                 

                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <Search size={16} color="gray" />
                            </InputAdornment>
                        ),
                        endAdornment: filters.query && (
                            <InputAdornment position="end">
                                <IconButton 
                                    edge="end" 
                                    size="small" 
                                    onClick={() => onFilterChange({ query: '' })}
                                >
                                    <X size={14} />
                                </IconButton>
                            </InputAdornment>
                        )
                    }}


                />
            </Box>

            {/* Основные фильтры */}
            <Box sx={{ p: 2 }}>
                <Typography variant="subtitle1" gutterBottom>{t('listings.filters.title')}</Typography>
                <Stack spacing={2}>
                    <Box>
                        <Typography gutterBottom>{t('listings.filters.price.label')}</Typography>
                        <Stack direction="row" spacing={1}>
                            <TextField
                                size="small"
                                type="number"
                                placeholder={t('listings.filters.price.min')}
                                value={filters.min_price || ''}
                                onChange={(e) => onFilterChange({ min_price: e.target.value })}
                            />
                            <TextField
                                size="small"
                                type="number"
                                placeholder={t('listings.filters.price.max')}
                                value={filters.max_price || ''}
                                onChange={(e) => onFilterChange({ max_price: e.target.value })}
                            />
                        </Stack>
                    </Box>

                    {/* Специфичные фильтры категории */}
                    {selectedCategory && (
                        <CategoryFilters
                            category={selectedCategory}
                            filters={filters}
                            onFilterChange={onFilterChange}
                        />
                    )}
                </Stack>
            </Box>

            {/* Категории */}
            <Box sx={{ 
                flex: 1, 
                overflow: 'auto',
                p: 2,
                backgroundColor: 'background.paper',
                borderTop: 1,
                borderColor: 'divider'
            }}>
                <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                {t('listings.create.сategories')}
                </Typography>
                <CompactCategoryTree
                    categories={categories}
                    selectedId={selectedCategoryId}
                    onSelectCategory={(id) => onFilterChange({ category_id: id })}
                />
            </Box>
        </Paper>
    );
};

export default CompactMarketplaceFilters;
