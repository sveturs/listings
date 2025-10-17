import type {
  IPaymentService,
  PaymentRequest,
  PaymentResponse,
  PaymentStatus,
  MockPayment,
} from '@/types/payment';
import type { MockConfig } from '@/config/payment';
import { generateMockPayment } from '@/utils/mockDataGenerator';

export class MockAllSecureService implements IPaymentService {
  private config: MockConfig;
  private storage: Map<string, MockPayment>;
  private webhookTimeouts: Map<string, NodeJS.Timeout>;

  constructor(config: MockConfig) {
    this.config = config;
    this.storage = new Map();
    this.webhookTimeouts = new Map();
  }

  async createPayment(data: PaymentRequest): Promise<PaymentResponse> {
    console.log('Mock: Creating payment', data);

    // Симуляция задержки API
    await this.delay(this.config.apiDelay);

    const payment = generateMockPayment();
    payment.amount = data.amount;
    payment.currency = data.currency;
    payment.listing_id = data.listing_id;
    payment.buyer_info = data.buyer_info;

    // Определяем нужен ли 3D Secure на основе номера карты из тестовых данных
    const requires3DS = Math.random() < this.config.require3DSRate;

    // Сохраняем в локальное хранилище
    this.storage.set(payment.id, payment);

    // Также сохраняем в localStorage для доступа между страницами
    if (typeof window !== 'undefined') {
      localStorage.setItem(
        `mock_payment_${payment.id}`,
        JSON.stringify(payment)
      );
    }

    // Симулируем webhook через заданное время (только если не требуется 3DS)
    if (!requires3DS) {
      this.scheduleWebhook(payment.id);
    }

    const locale = data.locale || 'ru';

    return {
      id: payment.id,
      redirectUrl: requires3DS
        ? `/${locale}/payment/mock?id=${payment.id}&require3ds=true`
        : `/${locale}/payment/mock?id=${payment.id}`,
      status: 'pending',
    };
  }

  async getPaymentStatus(paymentId: string): Promise<PaymentStatus> {
    console.log('Mock: Getting payment status', paymentId);

    await this.delay(500); // Небольшая задержка для реалистичности

    let payment = this.storage.get(paymentId);

    // Если не найдено в памяти, пытаемся загрузить из localStorage
    if (!payment && typeof window !== 'undefined') {
      const stored = localStorage.getItem(`mock_payment_${paymentId}`);
      if (stored) {
        try {
          payment = JSON.parse(stored);
          this.storage.set(paymentId, payment!);
        } catch {
          // Игнорируем некорректный JSON в localStorage
          console.warn(`Invalid payment data in localStorage for ${paymentId}`);
        }
      }
    }

    if (!payment) {
      throw new Error(`Payment ${paymentId} not found`);
    }

    return {
      id: payment.id,
      status: payment.status as any,
      amount: payment.amount,
      currency: payment.currency,
      createdAt: payment.createdAt,
      completedAt: payment.completedAt,
      card: payment.card,
    };
  }

  async handle3DSecure(paymentId: string, code: string): Promise<boolean> {
    console.log('Mock: Handling 3D Secure', paymentId, code);

    await this.delay(1500);

    // Простая валидация - код должен быть "123" для успеха
    const success = code === '123';

    if (success) {
      // Запускаем webhook после успешной 3DS аутентификации
      this.scheduleWebhook(paymentId);
    } else {
      // Помечаем как неуспешный
      this.updatePaymentStatus(paymentId, 'failed');
    }

    return success;
  }

  async simulateWebhook(paymentId: string, status: string): Promise<void> {
    console.log('Mock: Simulating webhook', paymentId, status);

    this.updatePaymentStatus(paymentId, status);

    // В реальной системе здесь был бы HTTP запрос к нашему webhook endpoint
    // Для mock версии просто обновляем статус
  }

  private scheduleWebhook(paymentId: string) {
    const timeoutId = setTimeout(async () => {
      const success = Math.random() < this.config.successRate;
      await this.simulateWebhook(paymentId, success ? 'captured' : 'failed');
      this.webhookTimeouts.delete(paymentId);
    }, this.config.webhookDelay);

    this.webhookTimeouts.set(paymentId, timeoutId);
  }

  private updatePaymentStatus(paymentId: string, status: string) {
    const payment = this.storage.get(paymentId);
    if (payment) {
      payment.status = status;
      if (['captured', 'failed', 'cancelled'].includes(status)) {
        payment.completedAt = new Date().toISOString();
      }

      this.storage.set(paymentId, payment);

      // Обновляем в localStorage
      if (typeof window !== 'undefined') {
        localStorage.setItem(
          `mock_payment_${paymentId}`,
          JSON.stringify(payment)
        );
      }
    }
  }

  private delay(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  // Метод для очистки таймаутов (для cleanup)
  cleanup() {
    for (const timeoutId of this.webhookTimeouts.values()) {
      clearTimeout(timeoutId);
    }
    this.webhookTimeouts.clear();
  }
}
