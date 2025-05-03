#!/bin/bash
set -e

# Скрипт для установки Docker и Harbor
# Запускать на сервере 207.180.197.172
# Выполнить: chmod +x install_harbor.sh && ./install_harbor.sh

echo "==== Установка Docker и Harbor для Sve Tu Platform ===="

# Обновление пакетов
echo "1. Обновление системы..."
sudo apt update
sudo apt upgrade -y

# Установка зависимостей
echo "2. Установка необходимых зависимостей..."
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common gnupg

# Установка Docker
echo "3. Добавление Docker репозитория..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "4. Установка Docker Engine..."
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Добавление текущего пользователя в группу docker
echo "5. Настройка прав Docker..."
sudo usermod -aG docker $USER
echo "Для применения изменений, пожалуйста, выйдите и снова войдите в систему или выполните: newgrp docker"

# Проверка установки Docker
echo "6. Проверка Docker..."
docker --version
docker compose version

# Создание директорий для Harbor
echo "7. Создание директорий для Harbor..."
sudo mkdir -p /opt/harbor
sudo chown $USER:$USER /opt/harbor
cd /opt/harbor

# Загрузка Harbor
echo "8. Загрузка Harbor..."
HARBOR_VERSION="v2.10.0"
wget https://github.com/goharbor/harbor/releases/download/${HARBOR_VERSION}/harbor-online-installer-${HARBOR_VERSION}.tgz

# Распаковка Harbor
echo "9. Распаковка Harbor..."
tar xvf harbor-online-installer-${HARBOR_VERSION}.tgz
cd harbor

# Копирование шаблона конфигурации
echo "10. Настройка конфигурации Harbor..."
cp harbor.yml.tmpl harbor.yml

# Настройка базовой конфигурации
echo "11. Редактирование harbor.yml..."
echo "Отредактируйте файл harbor.yml, установив:"
echo "- hostname: задайте ваш домен или IP-адрес сервера (207.180.197.172)"
echo "- harbor_admin_password: установите надежный пароль для admin"
echo "- data_volume: укажите путь для хранения данных (по умолчанию /data)"
echo ""
echo "Для HTTPS настройки добавьте пути к сертификатам:"
echo "- certificate: /your/certificate/path"
echo "- private_key: /your/private/key/path"
echo ""
echo "Выполните редактирование сейчас: sudo nano harbor.yml"
echo "Нажмите ENTER, чтобы продолжить..."
read -p ""

# Запуск установки Harbor
echo "12. Запуск установки Harbor..."
sudo ./install.sh

echo "==== Установка завершена ===="
echo "Harbor должен быть доступен по адресу: https://207.180.197.172"
echo "Логин: admin"
echo "Пароль: тот, который вы указали в harbor.yml"
echo ""
echo "Для проверки работы Harbor выполните в браузере переход по адресу https://207.180.197.172"
echo "Или проверьте работу контейнеров: docker ps"
echo ""
echo "Следующий шаг: настройка CI/CD с помощью setup_harbor_ci_cd.sh"