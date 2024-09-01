package config

import "gitlab.com/stackvista/demo/kubecon2024/poi/pkg/otel"

type Config struct {
	Port          int             `yaml:"port" viper:"port" env:"PORT" default:"3000"`
	StoreContents string          `yaml:"store_file" viper:"store_file" env:"STORE_FILE" default:"store.yaml"`
	OpenTelemetry otel.OtelConfig `yaml:"opentelemetry" viper:"opentelemetry"`
}

func Validate(cfg *Config) error {
	return cfg.OpenTelemetry.Validate()
}
