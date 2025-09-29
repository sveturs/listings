# 🔐 Инструкция по генерации валидных JWT токенов для тестирования

## Архитектура авторизации

В проекте используется микросервисная архитектура с отдельным Auth Service:

- **Auth Service**: https://authpreprod.svetu.rs (RS256 алгоритм)
- **Backend API**: http://localhost:3000 (проверяет токены через публичный ключ)
- **Frontend**: http://localhost:3001

## ⚠️ ВАЖНО: Почему локальный скрипт не работает

Скрипт `backend/scripts/create_test_jwt.go` создает токены с алгоритмом HS256 (HMAC), но система использует RS256 (RSA). Это разные алгоритмы:

- **HS256**: Симметричный, использует общий секрет (JWT_SECRET)
- **RS256**: Асимметричный, использует пару приватный/публичный ключ

Backend проверяет токены используя публичный ключ из `/data/hostel-booking-system/backend/keys/auth_service_public.pem`.

## ✅ Правильные способы получения валидного токена

### Способ 1: Использование Auth Service напрямую (Рекомендуется)

```bash
# 1. Подключиться к серверу с Auth Service
ssh svetu@svetu.rs

# 2. Перейти в директорию Auth Service
cd /opt/svetu-authpreprod

# 3. Создать токен используя существующий скрипт
go run scripts/create_admin_jwt.go

# Токен будет выведен в консоль
```

### Способ 2: Создание локального скрипта с правильным алгоритмом

Создайте файл `backend/scripts/create_rs256_jwt.go`:

```go
package main

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "io/ioutil"
    "log"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID  int      `json:"user_id"`
    Email   string   `json:"email"`
    Name    string   `json:"name"`
    Roles   []string `json:"roles"`
    IsAdmin bool     `json:"is_admin"`
    jwt.RegisteredClaims
}

func main() {
    // Читаем приватный ключ (нужно скопировать с сервера)
    privateKeyData, err := ioutil.ReadFile("/data/hostel-booking-system/backend/keys/auth_service_private.pem")
    if err != nil {
        log.Fatalf("Failed to read private key: %v", err)
    }

    block, _ := pem.Decode(privateKeyData)
    if block == nil {
        log.Fatal("Failed to decode PEM block")
    }

    privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        log.Fatalf("Failed to parse private key: %v", err)
    }

    rsaKey, ok := privateKey.(*rsa.PrivateKey)
    if !ok {
        log.Fatal("Not an RSA private key")
    }

    // Создаем claims для тестового пользователя
    now := time.Now()
    claims := Claims{
        UserID:  1,
        Email:   "test@example.com",
        Name:    "Test User",
        Roles:   []string{"user"},
        IsAdmin: false,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "https://auth.svetu.rs",
            Subject:   "1",
            Audience:  []string{"https://svetu.rs"},
            ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
        },
    }

    // Для админа
    // claims.UserID = 5
    // claims.Email = "admin@example.com"
    // claims.Roles = []string{"admin"}
    // claims.IsAdmin = true

    // Создаем токен
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

    // Подписываем приватным ключом
    tokenString, err := token.SignedString(rsaKey)
    if err != nil {
        log.Fatalf("Failed to sign token: %v", err)
    }

    fmt.Println(tokenString)
}
```

**Важно**: Для работы этого скрипта нужно:
1. Скопировать приватный ключ с сервера: `scp svetu@svetu.rs:/opt/svetu-authpreprod/keys/private.pem backend/keys/auth_service_private.pem`
2. **НЕ коммитить приватный ключ в репозиторий!**
3. Добавить `backend/keys/auth_service_private.pem` в `.gitignore`

### Способ 3: Получение токена через API Auth Service

```bash
# 1. Регистрация нового пользователя
curl -X POST https://authpreprod.svetu.rs/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPassword123!",
    "name": "Test User"
  }'

# 2. Вход и получение токена
curl -X POST https://authpreprod.svetu.rs/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPassword123!"
  }'

# Ответ будет содержать access_token и refresh_token
```

