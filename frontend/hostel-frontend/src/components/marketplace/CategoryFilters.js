// src/components/marketplace/CategoryFilters.js
import React from 'react';
import { ChevronDown } from 'lucide-react';

import {
    Box,
    Slider,
    FormControlLabel,
    Switch,
    TextField,
    Select,
    MenuItem,
    Typography,
    Accordion,
    AccordionSummary,
    AccordionDetails,
    Stack,
  } from '@mui/material';

const categoryFilters = {
'transport': {
    name: 'Транспорт',
    filters: [
        {
            type: 'range',
            key: 'year',
            label: 'Год выпуска',
            min: 1970,
            max: new Date().getFullYear(),
        },
        {
            type: 'select',
            key: 'transmission',
            label: 'Коробка передач',
            options: ['Механика', 'Автомат', 'Робот', 'Вариатор']
        },
        {
            type: 'range',
            key: 'mileage',
            label: 'Пробег',
            min: 0,
            max: 300000,
            step: 5000,
        },
        {
            type: 'select',
            key: 'fuel',
            label: 'Топливо',
            options: ['Бензин', 'Дизель', 'Электро', 'Гибрид']
        },
    ]
},
'real-estate': {
    name: 'Недвижимость',
    filters: [
        {
            type: 'range',
            key: 'area',
            label: 'Площадь, м²',
            min: 0,
            max: 500,
        },
        {
            type: 'range',
            key: 'rooms',
            label: 'Количество комнат',
            min: 1,
            max: 10,
        },
        {
            type: 'select',
            key: 'property_type',
            label: 'Тип недвижимости',
            options: ['Квартира', 'Дом', 'Комната', 'Участок']
        },
        {
            type: 'switch',
            key: 'has_parking',
            label: 'Парковка'
        },
    ]
},
'electronics': {
    name: 'Электроника',
    filters: [
        {
            type: 'select',
            key: 'brand',
            label: 'Бренд',
            options: ['Apple', 'Samsung', 'Xiaomi', 'Другие']
        },
        {
            type: 'range',
            key: 'warranty',
            label: 'Гарантия, мес.',
            min: 0,
            max: 24,
        },
        {
            type: 'switch',
            key: 'official',
            label: 'Официальная продукция'
        }
    ]
},
};

const CategoryFilters = ({ category, filters, onFilterChange }) => {
if (!category) return null;

const categoryConfig = categoryFilters[category.slug];
if (!categoryConfig) return null;

const renderFilter = (filter) => {
    switch (filter.type) {
        case 'range':
            return (
                <Box key={filter.key}>
                    <Typography gutterBottom>{filter.label}</Typography>
                    <Slider
                        value={filters[filter.key] || [filter.min, filter.max]}
                        onChange={(_, value) => onFilterChange({ [filter.key]: value })}
                        valueLabelDisplay="auto"
                        min={filter.min}
                        max={filter.max}
                        step={filter.step || 1}
                    />
                </Box>
            );

        case 'select':
            return (
                <Box key={filter.key}>
                    <Typography gutterBottom>{filter.label}</Typography>
                    <Select
                        fullWidth
                        value={filters[filter.key] || ''}
                        onChange={(e) => onFilterChange({ [filter.key]: e.target.value })}
                    >
                        <MenuItem value="">Все</MenuItem>
                        {filter.options.map(option => (
                            <MenuItem key={option} value={option}>{option}</MenuItem>
                        ))}
                    </Select>
                </Box>
            );

        case 'switch':
            return (
                <FormControlLabel
                    key={filter.key}
                    control={
                        <Switch
                            checked={filters[filter.key] || false}
                            onChange={(e) => onFilterChange({ [filter.key]: e.target.checked })}
                        />
                    }
                    label={filter.label}
                />
            );

        default:
            return null;
    }
};

return (
    <Box sx={{ mt: 2 }}>
        <Accordion defaultExpanded>
            <AccordionSummary expandIcon={<ChevronDown />}>
                <Typography variant="subtitle1">{categoryConfig.name}</Typography>
            </AccordionSummary>
            <AccordionDetails>
                <Stack spacing={2}>
                    {categoryConfig.filters.map(filter => renderFilter(filter))}
                </Stack>
            </AccordionDetails>
        </Accordion>
    </Box>
);
};

export default CategoryFilters;