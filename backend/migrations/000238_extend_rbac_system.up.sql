-- Extended RBAC System for Sve Tu Marketplace
-- This migration extends the existing role-based access control system with comprehensive roles and permissions

-- Add new roles for comprehensive marketplace management
INSERT INTO roles (name, display_name, description) VALUES
    -- Moderation roles
    ('content_moderator', 'Content Moderator', 'Moderates listings, products and content'),
    ('review_moderator', 'Review Moderator', 'Moderates reviews and ratings'),
    ('chat_moderator', 'Chat Moderator', 'Monitors and moderates chat conversations'),
    ('dispute_manager', 'Dispute Manager', 'Handles disputes between buyers and sellers'),
    
    -- Business roles
    ('vendor_manager', 'Vendor Manager', 'Manages vendor accounts and storefronts'),
    ('category_manager', 'Category Manager', 'Manages product categories and attributes'),
    ('marketing_manager', 'Marketing Manager', 'Manages promotions, discounts and marketing campaigns'),
    ('financial_manager', 'Financial Manager', 'Handles financial transactions and reports'),
    
    -- Operational roles
    ('warehouse_manager', 'Warehouse Manager', 'Manages warehouse operations and inventory'),
    ('warehouse_worker', 'Warehouse Worker', 'Handles order fulfillment and packaging'),
    ('pickup_manager', 'Pickup Point Manager', 'Manages pickup point operations'),
    ('pickup_worker', 'Pickup Point Worker', 'Handles order pickup and returns'),
    ('courier', 'Courier', 'Delivers orders to customers'),
    
    -- Support roles
    ('support_l1', 'Customer Support L1', 'Basic customer support and inquiries'),
    ('support_l2', 'Customer Support L2', 'Advanced support and issue resolution'),
    ('support_l3', 'Customer Support L3', 'Expert support and escalations'),
    ('legal_advisor', 'Legal Advisor', 'Handles legal matters and compliance'),
    ('compliance_officer', 'Compliance Officer', 'Ensures regulatory compliance'),
    
    -- Seller roles
    ('professional_vendor', 'Professional Vendor', 'Professional seller with storefront and advanced features'),
    ('individual_seller', 'Individual Seller', 'Individual seller with basic features'),
    ('storefront_staff', 'Storefront Staff', 'Staff member of a storefront with limited permissions'),
    
    -- Buyer roles
    ('verified_buyer', 'Verified Buyer', 'Verified customer with full buying privileges'),
    ('vip_customer', 'VIP Customer', 'VIP customer with exclusive benefits'),
    
    -- Analytics roles
    ('data_analyst', 'Data Analyst', 'Access to analytics and reporting tools'),
    ('business_analyst', 'Business Analyst', 'Business process analysis and optimization')
ON CONFLICT (name) DO NOTHING;

