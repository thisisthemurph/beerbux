package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load environment variables")
	}
}

type Config struct {
	Environment       Environment
	AuthServerAddress string
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
	Database          DBConfig
	Kafka             KafkaConfig
	Secrets           SecretsConfig
}

type DBConfig struct {
	Driver string
	URI    string
}

type KafkaConfig struct {
	Brokers []string
}

type SecretsConfig struct {
	JWTSecret string
}

func Load() *Config {
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

	return &Config{
		Environment:       environment,
		AuthServerAddress: mustGetenv("AUTH_SERVER_ADDRESS"),
		AccessTokenTTL:    time.Duration(accessTokenExpirationMinutes) * time.Minute,
		RefreshTokenTTL:   time.Duration(refreshTokenExpirationMinutes) * time.Minute,
		Database: DBConfig{
			Driver: mustGetenv("DB_DRIVER"),
			URI:    mustGetenv("DB_URI"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(mustGetenv("KAFKA_BROKERS"), ","),
		},
		Secrets: SecretsConfig{
			JWTSecret: mustGetenv("JWT_SECRET"),
		},
	}
}

func (c *Config) SlogLevel() slog.Level {
	if c.Environment.IsDevelopment() {
		return slog.LevelDebug
	}
	return slog.LevelInfo
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
