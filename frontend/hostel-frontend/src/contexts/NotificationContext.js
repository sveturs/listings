// frontend/hostel-frontend/src/contexts/NotificationContext.js
import React, { createContext, useContext, useState, useEffect } from 'react';
import axios from '../api/axios';
import { useAuth } from './AuthContext';

const NotificationContext = createContext(null);

export const NotificationProvider = ({ children }) => {
    const [settings, setSettings] = useState({});
    const [telegramConnected, setTelegramConnected] = useState(false);
    const { user } = useAuth();

    useEffect(() => {
        if (user) {
            const checkStatus = async () => {
                try {
                    const response = await axios.get('/api/v1/notifications/telegram');
                    setTelegramConnected(response.data?.data?.connected || false);
                    
                    // Также обновляем настройки
                    const settingsResponse = await axios.get('/api/v1/notifications/settings');
                    if (settingsResponse.data?.data) {
                        const settingsData = settingsResponse.data.data;
                        const formattedSettings = settingsData.reduce((acc, setting) => {
                            acc[setting.notification_type] = {
                                telegram_enabled: setting.telegram_enabled,
                                push_enabled: setting.push_enabled
                            };
                            return acc;
                        }, {});
                        setSettings(formattedSettings);
                    }
                } catch (err) {
                 }
            };
            checkStatus();
        } else {
            setTelegramConnected(false);
            setSettings({});
        }
    }, [user]);

    const [statusCheckInterval, setStatusCheckInterval] = useState(null);

    const startStatusCheck = () => {
        const interval = setInterval(async () => {
            try {
                const response = await axios.get('/api/v1/notifications/telegram');
                if (response.data.connected) {
                    clearInterval(interval);
                    setTelegramConnected(true);
                }
            } catch (err) {
                console.error('Error checking status:', err);
            }
        }, 3000);

        setStatusCheckInterval(interval);

        setTimeout(() => {
            clearInterval(interval);
        }, 120000);
    };
    useEffect(() => {
        return () => {
            if (statusCheckInterval) {
                clearInterval(statusCheckInterval);
            }
        };
    }, [statusCheckInterval]);
    useEffect(() => {
        const fetchSettings = async () => {
            try {
                console.log('Fetching notification settings...');
                const response = await axios.get('/api/v1/notifications/settings');
                console.log('Raw settings response:', response.data);
                
                if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
                    const settingsArray = response.data.data.data;
                    const formattedSettings = {};
                    
                    // Преобразуем массив настроек в объект для удобного доступа
                    settingsArray.forEach(setting => {
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
        
        fetchSettings();
    }, [user]);

    const updateSettings = async (type, channel, value) => {
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

    return (
        <NotificationContext.Provider value={{
            settings,
            telegramConnected,
            updateSettings,
            connectTelegram
        }}>
            {children}
        </NotificationContext.Provider>
    );
};
export const NOTIFICATION_TYPES = {
    NEW_MESSAGE: 'new_message',
    NEW_REVIEW: 'new_review',
    REVIEW_VOTE: 'review_vote',
    REVIEW_RESPONSE: 'review_response',
    LISTING_STATUS: 'listing_status',
    FAVORITE_PRICE: 'favorite_price'
};
export const useNotifications = () => {
    const context = useContext(NotificationContext);
    if (!context) {
        throw new Error('useNotifications must be used within NotificationProvider');
    }
    return context;
};