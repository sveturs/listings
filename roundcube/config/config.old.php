<?php
// Базовая конфигурация для работы самостоятельного (не в подпути) Roundcube
$config = [];

// Настройки почтового сервера
$config['imap_host'] = 'ssl://mailserver:993';
$config['smtp_host'] = 'tls://mailserver:587';
$config['smtp_port'] = 587;
$config['smtp_user'] = '%u';
$config['smtp_pass'] = '%p';

// Отключаем проверку SSL-сертификатов для внутренних соединений
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
$config['skin'] = 'elastic';

// Отладка
$config['debug_level'] = 1;
$config['smtp_debug'] = true;
$config['imap_debug'] = true;

// Заставляем работать на всех возможных URL
$config['use_https'] = true;

// Оставляем это в самом конце, чтобы Docker мог добавить свои настройки
include(__DIR__ . '/config.docker.inc.php');
