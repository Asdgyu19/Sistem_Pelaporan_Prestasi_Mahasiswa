package route

import (
	"database/sql"
	"prestasi-mahasiswa/database"
	"prestasi-mahasiswa/helper"
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *sql.DB, mongodb *database.MongoDB) {
	// Initialize services
	loginService := service.NewLoginService(db)
	registerService := service.NewRegisterService(db)
	achievementService := service.NewAchievementService(db, mongodb)
	fileService := service.NewFileService(mongodb)

	// Initialize helpers
	healthHelper := helper.NewHealthHelper(db, mongodb)
	authHelper := helper.NewAuthHelper(loginService, registerService)
	achievementHelper := helper.NewAchievementHelper(achievementService, fileService)
	userHelper := helper.NewUserHelper()

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
