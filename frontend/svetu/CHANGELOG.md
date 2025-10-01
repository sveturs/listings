# Frontend Changelog

All notable changes to the frontend will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-10-01

### Added

- **BFF Proxy**: Реализован Backend-for-Frontend прокси `/api/v2/[...path]`
  - Универсальный прокси для всех API запросов к backend
  - Автоматическое добавление JWT токенов из httpOnly cookies
  - Централизованная обработка авторизации
  - Поддержка всех HTTP методов: GET, POST, PUT, DELETE, PATCH
  - Логирование всех запросов для отладки

### Changed

- **API Client** (`src/services/api-client.ts`):

  - Все запросы теперь идут через `/api/v2` вместо прямого обращения к backend
  - Автоматический strip `/api/v1/` префикса для BFF совместимости
  - Добавлена опция `direct: true` для редких случаев прямого доступа

- **Admin Service** (`src/services/admin.ts`): Масштабная очистка (-527 строк)

  - Убраны все `/api/v1/` префиксы из эндпоинтов
  - Удалена функция `getAuthHeaders()` (рудимент)
  - Удалены все импорты `tokenManager`
  - Удален импорт `configManager` (больше не нужен)
  - Все прямые `fetch()` вызовы заменены на `apiClient.request()`
  - Упрощена обработка ошибок

- **Next.js Config** (`next.config.ts`):

  - Исключен `/api/v2` из rewrites для корректной работы BFF proxy
  - Rewrite pattern изменен: `/api/:path*` → `/api/:path((?!v2).*)*`

- **Version**: Updated from 0.1.1 to 0.2.0 (`package.json`)

### Removed

- Удален весь legacy код авторизации:
  - `getAuthHeaders()` function
  - Все импорты `tokenManager`
  - Прямые `fetch()` вызовы к backend
  - Ручное управление `Authorization` headers

### Security

- ✅ JWT токены теперь хранятся только в httpOnly cookies (недоступны JavaScript)
- ✅ Устранена возможность XSS атак на токены
- ✅ Централизованная авторизация через BFF proxy
- ✅ Нет CORS проблем (все на одном домене)

### Technical Details

- **New Files**:

  - `src/app/api/v2/[...path]/route.ts` (159 строк) - BFF proxy route handler

- **Modified Files**:
  - `src/services/api-client.ts` (+44 строки)
  - `src/services/admin.ts` (-527 строк)
  - `next.config.ts` (rewrite pattern)
  - `package.json` (version)

### Environment Variables

- **New**: `BACKEND_INTERNAL_URL` - URL для backend (server-side)
  - Default: `http://localhost:33423` (необычный порт для легкого обнаружения проблем конфигурации)
  - Recommended: `http://localhost:3000` (production)

### Migration Guide

#### До (0.1.1):

```typescript
// ❌ Старый код - НЕ используй!
import { tokenManager } from '@/utils/tokenManager';

const headers = await getAuthHeaders();
const response = await fetch(`${apiUrl}/api/v1/admin/categories`, {
  headers,
});
```

#### После (0.2.0):

```typescript
// ✅ Новый код
import { apiClient } from '@/services/api-client';

const response = await apiClient.get('/admin/categories');
```

### Breaking Changes

**НЕТ** - обратная совместимость сохранена

### Notes

- См. `CLAUDE.md` для полной документации BFF архитектуры
- Все компоненты автоматически используют новую архитектуру
- TypeScript компиляция без ошибок

## [0.1.1] - 2025-09-XX

### Initial Release

- Next.js 15 App Router
- React 19
- TypeScript
- Tailwind CSS
- Redux Toolkit
- i18n (en, ru, sr)
- Auth integration
- Admin panel
