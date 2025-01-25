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
       const fetchSettings = async () => {
           if (!user) return;
           try {
               const [settingsResponse, telegramResponse] = await Promise.all([
                   axios.get('/api/v1/notifications/settings'),
                   axios.get('/api/v1/notifications/telegram')
               ]);

               const settingsData = settingsResponse.data.data || {};
               setSettings(settingsData);
               setTelegramConnected(!!telegramResponse.data.connected);
           } catch (err) {
               console.error('Error fetching settings:', err);
           }
       };
       fetchSettings();
   }, [user]);

   const updateSettings = async (type, channel, value) => {
       if (!user) return false;
       try {
           await axios.put('/api/v1/notifications/settings', {
               notification_type: type,
               [channel + '_enabled']: value
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
           console.error('Error:', err);
           return false;
       }
   };

   const connectTelegram = async () => {
    try {
        const response = await axios.post('/api/v1/notifications/telegram/token');
        if (response.data.token) {
            const botLink = `https://t.me/SveTu_bot?start=${response.data.token}`;
            window.open(botLink, '_blank');
            startStatusCheck();
        }
    } catch (err) {
        console.error('Error:', err);
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

export const useNotifications = () => {
   const context = useContext(NotificationContext);
   if (!context) {
       throw new Error('useNotifications must be used within NotificationProvider');
   }
   return context;
};