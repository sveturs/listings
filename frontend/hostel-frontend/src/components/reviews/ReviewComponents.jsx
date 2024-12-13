//frontend/hostel-frontend/src/components/reviews/ReviewComponents.jsx
import React, { useState } from 'react';
import {
    Box,
    Typography,
    Rating,
    Button,
    Card,
    CardContent,
    Stack,
    Avatar,
    TextField,
    LinearProgress,
    IconButton,
    Chip,
    Menu,
    MenuItem
} from '@mui/material';
import { 
    ThumbsUp, 
    ThumbsDown,
    MessageSquare,
    MoreVertical,
    Camera,
    Flag,
    Edit,
    Trash2,
    CheckCircle2
} from 'lucide-react';

// Компонент формы создания/редактирования отзыва
const ReviewForm = ({ onSubmit, initialData = null, onCancel }) => {
    const [formData, setFormData] = useState({
        rating: initialData?.rating || 0,
        comment: initialData?.comment || '',
        pros: initialData?.pros || '',
        cons: initialData?.cons || '',
        photos: initialData?.photos || []
    });
    const [photoFiles, setPhotoFiles] = useState([]);

    const handlePhotoAdd = (event) => {
        const files = Array.from(event.target.files);
        if (files.length + photoFiles.length > 10) {
            alert('Максимум 10 фотографий');
            return;
        }
        setPhotoFiles(prev => [...prev, ...files]);
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        const formDataToSend = new FormData();
        Object.keys(formData).forEach(key => {
            if (key !== 'photos') {
                formDataToSend.append(key, formData[key]);
            }
        });
        photoFiles.forEach(file => {
            formDataToSend.append('photos', file);
        });
        onSubmit(formDataToSend);
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ p: 2 }}>
            <Stack spacing={3}>
                <Box>
                    <Typography gutterBottom>Общая оценка</Typography>
                    <Rating
                        value={formData.rating}
                        onChange={(_, value) => setFormData(prev => ({
                            ...prev,
                            rating: value
                        }))}
                        size="large"
                    />
                </Box>

                <TextField
                    label="Комментарий"
                    multiline
                    rows={4}
                    value={formData.comment}
                    onChange={(e) => setFormData(prev => ({
                        ...prev,
                        comment: e.target.value
                    }))}
                />

                <TextField
                    label="Достоинства"
                    multiline
                    rows={2}
                    value={formData.pros}
                    onChange={(e) => setFormData(prev => ({
                        ...prev,
                        pros: e.target.value
                    }))}
                />

                <TextField
                    label="Недостатки"
                    multiline
                    rows={2}
                    value={formData.cons}
                    onChange={(e) => setFormData(prev => ({
                        ...prev,
                        cons: e.target.value
                    }))}
                />

                <Box>
                    <Button
                        variant="outlined"
                        component="label"
                        startIcon={<Camera />}
                    >
                        Добавить фото
                        <input
                            type="file"
                            hidden
                            multiple
                            accept="image/*"
                            onChange={handlePhotoAdd}
                        />
                    </Button>
                    <Stack direction="row" spacing={1} sx={{ mt: 2 }}>
                        {photoFiles.map((file, index) => (
                            <Box
                                key={index}
                                sx={{
                                    width: 100,
                                    height: 100,
                                    position: 'relative'
                                }}
                            >
                                <img
                                    src={URL.createObjectURL(file)}
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
                                    onClick={() => setPhotoFiles(prev => 
                                        prev.filter((_, i) => i !== index)
                                    )}
                                >
                                    <Trash2 size={16} />
                                </IconButton>
                            </Box>
                        ))}
                    </Stack>
                </Box>

                <Stack direction="row" spacing={2}>
                    <Button
                        type="submit"
                        variant="contained"
                        disabled={!formData.rating || !formData.comment}
                        fullWidth
                    >
                        {initialData ? 'Сохранить изменения' : 'Опубликовать отзыв'}
                    </Button>
                    {onCancel && (
                        <Button
                            variant="outlined"
                            onClick={onCancel}
                            fullWidth
                        >
                            Отмена
                        </Button>
                    )}
                </Stack>
            </Stack>
        </Box>
    );
};

