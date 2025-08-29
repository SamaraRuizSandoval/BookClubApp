package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents all application configurations
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// ServerConfig represents HTTP server configurations
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Environment  string
}

// DatabaseConfig represents database configurations
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig represents Redis configurations
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig represents JWT configurations
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// Load loads configurations from environment
func Load() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Redis:    loadRedisConfig(),
		JWT:      loadJWTConfig(),
	}
}

// loadServerConfig loads server configurations
func loadServerConfig() ServerConfig {
	port := getEnv("PORT", ":8080")
	readTimeout := getEnvAsDuration("READ_TIMEOUT", 15*time.Second)
	writeTimeout := getEnvAsDuration("WRITE_TIMEOUT", 15*time.Second)
	idleTimeout := getEnvAsDuration("IDLE_TIMEOUT", 60*time.Second)
	environment := getEnv("ENVIRONMENT", "development")

	return ServerConfig{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Environment:  environment,
	}
}

// loadDatabaseConfig loads database configurations
func loadDatabaseConfig() DatabaseConfig {
	host := getEnv("DB_HOST", "localhost")
	port := getEnvAsInt("DB_PORT", 5432)
	user := getEnv("DB_USER", "bookclub")
	password := getEnv("DB_PASSWORD", "bookclub123")
	dbName := getEnv("DB_NAME", "bookclubapp")
	sslMode := getEnv("DB_SSLMODE", "disable")

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}

// loadRedisConfig loads Redis configurations
func loadRedisConfig() RedisConfig {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnvAsInt("REDIS_PORT", 6379)
	password := getEnv("REDIS_PASSWORD", "")
	db := getEnvAsInt("REDIS_DB", 0)

	return RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}

// loadJWTConfig loads JWT configurations
func loadJWTConfig() JWTConfig {
	secret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	expiration := getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour)

	return JWTConfig{
		Secret:     secret,
		Expiration: expiration,
	}
}

// getEnv gets environment variable or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as int or default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets environment variable as Duration or default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// IsDevelopment returns true if in development environment
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if in production environment
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + strconv.Itoa(c.Database.Port) +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port)
}
