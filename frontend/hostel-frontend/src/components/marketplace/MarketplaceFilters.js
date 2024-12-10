import React, { useState } from 'react';
import { Box, TextField, Button, Typography } from '@mui/material';

const MarketplaceFilters = ({ filters, onFilterChange }) => {
    const [localFilters, setLocalFilters] = useState(filters);

    const handleInputChange = (key, value) => {
        setLocalFilters((prev) => ({ ...prev, [key]: value }));
    };

    const applyFilters = () => {
        onFilterChange(localFilters);
    };

    return (
        <Box sx={{ p: 2, border: '1px solid #ddd', borderRadius: '8px', backgroundColor: '#f9f9f9' }}>
            <Typography variant="h6" gutterBottom>
                Фильтры
            </Typography>

            <Box sx={{ mb: 2 }}>
                <TextField
                    fullWidth
                    label="Поиск"
                    variant="outlined"
                    value={localFilters.query}
                    onChange={(e) => handleInputChange('query', e.target.value)}
                />
            </Box>

            <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
                <TextField
                    fullWidth
                    label="Минимальная цена"
                    variant="outlined"
                    type="number"
                    value={localFilters.min_price}
                    onChange={(e) => handleInputChange('min_price', e.target.value)}
                />
                <TextField
                    fullWidth
                    label="Максимальная цена"
                    variant="outlined"
                    type="number"
                    value={localFilters.max_price}
                    onChange={(e) => handleInputChange('max_price', e.target.value)}
                />
            </Box>

            <Button
                fullWidth
                variant="contained"
                color="primary"
                onClick={applyFilters}
            >
                Применить фильтры
            </Button>
        </Box>
    );
};

export default MarketplaceFilters;
