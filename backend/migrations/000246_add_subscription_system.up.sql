-- Создание системы подписок для витрин

-- Таблица тарифных планов
CREATE TABLE IF NOT EXISTS subscription_plans (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE, -- starter, professional, business, enterprise
    name VARCHAR(100) NOT NULL,
    price_monthly NUMERIC(10,2) DEFAULT 0, -- месячная цена в EUR
    price_yearly NUMERIC(10,2) DEFAULT 0, -- годовая цена в EUR
    
    -- Лимиты плана
    max_storefronts INTEGER DEFAULT 1, -- -1 для unlimited
    max_products_per_storefront INTEGER DEFAULT 50, -- -1 для unlimited  
    max_staff_per_storefront INTEGER DEFAULT 1, -- -1 для unlimited
    max_images_total INTEGER DEFAULT 100, -- -1 для unlimited
    
    -- Функции плана
    has_ai_assistant BOOLEAN DEFAULT FALSE,
    has_live_shopping BOOLEAN DEFAULT FALSE,
    has_export_data BOOLEAN DEFAULT FALSE,
    has_custom_domain BOOLEAN DEFAULT FALSE,
    has_analytics BOOLEAN DEFAULT TRUE,
    has_priority_support BOOLEAN DEFAULT FALSE,
    
    -- Дополнительные параметры
    commission_rate NUMERIC(5,2) DEFAULT 10.00, -- комиссия маркетплейса в %
    free_trial_days INTEGER DEFAULT 0,
    
    sort_order INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    is_recommended BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица подписок пользователей
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL REFERENCES subscription_plans(id),
    
    status VARCHAR(50) DEFAULT 'active', -- active, trial, expired, cancelled, suspended
    billing_cycle VARCHAR(20) DEFAULT 'monthly', -- monthly, yearly
    
    -- Даты подписки
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    trial_ends_at TIMESTAMP,
    current_period_start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    current_period_end TIMESTAMP NOT NULL,
    cancelled_at TIMESTAMP,
    expires_at TIMESTAMP,
    
    -- Платежная информация
    last_payment_id INTEGER REFERENCES payment_transactions(id),
    last_payment_at TIMESTAMP,
    next_payment_at TIMESTAMP,
    payment_method VARCHAR(50), -- card, balance, bank_transfer
    auto_renew BOOLEAN DEFAULT TRUE,
    
    -- Использование лимитов
    used_storefronts INTEGER DEFAULT 0,
    
    -- Метаданные
    metadata JSONB DEFAULT '{}',
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id) -- один пользователь = одна активная подписка
);

-- Таблица истории подписок (для отслеживания изменений)
CREATE TABLE IF NOT EXISTS subscription_history (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    action VARCHAR(50) NOT NULL, -- created, upgraded, downgraded, renewed, cancelled, expired
    from_plan_id INTEGER REFERENCES subscription_plans(id),
    to_plan_id INTEGER REFERENCES subscription_plans(id),
    
    reason TEXT,
    metadata JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id)
);

