# Сессия: Исправление алгоритма похожих объявлений
**Дата:** 2025-01-07
**Задача:** Исправить показ батареек вместо похожей недвижимости для виллы

## Проблема
Для объявления виллы (ID 31, цена $890,000) в разделе "Похожие объявления" показывались батарейки для телефонов вместо похожей недвижимости.

## Анализ проблемы

### 1. Неправильная привязка к витрине
```sql
-- Проверка в БД показала:
SELECT id, title, category_id, price, storefront_id 
FROM marketplace_listings WHERE id = 31;
-- Результат: storefront_id = 9 (витрина "345")
```

Объявление о вилле было ошибочно привязано к витрине, что приводило к поиску похожих товаров среди storefront products.

### 2. Жесткий фильтр по категории
Алгоритм искал похожие объявления только в той же категории (1300 - дома). При отсутствии других домов в базе, алгоритм переключался на товары витрин.

### 3. Обнуление скора для разных категорий
В `calculateCategoryScore` возвращалось 0.0 для разных категорий, что сильно снижало общий балл похожести.

## Выполненные исправления

### 1. Исправлена привязка к витрине
```sql
UPDATE marketplace_listings SET storefront_id = NULL WHERE id = 31;
```
Выполнена переиндексация OpenSearch.

### 2. Улучшен buildAdvancedSearchParams
```go
// На последней попытке поиск во всех категориях
if tryNumber < 3 {
    params.CategoryID = strconv.Itoa(listing.CategoryID)
    log.Printf("Попытка %d: поиск в категории %d", tryNumber, listing.CategoryID)
} else {
    log.Printf("Попытка %d: поиск во всех категориях", tryNumber)
}
```

### 3. Улучшен алгоритм скоринга

#### Гибкий расчет CategoryScore:
```go
func (sc *SimilarityCalculator) calculateCategoryScore(
    source, target *models.MarketplaceListing,
) float64 {
    if source.CategoryID == target.CategoryID {
        return 1.0
    }
    
    // Категории из одной группы (например, недвижимость)
    sourceGroup := source.CategoryID / 100
    targetGroup := target.CategoryID / 100
    
    if sourceGroup == targetGroup {
        return 0.6  // Вместо 0.0
    }
    
    return 0.1  // Минимальный балл для разных категорий
}
```

#### Адаптивные веса:
```go
if score.CategoryScore >= 1.0 {
    // Та же категория - стандартные веса
    categoryWeight = 0.3
    priceWeight = 0.15
} else if score.CategoryScore >= 0.6 {
    // Категории из одной группы - больше веса цене
    categoryWeight = 0.2
    priceWeight = 0.25
} else {
    // Разные категории - максимальный вес цене
    categoryWeight = 0.1
    priceWeight = 0.35
}
```

## Результаты

### До исправления:
- Baterija Sony Xperia E4G
- Baterija EG za Sony Xperia M2
- Ноутбук Lenovo ThinkPad X1 Carbon
- И другие батарейки...

### После исправления:
1. **Продается дом с садом** - $120,000 (категория 1300)
2. **Луксузан стан на Савском венцу** - $185,000 (категория 1100)
3. **Пентхаус у А Блоку** - $350,000 (категория 1100)
4. Комната для студентов - $200
5. Сдается 2-комнатная квартира - $500

Теперь приоритет отдается объявлениям из схожих категорий (недвижимость) и похожего ценового диапазона.

## Изменённые файлы
1. `/backend/internal/proj/marketplace/service/marketplace.go` - buildAdvancedSearchParams
2. `/backend/internal/proj/marketplace/service/similarity_scoring.go` - улучшен алгоритм скоринга
3. База данных - исправлен storefront_id для объявления 31

## Команды для применения изменений
```bash
# Переиндексация OpenSearch
cd /data/hostel-booking-system/backend
go run ./cmd/reindex/main.go

# Перезапуск backend
/home/dim/.local/bin/kill-port-3000.sh
screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'
```

## Статус: ✅ Завершено успешно