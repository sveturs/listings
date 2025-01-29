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
            console.log('Raw settings response:', response.data);
    
            // Правильный путь к массиву настроек
            const settingsArray = response.data?.data?.data;
            
            if (Array.isArray(settingsArray)) {
                const formattedSettings = {};
                
                settingsArray.forEach(setting => {
                    if (setting && setting.notification_type) {
                        formattedSettings[setting.notification_type] = {
                            telegram_enabled: Boolean(setting.telegram_enabled),
                            push_enabled: Boolean(setting.push_enabled)
                        };
                    }
                });
    
                console.log('Formatted settings:', formattedSettings);
                setSettings(formattedSettings);
            } else {
                console.warn('Settings data is not an array:', response.data);
                setSettings({});
            }
    
            // Отдельно проверяем статус подключения Telegram
            const telegramStatus = await axios.get('/api/v1/notifications/telegram');
            console.log('Telegram status raw:', telegramStatus.data);
            
            // Проверяем правильный путь к статусу подключения
            if (telegramStatus.data?.data?.connected === true) {
                setTelegramConnected(true);
            } else {
                setTelegramConnected(false);
            }
    
        } catch (err) {
            console.error('Error fetching settings:', err);
            setSettings({});
        }
    };

    useEffect(() => {
        let intervalId;
        
        const checkStatus = async () => {
            try {
                const telegramStatus = await axios.get('/api/v1/notifications/telegram');
                if (telegramStatus.data?.data?.connected) {
                    setTelegramConnected(true);
                    // Если подключено - очищаем интервал
                    clearInterval(intervalId);
                }
            } catch (err) {
                console.error('Error checking telegram status:', err);
            }
        };
    
        // Проверяем статус каждые 5 секунд вместо каждой секунды
        intervalId = setInterval(checkStatus, 5000);
    
        // Очищаем интервал при размонтировании
        return () => {
            clearInterval(intervalId);
        };
    }, []);

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
            console.log('Raw response:', response); // добавим для отладки

            // Исправляем извлечение токена из ответа
            const token = response.data?.data?.token || // старый вариант
                response.data?.token; // новый вариант

            if (!token) {
                console.error('Token structure:', response.data); // добавим для отладки
                throw new Error('Token not received');
            }

            const botLink = `https://t.me/SveTu_bot?start=${token}`;
            console.log('Opening bot link:', botLink);
            window.open(botLink, '_blank');

            startStatusCheck();
            return response.data;
        } catch (err) {
            console.error('Full error details:', err.response || err); // улучшаем логирование
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

        connectTelegram,
        fetchSettings

    };
};