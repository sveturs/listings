// frontend/hostel-frontend/src/components/marketplace/ListingCard.tsx
import React, { useState, useEffect, MouseEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { MapPin as LocationIcon, Clock as AccessTime, Camera, Store, Eye } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import axios from '../../api/axios';
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

interface ImageObject {
    id?: number | string;
    file_path?: string;
    public_url?: string;
    is_main?: boolean;
    storage_type?: string;
    [key: string]: any;
}

interface Attribute {
    attribute_id: number | string;
    attribute_name: string;
    display_name: string;
    display_value: string;
    numeric_value?: number;
    translations?: {
        [language: string]: string;
    };
    [key: string]: any;
}

interface Promotion {
    [type: string]: any;
}

interface DiscountInfo {
    previous_price?: number;
    has_price_history?: boolean;
    [key: string]: any;
}

interface ListingMetadata {
    discount?: DiscountInfo;
    promotions?: {
        [type: string]: Promotion;
    };
    [key: string]: any;
}

export interface Listing {
    id: number | string;
    title: string;
    price: number;
    old_price?: number;
    city?: string;
    location?: string;
    status?: string;
    created_at?: string;
    images?: (string | ImageObject)[];
    attributes?: Attribute[];
    rating?: number;
    reviews_count?: number;
    storefront_id?: number | string;
    views_count?: number;
    has_discount?: boolean;
    original_language?: string;
    translations?: {
        [language: string]: {
            [field: string]: string;
        };
    };
    metadata?: ListingMetadata;
    [key: string]: any;
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

    // Логирование для отладки
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

    const getDiscountInfo = (): DiscountData | null => {
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

    const renderDiscountBadge = (): React.ReactNode => {
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
        const fetchStoreName = async (): Promise<void> => {
            if (listing.storefront_id) {
                try {
                    const response = await axios.get(`/api/v1/public/storefronts/${listing.storefront_id}`);
                    if (response.data?.data?.name) {
                        setStoreName(response.data.data.name);
                    }
                } catch (err) {
                    console.error(t('listings.details.errors.loadFailed'), err);
                }
            }
        };
        fetchStoreName();
    }, [listing.storefront_id, t]);

    const getLocalizedText = (field: string): string => {
        if (!listing || !field) return '';
        if (i18n.language === listing.original_language) {
            return listing[field];
        }
        const translation = listing.translations?.[i18n.language]?.[field];
        return translation || listing[field];
    };

    const getTranslatedName = (attr: Attribute): string => {
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
        const attributeTranslations: Record<string, Record<string, string>> = {
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

    const formatPrice = (price: number | undefined): string => {
        return new Intl.NumberFormat('sr-RS', {
            style: 'currency',
            currency: 'RSD',
            maximumFractionDigits: 0
        }).format(price || 0);
    };

    const formatDate = (dateString?: string): string => {
        if (!dateString) return t('listings.details.seller.unknownDate');
        const date = new Intl.DateTimeFormat(i18n.language, {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        }).format(new Date(dateString));
        return t('listings.details.seller.memberSince', { date });
    };

    // Оптимизированная функция получения URL главного изображения
    const getMainImageUrl = (): string => {
        if (!listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
            console.log('No images found, using placeholder');
            return '/placeholder.jpg';
        }

        // Находим главное изображение или используем первое в списке
        let mainImage = listing.images.find(img => img && (img as ImageObject).is_main === true) || listing.images[0];

        // Используем переменную окружения из window.ENV вместо process.env
        const baseUrl = (window as any).ENV?.REACT_APP_MINIO_URL || (window as any).ENV?.REACT_APP_BACKEND_URL || '';
        console.log('Using baseUrl from env:', baseUrl);

        // 1. Строковые пути (для обратной совместимости)
        if (typeof mainImage === 'string') {
            console.log('Processing string image path:', mainImage);

            // Относительный путь MinIO
            if (mainImage.startsWith('/listings/')) {
                const url = `${baseUrl}${mainImage}`;
                console.log('Using MinIO relative path:', url);
                return url;
            }

            // ID/filename.jpg (прямой путь MinIO)
            if (mainImage.match(/^\d+\/[^\/]+$/)) {
                const url = `${baseUrl}/listings/${mainImage}`;
                console.log('Using direct MinIO path pattern:', url);
                return url;
            }

            // Локальное хранилище (обратная совместимость)
            const url = `${baseUrl}/uploads/${mainImage}`;
            console.log('Using local storage path:', url);
            return url;
        }

        // 2. Объекты с информацией о файле
        if (typeof mainImage === 'object' && mainImage !== null) {
            const imageObj = mainImage as ImageObject;
            console.log('Processing image object:', imageObj);

            // Приоритет 1: Используем PublicURL если он доступен
            if (imageObj.public_url && typeof imageObj.public_url === 'string' && imageObj.public_url.trim() !== '') {
                const publicUrl = imageObj.public_url;
                console.log('Found public_url string:', publicUrl);

                // Абсолютный URL
                if (publicUrl.startsWith('http')) {
                    console.log('Using absolute URL:', publicUrl);
                    return publicUrl;
                }
                // Относительный URL с /listings/
                else if (publicUrl.startsWith('/listings/')) {
                    const url = `${baseUrl}${publicUrl}`;
                    console.log('Using public_url with listings path:', url);
                    return url;
                }
                // Другой относительный URL
                else {
                    const url = `${baseUrl}${publicUrl}`;
                    console.log('Using general relative public_url:', url);
                    return url;
                }
            }

            // Приоритет 2: Формируем URL на основе типа хранилища и пути к файлу
            if (imageObj.file_path) {
                if (imageObj.storage_type === 'minio' || imageObj.file_path.includes('listings/')) {
                    // Учитываем возможность наличия префикса listings/ в пути
                    const filePath = imageObj.file_path.includes('listings/')
                        ? imageObj.file_path.replace('listings/', '')
                        : imageObj.file_path;

                    const url = `${baseUrl}/listings/${filePath}`;
                    console.log('Constructed MinIO URL from path:', url);
                    return url;
                }

                // Локальное хранилище
                const url = `${baseUrl}/uploads/${imageObj.file_path}`;
                console.log('Using local storage path from object:', url);
                return url;
            }
        }

        console.log('Could not determine image URL, using placeholder');
        return '/placeholder.jpg';
    };

    const handleCardClick = (e: React.MouseEvent<HTMLDivElement>): void => {
        if ((e.target as HTMLElement).closest('[data-shop-button="true"]') || (e.target as HTMLElement).closest('#detailsButton')) {
            return;
        }
        if (onClick) {
            onClick(listing);
        } else {
            navigate(`/marketplace/listings/${listing.id}`);
        }
    };

    const handleShopButtonClick = (e: React.MouseEvent<HTMLDivElement>): void => {
        e.preventDefault();
        e.stopPropagation();
        navigate(`/shop/${listing.storefront_id}`);
    };

    const handleDetailsButtonClick = (e: React.MouseEvent<HTMLButtonElement>): void => {
        e.preventDefault();
        e.stopPropagation();
        navigate(`/marketplace/listings/${listing.id}`);
    };

    const getDisplayLocation = (): string => {
        if (listing.city) return listing.city;

        if (listing.location) {
            const locationParts = listing.location.split(',');
            return locationParts.length > 0 ? locationParts[0].trim() : listing.location;
        }

        return t('listings.location.unknown');
    };

    const renderPriceSection = (): React.ReactNode => {
        const discount = getDiscountInfo();
        return (
            <Box>
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
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
                    
                    {/* View count indicator - always show even if zero */}
                    <Box sx={{ display: 'flex', alignItems: 'center', ml: 1, mt: 0.5 }}>
                        <Eye size={14} color="#666" />
                        <Typography variant="caption" color="text.secondary" sx={{ ml: 0.25 }}>
                            {listing.views_count || 0}
                        </Typography>
                    </Box>
                </Box>
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

    const renderAttributes = (): React.ReactNode => {
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
                    {getLocalizedText('title') || t('listings.create.name')}
                </Typography>

                {listing.rating && listing.rating > 0 && (
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
                <ModalWrapper
                    aria-labelledby="price-history-title"
                    open={isPriceHistoryOpen}
                    onClose={() => setIsPriceHistoryOpen(false)}
                    isMobile={isMobile}
                >
                    <Typography id="price-history-title" variant="h6" component="h2" gutterBottom>
                        {t('listings.details.priceHistory.title')}
                    </Typography>
                    <PriceHistoryChart listingId={listing.id} />
                    <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                        <Button onClick={() => setIsPriceHistoryOpen(false)}>
                            {t('gallery.close')}
                        </Button>
                    </Box>
                </ModalWrapper>
            )}
        </Card>
    );
};

export default ListingCard;