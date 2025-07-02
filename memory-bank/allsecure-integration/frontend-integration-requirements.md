# Технические требования для интеграции AllSecure на Frontend

## 1. Архитектура компонентов

### 1.1 AllSecureWidgetLoader
Компонент для асинхронной загрузки AllSecure JavaScript SDK.

```typescript
interface AllSecureWidgetLoaderProps {
  merchantId: string
  environment: 'sandbox' | 'production'
  onLoad: () => void
  onError: (error: Error) => void
}
```

**Функциональность:**
- Динамическая загрузка скрипта AllSecure
- Проверка доступности глобального объекта AllSecure
- Retry механизм при ошибке загрузки
- Cleanup при размонтировании

### 1.2 useAllSecurePayment Hook
React hook для управления платежным процессом.

```typescript
interface UseAllSecurePaymentOptions {
  autoCapture?: boolean
  webhookUrl?: string
  timeout?: number
}

interface UseAllSecurePaymentReturn {
  createPayment: (data: PaymentData) => Promise<Payment>
  getStatus: (paymentId: string) => Promise<PaymentStatus>
  capturePayment: (paymentId: string) => Promise<void>
  refundPayment: (paymentId: string, amount?: number) => Promise<void>
  isLoading: boolean
  error: PaymentError | null
}
```

### 1.3 AllSecureCheckout Component
Основной компонент для проведения платежа.

```typescript
interface AllSecureCheckoutProps {
  amount: number
  currency: string
  orderId: string
  customerInfo: CustomerInfo
  onSuccess: (payment: Payment) => void
  onError: (error: PaymentError) => void
  onCancel: () => void
  mode?: 'widget' | 'redirect'
}
```

## 2. Сценарии интеграции

### 2.1 Widget Flow (предпочтительный)

```typescript
// 1. Инициализация widget
const widget = new AllSecure.Widget({
  merchantId: MERCHANT_ID,
  amount: 1000.00,
  currency: 'RSD',
  language: 'sr',
  style: {
    primaryColor: '#570df8',
    fontFamily: 'system-ui'
  }
})

// 2. Монтирование в DOM
widget.mount('#payment-container')

// 3. Обработка событий
widget.on('payment.success', (data) => {
  // Обработка успешного платежа
})

widget.on('payment.error', (error) => {
  // Обработка ошибки
})

widget.on('3ds.required', (data) => {
  // Обработка 3D Secure
})
```

### 2.2 Redirect Flow (запасной вариант)

```typescript
// 1. Создание платежа на backend
const payment = await api.createPayment({
  listingId,
  amount,
  returnUrl: `${window.location.origin}/payment/return`
})

// 2. Сохранение состояния
sessionStorage.setItem('payment_id', payment.id)
sessionStorage.setItem('order_data', JSON.stringify(orderData))

// 3. Redirect на AllSecure
window.location.href = payment.redirectUrl
```

## 3. Обработка состояний

### 3.1 Состояния платежа

```typescript
enum PaymentState {
  IDLE = 'idle',
  INITIALIZING = 'initializing',
  PROCESSING = 'processing',
  AWAITING_3DS = 'awaiting_3ds',
  COMPLETING = 'completing',
  SUCCESS = 'success',
  ERROR = 'error',
  CANCELLED = 'cancelled'
}
```

### 3.2 State Machine

```typescript
const paymentMachine = {
  initial: PaymentState.IDLE,
  states: {
    [PaymentState.IDLE]: {
      on: { START: PaymentState.INITIALIZING }
    },
    [PaymentState.INITIALIZING]: {
      on: {
        INITIALIZED: PaymentState.PROCESSING,
        ERROR: PaymentState.ERROR
      }
    },
    [PaymentState.PROCESSING]: {
      on: {
        REQUIRE_3DS: PaymentState.AWAITING_3DS,
        SUCCESS: PaymentState.COMPLETING,
        ERROR: PaymentState.ERROR
      }
    },
    [PaymentState.AWAITING_3DS]: {
      on: {
        COMPLETE: PaymentState.COMPLETING,
        ERROR: PaymentState.ERROR
      }
    },
    [PaymentState.COMPLETING]: {
      on: {
        DONE: PaymentState.SUCCESS,
        ERROR: PaymentState.ERROR
      }
    }
  }
}
```

## 4. Безопасность

### 4.1 PCI DSS Compliance

**Уровень SAQ-A (при использовании widget):**
- Карточные данные не касаются наших серверов
- AllSecure Widget создает изолированный iframe
- Данные передаются напрямую в AllSecure

