'use client';

import { useState, useEffect, useCallback, useRef } from 'react';

interface AutocompleteData {
  popularValues: string[];
  recentValues: string[];
}

interface UseAttributeAutocompleteProps {
  attributeId: number;
  attributeName: string;
}

// Константы для оптимизации
const MAX_RECENT_VALUES = 5;
const MAX_POPULAR_VALUES = 10;
const DEBOUNCE_DELAY = 100; // мс для debounce сохранения в localStorage
const STORAGE_VERSION = 'v1'; // для миграций в будущем

export function useAttributeAutocomplete({
  attributeId,
  attributeName,
}: UseAttributeAutocompleteProps) {
  const [data, setData] = useState<AutocompleteData>({
    popularValues: [],
    recentValues: [],
  });

  // Debounce refs для оптимизации записи в localStorage
  const saveRecentTimeout = useRef<NodeJS.Timeout>();
  const savePopularTimeout = useRef<NodeJS.Timeout>();

  // Загрузка данных из localStorage при инициализации
  useEffect(() => {
    const loadData = () => {
      try {
        // Загрузка недавних значений
        const recentKey = `recent_${attributeId}`;
        const recentValues = JSON.parse(
          localStorage.getItem(recentKey) || '[]'
        );

        // Загрузка популярных значений (общих для атрибута по имени)
        const popularKey = `popular_${attributeName.toLowerCase()}`;
        const popularValues = JSON.parse(
          localStorage.getItem(popularKey) || '[]'
        );

        setData({
          popularValues: popularValues.slice(0, 10), // Максимум 10 популярных
          recentValues: recentValues.slice(0, 5), // Максимум 5 недавних
        });
      } catch (error) {
        console.warn('Error loading autocomplete data:', error);
        setData({ popularValues: [], recentValues: [] });
      }
    };

    loadData();
  }, [attributeId, attributeName]);

  // Добавление значения в недавние с debouncing
  const addRecentValue = useCallback(
    (value: string) => {
      if (!value.trim()) return;

      // Обновляем state сразу для UX
      const trimmedValue = value.trim();
      setData((prev) => {
        const newRecent = [
          trimmedValue,
          ...prev.recentValues.filter((v) => v !== trimmedValue),
        ].slice(0, MAX_RECENT_VALUES);

        // Debounced сохранение в localStorage
        if (saveRecentTimeout.current) {
          clearTimeout(saveRecentTimeout.current);
        }

        saveRecentTimeout.current = setTimeout(() => {
          try {
            const recentKey = `recent_${STORAGE_VERSION}_${attributeId}`;
            localStorage.setItem(recentKey, JSON.stringify(newRecent));
          } catch (error) {
            console.warn('Error saving recent value to localStorage:', error);
            // Очищаем localStorage если превышен лимит
            if (
              error instanceof DOMException &&
              error.name === 'QuotaExceededError'
            ) {
              clearOldStorageData();
            }
          }
        }, DEBOUNCE_DELAY);

        return {
          ...prev,
          recentValues: newRecent,
        };
      });
    },
    [attributeId]
  );

  // Функция очистки старых данных из localStorage
  const clearOldStorageData = useCallback(() => {
    try {
      const keysToRemove: string[] = [];

      // Проходим по всем ключам localStorage
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (
          key &&
          (key.startsWith('recent_') ||
            key.startsWith('popular_') ||
            key.startsWith('count_'))
        ) {
          // Удаляем старые ключи без версии или с неправильной версией
          if (!key.includes(STORAGE_VERSION)) {
            keysToRemove.push(key);
          }
        }
      }

      // Удаляем старые ключи
      keysToRemove.forEach((key) => localStorage.removeItem(key));

      console.log(`Cleaned ${keysToRemove.length} old localStorage keys`);
    } catch (error) {
      console.warn('Error cleaning old storage data:', error);
    }
  }, []);

  // Увеличение популярности значения с debouncing
  const incrementPopularity = useCallback(
    (value: string) => {
      if (!value.trim()) return;

      const trimmedValue = value.trim();

      // Обновляем state сразу для UX (оптимистичное обновление)
      setData((prev) => {
        const existingCounts: Record<string, number> = {};

        // Восстанавливаем счетчики из популярных значений
        prev.popularValues.forEach((val, index) => {
          existingCounts[val] = prev.popularValues.length - index;
        });

        // Увеличиваем счетчик
        existingCounts[trimmedValue] = (existingCounts[trimmedValue] || 0) + 1;

        // Пересортировываем по популярности
        const sortedValues = Object.entries(existingCounts)
          .sort(([, a], [, b]) => (b as number) - (a as number))
          .map(([val]) => val)
          .slice(0, MAX_POPULAR_VALUES);

        // Debounced сохранение в localStorage
        if (savePopularTimeout.current) {
          clearTimeout(savePopularTimeout.current);
        }

        savePopularTimeout.current = setTimeout(() => {
          try {
            const popularKey = `popular_${STORAGE_VERSION}_${attributeName.toLowerCase()}`;
            const countKey = `count_${STORAGE_VERSION}_${attributeName.toLowerCase()}`;

            localStorage.setItem(popularKey, JSON.stringify(sortedValues));
            localStorage.setItem(countKey, JSON.stringify(existingCounts));
          } catch (error) {
            console.warn('Error updating popularity in localStorage:', error);
            if (
              error instanceof DOMException &&
              error.name === 'QuotaExceededError'
            ) {
              clearOldStorageData();
            }
          }
        }, DEBOUNCE_DELAY);

        return {
          ...prev,
          popularValues: sortedValues,
        };
      });
    },
    [attributeName, clearOldStorageData]
  );

  // Комбинированная функция для сохранения значения
  const saveValue = useCallback(
    (value: string) => {
      addRecentValue(value);
      incrementPopularity(value);
    },
    [addRecentValue, incrementPopularity]
  );

  // Получение всех предложений для пустого поля
  const getAllSuggestions = useCallback(() => {
    const suggestions: { value: string; type: 'popular' | 'recent' }[] = [];

    // Добавляем популярные
    data.popularValues.forEach((value) => {
      suggestions.push({ value, type: 'popular' });
    });

    // Добавляем недавние (исключая дубликаты)
    data.recentValues.forEach((value) => {
      if (!data.popularValues.includes(value)) {
        suggestions.push({ value, type: 'recent' });
      }
    });

    return suggestions.slice(0, 8); // Максимум 8 предложений
  }, [data]);

  // Фильтрация предложений по запросу
  const getFilteredSuggestions = useCallback(
    (query: string) => {
      if (!query.trim()) return getAllSuggestions();

      const lowerQuery = query.toLowerCase();
      const allValues = [...data.popularValues, ...data.recentValues];
      const uniqueValues = [...new Set(allValues)]; // Убираем дубликаты

      const suggestions: {
        value: string;
        type: 'popular' | 'recent';
        relevance: number;
      }[] = [];

      uniqueValues.forEach((value) => {
        const lowerValue = value.toLowerCase();
        let relevance = 0;

        if (lowerValue === lowerQuery) {
          relevance = 1000; // Точное совпадение
        } else if (lowerValue.startsWith(lowerQuery)) {
          relevance = 800; // Начинается с запроса
        } else if (lowerValue.includes(lowerQuery)) {
          relevance = 600; // Содержит запрос
        }

        if (relevance > 0) {
          // Добавляем бонус за популярность
          if (data.popularValues.includes(value)) {
            relevance += 100;
          }

          suggestions.push({
            value,
            type: data.popularValues.includes(value) ? 'popular' : 'recent',
            relevance,
          });
        }
      });

      return suggestions
        .sort((a, b) => b.relevance - a.relevance)
        .slice(0, 6)
        .map(({ value, type }) => ({ value, type }));
    },
    [data, getAllSuggestions]
  );

  // Очистка данных (для отладки или сброса)
  const clearData = useCallback(() => {
    const recentKey = `recent_${attributeId}`;
    const popularKey = `popular_${attributeName.toLowerCase()}`;
    const countKey = `count_${attributeName.toLowerCase()}`;

    localStorage.removeItem(recentKey);
    localStorage.removeItem(popularKey);
    localStorage.removeItem(countKey);

    setData({ popularValues: [], recentValues: [] });
  }, [attributeId, attributeName]);

  return {
    popularValues: data.popularValues,
    recentValues: data.recentValues,
    saveValue,
    addRecentValue,
    incrementPopularity,
    getAllSuggestions,
    getFilteredSuggestions,
    clearData,
  };
}
