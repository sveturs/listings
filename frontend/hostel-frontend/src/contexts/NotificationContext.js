import React, { createContext, useState, useEffect, useCallback } from 'react';
import axios from '../api/axios';
import { useAuth } from './AuthContext';

export const NotificationContext = createContext();

export const NotificationProvider = ({ children }) => {
  const { user } = useAuth();
  const [notifications, setNotifications] = useState([]);
  const [settings, setSettings] = useState({});
  const [telegramConnected, setTelegramConnected] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [message, setMessage] = useState('');
  const [severity, setSeverity] = useState('success');
  const [statusCheckInterval, setStatusCheckInterval] = useState(null);

  const showNotification = (msg, sev = 'success') => {
    setMessage(msg);
    setSeverity(sev);
  };

  const fetchNotifications = useCallback(async () => {
    if (!user) return;
    
    try {
      setLoading(true);
      const response = await axios.get('/api/v1/notifications');
      const notificationsData = Array.isArray(response.data.data) ? response.data.data : [];
      setNotifications(notificationsData);
      
      // Считаем непрочитанные уведомления
      const unread = notificationsData.filter(n => !n.is_read).length;
      setUnreadCount(unread);
    } catch (err) {
      console.error('Error fetching notifications:', err);
      setError('Failed to load notifications');
    } finally {
      setLoading(false);
    }
  }, [user]);

  const fetchSettings = useCallback(async () => {
    if (!user) return;
    
    try {
      console.log('Fetching notification settings...');
      const response = await axios.get('/api/v1/notifications/settings');
      console.log('Raw settings response:', response.data);
      
      if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
        const settingsArray = response.data.data.data;
        const formattedSettings = {};
        
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
        console.warn('Settings data is not an array or is missing:', response.data);
      }
    } catch (err) {
      console.error('Error fetching settings:', err);
    }
  }, [user]);

  const checkTelegramStatus = useCallback(async () => {
    if (!user) return;
    
    try {
      const response = await axios.get('/api/v1/notifications/telegram');
      setTelegramConnected(!!response.data.data?.connected);
      return !!response.data.data?.connected;
    } catch (err) {
      console.error('Error checking Telegram status:', err);
      return false;
    }
  }, [user]);

  const startStatusCheck = useCallback(() => {
    if (statusCheckInterval) {
      clearInterval(statusCheckInterval);
    }
    
    console.log('Starting telegram status check');
    const interval = setInterval(async () => {
      const connected = await checkTelegramStatus();
      if (connected) {
        clearInterval(interval);
        setStatusCheckInterval(null);
        await fetchSettings();
        showNotification('Telegram успешно подключен');
      }
    }, 3000);
    
    setStatusCheckInterval(interval);
    
    setTimeout(() => {
      clearInterval(interval);
      setStatusCheckInterval(null);
    }, 120000);
  }, [statusCheckInterval, checkTelegramStatus, fetchSettings]);

  const connectTelegram = async () => {
    try {
      const response = await axios.get('/api/v1/notifications/telegram/token');
      console.log('Raw response:', response);
      
      const token = response.data?.data?.token || response.data?.token;
      
      if (!token) {
        console.error('Token not received:', response.data);
        throw new Error('Не удалось получить токен');
      }
      
      const botLink = `https://t.me/SveTu_bot?start=${token}`;
      console.log('Opening bot link:', botLink);
      window.open(botLink, '_blank');
      
      startStatusCheck();
      return response.data;
    } catch (err) {
      console.error('Error connecting Telegram:', err);
      throw err;
    }
  };

  const updateSettings = async (type, channel, value) => {
    try {
      console.log(`Updating ${channel} for ${type} to ${value}`);
      
      const payload = {
        notification_type: type,
        [`${channel}_enabled`]: value
      };
      
      console.log("Sending payload:", payload);
      
      const response = await axios.put('/api/v1/notifications/settings', payload);
      console.log("Server response:", response.data);
      
      if (response.data.success) {
        // После успешного обновления, получаем актуальные настройки
        await fetchSettings();
        return true;
      }
      
      return false;
    } catch (err) {
      console.error('Error updating settings:', err);
      return false;
    }
  };

  const markAsRead = async (notificationId) => {
    try {
      await axios.put(`/api/v1/notifications/${notificationId}/read`);
      
      // Обновляем локальный список уведомлений
      setNotifications(prev => prev.map(n => 
        n.id === notificationId ? { ...n, is_read: true } : n
      ));
      
      // Уменьшаем счетчик непрочитанных
      setUnreadCount(prev => Math.max(0, prev - 1));
      
      return true;
    } catch (err) {
      console.error('Error marking notification as read:', err);
      return false;
    }
  };

  useEffect(() => {
    if (user) {
      fetchNotifications();
      checkTelegramStatus();
      fetchSettings();
    }
  }, [user, fetchNotifications, checkTelegramStatus, fetchSettings]);

  useEffect(() => {
    return () => {
      if (statusCheckInterval) {
        clearInterval(statusCheckInterval);
      }
    };
  }, [statusCheckInterval]);

  const contextValue = {
    notifications,
    settings,
    telegramConnected,
    unreadCount,
    loading,
    error,
    message,
    severity,
    showNotification,
    fetchNotifications,
    fetchSettings,
    updateSettings,
    connectTelegram,
    markAsRead,
    setSettings
  };

  return (
    <NotificationContext.Provider value={contextValue}>
      {children}
    </NotificationContext.Provider>
  );
};

export default NotificationProvider;