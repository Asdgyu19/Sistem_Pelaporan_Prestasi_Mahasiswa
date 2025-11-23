package app

import (
	"database/sql"
	"prestasi-mahasiswa/config"
	"prestasi-mahasiswa/database"
	"prestasi-mahasiswa/helper"
	"prestasi-mahasiswa/route"
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config  *config.Config
	DB      *sql.DB
	MongoDB *database.MongoDB
	Router  *gin.Engine

	// Repositories
	userRepo        UserRepository
	achievementRepo AchievementRepository
	fileRepo        AchievementFileRepository

	// Usecases
	authUsecase        AuthUsecase
	achievementUsecase AchievementUsecase
	fileUsecase        FileUsecase
}

func NewApp(cfg *config.Config, db *sql.DB, mongodb *database.MongoDB) *App {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	app := &App{
		Config:  cfg,
		DB:      db,
		MongoDB: mongodb,
		Router:  gin.Default(),
	}

	// Initialize repositories
	app.initRepositories()

	// Initialize usecases
	app.initUsecases()

	// Initialize middleware
	app.initMiddleware()

	// Initialize routes
	app.initRoutes()

	return app
}

func (a *App) initRepositories() {
	// TODO: Initialize repository implementations
	// These will be implemented in separate files
}

func (a *App) initUsecases() {
	// TODO: Initialize usecase implementations
	// These will be implemented in separate files
}

func (a *App) initMiddleware() {
	// CORS Middleware
	a.Router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Request Logger
	a.Router.Use(gin.Logger())
	a.Router.Use(gin.Recovery())
}

func (a *App) initRoutes() {
	// Initialize services
	loginService := service.NewLoginService(a.DB)
	registerService := service.NewRegisterService(a.DB)
	achievementService := service.NewAchievementService(a.DB, a.MongoDB)
	fileService := service.NewFileService(a.MongoDB)

	// Initialize helpers
	healthHelper := helper.NewHealthHelper(a.DB, a.MongoDB)
	authHelper := helper.NewAuthHelper(loginService, registerService)
	achievementHelper := helper.NewAchievementHelper(achievementService, fileService)
	userHelper := helper.NewUserHelper()

	// Setup all routes using separate route files
	route.SetupRoutes(a.Router, healthHelper, authHelper, achievementHelper, userHelper)
}

func (a *App) Run() error {
	address := ":" + a.Config.Server.Port
	return a.Router.Run(address)
}
