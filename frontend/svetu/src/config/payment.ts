export interface MockConfig {
  successRate: number; // Вероятность успешного платежа (0-1)
  require3DSRate: number; // Вероятность 3D Secure (0-1)
  apiDelay: number; // Задержка API в мс
  webhookDelay: number; // Задержка webhook в мс
  debugMode?: boolean; // Режим отладки
}

export interface TestCard {
  number: string;
  type: string;
  description: string;
}

export interface PaymentConfig {
  mode: 'production' | 'sandbox' | 'mock';
  mock?: {
    enabled: boolean;
    config: MockConfig;
    testCards: TestCard[];
    debugMode: boolean;
  };
  allSecure?: {
    merchantId: string;
    apiUrl: string;
    widgetUrl: string;
    apiKey?: string;
  };
}

export const paymentConfig: PaymentConfig = {
  mode: (process.env.NEXT_PUBLIC_PAYMENT_MODE as any) || 'mock',
  mock: {
    enabled: true,
    config: {
      successRate: 0.8,
      require3DSRate: 0.3,
      apiDelay: 1500,
      webhookDelay: 3000,
    },
    testCards: [
      {
        number: '4111111111111111',
        type: 'success',
        description: 'Всегда успешно',
      },
      {
        number: '4000000000000002',
        type: 'declined',
        description: 'Всегда отклонено',
      },
      {
        number: '4000000000003220',
        type: '3ds_required',
        description: 'Требует 3D Secure',
      },
      {
        number: '4000000000000077',
        type: 'timeout',
        description: 'Истекает время',
      },
    ],
    debugMode: process.env.NODE_ENV === 'development',
  },
};
