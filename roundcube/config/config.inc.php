<?php
// Базовая конфигурация для Roundcube
$config = [];
// Настройки почтового сервера
$config['imap_host'] = 'ssl://mailserver:993';
$config['smtp_host'] = 'tls://mailserver:587';
$config['smtp_port'] = 587;
$config['smtp_user'] = '%u';
$config['smtp_pass'] = '%p';
// Отключение проверки SSL
$config['imap_conn_options'] = [
    'ssl' => [
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ]
];
$config['smtp_conn_options'] = [
    'ssl' => [
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ]
];
// Основные настройки
$config['product_name'] = 'Svetu.rs Mail';
$config['des_key'] = 'random-string-for-encryption';
$config['plugins'] = ['archive', 'zipdownload'];
//$config['plugins'] = ['archive', 'zipdownload','password', 'new_user_identity'];
$config['skin'] = 'elastic';
// Настройки защиты от CSRF
$config['csrf_protection'] = true; 
$config['session_lifetime'] = 60;
$config['session_storage'] = 'php';
$config['check_all_folders'] = true;
// Настройки сессий и доступа
$config['use_https'] = true;
$config['force_https'] = true;
// Настройки для корректной работы с абсолютными URL
$config['absolute_url'] = true;
// Базы данных - убедитесь, что директория существует и доступна для записи
$config['db_dsnw'] = 'sqlite:////var/roundcube/db/sqlite.db?mode=0646';

// Добавляем опции для регистрации новых пользователей
$config['new_user_dialog'] = true;
$config['support_url'] = 'mailto:admin@svetu.rs';

// Включаем конфигурацию Docker и перезаписываем ключевые настройки
include(__DIR__ . '/config.docker.inc.php');
// Важные настройки, которые должны быть перезаписаны
$config['request_path'] = '';
$config['base_path'] = '';
$config['login_url'] = '?_task=login';
