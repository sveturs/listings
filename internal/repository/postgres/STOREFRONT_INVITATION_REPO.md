# Storefront Invitation Repository

Repository layer для работы с системой приглашений в витрины (storefronts).

## Обзор

`StorefrontInvitationRepository` предоставляет методы для работы с приглашениями сотрудников в витрины. Поддерживает два типа приглашений:
- **Email invitations** - персональные приглашения на email (одноразовые)
- **Link invitations** - приглашения по ссылке (многоразовые, с лимитом использований)

## Создание repository

```go
import (
    "database/sql"
    "github.com/vondi-global/listings/internal/repository/postgres"
)

db, err := sql.Open("postgres", dsn)
if err != nil {
    // handle error
}

repo := postgres.NewStorefrontInvitationRepository(db)
```

## Основные методы

### Create - Создание приглашения

```go
import "github.com/vondi-global/listings/internal/domain"

// Email приглашение
email := "user@example.com"
inv := &domain.StorefrontInvitation{
    StorefrontID: 1,
    Role:         "staff",
    Type:         domain.InvitationTypeEmail,
    InvitedEmail: &email,
    InvitedByID:  100,
    Status:       domain.InvitationStatusPending,
}

err := repo.Create(ctx, inv)
// inv.ID, inv.CreatedAt, inv.UpdatedAt заполнены после успеха

// Link приглашение
code := "sf-abc123"
expires := time.Now().Add(7 * 24 * time.Hour)
maxUses := int32(10)

invLink := &domain.StorefrontInvitation{
    StorefrontID: 1,
    Role:         "manager",
    Type:         domain.InvitationTypeLink,
    InviteCode:   &code,
    ExpiresAt:    &expires,
    MaxUses:      &maxUses,
    InvitedByID:  100,
    Status:       domain.InvitationStatusPending,
}

err := repo.Create(ctx, invLink)
```

### GetByID - Получение по ID

```go
inv, err := repo.GetByID(ctx, 1)
if err != nil {
    // "invitation not found" если не найдено
}
```

### GetByCode - Получение по коду приглашения

```go
inv, err := repo.GetByCode(ctx, "sf-abc123")
if err != nil {
    // handle error
}
```

### GetByEmailAndStorefront - Проверка дубликатов

Проверяет есть ли уже pending приглашение для данного email в этой витрине:

```go
inv, err := repo.GetByEmailAndStorefront(ctx, "user@example.com", 1)
if err != nil {
    // handle error
}
if inv != nil {
    // Уже существует pending приглашение
}
```

### ListByStorefront - Список приглашений витрины

```go
filter := &domain.ListInvitationsFilter{
    StorefrontID: &storefrontID,
    Status:       &statusPending,  // optional
    Type:         &typeEmail,      // optional
    Page:         1,
    Limit:        10,
}

invitations, total, err := repo.ListByStorefront(ctx, storefrontID, filter)
// invitations - массив приглашений
// total - общее количество (для пагинации)
```

### ListByUser - Приглашения пользователя

Возвращает все pending приглашения для пользователя (по userID или email):

```go
invitations, err := repo.ListByUser(ctx, userID, "user@example.com")
// Автоматически фильтрует только pending и не истёкшие
```

### Update - Обновление приглашения

```go
inv.Status = domain.InvitationStatusAccepted
now := time.Now()
inv.AcceptedAt = &now

err := repo.Update(ctx, inv)
// inv.UpdatedAt обновится
```

### Delete - Удаление приглашения

```go
err := repo.Delete(ctx, invitationID)
// Hard delete из БД
```

### IncrementUses - Инкремент использований

Атомарно увеличивает счётчик использований для link приглашений:

```go
err := repo.IncrementUses(ctx, invitationID)
if err != nil {
    // "max uses reached" если лимит достигнут
}
```

### MarkAsExpired - Пометить как истёкшее

```go
err := repo.MarkAsExpired(ctx, invitationID)
// Устанавливает status = 'expired'
```

