package config

import "os"

type Config struct {
	MongoURI     string
	OTELEndpoint string
	ServiceName  string
	ServiceVersion string
	Port         string
}

func Load() *Config {
	return &Config{
		MongoURI:       getEnv("MONGODB_URI", "mongodb://admin:password@localhost:27017/apidb?authSource=admin"),
		OTELEndpoint:   getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		ServiceName:    getEnv("SERVICE_NAME", "api-sample"),
		ServiceVersion: getEnv("SERVICE_VERSION", "1.0.0"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
