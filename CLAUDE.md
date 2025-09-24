# CLAUDE.md

## 🎯 КРИТИЧЕСКИ ВАЖНЫЕ ПРАВИЛА

1. **НЕ ОСТАВЛЯЙ ТЕХНИЧЕСКИЙ ДОЛГ!**
2. **База данных: ТОЛЬКО через миграции** - см. [Database Guidelines](docs/CLAUDE_DATABASE_GUIDELINES.md)
3. **Коммиты: БЕЗ упоминания Claude** в авторах/соавторах
4. **Процессы: Всегда закрывай старые** перед запуском новых (kill-port скрипты + screen quit)
5. **НЕ ПОРАЖДАЙ РУДИМЕНТЫ** - проверяй наличие функций перед созданием новых
6. **Универсальные решения** - создавай код, который можно использовать везде

---

## 📚 БЫСТРЫЕ ССЫЛКИ НА ДОКУМЕНТАЦИЮ

### 🔧 Основные руководства
- [📋 TodoWrite Guidelines](docs/CLAUDE_TODOWRITE_GUIDELINES.md) - когда и как использовать TodoWrite
- [🔍 Pre-Check Guidelines](docs/CLAUDE_PRE_CHECK_GUIDELINES.md) - проверки перед коммитом
- [🗄️ Database Guidelines](docs/CLAUDE_DATABASE_GUIDELINES.md) - работа с БД через миграции
- [🆘 Troubleshooting](docs/CLAUDE_TROUBLESHOOTING.md) - типичные проблемы и решения
- [🤖 Parallel Agents](docs/CLAUDE_PARALLEL_AGENTS.md) - параллельное выполнение задач

### 📖 Документация по фичам
- [Категории и фильтры](docs/IMPLEMENTATION_CATEGORY_SELECTOR.md)
- [Витрины - статус](docs/STOREFRONTS_STATUS.md)
- [Автомобильный раздел](docs/AUTOMOTIVE_SECTION_STATUS_AND_PLAN.md)
- [Post Express интеграция](docs/POST_EXPRESS_INTEGRATION_COMPLETE.md)
- [Загрузка изображений](docs/IMAGE_UPLOAD_TESTING_GUIDE.md)

---

## 🚀 БЫСТРЫЙ СТАРТ

### Запуск сервисов
```bash
# Backend (порт 3000)
/home/dim/.local/bin/kill-port-3000.sh
screen -dmS backend-3000 bash -c 'go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'

# Frontend (порт 3001)
/home/dim/.local/bin/start-frontend-screen.sh

# Проверка
netstat -tlnp | grep ":3000\|:3001"
```

### Pre-check перед коммитом
```bash
# Backend
cd backend && make format && make lint

# Frontend
cd frontend/svetu && yarn format && yarn lint && yarn build

# Подробнее: docs/CLAUDE_PRE_CHECK_GUIDELINES.md
```

---

## 🔧 ЧАСТО ИСПОЛЬЗУЕМЫЕ КОМАНДЫ

### База данных
```bash
# Подключение
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"

# Применить миграции
cd backend && ./migrator up

# Подробнее: docs/CLAUDE_DATABASE_GUIDELINES.md
```

### OpenSearch переиндексация
```bash
# Полная переиндексация
python3 /data/hostel-booking-system/backend/reindex_full.py

# Проверка
curl -X GET "http://localhost:9200/marketplace_listings/_count" | jq '.'
```

### JWT токен для тестирования
```bash
# Получить токен
ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && sed 's|/data/auth_svetu/keys/private.pem|./keys/private.pem|g' scripts/create_admin_jwt.go > /tmp/create_jwt_fixed.go && go run /tmp/create_jwt_fixed.go" > /tmp/jwt_token.txt

# Проверить (используй bash -c в ZSH!)
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me | jq ".data.email"'

# Примеры использования
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/marketplace/favorites | jq "."'
```

### Очистка кэша
```bash
# Redis
docker exec hostel_redis redis-cli FLUSHALL

# Next.js
cd frontend/svetu && rm -rf .next
```

---

## 🔍 Поиск файлов и контента

**ВАЖНО:** Используй Glob tool, а НЕ bash find/grep!