**Требования:**
- HTTPS на всех страницах с платежами
- Content Security Policy headers
- Регулярные обновления зависимостей

### 4.2 Валидация данных

```typescript
// Валидация на frontend
const validatePaymentData = (data: PaymentData): ValidationResult => {
  const errors: ValidationError[] = []
  
  if (data.amount <= 0) {
    errors.push({ field: 'amount', message: 'Amount must be positive' })
  }
  
  if (!SUPPORTED_CURRENCIES.includes(data.currency)) {
    errors.push({ field: 'currency', message: 'Unsupported currency' })
  }
  
  // Email, phone validation...
  
  return { valid: errors.length === 0, errors }
}
```

### 4.3 Безопасная обработка результатов

```typescript
// Проверка подписи webhook
const verifyWebhookSignature = (
  payload: string,
  signature: string,
  secret: string
): boolean => {
  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(payload)
    .digest('hex')
    
  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(expectedSignature)
  )
}
```

## 5. UX/UI рекомендации

### 5.1 Индикаторы состояния

```typescript
const PaymentStateIndicator = ({ state }: { state: PaymentState }) => {
  const indicators = {
    [PaymentState.INITIALIZING]: {
      icon: <Loader className="animate-spin" />,
      text: t('payment.initializing'),
      color: 'text-blue-600'
    },
    [PaymentState.PROCESSING]: {
      icon: <CreditCard className="animate-pulse" />,
      text: t('payment.processing'),
      color: 'text-yellow-600'
    },
    [PaymentState.SUCCESS]: {
      icon: <CheckCircle />,
      text: t('payment.success'),
      color: 'text-green-600'
    },
    // ...
  }
  
  const indicator = indicators[state]
  
  return (
    <div className={`flex items-center gap-2 ${indicator.color}`}>
      {indicator.icon}
      <span>{indicator.text}</span>
    </div>
  )
}
```

### 5.2 Локализация

```typescript
// Сообщения для AllSecure widget
const widgetMessages = {
  sr: {
    cardNumber: 'Број картице',
    expiryDate: 'Датум истека',
    cvv: 'CVV код',
    cardholderName: 'Име на картици',
    pay: 'Платите {amount}',
    processing: 'Обрада плаћања...',
    secure: 'Безбедно плаћање'
  },
  en: {
    cardNumber: 'Card Number',
    expiryDate: 'Expiry Date',
    cvv: 'CVV Code',
    cardholderName: 'Cardholder Name',
    pay: 'Pay {amount}',
    processing: 'Processing payment...',
    secure: 'Secure payment'
  }
}
```

### 5.3 Мобильная адаптация

```css
/* Адаптивные стили для widget контейнера */
.allsecure-widget-container {
  width: 100%;
  max-width: 480px;
  margin: 0 auto;
}

@media (max-width: 640px) {
  .allsecure-widget-container {
    padding: 0 1rem;
  }
  
  /* Увеличенные touch targets */
  .payment-button {
    min-height: 48px;
    font-size: 16px; /* Предотвращает zoom на iOS */
  }
}
```

## 6. Обработка ошибок

### 6.1 Типы ошибок

```typescript
interface PaymentError {
  code: string
  message: string
  details?: Record<string, any>
  recoverable: boolean
}

const ERROR_CODES = {
  // Network errors
  NETWORK_ERROR: 'network_error',
  TIMEOUT: 'timeout',
  
  // Payment errors
  CARD_DECLINED: 'card_declined',
  INSUFFICIENT_FUNDS: 'insufficient_funds',
  EXPIRED_CARD: 'expired_card',
  INVALID_CVV: 'invalid_cvv',
  
  // Business errors
  LIMIT_EXCEEDED: 'limit_exceeded',
  DUPLICATE_TRANSACTION: 'duplicate_transaction',
  
  // Technical errors
  WIDGET_LOAD_FAILED: 'widget_load_failed',
  INVALID_CONFIGURATION: 'invalid_configuration'
}
```

### 6.2 Стратегии восстановления

```typescript
const errorRecoveryStrategies = {
  [ERROR_CODES.NETWORK_ERROR]: {
    action: 'retry',
    message: t('errors.network_retry'),
    maxRetries: 3
  },
  [ERROR_CODES.WIDGET_LOAD_FAILED]: {
    action: 'fallback',
    message: t('errors.widget_fallback'),
    fallbackTo: 'redirect'
  },
  [ERROR_CODES.CARD_DECLINED]: {
    action: 'user_action',
    message: t('errors.card_declined'),
    suggestion: t('errors.try_another_card')
  }
}
```

