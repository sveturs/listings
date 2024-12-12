// src/components/marketplace/ListingCard.js
import React from 'react';
import {
    Card,
    CardContent,
    CardMedia,
    Typography,
    Box,
    Chip,
    Button
} from '@mui/material';
import { MapPin as LocationIcon, Clock as AccessTime, Camera } from 'lucide-react';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000';

const ListingCard = ({ listing, isMobile }) => {
    const formatPrice = (price) => {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0
        }).format(price || 0);
    };

    const formatDate = (dateString) => {
        if (!dateString) return '';
        return new Date(dateString).toLocaleDateString('ru-RU', {
            day: 'numeric',
            month: 'long',
            year: 'numeric'
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
            '&:hover': {
                transform: 'translateY(-4px)',
                boxShadow: 3,
                transition: 'all 0.2s ease-in-out'
            }
        }}>
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
                    alt={listing.title || 'Изображение отсутствует'}
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
                    {listing.title || 'Без названия'}
                </Typography>

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
                            variant="contained"
                            fullWidth
                            sx={{ mt: 2 }}
                        >
                            Подробнее
                        </Button>
                    </>
                )}
            </CardContent>
        </Card>
    );
};

export default ListingCard;