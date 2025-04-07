package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load environment variables")
	}
}

type Config struct {
	ServiceName            string
	Environment            Environment
	StreamServerAddr       string
	ClientBaseURL          string
	HeartbeatTickerSeconds int64
	Kafka                  KafkaConfig
}

type KafkaConfig struct {
	Brokers []string
}

func Load() *Config {
	environment, err := NewEnvironment(getenvDefault("ENVIRONMENT", string(EnvProduction)))
	if err != nil {
		log.Fatal("Failed to determine environment", err)
	}

	heartbeatInterval := getenvDefault("HEARTBEAT_INTERVAL_SECONDS", "15")
	heartbeatIntervalSeconds, err := strconv.ParseInt(heartbeatInterval, 10, 64)
	if err != nil {
		log.Fatalf("Invalid heartbeat %s: %v", heartbeatInterval, err)
	}

	brokers := strings.Split(getenvDefault("KAFKA_BROKERS", "localhost:9092"), ",")
	if len(brokers) == 0 {
		log.Fatal("No Kafka brokers provided")
	}

	return &Config{
		ServiceName:            "stream-service",
		Environment:            environment,
		StreamServerAddr:       mustGetenv("STREAM_SERVER_ADDRESS"),
		ClientBaseURL:          mustGetenv("CLIENT_BASE_URL"),
		HeartbeatTickerSeconds: heartbeatIntervalSeconds,
		Kafka: KafkaConfig{
			Brokers: brokers,
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
