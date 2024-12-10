package tracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// FileExporter wraps the stdout exporter to also save logs to a file.
type FileExporter struct {
	Exporter sdktrace.SpanExporter // The wrapped stdout exporter
	File     *os.File              // File to write spans as logs
}

// ExportSpans writes filtered span information to both stdout and a log file.
func (fe *FileExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, span := range spans {
		// Extract the filtered information
		summary, err := createSpanSummary(span)
		if err != nil {
			fmt.Printf("Error creating span summary: %v\n", err)
			continue
		}

		// Write the concise summary to the command line
		fmt.Printf("Span: %s\n  Start Time: %s\n  End Time: %s\n  Time Taken: %s\n\n",
			summary["name"], summary["start_time"], summary["end_time"], summary["time_taken"])

		// Format the summary as a log entry
		logEntry := fmt.Sprintf(
			"[%s] %s - Start Time: %s, End Time: %s, Time Taken: %s, Attributes: %v\n",
			time.Now().Format(time.RFC3339),
			summary["name"],
			summary["start_time"],
			summary["end_time"],
			summary["time_taken"],
			summary["attributes"],
		)

		// Append the log entry to the file
		_, err = fe.File.WriteString(logEntry)
		if err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
			return err
		}
	}

	return nil
}

// createSpanSummary extracts and formats the required information from a span.
func createSpanSummary(span sdktrace.ReadOnlySpan) (map[string]interface{}, error) {
	// Calculate the time taken for the span
	start := span.StartTime()
	end := span.EndTime()
	duration := end.Sub(start)

	// Filter the attributes
	attributes := make(map[string]interface{})
	for _, attr := range span.Attributes() {
		attributes[string(attr.Key)] = attr.Value.AsInterface()
	}

	// Create a summary object
	summary := map[string]interface{}{
		"name":       span.Name(),
		"start_time": start.String(),
		"end_time":   end.String(),
		"time_taken": duration.String(),
		"attributes": attributes,
	}

	return summary, nil
}

// Shutdown ensures proper cleanup of the file exporter.
func (fe *FileExporter) Shutdown(ctx context.Context) error {
	// Shutdown the wrapped exporter
	if err := fe.Exporter.Shutdown(ctx); err != nil {
		return err
	}
	// Close the file to release resources
	return fe.File.Close()
}

// SetupTracerProvider initializes the OpenTelemetry tracer provider.
func SetupTracerProvider() (*sdktrace.TracerProvider, error) {
	// Create the stdout exporter for human-readable output
	stdoutExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// Open the file in append mode to preserve existing logs and add new entries
	file, err := os.OpenFile("telemetry.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open telemetry logs file: %w", err)
	}

	// Wrap the stdout exporter with FileExporter to handle log file output
	fileExporter := &FileExporter{
		Exporter: stdoutExporter,
		File:     file,
	}

	// Create the tracer provider and configure it with:
	// - Always sampling spans (useful for debugging)
	// - The default resource attributes
	// - A batcher to process spans efficiently and export them
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.Default()),
		sdktrace.WithBatcher(fileExporter),
	)

	// Set the global tracer provider so spans can be created and managed globally
	otel.SetTracerProvider(tp)
	return tp, nil
}

// StartSpan starts a new span with the given name.
// The span can be used to trace operations within the application.
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("travel-planner").Start(ctx, name)
}
