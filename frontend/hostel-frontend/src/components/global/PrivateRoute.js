import React from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";

const PrivateRoute = ({ children }) => {
  const { user, login } = useAuth();

  // Если пользователь не авторизован, перенаправляем на страницу логина
  if (!user) {
    // Вызываем login, если пользователь не авторизован
    login();
    return null; // Возвращаем null, чтобы не рендерить лишний компонент
  }

  // Если авторизован, показываем защищенную страницу
  return children;
};

export default PrivateRoute;
