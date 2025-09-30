# 🗄️ База данных - Правила работы

## ⚠️ КРИТИЧЕСКИ ВАЖНОЕ ПРАВИЛО

**ВСЕ изменения структуры и данных в базе данных должны производиться ТОЛЬКО через миграции!**

---

## ❌ ЗАПРЕЩЕНО

```bash
# НЕ делай так:
psql -c "ALTER TABLE users ADD COLUMN email VARCHAR(255);"
psql -c "UPDATE products SET price = price * 1.1;"
psql -c "INSERT INTO categories VALUES (...);"
```

## ✅ ПРАВИЛЬНО

```bash
# 1. Изучи БД
DATABASE_URL=postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable

# 2. Создай миграцию
cd /data/hostel-booking-system/backend/migrations
# Создай файлы: 000XXX_description.up.sql и 000XXX_description.down.sql

# 3. Примени миграцию
cd /data/hostel-booking-system/backend
./migrator up
```

---

## 📝 Создание миграций

### Структура:
```
backend/migrations/
├── 000001_initial.up.sql      # Применение
├── 000001_initial.down.sql    # Откат
├── 000002_add_users.up.sql
├── 000002_add_users.down.sql
└── ...
```

### Правила именования:
- `NNNNNN_description.up.sql` - применение
- `NNNNNN_description.down.sql` - откат
- NNNNNN - порядковый номер (6 цифр)

### Пример миграции:

**000042_add_email_to_users.up.sql:**
```sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);
CREATE INDEX idx_users_email ON users(email);
```

**000042_add_email_to_users.down.sql:**
```sql
DROP INDEX IF EXISTS idx_users_email;
ALTER TABLE users DROP COLUMN IF EXISTS email;
```

---

## 🔧 Работа с мигратором

### Команды:
```bash
cd /data/hostel-booking-system/backend

# Применить все миграции
./migrator up

# Откатить последнюю
./migrator down

# Применить с фикстурами (тестовые данные)
./migrator migrate --with-fixtures

# Только фикстуры
./migrator migrate --only-fixtures
```

---

## 🔗 Подключение к БД

### Быстрое подключение:
```bash
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
```

### Полезные команды:
```sql
-- Список таблиц
\dt

-- Структура таблицы
\d table_name

-- Текущее подключение
\conninfo

-- Размер БД
SELECT pg_size_pretty(pg_database_size('svetubd'));

-- Активные подключения
SELECT COUNT(*) FROM pg_stat_activity;
```

---

## ⚡ Troubleshooting

### "too many clients already"
```bash
# Проверить количество подключений
psql "..." -c "SELECT COUNT(*) FROM pg_stat_activity;"

# Если > 90 - перезапусти PostgreSQL
sudo systemctl restart postgresql

# Останови все backend процессы
/home/dim/.local/bin/kill-port-3000.sh
screen -ls | grep backend | awk '{print $1}' | xargs -I {} screen -S {} -X quit
```

---

## 📚 См. также

- [Правила миграций](.ai/migrations.md)
- [Backend правила](.ai/backend.md)
