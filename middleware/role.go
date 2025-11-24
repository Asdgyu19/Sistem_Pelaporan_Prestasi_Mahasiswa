package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole creates middleware that checks user role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "Please login first",
			})
			c.Abort()
			return
		}

		role := userRole.(string)

		// Check if user role is in allowed roles
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "Insufficient permissions",
				"message":       "Your role (" + role + ") is not authorized for this action",
				"allowed_roles": allowedRoles,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Role constants for easier use
const (
	RoleMahasiswa = "mahasiswa"
	RoleDosenWali = "dosen_wali"
	RoleAdmin     = "admin"
)

// Helper functions for common role combinations
func RequireMahasiswa() gin.HandlerFunc {
	return RequireRole(RoleMahasiswa)
}
func RequireDosenWali() gin.HandlerFunc {
	return RequireRole(RoleDosenWali)
}
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(RoleAdmin)
}
func RequireDosenOrAdmin() gin.HandlerFunc {
	return RequireRole(RoleDosenWali, RoleAdmin)
}
func RequireAnyAuthenticated() gin.HandlerFunc {
	return RequireRole(RoleMahasiswa, RoleDosenWali, RoleAdmin)
}
