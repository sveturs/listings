# Задача: Автоматическое создание настроек приватности при регистрации пользователей

## Дата: 16.06.2025

## Проблема
При добавлении контактов возникала ошибка 500 из-за отсутствия записей в таблице `user_privacy_settings` для новых пользователей.

## Решение
1. Обновлен метод `CreateUser` в файле `/backend/internal/proj/users/storage/postgres/user.go`:
   - Добавлено автоматическое создание настроек приватности после создания пользователя через email/password

2. Обновлен метод `GetOrCreateGoogleUser` в том же файле:
   - Добавлено автоматическое создание настроек приватности для пользователей, регистрирующихся через Google OAuth

3. Созданы настройки приватности для всех существующих пользователей через SQL:
   ```sql
   INSERT INTO user_privacy_settings (user_id) 
   SELECT id FROM users 
   WHERE id NOT IN (SELECT user_id FROM user_privacy_settings);
   ```

## Изменения в коде

### backend/internal/proj/users/storage/postgres/user.go

В методе `CreateUser` после создания пользователя добавлено:
```go
// Создаем настройки приватности для нового пользователя
_, err = s.pool.Exec(ctx, `
    INSERT INTO user_privacy_settings (user_id)
    VALUES ($1)
    ON CONFLICT (user_id) DO NOTHING
`, user.ID)
if err != nil {
    s.logger.Info("Failed to create privacy settings for user %d: %v", user.ID, err)
    // Не прерываем процесс, так как пользователь уже создан
}
```

Аналогичный код добавлен в метод `GetOrCreateGoogleUser` для OAuth пользователей.

## Результат
- Новые пользователи автоматически получают настройки приватности при регистрации
- Существующие пользователи получили настройки через миграцию данных
- Функция добавления контактов теперь работает корректно