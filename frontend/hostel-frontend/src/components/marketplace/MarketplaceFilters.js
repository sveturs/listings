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
    Slider,
} from '@mui/material';
import { Search, X } from 'lucide-react';
import CompactCategoryTree from './CategoryTree';

const CompactMarketplaceFilters = ({ filters, onFilterChange, categories, selectedCategoryId }) => {
    return (
        <Paper 
            variant="elevation" 
            elevation={3}
            sx={{ 
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                borderRadius: 2,
                overflow: 'hidden',
            }}
        >
            {/* Поиск */}
            <Box sx={{ 
                p: 2, 
                backgroundColor: 'background.default', 
                boxShadow: '0px 1px 2px rgba(0, 0, 0, 0.1)',
                zIndex: 1
            }}>
                <TextField
                    fullWidth
                    size="small"
                    placeholder="🔍 Поиск"
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
                    sx={{ 
                        '& .MuiInputBase-root': {
                            height: 40,
                            borderRadius: 2,
                            backgroundColor: 'white',
                            transition: 'all 0.3s ease',
                        },
                        '& .MuiInputBase-root:focus-within': {
                            boxShadow: '0 0 4px rgba(0, 123, 255, 0.5)',
                        }
                    }}
                />
            </Box>

            {/* Категории */}
            <Box sx={{ 
                flex: 1, 
                overflow: 'auto',
                p: 2,
                backgroundColor: 'background.paper',
                '&::-webkit-scrollbar': {
                    width: 8,
                },
                '&::-webkit-scrollbar-thumb': {
                    bgcolor: 'rgba(0, 0, 0, 0.2)',
                    borderRadius: 4,
                },
                '&:hover': {
                    backgroundColor: 'background.default',
                }
            }}>
                {/*категории полной версии */}
                <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                    Категории
                </Typography>
                <CompactCategoryTree
                    categories={categories}
                    selectedId={selectedCategoryId}
                    onSelectCategory={(id) => onFilterChange({ category_id: id })}
                />
            </Box>

            {/* Фильтры */}
            <Box sx={{ 
                p: 2, 
                backgroundColor: 'background.default', 
                borderTop: '1px solid', 
                borderColor: 'divider',
            }}>
                <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                    Цена
                </Typography>
                <Slider
                    value={[filters.min_price || 0, filters.max_price || 100000]}
                    onChange={(e, newValue) => onFilterChange({ min_price: newValue[0], max_price: newValue[1] })}
                    valueLabelDisplay="auto"
                    min={0}
                    max={100000}
                    sx={{ 
                        mt: 1,
                        '& .MuiSlider-thumb': {
                            backgroundColor: 'primary.main',
                            boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.2)',
                        },
                    }}
                />
                <Divider sx={{ my: 2 }} />
                <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                    Сортировка
                </Typography>
                <Select
                    fullWidth
                    size="small"
                    value={filters.sort_by || 'date_desc'}
                    onChange={(e) => onFilterChange({ sort_by: e.target.value })}
                    sx={{ 
                        '& .MuiSelect-select': {
                            py: 1,
                            fontSize: '1rem',
                            backgroundColor: 'white',
                            borderRadius: 2,
                        },
                        '& .MuiOutlinedInput-root': {
                            borderRadius: 2,
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
