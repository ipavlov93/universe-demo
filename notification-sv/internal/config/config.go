package config

import (
	"github.com/ipavlov93/universe-demo/universe-pkg/env"
)

// Config represents app config
type Config struct {
	MinLogLevel       string
	WorkersBufferSize int
	LocalStackCfg     LocalStackConfig
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
		MinLogLevel:       env.EnvironmentVariable("APP_MIN_LOG_LEVEL", "info"),
		WorkersBufferSize: env.ParseIntEnv("APP_WORKERS_BUFFER_SIZE", 10),
		LocalStackCfg: LocalStackConfig{
			Host:            env.EnvironmentVariable("LOCALSTACK_HOST", "localhost"),
			Port:            env.ParseIntEnv("LOCALSTACK_PORT", 4566),
			DefaultRegion:   env.EnvironmentVariable("AWS_DEFAULT_REGION", "us-east-1"),
			AccessKeyID:     env.EnvironmentVariable("AWS_ACCESS_KEY_ID", "universe-demo"),
			SecretAccessKey: env.EnvironmentVariable("AWS_SECRET_ACCESS_KEY", "universe-demo-secret"),
			Queue:           env.EnvironmentVariable("LOCALSTACK_PRODUCT_SV_QUEUE", "ProductSV"),
		},
	}
}
