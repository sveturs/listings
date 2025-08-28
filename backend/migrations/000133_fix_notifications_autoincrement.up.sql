-- Исправление автоинкремента для таблицы notifications
ALTER TABLE notifications 
ALTER COLUMN id SET DEFAULT nextval('notifications_id_seq');

ALTER SEQUENCE notifications_id_seq OWNED BY notifications.id;

SELECT setval('notifications_id_seq', COALESCE((SELECT MAX(id) FROM notifications), 1), true);