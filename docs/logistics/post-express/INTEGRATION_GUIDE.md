# WSP WebAPI Integration Guide

## Пошаговая инструкция по интеграции

### Шаг 1: Регистрация

1. Перейдите на https://onlinepostexpress.rs/registracija
2. Заполните регистрационную форму:
   - Данные компании (название, ПИБ, матичный номер)
   - Контактные данные (телефон, email)
   - Адрес компании
3. Дождитесь подтверждения регистрации
4. Получите учетные данные:
   - Username (имя пользователя)
   - Password (пароль)

### Шаг 2: Настройка окружения

#### Переменные окружения
```bash
# .env файл
POST_EXPRESS_API_URL=https://onlinepostexpress.rs/WSPWebApi/api/app/transakcija
POST_EXPRESS_USERNAME=your_username
POST_EXPRESS_PASSWORD=your_password
POST_EXPRESS_ENV=production  # или test для тестового окружения
```

#### Базовая конфигурация (Go)
```go
package postexpress

type Config struct {
    APIUrl   string
    Username string
    Password string
    Timeout  time.Duration
}

func NewConfig() *Config {
    return &Config{
        APIUrl:   os.Getenv("POST_EXPRESS_API_URL"),
        Username: os.Getenv("POST_EXPRESS_USERNAME"),
        Password: os.Getenv("POST_EXPRESS_PASSWORD"),
        Timeout:  30 * time.Second,
    }
}
```

### Шаг 3: Создание базового клиента

#### Go реализация
```go
package postexpress

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    config     *Config
    httpClient *http.Client
}

func NewClient(config *Config) *Client {
    return &Client{
        config: config,
        httpClient: &http.Client{
            Timeout: config.Timeout,
        },
    }
}

// Базовый метод для отправки запросов
func (c *Client) SendRequest(transactionID int, data interface{}) (*TransakcijaOut, error) {
    request := TransakcijaIn{
        TransakcijaId:      transactionID,
        DatumVremePosiljke: time.Now().Format(time.RFC3339),
        Klijent: Klijent{
            Username: c.config.Username,
            Password: c.config.Password,
        },
    }
    
    // Добавляем специфические данные
    requestBytes, _ := json.Marshal(request)
    var requestMap map[string]interface{}
    json.Unmarshal(requestBytes, &requestMap)
    
    dataBytes, _ := json.Marshal(data)
    var dataMap map[string]interface{}
    json.Unmarshal(dataBytes, &dataMap)
    
    for k, v := range dataMap {
        requestMap[k] = v
    }
    
    // Отправляем запрос
    jsonData, err := json.Marshal(requestMap)
    if err != nil {
        return nil, fmt.Errorf("marshal error: %w", err)
    }
    
    resp, err := c.httpClient.Post(
        c.config.APIUrl,
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, fmt.Errorf("request error: %w", err)
    }
    defer resp.Body.Close()
    
    var result TransakcijaOut
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("decode error: %w", err)
    }
    
    if !result.OK {
        return nil, fmt.Errorf("API error: %s", result.Poruka)
    }
    
    return &result, nil
}
```

### Шаг 4: Реализация основных методов

#### Поиск населенного пункта
```go
func (c *Client) GetNaselje(naziv string, ptt string) (*NaseljeOut, error) {
    data := map[string]interface{}{
        "NaseljeIn": NaseljeIn{
            Naziv: naziv,
            Ptt:   ptt,
        },
    }
    
    result, err := c.SendRequest(3, data)
    if err != nil {
        return nil, err
    }
    
    // Парсим результат
    var naseljeOut NaseljeOut
    if naseljeData, ok := result.Data["NaseljeOut"]; ok {
        dataBytes, _ := json.Marshal(naseljeData)
        json.Unmarshal(dataBytes, &naseljeOut)
    }
    
    return &naseljeOut, nil
}
```

#### Расчет стоимости доставки
```go
func (c *Client) CalculatePostage(params PostarinaIn) (*PostarinaOut, error) {
    data := map[string]interface{}{
        "PostarinaIn": params,
    }
    
    result, err := c.SendRequest(5, data)
    if err != nil {
        return nil, err
    }
    
    var postarinaOut PostarinaOut
    if postarinaData, ok := result.Data["PostarinaOut"]; ok {
        dataBytes, _ := json.Marshal(postarinaData)
        json.Unmarshal(dataBytes, &postarinaOut)
    }
    
    return &postarinaOut, nil
}
```

