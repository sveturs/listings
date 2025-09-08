# Инструкция по очистке кэша браузера для решения проблемы с OAuth

## Проблема
После изменения конфигурации с IP адреса `100.88.44.15` на `localhost`, браузер продолжает использовать старые URL из-за кэширования.

## Решение - Очистка кэша браузера

### Для Google Chrome:

1. **Откройте DevTools** (F12 или правый клик → Inspect)

2. **Очистите localStorage и sessionStorage:**
   - Перейдите во вкладку **Application** (или **Storage**)
   - В левой панели найдите **Local Storage** → `http://localhost:3001`
   - Нажмите правой кнопкой и выберите **Clear**
   - Повторите для **Session Storage**

3. **Очистите куки:**
   - В той же вкладке Application найдите **Cookies**
   - Очистите куки для:
     - `http://localhost:3001`
     - `http://localhost:3000`
     - `http://100.88.44.15:3001` (если есть)
     - `http://100.88.44.15:3000` (если есть)

4. **Полная очистка кэша:**
   - Нажмите **Ctrl+Shift+Delete** (или Cmd+Shift+Delete на Mac)
   - Выберите временной диапазон: **All time**
   - Отметьте:
     - ✓ Cookies and other site data
     - ✓ Cached images and files
   - Нажмите **Clear data**

5. **Альтернативный способ - Hard Reload:**
   - Откройте DevTools (F12)
   - Правый клик на кнопку обновления браузера
   - Выберите **Empty Cache and Hard Reload**

### Для Firefox:

1. **Очистка через DevTools:**
   - F12 → Storage → Clear All

2. **Полная очистка:**
   - Ctrl+Shift+Delete
   - Выберите все данные
   - Clear Now

### Для Safari:

1. **Включите меню разработчика:**
   - Preferences → Advanced → Show Develop menu

2. **Очистите кэш:**
   - Develop → Empty Caches

## После очистки кэша:

1. **Полностью закройте браузер**
2. **Откройте браузер заново**
3. **Перейдите на** `http://localhost:3001`
4. **Откройте DevTools Console** и убедитесь, что все запросы идут на `localhost`, а не на IP адрес

## Проверка в консоли браузера:

```javascript
// Проверьте текущие переменные окружения
console.log('API URL:', process.env.NEXT_PUBLIC_API_URL);
console.log('Frontend URL:', process.env.NEXT_PUBLIC_FRONTEND_URL);

// Очистите localStorage вручную
localStorage.clear();
sessionStorage.clear();

// Перезагрузите страницу
location.reload();
```

## Если проблема сохраняется:

1. **Используйте режим инкогнито** (Ctrl+Shift+N в Chrome)
2. **Проверьте файл hosts:**
   ```bash
   cat /etc/hosts | grep -E "localhost|100.88.44.15"
   ```
3. **Перезапустите все сервисы:**
   ```bash
   # Frontend
   /home/dim/.local/bin/kill-port-3001.sh && /home/dim/.local/bin/start-frontend-screen.sh
   
   # Backend
   /home/dim/.local/bin/kill-port-3000.sh
   screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'
   
   # Auth Service
   cd /data/auth_svetu && docker-compose restart
   ```

## Важно!

Frontend был полностью пересобран с новыми переменными окружения. Теперь все запросы должны идти на `localhost`, а не на IP адрес VPN.