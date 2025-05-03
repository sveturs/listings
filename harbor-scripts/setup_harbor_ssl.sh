#!/bin/bash
set -e

# Скрипт для настройки SSL для Harbor
# Для использования на сервере 207.180.197.172

# Настройки
DOMAIN="207.180.197.172" # Используем IP как домен
EMAIL="info@svetu.rs" # Email для Let's Encrypt
HARBOR_DIR="/opt/harbor"
CERT_DIR="/opt/harbor/certs"

echo "==== Настройка SSL для Harbor ===="

# Установка инструментов для работы с сертификатами
echo "1. Установка необходимых инструментов..."
sudo apt-get update
sudo apt-get install -y openssl

# Создание директории для сертификатов
echo "2. Создание директории для сертификатов..."
sudo mkdir -p $CERT_DIR
sudo chown $USER:$USER $CERT_DIR

# Создание самоподписанного сертификата (для начала)
echo "3. Создание самоподписанного сертификата..."
cd $CERT_DIR

# Создание конфигурационного файла для OpenSSL
cat > $CERT_DIR/openssl.conf << EOF
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[req_distinguished_name]
C = RS
ST = Serbia
L = Belgrade
O = Sve Tu Platform
OU = IT
CN = $DOMAIN

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
IP.1 = $DOMAIN
EOF

# Генерация ключа и сертификата
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout $CERT_DIR/server.key \
  -out $CERT_DIR/server.crt \
  -config $CERT_DIR/openssl.conf

echo "4. Обновление конфигурации Harbor..."
# Обновляем harbor.yml - раскомментируем и настраиваем HTTPS
sudo cp $HARBOR_DIR/harbor/harbor.yml $HARBOR_DIR/harbor/harbor.yml.bak
cat > $HARBOR_DIR/harbor/harbor.yml << EOF
# Configuration file of Harbor

# The IP address or hostname to access admin UI and registry service.
# DO NOT use localhost or 127.0.0.1, because Harbor needs to be accessed by external clients.
hostname: $DOMAIN

# http related config
http:
  # port for http, default is 80. If https enabled, this port will redirect to https port
  port: 80

# https related config
https:
  # https port for harbor, default is 443
  port: 443
  # The path of cert and key files for nginx
  certificate: $CERT_DIR/server.crt
  private_key: $CERT_DIR/server.key

# The initial password of Harbor admin
# It only works in first time to install harbor
# Remember Change the admin password from UI after launching Harbor.
harbor_admin_password: SveTu2025

# Harbor DB configuration
database:
  # The password for the root user of Harbor DB. Change this before any production use.
  password: root123
  # The maximum number of connections in the idle connection pool. If it <=0, no idle connections are retained.
  max_idle_conns: 100
  # The maximum number of open connections to the database. If it <= 0, then there is no limit on the number of open connections.
  # Note: the default number of connections is 1024 for postgres of harbor.
  max_open_conns: 900
  # The maximum amount of time a connection may be reused. Expired connections may be closed lazily before reuse. If it <= 0, connections are not closed due to a connection's age.
  # The value is a duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
  conn_max_lifetime: 5m
  # The maximum amount of time a connection may be idle. Expired connections may be closed lazily before reuse. If it <= 0, connections are not closed due to a connection's idle time.
  # The value is a duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
  conn_max_idle_time: 0

# The default data volume
data_volume: $HARBOR_DIR/data

# Trivy configuration
trivy:
  ignore_unfixed: false
  skip_update: false
  offline_scan: false
  security_check: vuln
  insecure: false

jobservice:
  max_job_workers: 10
  job_loggers:
    - STD_OUTPUT
    - FILE
  logger_sweeper_duration: 1

notification:
  webhook_job_max_retry: 3
  webhook_job_http_client_timeout: 3

# Log configurations
log:
  level: info
  local:
    rotate_count: 50
    rotate_size: 200M
    location: /var/log/harbor

# This attribute is for migrator to detect the version of the .cfg file, DO NOT MODIFY!
_version: 2.10.0

# Global proxy
proxy:
  http_proxy:
  https_proxy:
  no_proxy:
  components:
    - core
    - jobservice
    - trivy

# Upload purging
upload_purging:
  enabled: true
  age: 168h
  interval: 24h
  dryrun: false

# Cache settings
cache:
  enabled: false
  expire_hours: 24
EOF

echo "5. Перезапуск Harbor с поддержкой HTTPS..."
cd $HARBOR_DIR/harbor

# Перезапуск Harbor
sudo ./prepare
sudo docker compose down -v
sudo docker compose up -d

echo "6. Добавление сертификата в Docker..."
# Добавление сертификата в Docker
sudo mkdir -p /etc/docker/certs.d/$DOMAIN
sudo cp $CERT_DIR/server.crt /etc/docker/certs.d/$DOMAIN/

echo "==== Настройка SSL для Harbor завершена! ===="
echo "Теперь Harbor доступен по HTTPS: https://$DOMAIN"
echo ""
echo "Для использования Harbor с Docker выполните:"
echo "  docker login -u admin -p SveTu2025 $DOMAIN"
echo ""
echo "Важно: Поскольку сертификат самоподписанный, вам потребуется добавить его в доверенные на клиентских машинах"