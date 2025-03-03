import React, { useState, useEffect, useRef } from 'react';
import { TextField, InputAdornment, IconButton, Box, Paper, List, ListItem, ListItemText, Typography, CircularProgress } from '@mui/material';
import { Search, X } from 'lucide-react';
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


    // Функция для получения подсказок при вводе
    const fetchSuggestions = async (text) => {
        if (!text || text.length < 2) {
            setSuggestions([]);
            setCategorySuggestions([]);
            return;
        }

        try {
            setLoading(true);
            // Запрос на подсказки автодополнения
            const suggestResponse = await axios.get('/api/v1/marketplace/suggestions', {
                params: { q: text, size: 5 }
            });

            // Улучшенная обработка результатов
            if (suggestResponse.data && suggestResponse.data.data) {
                let suggestionData = suggestResponse.data.data;
                // Проверяем, является ли результат массивом
                if (!Array.isArray(suggestionData) && typeof suggestionData === 'object') {
                    // Пытаемся извлечь массив из вложенного объекта
                    suggestionData = suggestionData.data || [];
                }
                // Обеспечиваем уникальность результатов
                const uniqueSuggestions = [...new Set(suggestionData)];
                setSuggestions(uniqueSuggestions);
            } else {
                setSuggestions([]);
            }

            // Остальной код без изменений...
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
        setInputValue(suggestion);
        if (onChange) onChange(suggestion);
        if (onSearch) onSearch(suggestion);
        setShowSuggestions(false);
    };

    // Обработчик клика по категории
    const handleCategoryClick = (category) => {
        // Тут нужно реализовать логику фильтрации по категории
        // Например, передать событие наверх с указанием категории
        if (onSearch) {
            // Предположим, что у нас есть обработчик, который принимает текст и категорию
            onSearch(inputValue, category.id);
        }
        setShowSuggestions(false);
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
            {showSuggestions && (suggestions.length > 0 || categorySuggestions.length > 0) && (
                <Paper
                    elevation={3}
                    sx={{
                        position: 'absolute',
                        width: '100%',
                        zIndex: 1300,
                        mt: 0.5,
                        maxHeight: 350,
                        overflow: 'auto'
                    }}
                >
                    {/* Подсказки текста */}
                    {suggestions.length > 0 && (
                        <>

                            <List dense>
                                {suggestions.map((suggestion, index) => (
                                    <ListItem
                                        key={`suggestion-${index}`}
                                        button
                                        onClick={() => handleSuggestionClick(suggestion)}
                                    >
                                        <ListItemText primary={suggestion} />
                                    </ListItem>
                                ))}
                            </List>
                        </>
                    )}

                    {/* Подсказки категорий */}
                    {categorySuggestions.length > 0 && (
                        <>
                            <Typography variant="subtitle2" sx={{ px: 2, py: 1, bgcolor: 'grey.100' }}>
                                {t('categories')}
                            </Typography>
                            <List dense>
                                {categorySuggestions.map((category) => (
                                    <ListItem
                                        key={`category-${category.id}`}
                                        button
                                        onClick={() => handleCategoryClick(category)}
                                    >
                                        <ListItemText
                                            primary={category.name}
                                            secondary={`${category.listing_count} ${t('listings')}`}
                                        />
                                    </ListItem>
                                ))}
                            </List>
                        </>
                    )}
                </Paper>
            )}

        </Box>
    );
};

export default AutocompleteInput;