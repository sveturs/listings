// frontend/hostel-frontend/src/components/marketplace/MobileComponents.js
import React, { useState, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

import { Link } from 'react-router-dom';
import {
    Box, Button, IconButton, Typography, InputBase, Toolbar, TextField, Select, MenuItem,
    Paper, Grid, Drawer, Stack
} from '@mui/material';
import { Search as SearchIcon, Filter, X, Check, ArrowLeft, ChevronRight, Plus, Store  } from 'lucide-react';
import { debounce } from 'lodash';
 
// Компонент MobileHeader
export const MobileHeader = ({ onOpenFilters, filtersCount, onSearch, searchValue }) => {
    const { t } = useTranslation('marketplace', 'common');

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
                justifyContent: 'space-between',
                gap: 2
            }}>
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

                <Button
                    component={Link}
                    to="/marketplace/create"
                    variant="contained"
                    size="small"
                    startIcon={<Plus size={16} />}
                    sx={{
                        textTransform: 'none',
                        fontWeight: 500,
                        height: 32
                    }}
                >
                    {t('buttons.create', { ns: 'common' })}
                </Button>
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
                        placeholder={t('buttons.search', { ns: 'common' })}
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
    const { i18n } = useTranslation();

    const getTranslatedText = (field) => {
        if (!listing) return '';

        if (i18n.language === listing.original_language) {
            return listing[field];
        }

        return listing.translations?.[i18n.language]?.[field] || listing[field];
    };
    const formatPrice = (price) => {
        return new Intl.NumberFormat('sr-RS', {
            style: 'currency',
            currency: 'RSD',
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
                  {/*  бейдж магазина */}
                  {listing.storefront_id && (
                    <Box 
                        sx={{
                            position: 'absolute',
                            top: 5,
                            right: 5,
                            zIndex: 1,
                            bgcolor: 'primary.main',
                            color: 'white',
                            borderRadius: '4px',
                            px: 0.5,
                            py: 0.25,
                            display: 'flex',
                            alignItems: 'center',
                            gap: 0.5,
                            fontSize: '0.7rem',
                            fontWeight: 'bold'
                        }}
                    >
                        <Store size={12} />
                        Магазин
                    </Box>
                )}
                
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
                    {getTranslatedText('title')}
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
    const { t, i18n } = useTranslation('marketplace'); // Добавляем i18n

    const getTranslatedName = (category) => {
        if (!category) return '';

        // Если у категории есть переводы для текущего языка
        if (category.translations && category.translations[i18n.language]) {
            return category.translations[i18n.language].name || category.name;
        }

        // Если переводов нет или они не подходят, возвращаем исходное имя
        return category.name;
    };

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
        // Если нет текущей категории, показываем только корневые категории
        if (!currentCategory) {
            return categories.filter(category => !category.parent_id) || [];
        }

        // Иначе показываем только дочерние категории текущей категории
        return categories.filter(category =>
            category.parent_id && String(category.parent_id) === String(currentCategory.id)
        ) || [];
    };

    const handleCategoryClick = (category) => {
        const hasChildren = category.children && category.children.length > 0;
    
        // Устанавливаем выбранную категорию в фильтрах, даже если у неё есть дочерние элементы
        setTempFilters(prev => ({
            ...prev,
            category_id: category.id
        }));
    
        // Если есть дочерние категории, переходим в них
        if (hasChildren) {
            setNavigationHistory(prev => [...prev, currentCategory]);
            setCurrentCategory(category);
        }
        // Не закрываем фильтры автоматически
    };
    

    const handleBack = () => {
        if (navigationHistory.length > 0) {
            const newHistory = [...navigationHistory];
            const lastCategory = newHistory.pop();
            setNavigationHistory(newHistory);
            setCurrentCategory(lastCategory);
        }
    };

    const handleClearFilters = () => {
        setTempFilters({
            query: "",
            category_id: "",
            min_price: "",
            max_price: "",
            condition: "",
            sort_by: "date_desc"
        });
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
                            flex: 1,
                            fontWeight: 600,
                            color: 'text.primary'
                        }}
                    >
                        {currentCategory ? getTranslatedName(currentCategory) : t('listings.filters.title')}
                    </Typography>
                    <Button
                        variant="text"
                        size="small"
                        onClick={handleClearFilters}
                        sx={{ color: 'text.secondary' }}
                    >
                        {t('listings.filters.reset')}
                    </Button>
                </Box>

                {/* Основное содержимое с фильтрами */}
                <Box sx={{ flex: 1, overflow: 'auto' }}>
                    {/* Фильтр цены */}
                    <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
                        <Typography variant="subtitle2" gutterBottom>
                            {t('listings.filters.price.label')}
                        </Typography>
                        <Stack direction="row" spacing={1}>
                            <TextField
                                fullWidth
                                size="small"
                                placeholder={t('listings.filters.price.min')}
                                type="number"
                                value={tempFilters.min_price || ''}
                                onChange={(e) => setTempFilters(prev => ({
                                    ...prev,
                                    min_price: e.target.value
                                }))}
                            />
                            <TextField
                                fullWidth
                                size="small"
                                placeholder={t('listings.filters.price.max')}
                                type="number"
                                value={tempFilters.max_price || ''}
                                onChange={(e) => setTempFilters(prev => ({
                                    ...prev,
                                    max_price: e.target.value
                                }))}
                            />
                        </Stack>
                    </Box>

                    {/* Фильтр состояния */}
                    <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
                        <Typography variant="subtitle2" gutterBottom>
                            {t('listings.details.condition')}
                        </Typography>
                        <Select
                            fullWidth
                            size="small"
                            value={tempFilters.condition || ''}
                            onChange={(e) => setTempFilters(prev => ({
                                ...prev,
                                condition: e.target.value
                            }))}
                        >
                            <MenuItem value="">{t('listings.create.condition.any')}</MenuItem>
                            <MenuItem value="new">{t('listings.create.condition.new')}</MenuItem>
                            <MenuItem value="used">{t('listings.create.condition.used')}</MenuItem>
                        </Select>
                    </Box>

                    {/* Сортировка */}
                    <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
                        <Typography variant="subtitle2" gutterBottom>
                            {t('listings.filters.sort.label')}
                        </Typography>
                        <Select
                            fullWidth
                            size="small"
                            value={tempFilters.sort_by || 'date_desc'}
                            onChange={(e) => setTempFilters(prev => ({
                                ...prev,
                                sort_by: e.target.value
                            }))}
                        >
                            <MenuItem value="date_desc">{t('listings.filters.sort.newest')}</MenuItem>
                            <MenuItem value="price_asc">{t('listings.filters.sort.cheapest')}</MenuItem>
                            <MenuItem value="price_desc">{t('listings.filters.sort.expensive')}</MenuItem>
                        </Select>
                    </Box>

                    {/* Категории */}
                    <Box sx={{ borderBottom: '1px solid', borderColor: 'divider' }}>
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
                                                fontWeight: isSelected ? 600 : 400,
                                                color: isSelected ? 'primary.main' : 'text.primary'
                                            }}
                                        >
                                            {getTranslatedName(category)}
                                            {isSelected && (
                                                <Box component="span" sx={{ ml: 1, color: 'primary.main' }}>
                                                    ✓
                                                </Box>
                                            )}
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
                </Box>

                {/* Кнопки действий */}
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
                        }}
                    >
                        {t('listings.filters.apply')}
                    </Button>
                </Box>
            </Box>
        </Drawer>
    );

};