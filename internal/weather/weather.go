package weather

import (
	"context"
	"fmt"
	"time"
	"travel-planner/internal/tracing"

	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/exp/rand"
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

	// Simulate API Call with a random sleep between 1-5 seconds
	sleepDuration := time.Duration(rand.Intn(5)+1) * time.Second // Random duration between 1 and 5 seconds
	time.Sleep(sleepDuration)

	// Simulate API Call Result
	result := fmt.Sprintf("Sunny in %s from %s to %s", city, startDate, endDate)

	return result, nil
}
