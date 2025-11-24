-- Migration: Initial database schema for Sistem Pelaporan Prestasi Mahasiswa
-- Author: Aryo Prabowo (434231027)
-- Date: 2025-11-24

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nim VARCHAR(20) UNIQUE, -- For mahasiswa only
    nip VARCHAR(20) UNIQUE, -- For dosen_wali and admin
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('mahasiswa', 'dosen_wali', 'admin')),
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create achievements table
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mahasiswa_id UUID NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    achievement_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'verified', 'rejected')),
    verified_by UUID, -- Reference to dosen_wali or admin
    verified_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_achievements_mahasiswa FOREIGN KEY (mahasiswa_id) REFERENCES users(id),
    CONSTRAINT fk_achievements_verifier FOREIGN KEY (verified_by) REFERENCES users(id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_nim ON users(nim);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_active_deleted ON users(is_active, is_deleted);

CREATE INDEX IF NOT EXISTS idx_achievements_mahasiswa ON achievements(mahasiswa_id);
CREATE INDEX IF NOT EXISTS idx_achievements_status ON achievements(status);
CREATE INDEX IF NOT EXISTS idx_achievements_category ON achievements(category);
CREATE INDEX IF NOT EXISTS idx_achievements_date ON achievements(achievement_date);
CREATE INDEX IF NOT EXISTS idx_achievements_deleted ON achievements(is_deleted);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_achievements_updated_at 
    BEFORE UPDATE ON achievements 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default admin user (password: admin123)
INSERT INTO users (id, nip, name, email, password, role, is_active, is_deleted) 
VALUES (
    gen_random_uuid(),
    'ADMIN001',
    'System Administrator',
    'admin@prestasi.ac.id',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- bcrypt hash for 'admin123'
    'admin',
    TRUE,
    FALSE
) ON CONFLICT (email) DO NOTHING;

-- Success message
\echo 'Database migration completed successfully!'
\echo 'Default admin user: admin@prestasi.ac.id / admin123'