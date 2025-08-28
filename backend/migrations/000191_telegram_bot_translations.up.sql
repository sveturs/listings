-- Добавляем переводы для Telegram бота
-- Используем таблицу translations с entity_type = 'telegram_bot'

-- Базовые команды и сообщения
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_verified, is_machine_translated) VALUES
-- Приветственные сообщения
('telegram_bot', 1, 'ru', 'welcome', 'Добро пожаловать в Sve Tu бот! 🎉\n\nЯ помогу вам найти и разместить объявления о продаже автомобилей.', true, false),
('telegram_bot', 1, 'sr', 'welcome', 'Dobrodošli u Sve Tu bot! 🎉\n\nPomoći ću vam da pronađete i postavite oglase za prodaju automobila.', true, false),
('telegram_bot', 1, 'en', 'welcome', 'Welcome to Sve Tu bot! 🎉\n\nI will help you find and post car sale listings.', true, false),

-- Выбор языка
('telegram_bot', 2, 'ru', 'choose_language', 'Пожалуйста, выберите язык:', true, false),
('telegram_bot', 2, 'sr', 'choose_language', 'Molimo izaberite jezik:', true, false),
('telegram_bot', 2, 'en', 'choose_language', 'Please choose language:', true, false),

-- Язык сохранен
('telegram_bot', 3, 'ru', 'language_saved', 'Язык успешно изменен на русский 🇷🇺', true, false),
('telegram_bot', 3, 'sr', 'language_saved', 'Jezik je uspešno promenjen na srpski 🇷🇸', true, false),
('telegram_bot', 3, 'en', 'language_saved', 'Language successfully changed to English 🇬🇧', true, false),

-- Главное меню
('telegram_bot', 4, 'ru', 'main_menu', 'Главное меню', true, false),
('telegram_bot', 4, 'sr', 'main_menu', 'Glavni meni', true, false),
('telegram_bot', 4, 'en', 'main_menu', 'Main menu', true, false),

-- Команды
('telegram_bot', 5, 'ru', 'cmd_search', '🔍 Поиск автомобилей', true, false),
('telegram_bot', 5, 'sr', 'cmd_search', '🔍 Pretraga automobila', true, false),
('telegram_bot', 5, 'en', 'cmd_search', '🔍 Search cars', true, false),

('telegram_bot', 6, 'ru', 'cmd_my_listings', '📋 Мои объявления', true, false),
('telegram_bot', 6, 'sr', 'cmd_my_listings', '📋 Moji oglasi', true, false),
('telegram_bot', 6, 'en', 'cmd_my_listings', '📋 My listings', true, false),

('telegram_bot', 7, 'ru', 'cmd_create_listing', '➕ Создать объявление', true, false),
('telegram_bot', 7, 'sr', 'cmd_create_listing', '➕ Kreiraj oglas', true, false),
('telegram_bot', 7, 'en', 'cmd_create_listing', '➕ Create listing', true, false),

('telegram_bot', 8, 'ru', 'cmd_settings', '⚙️ Настройки', true, false),
('telegram_bot', 8, 'sr', 'cmd_settings', '⚙️ Podešavanja', true, false),
('telegram_bot', 8, 'en', 'cmd_settings', '⚙️ Settings', true, false),

('telegram_bot', 9, 'ru', 'cmd_help', '❓ Помощь', true, false),
('telegram_bot', 9, 'sr', 'cmd_help', '❓ Pomoć', true, false),
('telegram_bot', 9, 'en', 'cmd_help', '❓ Help', true, false),

-- Поиск автомобилей
('telegram_bot', 10, 'ru', 'search_prompt', 'Что вы ищете? Введите марку, модель или другие параметры:', true, false),
('telegram_bot', 10, 'sr', 'search_prompt', 'Šta tražite? Unesite marku, model ili druge parametre:', true, false),
('telegram_bot', 10, 'en', 'search_prompt', 'What are you looking for? Enter make, model or other parameters:', true, false),

('telegram_bot', 11, 'ru', 'search_results_found', 'Найдено объявлений: %d', true, false),
('telegram_bot', 11, 'sr', 'search_results_found', 'Pronađeno oglasa: %d', true, false),
('telegram_bot', 11, 'en', 'search_results_found', 'Listings found: %d', true, false),

