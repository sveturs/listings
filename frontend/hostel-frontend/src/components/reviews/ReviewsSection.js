import React, { useState, useEffect, useCallback } from 'react';
import {
    Box,
    Typography,
    Rating,
    Stack,
    LinearProgress,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    CircularProgress
} from '@mui/material';
import { PencilLine } from 'lucide-react';
import ReviewCard from './ReviewCard';
import ReviewForm from './ReviewForm';
import axios from '../../api/axios';

export const ReviewsSection = ({ 
    entityType,
    entityId,
    entityTitle,
    canAddReview = true
}) => {
    const [reviews, setReviews] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showReviewForm, setShowReviewForm] = useState(false);
    const [stats, setStats] = useState({
        averageRating: 0,
        totalReviews: 0,
        ratingDistribution: {
            1: 0, 2: 0, 3: 0, 4: 0, 5: 0
        }
    });

    const fetchReviews = useCallback(async () => {
        try {
            setLoading(true);
            const response = await axios.get('/api/v1/reviews', {
                params: {
                    entity_type: entityType,
                    entity_id: entityId
                }
            });
            setReviews(response.data.data || []);
        } catch (error) {
            console.error('Error fetching reviews:', error);
        } finally {
            setLoading(false);
        }
    }, [entityType, entityId]);

    const fetchStats = useCallback(async () => {
        try {
            const response = await axios.get('/api/v1/reviews/stats', {
                params: {
                    entity_type: entityType,
                    entity_id: entityId
                }
            });
            setStats(response.data.data || {
                averageRating: 0,
                totalReviews: 0,
                ratingDistribution: { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 }
            });
        } catch (error) {
            console.error('Error fetching review stats:', error);
        }
    }, [entityType, entityId]);

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" p={4}>
                <CircularProgress />
            </Box>
        );
    }

    return (
        <Box>
            {/* Статистика рейтингов */}
            <Stack direction="row" spacing={4} alignItems="center" sx={{ mb: 4 }}>
                <Box textAlign="center">
                    <Typography variant="h3" fontWeight="bold">
                        {stats.averageRating.toFixed(1)}
                    </Typography>
                    <Rating value={stats.averageRating} readOnly precision={0.1} />
                    <Typography color="text.secondary">
                        {stats.totalReviews} отзывов
                    </Typography>
                </Box>

                {/* Распределение рейтингов */}
                <Box flex={1}>
                    {Object.entries(stats.ratingDistribution).reverse().map(([rating, count]) => (
                        <Stack key={rating} direction="row" spacing={2} alignItems="center" sx={{ mb: 1 }}>
                            <Typography minWidth={20}>{rating}</Typography>
                            <LinearProgress
                                variant="determinate"
                                value={(count / stats.totalReviews) * 100 || 0}
                                sx={{ flex: 1, height: 8, borderRadius: 1 }}
                            />
                            <Typography minWidth={40}>{count}</Typography>
                        </Stack>
                    ))}
                </Box>

                {canAddReview && (
                    <Button
                        variant="contained"
                        onClick={() => setShowReviewForm(true)}
                        startIcon={<PencilLine />}
                    >
                        Написать отзыв
                    </Button>
                )}
            </Stack>

            {/* Список отзывов */}
            <Stack spacing={2}>
                {reviews.map(review => (
                    <ReviewCard key={review.id} review={review} />
                ))}
            </Stack>

            {/* Диалог добавления отзыва */}
            <Dialog
                open={showReviewForm}
                onClose={() => setShowReviewForm(false)}
                maxWidth="md"
                fullWidth
            >
                <DialogTitle>
                    Отзыв о {entityTitle}
                </DialogTitle>
                <DialogContent>
                    <ReviewForm 
                        onSubmit={async (reviewData) => {
                            try {
                                await axios.post('/api/v1/reviews', {
                                    ...reviewData,
                                    entity_type: entityType,
                                    entity_id: entityId
                                });
                                setShowReviewForm(false);
                                fetchReviews();
                                fetchStats();
                            } catch (error) {
                                console.error('Error submitting review:', error);
                            }
                        }}
                    />
                </DialogContent>
            </Dialog>
        </Box>
    );
};

export default ReviewsSection;