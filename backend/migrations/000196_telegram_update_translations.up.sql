-- Update existing language_saved to language_changed translations
UPDATE translations 
SET field_name = 'language_changed',
    translated_text = CASE 
        WHEN language = 'en' THEN 'Language successfully changed to English! 🇬🇧'
        WHEN language = 'ru' THEN 'Язык успешно изменен на русский! 🇷🇺'
        WHEN language = 'sr' THEN 'Језик је успешно промењен на српски! 🇷🇸'
    END
WHERE entity_type = 'telegram_bot' 
AND entity_id = 3 
AND field_name = 'language_saved';

-- Add missing translations that were not created yet
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
-- use_link
('telegram_bot', 50, 'en', 'use_link', 'Please use the connection link from the application'),
('telegram_bot', 50, 'ru', 'use_link', 'Пожалуйста, используйте ссылку для подключения из приложения'),
('telegram_bot', 50, 'sr', 'use_link', 'Молимо користите везу за повезивање из апликације'),

-- token_error
('telegram_bot', 51, 'en', 'token_error', 'Token validation error. Please try again.'),
('telegram_bot', 51, 'ru', 'token_error', 'Ошибка валидации токена. Пожалуйста, попробуйте снова.'),
('telegram_bot', 51, 'sr', 'token_error', 'Грешка валидације токена. Молимо покушајте поново.'),

-- connected
('telegram_bot', 52, 'en', 'connected', 'Notifications successfully connected!'),
('telegram_bot', 52, 'ru', 'connected', 'Уведомления успешно подключены!'),
('telegram_bot', 52, 'sr', 'connected', 'Обавештења су успешно повезана!'),

-- new_message (update existing ID 37 notification_new_message)
('telegram_bot', 53, 'en', 'new_message', '💬 New message'),
('telegram_bot', 53, 'ru', 'new_message', '💬 Новое сообщение'),
('telegram_bot', 53, 'sr', 'new_message', '💬 Нова порука'),

-- new_review
('telegram_bot', 54, 'en', 'new_review', '⭐ New review'),
('telegram_bot', 54, 'ru', 'new_review', '⭐ Новый отзыв'),
('telegram_bot', 54, 'sr', 'new_review', '⭐ Нова рецензија'),

-- listing_update
('telegram_bot', 55, 'en', 'listing_update', '📝 Listing update'),
('telegram_bot', 55, 'ru', 'listing_update', '📝 Обновление объявления'),
('telegram_bot', 55, 'sr', 'listing_update', '📝 Ажурирање огласа'),

-- price_change
('telegram_bot', 56, 'en', 'price_change', '💰 Price changed'),
('telegram_bot', 56, 'ru', 'price_change', '💰 Изменение цены'),
('telegram_bot', 56, 'sr', 'price_change', '💰 Промена цене'),

-- btn_contact
('telegram_bot', 57, 'en', 'btn_contact', 'Contact'),
('telegram_bot', 57, 'ru', 'btn_contact', 'Связаться'),
('telegram_bot', 57, 'sr', 'btn_contact', 'Контакт'),

-- btn_add_favorite
('telegram_bot', 58, 'en', 'btn_add_favorite', 'Add to favorites'),
('telegram_bot', 58, 'ru', 'btn_add_favorite', 'В избранное'),
('telegram_bot', 58, 'sr', 'btn_add_favorite', 'У омиљене'),

-- cmd_start
('telegram_bot', 59, 'en', 'cmd_start', 'Start the bot'),
('telegram_bot', 59, 'ru', 'cmd_start', 'Запустить бота'),
('telegram_bot', 59, 'sr', 'cmd_start', 'Покрени бота'),

-- cmd_language
('telegram_bot', 60, 'en', 'cmd_language', 'Change language'),
('telegram_bot', 60, 'ru', 'cmd_language', 'Изменить язык'),
('telegram_bot', 60, 'sr', 'cmd_language', 'Промени језик'),

-- view_in_browser
('telegram_bot', 61, 'en', 'view_in_browser', 'View in browser'),
('telegram_bot', 61, 'ru', 'view_in_browser', 'Посмотреть в браузере'),
('telegram_bot', 61, 'sr', 'view_in_browser', 'Погледај у прегледачу'),

-- car_search_help
('telegram_bot', 62, 'en', 'car_search_help', 'To search for cars, send the make and model. For example: BMW X5'),
('telegram_bot', 62, 'ru', 'car_search_help', 'Для поиска автомобилей отправьте марку и модель. Например: BMW X5'),
('telegram_bot', 62, 'sr', 'car_search_help', 'За претрагу аутомобила пошаљите марку и модел. На пример: BMW X5'),

-- car_found
('telegram_bot', 63, 'en', 'car_found', 'Found %d cars'),
('telegram_bot', 63, 'ru', 'car_found', 'Найдено автомобилей: %d'),
('telegram_bot', 63, 'sr', 'car_found', 'Пронађено аутомобила: %d'),

-- car_details
('telegram_bot', 64, 'en', 'car_details', '🚗 %s %s\n📅 Year: %d\n🛣️ Mileage: %d km\n💰 Price: %s'),
('telegram_bot', 64, 'ru', 'car_details', '🚗 %s %s\n📅 Год: %d\n🛣️ Пробег: %d км\n💰 Цена: %s'),
('telegram_bot', 64, 'sr', 'car_details', '🚗 %s %s\n📅 Година: %d\n🛣️ Километража: %d км\n💰 Цена: %s')
ON CONFLICT (entity_type, entity_id, language, field_name) DO UPDATE 
SET translated_text = EXCLUDED.translated_text,
    updated_at = CURRENT_TIMESTAMP;

-- Delete duplicates from entity_id 101-122
DELETE FROM translations 
WHERE entity_type = 'telegram_bot' 
AND entity_id BETWEEN 101 AND 122;