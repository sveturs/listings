// frontend/hostel-frontend/src/components/marketplace/ListingCard.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { MapPin as LocationIcon, Clock as AccessTime, Camera, Store } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
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
    Modal
} from '@mui/material';
import { Percent } from 'lucide-react';
import PriceHistoryChart from './PriceHistoryChart';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000';

const ListingCard = ({ listing, isMobile, onClick, showStatus = false }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [storeName, setStoreName] = useState(t('listings.details.seller.title')); // Заменено 'Магазин' на перевод
    const navigate = useNavigate();
    const [isPriceHistoryOpen, setIsPriceHistoryOpen] = useState(false);

    // Логирование для отладки (оставляем без изменений)
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
        } else if (listing && listing.has_discount) {
            console.log('Has discount flag is true, but no old_price or metadata');
        } else {
            console.log('No discount data found for listing:', listing.id);
        }
    }, [listing]);

    const getDiscountInfo = () => {
        if (listing.metadata && listing.metadata.discount) {
            const previousPrice = Number(listing.metadata.discount.previous_price);
            const calculatedPercent = Math.round((1 - listing.price / previousPrice) * 100);
            return {
                percent: calculatedPercent,
                oldPrice: previousPrice,
                hasPriceHistory: listing.metadata.discount.has_price_history || false
            };
        }
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

    useEffect(() => {
        const fetchStoreName = async () => {
            if (listing.storefront_id) {
                try {
                    const response = await axios.get(`/api/v1/public/storefronts/${listing.storefront_id}`);
                    if (response.data?.data?.name) {
                        setStoreName(response.data.data.name);
                    }
                } catch (err) {
                    console.error(t('listings.details.errors.loadFailed'), err); // Заменено сообщение об ошибке
                }
            }
        };
        fetchStoreName();
    }, [listing.storefront_id, t]);

    const getLocalizedText = (field) => {
        if (!listing || !field) return '';
        if (i18n.language === listing.original_language) {
            return listing[field];
        }
        const translation = listing.translations?.[i18n.language]?.[field];
        return translation || listing[field];
    };

    const getTranslatedName = (attr) => {
        if (!attr) return '';

        // Если текущий язык русский, возвращаем отображаемое имя без перевода
        if (i18n.language === 'ru') {
            return attr.display_name;
        }

        // Проверяем наличие прямого перевода имени атрибута
        if (attr.translations && attr.translations[i18n.language]) {
            return attr.translations[i18n.language];
        }

        // Стандартные переводы для атрибутов
        const attributeTranslations = {
            // Автомобили
            'make': { 'en': 'Make', 'sr': 'Marka' },
            'model': { 'en': 'Model', 'sr': 'Model' },
            'year': { 'en': 'Year', 'sr': 'Godina proizvodnje' },
            'mileage': { 'en': 'Mileage', 'sr': 'Kilometraža' },

            'engine_capacity': { 'en': 'Engine capacity', 'sr': 'Zapremina motora' },
            'fuel_type': { 'en': 'Fuel type', 'sr': 'Vrsta goriva' },
            'transmission': { 'en': 'Transmission', 'sr': 'Menjač' },
            'body_type': { 'en': 'Body type', 'sr': 'Tip karoserije' },
            'color': { 'en': 'Color', 'sr': 'Boja' },
            'power': { 'en': 'Power', 'sr': 'Snaga' },
            'drive_type': { 'en': 'Drive type', 'sr': 'Pogon' },
            'number_of_doors': { 'en': 'Number of doors', 'sr': 'Broj vrata' },
            'number_of_seats': { 'en': 'Number of seats', 'sr': 'Broj sedišta' },

            // Недвижимость
            'property_type': { 'en': 'Property type', 'sr': 'Tip nekretnine' },
            'rooms': { 'en': 'Rooms', 'sr': 'Broj soba' },
            'floor': { 'en': 'Floor', 'sr': 'Sprat' },
            'total_floors': { 'en': 'Total floors', 'sr': 'Ukupno spratova' },
            'area': { 'en': 'Area', 'sr': 'Površina' },
            'land_area': { 'en': 'Land area', 'sr': 'Površina zemljišta' },
            'building_type': { 'en': 'Building type', 'sr': 'Tip zgrade' },
            'has_balcony': { 'en': 'Balcony', 'sr': 'Balkon' },
            'has_elevator': { 'en': 'Elevator', 'sr': 'Lift' },
            'has_parking': { 'en': 'Parking', 'sr': 'Parking' },

            // Электроника
            'brand': { 'en': 'Brand', 'sr': 'Brend' },
            'model_phone': { 'en': 'Model', 'sr': 'Model' },
            'memory': { 'en': 'Memory', 'sr': 'Memorija' },
            'ram': { 'en': 'RAM', 'sr': 'RAM' },
            'os': { 'en': 'Operating system', 'sr': 'Operativni sistem' },
            'screen_size': { 'en': 'Screen size', 'sr': 'Veličina ekrana' },
            'camera': { 'en': 'Camera', 'sr': 'Kamera' },
            'has_5g': { 'en': '5G', 'sr': '5G' },

            // Компьютеры
            'pc_brand': { 'en': 'Brand', 'sr': 'Brend' },
            'pc_type': { 'en': 'Type', 'sr': 'Tip' },
            'cpu': { 'en': 'Processor', 'sr': 'Procesor' },
            'gpu': { 'en': 'Graphics card', 'sr': 'Grafička kartica' },
            'ram_pc': { 'en': 'RAM', 'sr': 'RAM' },
            'storage_type': { 'en': 'Storage type', 'sr': 'Tip skladišta' },
            'storage_capacity': { 'en': 'Storage capacity', 'sr': 'Kapacitet skladišta' },
            'os_pc': { 'en': 'Operating system', 'sr': 'Operativni sistem' }
        };

        // Проверяем наличие стандартного перевода для атрибута
        if (attributeTranslations[attr.attribute_name] &&
            attributeTranslations[attr.attribute_name][i18n.language]) {
            return attributeTranslations[attr.attribute_name][i18n.language];
        }

        // Если перевод не найден, возвращаем display_name
        return attr.display_name || attr.attribute_name;
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
        const date = new Intl.DateTimeFormat(i18n.language, {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        }).format(new Date(dateString));
        return t('listings.details.seller.memberSince', { date });
    };

const getMainImageUrl = () => {
        if (!listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
            return '/placeholder.jpg';
        }

        let mainImage = listing.images.find(img => img && img.is_main === true) || listing.images[0];

        if (mainImage && typeof mainImage === 'object') {
            // Если есть публичный URL, используем его напрямую
            if (mainImage.public_url) {
                return mainImage.public_url;
            }

            // Для MinIO формируем URL через специальный путь в nginx
            if (mainImage.storage_type === 'minio') {
                return `${process.env.REACT_APP_BACKEND_URL}/listings/${mainImage.file_path.split('/').pop()}`;
            }

            // Для локального хранилища используем старый формат
            return `${process.env.REACT_APP_BACKEND_URL}/uploads/${mainImage.file_path}`;
        }

        // Если изображение просто строка (обратная совместимость)
        if (mainImage && typeof mainImage === 'string') {
            return `${process.env.REACT_APP_BACKEND_URL}/uploads/${mainImage}`;
        }

        return '/placeholder.jpg';
    }





    const handleCardClick = (e) => {
        if (e.target.closest('[data-shop-button="true"]') || e.target.closest('#detailsButton')) {
            return;
        }
        if (onClick) {
            onClick(listing);
        } else {
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
        if (listing.city) return listing.city;

        if (listing.location) {
            const locationParts = listing.location.split(',');
            return locationParts.length > 0 ? locationParts[0].trim() : listing.location;
        }

        return t('listings.location.unknown');
    };

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
                {discount && (
                    <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ textDecoration: 'line-through', mt: 0.5 }}
                    >
                        {formatPrice(discount.oldPrice)}
                    </Typography>
                )}
            </Box>
        );
    };

    const renderAttributes = () => {
        if (!listing.attributes || listing.attributes.length === 0) {
            return null;
        }
        // Получаем важные атрибуты недвижимости
        const realEstateAttrs = listing.attributes.filter(attr =>
            ['rooms', 'floor', 'total_floors', 'area', 'land_area', 'property_type'].includes(attr.attribute_name)
        );
        if (realEstateAttrs.length === 0) {
            return null;
        }
        return (
            <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                {realEstateAttrs.map(attr => {
                    // Форматируем отображение атрибута
                    let displayValue = attr.display_value || '';

                    // Добавляем единицы измерения, если их нет
                    if (attr.attribute_name === 'area' && !displayValue.includes('м²')) {
                        displayValue = `${displayValue} м²`;
                    } else if (attr.attribute_name === 'land_area' && !displayValue.includes('сот')) {
                        displayValue = `${displayValue} сот`;
                    } else if (attr.attribute_name === 'rooms' && attr.numeric_value) {
                        // Форматируем комнаты правильно
                        const numRooms = attr.numeric_value;
                        const roomWord = numRooms === 1 ? 'комната' :
                            (numRooms >= 2 && numRooms <= 4) ? 'комнаты' : 'комнат';
                        displayValue = `${numRooms} ${roomWord}`;
                    }

                    return (
                        <Chip
                            key={attr.attribute_id}
                            label={`${getTranslatedName(attr)}: ${displayValue}`}
                            size="small"
                            variant="outlined"
                            sx={{ fontSize: '0.75rem', height: 24 }}
                        />
                    );
                })}
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
                        zIndex: 1200,
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

                    }}
                    onClick={handleShopButtonClick}
                    data-shop-button="true"
                >
                    <Store size={14} />
                    {t('listings.details.goToStore')}
                </Box>
            )}

            {/* Бейдж скидки */}
            {renderDiscountBadge()}

            {/* Контейнер для изображения */}
            <Box sx={{ position: 'relative' }}>
                <CardMedia
                    component="img"
                    sx={{
                        width: '100%',
                        height: isMobile ? 200 : 250, // Фиксированная высота для контроля размера
                        objectFit: 'contain', // Изображение полностью помещается
                        backgroundColor: '#f5f5f5', // Фон для пустых областей
                    }}
                    image={getMainImageUrl()}
                    alt={getLocalizedText('title') || t('listings.details.title.noImages')}
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
                    {getLocalizedText('title') || t('listings.create.name')} {/* Заменено 'Без названия' на 'Наименование товара' */}
                </Typography>

                {listing.rating > 0 && (
                    <Stack direction="row" spacing={0.5} alignItems="center" sx={{ mt: 1 }}>
                        <Rating value={listing.rating} readOnly size="small" precision={0.1} />
                        <Typography variant="body2" color="text.secondary">
                            ({listing.reviews_count})
                        </Typography>
                    </Stack>
                )}

                {renderPriceSection()}

                {/* Отображение атрибутов недвижимости */}
                {renderAttributes()}

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

                        {showStatus && (
                            <Box sx={{ mt: 1 }}>
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

                                {listing.metadata && listing.metadata.promotions && Object.keys(listing.metadata.promotions).length > 0 && (
                                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 0.5 }}>
                                        {Object.entries(listing.metadata.promotions).map(([type, data]) => (
                                            <Chip
                                                key={type}
                                                label={t(`listings.promotions.${type}`, { defaultValue: type })}
                                                color="primary"
                                                size="small"
                                                variant="outlined"
                                                sx={{ fontSize: '0.7rem', height: 20 }}
                                            />
                                        ))}
                                    </Box>
                                )}
                            </Box>
                        )}

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
                            {t('listings.details.priceHistory.title')}
                        </Typography>
                        <PriceHistoryChart listingId={listing.id} />
                        <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                            <Button onClick={() => setIsPriceHistoryOpen(false)}>
                                {t('gallery.close')}
                            </Button>
                        </Box>
                    </Box>
                </Modal>
            )}
        </Card>
    );
};

export default ListingCard;