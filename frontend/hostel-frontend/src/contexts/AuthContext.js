// frontend/hostel-frontend/src/contexts/AuthContext.js

import React, { createContext, useState, useContext, useEffect } from 'react';
import axios from '../api/axios';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  // Сохраняем сессию в localStorage вместо cookie
  const saveSession = (userData) => {
    localStorage.setItem('user_session', JSON.stringify(userData));
  };

  // Проверка и загрузка сессии из localStorage
  const loadSession = () => {
    try {
      const session = localStorage.getItem('user_session');
      if (session) {
        return JSON.parse(session);
      }
    } catch (error) {
      console.error('Error loading session:', error);
    }
    return null;
  };

  // Обновляем метод проверки авторизации
  const checkAuth = async () => {
    try {
      // Сначала проверяем локальную сессию
      const savedSession = loadSession();
      if (savedSession) {
        setUser(savedSession);
      }

      // Затем делаем запрос к серверу для подтверждения
      const response = await axios.get('/auth/session');
      if (response.data.authenticated) {
        setUser(response.data.user);
        saveSession(response.data.user);
      } else {
        // Если сервер говорит, что пользователь не авторизован, очищаем локальную сессию
        localStorage.removeItem('user_session');
        setUser(null);
      }
    } catch (error) {
      console.error('Error checking auth status:', error);
      // В случае ошибки не удаляем локальную сессию - это может быть просто проблема с сетью
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    // При загрузке сразу устанавливаем пользователя из localStorage
    const savedSession = loadSession();
    if (savedSession) {
      setUser(savedSession);
      setIsLoading(false);
    }
    
    // И затем проверяем статус авторизации
    checkAuth();
  }, []);

  const login = (params = '') => {
    window.location.href = `${process.env.REACT_APP_BACKEND_URL}/auth/google${params}`;
  };

  const logout = async () => {
    try {
      await axios.get('/auth/logout', { withCredentials: true });
      localStorage.removeItem('user_session');
      setUser(null);
    } catch (error) {
      console.error('Logout failed:', error);
      // Всё равно очищаем локальную сессию
      localStorage.removeItem('user_session');
      setUser(null);
    }
  };

  const value = {
    user,
    loading: isLoading,
    login,
    logout,
    checkAuth // Экспортируем метод, чтобы можно было вызвать его после успешной оплаты
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};