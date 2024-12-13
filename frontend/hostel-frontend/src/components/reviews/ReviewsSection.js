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
    canReview = true // может ли пользователь оставлять отзывы
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
    
            // Добавим проверку данных
            setReviews(reviewsResponse.data.data || []);  // Если нет данных, используем пустой массив
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

    // Обработка создания/редактирования отзыва
    const handleReviewSubmit = async (formData) => {
        try {
            if (editingReview) {
                await axios.put(`/api/v1/reviews/${editingReview.id}`, formData);
                setSnackbar({
                    open: true,
                    message: 'Отзыв успешно обновлен',
                    severity: 'success'
                });
            } else {
                await axios.post('/api/v1/reviews', {
                    ...formData,
                    entity_type: entityType,
                    entity_id: entityId
                });
                setSnackbar({
                    open: true,
                    message: 'Отзыв успешно опубликован',
                    severity: 'success'
                });
            }
            setShowReviewForm(false);
            setEditingReview(null);
            fetchData();
        } catch (err) {
            setSnackbar({
                open: true,
                message: 'Ошибка при сохранении отзыва',
                severity: 'error'
            });
        }
    };

    // Обработка голосования за отзыв
    const handleVote = async (reviewId, voteType) => {
        try {
            await axios.post(`/api/v1/reviews/${reviewId}/vote`, { vote_type: voteType });
            fetchData();
        } catch (err) {
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
            {Array.isArray(reviews) && reviews.map(review => (
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
            ))}

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