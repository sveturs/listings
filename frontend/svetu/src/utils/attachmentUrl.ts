import configManager from '@/config';

/**
 * Получает защищенный URL для скачивания файла вложения
 * Использует BFF proxy с автоматической авторизацией через cookies
 */
export function getSecureAttachmentUrl(attachmentId: number): string {
  // Используем BFF proxy - автоматически добавит cookies для авторизации
  return `/api/v2/c2c/chat/attachments/${attachmentId}/download`;
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
