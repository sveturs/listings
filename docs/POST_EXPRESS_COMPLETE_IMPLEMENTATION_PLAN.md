# 🎯 Post Express - План полной реализации БЕЗ обращения к Николе

**Дата:** 14 октября 2025
**Статус:** 🟢 ВСЁ УЖЕ ЕСТЬ В КОДЕ!
**Открытие:** Transaction IDs УЖЕ реализованы в client.go!

---

## 🎉 КРИТИЧЕСКОЕ ОТКРЫТИЕ!

**ВСЕ Transaction IDs УЖЕ ЕСТЬ В КОДЕ!**

Анализ `/backend/internal/proj/postexpress/service/client.go` показал:

| Transaction ID | Метод | Строка | Описание | Статус |
|---------------|-------|--------|----------|--------|
| **3** | `GetLocations()` | 318 | Поиск населённых пунктов | ✅ РЕАЛИЗОВАН |
| **10** | `GetOffices()` | 360 | Список отделений | ✅ РЕАЛИЗОВАН |
| **15** | `GetShipmentStatus()` | 444 | Отслеживание | ✅ РЕАЛИЗОВАН |
| **20** | `PrintLabel()` | 484 | Печать этикетки | ✅ РЕАЛИЗОВАН |
| **25** | `CancelShipment()` | 529 | Отмена отправления | ✅ РЕАЛИЗОВАН |
| **73** | `CreateShipmentViaManifest()` | - | Создание манифеста | ✅ РАБОТАЕТ |

**ВСЁ УЖЕ ГОТОВО! Нужно только подключить к handlers и service layer!**

---

## 📋 Что РЕАЛЬНО нужно сделать

### Phase 1: Подключение существующих методов WSP к Service Layer

#### 1.1. Tracking API (Transaction 15) ✅ УЖЕ ЕСТЬ

**WSP Client метод:** `GetShipmentStatus(ctx, trackingNumber)` - строка 434

**Что нужно:**
1. Обновить service метод `TrackShipment()` чтобы вызывать WSP client
2. Парсить ответ и мапить в `TrackingInfo`
3. Протестировать с реальным tracking number

**Код:**
```go
// В service.go
func (s *ServiceImpl) TrackShipment(ctx context.Context, trackingNumber string) ([]models.TrackingEvent, error) {
    // Вызываем WSP client
    wspResp, err := s.wspClient.GetShipmentStatus(ctx, trackingNumber)
    if err != nil {
        return nil, fmt.Errorf("WSP tracking failed: %w", err)
    }

    // Маппинг в TrackingEvent
    events := make([]models.TrackingEvent, 0, len(wspResp.Events))
    for _, e := range wspResp.Events {
        events = append(events, models.TrackingEvent{
            Timestamp:   e.Timestamp,
            Status:      e.Status,
            Description: e.Description,
            Location:    e.Location,
        })
    }

    return events, nil
}
```

**Приоритет:** 🔴 ВЫСОКИЙ
**Сложность:** ⭐ ЛЁГКО (WSP метод готов!)

---

#### 1.2. Cancel API (Transaction 25) ✅ УЖЕ ЕСТЬ

**WSP Client метод:** `CancelShipment(ctx, shipmentID)` - строка 518

**Что нужно:**
1. Обновить service метод `CancelShipment()` чтобы вызывать WSP client
2. Обработать ответ
3. Обновить статус в БД

**Код:**
```go
// В service.go
func (s *ServiceImpl) CancelShipment(ctx context.Context, id int, reason string) error {
    // Получить shipment из БД
    shipment, err := s.storage.GetShipment(ctx, id)
    if err != nil {
        return fmt.Errorf("shipment not found: %w", err)
    }

    // Вызвать WSP cancel
    if err := s.wspClient.CancelShipment(ctx, shipment.PostExpressID); err != nil {
        return fmt.Errorf("WSP cancel failed: %w", err)
    }

    // Обновить статус в БД
    return s.storage.UpdateShipmentStatus(ctx, id, models.StatusCancelled)
}
```

**Приоритет:** 🟡 СРЕДНИЙ
**Сложность:** ⭐ ЛЁГКО (WSP метод готов!)

---

#### 1.3. Label Printing (Transaction 20) ✅ УЖЕ ЕСТЬ

