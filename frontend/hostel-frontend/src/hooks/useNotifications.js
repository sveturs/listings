// /data/proj/hostel-booking-system/frontend/hostel-frontend/src/hooks/useNotifications.js
import { useState, useEffect } from 'react';
import axios from '../api/axios';

export const useNotifications = () => {
    const [notifications, setNotifications] = useState([]);
    const [settings, setSettings] = useState({});
    const [telegramConnected, setTelegramConnected] = useState(false);
    const [message, setMessage] = useState('');
    const [severity, setSeverity] = useState('success');
    const [statusCheckInterval, setStatusCheckInterval] = useState(null);

    const showNotification = (msg, sev = 'success') => {
        setMessage(msg);
        setSeverity(sev);
    };
    const startStatusCheck = () => {
        checkTelegramStatus();
    };

    useEffect(() => {
        fetchSettings();
        checkTelegramStatus();
    }, []);
    const fetchNotifications = async () => {
        try {
            const response = await axios.get('/api/v1/notifications');
            setNotifications(response.data.data || []);
        } catch (err) {
            console.error('Error:', err);
        }
    };
    const fetchSettings = async () => {
        try {
            const response = await axios.get('/api/v1/notifications/settings');
            setSettings(response.data.data || {});
        } catch (err) {
            console.error('Error fetching settings:', err);
        }
    };


    const updateSettings = async (type, channel, value) => {
        try {
            await axios.put('/api/v1/notifications/settings', {
                notification_type: type,
                [channel + '_enabled']: value
            });
            setSettings(prev => ({
                ...prev,
                [type]: { ...prev[type], [channel]: value }
            }));
            return true;
        } catch (err) {
            console.error('Error updating settings:', err);
            return false;
        }
    };

    const connectTelegram = async () => {
        try {
            const response = await axios.post('/api/v1/notifications/telegram/token');
            if (!response.data.token) {
                throw new Error('Токен не получен');
            }
            const botLink = `https://t.me/SveTu_bot?start=${response.data.token}`;
            window.open(botLink, '_blank');
            startStatusCheck();
        } catch (err) {
            showNotification(err.response?.data?.error || 'Ошибка подключения Telegram', 'error');
        }
    };
    const checkTelegramStatus = () => {
        if (!statusCheckInterval) {
            const interval = setInterval(async () => {
                try {
                    const response = await axios.get('/api/v1/notifications/telegram');
                    if (response.data.connected) {
                        clearInterval(interval);
                        setStatusCheckInterval(null);
                        setTelegramConnected(true);
                        showNotification('Telegram успешно подключен');
                    }
                } catch (err) {
                    console.error('Ошибка проверки статуса:', err);
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
        const init = async () => {
            try {
                const response = await axios.get('/api/v1/notifications/telegram');
                setTelegramConnected(!!response.data.connected);
            } catch (err) {
                console.error('Error checking initial Telegram status:', err);
            }
        };
        
        init();
        fetchSettings();
    }, []);
    return {
        notifications,
        settings,
        telegramConnected,
        updateSettings,
        message,
        severity,
        showNotification,

        connectTelegram

    };
};