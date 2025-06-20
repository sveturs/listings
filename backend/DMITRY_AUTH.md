# Авторизация пользователя Dmitry Voroshilov

## Данные пользователя
- **Email**: voroshilovdo@gmail.com
- **ID**: 14
- **Имя**: Dmitry Voroshilov
- **Google ID**: 102440686443518161778
- **Роль**: Администратор (is_admin: true)

## Способы авторизации

### 1. JWT Token (Рекомендуется для API)

#### Access Token (действителен 24 часа):
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNCwiZW1haWwiOiJ2b3Jvc2hpbG92ZG9AZ21haWwuY29tIiwiZXhwIjoxNzUwNDQzNTYxLCJuYmYiOjE3NTAzNTcxNjEsImlhdCI6MTc1MDM1NzE2MX0.GxVb3tziQOxfhbmiZ1OsIjUZdptEAVsBf-21BVNlyYE
```

#### Refresh Token (действителен 30 дней):
```
rRv_z02HCj6dcOpDrBr28mM4cCPsNB44KYxJ7Wq1T_0=
```

#### Использование:

**В cURL:**
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNCwiZW1haWwiOiJ2b3Jvc2hpbG92ZG9AZ21haWwuY29tIiwiZXhwIjoxNzUwNDQzNTYxLCJuYmYiOjE3NTAzNTcxNjEsImlhdCI6MTc1MDM1NzE2MX0.GxVb3tziQOxfhbmiZ1OsIjUZdptEAVsBf-21BVNlyYE" \
     http://localhost:3000/api/v1/users/profile
```

**В JavaScript (axios):**
```javascript
const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNCwiZW1haWwiOiJ2b3Jvc2hpbG92ZG9AZ21haWwuY29tIiwiZXhwIjoxNzUwNDQzNTYxLCJuYmYiOjE3NTAzNTcxNjEsImlhdCI6MTc1MDM1NzE2MX0.GxVb3tziQOxfhbmiZ1OsIjUZdptEAVsBf-21BVNlyYE';
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
```

**В JavaScript (fetch):**
```javascript
const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNCwiZW1haWwiOiJ2b3Jvc2hpbG92ZG9AZ21haWwuY29tIiwiZXhwIjoxNzUwNDQzNTYxLCJuYmYiOjE3NTAzNTcxNjEsImlhdCI6MTc1MDM1NzE2MX0.GxVb3tziQOxfhbmiZ1OsIjUZdptEAVsBf-21BVNlyYE';
fetch('http://localhost:3000/api/v1/users/profile', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

**В Postman:**
- Authorization type: Bearer Token
- Token: вставьте токен выше

### 2. Session Cookie (для браузера)

Сессионные токены хранятся в памяти сервера. Для создания настоящей сессии нужно:

1. Войти через Google OAuth:
   - Откройте http://localhost:3001
   - Нажмите "Войти через Google"
   - Выберите аккаунт voroshilovdo@gmail.com

2. Или используйте существующий JWT токен в localStorage:
   ```javascript
   localStorage.setItem('auth_token', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNCwiZW1haWwiOiJ2b3Jvc2hpbG92ZG9AZ21haWwuY29tIiwiZXhwIjoxNzUwNDQzNTYxLCJuYmYiOjE3NTAzNTcxNjEsImlhdCI6MTc1MDM1NzE2MX0.GxVb3tziQOxfhbmiZ1OsIjUZdptEAVsBf-21BVNlyYE');
   ```

### 3. Обновление токена

Когда access token истечет, используйте refresh token:

```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "rRv_z02HCj6dcOpDrBr28mM4cCPsNB44KYxJ7Wq1T_0="}'
```

### 4. Проверка авторизации

```bash
# Получить профиль пользователя
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNCwiZW1haWwiOiJ2b3Jvc2hpbG92ZG9AZ21haWwuY29tIiwiZXhwIjoxNzUwNDQzNTYxLCJuYmYiOjE3NTwM1NzE2MSwiaWF0IjoxNzUwMzU3MTYxfQ.GxVb3tziQOxfhbmiZ1OsIjUZdptEAVsBf-21BVNlyYE" \
     http://localhost:3000/api/v1/users/profile

# Проверить админ доступ
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNCwiZW1haWwiOiJ2b3Jvc2hpbG92ZG9AZ21haWwuY29tIiwiZXhwIjoxNzUwNDQzNTYxLCJuYmYiOjE3NTAzNTcxNjEsImlhdCI6MTc1MDM1NzE2MX0.GxVb3tziQOxfhbmiZ1OsIjUZdptEAVsBf-21BVNlyYE" \
     http://localhost:3000/api/v1/admin/users
```

## Утилиты для генерации новых токенов

Если токены истекли, используйте следующие утилиты:

### Генерация JWT токена:
```bash
cd /data/hostel-booking-system/backend
go run cmd/utils/generate_dmitry_jwt/simple.go
```

### Генерация с использованием прямого SQL:
```bash
cd /data/hostel-booking-system/backend
go run cmd/utils/generate_dmitry_token/direct_sql.go
```

## Важные замечания

1. JWT токены подписаны секретом `yoursecretkey` (из файла .env)
2. Пользователь имеет права администратора
3. Токены сохранены в таблице `refresh_tokens` в БД
4. Access token истекает через 24 часа
5. Refresh token истекает через 30 дней