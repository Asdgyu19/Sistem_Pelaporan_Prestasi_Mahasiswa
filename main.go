// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0.0
// @description Backend API untuk Student Achievement Reporting System dengan JWT authentication dan role-based access control
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @basePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Bearer <token>"

package main

import (
	"log"
	"prestasi-mahasiswa/app"
	"prestasi-mahasiswa/config"
	"prestasi-mahasiswa/database"

	_ "prestasi-mahasiswa/docs" // Swagger docs
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize databases
	log.Println("🔌 Connecting to PostgreSQL...")
	db, err := database.InitPostgreSQL(cfg)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer db.Close()
	log.Println("✅ PostgreSQL connected successfully!")

	// Database migrations (run silently if enabled)
	// err = database.RunMigrations(db, "./database/migrations")

	log.Println("🔌 Connecting to MongoDB...")
	mongodb, err := database.InitMongoDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongodb.Disconnect()
	log.Println("✅ MongoDB connected successfully!")

	// Initialize application
	application := app.NewApp(cfg, db, mongodb)

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	if err := application.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
