// frontend/hostel-frontend/src/components/reviews/ReviewComponents.jsx
import React, { useState } from 'react';
import { useMediaQuery } from '@mui/material';
import GalleryViewer from '../shared/GalleryViewer';
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
    CheckCircle2,
    PencilLine
} from 'lucide-react';

// Компонент формы создания/редактирования отзыва
const ReviewForm = ({ onSubmit, initialData = null, onCancel, entityType, entityId }) => {
    const [formData, setFormData] = useState({
        rating: initialData?.rating || 0,
        comment: initialData?.comment || '',
        pros: initialData?.pros || '',
        cons: initialData?.cons || ''
    });

    const [photoFiles, setPhotoFiles] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);

    const handleChange = (field) => (event) => {
        setFormData(prev => ({
            ...prev,
            [field]: event.target.value
        }));
    };

    const handlePhotoAdd = (event) => {
        const files = Array.from(event.target.files);

        const validFiles = files.filter(file => {
            const isValidType = file.type.startsWith('image/');
            const isValidSize = file.size <= 15 * 1024 * 1024; // 15MB
            return isValidType && isValidSize;
        });

        if (validFiles.length + photoFiles.length > 10) {
            alert('Можно загрузить максимум 10 фотографий');
            return;
        }

        setPhotoFiles(prev => [...prev, ...validFiles]);

        // Создаем URL для предпросмотра
        validFiles.forEach(file => {
            const url = URL.createObjectURL(file);
            setPreviewUrls(prev => [...prev, url]);
        });
    };

    const handleRemovePhoto = (index) => {
        setPhotoFiles(prev => prev.filter((_, i) => i !== index));
        setPreviewUrls(prev => prev.filter((_, i) => i !== index));
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        const reviewData = {
            entity_type: entityType,
            entity_id: entityId,
            rating: parseInt(formData.rating),
            comment: formData.comment,
            pros: formData.pros,
            cons: formData.cons
        };

        let photosFormData = null;
        if (photoFiles.length > 0) {
            photosFormData = new FormData();
            photoFiles.forEach(file => {
                photosFormData.append('photos', file);
            });
        }

        onSubmit({ reviewData, photosFormData });
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ p: 2 }}>
            <Stack spacing={3}>
                <Box>
                    <Typography gutterBottom>Общая оценка</Typography>
                    <Rating
                        value={formData.rating}
                        onChange={(_, newValue) => {
                            setFormData(prev => ({ ...prev, rating: newValue }));
                        }}
                        size="large"
                    />
                </Box>

                <TextField
                    label="Комментарий"
                    multiline
                    rows={4}
                    value={formData.comment}
                    onChange={handleChange('comment')}
                    fullWidth
                />

                <TextField
                    label="Достоинства"
                    multiline
                    rows={2}
                    value={formData.pros}
                    onChange={handleChange('pros')}
                    fullWidth
                />

                <TextField
                    label="Недостатки"
                    multiline
                    rows={2}
                    value={formData.cons}
                    onChange={handleChange('cons')}
                    fullWidth
                />

                <Box>
                    <Button
                        component="label"
                        startIcon={<Camera />}
                        variant="outlined"
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

                    {previewUrls.length > 0 && (
                        <Box sx={{ mt: 2 }}>
                            <GalleryViewer
                                images={previewUrls}
                                galleryMode="thumbnails"
                                thumbnailSize={{ width: '100%', height: '100px' }}
                                gridColumns={{ xs: 4, sm: 3, md: 2 }}
                            />
                        </Box>
                    )}
                </Box>

                <Stack direction="row" spacing={2} justifyContent="flex-end">
                    <Button onClick={onCancel}>
                        Отмена
                    </Button>
                    <Button
                        type="submit"
                        variant="contained"
                        disabled={!formData.rating || !formData.comment}
                    >
                        {initialData ? 'Сохранить' : 'Опубликовать'}
                    </Button>
                </Stack>
            </Stack>
        </Box>
    );
};

// Компонент отдельного отзыва
const ReviewCard = ({ review, currentUserId, onVote, onReply, onEdit, onDelete, onReport }) => {
    const [showGallery, setShowGallery] = useState(false);
    const [selectedImageIndex, setSelectedImageIndex] = useState(0);
    const [showReplyForm, setShowReplyForm] = useState(false);
    const [replyText, setReplyText] = useState('');
    const [menuAnchor, setMenuAnchor] = useState(null);

    const isMobile = useMediaQuery((theme) => theme.breakpoints.down('sm'));

    const handleReplySubmit = () => {
        onReply(review.id, replyText);
        setReplyText('');
        setShowReplyForm(false);
    };

    const handleVote = (voteType) => {
        // Проверяем авторизацию
        if (!currentUserId) {
            // Если пользователь не авторизован, ничего не делаем
            return;
        }

        // Если текущий голос такой же, снимаем его
        if (review.current_user_vote === voteType) {
            return;
        }

        onVote(review.id, voteType);
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
                                isMobile ? (
                                    <CheckCircle2 size={20} color="green" />
                                ) : (
                                    <Chip
                                        icon={<CheckCircle2 size={16} />}
                                        label="Проверенная покупка"
                                        size="small"
                                        color="success"
                                    />
                                )
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

                    {/* Галерея фотографий */}
                    {review.photos?.length > 0 && (
                        <>
                            {/* Превью фотографий */}
                            <Box sx={{ mt: 2 }}>
                                <GalleryViewer
                                    images={review.photos}
                                    galleryMode="thumbnails"
                                    thumbnailSize={{ width: '100%', height: '100px' }}
                                    gridColumns={{ xs: 4, sm: 3, md: 2 }}
                                    onClick={(index) => {
                                        setSelectedImageIndex(index);
                                        setShowGallery(true);
                                    }}
                                />
                            </Box>

                            {/* Полноэкранный просмотр */}
                            <GalleryViewer
                                images={review.photos}
                                open={showGallery}
                                onClose={() => setShowGallery(false)}
                                initialIndex={selectedImageIndex}
                                galleryMode="fullscreen"
                            />
                        </>
                    )}

                    {/* Кнопки голосования */}
                    <Stack direction="row" spacing={2}>
                        <Button
                            size="small"
                            onClick={() => handleVote('helpful')}
                            startIcon={<ThumbsUp />}
                            variant={review.current_user_vote === 'helpful' ? 'contained' : 'outlined'}
                            disabled={!currentUserId} // Добавляем disabled если пользователь не авторизован
                        >
                            {isMobile ? (
                                `(${review.votes_count?.helpful || 0})`
                            ) : (
                                `Полезно (${review.votes_count?.helpful || 0})`
                            )}
                        </Button>
                        <Button
                            size="small"
                            onClick={() => handleVote('not_helpful')}
                            startIcon={<ThumbsDown />}
                            variant={review.current_user_vote === 'not_helpful' ? 'contained' : 'outlined'}
                            disabled={!currentUserId} // Добавляем disabled если пользователь не авторизован
                        >
                            {isMobile ? (
                                `(${review.votes_count?.not_helpful || 0})`
                            ) : (
                                `Не полезно (${review.votes_count?.not_helpful || 0})`
                            )}
                        </Button>
                        {isMobile ? (
                            <>
                                <IconButton onClick={() => setShowReplyForm(!showReplyForm)}>
                                    <MessageSquare size={20} />
                                </IconButton>
                                <IconButton onClick={() => onReport(review.id)}>
                                    <Flag size={20} />
                                </IconButton>
                            </>
                        ) : (
                            <>
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
                            </>
                        )}
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