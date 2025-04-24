// frontend/hostel-frontend/src/components/marketplace/AttributeFilters.js
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
    // Добавляем состояние для диапазонов атрибутов
    const [attributeRanges, setAttributeRanges] = useState({});

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

        // Специальная обработка для атрибутов недвижимости
        if (['rooms', 'floor', 'total_floors', 'area', 'land_area', 'property_type'].includes(attributeName)) {
            // Если это пустая строка или undefined, удаляем фильтр
            if (value === '' || value === undefined) {
                setAttributeFilters(prev => {
                    const updated = { ...prev };
                    delete updated[attributeName];
                    debouncedFilterChange(updated);
                    return updated;
                });
                return;
            }

            // Проверяем, является ли значение строкой с запятой (диапазон)
            if (typeof value === 'string' && value.includes(',')) {
                // Это уже диапазон, просто обновляем
                setAttributeFilters(prev => {
                    const updated = { ...prev, [attributeName]: value };
                    debouncedFilterChange(updated);
                    return updated;
                });
            } else {
                // Пытаемся конвертировать в число
                let numericValue = parseFloat(value);
                if (!isNaN(numericValue)) {
                    // Это число, добавляем его без диапазона
                    setAttributeFilters(prev => {
                        const updated = { ...prev, [attributeName]: String(numericValue) };
                        debouncedFilterChange(updated);
                        return updated;
                    });
                } else {
                    // Это текст, пытаемся удалить все не-числовые символы
                    const cleanValue = value.replace(/[^\d.,\-]/g, '').replace(',', '.');
                    numericValue = parseFloat(cleanValue);
                    if (!isNaN(numericValue)) {
                        setAttributeFilters(prev => {
                            const updated = { ...prev, [attributeName]: String(numericValue) };
                            debouncedFilterChange(updated);
                            return updated;
                        });
                    } else {
                        // Не удалось преобразовать, добавляем как текст
                        setAttributeFilters(prev => {
                            const updated = { ...prev, [attributeName]: value };
                            debouncedFilterChange(updated);
                            return updated;
                        });
                    }
                }
            }
        } else {
            // Обычная обработка для других атрибутов
            setAttributeFilters(prev => {
                const updated = { ...prev };

                if (value === undefined || value === null || value === '') {
                    delete updated[attributeName];
                } else {
                    updated[attributeName] = value;
                }

                debouncedFilterChange(updated);
                return updated;
            });
        }
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
        
        // Загружаем диапазоны атрибутов
        const fetchAttributeRanges = async () => {
            try {
                const response = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attribute-ranges`);
                console.log("Получены диапазоны атрибутов:", response.data);
                if (response.data?.data) {
                    setAttributeRanges(response.data.data);
                }
            } catch (error) {
                console.error("Ошибка при загрузке диапазонов атрибутов:", error);
            }
        };
        
        // Загружаем атрибуты
        const fetchAttributes = async () => {
            try {
                console.log(`Запрос атрибутов для категории ${categoryId}`);
                
                // Параллельно загружаем диапазоны
                fetchAttributeRanges();
                
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
                    setError(" атрибуты отсутствуют");
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

    useEffect(() => {
        // Дополнительная обработка для правильного отображения фильтров недвижимости
        attributes.forEach(attr => {
            if (['rooms', 'floor', 'total_floors', 'area', 'land_area', 'property_type'].includes(attr.name)) {
                console.log(`Обрабатываем атрибут недвижимости: ${attr.name}`);
                
                // Обработка для numeric-атрибутов
                if (['rooms', 'floor', 'total_floors', 'area', 'land_area'].includes(attr.name)) {
                    // Если тип не number, меняем его
                    if (attr.attribute_type !== "number") {
                        console.warn(`Атрибут ${attr.name} имеет тип ${attr.attribute_type}, но должен быть 'number'`);
                        attr.attribute_type = "number";
                    }
                    
                    // Проверяем наличие диапазонов из реальных данных
                    if (attributeRanges[attr.name]) {
                        console.log(`Используем реальные диапазоны для ${attr.name}:`, attributeRanges[attr.name]);
                        
                        // Создаем или обновляем options
                        let options = {};
                        try {
                            if (attr.options) {
                                options = typeof attr.options === 'string' ?
                                    JSON.parse(attr.options) : attr.options;
                            }
                        } catch (e) {
                            options = {};
                        }
                        
                        // Обновляем значения из реальных данных
                        options.min = attributeRanges[attr.name].min;
                        options.max = attributeRanges[attr.name].max;
                        options.step = attributeRanges[attr.name].step;
                        
                        // Сохраняем обновленные options
                        attr.options = JSON.stringify(options);
                        console.log(`Обновлены опции для ${attr.name}: ${attr.options}`);
                    } else {
                        // Устанавливаем опции по умолчанию, если нет реальных данных
                        let defaultOptions = {};
                        
                        if (attr.name === 'rooms') {
                            defaultOptions = { min: 1, max: 10, step: 1 };
                        } else if (attr.name === 'floor' || attr.name === 'total_floors') {
                            defaultOptions = { min: 1, max: 50, step: 1 };
                        } else if (attr.name === 'area') {
                            defaultOptions = { min: 1, max: 500, step: 0.5 };
                        } else if (attr.name === 'land_area') {
                            defaultOptions = { min: 1, max: 1000, step: 0.5 };
                        }
                        
                        attr.options = JSON.stringify(defaultOptions);
                        console.log(`Установлены значения по умолчанию для атрибута ${attr.name}: ${attr.options}`);
                    }
                } 
                // Специальная обработка для property_type
                else if (attr.name === "property_type") {
                    // Всегда устанавливаем тип select
                    if (attr.attribute_type !== "select") {
                        console.warn(`Атрибут ${attr.name} имеет тип ${attr.attribute_type}, устанавливаем 'select'`);
                        attr.attribute_type = "select";
                    }
                    
                    // Всегда устанавливаем options для property_type
                    const propertyTypes = {
                        values: [
                            "квартира", 
                            "дом", 
                            "комната",
                            "земельный участок",
                            "гараж",
                            "коммерческая"
                        ]
                    };
                    attr.options = JSON.stringify(propertyTypes);
                    console.log(`Установлены значения для типа недвижимости: ${attr.options}`);
                }
                
                // Для атрибута rooms, если он select, а не number
                if (attr.name === "rooms" && attr.attribute_type === "select") {
                    const roomOptions = {
                        values: ["1", "2", "3", "4", "5+"]
                    };
                    attr.options = JSON.stringify(roomOptions);
                    console.log(`Установлены значения для комнат: ${attr.options}`);
                }
            }
        });
    }, [attributes, attributeRanges]);

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

                // Проверяем наличие диапазонов из реальных данных
                if (attributeRanges[attribute.name]) {
                    console.log(`Используем реальные диапазоны для ${attribute.name}:`, attributeRanges[attribute.name]);
                    options.min = attributeRanges[attribute.name].min;
                    options.max = attributeRanges[attribute.name].max;
                    options.step = attributeRanges[attribute.name].step;
                }

                // Задаем адекватные значения по умолчанию в зависимости от типа атрибута
                let min, max, step = options.step || 1;
                let valueSuffix = '';
                let inputAdornment = null;

                // Специальная обработка для атрибутов недвижимости
                if (['rooms', 'floor', 'total_floors', 'area', 'land_area'].includes(attribute.name)) {
                    switch (attribute.name) {
                        case 'rooms':
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 10;
                            break;
                        case 'floor':
                        case 'total_floors':
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 100;
                            break;
                        case 'area':
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 500;
                            step = options.step || 0.5;
                            valueSuffix = ' м²';
                            inputAdornment = "м²";
                            break;
                        case 'land_area':
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 1000;
                            step = options.step || 0.5;
                            valueSuffix = ' сот';
                            inputAdornment = "сот";
                            break;
                    }
                } else {
                    // Обработка других числовых атрибутов
                    switch (attribute.name) {
                        case 'year':
                            min = options.min !== undefined ? Number(options.min) : 1900;
                            max = options.max !== undefined ? Number(options.max) : new Date().getFullYear() + 1;
                            break;
                        case 'mileage':
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 500000;
                            valueSuffix = ' км';
                            inputAdornment = "км";
                            break;
                        case 'engine_capacity':
                            min = options.min !== undefined ? Number(options.min) : 0.1;
                            max = options.max !== undefined ? Number(options.max) : 10;
                            valueSuffix = ' л';
                            inputAdornment = "л";
                            step = 0.1;
                            break;
                        case 'power':
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 1000;
                            valueSuffix = ' л.с.';
                            inputAdornment = "л.с.";
                            break;
                        case 'price':
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 1000000;
                            break;
                        default:
                            min = options.min !== undefined ? Number(options.min) : 0;
                            max = options.max !== undefined ? Number(options.max) : 100;
                    }
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
                            marks={['rooms', 'floor', 'total_floors', 'area', 'land_area'].includes(attribute.name) ? [
                                { value: min, label: min.toString() + valueSuffix },
                                { value: max, label: max.toString() + valueSuffix }
                            ] : undefined}
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
                                    endAdornment: inputAdornment ? (
                                        <InputAdornment position="end">
                                            {inputAdornment}
                                        </InputAdornment>
                                    ) : null,
                                }}
                                sx={{ width: '45%' }}
                                inputProps={{
                                    min,
                                    max,
                                    step,
                                }}
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
                                    endAdornment: inputAdornment ? (
                                        <InputAdornment position="end">
                                            {inputAdornment}
                                        </InputAdornment>
                                    ) : null,
                                }}
                                sx={{ width: '45%' }}
                                inputProps={{
                                    min,
                                    max,
                                    step,
                                }}
                            />
                        </Box>
                    </Box>
                );
            }


            case 'select': {
                let options = [];
                try {
                    console.log(`Парсинг опций для ${attribute.name}, исходные данные:`, attribute.options);

                    // Улучшенная проверка формата options
                    if (typeof attribute.options === 'string') {
                        // Если options - строка JSON
                        try {
                            const parsedOptions = JSON.parse(attribute.options);
                            console.log("Распарсенные опции из строки:", parsedOptions);

                            if (Array.isArray(parsedOptions.values)) {
                                options = parsedOptions.values;
                            } else if (parsedOptions.values) {
                                // Если values существует, но не массив
                                options = [String(parsedOptions.values)];
                            }
                        } catch (parseError) {
                            console.error(`Ошибка при разборе JSON options для ${attribute.name}:`, parseError);
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

                    // Специальная обработка для типа недвижимости, если опции всё еще пусты
                    if (options.length === 0 && attribute.name === 'property_type') {
                        options = ["квартира", "дом", "комната", "земельный участок", "гараж", "коммерческая"];
                        console.log(`Установлены значения по умолчанию для типа недвижимости:`, options);
                    }

                    // Специальная обработка для комнат, если опции всё еще пусты
                    if (options.length === 0 && attribute.name === 'rooms') {
                        options = ["1", "2", "3", "4", "5+"];
                        console.log(`Установлены значения по умолчанию для комнат:`, options);
                    }

                    console.log(`Итоговые опции для ${attribute.name}:`, options);
                } catch (e) {
                    console.error(`Ошибка при обработке опций для ${attribute.name}:`, e);
                    // Устанавливаем значения по умолчанию для известных атрибутов
                    if (attribute.name === 'property_type') {
                        options = ["квартира", "дом", "комната", "земельный участок", "гараж", "коммерческая"];
                    } else if (attribute.name === 'rooms') {
                        options = ["1", "2", "3", "4", "5+"];
                    } else {
                        options = [];
                    }
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