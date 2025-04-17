<?php
// Включаем подробный отчет об ошибках для отладки
ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);

// Логируем входящий запрос
error_log("Входящий запрос: " . $_SERVER['REQUEST_METHOD']);
error_log("Данные POST: " . print_r($_POST, true));

header('Content-Type: application/json');

// Разрешаем CORS для локальной разработки
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Methods: POST, GET, OPTIONS");
header("Access-Control-Allow-Headers: Content-Type");

// Обработка OPTIONS запроса (preflight)
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// Проверка метода запроса
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    error_log("Метод не POST: " . $_SERVER['REQUEST_METHOD']);
    http_response_code(405);
    echo json_encode(['success' => false, 'message' => 'Method not allowed']);
    exit;
}

// Получение данных из формы
$name = isset($_POST['name']) ? filter_var($_POST['name'], FILTER_SANITIZE_FULL_SPECIAL_CHARS) : '';
$email = isset($_POST['email']) ? filter_var($_POST['email'], FILTER_SANITIZE_EMAIL) : '';
$message = isset($_POST['message']) ? filter_var($_POST['message'], FILTER_SANITIZE_FULL_SPECIAL_CHARS) : '';

error_log("Получены данные: name=$name, email=$email, message=$message");

// Проверка данных
if (empty($name) || empty($email) || empty($message)) {
    error_log("Отсутствуют обязательные данные");
    http_response_code(400);
    echo json_encode(['success' => false, 'message' => 'All fields are required']);
    exit;
}

if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
    error_log("Неверный формат email: $email");
    http_response_code(400);
    echo json_encode(['success' => false, 'message' => 'Invalid email format']);
    exit;
}

// Формирование письма
$to = 'klimagrad@svetu.rs';
$subject = 'Сообщение с сайта KlimaGrad';

// Безопасное форматирование данных
$body = "Имя: $name\n";
$body .= "Email: $email\n";
$body .= "Сообщение:\n$message";

// Дополнительные заголовки
$headers = "From: info@svetu.rs\r\n";
$headers .= "Reply-To: $email\r\n";
$headers .= "X-Mailer: PHP/" . phpversion();

error_log("Отправка письма на $to");

// Отправка письма
$success = mail($to, $subject, $body, $headers);

if ($success) {
    error_log("Письмо успешно отправлено");
    echo json_encode(['success' => true, 'message' => 'Your message has been sent successfully']);
} else {
    error_log("Ошибка отправки письма: " . error_get_last()['message']);
    http_response_code(500);
    echo json_encode(['success' => false, 'message' => 'Failed to send your message']);
}
?>