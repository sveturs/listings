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
import { TokenMigration } from '@/utils/tokenMigration';
// import { forceTokenCleanup } from '@/utils/forceTokenCleanup'; // Отключено - удалял валидные OAuth токены
import { logger } from '@/utils/logger';
import { decodeUserFromToken } from '@/utils/jwtDecode';

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
        logger.auth.warn(
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
        logger.auth.debug('Session refresh already in progress, skipping');
        return;
      }

      const now = Date.now();
      if (now - lastRefreshAttempt.current < REFRESH_COOLDOWN) {
        if (process.env.NODE_ENV === 'development') {
          logger.auth.debug('RefreshSession skipped due to cooldown');
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
            logger.auth.debug(
              '[AuthContext] Attempting to restore session via JWT...'
            );
            const session = await AuthService.restoreSession();
            if (session && session.authenticated && session.user) {
              logger.auth.debug('JWT session restored successfully');
              updateUser(session.user);
              setError(null);
            } else {
              // Если восстановление не удалось, значит нет валидной сессии
              logger.auth.debug('[AuthContext] No valid session to restore');
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
      logger.auth.debug('Checking for OAuth token in URL...');

      // Проверяем наличие токена в URL (для OAuth callback)
      if (typeof window !== 'undefined') {
        const urlParams = new URLSearchParams(window.location.search);
        logger.auth.debug(
          '[AuthContext] Current URL search:',
          window.location.search
        );

        // Backend отправляет токен как auth_token
        const authToken = urlParams.get('auth_token') || urlParams.get('token');

        if (authToken) {
          logger.auth.debug(
            '[AuthContext] Found OAuth token in URL:',
            authToken.substring(0, 30) + '...'
          );
          logger.auth.debug('Token length:', authToken.length);

          // Сохраняем токен
          tokenManager.setAccessToken(authToken);
          logger.auth.debug('Token saved to tokenManager');

          // Проверяем что токен действительно сохранен
          const savedToken = tokenManager.getAccessToken();
          logger.auth.debug(
            '[AuthContext] Verification - token retrieved:',
            savedToken ? 'Success' : 'Failed'
          );

          // Удаляем токен из URL для безопасности
          urlParams.delete('auth_token');
          urlParams.delete('token');
          const newUrl = `${window.location.pathname}${urlParams.toString() ? '?' + urlParams.toString() : ''}`;
          window.history.replaceState({}, document.title, newUrl);
          logger.auth.debug('Token removed from URL for security');

          // Сразу обновляем сессию после получения токена
          logger.auth.debug(
            '[AuthContext] Starting session refresh with new token...'
          );
          await refreshSession(1, false);
        } else {
          logger.auth.debug('No OAuth token found in URL');
        }
      }
    };

    handleOAuthToken();
  }, [refreshSession]); // Выполняется только при монтировании

  useEffect(() => {
    // ОТКЛЮЧЕНО: forceTokenCleanup удалял валидные токены после OAuth
    // const forceCleaned = forceTokenCleanup();
    // if (forceCleaned) {
    //   logger.auth.debug('Forced token cleanup performed');
    //   updateUser(null);
    //   tokenManager.clearTokens();
    // }

    // Проверяем и мигрируем старые токены
    const migrated = TokenMigration.runMigration();
    if (migrated) {
      logger.auth.debug(
        'Token migration performed, user needs to re-authenticate'
      );
      updateUser(null);
      setIsLoading(false);
      return;
    }

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
    let cachedUser = null;

    if (cachedData) {
      try {
        const parsedUser = JSON.parse(cachedData);
        hasValidCache =
          parsedUser &&
          typeof parsedUser === 'object' &&
          parsedUser.id &&
          parsedUser.email;

        if (hasValidCache) {
          cachedUser = parsedUser;
          // Немедленно устанавливаем пользователя из кеша
          updateUser(cachedUser);
          logger.auth.debug(
            '[AuthContext] Restored user from cache immediately'
          );
        }
      } catch {
        // Поврежденный кеш, очищаем его
        storageUtils.removeItem('svetu_user');
      }
    }

    // Проверяем наличие access token (может быть после OAuth)
    const currentToken = tokenManager.getAccessToken();

    if (currentToken && !tokenManager.isTokenExpired(currentToken)) {
      logger.auth.debug(
        '[AuthContext] Valid access token found, checking for user data'
      );
      
      // Если нет кешированного пользователя, пытаемся декодировать токен
      if (!hasValidCache) {
        const decodedUser = decodeUserFromToken(currentToken);
        if (decodedUser) {
          logger.auth.debug('[AuthContext] Decoded user from existing token:', decodedUser);
          updateUser(decodedUser);
          setIsLoading(false);
        }
      }
      
      // В любом случае обновляем сессию чтобы получить полные данные пользователя
      refreshSession();
    } else if (hasValidCache) {
      // Есть кешированный пользователь, но нужно проверить/обновить токен
      logger.auth.debug(
        '[AuthContext] User cache found, checking token validity'
      );
      setTimeout(() => refreshSession(3, true), 100);
    } else {
      // Нет ни токена, ни кеша - пытаемся восстановить через refresh token
      logger.auth.debug(
        '[AuthContext] No cache or token, attempting full session restore'
      );
      refreshSession();
    }

    // Cleanup on unmount
    return () => {
      AuthService.cleanup();
    };
  }, [refreshSession, storageUtils, updateUser]);

  // Listen for token changes from TokenManager
  useEffect(() => {
    const handleTokenChange = async (event: Event) => {
      const customEvent = event as CustomEvent;
      logger.auth.debug('Token changed event:', customEvent.detail);

      if (customEvent.detail.action === 'set') {
        // New token was set, refresh session to get user data
        logger.auth.debug('New token detected, refreshing session...');

        // Проверяем, есть ли уже пользователь в sessionStorage (например, после OAuth)
        const cachedUser = storageUtils.getItem('svetu_user');
        if (cachedUser) {
          try {
            const parsedUser = JSON.parse(cachedUser);
            if (parsedUser && parsedUser.id) {
              // Немедленно обновляем пользователя
              logger.auth.debug('[AuthContext] Found cached user data:', parsedUser);
              updateUser(parsedUser);
              setIsLoading(false);
              // Если есть кешированные данные, не нужно сразу делать запрос
              // Делаем его с задержкой для обновления данных
              setTimeout(() => refreshSession(3, true), 1000);
              return;
            }
          } catch (e) {
            logger.auth.error('Failed to parse cached user:', e);
          }
        }

        // Если нет кешированных данных, пытаемся декодировать токен
        logger.auth.debug('[AuthContext] No cached user, trying to decode token...');
        
        const token = tokenManager.getAccessToken();
        if (token) {
          const decodedUser = decodeUserFromToken(token);
          if (decodedUser) {
            logger.auth.debug('[AuthContext] Decoded user from token:', decodedUser);
            // Сохраняем в кеш и обновляем состояние
            updateUser(decodedUser);
            setIsLoading(false);
            // Запрашиваем полные данные с сервера с задержкой
            setTimeout(() => refreshSession(3, true), 500);
            return;
          }
        }
        
        // Если не удалось декодировать, запрашиваем сессию
        logger.auth.debug('[AuthContext] Could not decode token, fetching session...');
        await refreshSession(3, false); // Don't skip loading state
      } else if (customEvent.detail.action === 'cleared') {
        // Token was cleared, clear user state
        logger.auth.debug('Token cleared, clearing user state...');
        updateUser(null);
      }
    };

    // Add event listener
    if (typeof window !== 'undefined') {
      window.addEventListener('tokenChanged', handleTokenChange);
    }

    // Cleanup
    return () => {
      if (typeof window !== 'undefined') {
        window.removeEventListener('tokenChanged', handleTokenChange);
      }
    };
  }, [refreshSession, updateUser, storageUtils]);

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
          logger.auth.debug('Profile update response:', updatedProfile);

          // Обновляем пользователя с новыми данными от сервера
          const updatedUser = user ? { ...user, ...updatedProfile } : null;
          logger.auth.debug('Updated user data:', updatedUser);
          updateUser(updatedUser);

          // Принудительно получаем свежие данные пользователя с сервера
          // Сбрасываем cooldown для принудительного обновления
          lastRefreshAttempt.current = 0;
          try {
            const session = await AuthService.getSession();
            logger.auth.debug('Fresh session data:', session);
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
