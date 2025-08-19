import { useEffect, useRef } from 'react';

/**
 * Хук для безопасного создания и автоматической очистки URL объектов
 * Предотвращает утечки памяти при работе с URL.createObjectURL
 */
export function useObjectURL(
  file: File | Blob | null | undefined
): string | null {
  const urlRef = useRef<string | null>(null);

  useEffect(() => {
    // Очищаем предыдущий URL если он существует
    if (urlRef.current) {
      URL.revokeObjectURL(urlRef.current);
      urlRef.current = null;
    }

    // Создаем новый URL если файл существует
    if (file) {
      urlRef.current = URL.createObjectURL(file);
    }

    // Очистка при размонтировании или изменении файла
    return () => {
      if (urlRef.current) {
        URL.revokeObjectURL(urlRef.current);
        urlRef.current = null;
      }
    };
  }, [file]);

  return urlRef.current;
}

/**
 * Хук для управления массивом URL объектов
 * Автоматически очищает все URL при размонтировании
 */
export function useObjectURLs(files: (File | Blob)[]): string[] {
  const urlsRef = useRef<string[]>([]);

  useEffect(() => {
    // Очищаем предыдущие URLs
    urlsRef.current.forEach((url) => URL.revokeObjectURL(url));

    // Создаем новые URLs
    urlsRef.current = files.map((file) => URL.createObjectURL(file));

    // Очистка при размонтировании
    return () => {
      urlsRef.current.forEach((url) => URL.revokeObjectURL(url));
      urlsRef.current = [];
    };
  }, [files]);

  return urlsRef.current;
}

/**
 * Утилита для создания временного URL с автоматической очисткой через заданное время
 */
export function createTemporaryObjectURL(
  file: File | Blob,
  duration: number = 60000 // по умолчанию 1 минута
): { url: string; cleanup: () => void } {
  const url = URL.createObjectURL(file);

  const timeoutId = setTimeout(() => {
    URL.revokeObjectURL(url);
  }, duration);

  const cleanup = () => {
    clearTimeout(timeoutId);
    URL.revokeObjectURL(url);
  };

  return { url, cleanup };
}
