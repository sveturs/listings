# Миграция с относительных API путей на абсолютные URL

## 📋 Описание проблемы

Сейчас в проекте используются два способа обращения к API:
1. **Относительные пути**: `/api/v1/*` - проксируются через nginx с `dev.svetu.rs/api`
2. **Абсолютные URL**: `https://devapi.svetu.rs/api/v1/*` - прямые запросы на API сервер

Это создает путаницу и усложняет конфигурацию. Необходимо полностью перейти на использование абсолютных URL.

## 🎯 Цель миграции

- Унифицировать все API запросы через единый URL: `https://devapi.svetu.rs`
- Убрать ненужное проксирование через nginx
- Упростить конфигурацию и повысить производительность

## 📊 Текущее состояние

### ✅ Уже мигрировано:
- `src/services/api.ts` - основной API клиент
- `src/store/slices/storefrontSlice.ts` - Redux slice для витрин
- SSR запросы через `api-client-server.ts`

### ❌ Требуют миграции (30 файлов):
```
frontend/svetu/src/app/[locale]/admin/auth/page.tsx
frontend/svetu/src/app/[locale]/admin/postexpress/page.tsx
frontend/svetu/src/app/[locale]/admin/search/components/SearchWeights.tsx
frontend/svetu/src/app/[locale]/admin/search/components/WeightOptimization.tsx
frontend/svetu/src/app/[locale]/admin/variant-attributes/VariantAttributesClient.tsx
frontend/svetu/src/app/[locale]/create-listing-ai/page.tsx
frontend/svetu/src/app/[locale]/create-listing-smart/page.tsx
frontend/svetu/src/app/[locale]/docs/page.tsx
frontend/svetu/src/app/[locale]/examples/novi-sad-districts/manage/page.tsx
frontend/svetu/src/app/[locale]/user-contacts/page.tsx
frontend/svetu/src/components/GIS/hooks/useVisibleCities.ts
frontend/svetu/src/components/Storefront/ProductVariants/VariantGenerator.tsx
frontend/svetu/src/components/Storefront/ProductVariants/VariantManager.tsx
frontend/svetu/src/components/admin/translations/AITranslations.tsx
frontend/svetu/src/components/delivery/bexexpress/BEXAddressForm.tsx
frontend/svetu/src/components/delivery/bexexpress/BEXDeliverySelector.tsx
frontend/svetu/src/components/delivery/bexexpress/BEXDeliveryStep.tsx
frontend/svetu/src/components/delivery/postexpress/PostExpressDeliverySelector.tsx
frontend/svetu/src/components/delivery/postexpress/PostExpressRateCalculator.tsx
frontend/svetu/src/components/products/EnhancedVariantGenerator.tsx
frontend/svetu/src/components/products/SimplifiedVariantGenerator.tsx
frontend/svetu/src/components/search/QuerySuggestions.tsx
frontend/svetu/src/components/shared/ARProductViewer.tsx
frontend/svetu/src/components/shared/QRBarcodeScanner.tsx
frontend/svetu/src/hooks/useAnalytics.ts
frontend/svetu/src/services/abTestingService.ts
frontend/svetu/src/services/admin.ts
frontend/svetu/src/services/ai/claude.service.ts
frontend/svetu/src/services/biometricAuthService.ts
frontend/svetu/src/services/translationAdminApi.ts
```

## 🔧 Как мигрировать файл

### 1. Для компонентов и страниц

**Было:**
```typescript
const response = await fetch('/api/v1/some-endpoint', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(data)
});
```

**Стало:**
```typescript
import configManager from '@/config';

// В начале функции/компонента
const apiUrl = configManager.get('api.url');

const response = await fetch(`${apiUrl}/api/v1/some-endpoint`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(data)
});
```

### 2. Для сервисов (рекомендуется использовать api клиент)

**Было:**
```typescript
export const someService = {
  async getData() {
    const response = await fetch('/api/v1/data');
    return response.json();
  }
};
```

**Стало (вариант 1 - использовать готовый api клиент):**
```typescript
import api from '@/services/api';

export const someService = {
  async getData() {
    const response = await api.get('/api/v1/data');
    return response.data;
  }
};
```

**Стало (вариант 2 - если нужен чистый fetch):**
```typescript
import configManager from '@/config';

const API_URL = configManager.get('api.url');

export const someService = {
  async getData() {
    const response = await fetch(`${API_URL}/api/v1/data`);
    return response.json();
  }
};
```

