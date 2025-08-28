-- Миграция для добавления складской системы и самовывоза
-- Поддержка Fulfillment by Sve Tu (FBS) и складских операций

-- =====================================================
-- 1. ТАБЛИЦА СКЛАДСКИХ ПОМЕЩЕНИЙ
-- =====================================================
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    
    -- Основные данные
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) DEFAULT 'main',  -- main, temporary, partner
    
    -- Адрес
    address VARCHAR(500) NOT NULL,
    city VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(2) DEFAULT 'RS',
    
    -- Контакты
    phone VARCHAR(50),
    email VARCHAR(200),
    manager_name VARCHAR(200),
    manager_phone VARCHAR(50),
    
    -- Координаты
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- График работы (JSON)
    working_hours JSONB,  -- {"monday": {"open": "09:00", "close": "19:00"}, ...}
    
    -- Характеристики
    total_area_m2 DECIMAL(10,2),
    storage_area_m2 DECIMAL(10,2),
    max_capacity_m3 DECIMAL(10,2),
    current_occupancy_m3 DECIMAL(10,2) DEFAULT 0,
    
    -- Возможности
    supports_fbs BOOLEAN DEFAULT true,  -- Fulfillment by Sve Tu
    supports_pickup BOOLEAN DEFAULT true,  -- Самовывоз
    has_refrigeration BOOLEAN DEFAULT false,
    has_loading_dock BOOLEAN DEFAULT true,
    
    -- Статус
    is_active BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Добавляем основной склад
INSERT INTO warehouses (code, name, address, city, postal_code, phone, working_hours) VALUES (
    'NS-MAIN-01',
    'Главный склад Sve Tu - Нови Сад',
    'Улица Микија Манојловића 53',
    'Нови Сад',
    '21000',
    '+381 21 XXX-XXXX',
    '{"monday": {"open": "09:00", "close": "19:00"}, "tuesday": {"open": "09:00", "close": "19:00"}, "wednesday": {"open": "09:00", "close": "19:00"}, "thursday": {"open": "09:00", "close": "19:00"}, "friday": {"open": "09:00", "close": "19:00"}, "saturday": {"open": "10:00", "close": "16:00"}, "sunday": null}'
);

-- =====================================================
-- 2. ТАБЛИЦА СКЛАДСКИХ ТОВАРОВ
-- =====================================================
CREATE TABLE warehouse_inventory (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id) NOT NULL,
    
    -- Связь с товаром
    storefront_product_id BIGINT REFERENCES storefront_products(id),
    marketplace_listing_id INTEGER REFERENCES marketplace_listings(id),
    
    -- SKU и идентификация
    sku VARCHAR(100) NOT NULL,
    barcode VARCHAR(100),
    external_id VARCHAR(100),  -- ID в системе продавца
    
    -- Описание товара
    product_name VARCHAR(500) NOT NULL,
    product_description TEXT,
    
    -- Количество
    quantity_total INTEGER NOT NULL DEFAULT 0,  -- Общее количество
    quantity_available INTEGER NOT NULL DEFAULT 0,  -- Доступно для продажи
    quantity_reserved INTEGER NOT NULL DEFAULT 0,  -- Зарезервировано в заказах
    quantity_damaged INTEGER DEFAULT 0,  -- Поврежденные
    
    -- Размеры и вес единицы товара
    unit_weight_kg DECIMAL(10,3),
    unit_length_cm DECIMAL(10,2),
    unit_width_cm DECIMAL(10,2),
    unit_height_cm DECIMAL(10,2),
    unit_volume_m3 DECIMAL(10,4),
    
    -- Расположение на складе
    location_zone VARCHAR(50),  -- A, B, C...
    location_rack VARCHAR(50),  -- 01, 02, 03...
    location_shelf VARCHAR(50),  -- 1, 2, 3...
    location_bin VARCHAR(50),   -- A1, A2, B1...
    
    -- Даты
    received_at TIMESTAMP WITH TIME ZONE,
    expiry_date DATE,  -- Для товаров с ограниченным сроком
    
    -- Стоимость хранения
    storage_fee_daily DECIMAL(10,2) DEFAULT 0,
    
    -- Флаги
    is_fragile BOOLEAN DEFAULT false,
    requires_refrigeration BOOLEAN DEFAULT false,
    is_hazardous BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_warehouse_inventory_warehouse ON warehouse_inventory(warehouse_id);
