-- Add Post Express delivery providers to all existing storefronts
-- This migration adds comprehensive Post Express delivery options with proper settings

UPDATE storefronts 
SET settings = jsonb_set(
    COALESCE(settings, '{}'::jsonb),
    '{delivery_providers}',
    '[
        {
            "id": "pickup",
            "name": "🏪 Самовывоз",
            "icon": "🏪",
            "enabled": true,
            "description": "Покупатели могут забрать товар самостоятельно",
            "settings": {
                "pickup_address": "Novi Sad, Serbia",
                "working_hours": "9:00-20:00"
            }
        },
        {
            "id": "local_delivery",
            "name": "🚲 Локальная доставка",
            "icon": "🚲",
            "enabled": true,
            "description": "Доставка курьером в пределах города",
            "settings": {
                "base_rate": 0,
                "free_shipping_threshold": 0,
                "delivery_radius": 15,
                "estimated_days": "1-2"
            }
        },
        {
            "id": "post_express",
            "name": "📮 Post Express",
            "icon": "📮",
            "enabled": true,
            "description": "Национальная почта Сербии - доставка по всей стране с полной интеграцией",
            "settings": {
                "api_enabled": true,
                "estimated_days": 1,
                "weight_tiers": {
                    "0-2kg": 340,
                    "2-5kg": 450,
                    "5-10kg": 580,
                    "10-20kg": 790
                },
                "free_shipping_threshold": 5000,
                "cod_fee": 45,
                "insurance_included": 15000,
                "insurance_rate": 0.01
            }
        },
        {
            "id": "post_express_office",
            "name": "📬 Post Express - Почтовое отделение",
            "icon": "📬",
            "enabled": true,
            "description": "Доставка в ближайшее почтовое отделение Post Express",
            "settings": {
                "discount_percent": 10,
                "estimated_days": "1-2",
                "office_network": "180+ отделений по всей Сербии"
            }
        },
        {
            "id": "post_express_warehouse",
            "name": "📦 Post Express - Складской самовывоз",
            "icon": "📦",
            "enabled": true,
            "description": "Самовывоз из центрального склада Post Express в Нови-Саде",
            "settings": {
                "warehouse_address": "Novi Sad, Bulevar oslobođenja 127",
                "working_hours": "8:00-20:00",
                "free_shipping_threshold": 2000,
                "qr_code_enabled": true,
                "try_before_buy": true
            }
        },
        {
            "id": "post_express_express",
            "name": "⚡ Post Express - Экспресс доставка",
            "icon": "⚡",
            "enabled": true,
            "description": "Срочная доставка курьером за 1 день",
            "settings": {
                "estimated_hours": 24,
                "express_surcharge": 200,
                "available_cities": ["Belgrade", "Novi Sad", "Niš", "Kragujevac"],
                "cutoff_time": "14:00"
            }
        },
        {
            "id": "bex_courier",
            "name": "🚚 BEX Express - Курьерская доставка",
            "icon": "🚚",
            "enabled": false,
            "description": "Альтернативная курьерская служба BEX",
            "settings": {
                "base_rate": 350,
                "weight_based_pricing": true,
                "estimated_days": "2-3"
            }
        },
        {
            "id": "bex_pickup_point",
            "name": "📍 BEX Express - Пункт выдачи",
            "icon": "📍",
            "enabled": false,
            "description": "Самовывоз из пунктов выдачи BEX",
            "settings": {
                "discount_percent": 20,
                "pickup_points": 50
            }
        },
        {
            "id": "bex_warehouse_pickup",
            "name": "🏭 BEX Express - Склад",
            "icon": "🏭",
            "enabled": false,
            "description": "Самовывоз со склада BEX",
            "settings": {
                "warehouse_address": "Belgrade, Autoput 22",
                "free_pickup": true
            }
        }
    ]'::jsonb,
    true
)
WHERE settings IS NULL 
   OR settings->'delivery_providers' IS NULL 
   OR jsonb_array_length(COALESCE(settings->'delivery_providers', '[]'::jsonb)) = 0;

