# API Contract Management (OpenAPI v3)

Проект использует подход **Contract-First Development** на основе OpenAPI v3 схемы:

### Workflow взаимодействия Backend ↔ Frontend

1. **Backend (первичный источник)**:
    - Разработчики пишут Swagger аннотации в Go коде
    - Аннотации описывают все endpoints, их параметры и типы ответов

2. **Генерация OpenAPI схемы**:
   ```bash
   cd backend && make generate-types
   ```
    - Создается `docs/swagger.json` - OpenAPI v3 спецификация
    - Создает типизированные интерфейсы в `frontend/svetu/src/types/generated/api.ts`

### Важно при разработке

- **Backend разработчики**: Всегда документируйте endpoints с помощью Swagger аннотаций
- **Frontend разработчики**: Используйте сгенерированные типы из `@/types/generated/api`
- **При изменении API**: Обязательно перегенерируйте типы командой `make generate-types`
