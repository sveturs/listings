import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
    Box,
    Paper,
    Typography,
    Rating,
    LinearProgress,
    Divider,
    CircularProgress,
    Button
} from '@mui/material';
import { Link } from 'react-router-dom';
import axios from '../../api/axios';

// Определяем интерфейсы для типизации
interface RatingData {
    average_rating: number;
    total_reviews: number;
    rating_1: number;
    rating_2: number;
    rating_3: number;
    rating_4: number;
    rating_5: number;
}

interface Review {
    id: number | string;
    rating: number;
    comment: string;
    created_at: string;
    entity_type: string;
    entity_id: number | string;
}

interface UserRatingProps {
    userId: number | string;
}

const UserRating: React.FC<UserRatingProps> = ({ userId }) => {
    const { t } = useTranslation('marketplace');
    const [ratingData, setRatingData] = useState<RatingData | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [showAllReviews, setShowAllReviews] = useState<boolean>(false);
    const [reviews, setReviews] = useState<Review[]>([]);
    const [loadingReviews, setLoadingReviews] = useState<boolean>(false);

    useEffect(() => {
        const fetchRating = async () => {
            try {
                setLoading(true);
                const response = await axios.get(`/api/v1/users/${userId}/rating`);
                setRatingData(response.data.data);
            } catch (err) {
                console.error('Error fetching user rating:', err);
                setError(t('reviews.errors.loadFailed'));
            } finally {
                setLoading(false);
            }
        };

        if (userId) {
            fetchRating();
        }
    }, [userId, t]);

    const loadAllReviews = async () => {
        if (showAllReviews && reviews.length > 0) {
            setShowAllReviews(false);
            return;
        }

        try {
            setLoadingReviews(true);
            const response = await axios.get(`/api/v1/users/${userId}/reviews`);
            setReviews(response.data.data);
            setShowAllReviews(true);
        } catch (err) {
            console.error('Error fetching user reviews:', err);
            setError(t('reviews.errors.loadFailed'));
        } finally {
            setLoadingReviews(false);
        }
    };

    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                <CircularProgress />
            </Box>
        );
    }

    if (error) {
        return (
            <Box sx={{ p: 2, color: 'error.main' }}>
                <Typography>{error}</Typography>
            </Box>
        );
    }

    if (!ratingData || ratingData.total_reviews === 0) {
        return (
            <Paper sx={{ p: 3, mb: 3 }}>
                <Typography variant="subtitle1" gutterBottom>
                    {t('reviews.sellerRating')}
                </Typography>
                <Typography color="text.secondary">
                    {t('reviews.noRatingsYet')}
                </Typography>
            </Paper>
        );
    }

    const calculatePercentage = (value: number): number => {
        return (value / ratingData.total_reviews) * 100;
    };

    return (
        <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
                {t('reviews.sellerRating')}
            </Typography>

            <Box display="flex" alignItems="center" gap={1} mb={2}>
                <Typography variant="h4" fontWeight="bold">
                    {ratingData.average_rating.toFixed(1)}
                </Typography>
                <Box>
                    <Rating
                        value={ratingData.average_rating}
                        precision={0.1}
                        readOnly
                        size="large"
                    />
                    <Typography variant="body2" color="text.secondary">
                        {t('reviews.basedOn', { count: ratingData.total_reviews })}
                    </Typography>
                </Box>
            </Box>

            <Divider sx={{ my: 2 }} />

            <Box>
                {[5, 4, 3, 2, 1].map((rating) => (
                    <Box key={rating} display="flex" alignItems="center" mb={1} gap={2}>
                        <Box minWidth={20}>
                            <Typography>{rating}</Typography>
                        </Box>
                        <Box flex={1}>
                            <LinearProgress
                                variant="determinate"
                                value={calculatePercentage(ratingData[`rating_${rating}` as keyof RatingData])}
                                sx={{ height: 8, borderRadius: 1 }}
                            />
                        </Box>
                        <Box minWidth={30}>
                            <Typography align="right">
                                {ratingData[`rating_${rating}` as keyof RatingData]}
                            </Typography>
                        </Box>
                    </Box>
                ))}
            </Box>

            <Button
                onClick={loadAllReviews}
                variant="outlined"
                fullWidth
                sx={{ mt: 2 }}
                disabled={loadingReviews}
            >
                {loadingReviews ? (
                    <CircularProgress size={24} />
                ) : showAllReviews ? (
                    t('reviews.hideReviews')
                ) : (
                    t('reviews.seeAllReviews')
                )}
            </Button>

            {showAllReviews && reviews.length > 0 && (
                <Box mt={2}>
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="h6" gutterBottom>
                        {t('reviews.allUserReviews')}
                    </Typography>
                    {reviews.map((review) => (
                        <Box key={review.id} mb={2} p={2} bgcolor="background.subtle" borderRadius={1}>
                            <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                                <Rating value={review.rating} readOnly size="small" />
                                <Typography variant="caption" color="text.secondary">
                                    {new Date(review.created_at).toLocaleDateString()}
                                </Typography>
                            </Box>
                            <Typography variant="body2">{review.comment}</Typography>
                            {review.entity_type === 'listing' && (
                                <Box mt={1}>
                                    <Typography variant="caption" color="text.secondary">
                                        {t('reviews.forListing')}:
                                    </Typography>
                                    <Button
                                        component={Link}
                                        to={`/marketplace/listings/${review.entity_id}`}
                                        size="small"
                                        variant="text"
                                        sx={{ ml: 1 }}
                                    >
                                        {t('reviews.viewListing')}
                                    </Button>
                                </Box>
                            )}
                        </Box>
                    ))}
                </Box>
            )}
            <Box display="flex" justifyContent="center" mt={2}>
                <Button
                    component={Link}
                    to={`/user/${userId}/reviews`}
                    variant="outlined"
                    size="small"
                >
                    {t('reviews.seeAllUserReviews')}
                </Button>
            </Box>
        </Paper>
    );
};

export default UserRating;