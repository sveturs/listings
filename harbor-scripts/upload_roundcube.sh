#!/bin/bash
set -e

# Скрипт для загрузки образа Roundcube в Harbor
HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROJECT_NAME="svetu"

echo "==== Загрузка образа Roundcube в Harbor ===="

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Загрузка образа
echo "Загрузка образа roundcube/roundcubemail:latest..."
docker pull roundcube/roundcubemail:latest

# Переименование образа для Harbor
echo "Переименование образа для Harbor..."
docker tag roundcube/roundcubemail:latest $HARBOR_URL/$PROJECT_NAME/mail/webui:latest

# Отправка образа в Harbor
echo "Отправка образа в Harbor..."
docker push $HARBOR_URL/$PROJECT_NAME/mail/webui:latest

echo "==== Загрузка образа Roundcube в Harbor завершена ===="