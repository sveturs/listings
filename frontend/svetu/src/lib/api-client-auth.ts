import { apiClient } from '@/services/api-client';

/**
 * @deprecated Этот модуль больше не нужен. BFF proxy автоматически добавляет авторизацию через httpOnly cookies.
 * Используйте apiClient напрямую из './api-client'.
 *
 * Этот реэкспорт оставлен для обратной совместимости и будет удален в будущем.
 */
export const apiClientAuth = apiClient;
