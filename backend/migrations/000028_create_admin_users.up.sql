-- Create admin_users table to store admin emails
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by INTEGER REFERENCES users(id),
    notes TEXT
);

-- Add initial admin users
INSERT INTO admin_users (email, notes) 
VALUES 
    ('bevzenko.sergey@gmail.com', 'Added in migration 000028'),
    ('voroshilovdo@gmail.com', 'Added in migration 000028')
ON CONFLICT (email) DO NOTHING;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS admin_users_email_idx ON admin_users(email);

-- Grant permissions
COMMENT ON TABLE admin_users IS 'Stores list of admin users by email';