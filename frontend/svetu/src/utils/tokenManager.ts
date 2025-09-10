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
  private lastRefreshAttempt: number = 0;
  private refreshAttempts: number = 0;
  private rateLimitedUntil: number = 0; // Время до которого мы заблокированы из-за rate limit
  private readonly MIN_REFRESH_INTERVAL = 60000; // 60 секунд между попытками (увеличено для избежания rate limit)
  private readonly MAX_REFRESH_ATTEMPTS = 3;

  constructor(config: TokenManagerConfig = {}) {
    this.config = config;

    // Восстанавливаем токен из sessionStorage при инициализации
    if (typeof window !== 'undefined') {
      const savedToken = sessionStorage.getItem('svetu_access_token');

      if (savedToken) {
        // Check if it's an HS256 token and reject it
        try {
          const headerB64 = savedToken.split('.')[0];
          if (headerB64) {
            const headerJson = atob(
              headerB64.replace(/-/g, '+').replace(/_/g, '/')
            );
            const header = JSON.parse(headerJson);

            if (header.alg === 'HS256') {
              console.error(
                '[TokenManager] Constructor - found old HS256 token, removing'
              );
              sessionStorage.removeItem('svetu_access_token');
              // Also clear any other auth-related storage
              localStorage.removeItem('access_token');
              localStorage.removeItem('refresh_token');
              return;
            }
          }
        } catch (e) {
          console.warn(
            '[TokenManager] Constructor - could not parse token, removing',
            e
          );
          sessionStorage.removeItem('svetu_access_token');
          return;
        }

        this.accessToken = savedToken;
        // Проверяем, не истек ли токен
        if (!this.isTokenExpired()) {
          console.log(
            '[TokenManager] Constructor - RS256 token is valid, scheduling refresh'
          );
          this.scheduleTokenRefresh();
        } else {
          console.log(
            '[TokenManager] Constructor - RS256 token is expired, removing'
          );
          sessionStorage.removeItem('svetu_access_token');
          this.accessToken = null;
        }
      }

      // Также загружаем refresh токен если есть
      const savedRefreshToken = sessionStorage.getItem('svetu_refresh_token');
      if (savedRefreshToken) {
        console.log('[TokenManager] Constructor - loaded refresh token');
      }
    }
  }

  /**
   * Сохраняет refresh token в sessionStorage
   */
  setRefreshToken(token: string | null) {
    if (typeof window !== 'undefined') {
      if (token) {
        sessionStorage.setItem('svetu_refresh_token', token);
        console.log('[TokenManager] Refresh token saved');
      } else {
        sessionStorage.removeItem('svetu_refresh_token');
        console.log('[TokenManager] Refresh token cleared');
      }
    }
  }

  /**
   * Получает refresh token из sessionStorage
   */
  getRefreshToken(): string | null {
    if (typeof window !== 'undefined') {
      return sessionStorage.getItem('svetu_refresh_token');
    }
    return null;
  }

  /**
   * Сохраняет access token в памяти и sessionStorage
   */
  setAccessToken(token: string | null) {
    // Check if token is HS256 and reject it
    if (token) {
      try {
        const headerB64 = token.split('.')[0];
        const headerJson = atob(
          headerB64.replace(/-/g, '+').replace(/_/g, '/')
        );
        const header = JSON.parse(headerJson);

        if (header.alg === 'HS256') {
          console.error(
            '[TokenManager] Rejecting HS256 token - only RS256 tokens are supported'
          );
          this.clearTokens();
          return;
        }
      } catch (e) {
        console.error('[TokenManager] Failed to parse token header:', e);
      }
    }

    this.accessToken = token;

    // Сохраняем в sessionStorage для сохранения между перезагрузками
    if (typeof window !== 'undefined') {
      if (token) {
        sessionStorage.setItem('svetu_access_token', token);
      } else {
        sessionStorage.removeItem('svetu_access_token');
      }
    }

    if (process.env.NODE_ENV === 'development') {
      console.log(
        '[TokenManager] Access token set:',
        token
          ? `token received (length: ${token.length}, first 30 chars: ${token.substring(0, 30)}...)`
          : 'token cleared'
      );
    }

    // Перезапускаем таймер обновления при установке нового токена
    if (token) {
      this.scheduleTokenRefresh();
    } else {
      this.clearRefreshTimer();
    }

    // Генерируем событие для уведомления других компонентов об изменении токена
    if (typeof window !== 'undefined') {
      window.dispatchEvent(
        new CustomEvent('tokenChanged', {
          detail: {
            token: token,
            action: token ? 'set' : 'cleared',
          },
        })
      );
    }
  }

  /**
   * Возвращает текущий access token
   */
  getAccessToken(): string | null {
    if (process.env.NODE_ENV === 'development') {
      console.log(
        '[TokenManager] getAccessToken called, token:',
        this.accessToken
          ? `exists (length: ${this.accessToken.length})`
          : 'null'
      );
    }
    return this.accessToken;
  }

  /**
   * Очищает access и refresh токены и останавливает таймер обновления
   */
  clearTokens() {
    this.accessToken = null;
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('svetu_access_token');
      sessionStorage.removeItem('svetu_refresh_token');
    }
    this.clearRefreshTimer();
    this.refreshAttempts = 0;
    this.lastRefreshAttempt = 0;
    this.rateLimitedUntil = 0; // Сбрасываем rate limit
  }

  /**
   * Обновляет access token используя refresh token из httpOnly cookie
   */
  async refreshAccessToken(): Promise<string> {
    // Проверяем, не заблокированы ли мы из-за rate limit
    const now = Date.now();
    if (this.rateLimitedUntil > now) {
      const remainingTime = this.rateLimitedUntil - now;
      console.warn(
        `[TokenManager] Rate limited, waiting ${remainingTime}ms before retry`
      );
      return this.accessToken || '';
    }

    // Если уже идет процесс обновления, возвращаем существующий промис
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    // Проверяем, не слишком ли часто мы пытаемся обновить токен
    const timeSinceLastAttempt = now - this.lastRefreshAttempt;

    if (timeSinceLastAttempt < this.MIN_REFRESH_INTERVAL) {
      console.warn(
        `[TokenManager] Rate limit protection: ${timeSinceLastAttempt}ms since last attempt, waiting...`
      );
      // Возвращаем текущий токен если он еще валидный
      if (this.accessToken && !this.isTokenExpired()) {
        return this.accessToken;
      }
      // Иначе ждем перед повторной попыткой
      await new Promise((resolve) =>
        setTimeout(resolve, this.MIN_REFRESH_INTERVAL - timeSinceLastAttempt)
      );
    }

    // Проверяем количество попыток
    if (this.refreshAttempts >= this.MAX_REFRESH_ATTEMPTS) {
      console.error('[TokenManager] Max refresh attempts reached');
      this.clearTokens();
      throw new Error('Max refresh attempts reached');
    }

    this.lastRefreshAttempt = now;
    this.refreshPromise = this.performRefresh();

    try {
      const token = await this.refreshPromise;
      this.refreshAttempts = 0; // Сбрасываем счетчик при успехе
      this.rateLimitedUntil = 0; // Сбрасываем rate limit при успехе
      return token;
    } catch (error) {
      this.refreshAttempts++;
      throw error;
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

      // Получаем refresh токен
      const refreshToken = this.getRefreshToken();
      if (!refreshToken) {
        console.log('[TokenManager] No refresh token available');
        this.clearTokens();
        return '';
      }

      const response = await fetch(
        `${configManager.getApiUrl()}/api/v1/auth/refresh`,
        {
          method: 'POST',
          credentials: 'include', // Важно для cookies если используются
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${refreshToken}`, // Отправляем refresh токен в заголовке
          },
          body: JSON.stringify({
            refresh_token: refreshToken, // Также отправляем в теле запроса для совместимости
          }),
        }
      );

      if (!response.ok) {
        // Если 401 или 400, значит refresh token невалидный или отозван
        // Auth Service может вернуть 400 если токен отсутствует или некорректный
        if (response.status === 401 || response.status === 400) {
          console.log(
            '[TokenManager] Refresh token is invalid, missing or revoked'
          );
          this.clearTokens();
          return '';
        }

        // Если 429, добавляем экспоненциальную задержку
        if (response.status === 429) {
          const retryAfter = response.headers.get('Retry-After');
          const delay = retryAfter
            ? parseInt(retryAfter) * 1000
            : Math.min(
                this.MIN_REFRESH_INTERVAL * Math.pow(2, this.refreshAttempts),
                300000
              ); // макс 5 минут

          console.warn(
            `[TokenManager] Rate limited (429), retry after ${delay}ms`
          );

          // Устанавливаем время блокировки
          this.rateLimitedUntil = Date.now() + delay;

          // Планируем повторную попытку после задержки
          setTimeout(() => {
            this.rateLimitedUntil = 0; // Сбрасываем блокировку
            this.refreshAccessToken().catch(console.error);
          }, delay);

          return this.accessToken || '';
        }

        throw new Error(`Failed to refresh token: ${response.status}`);
      }

      const data = await response.json();
      if (process.env.NODE_ENV === 'development') {
        console.log('[TokenManager] Token refresh successful');
      }

      // Обрабатываем оба формата ответа
      let accessToken = data.access_token;
      let newRefreshToken = data.refresh_token;

      // Если ответ обернут в { data: {...}, success: true }
      if (!accessToken && data.data && data.data.access_token) {
        accessToken = data.data.access_token;
        newRefreshToken = data.data.refresh_token;
      }

      if (!accessToken) {
        console.error(
          '[TokenManager] No access token in refresh response:',
          data
        );
        throw new Error('No access token in refresh response');
      }

      // Сохраняем новые токены
      this.setAccessToken(accessToken);
      if (newRefreshToken) {
        this.setRefreshToken(newRefreshToken);
      }

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
   * Декодирует JWT токен
   */
  private decodeToken(
    token: string
  ): { exp?: number; iat?: number; sub?: string } | null {
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
   * Очищает таймер обновления токена
   */
  private clearRefreshTimer() {
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
      this.refreshTimer = null;
    }
  }

  /**
   * Сбрасывает счетчики rate limit (для отладки)
   */
  resetRateLimits() {
    this.lastRefreshAttempt = 0;
    this.refreshAttempts = 0;
    console.log('[TokenManager] Rate limits reset');
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

// Экспортируем отдельно для использования в static контексте
export const isTokenExpired = (token?: string) =>
  tokenManager.isTokenExpired(token);

// Экспортируем также класс для тестирования
export { TokenManager };

// Для отладки в development
if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
  (window as any).tokenManager = tokenManager;
}