```markdown
# ✅ Правильно - через Tools
Glob tool: pattern="**/*.go" path="/backend"
Grep tool: pattern="функция" path="/backend"

# ❌ Неправильно - через bash
find /backend -name "*.go"  # НЕ делай так!
grep -r "функция" /backend  # НЕ делай так!
```

Используй bash find/fd ТОЛЬКО для:
- Сложных условий (размер, дата модификации)
- Комбинации с другими командами через pipe

---

## 🔍 Проверка версий для подтверждения изменений

### Frontend версия
Frontend отображает версию из `package.json` в логотипе приложения.
Путь к файлу: `frontend/svetu/package.json`

### Backend версия
Backend возвращает версию при запросе на http://localhost:3000/
Версия хранится в файле: `backend/internal/version/version.go`

**Совет**: При внесении изменений поднимайте патч версию для подтверждения, что работает ваш код.

## 🔍 Проверка версий для подтверждения изменений

### Frontend версия
Frontend отображает версию из `package.json` в логотипе приложения.
Путь к файлу: `frontend/svetu/package.json`

### Backend версия
Backend возвращает версию при запросе на http://localhost:3000/
Версия хранится в файле: `backend/internal/version/version.go`

**Совет**: При внесении изменений поднимайте патч версию для подтверждения, что работает ваш код.

## 🚀 Развертывание на dev.svetu.rs

### Быстрое развертывание:
```bash
./deploy-to-dev.sh
```

Скрипт автоматически:
1. Коммитит и пушит изменения
2. Создаёт дамп БД
3. Загружает на dev сервер
4. Перезапускает сервисы

### Ручное развертывание:
```bash
# 1. Коммит и пуш
git add -A && git commit -m "сообщение" && git push

# 2. Дамп БД
PGPASSWORD=mX3g1XGhMRUZEX3l pg_dump -h localhost -U postgres -d svetubd --no-owner --no-acl --column-inserts --inserts -f /tmp/dump.sql

# 3. На сервере
ssh svetu@svetu.rs
cd /opt/svetu-dev
git pull
docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
docker exec -i svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db < /tmp/dump.sql

# 4. Перезапуск
cd backend && make dev-restart
cd ../frontend/svetu && make dev-restart
```

**Важно:**
- Сервер: https://dev.svetu.rs (frontend), https://devapi.svetu.rs (backend)
- SSH: `ssh svetu@svetu.rs` (без пароля)
- Директория: `/opt/svetu-dev`
- БД в Docker: `svetu-dev_db_1`, база `svetu_dev_db`

---

## 📋 Управление задачами (TodoWrite)

### ✅ Когда ОБЯЗАТЕЛЬНО использовать:
- Multi-step задачи (3+ шага)
- Сложная разработка (новые фичи)
- Рефакторинг (несколько файлов/модулей)
- По явному запросу пользователя

### ❌ Когда НЕ нужен:
- Простые одношаговые задачи
- Read-only исследование
- Тестирование без изменений
- Мелкие правки (1-2 строки)

📚 **Подробнее:** [CLAUDE_TODOWRITE_GUIDELINES.md](docs/CLAUDE_TODOWRITE_GUIDELINES.md)

---

## 📂 Структура проекта

### Backend:
```
backend/
├── cmd/api/              # Точка входа
├── internal/
│   ├── config/           # Конфигурация
│   ├── domain/           # Доменные модели
│   ├── middleware/       # Auth, CORS, Logger
│   ├── proj/             # Бизнес-логика модулей
│   │   ├── marketplace/
│   │   ├── storefronts/
│   │   ├── users/
│   │   └── payments/
│   ├── server/           # HTTP сервер (Fiber)
│   └── storage/          # Репозитории
│       ├── postgres/
│       ├── opensearch/
│       └── minio/
└── migrations/           # SQL миграции
```

### Frontend:
```
frontend/svetu/
├── src/
│   ├── app/[locale]/     # Next.js App Router
│   ├── components/       # React компоненты
│   ├── services/         # API клиенты
│   ├── store/            # Redux Toolkit
│   ├── messages/         # i18n переводы
│   │   ├── en/
│   │   ├── ru/
│   │   └── sr/
│   └── config/           # Конфигурация
```

---

