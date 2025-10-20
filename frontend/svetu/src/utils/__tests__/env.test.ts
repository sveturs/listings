// Mock next-runtime-env BEFORE imports
const mockEnvFunction = jest.fn((key: string) => {
  const mockEnv: Record<string, string> = {
    NEXT_PUBLIC_API_URL: 'http://test-api.com',
    NEXT_PUBLIC_MINIO_URL: 'http://test-minio.com',
    NEXT_PUBLIC_WEBSOCKET_URL: 'ws://test-websocket.com',
    NEXT_PUBLIC_IMAGE_HOSTS: 's3.test.com,cdn.test.com',
    NEXT_PUBLIC_ENABLE_PAYMENTS: 'true',
  };
  return mockEnv[key];
});

jest.mock('next-runtime-env', () => ({
  env: (key: string) => mockEnvFunction(key),
}));

import { getEnv, publicEnv } from '../env';

describe('env utilities', () => {
  const originalWindow = global.window;
  const originalProcessEnv = process.env;

  beforeEach(() => {
    mockEnvFunction.mockClear();
  });

  afterEach(() => {
    global.window = originalWindow;
    process.env = originalProcessEnv;
  });

  describe('getEnv', () => {
    describe('Client-side (browser)', () => {
      beforeEach(() => {
        // Mock browser environment
        global.window = {} as any;
      });

      test('возвращает значение из runtime env', () => {
        expect(getEnv('NEXT_PUBLIC_API_URL')).toBe('http://test-api.com');
      });

      test('возвращает значение для MINIO_URL', () => {
        expect(getEnv('NEXT_PUBLIC_MINIO_URL')).toBe('http://test-minio.com');
      });

      test('возвращает defaultValue если переменная не найдена', () => {
        expect(getEnv('NON_EXISTENT_VAR', 'default')).toBe('default');
      });

      test('возвращает undefined если нет defaultValue и переменной', () => {
        mockEnvFunction.mockReturnValueOnce(undefined);
        expect(getEnv('NON_EXISTENT_VAR')).toBeUndefined();
      });

      test('вызывает runtime env функцию', () => {
        getEnv('NEXT_PUBLIC_API_URL');
        expect(mockEnvFunction).toHaveBeenCalledWith('NEXT_PUBLIC_API_URL');
      });
    });

    // NOTE: Server-side тесты требуют сложной настройки с jest.resetModules(),
    // что конфликтует с нашим jest.mock('next-runtime-env'). Основная логика
    // server-side покрывается через интеграционные тесты Next.js SSR.

    describe('Edge cases (client-side)', () => {
      beforeEach(() => {
        // Mock browser environment
        global.window = {} as any;
      });

      test('обрабатывает пустую строку как значение', () => {
        mockEnvFunction.mockReturnValueOnce('');
        // Пустая строка в логике || считается falsy, поэтому вернется defaultValue
        expect(getEnv('EMPTY_VAR', 'default')).toBe('default');
      });

      test('обрабатывает undefined как отсутствие значения', () => {
        mockEnvFunction.mockReturnValueOnce(undefined);
        expect(getEnv('MISSING_VAR', 'fallback')).toBe('fallback');
      });
    });
  });

  describe('publicEnv', () => {
    beforeEach(() => {
      // Mock browser environment для всех publicEnv тестов
      global.window = {} as any;
    });

    describe('API_URL', () => {
      test('возвращает правильный API_URL', () => {
        expect(publicEnv.API_URL).toBe('http://test-api.com');
      });

      test('использует дефолтный API_URL если не задан', () => {
        // Mock empty return
        mockEnvFunction.mockReturnValueOnce(undefined);
        expect(publicEnv.API_URL).toBe('http://localhost:3000');
      });
    });

    describe('MINIO_URL', () => {
      test('возвращает правильный MINIO_URL', () => {
        expect(publicEnv.MINIO_URL).toBe('http://test-minio.com');
      });

      test('использует дефолтный MINIO_URL если не задан', () => {
        mockEnvFunction.mockReturnValueOnce(undefined);
        expect(publicEnv.MINIO_URL).toBe('http://localhost:9000');
      });
    });

    describe('WEBSOCKET_URL', () => {
      test('возвращает правильный WEBSOCKET_URL', () => {
        expect(publicEnv.WEBSOCKET_URL).toBe('ws://test-websocket.com');
      });

      test('возвращает undefined если WEBSOCKET_URL не задан', () => {
        mockEnvFunction.mockReturnValueOnce(undefined);
        expect(publicEnv.WEBSOCKET_URL).toBeUndefined();
      });
    });

    describe('IMAGE_HOSTS', () => {
      test('возвращает правильный IMAGE_HOSTS', () => {
        expect(publicEnv.IMAGE_HOSTS).toBe('s3.test.com,cdn.test.com');
      });

      test('возвращает undefined если IMAGE_HOSTS не задан', () => {
        mockEnvFunction.mockReturnValueOnce(undefined);
        expect(publicEnv.IMAGE_HOSTS).toBeUndefined();
      });
    });

    describe('ENABLE_PAYMENTS', () => {
      test('парсит "true" как boolean true', () => {
        expect(publicEnv.ENABLE_PAYMENTS).toBe(true);
      });

      test('парсит "false" как boolean false', () => {
        mockEnvFunction.mockReturnValueOnce('false');
        expect(publicEnv.ENABLE_PAYMENTS).toBe(false);
      });

      test('парсит пустую строку как boolean false', () => {
        mockEnvFunction.mockReturnValueOnce('');
        expect(publicEnv.ENABLE_PAYMENTS).toBe(false);
      });

      test('парсит undefined как boolean false', () => {
        mockEnvFunction.mockReturnValueOnce(undefined);
        expect(publicEnv.ENABLE_PAYMENTS).toBe(false);
      });

      test('парсит любое значение кроме "true" как false', () => {
        mockEnvFunction.mockReturnValueOnce('yes');
        expect(publicEnv.ENABLE_PAYMENTS).toBe(false);

        mockEnvFunction.mockReturnValueOnce('1');
        expect(publicEnv.ENABLE_PAYMENTS).toBe(false);

        mockEnvFunction.mockReturnValueOnce('TRUE');
        expect(publicEnv.ENABLE_PAYMENTS).toBe(false);
      });
    });

    describe('Геттеры вызываются каждый раз', () => {
      test('API_URL геттер вызывается при каждом обращении', () => {
        mockEnvFunction.mockClear();

        const url1 = publicEnv.API_URL;
        const url2 = publicEnv.API_URL;

        // Должно быть 2 вызова getEnv
        expect(mockEnvFunction).toHaveBeenCalledTimes(2);
      });

      test('ENABLE_PAYMENTS геттер вызывается при каждом обращении', () => {
        mockEnvFunction.mockClear();

        const enabled1 = publicEnv.ENABLE_PAYMENTS;
        const enabled2 = publicEnv.ENABLE_PAYMENTS;

        expect(mockEnvFunction).toHaveBeenCalledTimes(2);
      });
    });

    describe('Все поля доступны', () => {
      test('все поля publicEnv определены', () => {
        expect(publicEnv).toHaveProperty('API_URL');
        expect(publicEnv).toHaveProperty('MINIO_URL');
        expect(publicEnv).toHaveProperty('WEBSOCKET_URL');
        expect(publicEnv).toHaveProperty('IMAGE_HOSTS');
        expect(publicEnv).toHaveProperty('ENABLE_PAYMENTS');
      });

      test('все геттеры возвращают значения', () => {
        expect(publicEnv.API_URL).toBeDefined();
        expect(publicEnv.MINIO_URL).toBeDefined();
        expect(typeof publicEnv.ENABLE_PAYMENTS).toBe('boolean');
      });
    });
  });

  describe('Интеграция getEnv и publicEnv', () => {
    beforeEach(() => {
      global.window = {} as any;
    });

    test('publicEnv использует getEnv под капотом', () => {
      mockEnvFunction.mockClear();

      // Обращение к publicEnv должно вызвать getEnv, который вызовет mockEnvFunction
      const apiUrl = publicEnv.API_URL;

      expect(mockEnvFunction).toHaveBeenCalledWith('NEXT_PUBLIC_API_URL');
      expect(apiUrl).toBe('http://test-api.com');
    });

    test('defaultValue работает через publicEnv', () => {
      mockEnvFunction.mockReturnValueOnce(undefined);

      const apiUrl = publicEnv.API_URL;

      expect(apiUrl).toBe('http://localhost:3000'); // default value
    });
  });

  describe('Типы возвращаемых значений', () => {
    beforeEach(() => {
      global.window = {} as any;
    });

    test('getEnv возвращает string | undefined', () => {
      const result = getEnv('NEXT_PUBLIC_API_URL');
      expect(typeof result === 'string' || result === undefined).toBe(true);
    });

    test('publicEnv.API_URL всегда string (с default)', () => {
      expect(typeof publicEnv.API_URL).toBe('string');
    });

    test('publicEnv.WEBSOCKET_URL может быть undefined', () => {
      mockEnvFunction.mockReturnValueOnce(undefined);
      const result = publicEnv.WEBSOCKET_URL;
      expect(result === undefined || typeof result === 'string').toBe(true);
    });

    test('publicEnv.ENABLE_PAYMENTS всегда boolean', () => {
      expect(typeof publicEnv.ENABLE_PAYMENTS).toBe('boolean');
    });
  });
});
