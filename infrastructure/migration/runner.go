package migration

import (
	"fmt"
	"os"
	"strconv"
)

// Constants representing environment variable keys for migration configuration.
const (
	envKeyDsn          = "MIGRATION_DSN"
	envKeyOp           = "MIGRATION_OPERATION"
	envKeyForceVersion = "MIGRATION_FORCE_VERSION"
	envKeyFilesDir     = "MIGRATION_FILES_DIR"
)

func readForceVersion() (int, error) {
	forceVersionRaw, ok := os.LookupEnv(envKeyForceVersion)
	if !ok {
		return 0, nil
	}

	forceVersion, err := strconv.Atoi(forceVersionRaw)
	if err != nil {
		return 0, fmt.Errorf("convert forceVersion: %w", err)
	}

	return forceVersion, nil
}

// RunMigrationsFromEnv reads migration configuration from environment variables,
// creates a MigrationRunner, and runs the specified migration operation.
func RunMigrationsFromEnv(logger logger) error {
	dsn, ok := os.LookupEnv(envKeyDsn)
	if !ok {
		return fmt.Errorf("missing env: %s", envKeyDsn)
	}

	operation, ok := os.LookupEnv(envKeyOp)
	if !ok {
		return fmt.Errorf("missing env: %s", envKeyOp)
	}

	forceVersion, err := readForceVersion()
	if err != nil {
		return fmt.Errorf("read forceVersion: %w", err)
	}

	opts := []Option{WithLogger(logger)}
	if filesDir, exists := os.LookupEnv(envKeyFilesDir); exists {
		opts = append(opts, WithFilesDir(filesDir))
	}

	runner, err := NewRunner(dsn, opts...)
	if err != nil {
		return fmt.Errorf("new migrations runner: %w", err)
	}

	// Get the current migration version and log it.
	version, dirty, err := runner.Version()
	if err != nil {
		logger.Error(fmt.Sprintf("getting current migration version: %v", err))
	} else {
		logger.Info(fmt.Sprintf("migration version before operation: %d, dirty: %v", version, dirty))
	}

	if err = runner.Run(OperationData{
		ID:           operation,
		ForceVersion: forceVersion,
	}); err != nil {
		return fmt.Errorf("run operation %s: %w", operation, err)
	}

	logger.Info("successfully finished migration")

	// Get the migration version after the operation and log it.
	version, dirty, err = runner.Version()
	if err != nil {
		return fmt.Errorf("getting migration version after operation: %w", err)
	}

	logger.Info(fmt.Sprintf("migration version after operation: %d, dirty: %v", version, dirty))

	return nil
}