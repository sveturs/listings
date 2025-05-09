import React, { useState, ChangeEvent } from 'react';
import { useMediaQuery, Theme } from '@mui/material';
import GalleryViewer from '../shared/GalleryViewer';
import { useTranslation } from 'react-i18next';
import { Review } from './ReviewsSection';

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
    MenuItem,
    Collapse,
    Badge,
    Paper
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
    PencilLine,
    ChevronDown,
    ChevronUp
} from 'lucide-react';

interface ReviewFormData {
    rating: number;
    comment: string;
    pros: string;
    cons: string;
}

interface ReviewStats {
    average_rating: number;
    total_reviews: number;
    rating_distribution: Record<number, number>;
}

interface ReviewFormProps {
    onSubmit: (data: { reviewData: any; photosFormData: FormData | null }) => Promise<void>;
    initialData?: Review | null;
    onCancel: () => void;
    entityType: string;
    entityId: string | number;
}

interface ReviewCardProps {
    review: Review;
    currentUserId?: string | number;
    onVote: (reviewId: string | number, voteType: 'helpful' | 'not_helpful') => Promise<void>;
    onReply: (reviewId: string | number, text: string) => Promise<void>;
    onEdit: (review: Review) => void;
    onDelete: (reviewId: string | number) => Promise<void>;
    onReport: (reviewId: string | number) => Promise<void>;
}

interface RatingStatsProps {
    stats: ReviewStats;
}

// Компонент формы создания/редактирования отзыва
const ReviewForm: React.FC<ReviewFormProps> = ({ onSubmit, initialData = null, onCancel, entityType, entityId }) => {
    const { t, i18n } = useTranslation('marketplace');

    const [formData, setFormData] = useState<ReviewFormData>({
        rating: initialData?.rating || 0,
        comment: initialData?.comment || '',
        pros: initialData?.pros || '',
        cons: initialData?.cons || ''
    });
    const [showResponses, setShowResponses] = useState<boolean>(false);
    const [photoFiles, setPhotoFiles] = useState<File[]>([]);
    const [previewUrls, setPreviewUrls] = useState<string[]>([]);

    const handleChange = (field: keyof ReviewFormData) => (event: ChangeEvent<HTMLInputElement>) => {
        setFormData(prev => ({
            ...prev,
            [field]: event.target.value
        }));
    };

    const handlePhotoAdd = (event: ChangeEvent<HTMLInputElement>): void => {
        if (!event.target.files) return;

        const files = Array.from(event.target.files);

        const validFiles = files.filter(file => {
            // Явное приведение типа File
            const typedFile = file as File;
            const isValidType = typedFile.type.startsWith('image/');
            const isValidSize = typedFile.size <= 15 * 1024 * 1024; // 15MB
            return isValidType && isValidSize;
        });

        if (validFiles.length + photoFiles.length > 10) {
            alert(t('reviews.uploadupto10'));
            return;
        }

        setPhotoFiles(prev => [...prev, ...validFiles]);

        // Создаем URL для предпросмотра
        validFiles.forEach(file => {
            const url = URL.createObjectURL(file as Blob);
            setPreviewUrls(prev => [...prev, url]);
        });
    };

    const handleRemovePhoto = (index: number): void => {
        setPhotoFiles(prev => prev.filter((_, i) => i !== index));
        setPreviewUrls(prev => prev.filter((_, i) => i !== index));
    };

    const handleSubmit = (e: React.FormEvent): void => {
        e.preventDefault();

        const reviewData = {
            entity_type: entityType,
            entity_id: entityId,
            rating: parseInt(formData.rating.toString()),
            comment: formData.comment,
            pros: formData.pros,
            cons: formData.cons,
            original_language: i18n.language
        };

        let photosFormData: FormData | null = null;
        if (photoFiles.length > 0) {
            photosFormData = new FormData();
            photoFiles.forEach(file => {
                photosFormData!.append('photos', file);
            });
        }

        onSubmit({ reviewData, photosFormData });
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ p: 2 }}>
            <Stack spacing={3}>
                <Box>
                    <Typography gutterBottom>{t('reviews.rating')}</Typography>
                    <Rating
                        value={formData.rating}
                        onChange={(_, newValue) => {
                            setFormData(prev => ({ ...prev, rating: newValue || 0 }));
                        }}
                        size="large"
                    />
                </Box>

                <TextField
                    label={t('reviews.comment')}
                    multiline
                    rows={4}
                    value={formData.comment}
                    onChange={handleChange('comment')}
                    fullWidth
                />

                <TextField
                    label={t('reviews.pros')}
                    multiline
                    rows={2}
                    value={formData.pros}
                    onChange={handleChange('pros')}
                    fullWidth
                />

                <TextField
                    label={t('reviews.cons')}
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
                        {t('reviews.cancel')}
                    </Button>
                    <Button
                        type="submit"
                        variant="contained"
                        disabled={!formData.rating || !formData.comment}
                    >
                        {initialData ? t('reviews.save') : t('reviews.submit')}
                    </Button>
                </Stack>
            </Stack>
        </Box>
    );
};

