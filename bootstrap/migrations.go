package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// Migrate runs SQL-based migrations located in ./migrations/sql using golang-migrate.
func Migrate(db *gorm.DB) error {
	// Attempt to locate project root to build absolute path to migrations
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get working dir: %w", err)
	}

	migrationsPath := filepath.ToSlash(filepath.Join(cwd, "migrations", "sql"))
	sourceURL := "file://" + migrationsPath

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Apply all up migrations (no-op if already up)
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("migration up failed: %w", err)
	}

	fmt.Println("SQL migrations applied successfully")
	return nil
}
