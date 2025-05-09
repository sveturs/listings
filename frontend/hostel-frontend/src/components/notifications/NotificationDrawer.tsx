// frontend/hostel-frontend/src/components/notifications/NotificationDrawer.tsx
import React, { useState, useEffect } from 'react';
import {
    Drawer,
    Box,
    Typography,
    IconButton,
    List,
    ListItem,
    ListItemText,
    Divider,
    CircularProgress,
    Alert
} from '@mui/material';
import { X, MessageCircle, Star, FileText, Tag } from 'lucide-react';
import { useNotifications } from '../../hooks/useNotifications';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import axios from '../../api/axios';

interface NotificationDrawerProps {
    open: boolean;
    onClose: () => void;
}

interface NotificationData {
    listing_id?: string | number;
    [key: string]: any;
}

interface Notification {
    id: string | number;
    title: string;
    message: string;
    type: keyof typeof NOTIFICATION_ICONS;
    is_read: boolean;
    created_at: string;
    data?: NotificationData;
}

const NOTIFICATION_ICONS = {
    new_message: MessageCircle,
    new_review: FileText,
    review_vote: Star,
    review_response: MessageCircle,
    listing_status: Tag,
    favorite_price: Tag
} as const;
 
const NotificationDrawer: React.FC<NotificationDrawerProps> = ({ open, onClose }) => {
    const { loading, error, markAsRead } = useNotifications();
    const [notifications, setNotifications] = useState<Notification[]>([]);

    useEffect(() => {
        const fetchNotifications = async (): Promise<void> => {
            try {
                const response = await axios.get('/api/v1/notifications');
                setNotifications(response.data.data || []);
            } catch (error) {
                console.error('Error fetching notifications:', error);
            }
        };
        
        fetchNotifications();
    }, []);

    const handleClick = async (notification: Notification): Promise<void> => {
        if (!notification.is_read && markAsRead) {
            await markAsRead(notification.id);
        }
        
        // Навигация в зависимости от типа уведомления
        switch (notification.type) {
            case 'new_message':
                window.location.href = `/marketplace/chat`;
                break;
            case 'new_review':
            case 'review_vote':
            case 'review_response':
                if (notification.data?.listing_id) {
                    window.location.href = `/marketplace/listings/${notification.data.listing_id}#reviews`;
                }
                break;
            case 'listing_status':
            case 'favorite_price':
                if (notification.data?.listing_id) {
                    window.location.href = `/marketplace/listings/${notification.data.listing_id}`;
                }
                break;
            default:
                break;
        }
        onClose();
    };

    const renderContent = (): React.ReactNode => {
        if (loading) {
            return (
                <Box display="flex" justifyContent="center" p={3}>
                    <CircularProgress />
                </Box>
            );
        }

        if (error) {
            return <Alert severity="error">{error}</Alert>;
        }

        if (!notifications.length) {
            return (
                <Box p={3} textAlign="center">
                    <Typography color="text.secondary">
                        Нет уведомлений
                    </Typography>
                </Box>
            );
        }

        return (
            <List>
                {notifications.map((notification) => {
                    const Icon = NOTIFICATION_ICONS[notification.type] || MessageCircle;
                    return (
                        <React.Fragment key={notification.id}>
                            <ListItem
                                button
                                onClick={() => handleClick(notification)}
                                sx={{
                                    bgcolor: notification.is_read ? 'transparent' : 'action.hover',
                                    '&:hover': {
                                        bgcolor: 'action.hover'
                                    }
                                }}
                            >
                                <Icon size={20} style={{ marginRight: 16, opacity: 0.7 }} />
                                <ListItemText
                                    primary={notification.title}
                                    secondary={
                                        <React.Fragment>
                                            <Typography
                                                component="span"
                                                variant="body2"
                                                color="text.secondary"
                                                sx={{ display: 'block', mb: 0.5 }}
                                            >
                                                {notification.message}
                                            </Typography>
                                            <Typography
                                                component="span"
                                                variant="caption"
                                                color="text.secondary"
                                            >
                                                {formatDistanceToNow(new Date(notification.created_at), {
                                                    addSuffix: true,
                                                    locale: ru
                                                })}
                                            </Typography>
                                        </React.Fragment>
                                    }
                                    primaryTypographyProps={{
                                        variant: 'subtitle2',
                                        fontWeight: notification.is_read ? 'normal' : 'bold'
                                    }}
                                />
                            </ListItem>
                            <Divider component="li" />
                        </React.Fragment>
                    );
                })}
            </List>
        );
    };

    return (
        <Drawer
            anchor="right"
            open={open}
            onClose={onClose}
            PaperProps={{
                sx: { width: { xs: '100%', sm: 400 } }
            }}
        >
            <Box sx={{ 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'space-between',
                p: 2,
                borderBottom: 1,
                borderColor: 'divider'
            }}>
                <Typography variant="h6">Уведомления</Typography>
                <IconButton onClick={onClose}>
                    <X size={20} />
                </IconButton>
            </Box>
            {renderContent()}
        </Drawer>
    );
};

export default NotificationDrawer;