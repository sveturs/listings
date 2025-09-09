import configManager from '@/config';
import { tokenManager } from '@/utils/tokenManager';
import type {
  RegisterUserRequest,
  SessionResponse,
  UpdateProfileRequest,
  UserUpdate,
} from '@/types/auth';
import type { components } from '@/types/generated/api';

type ApiLoginRequest =
  components['schemas']['internal_proj_users_handler.LoginRequest'];
type ApiAuthResponse =
  components['schemas']['internal_proj_users_handler.AuthResponse'];
type ApiSuccessResponse<T> =
  components['schemas']['backend_pkg_utils.SuccessResponseSwag'] & {
    data?: T;
  };

const API_BASE = configManager.getApiUrl();

// Кастомная ошибка для аутентификации, которая не будет логироваться браузером
class AuthError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'AuthError';
    // Предотвращаем отображение в консоли DevTools
    Object.defineProperty(this, 'stack', {
      get() {
        return undefined;
      },
      configurable: true,
    });
  }
}

// Используем сгенерированные типы вместо дублирования
type LoginRequest = ApiLoginRequest;
type LoginResponse = ApiAuthResponse;

export class AuthService {
  private static abortControllers = new Map<string, AbortController>();
  private static rateLimiter = new Map<
    string,
    { count: number; resetTime: number }
  >();
  private static csrfToken: string | null = null;

  // Rate limiting configuration
  private static readonly RATE_LIMITS = {
    login: { maxAttempts: 5, windowMs: 15 * 60 * 1000 }, // 5 attempts per 15 minutes
    register: { maxAttempts: 3, windowMs: 60 * 60 * 1000 }, // 3 attempts per hour
  };

  static cleanup(): void {
    this.abortControllers.forEach((controller) => controller.abort());
    this.abortControllers.clear();
  }

