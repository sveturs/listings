/**
 * DOMPurify utility для работы в SSR и клиенте
 *
 * Безопасная санитизация HTML с поддержкой серверного рендеринга
 */

import DOMPurify, { type Config } from 'dompurify';

/**
 * Санитизация HTML строки
 * Работает только на клиенте, на сервере возвращает пустую строку
 *
 * @param dirty - HTML строка для санитизации
 * @param config - Конфигурация DOMPurify (опционально)
 * @returns Безопасная HTML строка
 */
export function sanitizeHTML(dirty: string, config?: Config): string {
  // На сервере возвращаем пустую строку
  if (typeof window === 'undefined') {
    return '';
  }

  // На клиенте используем DOMPurify
  return DOMPurify.sanitize(dirty, config);
}

/**
 * Санитизация для текстового контента (удаляет все HTML теги)
 *
 * @param dirty - Текст который может содержать HTML
 * @returns Чистый текст без HTML тегов
 */
export function sanitizeText(dirty: string): string {
  return sanitizeHTML(dirty, {
    ALLOWED_TAGS: [],
    KEEP_CONTENT: true,
  });
}
