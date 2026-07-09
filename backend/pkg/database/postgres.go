package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"

	"thinh/gin-app/config"
)

func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// ensureMigrationTable creates the schema_migrations tracking table if it doesn't exist.
func ensureMigrationTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS schema_migrations (
		filename VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	return err
}

// isMigrationApplied checks if a migration has already been applied.
func isMigrationApplied(db *sql.DB, filename string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE filename = $1", filename).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", filename, err)
	}
	return count > 0, nil
}

// markMigrationApplied records a migration as successfully applied.
func markMigrationApplied(db *sql.DB, filename string) error {
	_, err := db.Exec("INSERT INTO schema_migrations (filename) VALUES ($1) ON CONFLICT DO NOTHING", filename)
	return err
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Ensure the tracking table exists
	if err := ensureMigrationTable(db); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		// Skip if already applied
		applied, err := isMigrationApplied(db, f)
		if err != nil {
			return err
		}
		if applied {
			log.Printf("Migration already applied, skipping: %s", f)
			continue
		}

		path := filepath.Join(migrationsDir, f)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}

		// Execute the entire migration file as a single batch.
		// This supports multi-statement blocks (DO $$ ... END $$; etc.) that
		// would break if split by semicolons.
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("execute migration %s: %w\nSQL: %s", f, err, string(content))
		}

		// Mark as applied
		if err := markMigrationApplied(db, f); err != nil {
			return fmt.Errorf("mark migration %s as applied: %w", f, err)
		}

		log.Printf("Migration applied: %s", f)
	}

	return nil
}
