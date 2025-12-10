package helper

import (
	"github.com/gin-gonic/gin"
)

type UserHelper struct {
	// Add user-related services here when needed
}

func NewUserHelper() *UserHelper {
	return &UserHelper{}
}

// GetProfile godoc
// @Tags Users
// @Summary Get user profile
// @Description Get current authenticated user's profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object
// @Failure 401 {object} object
// @Router /users/profile [get]
func (h *UserHelper) GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get user profile endpoint - Coming soon",
		"status":  "not_implemented",
	})
}

// UpdateProfile godoc
// @Tags Users
// @Summary Update user profile
// @Description Update current authenticated user's profile information
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Updated profile information"
// @Success 200 {object} object
// @Failure 401 {object} object
// @Router /users/profile [put]
func (h *UserHelper) UpdateProfile(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Update user profile endpoint - Coming soon",
		"status":  "not_implemented",
	})
}

// GetUsers godoc
// @Tags Users
// @Summary List all users
// @Description Get list of all users (admin only)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object
// @Failure 403 {object} object
// @Router /users [get]
func (h *UserHelper) GetUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get all users endpoint (admin only) - Coming soon",
		"status":  "not_implemented",
	})
}
