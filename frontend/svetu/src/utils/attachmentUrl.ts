import configManager from '@/config';
import { tokenManager } from '@/utils/tokenManager';

/**
 * Получает защищенный URL для скачивания файла вложения
 * Использует новый эндпоинт с проверкой авторизации
 */
export function getSecureAttachmentUrl(attachmentId: number): string {
  const baseUrl = configManager.getApiUrl();
  const token = tokenManager.getAccessToken();

  // Формируем URL с токеном авторизации в query параметре
  // Это позволит браузеру скачивать файлы напрямую
  const url = `${baseUrl}/marketplace/chat/attachments/${attachmentId}/download`;

  // Для изображений и других файлов, которые нужно отображать inline,
  // добавляем токен в URL для прямого доступа
  if (token) {
    return `${url}?token=${encodeURIComponent(token)}`;
  }

  return url;
}

/**
 * Получает URL для отображения вложения
 * Для изображений возвращает защищенный URL, для остальных - публичный
 */
export function getAttachmentDisplayUrl(attachment: {
  id: number;
  public_url: string;
  file_type: string;
}): string {
  // Для изображений используем защищенный URL
  if (attachment.file_type === 'image') {
    return getSecureAttachmentUrl(attachment.id);
  }

  // Для других файлов используем публичный URL из MinIO
  if (attachment.public_url.startsWith('/chat-files/')) {
    return configManager.buildImageUrl(attachment.public_url);
  }

  return attachment.public_url;
}

/**
 * Проверяет, нужно ли использовать защищенный URL для данного типа файла
 */
export function requiresSecureUrl(_fileType: string): boolean {
  // Для всех типов файлов в чате теперь используем защищенный доступ
  return true;
}

/**
 * Получает URL для скачивания файла
 */
export function getDownloadUrl(attachmentId: number): string {
  return getSecureAttachmentUrl(attachmentId);
}
