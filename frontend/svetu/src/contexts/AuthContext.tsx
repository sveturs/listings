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
  isRefreshingSession: boolean;
  error: string | null;
  login: (returnTo?: string) => void;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
  updateProfile: (data: UpdateProfileRequest) => Promise<boolean>;
  updateUser: (userData: User | null) => void;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  // Безопасная инициализация с кешированным состоянием из sessionStorage
  const [user, setUser] = useState<User | null>(() => {
    if (typeof window !== 'undefined') {
      try {
        const cached = sessionStorage.getItem('svetu_user');
        if (cached) {
          const parsedUser = JSON.parse(cached);
          // Проверяем, что объект пользователя имеет минимально необходимые поля
          if (
            parsedUser &&
            typeof parsedUser === 'object' &&
            parsedUser.id &&
            parsedUser.email
          ) {
            return parsedUser;
          }
        }
      } catch (error) {
        console.warn(
          'Failed to parse cached user data, clearing cache:',
          error
        );
        // Очищаем поврежденный кеш
        try {
          sessionStorage.removeItem('svetu_user');
        } catch {
          // Игнорируем ошибки очистки
        }
      }
    }
    return null;
  });

  // Если есть валидный кешированный пользователь, начинаем с false, иначе с true
  const [isLoading, setIsLoading] = useState(() => {
    if (typeof window !== 'undefined') {
      try {
        const cached = sessionStorage.getItem('svetu_user');
        if (cached) {
          const parsedUser = JSON.parse(cached);
          // Проверяем валидность кешированных данных
          if (
            parsedUser &&
            typeof parsedUser === 'object' &&
            parsedUser.id &&
            parsedUser.email
          ) {
            return false; // Данные валидны, начинаем без загрузки
          }
        }
      } catch {
        // Если не можем прочитать/парсить - требуется загрузка
      }
    }
    return true; // По умолчанию показываем загрузку
  });
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [isUpdatingProfile, setIsUpdatingProfile] = useState(false);
  const [isRefreshingSession, setIsRefreshingSession] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Cooldown для предотвращения частых вызовов refreshSession
  const lastRefreshAttempt = useRef<number>(0);
  const REFRESH_COOLDOWN = 5000; // 5 секунд

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Безопасная работа с sessionStorage с fallback механизмами
  const storageUtils = useMemo(() => {
    // Проверяем доступность sessionStorage
    const isStorageAvailable = (() => {
      if (typeof window === 'undefined') return false;
      try {
        const testKey = '__storage_test__';
        sessionStorage.setItem(testKey, 'test');
        sessionStorage.removeItem(testKey);
        return true;
      } catch {
        return false;
      }
    })();

    return {
      isAvailable: isStorageAvailable,

      getItem: (key: string): string | null => {
        if (!isStorageAvailable) return null;
        try {
          return sessionStorage.getItem(key);
        } catch (error) {
          console.warn(`Failed to read from sessionStorage (${key}):`, error);
          return null;
        }
      },

      setItem: (key: string, value: string): boolean => {
        if (!isStorageAvailable) {
          console.warn('SessionStorage is not available, skipping cache');
          return false;
        }
        try {
          sessionStorage.setItem(key, value);
          return true;
        } catch (error) {
          console.warn(`Failed to write to sessionStorage (${key}):`, error);
          // Попытка очистить место, если ошибка связана с переполнением
          if (error instanceof Error && error.name === 'QuotaExceededError') {
            try {
              sessionStorage.clear();
              sessionStorage.setItem(key, value);
              console.info('Cleared sessionStorage and retried');
              return true;
            } catch {
              console.error('Failed to clear and retry sessionStorage');
            }
          }
          return false;
        }
      },

      removeItem: (key: string): boolean => {
        if (!isStorageAvailable) return false;
        try {
          sessionStorage.removeItem(key);
          return true;
        } catch (error) {
          console.warn(`Failed to remove from sessionStorage (${key}):`, error);
          return false;
        }
      },
    };
  }, []);

  // Функция для кеширования пользователя в sessionStorage
  const cacheUser = useCallback(
    (userData: User | null) => {
      if (userData) {
        const success = storageUtils.setItem(
          'svetu_user',
          JSON.stringify(userData)
        );
        if (!success) {
          console.warn('User data was not cached due to storage issues');
        }
      } else {
        storageUtils.removeItem('svetu_user');
      }
    },
    [storageUtils]
  );

  // Обертка для setUser с кешированием
  const updateUser = useCallback(
    (userData: User | null) => {
      setUser(userData);
      cacheUser(userData);
    },
    [cacheUser]
  );

  const refreshSession = useCallback(
    async (retries = 3, skipLoadingState = false) => {
      const now = Date.now();
      if (now - lastRefreshAttempt.current < REFRESH_COOLDOWN) {
        if (process.env.NODE_ENV === 'development') {
          console.log('RefreshSession skipped due to cooldown');
        }
        return;
      }
      lastRefreshAttempt.current = now;

      // Устанавливаем состояние загрузки только если это не фоновое обновление
      if (!skipLoadingState) {
        setIsRefreshingSession(true);
      }

      try {
        for (let i = 0; i < retries; i++) {
          try {
            // Сначала пытаемся восстановить сессию через JWT
            const session = await AuthService.restoreSession();
            if (session && session.authenticated && session.user) {
              updateUser(session.user);
              setError(null);
            } else {
              // Если восстановление не удалось, пытаемся получить сессию обычным способом
              const fallbackSession = await AuthService.getSession();
              if (fallbackSession.authenticated && fallbackSession.user) {
                updateUser(fallbackSession.user);
                setError(null);
              } else {
                updateUser(null);
              }
            }
            setIsLoading(false);
            return;
          } catch (error) {
            console.error(
              `Session refresh error (attempt ${i + 1}/${retries}):`,
              error
            );
            if (i === retries - 1) {
              setError(
                'Failed to load session. Please try refreshing the page.'
              );
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
      } finally {
        setIsRefreshingSession(false);
      }
    },
    [updateUser, REFRESH_COOLDOWN]
  );

  useEffect(() => {
    // Инициализируем TokenManager
    AuthService.initializeTokenManager();

    // Проверяем флаг logout
    const logoutFlag = sessionStorage.getItem('svetu_logout_flag');
    if (logoutFlag === 'true') {
      // Пользователь вышел из системы, очищаем флаг и не восстанавливаем сессию
      sessionStorage.removeItem('svetu_logout_flag');
      storageUtils.removeItem('svetu_user');
      setIsLoading(false);
      return;
    }

    // Проверяем наличие валидного кешированного пользователя при первой загрузке
    const cachedData = storageUtils.getItem('svetu_user');
    let hasValidCache = false;

    if (cachedData) {
      try {
        const parsedUser = JSON.parse(cachedData);
        hasValidCache =
          parsedUser &&
          typeof parsedUser === 'object' &&
          parsedUser.id &&
          parsedUser.email;
      } catch {
        // Поврежденный кеш, очищаем его
        storageUtils.removeItem('svetu_user');
      }
    }

    // Проверяем, есть ли refresh token cookie (возможно после Google OAuth)
    const hasRefreshToken = document.cookie.includes('refresh_token=');

    if (hasValidCache) {
      // Проверяем актуальность в фоне, не блокируя UI (skipLoadingState = true)
      setTimeout(() => refreshSession(3, true), 100);
    } else if (hasRefreshToken) {
      // Если есть refresh token, но нет кеша (возможно после Google OAuth)
      console.log(
        '[AuthContext] Detected refresh token, attempting to restore session'
      );
      refreshSession();
    } else {
      // Если нет ни кеша, ни refresh token, делаем полную проверку с loading state
      refreshSession();
    }

    // Cleanup on unmount
    return () => {
      AuthService.cleanup();
    };
  }, [refreshSession, storageUtils]);

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

      // Устанавливаем флаг в sessionStorage чтобы предотвратить автоматическое восстановление
      sessionStorage.setItem('svetu_logout_flag', 'true');

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
          console.log('Profile update response:', updatedProfile);

          // Обновляем пользователя с новыми данными от сервера
          const updatedUser = user ? { ...user, ...updatedProfile } : null;
          console.log('Updated user data:', updatedUser);
          updateUser(updatedUser);

          // Принудительно получаем свежие данные пользователя с сервера
          // Сбрасываем cooldown для принудительного обновления
          lastRefreshAttempt.current = 0;
          try {
            const session = await AuthService.getSession();
            console.log('Fresh session data:', session);
            if (session.authenticated && session.user) {
              updateUser(session.user);
            }
          } catch (error) {
            console.warn(
              'Failed to refresh session after profile update:',
              error
            );
            // Не показываем ошибку пользователю, так как основное обновление прошло успешно
          }

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
      isRefreshingSession,
      error,
      login,
      logout,
      refreshSession,
      updateProfile,
      updateUser,
      clearError,
    }),
    [
      user,
      isLoading,
      isLoggingOut,
      isUpdatingProfile,
      isRefreshingSession,
      error,
      login,
      logout,
      refreshSession,
      updateProfile,
      updateUser,
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
