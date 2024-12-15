//frontend/hostel-frontend/src/components/Review.js
import React, { useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Rating,
    Button,
    TextField,
    Avatar,
    Stack,
    Chip,
    IconButton,
    List,
    ListItem
} from '@mui/material';
import {
    ThumbsUp,
    ThumbsDown,
    MessageSquare,
    Flag,
    Check
} from 'lucide-react';

// Компонент для отображения отзыва
const ReviewCard = ({ review, onVote, onRespond, onReport, currentUserId }) => {
    const [showResponseForm, setShowResponseForm] = useState(false);
    const [response, setResponse] = useState('');

    const handleSubmitResponse = () => {
        onRespond(review.id, response);
        setShowResponseForm(false);
        setResponse('');
    };

    return (
        <Card sx={{ mb: 2, borderRadius: 2 }}>
            <CardContent>
                <Stack spacing={2}>
                    {/* Заголовок с информацией о пользователе */}
                    <Stack direction="row" spacing={2} alignItems="center">
                        <Avatar src={review.user?.picture_url} />
                        <Box>
                            <Typography variant="subtitle1" fontWeight="bold">
                                {review.user?.name}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                                {new Date(review.created_at).toLocaleDateString()}
                            </Typography>
                        </Box>
                        {review.is_verified_purchase && (
                            <Chip
                                icon={<Check size={16} />}
                                label="Проверенная покупка"
                                color="success"
                                size="small"
                            />
                        )}
                    </Stack>

                    {/* Рейтинг */}
                    <Box>
                        <Rating value={review.rating} readOnly precision={1} />
                    </Box>

                    {/* Основной контент */}
                    <Box>
                        {review.comment && (
                            <Typography variant="body1" paragraph>
                                {review.comment}
                            </Typography>
                        )}

                        {review.pros && (
                            <Box sx={{ mb: 1 }}>
                                <Typography variant="subtitle2" color="success.main">
                                    Достоинства:
                                </Typography>
                                <Typography variant="body2">{review.pros}</Typography>
                            </Box>
                        )}

                        {review.cons && (
                            <Box sx={{ mb: 1 }}>
                                <Typography variant="subtitle2" color="error.main">
                                    Недостатки:
                                </Typography>
                                <Typography variant="body2">{review.cons}</Typography>
                            </Box>
                        )}
                    </Box>

                    {/* Фотографии */}
                    {review.photos && review.photos.length > 0 && (
                        <Stack direction="row" spacing={1} sx={{ overflowX: 'auto' }}>
                            {review.photos.map((photo, index) => (
                                <Box
                                    key={index}
                                    component="img"
                                    src={photo}
                                    sx={{
                                        height: 100,
                                        width: 100,
                                        objectFit: 'cover',
                                        borderRadius: 1
                                    }}
                                />
                            ))}
                        </Stack>
                    )}

                    {/* Действия */}
                    <Stack direction="row" spacing={2} alignItems="center">
                        <Button
                            size="small"
                            startIcon={<ThumbsUp />}
                            onClick={() => onVote(review.id, 'helpful')}
                            color={review.current_user_vote === 'helpful' ? 'success' : 'inherit'}
                        >
                            Полезно ({review.votes_count?.helpful || 0})
                        </Button>
                        <Button
                            size="small"
                            startIcon={<ThumbsDown />}
                            onClick={() => onVote(review.id, 'not_helpful')}
                            color={review.current_user_vote === 'not_helpful' ? 'error' : 'inherit'}
                        >
                            Не полезно ({review.votes_count?.not_helpful || 0})
                        </Button>
                        <Button
                            size="small"
                            startIcon={<MessageSquare />}
                            onClick={() => setShowResponseForm(!showResponseForm)}
                        >
                            Ответить
                        </Button>
                        <IconButton onClick={() => onReport(review.id)}>
                            <Flag size={20} />
                        </IconButton>
                    </Stack>

                    {/* Форма ответа */}
                    {showResponseForm && (
                        <Box sx={{ mt: 2 }}>
                            <TextField
                                fullWidth
                                multiline
                                rows={3}
                                placeholder="Напишите ответ..."
                                value={response}
                                onChange={(e) => setResponse(e.target.value)}
                            />
                            <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
                                <Button
                                    variant="contained"
                                    size="small"
                                    onClick={handleSubmitResponse}
                                >
                                    Отправить
                                </Button>
                                <Button
                                    size="small"
                                    onClick={() => setShowResponseForm(false)}
                                >
                                    Отмена
                                </Button>
                            </Stack>
                        </Box>
                    )}

                    {/* Ответы */}
                    {review.responses && review.responses.length > 0 && (
                        <List>
                            {review.responses.map((response) => (
                                <ListItem key={response.id} sx={{ pl: 4 }}>
                                    <Stack spacing={1}>
                                        <Stack direction="row" spacing={2} alignItems="center">
                                            <Avatar
                                                src={response.user?.picture_url}
                                                sx={{ width: 24, height: 24 }}
                                            />
                                            <Typography variant="subtitle2">
                                                {response.user?.name}
                                            </Typography>
                                            <Typography variant="caption" color="text.secondary">
                                                {new Date(response.created_at).toLocaleDateString()}
                                            </Typography>
                                        </Stack>
                                        <Typography variant="body2">
                                            {response.response}
                                        </Typography>
                                    </Stack>
                                </ListItem>
                            ))}
                        </List>
                    )}
                </Stack>
            </CardContent>
        </Card>
    );
};


export default ReviewForm;