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
    Divider
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
    }, [categoryId, i18n.language]);

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
                    options = attribute.options ? JSON.parse(attribute.options) : {};
                } catch (e) {
                    options = {};
                }
                
                const min = options.min !== undefined ? options.min : 0;
                const max = options.max !== undefined ? options.max : 100;
                
                // Парсим текущее значение
                let value = [min, max];
                if (currentValue) {
                    try {
                        const parts = currentValue.split(',');
                        if (parts.length === 2) {
                            value = [parseFloat(parts[0]), parseFloat(parts[1])];
                        }
                    } catch (e) {
                        // Используем значения по умолчанию
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
                            min={min}
                            max={max}
                            aria-labelledby={`filter-${attribute.id}-label`}
                        />
                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                            <Typography variant="body2">{min}</Typography>
                            <Typography variant="body2">{max}</Typography>
                        </Box>
                    </Box>
                );
            }
                
            case 'select': {
                // Извлекаем options
                let options = [];
                try {
                    // Добавляем отладочный вывод
                    console.log(`Обработка опций для атрибута ${attribute.name}, исходные данные:`, attribute.options);
                    
                    // Проверяем формат options
                    if (typeof attribute.options === 'string') {
                        // Если options - строка JSON
                        const parsedOptions = JSON.parse(attribute.options);
                        console.log("Распарсенные опции из строки:", parsedOptions);
                        
                        if (Array.isArray(parsedOptions.values)) {
                            options = parsedOptions.values;
                        } else if (parsedOptions.values) {
                            // Если values существует, но не массив
                            options = [String(parsedOptions.values)];
                        }
                    } else if (attribute.options && typeof attribute.options === 'object') {
                        // Если options уже объект
                        console.log("Опции уже в формате объекта:", attribute.options);
                        
                        if (Array.isArray(attribute.options.values)) {
                            options = attribute.options.values;
                        } else if (attribute.options.values) {
                            // Если values существует, но не массив
                            options = [String(attribute.options.values)];
                        }
                    }
                    
                    console.log(`Итоговые опции для ${attribute.name}:`, options);
                } catch (e) {
                    console.error(`Ошибка при обработке опций для ${attribute.name}:`, e, attribute.options);
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
                                label={t('listings.filters.any')}
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
                                label={t('listings.filters.any')}
                            />
                            <FormControlLabel
                                value="true"
                                control={<Radio size="small" />}
                                label={t('common.yes')}
                            />
                            <FormControlLabel
                                value="false"
                                control={<Radio size="small" />}
                                label={t('common.no')}
                            />
                        </RadioGroup>
                    </FormControl>
                );
                
                case 'multiselect': {
                    // Извлекаем options
                    let options = [];
                    try {
                        console.log(`Обработка мульти-опций для атрибута ${attribute.name}, исходные данные:`, attribute.options);
                        
                        // Проверяем формат options
                        if (typeof attribute.options === 'string') {
                            // Если options - строка JSON
                            const parsedOptions = JSON.parse(attribute.options);
                            console.log("Распарсенные мульти-опции из строки:", parsedOptions);
                            
                            if (Array.isArray(parsedOptions.values)) {
                                options = parsedOptions.values;
                            } else if (parsedOptions.values) {
                                // Если values существует, но не массив
                                options = [String(parsedOptions.values)];
                            }
                        } else if (attribute.options && typeof attribute.options === 'object') {
                            // Если options уже объект
                            console.log("Мульти-опции уже в формате объекта:", attribute.options);
                            
                            if (Array.isArray(attribute.options.values)) {
                                options = attribute.options.values;
                            } else if (attribute.options.values) {
                                // Если values существует, но не массив
                                options = [String(attribute.options.values)];
                            }
                        }
                        
                        console.log(`Итоговые мульти-опции для ${attribute.name}:`, options);
                    } catch (e) {
                        console.error(`Ошибка при обработке мульти-опций для ${attribute.name}:`, e, attribute.options);
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
        return <Typography>{t('common.loading')}...</Typography>;
    }

    if (!categoryId || attributes.length === 0) {
        return null;
    }

    return (
        <Box>
            <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                {t('listings.filters.specific_attributes')}
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