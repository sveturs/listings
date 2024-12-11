import React, { useState, useEffect } from 'react';
import {
    Paper,
    Typography,
    TextField,
    Box,
    Slider,
    FormControl,
    Select,
    MenuItem,
    Divider,
    FormControlLabel,
    Switch,
    CircularProgress,
    Alert,
    InputLabel  
} from '@mui/material';
import CategoryTree from './CategoryTree';
import axios from '../../api/axios';

const MarketplaceFilters = ({ filters, onFilterChange }) => {
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchCategories = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await axios.get('/api/v1/marketplace/category-tree');
                setCategories(response.data.data || []);
            } catch (err) {
                console.error('Error fetching categories:', err);
                setError('Не удалось загрузить категории');
            } finally {
                setLoading(false);
            }
        };

        fetchCategories();
    }, []);

    const handlePriceChange = (event, newValue) => {
        onFilterChange({
            min_price: newValue[0],
            max_price: newValue[1]
        });
    };

    const conditions = [
        { value: '', label: 'Все' },
        { value: 'new', label: 'Новое' },
        { value: 'used', label: 'Б/у' }
    ];

    const sortOptions = [
        { value: 'date_desc', label: 'Сначала новые' },
        { value: 'date_asc', label: 'Сначала старые' },
        { value: 'price_asc', label: 'Сначала дешевле' },
        { value: 'price_desc', label: 'Сначала дороже' },
        { value: 'views', label: 'По популярности' }
    ];

    return (
        <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
                Фильтры
            </Typography>

            {/* Поиск */}
            <Box sx={{ mb: 3 }}>
                <TextField
                    label="Поиск по объявлениям"
                    fullWidth
                    value={filters.query || ''}
                    onChange={(e) => onFilterChange({ query: e.target.value })}
                    size="small"
                />
            </Box>

            <Divider sx={{ my: 2 }} />

            {/* Категории */}
            {loading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
                    <CircularProgress size={24} />
                </Box>
            ) : error ? (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            ) : (
                <CategoryTree
                    categories={categories}
                    selectedId={filters.category_id}
                    onSelectCategory={(id) => onFilterChange({ category_id: id })}
                />
            )}

            <Divider sx={{ my: 2 }} />

            {/* Цена */}
            <Box sx={{ mb: 3 }}>
                <Typography gutterBottom>Цена</Typography>
                <Box sx={{ px: 1 }}>
                    <Slider
                        value={[
                            Number(filters.min_price) || 0,
                            Number(filters.max_price) || 1000000
                        ]}
                        onChange={handlePriceChange}
                        valueLabelDisplay="auto"
                        min={0}
                        max={1000000}
                        step={1000}
                    />
                    <Box sx={{ display: 'flex', gap: 1 }}>
                        <TextField
                            label="От"
                            type="number"
                            size="small"
                            value={filters.min_price || ''}
                            onChange={(e) => onFilterChange({ min_price: e.target.value })}
                        />
                        <TextField
                            label="До"
                            type="number"
                            size="small"
                            value={filters.max_price || ''}
                            onChange={(e) => onFilterChange({ max_price: e.target.value })}
                        />
                    </Box>
                </Box>
            </Box>

            {/* Состояние */}
            <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Состояние</InputLabel>
                <Select
                    value={filters.condition || ''}
                    onChange={(e) => onFilterChange({ condition: e.target.value })}
                    label="Состояние"
                >
                    {conditions.map(({ value, label }) => (
                        <MenuItem key={value} value={value}>
                            {label}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>

            {/* Местоположение */}
            <Box sx={{ mb: 2 }}>
                <TextField
                    label="Город"
                    fullWidth
                    size="small"
                    value={filters.city || ''}
                    onChange={(e) => onFilterChange({ city: e.target.value })}
                    sx={{ mb: 2 }}
                />
                <TextField
                    label="Страна"
                    fullWidth
                    size="small"
                    value={filters.country || ''}
                    onChange={(e) => onFilterChange({ country: e.target.value })}
                />
            </Box>

            {/* Только с фото */}
            <FormControlLabel
                control={
                    <Switch
                        checked={filters.with_photos || false}
                        onChange={(e) => onFilterChange({ with_photos: e.target.checked })}
                    />
                }
                label="Только с фото"
            />

            <Divider sx={{ my: 2 }} />

            {/* Сортировка */}
            <FormControl fullWidth size="small">
                <InputLabel>Сортировка</InputLabel>
                <Select
                    value={filters.sort_by || 'date_desc'}
                    onChange={(e) => onFilterChange({ sort_by: e.target.value })}
                    label="Сортировка"
                >
                    {sortOptions.map(({ value, label }) => (
                        <MenuItem key={value} value={value}>
                            {label}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>
        </Paper>
    );
};

export default MarketplaceFilters;