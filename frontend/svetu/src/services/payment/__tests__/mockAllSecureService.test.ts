import { MockAllSecureService } from '../mockAllSecureService';
import type { MockConfig } from '@/config/payment';
import type { PaymentRequest } from '@/types/payment';

// Mock утилиты
jest.mock('@/utils/mockDataGenerator', () => ({
  generateMockPayment: jest.fn(() => ({
    id: 'mock-payment-123',
    status: 'pending',
    amount: 1000,
    currency: 'RSD',
    listing_id: '123',
    buyer_info: {},
    card: {
      last4: '1111',
      brand: 'visa',
    },
    createdAt: '2025-06-27T00:00:00.000Z',
    completedAt: null,
  })),
}));

describe('MockAllSecureService', () => {
  let service: MockAllSecureService;
  let mockConfig: MockConfig;
  let mockPaymentRequest: PaymentRequest;

  // Mock localStorage
  const mockLocalStorage = {
    getItem: jest.fn(),
    setItem: jest.fn(),
    removeItem: jest.fn(),
    clear: jest.fn(),
  };

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    jest.clearAllTimers();
    jest.useFakeTimers();

    // Mock localStorage
    Object.defineProperty(window, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    });

    // Mock Math.random для предсказуемых результатов
    jest.spyOn(Math, 'random').mockReturnValue(0.5);

    mockConfig = {
      successRate: 0.8,
      require3DSRate: 0.3,
      apiDelay: 100,
      webhookDelay: 2000,
      debugMode: false,
    };

    mockPaymentRequest = {
      amount: 1500,
      currency: 'RSD',
      listing_id: '456',
      buyer_info: {
        name: 'Test User',
        email: 'test@example.com',
      },
      locale: 'ru',
    };

    service = new MockAllSecureService(mockConfig);
  });

  afterEach(() => {
    jest.useRealTimers();
    jest.restoreAllMocks();
    service.cleanup();
  });

  describe('createPayment', () => {
    it('должен создать платеж с корректными данными', async () => {
      const response = await service.createPayment(mockPaymentRequest);

      expect(response).toEqual({
        id: 'mock-payment-123',
        redirectUrl: expect.stringContaining(
          '/payment/mock?id=mock-payment-123'
        ),
        status: 'pending',
      });
    });

    it('должен сохранить платеж в localStorage', async () => {
      await service.createPayment(mockPaymentRequest);

      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
        'mock_payment_mock-payment-123',
        expect.stringContaining('"id":"mock-payment-123"')
      );
    });

    it('должен добавить require3ds=true в URL если требуется 3DS', async () => {
      // Устанавливаем Math.random для требования 3DS (< 0.3)
      jest.spyOn(Math, 'random').mockReturnValue(0.2);

      const response = await service.createPayment(mockPaymentRequest);

      expect(response.redirectUrl).toContain('require3ds=true');
    });

    it('не должен добавлять require3ds если 3DS не требуется', async () => {
      // Устанавливаем Math.random для НЕ требования 3DS (>= 0.3)
      jest.spyOn(Math, 'random').mockReturnValue(0.5);

      const response = await service.createPayment(mockPaymentRequest);

      expect(response.redirectUrl).not.toContain('require3ds=true');
    });

    it('должен симулировать API задержку', async () => {
      const _startTime = Date.now();

      await service.createPayment(mockPaymentRequest);

      // Продвигаем таймеры
      jest.advanceTimersByTime(mockConfig.apiDelay);

      expect(setTimeout).toHaveBeenCalledWith(
        expect.any(Function),
        mockConfig.apiDelay
      );
    });

    it('должен планировать webhook если 3DS не требуется', async () => {
      // Устанавливаем Math.random для НЕ требования 3DS
      jest.spyOn(Math, 'random').mockReturnValue(0.5);

      await service.createPayment(mockPaymentRequest);

      // Проверяем что webhook запланирован
      expect(setTimeout).toHaveBeenCalledWith(
        expect.any(Function),
        mockConfig.webhookDelay
      );
    });
  });

  describe('getPaymentStatus', () => {
    it('должен возвращать статус существующего платежа', async () => {
      // Сначала создаем платеж
      const createResponse = await service.createPayment(mockPaymentRequest);

      const status = await service.getPaymentStatus(createResponse.id);

      expect(status).toEqual({
        id: 'mock-payment-123',
        status: 'pending',
        amount: 1500,
        currency: 'RSD',
        createdAt: '2025-06-27T00:00:00.000Z',
        completedAt: null,
        card: {
          last4: '1111',
          brand: 'visa',
        },
      });
    });

    it('должен загружать платеж из localStorage если не найден в памяти', async () => {
      const mockStoredPayment = JSON.stringify({
        id: 'stored-payment-123',
        status: 'captured',
        amount: 2000,
        currency: 'RSD',
        createdAt: '2025-06-27T00:00:00.000Z',
        completedAt: '2025-06-27T00:05:00.000Z',
        card: { last4: '1111', brand: 'visa' },
      });

      mockLocalStorage.getItem.mockReturnValue(mockStoredPayment);

      const status = await service.getPaymentStatus('stored-payment-123');

      expect(mockLocalStorage.getItem).toHaveBeenCalledWith(
        'mock_payment_stored-payment-123'
      );
      expect(status.id).toBe('stored-payment-123');
      expect(status.status).toBe('captured');
    });

    it('должен выбрасывать ошибку для несуществующего платежа', async () => {
      mockLocalStorage.getItem.mockReturnValue(null);

      await expect(service.getPaymentStatus('nonexistent-123')).rejects.toThrow(
        'Payment nonexistent-123 not found'
      );
    });

    it('должен симулировать задержку при получении статуса', async () => {
      await service.createPayment(mockPaymentRequest);

      const _promise = service.getPaymentStatus('mock-payment-123');

      // Проверяем что задержка настроена
      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 500);
    });
  });

  describe('handle3DSecure', () => {
    beforeEach(async () => {
      await service.createPayment(mockPaymentRequest);
    });

    it('должен вернуть true для правильного кода 3DS', async () => {
      const result = await service.handle3DSecure('mock-payment-123', '123');

      expect(result).toBe(true);
    });

    it('должен вернуть false для неправильного кода 3DS', async () => {
      const result = await service.handle3DSecure('mock-payment-123', '456');

      expect(result).toBe(false);
    });

    it('должен запланировать webhook при успешной 3DS аутентификации', async () => {
      // Очищаем предыдущие вызовы setTimeout
      jest.clearAllMocks();

      await service.handle3DSecure('mock-payment-123', '123');

      expect(setTimeout).toHaveBeenCalledWith(
        expect.any(Function),
        mockConfig.webhookDelay
      );
    });

    it('должен пометить платеж как failed при неуспешной 3DS', async () => {
      await service.handle3DSecure('mock-payment-123', '456');

      const status = await service.getPaymentStatus('mock-payment-123');
      expect(status.status).toBe('failed');
      expect(status.completedAt).toBeTruthy();
    });

    it('должен симулировать задержку 3DS обработки', async () => {
      const _promise = service.handle3DSecure('mock-payment-123', '123');

      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 1500);
    });
  });

  describe('simulateWebhook', () => {
    beforeEach(async () => {
      await service.createPayment(mockPaymentRequest);
    });

    it('должен обновить статус платежа', async () => {
      await service.simulateWebhook('mock-payment-123', 'captured');

      const status = await service.getPaymentStatus('mock-payment-123');
      expect(status.status).toBe('captured');
      expect(status.completedAt).toBeTruthy();
    });

    it('должен установить completedAt для финальных статусов', async () => {
      const finalStatuses = ['captured', 'failed', 'cancelled'];

      for (const status of finalStatuses) {
        await service.simulateWebhook('mock-payment-123', status);
        const paymentStatus =
          await service.getPaymentStatus('mock-payment-123');
        expect(paymentStatus.completedAt).toBeTruthy();
      }
    });

    it('должен обновить localStorage при изменении статуса', async () => {
      await service.simulateWebhook('mock-payment-123', 'captured');

      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
        'mock_payment_mock-payment-123',
        expect.stringContaining('"status":"captured"')
      );
    });
  });

  describe('webhook scheduling', () => {
    it('должен выполнить webhook с правильным статусом на основе successRate', async () => {
      // Mock успешного результата (0.5 < 0.8)
      jest
        .spyOn(Math, 'random')
        .mockReturnValueOnce(0.5) // для 3DS проверки
        .mockReturnValueOnce(0.5); // для success rate

      await service.createPayment(mockPaymentRequest);

      // Продвигаем время до выполнения webhook
      jest.advanceTimersByTime(mockConfig.webhookDelay);

      const status = await service.getPaymentStatus('mock-payment-123');
      expect(status.status).toBe('captured');
    });

    it('должен выполнить webhook с failed статусом при низком successRate', async () => {
      // Mock неуспешного результата (0.9 >= 0.8)
      jest
        .spyOn(Math, 'random')
        .mockReturnValueOnce(0.5) // для 3DS проверки
        .mockReturnValueOnce(0.9); // для success rate

      await service.createPayment(mockPaymentRequest);

      // Продвигаем время до выполнения webhook
      jest.advanceTimersByTime(mockConfig.webhookDelay);

      const status = await service.getPaymentStatus('mock-payment-123');
      expect(status.status).toBe('failed');
    });
  });

  describe('cleanup', () => {
    it('должен очистить все активные таймауты', async () => {
      // Создаем несколько платежей чтобы запланировать webhook'и
      await service.createPayment({ ...mockPaymentRequest, listing_id: '1' });
      await service.createPayment({ ...mockPaymentRequest, listing_id: '2' });

      // Мокаем clearTimeout
      const clearTimeoutSpy = jest.spyOn(global, 'clearTimeout');

      service.cleanup();

      expect(clearTimeoutSpy).toHaveBeenCalled();
    });

    it('должен очистить Map с таймаутами', async () => {
      await service.createPayment(mockPaymentRequest);

      service.cleanup();

      // После cleanup не должно выполняться webhook'ов
      jest.advanceTimersByTime(mockConfig.webhookDelay);

      const status = await service.getPaymentStatus('mock-payment-123');
      expect(status.status).toBe('pending'); // Статус не должен измениться
    });
  });

  describe('edge cases', () => {
    it('должен обрабатывать отсутствие window объекта (SSR)', async () => {
      // Временно удаляем window
      const originalWindow = global.window;
      // @ts-expect-error -- Testing with invalid config
      delete global.window;

      const response = await service.createPayment(mockPaymentRequest);

      expect(response.id).toBe('mock-payment-123');
      expect(mockLocalStorage.setItem).not.toHaveBeenCalled();

      // Восстанавливаем window
      global.window = originalWindow;
    });

    it('должен обрабатывать некорректные данные в localStorage', async () => {
      mockLocalStorage.getItem.mockReturnValue('invalid json');

      await expect(service.getPaymentStatus('test-123')).rejects.toThrow(
        'Payment test-123 not found'
      );
    });

    it('должен обрабатывать обновление статуса несуществующего платежа', async () => {
      // Это не должно выбрасывать ошибку
      await expect(
        service.simulateWebhook('nonexistent-123', 'captured')
      ).resolves.not.toThrow();
    });
  });
});
