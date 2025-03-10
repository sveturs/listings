import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
    Box,
    Typography,
    TextField,
    Slider,
    FormControl,
    FormControlLabel,
    Radio,
    RadioGroup,
    Checkbox,
    FormGroup,
    Accordion,
    AccordionSummary,
    AccordionDetails,
    Divider,
    Stack
} from '@mui/material';
import { ExpandMore } from '@mui/icons-material';
import axios from '../../api/axios';

const AttributeFilters = ({ categoryId, onFilterChange, filters = {} }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [attributes, setAttributes] = useState([]);
    const [loading, setLoading] = useState(false);
    const [attributeFilters, setAttributeFilters] = useState({});

    // Загрузка атрибутов при изменении категории
    useEffect(() => {
        if (!categoryId) {
            setAttributes([]);
            return;
        }

        const fetchAttributes = async () => {
            setLoading(true);
            try {
                const response = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attributes`);
                if (response.data?.data) {
                    // Фильтруем только те атрибуты, которые можно использовать для фильтрации
                    const filterableAttrs = response.data.data.filter(attr => attr.is_filterable);
                    setAttributes(filterableAttrs);
                    
                    // Инициализируем значения фильтров
                    const initialFilters = {};
                    filterableAttrs.forEach(attr => {
                        // Используем существующее значение фильтра, если есть
                        if (filters[attr.name]) {
                            initialFilters[attr.name] = filters[attr.name];
                        }
                    });
                    
                    setAttributeFilters(initialFilters);
                }
            } catch (error) {
                console.error('Error fetching attributes for filters:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchAttributes();
    }, [categoryId, i18n.language, filters]);

    // Обработчик изменения фильтра
    const handleFilterChange = (name, value) => {
        const updatedFilters = {
            ...attributeFilters,
            [name]: value
        };
        
        // Если значение пустое, удаляем фильтр
        if (value === '' || value === null || value === undefined) {
            delete updatedFilters[name];
        }
        
        setAttributeFilters(updatedFilters);
        if (onFilterChange) {
            onFilterChange(updatedFilters);
        }
    };

    // Получение переведенного имени атрибута
    const getTranslatedName = (attribute) => {
        if (!attribute) return '';
        
        if (attribute.translations && attribute.translations[i18n.language]) {
            return attribute.translations[i18n.language];
        }
        
        return attribute.display_name;
    };

    // Рендер фильтра в зависимости от типа атрибута
    const renderFilter = (attribute) => {
        const displayName = getTranslatedName(attribute);
        const attributeName = attribute.name;
        const currentValue = attributeFilters[attributeName] || '';
        
        switch (attribute.attribute_type) {
            case 'text':
                return (
                    <TextField
                        label={displayName}
                        fullWidth
                        size="small"
                        value={currentValue}
                        onChange={(e) => handleFilterChange(attributeName, e.target.value)}
                    />
                );
                
            case 'number': {
                // Извлекаем min, max из options
                let options = {};
                try {
                    if (typeof attribute.options === 'string') {
                        options = JSON.parse(attribute.options);
                    } else if (attribute.options && typeof attribute.options === 'object') {
                        options = attribute.options;
                    }
                } catch (e) {
                    console.error(`Ошибка при парсинге options для ${attribute.name}:`, e);
                    options = {};
                }
                
                const min = options.min !== undefined ? Number(options.min) : 0;
                const max = options.max !== undefined ? Number(options.max) : 100;
                
                // Устанавливаем значения по умолчанию в зависимости от типа атрибута
                let defaultMin = min;
                let defaultMax = max;

                switch (attribute.name) {
                    case 'mileage': // Пробег
                        defaultMin = min !== undefined ? min : 0;
                        defaultMax = max !== undefined ? max : 500000; // до 500,000 км
                        break;
                    case 'year': // Год выпуска
                        defaultMin = min !== undefined ? min : 1900;
                        defaultMax = max !== undefined ? max : new Date().getFullYear();
                        break;
                    case 'engine_capacity': // Объем двигателя
                        defaultMin = min !== undefined ? min : 0.1;
                        defaultMax = max !== undefined ? max : 8.0; // до 8 литров
                        break;
                    case 'price': // Цена
                        defaultMin = min !== undefined ? min : 0;
                        defaultMax = max !== undefined ? max : 10000000; // до 10 миллионов
                        break;
                    default:
                        defaultMin = min !== undefined ? min : 0;
                        defaultMax = max !== undefined ? max : 100;
                }
                
                // Парсим текущее значение из currentValue
                let value = [defaultMin, defaultMax];
                if (currentValue) {
                    try {
                        const parts = currentValue.split(',');
                        if (parts.length === 2) {
                            value = [parseFloat(parts[0]), parseFloat(parts[1])];
                        }
                    } catch (e) {
                        console.error("Ошибка при парсинге значения диапазона:", e);
                    }
                }
                
                return (
                    <Box sx={{ width: '100%', px: 1 }}>
                        <Typography id={`filter-${attribute.id}-label`} gutterBottom>
                            {displayName}
                        </Typography>
                        
                        <Slider
                            value={value}
                            onChange={(e, newValue) => handleFilterChange(
                                attributeName, 
                                `${newValue[0]},${newValue[1]}`
                            )}
                            valueLabelDisplay="auto"
                            min={defaultMin}
                            max={defaultMax}
                            aria-labelledby={`filter-${attribute.id}-label`}
                        />
                        
                        <Stack direction="row" spacing={2} sx={{ mt: 1 }}>
                            <TextField
                                label={t('filters.min', { defaultValue: 'От' })}
                                type="number"
                                size="small"
                                value={value[0]}
                                onChange={(e) => {
                                    const newValue = [Number(e.target.value), value[1]];
                                    handleFilterChange(attributeName, `${newValue[0]},${newValue[1]}`);
                                }}
                                inputProps={{ min: defaultMin, max: value[1] }}
                                fullWidth
                            />
                            <TextField
                                label={t('filters.max', { defaultValue: 'До' })}
                                type="number"
                                size="small"
                                value={value[1]}
                                onChange={(e) => {
                                    const newValue = [value[0], Number(e.target.value)];
                                    handleFilterChange(attributeName, `${newValue[0]},${newValue[1]}`);
                                }}
                                inputProps={{ min: value[0], max: defaultMax }}
                                fullWidth
                            />
                        </Stack>
                    </Box>
                );
            }
                
            case 'select': {
                // Извлекаем options
                let options = [];
                try {
                    // Проверяем формат options
                    if (typeof attribute.options === 'string') {
                        // Если options - строка JSON
                        const parsedOptions = JSON.parse(attribute.options);
                        
                        if (Array.isArray(parsedOptions.values)) {
                            options = parsedOptions.values;
                        } else if (parsedOptions.values) {
                            // Если values существует, но не массив
                            options = [String(parsedOptions.values)];
                        }
                    } else if (attribute.options && typeof attribute.options === 'object') {
                        // Если options уже объект
                        if (Array.isArray(attribute.options.values)) {
                            options = attribute.options.values;
                        } else if (attribute.options.values) {
                            // Если values существует, но не массив
                            options = [String(attribute.options.values)];
                        }
                    }
                } catch (e) {
                    console.error(`Ошибка при обработке опций для ${attribute.name}:`, e);
                    options = [];
                }
                
                return (
                    <FormControl component="fieldset">
                        <Typography>{displayName}</Typography>
                        <RadioGroup
                            value={currentValue}
                            onChange={(e) => handleFilterChange(attributeName, e.target.value)}
                        >
                            <FormControlLabel
                                value=""
                                control={<Radio size="small" />}
                                label={t('listings.filters.any', { defaultValue: 'Любой' })}
                            />
                            {options.length > 0 ? (
                                options.map((option) => (
                                    <FormControlLabel
                                        key={option}
                                        value={option}
                                        control={<Radio size="small" />}
                                        label={option}
                                    />
                                ))
                            ) : (
                                <FormControlLabel
                                    value=""
                                    disabled
                                    control={<Radio size="small" />}
                                    label={t('listings.filters.noOptions', { defaultValue: 'Нет доступных вариантов' })}
                                />
                            )}
                        </RadioGroup>
                    </FormControl>
                );
            }
                
            case 'boolean':
                return (
                    <FormControl component="fieldset">
                        <Typography>{displayName}</Typography>
                        <RadioGroup
                            value={currentValue}
                            onChange={(e) => handleFilterChange(attributeName, e.target.value)}
                        >
                            <FormControlLabel
                                value=""
                                control={<Radio size="small" />}
                                label={t('listings.filters.any', { defaultValue: 'Любой' })}
                            />
                            <FormControlLabel
                                value="true"
                                control={<Radio size="small" />}
                                label={t('common.yes', { defaultValue: 'Да' })}
                            />
                            <FormControlLabel
                                value="false"
                                control={<Radio size="small" />}
                                label={t('common.no', { defaultValue: 'Нет' })}
                            />
                        </RadioGroup>
                    </FormControl>
                );
                
            case 'multiselect': {
                // Извлекаем options
                let options = [];
                try {
                    // Проверяем формат options
                    if (typeof attribute.options === 'string') {
                        // Если options - строка JSON
                        const parsedOptions = JSON.parse(attribute.options);
                        
                        if (Array.isArray(parsedOptions.values)) {
                            options = parsedOptions.values;
                        } else if (parsedOptions.values) {
                            // Если values существует, но не массив
                            options = [String(parsedOptions.values)];
                        }
                    } else if (attribute.options && typeof attribute.options === 'object') {
                        // Если options уже объект
                        if (Array.isArray(attribute.options.values)) {
                            options = attribute.options.values;
                        } else if (attribute.options.values) {
                            // Если values существует, но не массив
                            options = [String(attribute.options.values)];
                        }
                    }
                } catch (e) {
                    console.error(`Ошибка при обработке мульти-опций для ${attribute.name}:`, e);
                    options = [];
                }
                
                // Парсим текущие выбранные значения
                let selectedValues = [];
                if (currentValue) {
                    try {
                        selectedValues = currentValue.split(',');
                    } catch (e) {
                        // Используем пустой массив
                    }
                }
                
                return (
                    <FormControl component="fieldset">
                        <Typography>{displayName}</Typography>
                        <FormGroup>
                            {options.length > 0 ? (
                                options.map((option) => (
                                    <FormControlLabel
                                        key={option}
                                        control={
                                            <Checkbox
                                                size="small"
                                                checked={selectedValues.includes(option)}
                                                onChange={(e) => {
                                                    const newSelected = e.target.checked
                                                        ? [...selectedValues, option]
                                                        : selectedValues.filter(val => val !== option);
                                                    
                                                    handleFilterChange(
                                                        attributeName,
                                                        newSelected.length > 0 ? newSelected.join(',') : ''
                                                    );
                                                }}
                                            />
                                        }
                                        label={option}
                                    />
                                ))
                            ) : (
                                <Typography variant="caption" color="text.secondary">
                                    {t('listings.filters.noOptions', { defaultValue: 'Нет доступных вариантов' })}
                                </Typography>
                            )}
                        </FormGroup>
                    </FormControl>
                );
            }
                                
            default:
                return null;
        }
    };

    if (loading) {
        return <Typography>{t('common.loading', { defaultValue: 'Загрузка' })}...</Typography>;
    }

    if (!categoryId || attributes.length === 0) {
        return null;
    }

    return (
        <Box>
            <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                {t('listings.filters.specific_attributes', { defaultValue: 'Дополнительные фильтры' })}
            </Typography>
            
            {attributes.map((attribute) => (
                <Accordion key={attribute.id} disableGutters>
                    <AccordionSummary expandIcon={<ExpandMore />}>
                        <Typography>{getTranslatedName(attribute)}</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        {renderFilter(attribute)}
                    </AccordionDetails>
                    <Divider />
                </Accordion>
            ))}
        </Box>
    );
};

export default AttributeFilters;