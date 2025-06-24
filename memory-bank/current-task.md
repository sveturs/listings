# Текущая задача: Исправление проблемы с переводом loadMore

## Статус: ЗАВЕРШЕНА
## Дата начала: 2025-06-24
## Последнее обновление: 2025-06-24

## Описание проблемы:
В компоненте InfiniteScrollTrigger кнопка "Загрузить еще" отображалась как `common.loadMore` вместо правильного перевода.

## Решение:

### 1. Изменен подход к передаче переводов:
- Компонент InfiniteScrollTrigger теперь принимает `loadMoreText` как prop
- Убрано использование `useTranslations` внутри компонента
- Переводы передаются из родительских компонентов

### 2. Обновлены файлы:
- `/frontend/svetu/src/components/common/InfiniteScrollTrigger.tsx` - добавлен prop loadMoreText
- `/frontend/svetu/src/components/products/ProductList.tsx` - передается перевод через prop
- `/frontend/svetu/src/messages/ru.json` - добавлен перевод в storefronts.products.loadMore
- `/frontend/svetu/src/messages/en.json` - добавлен перевод в storefronts.products.loadMore
- Обновлены все другие компоненты, использующие InfiniteScrollTrigger

### 3. Результат:
✅ Кнопка корректно отображает "Загрузить еще" на русском языке
✅ Ошибки IntlError больше не появляются в консоли
✅ Решение работает для всех компонентов с InfiniteScrollTrigger

## Связанные файлы:
- Предыдущая задача (bulk operations): `/memory-bank/analysis/bulk-operations-troubleshooting.md`
- Анализ bulk operations: `/memory-bank/analysis/bulk-operations-and-security-analysis.md`