('telegram_bot', 12, 'ru', 'search_no_results', 'По вашему запросу ничего не найдено. Попробуйте изменить параметры поиска.', true, false),
('telegram_bot', 12, 'sr', 'search_no_results', 'Nema rezultata za vašu pretragu. Pokušajte promeniti parametre pretrage.', true, false),
('telegram_bot', 12, 'en', 'search_no_results', 'No results found for your search. Try changing search parameters.', true, false),

-- Создание объявления
('telegram_bot', 13, 'ru', 'create_listing_start', 'Давайте создадим новое объявление! Сначала выберите марку автомобиля:', true, false),
('telegram_bot', 13, 'sr', 'create_listing_start', 'Hajde da kreiramo novi oglas! Prvo izaberite marku automobila:', true, false),
('telegram_bot', 13, 'en', 'create_listing_start', 'Let''s create a new listing! First, choose the car make:', true, false),

('telegram_bot', 14, 'ru', 'select_model', 'Отлично! Теперь выберите модель:', true, false),
('telegram_bot', 14, 'sr', 'select_model', 'Odlično! Sada izaberite model:', true, false),
('telegram_bot', 14, 'en', 'select_model', 'Great! Now choose the model:', true, false),

('telegram_bot', 15, 'ru', 'enter_year', 'Введите год выпуска:', true, false),
('telegram_bot', 15, 'sr', 'enter_year', 'Unesite godinu proizvodnje:', true, false),
('telegram_bot', 15, 'en', 'enter_year', 'Enter production year:', true, false),

('telegram_bot', 16, 'ru', 'enter_price', 'Введите цену в EUR:', true, false),
('telegram_bot', 16, 'sr', 'enter_price', 'Unesite cenu u EUR:', true, false),
('telegram_bot', 16, 'en', 'enter_price', 'Enter price in EUR:', true, false),

('telegram_bot', 17, 'ru', 'enter_mileage', 'Введите пробег в км:', true, false),
('telegram_bot', 17, 'sr', 'enter_mileage', 'Unesite kilometražu:', true, false),
('telegram_bot', 17, 'en', 'enter_mileage', 'Enter mileage in km:', true, false),

('telegram_bot', 18, 'ru', 'select_fuel_type', 'Выберите тип топлива:', true, false),
('telegram_bot', 18, 'sr', 'select_fuel_type', 'Izaberite tip goriva:', true, false),
('telegram_bot', 18, 'en', 'select_fuel_type', 'Select fuel type:', true, false),

('telegram_bot', 19, 'ru', 'select_transmission', 'Выберите тип трансмиссии:', true, false),
('telegram_bot', 19, 'sr', 'select_transmission', 'Izaberite tip menjača:', true, false),
('telegram_bot', 19, 'en', 'select_transmission', 'Select transmission type:', true, false),

('telegram_bot', 20, 'ru', 'upload_photos', 'Отправьте фотографии автомобиля (до 10 штук):', true, false),
('telegram_bot', 20, 'sr', 'upload_photos', 'Pošaljite fotografije automobila (do 10 komada):', true, false),
('telegram_bot', 20, 'en', 'upload_photos', 'Send car photos (up to 10):', true, false),

('telegram_bot', 21, 'ru', 'enter_description', 'Введите описание автомобиля:', true, false),
('telegram_bot', 21, 'sr', 'enter_description', 'Unesite opis automobila:', true, false),
('telegram_bot', 21, 'en', 'enter_description', 'Enter car description:', true, false),

('telegram_bot', 22, 'ru', 'listing_created', '✅ Объявление успешно создано!\n\nID: #%d\nПосмотреть: %s', true, false),
('telegram_bot', 22, 'sr', 'listing_created', '✅ Oglas je uspešno kreiran!\n\nID: #%d\nPogledaj: %s', true, false),
('telegram_bot', 22, 'en', 'listing_created', '✅ Listing created successfully!\n\nID: #%d\nView: %s', true, false),

