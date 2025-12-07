package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Create migrations table if it doesn't exist
	err := createMigrationsTable(db)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migrations
	migrationFiles, err := getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run migrations
	for _, file := range migrationFiles {
		err := runMigration(db, migrationsPath, file)
		if err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	log.Println("✅ All database migrations completed successfully!")
	return nil
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

func runMigration(db *sql.DB, migrationsPath, filename string) error {
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
	filepath := filepath.Join(migrationsPath, filename)
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
