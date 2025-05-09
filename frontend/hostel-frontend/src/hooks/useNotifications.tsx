// frontend/hostel-frontend/src/hooks/useNotifications.tsx
import axios from '../api/axios';
import { useState, useEffect, useCallback } from 'react';

interface NotificationTypeSettings {
    telegram_enabled: boolean;
    email_enabled: boolean;
}

interface NotificationSettings {
    [key: string]: NotificationTypeSettings;
}

interface Notification {
    id: string | number;
    title: string;
    message: string;
    type: string;
    is_read: boolean;
    created_at: string;
    data?: Record<string, any>;
}

interface UseNotificationsReturn {
    notifications: Notification[];
    settings: NotificationSettings;
    telegramConnected: boolean;
    updateSettings: (type: string, channel: string, value: boolean) => Promise<boolean>;
    message: string;
    severity: 'success' | 'error' | 'warning' | 'info';
    showNotification: (msg: string, sev?: 'success' | 'error' | 'warning' | 'info') => void;
    connectTelegram: () => Promise<any>;
    fetchSettings: () => Promise<void>;
    setSettings?: (settings: NotificationSettings) => void;
    unreadCount?: number;
    markAsRead?: (id: string | number) => Promise<void>;
    loading?: boolean;
    error?: string | null;
}

export const useNotifications = (): UseNotificationsReturn => {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [settings, setSettings] = useState<NotificationSettings>({});
    const [telegramConnected, setTelegramConnected] = useState<boolean>(false);
    const [message, setMessage] = useState<string>('');
    const [severity, setSeverity] = useState<'success' | 'error' | 'warning' | 'info'>('success');
    const [statusCheckInterval, setStatusCheckInterval] = useState<NodeJS.Timeout | null>(null);
    const [unreadCount, setUnreadCount] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const showNotification = (msg: string, sev: 'success' | 'error' | 'warning' | 'info' = 'success'): void => {
        setMessage(msg);
        setSeverity(sev);
    };
    
    const startStatusCheck = (): void => {
        checkTelegramStatus();
    };

    useEffect(() => {
        fetchSettings();
        checkTelegramStatus();
    }, []);

    const fetchNotifications = async (): Promise<void> => {
        try {
            const response = await axios.get('/api/v1/notifications');
            setNotifications(response.data.data || []);
            // Подсчитываем непрочитанные уведомления
            if (response.data.data && Array.isArray(response.data.data)) {
                const unread = response.data.data.filter((notification: Notification) => !notification.is_read).length;
                setUnreadCount(unread);
            }
        } catch (err) {
            console.error('Error fetching notifications:', err);
        }
    };

    useEffect(() => {
        let intervalId: NodeJS.Timeout;
        
        const checkStatus = async (): Promise<void> => {
            try {
                const telegramStatus = await axios.get('/api/v1/notifications/telegram');
                if (telegramStatus.data?.data?.connected) {
                    setTelegramConnected(true);
                    // Если подключено - очищаем интервал
                    clearInterval(intervalId);
                }
            } catch (err) {
                // Обработка ошибки
            }
        };
    
        // Проверяем статус каждые 5 секунд
        intervalId = setInterval(checkStatus, 5000);
    
        // Очищаем интервал при размонтировании
        return () => {
            clearInterval(intervalId);
        };
    }, []);

    // Функция для обновления настроек
    const updateSettings = async (type: string, channel: string, value: boolean): Promise<boolean> => {
        try {
            console.log(`Updating ${channel} for ${type} to ${value}`);
            
            // Отправляем только измененное поле
            const payload = {
                notification_type: type,
                [`${channel}_enabled`]: value
            };
            
            console.log("Sending payload:", payload);
            
            const response = await axios.put('/api/v1/notifications/settings', payload);
            console.log("Server response:", response.data);
            
            if (response.data.success) {
                // После успешного обновления запрашиваем свежие настройки
                await fetchSettings();
                return true;
            }
            
            return false;
        } catch (err) {
            console.error('Error updating settings:', err);
            return false;
        }
    };

    // Функция для получения настроек
    const fetchSettings = async (): Promise<void> => {
        try {
            console.log('Fetching notification settings...');
            const response = await axios.get('/api/v1/notifications/settings');
            console.log('Raw settings response:', response.data);
            
            if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
                const settingsArray = response.data.data.data;
                const formattedSettings: NotificationSettings = {};
                
                // Преобразуем массив настроек в объект для удобного доступа
                settingsArray.forEach((setting: any) => {
                    if (setting && setting.notification_type) {
                        formattedSettings[setting.notification_type] = {
                            telegram_enabled: Boolean(setting.telegram_enabled),
                            email_enabled: Boolean(setting.email_enabled)
                        };
                    }
                });
                
                console.log('Formatted settings:', formattedSettings);
                setSettings(formattedSettings);
            } else {
                console.warn('Settings data has unexpected format:', response.data);
            }
        } catch (err) {
            console.error('Error fetching settings:', err);
        }
    };
    
    const connectTelegram = async (): Promise<any> => {
        try {
            const response = await axios.post('/api/v1/notifications/telegram/token');
            console.log('Raw response:', response);

            // Исправляем извлечение токена из ответа
            const token = response.data?.data?.token || // старый вариант
                response.data?.token; // новый вариант

            if (!token) {
                console.error('Token structure:', response.data);
                throw new Error('Token not received');
            }

            const botLink = `https://t.me/SveTu_bot?start=${token}`;
            console.log('Opening bot link:', botLink);
            window.open(botLink, '_blank');

            startStatusCheck();
            return response.data;
        } catch (err) {
            console.error('Full error details:', err);
            throw err;
        }
    };

    const checkTelegramStatus = (): void => {
        if (!statusCheckInterval) {
            console.log('Starting telegram status check');
            const interval = setInterval(async () => {
                try {
                    const response = await axios.get('/api/v1/notifications/telegram');
                     if (response.data?.data?.connected) {
                        clearInterval(interval);
                        setStatusCheckInterval(null);
                        setTelegramConnected(true);
                        await fetchSettings();
                        showNotification('Telegram успешно подключен');
                    }
                } catch (err) {
                    // Обработка ошибки
                }
            }, 3000);

            setStatusCheckInterval(interval);

            setTimeout(() => {
                clearInterval(interval);
                setStatusCheckInterval(null);
            }, 120000);
        }
    };

    useEffect(() => {
        return () => {
            if (statusCheckInterval) {
                clearInterval(statusCheckInterval);
            }
        };
    }, [statusCheckInterval]);

    useEffect(() => {
        const init = async (): Promise<void> => {
            try {
                const response = await axios.get('/api/v1/notifications/telegram');
                setTelegramConnected(!!response.data.connected);
            } catch (err) {
                console.error('Error checking initial Telegram status:', err);
            }
        };

        init();
        fetchSettings();
        fetchNotifications(); // Загружаем уведомления при инициализации
    }, []);

    const markAsRead = useCallback(async (id: string | number): Promise<void> => {
        try {
            setLoading(true);
            await axios.post(`/api/v1/notifications/${id}/read`);
            await fetchNotifications(); // Обновляем список уведомлений
            setLoading(false);
        } catch (err) {
            console.error('Error marking notification as read:', err);
            setError('Failed to mark notification as read');
            setLoading(false);
        }
    }, []);

    return {
        notifications,
        settings,
        telegramConnected,
        updateSettings,
        message,
        severity,
        showNotification,
        connectTelegram,
        fetchSettings,
        setSettings,
        unreadCount,
        markAsRead,
        loading,
        error
    };
};