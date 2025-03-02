import React, { useState, useEffect, useRef } from 'react';
import { 
    TextField, 
    InputAdornment, 
    IconButton, 
    Box, 
    Paper, 
    List, 
    ListItem, 
    ListItemText,
    CircularProgress
} from '@mui/material';
import { Search, X } from 'lucide-react';
import axios from '../../api/axios';
import { debounce } from 'lodash';

const AutocompleteInput = ({ value, onChange, placeholder, onSearch, size = 'small' }) => {
    const [inputValue, setInputValue] = useState(value || '');
    const [suggestions, setSuggestions] = useState([]);
    const [loading, setLoading] = useState(false);
    const [open, setOpen] = useState(false);
    const inputRef = useRef(null);
    
    // Обновляем inputValue при изменении value
    useEffect(() => {
        setInputValue(value || '');
    }, [value]);
    
    // Создаем функцию с дебаунсом для запроса предложений
    const fetchSuggestions = useRef(
        debounce(async (text) => {
            if (!text || text.length < 2) {
                setSuggestions([]);
                setLoading(false);
                return;
            }
            
            try {
                setLoading(true);
                const response = await axios.get(`/api/v1/marketplace/suggestions?q=${encodeURIComponent(text)}`);
                if (response.data && response.data.data) {
                    setSuggestions(response.data.data);
                }
            } catch (error) {
                console.error('Error fetching suggestions:', error);
            } finally {
                setLoading(false);
            }
        }, 300)
    ).current;
    
    const handleInputChange = (e) => {
        const newValue = e.target.value;
        setInputValue(newValue);
        
        if (newValue.length >= 2) {
            fetchSuggestions(newValue);
            setOpen(true);
        } else {
            setSuggestions([]);
            setOpen(false);
        }
    };
    
    const handleSelectSuggestion = (suggestion) => {
        setInputValue(suggestion);
        onChange(suggestion);
        setSuggestions([]);
        setOpen(false);
        
        if (onSearch) {
            onSearch(suggestion);
        }
    };
    
    const handleClear = () => {
        setInputValue('');
        onChange('');
        setSuggestions([]);
        setOpen(false);
    };
    
    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            if (onSearch) {
                onSearch(inputValue);
            }
            setOpen(false);
        } else if (e.key === 'Escape') {
            setOpen(false);
        }
    };
    
    // Закрываем выпадающий список при клике вне компонента
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (inputRef.current && !inputRef.current.contains(event.target)) {
                setOpen(false);
            }
        };
        
        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);
    
    return (
        <Box ref={inputRef} sx={{ position: 'relative', width: '100%' }}>
            <TextField
                fullWidth
                size={size}
                value={inputValue}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown}
                placeholder={placeholder}
                InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                            <Search size={16} color="gray" />
                        </InputAdornment>
                    ),
                    endAdornment: (
                        <InputAdornment position="end">
                            {loading ? (
                                <CircularProgress size={16} />
                            ) : inputValue ? (
                                <IconButton
                                    edge="end"
                                    size="small"
                                    onClick={handleClear}
                                >
                                    <X size={16} />
                                </IconButton>
                            ) : null}
                        </InputAdornment>
                    )
                }}
            />
            
            {open && suggestions.length > 0 && (
                <Paper
                    sx={{
                        position: 'absolute',
                        width: '100%',
                        maxHeight: 300,
                        overflow: 'auto',
                        mt: 0.5,
                        zIndex: 1300,
                        boxShadow: 3
                    }}
                >
                    <List disablePadding>
                        {suggestions.map((suggestion, index) => (
                            <ListItem
                                key={index}
                                button
                                dense
                                onClick={() => handleSelectSuggestion(suggestion)}
                            >
                                <ListItemText
                                    primary={suggestion}
                                    primaryTypographyProps={{
                                        noWrap: true,
                                        variant: 'body2'
                                    }}
                                />
                            </ListItem>
                        ))}
                    </List>
                </Paper>
            )}
        </Box>
    );
};

export default AutocompleteInput;