import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
    TextField,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    FormControlLabel,
    Switch,
    FormHelperText,
    Checkbox,
    ListItemText,
    OutlinedInput,
    Box,
    Typography,
    Slider
} from '@mui/material';
import axios from '../../api/axios';
import { InputAdornment } from '@mui/material';

const AttributeFields = ({ categoryId, value = [], onChange, error }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [attributes, setAttributes] = useState([]);
    const [loading, setLoading] = useState(false);
    const [values, setValues] = useState(value);

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
                    console.log("Полученные атрибуты для категории:", categoryId, response.data.data);
                    setAttributes(response.data.data);

                    // Инициализируем значения атрибутов только если они еще не установлены
                    // или изменилась категория
                    if (values.length === 0 || values[0]?.attribute_id !== response.data.data[0]?.id) {
                        const initialValues = response.data.data.map(attr => {
                            // Логируем атрибуты для отладки
                            console.log(`Атрибут: ${attr.name}, тип: ${attr.attribute_type}, опции:`, attr.options);

                            // Создаем базовое значение атрибута
                            const attrValue = {
                                attribute_id: attr.id,
                                attribute_name: attr.name,
                                attribute_type: attr.attribute_type,
                                display_name: attr.display_name,
                                value: getDefaultValueForType(attr.attribute_type)
                            };

                            // В зависимости от типа атрибута, добавляем соответствующее поле
                            switch (attr.attribute_type) {
                                case 'text':
                                case 'select':
                                    attrValue.text_value = "";
                                    break;
                                case 'number':
                                    // Используем некоторые разумные значения по умолчанию, но не переопределяем 
                                    // их, если они уже установлены в существующих атрибутах
                                    let defaultValue = 0;

                                    if (attr.options) {
                                        try {
                                            const options = typeof attr.options === 'string'
                                                ? JSON.parse(attr.options)
                                                : attr.options;

                                            if (options.min !== undefined) {
                                                defaultValue = options.min;
                                            }
                                        } catch (e) {
                                            console.error(`Ошибка при разборе options для ${attr.name}:`, e);
                                        }
                                    }

                                    // Особые правила для некоторых атрибутов
                                    if (attr.name === 'year') {
                                        // Используем 2000 год как значение по умолчанию только для нового атрибута
                                        // Но не переопределяем существующее значение
                                        if (!attrValue.value) {
                                            defaultValue = 2000;
                                        } else {
                                            defaultValue = attrValue.value; // Сохраняем существующее значение
                                        }
                                    } else if (attr.name === 'mileage') {
                                        defaultValue = 0; // Пробег по умолчанию
                                    } else if (attr.name === 'engine_capacity') {
                                        defaultValue = 1.6; // Средний объем двигателя
                                    }

                                    attrValue.value = defaultValue;
                                    attrValue.numeric_value = defaultValue;
                                    break;
                                case 'boolean':
                                    attrValue.boolean_value = false;
                                    break;
                            }

                            return attrValue;
                        });

                        setValues(initialValues);
                        if (onChange) onChange(initialValues);
                    }
                }
            } catch (error) {
                console.error('Error fetching attributes:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchAttributes();
    }, [categoryId, i18n.language]); // Удаляем values и onChange из зависимостей

    const getDefaultValueForType = (type) => {
        switch (type) {
            case 'text':
            case 'select':
                return '';
            case 'number':
                return 0;
            case 'boolean':
                return false;
            case 'multiselect':
                return [];
            default:
                return '';
        }
    };

    // Обработчик изменения значения атрибута
    const handleAttributeChange = (attributeId, newValue) => {
        const updatedValues = values.map(attr => {
            if (attr.attribute_id === attributeId) {
                const attribute = attributes.find(a => a.id === attributeId);
                // Сначала создаем копию атрибута и сохраняем новое значение
                const updatedAttr = { ...attr, value: newValue };

                // Устанавливаем правильный тип значения в зависимости от типа атрибута
                if (attribute && attribute.attribute_type === 'number') {
                    // Используем более строгую проверку на число
                    let parsedValue = parseFloat(newValue);
                    if (!isNaN(parsedValue)) {
                        updatedAttr.numeric_value = parsedValue;
                        // Обязательно обновляем и типизированное значение, и общее значение
                        updatedAttr.value = parsedValue;

                        // Добавляем специальную логику для года выпуска для предотвращения сброса
                        if (attribute.name === 'year') {
                            // Логирование для отладки
                            console.log(`Меняем год выпуска на: ${parsedValue}`);
                        }
                    } else {
                        console.error(`Ошибка преобразования "${newValue}" в число для атрибута ${attribute.name}`);
                    }
                } else if (attribute && attribute.attribute_type === 'boolean') {
                    updatedAttr.boolean_value = Boolean(newValue);
                    updatedAttr.value = updatedAttr.boolean_value;
                } else {
                    updatedAttr.text_value = String(newValue);
                    updatedAttr.value = updatedAttr.text_value;
                }

                // Обновляем отображаемое значение
                if (attribute && attribute.attribute_type === 'boolean') {
                    updatedAttr.display_value = updatedAttr.boolean_value ? 'true' : 'false';
                } else {
                    updatedAttr.display_value = String(newValue);
                }

                return updatedAttr;
            }
            return attr;
        });

        // Сохраняем обновленные значения
        console.log("Обновленные значения атрибутов:", updatedValues);
        setValues(updatedValues);
        if (onChange) onChange(updatedValues);
    };

    // Получение переведенного имени атрибута
    const getTranslatedName = (attribute) => {
        if (!attribute) return '';

        if (attribute.translations && attribute.translations[i18n.language]) {
            return attribute.translations[i18n.language];
        }

        return attribute.display_name;
    };

    // Рендер поля в зависимости от типа атрибута
    const renderField = (attribute) => {
        const attr = attributes.find(a => a.id === attribute.attribute_id);
        if (!attr) return null;

        const displayName = getTranslatedName(attr);
        const isRequired = attr.is_required;

        // Получаем текущее значение
        const attrValue = attribute.value;

        switch (attr.attribute_type) {
            case 'text':
                return (
                    <TextField
                        label={displayName}
                        fullWidth
                        required={isRequired}
                        value={attrValue || ''}
                        onChange={(e) => handleAttributeChange(attr.id, e.target.value)}
                    />
                );

            case 'number': {
                // Извлекаем min, max, step из options
                let options = {};
                try {
                    options = attr.options ? JSON.parse(attr.options) : {};
                } catch (e) {
                    options = {};
                }

                const min = options.min !== undefined ? options.min : 0;
                const max = options.max !== undefined ? options.max : Number.MAX_SAFE_INTEGER;

                // Специальная обработка шага для объема двигателя
                let step = options.step || 1;
                if (attr.name === 'engine_capacity') {
                    step = 0.1; // Обязательно меняем шаг на 0.1 для объема двигателя
                }

                // Определяем, нужен ли слайдер
                // Для некоторых атрибутов слайдер не удобен
                const useSlider = attr.name === 'year' ||
                    (max - min <= 100) || // Только для небольших диапазонов
                    attr.name === 'engine_capacity';

                // Определяем специфические форматы полей
                let inputAdornment = null;
                let valueSuffix = '';

                if (attr.name === 'mileage') {
                    inputAdornment = "км";
                    valueSuffix = ' км';
                } else if (attr.name === 'engine_capacity') {
                    inputAdornment = "л";
                    valueSuffix = ' л';
                } else if (attr.name === 'power') {
                    inputAdornment = "л.с.";
                    valueSuffix = ' л.с.';
                }

                return (
                    <Box sx={{ width: '100%' }}>
                        <Typography id={`attribute-${attr.id}-label`} gutterBottom>
                            {displayName}{isRequired ? ' *' : ''}
                        </Typography>

                        {useSlider && (
                            <Slider
                                value={parseFloat(attrValue) || min}
                                onChange={(e, newValue) => handleAttributeChange(attr.id, newValue)}
                                aria-labelledby={`attribute-${attr.id}-label`}
                                min={min}
                                max={max}
                                step={step}
                                marks={[
                                    { value: min, label: min.toString() + valueSuffix },
                                    { value: max, label: max.toString() + valueSuffix }
                                ]}
                                valueLabelDisplay="auto"
                                valueLabelFormat={(value) => value + valueSuffix}
                                sx={{ mb: 2 }}
                            />
                        )}

                        <TextField
                            type="number"
                            fullWidth
                            required={isRequired}
                            value={attrValue || ''}
                            onChange={(e) => handleAttributeChange(attr.id, parseFloat(e.target.value) || 0)}
                            inputProps={{
                                min,
                                max,
                                step,
                                // Для объема двигателя позволяем дробные значения
                                inputMode: attr.name === 'engine_capacity' ? 'decimal' : 'numeric'
                            }}
                            InputProps={inputAdornment ? {
                                endAdornment: (
                                    <InputAdornment position="end">
                                        {inputAdornment}
                                    </InputAdornment>
                                ),
                            } : undefined}
                        />
                    </Box>
                );
            }

            // В методе renderField для select
            case 'select': {
                // Извлекаем options, улучшенная версия
                let options = [];

                try {
                    // Детальное логирование для отладки
                    console.log(`Парсинг опций для ${attr.name}, исходные данные:`, attr.options);

                    // Проверяем формат options
                    if (typeof attr.options === 'string') {
                        // Если options - строка JSON
                        const parsedOptions = JSON.parse(attr.options);
                        console.log("Распарсенные опции из строки:", parsedOptions);

                        if (Array.isArray(parsedOptions.values)) {
                            options = parsedOptions.values;
                        } else if (parsedOptions.values) {
                            // Если values существует, но не массив
                            options = [String(parsedOptions.values)];
                        }
                    } else if (attr.options && typeof attr.options === 'object') {
                        // Если options уже объект
                        console.log("Опции уже в формате объекта:", attr.options);

                        if (attr.options.values) {
                            if (Array.isArray(attr.options.values)) {
                                options = attr.options.values;
                            } else {
                                // Если values существует, но не массив
                                options = [String(attr.options.values)];
                            }
                        }
                    }

                    console.log(`Итоговые опции для ${attr.name}:`, options);
                } catch (e) {
                    console.error(`Ошибка при обработке опций для ${attr.name}:`, e);
                    options = [];
                }

                // Отображаем варианты или сообщение о их отсутствии
                return (
                    <FormControl fullWidth required={isRequired}>
                        <InputLabel>{displayName}</InputLabel>
                        <Select
                            value={attrValue || ''}
                            onChange={(e) => handleAttributeChange(attr.id, e.target.value)}
                            label={displayName}
                        >
                            {options.length > 0 ? (
                                options.map((option) => (
                                    <MenuItem key={option} value={option}>
                                        {option}
                                    </MenuItem>
                                ))
                            ) : (
                                <MenuItem value="" disabled>
                                    Нет доступных вариантов
                                </MenuItem>
                            )}
                        </Select>
                    </FormControl>
                );
            }


            case 'boolean':
                return (
                    <FormControlLabel
                        control={
                            <Switch
                                checked={Boolean(attrValue)}
                                onChange={(e) => handleAttributeChange(attr.id, e.target.checked)}
                            />
                        }
                        label={displayName + (isRequired ? ' *' : '')}
                    />
                );

            case 'multiselect': {
                // Извлекаем options
                let options = [];
                try {
                    const parsedOptions = attr.options ? JSON.parse(attr.options) : {};
                    options = parsedOptions.values || [];
                } catch (e) {
                    options = [];
                }

                // Убедимся, что attrValue - это массив
                const selectedValues = Array.isArray(attrValue) ? attrValue : [];

                return (
                    <FormControl fullWidth required={isRequired}>
                        <InputLabel>{displayName}</InputLabel>
                        <Select
                            multiple
                            value={selectedValues}
                            onChange={(e) => handleAttributeChange(attr.id, e.target.value)}
                            input={<OutlinedInput label={displayName} />}
                            renderValue={(selected) => selected.join(', ')}
                        >
                            {options.map((option) => (
                                <MenuItem key={option} value={option}>
                                    <Checkbox checked={selectedValues.indexOf(option) > -1} />
                                    <ListItemText primary={option} />
                                </MenuItem>
                            ))}
                        </Select>
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
        <Box sx={{ mt: 3 }}>
            <Typography variant="h6" gutterBottom>
                {t('listings.create.specific_attributes')}
            </Typography>

            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                {values.map((attribute) => (
                    <Box key={attribute.attribute_id}>
                        {renderField(attribute)}
                    </Box>
                ))}
            </Box>

            {error && (
                <FormHelperText error>
                    {error}
                </FormHelperText>
            )}
        </Box>
    );
};

export default AttributeFields;