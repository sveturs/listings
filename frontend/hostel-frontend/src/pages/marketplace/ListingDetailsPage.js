//frontend/hostel-frontend/src/pages/ListingDetailsPage.js
import ChatButton from '../../components/marketplace/chat/ChatButton';
import React, { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import ReviewsSection from '../../components/reviews/ReviewsSection';
import { useAuth } from '../../contexts/AuthContext';
import MiniMap from '../../components/maps/MiniMap';
import { PencilLine, Trash2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import ShareButton from '../../components/marketplace/ShareButton';

import { GoogleMap, Marker } from '@react-google-maps/api';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import CallButton from '../../components/marketplace/CallButton';

import {
    Container,
    Modal,
    Paper,
    Grid,
    Box,
    Typography,
    Button,
    Card,
    CardContent,
    Skeleton,
    Stack,
    Avatar,
    IconButton,
    useTheme,
    useMediaQuery,
    ImageList,
    ImageListItem
} from '@mui/material';
import {
    MapPin,
    Calendar,
    Heart,
    Share2,
    Phone,
    MessageCircle,
    ChevronLeft,
    Maximize2,
    ChevronRight
} from 'lucide-react';
import axios from '../../api/axios';


const ListingDetailsPage = () => {
    const navigate = useNavigate();

    const [isFavorite, setIsFavorite] = useState(false);
    const [isMapExpanded, setIsMapExpanded] = useState(false);
    const { id } = useParams();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const reviewsRef = useRef(null);

    const [listing, setListing] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [currentImageIndex, setCurrentImageIndex] = useState(0);
    const [reviewsCount, setReviewsCount] = useState(0);
    const { user, login } = useAuth();
    const [categoryPath, setCategoryPath] = useState([]);
    const [showPhoneNumber, setShowPhoneNumber] = useState(false);
    const [showPhone, setShowPhone] = useState(false);
    useEffect(() => {
        const fetchListing = async () => {
            try {
                setLoading(true);
                const [listingResponse, favoritesResponse] = await Promise.all([
                    axios.get(`/api/v1/marketplace/listings/${id}`),
                    axios.get('/api/v1/marketplace/favorites')
                ]);

                const isFavorite = favoritesResponse.data?.data?.some?.(
                    item => item.id === Number(id)
                ) || false;

                setListing({
                    ...listingResponse.data.data,
                    is_favorite: isFavorite
                });

                if (listingResponse.data.data.category_path) {
                    const path = listingResponse.data.data.category_path.map((name, index) => ({
                        id: listingResponse.data.data.category_path_ids[index],
                        name: name,
                        slug: listingResponse.data.data.category_path_slugs[index]
                    })).reverse();
                    setCategoryPath(path);
                }
            } catch (err) {
                console.error('Error fetching listing:', err);
                setError('Не удалось загрузить объявление');
            } finally {
                setLoading(false);
            }
        };

        fetchListing();
    }, [id]);
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
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0
        }).format(price);
    };
    const handleDelete = async () => {
        if (!window.confirm('Вы действительно хотите удалить объявление?')) {
            return;
        }

        try {
            await axios.delete(`/api/v1/marketplace/listings/${id}`);
            navigate('/marketplace');
        } catch (error) {
            setError('Ошибка при удалении объявления');
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
            // Оптимистичное обновление UI
            setListing(prev => ({
                ...prev,
                is_favorite: !prev.is_favorite
            }));

            if (listing?.is_favorite) {
                await axios.delete(`/api/v1/marketplace/listings/${id}/favorite`);
            } else {
                await axios.post(`/api/v1/marketplace/listings/${id}/favorite`);
            }

            // Получаем актуальные данные
            const [listingResponse, favoritesResponse] = await Promise.all([
                axios.get(`/api/v1/marketplace/listings/${id}`),
                axios.get('/api/v1/marketplace/favorites')
            ]);

            const isFavorite = favoritesResponse.data?.data?.some?.(
                item => item.id === Number(id)
            ) || false;

            setListing({
                ...listingResponse.data.data,
                is_favorite: isFavorite
            });
        } catch (err) {
            // Возвращаем предыдущее состояние при ошибке
            setListing(prev => ({
                ...prev,
                is_favorite: !prev.is_favorite
            }));
            console.error('Ошибка при обновлении избранного:', err);
        }
    };
    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Breadcrumbs paths={categoryPath} />
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
                <Breadcrumbs paths={categoryPath} />

                <Typography color="error">{error}</Typography>
            </Container>
        );
    }

    if (!listing) return null;

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Breadcrumbs paths={categoryPath} />
            <Grid container spacing={4}>
                {/* Галерея изображений */}
                <Grid item xs={12} md={8}>
                    <Box sx={{ position: 'relative' }}>
                        {listing.images && listing.images.length > 0 ? (
                            <>
                                <Box
                                    component="img"
                                    src={`${process.env.REACT_APP_BACKEND_URL}/uploads/${listing.images[currentImageIndex].file_path}`}
                                    alt={listing.title}
                                    sx={{
                                        width: '100%',
                                        height: isMobile ? '300px' : '500px',
                                        objectFit: 'cover',
                                        borderRadius: 2
                                    }}
                                />
                                {listing.images.length > 1 && (
                                    <>
                                        <IconButton
                                            sx={{
                                                position: 'absolute',
                                                left: 8,
                                                top: '50%',
                                                transform: 'translateY(-50%)',
                                                bgcolor: 'background.paper',
                                                '&:hover': { bgcolor: 'background.paper' }
                                            }}
                                            onClick={() => setCurrentImageIndex(prev =>
                                                prev > 0 ? prev - 1 : listing.images.length - 1
                                            )}
                                        >
                                            <ChevronLeft />
                                        </IconButton>
                                        <IconButton
                                            sx={{
                                                position: 'absolute',
                                                right: 8,
                                                top: '50%',
                                                transform: 'translateY(-50%)',
                                                bgcolor: 'background.paper',
                                                '&:hover': { bgcolor: 'background.paper' }
                                            }}
                                            onClick={() => setCurrentImageIndex(prev =>
                                                prev < listing.images.length - 1 ? prev + 1 : 0
                                            )}
                                        >
                                            <ChevronRight />
                                        </IconButton>
                                    </>
                                )}
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
                                    Нет изображений
                                </Typography>
                            </Box>
                        )}
                    </Box>

                    {listing.images && listing.images.length > 1 && (
                        <ImageList
                            sx={{ mt: 2, maxHeight: 100 }}
                            cols={Math.min(listing.images.length, 6)}
                            rowHeight={100}
                        >
                            {listing.images.map((image, index) => (
                                <ImageListItem
                                    key={image.id}
                                    sx={{
                                        cursor: 'pointer',
                                        opacity: currentImageIndex === index ? 1 : 0.6,
                                        transition: 'opacity 0.2s',
                                        '&:hover': { opacity: 1 }
                                    }}
                                    onClick={() => setCurrentImageIndex(index)}
                                >
                                    <img
                                        src={`${process.env.REACT_APP_BACKEND_URL}/uploads/${image.file_path}`}
                                        alt={`${listing.title} ${index + 1}`}
                                        style={{
                                            height: '100%',
                                            objectFit: 'cover'
                                        }}
                                    />
                                </ImageListItem>
                            ))}
                        </ImageList>
                    )}

                    {/* Описание объявления */}
                    <Box sx={{ mt: 4 }}>
                        <Typography variant="h4" gutterBottom>
                            {listing.title}
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
                                    '&:hover': {
                                        color: 'primary.dark'
                                    }
                                }}
                            >
                                {reviewsCount} отзывов
                            </Box>
                        </Stack>
                        <Typography variant="body1" sx={{ mb: 4 }}>
                            {listing.description}
                        </Typography>

                        {/* Отзывы */}
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

                {/* Правая панель */}
                <Grid item xs={12} md={4}>
                    <Box sx={{ position: 'sticky', top: 24 }}>
                        {/* Карточка с ценой и контактами */}
                        <Card elevation={2}>
                            <CardContent>
                                <Typography variant="h4" gutterBottom>
                                    {formatPrice(listing.price)}
                                </Typography>

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
                                        ) : listing?.is_favorite ? 'В избранном' : 'В избранное'}
                                    </Button>
                                    <ShareButton
                                        url={window.location.href}
                                        title={listing.title}
                                        isMobile={isMobile}
                                    />
                                </Stack>

                                {/* кнопки удаления и  редактирования */}
                                {user?.id === listing.user_id && (
                                    <Stack direction="row" spacing={1} sx={{ mt: 2 }}>
                                        <Button
                                            variant="outlined"
                                            fullWidth
                                            startIcon={!isMobile && <PencilLine />}
                                            onClick={() => navigate(`/marketplace/listings/${id}/edit`)}
                                        >
                                            {isMobile ? <PencilLine size={20} /> : 'Редактировать'}
                                        </Button>
                                        <Button
                                            variant="outlined"
                                            color="error"
                                            fullWidth
                                            startIcon={!isMobile && <Trash2 />}
                                            onClick={handleDelete}
                                        >
                                            {isMobile ? <Trash2 size={20} /> : 'Удалить'}
                                        </Button>
                                    </Stack>
                                )}
                            </CardContent>
                        </Card>

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
                                            <GoogleMap
                                                mapContainerStyle={{
                                                    width: '100%',
                                                    height: '80vh'
                                                }}
                                                center={{
                                                    lat: listing.latitude,
                                                    lng: listing.longitude
                                                }}
                                                zoom={15}
                                                options={{
                                                    zoomControl: true,
                                                    mapTypeControl: true,
                                                    streetViewControl: true,
                                                    gestureHandling: "greedy"
                                                }}
                                            >
                                                <Marker
                                                    position={{
                                                        lat: listing.latitude,
                                                        lng: listing.longitude
                                                    }}
                                                    title={listing.title}
                                                />
                                            </GoogleMap>
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
                        {/* Карточка продавца */}
                        <Card elevation={2} sx={{ mt: 2 }}>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Продавец
                                </Typography>
                                <Stack direction="row" spacing={2} alignItems="center">
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
                                            На сайте с {listing.user?.created_at ?
                                                new Intl.DateTimeFormat('ru-RU', {
                                                    year: 'numeric',
                                                    month: 'long',
                                                    day: 'numeric'
                                                }).format(new Date(listing.user?.created_at))
                                                : 'неизвестной даты'}
                                        </Typography>
                                    </Box>
                                </Stack>
                            </CardContent>
                        </Card>
                    </Box>
                </Grid>
            </Grid>
        </Container>
    );
};

export default ListingDetailsPage;