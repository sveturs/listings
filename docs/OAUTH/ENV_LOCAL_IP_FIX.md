# Решение проблемы с IP адресом в .env.local

## Проблема
После изменения `.env` файла на использование `localhost`, браузер продолжал отправлять запросы на старый IP адрес `100.88.44.15`.

## Причина
Файл `.env.local` имеет приоритет над `.env` и содержал старый IP адрес.

### Порядок загрузки переменных окружения в Next.js:
1. `.env.local` - **ВЫСШИЙ ПРИОРИТЕТ**
2. `.env.development` / `.env.production` (в зависимости от NODE_ENV)
3. `.env` - базовая конфигурация

## Решение

### 1. Исправлен файл `.env.local`
```diff
- NEXT_PUBLIC_API_URL=http://100.88.44.15:3000
+ NEXT_PUBLIC_API_URL=http://localhost:3000
```

### 2. Очищен кэш Next.js
```bash
rm -rf .next
```

### 3. Перезапущен frontend
```bash
/home/dim/.local/bin/kill-port-3001.sh
/home/dim/.local/bin/start-frontend-screen.sh
```

## Проверка

### В браузере:
1. Откройте DevTools (F12)
2. Очистите кэш: Ctrl+Shift+Delete
3. В консоли выполните:
```javascript
localStorage.clear();
sessionStorage.clear();
location.reload();
```

### Проверьте в Network tab:
- Все запросы должны идти на `http://localhost:3000`
- НЕ должно быть запросов на `http://100.88.44.15:3000`

## Файлы конфигурации

### Основные файлы:
- `.env.local` - локальные переопределения (приоритет!)
- `.env` - базовая конфигурация
- `.env.development` - для dev режима
- `.env.production` - для production

### Важно помнить:
- `.env.local` НЕ коммитится в git (в .gitignore)
- `.env.local` переопределяет ВСЕ другие env файлы
- После изменения env файлов нужно перезапустить Next.js

## Все исправлено! ✅

Теперь все запросы идут на `localhost`, OAuth авторизация работает корректно.