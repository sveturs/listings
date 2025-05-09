import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

import { Box, Button, Dialog, DialogTitle, DialogContent, Alert, Snackbar } from '@mui/material';
import { PencilLine } from 'lucide-react';
import { ReviewForm, ReviewCard, RatingStats } from './ReviewComponents';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';

interface ReviewStat {
    average_rating: number;
    total_reviews: number;
    rating_distribution: Record<number, number>;
}

interface ReviewResponse {
    id: string | number;
    user_id: string | number;
    content: string;
    created_at: string;
    user?: {
        name: string;
        picture_url?: string;
    };
}

interface ReviewVotes {
    helpful: number;
    not_helpful: number;
}

export interface Review {
    id: string | number;
    user_id: string | number;
    entity_type: string;
    entity_id: string | number;
    rating: number;
    comment: string;
    pros?: string;
    cons?: string;
    created_at: string;
    is_verified_purchase?: boolean;
    current_user_vote?: 'helpful' | 'not_helpful';
    helpful_votes?: number;
    not_helpful_votes?: number;
    votes_count: ReviewVotes;
    photos?: string[];
    responses?: ReviewResponse[];
    original_language?: string;
    translations?: Record<string, Record<string, string>>;
    user?: {
        name: string;
        picture_url?: string;
    };
}

interface ReviewFormData {
    reviewData: {
        entity_type: string;
        entity_id: string | number;
        rating: number;
        comment: string;
        pros?: string;
        cons?: string;
        original_language: string;
    };
    photosFormData: FormData | null;
}

interface SnackbarState {
    open: boolean;
    message: string;
    severity: 'success' | 'error' | 'warning' | 'info';
}

interface ReviewsSectionProps {
    entityType: string;
    entityId: string | number;
    entityTitle?: string;
    canReview?: boolean;
    onReviewsCountChange?: (count: number) => void;
}