### ExpirePendingInvitations - Массовое истечение

Помечает все pending приглашения с прошедшим expires_at как expired:

```go
affected, err := repo.ExpirePendingInvitations(ctx)
// affected - количество обновлённых записей
```

Полезно для cron job:

```go
// Запускать раз в час
func expireInvitationsJob(ctx context.Context) {
    affected, err := repo.ExpirePendingInvitations(ctx)
    if err != nil {
        log.Error().Err(err).Msg("failed to expire invitations")
        return
    }
    log.Info().Int64("affected", affected).Msg("expired invitations")
}
```

### MarkAsRevoked - Отзыв приглашения

```go
err := repo.MarkAsRevoked(ctx, invitationID)
// Устанавливает status = 'revoked'
```

## Вспомогательные методы

### GetStatsByStorefront - Статистика по статусам

```go
stats, err := repo.GetStatsByStorefront(ctx, storefrontID)
// map[string]int32{
//   "pending": 5,
//   "accepted": 10,
//   "declined": 2,
//   "expired": 3
// }
```

### CheckInviteCodeExists - Проверка существования кода

```go
exists, err := repo.CheckInviteCodeExists(ctx, "sf-abc123")
if exists {
    // Код уже используется, нужно сгенерировать новый
}
```

### GetActiveLinkInvitations - Активные link приглашения

Возвращает все активные link приглашения для витрины:
- status = 'pending'
- не истёкшие
- не достигшие лимита использований

```go
invitations, err := repo.GetActiveLinkInvitations(ctx, storefrontID)
// Полезно для отображения активных ссылок-приглашений
```

## Типичные use cases

### 1. Создание email приглашения с проверкой дубликатов

```go
// Проверить нет ли уже pending приглашения
existing, err := repo.GetByEmailAndStorefront(ctx, email, storefrontID)
if err != nil {
    return err
}
if existing != nil {
    return fmt.Errorf("invitation already exists")
}

// Создать новое приглашение
inv := &domain.StorefrontInvitation{
    StorefrontID: storefrontID,
    Role:         role,
    Type:         domain.InvitationTypeEmail,
    InvitedEmail: &email,
    InvitedByID:  inviterID,
    Status:       domain.InvitationStatusPending,
}

if err := repo.Create(ctx, inv); err != nil {
    return err
}

// Отправить email с кодом приглашения
sendInvitationEmail(email, inv.ID)
```

### 2. Создание link приглашения с уникальным кодом

```go
import "github.com/vondi-global/listings/internal/domain"

// Генерировать код пока не найдём свободный
var code string
for {
    var err error
    code, err = domain.GenerateInviteCode()
    if err != nil {
        return err
    }

    exists, err := repo.CheckInviteCodeExists(ctx, code)
    if err != nil {
        return err
    }
    if !exists {
        break
    }
}

// Создать приглашение
expires := time.Now().Add(7 * 24 * time.Hour)
maxUses := int32(10)

inv := &domain.StorefrontInvitation{
    StorefrontID: storefrontID,
    Role:         role,
    Type:         domain.InvitationTypeLink,
    InviteCode:   &code,
    ExpiresAt:    &expires,
    MaxUses:      &maxUses,
    InvitedByID:  inviterID,
    Status:       domain.InvitationStatusPending,
}

if err := repo.Create(ctx, inv); err != nil {
    return err
}

// Вернуть ссылку вида https://vondi.rs/invite/sf-abc123
```

### 3. Принятие приглашения по коду

```go
// Получить приглашение
inv, err := repo.GetByCode(ctx, code)
if err != nil {
    return err
}

// Проверить валидность через domain методы
if !inv.CanAccept() {
    return fmt.Errorf("invitation cannot be accepted")
}

// Обновить статус
inv.MarkAsAccepted() // Устанавливает Status и AcceptedAt

if err := repo.Update(ctx, inv); err != nil {
    return err
}

// Если link invitation - инкрементировать счётчик
if inv.IsLinkInvitation() {
    if err := repo.IncrementUses(ctx, inv.ID); err != nil {
        return err
    }
}

// Добавить пользователя в storefront_staff
// (это делается в service layer)
```

