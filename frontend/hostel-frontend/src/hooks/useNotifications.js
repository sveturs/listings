import { useState, useEffect } from 'react';
import axios from '../api/axios';

export const useNotifications = () => {
    const [notifications, setNotifications] = useState([]);
    const [settings, setSettings] = useState({});
    const [telegramConnected, setTelegramConnected] = useState(false);

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

    const checkTelegramStatus = async () => {
        try {
            const response = await axios.get('/api/v1/notifications/telegram');
            setTelegramConnected(response.data.connected || false);
        } catch (err) {
            setTelegramConnected(false);
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
        window.open('https://t.me/SveTu_bot', '_blank');
        await checkTelegramStatus();
    };

    return {
        notifications,
        settings,
        telegramConnected,
        updateSettings,
        connectTelegram
    };
};