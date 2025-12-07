package helper

import (
	"prestasi-mahasiswa/service"
	"prestasi-mahasiswa/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHelper struct {
	LoginService        *service.LoginService
	RegisterService     *service.RegisterService
	RefreshTokenService *service.RefreshTokenService
}

func NewAuthHelper(loginSvc *service.LoginService, registerSvc *service.RegisterService, refreshTokenSvc *service.RefreshTokenService) *AuthHelper {
	return &AuthHelper{
		LoginService:        loginSvc,
		RegisterService:     registerSvc,
		RefreshTokenService: refreshTokenSvc,
	}
}

func (h *AuthHelper) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.LoginService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// Generate token pair (access + refresh)
	tokenPair, err := h.LoginService.JWTUtil.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Store refresh token in database
	refreshTokenHash := utils.HashToken(tokenPair.RefreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	// Get client IP and User-Agent for security tracking
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err = h.RefreshTokenService.StoreRefreshToken(user.ID, refreshTokenHash, expiresAt, &clientIP, &userAgent)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"message":       "Login successful",
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"token_type":    tokenPair.TokenType,
		"expires_in":    tokenPair.ExpiresIn,
		"user":          user,
	})
}

func (h *AuthHelper) Register(c *gin.Context) {
	var req struct {
		NIM      string `json:"nim"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.RegisterService.ValidateUserData(req.NIM, req.Name, req.Email, req.Password, req.Role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.RegisterService.CreateUser(req.NIM, req.Name, req.Email, req.Password, req.Role); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "User registered successfully",
		"status":  "success",
	})
}

func (h *AuthHelper) Logout(c *gin.Context) {
	// Extract user ID from token (if available)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
		if err == nil {
			claims, err := h.LoginService.JWTUtil.ValidateAccessToken(tokenString)
			if err == nil {
				// Revoke all refresh tokens for this user (logout all devices)
				err = h.RefreshTokenService.RevokeUserRefreshTokens(claims.UserID)
				if err != nil {
					// Log but don't fail logout
					c.Header("X-Warning", "Failed to revoke refresh tokens")
				}
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "Logout successful",
		"status":  "success",
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHelper) RefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// Get client metadata for security tracking
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Refresh tokens
	tokenResponse, err := h.RefreshTokenService.RefreshTokens(req.RefreshToken, &clientIP, &userAgent)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tokenResponse)
}

// RevokeRefreshToken handles individual refresh token revocation
func (h *AuthHelper) RevokeRefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate and get refresh token details
	refreshToken, err := h.RefreshTokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// Revoke the token
	err = h.RefreshTokenService.RevokeRefreshToken(refreshToken.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to revoke refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Refresh token revoked successfully",
		"status":  "success",
	})
}

// LogoutAllDevices revokes all refresh tokens for the current user
func (h *AuthHelper) LogoutAllDevices(c *gin.Context) {
	// Get user ID from access token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.LoginService.JWTUtil.ValidateAccessToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid access token"})
		return
	}

	// Revoke all refresh tokens for this user
	err = h.RefreshTokenService.RevokeUserRefreshTokens(claims.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to logout from all devices"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Logged out from all devices successfully",
		"status":  "success",
	})
}

// GetActiveTokens returns list of active refresh tokens for the current user
func (h *AuthHelper) GetActiveTokens(c *gin.Context) {
	// Get user ID from access token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.LoginService.JWTUtil.ValidateAccessToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid access token"})
		return
	}

	// Get active refresh tokens
	tokens, err := h.RefreshTokenService.GetUserRefreshTokens(claims.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve active tokens"})
		return
	}

	// Remove sensitive token hash from response
	type TokenInfo struct {
		ID         int        `json:"id"`
		CreatedAt  time.Time  `json:"created_at"`
		LastUsedAt *time.Time `json:"last_used_at"`
		ExpiresAt  time.Time  `json:"expires_at"`
		IPAddress  *string    `json:"ip_address"`
		UserAgent  *string    `json:"user_agent"`
	}

	var tokenInfos []TokenInfo
	for _, token := range tokens {
		tokenInfos = append(tokenInfos, TokenInfo{
			ID:         token.ID,
			CreatedAt:  token.CreatedAt,
			LastUsedAt: token.LastUsedAt,
			ExpiresAt:  token.ExpiresAt,
			IPAddress:  token.IPAddress,
			UserAgent:  token.UserAgent,
		})
	}

	c.JSON(200, gin.H{
		"active_tokens": tokenInfos,
		"total":         len(tokenInfos),
	})
}
