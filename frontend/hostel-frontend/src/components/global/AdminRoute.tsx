// frontend/hostel-frontend/src/components/global/AdminRoute.tsx
import React, { ReactNode, useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { isAdmin, checkAdminStatus } from '../../utils/adminUtils';
import { CircularProgress, Box } from '@mui/material';

interface AdminRouteProps {
  children: ReactNode;
}

const AdminRoute: React.FC<AdminRouteProps> = ({ children }) => {
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    const verifyAdmin = async () => {
      if (!user) {
        setLoading(false);
        return;
      }

      // Первая проверка - по жесткому списку для быстрого доступа
      if (isAdmin(user.email)) {
        setIsAuthorized(true);
        setLoading(false);
        return;
      }

      // Вторая проверка - через API для проверки динамического списка админов
      try {
        const isAdminFromApi = await checkAdminStatus(user.email);
        setIsAuthorized(isAdminFromApi);
      } catch (error) {
        console.error('Error verifying admin status:', error);
        setIsAuthorized(false);
      } finally {
        setLoading(false);
      }
    };

    verifyAdmin();
  }, [user]);

  // Показываем загрузку во время проверки
  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  // Проверяем авторизацию после загрузки
  if (!user || !isAuthorized) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

export default AdminRoute;