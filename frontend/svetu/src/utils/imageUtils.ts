/**
 * Утилиты для безопасной работы с изображениями
 */

// Список разрешенных доменов для изображений
const ALLOWED_DOMAINS = [
  'localhost',
  'svetu.rs',
  '100.88.44.15',
  'lh3.googleusercontent.com',
];

// Список разрешенных протоколов
const ALLOWED_PROTOCOLS = ['http:', 'https:'];

/**
 * Проверяет, является ли URL безопасным для загрузки изображения
 */
export function isSafeImageUrl(url: string): boolean {
  if (!url) return false;

  try {
    const urlObj = new URL(url);

    // Проверяем протокол
    if (!ALLOWED_PROTOCOLS.includes(urlObj.protocol)) {
      return false;
    }

    // Проверяем домен
    const hostname = urlObj.hostname;
    const isAllowed = ALLOWED_DOMAINS.some((domain) => {
      // Проверяем точное совпадение или поддомен
      return hostname === domain || hostname.endsWith(`.${domain}`);
    });

    return isAllowed;
  } catch (error) {
    console.warn('Invalid URL:', url, error);
    return false;
  }
}

/**
 * Возвращает безопасный URL изображения или placeholder
 */
export function getSafeImageUrl(url: string | null | undefined): string | null {
  if (!url) return null;

  if (isSafeImageUrl(url)) {
    return url;
  }

  // Логируем попытки использования внешних изображений
  console.warn('Blocked external image URL:', url);

  // Возвращаем null, чтобы компонент показал placeholder
  return null;
}

/**
 * Фильтрует массив URL изображений, оставляя только безопасные
 */
export function filterSafeImageUrls(
  urls: (string | null | undefined)[]
): string[] {
  return urls
    .map((url) => getSafeImageUrl(url))
    .filter((url): url is string => url !== null);
}

/**
 * Placeholder изображение в формате base64 (1x1 прозрачный пиксель)
 */
export const PLACEHOLDER_IMAGE =
  'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';
