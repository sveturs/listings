# Паспорт User Models (Backend)

## User (основная модель)

### Go структура
```go
type User struct {
    ID         int       `json:"id" db:"id"`
    Name       string    `json:"name" db:"name"`
    Email      string    `json:"email" db:"email"`
    GoogleID   *string   `json:"google_id,omitempty" db:"google_id"`
    PictureURL *string   `json:"picture_url,omitempty" db:"picture_url"`
    Phone      *string   `json:"phone,omitempty" db:"phone"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
```

### Соответствие с другими слоями
- **PostgreSQL**: таблица `users`
- **Frontend**: тип `components['schemas']['models.User']`
- **OpenSearch**: не индексируется отдельно, только как часть listings

### Особенности
- **Email**: уникальный, обязательный
- **GoogleID**: nullable, используется для OAuth
- **Phone**: опциональный, для контактов

## UserBalance (баланс пользователя)

### Go структура
```go
type UserBalance struct {
    ID        int            `json:"id" db:"id"`
    UserID    int            `json:"user_id" db:"user_id"`
    Balance   decimal.Decimal `json:"balance" db:"balance"`
    Currency  string         `json:"currency" db:"currency"`
    CreatedAt time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}
```

### Особенности
- Использует `decimal.Decimal` для точности финансовых операций
- Поддержка мультивалютности через поле `currency`

## UserContact (контакты пользователя)

### Go структура
```go
type UserContact struct {
    ID          int64     `json:"id" db:"id"`
    UserID      int       `json:"user_id" db:"user_id"`
    ContactType string    `json:"contact_type" db:"contact_type"`
    ContactValue string   `json:"contact_value" db:"contact_value"`
    IsPublic    bool      `json:"is_public" db:"is_public"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### Типы контактов
- `phone` - телефон
- `email` - email (дополнительный)
- `telegram` - Telegram username
- `viber` - Viber
- `whatsapp` - WhatsApp

## UserTransaction (транзакции)

### Go структура
```go
type UserTransaction struct {
    ID            int64           `json:"id" db:"id"`
    UserID        int             `json:"user_id" db:"user_id"`
    Amount        decimal.Decimal `json:"amount" db:"amount"`
    Currency      string          `json:"currency" db:"currency"`
    Type          string          `json:"type" db:"type"`
    Status        string          `json:"status" db:"status"`
    Description   string          `json:"description" db:"description"`
    ReferenceID   *string         `json:"reference_id,omitempty" db:"reference_id"`
    CreatedAt     time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}
```

### Типы транзакций
- `deposit` - пополнение
- `withdrawal` - снятие
- `purchase` - покупка
- `refund` - возврат
- `commission` - комиссия

### Статусы
- `pending` - в обработке
- `completed` - завершена
- `failed` - неудачна
- `cancelled` - отменена

## API Endpoints для User

### Authentication
- `POST /auth/google` - Google OAuth
- `GET /auth/session` - получить текущую сессию
- `POST /auth/logout` - выход

### User Management
- `GET /users/profile` - профиль пользователя
- `PUT /users/profile` - обновить профиль
- `GET /users/contacts` - контакты пользователя
- `POST /users/contacts` - добавить контакт
- `PUT /users/contacts/{id}` - обновить контакт
- `DELETE /users/contacts/{id}` - удалить контакт

### Balance & Transactions
- `GET /users/balance` - баланс пользователя
- `GET /users/transactions` - история транзакций
- `POST /users/transactions` - создать транзакцию

## Валидация

### User
- `name`: минимум 2 символа, максимум 100
- `email`: валидный email формат
- `phone`: E164 формат (опционально)

### UserContact
- `contact_type`: один из допустимых типов
- `contact_value`: соответствует типу контакта
- `is_public`: boolean

### UserTransaction  
- `amount`: положительное число
- `currency`: поддерживаемая валюта
- `type`: один из допустимых типов

## Известные особенности

1. **ID типы**: User использует `int`, UserContact использует `int64`
2. **Decimal поля**: Финансовые суммы используют `decimal.Decimal`
3. **Nullable поля**: Google-специфичные поля опциональны
4. **JSON omitempty**: Применяется для опциональных полей