import { apiClient } from '@/services/api-client';
import configManager from '@/config';

interface TokenManagerConfig {
  onTokenRefreshed?: (accessToken: string) => void;
  onRefreshFailed?: (error: Error) => void;
}

class TokenManager {
  private accessToken: string | null = null;
  private refreshTimer: NodeJS.Timeout | null = null;
  private refreshPromise: Promise<string> | null = null;
  private config: TokenManagerConfig;

  constructor(config: TokenManagerConfig = {}) {
    this.config = config;
  }

  /**
   * Сохраняет access token в памяти
   */
  setAccessToken(token: string | null) {
    this.accessToken = token;
    if (process.env.NODE_ENV === 'development') {
      console.log(
        '[TokenManager] Access token set:',
        token ? 'token received' : 'token cleared'
      );
    }

    // Перезапускаем таймер обновления при установке нового токена
    if (token) {
      this.scheduleTokenRefresh();
    } else {
      this.clearRefreshTimer();
    }
  }

  /**
   * Возвращает текущий access token
   */
  getAccessToken(): string | null {
    return this.accessToken;
  }

  /**
   * Очищает access token и останавливает таймер обновления
   */
  clearTokens() {
    this.accessToken = null;
    this.clearRefreshTimer();
  }

  /**
   * Обновляет access token используя refresh token из httpOnly cookie
   */
  async refreshAccessToken(): Promise<string> {
    // Если уже идет процесс обновления, возвращаем существующий промис
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    this.refreshPromise = this.performRefresh();

    try {
      const token = await this.refreshPromise;
      return token;
    } finally {
      this.refreshPromise = null;
    }
  }

  /**
   * Выполняет запрос на обновление токена
   */
  private async performRefresh(): Promise<string> {
    try {
      console.log('[TokenManager] Attempting to refresh token...');
      const response = await fetch(
        `${configManager.getApiUrl()}/api/auth/refresh`,
        {
          method: 'POST',
          credentials: 'include', // Важно для отправки httpOnly cookie
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        // Если 401, значит refresh token невалидный или отозван
        if (response.status === 401) {
          console.log('[TokenManager] Refresh token is invalid or revoked');
          this.clearTokens();
          return '';
        }
        throw new Error(`Failed to refresh token: ${response.status}`);
      }

      const data = await response.json();
      if (process.env.NODE_ENV === 'development') {
        console.log('[TokenManager] Token refresh successful');
      }

      // Обрабатываем оба формата ответа
      let accessToken = data.access_token;

      // Если ответ обернут в { data: {...}, success: true }
      if (!accessToken && data.data && data.data.access_token) {
        accessToken = data.data.access_token;
      }

      if (!accessToken) {
        console.error(
          '[TokenManager] No access token in refresh response:',
          data
        );
        throw new Error('No access token in refresh response');
      }

      this.setAccessToken(accessToken);

      // Вызываем callback если он задан
      if (this.config.onTokenRefreshed) {
        this.config.onTokenRefreshed(accessToken);
      }

      return accessToken;
    } catch (error) {
      // Вызываем callback ошибки если он задан
      if (this.config.onRefreshFailed) {
        this.config.onRefreshFailed(error as Error);
      }

      // Очищаем токены при ошибке обновления
      this.clearTokens();

      throw error;
    }
  }

  /**
   * Планирует автоматическое обновление токена
   */
  private scheduleTokenRefresh() {
    this.clearRefreshTimer();

    // Декодируем токен для получения времени истечения
    if (!this.accessToken) return;

    try {
      const payload = this.decodeToken(this.accessToken);
      if (!payload || !payload.exp) return;

      // Вычисляем время до истечения токена
      const expiresAt = payload.exp * 1000; // exp в секундах, переводим в миллисекунды
      const now = Date.now();
      const timeUntilExpiry = expiresAt - now;

      // Обновляем токен за 5 минут до истечения
      const refreshTime = Math.max(0, timeUntilExpiry - 5 * 60 * 1000);

      if (refreshTime > 0) {
        this.refreshTimer = setTimeout(() => {
          this.refreshAccessToken().catch((error) => {
            console.error('Failed to refresh token:', error);
          });
        }, refreshTime);
      }
    } catch (error) {
      console.error('Failed to decode token:', error);
    }
  }

  /**
   * Очищает таймер обновления
   */
  private clearRefreshTimer() {
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
      this.refreshTimer = null;
    }
  }

  /**
   * Декодирует JWT токен
   */
  private decodeToken(
    token: string
  ): { exp?: number; [key: string]: unknown } | null {
    try {
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      );
      return JSON.parse(jsonPayload);
    } catch {
      return null;
    }
  }

  /**
   * Проверяет, истек ли токен
   */
  isTokenExpired(token?: string): boolean {
    const tokenToCheck = token || this.accessToken;
    if (!tokenToCheck) return true;

    try {
      const payload = this.decodeToken(tokenToCheck);
      if (!payload || !payload.exp) return true;

      const now = Date.now() / 1000;
      return payload.exp < now;
    } catch {
      return true;
    }
  }

  /**
   * Создает axios interceptor для автоматического добавления токена и обновления
   * @deprecated Не используется, так как apiClient имеет упрощенную реализацию
   */
  createAxiosInterceptor() {
    // Заглушка для совместимости
    // Реальная функциональность добавления токенов реализована
    // непосредственно в ApiClient через getAuthToken()
    console.log('[TokenManager] Interceptors not implemented in simplified apiClient');
  }
}

// Создаем единственный экземпляр TokenManager
export const tokenManager = new TokenManager({
  onTokenRefreshed: () => {
    if (process.env.NODE_ENV === 'development') {
      console.log('Token refreshed successfully');
    }
  },
  onRefreshFailed: (error) => {
    console.error('Token refresh failed:', error);
  },
});

// Экспортируем также класс для тестирования
export { TokenManager };