### Способ 4: Быстрый токен для тестирования через SSH (РЕКОМЕНДУЕТСЯ)

```bash
# Получить токен с исправленным путем к ключу
TOKEN=$(ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && sed 's|/data/auth_svetu/keys/private.pem|./keys/private.pem|g' scripts/create_admin_jwt.go > /tmp/create_jwt_fixed.go && go run /tmp/create_jwt_fixed.go 2>/dev/null")

# Проверить что токен получен
echo "Token obtained: ${TOKEN:0:50}..."

# Теперь можно использовать в запросах
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me | jq '.'
```

**Примечание**: В скрипте `create_admin_jwt.go` неверный путь к ключу, поэтому мы его исправляем на лету.

## 🧪 Тестирование токена

### Проверка валидности токена

```bash
# Установить токен в переменную
TOKEN="ваш_токен_здесь"

# Проверить профиль пользователя
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/users/me | jq '.'

# Проверить доступ к защищенным эндпоинтам
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/recommendations/user | jq '.'

# Для админских эндпоинтов нужен токен с is_admin=true
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/admin/users | jq '.'
```

### Декодирование токена для проверки

```bash
# Декодировать payload токена (без проверки подписи)
echo "$TOKEN" | cut -d. -f2 | base64 -d 2>/dev/null | jq '.'
```

## 📋 Чек-лист для успешного тестирования

1. ✅ Используйте RS256 алгоритм, а не HS256
2. ✅ Убедитесь, что публичный ключ в backend актуален
3. ✅ Проверьте issuer и audience в токене (должны соответствовать настройкам)
4. ✅ Проверьте время жизни токена (exp claim)
5. ✅ Для админских операций установите `is_admin: true` и `roles: ["admin"]`

## 🔧 Устранение проблем

### Ошибка "invalid_token"

Возможные причины:
- Использован неправильный алгоритм (HS256 вместо RS256)
- Истек срок действия токена
- Неправильный issuer или audience
- Несоответствие публичного ключа

### Ошибка "unauthorized"

Возможные причины:
- Отсутствует заголовок Authorization
- Неправильный формат: должен быть `Bearer TOKEN`
- Токен не содержит нужных прав (например, is_admin для админских эндпоинтов)

## 📝 Пример полного процесса тестирования

```bash
# 1. Получаем валидный токен с Auth Service
TOKEN=$(ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && go run scripts/create_admin_jwt.go 2>/dev/null")

# 2. Проверяем что токен работает
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me | jq '.data.email'

# 3. Тестируем нужный эндпоинт
curl -H "Authorization: Bearer $TOKEN" \
  -X GET "http://localhost:3000/api/v1/marketplace/recommendations/user?limit=10" | jq '.'

# 4. Для POST запросов
curl -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -X POST "http://localhost:3000/api/v1/marketplace/favorites" \
  -d '{"listing_id": 123}' | jq '.'
```

## 🚀 Автоматизация для команды

Создайте алиас в `.bashrc` или `.zshrc`:

```bash
alias get-test-token='ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && go run scripts/create_admin_jwt.go 2>/dev/null"'
alias test-with-token='TOKEN=$(get-test-token) && echo "Token obtained. Use \$TOKEN in your commands"'
```

После этого можно просто:
```bash
test-with-token
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me
```

## 📌 Важные замечания

1. **Безопасность**: Никогда не коммитьте приватные ключи в репозиторий
2. **Среды**: Продакшн использует другие ключи и URL (https://auth.svetu.rs)
3. **Кэширование**: Backend может кэшировать проверку токенов - при изменении ключей перезапустите сервер
4. **Логи**: При проблемах проверяйте логи backend: `tail -f /tmp/backend.log`

---

*Документ создан: 27.09.2025*
*Автор: Backend Team*