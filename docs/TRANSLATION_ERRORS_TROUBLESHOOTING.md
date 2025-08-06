# Инструкция по устранению ошибок переводов

## 🚨 Проблема: MISSING_MESSAGE и INSUFFICIENT_PATH ошибки

После 5 часов отладки найдены типичные проблемы и их решения:

## ✅ Диагностика ошибок

### 1. MISSING_MESSAGE: Could not resolve `key` in messages
**Причина:** Отсутствует файл модуля переводов или ключ не найден

**Решение:**
```bash
# 1. Найти где используется ключ
grep -r "useTranslations('moduleName')" src/

# 2. Проверить существует ли файл модуля
ls src/messages/ru/moduleName.json

# 3. Если файла нет - создать для всех локалей:
# src/messages/ru/moduleName.json
# src/messages/en/moduleName.json  
# src/messages/sr/moduleName.json

# 4. Добавить модуль в loadMessages.ts в switch case
# 5. Добавить в index.ts всех локалей
# 6. Добавить в тип TranslationModule в loadMessages.ts
```

### 2. INSUFFICIENT_PATH: Message resolved to an object
**Причина:** Конфликт имен - одинаковое имя для строки и объекта

**Пример проблемы:**
```json
{
  "profile": "Профиль",    // строка
  "profile": {             // объект - КОНФЛИКТ!
    "title": "Мой профиль"
  }
}
```

**Решение:** Переименовать объект
```json
{
  "profile": "Профиль",    // строка для кнопки
  "profilePage": {         // объект для страницы
    "title": "Мой профиль"
  }
}
```

## 🔧 Системные исправления

### 1. Исправить loadMessages.ts
Проблема с Object.assign - перезаписывает ключи:

```typescript
// ❌ ПЛОХО
Object.assign(messages, data);

// ✅ ХОРОШО  
for (const [key, value] of Object.entries(data)) {
  if (!messages[key]) {
    messages[key] = value;
  }
}
```

### 2. Двойной namespace в auth.json
Структура должна быть:
```json
{
  "profile": "Профиль",           // для useTranslations('auth')
  "myListings": "Мои объявления", // для useTranslations('auth')
  "auth": {                       // для useTranslations('auth.something')
    "profile": "Профиль",
    "myListings": "Мои объявления"
  }
}
```

### 3. Добавление нового модуля переводов
Полный чеклист:

1. **Создать JSON файлы:**
   - `src/messages/ru/moduleName.json`
   - `src/messages/en/moduleName.json`
   - `src/messages/sr/moduleName.json`

2. **Обновить loadMessages.ts:**
   ```typescript
   // Добавить в тип
   export type TranslationModule = 
     | 'existing'
     | 'moduleName' // <-- добавить
     
   // Добавить в switch
   case 'moduleName':
     moduleData = await import(`@/messages/${locale}/moduleName.json`);
     break;
   ```

3. **Обновить index.ts всех локалей:**
   ```typescript
   export type TranslationModule =
     | 'existing'
     | 'moduleName'; // <-- добавить
     
   export const moduleLoaders = {
     existing: () => import('./existing.json'),
     moduleName: () => import('./moduleName.json'), // <-- добавить
   };
   ```

4. **Добавить в layout.tsx если нужно глобально:**
   ```typescript
   const messages = await loadMessages(locale as any, [
     'common',
     'moduleName', // <-- добавить
   ]);
   ```

## 🚀 Быстрая проверка

```bash
# Очистить кэш и проверить
rm -rf .next && yarn build

# Если есть ошибки - смотреть в логи:
tail -50 /tmp/frontend.log | grep "MISSING_MESSAGE\|INSUFFICIENT_PATH"
```

## 📝 Часто встречающиеся модули

- `checkout` - оформление заказа
- `search` - поиск  
- `auth` - авторизация (watch out для конфликтов profile!)
- `common` - базовые переводы (всегда должен быть первым)

## 🔴 КРИТИЧЕСКАЯ ОШИБКА: Двойная обертка в модулях переводов

### Проблема
При использовании `useTranslations('admin')` возникали ошибки типа:
- `IntlError: MISSING_MESSAGE: Could not resolve 'admin.variantAttributes.types.color'`

Хотя файл `admin.json` содержал все ключи.

