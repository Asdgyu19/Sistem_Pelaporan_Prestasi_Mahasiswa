package helper

import (
	"net/http"
	"strconv"

	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

type AdminUserHelper struct {
	UserService *service.UserService
}

func NewAdminUserHelper(userService *service.UserService) *AdminUserHelper {
	return &AdminUserHelper{
		UserService: userService,
	}
}

// Get all users with filters
func (h *AdminUserHelper) GetAllUsers(c *gin.Context) {
	role := c.Query("role")
	search := c.Query("search")
	includeDeleted, _ := strconv.ParseBool(c.DefaultQuery("include_deleted", "false"))

	users, err := h.UserService.GetAllUsers(role, search, includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get users",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Users retrieved successfully",
		"data": gin.H{
			"users": users,
			"total": len(users),
		},
	})
}

// Get user by ID
func (h *AdminUserHelper) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// Create new user
func (h *AdminUserHelper) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	user, err := h.UserService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

// Update user
func (h *AdminUserHelper) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	user, err := h.UserService.UpdateUser(userID, req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to update user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// Delete user (soft delete)
func (h *AdminUserHelper) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.UserService.DeleteUser(userID)
	if err != nil {
		if err.Error() == "user not found or already deleted" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found or already deleted",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// Change user role
func (h *AdminUserHelper) ChangeUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Update user role using service
	updateReq := service.UpdateUserRequest{
		Role: &req.Role,
	}

	user, err := h.UserService.UpdateUser(userID, updateReq)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to change user role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User role changed successfully",
		"data":    user,
	})
}

// Assign advisor to mahasiswa
func (h *AdminUserHelper) AssignAdvisor(c *gin.Context) {
	var req service.AssignAdvisorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	err := h.UserService.AssignAdvisor(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to assign advisor",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Advisor assigned successfully",
	})
}

// Remove advisor from mahasiswa
func (h *AdminUserHelper) RemoveAdvisor(c *gin.Context) {
	mahasiswaID := c.Param("id")

	updateReq := service.UpdateUserRequest{
		AdvisorID: nil, // Set to null
	}

	user, err := h.UserService.UpdateUser(mahasiswaID, updateReq)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Mahasiswa not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to remove advisor",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Advisor removed successfully",
		"data":    user,
	})
}

// Get students by advisor
func (h *AdminUserHelper) GetStudentsByAdvisor(c *gin.Context) {
	advisorID := c.Param("advisor_id")

	students, err := h.UserService.GetStudentsByAdvisor(advisorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get students",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Students retrieved successfully",
		"data": gin.H{
			"students": students,
			"total":    len(students),
		},
	})
}

// Get available advisors
func (h *AdminUserHelper) GetAvailableAdvisors(c *gin.Context) {
	advisors, err := h.UserService.GetAvailableAdvisors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get advisors",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Advisors retrieved successfully",
		"data": gin.H{
			"advisors": advisors,
			"total":    len(advisors),
		},
	})
}

// Toggle user active status
func (h *AdminUserHelper) ToggleUserStatus(c *gin.Context) {
	userID := c.Param("id")

	// Get current user to check status
	currentUser, err := h.UserService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get user",
			"error":   err.Error(),
		})
		return
	}

	// Toggle status
	newStatus := !currentUser.IsActive
	updateReq := service.UpdateUserRequest{
		IsActive: &newStatus,
	}

	user, err := h.UserService.UpdateUser(userID, updateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to toggle user status",
			"error":   err.Error(),
		})
		return
	}

	action := "activated"
	if !newStatus {
		action = "deactivated"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User " + action + " successfully",
		"data":    user,
	})
}
