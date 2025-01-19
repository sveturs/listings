//frontend/hostel-frontend/src/components/marketplace/mobile/MobileComponents.js
import React, { useState } from 'react';
import {
    Box,
    Drawer,
    AppBar,
    Toolbar,
    IconButton,
    Typography,
    Badge,
    InputBase,
    Button,
    Chip,
    Stack,
    Divider,
    alpha,
    Paper
} from '@mui/material';
import {
    Search,
    Sliders, 
    ArrowLeft,
    X,
    Check,
} from 'lucide-react';

// Компактная карточка для мобильной версии
const MobileListingCard = ({ listing }) => {
    return (
        <Box
            component={Paper}
            variant="outlined"
            sx={{
                display: 'flex',
                flexDirection: 'column',
                height: '100%',
                overflow: 'hidden',
                borderRadius: 1,
            }}
        >
            {/* Изображение */}
            <Box
                sx={{
                    position: 'relative',
                    paddingTop: '100%',
                    overflow: 'hidden',
                    backgroundColor: 'grey.100'
                }}
            >
                <img
                    src={`${process.env.REACT_APP_BACKEND_URL}/uploads/${listing.images?.[0]?.file_path}`}
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
                {listing.images?.length > 1 && (
                    <Chip
                        label={`+${listing.images.length - 1}`}
                        size="small"
                        sx={{
                            position: 'absolute',
                            bottom: 8,
                            right: 8,
                            bgcolor: 'rgba(0,0,0,0.6)',
                            color: 'white',
                            height: 20,
                            '& .MuiChip-label': { px: 1 }
                        }}
                    />
                )}
            </Box>

            {/* Контент */}
            <Box sx={{ p: 1, flexGrow: 1 }}>
                <Typography
                    variant="h6"
                    sx={{
                        fontSize: '0.875rem',
                        fontWeight: 500,
                        mb: 0.5,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        display: '-webkit-box',
                        WebkitLineClamp: 2,
                        WebkitBoxOrient: 'vertical'
                    }}
                >
                    {listing.title}
                </Typography>

                <Typography
                    variant="subtitle1"
                    sx={{
                        fontSize: '1rem',
                        fontWeight: 600,
                        color: 'primary.main'
                    }}
                >
                    {new Intl.NumberFormat('ru-RU', {
                        style: 'currency',
                        currency: 'RUB',
                        maximumFractionDigits: 0
                    }).format(listing.price)}
                </Typography>

                <Typography
                    variant="caption"
                    color="text.secondary"
                    sx={{
                        display: 'block',
                        mt: 0.5,
                        fontSize: '0.75rem'
                    }}
                >
                    {listing.location}
                </Typography>
            </Box>
        </Box>
    );
};

// Мобильная шапка с поиском
const MobileHeader = ({ onOpenFilters, filtersCount }) => (
    <AppBar 
        position="sticky" 
        color="inherit" 
        elevation={0}
        sx={{ 
            borderBottom: 1,
            borderColor: 'divider'
        }}
    >
        <Toolbar sx={{ gap: 1, minHeight: 56 }}>
            <Box
                sx={{
                    flex: 1,
                    display: 'flex',
                    bgcolor: (theme) => alpha(theme.palette.common.black, 0.05),
                    borderRadius: 1,
                    px: 1,
                }}
            >
                <InputBase
                    placeholder="Поиск объявлений"
                    startAdornment={<Search size={18} style={{ marginRight: 8 }} />}
                    sx={{ flex: 1, fontSize: '0.875rem' }}
                />
            </Box>
            <IconButton 
    onClick={onOpenFilters}
    sx={{ position: 'relative' }}
>
    <Sliders size={20} />
    {filtersCount > 0 && (
        <Box
            sx={{
                position: 'absolute',
                top: 8,
                right: 8,
                width: 8,
                height: 8,
                bgcolor: 'primary.main',
                borderRadius: '50%'
            }}
        />
    )}
</IconButton>
        </Toolbar>
    </AppBar>
);

