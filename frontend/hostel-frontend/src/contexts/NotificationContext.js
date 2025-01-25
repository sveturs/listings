import React, { createContext, useContext, useState, useEffect } from 'react';
import axios from '../api/axios';
import { useAuth } from './AuthContext';

const NotificationContext = createContext(null);

export const NotificationProvider = ({ children }) => {
    const [settings, setSettings] = useState({});
    const [telegramConnected, setTelegramConnected] = useState(false);
    const { user } = useAuth();

    useEffect(() => {
        const fetchSettings = async () => {
            if (!user) return;

            try {
                const [settingsResponse, telegramResponse] = await Promise.all([
                    axios.get('/api/v1/notifications/settings'),
                    axios.get('/api/v1/notifications/telegram')
                ]);
                 console.log('Raw settings:', settingsResponse.data);
                const formattedSettings = {};
                (settingsResponse.data.data || []).forEach(setting => {  
                    formattedSettings[setting.notification_type] = {
                        telegram: setting.telegram_enabled,
                        push: setting.push_enabled
                    };
                });

                setSettings(formattedSettings);
                setTelegramConnected(!!telegramResponse.data.connected);
            } catch (err) {
                console.error('Error fetching notification settings:', err);
            }
        };

        fetchSettings();
    }, [user]);

    const updateSettings = async (type, channel, value) => {
        if (!user) return;

        try {
            await axios.put('/api/v1/notifications/settings', {
                notification_type: type,
                telegram_enabled: channel === 'telegram' ? value : settings[type]?.telegram || false,
                push_enabled: channel === 'push' ? value : settings[type]?.push || false
            });

            setSettings(prev => ({
                ...prev,
                [type]: {
                    ...prev[type],
                    [channel]: value
                }
            }));

            return true;
        } catch (err) {
            console.error('Error updating notification settings:', err);
            return false;
        }
    };

    const connectTelegram = async () => {
        if (!user) return false;

        // Открываем Telegram бот
        window.open('https://t.me/SveTu_bot', '_blank');

        // Ждем подключения
        try {
            const response = await axios.get('/api/v1/notifications/telegram/status');
            setTelegramConnected(response.data.connected);
            return response.data.connected;
        } catch (err) {
            console.error('Error checking telegram status:', err);
            return false;
        }
    };

    const value = {
        settings,
        telegramConnected,
        updateSettings,
        connectTelegram
    };

    return (
        <NotificationContext.Provider value={value}>
            {children}
        </NotificationContext.Provider>
    );
};


export const useNotifications = () => {
    const context = useContext(NotificationContext);
    if (!context) {
        throw new Error('useNotifications must be used within a NotificationProvider');
    }
    return context;
};

// Константы типов уведомлений
export const NOTIFICATION_TYPES = {
    NEW_MESSAGE: 'new_message',
    NEW_REVIEW: 'new_review',
    REVIEW_VOTE: 'review_vote',
    REVIEW_RESPONSE: 'review_response',
    LISTING_STATUS: 'listing_status',
    FAVORITE_PRICE: 'favorite_price'
};