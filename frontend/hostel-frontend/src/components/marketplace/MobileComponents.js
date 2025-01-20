import React, { useState, useCallback, useEffect } from 'react';
import { Link } from 'react-router-dom';
import {
    Box, Button, IconButton, Typography, InputBase, Toolbar,
    Paper, Grid, Drawer, Stack
} from '@mui/material';
import { Search as SearchIcon, Filter, X, Check, ArrowLeft, ChevronRight, Plus } from 'lucide-react';
import { debounce } from 'lodash';

// Компонент MobileHeader
export const MobileHeader = ({ onOpenFilters, filtersCount, onSearch, searchValue }) => {
    const [localSearchValue, setLocalSearchValue] = useState(searchValue || '');
    const debouncedSearch = useCallback(
        debounce((value) => onSearch(value), 500),
        [onSearch]
    );

    return (
        <Box sx={{ 
            borderBottom: 1, 
            borderColor: 'divider',
            position: 'sticky',
            top: 0,
            zIndex: 1100,
            bgcolor: 'background.paper'
        }}>
            <Toolbar sx={{ 
                minHeight: '56px !important', 
                px: 2,
                display: 'flex',
                justifyContent: 'space-between'
            }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="h6" component="div" sx={{ fontWeight: 'bold' }}>
                        Sve Tu
                    </Typography>
                </Box>
                <Box sx={{ display: 'flex', gap: 1 }}>
                    <IconButton
                        onClick={onOpenFilters}
                        sx={{
                            position: 'relative',
                            bgcolor: filtersCount > 0 ? 'action.selected' : 'transparent'
                        }}
                    >
                        <Filter size={20} />
                        {filtersCount > 0 && (
                            <Box
                                sx={{
                                    position: 'absolute',
                                    top: 4,
                                    right: 4,
                                    width: 16,
                                    height: 16,
                                    borderRadius: '50%',
                                    bgcolor: 'primary.main',
                                    color: 'primary.contrastText',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    fontSize: '0.75rem'
                                }}
                            >
                                {filtersCount}
                            </Box>
                        )}
                    </IconButton>
                </Box>
            </Toolbar>

            <Box sx={{ px: 2, pb: 2 }}>
                <Paper
                    elevation={0}
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        px: 2,
                        py: 1,
                        bgcolor: 'grey.100',
                        borderRadius: 2
                    }}
                >
                    <InputBase
                        fullWidth
                        placeholder="Поиск объявлений..."
                        value={localSearchValue}
                        onChange={(e) => {
                            setLocalSearchValue(e.target.value);
                            debouncedSearch(e.target.value);
                        }}
                        startAdornment={
                            <SearchIcon style={{ color: 'text.secondary', marginRight: 8 }} size={20} />
                        }
                        endAdornment={
                            localSearchValue && (
                                <IconButton
                                    size="small"
                                    onClick={() => {
                                        setLocalSearchValue('');
                                        onSearch('');
                                    }}
                                >
                                    <X size={16} />
                                </IconButton>
                            )
                        }
                        sx={{
                            '& input': {
                                p: '6px 0',
                                fontSize: '0.875rem'
                            }
                        }}
                    />
                </Paper>
            </Box>
        </Box>
    );
};

// Компонент MobileListingCard
export const MobileListingCard = ({ listing }) => {
    const formatPrice = (price) => {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0
        }).format(price || 0);
    };

    return (
        <Box sx={{ p: 1 }}>
            <Box
                sx={{
                    position: 'relative',
                    paddingTop: '100%',
                    borderRadius: 1,
                    overflow: 'hidden',
                    bgcolor: 'grey.100'
                }}
            >
                {listing.images && listing.images[0] && (
                    <img
                        src={`${process.env.REACT_APP_BACKEND_URL}/uploads/${listing.images[0].file_path}`}
                        alt={listing.title}
                        style={{
                            position: 'absolute',
                            top: 0,
                            left: 0,
                            width: '100%',
                            height: '100%',
                            objectFit: 'cover'
                        }}
                    />
                )}
            </Box>
            <Box sx={{ mt: 1 }}>
                <Typography
                    variant="subtitle2"
                    sx={{
                        fontSize: '0.875rem',
                        fontWeight: 500,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap'
                    }}
                >
                    {listing.title}
                </Typography>
                <Typography
                    variant="subtitle1"
                    sx={{
                        fontSize: '1rem',
                        fontWeight: 600,
                        color: 'primary.main',
                        mt: 0.5
                    }}
                >
                    {formatPrice(listing.price)}
                </Typography>
            </Box>
        </Box>
    );
};