const ReviewsSection: React.FC<ReviewsSectionProps> = ({
    entityType,
    entityId,
    entityTitle,
    canReview = true,
    onReviewsCountChange
}) => {
    const { t, i18n } = useTranslation('marketplace'); 

    const [reviews, setReviews] = useState<Review[]>([]);
    const [stats, setStats] = useState<ReviewStat | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [showReviewForm, setShowReviewForm] = useState<boolean>(false);
    const [editingReview, setEditingReview] = useState<Review | null>(null);
    const [snackbar, setSnackbar] = useState<SnackbarState>({ 
        open: false, 
        message: '', 
        severity: 'success' 
    });
    const { user } = useAuth();
    

    // Загрузка отзывов и статистики
    const fetchData = async (): Promise<void> => {
        try {
            setLoading(true);
            const [reviewsResponse, statsResponse] = await Promise.all([
                axios.get('/api/v1/reviews', {
                    params: {
                        entity_type: entityType,
                        entity_id: entityId
                    }
                }),
                axios.get(`/api/v1/entity/${entityType}/${entityId}/stats`)
            ]);
    
            // Преобразуем данные
            const reviews = (reviewsResponse.data.data.data || []).map((review: Review) => ({
                ...review,
                votes_count: {
                    helpful: review.helpful_votes || 0,
                    not_helpful: review.not_helpful_votes || 0
                }
            }));
    
            setReviews(reviews);
            setStats(statsResponse.data.data);
        } catch (err) {
            setError('Не удалось загрузить отзывы');
            console.error('Error fetching reviews:', err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [entityType, entityId]);

    useEffect(() => {
        if (reviews && onReviewsCountChange) {
            onReviewsCountChange(reviews.length);
        }
    }, [reviews, onReviewsCountChange]);

    const handleReviewSubmit = async ({ reviewData, photosFormData }: ReviewFormData): Promise<void> => {
        try {
            const response = await axios.post('/api/v1/reviews', {
                ...reviewData,
                original_language: i18n.language // Добавляем текущий язык
            });
            const reviewId = response.data.data.id;

            if (photosFormData && photosFormData.getAll('photos').length > 0) {
                try {
                    await axios.post(`/api/v1/reviews/${reviewId}/photos`, photosFormData, {
                        headers: {
                            'Content-Type': 'multipart/form-data'
                        }
                    });
                } catch (photoErr) {
                    console.error('Error uploading photos:', photoErr);
                    setSnackbar({
                        open: true,
                        message: 'Отзыв создан, но возникла ошибка при загрузке фотографий',
                        severity: 'warning'
                    });
                    return;
                }
            }

            setShowReviewForm(false);
            setEditingReview(null);
            fetchData();
            setSnackbar({
                open: true,
                message: 'Отзыв успешно создан',
                severity: 'success'
            });
        } catch (err: any) {
            console.error('Error submitting review:', err);
            setSnackbar({
                open: true,
                message: err.response?.data?.error || 'Ошибка при сохранении отзыва',
                severity: 'error'
            });
        }
    };

    const handleVote = async (reviewId: string | number, voteType: 'helpful' | 'not_helpful'): Promise<void> => {
        const oldReviews = [...reviews];
    
        try {
            // Оптимистично обновляем UI
            setReviews(prevReviews =>
                prevReviews.map(review => {
                    if (review.id === reviewId) {
                        const votes_count = { ...review.votes_count };
                        
                        // Если был предыдущий голос, убираем его
                        if (review.current_user_vote) {
                            votes_count[review.current_user_vote]--;
                        }
                        
                        // Добавляем новый голос
                        votes_count[voteType] = (votes_count[voteType] || 0) + 1;
    
                        return {
                            ...review,
                            votes_count,
                            current_user_vote: voteType
                        };
                    }
                    return review;
                })
            );
    
            // Отправляем запрос на сервер
            await axios.post(`/api/v1/reviews/${reviewId}/vote`, {
                vote_type: voteType
            });
    
            // Обновляем данные с сервера
            const response = await axios.get(`/api/v1/reviews/${reviewId}`);
            setReviews(prevReviews =>
                prevReviews.map(review =>
                    review.id === reviewId ? response.data.data : review
                )
            );
        } catch (err) {
            // В случае ошибки возвращаем предыдущее состояние
            setReviews(oldReviews);
            setSnackbar({
                open: true,
                message: 'Ошибка при голосовании',
                severity: 'error'
            });
        }
    };

    const handleReply = async (reviewId: string | number, response: string): Promise<void> => {
        try {
            await axios.post(`/api/v1/reviews/${reviewId}/response`, { response });
            
            // Получаем обновленные данные отзыва
            const updatedReviewResponse = await axios.get(`/api/v1/reviews/${reviewId}`);
            
            // Обновляем состояние, заменяя старый отзыв на новый с ответом
            setReviews(prevReviews => 
                prevReviews.map(review => 
                    review.id === reviewId ? updatedReviewResponse.data.data : review
                )
            );
    
            setSnackbar({
                open: true,
                message: 'Ответ успешно добавлен',
                severity: 'success'
            });
        } catch (err) {
            setSnackbar({
                open: true,
                message: 'Ошибка при добавлении ответа',
                severity: 'error'
            });
        }
    };

    const handleDelete = async (reviewId: string | number): Promise<void> => {
        try {
            await axios.delete(`/api/v1/reviews/${reviewId}`);
            fetchData();
            setSnackbar({
                open: true,
                message: 'Отзыв успешно удален',
                severity: 'success'
            });
        } catch (err) {
            setSnackbar({
                open: true,
                message: 'Ошибка при удалении отзыва',
                severity: 'error'
            });
        }
    };

    const handleReport = async (reviewId: string | number): Promise<void> => {
        try {
            await axios.post(`/api/v1/reviews/${reviewId}/report`);
            setSnackbar({
                open: true,
                message: 'Жалоба отправлена',
                severity: 'success'
            });
        } catch (err) {
            setSnackbar({
                open: true,
                message: 'Ошибка при отправке жалобы',
                severity: 'error'
            });
        }
    };

    return (
        <Box>
            {stats && <RatingStats stats={stats} />}

            {canReview && (
                <Button
                    id="reviewButton"
                    variant="contained"
                    onClick={() => setShowReviewForm(true)}
                    startIcon={<PencilLine />}
                    sx={{ mb: 3 }}
                >
                    {t('reviews.write')}
                </Button>
            )}

            {Array.isArray(reviews) && reviews.map(review => (
                <ReviewCard
                    key={review.id}
                    review={review}
                    currentUserId={user?.id} // Добавляем передачу ID пользователя
                    onVote={handleVote}
                    onReply={handleReply}
                    onEdit={(review) => {
                        setEditingReview(review);
                        setShowReviewForm(true);
                    }}
                    onDelete={handleDelete}
                    onReport={handleReport}
                />
            ))}

            <Dialog
                open={showReviewForm}
                onClose={() => {
                    setShowReviewForm(false);
                    setEditingReview(null);
                }}
                maxWidth="md"
                fullWidth
            >
                <DialogTitle>
                    {editingReview ? 'Редактирование отзыва' : 'Новый отзыв'}
                </DialogTitle>
                <DialogContent>
                    <ReviewForm
                        entityType={entityType}
                        entityId={entityId}
                        initialData={editingReview}
                        onSubmit={handleReviewSubmit}
                        onCancel={() => {
                            setShowReviewForm(false);
                            setEditingReview(null);
                        }}
                    />
                </DialogContent>
            </Dialog>

            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={() => setSnackbar({ ...snackbar, open: false })}
            >
                <Alert
                    onClose={() => setSnackbar({ ...snackbar, open: false })}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Box>
    );
};

export default ReviewsSection;