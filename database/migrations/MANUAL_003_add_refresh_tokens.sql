-- Execute this SQL manually in your PostgreSQL database to add refresh token support

-- Create refresh_tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP,
    ip_address VARCHAR(45), -- Support both IPv4 and IPv6
    user_agent TEXT,
    
    -- Indexes for performance
    CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_is_revoked ON refresh_tokens(is_revoked);

-- Optional: Create a cleanup function to remove expired tokens
CREATE OR REPLACE FUNCTION cleanup_expired_refresh_tokens()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM refresh_tokens
    WHERE expires_at < NOW();
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- You can call this function periodically to cleanup expired tokens:
-- SELECT cleanup_expired_refresh_tokens();

COMMENT ON TABLE refresh_tokens IS 'Stores refresh tokens for JWT authentication system';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA256 hash of the actual refresh token';
COMMENT ON COLUMN refresh_tokens.ip_address IS 'Client IP address where token was created/last used';
COMMENT ON COLUMN refresh_tokens.user_agent IS 'Client user agent string';







