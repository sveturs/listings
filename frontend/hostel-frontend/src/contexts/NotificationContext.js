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
       if (!user) return false;
       window.open('https://t.me/SveTu_bot', '_blank');
       try {
           const response = await axios.get('/api/v1/notifications/telegram/status');
           setTelegramConnected(response.data.connected);
           return response.data.connected;
       } catch (err) {
           console.error('Error:', err);
           return false;
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