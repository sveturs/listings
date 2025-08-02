# Анализ различий между PR #110 и текущей реализацией

## Что сломалось

В PR #110 система вариативных атрибутов работала корректно для женской одежды. При анализе выявлено, что основная проблема заключается в изменении логики фильтрации атрибутов.

## Ключевые различия

### PR #110 (рабочая версия)
```typescript
// Простая проверка по ключевым словам
const isVariantAttribute = [
  'color',
  'size',
  'цвет',
  'размер',
  'boja',
  'veličina',
].some((keyword) => name.includes(keyword));
```

### Текущая версия (до исправления)
```typescript
// Строгое сопоставление по именам из API
const variantAttributeNames = availableVariantAttributes.map((attr) =>
  attr.name.toLowerCase()
);
const isVariantAttribute = variantAttributeNames.includes(name);
```

## Выявленная проблема

Несоответствие названий атрибутов между таблицами:
- В `category_attributes`: атрибут называется `ram`
- В `product_variant_attributes`: атрибут называется `memory`

Это приводило к тому, что атрибуты `color` и `storage` не определялись как вариативные для категории smartphones.

## Решение

Реализован комбинированный подход с маппингом названий:

```typescript
const attributeNameMapping: Record<string, string> = {
  'ram': 'memory',  // В БД атрибут называется ram, а в вариантах - memory
  'color': 'color',
  'storage': 'storage',
  // ... другие атрибуты
};

// Проверка с учетом маппинга
let isVariantAttribute = variantAttributeNames.includes(attrName);

if (!isVariantAttribute) {
  const mappedName = attributeNameMapping[attrName];
  if (mappedName) {
    isVariantAttribute = variantAttributeNames.includes(mappedName);
  }
}

// Fallback на проверку по ключевым словам
if (!isVariantAttribute) {
  isVariantAttribute = [
    'color', 'size', 'цвет', 'размер', 'boja', 'veličina',
    'memory', 'storage', 'ram'
  ].some((keyword) => attrName.includes(keyword));
}
```

## Результат

После применения исправления система корректно определяет все вариативные атрибуты для категории smartphones:
- ✅ Color
- ✅ Memory (RAM) 
- ✅ Storage

## Рекомендации на будущее

1. **Синхронизировать названия атрибутов** между таблицами через миграцию
2. **Добавить unit тесты** для проверки корректности маппинга атрибутов
3. **Создать централизованную конфигурацию** для управления соответствием атрибутов
4. **Добавить валидацию** при создании новых атрибутов для предотвращения несоответствий