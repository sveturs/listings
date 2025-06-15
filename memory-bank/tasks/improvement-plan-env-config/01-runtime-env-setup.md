# Шаг 1: Внедрение next-runtime-env

## Цель
Установить и настроить библиотеку `next-runtime-env` для поддержки runtime конфигурации в Next.js приложении.

## Задачи

### 1.1 Установка библиотеки
```bash
cd frontend/svetu
yarn add next-runtime-env
```

### 1.2 Обновление корневого layout

Файл: `/frontend/svetu/src/app/[locale]/layout.tsx`

```typescript
import { PublicEnvScript } from 'next-runtime-env';

export default function RootLayout({ 
  children,
  params: { locale }
}: { 
  children: React.ReactNode;
  params: { locale: string };
}) {
  return (
    <html lang={locale}>
      <head>
        {/* Добавляем скрипт для runtime переменных */}
        <PublicEnvScript />
      </head>
      <body>
        {/* Существующий код */}
        {children}
      </body>
    </html>
  );
}
```

### 1.3 Создание wrapper для env функции

Создать файл: `/frontend/svetu/src/utils/env.ts`

```typescript
import { env as runtimeEnv } from 'next-runtime-env';

/**
 * Безопасный доступ к runtime переменным окружения
 * На сервере использует process.env, на клиенте - runtime env
 */
export function getEnv(key: string, defaultValue?: string): string | undefined {
  if (typeof window === 'undefined') {
    // Server-side: используем process.env
    return process.env[key] || defaultValue;
  }
  
  // Client-side: используем runtime env
  return runtimeEnv(key) || defaultValue;
}

/**
 * Типизированный доступ к публичным переменным
 */
export const publicEnv = {
  get API_URL() {
    return getEnv('NEXT_PUBLIC_API_URL', 'http://localhost:3000');
  },
  get MINIO_URL() {
    return getEnv('NEXT_PUBLIC_MINIO_URL', 'http://localhost:9000');
  },
  get WEBSOCKET_URL() {
    return getEnv('NEXT_PUBLIC_WEBSOCKET_URL');
  },
  get IMAGE_HOSTS() {
    return getEnv('NEXT_PUBLIC_IMAGE_HOSTS');
  },
  get ENABLE_PAYMENTS() {
    return getEnv('NEXT_PUBLIC_ENABLE_PAYMENTS') === 'true';
  },
};
```

## Проверка

### Тестовый компонент
Создать файл: `/frontend/svetu/src/app/[locale]/test-env/page.tsx`

```typescript
'use client';

import { env } from 'next-runtime-env';
import { publicEnv } from '@/utils/env';

export default function TestEnvPage() {
  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Runtime Environment Test</h1>
      
      <div className="space-y-2">
        <h2 className="text-xl font-semibold">Direct access:</h2>
        <p>API URL: {env('NEXT_PUBLIC_API_URL')}</p>
        <p>Minio URL: {env('NEXT_PUBLIC_MINIO_URL')}</p>
        
        <h2 className="text-xl font-semibold mt-4">Typed access:</h2>
        <p>API URL: {publicEnv.API_URL}</p>
        <p>Minio URL: {publicEnv.MINIO_URL}</p>
        <p>Payments enabled: {publicEnv.ENABLE_PAYMENTS ? 'Yes' : 'No'}</p>
      </div>
    </div>
  );
}
```

### Команды для проверки
```bash
# Запустить dev сервер
yarn dev -p 3001

# Открыть тестовую страницу
# http://localhost:3001/test-env

# Изменить переменные и перезапустить (без пересборки!)
NEXT_PUBLIC_API_URL=https://api.example.com yarn dev -p 3001
```

## Важные моменты

1. **PublicEnvScript** должен быть добавлен в `<head>` до любых клиентских скриптов
2. Переменные должны начинаться с `NEXT_PUBLIC_` чтобы быть доступными на клиенте
3. На сервере по-прежнему используется `process.env`
4. Runtime переменные работают только после гидратации React

## Возможные проблемы

### 1. TypeScript ошибки
Если TypeScript не видит типы `next-runtime-env`, добавить в `tsconfig.json`:
```json
{
  "compilerOptions": {
    "types": ["next-runtime-env"]
  }
}
```

### 2. Переменные не обновляются
Убедитесь, что:
- Переменные начинаются с `NEXT_PUBLIC_`
- Сервер был перезапущен после изменения переменных
- PublicEnvScript добавлен в layout

### 3. SSR/CSR несоответствие
При использовании переменных в компонентах, которые рендерятся и на сервере и на клиенте, используйте `useEffect`:
```typescript
const [apiUrl, setApiUrl] = useState(publicEnv.API_URL);

useEffect(() => {
  setApiUrl(publicEnv.API_URL);
}, []);
```

## Результат
После выполнения этого шага приложение будет поддерживать изменение публичных переменных окружения без пересборки.

## Критерии приемки

### 1. Установка библиотеки
- [x] `next-runtime-env` добавлен в package.json
- [x] yarn.lock обновлен после установки
- [x] Библиотека установлена локально (проверить через `yarn list next-runtime-env`)

### 2. Интеграция в layout
- [x] `PublicEnvScript` импортирован в `/src/app/[locale]/layout.tsx`
- [x] Скрипт добавлен в `<head>` секцию до других скриптов
- [x] Приложение запускается без ошибок (`yarn dev -p 3001`)

### 3. Создание утилиты env.ts
- [x] Файл `/src/utils/env.ts` создан
- [x] Функция `getEnv()` корректно определяет контекст (server/client)
- [x] Объект `publicEnv` предоставляет типизированный доступ к переменным
- [x] TypeScript не выдает ошибок при импорте и использовании

### 4. Проверка работоспособности
- [x] Тестовая страница `/test-env` создана и доступна по адресу http://localhost:3001/test-env
- [x] Страница отображает значения переменных окружения
- [x] При изменении переменных и перезапуске сервера значения обновляются:
  ```bash
  NEXT_PUBLIC_API_URL=https://test.com yarn dev -p 3001
  ```
- [x] В браузере в консоли доступен объект `window.__ENV` с переменными
- [x] Нет ошибок гидратации React в консоли браузера

### 5. Дополнительные проверки
- [x] Переменные без префикса `NEXT_PUBLIC_` не попадают в браузер
- [x] Runtime переменные работают как в dev, так и в production режиме
- [x] Существующий функционал приложения не сломан