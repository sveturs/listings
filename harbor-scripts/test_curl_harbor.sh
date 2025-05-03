#!/bin/bash
set -e

# Простой скрипт для тестирования доступности Harbor с использованием curl
# Этот скрипт можно запустить с любого сервера для проверки доступности Harbor

HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"

echo "==== Тестирование доступности Harbor ($HARBOR_URL) с помощью curl ===="

echo "1. Проверка доступности главной страницы:"
curl -s -I https://$HARBOR_URL | head -n 1

echo -e "\n2. Проверка доступности API:"
curl -s -I https://$HARBOR_URL/api/v2.0/health | head -n 1

echo -e "\n3. Проверка авторизации в API:"
curl -s -u $HARBOR_USER:$HARBOR_PASSWORD -X GET "https://$HARBOR_URL/api/v2.0/users/current" | grep "username"

echo -e "\n4. Проверка доступности проектов:"
curl -s -u $HARBOR_USER:$HARBOR_PASSWORD -X GET "https://$HARBOR_URL/api/v2.0/projects" | grep "svetu"

echo -e "\n5. Проверка списка репозиториев в проекте svetu:"
curl -s -u $HARBOR_USER:$HARBOR_PASSWORD -X GET "https://$HARBOR_URL/api/v2.0/projects/svetu/repositories" | grep "name"

echo -e "\n6. Проверка информации о сертификате:"
echo | openssl s_client -showcerts -servername $HARBOR_URL -connect $HARBOR_URL:443 2>/dev/null | openssl x509 -inform pem -noout -text | grep -A 2 "Issuer:" | head -n 3

echo -e "\n==== Тестирование завершено ===="