// frontend/hostel-frontend/src/components/marketplace/MarketplaceListingsList.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    Checkbox,
    Typography,
    IconButton,
    Chip,
    useTheme,
    useMediaQuery,
    TableSortLabel
} from '@mui/material';
import { Edit, Trash2, Eye, Calendar, Percent, ArrowUpDown, Star } from 'lucide-react';
import axios from '../../api/axios';

// Вспомогательные функции для форматирования
const formatPrice = (price) => {
    return new Intl.NumberFormat('sr-RS', {
        style: 'currency',
        currency: 'RSD',
        maximumFractionDigits: 0
    }).format(price || 0);
};

const formatDate = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString();
};

// Вспомогательная функция для форматирования числа отзывов
const formatReviewCount = (count) => {
    if (count === undefined || count === null) return '0';
    return count.toString();
};

const getImageUrl = (listing) => {
    if (!listing || !listing.images || listing.images.length === 0) {
        return '/placeholder.jpg';
    }

    const baseUrl = process.env.REACT_APP_BACKEND_URL || '';

    // Найдем главное изображение или используем первое
    const mainImage = listing.images.find(img => img.is_main) || listing.images[0];

    if (typeof mainImage === 'string') {
        return `${baseUrl}/uploads/${mainImage}`;
    }

    // Проверяем на Minio хранилище
    if (mainImage.storage_type === 'minio' || 
        (mainImage.file_path && mainImage.file_path.includes('listings/'))) {
        return `${baseUrl}/listings/${mainImage.file_path.split('/').pop()}`;
    }

    if (mainImage && mainImage.file_path) {
        return `${baseUrl}/uploads/${mainImage.file_path}`;
    }

    return '/placeholder.jpg';
};


