# Задача: Исправление критической уязвимости - утечка данных чата между сессиями

## Дата: 16.06.2025

## Проблема
При смене пользователей (logout -> login с другим аккаунтом) данные чата предыдущего пользователя оставались видимыми новому пользователю. Это критическая проблема безопасности для публичных компьютеров (интернет-кафе).

## Решение

### 1. Добавлен action clearAllData в chatSlice
Файл: `/frontend/svetu/src/store/slices/chatSlice.ts`
```typescript
clearAllData: (state) => {
  // Закрываем WebSocket если он открыт
  if (state.ws) {
    state.ws.close();
  }
  // Возвращаем состояние к начальному
  Object.assign(state, initialState);
}
```

### 2. Создан компонент AuthStateManager
Файл: `/frontend/svetu/src/components/AuthStateManager.tsx`
- Отслеживает изменение пользователя
- При смене пользователя очищает все данные чата
- Очищает localStorage и sessionStorage (сохраняя важные данные как locale)

### 3. Обновлен AuthContext
Файл: `/frontend/svetu/src/contexts/AuthContext.tsx`
- При logout очищает localStorage и sessionStorage
- Сохраняет важные данные (locale) перед очисткой

### 4. Обновлен useChat hook
Файл: `/frontend/svetu/src/hooks/useChat.ts`
- Добавлен метод clearChatData для очистки данных

## Результат
- При смене пользователя все данные чата очищаются
- Redux store возвращается к начальному состоянию
- WebSocket закрывается и переподключается для нового пользователя
- localStorage и sessionStorage очищаются, сохраняя только необходимые данные
- Данные предыдущего пользователя больше не доступны новому пользователю