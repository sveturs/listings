# Замена системы платежей Stripe на сербское решение

## Текущая архитектура (Stripe)

### Компоненты системы:
1. **Backend**: `/backend/internal/proj/payments/service/stripe.go` - интеграция со Stripe
2. **Модели**: 
   - `UserBalance` - баланс пользователя  
   - `BalanceTransaction` - транзакции пополнения/списания
   - `PaymentMethod` - методы оплаты
   - `PaymentSession` - сессии оплаты
3. **База данных**:
   - `user_balances` - балансы пользователей
   - `balance_transactions` - история транзакций
   - `payment_methods` - доступные методы оплаты
4. **Frontend**: `/frontend/hostel-frontend/src/components/balance/DepositDialog.tsx`

### Проблемы текущей системы:
- ❌ Stripe не поддерживает Сербию
- ❌ Система пополнения баланса - хранение денег пользователей
- ❌ Сложность с налогообложением предоплаченных средств

## Рекомендуемое решение для Сербии

### 1. Интеграция с банком "Поштанска штедионица"

**API банка:**
- **NBS IPS** (Национальная платежная система) для мгновенных переводов
- **eCommerce Gateway** - интеграция для интернет-платежей
- **QR код платежи** через IPS систему

### 2. Новая архитектура "Pay-per-service"

Вместо системы баланса прямые платежи:

```
к примеру услуга стоит 30 динар → Прямой платеж 30 динар → Услуга активируется
```

### 3. Интеграция с сербскими платежными системами

#### A. Поштанска штедионица eCommerce
```go
type PostanskaPaymentService struct {
    merchantID     string
    merchantSecret string
    apiEndpoint    string
    returnURL      string
}

func (p *PostanskaPaymentService) CreateDirectPayment(
    orderID string, 
    amount float64, 
    description string
) (*DirectPaymentSession, error) {
    // Создание прямого платежа без баланса
}
```

#### B. IPS QR код платежи
```go
type IPSQRService struct {
    merchantAccount string
    bankCode       string
}

func (i *IPSQRService) GenerateQRPayment(
    amount float64,
    reference string,
    description string
) (*QRCodePayment, error) {
    // Генерация QR кода для мгновенного платежа
}
```

#### C. Комфи платежи (Комерцијална банка)
```go
type ComfyPaymentService struct {
    terminalID string
    secretKey  string
}
```

### 4. Предлагаемая структура БД

```sql
-- Заменяем balance_transactions на direct_payments
CREATE TABLE direct_payments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    service_type VARCHAR(50) NOT NULL, -- 'listing_boost', 'translation', 'premium_feature'
    service_id INT, -- ID конкретной услуги
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RSD',
    payment_method VARCHAR(50) NOT NULL, -- 'postanska_card', 'ips_qr', 'comfy_pay'
    external_transaction_id VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending', -- pending, completed, failed, refunded
    payment_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Прайс-лист услуг
CREATE TABLE service_prices (
    id SERIAL PRIMARY KEY,
    service_type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RSD',
    is_active BOOLEAN DEFAULT TRUE
);

-- Поддерживаемые платежные системы
CREATE TABLE payment_providers (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'bank_card', 'qr_code', 'bank_transfer'
    config JSONB,
    is_active BOOLEAN DEFAULT TRUE
);
```

### 5. Примеры интеграции

#### Поштанска штедионица API:
```go
func (p *PostanskaPaymentService) CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
    payload := map[string]interface{}{
        "merchant_id": p.merchantID,
        "order_id":    req.OrderID,
        "amount":      fmt.Sprintf("%.2f", req.Amount),
        "currency":    "RSD",
        "return_url":  p.returnURL + "/success",
        "cancel_url":  p.returnURL + "/cancel",
        "description": req.Description,
        "timestamp":   time.Now().Unix(),
    }
    
    // Подпись запроса
    signature := p.generateSignature(payload)
    payload["signature"] = signature
    
    // HTTP запрос к банку
    resp, err := http.Post(p.apiEndpoint+"/create-payment", "application/json", bytes.NewBuffer(jsonPayload))
    // ...
}
```

#### IPS QR код:
```go
func (i *IPSQRService) GenerateQRCode(amount float64, reference string) (string, error) {
    // Формат IPS QR кода
    qrData := fmt.Sprintf("K:PR|V:01|C:1|R:%s|N:%s|I:RSD%.2f|SF:289|S:Sve Tu platforma|P:%s", 
        reference, i.merchantAccount, amount, reference)
    
    // Генерация QR кода
    qrCode, err := qr.Encode(qrData, qr.M, qr.Auto)
    return base64.StdEncoding.EncodeToString(qrCode.PNG()), err
}
```

### 6. Фронтенд изменения

```tsx
// Новый компонент прямого платежа
const DirectPaymentDialog: React.FC<{
  serviceType: string;
  amount: number;
  description: string;
  onSuccess: () => void;
}> = ({ serviceType, amount, description, onSuccess }) => {
  const [paymentMethod, setPaymentMethod] = useState('');
  
  const handlePayment = async () => {
    const response = await axios.post('/api/v1/payments/direct', {
      service_type: serviceType,
      amount: amount,
      payment_method: paymentMethod,
      description: description
    });
    
    // Перенаправление на платежную страницу банка
    window.location.href = response.data.payment_url;
  };
  
  return (
    // UI для выбора способа оплаты и прямого платежа
  );
};
```

### 7. План миграции

#### Этап 1: Подготовка
1. Получить договор с Поштанской штедионицей на eCommerce услуги
2. Настроить тестовую среду
3. Реализовать новые модели данных

#### Этап 2: Реализация
1. Создать сервисы для прямых платежей
2. Интегрировать с банковскими API
3. Обновить фронтенд на прямые платежи

#### Этап 3: Тестирование
1. Тестирование в sandbox среде банка
2. Тестирование QR код платежей
3. Проверка webhook'ов

#### Этап 4: Переход
1. Отключить Stripe
2. Включить новую систему
3. Мониторинг платежей

### 8. Преимущества нового решения

✅ **Соответствие законодательству Сербии**
✅ **Отсутствие предоплаченных средств** - упрощение налогообложения  
✅ **Прямые платежи** - пользователь платит только за то, что использует
✅ **IPS интеграция** - мгновенные переводы между банками Сербии
✅ **QR код платежи** - современный способ оплаты
✅ **Поддержка местных банков** - лучше для пользователей

### 9. Ориентировочные тарифы

- **Поштанска штедионица**: 1.5-2% + фиксированная комиссия
- **IPS переводы**: 0.5-1% 
- **QR код платежи**: обычно без комиссии для получателя

### 10. Необходимые документы

1. Договор с банком на интернет-платежи
2. Сертификат мерчанта
3. API ключи для интеграции
4. Настройка webhook URL для уведомлений о платежах