-- Add comprehensive permissions
INSERT INTO permissions (name, resource, action, description) VALUES
    -- Extended user management permissions
    ('users.block', 'users', 'block', 'Block/unblock users'),
    ('users.verify', 'users', 'verify', 'Verify user accounts'),
    ('users.assign_role', 'users', 'assign_role', 'Assign roles to users'),
    ('users.view_details', 'users', 'view_details', 'View detailed user information'),
    ('users.export', 'users', 'export', 'Export user data'),
    
    -- Extended listings permissions
    ('listings.view_all', 'listings', 'view_all', 'View all listings including private'),
    ('listings.feature', 'listings', 'feature', 'Feature/unfeature listings'),
    ('listings.bulk_edit', 'listings', 'bulk_edit', 'Bulk edit listings'),
    ('listings.approve', 'listings', 'approve', 'Approve pending listings'),
    ('listings.reject', 'listings', 'reject', 'Reject listings'),
    
    -- Orders management
    ('orders.view_all', 'orders', 'view_all', 'View all orders'),
    ('orders.view_own', 'orders', 'view_own', 'View own orders'),
    ('orders.process', 'orders', 'process', 'Process orders'),
    ('orders.cancel', 'orders', 'cancel', 'Cancel orders'),
    ('orders.refund', 'orders', 'refund', 'Process refunds'),
    ('orders.export', 'orders', 'export', 'Export order data'),
    
    -- Extended storefront permissions
    ('storefronts.verify', 'storefronts', 'verify', 'Verify storefronts'),
    ('storefronts.suspend', 'storefronts', 'suspend', 'Suspend storefronts'),
    ('storefronts.manage_staff', 'storefronts', 'manage_staff', 'Manage storefront staff'),
    ('storefronts.view_analytics', 'storefronts', 'view_analytics', 'View storefront analytics'),
    ('storefronts.bulk_upload', 'storefronts', 'bulk_upload', 'Bulk upload products'),
    
    -- Financial permissions
    ('finance.view_transactions', 'finance', 'view_transactions', 'View financial transactions'),
    ('finance.view_reports', 'finance', 'view_reports', 'View financial reports'),
    ('finance.manage_payouts', 'finance', 'manage_payouts', 'Manage seller payouts'),
    ('finance.set_commissions', 'finance', 'set_commissions', 'Set commission rates'),
    ('finance.export', 'finance', 'export', 'Export financial data'),
    
    -- Warehouse permissions
    ('warehouse.view_inventory', 'warehouse', 'view_inventory', 'View inventory levels'),
    ('warehouse.manage_inventory', 'warehouse', 'manage_inventory', 'Manage inventory'),
    ('warehouse.receive_goods', 'warehouse', 'receive_goods', 'Receive incoming goods'),
    ('warehouse.ship_goods', 'warehouse', 'ship_goods', 'Ship outgoing goods'),
    ('warehouse.transfer_stock', 'warehouse', 'transfer_stock', 'Transfer stock between locations'),
    
    -- Category management
    ('categories.view', 'categories', 'view', 'View categories'),
    ('categories.create', 'categories', 'create', 'Create categories'),
    ('categories.edit', 'categories', 'edit', 'Edit categories'),
    ('categories.delete', 'categories', 'delete', 'Delete categories'),
    ('categories.manage_attributes', 'categories', 'manage_attributes', 'Manage category attributes'),
    
    -- Marketing permissions
    ('marketing.create_campaigns', 'marketing', 'create_campaigns', 'Create marketing campaigns'),
    ('marketing.manage_promotions', 'marketing', 'manage_promotions', 'Manage promotions and discounts'),
    ('marketing.send_emails', 'marketing', 'send_emails', 'Send marketing emails'),
    ('marketing.manage_banners', 'marketing', 'manage_banners', 'Manage promotional banners'),
    
    -- Analytics permissions
    ('analytics.view_dashboard', 'analytics', 'view_dashboard', 'View analytics dashboard'),
    ('analytics.export_data', 'analytics', 'export_data', 'Export analytics data'),
    ('analytics.create_reports', 'analytics', 'create_reports', 'Create custom reports'),
    ('analytics.view_sensitive', 'analytics', 'view_sensitive', 'View sensitive business data'),
    
    -- System permissions
    ('system.view_logs', 'system', 'view_logs', 'View system logs'),
    ('system.manage_settings', 'system', 'manage_settings', 'Manage system settings'),
    ('system.manage_roles', 'system', 'manage_roles', 'Manage roles and permissions'),
    ('system.view_audit', 'system', 'view_audit', 'View audit logs'),
    ('system.manage_integrations', 'system', 'manage_integrations', 'Manage third-party integrations'),
    
    -- Support permissions
    ('support.view_tickets', 'support', 'view_tickets', 'View support tickets'),
    ('support.handle_tickets', 'support', 'handle_tickets', 'Handle support tickets'),
    ('support.escalate_tickets', 'support', 'escalate_tickets', 'Escalate tickets to higher level'),
    ('support.close_tickets', 'support', 'close_tickets', 'Close resolved tickets'),
    
    -- Dispute permissions
    ('disputes.view_all', 'disputes', 'view_all', 'View all disputes'),
    ('disputes.handle', 'disputes', 'handle', 'Handle disputes'),
    ('disputes.resolve', 'disputes', 'resolve', 'Resolve disputes'),
    ('disputes.escalate', 'disputes', 'escalate', 'Escalate disputes')
