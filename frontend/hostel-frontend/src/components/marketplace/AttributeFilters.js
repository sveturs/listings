import React, { useState, useEffect, useCallback, useRef } from 'react';
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
    Grid,
    Paper,
    CircularProgress,
    InputAdornment,
    MenuItem,
    Select,
    InputLabel
} from '@mui/material';
import axios from '../../api/axios';

const AttributeFilters = ({ categoryId, onFilterChange, filters = {}, onAttributesLoaded }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [attributes, setAttributes] = useState([]);
    const [loading, setLoading] = useState(false);
    const [attributeFilters, setAttributeFilters] = useState({ ...filters });
    const [error, setError] = useState(null);
    const isFirstRender = useRef(true);

    // Создаем отложенную функцию обновления фильтров
    const debouncedFilterChange = useCallback(
        debounce((newFilters) => {
            console.log(`Вызываем отложенное обновление фильтров:`, newFilters);
            if (onFilterChange) {
                onFilterChange(newFilters);
            }
        }, 500),
        [onFilterChange]
    );

    // Функция для получения переведенного значения опции атрибута
    const getTranslatedOptionValue = (attribute, value) => {
        // Для отладки
        console.log(`Получение перевода для опции: ${value} атрибута ${attribute?.name}`);
        console.log("Опции переводов:", attribute?.option_translations);

        if (!value || !attribute) return value;

        // Проверяем, есть ли переводы опций
        if (attribute.option_translations &&
            attribute.option_translations[i18n.language]) {

            // Ищем перевод для конкретного значения
            const optionKey = `option_${value}`;
            if (attribute.option_translations[i18n.language][optionKey]) {
                const translatedValue = attribute.option_translations[i18n.language][optionKey];
                console.log(`Найден перевод: ${value} -> ${translatedValue}`);
                return translatedValue;
            }
        }

        return value;
    };

    // Обработчик изменения фильтров по атрибутам
    const handleFilterChange = useCallback((attributeName, value) => {
        console.log(`Изменение атрибута ${attributeName}: ${value}`);

        // Обновляем локальное состояние сразу для отзывчивого UI
        setAttributeFilters(prev => {
            const updatedFilters = { ...prev };

            if (value === undefined || value === null || value === '') {
                delete updatedFilters[attributeName];
            } else {
                updatedFilters[attributeName] = value;
            }

            // Вызываем отложенное обновление родительского компонента
            debouncedFilterChange(updatedFilters);

            return updatedFilters;
        });
    }, [debouncedFilterChange]);

    // Обработчик изменения для диапазонных фильтров
    const handleRangeFilter = useCallback((attributeName, minValue, maxValue) => {
        console.log(`Изменение диапазона атрибута ${attributeName}: ${minValue}-${maxValue}`);

        // Формируем значение диапазона в формате "min,max"
        const rangeValue = `${minValue},${maxValue}`;

        // Обновляем фильтры
        handleFilterChange(attributeName, rangeValue);
    }, [handleFilterChange]);

    // Загрузка атрибутов при изменении категории
    useEffect(() => {
        if (!categoryId) {
            console.log("Нет ID категории, пропускаем запрос атрибутов");
            setAttributes([]);

            // Уведомляем родительский компонент об отсутствии атрибутов
            if (onAttributesLoaded) {
                onAttributesLoaded(false);
            }
            return;
        }

        // Сбрасываем состояние перед новым запросом
        setLoading(true);
        setError(null);

        const fetchAttributes = async () => {
            try {
                console.log(`Запрос атрибутов для категории ${categoryId}`);
                const response = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attributes`);

                if (response.status === 200 && response.data?.data) {
                    // Фильтруем только те атрибуты, которые можно использовать для фильтрации
                    const filterableAttrs = response.data.data.filter(attr => attr.is_filterable);
                    console.log(`Получено ${response.data.data.length} атрибутов, из них фильтруемых: ${filterableAttrs.length}`);

                    // Отладочный вывод для проверки переводов
                    filterableAttrs.forEach(attr => {
                        if (attr.option_translations) {
                            console.log(`Атрибут ${attr.name} имеет переводы опций:`, attr.option_translations);
                            if (attr.option_translations[i18n.language]) {
                                console.log(`Доступны переводы для языка ${i18n.language}:`,
                                    attr.option_translations[i18n.language]);

                                // Проверяем наличие ключей с префиксом "option_"
                                const hasOptionPrefix = Object.keys(attr.option_translations[i18n.language])
                                    .some(key => key.startsWith('option_'));

                                console.log(`Атрибут ${attr.name} имеет ключи с префиксом option_: ${hasOptionPrefix}`);
                            } else {
                                console.log(`Нет переводов для языка ${i18n.language}`);
                            }
                        } else {
                            console.log(`Атрибут ${attr.name} не имеет переводов опций`);
                        }
                    });

                    if (filterableAttrs.length > 0) {
                        setAttributes(filterableAttrs);

                        // Уведомляем родительский компонент о наличии атрибутов
                        if (onAttributesLoaded) {
                            onAttributesLoaded(true);
                        }

                        // Синхронизируем локальное состояние с переданными фильтрами
                        setAttributeFilters({ ...filters });
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
            } catch (error) {
                console.error(`Ошибка при запросе атрибутов для категории ${categoryId}:`, error);
                setError(`Ошибка: ${error.message}`);
                if (onAttributesLoaded) {
                    onAttributesLoaded(false);
                }
            } finally {
                setLoading(false);
            }
        };

        fetchAttributes();
    }, [categoryId, i18n.language, onAttributesLoaded, filters]);

    // Обновление локального состояния при изменении входных фильтров
    useEffect(() => {
        // Пропускаем первый рендер
        if (isFirstRender.current) {
            isFirstRender.current = false;
            return;
        }

        console.log("Обновление локальных фильтров из пропсов:", filters);
        setAttributeFilters({ ...filters });
    }, [filters]);

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
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
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
        return null;
    }

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
        console.log(`Рендеринг фильтра для атрибута: ${attribute.name} (${attribute.attribute_type})`);

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

                // Задаем адекватные значения по умолчанию в зависимости от типа атрибута
                let min, max;

                switch (attribute.name) {
                    case 'year':
                        min = options.min !== undefined ? Number(options.min) : 1900;
                        max = options.max !== undefined ? Number(options.max) : new Date().getFullYear() + 1;
                        break;
                    case 'mileage':
                        min = options.min !== undefined ? Number(options.min) : 0;
                        max = options.max !== undefined ? Number(options.max) : 500000;
                        break;
                    case 'engine_capacity':
                        min = options.min !== undefined ? Number(options.min) : 0.1;
                        max = options.max !== undefined ? Number(options.max) : 10;
                        break;
                    case 'power':
                        min = options.min !== undefined ? Number(options.min) : 0;
                        max = options.max !== undefined ? Number(options.max) : 1000;
                        break;
                    case 'price':
                        min = options.min !== undefined ? Number(options.min) : 0;
                        max = options.max !== undefined ? Number(options.max) : 1000000;
                        break;
                    case 'area':
                    case 'land_area':
                        min = options.min !== undefined ? Number(options.min) : 0;
                        max = options.max !== undefined ? Number(options.max) : 10000;
                        break;
                    default:
                        min = options.min !== undefined ? Number(options.min) : 0;
                        max = options.max !== undefined ? Number(options.max) : 100;
                }

                // Парсим текущее значение из currentValue
                let value = [min, max];
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
                            value={[parseFloat(value[0]) || min, parseFloat(value[1]) || max]}
                            onChange={(e, newValue) => {
                                // Только локальное обновление UI
                                setAttributeFilters(prev => ({
                                    ...prev,
                                    [attributeName]: `${newValue[0]},${newValue[1]}`
                                }));
                            }}
                            onChangeCommitted={(e, newValue) => {
                                // Вызываем обработчик диапазона при отпускании слайдера
                                handleRangeFilter(
                                    attributeName,
                                    newValue[0],
                                    newValue[1]
                                );
                            }}
                            valueLabelDisplay="auto"
                            valueLabelFormat={(value) => value + valueSuffix}
                            min={min}
                            max={max}
                            step={step}
                        />

                        <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
                            <TextField
                                size="small"
                                type="number"
                                value={value[0]}
                                onChange={(e) => {
                                    const newValue = [parseFloat(e.target.value) || min, value[1]];
                                    handleRangeFilter(attributeName, newValue[0], newValue[1]);
                                }}
                                InputProps={{
                                    endAdornment: valueSuffix ? (
                                        <InputAdornment position="end">
                                            {valueSuffix}
                                        </InputAdornment>
                                    ) : null,
                                }}
                                sx={{ width: '45%' }}
                            />
                            <TextField
                                size="small"
                                type="number"
                                value={value[1]}
                                onChange={(e) => {
                                    const newValue = [value[0], parseFloat(e.target.value) || max];
                                    handleRangeFilter(attributeName, newValue[0], newValue[1]);
                                }}
                                InputProps={{
                                    endAdornment: valueSuffix ? (
                                        <InputAdornment position="end">
                                            {valueSuffix}
                                        </InputAdornment>
                                    ) : null,
                                }}
                                sx={{ width: '45%' }}
                            />
                        </Box>
                    </Box>
                );
            }

            case 'select': {
                let options = [];
                try {
                    console.log(`Парсинг опций для ${attribute.name}, исходные данные:`, attribute.options);

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

                        if (attribute.options.values) {
                            if (Array.isArray(attribute.options.values)) {
                                options = attribute.options.values;
                            } else {
                                // Если values существует, но не массив
                                options = [String(attribute.options.values)];
                            }
                        }
                    }

                    console.log(`Итоговые опции для ${attribute.name}:`, options);

                    // Отладка переводов опций
                    options.forEach(option => {
                        const translated = getTranslatedOptionValue(attribute, option);
                        console.log(`Опция ${option} -> ${translated} (${option === translated ? 'не переведено' : 'переведено'})`);
                    });
                } catch (e) {
                    console.error(`Ошибка при обработке опций для ${attribute.name}:`, e);
                    options = [];
                }

                // Отображаем варианты или сообщение о их отсутствии
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
                                        label={getTranslatedOptionValue(attribute, option)}
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