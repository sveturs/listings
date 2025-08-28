-- Rollback extended RBAC system

-- Drop view
DROP VIEW IF EXISTS user_role_permissions;

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_log_role_change ON users;
DROP FUNCTION IF EXISTS log_role_change();

-- Drop audit table
DROP TABLE IF EXISTS role_audit_log;

-- Remove added columns from roles table
ALTER TABLE roles DROP COLUMN IF EXISTS is_system;
ALTER TABLE roles DROP COLUMN IF EXISTS is_assignable;
ALTER TABLE roles DROP COLUMN IF EXISTS priority;

-- Delete role-permission associations for new roles
DELETE FROM role_permissions 
WHERE role_id IN (
    SELECT id FROM roles 
    WHERE name IN (
        'content_moderator', 'review_moderator', 'chat_moderator', 'dispute_manager',
        'vendor_manager', 'category_manager', 'marketing_manager', 'financial_manager',
        'warehouse_manager', 'warehouse_worker', 'pickup_manager', 'pickup_worker', 'courier',
        'support_l1', 'support_l2', 'support_l3', 'legal_advisor', 'compliance_officer',
        'professional_vendor', 'individual_seller', 'storefront_staff',
        'verified_buyer', 'vip_customer', 'data_analyst', 'business_analyst'
    )
);

-- Delete new permissions
DELETE FROM permissions 
WHERE name IN (
    'users.block', 'users.verify', 'users.assign_role', 'users.view_details', 'users.export',
    'listings.view_all', 'listings.feature', 'listings.bulk_edit', 'listings.approve', 'listings.reject',
    'orders.view_all', 'orders.view_own', 'orders.process', 'orders.cancel', 'orders.refund', 'orders.export',
    'storefronts.verify', 'storefronts.suspend', 'storefronts.manage_staff', 'storefronts.view_analytics', 'storefronts.bulk_upload',
    'finance.view_transactions', 'finance.view_reports', 'finance.manage_payouts', 'finance.set_commissions', 'finance.export',
    'warehouse.view_inventory', 'warehouse.manage_inventory', 'warehouse.receive_goods', 'warehouse.ship_goods', 'warehouse.transfer_stock',
    'categories.view', 'categories.create', 'categories.edit', 'categories.delete', 'categories.manage_attributes',
    'marketing.create_campaigns', 'marketing.manage_promotions', 'marketing.send_emails', 'marketing.manage_banners',
    'analytics.view_dashboard', 'analytics.export_data', 'analytics.create_reports', 'analytics.view_sensitive',
    'system.view_logs', 'system.manage_settings', 'system.manage_roles', 'system.view_audit', 'system.manage_integrations',
    'support.view_tickets', 'support.handle_tickets', 'support.escalate_tickets', 'support.close_tickets',
    'disputes.view_all', 'disputes.handle', 'disputes.resolve', 'disputes.escalate'
);

-- Delete new roles
DELETE FROM roles 
WHERE name IN (
    'content_moderator', 'review_moderator', 'chat_moderator', 'dispute_manager',
    'vendor_manager', 'category_manager', 'marketing_manager', 'financial_manager',
    'warehouse_manager', 'warehouse_worker', 'pickup_manager', 'pickup_worker', 'courier',
    'support_l1', 'support_l2', 'support_l3', 'legal_advisor', 'compliance_officer',
    'professional_vendor', 'individual_seller', 'storefront_staff',
    'verified_buyer', 'vip_customer', 'data_analyst', 'business_analyst'
);