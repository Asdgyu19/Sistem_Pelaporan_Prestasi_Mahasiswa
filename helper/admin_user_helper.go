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

// GetAllUsers godoc
// @Tags Users (Admin)
// @Summary List all users
// @Description Get list of all users with optional filters (admin only)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role query string false "Filter by role (mahasiswa, dosen_wali, admin)"
// @Param search query string false "Search by name or email"
// @Param include_deleted query bool false "Include deleted users"
// @Success 200 {object} object
// @Failure 500 {object} object
// @Router /admin/users [get]
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

// GetUserByID godoc
// @Tags Users (Admin)
// @Summary Get user details
// @Description Get detailed information for a specific user (admin only)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /admin/users/{id} [get]
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

// CreateUser godoc
// @Tags Users (Admin)
// @Summary Create new user
// @Description Create a new user account (admin only)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "User creation details"
// @Success 201 {object} object
// @Failure 400 {object} object
// @Router /admin/users [post]
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

// UpdateUser godoc
// @Tags Users (Admin)
// @Summary Update user information
// @Description Update user details (admin only)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body object true "Updated user information"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /admin/users/{id} [put]
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

// DeleteUser godoc
// @Tags Users (Admin)
// @Summary Delete user
// @Description Soft delete a user account (admin only)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /admin/users/{id} [delete]
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

// ChangeUserRole godoc
// @Tags Users (Admin)
// @Summary Change user role
// @Description Change a user's role (admin only)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body object true "New role"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /admin/users/{id}/role [put]
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
