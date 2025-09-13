# Project Constitution: Translation System Fixes

## Core Principles

1. **100% Translation Coverage**: Все тексты должны быть переведены на все поддерживаемые языки (ru, en, sr)
2. **Zero MISSING_MESSAGE Errors**: В production не должно быть ошибок отсутствующих переводов
3. **Automated Validation**: Все изменения должны проходить автоматическую валидацию
4. **Professional Quality**: Переводы должны быть профессиональными и учитывать контекст

## Technology Constraints

- **Frontend**: Next.js 15.3.2 with next-intl
- **Languages**: English (en), Russian (ru), Serbian (sr - латиница)
- **Structure**: Modular JSON files per language and module
- **Location**: frontend/svetu/src/messages/{locale}/{module}.json

## Development Rules

1. Все переводы добавляются одновременно для всех языков
2. Структура JSON должна быть идентичной между языками
3. Сербские переводы используют латиницу, не кириллицу
4. Каждый модуль должен быть самодостаточным

## Quality Standards

- Минимальное покрытие переводами: 95%
- Максимальное время загрузки переводов: 100ms
- Все ключи должны следовать camelCase конвенции
- Максимальная вложенность: 3 уровня

## Prohibited Actions

- НЕ изменять архитектуру системы переводов
- НЕ использовать машинный перевод без проверки
- НЕ добавлять переводы только для одного языка
- НЕ игнорировать контекст использования