// Компонент отдельного отзыва
const ReviewCard: React.FC<ReviewCardProps> = ({ 
    review, 
    currentUserId, 
    onVote, 
    onReply, 
    onEdit, 
    onDelete, 
    onReport 
}) => {
    const { t, i18n } = useTranslation('marketplace');
    const [showGallery, setShowGallery] = useState<boolean>(false);
    const [selectedImageIndex, setSelectedImageIndex] = useState<number>(0);
    const [showReplyForm, setShowReplyForm] = useState<boolean>(false);
    const [replyText, setReplyText] = useState<string>('');
    const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
    const [showResponses, setShowResponses] = useState<boolean>(false);
    const isMobile = useMediaQuery((theme: Theme) => theme.breakpoints.down('sm'));

    const handleReplySubmit = (): void => {
        onReply(review.id, replyText);
        setReplyText('');
        setShowReplyForm(false);
    };

    const handleVote = (voteType: 'helpful' | 'not_helpful'): void => {
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

    // В функции getTranslatedContent добавить проверку:
    const getTranslatedContent = (field: keyof Pick<Review, 'comment' | 'pros' | 'cons'>): string => {
        if (!review || !field) return '';

        if (i18n.language === review.original_language) {
            return review[field] || '';
        }

        // Проверяем, что translations существует
        if (!review.translations || !review.translations[i18n.language]) {
            return review[field] || '';
        }

        const translation = review.translations[i18n.language][field];
        if (translation) {
            return translation;
        }

        return review[field] || '';
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
                                        label={t('reviews.verifiedPurchase')}
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
                                        {t('reviews.edit')}
                                    </MenuItem>
                                    <MenuItem onClick={() => {
                                        onDelete(review.id);
                                        setMenuAnchor(null);
                                    }}>
                                        <Trash2 size={16} style={{ marginRight: 8 }} />
                                        {t('reviews.delete')}
                                    </MenuItem>
                                </Menu>
                            </>
                        )}
                    </Stack>

                    <Rating value={review.rating} readOnly />

                    {review.comment && (
                        <Typography>{getTranslatedContent('comment')}</Typography>
                    )}
                    {review.pros && (
                        <Box>
                            <Typography color="success.main" variant="subtitle2">
                                {t('reviews.pros')}
                            </Typography>
                            <Typography>{getTranslatedContent('pros')}</Typography>
                        </Box>
                    )}

                    {review.cons && (
                        <Box>
                            <Typography color="error.main" variant="subtitle2">
                                {t('reviews.cons')}
                            </Typography>
                            <Typography>{getTranslatedContent('cons')}</Typography>
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
                            disabled={!currentUserId}
                        >
                            {isMobile ? `(${review.votes_count?.helpful || 0})` :
                                `${t('reviews.helpful')} (${review.votes_count?.helpful || 0})`}
                        </Button>
                        <Button
                            size="small"
                            onClick={() => handleVote('not_helpful')}
                            startIcon={<ThumbsDown />}
                            variant={review.current_user_vote === 'not_helpful' ? 'contained' : 'outlined'}
                            disabled={!currentUserId}
                        >
                            {isMobile ? `(${review.votes_count?.not_helpful || 0})` :
                                `${t('reviews.notHelpful')} (${review.votes_count?.not_helpful || 0})`}
                        </Button>

                        {/* НОВОЕ: Кнопка ответов */}
                        {review.responses && review.responses.length > 0 && (
                            isMobile ? (
                                <IconButton onClick={() => setShowResponses(!showResponses)}>
                                    <Badge badgeContent={review.responses.length} color="primary">
                                        <MessageSquare size={20} />
                                    </Badge>
                                </IconButton>
                            ) : (
                                <Button
                                    size="small"
                                    startIcon={showResponses ? <ChevronUp /> : <ChevronDown />}
                                    onClick={() => setShowResponses(!showResponses)}
                                >
                                    {`${t('reviews.responses')} (${review.responses.length})`}
                                </Button>
                            )
                        )}

                        {/* Кнопка "Ответить" */}
                        {isMobile ? (
                            <IconButton onClick={() => setShowReplyForm(!showReplyForm)}>
                                <MessageSquare size={20} />
                            </IconButton>
                        ) : (
                            <Button
                                size="small"
                                startIcon={<MessageSquare />}
                                onClick={() => setShowReplyForm(!showReplyForm)}
                            >
                                {t('reviews.reply')}
                            </Button>
                        )}

                        {isMobile ? (
                            <IconButton onClick={() => onReport(review.id)}>
                                <Flag size={20} />
                            </IconButton>
                        ) : (
                            <Button
                                size="small"
                                startIcon={<Flag />}
                                onClick={() => onReport(review.id)}
                            >
                                {t('reviews.report')}
                            </Button>
                        )}
                    </Stack>

                    {/* НОВОЕ: Блок с ответами */}
                    <Collapse in={showResponses}>
                        {review.responses && review.responses.map((response, index) => (
                            <Box key={index} sx={{ mt: 2, pl: isMobile ? 2 : 4 }}>
                                <Paper sx={{ p: isMobile ? 1.5 : 2, bgcolor: 'grey.50' }}>
                                    <Stack direction="row" spacing={isMobile ? 1 : 2} alignItems="center">
                                        <Avatar
                                            src={response.user?.picture_url}
                                            sx={{ width: isMobile ? 20 : 24, height: isMobile ? 20 : 24 }}
                                        />
                                        <Typography variant={isMobile ? "body2" : "subtitle2"}>
                                            {response.user?.name}
                                        </Typography>
                                        <Typography variant="caption" color="text.secondary">
                                            {new Date(response.created_at).toLocaleDateString()}
                                        </Typography>
                                    </Stack>
                                    <Typography variant={isMobile ? "body2" : "body1"} sx={{ mt: 1 }}>
                                        {response.content}
                                    </Typography>
                                </Paper>
                            </Box>
                        ))}
                    </Collapse>

                    {/* Форма ответа */}
                    {showReplyForm && (
                        <Box>
                            <TextField
                                fullWidth
                                multiline
                                rows={3}
                                placeholder={t('reviews.writeResponse')}
                                value={replyText}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setReplyText(e.target.value)}
                            />
                            <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
                                <Button
                                    variant="contained"
                                    size="small"
                                    onClick={handleReplySubmit}
                                    disabled={!replyText.trim()}
                                >
                                    {t('reviews.reply')}
                                </Button>
                                <Button
                                    size="small"
                                    onClick={() => {
                                        setShowReplyForm(false);
                                        setReplyText('');
                                    }}
                                >
                                    {t('reviews.cancel')}
                                </Button>
                            </Stack>
                        </Box>
                    )}
                </Stack>
            </CardContent>
        </Card>
    );
};


// Компонент статистики рейтингов
const RatingStats: React.FC<RatingStatsProps> = ({ stats }) => {
    const { t } = useTranslation('marketplace');
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
                    {t('listings.details.info.reviews.count', { count: stats.total_reviews || 0 })}
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