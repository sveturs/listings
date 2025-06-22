# Backend Rules

## Архитектура Backend

### Backend структура
```
backend/
├── cmd/api/          # Точка входа приложения
├── internal/
│   ├── config/       # Конфигурация
│   ├── domain/       # Доменные модели
│   ├── middleware/   # Auth, CORS, Logger
│   ├── proj/         # Бизнес-логика модулей
│   │   ├── marketplace/
│   │   ├── users/
│   │   ├── payments/
│   │   └── ...
│   ├── server/       # HTTP сервер (Fiber)
│   └── storage/      # Репозитории
│       ├── postgres/
│       ├── opensearch/
│       └── minio/
└── migrations/       # SQL миграции
```

### Сервисы инфраструктуры
- **PostgreSQL** - основная база данных
- **OpenSearch** - полнотекстовый поиск
- **MinIO** - S3-совместимое хранилище изображений
- **Nginx** - обратный прокси и статика
- **Harbor** - приватный Docker registry

## Конфигурация Backend

- Backend: `.env` в корне backend/

## Handlers

- Не возвращаем пользователю реальную ошибку. Нужно возвращать специальные placeholders из переводов. Саму ошибку нужно логировать.
- В swagger описании нужно возвращать структуры из пакета [@backend/pkg/utils/utils.go](backend/pkg/utils/utils.go)
  - type SuccessResponseSwag struct
  - type ErrorResponseSwag struct
  - Примеры описания swagger комментариев
    - // @Success 200 {object} utils.SuccessResponseSwag{data=models.UserBalance} "User balance information"
    - // @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"


## База данных

- Старайся не использовать SQL триггеры в базе и писать логику в коде

### Основные таблицы
- `users` - пользователи и аутентификация
- `marketplace_listings` - объявления
- `marketplace_categories` - категории с атрибутами
- `translations` - мультиязычные переводы
- `marketplace_images` - изображения объявлений
- `marketplace_chats` - чаты и сообщения
  TODO: Нужно доработать

## External Services

- **MinIO**: Object storage for images
    - Local: http://localhost:9000
    - Production: https://svetu.rs
    - Images are served from `/listings/` path



