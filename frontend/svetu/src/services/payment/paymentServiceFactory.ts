import { useMemo } from 'react';
import type { IPaymentService } from '@/types/payment';
import type { PaymentConfig } from '@/config/payment';
import { paymentConfig } from '@/config/payment';
import { MockAllSecureService } from './mockAllSecureService';

export function createPaymentService(config: PaymentConfig): IPaymentService {
  switch (config.mode) {
    case 'mock':
      if (!config.mock) {
        throw new Error('Mock configuration is required in mock mode');
      }
      return new MockAllSecureService(config.mock.config);

    case 'sandbox':
    case 'production':
      // Когда получим credentials, здесь будет создаваться реальный AllSecureService
      // return new AllSecureService(config.allSecure!);
      throw new Error(
        `Mode ${config.mode} is not implemented yet. Use mock mode for development.`
      );

    default:
      throw new Error(`Unknown payment mode: ${config.mode}`);
  }
}

// Hook для использования в компонентах
export function usePaymentService() {
  const service = useMemo(() => createPaymentService(paymentConfig), []);
  return service;
}
