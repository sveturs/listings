import { CarsService } from '../cars';
import type { CarMake, CarModel, CarGeneration } from '@/types/cars';

// Mock configManager
jest.mock('@/config', () => ({
  __esModule: true,
  default: {
    getApiUrl: () => 'http://localhost:3000',
  },
}));

describe('CarsService', () => {
  beforeEach(() => {
    // Mock global fetch
    global.fetch = jest.fn();
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  describe('getMakes', () => {
    test('возвращает список марок при успешном запросе', async () => {
      const mockMakes: CarMake[] = [
        { id: 1, name: 'BMW', slug: 'bmw' },
        { id: 2, name: 'Mercedes', slug: 'mercedes' },
      ];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: mockMakes }),
      } as Response);

      const result = await CarsService.getMakes();

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockMakes);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes',
        expect.objectContaining({
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
        })
      );
    });

    test('обрабатывает данные без обертки .data', async () => {
      const mockMakes: CarMake[] = [{ id: 1, name: 'BMW', slug: 'bmw' }];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockMakes, // Без обертки
      } as Response);

      const result = await CarsService.getMakes();

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockMakes);
    });

    test('обрабатывает HTTP ошибку', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 404,
      } as Response);

      const result = await CarsService.getMakes();

      expect(result.success).toBe(false);
      expect(result.error).toContain('404');
    });

    test('обрабатывает network ошибку', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      const result = await CarsService.getMakes();

      expect(result.success).toBe(false);
      expect(result.error).toBe('Network error');
    });

    test('обрабатывает неизвестную ошибку', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce('Unknown error object');

      const result = await CarsService.getMakes();

      expect(result.success).toBe(false);
      expect(result.error).toBe('Unknown error');
    });

    test('логирует ошибку в консоль', async () => {
      const consoleErrorSpy = jest
        .spyOn(console, 'error')
        .mockImplementation(() => {});

      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Test error')
      );

      await CarsService.getMakes();

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Error fetching car makes:',
        expect.any(Error)
      );

      consoleErrorSpy.mockRestore();
    });
  });

  describe('getModelsByMake', () => {
    test('возвращает модели для указанной марки', async () => {
      const mockModels: CarModel[] = [
        { id: 1, name: 'X5', make_id: 1 },
        { id: 2, name: 'X7', make_id: 1 },
      ];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: mockModels }),
      } as Response);

      const result = await CarsService.getModelsByMake('bmw');

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockModels);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes/bmw/models',
        expect.any(Object)
      );
    });

    test('правильно обрабатывает slug с дефисами', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      await CarsService.getModelsByMake('aston-martin');

      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes/aston-martin/models',
        expect.any(Object)
      );
    });

    test('обрабатывает данные без обертки .data', async () => {
      const mockModels: CarModel[] = [{ id: 1, name: 'X5', make_id: 1 }];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockModels,
      } as Response);

      const result = await CarsService.getModelsByMake('bmw');

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockModels);
    });

    test('обрабатывает HTTP ошибку', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 500,
      } as Response);

      const result = await CarsService.getModelsByMake('bmw');

      expect(result.success).toBe(false);
      expect(result.error).toContain('500');
    });

    test('обрабатывает network ошибку', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Connection failed')
      );

      const result = await CarsService.getModelsByMake('bmw');

      expect(result.success).toBe(false);
      expect(result.error).toBe('Connection failed');
    });

    test('логирует ошибку в консоль', async () => {
      const consoleErrorSpy = jest
        .spyOn(console, 'error')
        .mockImplementation(() => {});

      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Test error')
      );

      await CarsService.getModelsByMake('bmw');

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Error fetching car models:',
        expect.any(Error)
      );

      consoleErrorSpy.mockRestore();
    });
  });

  describe('getGenerationsByModel', () => {
    test('возвращает поколения для модели', async () => {
      const mockGenerations: CarGeneration[] = [
        { id: 1, name: 'F15 (2013-2018)', model_id: 10 },
        { id: 2, name: 'G05 (2018-present)', model_id: 10 },
      ];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: mockGenerations }),
      } as Response);

      const result = await CarsService.getGenerationsByModel(10);

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockGenerations);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/models/10/generations',
        expect.any(Object)
      );
    });

    test('правильно передает modelId в URL', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      await CarsService.getGenerationsByModel(42);

      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/models/42/generations',
        expect.any(Object)
      );
    });

    test('обрабатывает данные без обертки .data', async () => {
      const mockGenerations: CarGeneration[] = [
        { id: 1, name: 'F15', model_id: 10 },
      ];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockGenerations,
      } as Response);

      const result = await CarsService.getGenerationsByModel(10);

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockGenerations);
    });

    test('обрабатывает HTTP ошибку', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 404,
      } as Response);

      const result = await CarsService.getGenerationsByModel(999);

      expect(result.success).toBe(false);
      expect(result.error).toContain('404');
    });

    test('обрабатывает network ошибку', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce(new Error('Timeout'));

      const result = await CarsService.getGenerationsByModel(10);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Timeout');
    });

    test('логирует ошибку в консоль', async () => {
      const consoleErrorSpy = jest
        .spyOn(console, 'error')
        .mockImplementation(() => {});

      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Test error')
      );

      await CarsService.getGenerationsByModel(10);

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Error fetching car generations:',
        expect.any(Error)
      );

      consoleErrorSpy.mockRestore();
    });
  });

  describe('searchMakes', () => {
    test('ищет марки по запросу', async () => {
      const mockResults: CarMake[] = [{ id: 1, name: 'BMW', slug: 'bmw' }];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: mockResults }),
      } as Response);

      const result = await CarsService.searchMakes('BM');

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockResults);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes/search?q=BM',
        expect.any(Object)
      );
    });

    test('правильно кодирует спецсимволы в запросе', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      await CarsService.searchMakes('BMW & Mercedes');

      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes/search?q=BMW%20%26%20Mercedes',
        expect.any(Object)
      );
    });

    test('правильно кодирует пробелы', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      await CarsService.searchMakes('Land Rover');

      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes/search?q=Land%20Rover',
        expect.any(Object)
      );
    });

    test('правильно кодирует кириллицу', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      await CarsService.searchMakes('БМВ');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('makes/search?q='),
        expect.any(Object)
      );

      // Проверяем что query parameter закодирован
      const url = (global.fetch as jest.Mock).mock.calls[0][0] as string;
      expect(url).toContain('%');
    });

    test('обрабатывает пустой результат', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      const result = await CarsService.searchMakes('NonExistentBrand');

      expect(result.success).toBe(true);
      expect(result.data).toEqual([]);
    });

    test('обрабатывает данные без обертки .data', async () => {
      const mockResults: CarMake[] = [{ id: 1, name: 'BMW', slug: 'bmw' }];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResults,
      } as Response);

      const result = await CarsService.searchMakes('BM');

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockResults);
    });

    test('обрабатывает HTTP ошибку', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 400,
      } as Response);

      const result = await CarsService.searchMakes('test');

      expect(result.success).toBe(false);
      expect(result.error).toContain('400');
    });

    test('обрабатывает network ошибку', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('DNS lookup failed')
      );

      const result = await CarsService.searchMakes('test');

      expect(result.success).toBe(false);
      expect(result.error).toBe('DNS lookup failed');
    });

    test('логирует ошибку в консоль', async () => {
      const consoleErrorSpy = jest
        .spyOn(console, 'error')
        .mockImplementation(() => {});

      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Test error')
      );

      await CarsService.searchMakes('test');

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Error searching car makes:',
        expect.any(Error)
      );

      consoleErrorSpy.mockRestore();
    });
  });

  describe('baseUrl', () => {
    test('использует правильный базовый URL', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      await CarsService.getMakes();

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('http://localhost:3000/api/v1/cars'),
        expect.any(Object)
      );
    });
  });

  describe('Headers', () => {
    test('все методы используют правильные headers', async () => {
      const expectedHeaders = {
        'Content-Type': 'application/json',
      };

      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      await CarsService.getMakes();
      await CarsService.getModelsByMake('bmw');
      await CarsService.getGenerationsByModel(10);
      await CarsService.searchMakes('test');

      // Проверяем что все 4 вызова использовали правильные headers
      expect(global.fetch).toHaveBeenCalledTimes(4);
      (global.fetch as jest.Mock).mock.calls.forEach((call) => {
        expect(call[1]).toHaveProperty('headers', expectedHeaders);
      });
    });
  });

  describe('Edge cases', () => {
    test('обрабатывает пустой slug', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      const result = await CarsService.getModelsByMake('');

      expect(result.success).toBe(true);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes//models',
        expect.any(Object)
      );
    });

    test('обрабатывает modelId = 0', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      const result = await CarsService.getGenerationsByModel(0);

      expect(result.success).toBe(true);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/models/0/generations',
        expect.any(Object)
      );
    });

    test('обрабатывает отрицательный modelId', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      const result = await CarsService.getGenerationsByModel(-1);

      expect(result.success).toBe(true);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/models/-1/generations',
        expect.any(Object)
      );
    });

    test('обрабатывает пустой query в searchMakes', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [] }),
      } as Response);

      const result = await CarsService.searchMakes('');

      expect(result.success).toBe(true);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/cars/makes/search?q=',
        expect.any(Object)
      );
    });
  });
});
