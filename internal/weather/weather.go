package weather

import (
	"context"
	"fmt"
	"log"
	"travel-planner/internal/tracing"

	"go.opentelemetry.io/otel/attribute"
)

func GetWeather(ctx context.Context, city, startDate, endDate string) (string, error) {
	_, span := tracing.StartSpan(ctx, "GetWeather")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("city", city),
		attribute.String("startDate", startDate),
		attribute.String("endDate", endDate),
	)

	// Simulate an error scenario
	if city == "" {
		err := fmt.Errorf("city cannot be empty")
		span.RecordError(err)
		return "", err
	}

	log.Printf("Fetching weather for city: %s, dates: %s to %s", city, startDate, endDate)

	// Simulate API Call
	result := fmt.Sprintf("Sunny in %s from %s to %s", city, startDate, endDate)
	return result, nil
}
