# Frontend Rules

## API

- Используй сгенерированные типы при работе с API
  - Example usage:
    ```typescript
    import type { components } from "@/types/generated/api";
    type DocFile = components["schemas"]["handler.DocFile"];
    ```

### Development Commands

- `yarn dev -p 3001` - Start the development server with Turbopack at http://localhost:3001
- `yarn build` - Create an optimized production build
- `yarn start` - Run the production server
- `yarn lint` - Run ESLint for code quality checks
- `yarn lint:fix` - Fix ESLint errors automatically
- `yarn format` - Format all files with Prettier
- `yarn format:check` - Check formatting without changes

### Frontend Architecture

This is a Next.js 15.3.2 application using:

- React 19 with TypeScript
- Tailwind CSS v4 for styling
- App Router (located in `src/app/`)
- ESLint configured with Next.js and TypeScript rules
- Internationalization with next-intl (en/ru locales)
- Centralized configuration management in `src/config/`
- Google OAuth 2.0 authentication
- State management with Redux Toolkit (НЕ Zustand!)
  - Store: `src/store/`
  - Slices: `src/store/slices/`
  - Hooks: `src/store/hooks.ts` и `src/hooks/useChat.ts`

### Environment Variables for Frontend

Configuration is managed through environment variables.

#### Configuration Management

All environment variables are centralized in the src/config/ module:

- `src/config/types.ts` - TypeScript interfaces for configuration
- `src/config/index.ts` - Configuration manager with helper methods

### Frontend Key Dependencies

- next-intl: For internationalization support
- prettier: Code formatter with ESLint integration
- daisyui: Component library for Tailwind CSS
- @reduxjs/toolkit: State management (Redux Toolkit)
- react-redux: React bindings for Redux

### Важная информация о поиске и индексировании

**⚠️ ВАЖНО**: Главная страница маркетплейса получает данные из OpenSearch, а НЕ напрямую из PostgreSQL.

При изменении данных в базе PostgreSQL (например, user_id объявления) изменения НЕ отобразятся на главной странице до переиндексирования OpenSearch.

Для переиндексирования используйте:

```bash
#  через команду
cd backend && ./reindex
```

### Authentication System

Authentication flow:

1. User clicks "Sign in with Google" button
2. Redirects to Google OAuth consent
3. Google redirects back to backend callback
4. Backend creates session and redirects to frontend
5. Frontend fetches user data via `/auth/session`