**WSP Client метод:** `PrintLabel(ctx, shipmentID)` - строка 473

**Что нужно:**
1. Обновить service метод `GetShipmentLabel()` чтобы вызывать WSP client
2. Декодировать base64 PDF
3. Вернуть PDF bytes

**Код:**
```go
// В service.go
func (s *ServiceImpl) GetShipmentLabel(ctx context.Context, id int) ([]byte, error) {
    // Получить shipment из БД
    shipment, err := s.storage.GetShipment(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("shipment not found: %w", err)
    }

    // Вызвать WSP PrintLabel
    pdfContent, err := s.wspClient.PrintLabel(ctx, shipment.PostExpressID)
    if err != nil {
        return nil, fmt.Errorf("WSP print label failed: %w", err)
    }

    // Декодировать base64 если нужно
    // TODO: проверить формат ответа от WSP

    return pdfContent, nil
}
```

**Приоритет:** 🔴 ВЫСОКИЙ
**Сложность:** ⭐⭐ СРЕДНЕ (нужна проверка формата PDF)

---

#### 1.4. Office Locator (Transaction 10) ✅ УЖЕ ЕСТЬ

**WSP Client метод:** `GetOffices(ctx, locationID)` - строка 348

**Что нужно:**
1. Обновить service метод `SyncOffices()` чтобы вызывать WSP client
2. Сохранить офисы в БД
3. Обновить timestamp последней синхронизации

**Код:**
```go
// В service.go
func (s *ServiceImpl) SyncOffices(ctx context.Context) error {
    // Получить все населённые пункты из БД
    locations, err := s.storage.GetAllLocations(ctx)
    if err != nil {
        return fmt.Errorf("failed to get locations: %w", err)
    }

    // Для каждого населённого пункта получить офисы
    for _, loc := range locations {
        wspOffices, err := s.wspClient.GetOffices(ctx, loc.PostExpressID)
        if err != nil {
            s.logger.Error("Failed to get offices for location %d: %v", loc.ID, err)
            continue
        }

        // Сохранить офисы в БД
        for _, wspOffice := range wspOffices {
            office := &models.PostExpressOffice{
                Code:         wspOffice.Code,
                Name:         wspOffice.Name,
                Address:      wspOffice.Address,
                LocationID:   loc.ID,
                // ... other fields
            }

            if err := s.storage.UpsertOffice(ctx, office); err != nil {
                s.logger.Error("Failed to upsert office: %v", err)
            }
        }
    }

    return nil
}
```

**Приоритет:** 🟡 СРЕДНИЙ
**Сложность:** ⭐⭐ СРЕДНЕ (нужен цикл по локациям)

---

#### 1.5. Location Search (Transaction 3) ✅ УЖЕ ЕСТЬ

**WSP Client метод:** `GetLocations(ctx, search)` - строка 303

**Что нужно:**
1. Обновить service метод `SyncLocations()` чтобы вызывать WSP client
2. Сохранить населённые пункты в БД

**Код:**
```go
// В service.go
func (s *ServiceImpl) SyncLocations(ctx context.Context) error {
    // Получить ВСЕ населённые пункты (пустой поиск возвращает все)
    wspLocations, err := s.wspClient.GetLocations(ctx, "")
    if err != nil {
        return fmt.Errorf("failed to get locations from WSP: %w", err)
    }

    s.logger.Info("Fetched %d locations from Post Express", len(wspLocations))

    // Сохранить в БД
    for _, wspLoc := range wspLocations {
        location := &models.PostExpressLocation{
            Name:           wspLoc.Name,
            PostalCode:     wspLoc.PostalCode,
            PostExpressID:  wspLoc.ID,
            // ... other fields
        }

        if err := s.storage.UpsertLocation(ctx, location); err != nil {
            s.logger.Error("Failed to upsert location: %v", err)
        }
    }

    return nil
}
```

**Приоритет:** 🟡 СРЕДНИЙ
**Сложность:** ⭐ ЛЁГКО (прямой вызов WSP)

---

### Phase 2: Rate Calculator (нужна логика расчёта)

**Проблема:** Transaction ID для Rate Calculator НЕ найден в коде

**Возможные решения:**

