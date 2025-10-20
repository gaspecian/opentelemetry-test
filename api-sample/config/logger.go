package config

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"
)

func LogWithTrace(ctx context.Context, message string) {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()
	
	if spanCtx.IsValid() {
		log.Printf("[TraceID: %s] [SpanID: %s] %s", 
			spanCtx.TraceID().String(), 
			spanCtx.SpanID().String(), 
			message)
	} else {
		log.Println(message)
	}
}
