package route

import (
	"prestasi-mahasiswa/helper"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine,
	healthHelper *helper.HealthHelper,
	authHelper *helper.AuthHelper,
	achievementHelper *helper.AchievementHelper,
	userHelper *helper.UserHelper) {

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
		// Setup auth routes
		setupAuthRoutes(v1, authHelper)

		// Setup achievement routes
		setupAchievementRoutes(v1, achievementHelper)

		// Setup user routes
		setupUserRoutes(v1, userHelper)
	}
}

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(rg *gin.RouterGroup, authHelper *helper.AuthHelper) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", authHelper.Login)
		auth.POST("/register", authHelper.Register)
		auth.POST("/logout", authHelper.Logout)
	}
}

// setupAchievementRoutes configures achievement routes
func setupAchievementRoutes(rg *gin.RouterGroup, achievementHelper *helper.AchievementHelper) {
	achievements := rg.Group("/achievements")
	{
		achievements.GET("/", achievementHelper.GetAchievements)
		achievements.POST("/", achievementHelper.CreateAchievement)
		achievements.GET("/:id", achievementHelper.GetAchievement)
		achievements.PUT("/:id", achievementHelper.UpdateAchievement)
		achievements.DELETE("/:id", achievementHelper.DeleteAchievement)

		// File upload for achievements
		achievements.POST("/:id/files", achievementHelper.UploadFile)
		achievements.GET("/:id/files", achievementHelper.GetFiles)
		achievements.DELETE("/:id/files/:fileId", achievementHelper.DeleteFile)
	}
}

// setupUserRoutes configures user routes
func setupUserRoutes(rg *gin.RouterGroup, userHelper *helper.UserHelper) {
	users := rg.Group("/users")
	{
		users.GET("/profile", userHelper.GetProfile)
		users.PUT("/profile", userHelper.UpdateProfile)
		users.GET("/", userHelper.GetUsers) // For admin
	}
}
