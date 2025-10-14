'use client';

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
} from 'react';
import { useRouter } from '@/i18n/routing';
import { useDispatch } from 'react-redux';
import authService from '@/services/auth';
import type { User, UpdateProfileRequest } from '@/types/auth';
import { reset as resetChat, closeWebSocket } from '@/store/slices/chatSlice';
import { resetCart } from '@/store/slices/cartSlice';
import configManager from '@/config';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isLoggingOut: boolean;
  isUpdatingProfile: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  loginWithGoogle: () => void;
  register: (email: string, password: string, name?: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
  updateProfile: (data: UpdateProfileRequest) => Promise<boolean>;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const dispatch = useDispatch();
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [isUpdatingProfile, setIsUpdatingProfile] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Load session on mount
  const refreshSession = useCallback(async () => {
    try {
      const user = await authService.getSession();
      if (user) {
        // Add default provider if missing
        const userWithProvider: User = {
          id: user.id,
          email: user.email,
          name: user.name || user.email,
          provider: 'email', // Default provider since backend doesn't send it
          is_admin: user.is_admin, // Backend теперь возвращает is_admin из middleware
          roles: user.roles,
        };
        setUser(userWithProvider);
      } else {
        setUser(null);
      }
    } catch (error) {
      console.error('Session refresh error:', error);
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    refreshSession();
  }, [refreshSession]);

  const login = useCallback(
    async (email: string, password: string) => {
      setError(null);
      try {
        const response = await authService.login({ email, password });
        if (response && response.user) {
          // Add default provider if missing
          const userWithProvider: User = {
            id: response.user.id,
            email: response.user.email,
            name: response.user.name || response.user.email,
            provider: 'email', // Default provider since backend doesn't send it
            is_admin: response.user.is_admin,
            roles: response.user.roles,
          };
          setUser(userWithProvider);
          router.push('/');
        } else {
          const errorMessage = response?.error || 'Login failed';
          setError(errorMessage);
          throw new Error(errorMessage);
        }
      } catch (error: any) {
        console.error('Login error:', error);
        const errorMessage = error.message || 'Login failed';
        setError(errorMessage);
        throw error;
      }
    },
    [router]
  );

  const loginWithGoogle = useCallback(() => {
    // Redirect to backend OAuth endpoint
    const apiUrl = configManager.getApiUrl() || 'http://localhost:31876';
    window.location.href = `${apiUrl}/api/v1/auth/google`;
  }, []);

  const register = useCallback(
    async (email: string, password: string, name?: string) => {
      setError(null);
      try {
        const response = await authService.register({
          email,
          password,
          name: name || '',
        });
        if (response && response.user) {
          // Add default provider if missing
          const userWithProvider: User = {
            id: response.user.id,
            email: response.user.email,
            name: response.user.name || name || response.user.email,
            provider: 'email', // Default provider since backend doesn't send it
            is_admin: response.user.is_admin,
            roles: response.user.roles,
          };
          setUser(userWithProvider);
          router.push('/');
        } else {
          const errorMessage = response?.error || 'Registration failed';
          setError(errorMessage);
          throw new Error(errorMessage);
        }
      } catch (error: any) {
        console.error('Registration error:', error);
        const errorMessage = error.message || 'Registration failed';
        setError(errorMessage);
        throw error;
      }
    },
    [router]
  );

  const logout = useCallback(async () => {
    setIsLoggingOut(true);
    setError(null);
    try {
      // Закрываем WebSocket соединение
      dispatch(closeWebSocket());

      // Очищаем Redux store
      dispatch(resetChat()); // Чаты, сообщения, счетчики
      dispatch(resetCart()); // Корзина, заказы

      // Выполняем logout на backend
      await authService.logout();

      // Очищаем user state
      setUser(null);

      // Редирект на главную
      router.push('/');
    } catch (error) {
      console.error('Logout error:', error);
      setError('Failed to logout. Please try again.');
    } finally {
      setIsLoggingOut(false);
    }
  }, [router, dispatch]);

  const updateProfile = useCallback(
    async (data: UpdateProfileRequest): Promise<boolean> => {
      setIsUpdatingProfile(true);
      setError(null);
      try {
        await authService.updateProfile(data);
        // Refresh session to get updated user data
        await refreshSession();
        return true;
      } catch (error) {
        console.error('Profile update error:', error);
        const message =
          error instanceof Error
            ? error.message
            : 'Failed to update profile. Please try again.';
        setError(message);
        return false;
      } finally {
        setIsUpdatingProfile(false);
      }
    },
    [refreshSession]
  );

  const value = useMemo<AuthContextType>(
    () => ({
      user,
      isAuthenticated: !!user,
      isLoading,
      isLoggingOut,
      isUpdatingProfile,
      error,
      login,
      loginWithGoogle,
      register,
      logout,
      refreshSession,
      updateProfile,
      clearError,
    }),
    [
      user,
      isLoading,
      isLoggingOut,
      isUpdatingProfile,
      error,
      login,
      loginWithGoogle,
      register,
      logout,
      refreshSession,
      updateProfile,
      clearError,
    ]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

export function useAuthContext() {
  return useAuth();
}
