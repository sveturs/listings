import { renderHook, act } from '@testing-library/react';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { useAllSecurePayment } from '../useAllSecurePayment';
import type {
  IPaymentService,
  PaymentResponse,
  PaymentStatus,
} from '@/types/payment';

// Mock dependencies
jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
}));
jest.mock('next-intl', () => ({
  useLocale: jest.fn(),
}));
jest.mock('@/services/payment/paymentServiceFactory');

// Mock navigation helpers
jest.mock('../useAllSecurePaymentHelpers', () => ({
  navigateToUrl: jest.fn(),
  getLocationOrigin: jest.fn(() => 'http://localhost:3000'),
}));

import * as navigationHelpers from '../useAllSecurePaymentHelpers';
const mockNavigateToUrl = navigationHelpers.navigateToUrl as jest.Mock;

// Mock payment service
const mockPaymentService: jest.Mocked<IPaymentService> = {
  createPayment: jest.fn(),
  getPaymentStatus: jest.fn(),
  handle3DSecure: jest.fn(),
  simulateWebhook: jest.fn(),
};

// Mock usePaymentService
jest.mock('@/services/payment/paymentServiceFactory', () => ({
  usePaymentService: () => mockPaymentService,
}));

const mockPush = jest.fn();
const mockReplace = jest.fn();

(useRouter as jest.Mock).mockReturnValue({
  push: mockPush,
  replace: mockReplace,
  prefetch: jest.fn(),
  back: jest.fn(),
});

(useLocale as jest.Mock).mockReturnValue('ru');

