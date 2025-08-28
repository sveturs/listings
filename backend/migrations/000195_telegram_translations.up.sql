-- Migration: Add Telegram bot translations
-- Author: System
-- Date: 2025-08-01

-- Telegram bot translations use entity_type = 'telegram_bot'
-- entity_id is a numeric ID for grouping
-- field_name contains the actual translation key
-- language contains the language code

-- Start with entity_id = 100 to avoid conflicts with existing translations
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
-- English translations
('telegram_bot', 100, 'en', 'choose_language', 'Please choose your preferred language:'),
('telegram_bot', 101, 'en', 'language_changed', 'Language successfully changed to English! 🇬🇧'),
('telegram_bot', 102, 'en', 'welcome', 'Welcome to SveTu bot! Here you can receive notifications about new messages, price changes and other important events.'),
('telegram_bot', 103, 'en', 'use_link', 'Please use the connection link from the application'),
('telegram_bot', 104, 'en', 'token_error', 'Token validation error. Please try again.'),
('telegram_bot', 105, 'en', 'connected', 'Notifications successfully connected!'),
('telegram_bot', 106, 'en', 'new_message', '💬 New message'),
('telegram_bot', 107, 'en', 'new_review', '⭐ New review'),
('telegram_bot', 108, 'en', 'listing_update', '📝 Listing update'),
('telegram_bot', 109, 'en', 'price_change', '💰 Price changed'),
('telegram_bot', 110, 'en', 'btn_view_details', 'View details'),
('telegram_bot', 111, 'en', 'btn_contact', 'Contact'),
('telegram_bot', 112, 'en', 'btn_add_favorite', 'Add to favorites'),
('telegram_bot', 113, 'en', 'btn_share', 'Share'),
('telegram_bot', 114, 'en', 'cmd_start', 'Start the bot'),
('telegram_bot', 115, 'en', 'cmd_help', 'Show help'),
('telegram_bot', 116, 'en', 'cmd_search', 'Search listings'),
('telegram_bot', 117, 'en', 'cmd_settings', 'Notification settings'),
('telegram_bot', 118, 'en', 'cmd_language', 'Change language'),
('telegram_bot', 119, 'en', 'view_in_browser', 'View in browser'),
('telegram_bot', 120, 'en', 'car_search_help', 'To search for cars, send the make and model. For example: BMW X5'),
('telegram_bot', 121, 'en', 'car_found', 'Found %d cars'),
('telegram_bot', 122, 'en', 'car_details', '🚗 %s %s\n📅 Year: %d\n🛣️ Mileage: %d km\n💰 Price: %s'),

-- Russian translations
('telegram_bot', 100, 'ru', 'choose_language', 'Пожалуйста, выберите предпочитаемый язык:'),
('telegram_bot', 101, 'ru', 'language_changed', 'Язык успешно изменен на русский! 🇷🇺'),
('telegram_bot', 102, 'ru', 'welcome', 'Добро пожаловать в бот SveTu! Здесь вы можете получать уведомления о новых сообщениях, изменениях цен и других важных событиях.'),
('telegram_bot', 103, 'ru', 'use_link', 'Пожалуйста, используйте ссылку для подключения из приложения'),
('telegram_bot', 104, 'ru', 'token_error', 'Ошибка валидации токена. Пожалуйста, попробуйте снова.'),
('telegram_bot', 105, 'ru', 'connected', 'Уведомления успешно подключены!'),
('telegram_bot', 106, 'ru', 'new_message', '💬 Новое сообщение'),
('telegram_bot', 107, 'ru', 'new_review', '⭐ Новый отзыв'),
('telegram_bot', 108, 'ru', 'listing_update', '📝 Обновление объявления'),
('telegram_bot', 109, 'ru', 'price_change', '💰 Изменение цены'),
('telegram_bot', 110, 'ru', 'btn_view_details', 'Подробнее'),
('telegram_bot', 111, 'ru', 'btn_contact', 'Связаться'),
('telegram_bot', 112, 'ru', 'btn_add_favorite', 'В избранное'),
('telegram_bot', 113, 'ru', 'btn_share', 'Поделиться'),
('telegram_bot', 114, 'ru', 'cmd_start', 'Запустить бота'),
('telegram_bot', 115, 'ru', 'cmd_help', 'Показать помощь'),
('telegram_bot', 116, 'ru', 'cmd_search', 'Поиск объявлений'),
('telegram_bot', 117, 'ru', 'cmd_settings', 'Настройки уведомлений'),
('telegram_bot', 118, 'ru', 'cmd_language', 'Изменить язык'),
('telegram_bot', 119, 'ru', 'view_in_browser', 'Посмотреть в браузере'),
('telegram_bot', 120, 'ru', 'car_search_help', 'Для поиска автомобилей отправьте марку и модель. Например: BMW X5'),
('telegram_bot', 121, 'ru', 'car_found', 'Найдено автомобилей: %d'),
('telegram_bot', 122, 'ru', 'car_details', '🚗 %s %s\n📅 Год: %d\n🛣️ Пробег: %d км\n💰 Цена: %s'),

