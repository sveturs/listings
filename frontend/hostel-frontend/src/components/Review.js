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
    Divider,
    List,
    ListItem
} from '@mui/material';
import {
    ThumbsUp,
    ThumbsDown,
    MessageSquare,
    Flag,
    Image,
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
                            color={review.current_user_vote === 'helpful' ? 'primary' : 'inherit'}
                        >
                            Полезно ({review.votes_count.helpful})
                        </Button>
                        <Button
                            size="small"
                            startIcon={<ThumbsDown />}
                            onClick={() => onVote(review.id, 'not_helpful')}
                            color={review.current_user_vote === 'not_helpful' ? 'primary' : 'inherit'}
                        >
                            Не полезно ({review.votes_count.not_helpful})
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

// Компонент формы создания отзыва
const ReviewForm = ({ onSubmit, initialRating = 0 }) => {
    const [review, setReview] = useState({
        rating: initialRating,
        comment: '',
        pros: '',
        cons: '',
        photos: []
    });
    const [previewUrls, setPreviewUrls] = useState([]);

    const handleSubmit = (e) => {
        e.preventDefault();
        onSubmit(review);
    };

    const handlePhotoUpload = (e) => {
        const files = Array.from(e.target.files);
        
        // Проверяем ограничения
        if (files.length + review.photos.length > 10) {
            alert('Можно загрузить максимум 10 фотографий');
            return;
        }

        // Обрабатываем каждый файл
        const validFiles = files.filter(file => {
            if (!file.type.startsWith('image/')) {
                alert('Можно загружать только изображения');
                return false;
            }
            if (file.size > 5 * 1024 * 1024) {
                alert('Размер файла не должен превышать 5MB');
                return false;
            }
            return true;
        });

        setReview(prev => ({
            ...prev,
            photos: [...prev.photos, ...validFiles]
        }));

        // Создаем превью для каждого файла
        validFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreviewUrls(prev => [...prev, reader.result]);
            };
            reader.readAsDataURL(file);
        });
    };

    const removePhoto = (index) => {
        setReview(prev => ({
            ...prev,
            photos: prev.photos.filter((_, i) => i !== index)
        }));
        setPreviewUrls(prev => prev.filter((_, i) => i !== index));
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
            <Stack spacing={3}>
                {/* Рейтинг */}
                <Box>
                    <Typography gutterBottom>Общая оценка</Typography>
                    <Rating
                        value={review.rating}
                        onChange={(e, newValue) => setReview({ ...prev => ({ ...prev, rating: newValue })})}
                        size="large"
                        required
                    />
                </Box>

                {/* Комментарий */}
                <TextField
                    label="Комментарий"
                    multiline
                    rows={4}
                    value={review.comment}
                    onChange={(e) => setReview(prev => ({ ...prev, comment: e.target.value }))}
                    required
                />

                {/* Достоинства */}
                <TextField
                    label="Достоинства"
                    multiline
                    rows={2}
                    value={review.pros}
                    onChange={(e) => setReview(prev => ({ ...prev, pros: e.target.value }))}
                    placeholder="Что вам особенно понравилось?"
                />

                {/* Недостатки */}
                <TextField
                    label="Недостатки"
                    multiline
                    rows={2}
                    value={review.cons}
                    onChange={(e) => setReview(prev => ({ ...prev, cons: e.target.value }))}
                    placeholder="Что можно было бы улучшить?"
                />

                {/* Загрузка фотографий */}
                <Box>
                    <Button
                        variant="outlined"
                        component="label"
                        startIcon={<Image />}
                    >
                        Добавить фото
                        <input
                            type="file"
                            hidden
                            multiple
                            accept="image/*"
                            onChange={handlePhotoUpload}
                        />
                    </Button>
                    
                    {/* Превью фотографий */}
                    <Box sx={{ 
                        display: 'flex', 
                        gap: 1, 
                        flexWrap: 'wrap',
                        mt: 2 
                    }}>
                        {previewUrls.map((url, index) => (
                            <Box
                                key={index}
                                sx={{
                                    position: 'relative',
                                    width: 100,
                                    height: 100
                                }}
                            >
                                <img
                                    src={url}
                                    alt={`Preview ${index}`}
                                    style={{
                                        width: '100%',
                                        height: '100%',
                                        objectFit: 'cover',
                                        borderRadius: '4px'
                                    }}
                                />
                                <IconButton
                                    size="small"
                                    sx={{
                                        position: 'absolute',
                                        top: -10,
                                        right: -10,
                                        bgcolor: 'background.paper'
                                    }}
                                    onClick={() => removePhoto(index)}
                                >
                                    <X size={16} />
                                </IconButton>
                            </Box>
                        ))}
                    </Box>
                </Box>

                {/* Кнопка отправки */}
                <Button
                    type="submit"
                    variant="contained"
                    size="large"
                    disabled={!review.rating || !review.comment}
                >
                    Опубликовать отзыв
                </Button>
            </Stack>
        </Box>
    );
};

export default ReviewForm;