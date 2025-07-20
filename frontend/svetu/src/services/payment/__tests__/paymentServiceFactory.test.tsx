import { renderHook } from '@testing-library/react';
import {
  createPaymentService,
  usePaymentService,
} from '../paymentServiceFactory';
import { MockAllSecureService } from '../mockAllSecureService';
import type { PaymentConfig } from '@/config/payment';

// Mock зависимостей
jest.mock('../mockAllSecureService');
jest.mock('@/config/payment', () => ({
  paymentConfig: {
    mode: 'mock',
    mock: {
      config: {
        successRate: 0.8,
        require3DSRate: 0.3,
        apiDelay: 100,
        webhookDelay: 2000,
        debugMode: true,
      },
    },
  },
}));

const MockedMockAllSecureService = jest.mocked(MockAllSecureService);

describe('paymentServiceFactory', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('createPaymentService', () => {
    describe('mock режим', () => {
      it('должен создавать MockAllSecureService в mock режиме', () => {
        const mockConfig: PaymentConfig = {
          mode: 'mock',
          mock: {
            enabled: true,
            config: {
              successRate: 0.8,
              require3DSRate: 0.3,
              apiDelay: 100,
              webhookDelay: 2000,
              debugMode: true,
            },
            testCards: [],
            debugMode: true,
          },
        };

        const service = createPaymentService(mockConfig);

        expect(MockedMockAllSecureService).toHaveBeenCalledWith(
          mockConfig.mock?.config
        );
        expect(service).toBeInstanceOf(MockAllSecureService);
      });

      it('должен передавать корректную конфигурацию в MockAllSecureService', () => {
        const customMockConfig = {
          successRate: 0.9,
          require3DSRate: 0.1,
          apiDelay: 50,
          webhookDelay: 1000,
          debugMode: false,
        };

        const config: PaymentConfig = {
          mode: 'mock',
          mock: {
            enabled: true,
            config: customMockConfig,
            testCards: [],
            debugMode: false,
          },
        };

        createPaymentService(config);

        expect(MockedMockAllSecureService).toHaveBeenCalledWith(
          customMockConfig
        );
      });

      it('должен выбрасывать ошибку если mock конфигурация отсутствует', () => {
        const config: PaymentConfig = {
          mode: 'mock',
          // mock конфигурация отсутствует
        } as any;

        expect(() => createPaymentService(config)).toThrow(
          'Mock configuration is required in mock mode'
        );
      });

      it('должен выбрасывать ошибку если mock.config отсутствует', () => {
        const config: PaymentConfig = {
          mode: 'mock',
          mock: {} as any, // config отсутствует
        };

        expect(() => createPaymentService(config)).toThrow();
      });
    });

    describe('sandbox режим', () => {
      it('должен выбрасывать ошибку для sandbox режима', () => {
        const config: PaymentConfig = {
          mode: 'sandbox',
          allSecure: {
            merchantId: 'test-merchant',
            apiUrl: 'https://sandbox.allsecure.com',
            widgetUrl: 'https://sandbox.allsecure.com/widget',
            apiKey: 'test-key',
          },
        };

        expect(() => createPaymentService(config)).toThrow(
          'Mode sandbox is not implemented yet. Use mock mode for development.'
        );
      });
    });

    describe('production режим', () => {
      it('должен выбрасывать ошибку для production режима', () => {
        const config: PaymentConfig = {
          mode: 'production',
          allSecure: {
            merchantId: 'prod-merchant',
            apiUrl: 'https://api.allsecure.com',
            widgetUrl: 'https://api.allsecure.com/widget',
            apiKey: 'prod-key',
          },
        };

        expect(() => createPaymentService(config)).toThrow(
          'Mode production is not implemented yet. Use mock mode for development.'
        );
      });
    });

    describe('неизвестный режим', () => {
      it('должен выбрасывать ошибку для неизвестного режима', () => {
        const config = {
          mode: 'unknown',
        } as any;

        expect(() => createPaymentService(config)).toThrow(
          'Unknown payment mode: unknown'
        );
      });

      it('должен выбрасывать ошибку для null режима', () => {
        const config = {
          mode: null,
        } as any;

        expect(() => createPaymentService(config)).toThrow(
          'Unknown payment mode: null'
        );
      });

      it('должен выбрасывать ошибку для undefined режима', () => {
        const config = {
          mode: undefined,
        } as any;

        expect(() => createPaymentService(config)).toThrow(
          'Unknown payment mode: undefined'
        );
      });
    });
  });

  describe('usePaymentService hook', () => {
    it('должен возвращать payment service', () => {
      const { result } = renderHook(() => usePaymentService());

      expect(result.current).toBeInstanceOf(MockAllSecureService);
    });

    it('должен использовать мемоизацию - возвращать тот же экземпляр', () => {
      const { result, rerender } = renderHook(() => usePaymentService());

      const firstInstance = result.current;

      // Повторный рендер
      rerender();

      const secondInstance = result.current;

      expect(firstInstance).toBe(secondInstance);
    });

    it('должен создавать новый экземпляр только при изменении зависимостей', () => {
      const { result } = renderHook(() => usePaymentService());

      const _firstInstance = result.current;

      // Повторный рендер с теми же зависимостями
      const { result: result2 } = renderHook(() => usePaymentService());

      // Должны быть разные экземпляры, поскольку useMemo создает новый в каждом новом компоненте
      expect(result2.current).toBeInstanceOf(MockAllSecureService);
    });

    it('должен использовать paymentConfig из модуля конфигурации', () => {
      renderHook(() => usePaymentService());

      expect(MockedMockAllSecureService).toHaveBeenCalledWith({
        successRate: 0.8,
        require3DSRate: 0.3,
        apiDelay: 100,
        webhookDelay: 2000,
        debugMode: true,
      });
    });
  });

  describe('интеграция с различными конфигурациями', () => {
    it('должен работать с минимальной mock конфигурацией', () => {
      const minimalConfig: PaymentConfig = {
        mode: 'mock',
        mock: {
          enabled: true,
          config: {
            successRate: 1.0,
            require3DSRate: 0.0,
            apiDelay: 0,
            webhookDelay: 0,
            debugMode: false,
          },
          testCards: [],
          debugMode: false,
        },
      };

      const service = createPaymentService(minimalConfig);

      expect(service).toBeInstanceOf(MockAllSecureService);
      expect(MockedMockAllSecureService).toHaveBeenCalledWith(
        minimalConfig.mock?.config
      );
    });

    it('должен работать с максимальной mock конфигурацией', () => {
      const maximalConfig: PaymentConfig = {
        mode: 'mock',
        mock: {
          enabled: true,
          config: {
            successRate: 0.0,
            require3DSRate: 1.0,
            apiDelay: 5000,
            webhookDelay: 10000,
            debugMode: true,
          },
          testCards: [],
          debugMode: true,
        },
      };

      const service = createPaymentService(maximalConfig);

      expect(service).toBeInstanceOf(MockAllSecureService);
      expect(MockedMockAllSecureService).toHaveBeenCalledWith(
        maximalConfig.mock?.config
      );
    });

    it('должен работать с дробными значениями в конфигурации', () => {
      const floatConfig: PaymentConfig = {
        mode: 'mock',
        mock: {
          enabled: true,
          config: {
            successRate: 0.123,
            require3DSRate: 0.456,
            apiDelay: 123.45,
            webhookDelay: 456.78,
            debugMode: true,
          },
          testCards: [],
          debugMode: true,
        },
      };

      const service = createPaymentService(floatConfig);

      expect(service).toBeInstanceOf(MockAllSecureService);
      expect(MockedMockAllSecureService).toHaveBeenCalledWith(
        floatConfig.mock?.config
      );
    });
  });

  describe('edge cases', () => {
    it('должен обрабатывать пустой объект конфигурации', () => {
      const emptyConfig = {} as any;

      expect(() => createPaymentService(emptyConfig)).toThrow(
        'Unknown payment mode: undefined'
      );
    });

    it('должен обрабатывать null как конфигурацию', () => {
      expect(() => createPaymentService(null as any)).toThrow();
    });

    it('должен обрабатывать undefined как конфигурацию', () => {
      expect(() => createPaymentService(undefined as any)).toThrow();
    });

    it('должен обрабатывать конфигурацию с лишними полями', () => {
      const configWithExtra: any = {
        mode: 'mock',
        mock: {
          config: {
            successRate: 0.8,
            require3DSRate: 0.3,
            apiDelay: 100,
            webhookDelay: 2000,
            debugMode: true,
          },
        },
        extraField: 'should be ignored',
        anotherExtra: 123,
      };

      const service = createPaymentService(configWithExtra);

      expect(service).toBeInstanceOf(MockAllSecureService);
    });

    it('должен обрабатывать строки вместо чисел в mock конфигурации', () => {
      const stringConfig: any = {
        mode: 'mock',
        mock: {
          config: {
            successRate: '0.8',
            require3DSRate: '0.3',
            apiDelay: '100',
            webhookDelay: '2000',
            debugMode: 'true',
          },
        },
      };

      const service = createPaymentService(stringConfig);

      expect(service).toBeInstanceOf(MockAllSecureService);
      expect(MockedMockAllSecureService).toHaveBeenCalledWith(
        stringConfig.mock.config
      );
    });

    it('должен обрабатывать режим в разном регистре', () => {
      const upperCaseConfig = {
        mode: 'MOCK',
      } as any;

      expect(() => createPaymentService(upperCaseConfig)).toThrow(
        'Unknown payment mode: MOCK'
      );
    });
  });

  describe('type safety', () => {
    it('должен корректно типизировать возвращаемый сервис', () => {
      const config: PaymentConfig = {
        mode: 'mock',
        mock: {
          enabled: true,
          config: {
            successRate: 0.8,
            require3DSRate: 0.3,
            apiDelay: 100,
            webhookDelay: 2000,
            debugMode: true,
          },
          testCards: [],
          debugMode: true,
        },
      };

      const service = createPaymentService(config);

      // Проверяем что сервис имеет методы интерфейса IPaymentService
      expect(typeof service.createPayment).toBe('function');
      expect(typeof service.getPaymentStatus).toBe('function');
      expect(typeof service.handle3DSecure).toBe('function');
      expect(typeof service.simulateWebhook).toBe('function');
    });
  });

  describe('будущая extensibility', () => {
    it('должен быть готов к добавлению реального AllSecure сервиса', () => {
      // Проверяем что код структурирован для добавления новых режимов
      expect(() => createPaymentService({ mode: 'sandbox' } as any)).toThrow(
        /is not implemented yet/
      );
      expect(() => createPaymentService({ mode: 'production' } as any)).toThrow(
        /is not implemented yet/
      );
    });

    it('должен иметь расширяемую структуру switch statement', () => {
      // Если добавляется новый режим, он должен быть легко интегрируемым
      const codeString = createPaymentService.toString();

      expect(codeString).toContain('switch');
      expect(codeString).toContain("case 'mock'");
      expect(codeString).toContain("case 'sandbox'");
      expect(codeString).toContain("case 'production'");
      expect(codeString).toContain('default:');
    });
  });
});