## 7. Тестирование

### 7.1 Unit тесты

```typescript
describe('useAllSecurePayment', () => {
  it('should create payment successfully', async () => {
    const { result } = renderHook(() => useAllSecurePayment())
    
    const payment = await result.current.createPayment({
      amount: 1000,
      currency: 'RSD'
    })
    
    expect(payment.id).toBeDefined()
    expect(payment.status).toBe('pending')
  })
  
  it('should handle network errors', async () => {
    // Mock network error
    server.use(
      rest.post('/api/payments', (req, res, ctx) => {
        return res.networkError('Failed to connect')
      })
    )
    
    const { result } = renderHook(() => useAllSecurePayment())
    
    await expect(
      result.current.createPayment({ amount: 1000 })
    ).rejects.toThrow('Network error')
  })
})
```

### 7.2 E2E тесты

```typescript
describe('Payment Flow E2E', () => {
  it('should complete payment with 3D Secure', () => {
    cy.visit('/checkout')
    
    // Fill checkout form
    cy.fillCheckoutForm({
      name: 'Test User',
      email: 'test@example.com',
      phone: '+381601234567'
    })
    
    // Select card payment
    cy.get('[data-cy=payment-method-card]').click()
    
    // Wait for widget
    cy.get('[data-cy=allsecure-widget]').should('be.visible')
    
    // Fill card details in iframe
    cy.fillCardDetails({
      number: '4111111111111111',
      expiry: '12/25',
      cvv: '123'
    })
    
    // Submit payment
    cy.get('[data-cy=pay-button]').click()
    
    // Handle 3D Secure
    cy.handle3DSecure()
    
    // Verify success
    cy.url().should('include', '/payment/success')
    cy.contains('Payment successful')
  })
})
```

## 8. Мониторинг и аналитика

### 8.1 Метрики

```typescript
// Отслеживание событий платежей
const trackPaymentEvent = (event: string, data: any) => {
  // Google Analytics
  gtag('event', event, {
    event_category: 'Payment',
    event_label: data.paymentMethod,
    value: data.amount
  })
  
  // Custom analytics
  analytics.track(event, {
    ...data,
    timestamp: Date.now(),
    sessionId: getSessionId()
  })
}

// События для отслеживания
trackPaymentEvent('payment_initiated', { amount, currency })
trackPaymentEvent('payment_method_selected', { method })
trackPaymentEvent('3ds_required', { })
trackPaymentEvent('payment_completed', { paymentId, duration })
trackPaymentEvent('payment_failed', { error, stage })
```

### 8.2 Логирование

```typescript
// Структурированное логирование
const paymentLogger = {
  info: (message: string, data?: any) => {
    console.log('[Payment]', message, data)
    if (process.env.NODE_ENV === 'production') {
      Sentry.addBreadcrumb({
        message,
        category: 'payment',
        level: 'info',
        data
      })
    }
  },
  
  error: (message: string, error: Error, data?: any) => {
    console.error('[Payment Error]', message, error, data)
    if (process.env.NODE_ENV === 'production') {
      Sentry.captureException(error, {
        tags: { component: 'payment' },
        extra: data
      })
    }
  }
}
```

## 9. Конфигурация для разных окружений

```typescript
// config/payment.config.ts
export const getPaymentConfig = (env: string) => {
  const configs = {
    development: {
      merchantId: 'test_merchant',
      apiUrl: 'https://sandbox.allsecure.rs/api',
      widgetUrl: 'https://sandbox.allsecure.rs/widget.js',
      webhookSecret: 'test_secret',
      debug: true
    },
    staging: {
      merchantId: process.env.ALLSECURE_MERCHANT_ID,
      apiUrl: 'https://sandbox.allsecure.rs/api',
      widgetUrl: 'https://sandbox.allsecure.rs/widget.js',
      webhookSecret: process.env.ALLSECURE_WEBHOOK_SECRET,
      debug: true
    },
    production: {
      merchantId: process.env.ALLSECURE_MERCHANT_ID,
      apiUrl: 'https://api.allsecure.rs',
      widgetUrl: 'https://widget.allsecure.rs/v1/widget.js',
      webhookSecret: process.env.ALLSECURE_WEBHOOK_SECRET,
      debug: false
    }
  }
  
  return configs[env] || configs.development
}
```