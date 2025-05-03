# План перехода на Harbor в продакшн-окружении

## Текущий статус

Успешно выполнены следующие шаги:
1. Настроен Harbor с доменом harbor.svetu.rs и действительными SSL-сертификатами
2. Загружены следующие образы в Harbor:
   - svetu/db/postgres:15
   - svetu/backend/api:latest
   - svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
   - svetu/minio/mc:latest
   - svetu/tools/migrate:latest
   - svetu/nginx/nginx:latest
   - svetu/opensearch/dashboards:2.11.0
   - svetu/frontend/app:latest
3. Успешно запущена база данных с использованием образа из Harbor
4. Подготовлены скрипты для постепенного перехода на Harbor

## План постепенного перехода на Harbor

### Этап 1: Частичный деплой с Harbor (текущий)

1. Для минимизации рисков мы используем скрипт `deploy_limited_harbor.sh`, который:
   - Использует образы из Harbor для: db, backend, minio, createbuckets, migrate
   - Использует стандартные образы для остальных сервисов

2. Преимущества этого подхода:
   - Меньше рисков при переходе на Harbor
   - Возможность протестировать загрузку образов из Harbor в продакшн
   - Постепенный переход на новую инфраструктуру

### Этап 2: Миграция остальных образов в Harbor

1. Загрузить оставшиеся образы в Harbor:
   - opensearchproject/opensearch:2.11.0 → svetu/opensearch/opensearch:2.11.0
   - mailserver/docker-mailserver:latest → svetu/mail/server:latest
   - roundcube/roundcubemail:latest → svetu/mail/webui:latest
   - certbot/certbot → svetu/tools/certbot:latest

2. Обновить docker-compose.prod.yml.harbor для использования всех образов из Harbor

### Этап 3: Полный переход на Harbor

1. Перейти на полное использование Harbor:
   - Заменить docker-compose.prod.yml на docker-compose.prod.yml.harbor
   - Использовать скрипт deploy_with_harbor.sh для деплоя

2. Настроить CI/CD для автоматической сборки и отправки образов в Harbor:
   - Настроить GitHub Actions/Jenkins для сборки образов при пуше в main
   - Автоматическая отправка образов в Harbor с тегами на основе коммитов/версий

### Этап 4: Оптимизация и обеспечение безопасности

1. Создать обычного пользователя для запуска контейнеров вместо root:
   ```bash
   adduser svetu
   usermod -aG docker svetu
   mkdir -p /home/svetu/hostel-booking-system
   chown -R svetu:svetu /home/svetu/hostel-booking-system
   ```

2. Настроить логирование и мониторинг:
   - Настроить logrotate для Docker-логов
   - Добавить мониторинг доступности сервисов
   - Настроить оповещения при проблемах с сервисами

3. Настроить резервное копирование:
   - Автоматические бэкапы базы данных
   - Резервное копирование конфигурации и данных

4. Настроить обновление образов:
   - Регулярное обновление базовых образов
   - Политики удаления старых образов в Harbor

## Инструкции по выполнению Этапа 1

Для запуска частичного деплоя с использованием Harbor:

```bash
# На сервере svetu.rs
cd /opt/hostel-booking-system
./deploy_limited_harbor.sh
```

Этот скрипт:
1. Останавливает текущие сервисы
2. Авторизуется в Harbor
3. Загружает образы из Harbor
4. Создает бэкап базы данных
5. Запускает db, minio, backend из Harbor
6. Запускает остальные сервисы из стандартных образов

## Проверка после деплоя

После выполнения деплоя необходимо проверить:

1. Работоспособность всех сервисов:
   ```bash
   docker ps -a  # Все контейнеры должны быть в статусе Up
   ```

2. Доступность веб-интерфейса:
   ```bash
   curl -I https://svetu.rs
   ```

3. Логи контейнеров на наличие ошибок:
   ```bash
   docker logs backend
   docker logs hostel_db
   ```

## Откат в случае проблем

Если возникнут проблемы с деплоем, выполните откат:

```bash
# Остановка всех контейнеров
docker-compose -f docker-compose.prod.yml.harbor down

# Возврат к стандартному деплою
cd /opt/hostel-booking-system
./deploy.sh
```

## Дальнейшие шаги после успешного Этапа 1

1. Мониторинг работы сервисов в течение 24-48 часов
2. Сбор обратной связи от пользователей о работе системы
3. Планирование Этапа 2 - полной миграции на Harbor