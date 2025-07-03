# Паспорт компонента: CreateStorefrontContext

## Метаданные
- **Путь**: `/frontend/svetu/src/contexts/CreateStorefrontContext.tsx`
- **Роль**: Управление состоянием создания витрины
- **Тип**: React Context с состоянием и логикой
- **Размер**: 197 строк

## Назначение
Централизованное управление данными формы создания витрины. Предоставляет состояние, методы обновления и отправки данных для всех шагов мастера создания витрины.

## Props
```typescript
interface CreateStorefrontProviderProps {
  children: ReactNode;  // Дочерние компоненты
}
```

## Зависимости
### Внешние
- `react` - хуки и контекст
- `next-intl` - локализация

### Внутренние
- `@/utils/toast` - уведомления
- `@/services/storefrontApi` - API сервис
- `@/types/storefront` - типы данных

## Управление состоянием
### Состояние контекста
- `formData` - все данные формы создания
- `isSubmitting` - статус отправки

### Методы контекста
- `updateFormData` - частичное обновление данных
- `resetFormData` - сброс к начальным значениям
- `submitStorefront` - отправка на сервер

## Бизнес-логика
1. **Структура данных формы**:
   ```typescript
   StorefrontFormData {
     // Основная информация
     name, slug, description, businessType
     
     // Бизнес-детали
     registrationNumber?, taxNumber?, vatNumber?,
     website?, email?, phone?
     
     // Локация
     address, city, postalCode, country,
     latitude?, longitude?
     
     // Часы работы
     businessHours: Array<{
       dayOfWeek, openTime, closeTime, isClosed
     }>
     
     // Оплата и доставка
     paymentMethods: string[]
     deliveryOptions: Array<{
       providerName, deliveryTimeMinutes,
       deliveryCostRSD, freeDeliveryThresholdRSD?
     }>
     
     // Персонал
     staff: Array<{
       email, role, canManageProducts,
       canManageOrders, canManageSettings
     }>
   }
   ```

2. **Начальные значения**:
   - Тип бизнеса: retail
   - Страна: RS (Сербия)
   - Часы работы: пн-пт 9-18, сб 9-15, вс закрыто
   - Пустые массивы для методов оплаты, доставки, персонала

3. **Трансформация для API**:
   - Группировка location данных
   - Вложение настроек в объект settings
   - Преобразование в StorefrontCreateDTO

4. **Обработка результата**:
   - Успех: сброс формы, возврат ID
   - Ошибка: показ уведомления, возврат ошибки

## UI структура
```
<CreateStorefrontProvider>
  └── {children} // Все шаги мастера имеют доступ к контексту
      ├── BasicInfoStep
      ├── BusinessDetailsStep
      ├── LocationStep
      ├── BusinessHoursStep
      ├── PaymentDeliveryStep
      ├── StaffSetupStep
      └── PreviewPublishStep
```

## Примеры использования
```typescript
// Провайдер в корне мастера
<CreateStorefrontProvider>
  <CreateStorefrontWizard />
</CreateStorefrontProvider>

// Использование в компоненте шага
const { formData, updateFormData } = useCreateStorefrontContext();

// Обновление данных
updateFormData({ 
  name: "Моя витрина",
  slug: "my-store" 
});

// Отправка формы
const result = await submitStorefront();
if (result.success) {
  router.push(`/storefronts/${result.storefrontId}`);
}
```

## Известные особенности
### Позитивные
- Централизованное управление состоянием
- Типизированные данные формы
- Автоматическая трансформация для API
- Обработка ошибок с уведомлениями
- Мемоизация методов для оптимизации

### Технический долг
- Нет валидации данных перед отправкой
- Нет сохранения черновика в localStorage
- Простая обработка ошибок без деталей

### Возможные улучшения
- Добавить валидацию на уровне контекста
- Сохранение черновика при навигации
- Поддержка отмены изменений (undo)
- Прогресс отправки для больших данных
- Детальная обработка ошибок API
- Поддержка редактирования существующей витрины