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
import { AuthErrorBoundaryWrapper } from '@/components/AuthErrorBoundaryWrapper';
import type { User, UpdateProfileRequest } from '@/types/auth';
import { tokenManager } from '@/utils/tokenManager';

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
      // Если уже идет обновление сессии, не запускаем новое
      if (isRefreshingSession) {
        console.log(
          '[AuthContext] Session refresh already in progress, skipping'
        );
        return;
      }

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
      } else {
        // Даже для фоновых обновлений устанавливаем флаг, чтобы предотвратить параллельные вызовы
        setIsRefreshingSession(true);
      }

      try {
        for (let i = 0; i < retries; i++) {
          try {
            // Сначала пытаемся восстановить сессию через JWT
            console.log(
              '[AuthContext] Attempting to restore session via JWT...'
            );
            const session = await AuthService.restoreSession();
            if (session && session.authenticated && session.user) {
              console.log('[AuthContext] JWT session restored successfully');
              updateUser(session.user);
              setError(null);
            } else {
              // Если восстановление не удалось, пытаемся получить сессию обычным способом
              console.log(
                '[AuthContext] JWT restore failed, trying fallback session...'
              );
              const fallbackSession = await AuthService.getSession();
              if (fallbackSession.authenticated && fallbackSession.user) {
                console.log(
                  '[AuthContext] Fallback session restored successfully'
                );
                updateUser(fallbackSession.user);
                setError(null);
              } else {
                console.log('[AuthContext] No valid session found');
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
    [updateUser, REFRESH_COOLDOWN, isRefreshingSession]
  );

  // Обработка OAuth токена из URL
  useEffect(() => {
    const handleOAuthToken = async () => {
      console.log('[AuthContext] Checking for OAuth token in URL...');

      // Проверяем наличие токена в URL (для OAuth callback)
      if (typeof window !== 'undefined') {
        const urlParams = new URLSearchParams(window.location.search);
        console.log(
          '[AuthContext] Current URL search:',
          window.location.search
        );

        // Backend отправляет токен как auth_token
        const authToken = urlParams.get('auth_token') || urlParams.get('token');

        if (authToken) {
          console.log(
            '[AuthContext] Found OAuth token in URL:',
            authToken.substring(0, 30) + '...'
          );
          console.log('[AuthContext] Token length:', authToken.length);

          // Сохраняем токен
          tokenManager.setAccessToken(authToken);
          console.log('[AuthContext] Token saved to tokenManager');

          // Проверяем что токен действительно сохранен
          const savedToken = tokenManager.getAccessToken();
          console.log(
            '[AuthContext] Verification - token retrieved:',
            savedToken ? 'Success' : 'Failed'
          );

          // Удаляем токен из URL для безопасности
          urlParams.delete('auth_token');
          urlParams.delete('token');
          const newUrl = `${window.location.pathname}${urlParams.toString() ? '?' + urlParams.toString() : ''}`;
          window.history.replaceState({}, document.title, newUrl);
          console.log('[AuthContext] Token removed from URL for security');

          // Сразу обновляем сессию после получения токена
          console.log(
            '[AuthContext] Starting session refresh with new token...'
          );
          await refreshSession(1, false);
        } else {
          console.log('[AuthContext] No OAuth token found in URL');
        }
      }
    };

    handleOAuthToken();
  }, [refreshSession]); // Выполняется только при монтировании

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
    // const hasRefreshToken = document.cookie.includes('refresh_token=');

    if (hasValidCache) {
      // Проверяем токен в кеше
      const cachedUserJson = storageUtils.getItem('svetu_user');
      if (cachedUserJson) {
        try {
          const cachedUser = JSON.parse(cachedUserJson);
          if (cachedUser && cachedUser.accessToken) {
            // Проверяем токен
            if (tokenManager.isTokenExpired(cachedUser.accessToken)) {
              // Токен истек, обновляем в фоне
              setTimeout(() => refreshSession(3, true), 100);
            } else {
              console.log(
                '[AuthContext] Cached token is still valid, skipping refresh'
              );
            }
          }
        } catch (error) {
          console.warn('Failed to parse cached user:', error);
        }
      }
    } else {
      // Если нет валидного кеша, делаем полную проверку
      // Это покрывает случаи:
      // - есть refresh token (после Google OAuth)
      // - нет ни кеша, ни refresh token
      refreshSession();
    }

    // Cleanup on unmount
    return () => {
      AuthService.cleanup();
    };
  }, [refreshSession, storageUtils]);

  // Redirect to login page (matches the interface)
  const login = useCallback(
    (returnTo?: string) => {
      const returnPath = returnTo || window.location.pathname;
      router.push(`/auth/login?returnTo=${encodeURIComponent(returnPath)}`);
    },
    [router]
  );

  const _loginWithGoogle = useCallback((returnTo?: string) => {
    try {
      AuthService.loginWithGoogle(returnTo);
    } catch (error) {
      console.error('Google login error:', error);
      setError('Failed to initiate Google login. Please try again.');
    }
  }, []);

  const logout = useCallback(async () => {
    setIsLoggingOut(true);
    setError(null);
    try {
      await AuthService.logout();
      updateUser(null);

      // Очищаем localStorage и sessionStorage, но сохраняем важные данные
      const locale = localStorage.getItem('NEXT_LOCALE');
      // НЕ сохраняем корзину - она должна быть привязана к пользователю!

      // Очищаем токен через tokenManager чтобы он удалился из sessionStorage
      tokenManager.clearTokens();

      localStorage.clear();
      sessionStorage.clear();

      // Восстанавливаем важные данные
      if (locale) {
        localStorage.setItem('NEXT_LOCALE', locale);
      }
      // НЕ восстанавливаем корзину - при выходе корзина должна очищаться

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
    <AuthErrorBoundaryWrapper>
      <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
    </AuthErrorBoundaryWrapper>
  );
}

export function useAuthContext() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuthContext must be used within an AuthProvider');
  }
  return context;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