// Мобильные фильтры
const MobileFilters = ({ open, onClose, filters, onFilterChange, categories }) => {
    const [tempFilters, setTempFilters] = useState(filters);
    
    const handleApply = () => {
        onFilterChange(tempFilters);
        onClose();
    };

    const handleReset = () => {
        setTempFilters({
            query: '',
            category_id: '',
            min_price: '',
            max_price: '',
            condition: '',
            sort_by: 'date_desc'
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
            {/* Шапка */}
            <Box sx={{ 
                display: 'flex', 
                alignItems: 'center', 
                px: 2,
                py: 1,
                borderBottom: 1,
                borderColor: 'divider'
            }}>
                <IconButton onClick={onClose} edge="start">
                    <ArrowLeft />
                </IconButton>
                <Typography 
                    variant="subtitle1"
                    sx={{ 
                        flex: 1,
                        ml: 2,
                        fontWeight: 500
                    }}
                >
                    Фильтры
                </Typography>
                <Button 
                    variant="text" 
                    onClick={handleReset}
                    size="small"
                >
                    Сбросить
                </Button>
            </Box>

            {/* Контент */}
            <Box sx={{ 
                flex: 1,
                overflowY: 'auto',
                px: 2,
                py: 2
            }}>
                //Категории мобильной версии
                <Typography variant="subtitle2" gutterBottom>
                    Категории 
                </Typography>
                <Stack spacing={1} sx={{ mb: 3 }}>
                    {categories.map(category => (
                        <Button
                            key={category.id}
                            variant={tempFilters.category_id === category.id ? "contained" : "outlined"}
                            size="small"
                            onClick={() => setTempFilters(prev => ({
                                ...prev,
                                category_id: category.id
                            }))}
                            sx={{ 
                                justifyContent: 'flex-start',
                                px: 1.5,
                                py: 0.75
                            }}
                        >
                            {category.name}
                        </Button>
                    ))}
                </Stack>

                <Typography variant="subtitle2" gutterBottom>
                    Цена
                </Typography>
                <Stack direction="row" spacing={1} sx={{ mb: 3 }}>
                    <InputBase
                        placeholder="От"
                        type="number"
                        value={tempFilters.min_price}
                        onChange={(e) => setTempFilters(prev => ({
                            ...prev,
                            min_price: e.target.value
                        }))}
                        sx={{
                            flex: 1,
                            border: 1,
                            borderColor: 'divider',
                            borderRadius: 1,
                            px: 1,
                            py: 0.5
                        }}
                    />
                    <InputBase
                        placeholder="До"
                        type="number"
                        value={tempFilters.max_price}
                        onChange={(e) => setTempFilters(prev => ({
                            ...prev,
                            max_price: e.target.value
                        }))}
                        sx={{
                            flex: 1,
                            border: 1,
                            borderColor: 'divider',
                            borderRadius: 1,
                            px: 1,
                            py: 0.5
                        }}
                    />
                </Stack>

                <Typography variant="subtitle2" gutterBottom>
                    Состояние
                </Typography>
                <Stack direction="row" spacing={1} sx={{ mb: 3 }}>
                    {['new', 'used'].map((condition) => (
                        <Button
                            key={condition}
                            variant={tempFilters.condition === condition ? "contained" : "outlined"}
                            size="small"
                            onClick={() => setTempFilters(prev => ({
                                ...prev,
                                condition: condition
                            }))}
                            sx={{ flex: 1 }}
                        >
                            {condition === 'new' ? 'Новое' : 'Б/у'}
                        </Button>
                    ))}
                </Stack>

                <Typography variant="subtitle2" gutterBottom>
                    Сортировка
                </Typography>
                <Stack spacing={1}>
                    {[
                        { value: 'date_desc', label: 'Сначала новые' },
                        { value: 'price_asc', label: 'Сначала дешевле' },
                        { value: 'price_desc', label: 'Сначала дороже' }
                    ].map((option) => (
                        <Button
                            key={option.value}
                            variant={tempFilters.sort_by === option.value ? "contained" : "outlined"}
                            size="small"
                            onClick={() => setTempFilters(prev => ({
                                ...prev,
                                sort_by: option.value
                            }))}
                            sx={{ 
                                justifyContent: 'flex-start',
                                px: 1.5,
                                py: 0.75
                            }}
                        >
                            {option.label}
                        </Button>
                    ))}
                </Stack>
            </Box>

            {/* Футер */}
            <Box sx={{ 
                p: 2, 
                borderTop: 1,
                borderColor: 'divider'
            }}>
                <Button
                    variant="contained"
                    fullWidth
                    size="large"
                    onClick={handleApply}
                    startIcon={<Check />}
                >
                    Применить
                </Button>
            </Box>
        </Drawer>
    );
};

export { MobileListingCard, MobileHeader, MobileFilters };