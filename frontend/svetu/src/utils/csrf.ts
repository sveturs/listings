import config from '@/config';

let csrfToken: string | null = null;

/**
 * Получает CSRF токен из API
 */
export async function getCSRFToken(): Promise<string> {
  if (csrfToken) {
    return csrfToken;
  }

  try {
    const response = await fetch(`${config.getApiUrl()}/api/v1/csrf-token`, {
      method: 'GET',
      credentials: 'include', // Важно для получения cookies
    });

    if (!response.ok) {
      throw new Error(`Failed to get CSRF token: ${response.status}`);
    }

    const data = await response.json();
    csrfToken = data.csrf_token;

    return csrfToken!;
  } catch (error) {
    console.error('Error getting CSRF token:', error);
    throw error;
  }
}

/**
 * Очищает кэшированный CSRF токен
 */
export function clearCSRFToken(): void {
  csrfToken = null;
}
