// frontend/hostel-frontend/src/components/marketplace/ListingCard.tsx
import React, { useState, useEffect, MouseEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { MapPin as LocationIcon, Clock as AccessTime, Camera, Store, Eye } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import axios from '../../api/axios';
import { Listing, ListingMetadata, DiscountInfo, Attribute, ListingImage } from '../../types/listing';
// Import Modal wrapper component
import ModalWrapper from './ModalWrapper';

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
    Modal
} from '@mui/material';
import { Percent } from 'lucide-react';
import PriceHistoryChart from './PriceHistoryChart';

interface Promotion {
    [type: string]: any;
}

interface DiscountData {
    percent: number;
    oldPrice: number;
    hasPriceHistory: boolean;
}

interface ListingCardProps {
    listing: Listing;
    isMobile?: boolean;
    onClick?: (listing: Listing) => void;
    showStatus?: boolean;
}

const ListingCard: React.FC<ListingCardProps> = ({ listing, isMobile = false, onClick, showStatus = false }) => {
    const { t, i18n } = useTranslation('marketplace') as any;
    const [storeName, setStoreName] = useState<string>(t('listings.details.seller.title'));
    const navigate = useNavigate();
    const [isPriceHistoryOpen, setIsPriceHistoryOpen] = useState<boolean>(false);

    // Функция для определения, является ли объявление частью витрины
    const isStoreItem = (): boolean => {
        // Единственное надежное условие - наличие storefront_id
        return listing.storefront_id !== undefined && listing.storefront_id !== null;
    };

    // Получаем ID витрины
    const getStoreId = (): number | string | null => {
        // Просто возвращаем storefront_id, если он есть
        return listing.storefront_id || null;
    };

    // Расширенное логирование для отладки
    useEffect(() => {
        console.log('Listing data:', listing);

        // Проверка поля storefront_id для бирки "в магазин"
        console.log('Storefront ID:', listing.storefront_id,
            'Type:', typeof listing.storefront_id,
            'isPartOfStorefront:', listing.isPartOfStorefront,
            'storefront_name:', listing.storefront_name,
            'storefrontName:', listing.storefrontName);

        // Проверка данных о скидке
        console.log('Discount fields:', {
            has_discount: listing.has_discount,
            old_price: listing.old_price,
            price: listing.price
        });

        if (listing && listing.metadata) {
            console.log('Metadata found:', listing.metadata);
            if (listing.metadata.discount) {
                console.log('Discount found in metadata!', listing.metadata.discount);
                console.log('Discount data details:', {
                    percent: listing.metadata.discount.percent || listing.metadata.discount.discount_percent,
                    oldPrice: listing.metadata.discount.oldPrice || listing.metadata.discount.previous_price
                });
            } else {
                console.log('No discount in metadata');
            }
        } else {
            console.log('No metadata or listing is undefined');
        }

        // Проверка вычисленной скидки
        const discountData = getDiscountData();
        console.log('Calculated discount data:', discountData);
    }, [listing]);

    // Получение имени продавца
    useEffect(() => {
        // Сначала проверяем, есть ли имя витрины уже в данных листинга
        if (listing && (listing.storefront_name || listing.storefrontName)) {
            setStoreName(listing.storefront_name || listing.storefrontName);
            console.log('Using storefront name from listing data:', listing.storefront_name || listing.storefrontName);
            return;
        }

        // Если id витрины есть, но имени нет - запрашиваем с сервера
        if (isStoreItem()) {
            console.log('Fetching storefront name by ID:', getStoreId());

            axios.get(`/api/v1/public/storefronts/${getStoreId()}`)
                .then(response => {
                    if (response.data && response.data.data) {
                        const name = response.data.data.name || t('listings.details.seller.title');
                        console.log('Fetched storefront name:', name);
                        setStoreName(name);
                    }
                })
                .catch(error => {
                    console.error('Error fetching storefront:', error);
                    // В случае ошибки используем значение по умолчанию
                    setStoreName(t('listings.details.seller.title'));
                });
        } else if (listing && listing.isPartOfStorefront) {
            // Если есть флаг, но нет id - используем текст по умолчанию
            console.log('Using default storefront name for isPartOfStorefront flag');
            setStoreName(t('listings.details.seller.inStore'));
        }
    }, [listing, t]);

    const handleCardClick = (e: MouseEvent<HTMLDivElement>) => {
        // Предотвращаем переход, если клик был по кнопке
        if ((e.target as HTMLElement).closest('button')) {
            return;
        }

        if (onClick) {
            onClick(listing);
        } else {
            navigate(`/marketplace/listings/${listing.id}`);
        }
    };

    const getDiscountData = (): DiscountData | null => {
        // Проверка наличия полей для скидки в нескольких форматах

        // Вариант 1: Стандартные поля скидки (has_discount + old_price)
        if (listing.has_discount && listing.old_price && listing.price &&
            listing.old_price > listing.price) {
            // Если есть поле has_discount и old_price
            const percent = Math.round(((listing.old_price - listing.price) / listing.old_price) * 100);
            return {
                percent,
                oldPrice: listing.old_price,
                hasPriceHistory: false
            };
        }

        // Вариант 2: Прямое указание старой цены без флага has_discount
        else if (listing.old_price && listing.price &&
                listing.old_price > listing.price) {
            const percent = Math.round(((listing.old_price - listing.price) / listing.old_price) * 100);
            return {
                percent,
                oldPrice: listing.old_price,
                hasPriceHistory: false
            };
        }

        // Вариант 3: Метаданные о скидке
        else if (listing.metadata?.discount) {
            // Если есть метаданные скидки
            const discount = listing.metadata.discount;
            const oldPrice = discount.previous_price || discount.oldPrice || 0;
            const percent = discount.discount_percent || discount.percent || 0;

            // Проверяем, что есть значимые данные о скидке
            if (oldPrice > 0 && percent > 0) {
                return {
                    percent,
                    oldPrice,
                    hasPriceHistory: Boolean(discount.has_price_history)
                };
            }

            // Если есть старая цена, но нет процента - вычисляем процент
            if (oldPrice > 0 && listing.price && oldPrice > listing.price) {
                const calculatedPercent = Math.round(((oldPrice - listing.price) / oldPrice) * 100);
                return {
                    percent: calculatedPercent,
                    oldPrice,
                    hasPriceHistory: Boolean(discount.has_price_history)
                };
            }
        }

        return null;
    };

    const discountData = getDiscountData();

    const formatDate = (dateString?: string) => {
        if (!dateString) return '';
        try {
            const date = new Date(dateString);
            return date.toLocaleDateString(i18n.language);
        } catch (e) {
            return dateString;
        }
    };

    // Получение первого изображения
    const getImageUrl = () => {
        if (!listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
            return '/placeholder-listing.jpg';
        }

        const firstImage = listing.images[0];
        const baseUrl = (window as any)?.ENV?.REACT_APP_MINIO_URL || (window as any)?.ENV?.REACT_APP_BACKEND_URL || '';

        // Если это строка URL
        if (typeof firstImage === 'string') {
            // Проверяем, абсолютный или относительный URL
            if (firstImage.startsWith('http')) {
                return firstImage;
            } else {
                return `${baseUrl}/uploads/${firstImage}`;
            }
        }

        // Если это объект
        if (firstImage && typeof firstImage === 'object') {
            // Проверяем на наличие публичного URL
            if (firstImage.public_url) {
                // Проверяем, абсолютный или относительный URL
                if (firstImage.public_url.startsWith('http')) {
                    return firstImage.public_url;
                } else {
                    return `${baseUrl}${firstImage.public_url}`;
                }
            }

            // Проверяем на MinIO хранилище
            if (firstImage.storage_type === 'minio' ||
                (firstImage.file_path && firstImage.file_path.includes('listings/'))) {
                return `${baseUrl}${firstImage.public_url}`;
            }

            // Проверяем file_path
            if (firstImage.file_path) {
                return `${baseUrl}/uploads/${firstImage.file_path}`;
            }
        }

        return '/placeholder-listing.jpg';
    };

    // Обработчик открытия истории цен
    const handleOpenPriceHistory = (e: MouseEvent<HTMLButtonElement>) => {
        e.stopPropagation();
        setIsPriceHistoryOpen(true);
    };

    // Обработчик закрытия истории цен
    const handleClosePriceHistory = () => {
        setIsPriceHistoryOpen(false);
    };

    return (
        <>
            <Card
                sx={{
                    height: '100%',
                    display: 'flex',
                    flexDirection: 'column',
                    transition: 'transform 0.2s, box-shadow 0.2s',
                    cursor: 'pointer',
                    '&:hover': {
                        transform: 'translateY(-4px)',
                        boxShadow: 3
                    }
                }}
                onClick={handleCardClick}
            >
                <Box sx={{ position: 'relative' }}>
                    <CardMedia
                        component="img"
                        sx={{
                            height: isMobile ? 120 : 200,
                            objectFit: 'cover'
                        }}
                        image={getImageUrl()}
                        alt={listing.title}
                    />

                    {/* Контейнер для бирок (магазин и скидка) */}
                    <Box>
                        {/* Бирка "в магазин" - используем нашу функцию определения товаров из витрины */}
                        {isStoreItem() && (
                            <Box
                                sx={{
                                    position: 'absolute',
                                    top: 10,
                                    left: 10,
                                    bgcolor: 'primary.main',
                                    color: 'white',
                                    borderRadius: '4px',
                                    padding: '4px 8px',
                                    fontWeight: 'bold',
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 0.5,
                                    fontSize: '0.75rem',
                                    cursor: 'pointer',
                                    zIndex: 5
                                }}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    if (e.nativeEvent) {
                                        e.nativeEvent.stopImmediatePropagation();
                                    }
                                    window.location.href = `/shop/${getStoreId()}`;
                                    return false;
                                }}
                            >
                                <Store size={16} />
                                {t('listings.details.goToStore')}
                            </Box>
                        )}

                        {/* Показатель скидки */}
                        {discountData && discountData.percent > 0 && (
                            <Box
                                sx={{
                                    position: 'absolute',
                                    top: 10,
                                    right: 10,
                                    bgcolor: 'error.main',
                                    color: 'white',
                                    borderRadius: '4px',
                                    padding: '4px 8px',
                                    fontWeight: 'bold',
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 0.5,
                                    zIndex: 4 // Высокий z-index, но чуть ниже чем у бирки "в магазин"
                                }}
                            >
                                <Typography variant="body2" fontWeight="bold">
                                    -{discountData.percent}%
                                </Typography>
                            </Box>
                        )}
                    </Box>
                    
                    {/* Тег статуса */}
                    {showStatus && listing.status && (
                        <Chip
                            label={
                                listing.status === 'active'
                                    ? t('listings.status.active')
                                    : t('listings.status.inactive')
                            }
                            size="small"
                            color={listing.status === 'active' ? 'success' : 'default'}
                            sx={{
                                position: 'absolute',
                                top: 10,
                                left: 10
                            }}
                        />
                    )}
                    
                    {/* Индикатор количества фото */}
                    {listing.images && listing.images.length > 1 && (
                        <Box
                            sx={{
                                position: 'absolute',
                                bottom: 10,
                                right: 10,
                                bgcolor: 'rgba(0,0,0,0.6)',
                                color: 'white',
                                borderRadius: '4px',
                                padding: '2px 6px',
                                display: 'flex',
                                alignItems: 'center',
                                gap: 0.5
                            }}
                        >
                            <Camera size={16} />
                            <Typography variant="body2" fontWeight="medium">
                                {listing.images.length}
                            </Typography>
                        </Box>
                    )}
                </Box>
                
                <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
                    <Typography
                        variant="h6"
                        component="h2"
                        sx={{
                            mb: 1,
                            fontSize: isMobile ? '0.9rem' : '1rem',
                            fontWeight: 'medium',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            display: '-webkit-box',
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: 'vertical'
                        }}
                    >
                        {listing.title}
                    </Typography>
                    
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                        <Typography
                            variant="h6"
                            component="div"
                            color="primary"
                            fontWeight="bold"
                            sx={{ fontSize: isMobile ? '1rem' : '1.25rem' }}
                        >
                            {new Intl.NumberFormat(i18n.language, {
                                style: 'currency',
                                currency: 'RSD',
                                maximumFractionDigits: 0
                            }).format(listing.price)}
                        </Typography>
                        
                        {discountData && (
                            <Typography
                                variant="body2"
                                color="text.secondary"
                                sx={{
                                    ml: 1,
                                    textDecoration: 'line-through',
                                    fontSize: isMobile ? '0.75rem' : '0.875rem'
                                }}
                            >
                                {new Intl.NumberFormat(i18n.language, {
                                    style: 'currency',
                                    currency: 'RSD',
                                    maximumFractionDigits: 0
                                }).format(discountData.oldPrice)}
                            </Typography>
                        )}
                        
                        {discountData?.hasPriceHistory && (
                            <Button
                                size="small"
                                variant="text"
                                onClick={handleOpenPriceHistory}
                                sx={{ minWidth: 'auto', ml: 'auto', px: 1 }}
                            >
                                <AccessTime size={16} />
                            </Button>
                        )}
                    </Box>
                    
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            mb: 0.5,
                            fontSize: isMobile ? '0.75rem' : '0.875rem',
                            color: 'text.secondary'
                        }}
                    >
                        <LocationIcon size={isMobile ? 14 : 16} style={{ marginRight: 4 }} />
                        <Typography
                            variant="body2"
                            sx={{
                                overflow: 'hidden',
                                textOverflow: 'ellipsis',
                                whiteSpace: 'nowrap'
                            }}
                        >
                            {listing.city || listing.location || t('listings.locationNotSpecified')}
                        </Typography>
                    </Box>
                    
                    {/* Перенесли бирку "в магазин" поверх изображения */}
                    
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            mb: 0.5,
                            fontSize: isMobile ? '0.75rem' : '0.875rem',
                            color: 'text.secondary',
                            mt: 'auto'
                        }}
                    >
                        <AccessTime size={isMobile ? 14 : 16} style={{ marginRight: 4 }} />
                        <Typography variant="body2">
                            {formatDate(listing.created_at)}
                        </Typography>
                        
                        {listing.views_count !== undefined && (
                            <Box
                                sx={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    ml: 'auto',
                                    fontSize: isMobile ? '0.75rem' : '0.875rem',
                                    color: 'text.secondary'
                                }}
                            >
                                <Eye size={isMobile ? 14 : 16} style={{ marginRight: 4 }} />
                                <Typography variant="body2">{listing.views_count}</Typography>
                            </Box>
                        )}
                    </Box>
                    
                    {listing.rating !== undefined && listing.reviews_count !== undefined && (
                        <Box
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                mt: 1
                            }}
                        >
                            <Rating
                                value={listing.rating}
                                readOnly
                                precision={0.5}
                                size={isMobile ? 'small' : 'medium'}
                            />
                            <Typography variant="body2" color="text.secondary" sx={{ ml: 1 }}>
                                ({listing.reviews_count})
                            </Typography>
                        </Box>
                    )}
                </CardContent>
            </Card>
            
            {/* Модальное окно истории цен */}
            <ModalWrapper
                open={isPriceHistoryOpen}
                onClose={handleClosePriceHistory}
                title={t('listings.priceHistory.title')}
            >
                <PriceHistoryChart listingId={listing.id as number} />
            </ModalWrapper>
        </>
    );
};

export default ListingCard;