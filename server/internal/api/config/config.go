package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment       Environment
	LogLevel          slog.Level
	CORSClientBaseURL string
	Address           string
	Database          DBConfig
	Resend            ResendConfig
	Secrets           SecretConfig
	StreamService     StreamServiceConfig
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
}

type SecretConfig struct {
	JWTSecret string
}

type DBConfig struct {
	Driver       string
	URI          string
	MigrationDir string
}

type ResendConfig struct {
	Key                    string
	DevelopmentSendToEmail string
}

type StreamServiceConfig struct {
	HeartbeatTickerSeconds int64
}

func Load() (*Config, error) {
	if err := loadFirstEnvFile(".env", "/etc/secrets/.env"); err != nil {
		return nil, err
	}

	environment, err := NewEnvironment(getenvDefault("ENVIRONMENT", string(EnvironmentProduction)))
	if err != nil {
		return nil, err
	}

	accessTokenExpiration := getenvDefault("ACCESS_TOKEN_EXPIRATION", "15")
	refreshTokenExpiration := getenvDefault("REFRESH_TOKEN_EXPIRATION", "10080")

	accessTokenExpirationMinutes, err := strconv.ParseInt(accessTokenExpiration, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_EXPIRATION: %s", accessTokenExpiration)
	}

	refreshTokenExpirationMinutes, err := strconv.ParseInt(refreshTokenExpiration, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRATION: %s", refreshTokenExpiration)
	}

	hbIntervalSeconds := mustGetenv("HEARTBEAT_INTERVAL_SECONDS")
	heartbeatIntervalSeconds, err := strconv.ParseInt(hbIntervalSeconds, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid HEARTBEAT_INTERVAL_SECONDS: %s", hbIntervalSeconds)
	}

	return &Config{
		Environment:       environment,
		LogLevel:          getSlogLevel(),
		Address:           mustGetenv("API_ADDRESS"),
		CORSClientBaseURL: mustGetenv("CLIENT_BASE_URL"),
		Database: DBConfig{
			Driver:       mustGetenv("DB_DRIVER"),
			URI:          mustGetenv("DB_URI"),
			MigrationDir: mustGetenv("GOOSE_MIGRATION_DIR"),
		},
		Secrets: SecretConfig{
			JWTSecret: mustGetenv("JWT_SECRET"),
		},
		StreamService: StreamServiceConfig{
			HeartbeatTickerSeconds: heartbeatIntervalSeconds,
		},
		Resend: ResendConfig{
			Key:                    mustGetenv("RESEND_KEY"),
			DevelopmentSendToEmail: os.Getenv("RESEND_DEVELOPMENT_SEND_TO_EMAIL"),
		},
		AccessTokenTTL:  time.Duration(accessTokenExpirationMinutes) * time.Minute,
		RefreshTokenTTL: time.Duration(refreshTokenExpirationMinutes) * time.Minute,
	}, nil
}

func (c *Config) SlogLevel() slog.Level {
	if c.Environment.IsDevelopment() {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func (c *Config) GetAuthOptions() AuthOptions {
	return AuthOptions{
		JWTSecret:       c.Secrets.JWTSecret,
		AccessTokenTTL:  c.AccessTokenTTL,
		RefreshTokenTTL: c.RefreshTokenTTL,
	}
}

func getSlogLevel() slog.Level {
	switch strings.ToLower(getenvDefault("LOG_LEVEL", "debug")) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func loadFirstEnvFile(paths ...string) error {
	var err error
	for _, path := range paths {
		if err = godotenv.Load(path); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to load .env files: %s", strings.Join(paths, ", "))
}

func getenvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func mustGetenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	return value
}
