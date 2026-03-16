-- Add language column to users table
ALTER TABLE users ADD COLUMN language VARCHAR(10) DEFAULT 'en';

-- Create index for language queries
CREATE INDEX IF NOT EXISTS idx_users_language ON users(language);