#### Option A: Использовать локальные тарифы (ТЕКУЩЕЕ РЕШЕНИЕ)
- ✅ Уже работает
- ⚠️ Не real-time
- ⚠️ Нужно ручное обновление тарифов

#### Option B: Reverse-engineer из документации
- Поискать в WSP Help документации
- Попробовать разные Transaction IDs (50-60 диапазон обычно для тарифов)

#### Option C: Вычислять из создания манифеста
- Создать "пробный" манифест без реальной отправки
- Извлечь стоимость из ответа
- Не создавать реальную посылку

**Рекомендация:** Оставить локальные тарифы пока, это низкоприоритетная задача

**Приоритет:** 🟡 СРЕДНИЙ (работает локально)
**Сложность:** ⭐⭐⭐ СЛОЖНО (нет Transaction ID)

---

## 🚀 План выполнения (приоритизированный)

### Week 1: Критический функционал

#### Day 1-2: Tracking API
1. ✅ Изучить WSP response от `GetShipmentStatus()` (Transaction 15)
2. ✅ Обновить `service.go` метод `TrackShipment()`
3. ✅ Протестировать с реальным tracking number
4. ✅ Добавить UI компонент отслеживания

#### Day 3-4: Label Printing
1. ✅ Изучить WSP response от `PrintLabel()` (Transaction 20)
2. ✅ Проверить формат PDF (base64 или binary?)
3. ✅ Обновить `service.go` метод `GetShipmentLabel()`
4. ✅ Добавить UI кнопку "Печать этикетки"

#### Day 5: Cancel API
1. ✅ Обновить `service.go` метод `CancelShipment()`
2. ✅ Протестировать отмену созданной посылки
3. ✅ Добавить UI кнопку отмены

### Week 2: Синхронизация данных

#### Day 1-2: Location Sync
1. ✅ Обновить `service.go` метод `SyncLocations()`
2. ✅ Создать cron job для ежедневной синхронизации
3. ✅ Протестировать полный цикл синхронизации

#### Day 3-4: Office Sync
1. ✅ Обновить `service.go` метод `SyncOffices()`
2. ✅ Интегрировать с Location Sync
3. ✅ Протестировать получение офисов по городам

#### Day 5: UI для выбора офисов
1. ✅ Добавить карту с отделениями
2. ✅ Фильтр по городу/индексу
3. ✅ Выбор пакетомата

### Week 3: Дополнительные фичи

#### Day 1-2: Return Shipments UI
1. ✅ Добавить сценарий возврата на тестовую страницу
2. ✅ Документировать процесс возврата

#### Day 3-5: Bulk Operations UI
1. ✅ UI для загрузки CSV/Excel
2. ✅ Парсинг файлов
3. ✅ Массовое создание манифестов

---

## 📊 Сводная таблица (ОБНОВЛЁННАЯ)

| Функция | WSP Client | Transaction ID | Service | Handler | UI | Статус |
|---------|-----------|----------------|---------|---------|----|----|
| **Create Manifest** | ✅ | 73 | ✅ | ✅ | ✅ | ✅ DONE |
| **Tracking** | ✅ | 15 | 🔄 TODO | ✅ | 🔄 TODO | 🟡 PARTIAL |
| **Cancel** | ✅ | 25 | 🔄 TODO | ✅ | 🔄 TODO | 🟡 PARTIAL |
| **Label Printing** | ✅ | 20 | 🔄 TODO | ✅ | 🔄 TODO | 🟡 PARTIAL |
| **Office Locator** | ✅ | 10 | 🔄 TODO | ✅ | 🔄 TODO | 🟡 PARTIAL |
| **Location Search** | ✅ | 3 | 🔄 TODO | ✅ | 🔄 TODO | 🟡 PARTIAL |
| **Rate Calculator** | ❌ | ❓ | ✅ (local) | ✅ | 🔄 TODO | 🟡 LOCAL |
| **Return Shipments** | ✅ | 73 | ✅ | ✅ | 🔄 TODO | 🟡 PARTIAL |
| **Warehouse Pickup** | N/A | N/A | ✅ | ✅ | ❌ | ✅ LOCAL |
| **Statistics** | N/A | N/A | ✅ | ✅ | ❌ | ✅ LOCAL |

