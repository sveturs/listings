Инструкция по работе с проектом для AI агента

Проект называется Sve Tu Platform
сейчас в разработке маркетплейс

Рабочие окружения

1. Локальная разработка: /data/hostel-booking-system/
2. Боевой сервер: ssh root@svetu.rs, проект в /opt/hostel-booking-system/

Команды сборки

- Backend: go build -o main ./cmd/api
- Frontend: cd frontend/hostel-frontend && npm run build
- Запуск: docker-compose up -d

Деплой

- На сервер: ./harbor-scripts/blue_green_deploy_on_svetu.rs.sh [backend|frontend|all]
- Синий-зеленый деплой для нулевого простоя

Конфигурация

- Frontend-настройки: frontend/hostel-frontend/public/env.js
- Nginx: /opt/hostel-booking-system/nginx.conf (на сервере)
- Docker: docker-compose.yml, docker-compose.prod.yml

Технологии

- Backend: Go, PostgreSQL, OpenSearch, MinIO
- Frontend: React, MUI, Leaflet
- Инфраструктура: Docker, Nginx, Harbor, mail

и пиши мне на русском.
