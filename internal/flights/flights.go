package flights

import (
	"context"
	"fmt"
	"log"
	"travel-planner/internal/tracing"

	"go.opentelemetry.io/otel/attribute"
)

func GetFlights(ctx context.Context, city, startDate, endDate string, minBudget, maxBudget int) (string, error) {
	_, span := tracing.StartSpan(ctx, "GetFlights")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("city", city),
		attribute.String("startDate", startDate),
		attribute.String("endDate", endDate),
		attribute.Int("minBudget", minBudget),
		attribute.Int("maxBudget", maxBudget),
	)

	// Simulate an error scenario
	if city == "" {
		err := fmt.Errorf("city cannot be empty")
		span.RecordError(err)
		return "", err
	}

	log.Printf("Fetching flights for city: %s, budget: %d-%d", city, minBudget, maxBudget)

	// Simulate API Call
	result := fmt.Sprintf("Flights to %s: $500-$700", city)
	return result, nil
}