  static getAuthHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    // Добавляем JWT токен если есть
    const token = tokenManager.getAccessToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
      console.log(
        '[AuthService] Adding token to headers:',
        token.substring(0, 30) + '...'
      );
    } else {
      console.log('[AuthService] No token available for headers');
    }

    return headers;
  }

  // Инициализация TokenManager
  static initializeTokenManager(): void {
    // TokenManager инициализируется автоматически при импорте
  }

  // Попытка восстановить сессию при загрузке страницы
  static async restoreSession(): Promise<SessionResponse | null> {
    console.log('[AuthService] Attempting to restore session...');
    try {
      // Сначала проверяем, есть ли у нас валидный токен
      let accessToken = tokenManager.getAccessToken();

      // Если токен есть и он еще не истек, используем его
      if (accessToken && !tokenManager.isTokenExpired(accessToken)) {
        console.log('[AuthService] Using existing valid access token');
        try {
          const session = await this.getSession();
          if (session && session.authenticated) {
            return session;
          }
          // eslint-disable-next-line @typescript-eslint/no-unused-vars
        } catch (_sessionError) {
          console.log(
            '[AuthService] Session failed with existing token, will try refresh'
          );
        }
      }

      // Проверяем есть ли refresh токен
      const refreshToken = tokenManager.getRefreshToken();
      if (!refreshToken) {
        console.log(
          '[AuthService] No refresh token available, cannot restore session'
        );
        return null;
      }

      // Только если токен отсутствует или истек, пытаемся обновить
      console.log(
        '[AuthService] Access token expired or invalid, attempting refresh...'
      );
      accessToken = await tokenManager.refreshAccessToken();

      if (accessToken) {
        console.log('[AuthService] Access token obtained, fetching session...');
        // Если удалось получить access token, получаем сессию
        return await this.getSession();
      } else {
        console.log('[AuthService] No access token obtained from refresh');
      }
    } catch (error: any) {
      // Обрабатываем специфичные ошибки
      if (
        error.message?.includes('429') ||
        error.message?.includes('Rate limit')
      ) {
        console.warn('[AuthService] Rate limited, will retry later');
        // Не очищаем токены при rate limit, чтобы использовать текущий токен
        const currentToken = tokenManager.getAccessToken();
        if (currentToken && !tokenManager.isTokenExpired(currentToken)) {
          // Если токен еще валидный, используем его
          return await this.getSession();
        }
      } else if (error.message?.includes('Max refresh attempts')) {
        // Если достигнут лимит попыток, очищаем токены
        tokenManager.clearTokens();
        console.error('[AuthService] Max refresh attempts reached');
      } else {
        // Для других ошибок также очищаем токены
        tokenManager.clearTokens();
        console.log('[AuthService] Could not restore session:', error);
      }
    }

    return null;
  }

  // Get or fetch CSRF token
  static async getCsrfToken(): Promise<string> {
    // Сбрасываем токен, если запрос не удался ранее
    if (this.csrfToken && this.csrfToken.startsWith('client-')) {
      this.csrfToken = null;
    }

    if (this.csrfToken) {
      return this.csrfToken;
    }

    try {
      // Получаем JWT токен если есть
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = token
        ? {
            Authorization: `Bearer ${token}`,
          }
        : {};

      const response = await fetch(`${API_BASE}/api/v1/csrf-token`, {
        method: 'GET',
        credentials: 'include',
        headers,
      });

      if (response.ok) {
        const data = await response.json();
        this.csrfToken = data.csrf_token;
        return this.csrfToken || '';
      }
    } catch (error) {
      console.warn('Failed to fetch CSRF token:', error);
    }

    // Fallback: generate client-side token for basic protection
    this.csrfToken = `client-${Date.now()}-${Math.random().toString(36).substring(2)}`;
    return this.csrfToken;
  }

  // Check rate limiting
  private static checkRateLimit(action: 'login' | 'register'): boolean {
    const config = this.RATE_LIMITS[action];
    const key = `${action}_${this.getClientIdentifier()}`;
    const now = Date.now();

    const limit = this.rateLimiter.get(key);

    if (!limit || now > limit.resetTime) {
      // Reset or initialize
      this.rateLimiter.set(key, {
        count: 1,
        resetTime: now + config.windowMs,
      });
      return true;
    }

    if (limit.count >= config.maxAttempts) {
      return false;
    }

    limit.count++;
    return true;
  }

  // Get client identifier for rate limiting
  private static getClientIdentifier(): string {
    // Use sessionStorage for client identification
    let clientId = '';
    try {
      clientId = sessionStorage.getItem('client_id') || '';
      if (!clientId) {
        clientId = `client_${Date.now()}_${Math.random().toString(36).substring(2)}`;
        sessionStorage.setItem('client_id', clientId);
      }
    } catch {
      // Fallback if sessionStorage is not available
      clientId = `client_${Date.now()}_${Math.random().toString(36).substring(2)}`;
    }
    return clientId;
  }

  private static getAbortController(key: string): AbortController {
    // Cancel any existing request with the same key
    const existing = this.abortControllers.get(key);
    if (existing) {
      existing.abort();
    }

    // Create new controller
    const controller = new AbortController();
    this.abortControllers.set(key, controller);
    return controller;
  }

  static async getSession(): Promise<SessionResponse> {
    const controller = this.getAbortController('session');

    try {
      const response = await fetch(`${API_BASE}/api/v1/auth/session`, {
        method: 'GET',
        credentials: 'include',
        headers: this.getAuthHeaders(), // Добавляем JWT токен в заголовки
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error('Failed to fetch session');
      }

      const result =
        (await response.json()) as ApiSuccessResponse<SessionResponse>;
      this.abortControllers.delete('session');

      // Извлекаем данные из обертки
      const sessionData = result.data || (result as SessionResponse);
      return sessionData;
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Session request was cancelled');
        return { authenticated: false };
      }
      console.error('Session fetch error:', error);
      throw error;
    }
  }

  static async logout(): Promise<void> {
    try {
      // ВАЖНО: Сначала отправляем запрос с токеном для его отзыва на сервере
      await fetch(`${API_BASE}/api/v1/auth/logout`, {
        method: 'POST',
        credentials: 'include',
        headers: this.getAuthHeaders(), // Токен еще есть и будет отправлен
      });

      // Только после успешного logout на сервере очищаем токены локально
      tokenManager.clearTokens();
    } catch (error) {
      console.error('Logout error:', error);
      // В любом случае очищаем токены локально
      tokenManager.clearTokens();
    }
  }

  static async loginWithGoogle(
    returnTo?: string,
    redirect = true
  ): Promise<string> {
    // Get current locale from URL
    const locale =
      typeof window !== 'undefined'
        ? window.location.pathname.split('/')[1] || 'en'
        : 'en';

    // Build redirect URI with locale
    const redirectUri = `${window.location.origin}/${locale}/auth/oauth/google/callback`;

    // Save return URL for after OAuth
    if (returnTo && typeof window !== 'undefined') {
      sessionStorage.setItem('oauth_return_to', returnTo);
    }

    // Build OAuth URL with redirect_uri parameter
    const params = new URLSearchParams({
      redirect_uri: redirectUri,
    });

    const url = `${API_BASE}/api/v1/auth/google?${params.toString()}`;

    if (redirect && typeof window !== 'undefined') {
      window.location.href = url;
    }

    return url;
  }

  static async updateProfile(
    data: UpdateProfileRequest
  ): Promise<UserUpdate | null> {
    const controller = this.getAbortController('updateProfile');
    const csrfToken = await this.getCsrfToken();

    try {
      const response = await fetch(`${API_BASE}/api/v1/users/me`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          ...this.getAuthHeaders(),
          'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify(data),
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error('Failed to update profile');
      }

      const result = (await response.json()) as ApiSuccessResponse<UserUpdate>;
      this.abortControllers.delete('updateProfile');

      // Извлекаем данные из обертки
      return result.data || (result as UserUpdate);
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Update profile request was cancelled');
        return null;
      }
      console.error('Profile update error:', error);
      throw error;
    }
  }

  static async register(data: RegisterUserRequest): Promise<LoginResponse> {
    // Check rate limiting
    if (!this.checkRateLimit('register')) {
      throw new AuthError('users.errors.tooManyAttempts');
    }

    const controller = this.getAbortController('register');
    const csrfToken = await this.getCsrfToken();

    try {
      const response = await fetch(`${API_BASE}/api/v1/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify(data),
        credentials: 'include',
        signal: controller.signal,
      });

      if (!response.ok) {
        const errorData = await response.json();
        const errorMessage = errorData.error || 'Registration failed';

        // Создаем кастомную ошибку, которая не будет логироваться браузером
        throw new AuthError(errorMessage);
      }

      const result =
        (await response.json()) as ApiSuccessResponse<LoginResponse>;

      // Извлекаем данные из обертки
      const registerData = result.data || (result as LoginResponse);

      // Сохраняем JWT токен после успешной регистрации
      if (registerData.access_token) {
        tokenManager.setAccessToken(registerData.access_token);
      }

      this.abortControllers.delete('register');
      return registerData;
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Register request was cancelled');
        throw error;
      }
      // Не логируем AuthError и translation keys в консоль
      if (
        error instanceof Error &&
        error.name !== 'AuthError' &&
        !error.message.startsWith('users.')
      ) {
        console.error('Registration error:', error);
      }
      throw error;
    }
  }

  static async login(data: LoginRequest): Promise<LoginResponse> {
    // Check rate limiting
    if (!this.checkRateLimit('login')) {
      throw new AuthError('users.errors.tooManyAttempts');
    }

    const controller = this.getAbortController('login');
    const csrfToken = await this.getCsrfToken();

    try {
      const response = await fetch(`${API_BASE}/api/v1/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify(data),
        credentials: 'include',
        signal: controller.signal,
      });

      if (!response.ok) {
        const errorData = await response.json();
        const errorMessage = errorData.error || 'Login failed';

        // Создаем кастомную ошибку, которая не будет логироваться браузером
        throw new AuthError(errorMessage);
      }

      const result = await response.json();

      // Извлекаем данные из обертки - сервер возвращает {data: {...}, success: true}
      const loginData = result.data as LoginResponse;

      // Сохраняем JWT токены
      if (loginData) {
        if (loginData.access_token) {
          tokenManager.setAccessToken(loginData.access_token);
          // Также сохраняем в localStorage для translationAdminApi
          localStorage.setItem('access_token', loginData.access_token);
        } else {
          console.error('[AuthService] No access_token in login response');
        }

        // Сохраняем refresh_token если есть
        // Auth Service возвращает refresh_token в ответе
        if ((loginData as any).refresh_token) {
          tokenManager.setRefreshToken((loginData as any).refresh_token);
          console.log('[AuthService] Refresh token saved from login response');
        } else {
          console.warn('[AuthService] No refresh_token in login response');
        }
      }

      this.abortControllers.delete('login');
      return loginData;
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Login request was cancelled');
        throw error;
      }
      // Не логируем AuthError и translation keys в консоль
      if (
        error instanceof Error &&
        error.name !== 'AuthError' &&
        !error.message.startsWith('users.')
      ) {
        console.error('Login error:', error);
      }
      throw error;
    }
  }
}
