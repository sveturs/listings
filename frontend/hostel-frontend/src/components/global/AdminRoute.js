// frontend/hostel-frontend/src/components/global/AdminRoute.js
import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

// Компонент маршрута, доступного только для администратора
const AdminRoute = ({ children }) => {
  const { user } = useAuth();

  // Проверка, что пользователь авторизован
  if (!user) {
    return <Navigate to="/" replace />;
  }

  // Проверка, что email пользователя соответствует указанному
  if (user.email !== 'voroshilovdo@gmail.com') {
    // Если email не соответствует, перенаправляем на главную
    return <Navigate to="/" replace />;
  }

  // Если все проверки пройдены - отображаем дочерний компонент
  return children;
};

export default AdminRoute;