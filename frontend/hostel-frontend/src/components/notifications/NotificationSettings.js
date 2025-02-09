import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { QRCodeSVG } from 'qrcode.react';
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
    Grid,
    Tabs,
    Tab
} from '@mui/material';
import {
    MessageCircle,
    BellRing,
    FileText,
    Star,
    Tag,
    RefreshCw,
    QrCode
} from 'lucide-react';
import { useNotifications } from '../../hooks/useNotifications';
import axios from '../../api/axios';

const NotificationSettings = () => {
    const { t } = useTranslation('marketplace');
    const [activeTab, setActiveTab] = useState(0);
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
                        <Tabs value={activeTab} onChange={(e, v) => setActiveTab(v)}>
                            <Tab label={t('notifications.telegram.connect')} />
                            <Tab label={t('notifications.telegram.qrcode')} />
                        </Tabs>

                        <Box sx={{ mt: 2 }}>
                            {activeTab === 0 ? (
                                <Grid container spacing={2}>
                                    <Grid item xs={12} sm={6} md={4}>
                                        <Button
                                            variant={telegramConnected ? "outlined" : "contained"}
                                            onClick={handleTelegramConnect}
                                            startIcon={<MessageCircle />}
                                            disabled={loading}
                                            fullWidth
                                        >
                                            {loading ? t('notifications.telegram.connecting') :
                                                telegramConnected ? t('notifications.telegram.connected') :
                                                    t('notifications.telegram.connect')}
                                        </Button>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={4}>
                                        <Button
                                            variant="outlined"
                                            onClick={async () => {
                                                try {
                                                    await axios.post('/api/v1/notifications/test');
                                                    showSnackbar(t('notifications.testSent'));
                                                } catch (err) {
                                                    showSnackbar(t('notifications.testError'), 'error');
                                                }
                                            }}
                                            fullWidth
                                        >
                                            {t('notifications.test')}
                                        </Button>
                                    </Grid>
                                </Grid>
                            ) : (
                                <Box sx={{ 
                                    display: 'flex', 
                                    flexDirection: 'column', 
                                    alignItems: 'center',
                                    gap: 2 
                                }}>
                                    <Typography variant="body1" color="text.secondary">
                                        {t('notifications.telegram.scanQr')}
                                    </Typography>
                                    {qrToken && (
                                        <Box 
                                            sx={{ 
                                                p: 3, 
                                                bgcolor: 'white', 
                                                borderRadius: 1,
                                                cursor: 'pointer'
                                            }}
                                            onClick={() => window.open(qrToken, '_blank')}
                                        >
                                            <QRCodeSVG 
                                                value={qrToken}
                                                size={200}
                                                level="H"
                                                includeMargin
                                            />
                                        </Box>
                                    )}
                                </Box>
                            )}
                        </Box>
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