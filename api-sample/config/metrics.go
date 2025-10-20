package config

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Metrics struct {
	RequestCounter  metric.Int64Counter
	RequestDuration metric.Float64Histogram
	ErrorCounter    metric.Int64Counter
}

func InitMetrics() (*Metrics, error) {
	meter := otel.Meter("api-sample")

	requestCounter, err := meter.Int64Counter(
		"http.server.requests",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"http.server.duration",
		metric.WithDescription("HTTP request duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	errorCounter, err := meter.Int64Counter(
		"http.server.errors",
		metric.WithDescription("Total number of HTTP errors"),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		RequestCounter:  requestCounter,
		RequestDuration: requestDuration,
		ErrorCounter:    errorCounter,
	}, nil
}