const MarketplaceListingsList = ({
    listings,
    selectedItems = [],
    onSelectItem = null,
    onSelectAll = null,
    showSelection = false,
    onSortChange = null,
    filters = {},
    initialSortField = 'created_at',
    initialSortOrder = 'desc',
    loading = false 
}) => {

    console.log("MarketplaceListingsList: получены listings:", listings?.length || 0);

    const { t, i18n } = useTranslation('marketplace');
    const navigate = useNavigate();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const [order, setOrder] = useState(initialSortOrder);
    const [orderBy, setOrderBy] = useState(initialSortField);

    // Состояние для хранения данных о рейтингах
    const [listingsWithRatings, setListingsWithRatings] = useState([]);
    const [loadingRatings, setLoadingRatings] = useState(false);
    
    // Определяем колонки таблицы
    const [columns, setColumns] = useState(['image', 'title', 'price', 'reviews', 'date', 'status', 'promotions']);

    // Функция для получения локализованного текста
    const getLocalizedText = (listing, field) => {
        if (!listing) return '';
        
        // Получаем перевод для текущего языка интерфейса
        const currentLanguage = i18n.language;
        
        // Проверяем, совпадает ли текущий язык с оригинальным языком объявления
        if (currentLanguage === listing.original_language) {
            return listing[field];
        }
        
        // Проверяем наличие перевода
        if (listing.translations && 
            listing.translations[currentLanguage] && 
            listing.translations[currentLanguage][field]) {
            return listing.translations[currentLanguage][field];
        }
        
        // Если перевода нет, возвращаем оригинальный текст
        return listing[field];
    };

    // Загрузка рейтингов для объявлений
    useEffect(() => {
        if (!listings || listings.length === 0) {
            setListingsWithRatings([]);
            return;
        }

        // Функция для загрузки рейтингов
        const fetchRatings = async () => {
            setLoadingRatings(true);

            try {
                console.log(`Загрузка рейтингов для ${listings.length} объявлений`);

                // Создаем копию списка объявлений
                const enhancedListings = [...listings];

                // Получаем рейтинги для каждого объявления
                const ratingPromises = listings.map(listing =>
                    axios.get(`/api/v1/entity/listing/${listing.id}/stats`)
                        .then(response => {
                            const stats = response.data?.data;
                            console.log(`Рейтинг для объявления ${listing.id}:`, stats);
                            return {
                                id: listing.id,
                                review_count: stats?.total_reviews || 0,
                                average_rating: stats?.average_rating || 0
                            };
                        })
                        .catch(error => {
                            console.error(`Ошибка при загрузке рейтинга для объявления ${listing.id}:`, error);
                            return { id: listing.id, review_count: 0, average_rating: 0 };
                        })
                );

                // Ждем завершения всех запросов
                const ratingsResults = await Promise.all(ratingPromises);

                // Добавляем информацию о рейтингах к объявлениям
                ratingsResults.forEach(ratingData => {
                    const listingIndex = enhancedListings.findIndex(l => l.id === ratingData.id);
                    if (listingIndex !== -1) {
                        enhancedListings[listingIndex] = {
                            ...enhancedListings[listingIndex],
                            review_count: ratingData.review_count,
                            average_rating: ratingData.average_rating
                        };
                    }
                });

                console.log(`Объявления с рейтингами:`, enhancedListings.map(l => ({
                    id: l.id,
                    title: l.title,
                    reviews: l.review_count,
                    rating: l.average_rating
                })));

                setListingsWithRatings(enhancedListings);
            } catch (error) {
                console.error('Ошибка при загрузке рейтингов:', error);
                setListingsWithRatings(listings); // Используем исходные данные без рейтингов
            } finally {
                setLoadingRatings(false);
            }
        };

        fetchRatings();
    }, [listings]);
    
    // Синхронизация состояния сортировки с фильтрами
    useEffect(() => {
        // Если есть sort_by в формате field_direction
        if (filters && filters.sort_by) {
            const sortParts = filters.sort_by.split('_');
            if (sortParts.length >= 2) {
                // Получаем поле и направление
                let field = sortParts[0];
                let direction = sortParts.pop(); // Последний элемент - направление

                // Преобразуем поле API обратно в поле UI
                switch (field) {
                    case 'date':
                        setOrderBy('created_at');
                        break;
                    case 'price':
                    case 'title':
                    case 'rating':
                        setOrderBy(field === 'rating' ? 'reviews' : field);
                        break;
                    default:
                        setOrderBy('created_at');
                }

                if (direction === 'asc' || direction === 'desc') {
                    setOrder(direction);
                }
            }
        }
    }, [filters]);

    // Адаптация для мобильных: скрываем некоторые колонки


    // Проверяем, выбраны ли все элементы
    const isAllSelected = listings.length > 0 && selectedItems.length === listings.length;

    // Обработчик клика по строке - переход к объявлению
    const handleRowClick = (id) => {
        navigate(`/marketplace/listings/${id}`);
    };

    // Обработчик изменения сортировки
    const handleRequestSort = (property) => {
        const isAsc = orderBy === property && order === 'asc';
        const newOrder = isAsc ? 'desc' : 'asc';

        // Обновляем локальное состояние компонента
        setOrder(newOrder);
        setOrderBy(property);

        // Если предоставлен колбэк для внешней сортировки, вызываем его с корректными параметрами
        if (onSortChange) {
            onSortChange(property, newOrder);
        }
    };

    const createSortHandler = (property) => () => {
        const isAsc = orderBy === property && order === 'asc';
        const newOrder = isAsc ? 'desc' : 'asc';
        
        console.log(`SORT: Передаем в родительский компонент: поле=${property}, порядок=${newOrder}`);
        
        if (onSortChange) {
            // Вызываем колбэк с параметрами сортировки
            onSortChange(property, newOrder);
            // Обновляем локальное состояние
            setOrder(newOrder);
            setOrderBy(property);
        }
    };

    const getDiscountInfo = (listing) => {
        if (listing.metadata && listing.metadata.discount) {
            return {
                percent: listing.metadata.discount.discount_percent,
                oldPrice: listing.metadata.discount.previous_price
            };
        }
        if (listing.has_discount && listing.old_price) {
            const percent = Math.round((1 - listing.price / Number(listing.old_price)) * 100);
            return {
                percent: percent,
                oldPrice: listing.old_price
            };
        }
        return null;
    };

    const displayListings = listingsWithRatings.length > 0 ? listingsWithRatings : listings;

    return (
        <TableContainer component={Paper} elevation={0} variant="outlined">
            <Table sx={{ minWidth: 650 }}>
                <TableHead>
                    <TableRow>
                        {showSelection && onSelectAll && (
                            <TableCell padding="checkbox">
                                <Checkbox
                                    indeterminate={selectedItems.length > 0 && selectedItems.length < listings.length}
                                    checked={isAllSelected}
                                    onChange={(e) => onSelectAll(e.target.checked)}
                                />
                            </TableCell>
                        )}

                        {columns.includes('image') && (
                            <TableCell width={80}></TableCell>
                        )}

                        {columns.includes('title') && (
                            <TableCell>
                                <TableSortLabel
                                    active={orderBy === 'title'}
                                    direction={orderBy === 'title' ? order : 'asc'}
                                    onClick={createSortHandler('title')}
                                >
                                    {t('listings.table.title')}
                                </TableSortLabel>
                            </TableCell>
                        )}

                        {columns.includes('price') && (
                            <TableCell align="right">
                                <TableSortLabel
                                    active={orderBy === 'price'}
                                    direction={orderBy === 'price' ? order : 'asc'}
                                    onClick={createSortHandler('price')}
                                >
                                    {t('listings.table.price')}
                                </TableSortLabel>
                            </TableCell>
                        )}

                        {columns.includes('reviews') && (
                            <TableCell align="center">
                                <TableSortLabel
                                    active={orderBy === 'reviews'}
                                    direction={orderBy === 'reviews' ? order : 'asc'}
                                    onClick={createSortHandler('reviews')}
                                >
                                    {t('reviews.title', { defaultValue: 'reviews' })}

                                </TableSortLabel>
                            </TableCell>
                        )}

                        {columns.includes('date') && (
                            <TableCell>
                                <TableSortLabel
                                    active={orderBy === 'created_at'}
                                    direction={orderBy === 'created_at' ? order : 'asc'}
                                    onClick={createSortHandler('created_at')}
                                >
                                    {t('listings.table.date')}
                                </TableSortLabel>
                            </TableCell>
                        )}
                        
                        {columns.includes('status') && (
                            <TableCell align="center">
                                {t('listings.table.status')}
                            </TableCell>
                        )}
                        
                        {columns.includes('promotions') && (
                            <TableCell align="center">
                                {t('listings.table.promotions')}
                            </TableCell>
                        )}
                    </TableRow>
                </TableHead>

                <TableBody>
                    {displayListings.map((listing) => {
                        const isSelected = selectedItems.includes(listing.id);
                        const discount = getDiscountInfo(listing);

                        // Применяем локализацию текста
                        const localizedTitle = getLocalizedText(listing, 'title');
                        const localizedDescription = getLocalizedText(listing, 'description');

                        return (
                            <TableRow
                                key={listing.id}
                                hover
                                onClick={() => handleRowClick(listing.id)}
                                selected={isSelected}
                                sx={{
                                    cursor: 'pointer',
                                    '&:last-child td, &:last-child th': { border: 0 }
                                }}
                            >
                                {showSelection && onSelectItem && (
                                    <TableCell padding="checkbox" onClick={(e) => {
                                        e.stopPropagation();
                                        onSelectItem(listing.id);
                                    }}>
                                        <Checkbox checked={isSelected} />
                                    </TableCell>
                                )}

                                {columns.includes('image') && (
                                    <TableCell width={80}>
                                        <Box
                                            component="img"
                                            src={getImageUrl(listing)}
                                            alt={localizedTitle}
                                            sx={{
                                                width: 60,
                                                height: 60,
                                                objectFit: 'contain',
                                                backgroundColor: '#f5f5f5',
                                                borderRadius: 1
                                            }}
                                        />
                                    </TableCell>
                                )}

                                {columns.includes('title') && (
                                    <TableCell>
                                        <Typography variant="subtitle2">{localizedTitle}</Typography>
                                        {discount && (
                                            <Chip
                                                icon={<Percent size={12} />}
                                                label={`-${discount.percent}%`}
                                                color="warning"
                                                size="small"
                                                sx={{ mt: 0.5, height: 20, fontSize: '0.7rem' }}
                                            />
                                        )}
                                    </TableCell>
                                )}

                                {columns.includes('price') && (
                                    <TableCell align="right">
                                        <Typography variant="subtitle2" color="primary.main">
                                            {formatPrice(listing.price)}
                                        </Typography>
                                        {discount && (
                                            <Typography
                                                variant="caption"
                                                color="text.secondary"
                                                sx={{ textDecoration: 'line-through', display: 'block' }}
                                            >
                                                {formatPrice(discount.oldPrice)}
                                            </Typography>
                                        )}
                                    </TableCell>
                                )}

                                {columns.includes('reviews') && (
                                    <TableCell align="center">
                                        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                                            {(listing.review_count !== undefined && listing.review_count > 0) ? (
                                                <>
                                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                                        <Typography variant="body2" fontWeight="medium">
                                                            {listing.average_rating !== undefined ? listing.average_rating.toFixed(1) : '0.0'}
                                                        </Typography>
                                                        <Star
                                                            size={16}
                                                            fill={listing.average_rating > 0 ? theme.palette.warning.main : 'none'}
                                                            color={listing.average_rating > 0 ? theme.palette.warning.main : theme.palette.text.disabled}
                                                        />
                                                    </Box>
                                                    <Typography variant="caption" color="text.secondary">
                                                    {listing.review_count} {t('reviews.info.reviews.count', { defaultValue: 'reviews', count: listing.review_count })}
                                                    </Typography>
                                                </>
                                            ) : (
                                                <Typography variant="caption" color="text.secondary">
                                                    {t('reviews.noRatingsYet', { defaultValue: 'Нет отзывов' })}
                                                </Typography>
                                            )}
                                        </Box>
                                    </TableCell>
                                )}

                                {columns.includes('date') && (
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                            <Calendar size={16} />
                                            <Typography variant="body2">
                                                {formatDate(listing.created_at)}
                                            </Typography>
                                        </Box>
                                    </TableCell>
                                )}
                                
                                {columns.includes('status') && (
                                    <TableCell align="center">
                                        <Chip
                                            label={listing.status === 'active' 
                                                ? t('listings.status.active')
                                                : t('listings.status.inactive')
                                            }
                                            color={listing.status === 'active' ? 'success' : 'default'}
                                            size="small"
                                            variant={listing.status === 'active' ? "filled" : "outlined"}
                                            sx={{ minWidth: 80 }}
                                        />
                                    </TableCell>
                                )}
                                
                                {columns.includes('promotions') && (
                                    <TableCell align="center">
                                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5, alignItems: 'center' }}>
                                            {listing.metadata && listing.metadata.promotions ? (
                                                Object.entries(listing.metadata.promotions).map(([type, data]) => (
                                                    <Chip
                                                        key={type}
                                                        label={t(`listings.promotions.${type}`, { defaultValue: type })}
                                                        color="primary"
                                                        size="small"
                                                        variant="outlined"
                                                        sx={{ fontSize: '0.7rem', height: 20 }}
                                                    />
                                                ))
                                            ) : (
                                                <Typography variant="caption" color="text.secondary">
                                                    {t('listings.promotions.none', { defaultValue: 'Нет' })}
                                                </Typography>
                                            )}
                                        </Box>
                                    </TableCell>
                                )}
                            </TableRow>
                        );
                    })}

                    {displayListings.length === 0 && (
                        <TableRow>
                            <TableCell colSpan={showSelection ? columns.length + 1 : columns.length} align="center">
                                <Typography variant="body2" color="text.secondary" sx={{ py: 3 }}>
                                    {t('listings.table.noListings')}
                                </Typography>
                            </TableCell>
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </TableContainer>
    );
};

export default MarketplaceListingsList;