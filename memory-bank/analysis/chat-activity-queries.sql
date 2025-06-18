-- Анализ активности чатов в системе

-- 1. Структура таблицы marketplace_chats
-- id - уникальный идентификатор чата
-- listing_id - ID объявления (может быть NULL для прямых сообщений)
-- buyer_id - ID покупателя
-- seller_id - ID продавца  
-- last_message_at - время последнего сообщения
-- created_at - время создания чата
-- updated_at - время последнего обновления
-- is_archived - флаг архивации чата

-- 2. Структура таблицы marketplace_messages
-- id - уникальный идентификатор сообщения
-- chat_id - ID чата
-- listing_id - ID объявления (может быть NULL)
-- sender_id - ID отправителя
-- receiver_id - ID получателя
-- content - текст сообщения
-- is_read - прочитано ли сообщение
-- original_language - язык оригинала
-- created_at - время создания
-- updated_at - время обновления
-- has_attachments - есть ли вложения
-- attachments_count - количество вложений

-- 3. Связь с объявлениями
-- Чат связан с объявлением через listing_id
-- Если listing_id IS NULL - это прямое сообщение между пользователями
-- В чате участвуют buyer_id и seller_id

-- ЗАПРОСЫ ДЛЯ АНАЛИЗА АКТИВНОСТИ:

-- 1. Общая статистика по чатам
SELECT 
    COUNT(DISTINCT c.id) as total_chats,
    COUNT(DISTINCT CASE WHEN c.listing_id IS NOT NULL THEN c.id END) as chats_with_listings,
    COUNT(DISTINCT CASE WHEN c.listing_id IS NULL THEN c.id END) as direct_chats,
    COUNT(DISTINCT CASE WHEN c.is_archived THEN c.id END) as archived_chats
FROM marketplace_chats c;

-- 2. Количество сообщений по пользователю
SELECT 
    u.id as user_id,
    u.name,
    u.email,
    COUNT(DISTINCT CASE WHEN m.sender_id = u.id THEN m.id END) as sent_messages,
    COUNT(DISTINCT CASE WHEN m.receiver_id = u.id THEN m.id END) as received_messages,
    COUNT(DISTINCT CASE WHEN m.sender_id = u.id OR m.receiver_id = u.id THEN m.chat_id END) as active_chats,
    MAX(CASE WHEN m.sender_id = u.id THEN m.created_at END) as last_sent_at,
    MAX(CASE WHEN m.receiver_id = u.id THEN m.created_at END) as last_received_at
FROM users u
LEFT JOIN marketplace_messages m ON m.sender_id = u.id OR m.receiver_id = u.id
GROUP BY u.id, u.name, u.email
ORDER BY (COUNT(DISTINCT CASE WHEN m.sender_id = u.id THEN m.id END) + 
          COUNT(DISTINCT CASE WHEN m.receiver_id = u.id THEN m.id END)) DESC;

-- 3. Активность по чатам (топ 10 самых активных)
SELECT 
    c.id as chat_id,
    c.listing_id,
    l.title as listing_title,
    COUNT(m.id) as message_count,
    MIN(m.created_at) as first_message,
    MAX(m.created_at) as last_message,
    c.last_message_at,
    buyer.name as buyer_name,
    seller.name as seller_name,
    c.is_archived
FROM marketplace_chats c
LEFT JOIN marketplace_listings l ON c.listing_id = l.id
LEFT JOIN marketplace_messages m ON m.chat_id = c.id
LEFT JOIN users buyer ON c.buyer_id = buyer.id
LEFT JOIN users seller ON c.seller_id = seller.id
GROUP BY c.id, c.listing_id, l.title, c.last_message_at, buyer.name, seller.name, c.is_archived
ORDER BY message_count DESC
LIMIT 10;

-- 4. Непрочитанные сообщения по пользователям
SELECT 
    u.id as user_id,
    u.name,
    COUNT(m.id) as unread_messages,
    COUNT(DISTINCT m.chat_id) as chats_with_unread
