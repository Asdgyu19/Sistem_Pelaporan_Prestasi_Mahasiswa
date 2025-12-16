-- Add advisor_id column to users table if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS advisor_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- Add is_deleted and deleted_at to achievements table if they don't exist
ALTER TABLE achievements ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT false;

-- Add index for advisor_id
CREATE INDEX IF NOT EXISTS idx_users_advisor_id ON users(advisor_id);
CREATE INDEX IF NOT EXISTS idx_achievements_is_deleted ON achievements(is_deleted);
