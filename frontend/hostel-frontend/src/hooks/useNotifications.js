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
            console.log('Fetching settings...');
            const response = await axios.get('/api/v1/notifications/settings');
            console.log('Settings response:', response.data);
            
            if (response.data?.data) {
                const settingsData = Array.isArray(response.data.data) ? response.data.data : [];
                const formattedSettings = settingsData.reduce((acc, setting) => {
                    acc[setting.notification_type] = {
                        telegram_enabled: setting.telegram_enabled,
                        push_enabled: setting.push_enabled
                    };
                    return acc;
                }, {});
                console.log('Formatted settings:', formattedSettings);
                setSettings(formattedSettings);
            }
        } catch (err) {
            console.error('Error fetching settings:', err);
            setSettings({});
        }
    };
    


    const updateSettings = async (type, channel, value) => {
        try {
            const response = await axios.put('/api/v1/notifications/settings', {
                notification_type: type,
                [`${channel}_enabled`]: value
            });
            
            if (response.data.success) {
                // Правильно обновляем состояние
                setSettings(prev => ({
                    ...prev,
                    [type]: {
                        ...prev[type],
                        [`${channel}_enabled`]: value
                    }
                }));
                return true;
            }
            return false;
        } catch (err) {
            console.error('Error updating settings:', err);
            return false;
        }
    };
    const connectTelegram = async () => {
        try {
            const response = await axios.post('/api/v1/notifications/telegram/token');
            console.log('Full response data:', JSON.stringify(response.data));
                
            const token = response.data?.data?.token;
            if (!token) {
                throw new Error('Token not received');
            }
    
            const botLink = `https://t.me/SveTu_bot?start=${token}`;
            console.log('Opening bot link:', botLink);
            
            // Создаем и добавляем скрытый элемент <a> для совместимости со всеми браузерами
            const link = document.createElement('a');
            link.href = botLink;
            link.target = '_blank';
            link.rel = 'noopener noreferrer';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            
            // Добавляем задержку перед началом проверки статуса
            await new Promise(resolve => setTimeout(resolve, 1000));
            startStatusCheck();
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
                    if (response.data?.data?.connected) {
                        clearInterval(interval);
                        setStatusCheckInterval(null);
                        setTelegramConnected(true);
                        await fetchSettings();
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