-- Migration: Add advisor_id and profile_data columns to users table
-- This fixes the missing columns error

-- Add advisor_id column (for mahasiswa -> dosen_wali relationship)
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS advisor_id UUID;

-- Add profile_data column (for additional user information)
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS profile_data JSONB;

-- Add foreign key constraint for advisor relationship
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_users_advisor'
    ) THEN
        ALTER TABLE users 
        ADD CONSTRAINT fk_users_advisor 
        FOREIGN KEY (advisor_id) REFERENCES users(id);
    END IF;
END $$;

-- Add index for advisor_id for better performance
CREATE INDEX IF NOT EXISTS idx_users_advisor ON users(advisor_id);

-- Success message
\echo 'Database schema updated successfully!'
\echo 'Added advisor_id and profile_data columns to users table'