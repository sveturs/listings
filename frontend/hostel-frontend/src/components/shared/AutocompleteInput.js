import React, { useState, useEffect, useRef } from 'react';
import { TextField, InputAdornment, IconButton, Box, Paper, List, ListItem, ListItemText, Typography, CircularProgress, ListItemIcon } from '@mui/material';
import { Search, X, ShoppingCart, Folder } from 'lucide-react';
import axios from '../../api/axios';
import { useTranslation } from 'react-i18next';

const AutocompleteInput = ({ value, onChange, onSearch, placeholder, debounceTime = 300 }) => {
    const { t } = useTranslation('common');
    const [inputValue, setInputValue] = useState(value || '');
    const [suggestions, setSuggestions] = useState([]);
    const [loading, setLoading] = useState(false);
    const [showSuggestions, setShowSuggestions] = useState(false);
    const [categorySuggestions, setCategorySuggestions] = useState([]);
    const inputRef = useRef(null);
    const debounceRef = useRef(null);

    // Синхронизация с внешним значением
    useEffect(() => {
        setInputValue(value || '');
    }, [value]);

    // Функция для форматирования отображаемого текста категории
    const formatCategoryDisplay = (category, priority) => {
        if (!category || !category.name) return "";
        
        // В зависимости от приоритета используем разные префиксы
        if (priority === 2) {
            return `В категории: ${category.name}`; // Для непосредственных категорий
        } else if (priority === 3) {
            return `Раздел: ${category.name}`; // Для родительских категорий
        } else {
            return `Категория: ${category.name}`; // Для других категорий
        }
    };

    // Функция для получения подсказок при вводе
    const fetchSuggestions = async (text) => {
        if (!text || text.length < 2) {
            setSuggestions([]);
            setCategorySuggestions([]);
            return;
        }

        try {
            setLoading(true);

            // 1. Запрос на товары (через обычный поиск)
            const searchResponse = await axios.get('/api/v1/marketplace/search', {
                params: { q: text, size: 3 } // Ограничиваем количество товаров
            });

            console.log('Using search API for product suggestions, query:', text);

            let productSuggestions = [];

            // Извлекаем товары из результатов поиска
            if (searchResponse.data && searchResponse.data.data) {
                let listings = [];

                if (Array.isArray(searchResponse.data.data)) {
                    listings = searchResponse.data.data;
                } else if (searchResponse.data.data.data && Array.isArray(searchResponse.data.data.data)) {
                    listings = searchResponse.data.data.data;
                }

                // Фильтруем товары по релевантности
                const lowerQuery = text.toLowerCase();
                productSuggestions = listings
                    .filter(listing => listing && listing.title &&
                        listing.title.toLowerCase().includes(lowerQuery))
                    .map(listing => ({
                        id: listing.id,
                        type: 'product',
                        title: listing.title,
                        category_id: listing.category_id,
                        category_path_ids: listing.category_path_ids || []
                    }));
            }

            // 2. Запрос на категории 
            // Получаем все категории (или используем уже загруженные)
            let allCategories = [];
            try {
                const categoriesResponse = await axios.get('/api/v1/marketplace/category-tree');
                if (categoriesResponse.data?.data) {
                    allCategories = categoriesResponse.data.data;
                }
            } catch (err) {
                console.error('Ошибка при получении категорий:', err);
            }

            // Функция для плоского представления категорий (включая все подкатегории)
            const flattenCategories = (categories, parentPath = []) => {
                let result = [];

                for (const category of categories) {
                    const currentPath = [...parentPath, category];
                    result.push({
                        id: category.id,
                        name: category.name,
                        type: 'category',
                        depth: parentPath.length,
                        path: currentPath,
                        parent_id: category.parent_id
                    });

                    if (category.children && Array.isArray(category.children) && category.children.length > 0) {
                        result = [...result, ...flattenCategories(category.children, currentPath)];
                    }
                }

                return result;
            };

            // Получаем плоский список всех категорий
            const flatCategories = flattenCategories(allCategories);

            // Ищем категории, соответствующие запросу
            const lowerQuery = text.toLowerCase();
            const matchingCategories = flatCategories
                .filter(cat => cat.name.toLowerCase().includes(lowerQuery))
                .sort((a, b) => {
                    // Сначала сортируем по точности совпадения
                    const aStartsWith = a.name.toLowerCase().startsWith(lowerQuery) ? 0 : 1;
                    const bStartsWith = b.name.toLowerCase().startsWith(lowerQuery) ? 0 : 1;
                    if (aStartsWith !== bStartsWith) return aStartsWith - bStartsWith;

                    // Затем по глубине (более специфичные категории сначала)
                    return b.depth - a.depth;
                })
                .slice(0, 3); // Ограничиваем количество категорий

            // 3. Формируем массив подсказок в нужном порядке

            // Сначала добавляем конкретные товары
            let finalSuggestions = productSuggestions.map(product => ({
                ...product,
                display: product.title,
                priority: 1
            }));

            // Затем добавляем категории товаров из результатов поиска
            const productCategoryIds = new Set();

            for (const product of productSuggestions) {
                if (product.category_id) {
                    productCategoryIds.add(product.category_id);

                    // Находим категорию товара
                    const category = flatCategories.find(cat => cat.id === product.category_id);
                    if (category && !finalSuggestions.some(s => s.type === 'category' && s.id === category.id)) {
                        finalSuggestions.push({
                            id: category.id,
                            type: 'category',
                            title: category.name,
                            display: formatCategoryDisplay(category, 2),
                            priority: 2,
                            path: category.path
                        });
                    }

                    // Если есть родительская категория, добавляем её
                    if (category && category.parent_id) {
                        const parentCategory = flatCategories.find(cat => cat.id === category.parent_id);
                        if (parentCategory && !finalSuggestions.some(s => s.type === 'category' && s.id === parentCategory.id)) {
                            finalSuggestions.push({
                                id: parentCategory.id,
                                type: 'category',
                                title: parentCategory.name,
                                display: formatCategoryDisplay(parentCategory, 3),
                                priority: 3,
                                path: parentCategory.path
                            });
                        }
                    }
                }
            }

            // Добавляем категории из прямого поиска по категориям
            for (const category of matchingCategories) {
                if (!finalSuggestions.some(s => s.type === 'category' && s.id === category.id)) {
                    finalSuggestions.push({
                        id: category.id,
                        type: 'category',
                        title: category.name,
                        display: formatCategoryDisplay(category, category.depth === 0 ? 4 : 3),
                        priority: category.depth === 0 ? 4 : 3,
                        path: category.path
                    });
                }
            }

            // Сортируем по приоритету
            finalSuggestions.sort((a, b) => a.priority - b.priority);

            // Ограничиваем общее количество подсказок
            finalSuggestions = finalSuggestions.slice(0, 8);

            setSuggestions(finalSuggestions);
            console.log('Generated enhanced suggestions:', finalSuggestions);

        } catch (error) {
            console.error('Ошибка при получении подсказок:', error);
            setSuggestions([]);
            setCategorySuggestions([]);
        } finally {
            setLoading(false);
        }
    };

    // Обработка изменения ввода с дебаунсом
    const handleInputChange = (e) => {
        const newValue = e.target.value;
        setInputValue(newValue);

        // Вызываем внешний обработчик (для обновления state родителя)
        if (onChange) onChange(newValue);

        // Если строка пустая, очищаем подсказки
        if (!newValue) {
            setSuggestions([]);
            setCategorySuggestions([]);
            return;
        }

        // Включаем отображение подсказок
        setShowSuggestions(true);

        // Дебаунс для сокращения количества запросов
        if (debounceRef.current) {
            clearTimeout(debounceRef.current);
        }

        debounceRef.current = setTimeout(() => {
            fetchSuggestions(newValue);

            // Автоматически выполняем поиск по мере ввода
            if (onSearch && newValue.length >= 2) {
                onSearch(newValue);
            }
        }, debounceTime);
    };

    // Обработка отправки формы
    const handleSubmit = (e) => {
        e.preventDefault();
        if (onSearch && inputValue.trim()) {
            onSearch(inputValue.trim());
        }
        setShowSuggestions(false);
    };

    // Обработчик клика по подсказке
    const handleSuggestionClick = (suggestion) => {
        if (suggestion.type === 'product') {
            // Для товара - переходим на страницу товара
            setInputValue(suggestion.title);
            if (onChange) onChange(suggestion.title);
            if (onSearch) onSearch(suggestion.title);
            setShowSuggestions(false);

            // Если доступен ID товара, перенаправляем на его страницу
            if (suggestion.id) {
                window.location.href = `/marketplace/listings/${suggestion.id}`;
            }
        } else if (suggestion.type === 'category') {
            // Для категории - фильтруем по категории

            // Очищаем текст поиска при выборе категории, чтобы показать все товары категории
            setInputValue("");
            if (onChange) onChange("");

            // Вызываем поиск с указанием только категории
            // Передаем пустую строку как первый параметр, чтобы очистить поисковый запрос
            if (onSearch) onSearch("", suggestion.id);

            setShowSuggestions(false);

            // Можно также добавить сообщение для пользователя
            console.log(`Показаны все товары в категории: ${suggestion.title}`);
        }
    };

    // Очистка ввода
    const handleClear = () => {
        setInputValue('');
        setSuggestions([]);
        setCategorySuggestions([]);
        if (onChange) onChange('');
        if (onSearch) onSearch('');
        if (inputRef.current) inputRef.current.focus();
    };

    // Скрытие подсказок при клике вне компонента
    useEffect(() => {
        const handleClickOutside = (e) => {
            if (inputRef.current && !inputRef.current.contains(e.target)) {
                setShowSuggestions(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    return (
        <Box ref={inputRef} sx={{ position: 'relative', width: '100%' }}>
            <form onSubmit={handleSubmit}>
                <TextField
                    value={inputValue}
                    onChange={handleInputChange}
                    placeholder={placeholder || t('search')}
                    variant="outlined"
                    fullWidth
                    size="small"
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <Search size={20} />
                            </InputAdornment>
                        ),
                        endAdornment: inputValue ? (
                            <InputAdornment position="end">
                                {loading ? (
                                    <CircularProgress size={20} />
                                ) : (
                                    <IconButton size="small" onClick={handleClear}>
                                        <X size={16} />
                                    </IconButton>
                                )}
                            </InputAdornment>
                        ) : null
                    }}
                />
            </form>

            {/* Выпадающие подсказки */}
            {showSuggestions && suggestions.length > 0 && (
    <Paper
        elevation={4}
        sx={{
            position: 'absolute',
            width: '100%',
            zIndex: 1300,
            mt: 0.5,
            maxHeight: 400,
            overflow: 'auto',
            borderRadius: '8px',
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.08)'
        }}
    >
        <List dense sx={{ p: 0 }}>
            {suggestions.map((suggestion, index) => {
                // Определяем тип и приоритет для стилизации
                const isProduct = suggestion.type === 'product';
                const priority = suggestion.priority || 1;

                // Выбираем цвет полосы слева
                let borderColor = '#4682B4'; // Синий для товаров
                if (!isProduct) {
                    borderColor = '#6BAB33'; // Зеленый для категорий
                }

                return (
                    <ListItem
                        key={`suggestion-${index}`}
                        button
                        onClick={() => handleSuggestionClick(suggestion)}
                        sx={{
                            padding: '10px 16px',
                            borderBottom: index < suggestions.length - 1 ? '1px solid #f0f0f0' : 'none',
                            position: 'relative',
                            '&::before': {
                                content: '""',
                                position: 'absolute',
                                left: 0,
                                top: '20%',
                                height: '60%',
                                width: '4px',
                                backgroundColor: borderColor,
                                borderRadius: '0 4px 4px 0'
                            },
                            '&:hover': {
                                backgroundColor: 'rgba(0, 0, 0, 0.03)'
                            },
                            transition: 'background-color 0.2s ease'
                        }}
                    >
                        <ListItemText
                            primary={
                                <Typography
                                    variant="body2"
                                    sx={{
                                        fontWeight: isProduct ? 500 : 400,
                                        color: 'text.primary',
                                        fontSize: '0.95rem',
                                        lineHeight: '1.25',
                                        overflow: 'hidden',
                                        textOverflow: 'ellipsis',
                                        whiteSpace: 'nowrap'
                                    }}
                                >
                                    {/* Показываем только название без префикса "Категория/Раздел" */}
                                    {isProduct 
                                        ? suggestion.title 
                                        : suggestion.title}
                                </Typography>
                            }
                            secondary={
                                suggestion.type === 'category' && suggestion.path && suggestion.path.length > 1 ? (
                                    <Typography
                                        variant="caption"
                                        sx={{
                                            color: 'text.secondary',
                                            fontSize: '0.8rem',
                                            display: 'block',
                                            maxWidth: '100%',
                                            overflow: 'hidden',
                                            textOverflow: 'ellipsis',
                                            whiteSpace: 'nowrap'
                                        }}
                                    >
                                        {suggestion.path.slice(0, -1).map(cat => cat.name).join(' > ')}
                                    </Typography>
                                ) : null
                            }
                        />
                    </ListItem>
                );
            })}
        </List>
    </Paper>
)}

        </Box>
    );
};

export default AutocompleteInput;