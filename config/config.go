package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDB  MongoDBConfig
	RabbitMQ RabbitMQConfig
	Postgres PostgresConfig
	App      AppConfig
}

type MongoDBConfig struct {
	URI    string
	DBName string
}

type RabbitMQConfig struct {
	URL   string
	Queue string
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type AppConfig struct {
	LogLevel string
}

func LoadConfig() (*Config, error) {

	_ = godotenv.Load()

	config := &Config{
		MongoDB: MongoDBConfig{
			URI:    getEnv("MONGODB_URI", "mongodb://localhost:27017/campaigns"),
			DBName: getEnv("MONGODB_DB_NAME", "campaigns"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:   getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672"),
			Queue: getEnv("RABBITMQ_QUEUE", "campaign_messages"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			Database: getEnv("POSTGRES_DB", "campaigns"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		},
		App: AppConfig{
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
