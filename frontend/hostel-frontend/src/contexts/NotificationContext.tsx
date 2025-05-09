import React, { createContext, useState, useEffect, useCallback, ReactNode } from 'react';
import axios from '../api/axios';
import { useAuth } from './AuthContext';

// Типы для уведомлений
interface Notification {
  id: string | number;
  title: string;
  message: string;
  type: string;
  is_read: boolean;
  created_at: string;
  data?: Record<string, any>;
}

// Типы для настроек уведомлений
interface NotificationTypeSettings {
  telegram_enabled: boolean;
  email_enabled: boolean;
}

interface NotificationSettings {
  [key: string]: NotificationTypeSettings;
}

// Тип для контекста уведомлений
interface NotificationContextType {
  notifications: Notification[];
  settings: NotificationSettings;
  telegramConnected: boolean;
  unreadCount: number;
  loading: boolean;
  error: string | null;
  message: string;
  severity: 'success' | 'error' | 'warning' | 'info';
  showNotification: (msg: string, sev?: 'success' | 'error' | 'warning' | 'info') => void;
  fetchNotifications: () => Promise<void>;
  fetchSettings: () => Promise<void>;
  updateSettings: (type: string, channel: string, value: boolean) => Promise<boolean>;
  connectTelegram: () => Promise<any>;
  markAsRead: (notificationId: string | number) => Promise<boolean>;
  setSettings: React.Dispatch<React.SetStateAction<NotificationSettings>>;
}

// Свойства провайдера
interface NotificationProviderProps {
  children: ReactNode;
}

// Создание контекста с начальным значением null 
export const NotificationContext = createContext<NotificationContextType | null>(null);

export const NotificationProvider: React.FC<NotificationProviderProps> = ({ children }) => {
  const { user } = useAuth();
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [settings, setSettings] = useState<NotificationSettings>({});
  const [telegramConnected, setTelegramConnected] = useState<boolean>(false);
  const [unreadCount, setUnreadCount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string>('');
  const [severity, setSeverity] = useState<'success' | 'error' | 'warning' | 'info'>('success');
  const [statusCheckInterval, setStatusCheckInterval] = useState<NodeJS.Timeout | null>(null);

  const showNotification = (msg: string, sev: 'success' | 'error' | 'warning' | 'info' = 'success'): void => {
    setMessage(msg);
    setSeverity(sev);
  };

  const fetchNotifications = useCallback(async (): Promise<void> => {
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

  const fetchSettings = useCallback(async (): Promise<void> => {
    if (!user) return;
    
    try {
      console.log('Fetching notification settings...');
      const response = await axios.get('/api/v1/notifications/settings');
      console.log('Raw settings response:', response.data);
      
      if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
        const settingsArray = response.data.data.data;
        const formattedSettings: NotificationSettings = {};
        
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
        console.warn('Settings data is not an array or is missing:', response.data);
      }
    } catch (err) {
      console.error('Error fetching settings:', err);
    }
  }, [user]);

  const checkTelegramStatus = useCallback(async (): Promise<boolean> => {
    if (!user) return false;
    
    try {
      const response = await axios.get('/api/v1/notifications/telegram');
      setTelegramConnected(!!response.data.data?.connected);
      return !!response.data.data?.connected;
    } catch (err) {
      console.error('Error checking Telegram status:', err);
      return false;
    }
  }, [user]);

  const startStatusCheck = useCallback((): void => {
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

  const connectTelegram = async (): Promise<any> => {
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

  const updateSettings = async (type: string, channel: string, value: boolean): Promise<boolean> => {
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

  const markAsRead = async (notificationId: string | number): Promise<boolean> => {
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

  const contextValue: NotificationContextType = {
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