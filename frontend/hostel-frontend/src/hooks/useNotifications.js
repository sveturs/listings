// frontend/hostel-frontend/src/hooks/useNotifications.js
import axios from '../api/axios';
import { useState, useEffect, useContext } from 'react';
import { NotificationContext } from '../contexts/NotificationContext';

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
            console.log('Full response data:', JSON.stringify(response.data));
                
            const token = response.data?.data?.data?.token;
            if (token) {
                const botLink = `https://t.me/SveTu_bot?start=${token}`;
                console.log('Opening bot link:', botLink);
                window.open(botLink, '_blank');
                startStatusCheck();
            } else {
                throw new Error('Token not received');
            }
        } catch (err) {
            console.error('Error in connectTelegram:', err);
            throw err;
        }
    };

    const checkTelegramStatus = () => {
        if (!statusCheckInterval) {
            console.log('Starting telegram status check');
            const interval = setInterval(async () => {
                try {
                    const response = await axios.get('/api/v1/notifications/telegram');
                    console.log('Telegram status response:', response.data);
                    if (response.data.connected) {
                        clearInterval(interval);
                        setStatusCheckInterval(null);
                        setTelegramConnected(true);
                        showNotification('Telegram успешно подключен');
                    }
                } catch (err) {
                    console.error('Telegram status check error:', err);
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