## 🎨 Переводы (i18n)

**ВАЖНО:** Backend возвращает placeholder'ы, frontend переводит их!

```javascript
// Backend возвращает:
{ "error": "storefronts.no_image_file" }

// Frontend переводит через:
t('storefronts.no_image_file')  // → "Файл изображения не найден"
```

Файлы переводов: `frontend/svetu/src/messages/{en,ru,sr}/{module}.json`

**Исправление ошибок переводов:** См. секцию в CLAUDE.md (строка ~230)

---

## 🆘 Troubleshooting

См. [CLAUDE_TROUBLESHOOTING.md](docs/CLAUDE_TROUBLESHOOTING.md) для:
- Backend не запускается
- Frontend ошибки сборки
- База данных: too many connections
- JWT токен не работает
- Изображения не загружаются
- OpenSearch проблемы

---

## 🔐 Безопасность

**Defensive security ONLY:**
- ✅ Security analysis, detection rules
- ✅ Vulnerability explanations
- ✅ Defensive tools
- ❌ Offensive tools, malicious code
- ❌ Credential discovery/harvesting

---

## 📝 Git & Commits

### Правила коммитов:
```bash
# Conventional commits format
feat: add user profile page
fix: resolve login redirect issue
docs: update API documentation
refactor: optimize database queries
```

**ВАЖНО:** НЕ добавляй Claude в авторы!

- **Логирование в Backend**: Используй `backend/internal/logger` для глобального логирования. Если нужно передать логгер как объект - используй `github.com/rs/zerolog` (НЕ используй slog или другие логгеры!)

### Pre-commit hooks:
```bash
# Установка (один раз)
pre-commit install

# Проверка вручную
make pre-commit  # backend
yarn format && yarn lint  # frontend
```

📚 **Подробнее:** [.ai/git.md](.ai/git.md)

---

## 🗄️ Работа с API документацией (Swagger)

### Через JSON MCP (рекомендуется):
```bash
# 1. Запустить HTTP сервер
cd /data/hostel-booking-system/backend/docs && python3 -m http.server 8888 &

# 2. Использовать JSON MCP для поиска
JSON MCP query: "$.paths['/api/v1/auth/login']" from http://localhost:8888/swagger.json
JSON MCP query: "$.definitions['MarketplaceListing']" from http://localhost:8888/swagger.json

# 3. Остановить сервер
pkill -f "python3 -m http.server 8888"
```

**ВСЕГДА** сначала ищи информацию в swagger.json, потом анализируй код!

### Регенерация типов:
```bash
# ТОЛЬКО если изменял swagger аннотации
cd backend && make generate-types
```

---

## 📚 IMPORTANT WORKFLOW RULES

- **Язык общения:** Russian
- **Качество кода:** Pre-check перед коммитом обязательно!
- **Переводы:** Backend - placeholders, Frontend - переводы
- **Зависимости:** Обновляй "Key Dependencies" в .ai/*.md при добавлении
- **Handlers:** Не возвращай реальную ошибку, используй placeholders
- **Swagger:** Используй структуры из pkg/utils/utils.go

📚 **Подробнее:**
- [Frontend правила](.ai/frontend.md)
- [Backend правила](.ai/backend.md)
- [Миграции](.ai/migrations.md)

---

## 🔧 Key Technologies

- **Backend:** Go, Fiber, PostgreSQL, OpenSearch, MinIO, Redis
- **Frontend:** Next.js 15, React 19, TypeScript, Tailwind CSS, Redux Toolkit
- **Infra:** Docker, Nginx, Harbor

---

## 📌 Status Updates

- ✅ **Post Express Integration:** Production ready (waiting for credentials)
- 🚧 **Admin Variant Attributes:** Планируется (см. docs/ADMIN_VARIANT_ATTRIBUTES_EXTENSION_PLAN.md)
- ✅ **Image Upload System:** Полностью функционально (см. docs/IMAGE_UPLOAD_TESTING_GUIDE.md)
- 🚧 **Automotive Section:** В разработке (см. docs/AUTOMOTIVE_SECTION_STATUS_AND_PLAN.md)

---

**Дата последнего обновления:** 2025-09-29
**Backup оригинала:** CLAUDE.md.backup-*
