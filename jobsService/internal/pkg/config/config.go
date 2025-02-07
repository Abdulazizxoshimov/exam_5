package config

import (
	"os"
	"strings"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}

	ClientService struct{
		Host string
		Port string
	}

	OTLPCollector struct {
		Host string
		Port string
	}

	Kafka struct {
		Address []string
		Topic   struct {
			UserCreateTopic string
		}
	}
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":5001")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "db")
	config.DB.Port = getEnv("POSTGRES_PORT", "5332")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "4444")
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "examdb")

	// otlp collector configuration
	config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "otel-collector")
	config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4318")

	// kafka configuration
	config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "kafka:9091"), ",")
	config.Kafka.Topic.UserCreateTopic = getEnv("KAFKA_USER_CREATE_TOPIC", "api.create.user")

	config.ClientService.Host = getEnv("CLIENT_SERVICE_HOST", "clientService")
	config.ClientService.Port = getEnv("CLIENT_SERVICE_PORT", "5000")

	return &config
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
