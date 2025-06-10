package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	Security SecurityConfig
	Metrics  MetricsConfig
}

type ServerConfig struct {
	Port        int           `envconfig:"SERVER_PORT" default:"8080"`
	Mode        string        `envconfig:"SERVER_MODE" default:"development"`
	GRPCPort    int           `envconfig:"GRPC_PORT" default:"9090"`
	MetricsPort int           `envconfig:"METRICS_PORT" default:"9091"`
	Timeout     time.Duration `envconfig:"SERVER_TIMEOUT" default:"30s"`
}

type DatabaseConfig struct {
	Host               string        `envconfig:"DB_HOST" default:"localhost"`
	Port               int           `envconfig:"DB_PORT" default:"5432"`
	Name               string        `envconfig:"DB_NAME" default:"oms"`
	User               string        `envconfig:"DB_USER" default:"oms_user"`
	Password           string        `envconfig:"DB_PASSWORD" required:"true"`
	SSLMode            string        `envconfig:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns       int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns       int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime    time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"5m"`
	ConnMaxIdleTime    time.Duration `envconfig:"DB_CONN_MAX_IDLE_TIME" default:"30m"`
	MigrationDirectory string        `envconfig:"DB_MIGRATION_DIR" default:"./migrations"`
}

type RedisConfig struct {
	Host     string        `envconfig:"REDIS_HOST" default:"localhost"`
	Port     int           `envconfig:"REDIS_PORT" default:"6379"`
	Password string        `envconfig:"REDIS_PASSWORD"`
	DB       int           `envconfig:"REDIS_DB" default:"0"`
	TTL      time.Duration `envconfig:"REDIS_TTL" default:"5m"`
}

type KafkaConfig struct {
	Brokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	Topic   string   `envconfig:"KAFKA_TOPIC" default:"oms-events"`
	GroupID string   `envconfig:"KAFKA_GROUP_ID" default:"oms-service"`
}

type SecurityConfig struct {
	JWTSecret      string `envconfig:"JWT_SECRET" required:"true"`
	APIKeyHeader   string `envconfig:"API_KEY_HEADER" default:"X-API-Key"`
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS" default:"*"`
	TLSEnabled     bool   `envconfig:"TLS_ENABLED" default:"false"`
}

type MetricsConfig struct {
	Path          string `envconfig:"METRICS_PATH" default:"/metrics"`
	TraceEndpoint string `envconfig:"TRACE_ENDPOINT" default:"http://jaeger:14268/api/traces"`
	Enabled       bool   `envconfig:"METRICS_ENABLED" default:"true"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	var cfg Config

	// Load from environment
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.Password == "" && os.Getenv("DB_PASSWORD") == "" {
		return fmt.Errorf("database password is required")
	}

	if c.Security.JWTSecret == "" && os.Getenv("JWT_SECRET") == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}