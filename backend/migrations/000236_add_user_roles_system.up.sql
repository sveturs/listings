-- Create roles table
CREATE TABLE IF NOT EXISTS public.roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS public.permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS public.role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

-- Create user_roles junction table
CREATE TABLE IF NOT EXISTS public.user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    assigned_by INTEGER REFERENCES users(id),
    PRIMARY KEY (user_id, role_id)
);

-- Insert default roles
INSERT INTO roles (name, display_name, description) VALUES
    ('super_admin', 'Super Administrator', 'Full system access with all permissions'),
    ('admin', 'Administrator', 'Administrative access to manage users and content'),
    ('moderator', 'Moderator', 'Can moderate content and manage listings'),
    ('vendor', 'Vendor', 'Can create and manage storefronts and products'),
    ('user', 'User', 'Standard user with basic permissions')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (name, resource, action, description) VALUES
    -- User management
    ('users.view', 'users', 'view', 'View user profiles'),
    ('users.list', 'users', 'list', 'List all users'),
    ('users.edit', 'users', 'edit', 'Edit user profiles'),
    ('users.delete', 'users', 'delete', 'Delete users'),
    ('users.ban', 'users', 'ban', 'Ban/suspend users'),
    
    -- Admin access
    ('admin.access', 'admin', 'access', 'Access admin panel'),
    ('admin.users', 'admin', 'users', 'Manage users in admin panel'),
    ('admin.categories', 'admin', 'categories', 'Manage categories'),
    ('admin.attributes', 'admin', 'attributes', 'Manage attributes'),
    ('admin.translations', 'admin', 'translations', 'Manage translations'),
    
    -- Listings management
    ('listings.create', 'listings', 'create', 'Create listings'),
    ('listings.edit_own', 'listings', 'edit_own', 'Edit own listings'),
    ('listings.edit_any', 'listings', 'edit_any', 'Edit any listing'),
    ('listings.delete_own', 'listings', 'delete_own', 'Delete own listings'),
    ('listings.delete_any', 'listings', 'delete_any', 'Delete any listing'),
    ('listings.moderate', 'listings', 'moderate', 'Moderate listings'),
    
    -- Storefronts management
    ('storefronts.create', 'storefronts', 'create', 'Create storefronts'),
    ('storefronts.edit_own', 'storefronts', 'edit_own', 'Edit own storefronts'),
    ('storefronts.edit_any', 'storefronts', 'edit_any', 'Edit any storefront'),
    ('storefronts.delete_own', 'storefronts', 'delete_own', 'Delete own storefronts'),
    ('storefronts.delete_any', 'storefronts', 'delete_any', 'Delete any storefront'),
    
    -- Reviews management
    ('reviews.create', 'reviews', 'create', 'Create reviews'),
    ('reviews.edit_own', 'reviews', 'edit_own', 'Edit own reviews'),
    ('reviews.delete_own', 'reviews', 'delete_own', 'Delete own reviews'),
    ('reviews.delete_any', 'reviews', 'delete_any', 'Delete any review'),
    ('reviews.moderate', 'reviews', 'moderate', 'Moderate reviews'),
    
    -- Chat management
    ('chat.send', 'chat', 'send', 'Send messages'),
    ('chat.view_any', 'chat', 'view_any', 'View any chat'),
    ('chat.moderate', 'chat', 'moderate', 'Moderate chats')
ON CONFLICT (name) DO NOTHING;

-- Assign permissions to roles
-- Super Admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
ON CONFLICT DO NOTHING;

-- Admin gets most permissions except super admin specific
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin' 
    AND p.name NOT IN ('users.delete')
ON CONFLICT DO NOTHING;

-- Moderator permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'moderator' 
    AND p.name IN (
        'admin.access', 
        'listings.moderate', 
        'reviews.moderate', 
        'chat.moderate',
        'listings.edit_any',
        'reviews.delete_any'
    )
ON CONFLICT DO NOTHING;

-- Vendor permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'vendor' 
    AND p.name IN (
        'listings.create',
        'listings.edit_own',
        'listings.delete_own',
        'storefronts.create',
        'storefronts.edit_own',
        'storefronts.delete_own',
        'reviews.create',
        'reviews.edit_own',
        'chat.send'
    )
ON CONFLICT DO NOTHING;

-- User permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'user' 
    AND p.name IN (
        'listings.create',
        'listings.edit_own',
        'listings.delete_own',
        'reviews.create',
        'reviews.edit_own',
        'reviews.delete_own',
        'chat.send',
        'users.view'
    )
ON CONFLICT DO NOTHING;

-- Add role column to users table if not exists
ALTER TABLE users ADD COLUMN IF NOT EXISTS role_id INTEGER REFERENCES roles(id);

-- Set default role for users (after the role_id column is added)
UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user') WHERE role_id IS NULL;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

-- Update existing admin users to have admin role
UPDATE users u
SET role_id = (SELECT id FROM roles WHERE name = 'admin')
FROM admin_users au
WHERE u.email = au.email;

-- Create function to check user permission
CREATE OR REPLACE FUNCTION check_user_permission(
    p_user_id INTEGER,
    p_permission_name VARCHAR
) RETURNS BOOLEAN AS $$
DECLARE
    has_permission BOOLEAN;
BEGIN
    -- Check if user has permission through their role
    SELECT EXISTS (
        SELECT 1
        FROM users u
        JOIN roles r ON u.role_id = r.id
        JOIN role_permissions rp ON r.id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = p_user_id 
        AND p.name = p_permission_name
    ) INTO has_permission;
    
    -- Also check user_roles table for multiple roles
    IF NOT has_permission THEN
        SELECT EXISTS (
            SELECT 1
            FROM user_roles ur
            JOIN role_permissions rp ON ur.role_id = rp.role_id
            JOIN permissions p ON rp.permission_id = p.id
            WHERE ur.user_id = p_user_id 
            AND p.name = p_permission_name
        ) INTO has_permission;
    END IF;
    
    RETURN has_permission;
END;
$$ LANGUAGE plpgsql;

-- Note: update_updated_at trigger would be created if the function exists