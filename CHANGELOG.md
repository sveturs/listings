# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-10-01

### Added
- **BFF Proxy Architecture**: Реализован Backend-for-Frontend прокси `/api/v2`
  - Универсальный прокси для всех API запросов от frontend к backend
  - Автоматическое добавление JWT токенов из httpOnly cookies
  - Централизованная обработка авторизации
  - Исключен `/api/v2` из next.config.ts rewrites для корректной работы

### Changed
- **Frontend API Client**: Полностью переработан для работы через BFF proxy
  - Все запросы теперь идут через `/api/v2` вместо прямого обращения к backend
  - Убраны все `/api/v1/` префиксы из эндпоинтов
  - Удален рудимент `getAuthHeaders()` функция
  - Удалены все прямые `fetch()` вызовы к backend
  - Убран `tokenManager` из клиентского кода
  - Удален неиспользуемый импорт `configManager` из `admin.ts`

- **CLAUDE.md**: Добавлено критическое правило #8
  - Frontend → Backend: ВСЕГДА через BFF proxy `/api/v2`
  - Добавлена подробная секция "BFF Proxy Architecture"
  - Примеры правильного и неправильного использования
  - Документация по архитектуре и преимуществам BFF

### Security
- JWT токены теперь хранятся только в httpOnly cookies (недоступны JavaScript)
- Устранена возможность XSS атак на токены
- Централизованная авторизация через BFF proxy

### Technical Details
- **Backend Version**: 0.1.1 → 0.2.0
- **Frontend Version**: 0.1.1 → 0.2.0
- **New Files**:
  - `frontend/svetu/src/app/api/v2/[...path]/route.ts` - BFF proxy route handler
- **Modified Files**:
  - `frontend/svetu/src/services/api-client.ts` - поддержка BFF proxy
  - `frontend/svetu/src/services/admin.ts` - упрощение и очистка от legacy кода
  - `frontend/svetu/next.config.ts` - исключен `/api/v2` из rewrites
  - `CLAUDE.md` - новое правило и документация BFF
  - `backend/internal/version/version.go` - версия 0.2.0
  - `frontend/svetu/package.json` - версия 0.2.0

### Environment Variables
- **New**: `BACKEND_INTERNAL_URL` - URL для backend (server-side)
  - Default fallback: `http://localhost:33423` (необычный порт для легкого обнаружения проблем)

## [0.1.1] - 2025-09-XX

### Initial Release
- Базовая функциональность маркетплейса
- Авторизация через Auth Service
- Категории и листинги
- Admin панель
