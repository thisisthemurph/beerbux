package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load environment variables")
	}
}

type Config struct {
	Environment       Environment
	ClientBaseURL     string
	AuthServerAddress string
	GatewayAPIAddress string
	UserServerAddress string
	Secrets           SecretConfig
}

type SecretConfig struct {
	JWTSecret string
}

func Load() *Config {
	environment, err := NewEnvironment(getenvDefault("ENVIRONMENT", string(EnvProduction)))
	if err != nil {
		panic(err)
	}

	return &Config{
		Environment:       environment,
		ClientBaseURL:     mustGetenv("CLIENT_BASE_URL"),
		AuthServerAddress: mustGetenv("AUTH_SERVER_ADDRESS"),
		GatewayAPIAddress: mustGetenv("GATEWAY_API_ADDRESS"),
		UserServerAddress: mustGetenv("USER_SERVER_ADDRESS"),
		Secrets: SecretConfig{
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
