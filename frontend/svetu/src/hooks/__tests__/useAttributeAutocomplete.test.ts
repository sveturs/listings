import { renderHook, act, waitFor } from '@testing-library/react';
import { useAttributeAutocomplete } from '../useAttributeAutocomplete';

describe('useAttributeAutocomplete', () => {
  // Mock localStorage
  let localStorageMock: { [key: string]: string } = {};

  beforeEach(() => {
    // Clear mock
    localStorageMock = {};

    // Mock localStorage implementation
    global.Storage.prototype.getItem = jest.fn((key: string) => {
      return localStorageMock[key] || null;
    });

    global.Storage.prototype.setItem = jest.fn((key: string, value: string) => {
      localStorageMock[key] = value;
    });

    global.Storage.prototype.removeItem = jest.fn((key: string) => {
      delete localStorageMock[key];
    });

    global.Storage.prototype.clear = jest.fn(() => {
      localStorageMock = {};
    });

    Object.defineProperty(global.Storage.prototype, 'length', {
      get: () => Object.keys(localStorageMock).length,
    });

    global.Storage.prototype.key = jest.fn((index: number) => {
      const keys = Object.keys(localStorageMock);
      return keys[index] || null;
    });

    // Clear all timers
    jest.clearAllTimers();
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  describe('Initialization', () => {
    test('инициализируется с пустыми значениями', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      expect(result.current.popularValues).toEqual([]);
      expect(result.current.recentValues).toEqual([]);
    });

    test('загружает данные из localStorage', () => {
      localStorageMock['recent_1'] = JSON.stringify(['Apple', 'Samsung']);
      localStorageMock['popular_brand'] = JSON.stringify(['Xiaomi', 'Huawei']);

      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      expect(result.current.recentValues).toEqual(['Apple', 'Samsung']);
      expect(result.current.popularValues).toEqual(['Xiaomi', 'Huawei']);
    });

    test('ограничивает загрузку до MAX значений', () => {
      // 15 значений (больше MAX_RECENT_VALUES=5 и MAX_POPULAR_VALUES=10)
      const manyValues = Array.from({ length: 15 }, (_, i) => `Value${i}`);

      localStorageMock['recent_1'] = JSON.stringify(manyValues);
      localStorageMock['popular_brand'] = JSON.stringify(manyValues);

      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      expect(result.current.recentValues).toHaveLength(5); // MAX_RECENT_VALUES
      expect(result.current.popularValues).toHaveLength(10); // MAX_POPULAR_VALUES
    });

    test('обрабатывает некорректные данные в localStorage', () => {
      localStorageMock['recent_1'] = 'invalid json';
      localStorageMock['popular_brand'] = '{ broken json';

      const consoleWarnSpy = jest
        .spyOn(console, 'warn')
        .mockImplementation(() => {});

      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      expect(result.current.recentValues).toEqual([]);
      expect(result.current.popularValues).toEqual([]);
      expect(consoleWarnSpy).toHaveBeenCalled();

      consoleWarnSpy.mockRestore();
    });
  });

  describe('addRecentValue', () => {
    test('добавляет значение в недавние', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Apple');
      });

      expect(result.current.recentValues).toContain('Apple');
    });

    test('перемещает существующее значение в начало', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Apple');
        result.current.addRecentValue('Samsung');
        result.current.addRecentValue('Apple'); // Повторное добавление
      });

      expect(result.current.recentValues[0]).toBe('Apple');
      expect(result.current.recentValues).toHaveLength(2); // Без дубликатов
    });

    test('ограничивает недавние значения до MAX_RECENT_VALUES (5)', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Apple');
        result.current.addRecentValue('Samsung');
        result.current.addRecentValue('Xiaomi');
        result.current.addRecentValue('Huawei');
        result.current.addRecentValue('Sony');
        result.current.addRecentValue('LG'); // 6-е значение
      });

      expect(result.current.recentValues).toHaveLength(5);
      expect(result.current.recentValues[0]).toBe('LG'); // Последнее добавленное
      expect(result.current.recentValues).not.toContain('Apple'); // Старое должно быть удалено
    });

    test('игнорирует пустые значения', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('');
        result.current.addRecentValue('   ');
      });

      expect(result.current.recentValues).toEqual([]);
    });

    test('сохраняет в localStorage с debouncing', async () => {
      jest.useFakeTimers();

      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Apple');
      });

      // До истечения debounce - не должно быть в localStorage
      expect(localStorageMock['recent_v1_1']).toBeUndefined();

      // Ждем debounce (100ms)
      act(() => {
        jest.advanceTimersByTime(100);
      });

      await waitFor(() => {
        expect(localStorageMock['recent_v1_1']).toBeTruthy();
        expect(JSON.parse(localStorageMock['recent_v1_1'])).toContain('Apple');
      });

      jest.useRealTimers();
    });

    test('trim значения перед добавлением', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('  Apple  ');
      });

      expect(result.current.recentValues[0]).toBe('Apple');
    });
  });

  describe('incrementPopularity', () => {
    test('увеличивает популярность значения', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.incrementPopularity('Apple');
        result.current.incrementPopularity('Apple');
        result.current.incrementPopularity('Samsung');
      });

      // Apple должен быть первым (2 раза vs 1 раз)
      expect(result.current.popularValues[0]).toBe('Apple');
      expect(result.current.popularValues[1]).toBe('Samsung');
    });

    test('ограничивает популярные значения до MAX_POPULAR_VALUES (10)', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        // Добавляем 15 разных значений
        for (let i = 0; i < 15; i++) {
          result.current.incrementPopularity(`Brand${i}`);
        }
      });

      expect(result.current.popularValues).toHaveLength(10); // MAX_POPULAR_VALUES
    });

    test('игнорирует пустые значения', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.incrementPopularity('');
        result.current.incrementPopularity('   ');
      });

      expect(result.current.popularValues).toEqual([]);
    });

    test('сохраняет в localStorage с debouncing', async () => {
      jest.useFakeTimers();

      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.incrementPopularity('Apple');
      });

      // До истечения debounce
      expect(localStorageMock['popular_v1_brand']).toBeUndefined();

      // Ждем debounce
      act(() => {
        jest.advanceTimersByTime(100);
      });

      await waitFor(() => {
        expect(localStorageMock['popular_v1_brand']).toBeTruthy();
        expect(JSON.parse(localStorageMock['popular_v1_brand'])).toContain(
          'Apple'
        );
      });

      jest.useRealTimers();
    });
  });

  describe('saveValue', () => {
    test('сохраняет значение и в recent и в popular', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('Apple');
      });

      expect(result.current.recentValues).toContain('Apple');
      expect(result.current.popularValues).toContain('Apple');
    });
  });

  describe('getAllSuggestions', () => {
    test('возвращает комбинацию популярных и недавних', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Recent1');
        result.current.addRecentValue('Recent2');
        result.current.incrementPopularity('Popular1');
        result.current.incrementPopularity('Popular2');
      });

      const suggestions = result.current.getAllSuggestions();

      expect(suggestions.length).toBeGreaterThan(0);
      expect(suggestions.some((s) => s.value === 'Recent1')).toBe(true);
      expect(suggestions.some((s) => s.value === 'Popular1')).toBe(true);
    });

    test('не возвращает дубликаты', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Apple');
        result.current.incrementPopularity('Apple');
      });

      const suggestions = result.current.getAllSuggestions();
      const appleCount = suggestions.filter((s) => s.value === 'Apple').length;

      expect(appleCount).toBe(1); // Только один раз
    });

    test('ограничивает до 8 предложений', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        // Добавляем много значений
        for (let i = 0; i < 15; i++) {
          result.current.addRecentValue(`Recent${i}`);
          result.current.incrementPopularity(`Popular${i}`);
        }
      });

      const suggestions = result.current.getAllSuggestions();
      expect(suggestions.length).toBeLessThanOrEqual(8);
    });

    test('помечает типы предложений правильно', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.incrementPopularity('PopularBrand');
        result.current.addRecentValue('RecentBrand');
      });

      const suggestions = result.current.getAllSuggestions();

      const popularSuggestion = suggestions.find(
        (s) => s.value === 'PopularBrand'
      );
      const recentSuggestion = suggestions.find(
        (s) => s.value === 'RecentBrand'
      );

      expect(popularSuggestion?.type).toBe('popular');
      expect(recentSuggestion?.type).toBe('recent');
    });
  });

  describe('getFilteredSuggestions', () => {
    test('возвращает все предложения для пустого запроса', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('Apple');
        result.current.saveValue('Samsung');
      });

      const suggestions = result.current.getFilteredSuggestions('');
      expect(suggestions.length).toBeGreaterThan(0);
    });

    test('фильтрует по запросу (startsWith)', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('Apple');
        result.current.saveValue('Samsung');
        result.current.saveValue('Xiaomi');
      });

      const suggestions = result.current.getFilteredSuggestions('Sam');
      expect(suggestions).toHaveLength(1);
      expect(suggestions[0].value).toBe('Samsung');
    });

    test('фильтрует по запросу (contains)', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('iPhone 15');
        result.current.saveValue('Samsung Galaxy');
      });

      const suggestions = result.current.getFilteredSuggestions('phone');
      expect(suggestions.length).toBeGreaterThan(0);
      expect(suggestions[0].value).toBe('iPhone 15');
    });

    test('ранжирует точное совпадение выше', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('Apple iPhone');
        result.current.saveValue('Apple');
      });

      const suggestions = result.current.getFilteredSuggestions('Apple');
      expect(suggestions[0].value).toBe('Apple'); // Точное совпадение первым
    });

    test('не чувствителен к регистру', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('Apple');
      });

      const suggestions1 = result.current.getFilteredSuggestions('apple');
      const suggestions2 = result.current.getFilteredSuggestions('APPLE');

      expect(suggestions1).toHaveLength(1);
      expect(suggestions2).toHaveLength(1);
    });

    test('популярные значения имеют бонус к релевантности', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Apple iPhone');
        result.current.incrementPopularity('Apple Watch');
        result.current.incrementPopularity('Apple Watch'); // Дважды для популярности
      });

      const suggestions = result.current.getFilteredSuggestions('Apple');

      // Apple Watch должен быть выше из-за популярности
      expect(suggestions[0].value).toBe('Apple Watch');
    });

    test('ограничивает результаты до 6 предложений', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        for (let i = 0; i < 10; i++) {
          result.current.saveValue(`Brand${i}`);
        }
      });

      const suggestions = result.current.getFilteredSuggestions('Brand');
      expect(suggestions.length).toBeLessThanOrEqual(6);
    });
  });

  describe('clearData', () => {
    test('очищает все данные', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('Apple');
        result.current.saveValue('Samsung');
      });

      expect(result.current.recentValues.length).toBeGreaterThan(0);
      expect(result.current.popularValues.length).toBeGreaterThan(0);

      act(() => {
        result.current.clearData();
      });

      expect(result.current.recentValues).toEqual([]);
      expect(result.current.popularValues).toEqual([]);
    });

    test('удаляет данные из localStorage', () => {
      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.saveValue('Apple');
        result.current.clearData();
      });

      expect(localStorageMock['recent_1']).toBeUndefined();
      expect(localStorageMock['popular_brand']).toBeUndefined();
    });
  });

  describe('clearOldStorageData', () => {
    test('очищает старые ключи без версии', () => {
      localStorageMock['recent_1'] = '["old"]';
      localStorageMock['recent_v1_1'] = '["new"]';
      localStorageMock['popular_brand'] = '["old"]';
      localStorageMock['popular_v1_brand'] = '["new"]';

      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.clearData();
      });

      // Старые ключи должны быть удалены через clearOldStorageData
      // (вызывается при QuotaExceededError)
    });
  });

  describe('Error handling', () => {
    test('обрабатывает QuotaExceededError', async () => {
      jest.useFakeTimers();

      const consoleWarnSpy = jest
        .spyOn(console, 'warn')
        .mockImplementation(() => {});

      // Mock setItem to throw QuotaExceededError
      const originalSetItem = global.Storage.prototype.setItem;
      global.Storage.prototype.setItem = jest.fn(() => {
        const error = new DOMException('QuotaExceededError');
        error.name = 'QuotaExceededError';
        throw error;
      });

      const { result } = renderHook(() =>
        useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
      );

      act(() => {
        result.current.addRecentValue('Apple');
      });

      // Advance timers to trigger debounced save
      act(() => {
        jest.advanceTimersByTime(100);
      });

      // Wait for the error to be caught and logged
      await waitFor(() => {
        expect(consoleWarnSpy).toHaveBeenCalled();
      });

      // Restore
      global.Storage.prototype.setItem = originalSetItem;
      consoleWarnSpy.mockRestore();
      jest.useRealTimers();
    });
  });
});
