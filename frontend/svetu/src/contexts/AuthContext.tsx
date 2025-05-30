'use client';

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
  useRef,
} from 'react';
import { useRouter } from '@/i18n/routing';
import { AuthService } from '@/services/auth';
import { AuthErrorBoundary } from '@/components/ErrorBoundary';
import type { User, UpdateProfileRequest } from '@/types/auth';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isLoggingOut: boolean;
  isUpdatingProfile: boolean;
  error: string | null;
  login: (returnTo?: string) => void;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
  updateProfile: (data: UpdateProfileRequest) => Promise<boolean>;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  // Инициализация с кешированным состоянием из sessionStorage
  const [user, setUser] = useState<User | null>(() => {
    if (typeof window !== 'undefined') {
      try {
        const cached = sessionStorage.getItem('svetu_user');
        return cached ? JSON.parse(cached) : null;
      } catch {
        return null;
      }
    }
    return null;
  });

  // Если есть кешированный пользователь, начинаем с false, иначе с true
  const [isLoading, setIsLoading] = useState(() => {
    if (typeof window !== 'undefined') {
      try {
        const cached = sessionStorage.getItem('svetu_user');
        return !cached; // false если есть кеш, true если нет
      } catch {
        return true;
      }
    }
    return true;
  });
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [isUpdatingProfile, setIsUpdatingProfile] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Cooldown для предотвращения частых вызовов refreshSession
  const lastRefreshAttempt = useRef<number>(0);
  const REFRESH_COOLDOWN = 5000; // 5 секунд

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Функция для кеширования пользователя в sessionStorage
  const cacheUser = useCallback((userData: User | null) => {
    if (typeof window !== 'undefined') {
      try {
        if (userData) {
          sessionStorage.setItem('svetu_user', JSON.stringify(userData));
        } else {
          sessionStorage.removeItem('svetu_user');
        }
      } catch (error) {
        console.warn('Failed to cache user data:', error);
      }
    }
  }, []);

  // Обертка для setUser с кешированием
  const updateUser = useCallback(
    (userData: User | null) => {
      setUser(userData);
      cacheUser(userData);
    },
    [cacheUser]
  );

  const refreshSession = useCallback(
    async (retries = 3) => {
      const now = Date.now();
      if (now - lastRefreshAttempt.current < REFRESH_COOLDOWN) {
        if (process.env.NODE_ENV === 'development') {
          console.log('RefreshSession skipped due to cooldown');
        }
        return;
      }
      lastRefreshAttempt.current = now;

      for (let i = 0; i < retries; i++) {
        try {
          const session = await AuthService.getSession();
          if (session.authenticated && session.user) {
            updateUser(session.user);
            setError(null);
          } else {
            updateUser(null);
          }
          setIsLoading(false);
          return;
        } catch (error) {
          console.error(
            `Session refresh error (attempt ${i + 1}/${retries}):`,
            error
          );
          if (i === retries - 1) {
            setError('Failed to load session. Please try refreshing the page.');
            updateUser(null);
            setIsLoading(false);
          } else {
            // Exponential backoff
            await new Promise((resolve) =>
              setTimeout(resolve, 1000 * Math.pow(2, i))
            );
          }
        }
      }
    },
    [updateUser, REFRESH_COOLDOWN]
  );

  useEffect(() => {
    // Проверяем наличие кешированного пользователя при первой загрузке
    const hasCachedUser =
      typeof window !== 'undefined' && sessionStorage.getItem('svetu_user');
    if (hasCachedUser) {
      // Проверяем актуальность в фоне, не блокируя UI
      setTimeout(() => refreshSession(), 100);
    } else {
      // Если нет кешированного пользователя, делаем полную проверку
      refreshSession();
    }

    // Cleanup on unmount
    return () => {
      AuthService.cleanup();
    };
  }, [refreshSession]);

  const login = useCallback((returnTo?: string) => {
    try {
      AuthService.loginWithGoogle(returnTo);
    } catch (error) {
      console.error('Login error:', error);
      setError('Failed to initiate login. Please try again.');
    }
  }, []);

  const logout = useCallback(async () => {
    setIsLoggingOut(true);
    setError(null);
    try {
      await AuthService.logout();
      updateUser(null);
      router.push('/');
    } catch (error) {
      console.error('Logout error:', error);
      setError('Failed to logout. Please try again.');
    } finally {
      setIsLoggingOut(false);
    }
  }, [router, updateUser]);

  const updateProfile = useCallback(
    async (data: UpdateProfileRequest): Promise<boolean> => {
      setIsUpdatingProfile(true);
      setError(null);
      try {
        const updatedProfile = await AuthService.updateProfile(data);
        if (updatedProfile) {
          const updatedUser = user ? { ...user, ...updatedProfile } : null;
          updateUser(updatedUser);
          return true;
        }
        setError('Failed to update profile');
        return false;
      } catch (error) {
        console.error('Profile update error:', error);
        setError('Failed to update profile. Please try again.');
        return false;
      } finally {
        setIsUpdatingProfile(false);
      }
    },
    [updateUser, user]
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
      logout,
      refreshSession,
      updateProfile,
      clearError,
    ]
  );

  return (
    <AuthErrorBoundary>
      <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
    </AuthErrorBoundary>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
