package config

import "github.com/kelseyhightower/envconfig"

// Config holds all environment variables for the applications
type Config struct {
	DatabaseURL     string   `envconfig:"DATABASE_URL" required:"true"`
	KafkaBrokerURLs []string `envconfig:"KAFKA_BROKER_URLS" required:"true"`
	KafkaTopic      string   `envconfig:"KAFKA_TOPIC" required:"true"`
	KafkaGroupID    string   `envconfig:"KAFKA_GROUP_ID"` // Optional, only for processor
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}