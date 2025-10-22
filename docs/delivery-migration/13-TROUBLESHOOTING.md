# Troubleshooting

[Install]
WantedBy=multi-user.target
```

```bash
# Активация
sudo systemctl daemon-reload
sudo systemctl enable delivery-preprod.service
sudo systemctl start delivery-preprod.service
```

### 🔍 Мониторинг и отладка

#### Логи

```bash
# Все сервисы
docker-compose -f docker-compose.preprod.yml logs -f

# Только delivery-service
docker-compose -f docker-compose.preprod.yml logs -f delivery-service

# PostgreSQL
docker-compose -f docker-compose.preprod.yml logs -f delivery-postgres

# Redis
docker-compose -f docker-compose.preprod.yml logs -f delivery-redis
```

#### Проверка портов

```bash
# Занятые порты
sudo netstat -tlnp | grep -E "30051|35432|36379|38080|38081|39090"

# Процессы Docker
docker ps | grep svetudelivery
```

#### Подключение к базе данных

```bash
# Из хоста
psql "postgres://delivery_user:PASSWORD@localhost:35432/delivery_db"

# Или через docker exec
docker exec -it svetudelivery-postgres psql -U delivery_user -d delivery_db
```

#### Проверка Redis

```bash
# Ping
docker exec svetudelivery-redis redis-cli -a PASSWORD ping

# Мониторинг команд
docker exec svetudelivery-redis redis-cli -a PASSWORD monitor
```

### 🚨 Troubleshooting

#### Проблема: Порт 30051 занят

```bash
# Найти процесс
sudo lsof -i :30051

# Остановить конфликтующий сервис
docker-compose -f /opt/OTHER_SERVICE/docker-compose.yml stop
```

#### Проблема: БД не поднимается

```bash
# Проверка логов
docker logs svetudelivery-postgres

# Проверка прав доступа
docker exec svetudelivery-postgres ls -la /var/lib/postgresql/data

# Пересоздание volume
docker-compose -f docker-compose.preprod.yml down -v
docker-compose -f docker-compose.preprod.yml up -d
```

#### Проблема: gRPC недоступен

```bash
# Проверка Nginx конфигурации
sudo nginx -t

# Проверка SSL сертификата
sudo certbot certificates

# Проверка firewall
sudo ufw status
```

### 📊 Ресурсы сервера

**Текущее состояние** (2025-10-22):
- **Диск**: 22GB свободно из 193GB (90% использовано)
- **Docker**: версия 27.5.1
- **Go**: версия 1.25.0

**Рекомендации**:
1. ⚠️ Мониторить место на диске (осталось мало!)
2. Настроить ротацию логов Docker
3. Очистить старые образы: `docker system prune -a`

### 🔄 Интеграция с монолитом

После развертывания микросервиса, монолит будет обращаться к нему через:

**gRPC адрес (внутренний)**: `localhost:30051`
**gRPC адрес (внешний)**: `deliverypreprod.svetu.rs:443`

**Конфигурация в монолите** (`backend/internal/config/config.go`):
```go
type DeliveryConfig struct {
    GRPCAddress string `env:"DELIVERY_GRPC_ADDRESS" envDefault:"localhost:30051"`
    UseTLS      bool   `env:"DELIVERY_USE_TLS" envDefault:"false"`
}
```

