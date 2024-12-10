import React from 'react';
import {
    Card,
    CardMedia,
    CardContent,
    Typography,
    Box,
    Chip,
    Button
} from '@mui/material';
import {
    LocationOn as LocationIcon,
    AccessTime,
    PhotoCamera
} from '@mui/icons-material';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const ListingCard = ({ listing }) => {
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
            <Box sx={{ position: 'relative', pt: '75%' }}>
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
                    image={listing.images?.[0] ? 
                        `${process.env.REACT_APP_BACKEND_URL}/uploads/${listing.images[0].file_path}` : 
                        'https://via.placeholder.com/300x200'}
                    alt={listing.title || 'Изображение отсутствует'}
                />
                {listing.images?.length > 1 && (
                    <Chip
                        icon={<PhotoCamera />}
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

            <CardContent sx={{ flexGrow: 1 }}>
                <Typography variant="h6" noWrap>
                    {listing.title || 'Без названия'}
                </Typography>

                <Typography variant="h5" color="primary" sx={{ mt: 1 }}>
                    {formatPrice(listing.price)}
                </Typography>

                <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                    <LocationIcon sx={{ fontSize: 18, mr: 0.5 }} />
                    <Typography variant="body2" noWrap>
                        {listing.city || 'Местоположение не указано'}
                    </Typography>
                </Box>

                <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                    <AccessTime sx={{ fontSize: 18, mr: 0.5 }} />
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
            </CardContent>
        </Card>
    );
};

export default ListingCard;