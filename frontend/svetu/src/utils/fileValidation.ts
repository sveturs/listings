// Разрешенные типы файлов
export const ALLOWED_FILE_TYPES = {
  image: ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'],
  video: [
    'video/mp4',
    'video/mpeg',
    'video/quicktime',
    'video/x-msvideo',
    'video/webm',
  ],
  document: [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'text/plain',
  ],
};

// Максимальные размеры файлов в байтах
export const MAX_FILE_SIZES = {
  image: 10 * 1024 * 1024, // 10MB
  video: 100 * 1024 * 1024, // 100MB
  document: 20 * 1024 * 1024, // 20MB
};

// Определение типа файла
export function getFileCategory(
  mimeType: string
): 'image' | 'video' | 'document' | null {
  for (const [category, types] of Object.entries(ALLOWED_FILE_TYPES)) {
    if (types.includes(mimeType)) {
      return category as 'image' | 'video' | 'document';
    }
  }

  // Проверка по префиксу MIME типа
  if (mimeType.startsWith('image/')) return 'image';
  if (mimeType.startsWith('video/')) return 'video';
  if (mimeType.startsWith('application/') || mimeType.startsWith('text/'))
    return 'document';

  return null;
}

// Валидация типа файла
export function validateFileType(file: File): string | null {
  const category = getFileCategory(file.type);

  if (!category) {
    return `Тип файла "${file.type}" не поддерживается`;
  }

  const allowedTypes = ALLOWED_FILE_TYPES[category];
  if (!allowedTypes.includes(file.type)) {
    return `Тип файла "${file.type}" не разрешен для категории "${category}"`;
  }

  return null;
}

// Валидация размера файла
export function validateFileSize(file: File): string | null {
  const category = getFileCategory(file.type);

  if (!category) {
    return 'Невозможно определить категорию файла';
  }

  const maxSize = MAX_FILE_SIZES[category];
  if (file.size > maxSize) {
    const maxSizeMB = Math.round(maxSize / (1024 * 1024));
    const fileSizeMB = (file.size / (1024 * 1024)).toFixed(1);
    return `Файл слишком большой: ${fileSizeMB}MB (максимум ${maxSizeMB}MB)`;
  }

  if (file.size === 0) {
    return 'Файл пустой';
  }

  return null;
}

// Валидация имени файла
export function validateFileName(filename: string): string | null {
  if (!filename || filename.trim() === '') {
    return 'Имя файла не может быть пустым';
  }

  // Проверка на опасные расширения
  const dangerousExtensions = [
    '.exe',
    '.bat',
    '.cmd',
    '.com',
    '.scr',
    '.vbs',
    '.js',
    '.jar',
    '.app',
    '.deb',
    '.rpm',
    '.dmg',
    '.pkg',
    '.run',
  ];

  const ext = filename.substring(filename.lastIndexOf('.')).toLowerCase();
  if (dangerousExtensions.includes(ext)) {
    return `Расширение файла "${ext}" не разрешено из соображений безопасности`;
  }

  // Проверка на недопустимые символы
  const invalidChars = ['<', '>', ':', '"', '/', '\\', '|', '?', '*'];
  for (const char of invalidChars) {
    if (filename.includes(char)) {
      return `Имя файла содержит недопустимый символ: "${char}"`;
    }
  }

  // Проверка длины
  if (filename.length > 255) {
    return 'Имя файла слишком длинное (максимум 255 символов)';
  }

  return null;
}

// Общая валидация файла
export function validateFile(file: File): string | null {
  // Валидация имени
  const nameError = validateFileName(file.name);
  if (nameError) return nameError;

  // Валидация типа
  const typeError = validateFileType(file);
  if (typeError) return typeError;

  // Валидация размера
  const sizeError = validateFileSize(file);
  if (sizeError) return sizeError;

  return null;
}

// Валидация массива файлов
export function validateFiles(files: File[]): {
  valid: File[];
  errors: Map<string, string>;
} {
  const valid: File[] = [];
  const errors = new Map<string, string>();

  // Ограничение на количество файлов
  if (files.length > 10) {
    errors.set('_total', 'Можно загрузить максимум 10 файлов за раз');
    return { valid, errors };
  }

  for (const file of files) {
    const error = validateFile(file);
    if (error) {
      errors.set(file.name, error);
    } else {
      valid.push(file);
    }
  }

  return { valid, errors };
}

// Форматирование размера файла
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Получение иконки для типа файла
export function getFileIcon(mimeType: string): string {
  const category = getFileCategory(mimeType);

  switch (category) {
    case 'image':
      return '🖼️';
    case 'video':
      return '🎥';
    case 'document':
      if (mimeType.includes('pdf')) return '📄';
      if (mimeType.includes('word')) return '📝';
      if (mimeType.includes('excel')) return '📊';
      return '📋';
    default:
      return '📎';
  }
}
