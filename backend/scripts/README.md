# Backend Test Scripts

## Post Express Integration Test

### Описание

Скрипт `test_postexpress.go` выполняет комплексную проверку интеграции с Post Express API.

### Тесты

1. **Получение списка офисов** - запрашивает список отделений в Белграде
2. **Расчет стоимости** - рассчитывает стоимость доставки Белград → Нови Сад
3. **Создание отправления** - создает тестовое отправление с валидными данными
4. **Отслеживание** - отслеживает созданное отправление

### Запуск

```bash
# Из директории backend/scripts
cd /data/hostel-booking-system/backend/scripts
go run test_postexpress.go

# Или из корня backend через make
cd /data/hostel-booking-system/backend
make test-postexpress
```

### Требования

1. Файл `.env` в корне backend с переменными:
   ```bash
   POST_EXPRESS_API_URL=https://wsp-test.posta.rs/api
   POST_EXPRESS_USERNAME=b2b@svetu.rs
   POST_EXPRESS_PASSWORD=Sv5et@U!
   POST_EXPRESS_BRAND=SVETU
   POST_EXPRESS_WAREHOUSE=SVETU
   ```

2. Зависимости установлены (`go mod download`)

### Результаты

- Логи выводятся в консоль с цветным форматированием
- Tracking number последнего созданного отправления сохраняется в `/tmp/postexpress_tracking.txt`
- Полные JSON данные отслеживания выводятся для анализа

### Примечания

- Скрипт использует тестовую среду Post Express (wsp-test.posta.rs)
- Создает реальные отправления в тестовой БД Post Express
- Тестовые данные получателя: "Test Recipient", Белград
- Все номера отправлений имеют префикс `SVETU-TEST-{timestamp}`

### Troubleshooting

**Ошибка "Failed to initialize Post Express service":**
- Проверьте наличие `.env` файла в `backend/`
- Проверьте корректность credentials

**Ошибка "Failed to create shipment":**
- Проверьте доступность тестового API: `curl https://wsp-test.posta.rs/api`
- Проверьте валидность тестовых credentials
- Изучите детали ошибки в логах

**Ошибка "Failed to track shipment":**
- Возможно, отправление еще не зарегистрировано в системе (подождите несколько минут)
- Проверьте tracking number в `/tmp/postexpress_tracking.txt`
