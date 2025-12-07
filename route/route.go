package route

import (
	"prestasi-mahasiswa/helper"
	"prestasi-mahasiswa/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine,
	jwtSecret string,
	healthHelper *helper.HealthHelper,
	authHelper *helper.AuthHelper,
	achievementHelper *helper.AchievementHelper,
	userHelper *helper.UserHelper,
	adminUserHelper *helper.AdminUserHelper) {

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Sistem Pelaporan Prestasi Mahasiswa API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Health check route
	router.GET("/health", healthHelper.CheckHealth)

	// API v1 routes group
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		setupAuthRoutes(v1, authHelper)

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Setup achievement routes (role-based access)
			setupAchievementRoutes(protected, achievementHelper)

			// Setup user routes (role-based access)
			setupUserRoutes(protected, userHelper)

			// Setup admin user management routes
			setupAdminUserRoutes(protected, adminUserHelper)
		}
	}
}

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(rg *gin.RouterGroup, authHelper *helper.AuthHelper) {
	auth := rg.Group("/auth")
	{
		// Public authentication routes
		auth.POST("/login", authHelper.Login)
		auth.POST("/register", authHelper.Register)
		auth.POST("/refresh", authHelper.RefreshToken)      // Refresh access token using refresh token
		auth.POST("/revoke", authHelper.RevokeRefreshToken) // Revoke specific refresh token

		// Protected authentication routes (require valid access token)
		auth.POST("/logout", authHelper.Logout)               // Logout current session and revoke all tokens
		auth.POST("/logout-all", authHelper.LogoutAllDevices) // Logout from all devices
		auth.GET("/tokens", authHelper.GetActiveTokens)       // Get list of active refresh tokens
	}
}

// setupAchievementRoutes configures achievement routes with role-based access
func setupAchievementRoutes(rg *gin.RouterGroup, achievementHelper *helper.AchievementHelper) {
	achievements := rg.Group("/achievements")
	{
		// All authenticated users can view achievements (with role-based filtering in helper)
		achievements.GET("/", middleware.RequireAnyAuthenticated(), achievementHelper.GetAchievements)
		achievements.GET("/:id", middleware.RequireAnyAuthenticated(), achievementHelper.GetAchievement)
		achievements.GET("/:id/files", middleware.RequireAnyAuthenticated(), achievementHelper.GetFiles)

		// Only mahasiswa can create and manage their own achievements
		achievements.POST("/", middleware.RequireMahasiswa(), achievementHelper.CreateAchievement)
		achievements.PUT("/:id", middleware.RequireMahasiswa(), achievementHelper.UpdateAchievement)
		achievements.DELETE("/:id", middleware.RequireMahasiswa(), achievementHelper.DeleteAchievement)

		// Achievement workflow - mahasiswa submits, dosen/admin verify/reject
		achievements.POST("/:id/submit", middleware.RequireMahasiswa(), achievementHelper.SubmitAchievement)
		achievements.POST("/:id/verify", middleware.RequireDosenOrAdmin(), achievementHelper.VerifyAchievement)
		achievements.POST("/:id/reject", middleware.RequireDosenOrAdmin(), achievementHelper.RejectAchievement)

		// File management - mahasiswa can upload, all can view/download
		achievements.POST("/:id/files", middleware.RequireMahasiswa(), achievementHelper.UploadFile)
		achievements.GET("/:id/files/:fileId/download", middleware.RequireAnyAuthenticated(), achievementHelper.DownloadFile)
		achievements.DELETE("/:id/files/:fileId", middleware.RequireMahasiswa(), achievementHelper.DeleteFile)
	}
}

// setupUserRoutes configures user routes with role-based access
func setupUserRoutes(rg *gin.RouterGroup, userHelper *helper.UserHelper) {
	users := rg.Group("/users")
	{
		// All authenticated users can view and update their own profile
		users.GET("/profile", middleware.RequireAnyAuthenticated(), userHelper.GetProfile)
		users.PUT("/profile", middleware.RequireAnyAuthenticated(), userHelper.UpdateProfile)

		// Only admin can view all users
		users.GET("/", middleware.RequireAdmin(), userHelper.GetUsers)
	}
}

// setupAdminUserRoutes configures admin user management routes
func setupAdminUserRoutes(rg *gin.RouterGroup, adminUserHelper *helper.AdminUserHelper) {
	admin := rg.Group("/admin")
	admin.Use(middleware.RequireAdmin()) // All admin routes require admin role
	{
		// User management routes
		userMgmt := admin.Group("/users")
		{
			// Basic user CRUD
			userMgmt.GET("/", adminUserHelper.GetAllUsers)      // GET /api/v1/admin/users?role=mahasiswa&search=aryo
			userMgmt.GET("/:id", adminUserHelper.GetUserByID)   // GET /api/v1/admin/users/{id}
			userMgmt.POST("/", adminUserHelper.CreateUser)      // POST /api/v1/admin/users
			userMgmt.PUT("/:id", adminUserHelper.UpdateUser)    // PUT /api/v1/admin/users/{id}
			userMgmt.DELETE("/:id", adminUserHelper.DeleteUser) // DELETE /api/v1/admin/users/{id}

			// Role management
			userMgmt.PUT("/:id/role", adminUserHelper.ChangeUserRole)     // PUT /api/v1/admin/users/{id}/role
			userMgmt.PUT("/:id/status", adminUserHelper.ToggleUserStatus) // PUT /api/v1/admin/users/{id}/status

			// Advisor management
			userMgmt.POST("/assign-advisor", adminUserHelper.AssignAdvisor)                      // POST /api/v1/admin/users/assign-advisor
			userMgmt.DELETE("/:id/advisor", adminUserHelper.RemoveAdvisor)                       // DELETE /api/v1/admin/users/{id}/advisor
			userMgmt.GET("/advisors", adminUserHelper.GetAvailableAdvisors)                      // GET /api/v1/admin/users/advisors
			userMgmt.GET("/advisors/:advisor_id/students", adminUserHelper.GetStudentsByAdvisor) // GET /api/v1/admin/users/advisors/{id}/students
		}
	}
}
