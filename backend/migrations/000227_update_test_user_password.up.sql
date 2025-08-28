-- Update test user password to 'password123' (bcrypt hash)
-- Hash generated with cost 10 for password 'password123'
UPDATE users 
SET password = '$2a$10$YKqH9hEn.qZ8xVfL9TAcxOaRz3gGssM6DLvKz8h4rRrfqVXy7jJDe'
WHERE email = 'test@example.com';