#### Отслеживание отправления
```go
func (c *Client) TrackShipment(barcode string, language int) (*TTKretanjeOut, error) {
    data := map[string]interface{}{
        "TTKretanjeIn": TTKretanjeIn{
            BrojPosiljke: barcode,
            Jezik:        language,
        },
    }
    
    result, err := c.SendRequest(63, data)
    if err != nil {
        return nil, err
    }
    
    var kretanjeOut TTKretanjeOut
    if kretanjeData, ok := result.Data["TTKretanjeOut"]; ok {
        dataBytes, _ := json.Marshal(kretanjeData)
        json.Unmarshal(dataBytes, &kretanjeOut)
    }
    
    return &kretanjeOut, nil
}
```

### Шаг 5: Обработка ошибок

```go
type PostExpressError struct {
    Code    string
    Message string
}

func (e *PostExpressError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func parseError(poruka string) *PostExpressError {
    // Парсим сообщения об ошибках
    if strings.Contains(poruka, "корисничко име или лозинка") {
        return &PostExpressError{
            Code:    "AUTH_ERROR",
            Message: "Authentication failed",
        }
    }
    
    if strings.Contains(poruka, "минимум 3 карактера") {
        return &PostExpressError{
            Code:    "VALIDATION_ERROR",
            Message: "Minimum 3 characters required",
        }
    }
    
    return &PostExpressError{
        Code:    "UNKNOWN_ERROR",
        Message: poruka,
    }
}
```

### Шаг 6: Кеширование данных

```go
type Cache struct {
    settlements map[string]*NaseljeOut
    streets     map[string]*UlicaOut
    mu          sync.RWMutex
    ttl         time.Duration
}

func NewCache(ttl time.Duration) *Cache {
    return &Cache{
        settlements: make(map[string]*NaseljeOut),
        streets:     make(map[string]*UlicaOut),
        ttl:         ttl,
    }
}

func (c *Cache) GetSettlement(key string) (*NaseljeOut, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if data, ok := c.settlements[key]; ok {
        return data, true
    }
    return nil, false
}

func (c *Cache) SetSettlement(key string, data *NaseljeOut) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.settlements[key] = data
    
    // Удаляем через TTL
    go func() {
        time.Sleep(c.ttl)
        c.mu.Lock()
        delete(c.settlements, key)
        c.mu.Unlock()
    }()
}
```

### Шаг 7: Интеграция в существующую систему

#### Создание сервиса доставки
```go
package services

type DeliveryService struct {
    client *postexpress.Client
    cache  *postexpress.Cache
}

func NewDeliveryService() *DeliveryService {
    config := postexpress.NewConfig()
    client := postexpress.NewClient(config)
    cache := postexpress.NewCache(1 * time.Hour)
    
    return &DeliveryService{
        client: client,
        cache:  cache,
    }
}

// Метод для расчета доставки в корзине
func (s *DeliveryService) CalculateDeliveryForCart(cart *Cart, address *Address) (*DeliveryInfo, error) {
    // Получаем коды адреса
    settlement, err := s.findSettlement(address.City)
    if err != nil {
        return nil, err
    }
    
    street, err := s.findStreet(settlement.Sifra, address.Street)
    if err != nil {
        return nil, err
    }
    
    // Проверяем доступность услуги
    available, err := s.checkServiceAvailability(
        settlement.Sifra,
        street.Sifra,
        address.HouseNumber,
        58, // Express Today
    )
    if err != nil {
        return nil, err
    }
    
    if !available {
        return nil, fmt.Errorf("service not available for this address")
    }
    
    // Расчитываем стоимость
    weight := s.calculateTotalWeight(cart)
    value := s.calculateTotalValue(cart)
    
    postage, err := s.client.CalculatePostage(postexpress.PostarinaIn{
        SifraUsluge:        58,
        Masa:               weight,
        Vrednost:           value,
        SifraMestaPrimaoca: settlement.Sifra,
        BrojPosiljaka:      1,
    })
    
    if err != nil {
        return nil, err
    }
    
    return &DeliveryInfo{
        ServiceName:  "Post Express - Express Today",
        Price:        postage.Postarina.UkupnaCena,
        DeliveryTime: "Same day delivery",
        Available:    true,
    }, nil
}
```