-- Serbian translations
('telegram_bot', 100, 'sr', 'choose_language', 'Молимо изаберите жељени језик:'),
('telegram_bot', 101, 'sr', 'language_changed', 'Језик је успешно промењен на српски! 🇷🇸'),
('telegram_bot', 102, 'sr', 'welcome', 'Добродошли у SveTu бот! Овде можете примати обавештења о новим порукама, променама цена и другим важним догађајима.'),
('telegram_bot', 103, 'sr', 'use_link', 'Молимо користите везу за повезивање из апликације'),
('telegram_bot', 104, 'sr', 'token_error', 'Грешка валидације токена. Молимо покушајте поново.'),
('telegram_bot', 105, 'sr', 'connected', 'Обавештења су успешно повезана!'),
('telegram_bot', 106, 'sr', 'new_message', '💬 Нова порука'),
('telegram_bot', 107, 'sr', 'new_review', '⭐ Нова рецензија'),
('telegram_bot', 108, 'sr', 'listing_update', '📝 Ажурирање огласа'),
('telegram_bot', 109, 'sr', 'price_change', '💰 Промена цене'),
('telegram_bot', 110, 'sr', 'btn_view_details', 'Детаљи'),
('telegram_bot', 111, 'sr', 'btn_contact', 'Контакт'),
('telegram_bot', 112, 'sr', 'btn_add_favorite', 'У омиљене'),
('telegram_bot', 113, 'sr', 'btn_share', 'Подели'),
('telegram_bot', 114, 'sr', 'cmd_start', 'Покрени бота'),
('telegram_bot', 115, 'sr', 'cmd_help', 'Прикажи помоћ'),
('telegram_bot', 116, 'sr', 'cmd_search', 'Претрага огласа'),
('telegram_bot', 117, 'sr', 'cmd_settings', 'Подешавања обавештења'),
('telegram_bot', 118, 'sr', 'cmd_language', 'Промени језик'),
('telegram_bot', 119, 'sr', 'view_in_browser', 'Погледај у прегледачу'),
('telegram_bot', 120, 'sr', 'car_search_help', 'За претрагу аутомобила пошаљите марку и модел. На пример: BMW X5'),
('telegram_bot', 121, 'sr', 'car_found', 'Пронађено аутомобила: %d'),
('telegram_bot', 122, 'sr', 'car_details', '🚗 %s %s\n📅 Година: %d\n🛣️ Километража: %d км\n💰 Цена: %s')
ON CONFLICT (entity_type, entity_id, language, field_name) DO UPDATE 
SET translated_text = EXCLUDED.translated_text,
    updated_at = CURRENT_TIMESTAMP;