**Легенда:**
- ✅ Полностью реализовано и работает
- 🔄 TODO - Нужно реализовать
- 🟡 PARTIAL - Частично реализовано
- ❌ Не реализовано
- ❓ Transaction ID неизвестен

---

## 🎯 Immediate Action Plan

### 1. Протестировать существующие WSP методы

**Создать тестовый скрипт:**

```bash
cd /data/hostel-booking-system/backend
cat > scripts/test_wsp_methods.go <<'EOF'
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "backend/internal/proj/postexpress/service"
    "backend/pkg/logger"
)

func main() {
    logger := logger.NewLogger()

    config := &service.WSPConfig{
        Endpoint:   "http://212.62.32.201/WspWebApi/transakcija",
        Username:   "b2b@svetu.rs",
        Password:   "Sv5et@U!",
        Language:   "sr-Latn-RS",
        DeviceType: "2",
        Timeout:    30 * time.Second,
        MaxRetries: 3,
        RetryDelay: 2 * time.Second,
        TestMode:   true,
        PartnerID:  10109,
    }

    client := service.NewWSPClient(config, logger)
    ctx := context.Background()

    // Test 1: GetLocations (Transaction 3)
    fmt.Println("\n=== Test 1: GetLocations ===")
    locations, err := client.GetLocations(ctx, "Beograd")
    if err != nil {
        log.Printf("GetLocations failed: %v", err)
    } else {
        log.Printf("✅ Found %d locations", len(locations))
        if len(locations) > 0 {
            log.Printf("First location: %s (ID: %d)", locations[0].Name, locations[0].ID)
        }
    }

    // Test 2: GetOffices (Transaction 10)
    if len(locations) > 0 {
        fmt.Println("\n=== Test 2: GetOffices ===")
        offices, err := client.GetOffices(ctx, locations[0].ID)
        if err != nil {
            log.Printf("GetOffices failed: %v", err)
        } else {
            log.Printf("✅ Found %d offices", len(offices))
            if len(offices) > 0 {
                log.Printf("First office: %s - %s", offices[0].Code, offices[0].Name)
            }
        }
    }

    // Test 3: GetShipmentStatus (Transaction 15)
    // Нужен реальный tracking number из предыдущих тестов
    fmt.Println("\n=== Test 3: GetShipmentStatus ===")
    fmt.Println("⚠️ Нужен реальный tracking number для тестирования")

    // Test 4: PrintLabel (Transaction 20)
    fmt.Println("\n=== Test 4: PrintLabel ===")
    fmt.Println("⚠️ Нужен реальный shipment ID для тестирования")

    // Test 5: CancelShipment (Transaction 25)
    fmt.Println("\n=== Test 5: CancelShipment ===")
    fmt.Println("⚠️ Нужен реальный shipment ID для тестирования")

    fmt.Println("\n=== Tests Complete ===")
}
EOF

go run scripts/test_wsp_methods.go
```

### 2. Обновить Service Layer

Обновить файл `service.go` чтобы вызывать WSP client методы вместо заглушек.

### 3. Протестировать через frontend

Использовать тестовую страницу `/admin/postexpress/test` для проверки всех функций.

---

## 📝 Выводы

### ✅ Хорошие новости:
1. **ВСЕ Transaction IDs УЖЕ ЕСТЬ В КОДЕ!**
2. WSP Client методы полностью реализованы
3. Handlers и endpoints готовы
4. Осталось только подключить Service layer

### 🔄 Что нужно сделать:
1. Обновить Service методы (5 методов)
2. Протестировать WSP responses
3. Добавить UI компоненты (tracking, cancel, label)

### 🎯 Приоритет:
1. **🔴 HIGH**: Tracking, Label Printing
2. **🟡 MED**: Cancel, Office Sync, Location Sync
3. **🟢 LOW**: Rate Calculator (работает локально), Webhooks

### ⏱️ Оценка времени:
- Phase 1 (Tracking + Label + Cancel): 2-3 дня
- Phase 2 (Sync): 2-3 дня
- Phase 3 (UI polish): 2-3 дня
- **Итого: 6-9 дней до полной реализации**

---

**Created:** 14 октября 2025
**Status:** 🟢 READY TO IMPLEMENT
**Next:** Начать с Tracking API (Transaction 15)