CREATE INDEX idx_warehouse_inventory_sku ON warehouse_inventory(sku);
CREATE INDEX idx_warehouse_inventory_storefront_product ON warehouse_inventory(storefront_product_id);
CREATE INDEX idx_warehouse_inventory_marketplace_listing ON warehouse_inventory(marketplace_listing_id);
CREATE INDEX idx_warehouse_inventory_available ON warehouse_inventory(quantity_available);

-- =====================================================
-- 3. ТАБЛИЦА ДВИЖЕНИЙ ТОВАРА НА СКЛАДЕ
-- =====================================================
CREATE TABLE warehouse_movements (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id) NOT NULL,
    inventory_id INTEGER REFERENCES warehouse_inventory(id),
    
    -- Тип движения
    movement_type VARCHAR(50) NOT NULL,  -- inbound, outbound, transfer, adjustment, return
    movement_reason VARCHAR(100),  -- order_fulfillment, stock_replenishment, damaged, expired
    
    -- Количество
    quantity INTEGER NOT NULL,
    quantity_before INTEGER NOT NULL,
    quantity_after INTEGER NOT NULL,
    
    -- Связанные документы
    order_id INTEGER REFERENCES marketplace_orders(id),
    storefront_order_id BIGINT REFERENCES storefront_orders(id),
    shipment_id INTEGER REFERENCES post_express_shipments(id),
    
    -- Документы
    document_number VARCHAR(100),
    document_type VARCHAR(50),  -- invoice, delivery_note, return_note
    
    -- Исполнитель
    performed_by INTEGER REFERENCES users(id),
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_warehouse_movements_warehouse ON warehouse_movements(warehouse_id);
CREATE INDEX idx_warehouse_movements_inventory ON warehouse_movements(inventory_id);
CREATE INDEX idx_warehouse_movements_type ON warehouse_movements(movement_type);
CREATE INDEX idx_warehouse_movements_created ON warehouse_movements(created_at);

-- =====================================================
-- 4. ТАБЛИЦА ЗАКАЗОВ НА САМОВЫВОЗ
-- =====================================================
CREATE TABLE warehouse_pickup_orders (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id) NOT NULL,
    
    -- Связь с заказами
    marketplace_order_id INTEGER REFERENCES marketplace_orders(id),
    storefront_order_id BIGINT REFERENCES storefront_orders(id),
    
    -- Код самовывоза
    pickup_code VARCHAR(10) NOT NULL UNIQUE,  -- 6-значный код
    qr_code_url VARCHAR(500),  -- URL QR-кода
    
    -- Статусы
    status VARCHAR(50) DEFAULT 'pending',  -- pending, ready, picked_up, expired, cancelled
    ready_at TIMESTAMP WITH TIME ZONE,
    picked_up_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    
    -- Получатель
    customer_name VARCHAR(200) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    customer_email VARCHAR(200),
    
    -- Подтверждение получения
    pickup_confirmed_by VARCHAR(200),  -- Сотрудник склада
    id_document_type VARCHAR(50),  -- passport, id_card
    id_document_number VARCHAR(100),
    signature_url VARCHAR(500),  -- URL подписи получателя
    
    -- Уведомления
    notification_sent_at TIMESTAMP WITH TIME ZONE,
    reminder_sent_at TIMESTAMP WITH TIME ZONE,
    
    -- Дополнительно
    notes TEXT,
    pickup_photo_url VARCHAR(500),  -- Фото момента выдачи
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_warehouse_pickup_orders_code ON warehouse_pickup_orders(pickup_code);
CREATE INDEX idx_warehouse_pickup_orders_status ON warehouse_pickup_orders(status);
CREATE INDEX idx_warehouse_pickup_orders_expires ON warehouse_pickup_orders(expires_at);
CREATE INDEX idx_warehouse_pickup_orders_warehouse ON warehouse_pickup_orders(warehouse_id);

