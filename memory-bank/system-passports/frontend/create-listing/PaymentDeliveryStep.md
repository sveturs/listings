# PaymentDeliveryStep.tsx

## Метаданные
- **Путь**: `frontend/svetu/src/components/create-listing/steps/PaymentDeliveryStep.tsx`
- **Роль**: UI компонент настройки оплаты и доставки
- **Тип**: Step компонент мастера создания объявления
- **Размер**: 337 строк

## Назначение
Шестой шаг мастера создания объявления. Настраивает способы оплаты, службы доставки и дополнительные опции для региональной специфики торговли.

## Props структура
```typescript
interface PaymentDeliveryStepProps {
  onNext: () => void;   // Переход к следующему шагу
  onBack: () => void;   // Возврат к предыдущему шагу
}
```

## Зависимости
### Внешние зависимости
- `react` - управление состоянием
- `next-intl` - интернационализация

### Внутренние зависимости
- Пока не интегрирован с CreateListingContext (TODO)

## Управление состоянием
### Локальное состояние
```typescript
const [formData, setFormData] = useState({
  paymentMethods: ['cod'],          // Способы оплаты: cod, bank_transfer, cash
  codPrice: 250,                    // Стоимость наложенного платежа в РСД
  personalMeeting: true,            // Личная передача товара
  deliveryOptions: string[],        // Выбранные службы доставки
  negotiablePrice: boolean,         // Возможность торга
  bundleDeals: boolean,            // Комплектные предложения
});
```

### Региональные данные
```typescript
// TODO: Загружать список курьерских служб и тарифов из API
// ВРЕМЕННОЕ РЕШЕНИЕ: Хардкодные службы доставки для демонстрации
const deliveryServices = [
  { id: 'post_srbije', name: 'Пошта Србије', fee: 250 },
  { id: 'aks', name: 'AKS Express', fee: 300 },
  { id: 'bex', name: 'BEX Express', fee: 280 },
  { id: 'city_express', name: 'City Express', fee: 320 },
];
```

## Бизнес-логика
### Способы оплаты
- **COD (Cash on Delivery)**: наложенный платеж с фиксированной комиссией
- **Cash**: наличный расчет при встрече
- **Bank Transfer**: банковский перевод (предоплата)
- **Multiple selection**: возможность выбора нескольких способов

### Калькулятор стоимости
```typescript
const calculateCODTotal = (price: number) => {
  return price + formData.codPrice; // Товар + комиссия за наложенный платеж
};
```

### Валидация
- Минимум один способ оплаты должен быть выбран
- COD калькулятор показывается только при выборе наложенного платежа

## UI структура
### Способы оплаты
- **Checkbox cards**: визуальные карточки для каждого способа
- **Popular badges**: для популярных опций (COD)
- **Dynamic content**: COD калькулятор показывается условно

### Службы доставки
- **List selection**: checkbox для каждой службы
- **Pricing display**: стоимость доставки для каждой службы
- **Regional adaptation**: локальные сербские службы

### Дополнительные опции
- **Toggle switches**: для negotiable price и bundle deals
- **Compact layout**: toggle элементы в строку

## Примеры использования
```tsx
// В мастере создания объявления
<PaymentDeliveryStep
  onNext={() => setCurrentStep(6)}
  onBack={() => setCurrentStep(4)}
/>

// Обработка выбранных опций
const paymentMethods = formData.paymentMethods; // ['cod', 'cash']
const deliveryServices = formData.deliveryOptions; // ['post_srbije', 'aks']
```

## Известные особенности
### Позитивные
- ✅ Региональная адаптация (местные службы доставки)
- ✅ COD калькулятор для прозрачности стоимости
- ✅ Multiple selection для гибкости
- ✅ Visual feedback при выборе опций
- ✅ Локализация всех текстов

### Технический долг
- ⚠️ Не интегрирован с CreateListingContext (данные не сохраняются)
- ⚠️ Hardcoded список курьерских служб и тарифов
- ⚠️ Hardcoded стоимость COD (250 РСД)
- ⚠️ Отсутствует валидация логики (например, COD + доставка)
- ⚠️ Эмодзи иконки вместо SVG

### Потенциальные улучшения
- Интегрировать с глобальным контекстом создания объявления
- Создать API для динамической загрузки служб доставки
- Добавить калькулятор общей стоимости (товар + доставка + COD)
- Реализовать зависимости между опциями (например, COD требует доставку)
- Добавить региональные предпочтения (автовыбор популярных служб)
- Улучшить UX для comparison тарифов доставки