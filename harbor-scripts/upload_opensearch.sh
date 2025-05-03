#!/bin/bash
set -e

# Скрипт для загрузки образа opensearch в Harbor
HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROJECT_NAME="svetu"

echo "==== Загрузка образа OpenSearch в Harbor ===="

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Загрузка образа opensearch
echo "Загрузка образа opensearchproject/opensearch:2.11.0..."
docker pull opensearchproject/opensearch:2.11.0

# Переименование образа для Harbor
echo "Переименование образа для Harbor..."
docker tag opensearchproject/opensearch:2.11.0 $HARBOR_URL/$PROJECT_NAME/opensearch/opensearch:2.11.0

# Отправка образа в Harbor
echo "Отправка образа в Harbor..."
docker push $HARBOR_URL/$PROJECT_NAME/opensearch/opensearch:2.11.0

echo "==== Загрузка образа OpenSearch в Harbor завершена ===="