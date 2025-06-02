import configManager from '@/config';
import type {
  SessionResponse,
  UserUpdate,
  UpdateProfileRequest,
  RegisterUserRequest,
} from '@/types/auth';

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

interface LoginRequest {
  email: string;
  password: string;
}

interface LoginResponse {
  message: string;
  user: {
    id: number;
    name: string;
    email: string;
    provider: string;
    phone?: string;
    city?: string;
    country?: string;
    picture_url?: string;
    is_admin?: boolean;
  };
  token: string;
}

export class AuthService {
  private static abortControllers = new Map<string, AbortController>();
  private static rateLimiter = new Map<
    string,
    { count: number; resetTime: number }
  >();
  private static csrfToken: string | null = null;
  private static jwtToken: string | null = null;

  // Rate limiting configuration
  private static readonly RATE_LIMITS = {
    login: { maxAttempts: 5, windowMs: 15 * 60 * 1000 }, // 5 attempts per 15 minutes
    register: { maxAttempts: 3, windowMs: 60 * 60 * 1000 }, // 3 attempts per hour
  };

  static cleanup(): void {
    this.abortControllers.forEach((controller) => controller.abort());
    this.abortControllers.clear();
  }

  // JWT token management
  private static setJwtToken(token: string | null): void {
    this.jwtToken = token;
    if (token) {
      try {
        localStorage.setItem('jwt_token', token);
      } catch (error) {
        console.warn('Failed to save JWT token to localStorage:', error);
      }
    } else {
      try {
        localStorage.removeItem('jwt_token');
      } catch (error) {
        console.warn('Failed to remove JWT token from localStorage:', error);
      }
    }
  }

  private static getJwtToken(): string | null {
    if (this.jwtToken) {
      return this.jwtToken;
    }

    try {
      const token = localStorage.getItem('jwt_token');
      if (token) {
        this.jwtToken = token;
        return token;
      }
    } catch (error) {
      console.warn('Failed to read JWT token from localStorage:', error);
    }

    return null;
  }

  private static getAuthHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    const jwtToken = this.getJwtToken();
    if (jwtToken) {
      headers['Authorization'] = `Bearer ${jwtToken}`;
    }

    return headers;
  }

  // Get or fetch CSRF token
  private static async getCsrfToken(): Promise<string> {
    if (this.csrfToken) {
      return this.csrfToken;
    }

    try {
      const response = await fetch(`${API_BASE}/api/v1/csrf-token`, {
        method: 'GET',
        credentials: 'include',
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
      const headers: HeadersInit = {};
      const jwtToken = this.getJwtToken();
      if (jwtToken) {
        headers['Authorization'] = `Bearer ${jwtToken}`;
      }

      const response = await fetch(`${API_BASE}/auth/session`, {
        method: 'GET',
        credentials: 'include',
        headers,
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error('Failed to fetch session');
      }

      const data = await response.json();
      this.abortControllers.delete('session');
      return data;
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
      await fetch(`${API_BASE}/auth/logout`, {
        method: 'GET',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // Очищаем JWT токен при выходе
      this.setJwtToken(null);
    }
  }

  static async loginWithGoogle(
    returnTo?: string,
    redirect = true
  ): Promise<string> {
    const params = returnTo ? `?returnTo=${encodeURIComponent(returnTo)}` : '';
    const url = `${API_BASE}/auth/google${params}`;

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

      const result = await response.json();
      this.abortControllers.delete('updateProfile');
      return result;
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Update profile request was cancelled');
        return null;
      }
      console.error('Profile update error:', error);
      throw error;
    }
  }

  static async register(data: RegisterUserRequest): Promise<void> {
    // Check rate limiting
    if (!this.checkRateLimit('register')) {
      throw new AuthError('users.errors.tooManyAttempts');
    }

    const controller = this.getAbortController('register');
    const csrfToken = await this.getCsrfToken();

    try {
      const response = await fetch(`${API_BASE}/api/v1/users/register`, {
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

      this.abortControllers.delete('register');
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Register request was cancelled');
        return;
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
      const response = await fetch(`${API_BASE}/api/v1/users/login`, {
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

      // Сохраняем JWT токен если он есть в ответе
      if (result.token) {
        this.setJwtToken(result.token);
      }

      this.abortControllers.delete('login');
      return result;
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
