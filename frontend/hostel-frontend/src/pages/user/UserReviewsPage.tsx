import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
    Container,
    Typography,
    Box,
    CircularProgress,
    Grid,
    Paper,
    Rating,
    Divider,
    Button,
    Avatar,
    Card,
    CardContent,
    Alert
} from '@mui/material';
import axios from '../../api/axios';
import { ArrowLeft } from 'lucide-react';

// Интерфейсы для типизации данных
interface UserProfile {
    id: number | string;
    name: string;
    picture_url?: string;
    created_at: string;
    [key: string]: any;
}

interface ReviewUser {
    name?: string;
    picture_url?: string;
    [key: string]: any;
}

interface Review {
    id: number | string;
    rating: number;
    comment?: string;
    pros?: string;
    cons?: string;
    entity_type?: string;
    entity_id?: number | string;
    created_at: string;
    user?: ReviewUser;
    [key: string]: any;
}

interface RatingSummary {
    average_rating: number;
    total_reviews: number;
    rating_1: number;
    rating_2: number;
    rating_3: number;
    rating_4: number;
    rating_5: number;
    [key: string]: any;
}

interface RouteParams {
    id: string;
}

const UserReviewsPage: React.FC = () => {
    const { id } = useParams<"id">();
    const { t } = useTranslation('marketplace');
    const [userData, setUserData] = useState<UserProfile | null>(null);
    const [reviews, setReviews] = useState<Review[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [summary, setSummary] = useState<RatingSummary | null>(null);

    useEffect(() => {
        if (!id) {
            setError(t('reviews.userNotFound'));
            return;
        }

        const fetchData = async (): Promise<void> => {
            try {
                setLoading(true);
                // Используйте правильный эндпоинт для получения данных пользователя
                const userResponse = await axios.get(`/api/v1/users/${id}/profile`);
                if (!userResponse.data.data) {
                    throw new Error('User not found');
                }
                setUserData(userResponse.data.data);
                // Получаем обзоры пользователя
                const reviewsResponse = await axios.get(`/api/v1/users/${id}/reviews`);
                setReviews(reviewsResponse.data.data || []);

                // Получаем сводные данные о рейтинге
                const summaryResponse = await axios.get(`/api/v1/users/${id}/rating`);
                setSummary(summaryResponse.data.data);
            } catch (err) {
                console.error('Error fetching user reviews:', err);
                setError(t('reviews.errors.loadFailed'));
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [id, t]);

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
                    <CircularProgress />
                </Box>
            </Container>
        );
    }

    if (error) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Alert severity="error">{error}</Alert>
            </Container>
        );
    }

    if (!userData) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Alert severity="warning">{t('reviews.userNotFound')}</Alert>
            </Container>
        );
    }

    const calculatePercentage = (value: number): number => {
        return summary && summary.total_reviews > 0
            ? (value / summary.total_reviews) * 100
            : 0;
    };

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Button
                component={Link}
                to="/"
                startIcon={<ArrowLeft />}
                sx={{ mb: 3 }}
            >
                {t('common.backToHome')}
            </Button>

            <Grid container spacing={4}>
                <Grid item xs={12} md={4}>
                    <Card>
                        <CardContent>
                            <Box
                                display="flex"
                                alignItems="center"
                                gap={2}
                                mb={3}
                            >
                                <Avatar
                                    src={userData.picture_url}
                                    alt={userData.name}
                                    sx={{ width: 80, height: 80 }}
                                />
                                <Box>
                                    <Typography variant="h5">
                                        {userData.name}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        {t('reviews.memberSince', {
                                            date: new Date(userData.created_at).toLocaleDateString()
                                        })}
                                    </Typography>
                                </Box>
                            </Box>

                            {summary && summary.total_reviews > 0 ? (
                                <Box>
                                    <Divider sx={{ mb: 2 }} />
                                    <Typography variant="h6" gutterBottom>
                                        {t('reviews.sellerRating')}
                                    </Typography>

                                    <Box display="flex" alignItems="center" gap={1} mb={2}>
                                        <Typography variant="h4" fontWeight="bold">
                                            {summary.average_rating.toFixed(1)}
                                        </Typography>
                                        <Box>
                                            <Rating
                                                value={summary.average_rating}
                                                precision={0.1}
                                                readOnly
                                                size="large"
                                            />
                                            <Typography variant="body2" color="text.secondary">
                                                {t('reviews.basedOn', { count: summary.total_reviews })}
                                            </Typography>
                                        </Box>
                                    </Box>

                                    <Box>
                                        {[5, 4, 3, 2, 1].map((rating) => (
                                            <Box key={rating} display="flex" alignItems="center" mb={1} gap={2}>
                                                <Box minWidth={20}>
                                                    <Typography>{rating}</Typography>
                                                </Box>
                                                <Box flex={1} sx={{ backgroundColor: 'grey.200', height: 10, borderRadius: 1 }}>
                                                    <Box
                                                        sx={{
                                                            width: `${calculatePercentage(summary[`rating_${rating}`])}%`,
                                                            backgroundColor: 'primary.main',
                                                            height: '100%',
                                                            borderRadius: 1
                                                        }}
                                                    />
                                                </Box>
                                                <Box minWidth={30}>
                                                    <Typography align="right">
                                                        {summary[`rating_${rating}`]}
                                                    </Typography>
                                                </Box>
                                            </Box>
                                        ))}
                                    </Box>
                                </Box>
                            ) : (
                                <Box textAlign="center" py={3}>
                                    <Typography>{t('reviews.noRatingsYet')}</Typography>
                                </Box>
                            )}
                        </CardContent>
                    </Card>
                </Grid>

                <Grid item xs={12} md={8}>
                    <Paper sx={{ p: 3 }}>
                        <Typography variant="h5" gutterBottom>
                            {t('reviews.allReviews')}
                        </Typography>

                        {reviews.length === 0 ? (
                            <Box textAlign="center" py={4}>
                                <Typography variant="body1" color="text.secondary">
                                    {t('reviews.noReviewsYet')}
                                </Typography>
                            </Box>
                        ) : (
                            reviews.map((review) => (
                                <Box
                                    key={review.id}
                                    sx={{
                                        mb: 3,
                                        p: 2,
                                        border: 1,
                                        borderColor: 'divider',
                                        borderRadius: 1
                                    }}
                                >
                                    <Box
                                        display="flex"
                                        justifyContent="space-between"
                                        alignItems="center"
                                        mb={1}
                                    >
                                        <Box display="flex" alignItems="center" gap={1}>
                                            <Avatar
                                                src={review.user?.picture_url}
                                                alt={review.user?.name || 'User'}
                                                sx={{ width: 40, height: 40 }}
                                            />
                                            <Typography>{review.user?.name || 'User'}</Typography>
                                        </Box>
                                        <Typography variant="caption" color="text.secondary">
                                            {new Date(review.created_at).toLocaleDateString()}
                                        </Typography>
                                    </Box>

                                    <Rating value={review.rating} readOnly sx={{ mb: 1 }} />

                                    <Typography variant="body1" paragraph>
                                        {review.comment}
                                    </Typography>

                                    {review.pros && (
                                        <Box mb={1}>
                                            <Typography variant="caption" color="success.main" fontWeight="bold">
                                                {t('reviews.pros')}:
                                            </Typography>
                                            <Typography variant="body2">{review.pros}</Typography>
                                        </Box>
                                    )}

                                    {review.cons && (
                                        <Box mb={1}>
                                            <Typography variant="caption" color="error.main" fontWeight="bold">
                                                {t('reviews.cons')}:
                                            </Typography>
                                            <Typography variant="body2">{review.cons}</Typography>
                                        </Box>
                                    )}

                                    {review.entity_type === 'listing' && (
                                        <Box mt={2}>
                                            <Divider sx={{ mb: 1 }} />
                                            <Typography variant="caption" color="text.secondary">
                                                {t('reviews.forListing')}:
                                            </Typography>
                                            <Button
                                                component={Link}
                                                to={`/marketplace/listings/${review.entity_id}`}
                                                size="small"
                                                variant="outlined"
                                                sx={{ ml: 1 }}
                                            >
                                                {t('reviews.viewListing')}
                                            </Button>
                                        </Box>
                                    )}
                                </Box>
                            ))
                        )}
                    </Paper>
                </Grid>
            </Grid>
        </Container>
    );
};

export default UserReviewsPage;