package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Config holds the database configuration
type Config struct {
	Host                 string
	Port                 int
	Database             string
	User                 string
	Password             string
	SSLMode              string
	MaxConnections       int32
	MaxIdleConnections   int32
	MaxConnLifetime      time.Duration
	MaxConnIdleTime      time.Duration
	HealthCheckPeriod    time.Duration
	MinConns             int32
}

// DefaultConfig returns a default database configuration
func DefaultConfig() *Config {
	return &Config{
		Host:                 "localhost",
		Port:                 5432,
		Database:             "ga_ticketing",
		User:                 "ga_user",
		Password:             "ga_password",
		SSLMode:              "disable",
		MaxConnections:       50,
		MaxIdleConnections:   10,
		MaxConnLifetime:      3600 * time.Second,
		MaxConnIdleTime:      1800 * time.Second,
		HealthCheckPeriod:    60 * time.Second,
		MinConns:             5,
	}
}

// Pool wraps pgxpool.Pool with additional functionality
type Pool struct {
	*pgxpool.Pool
	logger *zap.Logger
}

// NewConnection creates a new database connection pool
func NewConnection(ctx context.Context, config *Config, logger *zap.Logger) (*Pool, error) {
	if config == nil {
		config = DefaultConfig()
	}

	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s",
		config.Host,
		config.Port,
		config.Database,
		config.User,
		config.Password,
		config.SSLMode,
		config.MaxConnections,
		config.MinConns,
		config.MaxConnLifetime,
		config.MaxConnIdleTime,
	)

	logger.Info("Connecting to database",
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.String("database", config.Database),
		zap.String("user", config.User),
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Error("Failed to parse database connection config", zap.Error(err))
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = config.MaxConnections
	poolConfig.MinConns = config.MinConns
	poolConfig.HealthCheckPeriod = config.HealthCheckPeriod
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	// Configure connection logger - tracer disabled for now
	// poolConfig.ConnConfig.Tracer = &Tracer{logger: logger}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error("Failed to create database connection pool", zap.Error(err))
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection pool created successfully",
		zap.Int32("max_connections", poolConfig.MaxConns),
		zap.Int32("min_connections", poolConfig.MinConns),
	)

	return &Pool{
		Pool:   pool,
		logger: logger,
	}, nil
}

// Close closes the database connection pool
func (p *Pool) Close() {
	p.logger.Info("Closing database connection pool")
	p.Pool.Close()
}

// Stats returns connection pool statistics
func (p *Pool) Stats() *pgxpool.Stat {
	return p.Pool.Stat()
}

// Health checks the database connection health
func (p *Pool) Health(ctx context.Context) error {
	return p.Ping(ctx)
}


// Tracer implements pgx.ConnTracer for detailed logging
type Tracer struct {
	logger *zap.Logger
}

func (t *Tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.logger.Debug("Query started",
		zap.String("sql", data.SQL),
		zap.Strings("args", formatArgs(data.Args)),
	)
	return ctx
}

func (t *Tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		t.logger.Error("Query failed",
			zap.Error(data.Err),
		)
	} else {
		t.logger.Debug("Query completed",
			zap.String("command", data.CommandTag.String()),
		)
	}
}

// Helper function to format query arguments for logging
func formatArgs(args []interface{}) []string {
	strArgs := make([]string, len(args))
	for i, arg := range args {
		strArgs[i] = fmt.Sprintf("%v", arg)
	}
	return strArgs
}