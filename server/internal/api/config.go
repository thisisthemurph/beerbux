package api

import (
	"beerbux/internal/auth/model"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load environment variables")
	}
}

type Config struct {
	Environment     Environment
	ClientBaseURL   string
	APIAddress      string
	Database        DBConfig
	Secrets         SecretConfig
	StreamService   StreamServiceConfig
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SecretConfig struct {
	JWTSecret string
}

type DBConfig struct {
	Driver string
	URI    string
}

type StreamServiceConfig struct {
	HeartbeatTickerSeconds int64
}

func LoadConfig() *Config {
	environment, err := NewEnvironment(getenvDefault("ENVIRONMENT", string(EnvProduction)))
	if err != nil {
		panic(err)
	}

	accessTokenExpiration := getenvDefault("ACCESS_TOKEN_EXPIRATION", "15")
	refreshTokenExpiration := getenvDefault("REFRESH_TOKEN_EXPIRATION", "10080")

	accessTokenExpirationMinutes, err := strconv.ParseInt(accessTokenExpiration, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid ACCESS_TOKEN_EXPIRATION: %s", accessTokenExpiration))
	}

	refreshTokenExpirationMinutes, err := strconv.ParseInt(refreshTokenExpiration, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid REFRESH_TOKEN_EXPIRATION: %s", refreshTokenExpiration))
	}

	hbIntervalSeconds := mustGetenv("HEARTBEAT_INTERVAL_SECONDS")
	heartbeatIntervalSeconds, err := strconv.ParseInt(hbIntervalSeconds, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid HEARTBEAT_INTERVAL_SECONDS: %s", hbIntervalSeconds))
	}

	return &Config{
		Environment:   environment,
		ClientBaseURL: mustGetenv("CLIENT_BASE_URL"),
		APIAddress:    mustGetenv("API_ADDRESS"),
		Database: DBConfig{
			Driver: mustGetenv("DB_DRIVER"),
			URI:    mustGetenv("DB_URI"),
		},
		Secrets: SecretConfig{
			JWTSecret: mustGetenv("JWT_SECRET"),
		},
		StreamService: StreamServiceConfig{
			HeartbeatTickerSeconds: heartbeatIntervalSeconds,
		},
		AccessTokenTTL:  time.Duration(accessTokenExpirationMinutes) * time.Minute,
		RefreshTokenTTL: time.Duration(refreshTokenExpirationMinutes) * time.Minute,
	}
}

func (c *Config) SlogLevel() slog.Level {
	if c.Environment.IsDevelopment() {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func (c *Config) GetAuthOptions() model.AuthOptions {
	return model.AuthOptions{
		JWTSecret:       c.Secrets.JWTSecret,
		AccessTokenTTL:  c.AccessTokenTTL,
		RefreshTokenTTL: c.RefreshTokenTTL,
	}
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
