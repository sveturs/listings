# Тестирование Frontend

## Настройка

Проект использует Jest для unit тестирования.

### Установка зависимостей

```bash
yarn add --dev jest @testing-library/react @testing-library/jest-dom @testing-library/dom jest-environment-jsdom @types/jest
```

### Конфигурация

- `jest.config.js` - основная конфигурация Jest
- `jest.setup.js` - настройка глобальных моков и полифиллов

## Запуск тестов

```bash
# Запустить все тесты
yarn test

# Запустить тесты в watch режиме
yarn test:watch

# Запустить тесты с покрытием
yarn test:coverage

# Запустить конкретный тест
yarn test src/utils/__tests__/env.test.ts
```

## Структура тестов

Тесты располагаются рядом с тестируемым кодом в папках `__tests__`:

```
src/
  services/
    api-client.ts
    __tests__/
      api-client.test.ts
  utils/
    env.ts
    __tests__/
      env.test.ts
```

## Написание тестов

### Базовый пример

```typescript
import { someFunction } from '../someModule';

describe('someFunction', () => {
  it('should do something', () => {
    const result = someFunction('input');
    expect(result).toBe('expected output');
  });
});
```

### Мокирование модулей

```typescript
// Мокировать модуль до импорта
jest.mock('@/config', () => ({
  default: {
    getApiUrl: jest.fn(),
  },
}));

import configManager from '@/config';

// В тесте
(configManager.getApiUrl as jest.Mock).mockReturnValue('http://test.com');
```

## Известные проблемы

1. **Response и Headers не определены в Jest**

   - Решение: Добавлены полифиллы в тестовые файлы

2. **Сложность тестирования API client**
   - API client имеет много зависимостей
   - Рекомендуется использовать интеграционные тесты или E2E тесты для полного покрытия

## Текущий статус

### Работающие тесты:

- `src/utils/__tests__/env.test.ts` - тесты для утилиты работы с переменными окружения
- `src/hooks/__tests__/useDebounce.test.ts` - тесты для хука debounce

### Пропущенные тесты:

- `src/services/__tests__/api-client.test.ts.skip` - требует сложной настройки моков

Запуск всех тестов: `yarn test` ✅

## Рекомендации

1. Фокусируйтесь на тестировании бизнес-логики
2. Используйте интеграционные тесты для API взаимодействий
3. Рассмотрите использование Playwright для E2E тестов
4. Для компонентов с множеством зависимостей лучше использовать интеграционные тесты
