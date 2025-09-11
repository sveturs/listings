# Управление SSL сертификатами для svetu.rs

## 📋 Обзор

Проект использует **Let's Encrypt** для бесплатных SSL сертификатов. Все домены защищены валидными сертификатами с автоматическим обновлением.

### Защищенные домены:
- **svetu.rs** - основной production домен (+ www.svetu.rs)
- **dev.svetu.rs** - development версия сайта
- **devapi.svetu.rs** - API для development версии
- **devs3.svetu.rs** - S3 хранилище для development

## 🏗️ Архитектура

### Структура файлов
```
/etc/letsencrypt/
├── live/                      # Символические ссылки на актуальные сертификаты
│   ├── svetu.rs/
│   │   ├── fullchain.pem     # Сертификат + промежуточные
│   │   ├── privkey.pem       # Приватный ключ
│   │   ├── cert.pem          # Только сертификат
│   │   └── chain.pem         # Цепочка промежуточных
│   ├── dev.svetu.rs/
│   ├── devapi.svetu.rs/
│   └── devs3.svetu.rs/
├── archive/                   # Все версии сертификатов
└── renewal/                   # Конфигурации для обновления

/opt/nginx-simple/
├── conf.d/                    # Конфигурации nginx
│   ├── svetu.rs.conf
│   ├── dev.svetu.rs.conf
│   ├── devapi.svetu.rs.conf
│   └── devs3.svetu.rs.conf
└── certbot/
    └── www/                   # Webroot для ACME challenge

/usr/local/bin/
├── renew-certificates.sh      # Скрипт автообновления
└── check-certificates.sh      # Скрипт проверки статуса

/var/log/letsencrypt/
└── renewal-YYYYMM.log        # Логи обновления
```

### Docker интеграция

Nginx работает в Docker контейнере с монтированными сертификатами:
```yaml
volumes:
  - /etc/letsencrypt:/etc/letsencrypt:ro
  - /opt/nginx-simple/certbot/www:/var/www/certbot:ro
```

## 🔄 Автоматическое обновление

### Три уровня защиты:

1. **Systemd Timer** (встроенный в Ubuntu)
   - Запускается дважды в день
   - Файл: `/usr/lib/systemd/system/certbot.timer`
   - Статус: `systemctl status certbot.timer`

2. **Cron задача** (основной механизм)
   - Файл: `/etc/cron.d/certbot-renew`
   - Расписание: 3:17 и 15:17 ежедневно
   - Скрипт: `/usr/local/bin/renew-certificates.sh`

3. **Ручное обновление** (при необходимости)
   - Команда: `sudo /usr/local/bin/renew-certificates.sh`

### Логика обновления

Сертификаты обновляются автоматически за 30 дней до истечения срока действия.

#### Скрипт renew-certificates.sh:
1. Проверяет статус всех сертификатов
2. Запускает `certbot renew` для обновления
3. При успехе перезапускает nginx контейнер
4. При ошибке пытается обновить каждый домен отдельно
5. Логирует все действия в `/var/log/letsencrypt/renewal-YYYYMM.log`
6. Ротирует логи старше 90 дней

## 📊 Мониторинг

### Проверка статуса сертификатов
```bash
sudo /usr/local/bin/check-certificates.sh
```

Вывод показывает:
- ✅ OK - более 60 дней до истечения
- ℹ️ INFO - 30-60 дней до истечения  
- ⚠️ WARNING - 7-30 дней до истечения
- ❌ CRITICAL - менее 7 дней до истечения

### Просмотр логов
```bash
# Последние записи логов обновления
sudo tail -f /var/log/letsencrypt/renewal-*.log

# Системные логи certbot
sudo journalctl -u certbot.service -f

# Логи Let's Encrypt
sudo tail -f /var/log/letsencrypt/letsencrypt.log
```

## 🔧 Основные команды

### Управление сертификатами
```bash
# Список всех сертификатов и их статус
sudo certbot certificates

# Тестовое обновление (dry-run)
sudo certbot renew --dry-run

# Принудительное обновление всех сертификатов
sudo certbot renew --force-renewal

# Обновление конкретного домена
sudo certbot renew --cert-name svetu.rs --force-renewal

# Удаление сертификата
sudo certbot delete --cert-name domain.com
```

### Получение новых сертификатов
```bash
# Webroot метод (nginx продолжает работать)
sudo certbot certonly --webroot \
  -w /opt/nginx-simple/certbot/www \
  -d newdomain.svetu.rs \
  --non-interactive \
  --agree-tos \
  --email admin@svetu.rs

# Standalone метод (требует остановки nginx)
sudo certbot certonly --standalone \
  -d newdomain.svetu.rs \
  --non-interactive \
  --agree-tos \
  --pre-hook "docker stop svetu_nginx" \
  --post-hook "docker start svetu_nginx"
```