### 3. Для хуков

**Было:**
```typescript
export function useCustomHook() {
  const fetchData = async () => {
    const res = await fetch('/api/v1/endpoint');
    // ...
  };
}
```

**Стало:**
```typescript
import { useMemo } from 'react';
import configManager from '@/config';

export function useCustomHook() {
  const apiUrl = useMemo(() => configManager.get('api.url'), []);
  
  const fetchData = async () => {
    const res = await fetch(`${apiUrl}/api/v1/endpoint`);
    // ...
  };
}
```

## 📝 Пошаговый план миграции

### Этап 1: Подготовка (выполнено ✅)
1. ✅ Обновить `src/services/api.ts` для использования `NEXT_PUBLIC_API_URL`
2. ✅ Проверить что `NEXT_PUBLIC_API_URL=https://devapi.svetu.rs` в `.env`
3. ✅ Обновить критичные файлы (storefrontSlice.ts)

### Этап 2: Массовая миграция
1. Обновить все оставшиеся файлы согласно списку выше
2. Использовать скрипт проверки: `./check-api-usage.sh`
3. Провести тестирование каждого модуля после изменений

### Этап 3: Тестирование
1. Запустить полное тестирование приложения
2. Проверить работу:
   - Административных страниц
   - Системы доставки (BEX, PostExpress)
   - AI функций
   - Аналитики
   - Всех форм создания/редактирования

### Этап 4: Отключение проксирования
1. Убедиться что все файлы мигрированы
2. Отредактировать `/opt/nginx-simple/conf.d/dev.svetu.rs.conf`:
   ```nginx
   # Удалить этот блок:
   location /api {
       proxy_pass http://172.17.0.1:3002;
       # ...
   }
   ```
3. Перезапустить nginx: `sudo nginx -s reload`
4. Провести финальное тестирование

## 🛠️ Утилиты для помощи

### Скрипт проверки использования API
```bash
#!/bin/bash
# check-api-usage.sh - уже создан в корне проекта
./check-api-usage.sh
```

### Команда для поиска файлов с относительными путями
```bash
# Найти все файлы с fetch('/api
grep -r "fetch(['\"]\/api" frontend/svetu/src --include="*.ts" --include="*.tsx"

# Найти все файлы с axios и относительными путями
grep -r "axios.*['\"]\/api" frontend/svetu/src --include="*.ts" --include="*.tsx"
```

### Автоматизированная замена (осторожно!)
```bash
# Пример sed команды для замены (ВСЕГДА делайте backup!)
# НЕ ИСПОЛЬЗОВАТЬ без проверки каждого изменения!
sed -i.bak "s/fetch('\\/api/fetch(\`\${apiUrl}\\/api/g" filename.tsx
```

## ⚠️ Важные замечания

1. **SSR запросы** - НЕ трогать! Они используют `INTERNAL_API_URL` для внутренней связи между контейнерами
2. **Токены авторизации** - При использовании `api` клиента токены добавляются автоматически
3. **CORS** - После отключения проксирования убедитесь что CORS настроен правильно на backend
4. **WebSocket** - Уже использует правильный URL: `wss://devapi.svetu.rs`

## 📊 Преимущества после миграции

- **Производительность**: Убираем лишний hop через nginx прокси
- **Простота**: Единый способ обращения к API
- **Масштабируемость**: Легче разделить frontend и backend на разные сервера
- **Отладка**: Проще отслеживать запросы в DevTools
- **Безопасность**: Явное указание API сервера, нет скрытого проксирования

## 🚨 Rollback план

Если что-то пошло не так после отключения проксирования:

1. Быстро вернуть проксирование в nginx:
   ```bash
   # Восстановить блок location /api в dev.svetu.rs.conf
   sudo nginx -s reload
   ```

2. Откатить изменения в коде:
   ```bash
   git revert <commit-hash>
   ```

3. Перезапустить frontend:
   ```bash
   docker-compose -f docker-compose.dev.yml restart frontend
   ```

## 📞 Контакты для вопросов

При возникновении вопросов или проблем во время миграции, обращайтесь в соответствующие каналы поддержки проекта.

---

*Документ создан: 2025-09-08*  
*Последнее обновление: 2025-09-08*