// Компонент отдельного отзыва
const ReviewCard = ({ review, currentUserId, onVote, onReply, onEdit, onDelete, onReport }) => {
    const [showReplyForm, setShowReplyForm] = useState(false);
    const [replyText, setReplyText] = useState('');
    const [menuAnchor, setMenuAnchor] = useState(null);

    const handleReplySubmit = () => {
        onReply(review.id, replyText);
        setReplyText('');
        setShowReplyForm(false);
    };

    return (
        <Card sx={{ mb: 2 }}>
            <CardContent>
                <Stack spacing={2}>
                    {/* Заголовок с информацией о пользователе */}
                    <Stack 
                        direction="row" 
                        alignItems="center" 
                        justifyContent="space-between"
                    >
                        <Stack direction="row" spacing={2} alignItems="center">
                            <Avatar src={review.user?.picture_url} />
                            <Box>
                                <Typography variant="subtitle1">
                                    {review.user?.name}
                                </Typography>
                                <Typography variant="caption" color="text.secondary">
                                    {new Date(review.created_at).toLocaleDateString()}
                                </Typography>
                            </Box>
                            {review.is_verified_purchase && (
                                <Chip
                                    icon={<CheckCircle2 size={16} />}
                                    label="Проверенная покупка"
                                    size="small"
                                    color="success"
                                />
                            )}
                        </Stack>
                        
                        {(currentUserId === review.user_id) && (
                            <>
                                <IconButton onClick={(e) => setMenuAnchor(e.currentTarget)}>
                                    <MoreVertical />
                                </IconButton>
                                <Menu
                                    anchorEl={menuAnchor}
                                    open={Boolean(menuAnchor)}
                                    onClose={() => setMenuAnchor(null)}
                                >
                                    <MenuItem onClick={() => {
                                        onEdit(review);
                                        setMenuAnchor(null);
                                    }}>
                                        <Edit size={16} style={{ marginRight: 8 }} />
                                        Редактировать
                                    </MenuItem>
                                    <MenuItem onClick={() => {
                                        onDelete(review.id);
                                        setMenuAnchor(null);
                                    }}>
                                        <Trash2 size={16} style={{ marginRight: 8 }} />
                                        Удалить
                                    </MenuItem>
                                </Menu>
                            </>
                        )}
                    </Stack>

                    <Rating value={review.rating} readOnly />

                    {review.comment && (
                        <Typography>{review.comment}</Typography>
                    )}

                    {review.pros && (
                        <Box>
                            <Typography color="success.main" variant="subtitle2">
                                Достоинства:
                            </Typography>
                            <Typography>{review.pros}</Typography>
                        </Box>
                    )}

                    {review.cons && (
                        <Box>
                            <Typography color="error.main" variant="subtitle2">
                                Недостатки:
                            </Typography>
                            <Typography>{review.cons}</Typography>
                        </Box>
                    )}

                    {review.photos && review.photos.length > 0 && (
                        <Stack direction="row" spacing={1} sx={{ overflowX: 'auto' }}>
                            {review.photos.map((photo, index) => (
                                <img
                                    key={index}
                                    src={photo}
                                    alt={`Review ${index + 1}`}
                                    style={{
                                        width: 100,
                                        height: 100,
                                        objectFit: 'cover',
                                        borderRadius: '4px'
                                    }}
                                />
                            ))}
                        </Stack>
                    )}

                    <Stack direction="row" spacing={2}>
                        <Button
                            size="small"
                            startIcon={<ThumbsUp />}
                            onClick={() => onVote(review.id, 'helpful')}
                            color={review.current_user_vote === 'helpful' ? 'primary' : 'inherit'}
                        >
                            Полезно ({review.votes_count?.helpful || 0})
                        </Button>
                        <Button
                            size="small"
                            startIcon={<ThumbsDown />}
                            onClick={() => onVote(review.id, 'not_helpful')}
                            color={review.current_user_vote === 'not_helpful' ? 'primary' : 'inherit'}
                        >
                            Не полезно ({review.votes_count?.not_helpful || 0})
                        </Button>
                        <Button
                            size="small"
                            startIcon={<MessageSquare />}
                            onClick={() => setShowReplyForm(!showReplyForm)}
                        >
                            Ответить
                        </Button>
                        <Button
                            size="small"
                            startIcon={<Flag />}
                            onClick={() => onReport(review.id)}
                        >
                            Пожаловаться
                        </Button>
                    </Stack>

                    {showReplyForm && (
                        <Box>
                            <TextField
                                fullWidth
                                multiline
                                rows={3}
                                placeholder="Напишите ответ..."
                                value={replyText}
                                onChange={(e) => setReplyText(e.target.value)}
                            />
                            <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
                                <Button
                                    variant="contained"
                                    size="small"
                                    onClick={handleReplySubmit}
                                    disabled={!replyText.trim()}
                                >
                                    Ответить
                                </Button>
                                <Button
                                    size="small"
                                    onClick={() => {
                                        setShowReplyForm(false);
                                        setReplyText('');
                                    }}
                                >
                                    Отмена
                                </Button>
                            </Stack>
                        </Box>
                    )}

                    {review.responses && review.responses.length > 0 && (
                        <Box sx={{ pl: 4 }}>
                            {review.responses.map((response, index) => (
                                <Box key={index} sx={{ mt: 2 }}>
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
                                    <Typography sx={{ mt: 1 }}>
                                        {response.response}
                                    </Typography>
                                </Box>
                            ))}
                        </Box>
                    )}
                </Stack>
            </CardContent>
        </Card>
    );
};

// Компонент статистики рейтингов
const RatingStats = ({ stats }) => {
    return (
        <Stack direction="row" spacing={4} alignItems="center" sx={{ mb: 4 }}>
            <Box textAlign="center">
                <Typography variant="h3" fontWeight="bold">
                    {stats.average_rating?.toFixed(1) || "0.0"}
                </Typography>
                <Rating 
                    value={stats.average_rating || 0} 
                    readOnly 
                    precision={0.1} 
                />
                <Typography color="text.secondary">
                    {stats.total_reviews || 0} отзывов
                </Typography>
            </Box>

            <Box flex={1}>
                {[5, 4, 3, 2, 1].map(rating => (
                    <Stack 
                        key={rating} 
                        direction="row" 
                        spacing={2} 
                        alignItems="center" 
                        sx={{ mb: 1 }}
                    >
                        <Typography minWidth={20}>{rating}</Typography>
                        <LinearProgress
                            variant="determinate"
                            value={((stats.rating_distribution?.[rating] || 0) / 
                                (stats.total_reviews || 1)) * 100}
                            sx={{ flex: 1, height: 8, borderRadius: 1 }}
                        />
                        <Typography minWidth={40}>
                            {stats.rating_distribution?.[rating] || 0}
                        </Typography>
                    </Stack>
                ))}
            </Box>
        </Stack>
    );
};

export { ReviewForm, ReviewCard, RatingStats };