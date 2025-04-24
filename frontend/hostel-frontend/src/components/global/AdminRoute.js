// frontend/hostel-frontend/src/components/global/AdminRoute.js
import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { isAdmin } from '../../utils/adminUtils';

const AdminRoute = ({ children }) => {
  const { user } = useAuth();

  // Проверка, что пользователь авторизован
  if (!user) {
    return <Navigate to="/" replace />;
  }

  // Проверка, что email пользователя находится в списке администраторов
  if (!isAdmin(user.email)) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export default AdminRoute;