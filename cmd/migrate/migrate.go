package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"prestasi-mahasiswa/config"
	"prestasi-mahasiswa/database"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := database.InitPostgreSQL(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	err = createMigrationsTable(db)
	if err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}

	// Get list of migrations
	migrationFiles, err := getMigrationFiles("./database/migrations")
	if err != nil {
		log.Fatal("Failed to get migration files:", err)
	}

	// Run migrations
	for _, file := range migrationFiles {
		err := runMigration(db, file)
		if err != nil {
			log.Fatalf("Failed to run migration %s: %v", file, err)
		}
	}

	log.Println("✅ All migrations completed successfully!")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}

	// Sort migrations to ensure they run in order
	sort.Strings(migrations)
	return migrations, nil
}

func runMigration(db *sql.DB, filename string) error {
	// Check if migration has already been run
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = $1", filename).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("⏭️  Skipping migration %s (already executed)", filename)
		return nil
	}

	// Read migration file
	filepath := filepath.Join("./database/migrations", filename)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Execute migration
	log.Printf("🔄 Running migration: %s", filename)
	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	// Mark migration as executed
	_, err = db.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)
	if err != nil {
		return err
	}

	log.Printf("✅ Migration %s completed successfully", filename)
	return nil
}