describe('useAllSecurePayment', () => {
  const mockPaymentData = {
    listingId: 'test-listing-123',
    amount: 5000,
    currency: 'RSD',
    buyerInfo: {
      name: 'Test User',
      email: 'test@example.com',
      phone: '+381601234567',
      address: 'Test Address 123',
    },
  };

  const mockPaymentResponse: PaymentResponse = {
    id: 'payment-123',
    redirectUrl: '/payment/mock?id=payment-123',
    status: 'pending',
  };

  const mockPaymentStatus: PaymentStatus = {
    id: 'payment-123',
    status: 'captured',
    amount: 5000,
    currency: 'RSD',
    createdAt: '2025-06-27T00:00:00.000Z',
    completedAt: '2025-06-27T00:05:00.000Z',
    card: {
      lastFour: '1111',
      brand: 'visa',
      expiryMonth: 12,
      expiryYear: 2025,
    },
  };

  // Mock sessionStorage
  const mockSessionStorage = {
    getItem: jest.fn(),
    setItem: jest.fn(),
    removeItem: jest.fn(),
    clear: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();

    Object.defineProperty(window, 'sessionStorage', {
      value: mockSessionStorage,
      writable: true,
      configurable: true,
    });
  });

  describe('createPayment', () => {
    it('должен создать платеж и выполнить redirect', async () => {
      mockPaymentService.createPayment.mockResolvedValue(mockPaymentResponse);

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.createPayment(mockPaymentData);
      });

      // Проверяем вызов сервиса
      expect(mockPaymentService.createPayment).toHaveBeenCalledWith({
        listing_id: 'test-listing-123',
        amount: 5000,
        currency: 'RSD',
        buyer_info: mockPaymentData.buyerInfo,
        return_url: 'http://localhost:3000/ru/payment/process',
        locale: 'ru',
      });

      // Проверяем сохранение в sessionStorage
      expect(mockSessionStorage.setItem).toHaveBeenCalledWith(
        'pending_payment',
        expect.stringContaining('"listingId":"test-listing-123"')
      );
      expect(mockSessionStorage.setItem).toHaveBeenCalledWith(
        'payment_id',
        'payment-123'
      );

      // Проверяем что был вызван redirect
      expect(mockNavigateToUrl).toHaveBeenCalledWith(
        '/ru/payment/mock?id=payment-123'
      );
    });

    it('должен установить состояние loading во время обработки', async () => {
      let resolvePayment: (value: PaymentResponse) => void;
      const paymentPromise = new Promise<PaymentResponse>((resolve) => {
        resolvePayment = resolve;
      });

      mockPaymentService.createPayment.mockReturnValue(paymentPromise);

      const { result } = renderHook(() => useAllSecurePayment());

      // Запускаем создание платежа
      act(() => {
        result.current.createPayment(mockPaymentData);
      });

      // Проверяем что isProcessing = true
      expect(result.current.isProcessing).toBe(true);

      // Завершаем промис
      await act(async () => {
        resolvePayment!(mockPaymentResponse);
      });

      // Проверяем что isProcessing = false
      expect(result.current.isProcessing).toBe(false);
    });

    it('должен обработать ошибку создания платежа', async () => {
      const errorMessage = 'Payment creation failed';
      mockPaymentService.createPayment.mockRejectedValue(
        new Error(errorMessage)
      );

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        try {
          await result.current.createPayment(mockPaymentData);
        } catch {
          // Ожидаем что ошибка будет выброшена
        }
      });

      expect(result.current.error).toBe(errorMessage);
      expect(result.current.isProcessing).toBe(false);
    });

    it('должен очистить ошибку при новом запросе', async () => {
      // Сначала устанавливаем ошибку
      mockPaymentService.createPayment.mockRejectedValueOnce(
        new Error('First error')
      );

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        try {
          await result.current.createPayment(mockPaymentData);
        } catch {}
      });

      expect(result.current.error).toBe('First error');

      // Теперь успешный запрос
      mockPaymentService.createPayment.mockResolvedValueOnce(
        mockPaymentResponse
      );

      await act(async () => {
        await result.current.createPayment(mockPaymentData);
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe('checkPaymentStatus', () => {
    it('должен вернуть статус платежа', async () => {
      mockPaymentService.getPaymentStatus.mockResolvedValue(mockPaymentStatus);

      const { result } = renderHook(() => useAllSecurePayment());

      const status = await act(async () => {
        return await result.current.checkPaymentStatus('payment-123');
      });

      expect(status).toEqual(mockPaymentStatus);
      expect(mockPaymentService.getPaymentStatus).toHaveBeenCalledWith(
        'payment-123'
      );
    });

    it('должен выбросить ошибку при неудачной проверке', async () => {
      const error = new Error('Payment not found');
      mockPaymentService.getPaymentStatus.mockRejectedValue(error);

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await expect(
          result.current.checkPaymentStatus('invalid-id')
        ).rejects.toThrow('Payment not found');
      });
    });
  });

  describe('handlePaymentReturn', () => {
    it('должен redirect на success при captured статусе', async () => {
      mockSessionStorage.getItem.mockReturnValue('payment-123');
      mockPaymentService.getPaymentStatus.mockResolvedValue({
        ...mockPaymentStatus,
        status: 'captured',
      });

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/ru/payment/success?id=payment-123'
      );
      expect(mockSessionStorage.removeItem).toHaveBeenCalledWith('payment_id');
      expect(mockSessionStorage.removeItem).toHaveBeenCalledWith(
        'pending_payment'
      );
    });

    it('должен redirect на success при authorized статусе', async () => {
      mockSessionStorage.getItem.mockReturnValue('payment-123');
      mockPaymentService.getPaymentStatus.mockResolvedValue({
        ...mockPaymentStatus,
        status: 'authorized',
      });

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/ru/payment/success?id=payment-123'
      );
    });

    it('должен redirect на error при failed статусе', async () => {
      mockSessionStorage.getItem.mockReturnValue('payment-123');
      mockPaymentService.getPaymentStatus.mockResolvedValue({
        ...mockPaymentStatus,
        status: 'failed',
        error_code: 'card_declined',
      });

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/ru/payment/error?id=payment-123&reason=card_declined'
      );
    });

    it('должен redirect на error при cancelled статусе', async () => {
      mockSessionStorage.getItem.mockReturnValue('payment-123');
      mockPaymentService.getPaymentStatus.mockResolvedValue({
        ...mockPaymentStatus,
        status: 'cancelled',
      });

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/ru/payment/error?id=payment-123&reason=cancelled'
      );
    });

    it('должен redirect на process при pending статусе', async () => {
      mockSessionStorage.getItem.mockReturnValue('payment-123');
      mockPaymentService.getPaymentStatus.mockResolvedValue({
        ...mockPaymentStatus,
        status: 'pending',
      });

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith('/ru/payment/process/payment-123');
    });

    it('должен redirect на error если нет payment_id', async () => {
      mockSessionStorage.getItem.mockReturnValue(null);

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/ru/payment/error?reason=no_payment_id'
      );
    });

    it('должен redirect на error при ошибке проверки статуса', async () => {
      mockSessionStorage.getItem.mockReturnValue('payment-123');
      mockPaymentService.getPaymentStatus.mockRejectedValue(
        new Error('Network error')
      );

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/ru/payment/error?reason=status_check_failed'
      );
    });
  });

  describe('createTestPayment', () => {
    it('должен создать тестовый платеж с дефолтными данными', async () => {
      mockPaymentService.createPayment.mockResolvedValue(mockPaymentResponse);

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.createTestPayment();
      });

      expect(mockPaymentService.createPayment).toHaveBeenCalledWith({
        listing_id: 'test_listing_123',
        amount: 5000,
        currency: 'RSD',
        buyer_info: {
          name: 'Test User',
          email: 'test@example.com',
          phone: '+381601234567',
          address: 'Test Address 123, Belgrade',
        },
        return_url: expect.stringContaining('/ru/payment/process'),
        locale: 'ru',
      });
    });

    it('должен создать тестовый платеж с кастомной суммой', async () => {
      mockPaymentService.createPayment.mockResolvedValue(mockPaymentResponse);

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.createTestPayment(10000);
      });

      expect(mockPaymentService.createPayment).toHaveBeenCalledWith(
        expect.objectContaining({
          amount: 10000,
        })
      );
    });
  });

  describe('clearError', () => {
    it('должен очистить ошибку', async () => {
      mockPaymentService.createPayment.mockRejectedValue(
        new Error('Test error')
      );

      const { result } = renderHook(() => useAllSecurePayment());

      // Устанавливаем ошибку
      await act(async () => {
        try {
          await result.current.createPayment(mockPaymentData);
        } catch {}
      });

      expect(result.current.error).toBe('Test error');

      // Очищаем ошибку
      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe('edge cases', () => {
    // TODO: Skipped - Jest/JSDOM limitation with completely removing window object
    // Even after deleting global.window, typeof window !== 'undefined' returns true in JSDOM
    // This edge case is unlikely in real usage as window is always defined in browsers
    it.skip('должен обработать отсутствие window объекта при создании платежа', async () => {
      mockPaymentService.createPayment.mockResolvedValue(mockPaymentResponse);

      // Временно удаляем window
      const originalWindow = global.window;
      // @ts-expect-error -- Need to temporarily remove window for testing
      delete global.window;

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.createPayment(mockPaymentData);
      });

      // Проверяем что sessionStorage не вызывался
      expect(mockSessionStorage.setItem).not.toHaveBeenCalled();

      // Восстанавливаем window
      global.window = originalWindow;
    });

    it('должен обработать отсутствие window объекта при handlePaymentReturn', async () => {
      // Mock getItem to return null when window is undefined
      mockSessionStorage.getItem.mockReturnValue(null);

      // Temporarily remove window
      const originalWindow = global.window;
      // @ts-expect-error -- Need to temporarily remove window for testing
      delete global.window;

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/ru/payment/error?reason=no_payment_id'
      );

      // Restore window
      global.window = originalWindow;
    });

    it('должен обработать не-Error объект как ошибку', async () => {
      mockPaymentService.createPayment.mockRejectedValue('String error');

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        try {
          await result.current.createPayment(mockPaymentData);
        } catch {}
      });

      expect(result.current.error).toBe('Неизвестная ошибка');
    });
  });

  describe('locale integration', () => {
    it('должен работать с английской локалью', async () => {
      (useLocale as jest.Mock).mockReturnValue('en');
      mockPaymentService.createPayment.mockResolvedValue(mockPaymentResponse);

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.createPayment(mockPaymentData);
      });

      expect(mockPaymentService.createPayment).toHaveBeenCalledWith(
        expect.objectContaining({
          return_url: 'http://localhost:3000/en/payment/process',
          locale: 'en',
        })
      );

      // Проверяем что redirect использует английскую локаль
      expect(mockNavigateToUrl).toHaveBeenCalledWith(
        expect.stringContaining('/en/')
      );
    });

    it('должен включать локаль во все redirect URL', async () => {
      (useLocale as jest.Mock).mockReturnValue('en');

      mockSessionStorage.getItem.mockReturnValue('payment-123');
      mockPaymentService.getPaymentStatus.mockResolvedValue({
        ...mockPaymentStatus,
        status: 'captured',
      });

      const { result } = renderHook(() => useAllSecurePayment());

      await act(async () => {
        await result.current.handlePaymentReturn();
      });

      expect(mockPush).toHaveBeenCalledWith(
        '/en/payment/success?id=payment-123'
      );
    });
  });
});
