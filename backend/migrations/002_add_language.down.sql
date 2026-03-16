-- Remove language column from users table
DROP INDEX IF EXISTS idx_users_language;
ALTER TABLE users DROP COLUMN IF EXISTS language;
