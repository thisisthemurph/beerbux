package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load environment variables")
	}
}

type Config struct {
	Environment       Environment
	UserServerAddress string
	Database          DBConfig
	Kafka             KafkaConfig
}

type DBConfig struct {
	Driver string
	URI    string
}

type KafkaConfig struct {
	Brokers []string
}

func Load() *Config {
	environment, err := NewEnvironment(getenvDefault("ENVIRONMENT", string(EnvProduction)))
	if err != nil {
		panic(err)
	}

	return &Config{
		Environment:       environment,
		UserServerAddress: mustGetenv("USER_SERVER_ADDRESS"),
		Database: DBConfig{
			Driver: mustGetenv("DB_DRIVER"),
			URI:    mustGetenv("DB_URI"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(mustGetenv("KAFKA_BROKERS"), ","),
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
