package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// FileExporter wraps the stdout exporter to also save JSON to a file.
type FileExporter struct {
	Exporter sdktrace.SpanExporter
	File     *os.File
}

// ExportSpans writes spans to both stdout and a file.
func (fe *FileExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	// Write spans to stdout
	if err := fe.Exporter.ExportSpans(ctx, spans); err != nil {
		return err
	}

	// Write spans to file
	for _, span := range spans {
		jsonData, err := spanToJSON(span)
		if err != nil {
			fmt.Printf("Error converting span to JSON: %v\n", err)
			continue
		}
		_, _ = fe.File.WriteString(jsonData + "\n")
	}
	return nil
}

// Shutdown ensures proper cleanup of the file exporter.
func (fe *FileExporter) Shutdown(ctx context.Context) error {
	if err := fe.Exporter.Shutdown(ctx); err != nil {
		return err
	}
	return fe.File.Close()
}

// Convert a span to JSON for file storage.
func spanToJSON(span sdktrace.ReadOnlySpan) (string, error) {
	data := map[string]interface{}{
		"name":       span.Name(),
		"start_time": span.StartTime().String(),
		"end_time":   span.EndTime().String(),
		"attributes": span.Attributes(),
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// SetupTracerProvider initializes the OpenTelemetry tracer provider.
func SetupTracerProvider() (*sdktrace.TracerProvider, error) {
	// Create the stdout exporter
	stdoutExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// Create the file for storing spans
	file, err := os.Create("telemetry.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry file: %w", err)
	}

	// Wrap the stdout exporter with FileExporter
	fileExporter := &FileExporter{
		Exporter: stdoutExporter,
		File:     file,
	}

	// Create the tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.Default()),
		sdktrace.WithBatcher(fileExporter),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)
	return tp, nil
}

// StartSpan starts a new span with the given name.
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("travel-planner").Start(ctx, name)
}
