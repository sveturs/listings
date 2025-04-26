// frontend/hostel-frontend/src/components/marketplace/ImageEnhancementOffer.js
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
    Box, Button, Typography, Card, CardContent,
    Grid, CircularProgress, Dialog,
    Divider, Alert, useTheme, useMediaQuery
} from '@mui/material';
import { Sparkles, Check } from 'lucide-react';
import axios from '../../api/axios';

const ImageEnhancementOffer = ({ images, onEnhanced }) => {
    const { t } = useTranslation('marketplace');
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    const [processing, setProcessing] = useState(false);
    const [previewOpen, setPreviewOpen] = useState(false);
    const [previewData, setPreviewData] = useState(null);

    const handleRequestPreview = async (imageIndex) => {
        if (images.length === 0 || !images[imageIndex]) return;

        setProcessing(true);

        try {
            const formData = new FormData();
            formData.append('image', images[imageIndex].file);

            const response = await axios.post('/api/v1/marketplace/enhance-preview', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            });

            setPreviewData({
                original: response.data.data.original || images[imageIndex].preview,
                enhanced: response.data.data.enhanced,
                price: response.data.data.price,
                index: imageIndex
            });

            setPreviewOpen(true);
        } catch (error) {
            console.error('Error getting enhancement preview:', error);
        } finally {
            setProcessing(false);
        }
    };

    const handleEnhanceAll = async () => {
        setProcessing(true);

        try {
            // Создаем FormData с множеством изображений
            const formData = new FormData();
            images.forEach((img) => {
                formData.append('images', img.file);
            });

            const response = await axios.post('/api/v1/marketplace/enhance-images', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            });

            // Получаем улучшенные изображения и передаем обратно в родительский компонент
            if (response.data.data.success && response.data.data.enhanced_images) {
                // Преобразуем формат данных для родительского компонента
                const enhancedImages = response.data.data.enhanced_images.map(img => ({
                    file: img.url, // URL к улучшенному изображению
                    preview: img.url,
                    cloudinary_id: img.public_id
                }));

                onEnhanced(enhancedImages);

                // Закрываем диалог предпросмотра
                setPreviewOpen(false);
            }
        } catch (error) {
            console.error('Error enhancing images:', error);

            // Если ошибка из-за недостатка средств, можно показать соответствующее сообщение
            if (error.response?.status === 402) {
                alert(t('listings.create.photos.enhance.insufficientFunds'));
            }
        } finally {
            setProcessing(false);
        }
    };

    if (images.length === 0) {
        return null;
    }

    return (
        <Box sx={{ mt: 3, mb: 3 }}>
            <Card variant="outlined">
                <CardContent>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                        <Sparkles size={20} color="#f59e0b" />
                        <Typography variant="h6" sx={{ ml: 1 }}>
                            {t('listings.create.photos.enhance.title', {
                                defaultValue: 'Улучшить фото профессионально'
                            })}
                        </Typography>
                    </Box>

                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                        {t('listings.create.photos.enhance.description', {
                            defaultValue: 'Автоматически улучшите качество ваших фотографий для привлечения большего числа покупателей. Профессиональные фотографии увеличивают шансы на продажу в 3 раза!'
                        })}
                    </Typography>

                    <Grid container spacing={2} sx={{ mb: 2 }}>
                        <Grid item xs={12} sm={6}>
                            <Box>
                                <Typography variant="subtitle2" gutterBottom>
                                    {t('listings.create.photos.enhance.benefits', {
                                        defaultValue: 'Преимущества'
                                    })}:
                                </Typography>
                                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                                    {[
                                        t('listings.create.photos.enhance.benefit1', {defaultValue: 'Улучшение яркости и контраста'}),
                                        t('listings.create.photos.enhance.benefit2', {defaultValue: 'Профессиональное качество фотографий'}),
                                        t('listings.create.photos.enhance.benefit3', {defaultValue: 'В 3 раза больше просмотров и откликов'})
                                    ].map((benefit, idx) => (
                                        <Box key={idx} sx={{ display: 'flex', alignItems: 'center' }}>
                                            <Check size={16} color="green" />
                                            <Typography variant="body2" sx={{ ml: 1 }}>
                                                {benefit}
                                            </Typography>
                                        </Box>
                                    ))}
                                </Box>
                            </Box>
                        </Grid>
                    </Grid>

                    <Divider sx={{ my: 2 }} />

                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
                        <Box>
                            <Typography variant="subtitle1">
                                {t('listings.create.photos.enhance.priceLabel', {defaultValue: 'Стоимость'})}:
                                <Typography
                                    component="span"
                                    variant="subtitle1"
                                    color="primary.main"
                                    sx={{ fontWeight: 'bold', ml: 1 }}
                                >
                                    {t('listings.create.photos.enhance.price', {
                                        count: images.length,
                                        pricePerImage: 30, // Цена за одно изображение
                                        defaultValue: '{{count}} фото x 30 RSD = {{totalPrice}} RSD',
                                        totalPrice: images.length * 30
                                    })}
                                </Typography>
                            </Typography>
                        </Box>
                        <Box sx={{ display: 'flex', gap: 2 }}>
                            <Button
                                variant="outlined"
                                onClick={() => handleRequestPreview(0)}
                                disabled={processing}
                            >
                                {t('listings.create.photos.enhance.preview', {
                                    defaultValue: 'Предпросмотр'
                                })}
                            </Button>
                            <Button
                                variant="contained"
                                startIcon={processing ? <CircularProgress size={20} /> : <Sparkles size={20} />}
                                onClick={handleEnhanceAll}
                                disabled={processing}
                            >
                                {processing
                                    ? t('listings.create.photos.enhance.processing', {defaultValue: 'Обработка...'})
                                    : t('listings.create.photos.enhance.enhanceAll', {defaultValue: 'Улучшить все фото'})
                                }
                            </Button>
                        </Box>
                    </Box>
                </CardContent>
            </Card>

            {/* Диалог с предпросмотром улучшенного изображения */}
            <Dialog
                open={previewOpen}
                onClose={() => setPreviewOpen(false)}
                maxWidth="md"
                fullWidth
            >
                {previewData && (
                    <Box sx={{ p: 3 }}>
                        <Typography variant="h6" gutterBottom>
                            {t('listings.create.photos.enhance.previewTitle', {
                                defaultValue: 'Сравнение до и после улучшения'
                            })}
                        </Typography>

                        <Grid container spacing={2} sx={{ mt: 1 }}>
                            <Grid item xs={12} sm={6}>
                                <Typography variant="subtitle2" gutterBottom align="center">
                                    {t('listings.create.photos.enhance.original', {
                                        defaultValue: 'Оригинал'
                                    })}
                                </Typography>
                                <Box sx={{
                                    width: '100%',
                                    height: 300,
                                    backgroundImage: `url(${previewData.original})`,
                                    backgroundSize: 'contain',
                                    backgroundRepeat: 'no-repeat',
                                    backgroundPosition: 'center',
                                    border: '1px solid #eee'
                                }} />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <Typography variant="subtitle2" gutterBottom align="center">
                                    {t('listings.create.photos.enhance.enhanced', {
                                        defaultValue: 'Улучшенное фото'
                                    })}
                                </Typography>
                                <Box sx={{
                                    width: '100%',
                                    height: 300,
                                    backgroundImage: `url(${previewData.enhanced})`,
                                    backgroundSize: 'contain',
                                    backgroundRepeat: 'no-repeat',
                                    backgroundPosition: 'center',
                                    border: '1px solid #eee'
                                }} />
                            </Grid>
                        </Grid>

                        <Alert severity="info" sx={{ mt: 3, mb: 2 }}>
                            {t('listings.create.photos.enhance.previewInfo', {
                                defaultValue: 'Профессиональное улучшение фотографий поможет привлечь больше внимания к вашему товару и ускорит продажу.'
                            })}
                        </Alert>

                        <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between' }}>
                            <Button
                                variant="outlined"
                                onClick={() => setPreviewOpen(false)}
                            >
                                {t('common.close', {defaultValue: 'Закрыть'})}
                            </Button>
                            <Button
                                variant="contained"
                                startIcon={<Sparkles size={20} />}
                                onClick={handleEnhanceAll}
                                disabled={processing}
                            >
                                {t('listings.create.photos.enhance.enhanceAll', {
                                    defaultValue: 'Улучшить все фото ({{price}} RSD)',
                                    price: images.length * 30
                                })}
                            </Button>
                        </Box>
                    </Box>
                )}
            </Dialog>
        </Box>
    );
};

export default ImageEnhancementOffer;