package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Get list of migrations
	migrationFiles, err := getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run migrations silently
	for _, file := range migrationFiles {
		runMigration(db, migrationsPath, file)
	}

	return nil
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
	// Read migration file
	filePath := filepath.Join(migrationsPath, filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Execute migration silently
	_, err = db.Exec(string(content))
	return err
}