export const MobileFilters = ({ open, onClose, filters, onFilterChange, categories }) => {
    const [tempFilters, setTempFilters] = useState(filters);
    const [currentCategory, setCurrentCategory] = useState(null);
    const [navigationHistory, setNavigationHistory] = useState([]);

    useEffect(() => {
        setTempFilters(filters);
    }, [filters, open]);

    const handleApply = () => {
        onFilterChange(tempFilters);
        onClose();
    };

    // Получаем текущие категории для отображения
    const getCurrentCategories = () => {
        if (!currentCategory) {
            return categories || [];
        }
        return currentCategory.children || [];
    };

    const handleCategoryClick = (category) => {
        const hasChildren = category.children && category.children.length > 0;
        
        if (hasChildren) {
            setNavigationHistory(prev => [...prev, currentCategory]);
            setCurrentCategory(category);
        } else {
            setTempFilters(prev => ({
                ...prev,
                category_id: category.id
            }));
        }
    };

    const handleBack = () => {
        if (navigationHistory.length > 0) {
            const newHistory = [...navigationHistory];
            const lastCategory = newHistory.pop();
            setNavigationHistory(newHistory);
            setCurrentCategory(lastCategory);
        }
    };

    return (
        <Drawer
            anchor="right"
            open={open}
            onClose={onClose}
            PaperProps={{
                sx: { width: '100%', maxWidth: 400 }
            }}
        >
            <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
                {/* Шапка с навигацией */}
                <Box sx={{ 
                    p: 2,
                    borderBottom: '1px solid',
                    borderColor: 'divider',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1
                }}>
                    {navigationHistory.length > 0 && (
                        <IconButton 
                            onClick={handleBack}
                            sx={{ 
                                color: 'text.secondary',
                                '&:hover': { color: 'primary.main' }
                            }}
                        >
                            <ArrowLeft size={20} />
                        </IconButton>
                    )}
                    <Typography 
                        variant="subtitle1" 
                        sx={{ 
                            fontWeight: 600,
                            color: 'text.primary'
                        }}
                    >
                        {currentCategory ? currentCategory.name : 'Категории'}
                    </Typography>
                </Box>

                {/* Список категорий */}
                <Box sx={{ flex: 1, overflow: 'auto' }}>
                    {getCurrentCategories().map((category) => {
                        const hasChildren = category.children && category.children.length > 0;
                        const isSelected = tempFilters.category_id === category.id;
                        
                        return (
                            <Box
                                key={category.id}
                                sx={{
                                    borderBottom: '1px solid',
                                    borderColor: 'divider',
                                    '&:last-child': {
                                        borderBottom: 'none'
                                    }
                                }}
                            >
                                <Button
                                    onClick={() => handleCategoryClick(category)}
                                    sx={{ 
                                        width: '100%',
                                        justifyContent: 'flex-start',
                                        textTransform: 'none',
                                        py: 2,
                                        px: 2,
                                        color: isSelected ? 'primary.main' : 'text.primary',
                                        backgroundColor: isSelected ? 'action.selected' : 'transparent',
                                        borderRadius: 0,
                                        '&:hover': {
                                            backgroundColor: isSelected ? 'action.selected' : 'action.hover'
                                        }
                                    }}
                                >
                                    <Typography 
                                        sx={{ 
                                            flex: 1,
                                            textAlign: 'left',
                                            fontWeight: isSelected ? 500 : 400
                                        }}
                                    >
                                        {category.name}
                                    </Typography>
                                    {hasChildren && (
                                        <ChevronRight 
                                            size={20} 
                                            style={{ 
                                                opacity: 0.5,
                                                marginLeft: 8
                                            }}
                                        />
                                    )}
                                </Button>
                            </Box>
                        );
                    })}
                </Box>

                {/* Кнопка применения фильтров */}
                <Box sx={{ p: 2, borderTop: '1px solid', borderColor: 'divider' }}>
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={handleApply}
                        startIcon={<Check size={20} />}
                        sx={{ 
                            py: 1.5,
                            textTransform: 'none',
                            fontWeight: 500,
                            boxShadow: 'none',
                            '&:hover': {
                                boxShadow: 'none'
                            }
                        }}
                    >
                        Применить фильтры
                    </Button>
                </Box>
            </Box>
        </Drawer>
    );
};