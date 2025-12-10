package app

import (
	"database/sql"
	"prestasi-mahasiswa/config"
	"prestasi-mahasiswa/database"
	"prestasi-mahasiswa/helper"
	"prestasi-mahasiswa/middleware"
	"prestasi-mahasiswa/route"
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// General middleware
	a.Router.Use(middleware.CORS())
	a.Router.Use(middleware.RequestLogger())
	a.Router.Use(middleware.Recovery())

	// Setup Swagger documentation
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (a *App) initRoutes() {
	// Initialize services with JWT configuration
	loginService := service.NewLoginService(a.DB, a.Config.JWT.Secret, a.Config.JWT.ExpireHours)
	registerService := service.NewRegisterService(a.DB)
	achievementService := service.NewAchievementService(a.DB, a.MongoDB)
	fileService := service.NewFileService(a.MongoDB)
	userService := service.NewUserService(a.DB)
	refreshTokenService := service.NewRefreshTokenService(a.DB, loginService.JWTUtil)
	reportService := service.NewReportService(a.DB)

	// Initialize helpers
	healthHelper := helper.NewHealthHelper(a.DB, a.MongoDB)
	authHelper := helper.NewAuthHelper(loginService, registerService, refreshTokenService)
	achievementHelper := helper.NewAchievementHelper(achievementService, fileService)
	userHelper := helper.NewUserHelper()
	adminUserHelper := helper.NewAdminUserHelper(userService)
	studentHelper := helper.NewStudentHelper(userService, achievementService)
	lecturerHelper := helper.NewLecturerHelper(userService)
	reportHelper := helper.NewReportHelper(reportService)

	// Setup all routes using separate route files with JWT secret
	route.SetupRoutes(a.Router, a.Config.JWT.Secret, healthHelper, authHelper, achievementHelper, userHelper, adminUserHelper, studentHelper, lecturerHelper, reportHelper)
}

func (a *App) Run() error {
	address := ":" + a.Config.Server.Port
	return a.Router.Run(address)
}
