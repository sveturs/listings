# Исправление конфигурации mail-webui в docker-compose

## Проблема
Mail-webui не мог подключиться к mailserver из-за неправильной конфигурации IMAP/SMTP хостов и требований SSL.

## Решение

Замените секцию `mail-webui` в `docker-compose.prod.yml` на следующую:

```yaml
  mail-webui:
    image: harbor.svetu.rs/svetu/mail/webui:latest
    container_name: mail-webui
    volumes:
      - ./roundcube/config:/var/www/html/config
      - roundcube_data:/var/roundcube
    environment:
      - ROUNDCUBEMAIL_DEFAULT_HOST=ssl://svetu.rs:993
      - ROUNDCUBEMAIL_SMTP_SERVER=tls://svetu.rs:587
      - ROUNDCUBEMAIL_SMTP_PORT=587
      - ROUNDCUBEMAIL_PLUGINS=archive,zipdownload
      - ROUNDCUBEMAIL_DES_KEY=random-string-for-encryption
      - ROUNDCUBEMAIL_SKIN=elastic
    expose:
      - "80"
    networks:
      hostel_network:
        aliases:
          - mail-webui
    restart: unless-stopped
    depends_on:
      - mailserver
```

## Ключевые изменения

1. **IMAP хост**: `ssl://svetu.rs:993` вместо `ssl://mailserver:993`
2. **SMTP хост**: `tls://svetu.rs:587` вместо `tls://mailserver:587`

## Почему это работает

- Использует внешний хост `svetu.rs` вместо внутреннего `mailserver`
- SSL/TLS подключения работают корректно через внешний интерфейс
- Избегает проблем с внутренней сетевой связностью Docker контейнеров

## Автоматические исправления

Также настроены cron задачи для автоматического исправления конфигурации:

```bash
# Каждые 5 минут - исправление прав на БД
*/5 * * * * docker exec mail-webui chown www-data:www-data /var/roundcube/db/sqlite.db 2>/dev/null || true

# Каждые 10 минут - исправление конфигурации хостов
*/10 * * * * docker exec mail-webui sed -i "s#ssl://mailserver:993:143#ssl://svetu.rs:993#g; s#tls://mailserver:587:587#tls://svetu.rs:587#g; s#ssl://mailserver:993#ssl://svetu.rs:993#g; s#tls://mailserver:587#tls://svetu.rs:587#g" /var/www/html/config/config.docker.inc.php /var/www/html/config/config.inc.php 2>/dev/null || true
```

## Применение изменений

После обновления docker-compose.prod.yml:

```bash
cd /opt/hostel-booking-system
docker-compose -f docker-compose.prod.yml down mail-webui
docker-compose -f docker-compose.prod.yml up -d mail-webui
```

## Проверка работы

1. Откройте https://mail.svetu.rs
2. Войдите с полным email адресом (например: `user@svetu.rs`) и паролем
3. Проверьте логи: `docker logs mail-webui`