### Причина
Файлы модулей переводов имели двойную вложенность:
```json
// ❌ НЕПРАВИЛЬНО - admin.json
{
  "admin": {
    "variantAttributes": {
      "types": {
        "color": "Цвет"
      }
    }
  }
}
```

При загрузке через `loadMessages` получался путь: `admin.admin.variantAttributes.types.color`

### Решение
Убрать лишнюю обертку с именем модуля:
```json
// ✅ ПРАВИЛЬНО - admin.json
{
  "variantAttributes": {
    "types": {
      "color": "Цвет"
    }
  }
}
```

### Автоматическое исправление
```javascript
// scripts/fix-admin-structure.js
const data = JSON.parse(content);
if (data.admin && typeof data.admin === 'object') {
  const adminContent = data.admin;
  const newData = { ...adminContent, ...otherKeys };
  fs.writeFileSync(filePath, JSON.stringify(newData, null, 2));
}
```

### Важно после исправления
1. Остановить сервер: `/home/dim/.local/bin/kill-port-3001.sh`
2. Очистить кэш: `rm -rf .next`
3. Запустить заново: `/home/dim/.local/bin/start-frontend-screen.sh`
4. Жесткое обновление страницы: Ctrl+F5

## 🎯 Результат

После применения всех исправлений:
- ✅ `yarn build` без ошибок MISSING_MESSAGE
- ✅ Все локали работают (ru, en, sr)
- ✅ Development и production режимы стабильны
- ✅ Модульная система переводов работает корректно

**Время решения:** 5 часов → 15 минут с этой инструкцией! 🎉

## 🔴 INSUFFICIENT_PATH ошибка: конфликт между строкой и объектом

### Проблема
Ошибка `IntlError: INSUFFICIENT_PATH: Message at 'admin.attributes' resolved to an object` возникает когда компонент ожидает строку, но в JSON файле по этому пути находится объект.

### Пример проблемы:
```typescript
// Компонент ожидает строку
{t('attributes')} // Ошибка!

// JSON файл содержит объект
{
  "attributes": {
    "types": { ... }
  }
}
```

### Решение: использовать разные ключи
```json
{
  "attributesTitle": "Атрибуты",        // строка для заголовка
  "attributeGroupsTitle": "Группы атрибутов", // строка для заголовка
  "attributes": {                        // объект с вложенными переводами
    "types": {
      "multiselect": "Множественный выбор"
    }
  }
}
```

```typescript
// В компоненте
{t('attributesTitle')}          // для заголовка
{t('attributes.types.multiselect')} // для вложенных значений
```

## 🔐 Ошибки авторизации 401 для API запросов

### Проблема
При попытке обращения к защищенным эндпоинтам (например, dashboard API) возникают ошибки 401 Unauthorized, хотя пользователь авторизован.

### Причина
Токен авторизации не передается в заголовках запроса. Система хранит токен в `sessionStorage` через `tokenManager`, но API клиент не знает где его искать.

### Решение

1. **Создать обертку для API клиента с авторизацией:**

```typescript
// src/lib/api-client-auth.ts
import { apiClient } from './api-client';
import { tokenManager } from '@/utils/tokenManager';

// Функция для получения токена
function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null;
  
  // Сначала пробуем получить токен из tokenManager
  const token = tokenManager.getAccessToken();
  if (token) {
    return token;
  }
  
  // Если tokenManager не инициализирован, пробуем sessionStorage напрямую
  return sessionStorage.getItem('svetu_access_token');
}

// Расширенный API клиент с автоматической авторизацией
export const apiClientAuth = {
  async get(path: string, options?: any) {
    const token = getAuthToken();
    const headers = {
      ...options?.headers,
    };
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    return apiClient.get(path, { ...options, headers });
  },
  // ... аналогично для post, put, delete
};
```

2. **Использовать apiClientAuth для защищенных запросов:**

```typescript
// Вместо apiClient используем apiClientAuth
import { apiClientAuth } from '@/lib/api-client-auth';

// В Redux thunk
const response = await apiClientAuth.get(
  `/api/v1/storefronts/${slug}/dashboard/stats`
);
```

### Важно знать
- Токен сохраняется в `sessionStorage` под ключом `svetu_access_token`
- `TokenManager` автоматически обновляет токен перед истечением
- При logout токен очищается из всех хранилищ

### Результат
✅ Все защищенные API запросы проходят успешно с правильными заголовками авторизации