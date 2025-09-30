# 🆘 Troubleshooting - Типичные проблемы и решения

## 🚫 Backend не запускается

### Симптомы:
```
curl: (7) Failed to connect to localhost port 3000
```

### Решение:
```bash
# 1. Проверь порт
netstat -tlnp | grep :3000

# 2. Если занят - останови
/home/dim/.local/bin/kill-port-3000.sh

# 3. Закрой screen сессии
screen -S backend-3000 -X quit

# 4. Запусти заново
screen -dmS backend-3000 bash -c 'go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'

# 5. Проверь логи
tail -f /tmp/backend.log
```

---

## 🎨 Frontend ошибки сборки

### Ошибка: "Module not found"
```bash
# Переустановить зависимости
cd /data/hostel-booking-system/frontend/svetu
rm -rf node_modules package-lock.json
yarn install
```

### Ошибка: "Port 3001 already in use"
```bash
# Правильный способ остановки
/home/dim/.local/bin/kill-port-3001.sh
screen -S frontend-3001 -X quit

# Запуск
/home/dim/.local/bin/start-frontend-screen.sh
```

---

## 🗄️ База данных

### "too many clients already"
```bash
# Проверить подключения
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable" -c "SELECT COUNT(*) FROM pg_stat_activity;"

# Если > 90 - перезапуск PostgreSQL
sudo systemctl restart postgresql

# Остановка всех backend процессов
/home/dim/.local/bin/kill-port-3000.sh
```

### Миграции не применяются
```bash
# Проверить текущую версию
cd /data/hostel-booking-system/backend
./migrator version

# Откатить и применить заново
./migrator down
./migrator up
```

---

## 🔐 JWT токен не работает

### "401 Unauthorized"

**Причины:**
1. Токен истёк (живёт 1 день)
2. Backend не запущен
3. Неправильный формат токена

**Решение:**
```bash
# Получить свежий токен
ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && sed 's|/data/auth_svetu/keys/private.pem|./keys/private.pem|g' scripts/create_admin_jwt.go > /tmp/create_jwt_fixed.go && go run /tmp/create_jwt_fixed.go" > /tmp/jwt_token.txt

# Проверить
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me | jq ".data.email"'
```

---

## 🖼️ Изображения не загружаются

### 400 Bad Request при загрузке изображения

**Причины:**
1. Неправильное имя поля формы
2. Размер файла превышает лимит
3. Неправильный content-type

**Решение:**
```bash
# Проверить размер файла (должен быть < 10MB)
ls -lh /tmp/test_image.jpg

# Правильное имя поля:
# - Для storefront products: "image" (НЕ "images")
# - Для marketplace listings: "images"

# Пример:
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -X POST "http://localhost:3000/api/v1/storefronts/slug/shop/products/123/images" -H "Authorization: Bearer $TOKEN" -F "image=@/tmp/test.jpg"'
```

---

## 🔍 OpenSearch / Поиск

### Товары не находятся в поиске

**Причина:** Индекс OpenSearch не синхронизирован с БД

**Решение:**
```bash
# Полная переиндексация
python3 /data/hostel-booking-system/backend/reindex_full.py

# Проверка количества документов
curl -X GET "http://localhost:9200/marketplace_listings/_count" | jq '.'
```

---

## 🧹 Очистка кэша

### Redis кэш
```bash
docker exec hostel_redis redis-cli FLUSHALL
```

### Next.js кэш
```bash
cd /data/hostel-booking-system/frontend/svetu
rm -rf .next
yarn dev -p 3001
```

---

## 📚 См. также

- [Pre-check Guidelines](CLAUDE_PRE_CHECK_GUIDELINES.md)
- [Database Guidelines](CLAUDE_DATABASE_GUIDELINES.md)
- [Основное руководство](../CLAUDE.md)
