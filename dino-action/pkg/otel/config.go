package otel

import "fmt"

type OtelConfig struct {
	Trace   TraceConfig   `yaml:"trace" viper:"trace"`
	Metrics MetricsConfig `yaml:"metrics" viper:"metrics"`
}

type TraceConfig struct {
	Enabled         bool   `yaml:"enabled" viper:"enabled" env:"OTEL_EXPORTER_OTLP_TRACES_ENABLED"`
	TracerName      string `yaml:"tracer-name" viper:"tracer-name" env:"OTEL_EXPORTER_OTLP_TRACES_TRACER_NAME"`
	HttpEndpoint    string `yaml:"http-endpoint" viper:"http-endpoint" env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`
	HttpEndpointURL string `yaml:"http-endpoint-url" viper:"http-endpoint-url" env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT_URL"`
	GrpcEndpoint    string `yaml:"grpc-endpoint" viper:"grpc-endpoint" env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`
	GrpcEndpointURL string `yaml:"grpc-endpoint-url" viper:"grpc-endpoint-url" env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT_URL"`
	Insecure        bool   `yaml:"insecure" viper:"insecure" env:"OTEL_EXPORTER_OTLP_TRACES_INSECURE"`
}

type MetricsConfig struct {
	Enabled         bool   `yaml:"enabled" viper:"enabled" env:"OTEL_EXPORTER_OTLP_METRICS_ENABLED"`
	HttpEndpoint    string `yaml:"http-endpoint" viper:"http-endpoint" env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"`
	HttpEndpointURL string `yaml:"http-endpoint-url" viper:"http-endpoint-url" env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT_URL"`
	GrpcEndpoint    string `yaml:"grpc-endpoint" viper:"grpc-endpoint" env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"`
	GrpcEndpointURL string `yaml:"grpc-endpoint-url" viper:"grpc-endpoint-url" env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT_URL"`
	Insecure        bool   `yaml:"insecure" viper:"insecure" env:"OTEL_EXPORTER_OTLP_METRICS_INSECURE"`
}

func (c OtelConfig) Validate() error {
	if c.Trace.Enabled && countSet(c.Trace.HttpEndpointURL, c.Trace.GrpcEndpointURL, c.Trace.HttpEndpoint, c.Trace.GrpcEndpoint) != 1 {
		return fmt.Errorf("exactly one http or grpc endpoint is required when opentelemetry tracing is enabled")
	}

	if c.Metrics.Enabled && countSet(c.Metrics.HttpEndpointURL, c.Metrics.GrpcEndpointURL, c.Metrics.HttpEndpoint, c.Metrics.GrpcEndpoint) != 1 {
		return fmt.Errorf("exactly one http or grpc endpoint is required when opentelemetry metrics is enabled")
	}

	return nil
}

func countSet(s ...string) int {
	count := 0
	for _, v := range s {
		if v != "" {
			count++
		}
	}

	return count
}
