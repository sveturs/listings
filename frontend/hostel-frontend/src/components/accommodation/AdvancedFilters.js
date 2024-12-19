import React, { useState, useEffect } from 'react';
import {
    Box,
    Paper,
    Grid,
    TextField,
    Button,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Chip,
    Typography,
    Collapse,
} from '@mui/material';
import {
    FilterList as FilterIcon,
    ExpandMore as ExpandMoreIcon,
    ExpandLess as ExpandLessIcon,
    Clear as ClearIcon,
} from '@mui/icons-material';

const AdvancedFilters = ({ 
    initialFilters = {}, 
    onFilterChange, 
    onSortChange,
    isLoading 
}) => {
    const [expanded, setExpanded] = useState(false);
    const [filters, setFilters] = useState({
        start_date: '',
        end_date: '',
        city: '',
        country: '',
        min_price: '',
        max_price: '',
        capacity: '',
        accommodation_type: '',
        sort_by: 'created_at',
        sort_direction: 'desc',
        ...initialFilters
    });

    const [activeFilters, setActiveFilters] = useState([]);

    useEffect(() => {
        // Обновляем список активных фильтров при изменении filters
        const newActiveFilters = Object.entries(filters)
            .filter(([key, value]) => value && key !== 'sort_by' && key !== 'sort_direction')
            .map(([key, value]) => ({
                key,
                label: `${getFilterLabel(key)}: ${value}`,
                value
            }));
        setActiveFilters(newActiveFilters);
    }, [filters]);

    const handleFilterChange = (field, value) => {
        const newFilters = { ...filters, [field]: value };
        setFilters(newFilters);
        if (onFilterChange) {
            onFilterChange(newFilters);
        }
    };

    const clearFilters = () => {
        const defaultFilters = {
            start_date: '',
            end_date: '',
            city: '',
            country: '',
            min_price: '',
            max_price: '',
            capacity: '',
            accommodation_type: '',
            sort_by: 'created_at',
            sort_direction: 'desc'
        };
        setFilters(defaultFilters);
        if (onFilterChange) {
            onFilterChange(defaultFilters);
        }
    };

    const getFilterLabel = (key) => {
        const labels = {
            start_date: 'Дата заезда',
            end_date: 'Дата выезда',
            city: 'Город',
            country: 'Страна',
            min_price: 'Мин. цена',
            max_price: 'Макс. цена',
            capacity: 'Количество мест',
            accommodation_type: 'Тип жилья'
        };
        return labels[key] || key;
    };

    const today = new Date().toISOString().split('T')[0];

    return (
        <Paper sx={{ p: 2, mb: 3 }}>
            {/* Основные фильтры */}
            <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} sm={6} md={2}>
                    <TextField
                        label="Дата заезда"
                        type="date"
                        fullWidth
                        size="small"
                        InputLabelProps={{ shrink: true }}
                        value={filters.start_date}
                        onChange={(e) => handleFilterChange('start_date', e.target.value)}
                        inputProps={{ min: today }}
                        disabled={isLoading}
                    />
                </Grid>
                <Grid item xs={12} sm={6} md={2}>
                    <TextField
                        label="Дата выезда"
                        type="date"
                        fullWidth
                        size="small"
                        InputLabelProps={{ shrink: true }}
                        value={filters.end_date}
                        onChange={(e) => handleFilterChange('end_date', e.target.value)}
                        inputProps={{ min: filters.start_date || today }}
                        disabled={isLoading}
                    />
                </Grid>
                <Grid item xs={12} sm={6} md={2}>
                    <TextField
                        label="Город"
                        fullWidth
                        size="small"
                        value={filters.city}
                        onChange={(e) => handleFilterChange('city', e.target.value)}
                        disabled={isLoading}
                    />
                </Grid>

                <Grid item xs={12} md={4}>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                        <Button
                            variant="outlined"
                            onClick={() => setExpanded(!expanded)}
                            startIcon={expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                            disabled={isLoading}
                        >
                            Расширенный поиск
                        </Button>
                        {activeFilters.length > 0 && (
                            <Button
                                variant="outlined"
                                color="error"
                                onClick={clearFilters}
                                startIcon={<ClearIcon />}
                                disabled={isLoading}
                            >
                                Сбросить
                            </Button>
                        )}
                    </Box>
                </Grid>
            </Grid>

            {/* Расширенные фильтры */}
            <Collapse in={expanded} timeout="auto">
                <Grid container spacing={2} sx={{ mt: 2 }}>
                    <Grid item xs={12} sm={6} md={3}>
                        <TextField
                            label="Количество мест"
                            type="number"
                            fullWidth
                            size="small"
                            value={filters.capacity}
                            onChange={(e) => handleFilterChange('capacity', e.target.value)}
                            InputProps={{ inputProps: { min: 1 } }}
                            disabled={isLoading}
                        />
                    </Grid>
                    <Grid item xs={12} sm={6} md={3}>
                        <FormControl fullWidth size="small">
                            <InputLabel>Тип жилья</InputLabel>
                            <Select
                                value={filters.accommodation_type}
                                onChange={(e) => handleFilterChange('accommodation_type', e.target.value)}
                                label="Тип жилья"
                                disabled={isLoading}
                            >
                                <MenuItem value="">Все типы</MenuItem>
                                <MenuItem value="apartment">Апартаменты</MenuItem>
                                <MenuItem value="room">Комната</MenuItem>
                                <MenuItem value="bed">Койко-место</MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                    <Grid item xs={12} sm={6} md={3}>
                        <TextField
                            label="Минимальная цена"
                            type="number"
                            fullWidth
                            size="small"
                            value={filters.min_price}
                            onChange={(e) => handleFilterChange('min_price', e.target.value)}
                            InputProps={{ inputProps: { min: 0 } }}
                            disabled={isLoading}
                        />
                    </Grid>
                    <Grid item xs={12} sm={6} md={3}>
                        <TextField
                            label="Максимальная цена"
                            type="number"
                            fullWidth
                            size="small"
                            value={filters.max_price}
                            onChange={(e) => handleFilterChange('max_price', e.target.value)}
                            InputProps={{ inputProps: { min: filters.min_price || 0 } }}
                            disabled={isLoading}
                        />
                    </Grid>
                </Grid>

                {/* Сортировка */}
                <Box sx={{ mt: 2 }}>
                    <FormControl fullWidth size="small">
                        <InputLabel>Сортировка</InputLabel>
                        <Select
                            value={filters.sort_by}
                            onChange={(e) => {
                                const newFilters = {
                                    ...filters,
                                    sort_by: e.target.value
                                };
                                setFilters(newFilters);
                                if (onSortChange) {
                                    onSortChange(newFilters.sort_by, newFilters.sort_direction);
                                }
                            }}
                            label="Сортировка"
                            disabled={isLoading}
                        >
                            <MenuItem value="created_at">По дате добавления</MenuItem>
                            <MenuItem value="price_per_night">По цене</MenuItem>
                            <MenuItem value="rating">По рейтингу</MenuItem>
                        </Select>
                    </FormControl>
                </Box>
            </Collapse>

            {/* Активные фильтры */}
            {activeFilters.length > 0 && (
                <Box sx={{ mt: 2, display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                    {activeFilters.map((filter, index) => (
                        <Chip
                            key={index}
                            label={filter.label}
                            onDelete={() => handleFilterChange(filter.key, '')}
                            size="small"
                            disabled={isLoading}
                        />
                    ))}
                </Box>
            )}
        </Paper>
    );
};

export default AdvancedFilters;