### Шаг 8: Webhook для обновления статусов

```go
// Endpoint для приема webhook от Post Express
func (h *Handler) HandlePostExpressWebhook(c *fiber.Ctx) error {
    var webhook PostExpressWebhook
    if err := c.BodyParser(&webhook); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid webhook data",
        })
    }
    
    // Проверяем подпись
    if !h.verifyWebhookSignature(c) {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid signature",
        })
    }
    
    // Обновляем статус заказа
    order, err := h.orderService.GetByTrackingNumber(webhook.BrojPosiljke)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "Order not found",
        })
    }
    
    // Обновляем статус
    order.DeliveryStatus = webhook.Status
    order.LastStatusUpdate = time.Now()
    
    if err := h.orderService.Update(order); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to update order",
        })
    }
    
    // Отправляем уведомление пользователю
    h.notificationService.SendDeliveryUpdate(order)
    
    return c.JSON(fiber.Map{
        "success": true,
    })
}
```

### Шаг 9: Тестирование

#### Unit тесты
```go
func TestCalculatePostage(t *testing.T) {
    client := NewMockClient()
    
    result, err := client.CalculatePostage(PostarinaIn{
        SifraUsluge: 58,
        Masa:        1000,
        Vrednost:    5000,
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Greater(t, result.Postarina.UkupnaCena, 0.0)
}
```

#### Интеграционные тесты
```go
func TestRealAPIIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    config := NewConfig()
    client := NewClient(config)
    
    // Тестируем поиск населенного пункта
    naselja, err := client.GetNaselje("Београд", "")
    assert.NoError(t, err)
    assert.NotEmpty(t, naselja.Naselja)
    
    // Тестируем расчет стоимости
    postage, err := client.CalculatePostage(PostarinaIn{
        SifraUsluge:        58,
        Masa:               1000,
        SifraMestaPrimaoca: naselja.Naselja[0].Sifra,
    })
    assert.NoError(t, err)
    assert.Greater(t, postage.Postarina.UkupnaCena, 0.0)
}
```

### Шаг 10: Мониторинг и логирование

```go
// Middleware для логирования API вызовов
func (c *Client) logAPICall(transactionID int, duration time.Duration, err error) {
    status := "success"
    if err != nil {
        status = "error"
    }
    
    log.Printf(
        "PostExpress API Call: TransactionID=%d, Duration=%v, Status=%s, Error=%v",
        transactionID,
        duration,
        status,
        err,
    )
    
    // Отправляем метрики
    metrics.PostExpressAPICallDuration.Observe(duration.Seconds())
    metrics.PostExpressAPICallTotal.WithLabelValues(
        fmt.Sprintf("%d", transactionID),
        status,
    ).Inc()
}
```

## Best Practices

### 1. Безопасность
- Никогда не храните пароли в коде
- Используйте переменные окружения
- Регулярно обновляйте пароли
- Логируйте все операции

### 2. Производительность
- Кешируйте данные о населенных пунктах и улицах
- Используйте connection pooling
- Реализуйте retry логику
- Установите разумные таймауты

### 3. Обработка ошибок
- Всегда проверяйте поле OK в ответе
- Логируйте все ошибки
- Предоставляйте понятные сообщения пользователям
- Реализуйте fallback механизмы

### 4. Тестирование
- Используйте тестовое окружение для разработки
- Пишите unit и интеграционные тесты
- Тестируйте граничные случаи
- Мокируйте внешние зависимости

## Частые проблемы и решения

### Проблема: "Неисправно корисничко име или лозинка"
**Решение**: Проверьте правильность учетных данных

### Проблема: "Сервис тренутно недоступан"
**Решение**: Реализуйте retry логику с экспоненциальной задержкой

### Проблема: Медленные ответы API
**Решение**: Используйте кеширование для справочных данных

### Проблема: Некорректные коды адресов
**Решение**: Всегда используйте методы GetNaselje и GetUlica для получения актуальных кодов

## Контакты поддержки

- **Техническая поддержка**: wsp@posta.rs
- **Телефон**: 0700 100 300
- **Рабочее время**: Пн-Пт 08:00-16:00