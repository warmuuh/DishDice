-- Add role column
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';
ALTER TABLE users ADD CONSTRAINT check_role CHECK (role IN ('user', 'admin'));

-- Add status column
ALTER TABLE users ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending';
ALTER TABLE users ADD CONSTRAINT check_status CHECK (status IN ('pending', 'approved', 'rejected'));

-- Create indexes for admin queries
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Grandfather existing users to approved status
UPDATE users SET status = 'approved';
