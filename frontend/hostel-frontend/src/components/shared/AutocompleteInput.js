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
    
            console.log('Received suggestions response:', suggestResponse.data);
    
            // Отладочный вывод структуры данных
            console.log('Response data structure:', JSON.stringify(suggestResponse.data, null, 2));
    
            let suggestionsList = [];
    
            // Проверяем структуру ответа и извлекаем подсказки
            if (suggestResponse.data) {
                // Попытка 1: Извлечение из data.suggestions (предполагаемая структура API)
                if (suggestResponse.data.data && Array.isArray(suggestResponse.data.data.suggestions)) {
                    suggestionsList = suggestResponse.data.data.suggestions;
                } 
                // Попытка 2: Извлечение напрямую из data, если это массив
                else if (suggestResponse.data.data && Array.isArray(suggestResponse.data.data)) {
                    suggestionsList = suggestResponse.data.data;
                }
                // Попытка 3: Если data.data содержит объект с hits/options
                else if (suggestResponse.data.data && typeof suggestResponse.data.data === 'object') {
                    // Поиск в различных стандартных местах ответа OpenSearch
                    if (suggestResponse.data.data.hits && 
                        suggestResponse.data.data.hits.hits && 
                        Array.isArray(suggestResponse.data.data.hits.hits)) {
                        // Получение из hits.hits._source.title (обычный поисковый результат)
                        suggestionsList = suggestResponse.data.data.hits.hits
                            .filter(hit => hit._source && hit._source.title)
                            .map(hit => hit._source.title);
                    } 
                    // Поиск в структуре suggest API
                    else if (suggestResponse.data.data.suggest) {
                        // Обработка различных полей suggest API
                        Object.keys(suggestResponse.data.data.suggest).forEach(key => {
                            const suggestField = suggestResponse.data.data.suggest[key];
                            if (Array.isArray(suggestField) && suggestField.length > 0) {
                                // Извлекаем options из каждого элемента suggest
                                suggestField.forEach(item => {
                                    if (item.options && Array.isArray(item.options)) {
                                        const texts = item.options
                                            .filter(opt => opt.text)
                                            .map(opt => opt.text);
                                        suggestionsList = [...suggestionsList, ...texts];
                                    }
                                });
                            }
                        });
                    }
                }
                
                // Последняя попытка: прямой поиск в корне ответа
                else if (typeof suggestResponse.data === 'object') {
                    // Поиск полей, которые могут содержать массивы строк
                    Object.keys(suggestResponse.data).forEach(key => {
                        const value = suggestResponse.data[key];
                        if (Array.isArray(value) && value.length > 0 && typeof value[0] === 'string') {
                            suggestionsList = [...suggestionsList, ...value];
                        }
                    });
                }
            }
    
            // Обеспечиваем уникальность и удаляем пустые значения
            const uniqueSuggestions = [...new Set(suggestionsList.filter(item => item && typeof item === 'string'))]
                .map(text => text.trim());
                
            setSuggestions(uniqueSuggestions);
            console.log('Processed suggestions:', uniqueSuggestions);
            
            // Если мы всё ещё не нашли подсказок, но поисковый запрос возвращает результаты,
            // используем заголовки объявлений как подсказки
            if (uniqueSuggestions.length === 0) {
                try {
                    const searchResponse = await axios.get('/api/v1/marketplace/search', {
                        params: { q: text, size: 3 }
                    });
                    
                    console.log('Fallback to search API for suggestions');
                    
                    if (searchResponse.data && searchResponse.data.data) {
                        let listings = [];
                        
                        if (Array.isArray(searchResponse.data.data)) {
                            listings = searchResponse.data.data;
                        } else if (searchResponse.data.data.data && Array.isArray(searchResponse.data.data.data)) {
                            listings = searchResponse.data.data.data;
                        }
                        
                        if (listings.length > 0) {
                            const titles = listings
                                .filter(listing => listing.title && listing.title.toLowerCase().includes(text.toLowerCase()))
                                .map(listing => listing.title);
                                
                            setSuggestions([...new Set(titles)]);
                            console.log('Fallback suggestions from search results:', titles);
                        }
                    }
                } catch (error) {
                    console.error('Ошибка запроса к fallback API:', error);
                }
            }
            
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