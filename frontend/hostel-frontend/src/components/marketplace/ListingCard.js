// src/components/marketplace/ListingCard.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { MapPin as LocationIcon, Clock as AccessTime, Camera, Store } from 'lucide-react';
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
    Tooltip
} from '@mui/material';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000';

const ListingCard = ({ listing, isMobile }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [storeName, setStoreName] = useState('Магазин');
    
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
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString(i18n.language, {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    const getMainImageUrl = () => {
        if (!listing.images || listing.images.length === 0) {
            return '/placeholder.jpg';
        }

        const mainImage = listing.images.find(img => img.is_main) || listing.images[0];
        if (!mainImage || !mainImage.file_path) {
            return '/placeholder.jpg';
        }

        return `${BACKEND_URL}/uploads/${mainImage.file_path}`;
    };

    return (
        <Card sx={{
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
            position: 'relative', // Важно для абсолютного позиционирования кнопки
            '&:hover': {
                transform: 'translateY(-4px)',
                boxShadow: 3,
                transition: 'all 0.2s ease-in-out'
            }
        }}>
            {/* Кнопка "Магазин тут" с всплывающей подсказкой */}
            {listing.storefront_id && (
                <Tooltip 
                    title={storeName}
                    placement="top"
                    arrow
                >
                    <Box
                        sx={{
                            position: 'absolute',
                            top: 10,
                            right: 10,
                            zIndex: 10, // Высокий z-index для отображения поверх всех элементов
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
                            pointerEvents: 'auto'
                        }}
                        onClick={(e) => {
                            e.preventDefault();
                            e.stopPropagation();
                            window.location.href = `/shop/${listing.storefront_id}`;
                        }}
                        data-shop-button="true"
                    >
                        <Store size={14} />
                        в магазин
                    </Box>
                </Tooltip>
            )}
            
            {/* Остальной код без изменений */}
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
                        label={`${listing.images.length} фото`}
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

                {!isMobile && (
                    <>
                        <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                            <LocationIcon size={18} style={{ marginRight: 4 }} />
                            <Typography variant="body2" noWrap>
                                {listing.city || 'Местоположение не указано'}
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
                        >
                            {t('listings.details.moreDetails')}
                        </Button>
                    </>
                )}
            </CardContent>
        </Card>
    );
};

export default ListingCard;