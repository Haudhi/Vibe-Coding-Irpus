package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds all application configuration
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	OAuth2     OAuth2Config     `mapstructure:"oauth2"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	CORS       CORSConfig       `mapstructure:"cors"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	SMTP       SMTPConfig       `mapstructure:"smtp"`
	Security   SecurityConfig   `mapstructure:"security"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	App        AppConfig        `mapstructure:"app"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host                 string        `mapstructure:"host"`
	Port                 int           `mapstructure:"port"`
	Database             string        `mapstructure:"database"`
	User                 string        `mapstructure:"user"`
	Password             string        `mapstructure:"password"`
	SSLMode              string        `mapstructure:"ssl_mode"`
	MaxConnections       int32         `mapstructure:"max_connections"`
	MaxIdleConnections   int32         `mapstructure:"max_idle_connections"`
	MaxConnLifetime      time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime      time.Duration `mapstructure:"max_conn_idle_time"`
	HealthCheckPeriod    time.Duration `mapstructure:"health_check_period"`
	MinConns             int32         `mapstructure:"min_conns"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	Expiry        time.Duration `mapstructure:"expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
	Issuer        string        `mapstructure:"issuer"`
}

// OAuth2Config holds OAuth2 configuration
type OAuth2Config struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	AuthURL      string `mapstructure:"auth_url"`
	TokenURL     string `mapstructure:"token_url"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FileName string `mapstructure:"file_name"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	MaxAge         int      `mapstructure:"max_age"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	Burst             int `mapstructure:"burst"`
}

// SMTPConfig holds SMTP configuration for emails
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	BCRYPTCost   int    `mapstructure:"bcrypt_cost"`
	SessionSecret string `mapstructure:"session_secret"`
	CSRFSecret   string `mapstructure:"csrf_secret"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	EnableMetrics bool   `mapstructure:"enable_metrics"`
	MetricsPort   int    `mapstructure:"metrics_port"`
	EnableProfiling bool `mapstructure:"enable_profiling"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig(logger *zap.Logger) (*Config, error) {
	config := &Config{}

	// Set default values
	setDefaults()

	// Load .env file using godotenv
	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file", zap.Error(err))
	} else {
		logger.Info("Successfully loaded .env file")
	}

	// Setup viper for environment variable mapping
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	logger.Info("Configuration loaded successfully",
		zap.String("app_name", config.App.Name),
		zap.String("app_version", config.App.Version),
		zap.String("app_env", config.App.Env),
	)

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("server.shutdown_timeout", "10s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "ga_ticketing")
	viper.SetDefault("database.user", "ga_user")
	viper.SetDefault("database.password", "ga_password")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_connections", 50)
	viper.SetDefault("database.max_idle_connections", 10)
	viper.SetDefault("database.max_conn_lifetime", "3600s")
	viper.SetDefault("database.max_conn_idle_time", "1800s")
	viper.SetDefault("database.health_check_period", "60s")
	viper.SetDefault("database.min_conns", 5)

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-this-in-production")
	viper.SetDefault("jwt.expiry", "24h")
	viper.SetDefault("jwt.refresh_expiry", "168h")
	viper.SetDefault("jwt.issuer", "ga-ticketing-system")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Content-Type", "Authorization"})

	// Rate limiting defaults
	viper.SetDefault("rate_limit.requests_per_minute", 100)
	viper.SetDefault("rate_limit.burst", 200)

	// Security defaults
	viper.SetDefault("security.bcrypt_cost", 12)

	// Monitoring defaults
	viper.SetDefault("monitoring.enable_metrics", true)
	viper.SetDefault("monitoring.metrics_port", 9090)
	viper.SetDefault("monitoring.enable_profiling", false)

	// App defaults
	viper.SetDefault("app.name", "GA Ticketing System")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.env", "development")
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server port
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate database configuration
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}
	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}

	// Validate JWT secret
	if config.JWT.Secret == "" || config.JWT.Secret == "your-super-secret-jwt-key-change-this-in-production" {
		if config.App.Env == "production" {
			return fmt.Errorf("JWT secret must be set in production")
		}
	}

	// Validate Redis configuration
	if config.Redis.Port <= 0 || config.Redis.Port > 65535 {
		return fmt.Errorf("invalid redis port: %d", config.Redis.Port)
	}

	// Validate security configurations
	if config.Security.BCRYPTCost < 4 || config.Security.BCRYPTCost > 31 {
		return fmt.Errorf("invalid bcrypt cost: %d (must be between 4 and 31)", config.Security.BCRYPTCost)
	}

	return nil
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.SSLMode,
	)
}

// GetRedisAddress returns the Redis connection address
func (c *Config) GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// IsDevelopment returns true if the app is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the app is running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// IsTest returns true if the app is running in test mode
func (c *Config) IsTest() bool {
	return c.App.Env == "test"
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt gets an environment variable as an integer with a default value
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBool gets an environment variable as a boolean with a default value
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}