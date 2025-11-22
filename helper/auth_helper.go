package helper

import (
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

type AuthHelper struct {
	LoginService    *service.LoginService
	RegisterService *service.RegisterService
}

func NewAuthHelper(loginSvc *service.LoginService, registerSvc *service.RegisterService) *AuthHelper {
	return &AuthHelper{
		LoginService:    loginSvc,
		RegisterService: registerSvc,
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

	token, err := h.LoginService.GenerateToken(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
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
	c.JSON(200, gin.H{
		"message": "Logout successful",
		"status":  "success",
	})
}
