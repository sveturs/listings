# Отчет по анализу переводов фронтенда

## 📊 Общая статистика

- **Английские переводы (en.json)**: 2634 ключа
- **Русские переводы (ru.json)**: 2612 ключей
- **Разница**: 22 ключа отсутствует в русском переводе
- **Файлов использующих t()**: 297

## 🚨 Критические проблемы

### 1. Hardcoded русский текст

Найдено множество файлов с хардкодным русским текстом, который должен быть вынесен в переводы:

#### `/src/components/marketplace/listing/SellerInfo.tsx`

```tsx
// Строки 76, 106, 110, 123, 131, 146, 158, 184, 207, 228, 236, 256, 263, 270, 295
{
  locale === 'ru' ? 'Продавец' : 'Seller';
}
{
  locale === 'ru' ? 'Нет отзывов' : 'No reviews yet';
}
{
  locale === 'ru' ? 'На сайте с' : 'Member since';
}
{
  locale === 'ru' ? 'Процент ответов' : 'Response rate';
}
{
  locale === 'ru' ? 'Время ответа' : 'Response time';
}
// ... и другие
```

#### `/src/components/reviews/RatingInput.tsx`

```tsx
// Строки 93-97
const ratingLabels = {
  1: { text: 'Ужасно', color: 'text-error' },
  2: { text: 'Плохо', color: 'text-warning' },
  3: { text: 'Нормально', color: 'text-info' },
  4: { text: 'Хорошо', color: 'text-success' },
  5: { text: 'Отлично', color: 'text-success' },
};
```

#### `/src/components/IconPicker.tsx`

```tsx
// Строки 14, 39, 64, 89, 114, 139, 164, 189, 214, 239, 264, 289, 318
const iconCategories = [
  { name: 'Транспорт', icons: [...] },
  { name: 'Электроника', icons: [...] },
  { name: 'Дом и быт', icons: [...] },
  // ... и другие
];
placeholder = 'Выберите иконку';
```

#### `/src/components/marketplace/listing/ListingActions.tsx`

```tsx
// Множественные условные рендеры
locale === 'ru'
  ? 'Войдите, чтобы добавить в избранное'
  : 'Sign in to add to favorites';
locale === 'ru' ? 'Удалено из избранного' : 'Removed from favorites';
locale === 'ru' ? 'Добавлено в избранное' : 'Added to favorites';
// ... и другие
```

### 2. Hardcoded английский текст

Найдено много файлов с hardcoded английским текстом:

#### `/src/contexts/AuthContext.tsx`

```tsx
'Failed to parse cached user data, clearing cache:';
'SessionStorage is not available, skipping cache';
'Failed to load session. Please try refreshing the page.';
'Failed to initiate login. Please try again.';
// ... и другие
```

#### Файлы с техническими сообщениями

- `src/lib/api-client.ts` - заголовки HTTP запросов
- `src/services/` - множество файлов с английскими error messages
- `src/components/import/` - сообщения об ошибках импорта

## 🔍 Отсутствующие переводы

### Ключи отсутствующие в ru.json (первые 10):

1. `acceptedPaymentMethods`
2. `active`
3. `address`
4. `all`
5. `banner`
6. `bannerRequirements`
7. `basicInfo`
8. `bulk`
9. `card`
10. `categoriesDescription`

### Ключи отсутствующие в en.json (первые 10):

1. `applyChanges`
2. `aspectRatio`
3. `attributes`
4. `bulkEditDescription`
5. `bulkEditTitle`
6. `cancelOperation`
7. `categoryAttributesDescription`
8. `categorySelected`
9. `categorySelection`
10. `chooseFiles`

## ✅ Рекомендации по исправлению

### 1. Немедленные действия (высокий приоритет)

#### Вынести hardcoded русский текст в переводы:

**SellerInfo.tsx:**

```json
// Добавить в ru.json и en.json
"seller": {
  "title": "Продавец" / "Seller",
  "noReviews": "Нет отзывов" / "No reviews yet",
  "memberSince": "На сайте с" / "Member since",
  "responseRate": "Процент ответов" / "Response rate",
  "responseTime": "Время ответа" / "Response time",
  "verified": "Проверен" / "Verified",
  "experienced": "Опытный продавец" / "Experienced",
  "sendMessage": "Написать сообщение" / "Send Message",
  "showPhone": "Показать телефон" / "Show Phone",
  "allItems": "Все товары продавца" / "All seller items",
  "yourListing": "Это ваше объявление" / "This is your listing",
  "edit": "Редактировать" / "Edit",
  "signInToContact": "Войдите, чтобы связаться с продавцом" / "Sign in to contact seller",
  "signIn": "Войти" / "Sign In",
  "platformProtection": "Все сделки защищены правилами платформы" / "All transactions protected by platform rules"
}
```

**RatingInput.tsx:**

```json
"rating": {
  "labels": {
    "terrible": "Ужасно" / "Terrible",
    "bad": "Плохо" / "Bad",
    "normal": "Нормально" / "Normal",
    "good": "Хорошо" / "Good",
    "excellent": "Отлично" / "Excellent"
  }
}
```

**IconPicker.tsx:**

```json
"iconCategories": {
  "transport": "Транспорт" / "Transport",
  "electronics": "Электроника" / "Electronics",
  "homeAndLife": "Дом и быт" / "Home & Life",
  "clothing": "Одежда" / "Clothing",
  "foodAndDrinks": "Еда и напитки" / "Food & Drinks",
  "sportsAndLeisure": "Спорт и отдых" / "Sports & Leisure",
  "beautyAndHealth": "Красота и здоровье" / "Beauty & Health",
  "booksAndEducation": "Книги и обучение" / "Books & Education",
  "natureAndAnimals": "Природа и животные" / "Nature & Animals",
  "tools": "Инструменты" / "Tools",
  "numbersAndSymbols": "Числа и символы" / "Numbers & Symbols",
  "attributes": "Атрибуты" / "Attributes",
  "placeholder": "Выберите иконку" / "Select an icon"
}
```

### 2. Средний приоритет

#### Заменить условную логику locale === 'ru' на useTranslations:

```tsx
// Вместо:
{
  locale === 'ru' ? 'Русский текст' : 'English text';
}

// Использовать:
const t = useTranslations('appropriate.namespace');
{
  t('key');
}
```

#### Добавить отсутствующие переводы:

- Добавить 22 ключа в ru.json
- Добавить отсутствующие ключи в en.json

### 3. Низкий приоритет

#### Убрать hardcoded английские технические сообщения:

- Перенести error messages в переводы
- Создать namespace для технических сообщений
- Заменить console.log/console.error сообщения на переводы

## 📋 План выполнения

1. **Этап 1**: Исправить SellerInfo.tsx, RatingInput.tsx, IconPicker.tsx
2. **Этап 2**: Исправить ListingActions.tsx и другие компоненты marketplace
3. **Этап 3**: Добавить отсутствующие переводы в файлы messages
4. **Этап 4**: Заменить условную логику на useTranslations
5. **Этап 5**: Убрать hardcoded технические сообщения

## 🎯 Ожидаемый результат

После исправления:

- Полная интернационализация UI компонентов
- Единообразие в использовании переводов
- Удаление условной логики для языков
- Лучшая поддерживаемость кода
- Готовность к добавлению новых языков
