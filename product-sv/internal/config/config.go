package config

import (
	"github.com/ipavlov93/universe-demo/universe-pkg/env"
)

// Config represents app config
type Config struct {
	ServerPort    int
	MinLogLevel   string
	PostgresCfg   PostgresConfig
	LocalStackCfg LocalStackConfig
}

type PostgresConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	Options      string
}

type LocalStackConfig struct {
	Host            string
	Port            int
	DefaultRegion   string
	AccessKeyID     string
	SecretAccessKey string
	Queue           string
}

// LoadConfigEnv sets Config with environment variables values
func LoadConfigEnv() Config {
	return Config{
		ServerPort:  env.ParseIntEnv("PRODUCT_SV_HTTP_SERVER_PORT", 0),
		MinLogLevel: env.EnvironmentVariable("APP_MIN_LOG_LEVEL", "info"),
		LocalStackCfg: LocalStackConfig{
			Host:            env.EnvironmentVariable("LOCALSTACK_HOST", "localhost"),
			Port:            env.ParseIntEnv("LOCALSTACK_PORT", 4566),
			DefaultRegion:   env.EnvironmentVariable("AWS_DEFAULT_REGION", "us-east-1"),
			AccessKeyID:     env.EnvironmentVariable("AWS_ACCESS_KEY_ID", "universe-demo"),
			SecretAccessKey: env.EnvironmentVariable("AWS_SECRET_ACCESS_KEY", "universe-demo-secret"),
			Queue:           env.EnvironmentVariable("LOCALSTACK_PRODUCT_SV_QUEUE", "ProductSV"),
		},
		PostgresCfg: PostgresConfig{
			Host:         env.EnvironmentVariable("POSTGRES_HOST", "localhost"),
			Port:         env.ParseIntEnv("POSTGRES_PORT", 5432),
			User:         env.EnvironmentVariable("POSTGRES_USER", "universe_demo"),
			Password:     env.EnvironmentVariable("POSTGRES_PASSWORD", "secret"),
			DatabaseName: env.EnvironmentVariable("POSTGRES_DB_NAME", "product-sv-db"),
			Options:      env.EnvironmentVariable("POSTGRES_OPTIONS", "sslmode=disable"),
		},
	}
}