-- =====================================================
-- 5. ТАБЛИЦА НАСТРОЕК FBS ДЛЯ ВИТРИН
-- =====================================================
CREATE TABLE storefront_fbs_settings (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER REFERENCES storefronts(id) UNIQUE NOT NULL,
    
    -- Статус FBS
    fbs_enabled BOOLEAN DEFAULT false,
    auto_fulfillment BOOLEAN DEFAULT true,  -- Автоматическая обработка
    
    -- Настройки хранения
    storage_tier VARCHAR(50) DEFAULT 'standard',  -- standard, premium, economy
    max_storage_volume_m3 DECIMAL(10,2),
    free_storage_days INTEGER DEFAULT 30,
    
    -- Настройки упаковки
    default_packaging VARCHAR(50) DEFAULT 'standard',
    use_branded_packaging BOOLEAN DEFAULT false,
    include_invoice BOOLEAN DEFAULT true,
    include_marketing_materials BOOLEAN DEFAULT false,
    
    -- Правила обработки
    min_processing_hours INTEGER DEFAULT 24,  -- Минимальное время обработки
    cutoff_time TIME DEFAULT '15:00',  -- Время отсечки для отправки в тот же день
    process_on_weekends BOOLEAN DEFAULT false,
    
    -- Биллинг
    billing_cycle VARCHAR(50) DEFAULT 'monthly',  -- weekly, monthly
    storage_fee_per_m3_daily DECIMAL(10,2) DEFAULT 50,
    picking_fee_per_item DECIMAL(10,2) DEFAULT 30,
    packing_fee_per_order DECIMAL(10,2) DEFAULT 50,
    
    -- Статистика
    total_orders_fulfilled INTEGER DEFAULT 0,
    total_items_stored INTEGER DEFAULT 0,
    current_storage_used_m3 DECIMAL(10,2) DEFAULT 0,
    last_billed_at TIMESTAMP WITH TIME ZONE,
    current_charges DECIMAL(12,2) DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =====================================================
-- 6. ТАБЛИЦА СЧЕТОВ ЗА СКЛАДСКИЕ УСЛУГИ
-- =====================================================
CREATE TABLE warehouse_invoices (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER REFERENCES storefronts(id) NOT NULL,
    warehouse_id INTEGER REFERENCES warehouses(id) NOT NULL,
    
    -- Период
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    -- Детали счета
    invoice_number VARCHAR(100) UNIQUE NOT NULL,
    invoice_date DATE NOT NULL,
    due_date DATE NOT NULL,
    
    -- Позиции (JSON)
    line_items JSONB NOT NULL,  -- [{description, quantity, unit_price, total}, ...]
    
    -- Суммы
    subtotal DECIMAL(12,2) NOT NULL,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL,
    
    -- Статус
    status VARCHAR(50) DEFAULT 'draft',  -- draft, sent, paid, overdue, cancelled
    paid_at TIMESTAMP WITH TIME ZONE,
    payment_method VARCHAR(50),
    payment_reference VARCHAR(100),
    
    -- Документы
    pdf_url VARCHAR(500),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_warehouse_invoices_storefront ON warehouse_invoices(storefront_id);
CREATE INDEX idx_warehouse_invoices_status ON warehouse_invoices(status);
CREATE INDEX idx_warehouse_invoices_period ON warehouse_invoices(period_start, period_end);

-- =====================================================
-- 7. ФУНКЦИИ И ТРИГГЕРЫ
-- =====================================================

-- Функция генерации кода самовывоза
CREATE OR REPLACE FUNCTION generate_pickup_code()
RETURNS VARCHAR AS $$
DECLARE
    chars VARCHAR := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result VARCHAR := '';
    i INTEGER;
BEGIN
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Функция автоматического резервирования товара при создании заказа
CREATE OR REPLACE FUNCTION reserve_inventory_on_order()
RETURNS TRIGGER AS $$
DECLARE
    inv_record RECORD;
BEGIN
    -- Находим товар на складе
    SELECT * INTO inv_record
    FROM warehouse_inventory
    WHERE (storefront_product_id = NEW.product_id OR marketplace_listing_id = NEW.listing_id)
    AND quantity_available >= NEW.quantity
    LIMIT 1;
    
    IF FOUND THEN
        -- Резервируем товар
        UPDATE warehouse_inventory
        SET 
            quantity_available = quantity_available - NEW.quantity,
            quantity_reserved = quantity_reserved + NEW.quantity,
            updated_at = NOW()
        WHERE id = inv_record.id;
        
        -- Создаем запись о движении
        INSERT INTO warehouse_movements (
            warehouse_id, inventory_id, movement_type, movement_reason,
            quantity, quantity_before, quantity_after,
            order_id, storefront_order_id
        ) VALUES (
            inv_record.warehouse_id, inv_record.id, 'outbound', 'order_fulfillment',
            -NEW.quantity, inv_record.quantity_available, inv_record.quantity_available - NEW.quantity,
            NEW.marketplace_order_id, NEW.storefront_order_id
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Функция отслеживания изменений inventory
CREATE OR REPLACE FUNCTION track_inventory_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.quantity_total != NEW.quantity_total OR 
       OLD.quantity_available != NEW.quantity_available OR
       OLD.quantity_reserved != NEW.quantity_reserved THEN
        
        INSERT INTO warehouse_movements (
            warehouse_id, inventory_id, movement_type, movement_reason,
            quantity, quantity_before, quantity_after
        ) VALUES (
            NEW.warehouse_id, NEW.id, 'adjustment', 'manual_adjustment',
            NEW.quantity_total - OLD.quantity_total,
            OLD.quantity_total, NEW.quantity_total
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер на изменение inventory
CREATE TRIGGER track_warehouse_inventory_changes
    AFTER UPDATE ON warehouse_inventory
    FOR EACH ROW EXECUTE FUNCTION track_inventory_changes();

-- Функция обновления updated_at для складских таблиц
CREATE OR REPLACE FUNCTION update_warehouse_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Применяем триггеры updated_at
CREATE TRIGGER update_warehouses_updated_at 
    BEFORE UPDATE ON warehouses 
    FOR EACH ROW EXECUTE FUNCTION update_warehouse_updated_at();

CREATE TRIGGER update_warehouse_inventory_updated_at 
    BEFORE UPDATE ON warehouse_inventory 
    FOR EACH ROW EXECUTE FUNCTION update_warehouse_updated_at();

CREATE TRIGGER update_warehouse_pickup_orders_updated_at 
    BEFORE UPDATE ON warehouse_pickup_orders 
    FOR EACH ROW EXECUTE FUNCTION update_warehouse_updated_at();

CREATE TRIGGER update_storefront_fbs_settings_updated_at 
    BEFORE UPDATE ON storefront_fbs_settings 
    FOR EACH ROW EXECUTE FUNCTION update_warehouse_updated_at();

CREATE TRIGGER update_warehouse_invoices_updated_at 
    BEFORE UPDATE ON warehouse_invoices 
    FOR EACH ROW EXECUTE FUNCTION update_warehouse_updated_at();

-- =====================================================
-- 8. КОММЕНТАРИИ К ТАБЛИЦАМ
-- =====================================================

COMMENT ON TABLE warehouses IS 'Складские помещения платформы';
COMMENT ON TABLE warehouse_inventory IS 'Товары на складах';
COMMENT ON TABLE warehouse_movements IS 'История движения товаров на складе';
COMMENT ON TABLE warehouse_pickup_orders IS 'Заказы на самовывоз со склада';
COMMENT ON TABLE storefront_fbs_settings IS 'Настройки Fulfillment by Sve Tu для витрин';
COMMENT ON TABLE warehouse_invoices IS 'Счета за складские услуги';