-- Ошибки
('telegram_bot', 23, 'ru', 'error_generic', '❌ Произошла ошибка. Пожалуйста, попробуйте позже.', true, false),
('telegram_bot', 23, 'sr', 'error_generic', '❌ Došlo je do greške. Molimo pokušajte kasnije.', true, false),
('telegram_bot', 23, 'en', 'error_generic', '❌ An error occurred. Please try again later.', true, false),

('telegram_bot', 24, 'ru', 'error_invalid_input', '❌ Неверный ввод. Пожалуйста, попробуйте еще раз.', true, false),
('telegram_bot', 24, 'sr', 'error_invalid_input', '❌ Neispravan unos. Molimo pokušajte ponovo.', true, false),
('telegram_bot', 24, 'en', 'error_invalid_input', '❌ Invalid input. Please try again.', true, false),

('telegram_bot', 25, 'ru', 'error_not_connected', '❌ Ваш Telegram аккаунт не связан с учетной записью Sve Tu. Используйте команду /start с кодом привязки.', true, false),
('telegram_bot', 25, 'sr', 'error_not_connected', '❌ Vaš Telegram nalog nije povezan sa Sve Tu nalogom. Koristite komandu /start sa kodom za povezivanje.', true, false),
('telegram_bot', 25, 'en', 'error_not_connected', '❌ Your Telegram account is not connected to Sve Tu account. Use /start command with connection code.', true, false),

-- Кнопки действий
('telegram_bot', 26, 'ru', 'btn_back', '◀️ Назад', true, false),
('telegram_bot', 26, 'sr', 'btn_back', '◀️ Nazad', true, false),
('telegram_bot', 26, 'en', 'btn_back', '◀️ Back', true, false),

('telegram_bot', 27, 'ru', 'btn_cancel', '❌ Отмена', true, false),
('telegram_bot', 27, 'sr', 'btn_cancel', '❌ Otkaži', true, false),
('telegram_bot', 27, 'en', 'btn_cancel', '❌ Cancel', true, false),

('telegram_bot', 28, 'ru', 'btn_next', 'Далее ▶️', true, false),
('telegram_bot', 28, 'sr', 'btn_next', 'Dalje ▶️', true, false),
('telegram_bot', 28, 'en', 'btn_next', 'Next ▶️', true, false),

('telegram_bot', 29, 'ru', 'btn_skip', 'Пропустить ⏭️', true, false),
('telegram_bot', 29, 'sr', 'btn_skip', 'Preskoči ⏭️', true, false),
('telegram_bot', 29, 'en', 'btn_skip', 'Skip ⏭️', true, false),

-- Топливо
('telegram_bot', 30, 'ru', 'fuel_petrol', 'Бензин', true, false),
('telegram_bot', 30, 'sr', 'fuel_petrol', 'Benzin', true, false),
('telegram_bot', 30, 'en', 'fuel_petrol', 'Petrol', true, false),

('telegram_bot', 31, 'ru', 'fuel_diesel', 'Дизель', true, false),
('telegram_bot', 31, 'sr', 'fuel_diesel', 'Dizel', true, false),
('telegram_bot', 31, 'en', 'fuel_diesel', 'Diesel', true, false),

('telegram_bot', 32, 'ru', 'fuel_electric', 'Электро', true, false),
('telegram_bot', 32, 'sr', 'fuel_electric', 'Električni', true, false),
('telegram_bot', 32, 'en', 'fuel_electric', 'Electric', true, false),

('telegram_bot', 33, 'ru', 'fuel_hybrid', 'Гибрид', true, false),
('telegram_bot', 33, 'sr', 'fuel_hybrid', 'Hibrid', true, false),
('telegram_bot', 33, 'en', 'fuel_hybrid', 'Hybrid', true, false),

('telegram_bot', 34, 'ru', 'fuel_lpg', 'Газ', true, false),
('telegram_bot', 34, 'sr', 'fuel_lpg', 'Gas', true, false),
('telegram_bot', 34, 'en', 'fuel_lpg', 'LPG', true, false),

-- Трансмиссия
('telegram_bot', 35, 'ru', 'trans_manual', 'Механика', true, false),
('telegram_bot', 35, 'sr', 'trans_manual', 'Manuelni', true, false),
('telegram_bot', 35, 'en', 'trans_manual', 'Manual', true, false),

