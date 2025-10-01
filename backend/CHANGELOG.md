# Backend Changelog

All notable changes to the backend will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-10-01

### Changed
- **Version**: Updated from 0.1.1 to 0.2.0
- **Architecture**: Backend готов к работе с BFF proxy от frontend
  - Все API endpoints остаются на `/api/v1/*`
  - Поддержка Authorization header от BFF proxy
  - Нет breaking changes в API

### Technical Details
- `internal/version/version.go`: Version = "0.2.0"

### Notes
- Backend API остается стабильным
- Изменения в основном на стороне frontend (BFF proxy)
- Полная обратная совместимость

## [0.1.1] - 2025-09-XX

### Initial Release
- Базовая функциональность маркетплейса
- Интеграция с Auth Service
- REST API `/api/v1/*`
- PostgreSQL, OpenSearch, MinIO, Redis