### Nginx управление
```bash
# Проверка конфигурации
docker exec svetu_nginx nginx -t

# Перезагрузка конфигурации
docker exec svetu_nginx nginx -s reload

# Перезапуск контейнера
docker restart svetu_nginx
```

## 🚨 Устранение проблем

### Проблема: Сертификат не обновляется автоматически

1. Проверьте cron задачу:
```bash
cat /etc/cron.d/certbot-renew
sudo crontab -l | grep certbot
```

2. Проверьте systemd timer:
```bash
systemctl status certbot.timer
systemctl list-timers | grep certbot
```

3. Запустите обновление вручную с подробным выводом:
```bash
sudo certbot renew -v
```

### Проблема: Ошибка ACME challenge

1. Проверьте доступность директории webroot:
```bash
ls -la /opt/nginx-simple/certbot/www/
docker exec svetu_nginx ls -la /var/www/certbot/
```

2. Проверьте nginx конфигурацию для .well-known:
```bash
grep -r "well-known" /opt/nginx-simple/conf.d/
```

3. Тест доступности:
```bash
echo "test" | sudo tee /opt/nginx-simple/certbot/www/test.txt
curl http://yourdomain.com/.well-known/acme-challenge/test.txt
```

### Проблема: Certbot зависает

1. Убейте зависшие процессы:
```bash
sudo pkill -9 certbot
```

2. Очистите lock файлы:
```bash
sudo rm -f /var/lib/letsencrypt/.certbot.lock
```

3. Используйте альтернативный метод обновления:
```bash
# Вместо standalone используйте webroot
sudo certbot renew --authenticator webroot \
  --webroot-path /opt/nginx-simple/certbot/www
```

### Проблема: Nginx не видит обновленные сертификаты

1. Проверьте монтирование volumes:
```bash
docker inspect svetu_nginx | grep -A 10 Mounts
```

2. Перезапустите контейнер:
```bash
docker restart svetu_nginx
```

3. Проверьте права доступа:
```bash
sudo ls -la /etc/letsencrypt/live/
sudo ls -la /etc/letsencrypt/archive/
```

## 📝 Добавление нового домена

### Шаг 1: Создайте nginx конфигурацию

```nginx
# /opt/nginx-simple/conf.d/newdomain.conf
server {
    listen 80;
    server_name newdomain.svetu.rs;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$server_name$request_uri;
    }
}
```

### Шаг 2: Перезагрузите nginx
```bash
docker exec svetu_nginx nginx -s reload
```

### Шаг 3: Получите сертификат
```bash
sudo certbot certonly --webroot \
  -w /opt/nginx-simple/certbot/www \
  -d newdomain.svetu.rs \
  --non-interactive \
  --agree-tos
```

### Шаг 4: Добавьте HTTPS в конфигурацию
```nginx
server {
    listen 443 ssl;
    http2 on;
    server_name newdomain.svetu.rs;

    ssl_certificate /etc/letsencrypt/live/newdomain.svetu.rs/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/newdomain.svetu.rs/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # ... остальная конфигурация ...
}
```

### Шаг 5: Финальная перезагрузка
```bash
docker exec svetu_nginx nginx -s reload
```

## 📅 Расписание обслуживания

- **Ежедневно 3:17 и 15:17** - автоматическая проверка и обновление
- **За 30 дней до истечения** - автоматическое обновление сертификата
- **Каждые 90 дней** - ротация старых логов
- **При обновлении** - автоматический перезапуск nginx

## 🔐 Безопасность

### Рекомендации:
1. Регулярно проверяйте статус сертификатов
2. Подпишитесь на уведомления Let's Encrypt об истечении
3. Храните резервные копии `/etc/letsencrypt/`
4. Используйте мониторинг внешних сервисов (например, UptimeRobot)
5. Настройте оповещения в логах при ошибках обновления

### Резервное копирование
```bash
# Backup сертификатов
sudo tar -czf letsencrypt-backup-$(date +%Y%m%d).tar.gz /etc/letsencrypt/

# Восстановление
sudo tar -xzf letsencrypt-backup-YYYYMMDD.tar.gz -C /
```

## 📞 Контакты и поддержка

При критических проблемах с сертификатами:
1. Проверьте статус: `sudo /usr/local/bin/check-certificates.sh`
2. Запустите ручное обновление: `sudo /usr/local/bin/renew-certificates.sh`
3. Проверьте логи: `sudo tail -100 /var/log/letsencrypt/renewal-*.log`

---

*Последнее обновление: 8 сентября 2025*
*Автор: System Administrator*