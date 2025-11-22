package main

import (
	"log"
	"prestasi-mahasiswa/app"
	"prestasi-mahasiswa/config"
	"prestasi-mahasiswa/database"
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
