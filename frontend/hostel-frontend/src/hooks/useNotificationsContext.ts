import { useContext } from 'react';
import { NotificationContext } from '../contexts/NotificationContext';

// Хук для удобного использования контекста уведомлений
export const useNotificationsContext = () => {
  const context = useContext(NotificationContext);
  
  if (!context) {
    throw new Error('useNotificationsContext must be used within a NotificationProvider');
  }
  
  return context;
};

export default useNotificationsContext;