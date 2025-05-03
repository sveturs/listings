# Скрипты для интеграции с Harbor

Эта директория содержит скрипты, необходимые для установки, настройки и миграции на Harbor.

## Установка и настройка

- **install_harbor.sh** - установка Docker и Harbor на сервер
- **setup_harbor_ssl.sh** - настройка SSL для Harbor
- **setup_cert_trust.sh** - настройка доверия к сертификатам Harbor
- **update_harbor_domain.sh** - обновление Harbor для использования доменного имени
- **test_harbor_connection.sh** - проверка соединения с Harbor
- **test_curl_harbor.sh** - проверка доступности API Harbor

## Миграция образов

### Основные скрипты миграции

- **build_and_push.sh** - сборка и загрузка образов в Harbor
- **fixed_migrate_images.sh** - миграция нескольких образов Docker в Harbor
- **upload_opensearch.sh** - загрузка образов OpenSearch в Harbor
- **upload_roundcube.sh** - загрузка образа Roundcube в Harbor

### Скрипты для загрузки оставшихся образов

- **upload_mailserver.sh** - загрузка образа mailserver в Harbor с механизмом повторных попыток
- **upload_mail_webui.sh** - загрузка образа mail-webui в Harbor с механизмом повторных попыток
- **upload_remaining_mail_images.sh** - объединенный скрипт для загрузки обоих почтовых образов

## Скрипты развертывания

Скрипты развертывания предоставляют различные этапы миграции на Harbor:

1. **Частичное развертывание с Harbor**:
   - **deploy_limited_harbor.sh** - начальный скрипт для ограниченного развертывания с Harbor
   - **deploy_updated_harbor.sh** - обновленный скрипт для использования доступных образов из Harbor

2. **Полное развертывание с Harbor**:
   - **deploy_full_harbor.sh** - полное развертывание с использованием только образов из Harbor

## Процесс миграции

### Фаза 1: Настройка (Завершена)
1. Установка Harbor на выделенном сервере (207.180.197.172)
2. Настройка SSL и домена (harbor.svetu.rs)
3. Проверка соединения с производственного сервера

### Фаза 2: Миграция образов
1. Миграция основных образов (завершена)
2. Миграция оставшихся образов:
   - Выполнить `upload_mailserver.sh` и `upload_mail_webui.sh`
   - Альтернативно, использовать объединенный скрипт `upload_remaining_mail_images.sh`

### Фаза 3: Развертывание
1. Начальное/Частичное развертывание:
   - Использовать `deploy_updated_harbor.sh` для развертывания с доступными образами Harbor
   - Мониторинг поведения системы и обеспечение стабильности

2. Полная миграция:
   - После загрузки и тестирования всех образов использовать `deploy_full_harbor.sh`
   - Это развернет все сервисы с использованием образов из Harbor

## Статус образов

Успешно перенесены в Harbor:
- ✅ svetu/db/postgres:15
- ✅ svetu/opensearch/opensearch:2.11.0
- ✅ svetu/opensearch/dashboards:2.11.0
- ✅ svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
- ✅ svetu/minio/mc:latest
- ✅ svetu/tools/migrate:latest
- ✅ svetu/tools/certbot:latest
- ✅ svetu/backend/api:latest
- ✅ svetu/frontend/app:latest
- ✅ svetu/nginx/nginx:latest
- ✅ svetu/mail/server:latest
- ✅ svetu/mail/webui:latest

Все образы успешно загружены в Harbor! Можно переходить к полному развертыванию.

## Рекомендации по использованию

1. **Загрузка образов**:
   ```bash
   # Загрузка образа mailserver
   ./upload_mailserver.sh
   
   # Загрузка образа mail-webui
   ./upload_mail_webui.sh
   
   # Загрузка обоих оставшихся образов
   ./upload_remaining_mail_images.sh
   ```

2. **Развертывание**:
   ```bash
   # Частичное развертывание с доступными образами Harbor
   ./deploy_updated_harbor.sh
   
   # Полное развертывание со всеми образами Harbor (после загрузки всех образов)
   ./deploy_full_harbor.sh
   ```

## Устранение неполадок

- Если загрузка образов прерывается из-за тайм-аута, попробуйте специализированные скрипты с механизмом повторных попыток
- При проблемах с соединением проверьте статус сервера Harbor и настройки DNS
- При проблемах с развертыванием убедитесь, что все необходимые образы доступны в Harbor