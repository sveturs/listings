// frontend/hostel-frontend/src/components/marketplace/AttributeFilters.js
// Изменения для логирования и отладки
import React, { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { debounce } from 'lodash';
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
    Grid,
    Paper,
    Divider,
    Stack,
    CircularProgress
} from '@mui/material';
import axios from '../../api/axios';

const AttributeFilters = ({ categoryId, onFilterChange, filters = {}, onAttributesLoaded }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [attributes, setAttributes] = useState([]);
    const [loading, setLoading] = useState(false);
    const [attributeFilters, setAttributeFilters] = useState(() => ({ ...filters }));
    const [attributeCount, setAttributeCount] = useState(0);
    const [error, setError] = useState(null);
    console.log(`AttributeFilters: Запрашиваем атрибуты для категории ${categoryId}`);
    console.log(`AttributeFilters: Текущие фильтры:`, filters);
    const debouncedFilterChange = useCallback(
        debounce((name, value) => {
            console.log(`Применяем отложенный фильтр ${name}: ${value}`);
            
            // Обновляем локальное состояние фильтров
            setAttributeFilters(prevFilters => {
                const updatedFilters = { ...prevFilters };
                
                if (value === undefined || value === null || value === '') {
                    delete updatedFilters[name];
                } else {
                    updatedFilters[name] = value;
                }

                // Вызываем основной обработчик фильтров
                if (onFilterChange) {
                    onFilterChange(updatedFilters);
                }
                
                return updatedFilters;
            });
        }, 500),
        [onFilterChange]
    );

    // Обработчик изменения фильтров по атрибутам
    const handleFilterChange = useCallback((attributeName, value) => {
        console.log(`Изменение атрибута ${attributeName}: ${value}`);

        // Обновляем UI немедленно
        setAttributeFilters(prevFilters => {
            const updatedFilters = { ...prevFilters };
            
            if (value === undefined || value === null || value === '') {
                delete updatedFilters[attributeName];
            } else {
                updatedFilters[attributeName] = value;
            }
            
            return updatedFilters;
        });
        
        // Отложенно применяем фильтр
        debouncedFilterChange(attributeName, value);
    }, [debouncedFilterChange]);

    // Загрузка атрибутов при изменении категории
    useEffect(() => {
        if (!categoryId) {
            console.log("Нет ID категории, пропускаем запрос атрибутов");
            setAttributes([]);
            setAttributeCount(0);
            // Вызываем колбэк с false, так как атрибутов нет
            if (onAttributesLoaded) {
                onAttributesLoaded(false);
            }
            return;
        }

        console.log(`Загружаем атрибуты для категории ${categoryId}`);
        setLoading(true);
        setError(null);

        const fetchAttributes = async () => {
            console.log(`AttributeFilters: Начинаем fetchAttributes для категории ${categoryId}`);
            try {
                console.log(`Начинаю запрос атрибутов для категории ${categoryId}`);
                const response = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attributes`);
                console.log(`Ответ API для категории ${categoryId}:`, response);
                console.log(`AttributeFilters: Отправляем запрос к API: /api/v1/marketplace/categories/${categoryId}/attributes`);
                if (response.status === 200 && response.data?.data) {
                    // Фильтруем только те атрибуты, которые можно использовать для фильтрации
                    const filterableAttrs = response.data.data.filter(attr => attr.is_filterable);
                    console.log(`Получено ${response.data.data.length} атрибутов, из них фильтруемых: ${filterableAttrs.length}`);
                    
                    if (filterableAttrs.length > 0) {
                        setAttributes(filterableAttrs);
                        setAttributeCount(filterableAttrs.length);
                        console.log(`Атрибуты для категории ${categoryId}:`, filterableAttrs.map(a => a.name).join(', '));
                        
                        // Уведомляем родительский компонент о наличии атрибутов
                        if (onAttributesLoaded) {
                            onAttributesLoaded(true);
                        }

                        // Сохраняем текущие значения атрибутов
                        const currentFilters = {...attributeFilters};
                        
                        // Добавляем значения из переданных фильтров, которые еще не учтены
                        filterableAttrs.forEach(attr => {
                            // Если в переданных фильтрах есть значение и оно еще не сохранено
                            if (filters[attr.name] && currentFilters[attr.name] !== filters[attr.name]) {
                                currentFilters[attr.name] = filters[attr.name];
                            }
                        });

                        // Обновляем состояние только если изменились фильтры
                        if (JSON.stringify(currentFilters) !== JSON.stringify(attributeFilters)) {
                            setAttributeFilters(currentFilters);
                            console.log("Обновлены атрибутные фильтры:", currentFilters);
                        }
                    } else {
                        console.log(`Для категории ${categoryId} нет фильтруемых атрибутов`);
                        if (onAttributesLoaded) {
                            onAttributesLoaded(false);
                        }
                    }
                } else {
                    console.log(`Ответ не содержит данных или данные некорректны`, response.data);
                    setError("Не удалось загрузить атрибуты");
                    if (onAttributesLoaded) {
                        onAttributesLoaded(false);
                    }
                }
                console.log(`AttributeFilters: Получен ответ:`, response);
            } catch (error) {
                console.error(`Ошибка при запросе атрибутов для категории ${categoryId}:`, error);
                setError(`Ошибка: ${error.message}`);
                if (onAttributesLoaded) {
                    onAttributesLoaded(false);
                }
            } finally {
                setLoading(false);
            }
            console.log(`AttributeFilters: useEffect выполнен для категории ${categoryId}`);
        };

        fetchAttributes();
    }, [categoryId, i18n.language, onAttributesLoaded]);

    // Получение переведенного имени атрибута
    const getTranslatedName = (attribute) => {
        if (!attribute) return '';

        if (attribute.translations && attribute.translations[i18n.language]) {
            return attribute.translations[i18n.language];
        }

        return attribute.display_name;
    };

    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
                <CircularProgress size={24} />
                <Typography sx={{ ml: 2 }}>{t('common.loading', { defaultValue: 'Загрузка' })}...</Typography>
            </Box>
        );
    }

    if (error) {
        return (
            <Box sx={{ p: 2, color: 'error.main' }}>
                <Typography>{error}</Typography>
            </Box>
        );
    }

    if (!categoryId || attributes.length === 0) {
        // Вызываем колбэк с false, так как атрибутов нет
        if (onAttributesLoaded) {
            onAttributesLoaded(false);
        }
        return null;
    }
    console.log(`AttributeFilters: Начало useEffect для категории ${categoryId}`);
    // Группировка атрибутов по типам для лучшего представления
    const numberAttrs = attributes.filter(attr => attr.attribute_type === 'number');
    const selectAttrs = attributes.filter(attr => attr.attribute_type === 'select');
    const textAttrs = attributes.filter(attr =>
        attr.attribute_type === 'text' ||
        attr.attribute_type === 'boolean' ||
        !['number', 'select'].includes(attr.attribute_type)
    );

    // Рендер фильтра в зависимости от типа атрибута
    const renderFilter = (attribute) => {
        const displayName = getTranslatedName(attribute);
        const attributeName = attribute.name;
        const currentValue = attributeFilters[attributeName] || '';

        // Остается без изменений, так как эта логика работает правильно
        switch (attribute.attribute_type) {
            case 'text':
                // Код для текстовых атрибутов
                return (
                    <TextField
                        label={displayName}
                        fullWidth
                        size="small"
                        value={currentValue}
                        onChange={(e) => handleFilterChange(attributeName, e.target.value)}
                        onBlur={(e) => debouncedFilterChange(attributeName, e.target.value)}
                    />
                );

            case 'number': {
                // Код для числовых атрибутов
                // Остается без изменений
                // ...
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

                // Определяем специфические форматы полей
                let valueSuffix = '';
                let step = options.step || 1;

                if (attribute.name === 'mileage') {
                    valueSuffix = ' км';
                } else if (attribute.name === 'engine_capacity') {
                    valueSuffix = ' л';
                    step = 0.1;
                } else if (attribute.name === 'power') {
                    valueSuffix = ' л.с.';
                }

                return (
                    <Box sx={{ width: '100%', px: 1 }}>
                        <Typography id={`filter-${attribute.id}-label`} gutterBottom>
                            {displayName}
                        </Typography>

                        <Slider
                            value={value}
                            onChange={(e, newValue) => {
                                // Обновляем только локальный UI
                                setAttributeFilters(prev => ({
                                    ...prev,
                                    [attributeName]: `${newValue[0]},${newValue[1]}`
                                }));
                            }}
                            onChangeCommitted={(e, newValue) => {
                                // Применяем фильтр только после отпускания ползунка
                                handleFilterChange(
                                    attributeName,
                                    `${newValue[0]},${newValue[1]}`
                                );
                            }}
                            valueLabelDisplay="auto"
                            valueLabelFormat={(value) => value + valueSuffix}
                            min={defaultMin}
                            max={defaultMax}
                            step={step}
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
                                    // Обновляем только локальный UI
                                    setAttributeFilters(prev => ({
                                        ...prev,
                                        [attributeName]: `${newValue[0]},${newValue[1]}`
                                    }));
                                }}
                                onBlur={(e) => {
                                    // Применяем фильтр при потере фокуса
                                    const newValue = [Number(e.target.value), value[1]];
                                    handleFilterChange(attributeName, `${newValue[0]},${newValue[1]}`);
                                }}
                                inputProps={{
                                    min: defaultMin,
                                    max: value[1],
                                    step: attribute.name === 'engine_capacity' ? 0.1 : 1
                                }}
                                fullWidth
                            />
                            <TextField
                                label={t('filters.max', { defaultValue: 'До' })}
                                type="number"
                                size="small"
                                value={value[1]}
                                onChange={(e) => {
                                    const newValue = [value[0], Number(e.target.value)];
                                    // Обновляем только локальный UI
                                    setAttributeFilters(prev => ({
                                        ...prev,
                                        [attributeName]: `${newValue[0]},${newValue[1]}`
                                    }));
                                }}
                                onBlur={(e) => {
                                    // Применяем фильтр при потере фокуса
                                    const newValue = [value[0], Number(e.target.value)];
                                    handleFilterChange(attributeName, `${newValue[0]},${newValue[1]}`);
                                }}
                                inputProps={{
                                    min: value[0],
                                    max: defaultMax,
                                    step: attribute.name === 'engine_capacity' ? 0.1 : 1
                                }}
                                fullWidth
                            />
                        </Stack>
                    </Box>
                );
            }

            case 'select': {
                // Код для атрибутов с выбором из списка
                // ...
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
                // Код для булевых атрибутов
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

            default:
                return null;
        }
    };

    return (
        <Box sx={{ mt: 1 }}>
            <Typography variant="h6" sx={{ fontSize: '1rem', mb: 2 }}>
                {t('listings.filters.specific_attributes', { defaultValue: 'Параметры' })} ({attributeCount})
            </Typography>

            {/* Отображаем атрибуты в сетке, сгруппированные по типам */}
            {numberAttrs.length > 0 && (
                <Box sx={{ mb: 3 }}>
                    <Typography variant="subtitle2" gutterBottom>
                        {t('listings.filters.numeric_attributes', { defaultValue: 'Числовые параметры' })}
                    </Typography>
                    <Grid container spacing={2}>
                        {numberAttrs.map(attribute => (
                            <Grid item xs={12} key={attribute.id}>
                                <Paper sx={{ p: 2 }}>
                                    {renderFilter(attribute)}
                                </Paper>
                            </Grid>
                        ))}
                    </Grid>
                </Box>
            )}

            {selectAttrs.length > 0 && (
                <Box sx={{ mb: 3 }}>
                    <Typography variant="subtitle2" gutterBottom>
                        {t('listings.filters.select_attributes', { defaultValue: 'Выбор из списка' })}
                    </Typography>
                    <Grid container spacing={2}>
                        {selectAttrs.map(attribute => (
                            <Grid item xs={12} sm={6} lg={4} key={attribute.id}>
                                <Paper sx={{ p: 2, height: '100%' }}>
                                    {renderFilter(attribute)}
                                </Paper>
                            </Grid>
                        ))}
                    </Grid>
                </Box>
            )}

            {textAttrs.length > 0 && (
                <Box sx={{ mb: 3 }}>
                    <Typography variant="subtitle2" gutterBottom>
                        {t('listings.filters.text_attributes', { defaultValue: 'Текстовые параметры' })}
                    </Typography>
                    <Grid container spacing={2}>
                        {textAttrs.map(attribute => (
                            <Grid item xs={12} sm={6} key={attribute.id}>
                                <Paper sx={{ p: 2, height: '100%' }}>
                                    {renderFilter(attribute)}
                                </Paper>
                            </Grid>
                        ))}
                    </Grid>
                </Box>
            )}
        </Box>
    );
};

export default AttributeFilters;