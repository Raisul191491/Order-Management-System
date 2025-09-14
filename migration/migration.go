package migrations

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	sqlDir           = "./sql"
	migrationTimeout = 30 * time.Minute
)

// Migration represents a single database migration
type Migration struct {
	Version string
	Up      func(context.Context, *gorm.DB) error
}

// MigrationRecord tracks applied migrations in the database
type MigrationRecord struct {
	Version   string    `gorm:"primaryKey;size:255"`
	AppliedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the table name for MigrationRecord
func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db     *gorm.DB
	sqlDir string
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:     db,
		sqlDir: sqlDir,
	}
}

// SetSQLDir allows overriding the default SQL directory
func (m *MigrationManager) SetSQLDir(dir string) {
	m.sqlDir = dir
}

// createMigrationFunc creates a migration function for a given version
func (m *MigrationManager) createMigrationFunc(version string) func(context.Context, *gorm.DB) error {
	return func(ctx context.Context, db *gorm.DB) error {
		return m.applyMigration(ctx, version, db)
	}
}

// applyMigration applies a single migration by reading and executing SQL file
func (m *MigrationManager) applyMigration(ctx context.Context, version string, db *gorm.DB) error {
	statements, err := m.readSQLFile(ctx, version)
	if err != nil {
		return fmt.Errorf("failed to read migration %s: %w", version, err)
	}

	for i, stmt := range statements {
		if err := m.execSQL(ctx, db, stmt); err != nil {
			return fmt.Errorf("failed to execute statement %d in migration %s: %w", i+1, version, err)
		}
	}

	return nil
}

// readSQLFile reads and parses SQL file into individual statements
func (m *MigrationManager) readSQLFile(ctx context.Context, version string) ([]string, error) {
	filename := fmt.Sprintf("%s.sql", version)
	filePath := filepath.Join(m.sqlDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("migration file %s does not exist", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file %s: %w", filePath, err)
	}

	// Parse SQL statements
	statements := m.parseSQLStatements(string(content))
	if len(statements) == 0 {
		return nil, fmt.Errorf("no valid SQL statements found in %s", filename)
	}

	return statements, nil
}

// parseSQLStatements splits SQL content into individual statements
func (m *MigrationManager) parseSQLStatements(content string) []string {
	// Split by semicolon, but be smarter about it
	rawStatements := strings.Split(content, ";")
	var statements []string

	for _, stmt := range rawStatements {
		// Clean up the statement
		cleaned := strings.TrimSpace(stmt)

		// Skip empty statements and comments
		if cleaned == "" || strings.HasPrefix(cleaned, "--") {
			continue
		}

		// Remove single-line comments from the statement
		lines := strings.Split(cleaned, "\n")
		var cleanLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "--") {
				cleanLines = append(cleanLines, line)
			}
		}

		if len(cleanLines) > 0 {
			statements = append(statements, strings.Join(cleanLines, "\n"))
		}
	}

	return statements
}

// execSQL executes a SQL statement with proper error handling
func (m *MigrationManager) execSQL(ctx context.Context, db *gorm.DB, sql string) error {
	// Create a timeout context for the SQL execution
	execCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	result := db.WithContext(execCtx).Exec(sql)
	if result.Error != nil {
		var pqErr *pq.Error
		if errors.As(result.Error, &pqErr) {
			// Handle specific PostgreSQL errors
			switch pqErr.Code {
			case "42P07": // relation already exists
				log.Printf("Relation already exists, skipping: %s", sql[:min(50, len(sql))])
				return nil
			case "42701": // duplicate column
				log.Printf("Column already exists, skipping: %s", sql[:min(50, len(sql))])
				return nil
			case "42P16": // index already exists
				log.Printf("Index already exists, skipping: %s", sql[:min(50, len(sql))])
				return nil
			default:
				return fmt.Errorf("postgres error %s: %s", pqErr.Code, pqErr.Message)
			}
		}
		return fmt.Errorf("SQL execution failed: %w", result.Error)
	}

	log.Printf("Successfully executed SQL statement")
	return nil
}

// getAppliedMigrations retrieves all applied migrations from database
func (m *MigrationManager) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	var appliedMigrations []MigrationRecord

	result := m.db.WithContext(ctx).Find(&appliedMigrations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch applied migrations: %w", result.Error)
	}

	appliedVersions := make(map[string]bool, len(appliedMigrations))
	for _, migration := range appliedMigrations {
		appliedVersions[migration.Version] = true
	}

	return appliedVersions, nil
}

// recordMigration saves a migration record to the database
func (m *MigrationManager) recordMigration(ctx context.Context, tx *gorm.DB, version string) error {
	record := MigrationRecord{
		Version:   version,
		AppliedAt: time.Now(),
	}

	if err := tx.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("failed to record migration %s: %w", version, err)
	}

	return nil
}

// Migrate runs all pending migrations
func (m *MigrationManager) Migrate(ctx context.Context) error {
	// Set up timeout for the entire migration process
	migrationCtx, cancel := context.WithTimeout(ctx, migrationTimeout)
	defer cancel()

	// Ensure migrations table exists
	if err := m.db.WithContext(migrationCtx).AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedVersions, err := m.getAppliedMigrations(migrationCtx)
	if err != nil {
		return err
	}

	// Sort migrations by version
	migrations := make([]Migration, len(AllMigrations))
	copy(migrations, AllMigrations)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Apply pending migrations
	for _, migration := range migrations {
		if appliedVersions[migration.Version] {
			log.Printf("Skipping already applied migration: %s", migration.Version)
			continue
		}

		log.Printf("Applying migration: %s", migration.Version)

		// Execute migration in transaction
		err := m.db.WithContext(migrationCtx).Transaction(func(tx *gorm.DB) error {
			// Apply the migration
			if err := migration.Up(migrationCtx, tx); err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}

			// Record successful migration
			if err := m.recordMigration(migrationCtx, tx, migration.Version); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		log.Printf("Successfully applied migration: %s", migration.Version)
	}

	log.Printf("All migrations completed successfully")
	return nil
}

// Migrate is a convenience function that uses the default migration manager
func Migrate(ctx context.Context, db *gorm.DB) error {
	manager := NewMigrationManager(db)
	return manager.Migrate(ctx)
}

// AllMigrations holds all migrations in chronological order
var AllMigrations = []Migration{
	{
		Version: "0001_initial_database_setup",
		Up: func(ctx context.Context, db *gorm.DB) error {
			manager := NewMigrationManager(db)
			return manager.applyMigration(ctx, "0001_initial_database_setup", db)
		},
	},
}