('telegram_bot', 36, 'ru', 'trans_automatic', 'Автомат', true, false),
('telegram_bot', 36, 'sr', 'trans_automatic', 'Automatik', true, false),
('telegram_bot', 36, 'en', 'trans_automatic', 'Automatic', true, false),

-- Уведомления
('telegram_bot', 37, 'ru', 'notification_new_message', '💬 Новое сообщение по объявлению #%d:\n\nОт: %s\n\n%s', true, false),
('telegram_bot', 37, 'sr', 'notification_new_message', '💬 Nova poruka za oglas #%d:\n\nOd: %s\n\n%s', true, false),
('telegram_bot', 37, 'en', 'notification_new_message', '💬 New message for listing #%d:\n\nFrom: %s\n\n%s', true, false),

('telegram_bot', 38, 'ru', 'notification_price_alert', '💰 Изменение цены!\n\n%s %s %d\nНовая цена: €%d (было €%d)', true, false),
('telegram_bot', 38, 'sr', 'notification_price_alert', '💰 Promena cene!\n\n%s %s %d\nNova cena: €%d (bila €%d)', true, false),
('telegram_bot', 38, 'en', 'notification_price_alert', '💰 Price change!\n\n%s %s %d\nNew price: €%d (was €%d)', true, false),

-- Помощь
('telegram_bot', 39, 'ru', 'help_text', '📖 Доступные команды:\n\n/start - Начать работу\n/search - Поиск автомобилей\n/create - Создать объявление\n/my_listings - Мои объявления\n/language - Изменить язык\n/help - Эта справка\n\nПо всем вопросам: @svetu_support', true, false),
('telegram_bot', 39, 'sr', 'help_text', '📖 Dostupne komande:\n\n/start - Počni\n/search - Pretraga automobila\n/create - Kreiraj oglas\n/my_listings - Moji oglasi\n/language - Promeni jezik\n/help - Ova pomoć\n\nZa sva pitanja: @svetu_support', true, false),
('telegram_bot', 39, 'en', 'help_text', '📖 Available commands:\n\n/start - Start\n/search - Search cars\n/create - Create listing\n/my_listings - My listings\n/language - Change language\n/help - This help\n\nFor all questions: @svetu_support', true, false),

-- Inline кнопки для объявлений
('telegram_bot', 40, 'ru', 'btn_view_details', '👁️ Подробнее', true, false),
('telegram_bot', 40, 'sr', 'btn_view_details', '👁️ Detalji', true, false),
('telegram_bot', 40, 'en', 'btn_view_details', '👁️ View details', true, false),

('telegram_bot', 41, 'ru', 'btn_contact_seller', '📞 Связаться', true, false),
('telegram_bot', 41, 'sr', 'btn_contact_seller', '📞 Kontakt', true, false),
('telegram_bot', 41, 'en', 'btn_contact_seller', '📞 Contact', true, false),

('telegram_bot', 42, 'ru', 'btn_share', '📤 Поделиться', true, false),
('telegram_bot', 42, 'sr', 'btn_share', '📤 Podeli', true, false),
('telegram_bot', 42, 'en', 'btn_share', '📤 Share', true, false);

-- Создаем функцию для получения переводов Telegram бота
CREATE OR REPLACE FUNCTION get_telegram_translation(
    p_key VARCHAR,
    p_language VARCHAR DEFAULT 'ru'
) RETURNS TEXT AS $$
DECLARE
    v_translation TEXT;
BEGIN
    SELECT translated_text INTO v_translation
    FROM translations
    WHERE entity_type = 'telegram_bot'
    AND field_name = p_key
    AND language = p_language
    LIMIT 1;
    
    -- Если перевод не найден, пробуем английский
    IF v_translation IS NULL AND p_language != 'en' THEN
        SELECT translated_text INTO v_translation
        FROM translations
        WHERE entity_type = 'telegram_bot'
        AND field_name = p_key
        AND language = 'en'
        LIMIT 1;
    END IF;
    
    -- Если и английский не найден, возвращаем ключ
    IF v_translation IS NULL THEN
        RETURN p_key;
    END IF;
    
    RETURN v_translation;
END;
$$ LANGUAGE plpgsql;