package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

// Supported migrate operations.
const (
	defaultFilesDir = "dbmigrations"
	operationUp     = "up"
	operationDown   = "down"
	operationForce  = "force"
)

type operationFn func(*Runner, OperationData) error

var supportedOperations = map[string]operationFn{
	operationUp:    runUp,
	operationDown:  runDown,
	operationForce: runForce,
}

// OperationData contains information about the migration operation to be performed.
type OperationData struct {
	ID           string
	ForceVersion int
}

// Runner is responsible for managing and running database migrations.
type Runner struct {
	mgr      *migrate.Migrate
	filesDir string
	logger   logger
}

// Option represents a function that configures a Runner.
type Option func(runner *Runner)

// WithLogger sets a custom logger for the Runner.
// If not provided, a noopLogger will be used by default.
func WithLogger(logger logger) Option {
	return func(runner *Runner) {
		runner.logger = logger
	}
}

// WithFilesDir sets a custom directory containing migration files for the Runner.
// If not provided, the default directory "dbmigrations" will be used.
func WithFilesDir(filesDir string) Option {
	return func(runner *Runner) {
		runner.filesDir = filesDir
	}
}

// NewRunner creates a new Runner with the given database connection string (dsn) and options.
func NewRunner(dsn string, opts ...Option) (*Runner, error) {
	runner := &Runner{
		filesDir: defaultFilesDir,
		logger:   &noopLogger{},
	}

	for _, opt := range opts {
		opt(runner)
	}

	mgr, err := migrate.New("file://"+runner.filesDir, dsn)
	if err != nil {
		return nil, fmt.Errorf("creating Migrate object: %w", err)
	}

	mgr.Log = toMigrationsLogger(runner.logger)
	runner.mgr = mgr

	return runner, nil
}

// Run executes the migration operation specified by the OperationData.
func (m *Runner) Run(operation OperationData) error {
	operationName := operation.ID
	operationFunc, found := supportedOperations[operationName]
	if !found {
		return fmt.Errorf("unsupported migration operation: %s", operationName)
	}

	if err := operationFunc(m, operation); err != nil {
		return fmt.Errorf("operation %s failed: %w", operationName, err)
	}

	return nil
}

// Version returns the current migration version, a dirty flag, and an error if any.
func (m *Runner) Version() (uint, bool, error) {
	return m.mgr.Version()
}

// runUp runs the "up" migration operation, applying new migrations to the database.
func runUp(m *Runner, _ OperationData) error {
	m.logger.Info("running migrate UP")

	err := m.mgr.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info("no new migrations found in: ", m.filesDir)
		return nil
	}
	if err != nil {
		return fmt.Errorf("running migrations UP failed: %w", err)
	}
	return nil
}

// runDown runs the "down" migration operation, rolling back the latest applied migration.
func runDown(m *Runner, _ OperationData) error {
	m.logger.Info(fmt.Sprintf("running migrate DOWN with STEPS=%d", 1))

	// always rollback the latest applied migration only
	err := m.mgr.Steps(-1)
	if err != nil {
		return fmt.Errorf("running migrations DOWN failed: %w", err)
	}
	return nil
}

// runForce runs the "force" migration operation,
// forcibly setting the migration version without running the actual migrations.
func runForce(m *Runner, op OperationData) error {
	m.logger.Info(fmt.Sprintf("running FORCE with VERSION %d", op.ForceVersion))

	err := m.mgr.Force(op.ForceVersion)
	if err != nil {
		return fmt.Errorf("running migrations FORCE with VERSION %d failed: %w", op.ForceVersion, err)
	}
	return nil
}

type logger interface {
	Info(args ...any)
	Error(args ...any)
	Printf(format string, v ...any)
}

type noopLogger struct{}

func (l *noopLogger) Info(...any)           {}
func (l *noopLogger) Error(...any)          {}
func (l *noopLogger) Printf(string, ...any) {}
func (l *noopLogger) Verbose() bool         { return false }

// Adapter to use logger like logrus.
func toMigrationsLogger(logger logger) *migrationsLogger {
	return &migrationsLogger{logger: logger}
}

// To be able to log not only errors, but also Info level logs from golang-migrate,
// we have to implement migrate.Logger interface.
type migrationsLogger struct {
	logger logger
}

// Printf is like fmt.Printf.
func (m *migrationsLogger) Printf(format string, v ...any) {
	m.logger.Printf(format, v...)
}

// Verbose should return true when verbose logging output is wanted.
func (m *migrationsLogger) Verbose() bool {
	return true
}

func (m *migrationsLogger) Info(args ...any) {
	m.logger.Info(args...)
}

func (m *migrationsLogger) Error(args ...any) {
	m.logger.Error(args...)
}