package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// SetupTracerProvider initializes the OpenTelemetry tracer provider
func SetupTracerProvider() (*sdktrace.TracerProvider, error) {
	// Create a stdout exporter for tracing data
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// Create a new tracer provider with the stdout exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.Default()),
		sdktrace.WithBatcher(exporter),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)
	return tp, nil
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("travel-planner").Start(ctx, name)
}
