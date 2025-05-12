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
  const { user, checkAuth } = useAuth();
  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    const verifyAdmin = async () => {
      setLoading(true);

      if (!user) {
        setLoading(false);
        return;
      }

      // Проверяем флаг is_admin напрямую сначала
      if (user.is_admin) {
        setIsAuthorized(true);
        setLoading(false);
        return;
      }

      // Если нет флага, проверяем через API
      try {
        // Обновим данные пользователя для получения актуального статуса
        await checkAuth();

        // Если после обновления появился флаг is_admin
        if (user?.is_admin) {
          setIsAuthorized(true);
        } else {
          // Только если флага все еще нет, делаем отдельный запрос
          const isAdminFromApi = await checkAdminStatus(user.email);
          setIsAuthorized(isAdminFromApi);
        }
      } catch (error) {
        console.error('Error verifying admin status:', error);
        setIsAuthorized(false);
      } finally {
        setLoading(false);
      }
    };

    verifyAdmin();
  }, [user, checkAuth]);


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