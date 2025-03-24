import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import axios from '../../api/axios';
import {
    Box,
    Typography,
    Grid,
    Skeleton,
    useTheme,
    useMediaQuery,
    Card,
    CardContent,
    CardMedia,
    Chip,
} from '@mui/material';
import { MapPin, Percent } from 'lucide-react';

const SimilarListings = ({ listingId }) => {
    const { t } = useTranslation('marketplace');
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const navigate = useNavigate();
    
    const [similarListings, setSimilarListings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    
    useEffect(() => {
        const fetchSimilarListings = async () => {
            try {
                setLoading(true);
                const response = await axios.get(`/api/v1/marketplace/listings/${listingId}/similar?limit=8`);
                if (response.data && response.data.data) {
                    setSimilarListings(response.data.data);
                }
            } catch (err) {
                console.error('Error fetching similar listings:', err);
                setError(t('listings.similar.error'));
            } finally {
                setLoading(false);
            }
        };
        
        if (listingId) {
            fetchSimilarListings();
        }
    }, [listingId, t]);
    
    const formatPrice = (price) => {
        return new Intl.NumberFormat('sr-RS', {
            style: 'currency',
            currency: 'RSD',
            maximumFractionDigits: 0
        }).format(price || 0);
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
        
        if (mainImage && mainImage.file_path) {
            return `${baseUrl}/uploads/${mainImage.file_path}`;
        }
        
        return '/placeholder.jpg';
    };
    
    const renderDiscountBadge = (listing) => {
        if (!listing || !listing.metadata || !listing.metadata.discount) return null;
        
        const discount = listing.metadata.discount;
        
        return (
            <Chip
                icon={<Percent size={14} />}
                label={`-${discount.discount_percent}%`}
                color="warning"
                size="small"
                sx={{
                    position: 'absolute',
                    top: 8,
                    left: 8,
                    zIndex: 2,
                    fontWeight: 'bold',
                    fontSize: '0.75rem'
                }}
            />
        );
    };
    
    const handleCardClick = (id) => {
        navigate(`/marketplace/listings/${id}`);
    };
    
    if (loading) {
        return (
            <Box sx={{ mt: 6, mb: 4 }}>
                <Typography variant="h5" gutterBottom>
                    {t('listings.similar.title', { defaultValue: 'Похожие объявления' })}
                </Typography>
                <Grid container spacing={2}>
                    {[...Array(4)].map((_, index) => (
                        <Grid item xs={6} sm={4} md={3} key={index}>
                            <Skeleton variant="rectangular" height={200} />
                            <Skeleton width="60%" sx={{ mt: 1 }} />
                            <Skeleton width="40%" />
                        </Grid>
                    ))}
                </Grid>
            </Box>
        );
    }
    
    if (error) {
        return null; // Скрываем секцию при ошибке
    }
    
    if (!similarListings || similarListings.length === 0) {
        return null; // Скрываем секцию, если нет похожих объявлений
    }
    
    return (
        <Box sx={{ mt: 6, mb: 4 }}>
            <Typography variant="h5" gutterBottom>
                {t('listings.similar.title', { defaultValue: 'Похожие объявления' })}
            </Typography>
            <Grid container spacing={2}>
                {similarListings.map((listing) => (
                    <Grid item xs={6} sm={4} md={3} key={listing.id}>
                        <Card 
                            sx={{
                                height: '100%',
                                display: 'flex',
                                flexDirection: 'column',
                                position: 'relative',
                                cursor: 'pointer',
                                '&:hover': {
                                    transform: 'translateY(-4px)',
                                    boxShadow: 3,
                                    transition: 'all 0.2s ease-in-out'
                                }
                            }}
                            onClick={() => handleCardClick(listing.id)}
                        >
                            {renderDiscountBadge(listing)}
                            <Box sx={{ position: 'relative', pt: '75%' /* 4:3 Aspect Ratio */ }}>
                                <CardMedia
                                    component="img"
                                    image={getImageUrl(listing)}
                                    alt={listing.title}
                                    sx={{
                                        position: 'absolute',
                                        top: 0,
                                        left: 0,
                                        width: '100%',
                                        height: '100%',
                                        objectFit: 'contain',
                                        backgroundColor: '#f5f5f5'
                                    }}
                                />
                            </Box>
                            <CardContent sx={{ flexGrow: 1, p: 1.5, '&:last-child': { pb: 1.5 } }}>
                                <Typography 
                                    variant="subtitle2" 
                                    sx={{ 
                                        fontWeight: 'medium',
                                        overflow: 'hidden',
                                        textOverflow: 'ellipsis',
                                        display: '-webkit-box',
                                        WebkitLineClamp: 2,
                                        WebkitBoxOrient: 'vertical',
                                        height: '2.5em',
                                        lineHeight: 1.2
                                    }}
                                >
                                    {listing.title}
                                </Typography>
                                
                                <Box sx={{ mt: 1 }}>
                                    <Typography 
                                        variant="body2" 
                                        color="primary.main" 
                                        sx={{ fontWeight: 'bold' }}
                                    >
                                        {formatPrice(listing.price)}
                                    </Typography>
                                    
                                    {listing.metadata && listing.metadata.discount && (
                                        <Typography
                                            variant="caption"
                                            color="text.secondary"
                                            sx={{ textDecoration: 'line-through', ml: 1 }}
                                        >
                                            {formatPrice(listing.metadata.discount.previous_price)}
                                        </Typography>
                                    )}
                                </Box>
                                
                                {!isMobile && listing.city && (
                                    <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                                        <MapPin size={14} style={{ marginRight: 4 }} />
                                        <Typography 
                                            variant="caption" 
                                            color="text.secondary"
                                            noWrap
                                        >
                                            {listing.city}
                                        </Typography>
                                    </Box>
                                )}
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>
        </Box>
    );
};

export default SimilarListings;