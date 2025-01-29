// frontend/hostel-frontend/src/components/notifications/NotificationSettings.js
import React, { useState, useEffect } from 'react';
import {
    Box,
    Typography,
    Switch,
    FormControlLabel,
    Paper,
    Stack,
    Button,
    Alert,
    Snackbar,
    Divider,
    Grid
} from '@mui/material';
import {
    MessageCircle,
    BellRing,
    FileText,
    Star,
    Tag,
    RefreshCw
} from 'lucide-react';
import { useNotifications } from '../../hooks/useNotifications';
import axios from '../../api/axios';

const NOTIFICATION_TYPES = {
    new_message: {
        label: 'Новые сообщения',
        icon: MessageCircle,
        description: 'Уведомления о новых сообщениях в чате',
        implemented: true
    },
    new_review: {
        label: 'Новые отзывы',
        icon: FileText,
        description: 'Уведомления о новых отзывах на ваши объявления',
        implemented: true
    },
    review_vote: {
        label: 'Оценка отзыва',
        icon: Star,
        description: 'Уведомления об оценках ваших отзывов',
        implemented: true
    },
    review_response: {
        label: 'Ответы на отзывы',
        icon: MessageCircle,
        description: 'Уведомления об ответах на ваши отзывы',
        implemented: true
    },
    listing_status: {
        label: 'Статус объявлений',
        icon: RefreshCw,
        description: 'Уведомления об изменениях статуса ваших объявлений',
        implemented: true
    },
    favorite_price: {
        label: 'Цены в избранном',
        icon: Tag,
        description: 'Уведомления об изменении цен в избранных объявлениях',
        implemented: true
    }
};

const NotificationSettings = () => {
    const {
        settings,
        telegramConnected,
        updateSettings,
        connectTelegram,
        fetchSettings
    } = useNotifications();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

    useEffect(() => {
        if (telegramConnected) {
            fetchSettings();
        }
    }, [telegramConnected, fetchSettings]);

    const showSnackbar = (message, severity = 'success') => {
        setSnackbar({ open: true, message, severity });
    };



    const handleTelegramConnect = async () => {
        try {
            setLoading(true);
            setError(null);

            // Добавляем больше логирования
            console.log('Initiating Telegram connection...');

            const response = await connectTelegram();
            console.log('Connect Telegram response:', response);

        } catch (err) {
            console.error('Telegram connection error:', err);
            setError(err.message || 'Ошибка подключения к Telegram');

            // Показываем пользователю более информативное сообщение
            showSnackbar('Не удалось подключить Telegram. Пожалуйста, попробуйте позже.', 'error');
        } finally {
            setLoading(false);
        }
    };

    const handleSettingChange = async (type, channel, value) => {
        // Проверяем, реализован ли данный тип уведомлений
        if (!NOTIFICATION_TYPES[type]?.implemented) {
            showSnackbar('Этот тип уведомлений пока недоступен', 'warning');
            return;
        }

        try {
            const success = await updateSettings(type, channel, value);
            if (success) {
                showSnackbar('Настройки успешно обновлены');
            } else {
                showSnackbar('Ошибка при обновлении настроек', 'error');
            }
        } catch (error) {
            showSnackbar('Ошибка при обновлении настроек', 'error');
        }
    };

    return (
        <Box>
            <Typography variant="h6" gutterBottom>
                Настройки уведомлений
            </Typography>

            <Paper sx={{ p: 3, mb: 3 }}>
                <Stack spacing={3}>
                    <Box>
                        <Typography variant="subtitle1" gutterBottom>
                            Каналы уведомлений
                        </Typography>
                        <Grid container spacing={2}>
                            <Grid item xs={12} sm={6} md={4}>
                                <Button
                                    variant={telegramConnected ? "outlined" : "contained"}
                                    onClick={handleTelegramConnect}
                                    startIcon={<MessageCircle />}
                                    disabled={loading}
                                    fullWidth
                                >
                                    {loading ? 'Подключение...' :
                                        telegramConnected ? 'Telegram подключен' :
                                            'Подключить Telegram'}
                                </Button>
                            </Grid>
                             <Grid item xs={12} sm={6} md={4}>
                                <Button
                                    variant="outlined"
                                    onClick={async () => {
                                        try {
                                            await axios.post('/api/v1/notifications/test');
                                            showSnackbar('Тестовое уведомление отправлено');
                                        } catch (err) {
                                            showSnackbar('Ошибка отправки уведомления', 'error');
                                        }
                                    }}
                                    fullWidth
                                >
                                    Отправить тестовое уведомление
                                </Button>
                            </Grid>
                        </Grid>
                    </Box>

                    <Divider />

                    <Stack spacing={2}>
                        {Object.entries(NOTIFICATION_TYPES).map(([type, { label, icon: Icon, description, implemented }]) => (
                            <Box key={type}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                                    <Icon size={20} />
                                    <Typography variant="subtitle2">{label}</Typography>
                                    {!implemented && (
                                        <Typography
                                            variant="caption"
                                            sx={{
                                                ml: 1,
                                                color: 'text.secondary',
                                                bgcolor: 'action.hover',
                                                px: 1,
                                                py: 0.5,
                                                borderRadius: 1
                                            }}
                                        >
                                            В разработке
                                        </Typography>
                                    )}
                                </Box>
                                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                                    {description}
                                </Typography>
                                <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                                    <FormControlLabel
                                        control={
                                            <Switch
                                                checked={Boolean(settings[type]?.telegram_enabled)}
                                                onChange={(e) => handleSettingChange(type, 'telegram', e.target.checked)}
                                                disabled={!telegramConnected || !implemented}
                                                color="primary"
                                            />
                                        }
                                        label="Telegram"
                                    />

 
                                </Stack>
                                {!implemented && (
                                    <Typography
                                        variant="caption"
                                        color="text.secondary"
                                        sx={{ display: 'block', mt: 1 }}
                                    >
                                        Этот тип уведомлений пока недоступен
                                    </Typography>
                                )}
                            </Box>
                        ))}
                    </Stack>
                </Stack>
            </Paper>

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

export default NotificationSettings;
