// frontend/hostel-frontend/src/pages/marketplace/ListingDetailsPage.js
import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useLanguage } from '../../contexts/LanguageContext';
import { useParams } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import MiniMap from '../../components/maps/MiniMap';
import { PencilLine, Trash2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import ShareButton from '../../components/marketplace/ShareButton';
import ChatButton from '../../components/marketplace/chat/ChatButton';
import ReviewsSection from '../../components/reviews/ReviewsSection';
import FullscreenMap from '../../components/maps/FullscreenMap';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import CallButton from '../../components/marketplace/CallButton';
import { Link } from 'react-router-dom';
import { Store } from 'lucide-react';
import GalleryViewer from '../../components/shared/GalleryViewer';
import UserRating from '../../components/user/UserRating';
import AutoDetails from '../../components/marketplace/AutoDetails';
import SimilarListings from '../../components/marketplace/SimilarListings';

import {
    MapPin,
    Calendar,
    Heart,
    ChevronLeft,
    ChevronRight,
    Info,
    Tag,
    Settings,
    BarChart,
    Check,
    DollarSign,
    Filter,
    Zap
} from 'lucide-react';
import axios from '../../api/axios';
import {
    Container, Modal, Paper, Grid, Box, Typography,
    Button, Card, CardContent, Skeleton, Stack,
    Avatar, IconButton, useTheme, useMediaQuery,
    ImageList, ImageListItem, Chip
} from '@mui/material';
import { Percent } from 'lucide-react';
import PriceHistoryChart from '../../components/marketplace/PriceHistoryChart';

const ListingDetailsPage = () => {
    const { t, i18n } = useTranslation('marketplace');
    const { language } = useLanguage();
    const currentLanguage = useRef(language);
    const { id } = useParams();
    const theme = useTheme();
    const navigate = useNavigate();
    const { user, login } = useAuth();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const reviewsRef = useRef(null);

    // State declarations
    const [listing, setListing] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [currentImageIndex, setCurrentImageIndex] = useState(0);
    const [reviewsCount, setReviewsCount] = useState(0);
    const [categoryPath, setCategoryPath] = useState([]);
    const [isMapExpanded, setIsMapExpanded] = useState(false);
    const [categories, setCategories] = useState([]);
    const [galleryOpen, setGalleryOpen] = useState(false);
    const [isPriceHistoryOpen, setIsPriceHistoryOpen] = useState(false);

    const getAttributeIcon = (attributeName) => {
        // Определяем иконку в зависимости от типа атрибута
        switch (attributeName) {
            // Авто
            case 'make': return <Tag size={20} />;
            case 'model': return <Tag size={20} />;
            case 'year': return <Calendar size={20} />;
            case 'mileage': return <BarChart size={20} />;
            case 'engine_capacity': return <Zap size={20} />;
            case 'fuel_type': return <Filter size={20} />;
            case 'transmission': return <Settings size={20} />;
            case 'body_type': return <Settings size={20} />;
            case 'color': return <Info size={20} />;

            // Недвижимость
            case 'property_type': return <Tag size={20} />;
            case 'rooms': return <Tag size={20} />;
            case 'floor': return <BarChart size={20} />;
            case 'total_floors': return <BarChart size={20} />;
            case 'area': return <BarChart size={20} />;
            case 'land_area': return <BarChart size={20} />;
            case 'building_type': return <Tag size={20} />;
            case 'has_balcony': return <Check size={20} />;
            case 'has_elevator': return <Check size={20} />;
            case 'has_parking': return <Check size={20} />;

            // Телефоны
            case 'brand': return <Tag size={20} />;
            case 'model_phone': return <Tag size={20} />;
            case 'memory': return <BarChart size={20} />;
            case 'ram': return <BarChart size={20} />;
            case 'os': return <Settings size={20} />;
            case 'screen_size': return <BarChart size={20} />;
            case 'camera': return <BarChart size={20} />;
            case 'has_5g': return <Check size={20} />;

            // Компьютеры
            case 'pc_brand': return <Tag size={20} />;
            case 'pc_type': return <Tag size={20} />;
            case 'cpu': return <Settings size={20} />;
            case 'gpu': return <Settings size={20} />;
            case 'ram_pc': return <BarChart size={20} />;
            case 'storage_type': return <Tag size={20} />;
            case 'storage_capacity': return <BarChart size={20} />;
            case 'os_pc': return <Settings size={20} />;

            // По умолчанию
            default: return <Info size={20} />;
        }
    };
    // Замените функцию renderDiscountInfo:
    const renderDiscountInfo = () => {
        if (!listing || !listing.metadata || !listing.metadata.discount) return null;

        const discount = listing.metadata.discount;

        // Добавляем вывод в консоль для отладки
        console.log("Render discount info:", discount);

        return (
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    mt: 1,
                    mb: 2,
                    cursor: discount.has_price_history ? 'pointer' : 'default'
                }}
                onClick={() => {
                    console.log("Discount badge clicked, has_price_history:", discount.has_price_history);
                    if (discount.has_price_history) {
                        // Явно устанавливаем состояние модального окна в true
                        setIsPriceHistoryOpen(true);
                        console.log("Setting isPriceHistoryOpen to true");
                    }
                }}
            >
                <Chip
                    icon={<Percent size={16} />}
                    label={`-${discount.discount_percent}%`}
                    color="warning"
                    size="medium"
                    sx={{ mr: 1, fontWeight: 'bold' }}
                />
                <Typography variant="body2" color="text.secondary" sx={{ textDecoration: 'line-through' }}>
                    {formatPrice(discount.previous_price)}
                </Typography>
                {discount.has_price_history && (
                    <Typography
                        variant="caption"
                        color="primary.main"
                        sx={{ ml: 1 }}
                    >
                        {t('listings.details.priceHistory.showHistory')}
                    </Typography>
                )}
            </Box>
        );
    };
    const CAR_CATEGORY_ID = 2100;
    // Функция форматирования значения атрибута
    const formatAttributeValue = (attr) => {
        // Если атрибут не имеет значения - возвращаем "Не указано"
        if (!attr.display_value && attr.display_value !== 0 && attr.display_value !== false) {
            return t('common.not_specified', { defaultValue: 'Не указано' });
        }

        if (attr.attribute_type === 'boolean') {
            return attr.display_value === 'true' || attr.display_value === true ?
                t('common.yes') : t('common.no');
        }

        if (attr.attribute_name === 'price') {
            return formatPrice(attr.display_value);
        }

        // Особая обработка для некоторых числовых атрибутов
        if (attr.attribute_type === 'number') {
            let numValue;

            // Проверяем и преобразуем любое числовое значение
            if (typeof attr.display_value === 'string') {
                numValue = parseFloat(attr.display_value);
            } else if (typeof attr.display_value === 'number') {
                numValue = attr.display_value;
            } else if (attr.numeric_value !== undefined && attr.numeric_value !== null) {
                numValue = attr.numeric_value;
            } else {
                return t('common.not_specified', { defaultValue: 'Не указано' });
            }

            if (!isNaN(numValue)) {
                if (attr.attribute_name === 'year') {
                    return Math.round(numValue).toString(); // Убираем десятичные части для года
                }
                if (attr.attribute_name === 'mileage') {
                    return `${Math.round(numValue).toLocaleString()} км`; // Форматируем пробег с разделителями тысяч
                }
                if (attr.attribute_name === 'engine_capacity') {
                    return `${numValue.toFixed(1)} л`; // Форматируем объем с одним знаком после запятой
                }
                if (attr.attribute_name === 'power') {
                    return `${Math.round(numValue)} л.с.`; // Форматируем мощность двигателя
                }
            }
        }

        return attr.display_value;
    };

    // Преобразование атрибутов в формат для AutoDetails компонента
    const convertAttributesToAutoProps = (attributes) => {
        const props = {};

        console.log("Преобразуемые атрибуты:", attributes);

        attributes.forEach(attr => {
            if (!attr) return; // Пропускаем undefined и null

            let value;
            // Определяем, какое значение использовать в зависимости от типа атрибута
            if (attr.attribute_type === 'number' && attr.numeric_value !== undefined) {
                // Для числовых атрибутов берем numeric_value
                value = attr.numeric_value;

                // Для некоторых атрибутов преобразуем значение
                if (attr.attribute_name === 'mileage' && typeof value === 'number') {
                    // Округляем пробег до целого числа
                    value = Math.round(value);
                }
            } else if (attr.attribute_type === 'select' || attr.attribute_type === 'text') {
                // Для текстовых атрибутов берем text_value
                value = attr.text_value;
            } else if (attr.attribute_type === 'boolean') {
                // Для логических атрибутов берем boolean_value
                value = attr.boolean_value;
            } else {
                // В остальных случаях берем display_value
                value = attr.display_value;
            }

            switch (attr.attribute_name) {
                case 'make':
                    props.brand = value;
                    break;
                case 'model':
                    props.model = value;
                    break;
                case 'year':
                    props.year = parseInt(value) || null;
                    break;
                case 'mileage':
                    props.mileage = typeof value === 'number' ? value : parseInt(value) || 0;
                    break;
                case 'engine_capacity':
                    props.engine_capacity = typeof value === 'number' ? value : parseFloat(value) || null;
                    break;
                case 'power':
                    props.power = typeof value === 'number' ? value : parseInt(value) || null;
                    break;
                case 'fuel_type':
                    props.fuel_type = value;
                    break;
                case 'transmission':
                    props.transmission = value;
                    break;
                case 'body_type':
                    props.body_type = value;
                    break;
                case 'color':
                    props.color = value;
                    break;
                case 'drive_type':
                    props.drive_type = value;
                    break;
                case 'number_of_doors':
                    props.number_of_doors = value;
                    break;
                case 'number_of_seats':
                    props.number_of_seats = value;
                    break;
                default:
                    break;
            }
        });

        console.log("Результат преобразования атрибутов:", props);

        return props;
    };


    const findCategoryPath = useCallback((categoryId, categoriesTree) => {
        const path = [];

        const findPath = (id, categories) => {
            for (const category of categories) {
                if (String(category.id) === String(id)) {
                    path.unshift({
                        id: category.id,
                        name: category.name,
                        slug: category.slug,
                        translations: category.translations
                    });
                    return true;
                }

                if (category.children && findPath(id, category.children)) {
                    path.unshift({
                        id: category.id,
                        name: category.name,
                        slug: category.slug,
                        translations: category.translations
                    });
                    return true;
                }
            }
            return false;
        };

        findPath(categoryId, categoriesTree);
        return path;
    }, []);

    const getTranslatedText = (field) => {
        if (!listing || !field) return '';

        if (i18n.language === listing.original_language) {
            return listing[field];
        }

        const translation = listing.translations?.[i18n.language]?.[field];
        if (translation) {
            return translation;
        }

        return listing[field];
    };

    const fetchListing = useCallback(async () => {
        try {
            setLoading(true);
            const [listingResponse, favoritesResponse] = await Promise.all([
                axios.get(`/api/v1/marketplace/listings/${id}`),
                axios.get('/api/v1/marketplace/favorites')
            ]);

            const listingData = listingResponse.data.data;

            if (!listingData.images) {
                listingData.images = [];
            }

            setListing({
                ...listingData,
                is_favorite: favoritesResponse.data?.data?.some?.(
                    item => item.id === Number(id)
                ) || false,
                images: listingData.images || []
            });
            console.log("Listing data:", listingData);
            console.log("Storefront ID:", listingData.storefront_id);
        } catch (err) {
            console.error('Error fetching listing:', err);
            setError(t('listings.details.errors.loadFailed'));
        } finally {
            setLoading(false);
        }
    }, [id, t]);

    // All useEffects grouped together
    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const categoriesResponse = await axios.get('/api/v1/marketplace/category-tree');
                if (categoriesResponse.data?.data) {
                    setCategories(categoriesResponse.data.data);
                }
            } catch (err) {
                console.error('Error fetching categories:', err);
            }
        };

        fetchCategories();
    }, []);

    useEffect(() => {
        fetchListing();
    }, [fetchListing]);

    useEffect(() => {
        if (listing) {
            if (listing.category_path_ids && listing.category_path_ids.length > 0) {
                // Если с сервера пришла полная цепочка категорий, используем её
                const fullPath = listing.category_path_ids.map((id, index) => ({
                    id,
                    name: listing.category_path_names[index],
                    slug: listing.category_path_slugs[index],
                    translations: {} // По умолчанию пустые переводы
                }));
                setCategoryPath(fullPath);
            } else if (listing.category_id && categories.length > 0) {
                // Иначе используем старый метод построения пути
                const path = findCategoryPath(listing.category_id, categories);
                setCategoryPath(path);
            }
        }
    }, [listing, categories, findCategoryPath]);

    const scrollToReviews = () => {
        const reviewsSection = document.getElementById('reviews-section');
        if (reviewsSection) {
            reviewsSection.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('sr-RS', {
            style: 'currency',
            currency: 'RSD',
            maximumFractionDigits: 0
        }).format(price);
    };

    const handleDelete = async () => {
        if (!window.confirm(t('listings.details.actions.deleteConfirm'))) {
            return;
        }

        try {
            await axios.delete(`/api/v1/marketplace/listings/${id}`);
            navigate('/marketplace');
        } catch (error) {
            setError(t('listings.details.errors.deleteFailed'));
        }
    };

    const handleFavoriteClick = async () => {
        if (!user) {
            const returnUrl = window.location.pathname;
            const encodedReturnUrl = encodeURIComponent(returnUrl);
            login(`?returnTo=${encodedReturnUrl}`);
            return;
        }

        try {
            const newFavoriteState = !listing.is_favorite;
            setListing(prev => ({
                ...prev,
                is_favorite: newFavoriteState
            }));

            if (listing.is_favorite) {
                await axios.delete(`/api/v1/marketplace/listings/${id}/favorite`);
            } else {
                await axios.post(`/api/v1/marketplace/listings/${id}/favorite`);
            }

            await fetchListing();
        } catch (err) {
            setListing(prev => ({
                ...prev,
                is_favorite: !prev.is_favorite
            }));
            setError(t('listings.details.errors.updateFailed'));
        }
    };

    const getImageUrl = (image) => {
        if (!image) {
            return '';
        }

        const baseUrl = process.env.REACT_APP_BACKEND_URL;
        if (!baseUrl) {
            console.error('REACT_APP_BACKEND_URL is not defined!');
            return '';
        }

        if (typeof image === 'string') {
            return `${baseUrl}/uploads/${image}`;
        }

        if (image.file_path) {
            return `${baseUrl}/uploads/${image.file_path}`;
        }

        return '';
    };

    // Функция для получения всех путей к изображениям для галереи
    const getImagePaths = () => {
        if (!listing || !listing.images || listing.images.length === 0) {
            return [];
        }

        // Возвращаем пути file_path для передачи в GalleryViewer
        return listing.images.map(img => img.file_path);
    };

    // Обработчик клика по основному изображению
    const handleMainImageClick = () => {
        setGalleryOpen(true);
    };

    const formatMemberDate = (dateString) => {
        if (!dateString) return t('listings.details.seller.unknownDate');

        // Форматируем дату в соответствии с текущим языком
        const date = new Intl.DateTimeFormat(i18n.language, {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        }).format(new Date(dateString));

        // Используем переведенный шаблон из JSON с отформатированной датой
        return t('listings.details.seller.memberSince', { date });
    };

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Breadcrumbs paths={categoryPath} categories={categories} />
                <Grid container spacing={4}>
                    <Grid item xs={12} md={8}>
                        <Skeleton variant="rectangular" height={400} />
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <Skeleton variant="rectangular" height={200} />
                    </Grid>
                </Grid>
            </Container>
        );
    }

    if (error) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Breadcrumbs paths={categoryPath} categories={categories} />
                <Typography color="error">{error}</Typography>
            </Container>
        );
    }

    if (!listing) return null;

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Breadcrumbs paths={categoryPath} categories={categories} />
            <Grid container spacing={4}>
                {/* Images Gallery */}
                <Grid item xs={12} md={8}>
                    <Box sx={{ position: 'relative' }}>
                        {listing.images && Array.isArray(listing.images) && listing.images.length > 0 ? (
                            <>
                                <Box
                                    component="div"
                                    sx={{
                                        width: '100%',
                                        height: isMobile ? '300px' : '500px',
                                        borderRadius: 2,
                                        overflow: 'hidden',
                                        position: 'relative',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        backgroundColor: '#f5f5f5'
                                    }}
                                >
                                    <Box
                                        component="img"
                                        src={getImageUrl(listing.images[currentImageIndex])}
                                        alt={listing.title}
                                        sx={{
                                            maxWidth: '100%',
                                            maxHeight: '100%',
                                            objectFit: 'contain',
                                            cursor: 'pointer',
                                            width: 'auto',
                                            height: 'auto'
                                        }}
                                        onClick={handleMainImageClick}
                                    />
                                    {listing.images.length > 1 && (
                                        <>
                                            <IconButton
                                                aria-label={t('listings.details.image.prev')}
                                                sx={{
                                                    position: 'absolute',
                                                    left: 8,
                                                    top: '50%',
                                                    transform: 'translateY(-50%)',
                                                    bgcolor: 'background.paper',
                                                    '&:hover': { bgcolor: 'background.paper' }
                                                }}
                                                onClick={(e) => {
                                                    e.stopPropagation(); // Предотвращаем открытие галереи
                                                    setCurrentImageIndex(prev =>
                                                        prev > 0 ? prev - 1 : listing.images.length - 1
                                                    );
                                                }}
                                            >
                                                <ChevronLeft />
                                            </IconButton>
                                            <IconButton
                                                aria-label={t('listings.details.image.next')}
                                                sx={{
                                                    position: 'absolute',
                                                    right: 8,
                                                    top: '50%',
                                                    transform: 'translateY(-50%)',
                                                    bgcolor: 'background.paper',
                                                    '&:hover': { bgcolor: 'background.paper' }
                                                }}
                                                onClick={(e) => {
                                                    e.stopPropagation(); // Предотвращаем открытие галереи
                                                    setCurrentImageIndex(prev =>
                                                        prev < listing.images.length - 1 ? prev + 1 : 0
                                                    );
                                                }}
                                            >
                                                <ChevronRight />
                                            </IconButton>
                                        </>
                                    )}
                                </Box>
                            </>
                        ) : (
                            <Box
                                sx={{
                                    width: '100%',
                                    height: isMobile ? '300px' : '500px',
                                    bgcolor: 'grey.200',
                                    borderRadius: 2,
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center'
                                }}
                            >
                                <Typography color="text.secondary">
                                    {t('listings.details.title.noImages')}
                                </Typography>
                            </Box>
                        )}
                    </Box>

                    {/* Thumbnails */}
                    {listing.images && listing.images.length > 1 && (
                        <Box sx={{ mt: 2, overflowX: 'auto' }}>
                            <Stack
                                direction="row"
                                spacing={1}
                                sx={{
                                    py: 1,
                                    minHeight: 100
                                }}
                            >
                                {listing.images.map((image, index) => (
                                    <Box
                                        key={image.id || index}
                                        sx={{
                                            width: 100,
                                            height: 100,
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                            cursor: 'pointer',
                                            backgroundColor: '#f5f5f5',
                                            borderRadius: 1,
                                            overflow: 'hidden',
                                            opacity: currentImageIndex === index ? 1 : 0.6,
                                            border: currentImageIndex === index ? '2px solid #1976d2' : 'none',
                                            transition: 'opacity 0.2s',
                                            '&:hover': { opacity: 1 },
                                            flexShrink: 0
                                        }}
                                        onClick={() => setCurrentImageIndex(index)}
                                    >
                                        <Box
                                            component="img"
                                            src={getImageUrl(image)}
                                            alt={t('listings.details.image.preview', { number: index + 1 })}
                                            sx={{
                                                maxWidth: '100%',
                                                maxHeight: '100%',
                                                width: 'auto',
                                                height: 'auto',
                                                objectFit: 'contain'
                                            }}
                                        />
                                    </Box>
                                ))}
                            </Stack>
                        </Box>
                    )}

                    {/* Listing Description */}
                    <Box sx={{ mt: 4 }}>
                        <Typography variant="h4" gutterBottom>
                            {getTranslatedText('title')}
                        </Typography>

                        <Stack direction="row" spacing={2} sx={{ mb: 2 }}>
                            <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                                <MapPin size={18} style={{ marginRight: 4 }} />
                                <Typography>
                                    {listing.location || `${listing.city}, ${listing.country}`}
                                </Typography>
                            </Box>
                            <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                                <Calendar size={18} style={{ marginRight: 4 }} />
                                <Typography>
                                    {new Date(listing.created_at).toLocaleDateString()}
                                </Typography>
                            </Box>
                            <Box
                                component="button"
                                onClick={scrollToReviews}
                                sx={{
                                    background: 'none',
                                    border: 'none',
                                    color: 'primary.main',
                                    cursor: 'pointer',
                                    textDecoration: 'underline',
                                    padding: 0,
                                    display: 'flex',
                                    alignItems: 'center',
                                    '&:hover': { color: 'primary.dark' }
                                }}
                            >
                                {t('listings.details.info.reviews.count', { count: reviewsCount })}
                            </Box>
                        </Stack>

                        {/* Добавляем основные характеристики под основной информацией */}
                        {listing.attributes && listing.attributes.length > 0 && (
                            <Box sx={{ mb: 3 }}>
                                <Paper
                                    elevation={0}
                                    sx={{
                                        p: 2,
                                        borderRadius: 2,
                                        bgcolor: 'background.paper',
                                        border: '1px solid',
                                        borderColor: 'divider'
                                    }}
                                >
                                    <Grid container spacing={2}>
                                        {listing.attributes.slice(0, 3).map((attr) => (
                                            <Grid item xs={4} key={`top-${attr.attribute_id}`}>
                                                <Box sx={{
                                                    display: 'flex',
                                                    flexDirection: 'column',
                                                    alignItems: 'center',
                                                    textAlign: 'center'
                                                }}>
                                                    <Box sx={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        justifyContent: 'center',
                                                        width: 36,
                                                        height: 36,
                                                        borderRadius: '50%',
                                                        bgcolor: 'primary.main',
                                                        color: 'white',
                                                        mb: 1
                                                    }}>
                                                        {getAttributeIcon(attr.attribute_name)}
                                                    </Box>
                                                    <Typography variant="body2" color="text.secondary" gutterBottom>
                                                        {attr.display_name}
                                                    </Typography>
                                                    <Typography variant="body1" fontWeight="bold">
                                                        {formatAttributeValue(attr)}
                                                    </Typography>
                                                </Box>
                                            </Grid>
                                        ))}
                                    </Grid>
                                </Paper>
                            </Box>
                        )}

                        <Typography
                            variant="body1"
                            sx={{ mb: 4 }}
                            dangerouslySetInnerHTML={{ __html: getTranslatedText('description') }}
                        />

                        {/* Отображаем все атрибуты */}
                        {listing.attributes && listing.attributes.length > 0 && (
                            <Box sx={{ mt: 3, mb: 4 }}>
                                <Typography variant="h5" gutterBottom sx={{ mb: 2 }}>
                                    {t('listings.details.attributes.title', { defaultValue: 'Характеристики' })}
                                </Typography>

                                {/* Специальное отображение для автомобилей */}

                                {listing.category_id === CAR_CATEGORY_ID && (
                                    <AutoDetails autoProperties={convertAttributesToAutoProps(listing.attributes)} />
                                )}

                                {/* Стандартное отображение атрибутов */}
                                {listing.category_id !== 2100 && (
                                    <Paper sx={{ p: 2, borderRadius: 2 }}>
                                        <Grid container spacing={3}>
                                            {listing.attributes.map((attr) => (
                                                <Grid item xs={12} sm={6} key={attr.attribute_id}>
                                                    <Box sx={{
                                                        display: 'flex',
                                                        gap: 2,
                                                        alignItems: 'center',
                                                        p: 1,
                                                        borderRadius: 1,
                                                        '&:hover': { bgcolor: 'action.hover' }
                                                    }}>
                                                        <Box sx={{
                                                            display: 'flex',
                                                            alignItems: 'center',
                                                            justifyContent: 'center',
                                                            width: 40,
                                                            height: 40,
                                                            borderRadius: '50%',
                                                            bgcolor: 'primary.light',
                                                            color: 'primary.contrastText'
                                                        }}>
                                                            {getAttributeIcon(attr.attribute_name)}
                                                        </Box>
                                                        <Box>
                                                            <Typography variant="body2" color="text.secondary">
                                                                {attr.display_name}
                                                            </Typography>
                                                            <Typography variant="subtitle1" fontWeight="medium">
                                                                {formatAttributeValue(attr)}
                                                            </Typography>
                                                        </Box>
                                                    </Box>
                                                </Grid>
                                            ))}
                                        </Grid>
                                    </Paper>
                                )}
                            </Box>
                        )}

                        {/* Reviews section */}
                        <Box id="reviews-section" ref={reviewsRef} sx={{ mt: 4 }}>
                            <ReviewsSection
                                entityType="listing"
                                entityId={parseInt(id)}
                                entityTitle={listing.title}
                                canReview={user && user.id !== listing.user_id}
                                onReviewsCountChange={setReviewsCount}
                            />
                        </Box>
                    </Box>
                </Grid>

                {/* Right panel */}
                <Grid item xs={12} md={4}>
                    <Box sx={{ position: 'sticky', top: 24 }}>
                        {/* Price and contact card */}
                        <Card elevation={2}>
                            <CardContent>
                                <Typography variant="h4" gutterBottom>
                                    {formatPrice(listing.price)}
                                </Typography>
                                {renderDiscountInfo()}
                                {listing.storefront_id && (
                                    <Card elevation={2} sx={{ mt: 2 }}>
                                        <CardContent>
                                            <Typography variant="h6" gutterBottom>
                                                {t('listings.details.storeProduct')}
                                            </Typography>
                                            <Button
                                                variant="contained"
                                                color="primary"
                                                startIcon={<Store />}
                                                component={Link}
                                                to={`/shop/${listing.storefront_id}`}
                                                sx={{ mt: 1 }}
                                            >
                                                {t('listings.details.goToStore')}
                                            </Button>
                                        </CardContent>
                                    </Card>
                                )}

                                <Stack direction="row" spacing={1}>
                                    <Box sx={{ flex: 1 }}>
                                        <CallButton phone={listing.user?.phone} isMobile={isMobile} />
                                    </Box>
                                    <Box sx={{ flex: 1 }}>
                                        <ChatButton listing={listing} isMobile={isMobile} />
                                    </Box>
                                </Stack>

                                <Stack direction="row" spacing={1}>
                                    <Button
                                        variant="outlined"
                                        fullWidth
                                        startIcon={!isMobile && <Heart fill={listing?.is_favorite ? 'currentColor' : 'none'} />}
                                        onClick={handleFavoriteClick}
                                    >
                                        {isMobile ? (
                                            <Heart
                                                size={20}
                                                fill={listing?.is_favorite ? 'currentColor' : 'none'}
                                            />
                                        ) : t(`listings.details.favorite.${listing?.is_favorite ? 'remove' : 'add'}`)}
                                    </Button>
                                    <ShareButton
                                        url={window.location.href}
                                        title={listing.title}
                                        isMobile={isMobile}
                                    />
                                </Stack>

                                {/* Edit and delete buttons */}
                                {user?.id === listing.user_id && (
                                    <Stack direction="row" spacing={1} sx={{ mt: 2 }}>
                                        <Button
                                            variant="outlined"
                                            fullWidth
                                            startIcon={!isMobile && <PencilLine />}
                                            onClick={() => navigate(`/marketplace/listings/${id}/edit`)}
                                        >
                                            {isMobile ? <PencilLine size={20} /> : t('listings.details.actions.edit')}
                                        </Button>
                                        <Button
                                            variant="outlined"
                                            color="error"
                                            fullWidth
                                            startIcon={!isMobile && <Trash2 />}
                                            onClick={handleDelete}
                                        >
                                            {isMobile ? <Trash2 size={20} /> : t('listings.details.actions.delete')}
                                        </Button>
                                    </Stack>
                                )}
                            </CardContent>
                        </Card>

                        {/* Map card */}
                        {listing.latitude && listing.longitude ? (
                            listing.show_on_map ? (
                                <>
                                    <Card elevation={2} sx={{ mt: 2 }}>
                                        <CardContent sx={{ p: 1 }}>
                                            <MiniMap
                                                latitude={listing.latitude}
                                                longitude={listing.longitude}
                                                title={listing.title}
                                                address={listing.location}
                                                onClick={() => setIsMapExpanded(true)}
                                                onExpand={() => setIsMapExpanded(true)}
                                            />
                                        </CardContent>
                                    </Card>

                                    <Modal
                                        open={isMapExpanded}
                                        onClose={() => setIsMapExpanded(false)}
                                        sx={{
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                            p: 2
                                        }}
                                    >
                                        <Paper
                                            sx={{
                                                position: 'relative',
                                                width: '100%',
                                                maxWidth: 1200,
                                                maxHeight: '90vh',
                                                overflow: 'hidden'
                                            }}
                                        >
                                            <FullscreenMap
                                                latitude={listing.latitude}
                                                longitude={listing.longitude}
                                                title={listing.title}
                                            />
                                        </Paper>
                                    </Modal>
                                </>
                            ) : (
                                <Card elevation={2} sx={{ mt: 2 }}>
                                    <CardContent>
                                        <Stack direction="row" spacing={1} alignItems="center">
                                            <MapPin size={18} />
                                            <Typography>
                                                {`${listing.city}, ${listing.country}`}
                                            </Typography>
                                        </Stack>
                                    </CardContent>
                                </Card>
                            )
                        ) : null}

                        {/* Seller card */}
                        <Card elevation={2} sx={{ mt: 2 }}>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    {t('listings.details.seller.title')}
                                </Typography>
                                <Stack direction="row" spacing={2} alignItems="center" mb={2}>
                                    <Avatar
                                        src={listing.user?.picture_url}
                                        alt={listing.user?.name}
                                        sx={{ width: 56, height: 56 }}
                                    />
                                    <Box>
                                        <Typography variant="subtitle1">
                                            {listing.user?.name}
                                        </Typography>
                                        <Typography variant="body2" color="text.secondary">
                                            {formatMemberDate(listing.user?.created_at)}
                                        </Typography>
                                        <Button
                                            component={Link}
                                            to={`/user/${listing.user_id}/reviews`}
                                            size="small"
                                            variant="text"
                                            sx={{ p: 0, minWidth: 'auto' }}
                                        >
                                            {t('listings.details.seller.seeAllReviews')}
                                        </Button>
                                    </Box>
                                </Stack>
                                <UserRating userId={listing.user_id} />
                            </CardContent>
                        </Card>
                    </Box>
                </Grid>
            </Grid>

            {/* Полноэкранная галерея */}
            {listing.images && listing.images.length > 0 && (
                <GalleryViewer
                    images={getImagePaths()}
                    open={galleryOpen}
                    onClose={() => setGalleryOpen(false)}
                    initialIndex={currentImageIndex}
                    galleryMode="fullscreen"
                />
            )}

            {/* Модальное окно истории цен */}
            <Modal
                open={isPriceHistoryOpen} // Убедитесь, что здесь используется правильное имя переменной состояния
                onClose={() => {
                    console.log("Closing price history modal");
                    setIsPriceHistoryOpen(false);
                }}
                aria-labelledby="price-history-modal-title"
            >
                <Box sx={{
                    position: 'absolute',
                    top: '50%',
                    left: '50%',
                    transform: 'translate(-50%, -50%)',
                    width: isMobile ? '90%' : 800,
                    maxWidth: '100%',
                    bgcolor: 'background.paper',
                    borderRadius: 2,
                    boxShadow: 24,
                    p: 4,
                }}>
                    <Typography id="price-history-modal-title" variant="h6" component="h2" gutterBottom>
                        {t('listings.details.priceHistory.title')}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" paragraph>
                        {t('listings.details.priceHistory.description')}
                    </Typography>

                    <PriceHistoryChart listingId={id} />

                    <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 3 }}>
                        <Button onClick={() => setIsPriceHistoryOpen(false)}>
                            {t('common.close', { defaultValue: 'Закрыть' })}
                        </Button>
                    </Box>
                </Box>
            </Modal>
            <SimilarListings listingId={id} />
        </Container>
    );
};

export default ListingDetailsPage;