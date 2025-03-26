import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { debounce } from 'lodash';
import {
    Box,
    Typography,
    TextField,
    Slider,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    FormControlLabel,
    Radio,
    RadioGroup,
    Switch,
    FormHelperText,
    Checkbox,
    ListItemText,
    OutlinedInput,
    FormGroup,
    Grid,
    Paper,
    Stack,
    CircularProgress,
    InputAdornment
} from '@mui/material';
import axios from '../../api/axios';

const AttributeFields = ({ categoryId, value = [], onChange, error }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [attributes, setAttributes] = useState([]);
    const [loading, setLoading] = useState(false);
    const [values, setValues] = useState(value);
    const [requestInProgress, setRequestInProgress] = useState(false);

    // Отслеживаем предыдущую категорию для предотвращения повторных запросов
    const prevCategoryIdRef = useRef(null);
    // Отслеживаем, были ли установлены внешние значения
    const hasSetExternalValues = useRef(false);

    // Загрузка атрибутов при изменении категории
    useEffect(() => {
        // Проверка изменения категории и избежание повторных запросов
        if (!categoryId || (prevCategoryIdRef.current === categoryId && attributes.length > 0) || requestInProgress) {
            return;
        }

        const fetchAttributes = async () => {
            setLoading(true);
            setRequestInProgress(true);

            try {
                console.log(`Запрос атрибутов для категории ${categoryId}`);

                const response = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attributes`);

                if (response.data?.data) {
                    console.log(`Получено ${response.data.data.length} атрибутов для категории ${categoryId}`);
                    setAttributes(response.data.data);

                    // Проверяем, нужно ли сбросить значения атрибутов
                    const needReset = !values.length ||
                        !value.length ||
                        !hasSetExternalValues.current ||
                        !response.data.data.some(attr =>
                            values.some(val => val.attribute_id === attr.id));

                    if (needReset) {
                        console.log("Инициализация значений атрибутов по умолчанию");
                        const initialValues = response.data.data.map(attr => {
                            // Инициализация значения по умолчанию в зависимости от типа
                            let defaultValue = getDefaultValueForType(attr.attribute_type);

                            // Создаем базовую структуру атрибута
                            const attrValue = {
                                attribute_id: attr.id,
                                attribute_name: attr.name,
                                attribute_type: attr.attribute_type,
                                display_name: attr.display_name,
                                value: defaultValue
                            };

                            // Обработка разных типов атрибутов
                            switch (attr.attribute_type) {
                                case 'text':
                                case 'select':
                                    attrValue.text_value = "";
                                    attrValue.display_value = "";
                                    break;

                                case 'number':
                                    let numDefaultValue = 0;

                                    // Пытаемся получить минимальное значение из опций
                                    if (attr.options) {
                                        try {
                                            const options = typeof attr.options === 'string'
                                                ? JSON.parse(attr.options)
                                                : attr.options;

                                            if (options.min !== undefined) {
                                                numDefaultValue = parseFloat(options.min);
                                            }
                                        } catch (e) {
                                            console.error(`AttributeFields: Ошибка при разборе options для ${attr.name}:`, e);
                                        }
                                    }

                                    // Специальные правила по умолчанию для разных атрибутов
                                    if (attr.name === 'year') {
                                        numDefaultValue = new Date().getFullYear();
                                    } else if (attr.name === 'mileage') {
                                        numDefaultValue = 0;
                                    } else if (attr.name === 'engine_capacity') {
                                        numDefaultValue = 1.6;
                                    }

                                    attrValue.value = numDefaultValue;
                                    attrValue.numeric_value = numDefaultValue;

                                    // Формируем отображаемое значение
                                    let displayValue = numDefaultValue.toString();
                                    if (attr.name === 'mileage') {
                                        displayValue += ' км';
                                    } else if (attr.name === 'engine_capacity') {
                                        displayValue += ' л';
                                    } else if (attr.name === 'power') {
                                        displayValue += ' л.с.';
                                    }

                                    attrValue.display_value = displayValue;
                                    break;

                                case 'boolean':
                                    attrValue.boolean_value = false;
                                    attrValue.display_value = 'Нет';
                                    break;

                                default:
                                    attrValue.display_value = "";
                            }

                            console.log(`AttributeFields: Инициализирован атрибут ${attr.name} (${attr.attribute_type}) со значением:`, attrValue.value);
                            return attrValue;
                        });

                        setValues(initialValues);

                        // Отложенный вызов onChange для избежания race condition
                        setTimeout(() => {
                            console.log("AttributeFields: Применяем начальные значения атрибутов:", initialValues);
                            if (onChange) onChange(initialValues);
                        }, 0);
                    }
                }
            } catch (error) {
                console.error('Error fetching attributes:', error);
            } finally {
                setLoading(false);
                setRequestInProgress(false);
                prevCategoryIdRef.current = categoryId;
            }
        };

        fetchAttributes();
    }, [categoryId, i18n.language]);

    // Обработка внешних значений value
    // Проверяем, есть ли входящие данные и отличаются ли они от текущих
    useEffect(() => {
        // Проверяем, есть ли входящие данные и отличаются ли они от текущих
        if (value && value.length > 0) {
            console.log("AttributeFields: Получены внешние значения атрибутов:", value);
            
            // Дополнительная проверка и коррекция входящих значений
            const processedValues = [];
            const seen = {}; // Для отслеживания дубликатов

            value.forEach(attr => {
                // Пропускаем дубликаты атрибутов
                if (seen[attr.attribute_id]) {
                    console.log(`AttributeFields: Пропущен дубликат атрибута ${attr.attribute_name} (ID: ${attr.attribute_id})`);
                    return;
                }
                seen[attr.attribute_id] = true;

                const processed = { ...attr };

                // Исправляем проблемы с полем value
                if (processed.value === undefined || processed.value === null) {
                    // Восстанавливаем значение на основе типизированных полей
                    if (processed.attribute_type === 'text' || processed.attribute_type === 'select') {
                        processed.value = processed.text_value || processed.display_value || '';
                    } else if (processed.attribute_type === 'number') {
                        processed.value = processed.numeric_value !== null ? processed.numeric_value :
                            (processed.display_value ? parseFloat(processed.display_value) : 0);
                    } else if (processed.attribute_type === 'boolean') {
                        processed.value = processed.boolean_value !== null ? processed.boolean_value : false;
                    }

                    console.log(`AttributeFields: Восстановлено отсутствующее значение для атрибута ${processed.attribute_name}: ${processed.value}`);
                }

                processedValues.push(processed);
            });

            console.log("AttributeFields: Обработанные значения атрибутов:", processedValues);
            setValues(processedValues);
            hasSetExternalValues.current = true;
            // }
        }
    }, [value]);

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
        console.log(`AttributeFields: handleAttributeChange вызван для атрибута ${attributeId}, новое значение:`, newValue);

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
                        updatedAttr.value = parsedValue;

                        // Особая обработка для атрибута year (год выпуска)
                        if (attribute.name === 'year') {
                            // Проверяем, что год в разумных пределах
                            const currentYear = new Date().getFullYear();
                            if (parsedValue < 1900 || parsedValue > currentYear + 1) {
                                parsedValue = currentYear;
                                updatedAttr.numeric_value = parsedValue;
                                updatedAttr.value = parsedValue;
                                console.log(`AttributeFields: Корректируем год выпуска на ${parsedValue}`);
                            }
                        }

                        // Для дробных значений округляем до 1-2 знаков после запятой
                        if (attribute.name === 'engine_capacity') {
                            parsedValue = Math.round(parsedValue * 10) / 10; // Округление до 0.1
                            updatedAttr.numeric_value = parsedValue;
                            updatedAttr.value = parsedValue;
                        }
                    } else {
                        console.error(`AttributeFields: Ошибка преобразования "${newValue}" в число для атрибута ${attribute.name}`);
                        // В случае ошибки устанавливаем значение по умолчанию
                        updatedAttr.numeric_value = 0;
                        updatedAttr.value = 0;
                    }
                } else if (attribute && attribute.attribute_type === 'boolean') {
                    updatedAttr.boolean_value = Boolean(newValue);
                    updatedAttr.value = updatedAttr.boolean_value;
                } else {
                    updatedAttr.text_value = String(newValue || '');
                    updatedAttr.value = updatedAttr.text_value;
                }

                // Обновляем отображаемое значение
                if (attribute && attribute.attribute_type === 'boolean') {
                    updatedAttr.display_value = updatedAttr.boolean_value ? 'Да' : 'Нет';
                } else if (attribute && attribute.attribute_type === 'number') {
                    updatedAttr.display_value = String(updatedAttr.numeric_value);

                    // Добавляем единицы измерения для улучшения читаемости
                    if (attribute.name === 'mileage') {
                        updatedAttr.display_value += ' км';
                    } else if (attribute.name === 'engine_capacity') {
                        updatedAttr.display_value += ' л';
                    } else if (attribute.name === 'power') {
                        updatedAttr.display_value += ' л.с.';
                    }
                } else {
                    updatedAttr.display_value = String(newValue || '');
                }

                console.log(`AttributeFields: Обновлен атрибут ${attribute ? attribute.name : attributeId}:`, updatedAttr);
                return updatedAttr;
            }
            return attr;
        });

        // Сохраняем обновленные значения
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
        if (!attr) {
            console.error(`AttributeFields: Не найден атрибут с ID=${attribute.attribute_id} для отображения`);
            return null;
        }
        let attrValue = attribute.value;
        if (attrValue === undefined || attrValue === null) {
            // Пытаемся извлечь значение из типизированных полей
            if (attribute.attribute_type === 'text' || attribute.attribute_type === 'select') {
                attrValue = attribute.text_value || '';
            } else if (attribute.attribute_type === 'number') {
                attrValue = attribute.numeric_value !== null ? attribute.numeric_value : 0;
            } else if (attribute.attribute_type === 'boolean') {
                attrValue = attribute.boolean_value !== null ? attribute.boolean_value : false;
            } else {
                attrValue = '';
            }
            console.log(`AttributeFields: Исправлено отсутствующее значение для ${attribute.attribute_name}: ${attrValue}`);
        }
        const displayName = getTranslatedName(attr);
        const isRequired = attr.is_required;

        // Получаем текущее значение

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

                // Устанавливаем разумные максимальные значения в зависимости от типа атрибута
                let max;
                if (attr.name === 'engine_capacity') {
                    max = options.max !== undefined ? options.max : 10; // Максимум 10 литров
                } else if (attr.name === 'year') {
                    max = options.max !== undefined ? options.max : new Date().getFullYear() + 1; // Текущий год + 1
                } else if (attr.name === 'mileage') {
                    max = options.max !== undefined ? options.max : 1000000; // Максимум 1 млн километров
                } else if (attr.name === 'power') {
                    max = options.max !== undefined ? options.max : 2000; // Максимум 2000 л.с.
                } else if (attr.name === 'screen_size') {
                    max = options.max !== undefined ? options.max : 30; // Максимум 30 дюймов
                } else if (attr.name === 'camera') {
                    max = options.max !== undefined ? options.max : 200; // Максимум 200 МП
                } else if (attr.name === 'area' || attr.name === 'land_area') {
                    max = options.max !== undefined ? options.max : 10000; // Максимум 10000 м²/соток
                } else if (attr.name === 'floor' || attr.name === 'total_floors') {
                    max = options.max !== undefined ? options.max : 200; // Максимум 200 этажей
                } else {
                    max = options.max !== undefined ? options.max : 1000000; // Значение по умолчанию
                }

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

            case 'select': {
                // Извлекаем options, улучшенная версия
                let options = [];

                try {
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
                    <Box key={`attr-${attribute.attribute_id}`}>
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