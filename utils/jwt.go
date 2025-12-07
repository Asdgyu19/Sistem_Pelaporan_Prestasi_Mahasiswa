package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secretKey         []byte
	expireHours       int
	refreshExpireDays int
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

func NewJWTUtil(secret string, expireHours int) *JWTUtil {
	return &JWTUtil{
		secretKey:         []byte(secret),
		expireHours:       expireHours,
		refreshExpireDays: 7, // 7 days for refresh token
	}
}

// GenerateToken creates access token (backward compatibility)
func (j *JWTUtil) GenerateToken(userID, email, role string) (string, error) {
	return j.GenerateAccessToken(userID, email, role)
}

// GenerateAccessToken creates short-lived access token (15 minutes)
func (j *JWTUtil) GenerateAccessToken(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Short-lived for security

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "prestasi-mahasiswa-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken creates long-lived refresh token (7 days)
func (j *JWTUtil) GenerateRefreshToken(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.refreshExpireDays) * 24 * time.Hour)

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "prestasi-mahasiswa-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateTokenPair creates both access and refresh tokens
func (j *JWTUtil) GenerateTokenPair(userID, email, role string) (*TokenPair, error) {
	accessToken, err := j.GenerateAccessToken(userID, email, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.GenerateRefreshToken(userID, email, role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	}, nil
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateAccessToken validates access token specifically
func (j *JWTUtil) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type: expected access token")
	}

	return claims, nil
}

// ValidateRefreshToken validates refresh token specifically
func (j *JWTUtil) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type: expected refresh token")
	}

	return claims, nil
}

// GenerateRandomToken creates secure random token for refresh token hash
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken creates SHA256 hash of token for storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ExtractTokenFromHeader extracts Bearer token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[7:], nil
}
