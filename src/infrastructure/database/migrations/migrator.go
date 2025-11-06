package migrations

import (
	"context"
	"database/sql/driver"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrator handles database migrations
type Migrator struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(pool *pgxpool.Pool, logger *zap.Logger) *Migrator {
	return &Migrator{
		pool:   pool,
		logger: logger,
	}
}

// Up runs all available migrations
func (m *Migrator) Up(ctx context.Context) error {
	m.logger.Info("Starting database migrations")

	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Info("No migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := migrator.Version()
	if err != nil {
		m.logger.Warn("Could not get migration version", zap.Error(err))
	} else {
		m.logger.Info("Migrations applied successfully",
			zap.Uint32("version", uint32(version)),
			zap.Bool("dirty", dirty),
		)
	}

	return nil
}

// Down rolls back all migrations
func (m *Migrator) Down(ctx context.Context) error {
	m.logger.Info("Rolling back database migrations")

	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Down(); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Info("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	m.logger.Info("Migrations rolled back successfully")
	return nil
}

// Steps applies n migrations forward
func (m *Migrator) Steps(ctx context.Context, n int) error {
	m.logger.Info("Applying migrations steps", zap.Int("steps", n))

	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Steps(n); err != nil {
		return fmt.Errorf("failed to apply migration steps: %w", err)
	}

	m.logger.Info("Migration steps applied successfully")
	return nil
}

// Version returns the current migration version
func (m *Migrator) Version(ctx context.Context) (uint32, bool, error) {
	migrator, err := m.createMigrator()
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return uint32(version), dirty, nil
}

// Force sets the migration version (use with caution)
func (m *Migrator) Force(ctx context.Context, version int) error {
	m.logger.Warn("Forcing migration version", zap.Int("version", version))

	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	m.logger.Info("Migration version forced successfully", zap.Int("version", version))
	return nil
}

// createMigrator creates a new migrate instance
func (m *Migrator) createMigrator() (*migrate.Migrate, error) {
	// Get connection string from pool config
	connString := m.pool.Config().ConnString()

	// Get migration files path
	migrationsPath, err := m.getMigrationsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get migrations path: %w", err)
	}

	// Create migrate instance with database connection string
	migrator, err := migrate.New(
		"file://"+migrationsPath,
		connString,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return migrator, nil
}

// getMigrationsPath returns the path to migration files
func (m *Migrator) getMigrationsPath() (string, error) {
	// Look for migrations in standard locations
	paths := []string{
		"migrations",
		"./migrations",
		"../migrations",
		"../../migrations",
	}

	for _, path := range paths {
		if m.pathExists(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("no migrations directory found")
}

// pathExists checks if a path exists
func (m *Migrator) pathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// sqlDB implements the database/sql driver interface for pgxpool
type sqlDB struct {
	pool *pgxpool.Pool
}

func (db *sqlDB) Begin() (driver.Tx, error) {
	// This is a simplified implementation
	// In production, you might need a more sophisticated approach
	return nil, fmt.Errorf("transactions not supported in this implementation")
}

func (db *sqlDB) Close() error {
	return nil
}

// CreateMigration creates a new migration file with timestamp
func CreateMigration(name string) error {
	timestamp := time.Now().Format("20060102150405")
	upFile := fmt.Sprintf("%s_%s.up.sql", timestamp, name)
	downFile := fmt.Sprintf("%s_%s.down.sql", timestamp, name)

	// Create up migration file
	upContent := fmt.Sprintf(`-- Migration: %s
-- Created: %s
-- Description: %s

-- Add your UP migration SQL here

`, name, time.Now().Format(time.RFC3339), name)

	if err := writeFile("migrations/"+upFile, upContent); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration file
	downContent := fmt.Sprintf(`-- Migration: %s
-- Created: %s
-- Description: %s

-- Add your DOWN migration SQL here

`, name, time.Now().Format(time.RFC3339), name)

	if err := writeFile("migrations/"+downFile, downContent); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
	return nil
}

// GetMigrationFiles returns a list of migration files sorted by version
func GetMigrationFiles() ([]string, error) {
	migrationsDir := "migrations"

	// Read migration directory
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		files = append(files, filepath.Join(migrationsDir, entry.Name()))
	}

	// Sort files by version (timestamp)
	sort.Slice(files, func(i, j int) bool {
		versionI := extractVersion(files[i])
		versionJ := extractVersion(files[j])
		return versionI < versionJ
	})

	return files, nil
}

// extractVersion extracts version number from migration filename
func extractVersion(filename string) int {
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	if len(parts) == 0 {
		return 0
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return version
}

// writeFile writes content to a file
func writeFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}