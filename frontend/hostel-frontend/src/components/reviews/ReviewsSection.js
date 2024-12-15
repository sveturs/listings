//frontend/hostel-frontend/src/components/reviews/ReviewsSection.js
import React, { useState, useEffect } from 'react';
import { Box, Button, Dialog, DialogTitle, DialogContent, Alert, Snackbar } from '@mui/material';
import { PencilLine } from 'lucide-react';
import { ReviewForm, ReviewCard, RatingStats } from './ReviewComponents';
import axios from '../../api/axios';

const ReviewsSection = ({
    entityType, // тип сущности (listing, room, car)
    entityId,   // ID сущности
    entityTitle, // название сущности для отображения
    canReview = true,
    onReviewsCountChange
}) => {
    const [reviews, setReviews] = useState([]);
    const [stats, setStats] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showReviewForm, setShowReviewForm] = useState(false);
    const [editingReview, setEditingReview] = useState(null);
    const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });
    

    // Загрузка отзывов и статистики
    const fetchData = async () => {
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
            console.log('Reviews response:', reviewsResponse.data);
            console.log('Stats response:', statsResponse.data);
            
            setReviews(reviewsResponse.data.data.data || []); // Изменена эта строка
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
    // Обработка создания/редактирования отзыва
    const handleReviewSubmit = async ({ reviewData, photosFormData }) => {
        try {
            console.log('Sending review data:', reviewData);

            // Сначала создаем отзыв
            const response = await axios.post('/api/v1/reviews', reviewData);

            // Если есть фотографии и отзыв успешно создан - загружаем их
            if (photosFormData && response.data && response.data.id) {
                await axios.post(`/api/v1/reviews/${response.data.id}/photos`, photosFormData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            }

            setShowReviewForm(false);
            setEditingReview(null);
            fetchData();
            setSnackbar({
                open: true,
                message: 'Отзыв успешно создан',
                severity: 'success'
            });
        } catch (err) {
            console.error('Error submitting review:', err);
            setSnackbar({
                open: true,
                message: err.response?.data?.error || 'Ошибка при сохранении отзыва',
                severity: 'error'
            });
        }
    };
    // Обработка голосования за отзыв
    const handleVote = async (reviewId, voteType) => {
        try {
            console.log('Voting for review:', reviewId, voteType);
            await axios.post(`/api/v1/reviews/${reviewId}/vote`, { vote_type: voteType });
            console.log('Vote successful, fetching new data...');
            await fetchData();
            console.log('Data updated after vote');
        } catch (err) {
            console.error('Vote error:', err);
            setSnackbar({
                open: true,
                message: 'Ошибка при голосовании',
                severity: 'error'
            });
        }
    };

    // Обработка ответа на отзыв
    const handleReply = async (reviewId, response) => {
        try {
            await axios.post(`/api/v1/reviews/${reviewId}/response`, { response });
            fetchData();
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

    // Обработка удаления отзыва
    const handleDelete = async (reviewId) => {
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

    // Обработка жалобы на отзыв
    const handleReport = async (reviewId) => {
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
            {/* Статистика рейтингов */}
            {stats && <RatingStats stats={stats} />}

            {/* Кнопка добавления отзыва */}
            {canReview && (
                <Button
                    variant="contained"
                    onClick={() => setShowReviewForm(true)}
                    startIcon={<PencilLine />}
                    sx={{ mb: 3 }}
                >
                    Написать отзыв
                </Button>
            )}

            {/* Список отзывов */}
            {Array.isArray(reviews) && reviews.map(review => {
                console.log('Review data:', review); // Добавляем это
                return (
                    <ReviewCard
                        key={review.id}
                        review={review}
                        onVote={handleVote}
                        onReply={handleReply}
                        onEdit={(review) => {
                            setEditingReview(review);
                            setShowReviewForm(true);
                        }}
                        onDelete={handleDelete}
                        onReport={handleReport}
                    />
                );
            })}

            {/* Диалог создания/редактирования отзыва */}
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

            {/* Уведомления */}
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