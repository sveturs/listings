import { useRef, useCallback } from 'react';

interface CacheEntry<T> {
  data: T;
  timestamp: number;
  key: string;
}

interface MapCacheOptions {
  ttl?: number; // Time to live в миллисекундах
  maxEntries?: number; // Максимальное количество записей в кэше
}

/**
 * Хук для кэширования данных карты
 * Реализует LRU (Least Recently Used) стратегию для управления памятью
 * @param options - Опции кэша (ttl и maxEntries)
 */
export function useMapCache<T = any>(options: MapCacheOptions = {}) {
  const { ttl = 300000, maxEntries = 100 } = options; // По умолчанию 5 минут TTL и 100 записей
  const cache = useRef<Map<string, CacheEntry<T>>>(new Map());

  /**
   * Получить данные из кэша
   * @param key - Ключ для поиска в кэше
   * @returns Данные из кэша или null если не найдено/устарело
   */
  const getCached = useCallback(
    (key: string): T | null => {
      const entry = cache.current.get(key);
      if (!entry) return null;

      // Проверяем не истек ли TTL
      if (Date.now() - entry.timestamp > ttl) {
        cache.current.delete(key);
        return null;
      }

      // Перемещаем запись в конец (LRU)
      cache.current.delete(key);
      cache.current.set(key, entry);

      return entry.data;
    },
    [ttl]
  );

  /**
   * Сохранить данные в кэш
   * @param key - Ключ для сохранения
   * @param data - Данные для кэширования
   */
  const setCached = useCallback(
    (key: string, data: T) => {
      // Удаляем старые записи если превышен лимит (LRU)
      if (cache.current.size >= maxEntries) {
        const firstKey = cache.current.keys().next().value;
        if (firstKey) {
          cache.current.delete(firstKey);
        }
      }

      cache.current.set(key, {
        data,
        timestamp: Date.now(),
        key,
      });
    },
    [maxEntries]
  );

  /**
   * Проверить наличие актуальных данных в кэше
   * @param key - Ключ для проверки
   * @returns true если данные есть и актуальны
   */
  const hasCached = useCallback(
    (key: string): boolean => {
      const entry = cache.current.get(key);
      if (!entry) return false;

      // Проверяем не истек ли TTL
      if (Date.now() - entry.timestamp > ttl) {
        cache.current.delete(key);
        return false;
      }

      return true;
    },
    [ttl]
  );

  /**
   * Удалить конкретную запись из кэша
   * @param key - Ключ для удаления
   */
  const removeCached = useCallback((key: string) => {
    cache.current.delete(key);
  }, []);

  /**
   * Очистить весь кэш
   */
  const clearCache = useCallback(() => {
    cache.current.clear();
  }, []);

  /**
   * Получить размер кэша
   * @returns Количество записей в кэше
   */
  const getCacheSize = useCallback((): number => {
    return cache.current.size;
  }, []);

  /**
   * Получить статистику кэша
   * @returns Объект со статистикой
   */
  const getCacheStats = useCallback(() => {
    const entries = Array.from(cache.current.values());
    const now = Date.now();

    return {
      size: cache.current.size,
      maxEntries,
      ttl,
      entries: entries.map((entry) => ({
        key: entry.key,
        age: now - entry.timestamp,
        expired: now - entry.timestamp > ttl,
      })),
    };
  }, [maxEntries, ttl]);

  /**
   * Создать ключ кэша из параметров
   * Полезно для создания уникальных ключей из объектов параметров
   * @param params - Объект с параметрами
   * @returns Строковый ключ
   */
  const createCacheKey = useCallback((params: Record<string, any>): string => {
    // Сортируем ключи для консистентности
    const sortedKeys = Object.keys(params).sort();
    const keyParts = sortedKeys
      .map((key) => {
        const value = params[key];
        if (value === undefined || value === null) return '';
        if (typeof value === 'object') {
          return `${key}:${JSON.stringify(value)}`;
        }
        return `${key}:${value}`;
      })
      .filter(Boolean);

    return keyParts.join('|');
  }, []);

  /**
   * Обновить TTL для существующей записи
   * @param key - Ключ записи
   */
  const refreshTTL = useCallback((key: string) => {
    const entry = cache.current.get(key);
    if (entry) {
      entry.timestamp = Date.now();
      // Перемещаем в конец (LRU)
      cache.current.delete(key);
      cache.current.set(key, entry);
    }
  }, []);

  return {
    getCached,
    setCached,
    hasCached,
    removeCached,
    clearCache,
    getCacheSize,
    getCacheStats,
    createCacheKey,
    refreshTTL,
  };
}