ON CONFLICT (name) DO NOTHING;

-- Assign permissions to new roles

-- Content Moderator
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'content_moderator' 
    AND p.name IN (
        'admin.access',
        'listings.view_all',
        'listings.moderate',
        'listings.approve',
        'listings.reject',
        'listings.edit_any',
        'listings.delete_any',
        'categories.view'
    )
ON CONFLICT DO NOTHING;

-- Review Moderator
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'review_moderator' 
    AND p.name IN (
        'admin.access',
        'reviews.moderate',
        'reviews.delete_any',
        'users.view'
    )
ON CONFLICT DO NOTHING;

-- Chat Moderator
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'chat_moderator' 
    AND p.name IN (
        'admin.access',
        'chat.view_any',
        'chat.moderate',
        'users.view',
        'users.block'
    )
ON CONFLICT DO NOTHING;

-- Dispute Manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'dispute_manager' 
    AND p.name IN (
        'admin.access',
        'disputes.view_all',
        'disputes.handle',
        'disputes.resolve',
        'disputes.escalate',
        'orders.view_all',
        'orders.refund',
        'chat.view_any',
        'users.view'
    )
ON CONFLICT DO NOTHING;

-- Vendor Manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'vendor_manager' 
    AND p.name IN (
        'admin.access',
        'storefronts.edit_any',
        'storefronts.verify',
        'storefronts.suspend',
        'storefronts.view_analytics',
        'users.view',
        'users.view_details',
        'finance.view_transactions',
        'finance.set_commissions'
    )
ON CONFLICT DO NOTHING;

-- Category Manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'category_manager' 
    AND p.name IN (
        'admin.access',
        'categories.view',
        'categories.create',
        'categories.edit',
        'categories.delete',
        'categories.manage_attributes',
        'admin.categories',
        'admin.attributes'
    )
ON CONFLICT DO NOTHING;

-- Marketing Manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'marketing_manager' 
    AND p.name IN (
        'admin.access',
        'marketing.create_campaigns',
        'marketing.manage_promotions',
        'marketing.send_emails',
        'marketing.manage_banners',
        'analytics.view_dashboard',
        'users.view'
    )
ON CONFLICT DO NOTHING;

-- Financial Manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'financial_manager' 
    AND p.name IN (
        'admin.access',
        'finance.view_transactions',
        'finance.view_reports',
        'finance.manage_payouts',
        'finance.set_commissions',
        'finance.export',
        'analytics.view_dashboard',
        'orders.view_all'
    )
ON CONFLICT DO NOTHING;

-- Warehouse Manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'warehouse_manager' 
    AND p.name IN (
        'warehouse.view_inventory',
        'warehouse.manage_inventory',
        'warehouse.receive_goods',
        'warehouse.ship_goods',
        'warehouse.transfer_stock',
        'orders.view_all',
        'orders.process'
    )
ON CONFLICT DO NOTHING;

-- Warehouse Worker
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'warehouse_worker' 
    AND p.name IN (
        'warehouse.view_inventory',
        'warehouse.receive_goods',
        'warehouse.ship_goods',
        'orders.view_own',
        'orders.process'
    )
ON CONFLICT DO NOTHING;

-- Professional Vendor
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'professional_vendor' 
    AND p.name IN (
        'listings.create',
        'listings.edit_own',
        'listings.delete_own',
        'listings.bulk_edit',
        'storefronts.create',
        'storefronts.edit_own',
        'storefronts.delete_own',
        'storefronts.manage_staff',
        'storefronts.view_analytics',
        'storefronts.bulk_upload',
        'orders.view_own',
        'orders.process',
        'finance.view_transactions',
        'analytics.view_dashboard',
        'chat.send',
        'reviews.create'
    )
