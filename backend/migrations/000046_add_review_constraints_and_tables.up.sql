-- Добавляем уникальный индекс для защиты от спама
-- Один пользователь может оставить только один отзыв на одну сущность
CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_user_entity_unique 
ON reviews(user_id, entity_type, entity_id) 
WHERE status != 'deleted';

-- Создаем таблицу для подтверждений отзывов продавцами
CREATE TABLE review_confirmations (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    confirmed_by INTEGER NOT NULL REFERENCES users(id),
    confirmation_status VARCHAR(50) NOT NULL CHECK (confirmation_status IN ('confirmed', 'disputed')),
    confirmed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT,
    
    -- Уникальный индекс - один отзыв может быть подтвержден только один раз
    UNIQUE(review_id)
);

-- Создаем таблицу для споров по отзывам
CREATE TABLE review_disputes (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    disputed_by INTEGER NOT NULL REFERENCES users(id),
    dispute_reason VARCHAR(100) NOT NULL CHECK (dispute_reason IN (
        'not_a_customer', -- не является покупателем
        'false_information', -- ложная информация
        'deal_cancelled', -- сделка не состоялась
        'spam', -- спам
        'other' -- другое
    )),
    dispute_description TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending', -- ожидает рассмотрения
        'in_review', -- на рассмотрении
        'resolved_keep_review', -- отзыв остается
        'resolved_remove_review', -- отзыв удаляется
        'resolved_remove_verification', -- отзыв остается, но без верификации
        'cancelled' -- спор отменен
    )),
    admin_id INTEGER REFERENCES users(id), -- администратор, рассматривающий спор
    admin_notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP
);

-- Создаем индексы для таблицы споров
CREATE INDEX IF NOT EXISTS idx_disputes_status ON review_disputes(status);
CREATE INDEX IF NOT EXISTS idx_disputes_review_id ON review_disputes(review_id);

-- Создаем таблицу для сообщений в спорах
CREATE TABLE review_dispute_messages (
    id SERIAL PRIMARY KEY,
    dispute_id INTEGER NOT NULL REFERENCES review_disputes(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    message TEXT NOT NULL,
    attachments JSONB, -- массив ссылок на прикрепленные файлы
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создаем индекс для таблицы сообщений
CREATE INDEX IF NOT EXISTS idx_dispute_messages_dispute_id ON review_dispute_messages(dispute_id);

-- Добавляем поля в таблицу reviews для поддержки новой функциональности
ALTER TABLE reviews
ADD COLUMN IF NOT EXISTS seller_confirmed BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS has_active_dispute BOOLEAN DEFAULT false;

-- Создаем индексы для новых полей
CREATE INDEX IF NOT EXISTS idx_reviews_seller_confirmed ON reviews(seller_confirmed) WHERE seller_confirmed = true;
CREATE INDEX IF NOT EXISTS idx_reviews_has_dispute ON reviews(has_active_dispute) WHERE has_active_dispute = true;

-- Добавляем комментарии для документации
COMMENT ON TABLE review_confirmations IS 'Подтверждения отзывов продавцами';
COMMENT ON TABLE review_disputes IS 'Споры по отзывам между покупателями и продавцами';
COMMENT ON TABLE review_dispute_messages IS 'Сообщения в рамках рассмотрения споров';
COMMENT ON COLUMN reviews.seller_confirmed IS 'Флаг подтверждения отзыва продавцом';
COMMENT ON COLUMN reviews.has_active_dispute IS 'Флаг наличия активного спора по отзыву';