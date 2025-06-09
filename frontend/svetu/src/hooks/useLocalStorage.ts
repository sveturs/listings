import { useState, useEffect, useCallback } from 'react';

interface UseLocalStorageOptions {
  userId?: string | number;
  namespace?: string;
}

/**
 * Hook для работы с localStorage с поддержкой многопользовательского режима
 * Автоматически привязывает данные к userId для изоляции между пользователями
 */
export function useLocalStorage<T>(
  key: string,
  initialValue: T,
  options: UseLocalStorageOptions = {}
) {
  const { userId, namespace = 'svetu' } = options;

  // Генерируем уникальный ключ с учетом userId
  const getStorageKey = useCallback(() => {
    const parts = [namespace];
    if (userId) {
      parts.push(`user_${userId}`);
    }
    parts.push(key);
    return parts.join('_');
  }, [namespace, userId, key]);

  // Функция для чтения из localStorage
  const readValue = useCallback((): T => {
    if (typeof window === 'undefined') {
      return initialValue;
    }

    try {
      const storageKey = getStorageKey();
      const item = window.localStorage.getItem(storageKey);

      if (!item) {
        return initialValue;
      }

      const parsed = JSON.parse(item);

      // Проверяем, что данные принадлежат текущему пользователю
      if (parsed.userId && userId && parsed.userId !== userId) {
        return initialValue;
      }

      return parsed.value;
    } catch (error) {
      console.warn(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  }, [initialValue, key, userId, getStorageKey]);

  // State для хранения значения
  const [storedValue, setStoredValue] = useState<T>(readValue);

  // Функция для записи в localStorage
  const setValue = useCallback(
    (value: T | ((val: T) => T)) => {
      if (typeof window === 'undefined') {
        console.warn(
          `Tried setting localStorage key "${key}" even though environment is not a client`
        );
        return;
      }

      try {
        // Вычисляем новое значение
        const newValue = value instanceof Function ? value(storedValue) : value;

        // Сохраняем в state
        setStoredValue(newValue);

        // Сохраняем в localStorage с метаданными
        const storageKey = getStorageKey();
        const dataToStore = {
          value: newValue,
          userId: userId || null,
          timestamp: new Date().toISOString(),
          key: key,
        };

        window.localStorage.setItem(storageKey, JSON.stringify(dataToStore));

        // Диспатчим событие для синхронизации между вкладками
        window.dispatchEvent(
          new StorageEvent('storage', {
            key: storageKey,
            newValue: JSON.stringify(dataToStore),
            storageArea: window.localStorage,
          })
        );
      } catch (error) {
        console.warn(`Error setting localStorage key "${key}":`, error);
      }
    },
    [key, storedValue, userId, getStorageKey]
  );

  // Функция для удаления значения
  const removeValue = useCallback(() => {
    if (typeof window === 'undefined') {
      return;
    }

    try {
      const storageKey = getStorageKey();
      window.localStorage.removeItem(storageKey);
      setStoredValue(initialValue);
    } catch (error) {
      console.warn(`Error removing localStorage key "${key}":`, error);
    }
  }, [key, initialValue, getStorageKey]);

  // Обновляем значение при изменении userId
  useEffect(() => {
    setStoredValue(readValue());
  }, [readValue, userId]);

  // Слушаем изменения в других вкладках
  useEffect(() => {
    const handleStorageChange = (e: StorageEvent) => {
      const storageKey = getStorageKey();
      if (e.key !== storageKey || e.storageArea !== window.localStorage) {
        return;
      }

      try {
        if (e.newValue) {
          const parsed = JSON.parse(e.newValue);
          // Проверяем, что данные принадлежат текущему пользователю
          if (!parsed.userId || !userId || parsed.userId === userId) {
            setStoredValue(parsed.value);
          }
        } else {
          setStoredValue(initialValue);
        }
      } catch {
        console.warn('Error parsing storage event data');
      }
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, [initialValue, key, userId, getStorageKey]);

  return [storedValue, setValue, removeValue] as const;
}

/**
 * Hook для получения всех ключей localStorage для конкретного пользователя
 */
export function useUserStorageKeys(
  userId?: string | number,
  namespace = 'svetu'
) {
  const getUserKeys = useCallback(() => {
    if (typeof window === 'undefined' || !userId) {
      return [];
    }

    const prefix = `${namespace}_user_${userId}_`;
    const keys: string[] = [];

    for (let i = 0; i < window.localStorage.length; i++) {
      const key = window.localStorage.key(i);
      if (key && key.startsWith(prefix)) {
        keys.push(key.replace(prefix, ''));
      }
    }

    return keys;
  }, [userId, namespace]);

  return getUserKeys;
}

/**
 * Hook для очистки всех данных пользователя из localStorage
 */
export function useClearUserStorage(
  userId?: string | number,
  namespace = 'svetu'
) {
  const clearUserData = useCallback(() => {
    if (typeof window === 'undefined' || !userId) {
      return;
    }

    const prefix = `${namespace}_user_${userId}_`;
    const keysToRemove: string[] = [];

    // Собираем все ключи пользователя
    for (let i = 0; i < window.localStorage.length; i++) {
      const key = window.localStorage.key(i);
      if (key && key.startsWith(prefix)) {
        keysToRemove.push(key);
      }
    }

    // Удаляем их
    keysToRemove.forEach((key) => {
      window.localStorage.removeItem(key);
    });

    console.log(`Cleared ${keysToRemove.length} items for user ${userId}`);
  }, [userId, namespace]);

  return clearUserData;
}