-- Таблица платежей за подписки
CREATE TABLE IF NOT EXISTS subscription_payments (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payment_id INTEGER REFERENCES payment_transactions(id),
    
    amount NUMERIC(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'EUR',
    
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed, refunded
    payment_method VARCHAR(50),
    
    transaction_data JSONB DEFAULT '{}',
    
    paid_at TIMESTAMP,
    failed_at TIMESTAMP,
    refunded_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица использования лимитов
CREATE TABLE IF NOT EXISTS subscription_usage (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    storefront_id INTEGER REFERENCES storefronts(id) ON DELETE CASCADE,
    
    resource_type VARCHAR(50) NOT NULL, -- storefront, product, staff, image
    resource_id INTEGER,
    resource_count INTEGER DEFAULT 1,
    
    action VARCHAR(50) NOT NULL, -- created, deleted
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Вставка тарифных планов
INSERT INTO subscription_plans (code, name, price_monthly, price_yearly, 
    max_storefronts, max_products_per_storefront, max_staff_per_storefront, max_images_total,
    has_ai_assistant, has_live_shopping, has_export_data, has_custom_domain,
    commission_rate, free_trial_days, sort_order, is_recommended)
VALUES 
    ('starter', 'Starter', 0, 0, 
     1, 50, 1, 100,
     FALSE, FALSE, FALSE, FALSE,
     15.00, 0, 1, FALSE),
     
    ('professional', 'Professional', 29.00, 290.00,
     3, 500, 5, 1000,
     FALSE, FALSE, TRUE, TRUE,
     10.00, 14, 2, TRUE),
     
    ('business', 'Business', 99.00, 990.00,
     10, 5000, 20, 10000,
     TRUE, TRUE, TRUE, TRUE,
     7.00, 14, 3, FALSE),
     
    ('enterprise', 'Enterprise', NULL, NULL,
     -1, -1, -1, -1,
     TRUE, TRUE, TRUE, TRUE,
     5.00, 30, 4, FALSE);

-- Добавление полей подписки к витринам
ALTER TABLE storefronts
ADD COLUMN IF NOT EXISTS subscription_id INTEGER REFERENCES user_subscriptions(id),
ADD COLUMN IF NOT EXISTS is_subscription_active BOOLEAN DEFAULT TRUE;

-- Функция для проверки лимитов подписки
CREATE OR REPLACE FUNCTION check_subscription_limits(
    p_user_id INTEGER,
    p_resource_type VARCHAR,
    p_count INTEGER DEFAULT 1
) RETURNS BOOLEAN AS $$
DECLARE
    v_subscription RECORD;
    v_plan RECORD;
    v_current_usage INTEGER;
BEGIN
    -- Получаем активную подписку пользователя
    SELECT us.*, sp.*
    INTO v_subscription
    FROM user_subscriptions us
    JOIN subscription_plans sp ON us.plan_id = sp.id
    WHERE us.user_id = p_user_id 
    AND us.status IN ('active', 'trial')
    LIMIT 1;
    
    IF NOT FOUND THEN
        -- Если нет подписки, используем бесплатный план
        SELECT * INTO v_plan
        FROM subscription_plans
        WHERE code = 'starter'
        LIMIT 1;
        
        v_subscription.id = NULL;
    ELSE
        v_plan = v_subscription;
    END IF;
    
    -- Проверяем лимиты в зависимости от типа ресурса
    CASE p_resource_type
        WHEN 'storefront' THEN
            -- Для unlimited возвращаем true
            IF v_plan.max_storefronts = -1 THEN
                RETURN TRUE;
            END IF;
            
            -- Считаем текущие витрины
            SELECT COUNT(*) INTO v_current_usage
            FROM storefronts
            WHERE user_id = p_user_id
            AND is_active = TRUE;
            
            RETURN (v_current_usage + p_count) <= v_plan.max_storefronts;
            
        WHEN 'product' THEN
            IF v_plan.max_products_per_storefront = -1 THEN
                RETURN TRUE;
            END IF;
            -- Логика для продуктов будет добавлена позже
            RETURN TRUE;
            
        WHEN 'staff' THEN
            IF v_plan.max_staff_per_storefront = -1 THEN
                RETURN TRUE;
            END IF;
            -- Логика для персонала будет добавлена позже
            RETURN TRUE;
            
        ELSE
            RETURN FALSE;
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- Функция для получения текущей подписки пользователя
CREATE OR REPLACE FUNCTION get_user_subscription(p_user_id INTEGER)
RETURNS TABLE (
    subscription_id INTEGER,
    plan_code VARCHAR,
    plan_name VARCHAR,
    status VARCHAR,
    expires_at TIMESTAMP,
    max_storefronts INTEGER,
    used_storefronts INTEGER,
    max_products INTEGER,
    max_staff INTEGER,
    max_images INTEGER,
    has_ai BOOLEAN,
    has_live BOOLEAN,
    has_export BOOLEAN,
    has_custom_domain BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        us.id as subscription_id,
        sp.code as plan_code,
        sp.name as plan_name,
        us.status,
        us.expires_at,
        sp.max_storefronts,
        us.used_storefronts,
        sp.max_products_per_storefront as max_products,
        sp.max_staff_per_storefront as max_staff,
        sp.max_images_total as max_images,
        sp.has_ai_assistant as has_ai,
        sp.has_live_shopping as has_live,
        sp.has_export_data as has_export,
        sp.has_custom_domain
    FROM user_subscriptions us
    JOIN subscription_plans sp ON us.plan_id = sp.id
    WHERE us.user_id = p_user_id
    AND us.status IN ('active', 'trial')
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- Триггер для обновления used_storefronts при создании/удалении витрины
CREATE OR REPLACE FUNCTION update_subscription_usage() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- При создании витрины увеличиваем счетчик
        UPDATE user_subscriptions
        SET used_storefronts = used_storefronts + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = NEW.user_id;
        
        -- Записываем в историю использования
        INSERT INTO subscription_usage (subscription_id, storefront_id, resource_type, action)
        SELECT id, NEW.id, 'storefront', 'created'
        FROM user_subscriptions
        WHERE user_id = NEW.user_id
        LIMIT 1;
        
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении витрины уменьшаем счетчик
        UPDATE user_subscriptions
        SET used_storefronts = GREATEST(0, used_storefronts - 1),
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = OLD.user_id;
        
        -- Записываем в историю использования
        INSERT INTO subscription_usage (subscription_id, storefront_id, resource_type, action)
        SELECT id, OLD.id, 'storefront', 'deleted'
        FROM user_subscriptions
        WHERE user_id = OLD.user_id
        LIMIT 1;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_storefront_usage
    AFTER INSERT OR DELETE ON storefronts
    FOR EACH ROW
    EXECUTE FUNCTION update_subscription_usage();

-- Индексы для производительности
CREATE INDEX idx_user_subscriptions_user_status ON user_subscriptions(user_id, status);
CREATE INDEX idx_user_subscriptions_expires ON user_subscriptions(expires_at);
CREATE INDEX idx_subscription_payments_subscription ON subscription_payments(subscription_id);
CREATE INDEX idx_subscription_payments_status ON subscription_payments(status);
CREATE INDEX idx_subscription_history_subscription ON subscription_history(subscription_id);
CREATE INDEX idx_subscription_history_user ON subscription_history(user_id);
CREATE INDEX idx_subscription_usage_subscription ON subscription_usage(subscription_id);
CREATE INDEX idx_subscription_usage_storefront ON subscription_usage(storefront_id);
CREATE INDEX idx_subscription_usage_created ON subscription_usage(created_at);

-- Комментарии к таблицам
COMMENT ON TABLE subscription_plans IS 'Тарифные планы для витрин';
COMMENT ON TABLE user_subscriptions IS 'Активные подписки пользователей';
COMMENT ON TABLE subscription_history IS 'История изменений подписок';
COMMENT ON TABLE subscription_payments IS 'Платежи за подписки';
COMMENT ON TABLE subscription_usage IS 'Использование ресурсов подписки';