ON CONFLICT DO NOTHING;

-- Support L1
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'support_l1' 
    AND p.name IN (
        'support.view_tickets',
        'support.handle_tickets',
        'chat.send',
        'users.view',
        'orders.view_all',
        'listings.view_all'
    )
ON CONFLICT DO NOTHING;

-- Support L2
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'support_l2' 
    AND p.name IN (
        'support.view_tickets',
        'support.handle_tickets',
        'support.escalate_tickets',
        'support.close_tickets',
        'orders.view_all',
        'orders.cancel',
        'orders.refund',
        'users.view',
        'users.view_details',
        'chat.send'
    )
ON CONFLICT DO NOTHING;

-- Data Analyst
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'data_analyst' 
    AND p.name IN (
        'analytics.view_dashboard',
        'analytics.export_data',
        'analytics.create_reports',
        'finance.view_reports',
        'orders.export',
        'users.export'
    )
ON CONFLICT DO NOTHING;

-- Add columns to track role metadata
ALTER TABLE roles ADD COLUMN IF NOT EXISTS is_system BOOLEAN DEFAULT FALSE;
ALTER TABLE roles ADD COLUMN IF NOT EXISTS is_assignable BOOLEAN DEFAULT TRUE;
ALTER TABLE roles ADD COLUMN IF NOT EXISTS priority INTEGER DEFAULT 100;

-- Mark system roles
UPDATE roles SET is_system = TRUE WHERE name IN ('super_admin', 'admin', 'user');

-- Set role priorities (lower number = higher priority)
UPDATE roles SET priority = CASE
    WHEN name = 'super_admin' THEN 1
    WHEN name = 'admin' THEN 10
    WHEN name = 'financial_manager' THEN 20
    WHEN name = 'vendor_manager' THEN 30
    WHEN name = 'dispute_manager' THEN 40
    WHEN name = 'content_moderator' THEN 50
    WHEN name = 'warehouse_manager' THEN 60
    WHEN name = 'professional_vendor' THEN 70
    WHEN name = 'vendor' THEN 80
    WHEN name = 'user' THEN 100
    ELSE 90
END;

-- Create audit table for role changes
CREATE TABLE IF NOT EXISTS role_audit_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    target_user_id INTEGER REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    old_role_id INTEGER REFERENCES roles(id),
    new_role_id INTEGER REFERENCES roles(id),
    details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create function to log role changes
CREATE OR REPLACE FUNCTION log_role_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND OLD.role_id IS DISTINCT FROM NEW.role_id THEN
        INSERT INTO role_audit_log (
            target_user_id,
            action,
            old_role_id,
            new_role_id,
            details
        ) VALUES (
            NEW.id,
            'role_changed',
            OLD.role_id,
            NEW.role_id,
            jsonb_build_object(
                'old_role', (SELECT name FROM roles WHERE id = OLD.role_id),
                'new_role', (SELECT name FROM roles WHERE id = NEW.role_id),
                'timestamp', CURRENT_TIMESTAMP
            )
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for role changes
DROP TRIGGER IF EXISTS trigger_log_role_change ON users;
CREATE TRIGGER trigger_log_role_change
    AFTER UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION log_role_change();

-- Create view for user roles with permissions
CREATE OR REPLACE VIEW user_role_permissions AS
SELECT 
    u.id as user_id,
    u.email,
    u.name as user_name,
    r.id as role_id,
    r.name as role_name,
    r.display_name as role_display_name,
    p.id as permission_id,
    p.name as permission_name,
    p.resource,
    p.action
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.account_status = 'active';

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_role_audit_log_user_id ON role_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_role_audit_log_target_user_id ON role_audit_log(target_user_id);
CREATE INDEX IF NOT EXISTS idx_role_audit_log_created_at ON role_audit_log(created_at DESC);

-- Add comment to tables
COMMENT ON TABLE role_audit_log IS 'Audit log for tracking role changes';
COMMENT ON VIEW user_role_permissions IS 'View combining users with their roles and permissions';