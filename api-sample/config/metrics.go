package config

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Metrics struct {
	// HTTP Metrics
	RequestCounter   metric.Int64Counter
	RequestDuration  metric.Float64Histogram
	ErrorCounter     metric.Int64Counter
	ActiveRequests   metric.Int64UpDownCounter
	
	// Database Metrics
	DBQueryDuration  metric.Float64Histogram
	DBQueryCounter   metric.Int64Counter
	DBErrorCounter   metric.Int64Counter
	
	// Business Metrics
	UserCreatedCounter metric.Int64Counter
	UserDeletedCounter metric.Int64Counter
	
	// Runtime tracking
	activeRequestsCount int64
}

func InitMetrics() (*Metrics, error) {
	meter := otel.Meter("api-sample")
	
	requestCounter, err := meter.Int64Counter(
		"http_server_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}
	
	requestDuration, err := meter.Float64Histogram(
		"http_server_duration_milliseconds",
		metric.WithDescription("HTTP request duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}
	
	errorCounter, err := meter.Int64Counter(
		"http_server_errors_total",
		metric.WithDescription("Total number of HTTP errors"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return nil, err
	}
	
	activeRequests, err := meter.Int64UpDownCounter(
		"http_server_active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}
	
	dbQueryDuration, err := meter.Float64Histogram(
		"db_query_duration_milliseconds",
		metric.WithDescription("Database query duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}
	
	dbQueryCounter, err := meter.Int64Counter(
		"db_queries_total",
		metric.WithDescription("Total number of database queries"),
		metric.WithUnit("{query}"),
	)
	if err != nil {
		return nil, err
	}
	
	dbErrorCounter, err := meter.Int64Counter(
		"db_errors_total",
		metric.WithDescription("Total number of database errors"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return nil, err
	}
	
	userCreatedCounter, err := meter.Int64Counter(
		"users_created_total",
		metric.WithDescription("Total number of users created"),
		metric.WithUnit("{user}"),
	)
	if err != nil {
		return nil, err
	}
	
	userDeletedCounter, err := meter.Int64Counter(
		"users_deleted_total",
		metric.WithDescription("Total number of users deleted"),
		metric.WithUnit("{user}"),
	)
	if err != nil {
		return nil, err
	}
	
	m := &Metrics{
		RequestCounter:     requestCounter,
		RequestDuration:    requestDuration,
		ErrorCounter:       errorCounter,
		ActiveRequests:     activeRequests,
		DBQueryDuration:    dbQueryDuration,
		DBQueryCounter:     dbQueryCounter,
		DBErrorCounter:     dbErrorCounter,
		UserCreatedCounter: userCreatedCounter,
		UserDeletedCounter: userDeletedCounter,
	}
	
	// Register active requests gauge callback
	_, err = meter.Int64ObservableGauge(
		"http_server_active_requests_gauge",
		metric.WithDescription("Current number of active HTTP requests"),
		metric.WithUnit("{request}"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(atomic.LoadInt64(&m.activeRequestsCount))
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}
	
	return m, nil
}

func (m *Metrics) IncrementActiveRequests() {
	atomic.AddInt64(&m.activeRequestsCount, 1)
}

func (m *Metrics) DecrementActiveRequests() {
	atomic.AddInt64(&m.activeRequestsCount, -1)
}