-- Also update storefronts that may have incomplete delivery providers
UPDATE storefronts 
SET settings = jsonb_set(
    settings,
    '{delivery_providers}',
    '[
        {
            "id": "pickup",
            "name": "🏪 Самовывоз",
            "icon": "🏪",
            "enabled": true,
            "description": "Покупатели могут забрать товар самостоятельно",
            "settings": {
                "pickup_address": "Novi Sad, Serbia",
                "working_hours": "9:00-20:00"
            }
        },
        {
            "id": "local_delivery",
            "name": "🚲 Локальная доставка",
            "icon": "🚲",
            "enabled": true,
            "description": "Доставка курьером в пределах города",
            "settings": {
                "base_rate": 0,
                "free_shipping_threshold": 0,
                "delivery_radius": 15,
                "estimated_days": "1-2"
            }
        },
        {
            "id": "post_express",
            "name": "📮 Post Express",
            "icon": "📮",
            "enabled": true,
            "description": "Национальная почта Сербии - доставка по всей стране с полной интеграцией",
            "settings": {
                "api_enabled": true,
                "estimated_days": 1,
                "weight_tiers": {
                    "0-2kg": 340,
                    "2-5kg": 450,
                    "5-10kg": 580,
                    "10-20kg": 790
                },
                "free_shipping_threshold": 5000,
                "cod_fee": 45,
                "insurance_included": 15000,
                "insurance_rate": 0.01
            }
        },
        {
            "id": "post_express_office",
            "name": "📬 Post Express - Почтовое отделение",
            "icon": "📬",
            "enabled": true,
            "description": "Доставка в ближайшее почтовое отделение Post Express",
            "settings": {
                "discount_percent": 10,
                "estimated_days": "1-2",
                "office_network": "180+ отделений по всей Сербии"
            }
        },
        {
            "id": "post_express_warehouse",
            "name": "📦 Post Express - Складской самовывоз",
            "icon": "📦",
            "enabled": true,
            "description": "Самовывоз из центрального склада Post Express в Нови-Саде",
            "settings": {
                "warehouse_address": "Novi Sad, Bulevar oslobođenja 127",
                "working_hours": "8:00-20:00",
                "free_shipping_threshold": 2000,
                "qr_code_enabled": true,
                "try_before_buy": true
            }
        },
        {
            "id": "post_express_express",
            "name": "⚡ Post Express - Экспресс доставка",
            "icon": "⚡",
            "enabled": true,
            "description": "Срочная доставка курьером за 1 день",
            "settings": {
                "estimated_hours": 24,
                "express_surcharge": 200,
                "available_cities": ["Belgrade", "Novi Sad", "Niš", "Kragujevac"],
                "cutoff_time": "14:00"
            }
        },
        {
            "id": "bex_courier",
            "name": "🚚 BEX Express - Курьерская доставка",
            "icon": "🚚",
            "enabled": false,
            "description": "Альтернативная курьерская служба BEX",
            "settings": {
                "base_rate": 350,
                "weight_based_pricing": true,
                "estimated_days": "2-3"
            }
        },
        {
            "id": "bex_pickup_point",
            "name": "📍 BEX Express - Пункт выдачи",
            "icon": "📍",
            "enabled": false,
            "description": "Самовывоз из пунктов выдачи BEX",
            "settings": {
                "discount_percent": 20,
                "pickup_points": 50
            }
        },
        {
            "id": "bex_warehouse_pickup",
            "name": "🏭 BEX Express - Склад",
            "icon": "🏭",
            "enabled": false,
            "description": "Самовывоз со склада BEX",
            "settings": {
                "warehouse_address": "Belgrade, Autoput 22",
                "free_pickup": true
            }
        }
    ]'::jsonb,
    true
)
WHERE jsonb_array_length(COALESCE(settings->'delivery_providers', '[]'::jsonb)) < 5;