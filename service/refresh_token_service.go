package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"prestasi-mahasiswa/utils"
)

type RefreshTokenService struct {
	db      *sql.DB
	jwtUtil *utils.JWTUtil
}

type RefreshToken struct {
	ID         int        `json:"id"`
	UserID     string     `json:"user_id"`
	TokenHash  string     `json:"token_hash"`
	ExpiresAt  time.Time  `json:"expires_at"`
	IsRevoked  bool       `json:"is_revoked"`
	RevokedAt  *time.Time `json:"revoked_at"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	IPAddress  *string    `json:"ip_address"`
	UserAgent  *string    `json:"user_agent"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func NewRefreshTokenService(db *sql.DB, jwtUtil *utils.JWTUtil) *RefreshTokenService {
	return &RefreshTokenService{
		db:      db,
		jwtUtil: jwtUtil,
	}
}

// StoreRefreshToken stores refresh token in database
func (rts *RefreshTokenService) StoreRefreshToken(userID, tokenHash string, expiresAt time.Time, ipAddress, userAgent *string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := rts.db.Exec(query, userID, tokenHash, expiresAt, ipAddress, userAgent)
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

// ValidateRefreshToken validates and retrieves refresh token from database
func (rts *RefreshTokenService) ValidateRefreshToken(tokenString string) (*RefreshToken, error) {
	// First validate JWT structure and signature
	claims, err := rts.jwtUtil.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate hash to look up in database
	tokenHash := utils.HashToken(tokenString)

	// Check if token exists and is not revoked
	var refreshToken RefreshToken
	query := `
		SELECT id, user_id, token_hash, expires_at, is_revoked, revoked_at, 
		       created_at, last_used_at, ip_address, user_agent
		FROM refresh_tokens
		WHERE token_hash = $1 AND user_id = $2
	`

	err = rts.db.QueryRow(query, tokenHash, claims.UserID).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.IsRevoked,
		&refreshToken.RevokedAt,
		&refreshToken.CreatedAt,
		&refreshToken.LastUsedAt,
		&refreshToken.IPAddress,
		&refreshToken.UserAgent,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("refresh token not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if token is revoked
	if refreshToken.IsRevoked {
		return nil, errors.New("refresh token has been revoked")
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	return &refreshToken, nil
}

// RefreshTokens validates refresh token and returns new token pair
func (rts *RefreshTokenService) RefreshTokens(tokenString string, ipAddress, userAgent *string) (*TokenRefreshResponse, error) {
	// Validate refresh token
	refreshToken, err := rts.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Get user details from claims
	claims, err := rts.jwtUtil.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Update last used timestamp
	err = rts.UpdateLastUsed(refreshToken.ID, ipAddress, userAgent)
	if err != nil {
		// Log error but don't fail the refresh
		fmt.Printf("Warning: failed to update last used timestamp: %v\n", err)
	}

	// Generate new token pair
	tokenPair, err := rts.jwtUtil.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Store new refresh token
	newTokenHash := utils.HashToken(tokenPair.RefreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	err = rts.StoreRefreshToken(claims.UserID, newTokenHash, expiresAt, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	// Revoke old refresh token for security
	err = rts.RevokeRefreshToken(refreshToken.ID)
	if err != nil {
		// Log error but don't fail the refresh
		fmt.Printf("Warning: failed to revoke old refresh token: %v\n", err)
	}

	return &TokenRefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// RevokeRefreshToken marks refresh token as revoked
func (rts *RefreshTokenService) RevokeRefreshToken(tokenID int) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true, revoked_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := rts.db.Exec(query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// RevokeUserRefreshTokens revokes all refresh tokens for a user (logout all devices)
func (rts *RefreshTokenService) RevokeUserRefreshTokens(userID string) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true, revoked_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND is_revoked = false
	`
	result, err := rts.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user refresh tokens: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Revoked %d refresh tokens for user %s\n", rowsAffected, userID)

	return nil
}

// UpdateLastUsed updates the last used timestamp and metadata
func (rts *RefreshTokenService) UpdateLastUsed(tokenID int, ipAddress, userAgent *string) error {
	query := `
		UPDATE refresh_tokens
		SET last_used_at = CURRENT_TIMESTAMP, ip_address = $2, user_agent = $3
		WHERE id = $1
	`
	_, err := rts.db.Exec(query, tokenID, ipAddress, userAgent)
	if err != nil {
		return fmt.Errorf("failed to update last used timestamp: %w", err)
	}

	return nil
}

// GetUserRefreshTokens gets all active refresh tokens for a user
func (rts *RefreshTokenService) GetUserRefreshTokens(userID string) ([]RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, is_revoked, revoked_at,
		       created_at, last_used_at, ip_address, user_agent
		FROM refresh_tokens
		WHERE user_id = $1 AND is_revoked = false
		ORDER BY created_at DESC
	`

	rows, err := rts.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user refresh tokens: %w", err)
	}
	defer rows.Close()

	var tokens []RefreshToken
	for rows.Next() {
		var token RefreshToken
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.TokenHash,
			&token.ExpiresAt,
			&token.IsRevoked,
			&token.RevokedAt,
			&token.CreatedAt,
			&token.LastUsedAt,
			&token.IPAddress,
			&token.UserAgent,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan refresh token: %w", err)
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// CleanupExpiredTokens removes expired refresh tokens from database
func (rts *RefreshTokenService) CleanupExpiredTokens() error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < CURRENT_TIMESTAMP
	`
	result, err := rts.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	rowsDeleted, _ := result.RowsAffected()
	fmt.Printf("Cleaned up %d expired refresh tokens\n", rowsDeleted)

	return nil
}
