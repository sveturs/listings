import React, { useState, useEffect, useMemo } from 'react';
import { debounce } from 'lodash';
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
} from '@mui/material';
import {
    FilterList as FilterIcon,
    CalendarMonth as CalendarIcon,
    LocationOn as LocationIcon,
    Clear as ClearIcon,
} from '@mui/icons-material';

const AdvancedFilters = ({
    initialFilters = {},
    onFilterChange,
    onSortChange,
    isLoading
}) => {
    const [expanded, setExpanded] = useState(false);
    const [localFilters, setLocalFilters] = useState({
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

    const debouncedFilterChange = useMemo(
        () => debounce((newFilters) => {
            if (onFilterChange) {
                onFilterChange(newFilters);
            }
        }, 500),
        [onFilterChange]
    );

    useEffect(() => {
        return () => {
            debouncedFilterChange.cancel();
        };
    }, [debouncedFilterChange]);

    const handleFilterChange = (field, value) => {
        const newFilters = { ...localFilters, [field]: value };
        setLocalFilters(newFilters);

        if (['start_date', 'end_date', 'accommodation_type', 'sort_by'].includes(field)) {
            if (onFilterChange) {
                onFilterChange(newFilters);
            }
        } else {
            debouncedFilterChange(newFilters);
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
        setLocalFilters(defaultFilters);
        if (onFilterChange) {
            onFilterChange(defaultFilters);
        }
    };
    const today = new Date().toISOString().split('T')[0];

    return (
        <Paper 
            elevation={2} 
            sx={{ 
                p: 3, 
                mb: 3, 
                borderRadius: 4,
                background: 'white' 
            }}
        >
            <Grid container spacing={3}>
                {/* Город */}
                <Grid item xs={12} md={4}>
                    <TextField
                        fullWidth
                        label="Город"
                        variant="outlined"
                        value={localFilters.city}
                        onChange={(e) => handleFilterChange('city', e.target.value)}
                        disabled={isLoading}
                        InputProps={{
                            startAdornment: <LocationIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                            sx: { borderRadius: 3 }
                        }}
                    />
                </Grid>

                {/* Даты */}
                <Grid item xs={12} md={5}>
                    <Grid container spacing={2}>
                        <Grid item xs={6}>
                            <TextField
                                fullWidth
                                label="Дата заезда"
                                type="date"
                                value={localFilters.start_date}
                                onChange={(e) => handleFilterChange('start_date', e.target.value)}
                                disabled={isLoading}
                                InputProps={{
                                    startAdornment: <CalendarIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                                    sx: { borderRadius: 3 }
                                }}
                                InputLabelProps={{ shrink: true }}
                            />
                        </Grid>
                        <Grid item xs={6}>
                            <TextField
                                fullWidth
                                label="Дата выезда"
                                type="date"
                                value={localFilters.end_date}
                                onChange={(e) => handleFilterChange('end_date', e.target.value)}
                                disabled={isLoading}
                                InputProps={{
                                    startAdornment: <CalendarIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                                    sx: { borderRadius: 3 }
                                }}
                                InputLabelProps={{ shrink: true }}
                            />
                        </Grid>
                    </Grid>
                </Grid>

                {/* Кнопки */}
                <Grid item xs={12} md={3}>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                        <Button
                            fullWidth
                            variant="contained"
                            onClick={() => setExpanded(!expanded)}
                            startIcon={<FilterIcon />}
                            disabled={isLoading}
                            sx={{ 
                                borderRadius: 3,
                                height: '56px'
                            }}
                        >
                            {expanded ? 'Скрыть фильтры' : 'Показать фильтры'}
                        </Button>
                        {Object.values(localFilters).some(Boolean) && (
                            <Button
                                variant="outlined"
                                color="error"
                                onClick={clearFilters}
                                disabled={isLoading}
                                sx={{ 
                                    borderRadius: 3,
                                    height: '56px',
                                    minWidth: '56px',
                                    p: 0
                                }}
                            >
                                <ClearIcon />
                            </Button>
                        )}
                    </Box>
                </Grid>

                {/* Расширенные фильтры */}
                {expanded && (
                    <>
                        <Grid item xs={12} md={3}>
                            <TextField
                                fullWidth
                                label="Количество гостей"
                                type="number"
                                value={localFilters.capacity}
                                onChange={(e) => handleFilterChange('capacity', e.target.value)}
                                disabled={isLoading}
                                InputProps={{ sx: { borderRadius: 3 } }}
                            />
                        </Grid>
                        <Grid item xs={12} md={3}>
                            <FormControl fullWidth>
                                <InputLabel>Тип жилья</InputLabel>
                                <Select
                                    value={localFilters.accommodation_type}
                                    onChange={(e) => handleFilterChange('accommodation_type', e.target.value)}
                                    label="Тип жилья"
                                    disabled={isLoading}
                                    sx={{ borderRadius: 3 }}
                                >
                                    <MenuItem value="">Все типы</MenuItem>
                                    <MenuItem value="apartment">Апартаменты</MenuItem>
                                    <MenuItem value="room">Комната</MenuItem>
                                    <MenuItem value="bed">Койко-место</MenuItem>
                                </Select>
                            </FormControl>
                        </Grid>
                        <Grid item xs={12} md={3}>
                            <TextField
                                fullWidth
                                label="Минимальная цена"
                                type="number"
                                value={localFilters.min_price}
                                onChange={(e) => handleFilterChange('min_price', e.target.value)}
                                disabled={isLoading}
                                InputProps={{ sx: { borderRadius: 3 } }}
                            />
                        </Grid>
                        <Grid item xs={12} md={3}>
                            <TextField
                                fullWidth
                                label="Максимальная цена"
                                type="number"
                                value={localFilters.max_price}
                                onChange={(e) => handleFilterChange('max_price', e.target.value)}
                                disabled={isLoading}
                                InputProps={{ sx: { borderRadius: 3 } }}
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <FormControl fullWidth>
                                <InputLabel>Сортировка</InputLabel>
                                <Select
                                    value={localFilters.sort_by}
                                    onChange={(e) => {
                                        handleFilterChange('sort_by', e.target.value);
                                        if (onSortChange) {
                                            onSortChange(e.target.value, localFilters.sort_direction);
                                        }
                                    }}
                                    label="Сортировка"
                                    disabled={isLoading}
                                    sx={{ 
                                        borderRadius: 3,
                                        maxWidth: '300px'
                                    }}
                                >
                                    <MenuItem value="created_at">По дате добавления</MenuItem>
                                    <MenuItem value="price_per_night">По цене</MenuItem>
                                    <MenuItem value="rating">По рейтингу</MenuItem>
                                </Select>
                            </FormControl>
                        </Grid>
                    </>
                )}
            </Grid>
{/* Активные фильтры */}
{Object.values(localFilters).some(Boolean) && (
                <Box sx={{ mt: 2, display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                    {Object.entries(localFilters).map(([key, value]) => {
                        if (!value || key === 'sort_by' || key === 'sort_direction') return null;
                        const labels = {
                            start_date: 'Заезд',
                            end_date: 'Выезд',
                            city: 'Город',
                            country: 'Страна',
                            min_price: 'От',
                            max_price: 'До',
                            capacity: 'Гости',
                            accommodation_type: 'Тип жилья'
                        };
                        return (
                            <Chip
                                key={key}
                                label={`${labels[key]}: ${value}`}
                                onDelete={() => handleFilterChange(key, '')}
                                disabled={isLoading}
                            />
                        );
                    })}
                </Box>
            )}
        </Paper>
    );
};

export default AdvancedFilters;