FROM users u
JOIN marketplace_messages m ON m.receiver_id = u.id AND m.is_read = false
GROUP BY u.id, u.name
ORDER BY unread_messages DESC;

-- 5. Активность по дням
SELECT 
    DATE(m.created_at) as message_date,
    COUNT(m.id) as messages_count,
    COUNT(DISTINCT m.sender_id) as unique_senders,
    COUNT(DISTINCT m.chat_id) as active_chats
FROM marketplace_messages m
WHERE m.created_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE(m.created_at)
ORDER BY message_date DESC;

-- 6. Чаты с удаленными объявлениями
SELECT 
    c.id as chat_id,
    c.listing_id,
    COUNT(m.id) as message_count,
    MAX(m.created_at) as last_message_at,
    buyer.name as buyer_name,
    seller.name as seller_name
FROM marketplace_chats c
LEFT JOIN marketplace_listings l ON c.listing_id = l.id
LEFT JOIN marketplace_messages m ON m.chat_id = c.id
LEFT JOIN users buyer ON c.buyer_id = buyer.id
LEFT JOIN users seller ON c.seller_id = seller.id
WHERE c.listing_id IS NOT NULL AND l.id IS NULL
GROUP BY c.id, c.listing_id, buyer.name, seller.name
ORDER BY last_message_at DESC;

-- 7. Статистика по пользователю (параметризованный запрос)
-- Замените :user_id на нужный ID пользователя
SELECT 
    u.id,
    u.name,
    u.email,
    (SELECT COUNT(*) FROM marketplace_messages WHERE sender_id = u.id) as total_sent,
    (SELECT COUNT(*) FROM marketplace_messages WHERE receiver_id = u.id) as total_received,
    (SELECT COUNT(*) FROM marketplace_messages WHERE receiver_id = u.id AND is_read = false) as unread_count,
    (SELECT COUNT(DISTINCT chat_id) FROM marketplace_messages WHERE sender_id = u.id OR receiver_id = u.id) as total_chats,
    (SELECT COUNT(*) FROM marketplace_chats WHERE (buyer_id = u.id OR seller_id = u.id) AND is_archived = true) as archived_chats,
    (SELECT MAX(created_at) FROM marketplace_messages WHERE sender_id = u.id) as last_sent_at,
    (SELECT MAX(created_at) FROM marketplace_messages WHERE receiver_id = u.id) as last_received_at
FROM users u
WHERE u.id = :user_id;

-- 8. Детальная информация по чату
-- Замените :chat_id на нужный ID чата
SELECT 
    c.id as chat_id,
    c.listing_id,
    l.title as listing_title,
    l.status as listing_status,
    c.created_at as chat_created,
    c.last_message_at,
    c.is_archived,
    buyer.id as buyer_id,
    buyer.name as buyer_name,
    seller.id as seller_id,
    seller.name as seller_name,
    (SELECT COUNT(*) FROM marketplace_messages WHERE chat_id = c.id) as total_messages,
    (SELECT MIN(created_at) FROM marketplace_messages WHERE chat_id = c.id) as first_message_at,
    (SELECT MAX(created_at) FROM marketplace_messages WHERE chat_id = c.id) as last_message_at,
    (SELECT COUNT(*) FROM marketplace_messages WHERE chat_id = c.id AND has_attachments = true) as messages_with_attachments
FROM marketplace_chats c
LEFT JOIN marketplace_listings l ON c.listing_id = l.id
LEFT JOIN users buyer ON c.buyer_id = buyer.id
LEFT JOIN users seller ON c.seller_id = seller.id
WHERE c.id = :chat_id;

-- 9. Сообщения с вложениями
SELECT 
    m.id as message_id,
    m.chat_id,
    m.content,
    m.has_attachments,
    m.attachments_count,
    m.created_at,
    sender.name as sender_name,
    ca.file_name,
    ca.file_type,
    ca.file_size
FROM marketplace_messages m
JOIN users sender ON m.sender_id = sender.id
LEFT JOIN chat_attachments ca ON ca.message_id = m.id
WHERE m.has_attachments = true
ORDER BY m.created_at DESC
LIMIT 20;