//frontend/hostel-frontend/src/components/marketplace/MarketplaceFilters.js
import React from 'react';
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
} from '@mui/material';
import { Search, X } from 'lucide-react';
import CompactCategoryTree from './CategoryTree';

const CompactMarketplaceFilters = ({ filters, onFilterChange, categories, selectedCategoryId }) => {
    return (
        <Paper 
            variant="outlined" 
            sx={{ 
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
            }}
        >
            {/* Поиск */}
            <Box sx={{ p: 1, borderBottom: 1, borderColor: 'divider' }}>
                <TextField
                    fullWidth
                    size="small"
                    placeholder="Поиск"
                    value={filters.query || ''}
                    onChange={(e) => onFilterChange({ query: e.target.value })}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <Search size={16} />
                            </InputAdornment>
                        ),
                        endAdornment: filters.query && (
                            <InputAdornment position="end">
                                <IconButton edge="end" size="small" onClick={() => onFilterChange({ query: '' })}>
                                    <X size={14} />
                                </IconButton>
                            </InputAdornment>
                        )
                    }}
                    sx={{ 
                        '& .MuiInputBase-root': {
                            height: 36
                        }
                    }}
                />
            </Box>

            {/* Категории */}
            <Box sx={{ 
                flex: 1, 
                overflow: 'auto',
                '&::-webkit-scrollbar': {
                    width: 6,
                },
                '&::-webkit-scrollbar-thumb': {
                    bgcolor: 'rgba(0,0,0,0.1)',
                    borderRadius: 3,
                }
            }}>
                <CompactCategoryTree
                    categories={categories}
                    selectedId={selectedCategoryId}
                    onSelectCategory={(id) => onFilterChange({ category_id: id })}
                />
            </Box>

            {/* Фильтры */}
            <Box sx={{ p: 1.5, borderTop: 1, borderColor: 'divider' }}>
                <Typography variant="caption" color="text.secondary" gutterBottom>
                    Цена
                </Typography>
                <Box sx={{ display: 'flex', gap: 1, mb: 1.5 }}>
                    <TextField
                        size="small"
                        placeholder="От"
                        type="number"
                        value={filters.min_price || ''}
                        onChange={(e) => onFilterChange({ min_price: e.target.value })}
                        InputProps={{
                            startAdornment: <InputAdornment position="start">₽</InputAdornment>,
                        }}
                        sx={{ 
                            '& .MuiInputBase-root': {
                                height: 32,
                                fontSize: '0.875rem'
                            }
                        }}
                    />
                    <TextField
                        size="small"
                        placeholder="До"
                        type="number"
                        value={filters.max_price || ''}
                        onChange={(e) => onFilterChange({ max_price: e.target.value })}
                        InputProps={{
                            startAdornment: <InputAdornment position="start">₽</InputAdornment>,
                        }}
                        sx={{ 
                            '& .MuiInputBase-root': {
                                height: 32,
                                fontSize: '0.875rem'
                            }
                        }}
                    />
                </Box>

                <Typography variant="caption" color="text.secondary" gutterBottom>
                    Сортировка
                </Typography>
                <Select
                    fullWidth
                    size="small"
                    value={filters.sort_by || 'date_desc'}
                    onChange={(e) => onFilterChange({ sort_by: e.target.value })}
                    sx={{ 
                        '& .MuiSelect-select': {
                            py: 0.75,
                            fontSize: '0.875rem'
                        }
                    }}
                >
                    <MenuItem value="date_desc">Сначала новые</MenuItem>
                    <MenuItem value="price_asc">Сначала дешевле</MenuItem>
                    <MenuItem value="price_desc">Сначала дороже</MenuItem>
                    <MenuItem value="views">По популярности</MenuItem>
                </Select>
            </Box>
        </Paper>
    );
};

export default CompactMarketplaceFilters;