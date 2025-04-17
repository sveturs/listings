import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { QRCodeSVG } from 'qrcode.react';
import { Check } from 'lucide-react';
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
import { useNotifications } from '../../hooks/useNotifications';
import axios from '../../api/axios';
import { MessageCircle, FileText, Star, Tag, RefreshCw, QrCode, Mail } from 'lucide-react';

const NotificationSettings = () => {
    const { t } = useTranslation('marketplace');
    const [qrToken, setQrToken] = useState('');
    
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

    const NOTIFICATION_TYPES = {
        new_message: {
            label: t('notifications.types.newMessage'),
            icon: MessageCircle,
            description: t('notifications.types.newMessageDescription'),
            implemented: true
        },
        new_review: {
            label: t('notifications.types.newReview'),
            icon: FileText,
            description: t('notifications.types.newReviewDescription'),
            implemented: true
        },
        review_vote: {
            label: t('notifications.types.reviewVote'),
            icon: Star,
            description: t('notifications.types.reviewVoteDescription'),
            implemented: true
        },
        review_response: {
            label: t('notifications.types.reviewResponse'),
            icon: MessageCircle,
            description: t('notifications.types.reviewResponseDescription'),
            implemented: true
        },
        listing_status: {
            label: t('notifications.types.listingStatus'),
            icon: RefreshCw,
            description: t('notifications.types.listingStatusDescription'),
            implemented: true
        },
        favorite_price: {
            label: t('notifications.types.favoritePrice'),
            icon: Tag,
            description: t('notifications.types.favoritePriceDescription'),
            implemented: true
        }
    };

    useEffect(() => {
        if (telegramConnected) {
            fetchSettings();
        }
    }, [telegramConnected, fetchSettings]);
    
    useEffect(() => {
        const fetchQrToken = async () => {
            try {
                const response = await axios.get('/api/v1/notifications/telegram/token');
                if (response.data?.data?.token) {
                    setQrToken(`https://t.me/SveTu_bot?start=${response.data.data.token}`);
                }
            } catch (err) {
                console.error('Error fetching QR token:', err);
            }
        };
    
        if (!telegramConnected) {
            fetchQrToken();
        }
    }, [telegramConnected]);

    const showSnackbar = (message, severity = 'success') => {
        setSnackbar({ open: true, message, severity });
    };

    const handleTelegramConnect = async () => {
        try {
            setLoading(true);
            setError(null);
            await connectTelegram();
        } catch (err) {
            setError(err.message || t('notifications.telegram.error'));
            showSnackbar(t('notifications.telegram.error'), 'error');
        } finally {
            setLoading(false);
        }
    };

    const handleSettingChange = async (type, channel, value) => {
        if (!NOTIFICATION_TYPES[type]?.implemented) {
            showSnackbar(t('notifications.inDevelopment'), 'warning');
            return;
        }

        try {
            const success = await updateSettings(type, channel, value);
            if (success) {
                showSnackbar(t('notifications.settingsUpdated'));
            } else {
                showSnackbar(t('notifications.updateError'), 'error');
            }
        } catch (error) {
            showSnackbar(t('notifications.updateError'), 'error');
        }
    };

    return (
<Box>
    <Typography variant="h6" gutterBottom>
        {t('notifications.title')}
    </Typography>

    <Paper sx={{ p: 3, mb: 3 }}>
        <Stack spacing={3}>
            <Box>
                {/* Заменяем Tabs на единый блок подключения */}
                <Box sx={{ 
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    gap: 3,
                    mb: 3
                }}>
                    {telegramConnected ? (
                        <Alert 
                            icon={<Check />}
                            severity="success"
                            sx={{ width: '100%' }}
                        >
                            {t('notifications.telegram.connected')}
                        </Alert>
                    ) : (
                        <>
                            <Typography variant="body1" color="text.secondary" align="center">
                                {t('notifications.telegram.description')}
                            </Typography>
                            
                            <Stack 
                                direction={{ xs: 'column', sm: 'row' }} 
                                spacing={2}
                                alignItems="center"
                            >
                                <Button
                                    variant="contained"
                                    onClick={() => window.open(qrToken, '_blank')}
                                    startIcon={<QrCode />}
                                    disabled={loading || !qrToken}
                                >
                                    {t('notifications.telegram.scanQr')}
                                </Button>

                                <Typography color="text.secondary">
                                    {t('notifications.telegram.or')}
                                </Typography>

                                <Button
                                    variant="contained"
                                    onClick={handleTelegramConnect}
                                    startIcon={<MessageCircle />}
                                    disabled={loading}
                                >
                                    {loading ? t('notifications.telegram.connecting') : 
                                        t('notifications.telegram.connect')}
                                </Button>
                            </Stack>

                            {qrToken && (
                                <Box 
                                    sx={{ 
                                        p: 3, 
                                        bgcolor: 'background.paper', 
                                        borderRadius: 1,
                                        border: '1px solid',
                                        borderColor: 'divider',
                                        cursor: 'pointer'
                                    }}
                                    onClick={() => window.open(qrToken, '_blank')}
                                    title={t('notifications.telegram.clickToOpen')}
                                >
                                    <QRCodeSVG 
                                        value={qrToken}
                                        size={200}
                                        level="H"
                                        includeMargin
                                    />
                                </Box>
                            )}
                        </>
                    )}
                </Box>

                <Divider />

                <Stack spacing={2} sx={{ mt: 3 }}>
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
                                            {t('notifications.inDevelopment')}
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
                                    <FormControlLabel
                                        control={
                                            <Switch
                                                checked={Boolean(settings[type]?.email_enabled)}
                                                onChange={(e) => handleSettingChange(type, 'email', e.target.checked)}
                                                disabled={!implemented}
                                                color="primary"
                                            />
                                        }
                                        label="Email"
                                    />
                                </Stack>
                            </Box>
                        ))}
                    </Stack>

            </Box>
        </Stack>
    </Paper>
</Box>
    );
};

export default NotificationSettings;