### 4. Отображение списка приглашений с фильтрацией

```go
// Фильтр для админ панели
status := domain.InvitationStatusPending
filter := &domain.ListInvitationsFilter{
    Status: &status,
    Page:   1,
    Limit:  20,
}

invitations, total, err := repo.ListByStorefront(ctx, storefrontID, filter)
if err != nil {
    return err
}

// Рендер списка с пагинацией
renderInvitations(invitations, total)
```

### 5. Автоматическое истечение приглашений (cron)

```go
// Запускать каждый час
func expireInvitationsCron(repo *StorefrontInvitationRepository) {
    ctx := context.Background()

    affected, err := repo.ExpirePendingInvitations(ctx)
    if err != nil {
        log.Error().Err(err).Msg("failed to expire invitations")
        return
    }

    if affected > 0 {
        log.Info().Int64("affected", affected).Msg("expired invitations")
    }
}
```

## Обработка ошибок

Repository возвращает следующие типы ошибок:

### "invitation not found"
- `GetByID` - приглашение не существует
- `GetByCode` - код не найден
- `Delete` - попытка удалить несуществующее приглашение

### "max uses reached"
- `IncrementUses` - достигнут лимит использований link приглашения

### "invitation not found or not pending"
- `MarkAsExpired` - попытка пометить как expired не-pending приглашение
- `MarkAsRevoked` - попытка отозвать не-pending приглашение

### Database errors
Все остальные ошибки обёрнуты с контекстом:
```go
if err != nil {
    if strings.Contains(err.Error(), "invitation not found") {
        // Handle not found
    } else {
        // Database or other error
    }
}
```

## Тестирование

Запуск тестов:

```bash
cd /p/github.com/sveturs/listings
go test ./internal/repository/postgres -run TestStorefrontInvitation -v
```

Все тесты используют `sqlmock` для изоляции от реальной БД.

## Связь с другими компонентами

### Domain Layer
- `internal/domain/storefront_invitation.go` - модели и бизнес-логика
- `domain.StorefrontInvitation` - основная модель
- `domain.GenerateInviteCode()` - генерация уникальных кодов

### Migration
- `migrations/20251205000001_storefront_invitations.up.sql` - создание таблицы

### Service Layer (TODO)
- `internal/service/invitation_service.go` - бизнес-логика приглашений
- Интеграция с email отправкой
- Интеграция с storefront_staff

## Performance considerations

### Индексы

Таблица имеет следующие индексы:
- `idx_storefront_invitations_storefront` - для ListByStorefront
- `idx_storefront_invitations_email` - для GetByEmailAndStorefront
- `idx_storefront_invitations_user` - для ListByUser
- `idx_storefront_invitations_code` - для GetByCode (UNIQUE)
- `idx_storefront_invitations_status` - для фильтрации по статусу
- `idx_storefront_invitations_expires` - для ExpirePendingInvitations

### Атомарные операции

`IncrementUses` использует атомарный UPDATE для предотвращения race conditions:
```sql
UPDATE storefront_invitations
SET current_uses = current_uses + 1
WHERE id = $1 AND (max_uses IS NULL OR current_uses < max_uses)
```

### Batch operations

`ExpirePendingInvitations` обновляет множество записей за один запрос:
```sql
UPDATE storefront_invitations
SET status = 'expired'
WHERE status = 'pending'
  AND expires_at IS NOT NULL
  AND expires_at < NOW()
```

## TODO

- [ ] Добавить метод `BulkCreate` для массового создания приглашений
- [ ] Добавить метод `GetExpiringSoon` для напоминаний
- [ ] Интеграция с Redis для кеширования активных link приглашений
- [ ] Метрики использования (prometheus)

---

**Версия:** 1.0
**Дата:** 2025-12-05
**Автор:** Repository layer для Listings Microservice
