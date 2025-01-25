import React, { useState } from 'react';
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
    Divider
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
import { urlBase64ToUint8Array } from '../../utils/webPush';

const NOTIFICATION_TYPES = {
    new_message: {
        label: 'Новые сообщения',
        icon: MessageCircle,
        description: 'Уведомления о новых сообщениях в чате'
    },
    new_review: {
        label: 'Новые отзывы',
        icon: FileText,
        description: 'Уведомления о новых отзывах на ваши объявления'
    },
    review_vote: {
        label: 'Оценка отзыва',
        icon: Star,
        description: 'Уведомления об оценках ваших отзывов'
    },
    review_response: {
        label: 'Ответы на отзывы',
        icon: MessageCircle,
        description: 'Уведомления об ответах на ваши отзывы'
    },
    listing_status: {
        label: 'Статус объявлений',
        icon: RefreshCw,
        description: 'Уведомления об изменениях статуса ваших объявлений'
    },
    favorite_price: {
        label: 'Цены в избранном',
        icon: Tag,
        description: 'Уведомления об изменении цен в избранных объявлениях'
    }
};

const NotificationSettings = () => {
    const {
        settings,
        telegramConnected,
        updateSettings,
        connectTelegram
    } = useNotifications();

    const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

    const showSnackbar = (message, severity = 'success') => {
        setSnackbar({ open: true, message, severity });
    };

    const handlePushSubscription = async () => {
        try {
            const permission = await Notification.requestPermission();
            if (permission === 'granted') {
                const registration = await navigator.serviceWorker.register('/service-worker.js');
                const convertedKey = urlBase64ToUint8Array(process.env.REACT_APP_VAPID_PUBLIC_KEY); // Добавить функцию конвертации
                const subscription = await registration.pushManager.subscribe({
                    userVisibleOnly: true,
                    applicationServerKey: convertedKey
                });
                await axios.post('/api/v1/notifications/push/subscribe', subscription);
            }
        } catch (err) {
            console.error('Error enabling push notifications:', err);
        }
    };
};

const handleSettingChange = async (type, channel, value) => {
    try {
        await updateSettings(type, channel, value);
        showSnackbar('Настройки успешно обновлены');
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
                    <Stack direction="row" spacing={2}>
                        <Button
                            variant={telegramConnected ? "outlined" : "contained"}
                            onClick={connectTelegram}
                            startIcon={<MessageCircle />}
                        >
                            {telegramConnected ? 'Telegram подключен' : 'Подключить Telegram'}
                        </Button>

                        <Button
                            variant="contained"
                            onClick={handlePushSubscription}
                            startIcon={<BellRing />}
                        >
                            Включить Push-уведомления
                        </Button>
                    </Stack>
                </Box>

                <Divider />

                <Stack spacing={2}>
                    {Object.entries(NOTIFICATION_TYPES).map(([type, { label, icon: Icon, description }]) => (
                        <Box key={type}>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                                <Icon size={20} />
                                <Typography variant="subtitle2">{label}</Typography>
                            </Box>
                            <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                                {description}
                            </Typography>
                            <Stack direction="row" spacing={2}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings[type]?.telegram || false}
                                            onChange={(e) => handleSettingChange(type, 'telegram', e.target.checked)}
                                            disabled={!telegramConnected}
                                            color="primary"
                                        />
                                    }
                                    label="Telegram"
                                />
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={settings[type]?.push || false}
                                            onChange={(e) => handleSettingChange(type, 'push', e.target.checked)}
                                            color="primary"
                                        />
                                    }
                                    label="Push"
                                />
                            </Stack>
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