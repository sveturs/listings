// frontend/hostel-frontend/src/components/marketplace/ListingCard.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { MapPin as LocationIcon, Clock as AccessTime, Camera, Store } from 'lucide-react';
import { useNavigate } from 'react-router-dom'; // Используем хук для навигации
import axios from '../../api/axios';

import {
    Card,
    CardContent,
    CardMedia,
    Typography,
    Box,
    Chip,
    Button,
    Rating,
    Stack,
    Tooltip,
    Modal
} from '@mui/material';
import { Percent } from 'lucide-react';
import PriceHistoryChart from './PriceHistoryChart';
const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000';

const ListingCard = ({ listing, isMobile, onClick }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [storeName, setStoreName] = useState('Магазин');
    const navigate = useNavigate();
    const [isPriceHistoryOpen, setIsPriceHistoryOpen] = useState(false);
    
    // Добавляем логирование для отладки скидок
    useEffect(() => {
        console.log('Listing data:', listing);
        if (listing && listing.metadata) {
            console.log('Metadata found:', listing.metadata);
            if (listing.metadata.discount) {
                console.log('Discount found!', listing.metadata.discount);
            } else {
                console.log('No discount in metadata');
            }
        } else if (listing && listing.old_price) {
            console.log('Old price found:', listing.old_price);
        } else {
            console.log('No discount data found for listing:', listing.id);
        }
    }, [listing]);
    const getDiscountInfo = () => {
        // Проверяем наличие метаданных о скидке
        if (listing.metadata && listing.metadata.discount) {
            const previousPrice = Number(listing.metadata.discount.previous_price);
            
            // Всегда сами пересчитываем процент скидки на базе актуальной цены
            const calculatedPercent = Math.round((1 - listing.price / previousPrice) * 100);
            
            return {
                percent: calculatedPercent, // используем рассчитанный процент, а не сохраненный
                oldPrice: previousPrice,
                hasPriceHistory: listing.metadata.discount.has_price_history || false
            };
        }
        
        // Если есть флаг и старая цена
        if (listing.has_discount && listing.old_price) {
            const percent = Math.round((1 - listing.price / Number(listing.old_price)) * 100);
            return {
                percent: percent,
                oldPrice: listing.old_price,
                hasPriceHistory: false
            };
        }
        
        return null;
    };
    
    
    const renderDiscountBadge = () => {
        const discount = getDiscountInfo();
        
        if (discount) {
            return (
                <Box
                    sx={{
                        position: 'absolute',
                        top: 10,
                        left: 10,
                        zIndex: 2,
                        bgcolor: 'warning.main',
                        color: 'warning.contrastText',
                        borderRadius: '4px',
                        px: 1,
                        py: 0.5,
                        display: 'flex',
                        alignItems: 'center',
                        gap: 0.5,
                        fontSize: '0.875rem',
                        fontWeight: 'bold'
                    }}
                >
                    <Percent size={14} />
                    {`-${discount.percent}%`}
                </Box>
            );
        }
        
        return null;
    };
    
    
    // Загружаем название магазина при монтировании компонента
    useEffect(() => {
        const fetchStoreName = async () => {
            if (listing.storefront_id) {
                try {
                    const response = await axios.get(`/api/v1/public/storefronts/${listing.storefront_id}`);
                    if (response.data?.data?.name) {
                        setStoreName(response.data.data.name);
                    }
                } catch (err) {
                    console.error('Ошибка загрузки информации о магазине:', err);
                }
            }
        };

        fetchStoreName();
    }, [listing.storefront_id]);

    const getLocalizedText = (field) => {
        if (!listing || !field) return '';

        // Если текущий язык совпадает с языком оригинала
        if (i18n.language === listing.original_language) {
            return listing[field];
        }

        // Пытаемся получить перевод
        const translation = listing.translations?.[i18n.language]?.[field];
        if (translation) {
            return translation;
        }

        // Если перевод не найден, возвращаем оригинальный текст
        return listing[field];
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('sr-RS', {
            style: 'currency',
            currency: 'RSD',
            maximumFractionDigits: 0
        }).format(price || 0);
    };

    const formatDate = (dateString) => {
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

    const getMainImageUrl = () => {
        if (!listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
            return '/placeholder.jpg';
        }
        
        // Ищем главное изображение
        let mainImage = null;
        
        // Проверяем все возможные форматы хранения изображений
        if (Array.isArray(listing.images)) {
            // Ищем изображение с флагом is_main
            mainImage = listing.images.find(img => img && img.is_main === true);
            
            // Если главное изображение не найдено, берем первое
            if (!mainImage) {
                mainImage = listing.images[0];
            }
        }
        
        // Если путь содержит "pending", отображаем placeholder
        if (mainImage && typeof mainImage === 'object' && mainImage.file_path && 
            mainImage.file_path.includes('pending')) {
            return '/placeholder.jpg';
        }
        
        // Проверяем, что у нас есть объект с file_path
        if (mainImage && typeof mainImage === 'object' && mainImage.file_path) {
            return `${BACKEND_URL}/uploads/${mainImage.file_path}`;
        }
        
        // Если изображение передано как строка
        if (mainImage && typeof mainImage === 'string') {
            return `${BACKEND_URL}/uploads/${mainImage}`;
        }
        
        // Если ничего не подошло, возвращаем placeholder
        return '/placeholder.jpg';
    };


    const handleCardClick = (e) => {
        // Если клик был внутри элемента с атрибутом data-shop-button или detailsButton, не выполняем навигацию
        if (e.target.closest('[data-shop-button="true"]') || e.target.closest('#detailsButton')) {
            return;
        }

        if (onClick) {
            onClick(listing);
        } else {
            // Используем navigate вместо window.location
            navigate(`/marketplace/listings/${listing.id}`);
        }
    };

    const handleShopButtonClick = (e) => {
        e.preventDefault();
        e.stopPropagation();
        navigate(`/shop/${listing.storefront_id}`);
    };

    const handleDetailsButtonClick = (e) => {
        e.preventDefault();
        e.stopPropagation();
        navigate(`/marketplace/listings/${listing.id}`);
    };
    
    const getDisplayLocation = () => {
        // Если есть город, используем его
        if (listing.city) {
            return listing.city;
        }

        // Если есть полный адрес, извлекаем из него город или первую часть
        if (listing.location) {
            // Пытаемся извлечь город из адреса (обычно он в начале)
            const locationParts = listing.location.split(',');
            if (locationParts.length > 0) {
                return locationParts[0].trim();
            }
            return listing.location;
        }

        // Если нет ни города, ни адреса
        return 'Местоположение не указано';
    };

    // Метод для рендеринга секции с ценой и скидкой
    const renderPriceSection = () => {
        const discount = getDiscountInfo();
        
        return (
            <Box>
                <Typography
                    variant={isMobile ? "body2" : "h5"}
                    color="primary"
                    sx={{
                        mt: 0.5,
                        fontSize: isMobile ? '0.875rem' : undefined,
                        fontWeight: 'bold'
                    }}
                >
                    {formatPrice(listing.price)}
                </Typography>
                
                {/* Отображение старой цены, если есть информация о скидке */}
                {discount && (
                    <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{
                            textDecoration: 'line-through',
                            mt: 0.5
                        }}
                    >
                        {formatPrice(discount.oldPrice)}
                    </Typography>
                )}
            </Box>
        );
    };
    

    return (
        <Card
            sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                position: 'relative',
                '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: 3,
                    transition: 'all 0.2s ease-in-out'
                },
                cursor: 'pointer'
            }}
            onClick={handleCardClick}
        >
            {/* Бейдж магазина */}
            {listing.storefront_id && (
                <Box
                    sx={{
                        position: 'absolute',
                        top: 10,
                        right: 10,
                        zIndex: 9999,
                        bgcolor: 'primary.main',
                        color: 'white',
                        borderRadius: '4px',
                        px: 1,
                        py: 0.5,
                        display: 'flex',
                        alignItems: 'center',
                        gap: 0.5,
                        fontSize: '0.75rem',
                        fontWeight: 'bold',
                        cursor: 'pointer',
                        pointerEvents: 'auto',
                        opacity: 1,
                        '&::before': {},
                    }}
                    onClick={handleShopButtonClick}
                    data-shop-button="true"
                >
                    <Store size={14} />
                    в магазин
                </Box>
            )}
            
            {/* Бейдж скидки */}
            {renderDiscountBadge()}

            <Box sx={{ position: 'relative', pt: isMobile ? '100%' : '75%' }}>
                <CardMedia
                    component="img"
                    sx={{
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        width: '100%',
                        height: '100%',
                        objectFit: 'cover'
                    }}
                    image={getMainImageUrl()}
                    alt={getLocalizedText('title') || 'Изображение отсутствует'}
                />
                {listing.images && listing.images.length > 1 && !isMobile && (
                    <Chip
                        icon={<Camera size={16} />}
                        label={`${listing.images.length} ${t('listings.details.title.photoCount', { count: listing.images.length }).split(' ')[1]}`}
                        size="small"
                        sx={{
                            position: 'absolute',
                            bottom: 8,
                            right: 8,
                            bgcolor: 'rgba(0,0,0,0.6)',
                            color: 'white'
                        }}
                    />
                )}
            </Box>

            <CardContent sx={{
                flexGrow: 1,
                p: isMobile ? 1 : 2,
                '&:last-child': { pb: isMobile ? 1 : 2 }
            }}>
                <Typography
                    variant={isMobile ? "body2" : "h6"}
                    noWrap
                    sx={{
                        fontSize: isMobile ? '0.875rem' : undefined,
                        fontWeight: 'medium'
                    }}
                >
                    {getLocalizedText('title') || 'Без названия'}
                </Typography>

                {listing.rating > 0 && (
                    <Stack direction="row" spacing={0.5} alignItems="center" sx={{ mt: 1 }}>
                        <Rating
                            value={listing.rating}
                            readOnly
                            size="small"
                            precision={0.1}
                        />
                        <Typography
                            variant="body2"
                            color="text.secondary"
                        >
                            ({listing.reviews_count})
                        </Typography>
                    </Stack>
                )}

                {/* Использование обновленного метода для отображения цены и скидки */}
                {renderPriceSection()}

                {!isMobile && (
                    <>
                        <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                            <LocationIcon size={18} style={{ marginRight: 4 }} />
                            <Typography variant="body2" noWrap>
                                {getDisplayLocation()}
                            </Typography>
                        </Box>

                        <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                            <AccessTime size={18} style={{ marginRight: 4 }} />
                            <Typography variant="body2">
                                {formatDate(listing.created_at)}
                            </Typography>
                        </Box>

                        <Button
                            id="detailsButton"
                            variant="contained"
                            fullWidth
                            sx={{ mt: 2 }}
                            onClick={handleDetailsButtonClick}
                        >
                            {t('listings.details.moreDetails')}
                        </Button>
                    </>
                )}
            </CardContent>
            
            {/* Модальное окно для истории цен */}
            {isPriceHistoryOpen && listing.metadata && listing.metadata.discount && listing.metadata.discount.has_price_history && (
                <Modal
                    open={isPriceHistoryOpen}
                    onClose={() => setIsPriceHistoryOpen(false)}
                    aria-labelledby="price-history-title"
                >
                    <Box sx={{
                        position: 'absolute',
                        top: '50%',
                        left: '50%',
                        transform: 'translate(-50%, -50%)',
                        width: isMobile ? '90%' : 600,
                        bgcolor: 'background.paper',
                        borderRadius: 2,
                        boxShadow: 24,
                        p: 4,
                    }}>
                        <Typography id="price-history-title" variant="h6" component="h2" gutterBottom>
                            История изменения цены
                        </Typography>
                        <PriceHistoryChart listingId={listing.id} />
                        <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                            <Button onClick={() => setIsPriceHistoryOpen(false)}>
                                Закрыть
                            </Button>
                        </Box>
                    </Box>
                </Modal>
            )}
        </Card>
    );
};

export default ListingCard;